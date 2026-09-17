package components

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
	resp := apiRequest(t, server, http.MethodGet, "/api/unknown", "", nil)
	if resp.Code != http.StatusNotFound {
		t.Fatalf("unknown code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", "", nil)
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
	resp = apiRequest(t, server, http.MethodGet, "/api/status", token, map[string]string{"unexpected": "body"})
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("body-forbidden code=%d body=%s", resp.Code, resp.Body.String())
	}
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString("{"))
	req.RemoteAddr = "192.0.2.1:1234"
	resp = httptest.NewRecorder()
	server.Handler().ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("invalid json code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPIPhase3StatusFixtures(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")
	queueEntry := map[string]any{"id": "q20260917010101", "trigger": "manual"}

	writeTestJSON(t, filepath.Join(state, ".build_state"), apiBuildState{Queued: []map[string]any{queueEntry}})
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260917010101","finished_at":"2026-09-17T01:01:03Z","status":"success","trigger":"manual","duration_seconds":2,"blob_sha":"blob-1"}`)
	before := snapshotFiles(t, state, []string{".build_state", ".build_history", ".build_status.json"})
	resp := apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("fallback status code=%d body=%s", resp.Code, resp.Body.String())
	}
	var status map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &status)
	if status["last_sha"] != "blob-1" || status["last_build_status"] != "success" || status["running"] != false || len(status["queued"].([]any)) != 1 {
		t.Fatalf("unexpected fallback status: %#v", status)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_state", ".build_history", ".build_status.json"}))

	lastBlob := "blob-primary"
	lastStatus := "success"
	lastTrigger := "manual"
	writeTestJSON(t, filepath.Join(state, ".build_status.json"), apiBuildStatus{
		LastBlobSHA: &lastBlob, LastTargetStatus: &lastStatus, LastTrigger: &lastTrigger,
		PendingTransfersCount: 2, NotifyPendingCount: 1, CircuitOpen: true,
	})
	if err := os.WriteFile(filepath.Join(state, ".build_lock"), []byte("locked"), 0600); err != nil {
		t.Fatal(err)
	}
	before = snapshotFiles(t, state, []string{".build_state", ".build_history", ".build_status.json", ".build_lock"})
	resp = apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("primary status code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &status)
	if status["last_sha"] != "blob-primary" || status["pending_transfers_count"].(float64) != 2 || status["running"] != true || len(status["queued"].([]any)) != 1 {
		t.Fatalf("unexpected primary status: %#v", status)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_state", ".build_history", ".build_status.json", ".build_lock"}))

	if err := os.WriteFile(filepath.Join(state, ".build_status.json"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	before = snapshotFiles(t, state, []string{".build_status.json"})
	resp = apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("corrupt status code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_status.json"}))
}

func TestAPIPhase3HistoryQueueAndLogFixtures(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b1","finished_at":"2026-09-17T01:01:03Z","status":"success","trigger":"manual","duration_seconds":2}`)
	appendLine(t, filepath.Join(state, ".build_history"), ``)
	appendLine(t, filepath.Join(state, ".build_history"), `{`)
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"missing-status"}`)
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b2","finished_at":"2026-09-17T01:01:04Z","status":"failure","trigger":"manual","duration_seconds":3}`)
	before := snapshotFiles(t, state, []string{".build_history"})
	resp := apiRequest(t, server, http.MethodGet, "/api/history?page=1&per_page=10", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("history code=%d body=%s", resp.Code, resp.Body.String())
	}
	var history map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &history)
	if history["total"].(float64) != 2 || history["pages"].(float64) != 1 || len(history["history"].([]any)) != 2 {
		t.Fatalf("unexpected history: %#v", history)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_history"}))

	resp = apiRequest(t, server, http.MethodGet, "/api/queue", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("missing queue code=%d body=%s", resp.Code, resp.Body.String())
	}
	var queue map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &queue)
	if len(queue["queued"].([]any)) != 0 {
		t.Fatalf("unexpected missing queue: %#v", queue)
	}
	if _, err := os.Stat(filepath.Join(state, ".build_state")); !os.IsNotExist(err) {
		t.Fatalf("GET queue must not create .build_state, err=%v", err)
	}
	if err := os.WriteFile(filepath.Join(state, ".build_state"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/queue", token, nil)
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("corrupt queue code=%d body=%s", resp.Code, resp.Body.String())
	}

	if err := os.MkdirAll(filepath.Join(state, ".build_logs", "archive"), 0755); err != nil {
		t.Fatal(err)
	}
	writeGzipJSON(t, filepath.Join(state, ".build_logs", "archive", "archived.json.gz"), apiBuildLog{
		ID: "archived", StartedAt: "2026-09-17T01:01:01Z", FinishedAt: "2026-09-17T01:01:02Z",
		TargetStatus: "success", Pipeline: apiPipelineLog{Stdout: "[INFO] archived"},
	})
	resp = apiRequest(t, server, http.MethodGet, "/api/history/archived/log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("archive log code=%d body=%s", resp.Code, resp.Body.String())
	}
	if _, err := os.Stat(filepath.Join(state, ".build_logs", "archived.json")); !os.IsNotExist(err) {
		t.Fatalf("archive lookup must not restore normal log, err=%v", err)
	}
	if err := os.WriteFile(filepath.Join(state, ".build_logs", "broken.json"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/history/broken/log", token, nil)
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("broken log code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/history/missing/log", token, nil)
	if resp.Code != http.StatusNotFound {
		t.Fatalf("missing log code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPIPhase3BuildConflictFixtures(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	if err := os.WriteFile(filepath.Join(state, ".build_lock"), []byte("locked"), 0600); err != nil {
		t.Fatal(err)
	}
	before := snapshotFiles(t, state, []string{".build_state", ".build_history", ".build_status.json", ".build_lock"})
	resp := apiRequest(t, server, http.MethodPost, "/api/build", token, nil)
	if resp.Code != http.StatusConflict {
		t.Fatalf("lock conflict code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_state", ".build_history", ".build_status.json", ".build_lock"}))
	if err := os.Remove(filepath.Join(state, ".build_lock")); err != nil {
		t.Fatal(err)
	}

	writeTestJSON(t, filepath.Join(state, ".build_circuit_state"), apiCircuitState{Open: true, ConsecutiveFailures: 2})
	before = snapshotFiles(t, state, []string{".build_state", ".build_history", ".build_status.json", ".build_circuit_state"})
	resp = apiRequest(t, server, http.MethodPost, "/api/build", token, nil)
	if resp.Code != http.StatusConflict {
		t.Fatalf("circuit conflict code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_state", ".build_history", ".build_status.json", ".build_circuit_state"}))

	resp = apiRequest(t, server, http.MethodPost, "/api/circuit-breaker/reset", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("circuit reset code=%d body=%s", resp.Code, resp.Body.String())
	}
	var circuit apiCircuitState
	readTestJSON(t, filepath.Join(state, ".build_circuit_state"), &circuit)
	if circuit.Open || circuit.ConsecutiveFailures != 0 {
		t.Fatalf("circuit reset state: %#v", circuit)
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

func TestAPIAccessControlEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	resp := apiRequest(t, server, http.MethodGet, "/api/access-control", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("default access-control code=%d body=%s", resp.Code, resp.Body.String())
	}
	var cfg apiAccessControl
	decodeTestJSON(t, resp.Body.Bytes(), &cfg)
	if len(cfg.Allow) != 0 {
		t.Fatalf("unexpected default access control: %#v", cfg)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/access-control", token, map[string]any{"allow": []string{" 192.0.2.0/24 ", "10.0.0.1", "10.0.0.1"}})
	if resp.Code != http.StatusOK {
		t.Fatalf("set access-control code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".access_control"), &cfg)
	expected := []string{"10.0.0.1", "192.0.2.0/24"}
	if !stringSlicesEqual(cfg.Allow, expected) {
		t.Fatalf("unexpected normalized allow list: %#v", cfg.Allow)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/access-control", token, map[string]any{"allow": []string{"192.0.2.0/24", "10.0.0.1"}})
	if resp.Code != http.StatusOK {
		t.Fatalf("access-control no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["message"] != "No changes" {
		t.Fatalf("expected no changes: %#v", body)
	}

	resp = apiRequestFrom(t, server, http.MethodGet, "/api/config", token, nil, "203.0.113.9:1234")
	if resp.Code != http.StatusForbidden {
		t.Fatalf("disallowed remote code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequestFrom(t, server, http.MethodGet, "/api/health", "", nil, "203.0.113.9:1234")
	if resp.Code != http.StatusOK {
		t.Fatalf("health must bypass access control: code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequestFrom(t, server, http.MethodGet, "/api/config", token, nil, "not-a-remote")
	if resp.Code != http.StatusForbidden {
		t.Fatalf("invalid remote must be forbidden: code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/access-control", token, map[string]any{"allow": []string{"2001:db8::1"}})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("ipv6 must fail: code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/config-log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("config log code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["total"].(float64) != 1 {
		t.Fatalf("unexpected access-control config log: %#v", body)
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

func TestAPIPhase4BackupWebhookSnapshotAndTokenFixtures(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	secret := "phase4-webhook-secret"
	resp := apiRequest(t, server, http.MethodPost, "/api/webhook-config", token, map[string]string{"secret": secret})
	if resp.Code != http.StatusOK {
		t.Fatalf("webhook config code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/backup", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("backup code=%d body=%s", resp.Code, resp.Body.String())
	}
	var backup map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &backup)
	if body := resp.Body.String(); strings.Contains(body, secret) {
		t.Fatalf("backup leaked secret: %s", body)
	}
	if backup["config"] == nil {
		t.Fatalf("backup missing config: %#v", backup)
	}

	before := snapshotFiles(t, state, []string{".webhook_events.json", ".build_state", ".build_history"})
	resp = signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"abc123"}`), "delivery-1", "sha256=bad")
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("invalid webhook code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".webhook_events.json", ".build_state", ".build_history"}))
	resp = signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"abc123"}`), "delivery-1", "")
	if resp.Code != http.StatusAccepted {
		t.Fatalf("signed webhook code=%d body=%s", resp.Code, resp.Body.String())
	}
	var webhookResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &webhookResp)
	if webhookResp["queued"] != true || webhookResp["queue_id"] == nil {
		t.Fatalf("unexpected webhook response: %#v", webhookResp)
	}
	var buildState apiBuildState
	readTestJSON(t, filepath.Join(state, ".build_state"), &buildState)
	if len(buildState.Queued) != 1 || buildState.Queued[0]["trigger"] != "webhook" {
		t.Fatalf("webhook did not queue build: %#v", buildState)
	}

	snapshotDir := filepath.Join(state, ".snapshots", "snap1")
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(snapshotDir, "site.tar.gz"), []byte("archive"), 0600); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(snapshotDir, "meta.json"), map[string]any{"saved_at": "2026-09-17T01:01:01Z", "size_bytes": 7})
	resp = apiRequest(t, server, http.MethodGet, "/api/snapshots", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("snapshots code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/snapshots/snap1/download", token, nil)
	if resp.Code != http.StatusOK || resp.Header().Get("Content-Type") != "application/octet-stream" {
		t.Fatalf("download code=%d type=%q body=%s", resp.Code, resp.Header().Get("Content-Type"), resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/snapshots/snap1", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("snapshot delete code=%d body=%s", resp.Code, resp.Body.String())
	}
	if _, err := os.Stat(snapshotDir); !os.IsNotExist(err) {
		t.Fatalf("snapshot dir should be deleted, err=%v", err)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/tokens", token, map[string]any{"label": "automation", "scopes": []string{"read"}})
	if resp.Code != http.StatusCreated {
		t.Fatalf("token issue code=%d body=%s", resp.Code, resp.Body.String())
	}
	var issued map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &issued)
	rawToken, _ := issued["token"].(string)
	if rawToken == "" || issued["label"] != "automation" || issued["created_at"] == "" {
		t.Fatalf("missing issued token: %#v", issued)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/tokens", token, nil)
	if resp.Code != http.StatusOK || strings.Contains(resp.Body.String(), rawToken) {
		t.Fatalf("token list leaked raw token: code=%d body=%s", resp.Code, resp.Body.String())
	}
	var stored []apiTokenRecord
	readTestJSON(t, filepath.Join(state, ".api_tokens"), &stored)
	if len(stored) != 1 || stored[0].TokenHash == "" || strings.Contains(stored[0].TokenHash, rawToken) {
		t.Fatalf("token hash not stored safely: %#v", stored)
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/tokens/missing", token, nil)
	if resp.Code != http.StatusNotFound {
		t.Fatalf("missing token revoke code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPIPhase4RulePipelineNotesAndLayoutFixtures(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	rule := map[string]any{"name": "failures", "status": "failure"}
	resp := apiRequest(t, server, http.MethodPost, "/api/alert-rules", token, rule)
	if resp.Code != http.StatusCreated {
		t.Fatalf("alert rule code=%d body=%s", resp.Code, resp.Body.String())
	}
	var createdRule map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &createdRule)
	if createdRule["id"] == "" || createdRule["rule"] != nil || createdRule["name"] != "failures" {
		t.Fatalf("alert rule response must be the rule record: %#v", createdRule)
	}
	before := snapshotFiles(t, state, []string{".alert_rules", ".config_log"})
	resp = apiRequest(t, server, http.MethodPost, "/api/alert-rules", token, rule)
	if resp.Code != http.StatusConflict {
		t.Fatalf("duplicate alert rule code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".alert_rules", ".config_log"}))

	before = snapshotFiles(t, state, []string{".pipeline_config"})
	resp = apiRequest(t, server, http.MethodPost, "/api/pipeline-config", token, map[string]any{"extra_args": []string{"--src"}})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("reserved pipeline arg code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".pipeline_config"}))

	resp = apiRequest(t, server, http.MethodPost, "/api/notes", token, map[string]string{"content": "release note"})
	if resp.Code != http.StatusOK {
		t.Fatalf("notes code=%d body=%s", resp.Code, resp.Body.String())
	}
	before = snapshotFiles(t, state, []string{".notes", ".config_log"})
	resp = apiRequest(t, server, http.MethodPost, "/api/notes", token, map[string]string{"content": "release note"})
	if resp.Code != http.StatusOK {
		t.Fatalf("notes no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["message"] != "No changes" {
		t.Fatalf("unexpected notes no-op: %#v", body)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".notes", ".config_log"}))

	resp = apiRequest(t, server, http.MethodPost, "/api/dashboard-layout", token, map[string]any{"widgets": []string{"status", "status"}})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("dashboard duplicate code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/smtp-test", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("smtp disabled code=%d body=%s", resp.Code, resp.Body.String())
	}
	if _, err := os.Stat(filepath.Join(state, ".notify_log")); !os.IsNotExist(err) {
		t.Fatalf("smtp disabled must not write notify log, err=%v", err)
	}
}

func TestAPICompletionEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "admin")

	appendLine(t, filepath.Join(state, ".audit_log"), `{"at":"2026-09-17T01:01:02Z","actor":"admin","action":"config_update","result":"success"}`)
	resp := apiRequest(t, server, http.MethodGet, "/api/audit-log?action=config_update", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("audit-log code=%d body=%s", resp.Code, resp.Body.String())
	}
	var logResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &logResp)
	if logResp["total"].(float64) != 1 {
		t.Fatalf("unexpected audit log: %#v", logResp)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/auth/totp-status", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("totp status code=%d body=%s", resp.Code, resp.Body.String())
	}
	var status map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &status)
	if status["enabled"] != false {
		t.Fatalf("unexpected initial totp status: %#v", status)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/auth/totp-setup", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("totp setup code=%d body=%s", resp.Code, resp.Body.String())
	}
	var setup map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &setup)
	secret, _ := setup["secret"].(string)
	code := totpCode(secret, server.cfg.Now().UTC().Unix()/30)
	resp = apiRequest(t, server, http.MethodPost, "/api/auth/totp-confirm", token, map[string]string{"code": code})
	if resp.Code != http.StatusOK {
		t.Fatalf("totp confirm code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "admin"})
	if resp.Code != http.StatusOK {
		t.Fatalf("totp login stage1 code=%d body=%s", resp.Code, resp.Body.String())
	}
	var loginResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &loginResp)
	ticket, _ := loginResp["ticket"].(string)
	if ticket == "" || loginResp["totp_required"] != true {
		t.Fatalf("login did not require totp: %#v", loginResp)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/sessions", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("sessions after totp ticket code=%d body=%s", resp.Code, resp.Body.String())
	}
	var sessionsAfterTicket map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &sessionsAfterTicket)
	if len(sessionsAfterTicket["sessions"].([]any)) != 1 {
		t.Fatalf("totp stage1 must not create a session: %#v", sessionsAfterTicket)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/login/totp", "", map[string]string{"ticket": ticket, "code": code})
	if resp.Code != http.StatusOK {
		t.Fatalf("totp login stage2 code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/auth/totp", token, map[string]string{"code": code})
	if resp.Code != http.StatusOK {
		t.Fatalf("totp disable code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/log-level", token, map[string]string{"level": "debug"})
	if resp.Code != http.StatusOK {
		t.Fatalf("log-level code=%d body=%s", resp.Code, resp.Body.String())
	}
	var cfg apiServerConfig
	readTestJSON(t, filepath.Join(state, ".server_config"), &cfg)
	if cfg.LogLevel != "DEBUG" {
		t.Fatalf("log level was not saved: %#v", cfg)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/api-rate-limit", token, map[string]any{"enabled": true, "groups": []any{}})
	if resp.Code != http.StatusOK {
		t.Fatalf("api-rate-limit post code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/api-rate-limit", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("api-rate-limit get code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/pat-verify", token, nil)
	if resp.Code != http.StatusNotImplemented {
		t.Fatalf("pat verify without token code=%d body=%s", resp.Code, resp.Body.String())
	}
	if err := os.WriteFile(filepath.Join(state, ".github_token"), []byte("ghp_exampletoken\n"), 0600); err != nil {
		t.Fatal(err)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/pat-verify", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("pat verify code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/rate-limit", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("rate-limit code=%d body=%s", resp.Code, resp.Body.String())
	}

	hook := map[string]any{"phase": "after_build", "command_args": []string{"/bin/echo", "ok"}, "abort_on_failure": false}
	resp = apiRequest(t, server, http.MethodPost, "/api/hooks", token, hook)
	if resp.Code != http.StatusCreated {
		t.Fatalf("hook create code=%d body=%s", resp.Code, resp.Body.String())
	}
	var hookResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &hookResp)
	hookID, _ := hookResp["id"].(string)
	if hookID == "" || hookResp["phase"] != "after_build" || hookResp["enabled"] != true || hookResp["timeout_seconds"].(float64) != 300 {
		t.Fatalf("hook response must be the hook record: %#v", hookResp)
	}
	writeTestJSON(t, filepath.Join(state, ".build_logs", "b1_hook_"+hookID+".json"), map[string]any{"at": "2026-09-17T01:02:00Z", "result": "success"})
	resp = apiRequest(t, server, http.MethodGet, "/api/hooks/"+hookID+"/log", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("hook log code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/hooks/"+hookID, token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("hook delete code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/tag-rules", token, map[string]any{"condition": map[string]any{"status": "failure"}, "tags": []string{"failure"}})
	if resp.Code != http.StatusCreated {
		t.Fatalf("tag rule create code=%d body=%s", resp.Code, resp.Body.String())
	}
	var tagRule map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &tagRule)
	if tagRule["id"] == "" || tagRule["condition"] == nil || tagRule["rule"] != nil {
		t.Fatalf("tag rule response must be the tag rule record: %#v", tagRule)
	}

	siteDir := filepath.Join(state, "site")
	if err := os.MkdirAll(siteDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(siteDir, "index.html"), []byte("<h1>Adlaire</h1>\n"), 0600); err != nil {
		t.Fatal(err)
	}
	meta, err := inspectOutput(siteDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(state, ".build_logs"), 0755); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(state, ".build_logs", "bverify.json"), apiBuildLog{ID: "bverify", StartedAt: "2026-09-17T01:01:01Z", FinishedAt: "2026-09-17T01:01:03Z", TargetStatus: "success", OutputSHA256: meta.SHA256})
	resp = apiRequest(t, server, http.MethodPost, "/api/verify-output", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("verify-output code=%d body=%s", resp.Code, resp.Body.String())
	}
	var verify map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &verify)
	if verify["match"] != true {
		t.Fatalf("unexpected verify result: %#v", verify)
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
	return apiRequestFrom(t, server, method, path, token, body, "192.0.2.1:1234")
}

func apiRequestFrom(t *testing.T, server *APIServer, method, path, token string, body any, remoteAddr string) *httptest.ResponseRecorder {
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
	req.RemoteAddr = remoteAddr
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	return rec
}

func signedWebhookRequest(t *testing.T, server *APIServer, secret string, body []byte, deliveryID, signature string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/webhook", bytes.NewReader(body))
	req.RemoteAddr = "192.0.2.1:1234"
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", deliveryID)
	if signature == "" {
		mac := hmac.New(sha256.New, []byte(secret))
		_, _ = mac.Write(body)
		signature = "sha256=" + hex.EncodeToString(mac.Sum(nil))
	}
	req.Header.Set("X-Hub-Signature-256", signature)
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

type fileSnapshot struct {
	exists bool
	mode   os.FileMode
	mtime  time.Time
	data   string
}

func snapshotFiles(t *testing.T, root string, names []string) map[string]fileSnapshot {
	t.Helper()
	out := map[string]fileSnapshot{}
	for _, name := range names {
		path := filepath.Join(root, name)
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			out[name] = fileSnapshot{}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		out[name] = fileSnapshot{exists: true, mode: info.Mode(), mtime: info.ModTime(), data: string(data)}
	}
	return out
}

func assertSnapshotEqual(t *testing.T, before, after map[string]fileSnapshot) {
	t.Helper()
	if len(before) != len(after) {
		t.Fatalf("snapshot size changed: before=%#v after=%#v", before, after)
	}
	for name, want := range before {
		got, ok := after[name]
		if !ok {
			t.Fatalf("snapshot missing path %s", name)
		}
		if want != got {
			t.Fatalf("snapshot changed for %s\nbefore=%#v\nafter=%#v", name, want, got)
		}
	}
}

func writeGzipJSON(t *testing.T, path string, value any) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	writer := gzip.NewWriter(file)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
}
