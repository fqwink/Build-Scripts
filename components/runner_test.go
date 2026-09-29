package components

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunnerFixtureR1CLI(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := RunRunner([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("help exit=%d", code)
	}
	if strings.TrimSpace(stdout.String()) != "Usage: adlaire-ci-runner [--state-dir path] [--once] [--dry-run] [--version] [--help]" || stderr.Len() != 0 {
		t.Fatalf("unexpected help stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunRunner([]string{"--state-dir", "relative"}, &stdout, &stderr); code != 2 {
		t.Fatalf("relative exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "state directory must be absolute: relative") {
		t.Fatalf("missing relative error: %s", stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunRunner([]string{"--unknown"}, &stdout, &stderr); code != 2 {
		t.Fatalf("unknown exit=%d", code)
	}
	if !strings.Contains(stderr.String(), "unknown option: --unknown") {
		t.Fatalf("missing unknown error: %s", stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunRunner([]string{"--version", "--state-dir", "relative"}, &stdout, &stderr); code != 0 {
		t.Fatalf("version exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout.String()), "adlaire-ci-runner V.0.0-dev go=") || stderr.Len() != 0 {
		t.Fatalf("unexpected version stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunnerFixtureR2NoChange(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"blob-1\"}\n" {
			t.Fatalf("sha changed")
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
			t.Fatalf("history must not exist")
		}
		if !strings.Contains(stdout.String(), "NO_CHANGE: branch=main target=docs sha=blob-1") {
			t.Fatalf("missing no change log: %s", stdout.String())
		}
	})
}

func TestRunnerFixtureR3BuildSuccessNoDeploy(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"new-blob\"}\n" {
			t.Fatalf("sha not updated: %s", readFile(t, filepath.Join(state, ".last_sha")))
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "success" || log.Pipeline.ExitCode == nil || *log.Pipeline.ExitCode != 0 || log.Report == nil || log.Report.Pages != 1 || len(log.Deploy) != 0 || log.Error != nil {
			t.Fatalf("unexpected log: %+v", log)
		}
		if log.SnapshotID == nil {
			t.Fatalf("snapshot id must be recorded")
		}
		if log.Status != "success" || log.OutputSHA256 == nil || log.OutputSizeBytes == nil || log.Trigger != "polling" {
			t.Fatalf("normalized log fields missing: %+v", log)
		}
		if log.Environment.RunnerVersion != "V.0.0-dev" || log.Environment.GoVersion == "" || log.Environment.OS == "" || log.Environment.Arch == "" {
			t.Fatalf("environment record missing: %+v", log.Environment)
		}
		if _, err := os.Stat(filepath.Join(state, ".snapshots", *log.SnapshotID, "site", "index.html")); err != nil {
			t.Fatalf("snapshot missing: %v", err)
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".build_history")), `"status":"success"`) {
			t.Fatalf("history missing success")
		}
		var summary buildStatusSummary
		readJSON(t, filepath.Join(state, ".build_status.json"), &summary)
		if summary.Status != "success" || summary.Running || summary.LastBuildID == nil || *summary.LastBuildID != log.ID || summary.OutputSHA256 == nil {
			t.Fatalf("unexpected build status summary: %+v", summary)
		}
		var trends buildTrendsFile
		readJSON(t, filepath.Join(state, ".build_trends.json"), &trends)
		if trends.SchemaVersion != 1 || len(trends.Samples) != 1 || trends.Samples[0].BuildID != log.ID || trends.Summary.Count != 1 || trends.Summary.AvgSeconds == nil {
			t.Fatalf("unexpected build trends: %+v", trends)
		}
		var bs buildState
		readJSON(t, filepath.Join(state, ".build_state"), &bs)
		if bs.Running || bs.CurrentBuildID != nil {
			t.Fatalf("state not finalized: %+v", bs)
		}
		if _, err := os.Stat(filepath.Join(state, ".build_lock")); !os.IsNotExist(err) {
			t.Fatalf("lock remains")
		}
	})
}

func TestRunnerFixtureR4PipelineFailure(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 7, "before fail", "failed")
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"old-blob\"}\n" {
			t.Fatalf("sha changed on failure")
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "failure_build" || log.Pipeline.ExitCode == nil || *log.Pipeline.ExitCode != 7 || log.Pipeline.Stdout != "before fail" || log.Pipeline.Stderr != "failed" || log.Error == nil || *log.Error != "pipeline failed" {
			t.Fatalf("unexpected failure log: %+v", log)
		}
		if log.FailureCategory == nil || *log.FailureCategory != "build_failure" || log.Status != "failure" {
			t.Fatalf("failure category not recorded: %+v", log)
		}
	})
}

func TestRunnerCompletionServerConfigSparseSchema(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"log_max_lines":              750,
		"history_max_count":          500,
		"history_retention":          map[string]any{"enabled": false, "max_count": 1000, "max_age_days": nil, "updated_at": nil},
		"build_timeout_seconds":      120,
		"log_retention_days":         30,
		"log_level":                  "INFO",
		"pat_expires_at":             nil,
		"snapshots_keep":             0,
		"queue_max_size":             3,
		"build_retry_max":            1,
		"build_retry_base_seconds":   5,
		"commit_status_enabled":      false,
		"commit_status_context":      "Adlaire CI",
		"commit_status_target_url":   nil,
		"log_archive_after_days":     0,
		"build_trend_keep_count":     1000,
		"duration_anomaly":           map[string]any{"enabled": false, "min_samples": 20, "avg_multiplier": 2.0, "p95_multiplier": 1.5},
		"watch_mode":                 "github",
		"tag_filter":                 map[string]any{"enabled": false, "patterns": []any{}},
		"build_cache_enabled":        false,
		"deploy_parallelism":         1,
		"remote_build":               map[string]any{"enabled": false, "host": nil, "user": nil, "work_dir": nil, "command_args": []any{}, "artifact_path": nil},
		"approval_timeout_seconds":   86400,
		"force_build_interval_hours": 0,
		"build_cooldown_seconds":     0,
		"schedule_interval_seconds":  300,
		"schedule_paused":            false,
		"allowed_hours":              nil,
		"session_timeout_seconds":    28800,
		"api_rate_limit": map[string]any{"enabled": true, "groups": map[string]any{
			"login": map[string]any{"window_seconds": 60, "max_requests": 10},
		}},
	})
	cfg := defaultRunnerConfig(state)
	if err := loadRunnerConfig(&cfg, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))); err != nil {
		t.Fatalf("valid sparse server config rejected: %v", err)
	}
	if cfg.BuildTimeoutSeconds != 120 || cfg.SnapshotsKeep != 0 || cfg.BuildRetryMax != 1 || cfg.WatchMode != "github" || cfg.DeployParallelism != 1 {
		t.Fatalf("server config not applied: %+v", cfg)
	}
}

