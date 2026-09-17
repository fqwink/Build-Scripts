package components

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAPILoginStatusAndQueue(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	writeTestJSON(t, filepath.Join(state, ".build_state"), apiBuildState{
		Running: true,
		Queued:  []map[string]any{},
	})

	resp := apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("status code=%d body=%s", resp.Code, resp.Body.String())
	}
	var status map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &status)
	if status["running"] != true || status["last_build_status"] != "none" {
		t.Fatalf("unexpected status: %#v", status)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/build", token, nil)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("build code=%d body=%s", resp.Code, resp.Body.String())
	}
	var build map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &build)
	if build["message"] != "Build queued" || build["queued"] != true {
		t.Fatalf("unexpected build response: %#v", build)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/queue", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("queue code=%d body=%s", resp.Code, resp.Body.String())
	}
	var queue map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &queue)
	if len(queue["queued"].([]any)) != 1 || queue["max_size"].(float64) != 3 {
		t.Fatalf("unexpected queue: %#v", queue)
	}

	resp = apiRequest(t, server, http.MethodDelete, "/api/queue", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("clear code=%d body=%s", resp.Code, resp.Body.String())
	}
	var clear map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &clear)
	if clear["cleared_count"].(float64) != 1 {
		t.Fatalf("unexpected clear: %#v", clear)
	}
}

