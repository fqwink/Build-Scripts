package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	if code := runRunner([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("help exit=%d", code)
	}
	if strings.TrimSpace(stdout.String()) != "Usage: adlaire-ci-runner [--state-dir path] [--once] [--version] [--help]" || stderr.Len() != 0 {
		t.Fatalf("unexpected help stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := runRunner([]string{"--state-dir", "relative"}, &stdout, &stderr); code != 2 {
		t.Fatalf("relative exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "state directory must be absolute: relative") {
		t.Fatalf("missing relative error: %s", stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := runRunner([]string{"--unknown"}, &stdout, &stderr); code != 2 {
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
		code := runRunner([]string{"--state-dir", state}, &stdout, &stderr)
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
		code := runRunner([]string{"--state-dir", state}, &stdout, &stderr)
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
		code := runRunner([]string{"--state-dir", state}, &stdout, &stderr)
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

func TestRunnerFixtureR5DeployPending(t *testing.T) {
	state := newRunnerState(t, "old-blob", []DeployTarget{{Host: "example.invalid", User: "deploy", DestDir: "/var/www/html"}})
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	fakeBin := t.TempDir()
	writeExecutable(t, filepath.Join(fakeBin, "ssh"), "#!/bin/sh\nexit 255\n")
	t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := runRunner([]string{"--state-dir", state}, &stdout, &stderr)
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
	if code := runRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
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
		if code := runRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
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
	writeJSON(t, filepath.Join(state, ".branch_config"), branchConfigFile{BranchTargets: []BranchTarget{target}})
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

func writeJSON(t *testing.T, path string, v any) {
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