func TestRunnerCompletionServerConfigRejectsUnknownKey(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{"unknown": true})
	cfg := defaultRunnerConfig(state)
	err := loadRunnerConfig(&cfg, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	if err == nil || !strings.Contains(err.Error(), "SERVER_CONFIG_INVALID") {
		t.Fatalf("unknown server config key accepted: %v", err)
	}
}

func TestRunnerCompletionDryRunNoWrites(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	for _, path := range []string{filepath.Join(state, "repo"), filepath.Join(state, "dist"), filepath.Join(state, ".build_logs"), filepath.Join(state, ".build_lock")} {
		_ = os.RemoveAll(path)
	}
	var stdout, stderr bytes.Buffer
	if code := RunRunner([]string{"--state-dir", state, "--dry-run"}, &stdout, &stderr); code != 0 {
		t.Fatalf("dry-run exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	var result map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &result); err != nil {
		t.Fatalf("invalid dry-run json: %v stdout=%s", err, stdout.String())
	}
	if result["dry_run"] != true {
		t.Fatalf("dry_run flag missing: %+v", result)
	}
	for _, path := range []string{filepath.Join(state, "repo"), filepath.Join(state, "dist"), filepath.Join(state, ".build_logs"), filepath.Join(state, ".build_lock"), filepath.Join(state, ".build_history")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("dry-run wrote path %s", path)
		}
	}
}

func TestRunnerCompletionLocalWatchDoesNotRequireToken(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	if err := os.Remove(filepath.Join(state, ".github_token")); err != nil {
		t.Fatal(err)
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"watch_mode":             "local",
		"build_cooldown_seconds": 0,
	})
	src := filepath.Join(state, "repo", "docs")
	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "source.md"), []byte("# Local\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	log := onlyBuildLog(t, state)
	if log.Environment.WatchMode != "local" || log.BlobSHA == nil || *log.BlobSHA == "old-blob" {
		t.Fatalf("local watch log mismatch: %+v", log)
	}
}

func TestRunnerCompletionMultiFileChangedTargets(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipelineWithExtra(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "", `if [ "${ADLAIRE_CHANGED_TARGETS:-}" != '["docs/a.md","docs/b.md"]' ]; then exit 8; fi`)
	target := BranchTarget{
		Branch: "main", TargetFile: "docs/a.md", TargetFiles: []string{"docs/b.md", "docs/a.md"},
		SHAFile: filepath.Join(state, ".last_sha"), Src: filepath.Join(state, "repo", "docs"), Out: filepath.Join(state, "dist", "site"),
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".branch_config"), branchConfigFile{BranchTargets: []BranchTarget{target}})
	encodedA := base64.StdEncoding.EncodeToString([]byte("# A\n"))
	encodedB := base64.StdEncoding.EncodeToString([]byte("# B\n"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{
				{"path": "docs/a.md", "type": "blob", "sha": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
				{"path": "docs/b.md", "type": "blob", "sha": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
			}})
		case strings.Contains(r.URL.Path, "/git/blobs/aaaaaaaa"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encodedA, "encoding": "base64"})
		case strings.Contains(r.URL.Path, "/git/blobs/bbbbbbbb"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encodedB, "encoding": "base64"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if strings.Join(log.TargetFiles, ",") != "docs/a.md,docs/b.md" || strings.Join(log.ChangedTargets, ",") != "docs/a.md,docs/b.md" {
			t.Fatalf("target tracking mismatch: %+v", log)
		}
		if _, err := os.Stat(filepath.Join(state, "repo", "docs", "docs", "a.md")); err != nil {
			t.Fatalf("source a not materialized: %v", err)
		}
		if readFile(t, filepath.Join(state, ".last_sha")) == "{\"sha\":\"old-blob\"}\n" {
			t.Fatalf("multi-file digest not updated")
		}
	})
}

