package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRunAPICLI(t *testing.T) {
	t.Run("help", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunAPI([]string{"--help", "--state-dir", "relative"}, strings.NewReader(""), &stdout, &stderr)
		if code != 0 {
			t.Fatalf("expected help exit 0, got %d", code)
		}
		if got := strings.TrimSpace(stdout.String()); got != "Usage: adlaire-ci-api --state-dir path [--addr 127.0.0.1:port] [--init-credentials] [--version] [--help]" {
			t.Fatalf("unexpected help: %q", got)
		}
		if stderr.Len() != 0 {
			t.Fatalf("expected empty stderr, got %q", stderr.String())
		}
	})

	t.Run("version", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunAPI([]string{"--version", "--state-dir", "relative"}, strings.NewReader(""), &stdout, &stderr)
		if code != 0 {
			t.Fatalf("expected version exit 0, got %d", code)
		}
		if got := strings.TrimSpace(stdout.String()); !strings.HasPrefix(got, "adlaire-ci-api V.0.0-dev go=") {
			t.Fatalf("unexpected version: %q", got)
		}
		if stderr.Len() != 0 {
			t.Fatalf("expected empty stderr, got %q", stderr.String())
		}
	})

	t.Run("init credentials", func(t *testing.T) {
		state := t.TempDir()
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunAPI([]string{"--state-dir", state, "--init-credentials"}, strings.NewReader("password123\n"), &stdout, &stderr)
		if code != 0 {
			t.Fatalf("expected init exit 0, got %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
		if strings.TrimSpace(stdout.String()) != "credentials initialized" || stderr.Len() != 0 {
			t.Fatalf("unexpected init output stdout=%q stderr=%q", stdout.String(), stderr.String())
		}
		if _, err := readCredentials(filepath.Join(state, ".admin_credentials")); err != nil {
			t.Fatalf("credentials should be readable: %v", err)
		}

		stdout.Reset()
		stderr.Reset()
		code = RunAPI([]string{"--state-dir", state, "--init-credentials"}, strings.NewReader("another-password\n"), &stdout, &stderr)
		if code != 2 || stdout.Len() != 0 || strings.TrimSpace(stderr.String()) != "credentials already exist" {
			t.Fatalf("unexpected existing credentials result code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
	})

	t.Run("validation", func(t *testing.T) {
		state := t.TempDir()
		cases := []struct {
			name string
			args []string
			err  string
		}{
			{name: "missing state", args: []string{}, err: "state directory is required"},
			{name: "empty state", args: []string{"--state-dir", ""}, err: "state directory must not be empty"},
			{name: "relative state", args: []string{"--state-dir", "relative"}, err: "state directory must be absolute: relative"},
			{name: "addr with init", args: []string{"--state-dir", state, "--init-credentials", "--addr", "127.0.0.1:8765"}, err: "--addr is not allowed with --init-credentials"},
			{name: "bad addr", args: []string{"--state-dir", state, "--addr", "0.0.0.0:8765"}, err: "invalid listen address: 0.0.0.0:8765"},
			{name: "credentials missing", args: []string{"--state-dir", state}, err: "credentials are not initialized"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				code := RunAPI(tc.args, strings.NewReader("password123\n"), &stdout, &stderr)
				if code != 2 || stdout.Len() != 0 || strings.TrimSpace(stderr.String()) != tc.err {
					t.Fatalf("unexpected result code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
				}
			})
		}
	})

	t.Run("invalid initial password", func(t *testing.T) {
		state := t.TempDir()
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunAPI([]string{"--state-dir", state, "--init-credentials"}, strings.NewReader("short\n"), &stdout, &stderr)
		if code != 2 || stdout.Len() != 0 || strings.TrimSpace(stderr.String()) != "invalid initial password" {
			t.Fatalf("unexpected invalid password result code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".admin_credentials")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("credentials must not be created, stat err=%v", err)
		}
	})

	t.Run("invalid init credentials direct input", func(t *testing.T) {
		state := t.TempDir()
		now := time.Date(2026, 9, 17, 1, 0, 0, 0, time.UTC)
		if err := InitCredentials(state, "short", now); err == nil || err.Error() != "invalid initial password" {
			t.Fatalf("expected invalid direct password, got %v", err)
		}
		if _, err := os.Stat(filepath.Join(state, ".admin_credentials")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("credentials must not be created, stat err=%v", err)
		}
	})
}

func TestAPIAdminStaticServingContract(t *testing.T) {
	state := newAPIState(t)
	adminDir := filepath.Join(state, "admin")
	if err := os.Mkdir(adminDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "index.html"), []byte("<main>admin</main>\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "adlaire-ci-sdk.js"), []byte("export const ready = true;\n"), 0644); err != nil {
		t.Fatal(err)
	}
	server, err := NewAPIServer(APIConfig{StateDir: state})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		method       string
		path         string
		status       int
		contentType  string
		cacheControl string
		body         string
	}{
		{http.MethodGet, "/", http.StatusOK, "text/html; charset=utf-8", "no-store", "<main>admin</main>\n"},
		{http.MethodGet, "/admin/", http.StatusOK, "text/html; charset=utf-8", "no-store", "<main>admin</main>\n"},
		{http.MethodGet, "/admin/index.html", http.StatusOK, "text/html; charset=utf-8", "no-store", "<main>admin</main>\n"},
		{http.MethodHead, "/admin/index.html", http.StatusOK, "text/html; charset=utf-8", "no-store", ""},
		{http.MethodGet, "/admin/adlaire-ci-sdk.js", http.StatusOK, "text/javascript; charset=utf-8", "no-cache", "export const ready = true;\n"},
	}
	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		resp := httptest.NewRecorder()
		server.Handler().ServeHTTP(resp, req)
		if resp.Code != tc.status || resp.Header().Get("Content-Type") != tc.contentType || resp.Header().Get("Cache-Control") != tc.cacheControl || resp.Body.String() != tc.body {
			t.Fatalf("%s %s code=%d type=%q cache=%q body=%q", tc.method, tc.path, resp.Code, resp.Header().Get("Content-Type"), resp.Header().Get("Cache-Control"), resp.Body.String())
		}
	}

	for _, path := range []string{"/admin/missing", "/admin/../.admin_credentials", "/.admin_credentials"} {
		resp := httptest.NewRecorder()
		server.Handler().ServeHTTP(resp, httptest.NewRequest(http.MethodGet, path, nil))
		if resp.Code != http.StatusNotFound {
			t.Fatalf("unsafe path %s returned %d", path, resp.Code)
		}
	}
	resp := httptest.NewRecorder()
	server.Handler().ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/admin/index.html", strings.NewReader("ignored")))
	if resp.Code != http.StatusMethodNotAllowed || resp.Header().Get("Allow") != "GET, HEAD" {
		t.Fatalf("method contract code=%d allow=%q", resp.Code, resp.Header().Get("Allow"))
	}

	if err := os.Remove(filepath.Join(adminDir, "adlaire-ci-sdk.js")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(state, ".admin_credentials"), filepath.Join(adminDir, "adlaire-ci-sdk.js")); err != nil {
		t.Fatal(err)
	}
	resp = httptest.NewRecorder()
	server.Handler().ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/admin/adlaire-ci-sdk.js", nil))
	if resp.Code != http.StatusNotFound {
		t.Fatalf("symlink asset returned %d", resp.Code)
	}
}

func TestAPILoginStatusAndQueue(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")

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

	activeEntry := map[string]any{"id": "q20260917010000", "trigger": "manual", "priority": "normal", "created_seq": 1, "payload": map[string]any{"force": false}}
	writeTestJSON(t, filepath.Join(state, ".build_state"), apiBuildState{
		Running:          true,
		ActiveQueueEntry: activeEntry,
		Queued:           []map[string]any{{"id": "q20260917010102", "trigger": "manual", "priority": "normal", "created_seq": 2, "payload": map[string]any{"force": true}}},
	})
	resp = apiRequest(t, server, http.MethodDelete, "/api/queue", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("active-preserving clear code=%d body=%s", resp.Code, resp.Body.String())
	}
	var stateAfterClear apiBuildState
	readTestJSON(t, filepath.Join(state, ".build_state"), &stateAfterClear)
	if len(stateAfterClear.Queued) != 0 || stateAfterClear.ActiveQueueEntry["id"] != activeEntry["id"] {
		t.Fatalf("queue clear must preserve active entry: %#v", stateAfterClear)
	}

	writeTestJSON(t, filepath.Join(state, ".build_state"), apiBuildState{
		Running: true,
		Queued:  []map[string]any{{"id": "q20260917010101", "trigger": "webhook", "payload": map[string]any{"branch": "main", "sha": strings.Repeat("1", 40)}}},
	})
	resp = apiRequest(t, server, http.MethodPost, "/api/build/force", token, nil)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("build force collision code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &build)
	if build["queue_id"] != "q20260917010101-001" {
		t.Fatalf("unexpected collision queue id: %#v", build)
	}
}

func TestAPILogsHistoryAndCircuitReset(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")

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
	var message map[string]string
	decodeTestJSON(t, resp.Body.Bytes(), &message)
	if message["message"] != "Comment saved" {
		t.Fatalf("unexpected comment response: %#v", message)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/history/b20260917010101/flag", token, map[string]bool{"flagged": true})
	if resp.Code != http.StatusOK {
		t.Fatalf("flag code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &message)
	if message["message"] != "Flag updated" {
		t.Fatalf("unexpected flag response: %#v", message)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/history/b20260917010101/tags", token, map[string][]string{"tags": []string{"release", "manual"}})
	if resp.Code != http.StatusOK {
		t.Fatalf("tags code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &message)
	if message["message"] != "Tags updated" {
		t.Fatalf("unexpected tags response: %#v", message)
	}
	var updated apiBuildLog
	readTestJSON(t, filepath.Join(state, ".build_logs", "b20260917010101.json"), &updated)
	if updated.Comment == nil || *updated.Comment != "reviewed" || !updated.Flagged || len(updated.Tags) != 2 {
		t.Fatalf("history log was not updated: %#v", updated)
	}
	beforeHistoryMeta := snapshotFiles(t, state, []string{".build_logs/b20260917010101.json", ".build_history", ".config_log", ".audit_log"})
	resp = apiRequest(t, server, http.MethodPost, "/api/history/b20260917010101/comment", token, map[string]string{"comment": " reviewed "})
	if resp.Code != http.StatusOK {
		t.Fatalf("comment no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/history/b20260917010101/flag", token, map[string]bool{"flagged": true})
	if resp.Code != http.StatusOK {
		t.Fatalf("flag no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/history/b20260917010101/tags", token, map[string][]string{"tags": []string{"manual", "release", "manual"}})
	if resp.Code != http.StatusOK {
		t.Fatalf("tags no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, beforeHistoryMeta, snapshotFiles(t, state, []string{".build_logs/b20260917010101.json", ".build_history", ".config_log", ".audit_log"}))
	resp = apiRequest(t, server, http.MethodGet, "/api/history?trigger=manual&tag=release&flagged=true", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("filtered history code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &history)
	if history["total"].(float64) != 1 {
		t.Fatalf("unexpected filtered history: %#v", history)
	}

	writeTestJSON(t, filepath.Join(state, ".build_circuit_state"), apiCircuitState{Open: true, ConsecutiveFailures: 3})
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
	token := login(t, server, "password123")
	resp = apiRequest(t, server, http.MethodPost, "/api/status", token, nil)
	if resp.Code != http.StatusMethodNotAllowed {
		t.Fatalf("method code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/history?page=0", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("validation code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/stats?days=7&days=8", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("duplicate query code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/history/bad.id/log", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid history id code=%d body=%s", resp.Code, resp.Body.String())
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
	token := login(t, server, "password123")
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
	if status["last_sha"] != "blob-1" || status["last_build_status"] != "success" || status["last_target_status"] != "success" || status["running"] != false {
		t.Fatalf("unexpected fallback status: %#v", status)
	}
	if _, exists := status["queued"]; exists {
		t.Fatalf("status must not expose queue state: %#v", status)
	}
	if _, exists := status["output_url"]; exists {
		t.Fatalf("status must not expose output url: %#v", status)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_state", ".build_history", ".build_status.json"}))

	lastBlob := "blob-primary"
	normalizedStatus := "success"
	lastStatus := "success"
	lastTrigger := "manual"
	writeTestJSON(t, filepath.Join(state, ".build_status.json"), apiBuildStatus{
		Status: &normalizedStatus, LastBlobSHA: &lastBlob, LastTargetStatus: &lastStatus, LastTrigger: &lastTrigger,
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
	if status["last_sha"] != "blob-primary" || status["last_build_status"] != "success" || status["last_target_status"] != "success" || status["pending_transfers_count"].(float64) != 2 || status["running"] != true {
		t.Fatalf("unexpected primary status: %#v", status)
	}
	if _, exists := status["queued"]; exists {
		t.Fatalf("primary status must not expose queue state: %#v", status)
	}
	if _, exists := status["output_url"]; exists {
		t.Fatalf("primary status must not expose output url: %#v", status)
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
	token := login(t, server, "password123")

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
	token := login(t, server, "password123")

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
	if len(readJSONLines(filepath.Join(state, ".config_log"))) != 1 || len(readJSONLines(filepath.Join(state, ".audit_log"))) != 1 {
		t.Fatalf("circuit reset must append config and audit logs")
	}
	before = snapshotFiles(t, state, []string{".build_circuit_state", ".config_log", ".audit_log"})
	resp = apiRequest(t, server, http.MethodPost, "/api/circuit-breaker/reset", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("circuit reset no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_circuit_state", ".config_log", ".audit_log"}))
}

func TestAPIPhase3RequestIDAndAccessLogContract(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")
	if err := os.Remove(filepath.Join(state, ".api_access_log")); err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/logs?n=2&q=raw-token", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Phase3Test/1")
	req.RemoteAddr = "192.0.2.1:1234"
	authenticated := httptest.NewRecorder()
	server.Handler().ServeHTTP(authenticated, req)
	if authenticated.Code != http.StatusOK {
		t.Fatalf("authenticated code=%d body=%s", authenticated.Code, authenticated.Body.String())
	}
	requestID := authenticated.Header().Get("X-Request-Id")
	if len(requestID) != 32 || strings.ToLower(requestID) != requestID {
		t.Fatalf("unexpected request id: %q", requestID)
	}
	if authenticated.Header().Get("Cache-Control") != "no-store" || authenticated.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("missing security headers: %#v", authenticated.Header())
	}

	records := readJSONLines(filepath.Join(state, ".api_access_log"))
	if len(records) != 1 {
		t.Fatalf("unexpected access log records: %#v", records)
	}
	record := records[0]
	if record["request_id"] != requestID || record["method"] != http.MethodGet || record["path"] != "/api/logs" || record["status"].(float64) != 200 {
		t.Fatalf("unexpected access log record: %#v", record)
	}
	if record["auth_type"] != "session" || record["actor"] != "admin" || record["remote_addr"] != "192.0.2.1" || record["user_agent"] != "Phase3Test/1" {
		t.Fatalf("unexpected access log actor/source: %#v", record)
	}
	if record["duration_ms"].(float64) < 0 {
		t.Fatalf("duration must be non-negative: %#v", record)
	}
	query := record["query"].(map[string]any)
	if query["n"] != "2" {
		t.Fatalf("validated query value missing: %#v", query)
	}
	if query["q"] != "***" {
		t.Fatalf("free-form query value must be redacted: %#v", query)
	}
	if _, ok := query["secret"]; ok {
		t.Fatalf("secret query must not be logged: %#v", query)
	}
	serialized, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(serialized, []byte(token)) || bytes.Contains(serialized, []byte("raw-token")) || bytes.Contains(serialized, []byte("Authorization")) {
		t.Fatalf("access log leaked secret material: %s", serialized)
	}

	unauthorized := apiRequest(t, server, http.MethodGet, "/api/status", "", nil)
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized code=%d body=%s", unauthorized.Code, unauthorized.Body.String())
	}
	unauthorizedID := unauthorized.Header().Get("X-Request-Id")
	if len(unauthorizedID) != 32 {
		t.Fatalf("unauthorized request id missing: %q", unauthorizedID)
	}
	records = readJSONLines(filepath.Join(state, ".api_access_log"))
	last := records[len(records)-1]
	if last["request_id"] != unauthorizedID || last["status"].(float64) != 401 || last["auth_type"] != "none" || last["actor"] != nil {
		t.Fatalf("unexpected unauthorized access log: %#v", last)
	}
}

func TestAPIPhase3RequestIDFailureStopsRequest(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	server.cfg.RequestID = func() (string, error) {
		return "", errors.New("entropy unavailable")
	}

	resp := apiRequest(t, server, http.MethodGet, "/api/status", "", nil)
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("request id failure code=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp.Header().Get("X-Request-Id") != "" {
		t.Fatalf("request id header must not be set: %#v", resp.Header())
	}
	if records := readJSONLines(filepath.Join(state, ".api_access_log")); len(records) != 0 {
		t.Fatalf("request id failure must not append api access log: %#v", records)
	}
}

func TestAPIPhase3ServerAndStatusWriterContract(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	httpServer := server.HTTPServer()
	if httpServer.Addr != defaultAPIAddr || httpServer.Handler == nil {
		t.Fatalf("unexpected http server: %#v", httpServer)
	}
	if httpServer.ReadHeaderTimeout != 5*time.Second || httpServer.ReadTimeout != 30*time.Second || httpServer.WriteTimeout != 0 || httpServer.IdleTimeout != 60*time.Second || httpServer.MaxHeaderBytes != 32768 {
		t.Fatalf("unexpected timeout config: %#v", httpServer)
	}
	if httpServer.ErrorLog == nil {
		t.Fatal("error log must be redacted logger")
	}

	base := &plainResponseWriter{header: http.Header{}}
	rw, rec := statusResponseWriter(base)
	rw.WriteHeader(http.StatusAccepted)
	rw.WriteHeader(http.StatusInternalServerError)
	if rec.status != http.StatusAccepted || base.status != http.StatusAccepted {
		t.Fatalf("first explicit status must win: rec=%d base=%d", rec.status, base.status)
	}

	base = &plainResponseWriter{header: http.Header{}}
	rw, rec = statusResponseWriter(base)
	n, err := rw.Write([]byte("ok"))
	if err != nil || n != 2 || rec.status != http.StatusOK || base.body.String() != "ok" {
		t.Fatalf("implicit write mismatch: n=%d err=%v rec=%d body=%q", n, err, rec.status, base.body.String())
	}
	if _, ok := rw.(http.Flusher); ok {
		t.Fatal("non-flusher writer must not be reported as flusher")
	}

	flushing := &flushingResponseWriter{plainResponseWriter: plainResponseWriter{header: http.Header{}}}
	rw, _ = statusResponseWriter(flushing)
	flusher, ok := rw.(http.Flusher)
	if !ok {
		t.Fatal("flusher writer must pass through flusher")
	}
	flusher.Flush()
	if flushing.flushes != 1 {
		t.Fatalf("flush count=%d", flushing.flushes)
	}
}

func TestAPIPhase3AccessControlCorruptFailsClosed(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")
	if err := os.WriteFile(filepath.Join(state, ".access_control"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}

	resp := apiRequest(t, server, http.MethodGet, "/api/unknown", "", nil)
	if resp.Code != http.StatusNotFound {
		t.Fatalf("unknown path must precede access control: code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/status", "", nil)
	if resp.Code != http.StatusMethodNotAllowed {
		t.Fatalf("method mismatch must precede access control/auth: code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/config", token, nil)
	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("corrupt access control code=%d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]string
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["error"] != "Access control unavailable" {
		t.Fatalf("unexpected body: %#v", body)
	}

	resp = apiRequestFrom(t, server, http.MethodGet, "/api/health", "", nil, "203.0.113.9:1234")
	if resp.Code != http.StatusOK {
		t.Fatalf("health must bypass corrupt access control: code=%d body=%s", resp.Code, resp.Body.String())
	}
	records := readJSONLines(filepath.Join(state, ".api_access_log"))
	found := false
	for _, record := range records {
		if record["path"] == "/api/config" && record["status"].(float64) == 503 {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing 503 access log record: %#v", records)
	}
}

func TestAPIPhase3BodyValidationContract(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")

	oversized := `{"queue_max_size":` + strings.Repeat("1", 1<<20) + `}`
	req := httptest.NewRequest(http.MethodPost, "/api/config/validate", strings.NewReader(oversized))
	req.Header.Set("Authorization", "Bearer "+token)
	req.RemoteAddr = "192.0.2.1:1234"
	resp := httptest.NewRecorder()
	server.Handler().ServeHTTP(resp, req)
	if resp.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized JSON code=%d body=%s", resp.Code, resp.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/api/config/validate", strings.NewReader(`{} {}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.RemoteAddr = "192.0.2.1:1234"
	resp = httptest.NewRecorder()
	server.Handler().ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("multiple JSON values code=%d body=%s", resp.Code, resp.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/status", strings.NewReader(strings.Repeat("x", 1<<20+1)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.RemoteAddr = "192.0.2.1:1234"
	resp = httptest.NewRecorder()
	server.Handler().ServeHTTP(resp, req)
	if resp.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized forbidden body code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPIMaintenanceEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")

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
	var body map[string]any
	if len(readJSONLines(filepath.Join(state, ".config_log"))) != 1 || len(readJSONLines(filepath.Join(state, ".audit_log"))) != 1 {
		t.Fatalf("maintenance enable must append config and audit logs")
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/build", token, nil)
	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("build during maintenance code=%d body=%s", resp.Code, resp.Body.String())
	}

	before := snapshotFiles(t, state, []string{".maintenance", ".config_log", ".audit_log"})
	resp = apiRequest(t, server, http.MethodPost, "/api/maintenance/enable", token, map[string]string{"reason": "deploy window"})
	if resp.Code != http.StatusOK {
		t.Fatalf("enable no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["message"] != "No changes" {
		t.Fatalf("unexpected enable no-op: %#v", body)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".maintenance", ".config_log", ".audit_log"}))

	resp = apiRequest(t, server, http.MethodPost, "/api/maintenance/disable", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("disable code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".maintenance"), &maintenance)
	if maintenance.Enabled || maintenance.Reason != nil || maintenance.Since != nil {
		t.Fatalf("maintenance was not disabled: %#v", maintenance)
	}
	if len(readJSONLines(filepath.Join(state, ".config_log"))) != 2 || len(readJSONLines(filepath.Join(state, ".audit_log"))) != 2 {
		t.Fatalf("maintenance disable must append config and audit logs")
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/build", token, nil)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("build after maintenance code=%d body=%s", resp.Code, resp.Body.String())
	}

	if err := os.WriteFile(filepath.Join(state, ".maintenance"), []byte(`{"enabled":`), 0600); err != nil {
		t.Fatal(err)
	}
	before = snapshotFiles(t, state, []string{".maintenance", ".build_state", ".audit_log"})
	resp = apiRequest(t, server, http.MethodPost, "/api/build", token, nil)
	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("build with corrupt maintenance code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["error"] != "maintenance_unavailable" {
		t.Fatalf("unexpected maintenance error: %#v", body)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".maintenance", ".build_state", ".audit_log"}))
	if err := os.Remove(filepath.Join(state, ".maintenance")); err != nil {
		t.Fatal(err)
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
	token := login(t, server, "password123")

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
	var body map[string]any

	resp = apiRequest(t, server, http.MethodPost, "/api/repo-config", token, map[string]string{"branch": "release"})
	if resp.Code != http.StatusOK {
		t.Fatalf("repo-config no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
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
	if err := os.WriteFile(filepath.Join(state, ".server_config"), []byte("{\"watch_mode\":\"local\"}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	server := newTestAPI(t, state)
	token := login(t, server, "password123")
	systemdCalls := []string{}
	server.cfg.SystemdDir = filepath.Join(state, "systemd")
	server.cfg.CommandRunner = func(ctx context.Context, name string, args ...string) (apiCommandResult, error) {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatalf("schedule systemd command must set a deadline")
		}
		call := name + " " + strings.Join(args, " ")
		systemdCalls = append(systemdCalls, call)
		switch call {
		case "systemctl daemon-reload", "systemctl restart adlaire-ci.timer":
			return apiCommandResult{ExitCode: 0}, nil
		case "systemctl show adlaire-ci.timer -p OnUnitActiveSec":
			return apiCommandResult{Stdout: "OnUnitActiveSec=600s\n", ExitCode: 0}, nil
		default:
			t.Fatalf("unexpected schedule systemd command: %s", call)
		}
		return apiCommandResult{ExitCode: -1}, errors.New("unexpected command")
	}

	resp := apiRequest(t, server, http.MethodGet, "/api/schedule", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("schedule code=%d body=%s", resp.Code, resp.Body.String())
	}
	var schedule map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &schedule)
	if schedule["interval"] != "5min" || schedule["interval_seconds"].(float64) != 300 || schedule["paused"] != false {
		t.Fatalf("unexpected schedule defaults: %#v", schedule)
	}

	var body map[string]any
	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/interval", token, map[string]int{"interval_seconds": 600})
	if resp.Code != http.StatusOK {
		t.Fatalf("interval code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["message"] != "Schedule interval updated" || body["interval_seconds"].(float64) != 600 {
		t.Fatalf("unexpected interval response: %#v", body)
	}
	var cfg apiServerConfig
	readTestJSON(t, filepath.Join(state, ".server_config"), &cfg)
	if cfg.ScheduleIntervalSeconds != 600 {
		t.Fatalf("interval was not saved: %#v", cfg)
	}
	dropIn, err := os.ReadFile(filepath.Join(server.cfg.SystemdDir, "adlaire-ci.timer.d", "override.conf"))
	if err != nil || string(dropIn) != "[Timer]\nOnUnitActiveSec=600s\nPersistent=true\n" {
		t.Fatalf("unexpected schedule drop-in: %q err=%v", string(dropIn), err)
	}
	expectedSystemdCalls := []string{"systemctl daemon-reload", "systemctl restart adlaire-ci.timer", "systemctl show adlaire-ci.timer -p OnUnitActiveSec"}
	if len(systemdCalls) != len(expectedSystemdCalls) {
		t.Fatalf("unexpected schedule systemd calls: %#v", systemdCalls)
	}
	for i, expected := range expectedSystemdCalls {
		if systemdCalls[i] != expected {
			t.Fatalf("unexpected schedule systemd calls: %#v", systemdCalls)
		}
	}
	before := snapshotFiles(t, state, []string{".server_config", ".config_log", ".audit_log", "systemd/adlaire-ci.timer.d/override.conf"})
	resp = apiRequest(t, server, http.MethodPost, "/api/schedule/interval", token, map[string]int{"interval_seconds": 600})
	if resp.Code != http.StatusOK {
		t.Fatalf("interval no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["message"] != "No changes" || body["interval_seconds"].(float64) != 600 {
		t.Fatalf("unexpected interval no-op response: %#v", body)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".server_config", ".config_log", ".audit_log", "systemd/adlaire-ci.timer.d/override.conf"}))
	if len(systemdCalls) != 3 {
		t.Fatalf("no-op interval must not call systemd: %#v", systemdCalls)
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
	if cfg.SchedulePaused || cfg.ForceBuildIntervalHours != 24 || cfg.BuildCooldownSeconds != 120 || cfg.AllowedHours != nil || cfg.WatchMode != "local" {
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
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["total"].(float64) != 7 {
		t.Fatalf("unexpected schedule config log: %#v", body)
	}
}

func TestAPIAccessControlEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")

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
	token := login(t, server, "password123")

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
	token := login(t, server, "password123")

	if err := os.MkdirAll(filepath.Join(state, ".build_logs"), 0755); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(state, ".build_logs", "b20260917010101.json"), apiBuildLog{
		ID: "b20260917010101", StartedAt: "2026-09-17T01:01:01Z", FinishedAt: "2026-09-17T01:01:03Z",
		TargetStatus: "success", Pipeline: apiPipelineLog{Stdout: "[INFO] streamed", Stderr: ""},
	})

	resp := apiRequest(t, server, http.MethodGet, "/api/build/stream", token, nil)
	if resp.Code != http.StatusOK || resp.Header().Get("Content-Type") != "text/event-stream; charset=utf-8" {
		t.Fatalf("stream code=%d content-type=%q body=%s", resp.Code, resp.Header().Get("Content-Type"), resp.Body.String())
	}
	streamBody := resp.Body.String()
	if !strings.Contains(streamBody, `data: {"type":"log","line":"[INFO] streamed","at":"2026-09-17T01:01:03Z"}`) {
		t.Fatalf("stream missing log frame: %s", streamBody)
	}
	if !strings.Contains(streamBody, `data: {"type":"end","status":"success","duration_seconds":0}`) {
		t.Fatalf("stream missing end frame: %s", streamBody)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/change-password", token, map[string]string{
		"current_password": "password123",
		"new_password":     "changed-password",
	})
	if resp.Code != http.StatusOK {
		t.Fatalf("change-password code=%d body=%s", resp.Code, resp.Body.String())
	}
	passwordAccess := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{"password_change": true})
	if len(passwordAccess) != 1 || passwordAccess[0]["result"] != "success" || passwordAccess[0]["reason"] != nil {
		t.Fatalf("unexpected password change access log: %#v", passwordAccess)
	}
	passwordAudit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"password_change": true})
	if len(passwordAudit) != 1 || passwordAudit[0]["actor_type"] != "admin" || passwordAudit[0]["actor_id"] != "admin" || passwordAudit[0]["target_type"] != "auth" || passwordAudit[0]["target_id"] != "password" || passwordAudit[0]["result"] != "success" {
		t.Fatalf("unexpected password change audit log: %#v", passwordAudit)
	}
	_ = login(t, server, "changed-password")

	resp = apiRequest(t, server, http.MethodPost, "/api/logout", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("logout code=%d body=%s", resp.Code, resp.Body.String())
	}
	logoutAccess := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{"logout": true})
	if len(logoutAccess) != 1 || logoutAccess[0]["result"] != "success" || logoutAccess[0]["reason"] != nil {
		t.Fatalf("unexpected logout access log: %#v", logoutAccess)
	}
	logoutAudit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"logout": true})
	if len(logoutAudit) != 1 || logoutAudit[0]["actor_type"] != "admin" || logoutAudit[0]["actor_id"] != "admin" || logoutAudit[0]["target_type"] != "auth" || logoutAudit[0]["target_id"] != "logout" || logoutAudit[0]["result"] != "success" {
		t.Fatalf("unexpected logout audit log: %#v", logoutAudit)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("post-logout code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPIOperationConfigSessionsAndLogs(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")
	second := login(t, server, "password123")

	resp := apiRequest(t, server, http.MethodGet, "/api/health", "", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("health code=%d body=%s", resp.Code, resp.Body.String())
	}
	var health map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &health)
	if health["status"] != "degraded" || health["last_build_status"] != "none" || health["pending_transfers"].(float64) != 0 {
		t.Fatalf("unexpected initial health: %#v", health)
	}
	checks, ok := health["checks"].([]any)
	if !ok || len(checks) != 1 || checks[0] != "build_status_missing" {
		t.Fatalf("unexpected initial health checks: %#v", health["checks"])
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
	var validation map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &validation)
	if validation["valid"] != true || len(validation["errors"].([]any)) != 0 {
		t.Fatalf("unexpected valid response: %#v", validation)
	}
	if _, err := os.Stat(filepath.Join(state, ".server_config")); !os.IsNotExist(err) {
		t.Fatalf("validate must not write server config")
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/config/validate", token, map[string]any{
		"build_retry_max":          11,
		"commit_status_context":    "",
		"commit_status_target_url": "http://example.test/status",
	})
	if resp.Code != http.StatusOK {
		t.Fatalf("validate invalid code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &validation)
	if validation["valid"] != false {
		t.Fatalf("validate should be invalid: %#v", validation)
	}
	errorsList := validation["errors"].([]any)
	if len(errorsList) != 2 || errorsList[0].(map[string]any)["field"] != "build_retry_max" || errorsList[1].(map[string]any)["field"] != "commit_status_context" {
		t.Fatalf("validation errors must follow request order: %#v", errorsList)
	}
	warningsList := validation["warnings"].([]any)
	if len(warningsList) != 1 || warningsList[0].(map[string]any)["field"] != "commit_status_target_url" || warningsList[0].(map[string]any)["code"] != "http_url" {
		t.Fatalf("unexpected validation warnings: %#v", warningsList)
	}
	if _, err := os.Stat(filepath.Join(state, ".server_config")); !os.IsNotExist(err) {
		t.Fatalf("invalid validate must not write server config")
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/config/validate", token, map[string]any{"unknown_key": true})
	if resp.Code != http.StatusUnprocessableEntity || !strings.Contains(resp.Body.String(), "Unknown config key") {
		t.Fatalf("validate unknown key code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/config/validate", token, map[string]any{"webhook_secret": "dont-log-this"})
	if resp.Code != http.StatusUnprocessableEntity || strings.Contains(resp.Body.String(), "dont-log-this") {
		t.Fatalf("validate secret key code=%d body=%s", resp.Code, resp.Body.String())
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
	revokeAccess := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{"session_revoke_all": true})
	if len(revokeAccess) != 1 || revokeAccess[0]["result"] != "success" || revokeAccess[0]["reason"] != nil {
		t.Fatalf("unexpected session revoke access log: %#v", revokeAccess)
	}
	revokeAudit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"session_revoke_all": true})
	if len(revokeAudit) != 1 || revokeAudit[0]["actor_type"] != "admin" || revokeAudit[0]["actor_id"] != "admin" || revokeAudit[0]["target_type"] != "auth" || revokeAudit[0]["target_id"] != "sessions" || revokeAudit[0]["result"] != "success" {
		t.Fatalf("unexpected session revoke audit log: %#v", revokeAudit)
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

func TestAPIHealthContract(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)

	status := "success"
	targetStatus := "success_deploy_pending"
	finishedAt := "2026-09-17T01:00:00Z"
	deployAt := "2026-09-17T01:00:30Z"
	deployStatus := "pending"
	writeTestJSON(t, filepath.Join(state, ".build_status.json"), apiBuildStatus{
		Status: &status, LastTargetStatus: &targetStatus, LastFinishedAt: &finishedAt,
		LastDeployAt: &deployAt, LastDeployStatus: &deployStatus,
	})
	before := snapshotFiles(t, state, []string{".build_status.json", ".build_history", ".pending_transfers", ".notify_pending"})
	resp := apiRequest(t, server, http.MethodGet, "/api/health", "", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("health ok code=%d body=%s", resp.Code, resp.Body.String())
	}
	var health map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &health)
	if health["status"] != "ok" || health["last_build_at"] != finishedAt || health["last_build_status"] != "success" || health["last_deploy_at"] != deployAt || health["last_deploy_status"] != deployStatus {
		t.Fatalf("unexpected healthy response: %#v", health)
	}
	if checks := health["checks"].([]any); len(checks) != 0 {
		t.Fatalf("healthy checks must be empty: %#v", checks)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_status.json", ".build_history", ".pending_transfers", ".notify_pending"}))

	if err := os.Remove(filepath.Join(state, ".build_status.json")); err != nil {
		t.Fatal(err)
	}
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"appr1","finished_at":"2026-09-17T01:02:00Z","status":"approval_rejected","trigger":"approval","duration_seconds":0}`)
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260917010000","finished_at":"2026-09-17T01:00:00Z","status":"failure_build","trigger":"manual","duration_seconds":3}`)
	writeTestJSON(t, filepath.Join(state, ".pending_transfers"), []map[string]any{{"id": "p1"}})
	if err := os.WriteFile(filepath.Join(state, ".notify_pending"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	before = snapshotFiles(t, state, []string{".build_status.json", ".build_history", ".pending_transfers", ".notify_pending"})
	resp = apiRequest(t, server, http.MethodGet, "/api/health", "", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("health degraded code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &health)
	if health["status"] != "degraded" || health["last_build_status"] != "failure" || health["pending_transfers"].(float64) != 1 || health["last_deploy_at"] != nil || health["last_deploy_status"] != nil {
		t.Fatalf("unexpected degraded health: %#v", health)
	}
	expectedChecks := []any{"build_status_missing", "pending_transfers_present", "notify_pending_read_error"}
	if got := health["checks"].([]any); !anySlicesEqual(got, expectedChecks) {
		t.Fatalf("unexpected degraded checks: %#v", got)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_status.json", ".build_history", ".pending_transfers", ".notify_pending"}))
}

func TestAPIReadOnlyAggregateEndpoints(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")

	siteDir := filepath.Join(state, "site")
	if err := os.MkdirAll(siteDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(siteDir, "index.html"), []byte("<h1>Adlaire</h1>\n"), 0600); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(state, ".branch_config"), apiBranchConfigFile{BranchTargets: []apiBranchTarget{{Branch: "main", TargetFile: "docs", SHAFile: filepath.Join(state, ".last_sha"), Src: filepath.Join(state, "repo", "docs"), Out: siteDir}}})
	if err := os.MkdirAll(filepath.Join(state, ".build_logs"), 0755); err != nil {
		t.Fatal(err)
	}
	meta, err := inspectOutput(siteDir)
	if err != nil {
		t.Fatal(err)
	}
	outputSize := meta.SizeBytes
	outputSHA := meta.SHA256
	writeTestJSON(t, filepath.Join(state, ".build_logs", "b20260917010101.json"), apiBuildLog{
		ID: "b20260917010101", StartedAt: "2026-09-17T01:01:01Z", FinishedAt: "2026-09-17T01:01:06Z",
		TargetStatus: "success", Commit: map[string]any{"sha": "abc123"}, Warnings: []string{"[WARNING] slow"}, DurationSeconds: 5,
	})
	historyData, err := json.Marshal(apiHistoryRecord{ID: "b20260917010101", Branch: "main", TargetFile: "docs", FinishedAt: "2026-09-17T01:01:06Z", Status: "success", Trigger: "manual", DurationSeconds: 5, OutputSizeBytes: &outputSize, OutputSHA256: &outputSHA})
	if err != nil {
		t.Fatal(err)
	}
	appendLine(t, filepath.Join(state, ".build_history"), string(historyData))
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260916010101","finished_at":"2026-09-16T01:01:06Z","status":"failure","trigger":"manual","duration_seconds":7,"commit_sha":"def456","failure_category":"pipeline_exit"}`)
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"b20260916010101","finished_at":"2026-09-18T01:01:06Z","status":"success","trigger":"manual","duration_seconds":1}`)
	appendLine(t, filepath.Join(state, ".notify_log"), `{"at":"2026-09-17T01:01:07Z","type":"build","result":"sent"}`)

	resp := apiRequest(t, server, http.MethodGet, "/api/sysinfo", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("sysinfo code=%d body=%s", resp.Code, resp.Body.String())
	}
	var sysinfo map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &sysinfo)
	if sysinfo["output_size_bytes"].(float64) == 0 || sysinfo["output_mtime"] == nil || sysinfo["uptime_seconds"].(float64) < 0 {
		t.Fatalf("unexpected sysinfo: %#v", sysinfo)
	}
	if _, exists := sysinfo["output_exists"]; exists {
		t.Fatalf("sysinfo must not expose output_exists: %#v", sysinfo)
	}
	if _, exists := sysinfo["state_dir"]; exists {
		t.Fatalf("sysinfo must not expose state_dir: %#v", sysinfo)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/sysinfo?extra=1", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("sysinfo unknown query code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/stats?days=7", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("stats code=%d body=%s", resp.Code, resp.Body.String())
	}
	var stats map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &stats)
	if stats["total_builds"].(float64) != 2 || stats["success_count"].(float64) != 1 || stats["failure_count"].(float64) != 1 || stats["success_rate"].(float64) != 0.5 {
		t.Fatalf("unexpected stats: %#v", stats)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/stats?days=7&days=8", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("stats duplicate query code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/stats/timeline?days=7", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("timeline code=%d body=%s", resp.Code, resp.Body.String())
	}
	var timeline map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &timeline)
	timelineItems := timeline["timeline"].([]any)
	if len(timelineItems) != 2 || timelineItems[0].(map[string]any)["date"] != "2026-09-17" {
		t.Fatalf("unexpected timeline: %#v", timeline)
	}
	if _, exists := timelineItems[0].(map[string]any)["total"]; exists {
		t.Fatalf("timeline item must not expose total: %#v", timelineItems[0])
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
	recent := duration["recent"].([]any)
	if recent[0].(map[string]any)["target_status"] != "success" {
		t.Fatalf("duration recent must expose target_status: %#v", recent)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/logs/search?q=slow&level=warn", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("logs search code=%d body=%s", resp.Code, resp.Body.String())
	}
	var search map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &search)
	if search["query"] != "slow" || search["from"] != nil || search["to"] != nil || search["level"] != "WARNING" || len(search["results"].([]any)) != 1 {
		t.Fatalf("unexpected logs search: %#v", search)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/logs/search?limit=1", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("logs search unknown query code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/logs/search?q=slow&q=done", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("logs search duplicate query code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/logs?secret=value", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("logs unknown query code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/history?failure_category=pipeline_exit", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("history failure_category code=%d body=%s", resp.Code, resp.Body.String())
	}
	var history map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &history)
	historyItems := history["history"].([]any)
	if history["total"].(float64) != 1 || historyItems[0].(map[string]any)["id"] != "b20260916010101" || historyItems[0].(map[string]any)["sha"] != "def456" || historyItems[0].(map[string]any)["build_at"] == "" {
		t.Fatalf("unexpected history failure_category result: %#v", history)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/history?failure_category=invalid", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("history invalid failure_category code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/output-meta", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("output-meta code=%d body=%s", resp.Code, resp.Body.String())
	}
	var outputMeta map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &outputMeta)
	if outputMeta["sha256"] != outputSHA || outputMeta["build_id"] != "b20260917010101" || outputMeta["commit_sha"] != "abc123" {
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
	if _, ok := dashboard["alerts"].([]any); !ok {
		t.Fatalf("dashboard alerts must be an array: %#v", dashboard["alerts"])
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
	systemdStartCalls := []string{}
	server.cfg.CommandRunner = func(ctx context.Context, name string, args ...string) (apiCommandResult, error) {
		call := strings.TrimSpace(name + " " + strings.Join(args, " "))
		systemdStartCalls = append(systemdStartCalls, call)
		if call == "systemctl start --no-block adlaire-ci.service" {
			return apiCommandResult{ExitCode: 0}, nil
		}
		return apiCommandResult{ExitCode: -1}, errors.New("unexpected command")
	}
	token := login(t, server, "password123")

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
	if _, exists := backup["config"]; exists {
		t.Fatalf("backup must not use nested config: %#v", backup)
	}
	for _, key := range []string{"exported_at", "server_config", "notify_config", "repo_config", "branch_config", "access_control", "hooks", "alert_rules", "tag_rules", "pipeline_config", "dashboard_layout", "smtp_config", "webhook_secret_set", "smtp_password_set"} {
		if _, exists := backup[key]; !exists {
			t.Fatalf("backup missing %q: %#v", key, backup)
		}
	}
	if backup["webhook_secret_set"] != true || backup["smtp_password_set"] != false {
		t.Fatalf("unexpected backup secret markers: %#v", backup)
	}
	webhookSHA := "1111111111111111111111111111111111111111"

	if err := os.WriteFile(filepath.Join(state, ".maintenance"), []byte(`{"enabled":`), 0600); err != nil {
		t.Fatal(err)
	}
	before := snapshotFiles(t, state, []string{".webhook_events.json", ".build_state", ".maintenance"})
	resp = signedWebhookRequest(t, server, secret, []byte(`{`), "delivery-invalid-json", "")
	if resp.Code != http.StatusBadRequest || !strings.Contains(resp.Body.String(), `"Invalid JSON"`) {
		t.Fatalf("webhook invalid json code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".webhook_events.json", ".build_state", ".maintenance"}))
	resp = signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"`+webhookSHA+`"}`), "delivery-maintenance", "")
	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("webhook with corrupt maintenance code=%d body=%s", resp.Code, resp.Body.String())
	}
	if !strings.Contains(resp.Body.String(), `"maintenance_unavailable"`) {
		t.Fatalf("unexpected webhook maintenance error: %s", resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".webhook_events.json", ".build_state", ".maintenance"}))
	if err := os.Remove(filepath.Join(state, ".maintenance")); err != nil {
		t.Fatal(err)
	}

	writeTestJSON(t, filepath.Join(state, ".branch_config"), apiBranchConfigFile{BranchTargets: []apiBranchTarget{{Branch: "release", TargetFile: "docs", SHAFile: filepath.Join(state, ".last_sha"), Src: filepath.Join(state, "repo", "docs"), Out: filepath.Join(state, "dist", "site")}}})
	before = snapshotFiles(t, state, []string{".build_state"})
	resp = signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"`+webhookSHA+`"}`), "delivery-ignored", "")
	if resp.Code != http.StatusAccepted {
		t.Fatalf("ignored branch webhook code=%d body=%s", resp.Code, resp.Body.String())
	}
	var ignored map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &ignored)
	if ignored["queued"] != false {
		t.Fatalf("ignored branch queued build: %#v", ignored)
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_state"}))
	if err := os.Remove(filepath.Join(state, ".branch_config")); err != nil {
		t.Fatal(err)
	}

	before = snapshotFiles(t, state, []string{".webhook_events.json", ".build_state", ".build_history"})
	resp = signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"`+webhookSHA+`"}`), "delivery-1", "sha256=bad")
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("invalid webhook code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".webhook_events.json", ".build_state", ".build_history"}))
	resp = signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"`+webhookSHA+`"}`), "delivery-1", "")
	if resp.Code != http.StatusAccepted {
		t.Fatalf("signed webhook code=%d body=%s", resp.Code, resp.Body.String())
	}
	var webhookResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &webhookResp)
	if webhookResp["queued"] != true || webhookResp["queue_id"] == nil || webhookResp["dispatch"] != "requested" {
		t.Fatalf("unexpected webhook response: %#v", webhookResp)
	}
	if len(systemdStartCalls) != 1 || systemdStartCalls[0] != "systemctl start --no-block adlaire-ci.service" {
		t.Fatalf("webhook dispatch calls mismatch: %#v", systemdStartCalls)
	}
	var buildState apiBuildState
	readTestJSON(t, filepath.Join(state, ".build_state"), &buildState)
	if len(buildState.Queued) != 1 || buildState.Queued[0]["trigger"] != "webhook" {
		t.Fatalf("webhook did not queue build: %#v", buildState)
	}
	if buildState.Queued[0]["id"] != "q20260917010101" {
		t.Fatalf("unexpected webhook queue id: %#v", buildState.Queued[0])
	}
	webhookAudits := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"build_trigger": true})
	if len(webhookAudits) != 1 {
		t.Fatalf("expected one webhook build audit record: %#v", webhookAudits)
	}
	if webhookAudits[0]["timestamp"] == "" || webhookAudits[0]["request_id"] == "" || webhookAudits[0]["actor_type"] != "webhook" || webhookAudits[0]["actor_id"] != "webhook" || webhookAudits[0]["target_type"] != "build" || webhookAudits[0]["target_id"] != "q20260917010101" || webhookAudits[0]["result"] != "success" || webhookAudits[0]["remote_addr"] != "192.0.2.1" || webhookAudits[0]["message"] != nil {
		t.Fatalf("unexpected webhook audit record: %#v", webhookAudits[0])
	}
	webhookEntry := buildState.Queued[0]
	buildState.ActiveQueueEntry = webhookEntry
	buildState.Queued = []map[string]any{}
	writeTestJSON(t, filepath.Join(state, ".build_state"), buildState)
	resp = signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"`+webhookSHA+`"}`), "delivery-1", "")
	if resp.Code != http.StatusAccepted {
		t.Fatalf("duplicate active webhook code=%d body=%s", resp.Code, resp.Body.String())
	}
	decodeTestJSON(t, resp.Body.Bytes(), &webhookResp)
	if webhookResp["message"] != "Webhook already queued" || webhookResp["queue_id"] != webhookEntry["id"] {
		t.Fatalf("unexpected duplicate webhook response: %#v", webhookResp)
	}
	if len(systemdStartCalls) != 2 || systemdStartCalls[1] != "systemctl start --no-block adlaire-ci.service" {
		t.Fatalf("duplicate webhook dispatch calls mismatch: %#v", systemdStartCalls)
	}
	readTestJSON(t, filepath.Join(state, ".build_state"), &buildState)
	if len(buildState.Queued) != 0 || buildState.ActiveQueueEntry["id"] != webhookEntry["id"] {
		t.Fatalf("duplicate webhook changed queue state: %#v", buildState)
	}
	before = snapshotFiles(t, state, []string{".build_state"})
	resp = signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"2222222222222222222222222222222222222222"}`), "delivery-1", "")
	if resp.Code != http.StatusConflict || !strings.Contains(resp.Body.String(), `"Conflicting delivery"`) {
		t.Fatalf("conflicting delivery code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_state"}))

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
	foundSnapshotConfig := false
	for _, record := range readJSONLines(filepath.Join(state, ".config_log")) {
		if record["type"] == "snapshot_delete" {
			changes := record["changes"].(map[string]any)
			if changes["id"] != "snap1" {
				t.Fatalf("unexpected snapshot config log: %#v", record)
			}
			foundSnapshotConfig = true
		}
	}
	if !foundSnapshotConfig {
		t.Fatalf("missing snapshot delete config log")
	}
	snapshotAudits := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"snapshot_delete": true})
	if len(snapshotAudits) != 1 || snapshotAudits[0]["actor_type"] != "admin" || snapshotAudits[0]["actor_id"] != "admin" || snapshotAudits[0]["target_type"] != "snapshot" || snapshotAudits[0]["target_id"] != "snap1" || snapshotAudits[0]["result"] != "success" {
		t.Fatalf("unexpected snapshot delete audit log: %#v", snapshotAudits)
	}

	resp = apiRequest(t, server, http.MethodPost, "/api/tokens", token, map[string]any{"label": "automation", "scopes": []string{"read"}})
	if resp.Code != http.StatusCreated {
		t.Fatalf("token issue code=%d body=%s", resp.Code, resp.Body.String())
	}
	var issued map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &issued)
	rawToken, _ := issued["token"].(string)
	if issued["id"] != "tok000001" || !strings.HasPrefix(rawToken, "act_") || len(rawToken) != 47 || issued["label"] != "automation" || issued["created_at"] == "" || issued["expires_at"] != nil || issued["last_used_at"] != nil || issued["revoked_at"] != nil {
		t.Fatalf("missing issued token: %#v", issued)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/tokens", token, nil)
	if resp.Code != http.StatusOK || strings.Contains(resp.Body.String(), rawToken) || strings.Contains(resp.Body.String(), "token_hash") || strings.Contains(resp.Body.String(), `"name"`) {
		t.Fatalf("token list leaked raw token: code=%d body=%s", resp.Code, resp.Body.String())
	}
	var stored apiTokenFile
	readTestJSON(t, filepath.Join(state, ".api_tokens"), &stored)
	if len(stored.Tokens) != 1 || stored.Tokens[0].ID != "tok000001" || stored.Tokens[0].Label != "automation" || stored.Tokens[0].TokenHash == "" || strings.Contains(stored.Tokens[0].TokenHash, rawToken) || stored.Tokens[0].LastUsedAt != nil || stored.Tokens[0].ExpiresAt != nil {
		t.Fatalf("token hash not stored safely: %#v", stored)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", rawToken, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("read-scope token status code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, filepath.Join(state, ".api_tokens"), &stored)
	if len(stored.Tokens) != 1 || stored.Tokens[0].LastUsedAt == nil {
		t.Fatalf("token last_used_at was not updated: %#v", stored)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/config", rawToken, map[string]any{"log_level": "DEBUG"})
	if resp.Code != http.StatusForbidden {
		t.Fatalf("read-scope token config code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/tokens", token, map[string]any{"label": "bad", "scopes": []string{}})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("empty scopes code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/tokens/missing", token, nil)
	if resp.Code != http.StatusNotFound {
		t.Fatalf("missing token revoke code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/tokens/tok000001", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("token revoke code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/tokens/tok000001", token, nil)
	if resp.Code != http.StatusNotFound {
		t.Fatalf("revoked token rerevoke code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", rawToken, nil)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("revoked token status code=%d body=%s", resp.Code, resp.Body.String())
	}
	tokenAccessLog := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{
		"token_create":         true,
		"token_auth":           true,
		"permission_denied":    true,
		"token_revoke":         true,
		"token_revoked_reject": true,
	})
	if actionsOf(tokenAccessLog) != "token_create,token_auth,permission_denied,token_revoke,token_revoked_reject" {
		t.Fatalf("unexpected token access log actions: %#v", tokenAccessLog)
	}
	for _, record := range tokenAccessLog {
		if record["token_id"] != "tok000001" || record["remote_addr"] != "192.0.2.1" || record["session_id"] != nil {
			t.Fatalf("unexpected token access log record: %#v", record)
		}
		if record["action"] == "permission_denied" && record["reason"] != "scope_denied" {
			t.Fatalf("permission denial reason not fixed: %#v", record)
		}
	}
	tokenAuditLog := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{
		"token_create":         true,
		"token_auth":           true,
		"permission_denied":    true,
		"token_revoke":         true,
		"token_revoked_reject": true,
	})
	if actionsOf(tokenAuditLog) != "token_create,token_auth,permission_denied,token_revoke,token_revoked_reject" {
		t.Fatalf("unexpected token audit log actions: %#v", tokenAuditLog)
	}
	for _, record := range tokenAuditLog {
		if record["timestamp"] == "" || record["request_id"] == "" || record["remote_addr"] != "192.0.2.1" || record["message"] != nil {
			t.Fatalf("token audit log must use strict schema fields: %#v", record)
		}
		switch record["action"] {
		case "token_create", "token_revoke":
			if record["actor_type"] != "admin" || record["actor_id"] != "admin" || record["target_type"] != "api_token" || record["target_id"] != "tok000001" || record["result"] != "success" {
				t.Fatalf("unexpected admin token audit record: %#v", record)
			}
		case "token_auth":
			if record["actor_type"] != "api_token" || record["actor_id"] != "tok000001" || record["target_type"] != "api_token" || record["target_id"] != "tok000001" || record["result"] != "success" {
				t.Fatalf("unexpected token auth audit record: %#v", record)
			}
		case "permission_denied":
			if record["actor_type"] != "api_token" || record["actor_id"] != "tok000001" || record["target_type"] != "endpoint" || record["target_id"] != "POST /api/config" || record["result"] != "denied" {
				t.Fatalf("unexpected permission audit record: %#v", record)
			}
		case "token_revoked_reject":
			if record["actor_type"] != "api_token" || record["actor_id"] != "tok000001" || record["target_type"] != "api_token" || record["target_id"] != "tok000001" || record["result"] != "failure" {
				t.Fatalf("unexpected revoked audit record: %#v", record)
			}
		}
	}
	if serialized := string(mustMarshalTest(t, tokenAuditLog)); strings.Contains(serialized, rawToken) || strings.Contains(serialized, "Authorization") {
		t.Fatalf("token audit log leaked secret material: %s", serialized)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/tokens", token, map[string]any{"label": "trigger", "scopes": []string{"trigger"}})
	if resp.Code != http.StatusCreated {
		t.Fatalf("trigger token issue code=%d body=%s", resp.Code, resp.Body.String())
	}
	var triggerIssued map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &triggerIssued)
	triggerToken, _ := triggerIssued["token"].(string)
	if triggerIssued["id"] != "tok000002" || triggerToken == "" {
		t.Fatalf("unexpected trigger token: %#v", triggerIssued)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/build", triggerToken, nil)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("trigger token build code=%d body=%s", resp.Code, resp.Body.String())
	}
	var triggerBuild map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &triggerBuild)
	triggerQueueID, _ := triggerBuild["queue_id"].(string)
	if triggerQueueID == "" {
		t.Fatalf("trigger build missing queue id: %#v", triggerBuild)
	}
	readTestJSON(t, filepath.Join(state, ".build_state"), &buildState)
	foundTriggerQueue := false
	for _, entry := range buildState.Queued {
		if entry["id"] == triggerQueueID {
			foundTriggerQueue = entry["requested_by"] == "tok000002" && entry["trigger"] == "manual"
		}
	}
	if !foundTriggerQueue {
		t.Fatalf("trigger token queue entry missing actor: %#v", buildState.Queued)
	}
	buildAudits := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"build_trigger": true})
	if len(buildAudits) == 0 {
		t.Fatalf("missing build trigger audit")
	}
	lastBuildAudit := buildAudits[len(buildAudits)-1]
	if lastBuildAudit["actor_type"] != "api_token" || lastBuildAudit["actor_id"] != "tok000002" || lastBuildAudit["target_type"] != "build" || lastBuildAudit["target_id"] != triggerQueueID || lastBuildAudit["result"] != "success" {
		t.Fatalf("unexpected trigger token build audit: %#v", lastBuildAudit)
	}

	restoreBody := cloneTestMap(t, backup)
	restoreBody["webhook_secret"] = "restored-webhook-secret"
	restoreBody["smtp_password"] = "restored-smtp-secret"
	serverConfig := restoreBody["server_config"].(map[string]any)
	serverConfig["log_level"] = "DEBUG"
	repoConfig := restoreBody["repo_config"].(map[string]any)
	repoConfig["owner"] = "example"
	repoConfig["repo"] = "docs"
	repoConfig["updated_at"] = nil
	resp = apiRequest(t, server, http.MethodPost, "/api/restore", token, restoreBody)
	if resp.Code != http.StatusOK {
		t.Fatalf("restore code=%d body=%s", resp.Code, resp.Body.String())
	}
	var restoredServer map[string]any
	readTestJSON(t, filepath.Join(state, ".server_config"), &restoredServer)
	if restoredServer["log_level"] != "DEBUG" {
		t.Fatalf("server config was not restored: %#v", restoredServer)
	}
	var restoredRepo map[string]any
	readTestJSON(t, filepath.Join(state, ".repo_config"), &restoredRepo)
	if restoredRepo["owner"] != "example" || restoredRepo["repo"] != "docs" || restoredRepo["updated_at"] == nil {
		t.Fatalf("repo config was not restored: %#v", restoredRepo)
	}
	restoredWebhookSecret, err := os.ReadFile(filepath.Join(state, ".webhook_secret"))
	if err != nil || string(restoredWebhookSecret) != "restored-webhook-secret" {
		t.Fatalf("webhook secret restore mismatch: value=%q err=%v", string(restoredWebhookSecret), err)
	}
	restoredSMTPSecret, err := os.ReadFile(filepath.Join(state, ".smtp_secret"))
	if err != nil || string(restoredSMTPSecret) != "restored-smtp-secret" {
		t.Fatalf("smtp secret restore mismatch: value=%q err=%v", string(restoredSMTPSecret), err)
	}
	configLog, err := os.ReadFile(filepath.Join(state, ".config_log"))
	if err != nil || !strings.Contains(string(configLog), `"type":"restore"`) || strings.Contains(string(configLog), "restored-webhook-secret") {
		t.Fatalf("restore config log mismatch: log=%s err=%v", string(configLog), err)
	}
	auditLog, err := os.ReadFile(filepath.Join(state, ".audit_log"))
	if err != nil || !strings.Contains(string(auditLog), `"action":"config_update"`) {
		t.Fatalf("restore audit log mismatch: log=%s err=%v", string(auditLog), err)
	}
	before = snapshotFiles(t, state, []string{".server_config", ".repo_config", ".webhook_secret", ".smtp_secret", ".config_log", ".audit_log"})
	resp = apiRequest(t, server, http.MethodPost, "/api/restore", token, restoreBody)
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"No changes"`) {
		t.Fatalf("restore no-op code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".server_config", ".repo_config", ".webhook_secret", ".smtp_secret", ".config_log", ".audit_log"}))
	badRestore := cloneTestMap(t, restoreBody)
	badRestore["config"] = map[string]any{}
	resp = apiRequest(t, server, http.MethodPost, "/api/restore", token, badRestore)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("restore unknown key code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".server_config", ".repo_config", ".webhook_secret", ".smtp_secret", ".config_log", ".audit_log"}))
}

func TestAPIWebhookAuditFailureDoesNotDispatch(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	dispatchCalls := 0
	server.cfg.CommandRunner = func(ctx context.Context, name string, args ...string) (apiCommandResult, error) {
		dispatchCalls++
		return apiCommandResult{ExitCode: 0}, nil
	}
	secret := "phase4-webhook-secret"
	if err := os.WriteFile(filepath.Join(state, ".webhook_secret"), []byte(secret), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(state, ".audit_log"), 0700); err != nil {
		t.Fatal(err)
	}
	webhookSHA := "1111111111111111111111111111111111111111"
	resp := signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"`+webhookSHA+`"}`), "delivery-audit-fail", "")
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("webhook audit failure code=%d body=%s", resp.Code, resp.Body.String())
	}
	if dispatchCalls != 0 {
		t.Fatalf("webhook dispatched despite audit failure: %d", dispatchCalls)
	}
	var buildState apiBuildState
	readTestJSON(t, filepath.Join(state, ".build_state"), &buildState)
	if len(buildState.Queued) != 1 || buildState.Queued[0]["id"] != "q20260917010101" {
		t.Fatalf("queue state should remain for audit partial failure: %#v", buildState)
	}
}

func TestAPIPhase4TrendChainAndApprovalFixtures(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")

	writeTestJSON(t, filepath.Join(state, ".build_trends.json"), apiBuildTrendFile{
		SchemaVersion: 1,
		Samples: []apiBuildTrendSample{
			{BuildID: "b1", FinishedAt: "2026-09-17T01:00:00Z", Branch: "main", Trigger: "manual", DurationSeconds: 10, Status: "success", TargetStatus: "success"},
			{BuildID: "b2", FinishedAt: "2026-09-17T01:01:00Z", Branch: "main", Trigger: "polling", DurationSeconds: 20, Status: "failure", TargetStatus: "failure_build", Anomaly: true},
			{BuildID: "b3", FinishedAt: "2026-09-17T01:02:00Z", Branch: "main", Trigger: "manual", DurationSeconds: 30, Status: "success", TargetStatus: "success"},
		},
	})
	resp := apiRequest(t, server, http.MethodGet, "/api/stats/build-trends?n=2", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("build trends code=%d body=%s", resp.Code, resp.Body.String())
	}
	var trends struct {
		Count   int                   `json:"count"`
		Samples []apiBuildTrendSample `json:"samples"`
		Summary map[string]any        `json:"summary"`
	}
	decodeTestJSON(t, resp.Body.Bytes(), &trends)
	if trends.Count != 2 || len(trends.Samples) != 2 || trends.Samples[0].BuildID != "b2" || trends.Samples[1].BuildID != "b3" {
		t.Fatalf("unexpected trends response: %#v", trends)
	}
	if trends.Summary["avg_seconds"].(float64) != 25 || trends.Summary["anomaly_count"].(float64) != 1 {
		t.Fatalf("unexpected trend summary: %#v", trends.Summary)
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/build-chain-config", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("build chain default code=%d body=%s", resp.Code, resp.Body.String())
	}
	var chain apiBuildChainConfig
	decodeTestJSON(t, resp.Body.Bytes(), &chain)
	if len(chain.Chains) != 0 {
		t.Fatalf("unexpected default chain: %#v", chain)
	}
	validChains := map[string]any{"chains": []map[string]any{
		{"id": "base", "branch": "main", "target_file": "README.md", "depends_on": []string{}, "required": true, "enabled": true},
		{"id": "docs", "branch": "main", "target_file": "docs/SPEC.md", "depends_on": []string{"base"}, "required": false, "enabled": true},
	}}
	resp = apiRequest(t, server, http.MethodPost, "/api/build-chain-config", token, validChains)
	if resp.Code != http.StatusOK {
		t.Fatalf("build chain update code=%d body=%s", resp.Code, resp.Body.String())
	}
	before := snapshotFiles(t, state, []string{".build_chain_config", ".config_log"})
	cycle := map[string]any{"chains": []map[string]any{
		{"id": "base", "branch": "main", "target_file": "README.md", "depends_on": []string{"docs"}, "required": true, "enabled": true},
		{"id": "docs", "branch": "main", "target_file": "docs/SPEC.md", "depends_on": []string{"base"}, "required": false, "enabled": true},
	}}
	resp = apiRequest(t, server, http.MethodPost, "/api/build-chain-config", token, cycle)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("build chain cycle code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".build_chain_config", ".config_log"}))

	appendLine(t, filepath.Join(state, ".approval_queue"), `{"id":"appr1","status":"pending","branch":"main","sha":"1111111111111111111111111111111111111111","target":"README.md","requested_trigger":"manual","requested_force":true,"requested_by":"runner","created_at":"2026-09-17T01:01:01Z","expires_at":"2026-09-17T02:01:01Z"}`)
	appendLine(t, filepath.Join(state, ".approval_queue"), `{`)
	resp = apiRequest(t, server, http.MethodGet, "/api/approvals", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("approvals list code=%d body=%s", resp.Code, resp.Body.String())
	}
	var approvals struct {
		Approvals []apiApprovalRecord `json:"approvals"`
	}
	decodeTestJSON(t, resp.Body.Bytes(), &approvals)
	if len(approvals.Approvals) != 1 || approvals.Approvals[0].ID != "appr1" || approvals.Approvals[0].Status != "pending" {
		t.Fatalf("unexpected approvals list: %#v", approvals)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/approvals/appr1/approve", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("approve code=%d body=%s", resp.Code, resp.Body.String())
	}
	var approveResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &approveResp)
	if approveResp["queued"] != true || approveResp["queue_id"] == nil || approveResp["dispatch"] == nil {
		t.Fatalf("unexpected approve response: %#v", approveResp)
	}
	var buildState apiBuildState
	readTestJSON(t, filepath.Join(state, ".build_state"), &buildState)
	if len(buildState.Queued) != 1 || buildState.Queued[0]["trigger"] != "approval" {
		t.Fatalf("approval did not queue build: %#v", buildState)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/approvals/appr1/approve", token, nil)
	if resp.Code != http.StatusConflict {
		t.Fatalf("duplicate approve code=%d body=%s", resp.Code, resp.Body.String())
	}

	appendLine(t, filepath.Join(state, ".approval_queue"), `{"id":"appr2","status":"pending","branch":"main","sha":"2222222222222222222222222222222222222222","target":"README.md","requested_trigger":"webhook","requested_force":false,"requested_by":"runner","created_at":"2026-09-17T01:03:01Z","expires_at":"2026-09-17T02:03:01Z"}`)
	resp = apiRequest(t, server, http.MethodPost, "/api/approvals/appr2/reject", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("reject code=%d body=%s", resp.Code, resp.Body.String())
	}
	historyData, err := os.ReadFile(filepath.Join(state, ".build_history"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(historyData), `"status":"approval_rejected"`) {
		t.Fatalf("approval rejection history missing: %s", string(historyData))
	}
}

func TestAPIPhase4RulePipelineNotesAndLayoutFixtures(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")

	rule := map[string]any{"name": "failures", "status": "failure"}
	resp := apiRequest(t, server, http.MethodPost, "/api/alert-rules", token, rule)
	if resp.Code != http.StatusCreated {
		t.Fatalf("alert rule code=%d body=%s", resp.Code, resp.Body.String())
	}
	var createdRule map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &createdRule)
	if createdRule["id"] != "r20260917010101" || createdRule["rule"] != nil || createdRule["name"] != "failures" {
		t.Fatalf("alert rule response must be the rule record: %#v", createdRule)
	}
	before := snapshotFiles(t, state, []string{".alert_rules", ".config_log"})
	resp = apiRequest(t, server, http.MethodPost, "/api/alert-rules", token, rule)
	if resp.Code != http.StatusConflict {
		t.Fatalf("duplicate alert rule code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".alert_rules", ".config_log"}))

	resp = apiRequest(t, server, http.MethodGet, "/api/pipeline-config", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("pipeline config default code=%d body=%s", resp.Code, resp.Body.String())
	}
	var pipelineDefault map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &pipelineDefault)
	if len(pipelineDefault) != 2 || len(pipelineDefault["extra_args"].([]any)) != 0 || len(pipelineDefault["env"].(map[string]any)) != 0 {
		t.Fatalf("unexpected pipeline config default: %#v", pipelineDefault)
	}
	before = snapshotFiles(t, state, []string{".pipeline_config"})
	resp = apiRequest(t, server, http.MethodPost, "/api/pipeline-config", token, map[string]any{"extra_args": []string{"--src"}, "env": map[string]string{}})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("reserved pipeline arg code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".pipeline_config"}))
	resp = apiRequest(t, server, http.MethodPost, "/api/pipeline-config", token, map[string]any{"extra_args": []string{}, "env": map[string]string{}, "unknown_key": nil})
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unknown pipeline config key code=%d body=%s", resp.Code, resp.Body.String())
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
	token := login(t, server, "password123")

	appendLine(t, filepath.Join(state, ".audit_log"), `{"at":"2026-09-17T01:01:02Z","actor":"admin","action":"config_update","result":"success"}`)
	appendLine(t, filepath.Join(state, ".audit_log"), `{"timestamp":"2026-09-17T01:01:03Z","request_id":"11111111111111111111111111111111","actor_type":"admin","actor_id":"admin","action":"config_update","target_type":"config","target_id":"server_config","result":"success","remote_addr":"192.0.2.1","message":null}`)
	resp := apiRequest(t, server, http.MethodGet, "/api/audit-log?action=config_update&actor=admin", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("audit-log code=%d body=%s", resp.Code, resp.Body.String())
	}
	var logResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &logResp)
	if logResp["total"].(float64) != 2 {
		t.Fatalf("unexpected audit log: %#v", logResp)
	}
	logItems := logResp["log"].([]any)
	firstAudit := logItems[0].(map[string]any)
	if firstAudit["timestamp"] != "2026-09-17T01:01:03Z" || firstAudit["actor_id"] != "admin" {
		t.Fatalf("audit log must sort by timestamp and filter actor_id: %#v", logResp)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/audit-log?action=unknown_action", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unknown audit action code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/audit-log?result=unknown", token, nil)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unknown audit result code=%d body=%s", resp.Code, resp.Body.String())
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
	resp = apiRequest(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "password123"})
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
	var totpSuccess map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &totpSuccess)
	totpToken, _ := totpSuccess["token"].(string)
	if totpToken == "" {
		t.Fatalf("totp login did not return token: %#v", totpSuccess)
	}
	replayResp := apiRequest(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "password123"})
	if replayResp.Code != http.StatusOK {
		t.Fatalf("totp replay stage1 code=%d body=%s", replayResp.Code, replayResp.Body.String())
	}
	var replayLogin map[string]any
	decodeTestJSON(t, replayResp.Body.Bytes(), &replayLogin)
	replayTicket, _ := replayLogin["ticket"].(string)
	replayResp = apiRequest(t, server, http.MethodPost, "/api/login/totp", "", map[string]string{"ticket": replayTicket, "code": code})
	if replayResp.Code != http.StatusUnauthorized {
		t.Fatalf("totp replay code should be rejected, got %d body=%s", replayResp.Code, replayResp.Body.String())
	}
	totpActions := map[string]bool{"totp_required": true, "login_success": true, "totp_failure": true}
	totpAccess := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), totpActions)
	if actionsOf(totpAccess) != "totp_required,login_success,totp_required,totp_failure" {
		t.Fatalf("unexpected totp access log actions: %#v", totpAccess)
	}
	totpAudit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), totpActions)
	if actionsOf(totpAudit) != "totp_required,login_success,totp_required,totp_failure" {
		t.Fatalf("unexpected totp audit log actions: %#v", totpAudit)
	}
	for _, record := range totpAudit {
		if record["actor_type"] != "anonymous" || record["actor_id"] != nil || record["target_type"] != "auth" || record["request_id"] == "" {
			t.Fatalf("unexpected totp audit schema: %#v", record)
		}
	}
	serializedTOTPLogs := string(mustMarshalTest(t, []any{totpAccess, totpAudit}))
	if strings.Contains(serializedTOTPLogs, ticket) || strings.Contains(serializedTOTPLogs, replayTicket) || strings.Contains(serializedTOTPLogs, totpToken) || strings.Contains(serializedTOTPLogs, secret) {
		t.Fatalf("totp login logs leaked secret material: %s", serializedTOTPLogs)
	}
	nextCode := totpCode(secret, server.cfg.Now().UTC().Unix()/30+1)
	invalidCode := "000000"
	if invalidCode == code || invalidCode == nextCode {
		invalidCode = "111111"
	}
	if invalidCode == code || invalidCode == nextCode {
		invalidCode = "222222"
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/auth/totp", token, map[string]string{"code": invalidCode})
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("invalid totp disable code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/auth/totp", token, map[string]string{"code": nextCode})
	if resp.Code != http.StatusOK {
		t.Fatalf("totp disable code=%d body=%s", resp.Code, resp.Body.String())
	}
	totpLifecycleActions := map[string]bool{"totp_setup": true, "totp_failure": true, "totp_enabled": true, "totp_disabled": true}
	totpLifecycleAccess := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), totpLifecycleActions)
	if actionsOf(totpLifecycleAccess) != "totp_setup,totp_enabled,totp_failure,totp_failure,totp_disabled" {
		t.Fatalf("unexpected totp lifecycle access log actions: %#v", totpLifecycleAccess)
	}
	totpLifecycleAudit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), totpLifecycleActions)
	if actionsOf(totpLifecycleAudit) != "totp_setup,totp_enabled,totp_failure,totp_failure,totp_disabled" {
		t.Fatalf("unexpected totp lifecycle audit log actions: %#v", totpLifecycleAudit)
	}
	expectedTargets := []string{"totp_setup", "totp", "login_totp", "totp_disable", "totp"}
	for i, record := range totpLifecycleAudit {
		if record["actor_type"] != "admin" && i != 2 {
			t.Fatalf("unexpected admin totp lifecycle audit actor: %#v", record)
		}
		if i == 2 && (record["actor_type"] != "anonymous" || record["actor_id"] != nil) {
			t.Fatalf("unexpected login totp failure audit actor: %#v", record)
		}
		if record["target_type"] != "auth" || record["target_id"] != expectedTargets[i] || record["request_id"] == "" {
			t.Fatalf("unexpected totp lifecycle audit schema: %#v", record)
		}
	}
	serializedTOTPLifecycleLogs := string(mustMarshalTest(t, []any{totpLifecycleAccess, totpLifecycleAudit}))
	if strings.Contains(serializedTOTPLifecycleLogs, secret) || strings.Contains(serializedTOTPLifecycleLogs, setup["otpauth_uri"].(string)) {
		t.Fatalf("totp lifecycle logs leaked setup secret material: %s", serializedTOTPLifecycleLogs)
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
	resp = apiRequest(t, server, http.MethodPost, "/api/api-rate-limit", token, testAPIRateLimitPolicy(nil))
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
	githubUserCalls := 0
	githubRateCalls := 0
	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected GitHub verification request: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer ghp_exampletoken" {
			t.Fatalf("unexpected authorization header: %q", r.Header.Get("Authorization"))
		}
		switch r.URL.Path {
		case "/user":
			githubUserCalls++
			w.Header().Set("X-OAuth-Scopes", "repo, contents:read, repo")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"login":"octocat"}`))
		case "/rate_limit":
			githubRateCalls++
			reset := time.Date(2026, 9, 17, 2, 0, 0, 0, time.UTC).Unix()
			_ = json.NewEncoder(w).Encode(map[string]any{"resources": map[string]any{"core": map[string]any{"limit": 5000, "remaining": 4823, "reset": reset, "used": 177}}})
		default:
			t.Fatalf("unexpected GitHub path: %s", r.URL.Path)
		}
	}))
	defer github.Close()
	server.cfg.GitHubAPIBaseURL = github.URL
	server.cfg.HTTPClient = github.Client()
	expires := "2026-10-01"
	writeTestJSON(t, filepath.Join(state, ".server_config"), apiServerConfig{PATExpiresAt: &expires})
	resp = apiRequest(t, server, http.MethodGet, "/api/pat-status", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("pat status code=%d body=%s", resp.Code, resp.Body.String())
	}
	var patStatus map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &patStatus)
	if patStatus["configured"] != true || patStatus["expires_at"] != expires || patStatus["expires_in_days"].(float64) != 14 {
		t.Fatalf("unexpected pat status: %#v", patStatus)
	}
	if _, exists := patStatus["mode"]; exists {
		t.Fatalf("pat status must not expose mode: %#v", patStatus)
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/pat-verify", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("pat verify code=%d body=%s", resp.Code, resp.Body.String())
	}
	var patVerify map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &patVerify)
	if patVerify["valid"] != true || patVerify["checked_at"] != "2026-09-17T01:01:01Z" || githubUserCalls != 1 {
		t.Fatalf("unexpected pat verify response/calls: body=%#v calls=%d", patVerify, githubUserCalls)
	}
	if !anySlicesEqual(patVerify["scopes"].([]any), []any{"contents:read", "repo"}) {
		t.Fatalf("unexpected scopes: %#v", patVerify["scopes"])
	}
	if _, exists := patVerify["configured"]; exists {
		t.Fatalf("pat verify must not expose configured: %#v", patVerify)
	}
	if _, exists := patVerify["remote_checked"]; exists {
		t.Fatalf("pat verify must not expose remote_checked: %#v", patVerify)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/rate-limit", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("rate-limit code=%d body=%s", resp.Code, resp.Body.String())
	}
	var rateLimit map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &rateLimit)
	if rateLimit["limit"].(float64) != 5000 || rateLimit["remaining"].(float64) != 4823 || rateLimit["used"].(float64) != 177 || rateLimit["reset_at"] != "2026-09-17T02:00:00Z" || githubRateCalls != 1 {
		t.Fatalf("unexpected rate limit response/calls: body=%#v calls=%d", rateLimit, githubRateCalls)
	}
	if _, exists := rateLimit["remote_checked"]; exists {
		t.Fatalf("rate limit must not expose remote_checked: %#v", rateLimit)
	}

	hook := map[string]any{"phase": "after_build", "command_args": []string{"/bin/echo", "ok"}, "abort_on_failure": false}
	resp = apiRequest(t, server, http.MethodPost, "/api/hooks", token, hook)
	if resp.Code != http.StatusCreated {
		t.Fatalf("hook create code=%d body=%s", resp.Code, resp.Body.String())
	}
	var hookResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &hookResp)
	hookID, _ := hookResp["id"].(string)
	if hookID != "h20260917010101" || hookResp["phase"] != "after_build" || hookResp["enabled"] != true || hookResp["timeout_seconds"].(float64) != 300 {
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
	if tagRule["id"] != "t20260917010101" || tagRule["condition"] == nil || tagRule["rule"] != nil {
		t.Fatalf("tag rule response must be the tag rule record: %#v", tagRule)
	}

	siteDir := filepath.Join(state, "site")
	if err := os.MkdirAll(siteDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(siteDir, "index.html"), []byte("<h1>Adlaire</h1>\n"), 0600); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(state, ".branch_config"), apiBranchConfigFile{BranchTargets: []apiBranchTarget{{Branch: "main", TargetFile: "docs", SHAFile: filepath.Join(state, ".last_sha"), Src: filepath.Join(state, "repo", "docs"), Out: siteDir}}})
	meta, err := inspectOutput(siteDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(state, ".build_logs"), 0755); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(state, ".build_logs", "bverify.json"), apiBuildLog{ID: "bverify", StartedAt: "2026-09-17T01:01:01Z", FinishedAt: "2026-09-17T01:01:03Z", TargetStatus: "success", OutputSHA256: meta.SHA256})
	outputSize := meta.SizeBytes
	appendLine(t, filepath.Join(state, ".build_history"), `{"id":"bverify","branch":"main","target_file":"docs","finished_at":"2026-09-17T01:01:03Z","status":"success","trigger":"manual","duration_seconds":2,"output_size_bytes":`+strconv.FormatInt(outputSize, 10)+`,"output_sha256":"`+meta.SHA256+`"}`)
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

func TestAPISecurityCredentialsAndForcedGate(t *testing.T) {
	state := newAPIState(t)
	credPath := filepath.Join(state, ".admin_credentials")
	cred, err := readCredentials(credPath)
	if err != nil {
		t.Fatal(err)
	}
	cred.MustChange = true
	cred.LoginCount = 4
	cred.LastLoginAt = nil
	if err := atomicWriteJSON(credPath, cred, 0600); err != nil {
		t.Fatal(err)
	}
	server := newTestAPI(t, state)
	resp := apiRequest(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "password123"})
	if resp.Code != http.StatusOK {
		t.Fatalf("forced login code=%d body=%s", resp.Code, resp.Body.String())
	}
	var loginBody map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &loginBody)
	token, _ := loginBody["token"].(string)
	if token == "" || loginBody["must_change"] != "forced" {
		t.Fatalf("expected forced login response: %#v", loginBody)
	}
	readTestJSON(t, credPath, &cred)
	if cred.LoginCount != 5 || cred.LastLoginAt == nil {
		t.Fatalf("login_count was not advanced: %#v", cred)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusForbidden || !strings.Contains(resp.Body.String(), "Password change required") {
		t.Fatalf("forced session must be gated, code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/change-password", token, map[string]string{
		"current_password": "password123",
		"new_password":     "changed-password",
	})
	if resp.Code != http.StatusOK {
		t.Fatalf("forced change-password code=%d body=%s", resp.Code, resp.Body.String())
	}
	readTestJSON(t, credPath, &cred)
	if cred.MustChange || cred.Salt == "" || len(cred.Salt) != 64 {
		t.Fatalf("credentials were not normalized after change-password: %#v", cred)
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("forced gate should clear after password change, code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPICredentialsLockedUpdateSerializes(t *testing.T) {
	state := newAPIState(t)
	credPath := filepath.Join(state, ".admin_credentials")
	const total = 12
	var wg sync.WaitGroup
	errCh := make(chan error, total)
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := updateCredentialsLocked(credPath, func(current apiCredentials) (apiCredentials, bool, error) {
				advanceCredentialsLogin(&current, fmt.Sprintf("2026-09-17T01:01:%02dZ", i))
				return current, true, nil
			})
			if err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}
	var cred apiCredentials
	readTestJSON(t, credPath, &cred)
	if cred.LoginCount != total || cred.LastLoginAt == nil {
		t.Fatalf("locked credentials update lost increments: %#v", cred)
	}
}

func TestAPITokenLockedUpdateSerializes(t *testing.T) {
	state := newAPIState(t)
	tokenPath := filepath.Join(state, ".api_tokens")
	const total = 12
	var wg sync.WaitGroup
	errCh := make(chan error, total)
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := updateAPITokensLocked(tokenPath, func(tokens []apiTokenRecord) ([]apiTokenRecord, bool, error) {
				id, err := nextTokenID(tokens)
				if err != nil {
					return nil, false, err
				}
				tokens = append(tokens, apiTokenRecord{
					ID:        id,
					Label:     fmt.Sprintf("automation-%02d", i),
					Scopes:    []string{"read"},
					TokenHash: fmt.Sprintf("%064x", i+1),
					CreatedAt: fmt.Sprintf("2026-09-17T01:01:%02dZ", i),
				})
				return tokens, true, nil
			})
			if err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}
	var stored apiTokenFile
	readTestJSON(t, tokenPath, &stored)
	if len(stored.Tokens) != total {
		t.Fatalf("locked token update lost records: %#v", stored.Tokens)
	}
	seen := map[string]bool{}
	for _, record := range stored.Tokens {
		if seen[record.ID] {
			t.Fatalf("duplicate token id after locked updates: %#v", stored.Tokens)
		}
		seen[record.ID] = true
	}
}

func TestAPITokenStateLockConflictReturns409(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	sessionToken := login(t, server, "password123")
	resp := apiRequest(t, server, http.MethodPost, "/api/tokens", sessionToken, map[string]any{"label": "automation", "scopes": []string{"read"}})
	if resp.Code != http.StatusCreated {
		t.Fatalf("token issue code=%d body=%s", resp.Code, resp.Body.String())
	}
	var issued map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &issued)
	rawToken, _ := issued["token"].(string)
	if rawToken == "" {
		t.Fatalf("missing issued token: %#v", issued)
	}

	oldSleep := runnerSleep
	runnerSleep = func(time.Duration) {}
	defer func() { runnerSleep = oldSleep }()
	release, err := acquireStateFileLock(filepath.Join(state, ".api_tokens"))
	if err != nil {
		t.Fatal(err)
	}
	defer release()

	resp = apiRequest(t, server, http.MethodPost, "/api/tokens", sessionToken, map[string]any{"label": "locked", "scopes": []string{"read"}})
	if resp.Code != http.StatusConflict {
		t.Fatalf("token create lock conflict code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", rawToken, nil)
	if resp.Code != http.StatusConflict {
		t.Fatalf("token auth lock conflict code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodDelete, "/api/tokens/tok000001", sessionToken, nil)
	if resp.Code != http.StatusConflict {
		t.Fatalf("token revoke lock conflict code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestAPILoginSecurityLogs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		state := newAPIState(t)
		server := newTestAPI(t, state)
		resp := apiRequest(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "password123"})
		if resp.Code != http.StatusOK {
			t.Fatalf("login code=%d body=%s", resp.Code, resp.Body.String())
		}
		var body map[string]any
		decodeTestJSON(t, resp.Body.Bytes(), &body)
		token, _ := body["token"].(string)
		if token == "" {
			t.Fatalf("missing token: %#v", body)
		}
		access := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{"login_success": true})
		if len(access) != 1 || access[0]["result"] != "success" || access[0]["token_id"] != nil || access[0]["session_id"] != nil || access[0]["remote_addr"] != "192.0.2.1" || access[0]["reason"] != nil {
			t.Fatalf("unexpected login success access log: %#v", access)
		}
		audit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"login_success": true})
		if len(audit) != 1 || audit[0]["actor_type"] != "anonymous" || audit[0]["actor_id"] != nil || audit[0]["target_type"] != "auth" || audit[0]["target_id"] != "login" || audit[0]["result"] != "success" {
			t.Fatalf("unexpected login success audit log: %#v", audit)
		}
		serialized := string(mustMarshalTest(t, []any{access, audit}))
		if strings.Contains(serialized, token) || strings.Contains(serialized, "password123") {
			t.Fatalf("login success logs leaked secret material: %s", serialized)
		}
	})

	t.Run("failure", func(t *testing.T) {
		state := newAPIState(t)
		server := newTestAPI(t, state)
		resp := apiRequest(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "wrong-password"})
		if resp.Code != http.StatusUnauthorized {
			t.Fatalf("login failure code=%d body=%s", resp.Code, resp.Body.String())
		}
		access := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{"login_failure": true})
		if len(access) != 1 || access[0]["result"] != "failure" || access[0]["reason"] != "password_mismatch" || access[0]["remote_addr"] != "192.0.2.1" {
			t.Fatalf("unexpected login failure access log: %#v", access)
		}
		audit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"login_failure": true})
		if len(audit) != 1 || audit[0]["actor_type"] != "anonymous" || audit[0]["actor_id"] != nil || audit[0]["target_type"] != "auth" || audit[0]["target_id"] != "login" || audit[0]["result"] != "failure" {
			t.Fatalf("unexpected login failure audit log: %#v", audit)
		}
	})

	t.Run("memory login lock", func(t *testing.T) {
		state := newAPIState(t)
		policy := defaultAPIRateLimitPolicy()
		loginGroup := policy.Groups["login"]
		loginGroup.MaxRequests = 100
		policy.Groups["login"] = loginGroup
		writeTestJSON(t, filepath.Join(state, ".server_config"), apiServerConfig{APIRateLimit: apiRateLimitPolicyToMap(policy)})
		server := newTestAPI(t, state)
		for i := 0; i < 10; i++ {
			resp := apiRequestFrom(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "wrong-password"}, "198.51.100.7:1234")
			if resp.Code != http.StatusUnauthorized {
				t.Fatalf("login failure %d code=%d body=%s", i+1, resp.Code, resp.Body.String())
			}
		}
		resp := apiRequestFrom(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "password123"}, "198.51.100.7:1234")
		if resp.Code != http.StatusTooManyRequests || !strings.Contains(resp.Body.String(), "Too many attempts") {
			t.Fatalf("locked login code=%d body=%s", resp.Code, resp.Body.String())
		}
		var cred apiCredentials
		readTestJSON(t, filepath.Join(state, ".admin_credentials"), &cred)
		if cred.LoginCount != 0 || cred.LastLoginAt != nil {
			t.Fatalf("locked login must not verify/update credentials: %#v", cred)
		}
		access := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{"login_failure": true, "login_locked": true})
		if len(access) != 11 || actionsOf(access[9:]) != "login_failure,login_locked" || access[10]["reason"] != "too_many_attempts" {
			t.Fatalf("unexpected login lock access log: %#v", access)
		}
		audit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"login_failure": true, "permission_denied": true})
		if len(audit) != 11 || actionsOf(audit[9:]) != "login_failure,permission_denied" || audit[10]["target_type"] != "auth" || audit[10]["target_id"] != "login" {
			t.Fatalf("unexpected login lock audit log: %#v", audit)
		}
		server.cfg.Now = func() time.Time {
			return time.Date(2026, 9, 17, 1, 11, 2, 0, time.UTC)
		}
		resp = apiRequestFrom(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "password123"}, "198.51.100.7:1234")
		if resp.Code != http.StatusOK {
			t.Fatalf("login after lock expiry code=%d body=%s", resp.Code, resp.Body.String())
		}
		readTestJSON(t, filepath.Join(state, ".admin_credentials"), &cred)
		if cred.LoginCount != 1 || cred.LastLoginAt == nil {
			t.Fatalf("login after lock expiry must update credentials: %#v", cred)
		}
	})

	t.Run("auth transaction serializes logins", func(t *testing.T) {
		state := newAPIState(t)
		policy := defaultAPIRateLimitPolicy()
		loginGroup := policy.Groups["login"]
		loginGroup.MaxRequests = 100
		policy.Groups["login"] = loginGroup
		writeTestJSON(t, filepath.Join(state, ".server_config"), apiServerConfig{APIRateLimit: apiRateLimitPolicyToMap(policy)})
		server := newTestAPI(t, state)
		const total = 8
		var wg sync.WaitGroup
		errCh := make(chan string, total)
		for i := 0; i < total; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				resp := apiRequestFrom(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "password123"}, "198.51.100."+strconv.Itoa(20+i)+":1234")
				if resp.Code != http.StatusOK {
					errCh <- fmt.Sprintf("login %d code=%d body=%s", i, resp.Code, resp.Body.String())
				}
			}(i)
		}
		wg.Wait()
		close(errCh)
		for errText := range errCh {
			t.Fatal(errText)
		}
		var cred apiCredentials
		readTestJSON(t, filepath.Join(state, ".admin_credentials"), &cred)
		if cred.LoginCount != total || cred.LastLoginAt == nil {
			t.Fatalf("auth transaction must serialize login_count updates: %#v", cred)
		}
	})

	t.Run("auth transaction timeout", func(t *testing.T) {
		state := newAPIState(t)
		server := newTestAPI(t, state)
		server.cfg.AuthTransactionTimeout = time.Millisecond
		server.authTransactions <- struct{}{}
		defer func() { <-server.authTransactions }()
		req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader("{"))
		req.RemoteAddr = "192.0.2.1:1234"
		rec := httptest.NewRecorder()
		server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusConflict {
			t.Fatalf("auth transaction timeout code=%d body=%s", rec.Code, rec.Body.String())
		}
		if len(filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{"login_success": true, "login_failure": true, "login_locked": true})) != 0 {
			t.Fatalf("auth transaction timeout must not write access log")
		}
		if len(readJSONLines(filepath.Join(state, ".audit_log"))) != 0 {
			t.Fatalf("auth transaction timeout must not write audit log")
		}
		var cred apiCredentials
		readTestJSON(t, filepath.Join(state, ".admin_credentials"), &cred)
		if cred.LoginCount != 0 || cred.LastLoginAt != nil {
			t.Fatalf("auth transaction timeout must not update credentials: %#v", cred)
		}
	})

	t.Run("auth transaction covers authenticated session touch", func(t *testing.T) {
		state := newAPIState(t)
		server := newTestAPI(t, state)
		token := login(t, server, "password123")
		before, ok := apiSessionForToken(t, server, token)
		if !ok {
			t.Fatal("missing session before authenticated request")
		}
		server.cfg.Now = func() time.Time {
			return time.Date(2026, 9, 17, 1, 2, 3, 0, time.UTC)
		}
		resp := apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
		if resp.Code != http.StatusOK {
			t.Fatalf("authenticated status code=%d body=%s", resp.Code, resp.Body.String())
		}
		after, ok := apiSessionForToken(t, server, token)
		if !ok {
			t.Fatal("missing session after authenticated request")
		}
		if after.LastUsedAt != "2026-09-17T01:02:03Z" {
			t.Fatalf("session last_used_at was not touched: before=%#v after=%#v", before, after)
		}
		if after.ExpiresAt != before.ExpiresAt {
			t.Fatalf("authenticated request must not slide expires_at: before=%#v after=%#v", before, after)
		}
	})

	t.Run("auth transaction timeout for authenticated endpoint", func(t *testing.T) {
		state := newAPIState(t)
		server := newTestAPI(t, state)
		token := login(t, server, "password123")
		before, ok := apiSessionForToken(t, server, token)
		if !ok {
			t.Fatal("missing session before timeout")
		}
		server.cfg.AuthTransactionTimeout = time.Millisecond
		server.cfg.Now = func() time.Time {
			return time.Date(2026, 9, 17, 1, 2, 3, 0, time.UTC)
		}
		server.authTransactions <- struct{}{}
		defer func() { <-server.authTransactions }()
		resp := apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
		if resp.Code != http.StatusConflict {
			t.Fatalf("authenticated auth transaction timeout code=%d body=%s", resp.Code, resp.Body.String())
		}
		after, ok := apiSessionForToken(t, server, token)
		if !ok {
			t.Fatal("timeout must not delete session")
		}
		if after.LastUsedAt != before.LastUsedAt {
			t.Fatalf("timeout must not touch session: before=%#v after=%#v", before, after)
		}
	})

	t.Run("auth transaction deletes expired session before endpoint", func(t *testing.T) {
		state := newAPIState(t)
		server := newTestAPI(t, state)
		token := login(t, server, "password123")
		server.cfg.Now = func() time.Time {
			return time.Date(2026, 9, 17, 10, 1, 2, 0, time.UTC)
		}
		resp := apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
		if resp.Code != http.StatusUnauthorized {
			t.Fatalf("expired session code=%d body=%s", resp.Code, resp.Body.String())
		}
		if _, ok := apiSessionForToken(t, server, token); ok {
			t.Fatal("expired session was not deleted")
		}
	})

	t.Run("totp confirm failure", func(t *testing.T) {
		state := newAPIState(t)
		server := newTestAPI(t, state)
		token := login(t, server, "password123")
		resp := apiRequest(t, server, http.MethodPost, "/api/auth/totp-setup", token, nil)
		if resp.Code != http.StatusOK {
			t.Fatalf("totp setup code=%d body=%s", resp.Code, resp.Body.String())
		}
		var setup map[string]any
		decodeTestJSON(t, resp.Body.Bytes(), &setup)
		secret, _ := setup["secret"].(string)
		code := totpCode(secret, server.cfg.Now().UTC().Unix()/30)
		invalidCode := "000000"
		if invalidCode == code {
			invalidCode = "111111"
		}
		resp = apiRequest(t, server, http.MethodPost, "/api/auth/totp-confirm", token, map[string]string{"code": invalidCode})
		if resp.Code != http.StatusUnauthorized {
			t.Fatalf("totp confirm failure code=%d body=%s", resp.Code, resp.Body.String())
		}
		access := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{"totp_setup": true, "totp_failure": true})
		if actionsOf(access) != "totp_setup,totp_failure" || access[1]["reason"] != "invalid_totp" {
			t.Fatalf("unexpected totp confirm failure access log: %#v", access)
		}
		audit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"totp_setup": true, "totp_failure": true})
		if actionsOf(audit) != "totp_setup,totp_failure" || audit[1]["target_type"] != "auth" || audit[1]["target_id"] != "totp_confirm" || audit[1]["result"] != "failure" {
			t.Fatalf("unexpected totp confirm failure audit log: %#v", audit)
		}
		serialized := string(mustMarshalTest(t, []any{access, audit}))
		if strings.Contains(serialized, secret) || strings.Contains(serialized, setup["otpauth_uri"].(string)) {
			t.Fatalf("totp confirm failure logs leaked secret material: %s", serialized)
		}
	})

	t.Run("totp required", func(t *testing.T) {
		state := newAPIState(t)
		secret := "JBSWY3DPEHPK3PXP"
		if err := atomicWriteJSON(filepath.Join(state, ".totp_secret"), apiTOTPSecret{Enabled: true, SecretBase32: &secret}, 0600); err != nil {
			t.Fatal(err)
		}
		var before apiCredentials
		readTestJSON(t, filepath.Join(state, ".admin_credentials"), &before)
		server := newTestAPI(t, state)
		resp := apiRequest(t, server, http.MethodPost, "/api/login", "", map[string]string{"password": "password123"})
		if resp.Code != http.StatusOK {
			t.Fatalf("totp required code=%d body=%s", resp.Code, resp.Body.String())
		}
		var body map[string]any
		decodeTestJSON(t, resp.Body.Bytes(), &body)
		ticket, _ := body["ticket"].(string)
		if body["totp_required"] != true || ticket == "" {
			t.Fatalf("missing totp ticket: %#v", body)
		}
		var after apiCredentials
		readTestJSON(t, filepath.Join(state, ".admin_credentials"), &after)
		if after.LoginCount != before.LoginCount || after.LastLoginAt != before.LastLoginAt {
			t.Fatalf("totp-required password phase must not advance credentials: before=%#v after=%#v", before, after)
		}
		access := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".access_log")), map[string]bool{"totp_required": true})
		if len(access) != 1 || access[0]["result"] != "success" || access[0]["reason"] != nil {
			t.Fatalf("unexpected totp required access log: %#v", access)
		}
		audit := filterJSONLinesByActions(readJSONLines(filepath.Join(state, ".audit_log")), map[string]bool{"totp_required": true})
		if len(audit) != 1 || audit[0]["actor_type"] != "anonymous" || audit[0]["actor_id"] != nil || audit[0]["target_type"] != "auth" || audit[0]["target_id"] != "login" || audit[0]["result"] != "success" {
			t.Fatalf("unexpected totp required audit log: %#v", audit)
		}
		serialized := string(mustMarshalTest(t, []any{access, audit}))
		if strings.Contains(serialized, ticket) || strings.Contains(serialized, secret) {
			t.Fatalf("totp required logs leaked secret material: %s", serialized)
		}
	})
}