func TestAPILogsHistoryAndCircuitReset(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	if err := os.MkdirAll(filepath.Join(state, ".build_logs"), 0755); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(state, ".build_logs", "b20260917010101.json"), apiBuildLog{
		ID: "b20260917010101", StartedAt: "2026-09-17T01:01:01Z", FinishedAt: "2026-09-17T01:01:03Z",
		TargetStatus: "success", Pipeline: apiPipelineLog{Stdout: "[INFO] start\n[INFO] done", Stderr: ""}, Warnings: []string{"[WARNING] slow"},
	})
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260917010101","finished_at":"2026-09-17T01:01:03Z","status":"success","trigger":"manual","duration_seconds":2}`)
	writeTestJSON(t, filepath.Join(state, ".build_circuit_state"), apiCircuitState{Open: true, ConsecutiveFailures: 3})

	resp := apiRequest(t, server, http.MethodGet, "/api/logs?n=2", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("logs code=%d body=%s", resp.Code, resp.Body.String())
	}
	var logs map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &logs)
	if len(logs["lines"].([]any)) != 2 {
		t.Fatalf("unexpected logs: %#v", logs)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/history", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("history code=%d body=%s", resp.Code, resp.Body.String())
	}
	var history map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &history)
	if history["total"].(float64) != 1 || history["pages"].(float64) != 1 {
		t.Fatalf("unexpected history: %#v", history)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/history/b20260917010101/comment", token, map[string]string{"comment": "reviewed"})
	if resp.Code != http.StatusOK {
		t.Fatalf("comment code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/history/b20260917010101/flag", token, map[string]bool{"flagged": true})
	if resp.Code != http.StatusOK {
		t.Fatalf("flag code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/history/b20260917010101/tags", token, map[string][]string{"tags": []string{"release", "manual"}})
	if resp.Code != http.StatusOK {
		t.Fatalf("tags code=%d body=%s", resp.Code, resp.Body.String())
	}
	var updated apiBuildLog
	readTestJSON(t, filepath.Join(state, ".build_logs", "b20260917010101.json"), &updated)
	if updated.Comment == nil || *updated.Comment != "reviewed" || !updated.Flagged || len(updated.Tags) != 2 {
		t.Fatalf("history log was not updated: %#v", updated)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/history?trigger=manual&tag=release&flagged=true", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("filtered history code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &history)
	if history["total"].(float64) != 1 {
		t.Fatalf("unexpected filtered history: %#v", history)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/circuit-breaker/reset", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("circuit code=%d body=%s", resp.Code, resp.Body.String())
	}
	var circuit apiCircuitState
	readTestJSON(t, filepath.Join(state, ".build_circuit_state"), &circuit)
	if circuit.Open || circuit.ConsecutiveFailures != 0 {
		t.Fatalf("circuit was not reset: %#v", circuit)
	}
}

func TestAPICommonErrors(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	resp := apiRequest(t, server, http.MethodGet, "/api/status", "", nil)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized code=%d body=%s", resp.Code, resp.Body.String())
	}
	token := login(t, server, "admin")
	resp = apiRequest(t, server, http.MethodPost, "/api/status", token, nil)
	if resp.Code != http.StatusMethodNotAllowed {
		t.Fatalf("method code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/history?page=0", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("validation code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPIMaintenanceEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	resp := apiRequest(t, server, http.MethodGet, "/api/maintenance", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("maintenance code=%d body=%s", resp.Code, resp.Body.String())
	}
	var maintenance apiMaintenanceState
	decodeTestJSON(t, resp.Body.Bytes(), &maintenance)
	if maintenance.Enabled || maintenance.Reason != nil || maintenance.Since != nil {
		t.Fatalf("unexpected default maintenance: %#v", maintenance)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/maintenance/enable", token, map[string]string{"reason": "  deploy window  "})
	if resp.Code != http.StatusOK {
		t.Fatalf("enable code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".maintenance"), &maintenance)
	if !maintenance.Enabled || stringPtrValue(maintenance.Reason) != "deploy window" || maintenance.Since == nil {
		t.Fatalf("maintenance was not enabled: %#v", maintenance)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/build", token, nil)
	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("build during maintenance code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/maintenance/enable", token, map[string]string{"reason": "deploy window"})
	if resp.Code != http.StatusOK {
		t.Fatalf("enable no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["message"] != "No changes" {
		t.Fatalf("unexpected enable no-op: %#v", body)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/maintenance/disable", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("disable code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".maintenance"), &maintenance)
	if maintenance.Enabled || maintenance.Reason != nil || maintenance.Since != nil {
		t.Fatalf("maintenance was not disabled: %#v", maintenance)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/build", token, nil)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("build after maintenance code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/config-log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("config log code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["total"].(float64) != 2 {
		t.Fatalf("unexpected config log: %#v", body)
	}
}

func TestAPIRepoAndBranchConfigEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	resp := apiRequest(t, server, http.MethodGet, "/api/repo-info", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("repo-info code=%d body=%s", resp.Code, resp.Body.String())
	}
	var repo apiRepoConfig
	decodeTestJSON(t, resp.Body.Bytes(), &repo)
	if repo.Owner == "" || repo.Repo == "" || repo.Branch != "main" || repo.TargetFile != "docs" {
		t.Fatalf("unexpected repo defaults: %#v", repo)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/repo-config", token, map[string]string{
		"owner": " fqwink ", "repo": "Build-Scripts", "branch": "release", "target_file": "docs",
	})
	if resp.Code != http.StatusOK {
		t.Fatalf("repo-config code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".repo_config"), &repo)
	if repo.Owner != "fqwink" || repo.Repo != "Build-Scripts" || repo.Branch != "release" || repo.UpdatedAt == "" {
		t.Fatalf("repo config was not saved: %#v", repo)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/repo-config", token, map[string]string{"branch": "release"})
	if resp.Code != http.StatusOK {
		t.Fatalf("repo-config no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["message"] != "No changes" {
		t.Fatalf("unexpected repo no-op: %#v", body)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/repo-config", token, map[string]string{"target_file": "../secret"})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("repo invalid code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/branch-config", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("branch-config get code=%d body=%s", resp.Code, resp.Body.String())
	}
	var branchResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &branchResp)
	if branchResp["source"] != "default" || len(branchResp["branches"].([]any)) != 1 {
		t.Fatalf("unexpected branch defaults: %#v", branchResp)
	}

	targets := []apiBranchTarget{{
		Branch:     "release",
		TargetFile: "docs",
		SHAFile:    filepath.Join(state, ".last_sha.release"),
		Src:        filepath.Join(state, "repo", "release", "docs"),
		Out:        filepath.Join(state, "dist", "release"),
		DeployTargets: []apiDeployTarget{{
			Host: "192.0.2.1", User: "deploy", DestDir: "/var/www/html/",
		}},
	}}
	resp = apiRequest(t, server, http.MethodPost, "/api/branch-config", token, map[string]any{"branches": targets})
	if resp.Code != http.StatusOK {
		t.Fatalf("branch-config post code=%d body=%s", resp.Code, resp.Body.String())
	}
	var stored apiBranchConfigFile
	readTestJSON(t, filepath.Join(state, ".branch_config"), &stored)
	if len(stored.BranchTargets) != 1 || stored.BranchTargets[0].Branch != "release" {
		t.Fatalf("branch config was not saved: %#v", stored)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/branch-config", token, map[string]any{"branches": targets})
	if resp.Code != http.StatusOK {
		t.Fatalf("branch-config no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["message"] != "No changes" {
		t.Fatalf("unexpected branch no-op: %#v", body)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/branch-config", token, map[string]any{"branches": []apiBranchTarget{}})
	if resp.Code != http.StatusOK {
		t.Fatalf("branch-config clear code=%d body=%s", resp.Code, resp.Body.String())
	}
	if _, err := os.Stat(filepath.Join(state, ".branch_config")); !os.IsNotExist(err) {
		t.Fatalf("branch config should be removed, err=%v", err)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/branch-config", token, map[string]any{"branches": []apiBranchTarget{{
		Branch: "bad", TargetFile: "../docs", SHAFile: filepath.Join(state, ".sha"), Src: filepath.Join(state, "src"), Out: filepath.Join(state, "out"),
	}}})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("branch invalid code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/config-log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("config log code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["total"].(float64) != 3 {
		t.Fatalf("unexpected config log: %#v", body)
	}
}

func TestAPIScheduleEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	resp := apiRequest(t, server, http.MethodGet, "/api/schedule", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("schedule code=%d body=%s", resp.Code, resp.Body.String())
	}
	var schedule map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &schedule)
	if schedule["interval"] != "5min" || schedule["interval_seconds"].(float64) != 300 || schedule["paused"] != false {
		t.Fatalf("unexpected schedule defaults: %#v", schedule)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/interval", token, map[string]int{"interval_seconds": 600})
	if resp.Code != http.StatusOK {
		t.Fatalf("interval code=%d body=%s", resp.Code, resp.Body.String())
	}
	var cfg apiServerConfig
	readTestJSON(t, filepath.Join(state, ".server_config"), &cfg)
	if cfg.ScheduleIntervalSeconds != 600 {
		t.Fatalf("interval was not saved: %#v", cfg)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/allowed-hours", token, map[string]any{"from": 9, "to": 18})
	if resp.Code != http.StatusOK {
		t.Fatalf("allowed-hours code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".server_config"), &cfg)
	if cfg.AllowedHours == nil || cfg.AllowedHours.From != 9 || cfg.AllowedHours.To != 18 {
		t.Fatalf("allowed hours were not saved: %#v", cfg)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/pause", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("pause code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/pause", token, nil)
	if resp.Code != http.StatusConflict {
		t.Fatalf("pause conflict code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/schedule", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("schedule paused code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &schedule)
	if schedule["paused"] != true || schedule["next_run_at"] != nil {
		t.Fatalf("unexpected paused schedule: %#v", schedule)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/resume", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("resume code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/force-interval", token, map[string]int{"hours": 24})
	if resp.Code != http.StatusOK {
		t.Fatalf("force interval code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/cooldown", token, map[string]int{"seconds": 120})
	if resp.Code != http.StatusOK {
		t.Fatalf("cooldown code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/allowed-hours", token, map[string]any{"from": nil, "to": nil})
	if resp.Code != http.StatusOK {
		t.Fatalf("clear allowed-hours code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".server_config"), &cfg)
	if cfg.SchedulePaused || cfg.ForceBuildIntervalHours != 24 || cfg.BuildCooldownSeconds != 120 || cfg.AllowedHours != nil {
		t.Fatalf("schedule values were not saved: %#v", cfg)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/interval", token, map[string]int{"interval_seconds": 1})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid interval code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/config", token, map[string]int{"schedule_interval_seconds": 900})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("direct schedule config code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/config-log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("config log code=%d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["total"].(float64) != 7 {
		t.Fatalf("unexpected schedule config log: %#v", body)
	}
}

func TestAPINotifyEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	requests := 0
	signatures := []string{}
	events := []string{}
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		signatures = append(signatures, r.Header.Get("X-Adlaire-Signature"))
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("webhook body decode failed: %v", err)
		}
		if event, _ := body["event"].(string); event != "" {
			events = append(events, event)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer webhook.Close()

	resp := apiRequest(t, server, http.MethodGet, "/api/notify-config", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("default notify config code=%d body=%s", resp.Code, resp.Body.String())
	}
	var cfg apiNotifyConfig
	decodeTestJSON(t, resp.Body.Bytes(), &cfg)
	if len(cfg.Webhooks) != 0 || len(cfg.Channels) != 0 || cfg.Summary.Interval != "weekly" || cfg.Summary.Hour != 9 || cfg.Summary.DayOfWeek != 1 {
		t.Fatalf("unexpected default notify config: %#v", cfg)
	}

	secret := "notify-secret"
	notifyCfg := apiNotifyConfig{
		Webhooks: []apiNotifyWebhook{{
			URL:                  webhook.URL,
			Label:                "main",
			Enabled:              true,
			On:                   []string{"failure", "weekly_summary"},
			RetryCount:           2,
			RetryIntervalSeconds: 30,
			Secret:               &secret,
		}},
		On:      []string{"failure", "weekly_summary"},
		Summary: apiNotifySummary{Enabled: true, Interval: "weekly", Hour: 9, DayOfWeek: 1},
		Email:   apiNotifyEmail{Enabled: false, To: []string{}, On: []string{}},
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/notify-config", token, notifyCfg)
	if resp.Code != http.StatusOK {
		t.Fatalf("set notify config code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".notify_config"), &cfg)
	if len(cfg.Webhooks) != 1 || cfg.Webhooks[0].Secret == nil || *cfg.Webhooks[0].Secret != secret {
		t.Fatalf("notify config was not saved with secret: %#v", cfg)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/notify-config", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("masked notify config code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &cfg)
	if cfg.Webhooks[0].Secret == nil || *cfg.Webhooks[0].Secret != "***" {
		t.Fatalf("secret must be masked: %#v", cfg.Webhooks[0].Secret)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/notify-test", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("notify test code=%d body=%s", resp.Code, resp.Body.String())
	}
	if requests != 1 || events[0] != "test" || signatures[0] == "" {
		t.Fatalf("test webhook not delivered correctly: requests=%d events=%#v signatures=%#v", requests, events, signatures)
	}

	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260917010001","finished_at":"2026-09-17T01:00:06Z","status":"success","trigger":"manual","duration_seconds":5}`)
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260916010101","finished_at":"2026-09-16T01:01:06Z","status":"failure","trigger":"manual","duration_seconds":7}`)
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260901010101","finished_at":"2026-09-01T01:01:06Z","status":"failure","trigger":"manual","duration_seconds":9}`)
	resp = apiRequest(t, server, http.MethodPost, "/api/notify/weekly-summary", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("weekly summary code=%d body=%s", resp.Code, resp.Body.String())
	}
	var summary map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &summary)
	if summary["success_count"].(float64) != 1 || summary["failure_count"].(float64) != 1 || summary["success_rate"].(float64) != 50 {
		t.Fatalf("unexpected weekly summary: %#v", summary)
	}
	if requests != 2 || events[1] != "weekly_summary" {
		t.Fatalf("weekly webhook not delivered: requests=%d events=%#v", requests, events)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/notify-log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("notify log code=%d body=%s", resp.Code, resp.Body.String())
	}
	var notifyLog map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &notifyLog)
	if notifyLog["total"].(float64) != 2 {
		t.Fatalf("unexpected notify log: %#v", notifyLog)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/notify-config", token, map[string]any{
		"webhooks": []map[string]any{{
			"url":                    webhook.URL,
			"on":                     []string{"failure"},
			"retry_count":            1,
			"retry_interval_seconds": 10,
		}},
		"summary": map[string]any{"interval": "weekly"},
	})
	if resp.Code != http.StatusOK {
		t.Fatalf("notify config omitted enabled code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".notify_config"), &cfg)
	if len(cfg.Webhooks) != 1 || !cfg.Webhooks[0].Enabled {
		t.Fatalf("omitted enabled must default true: %#v", cfg.Webhooks)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/notify-config", token, map[string]any{"webhooks": []map[string]any{{"url": "ftp://example.com", "enabled": true}}, "summary": map[string]any{"interval": "weekly"}})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid notify config must fail: code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPISessionPasswordAndStream(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	if err := os.MkdirAll(filepath.Join(state, ".build_logs"), 0755); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(state, ".build_logs", "b20260917010101.json"), apiBuildLog{
		ID: "b20260917010101", StartedAt: "2026-09-17T01:01:01Z", FinishedAt: "2026-09-17T01:01:03Z",
		TargetStatus: "success", Pipeline: apiPipelineLog{Stdout: "[INFO] streamed", Stderr: ""},
	})

	resp := apiRequest(t, server, http.MethodGet, "/api/build/stream", token, nil)
	if resp.Code != http.StatusOK || resp.Header().Get("Content-Type") != "text/event-stream" {
		t.Fatalf("stream code=%d content-type=%q body=%s", resp.Code, resp.Header().Get("Content-Type"), resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/change-password", token, map[string]string{
		"current_password": "admin",
		"new_password":     "changed-password",
	})
	if resp.Code != http.StatusOK {
		t.Fatalf("change-password code=%d body=%s", resp.Code, resp.Body.String())
	}
	_ = login(t, server, "changed-password")

	resp = apiRequest(t, server, http.MethodPost, "/api/logout", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("logout code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("post-logout code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPIOperationConfigSessionsAndLogs(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")
	second := login(t, server, "admin")

	resp := apiRequest(t, server, http.MethodGet, "/api/health", "", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("health code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/config", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("config code=%d body=%s", resp.Code, resp.Body.String())
	}
	var cfg apiServerConfig
	decodeTestJSON(t, resp.Body.Bytes(), &cfg)
	if cfg.QueueMaxSize != 3 || cfg.LogMaxLines != 500 || cfg.SessionTimeoutSeconds != 28800 {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/config/validate", token, map[string]any{"queue_max_size": 0, "log_level": "DEBUG"})
	if resp.Code != http.StatusOK {
		t.Fatalf("validate code=%d body=%s", resp.Code, resp.Body.String())
	}
	if _, err := os.Stat(filepath.Join(state, ".server_config")); !os.IsNotExist(err) {
		t.Fatalf("validate must not write server config")
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/config", token, map[string]any{"queue_max_size": 0, "log_level": "DEBUG"})
	if resp.Code != http.StatusOK {
		t.Fatalf("set config code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".server_config"), &cfg)
	if cfg.QueueMaxSize != 0 || cfg.LogLevel != "DEBUG" {
		t.Fatalf("config was not saved: %#v", cfg)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/sessions", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("sessions code=%d body=%s", resp.Code, resp.Body.String())
	}
	var sessions map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &sessions)
	if len(sessions["sessions"].([]any)) != 2 {
		t.Fatalf("unexpected sessions: %#v", sessions)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/sessions/revoke-all", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("revoke code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", second, nil)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("second token must be revoked: code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/config-log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("config log code=%d body=%s", resp.Code, resp.Body.String())
	}
	var configLog map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &configLog)
	if len(configLog["log"].([]any)) != 1 {
		t.Fatalf("unexpected config log: %#v", configLog)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/access-log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("access log code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/api-access-log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("api access log code=%d body=%s", resp.Code, resp.Body.String())
	}
	var apiAccess map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &apiAccess)
	if apiAccess["total"].(float64) == 0 {
		t.Fatalf("api access log should contain requests: %#v", apiAccess)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/api-access-log?method=GET&path=/api/config&status=200", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("filtered api access log code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &apiAccess)
	if apiAccess["total"].(float64) == 0 {
		t.Fatalf("filtered api access log should contain config requests: %#v", apiAccess)
	}
}

func TestAPIReadOnlyAggregateEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	siteDir := filepath.Join(state, "site")
	if err := os.MkdirAll(siteDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(siteDir, "index.html"), []byte("<h1>Adlaire</h1>\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(state, ".build_logs"), 0755); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(state, ".build_logs", "b20260917010101.json"), apiBuildLog{
		ID: "b20260917010101", StartedAt: "2026-09-17T01:01:01Z", FinishedAt: "2026-09-17T01:01:06Z",
		TargetStatus: "success", Commit: map[string]any{"sha": "abc123"}, Warnings: []string{"[WARNING] slow"}, DurationSeconds: 5,
	})
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260917010101","finished_at":"2026-09-17T01:01:06Z","status":"success","trigger":"manual","duration_seconds":5}`)
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260916010101","finished_at":"2026-09-16T01:01:06Z","status":"failure","trigger":"manual","duration_seconds":7}`)
	appendLine(t, filepath.Join(state, ".notify_log"), `{"at":"2026-09-17T01:01:07Z","type":"build","result":"sent"}`)

	resp := apiRequest(t, server, http.MethodGet, "/api/sysinfo", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("sysinfo code=%d body=%s", resp.Code, resp.Body.String())
	}
	var sysinfo map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &sysinfo)
	if sysinfo["output_exists"] != true || sysinfo["output_size_bytes"].(float64) == 0 {
		t.Fatalf("unexpected sysinfo: %#v", sysinfo)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/stats?days=7", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("stats code=%d body=%s", resp.Code, resp.Body.String())
	}
	var stats map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &stats)
	if stats["total"].(float64) != 2 || stats["success"].(float64) != 1 || stats["failure"].(float64) != 1 {
		t.Fatalf("unexpected stats: %#v", stats)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/stats/timeline?days=7", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("timeline code=%d body=%s", resp.Code, resp.Body.String())
	}
	var timeline map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &timeline)
	if len(timeline["timeline"].([]any)) != 2 {
		t.Fatalf("unexpected timeline: %#v", timeline)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/stats/build-duration?n=1", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("duration code=%d body=%s", resp.Code, resp.Body.String())
	}
	var duration map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &duration)
	if duration["count"].(float64) != 1 || duration["avg_seconds"].(float64) != 5 {
		t.Fatalf("unexpected duration: %#v", duration)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/output-meta", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("output-meta code=%d body=%s", resp.Code, resp.Body.String())
	}
	var outputMeta map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &outputMeta)
	if outputMeta["sha256"] == "" || outputMeta["build_id"] != "b20260917010101" || outputMeta["commit_sha"] != "abc123" {
		t.Fatalf("unexpected output meta: %#v", outputMeta)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/dashboard", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("dashboard code=%d body=%s", resp.Code, resp.Body.String())
	}
	var dashboard map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &dashboard)
	for _, key := range []string{"status", "sysinfo", "stats", "schedule", "alerts"} {
		if dashboard[key] == nil {
			t.Fatalf("missing dashboard key %q: %#v", key, dashboard)
		}
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/notify-log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("notify-log code=%d body=%s", resp.Code, resp.Body.String())
	}
	var notify map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &notify)
	if notify["total"].(float64) != 1 || len(notify["log"].([]any)) != 1 {
		t.Fatalf("unexpected notify log: %#v", notify)
	}
}

func newAPIState(t *testing.T) string {
	t.Helper()
	state := t.TempDir()
	now := time.Date(2026, 9, 17, 1, 0, 0, 0, time.UTC)
	if err := InitCredentials(state, "admin", now); err != nil {
		t.Fatal(err)
	}
	return state
}

func newTestAPI(t *testing.T, state string) *APIServer {
	t.Helper()
	server, err := NewAPIServer(APIConfig{
		StateDir: state,
		Now: func() time.Time {
			return time.Date(2026, 9, 17, 1, 1, 1, 0, time.UTC)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func login(t *testing.T, server *APIServer, password string) string {
	t.Helper()
	resp := apiRequest(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": password})
	if resp.Code != http.StatusOK {
		t.Fatalf("login code=%d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	token, ok := body["token"].(string)
	if !ok || token == "" {
		t.Fatalf("missing token: %#v", body)
	}
	return token
}

func apiRequest(t *testing.T, server *APIServer, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(data)
	}
	req := httptest.NewRequest(method, path, reader)
	if body == nil {
		req.Body = http.NoBody
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	return rec
}

func writeTestJSON(t *testing.T, path string, value any) {
	t.Helper()
	if err := atomicWriteJSON(path, value, 0600); err != nil {
		t.Fatal(err)
	}
}

func readTestJSON(t *testing.T, path string, out any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	decodeTestJSON(t, data, out)
}

func decodeTestJSON(t *testing.T, data []byte, out any) {
	t.Helper()
	if err := json.Unmarshal(data, out); err != nil {
		t.Fatalf("decode %s: %v", string(data), err)
	}
}

func appendLine(t *testing.T, path, line string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		t.Fatal(err)
	}
}