func TestRunnerCompletionApprovalRequiredSkipsBuild(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	target := BranchTarget{
		Branch: "main", TargetFile: "docs", TargetFiles: []string{"docs"},
		SHAFile: filepath.Join(state, ".last_sha"), Src: filepath.Join(state, "repo", "docs"), Out: filepath.Join(state, "dist", "site"),
		ApprovalRequired: true,
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".branch_config"), branchConfigFile{BranchTargets: []BranchTarget{target}})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
			t.Fatalf("approval pending must not create history")
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"old-blob\"}\n" {
			t.Fatalf("approval pending updated sha")
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".approval_queue")), `"status":"pending"`) {
			t.Fatalf("approval queue missing")
		}
		var summary buildStatusSummary
		readJSON(t, filepath.Join(state, ".build_status.json"), &summary)
		if summary.Status != "skipped" || summary.LastTargetStatus == nil || *summary.LastTargetStatus != "pending_approval" {
			t.Fatalf("approval status mismatch: %+v", summary)
		}
	})
}

func TestRunnerCompletionReportDuplicateUsesFirst(t *testing.T) {
	report, warnings := parseRunnerReport("[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=first\n[REPORT] pages=2 headings=2 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=second\n")
	if report == nil || report.Pages != 1 || report.Theme != "first" {
		t.Fatalf("first report not preserved: %+v", report)
	}
	if len(warnings) != 1 || warnings[0] != "REPORT_DUPLICATE" {
		t.Fatalf("missing duplicate warning: %+v", warnings)
	}
}

func TestRunnerCompletionTrimLogUTF8Safe(t *testing.T) {
	raw := strings.Repeat("a", 1024*1024-1) + "あいう"
	trimmed, truncated := trimLog(raw)
	if !truncated {
		t.Fatalf("expected truncation")
	}
	if !strings.Contains(trimmed, "あいう") || strings.Contains(trimmed, "\uFFFD") {
		t.Fatalf("trimmed log is not utf8 safe: suffix=%q", trimmed[len(trimmed)-12:])
	}
}

func TestRunnerCompletionArchivesOldBuildLogs(t *testing.T) {
	state := t.TempDir()
	dir := filepath.Join(state, ".build_logs")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(dir, "old.json")
	newer := filepath.Join(dir, "new.json")
	if err := os.WriteFile(old, []byte(`{"id":"old"}`+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newer, []byte(`{"id":"new"}`+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	oldNow := runnerNow
	runnerNow = func() time.Time { return time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC) }
	defer func() { runnerNow = oldNow }()
	oldTime := runnerNow().Add(-48 * time.Hour)
	if err := os.Chtimes(old, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}
	newTime := runnerNow()
	if err := os.Chtimes(newer, newTime, newTime); err != nil {
		t.Fatal(err)
	}
	cleanupBuildLogs(dir, 1, 1, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatalf("old log not removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "archive", "old.json.gz")); err != nil {
		t.Fatalf("old log not archived: %v", err)
	}
	if _, err := os.Stat(newer); err != nil {
		t.Fatalf("new log removed: %v", err)
	}
}

func TestRunnerHardeningRetriesGitHubAPI(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	hits := 0
	encoded := base64.StdEncoding.EncodeToString([]byte("# Title\n"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/git/trees/") {
			hits++
			if hits == 1 {
				http.Error(w, "temporary", http.StatusInternalServerError)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": "docs", "type": "blob", "sha": "retry-blob"}}})
			return
		}
		if strings.Contains(r.URL.Path, "/git/blobs/") {
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
			return
		}
		http.NotFound(w, r)
	}))
	oldSleep := runnerSleep
	runnerSleep = func(time.Duration) {}
	defer func() { runnerSleep = oldSleep }()
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if hits != 2 {
			t.Fatalf("expected retry hits=2 got %d", hits)
		}
	})
}

func TestRunnerHardeningPrecheckFailure(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	t.Setenv("ADLAIRE_CI_BUILD_BIN", filepath.Join(state, "missing-build-bin"))
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "failure_precheck" || log.Error == nil || *log.Error != "precheck failed" {
			t.Fatalf("unexpected precheck log: %+v", log)
		}
	})
}