func TestAPIRateLimitEnforcement(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")

	resp := apiRequest(t, server, http.MethodPost, "/api/tokens", token, map[string]any{"label": "reader", "scopes": []string{"read"}})
	if resp.Code != http.StatusCreated {
		t.Fatalf("token create code=%d body=%s", resp.Code, resp.Body.String())
	}
	var tokenResp map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &tokenResp)
	apiToken, _ := tokenResp["token"].(string)
	if apiToken == "" {
		t.Fatalf("missing api token: %#v", tokenResp)
	}

	secret := "rate-webhook-secret"
	resp = apiRequest(t, server, http.MethodPost, "/api/webhook-config", token, map[string]string{"secret": secret})
	if resp.Code != http.StatusOK {
		t.Fatalf("webhook config code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodPost, "/api/api-rate-limit", token, testAPIRateLimitPolicy(map[string]apiRateLimitGroup{
		"read":    {WindowSeconds: 60, MaxRequests: 1},
		"trigger": {WindowSeconds: 60, MaxRequests: 1},
	}))
	if resp.Code != http.StatusOK {
		t.Fatalf("api-rate-limit config code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("first session status code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequest(t, server, http.MethodGet, "/api/status", token, nil)
	if resp.Code != http.StatusTooManyRequests || !strings.Contains(resp.Body.String(), `"Too many requests"`) {
		t.Fatalf("second session status code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = apiRequestFrom(t, server, http.MethodGet, "/api/status", apiToken, nil, "198.51.100.9:4321")
	if resp.Code != http.StatusOK {
		t.Fatalf("first token status code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = apiRequestFrom(t, server, http.MethodGet, "/api/status", apiToken, nil, "198.51.100.9:4321")
	if resp.Code != http.StatusTooManyRequests || !strings.Contains(resp.Body.String(), `"Too many requests"`) {
		t.Fatalf("second token status code=%d body=%s", resp.Code, resp.Body.String())
	}

	before := snapshotFiles(t, state, []string{".webhook_events.json", ".build_state"})
	resp = signedWebhookRequest(t, server, secret, []byte(`{`), "rate-delivery-1", "")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("first webhook code=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = signedWebhookRequest(t, server, secret, []byte(`{"ref":"refs/heads/main","after":"1111111111111111111111111111111111111111"}`), "rate-delivery-2", "")
	if resp.Code != http.StatusTooManyRequests || !strings.Contains(resp.Body.String(), `"Too many requests"`) {
		t.Fatalf("second webhook code=%d body=%s", resp.Code, resp.Body.String())
	}
	assertSnapshotEqual(t, before, snapshotFiles(t, state, []string{".webhook_events.json", ".build_state"}))

	var rateState apiRateLimitState
	readTestJSON(t, filepath.Join(state, ".api_rate_state"), &rateState)
	for _, key := range []string{"session:admin:read", "ip:192.0.2.1:read", "token:tok000001:read", "ip:198.51.100.9:read", "ip:192.0.2.1:trigger"} {
		window, ok := rateState.Windows[key]
		if !ok || window.Count != 1 {
			t.Fatalf("unexpected rate window %s: %#v", key, rateState.Windows)
		}
	}
	auditRecords := readJSONLines(filepath.Join(state, ".audit_log"))
	assertAuditRecord(t, auditRecords, "permission_denied", "admin", "admin", "denied")
	assertAuditRecord(t, auditRecords, "permission_denied", "api_token", "tok000001", "denied")
	assertAuditRecord(t, auditRecords, "permission_denied", "webhook", "webhook", "denied")

	resp = apiRequest(t, server, http.MethodGet, "/api/api-rate-limit", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("api-rate-limit get code=%d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if _, ok := body["state_summary"].([]any); !ok {
		t.Fatalf("state_summary must be an array: %#v", body)
	}
}

func TestAPIDiagnosticsContract(t *testing.T) {
	state := newAPIState(t)
	server := newTestAPI(t, state)
	token := login(t, server, "password123")
	if err := os.WriteFile(filepath.Join(state, ".github_token"), []byte("ghp_diagnostic\n"), 0600); err != nil {
		t.Fatal(err)
	}
	expires := "2026-10-01"
	writeTestJSON(t, filepath.Join(state, ".server_config"), apiServerConfig{PATExpiresAt: &expires})
	if err := os.WriteFile(filepath.Join(state, ".webhook_secret"), []byte("secret\n"), 0600); err != nil {
		t.Fatal(err)
	}
	siteDir := filepath.Join(state, "site")
	if err := os.MkdirAll(siteDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(siteDir, "index.html"), []byte(strings.Repeat("a", 2048)), 0600); err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(state, ".branch_config"), apiBranchConfigFile{BranchTargets: []apiBranchTarget{{Branch: "main", TargetFile: "docs", SHAFile: filepath.Join(state, ".last_sha"), Src: filepath.Join(state, "repo", "docs"), Out: siteDir}}})
	githubRateCalls := 0
	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rate_limit" {
			t.Fatalf("unexpected GitHub diagnostics request: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer ghp_diagnostic" {
			t.Fatalf("unexpected authorization header: %q", r.Header.Get("Authorization"))
		}
		githubRateCalls++
		reset := time.Date(2026, 9, 17, 2, 0, 0, 0, time.UTC).Unix()
		_ = json.NewEncoder(w).Encode(map[string]any{"resources": map[string]any{"core": map[string]any{"limit": 5000, "remaining": 4823, "reset": reset, "used": 177}}})
	}))
	defer github.Close()
	server.cfg.GitHubAPIBaseURL = github.URL
	server.cfg.HTTPClient = github.Client()
	systemdCalls := []string{}
	server.cfg.CommandRunner = func(ctx context.Context, name string, args ...string) (apiCommandResult, error) {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatalf("diagnostics systemd command must set a deadline")
		}
		call := name + " " + strings.Join(args, " ")
		systemdCalls = append(systemdCalls, call)
		switch call {
		case "systemctl is-active adlaire-ci.timer":
			return apiCommandResult{Stdout: "active\n", ExitCode: 0}, nil
		case "systemctl cat adlaire-ci.service":
			return apiCommandResult{Stdout: "[Service]\n", ExitCode: 0}, nil
		default:
			t.Fatalf("unexpected systemd command: %s", call)
		}
		return apiCommandResult{ExitCode: -1}, errors.New("unexpected command")
	}

	resp := apiRequest(t, server, http.MethodGet, "/api/diagnostics", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("diagnostics code=%d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]any
	decodeTestJSON(t, resp.Body.Bytes(), &body)
	if body["checked_at"] != "2026-09-17T01:01:01Z" {
		t.Fatalf("unexpected checked_at: %#v", body)
	}
	items := body["items"].([]any)
	expectedNames := []string{"pat", "github_api", "output_file", "systemd", "webhook"}
	if len(items) != len(expectedNames) {
		t.Fatalf("unexpected diagnostics items: %#v", items)
	}
	for i, expectedName := range expectedNames {
		item := items[i].(map[string]any)
		if item["name"] != expectedName || item["status"] != "ok" {
			t.Fatalf("unexpected diagnostics item %d: %#v", i, item)
		}
		if item["message"] == "" {
			t.Fatalf("diagnostics item must include message: %#v", item)
		}
	}
	if githubRateCalls != 1 {
		t.Fatalf("unexpected GitHub diagnostics calls: %d", githubRateCalls)
	}
	if len(systemdCalls) != 2 || systemdCalls[0] != "systemctl is-active adlaire-ci.timer" || systemdCalls[1] != "systemctl cat adlaire-ci.service" {
		t.Fatalf("unexpected systemd calls: %#v", systemdCalls)
	}
}

func newAPIState(t *testing.T) string {
	t.Helper()
	state := t.TempDir()
	now := time.Date(2026, 9, 17, 1, 0, 0, 0, time.UTC)
	if err := InitCredentials(state, "password123", now); err != nil {
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

func testAPIRateLimitPolicy(overrides map[string]apiRateLimitGroup) map[string]any {
	policy := defaultAPIRateLimitPolicy()
	for group, value := range overrides {
		policy.Groups[group] = value
	}
	return apiRateLimitPolicyToMap(policy)
}

func assertAuditRecord(t *testing.T, records []map[string]any, action, actorType, actorID, result string) {
	t.Helper()
	for _, record := range records {
		if record["action"] == action && record["actor_type"] == actorType && record["actor_id"] == actorID && record["result"] == result {
			return
		}
	}
	t.Fatalf("missing audit record action=%s actor_type=%s actor_id=%s result=%s in %#v", action, actorType, actorID, result, records)
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
	for _, name := range []string{".access_log", ".audit_log"} {
		if err := os.Remove(filepath.Join(server.cfg.StateDir, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
	}
	return token
}

func apiSessionForToken(t *testing.T, server *APIServer, token string) (apiSession, bool) {
	t.Helper()
	server.mu.Lock()
	defer server.mu.Unlock()
	session, ok := server.sessions[tokenHash(token)]
	return session, ok
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

func filterJSONLinesByActions(records []map[string]any, allowed map[string]bool) []map[string]any {
	out := []map[string]any{}
	for _, record := range records {
		action, _ := record["action"].(string)
		if allowed[action] {
			out = append(out, record)
		}
	}
	return out
}

func actionsOf(records []map[string]any) string {
	actions := []string{}
	for _, record := range records {
		action, _ := record["action"].(string)
		actions = append(actions, action)
	}
	return strings.Join(actions, ",")
}

func mustMarshalTest(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func cloneTestMap(t *testing.T, value map[string]any) map[string]any {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("clone marshal: %v", err)
	}
	var out map[string]any
	decodeTestJSON(t, data, &out)
	return out
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

type plainResponseWriter struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func (w *plainResponseWriter) Header() http.Header {
	return w.header
}

func (w *plainResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *plainResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.body.Write(data)
}

type flushingResponseWriter struct {
	plainResponseWriter
	flushes int
}

func (w *flushingResponseWriter) Flush() {
	w.flushes++
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

func anySlicesEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
