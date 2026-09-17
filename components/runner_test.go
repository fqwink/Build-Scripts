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
	if strings.TrimSpace(stdout.String()) != "Usage: adlaire-ci-runner [--state-dir path] [--once] [--version] [--help]" || stderr.Len() != 0 {
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
		if _, err := os.Stat(filepath.Join(state, ".snapshots", *log.SnapshotID, "site", "index.html")); err != nil {
			t.Fatalf("snapshot missing: %v", err)
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".build_history")), `"status":"success"`) {
			t.Fatalf("history missing success")
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
	})
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
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 1 {
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
	writeRunnerTestJSON(t, filepath.Join(state, ".pending_transfers"), []pendingTransfer{{Out: out, Host: "host", User: "deploy", DestDir: "/var/www/html", RetryCount: 1}})
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
	dir := filepath.Join(state, "repo", ".ci")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	script := "#!/usr/bin/env bash\nprintf '%s' '[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default'\nmkdir -p \"$ADLAIRE_CI_OUT\"\ndd if=/dev/zero of=\"$ADLAIRE_CI_OUT/big.bin\" bs=1048576 count=2 2>/dev/null\nexit 0\n"
	writeExecutable(t, filepath.Join(dir, "pipeline.sh"), script)
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
		if len(pending) != 1 || pending[0].RetryCount != 1 || pending[0].DestDir != "/var/www/html" {
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
	writeExecutable(t, buildBin, "#!/bin/sh\nprintf 'adlaire-ci-build ADLAIRE_CI_SPEC go=fake\\n'\n")
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
	dir := filepath.Join(state, "repo", ".ci")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	script := fmt.Sprintf("#!/usr/bin/env bash\nprintf '%%s' %q\nprintf '%%s' %q >&2\nmkdir -p \"$ADLAIRE_CI_OUT\"\nprintf '<html></html>' > \"$ADLAIRE_CI_OUT/index.html\"\nexit %d\n", stdout, stderr, code)
	writeExecutable(t, filepath.Join(dir, "pipeline.sh"), script)
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