func TestRunnerHardeningNotifyPendingRetry(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	received := 0
	hook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received++
		w.WriteHeader(http.StatusNoContent)
	}))
	defer hook.Close()
	writeRunnerTestJSON(t, filepath.Join(state, ".notify_pending"), []notifyPendingEntry{{
		Event: "success", URL: hook.URL, Payload: map[string]any{"ok": true},
		QueuedAt: "2026-09-16T00:00:00Z", RetryCount: 1,
	}})
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if received != 1 {
			t.Fatalf("expected one notify retry, got %d", received)
		}
		if strings.TrimSpace(readFile(t, filepath.Join(state, ".notify_pending"))) != "[]" {
			t.Fatalf("notify pending must be empty")
		}
	})
}

func TestRunnerHardeningCircuitOpenSkipsPolling(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	now := "2026-09-16T00:00:00Z"
	msg := "github api failed"
	writeRunnerTestJSON(t, filepath.Join(state, ".build_circuit_state"), buildCircuitState{Open: true, ConsecutiveFailures: 3, OpenedAt: &now, LastFailureAt: &now, LastError: &msg})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
			t.Fatalf("history must not be touched while circuit is open")
		}
		if !strings.Contains(stdout.String(), "CIRCUIT_OPEN") {
			t.Fatalf("missing circuit log: %s", stdout.String())
		}
	})
}

func TestRunnerCompletionPendingTransferRetrySuccess(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	out := filepath.Join(state, "dist", "site")
	if err := os.MkdirAll(out, 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(out, "index.html")
	if err := os.WriteFile(file, []byte("<html></html>"), 0644); err != nil {
		t.Fatal(err)
	}
	sha, err := runnerFileSHA256(file)
	if err != nil {
		t.Fatal(err)
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".pending_transfers"), []pendingTransfer{{BuildID: "b20260916000000", Trigger: "deploy", SourceKind: "output", Branch: "main", TargetID: "deploy-1", Out: &out, Host: "host", User: "deploy", DestDir: "/var/www/html", OutputSHA256: sha, FailedAt: "2026-09-16T00:00:00Z", RetryCount: 1, LastError: "previous failure"}})
	fakeBin := t.TempDir()
	writeExecutable(t, filepath.Join(fakeBin, "ssh"), "#!/bin/sh\nif [ \"$2\" = \"sha256sum\" ]; then printf '"+sha+"  file\\n'; exit 0; fi\ncat >/dev/null\nexit 0\n")
	t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if strings.TrimSpace(readFile(t, filepath.Join(state, ".pending_transfers"))) != "[]" {
			t.Fatalf("pending transfers not cleared")
		}
	})
}

func TestRunnerCompletionDeployUploadsAndVerifies(t *testing.T) {
	out := t.TempDir()
	file := filepath.Join(out, "index.html")
	if err := os.WriteFile(file, []byte("<html></html>"), 0644); err != nil {
		t.Fatal(err)
	}
	sha, err := runnerFileSHA256(file)
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "uploaded")
	fakeBin := t.TempDir()
	script := fmt.Sprintf("#!/bin/sh\nif [ \"$2\" = \"sha256sum\" ]; then if [ -f %q ]; then printf '%s  %%s\\n' \"$3\"; else printf 'old  %%s\\n' \"$3\"; fi; exit 0; fi\nif [ \"$2\" = \"mkdir\" ]; then cat >/dev/null; touch %q; exit 0; fi\nexit 1\n", marker, sha, marker)
	writeExecutable(t, filepath.Join(fakeBin, "ssh"), script)
	t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	total, uploaded, skipped, bytesUploaded, err := deploySite(out, DeployTarget{Host: "host", User: "deploy", DestDir: "/var/www/html"})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || uploaded != 1 || skipped != 0 || bytesUploaded != int64(len("<html></html>")) {
		t.Fatalf("unexpected deploy stats total=%d uploaded=%d skipped=%d bytes=%d", total, uploaded, skipped, bytesUploaded)
	}
}

func TestRunnerCompletionCommitInfoAndNotification(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	notified := 0
	hook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notified++
		w.WriteHeader(http.StatusNoContent)
	}))
	defer hook.Close()
	writeRunnerTestJSON(t, filepath.Join(state, ".notify_config"), notifyConfigFile{Webhooks: []notifyWebhook{{URL: hook.URL, Enabled: true, On: []string{"success"}}}})
	server := fakeGitHubWithCommit(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.Commit.SHA == nil || *log.Commit.SHA != "commit-sha" || log.Commit.Message == nil || *log.Commit.Message != "Update docs" {
			t.Fatalf("commit info not recorded: %+v", log.Commit)
		}
		if notified != 1 {
			t.Fatalf("expected one notification, got %d", notified)
		}
	})
}

func TestRunnerPhase2TagFilterBuildsMatchingTag(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"tag_filter": map[string]any{"enabled": true, "patterns": []string{"v*"}},
	})
	server := fakeGitHubWithTags(t, "docs", "new-blob", "# Title\n", []map[string]string{{"ref": "refs/tags/v1.0.0", "sha": "commit-sha"}})
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "success" || strings.Join(log.MatchedTags, ",") != "v1.0.0" {
			t.Fatalf("tag filter log mismatch: %+v", log)
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"new-blob\"}\n" {
			t.Fatalf("sha not updated after matching tag")
		}
	})
}

func TestRunnerPhase2TagFilterSkipsNonMatchingTag(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"tag_filter": map[string]any{"enabled": true, "patterns": []string{"release-*"}},
	})
	server := fakeGitHubWithTags(t, "docs", "new-blob", "# Title\n", []map[string]string{{"ref": "refs/tags/v1.0.0", "sha": "commit-sha"}})
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
			t.Fatalf("tag-filter skip must not create history")
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"old-blob\"}\n" {
			t.Fatalf("sha updated despite non-matching tag")
		}
		var summary buildStatusSummary
		readJSON(t, filepath.Join(state, ".build_status.json"), &summary)
		if summary.Status != "skipped" || summary.LastTargetStatus == nil || *summary.LastTargetStatus != "skipped_tag_filter" {
			t.Fatalf("tag-filter skip status mismatch: %+v", summary)
		}
	})
}

func TestRunnerPhase2RemoteBuild(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	workDir := t.TempDir()
	fakeBin := t.TempDir()
	writeExecutable(t, filepath.Join(fakeBin, "ssh"), "#!/usr/bin/env bash\nset -euo pipefail\neval \"$2\"\n")
	t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"remote_build": map[string]any{
			"enabled": true, "host": "host", "user": "deploy", "work_dir": workDir,
			"command_args": []string{"sh", "-c", "mkdir -p \"$ADLAIRE_CI_OUT/assets\" && printf '<html></html>' > \"$ADLAIRE_CI_OUT/index.html\" && printf 'body{}' > \"$ADLAIRE_CI_OUT/assets/style.css\" && printf 'console.log(\"ok\")' > \"$ADLAIRE_CI_OUT/assets/app.js\" && printf '[]' > \"$ADLAIRE_CI_OUT/assets/search-index.json\" && printf '[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=remote'"},
		},
	})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if !log.Environment.RemoteBuild || log.Report == nil || log.Report.Theme != "remote" || log.TargetStatus != "success" {
			t.Fatalf("remote build log mismatch: %+v", log)
		}
	})
}

func TestRunnerPhase2WeeklySummaryChannel(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	received := 0
	hook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received++
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("summary payload decode: %v", err)
		}
		if payload["event"] != "weekly_summary" || payload["success_count"] == nil {
			t.Fatalf("summary payload mismatch: %+v", payload)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer hook.Close()
	writeRunnerTestJSON(t, filepath.Join(state, ".notify_config"), notifyConfigFile{
		Channels: []notifyChannel{{ID: "weekly", Type: "webhook", Enabled: true, On: []string{"weekly_summary"}, Config: map[string]any{"url": hook.URL}}},
		Summary:  notifySummary{Enabled: true, Interval: "weekly", Hour: 9, DayOfWeek: 0},
	})
	if err := appendRunnerJSONLine(filepath.Join(state, ".build_history"), historyRecord{ID: "b1", Branch: "main", TargetFile: "docs", Status: "success", FinishedAt: "2026-09-27T08:00:00Z", DurationSeconds: 4}, 0600); err != nil {
		t.Fatal(err)
	}
	oldNow := runnerNow
	runnerNow = func() time.Time { return time.Date(2026, 9, 27, 9, 0, 0, 0, time.UTC) }
	defer func() { runnerNow = oldNow }()
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if received != 1 {
			t.Fatalf("weekly summary not sent: %d", received)
		}
		var bs buildState
		readJSON(t, filepath.Join(state, ".build_state"), &bs)
		if bs.WeeklySummarySentDate == nil || *bs.WeeklySummarySentDate != "2026-09-27" {
			t.Fatalf("summary date not recorded: %+v", bs)
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".notify_log")), `"event":"weekly_summary"`) {
			t.Fatalf("notify log missing weekly summary")
		}
	})
}

func TestRunnerPhase2DurationAnomaly(t *testing.T) {
	state := t.TempDir()
	cfg := defaultRunnerConfig(state)
	cfg.DurationAnomaly = DurationAnomalyConfig{Enabled: true, MinSamples: 2, AvgMultiplier: 2, P95Multiplier: 2}
	writeRunnerTestJSON(t, filepath.Join(state, ".build_trends.json"), buildTrendsFile{SchemaVersion: 1, Samples: []buildTrendSample{
		{BuildID: "a", FinishedAt: "2026-09-16T00:00:00Z", DurationSeconds: 10, Status: "success", TargetStatus: "success"},
		{BuildID: "b", FinishedAt: "2026-09-16T00:01:00Z", DurationSeconds: 10, Status: "success", TargetStatus: "success"},
	}})
	if err := updateBuildTrends(cfg, buildLog{ID: "c", FinishedAt: "2026-09-16T00:02:00Z", DurationSeconds: 100, Status: "success", TargetStatus: "success"}); err != nil {
		t.Fatal(err)
	}
	var trends buildTrendsFile
	readJSON(t, filepath.Join(state, ".build_trends.json"), &trends)
	if len(trends.Samples) != 3 || !trends.Samples[2].Anomaly || trends.Summary.AnomalyCount != 1 {
		t.Fatalf("duration anomaly not recorded: %+v", trends)
	}
}

func TestRunnerPhase2PriorityQueueRunsUrgentBeforeNormal(t *testing.T) {
	state := newRunnerState(t, "same-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"schedule_paused": true,
	})
	writeRunnerTestJSON(t, filepath.Join(state, ".build_state"), buildState{Queued: []map[string]any{
		{"id": "q-normal", "trigger": "manual", "queued_at": "2026-09-16T00:00:00Z", "requested_by": "admin", "priority": "normal", "created_seq": 1, "payload": map[string]any{"force": true}},
		{"id": "q-urgent", "trigger": "manual", "queued_at": "2026-09-16T00:00:01Z", "requested_by": "admin", "priority": "urgent", "created_seq": 2, "payload": map[string]any{"force": true}},
	}})
	server := fakeGitHub(t, "docs", "same-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.Trigger != "manual" || log.TargetStatus != "success" {
			t.Fatalf("queue build log mismatch: %+v", log)
		}
		var bs buildState
		readJSON(t, filepath.Join(state, ".build_state"), &bs)
		if bs.ActiveQueueEntry != nil || len(bs.Queued) != 1 || fmt.Sprint(bs.Queued[0]["id"]) != "q-normal" {
			t.Fatalf("queue state mismatch: %+v", bs)
		}
	})
}

func TestRunnerCompletionCommitStatus(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"build_cooldown_seconds": 0,
		"commit_status_enabled":  true,
		"commit_status_context":  "Adlaire CI",
	})
	encoded := base64.StdEncoding.EncodeToString([]byte("# Title\n"))
	statuses := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": "docs", "type": "blob", "sha": "new-blob"}}})
		case strings.Contains(r.URL.Path, "/git/blobs/"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
		case strings.Contains(r.URL.Path, "/commits"):
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"sha": "abcdef0123456789abcdef0123456789abcdef01",
				"commit": map[string]any{
					"message": "Update docs",
					"author":  map[string]string{"name": "A. Developer", "date": "2026-09-16T00:00:00Z"},
				},
			}})
		case strings.Contains(r.URL.Path, "/statuses/"):
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("status payload decode: %v", err)
			}
			statuses = append(statuses, fmt.Sprint(payload["state"]))
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if got := strings.Join(statuses, ","); got != "pending,success" {
			t.Fatalf("unexpected commit statuses: %s", got)
		}
		log := onlyBuildLog(t, state)
		if log.CommitStatus == nil || *log.CommitStatus != "success" {
			t.Fatalf("commit status not recorded in log: %+v", log.CommitStatus)
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".build_history")), `"commit_status_state":"success"`) {
			t.Fatalf("commit status not recorded in history")
		}
	})
}

func TestRunnerCompletionBuildRetrySucceeds(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	counter := filepath.Join(state, "retry-count")
	buildBin := filepath.Join(t.TempDir(), "adlaire-ci-build")
	writeExecutable(t, buildBin, fmt.Sprintf(`#!/usr/bin/env bash
set -u
if [ "${1:-}" = "--version" ]; then
  printf 'adlaire-ci-build V.0.0-dev go=fake\n'
  exit 0
fi
out=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --out)
      shift
      out="${1:-}"
      ;;
  esac
  shift || true
done
if [ ! -f %q ]; then
  printf '1' > %q
  printf 'first failure'
  exit 7
fi
mkdir -p "$out/assets"
printf '<html></html>' > "$out/index.html"
printf 'body{}' > "$out/assets/style.css"
printf 'console.log("ok")' > "$out/assets/app.js"
printf '[]' > "$out/assets/search-index.json"
printf '[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default'
exit 0
`, counter, counter))
	t.Setenv("ADLAIRE_CI_BUILD_BIN", buildBin)
	oldSleep := runnerSleep
	runnerSleep = func(time.Duration) {}
	defer func() { runnerSleep = oldSleep }()
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		cfg := defaultRunnerConfig(state)
		cfg.BuildRetryMax = 1
		cfg.BuildRetryBaseSeconds = 1
		cfg.BuildCooldownSeconds = 0
		var stdout, stderr bytes.Buffer
		if code := executeRunner(cfg, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.RetryCount != 1 || len(log.Attempts) != 2 || log.Attempts[0].ExitCode == nil || *log.Attempts[0].ExitCode != 7 || log.TargetStatus != "success" {
			t.Fatalf("retry log mismatch: %+v", log)
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".build_history")), `"retry_count":1`) {
			t.Fatalf("history missing retry count")
		}
	})
}

func TestRunnerCompletionCooldownSkips(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	now := time.Date(2026, 9, 16, 1, 2, 3, 0, time.UTC)
	nowText := now.Format(time.RFC3339)
	writeRunnerTestJSON(t, filepath.Join(state, ".build_state"), buildState{Queued: []map[string]any{}, LastFinishedAt: &nowText})
	oldNow := runnerNow
	runnerNow = func() time.Time { return now.Add(10 * time.Second) }
	defer func() { runnerNow = oldNow }()
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		http.NotFound(w, r)
	}))
	defer server.Close()
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if hits != 0 {
			t.Fatalf("cooldown should skip API calls, hits=%d", hits)
		}
	})
}

func TestRunnerCompletionForceIntervalBuildsSameSHA(t *testing.T) {
	state := newRunnerState(t, "same-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	old := time.Date(2026, 9, 15, 1, 0, 0, 0, time.UTC).Format(time.RFC3339)
	writeRunnerTestJSON(t, filepath.Join(state, ".build_state"), buildState{Queued: []map[string]any{}, LastFinishedAt: &old})
	server := fakeGitHub(t, "docs", "same-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		cfg := defaultRunnerConfig(state)
		cfg.ForceBuildIntervalHours = 1
		cfg.BuildCooldownSeconds = 0
		var stdout, stderr bytes.Buffer
		if code := executeRunner(cfg, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); err != nil {
			t.Fatalf("force interval did not build: %v", err)
		}
	})
}

func TestRunnerCompletionSizeWarning(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipelineWithExtra(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "", `dd if=/dev/zero of="$out/big.bin" bs=1048576 count=2 2>/dev/null`)
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		cfg := defaultRunnerConfig(state)
		cfg.OutputSizeWarnMB = 1
		cfg.BuildCooldownSeconds = 0
		var stdout, stderr bytes.Buffer
		if code := executeRunner(cfg, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.Report == nil || !log.Report.SizeWarn {
			t.Fatalf("size warn should be true when threshold exceeded: %+v", log.Report)
		}
		if len(log.Warnings) != 1 || !strings.Contains(log.Warnings[0], "OUTPUT_SIZE_WARN") {
			t.Fatalf("missing size warning: %+v", log.Warnings)
		}
	})
}

func TestRunnerCompletionPATExpiryWarning(t *testing.T) {
	oldNow := runnerNow
	runnerNow = func() time.Time { return time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC) }
	defer func() { runnerNow = oldNow }()

	var out bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&out, &slog.HandlerOptions{Level: slog.LevelDebug}))
	resp := &http.Response{Header: make(http.Header)}
	resp.Header.Set("GitHub-Authentication-Token-Expiration", "2026-09-20T00:00:00Z")
	logPATExpiryWarning(resp, logger)
	if !strings.Contains(out.String(), "PAT_EXPIRY_WARN") || !strings.Contains(out.String(), "remaining_days=4") {
		t.Fatalf("missing PAT expiry warning: %s", out.String())
	}
}

func TestRunnerFixtureR5DeployPending(t *testing.T) {
	state := newRunnerState(t, "old-blob", []DeployTarget{{Host: "example.invalid", User: "deploy", DestDir: "/var/www/html"}})
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	fakeBin := t.TempDir()
	writeExecutable(t, filepath.Join(fakeBin, "ssh"), "#!/bin/sh\nexit 255\n")
	t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"new-blob\"}\n" {
			t.Fatalf("sha not updated before deploy pending")
		}
		var pending []pendingTransfer
		readJSON(t, filepath.Join(state, ".pending_transfers"), &pending)
		if len(pending) != 1 || pending[0].RetryCount != 0 || pending[0].DestDir != "/var/www/html" || pending[0].TargetID == "" || pending[0].Out == nil {
			t.Fatalf("unexpected pending: %+v", pending)
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "success_deploy_pending" || len(log.Deploy) != 1 || log.Deploy[0].Status != "pending" || log.Deploy[0].TransferVerified || log.Error == nil || *log.Error != "deploy pending" {
			t.Fatalf("unexpected deploy log: %+v", log)
		}
	})
}

func TestRunnerFixtureR6LockConflict(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	lock := fmt.Sprintf("pid=%d\nstarted_at=2026-09-16T00:00:00Z\n", os.Getpid())
	if err := os.WriteFile(filepath.Join(state, ".build_lock"), []byte(lock), 0600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d", code)
	}
	if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
		t.Fatalf("history must not be touched")
	}
	if !strings.Contains(stdout.String(), "BUILD_SKIP: already running") {
		t.Fatalf("missing lock log: %s", stdout.String())
	}
}

func TestRunnerFixtureR7CorruptNotifyPending(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	if err := os.WriteFile(filepath.Join(state, ".notify_pending"), []byte("{bad json"), 0600); err != nil {
		t.Fatal(err)
	}
	fixedNow := time.Date(2026, 9, 16, 1, 2, 3, 0, time.UTC)
	oldNow := runnerNow
	runnerNow = func() time.Time { return fixedNow }
	defer func() { runnerNow = oldNow }()
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".notify_pending.corrupt.20260916010203.bak")); err != nil {
			t.Fatalf("missing corrupt backup: %v", err)
		}
		if strings.TrimSpace(readFile(t, filepath.Join(state, ".notify_pending"))) != "[]" {
			t.Fatalf("notify pending not reset")
		}
	})
}

func newRunnerState(t *testing.T, sha string, deploy []DeployTarget) string {
	t.Helper()
	state := t.TempDir()
	buildBin := filepath.Join(t.TempDir(), "adlaire-ci-build")
	writeExecutable(t, buildBin, fakeBuilderScript(0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "", ""))
	t.Setenv("ADLAIRE_CI_BUILD_BIN", buildBin)
	if err := os.WriteFile(filepath.Join(state, ".github_token"), []byte("token\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(state, ".last_sha"), []byte(fmt.Sprintf("{\"sha\":\"%s\"}\n", sha)), 0600); err != nil {
		t.Fatal(err)
	}
	target := BranchTarget{
		Branch: "main", TargetFile: "docs", SHAFile: filepath.Join(state, ".last_sha"),
		Src: filepath.Join(state, "repo", "docs"), Out: filepath.Join(state, "dist", "site"),
		DeployTargets: deploy,
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".branch_config"), branchConfigFile{BranchTargets: []BranchTarget{target}})
	return state
}

func writePipeline(t *testing.T, state string, code int, stdout string, stderr string) {
	t.Helper()
	writePipelineWithExtra(t, state, code, stdout, stderr, "")
}

func writePipelineWithExtra(t *testing.T, state string, code int, stdout string, stderr string, extra string) {
	t.Helper()
	buildBin := filepath.Join(t.TempDir(), "adlaire-ci-build")
	writeExecutable(t, buildBin, fakeBuilderScript(code, stdout, stderr, extra))
	t.Setenv("ADLAIRE_CI_BUILD_BIN", buildBin)
}

func fakeBuilderScript(code int, stdout string, stderr string, extra string) string {
	return fmt.Sprintf(`#!/usr/bin/env bash
set -u
if [ "${1:-}" = "--version" ]; then
  printf 'adlaire-ci-build V.0.0-dev go=fake\n'
  exit 0
fi
out=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --out)
      shift
      out="${1:-}"
      ;;
  esac
  shift || true
done
printf '%%s' %q
printf '%%s' %q >&2
if [ %d -eq 0 ]; then
  mkdir -p "$out/assets"
  printf '<html></html>' > "$out/index.html"
  printf 'body{}' > "$out/assets/style.css"
  printf 'console.log("ok")' > "$out/assets/app.js"
  printf '[]' > "$out/assets/search-index.json"
  %s
fi
exit %d
`, stdout, stderr, code, extra, code)
}

func writeExecutable(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
}

func fakeGitHub(t *testing.T, target, sha, blob string) *httptest.Server {
	t.Helper()
	encoded := base64.StdEncoding.EncodeToString([]byte(blob))
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": target, "type": "blob", "sha": sha}}})
		case strings.Contains(r.URL.Path, "/git/blobs/"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
		default:
			http.NotFound(w, r)
		}
	}))
}

func fakeGitHubWithCommit(t *testing.T, target, sha, blob string) *httptest.Server {
	t.Helper()
	encoded := base64.StdEncoding.EncodeToString([]byte(blob))
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": target, "type": "blob", "sha": sha}}})
		case strings.Contains(r.URL.Path, "/git/blobs/"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
		case strings.Contains(r.URL.Path, "/commits"):
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"sha": "commit-sha",
				"commit": map[string]any{
					"message": "Update docs\n\nBody",
					"author":  map[string]string{"name": "A. Developer", "date": "2026-09-16T00:00:00Z"},
				},
			}})
		default:
			http.NotFound(w, r)
		}
	}))
}

func fakeGitHubWithTags(t *testing.T, target, sha, blob string, refs []map[string]string) *httptest.Server {
	t.Helper()
	encoded := base64.StdEncoding.EncodeToString([]byte(blob))
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": target, "type": "blob", "sha": sha}}})
		case strings.Contains(r.URL.Path, "/git/blobs/"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
		case strings.Contains(r.URL.Path, "/commits"):
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"sha": "commit-sha",
				"commit": map[string]any{
					"message": "Release docs",
					"author":  map[string]string{"name": "A. Developer", "date": "2026-09-16T00:00:00Z"},
				},
			}})
		case strings.Contains(r.URL.Path, "/git/matching-refs/tags"):
			out := []map[string]any{}
			for _, ref := range refs {
				out = append(out, map[string]any{"ref": ref["ref"], "object": map[string]string{"sha": ref["sha"], "type": "commit"}})
			}
			_ = json.NewEncoder(w).Encode(out)
		default:
			http.NotFound(w, r)
		}
	}))
}

func withRunnerServer(t *testing.T, url string, fn func()) {
	t.Helper()
	old := runnerGitHubAPIBase
	runnerGitHubAPIBase = url
	t.Setenv("ADLAIRE_CI_REPOSITORY", "owner/repo")
	defer func() { runnerGitHubAPIBase = old }()
	fn()
}

func onlyBuildLog(t *testing.T, state string) buildLog {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(state, ".build_logs"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one build log, got %d", len(entries))
	}
	var log buildLog
	readJSON(t, filepath.Join(state, ".build_logs", entries[0].Name()), &log)
	return log
}

func writeRunnerTestJSON(t *testing.T, path string, v any) {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
}

func readJSON(t *testing.T, path string, v any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatal(err)
	}
}
