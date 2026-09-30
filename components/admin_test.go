package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runAdminForTest(args ...string) (int, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunAdmin(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestRunAdminHelpVersionPrecedence(t *testing.T) {
	code, stdout, stderr := runAdminForTest("--help", "--api-url", "bad\nurl", "unknown")
	if code != 0 || stdout != "Usage: adlaire-ci-admin --api-url url --token token [--json] command [command-args]\n" || stderr != "" {
		t.Fatalf("help precedence mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	code, stdout, stderr = runAdminForTest("--version", "--api-url", "relative", "--token")
	if code != 0 || !strings.HasPrefix(stdout, "adlaire-ci-admin V.0.0-dev go=") || stderr != "" {
		t.Fatalf("version precedence mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestRunAdminParseFailures(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		stderr string
	}{
		{name: "unsafe", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "bad\ncommand"}, stderr: "invalid command line token\n"},
		{name: "missing api url value", args: []string{"--api-url"}, stderr: "missing value: --api-url\n"},
		{name: "missing token value", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "--json"}, stderr: "missing value: --token\n"},
		{name: "unknown option", args: []string{"--unknown"}, stderr: "unknown option: --unknown\n"},
		{name: "equals option", args: []string{"--api-url=http://127.0.0.1:8765"}, stderr: "unknown option: --api-url=http://127.0.0.1:8765\n"},
		{name: "short option", args: []string{"-u"}, stderr: "unknown option: -u\n"},
		{name: "duplicate api url", args: []string{"--api-url", "http://127.0.0.1:8765", "--api-url", "http://127.0.0.1:8766", "--token", "tok", "status"}, stderr: "usage error\n"},
		{name: "duplicate token", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "--token", "tok2", "status"}, stderr: "usage error\n"},
		{name: "duplicate json", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "--json", "--json", "status"}, stderr: "usage error\n"},
		{name: "missing api url", args: []string{"--token", "tok", "status"}, stderr: "usage error\n"},
		{name: "missing token", args: []string{"--api-url", "http://127.0.0.1:8765", "status"}, stderr: "usage error\n"},
		{name: "missing command", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok"}, stderr: "usage error\n"},
		{name: "invalid api url", args: []string{"--api-url", "http://user@127.0.0.1:8765", "--token", "tok", "status"}, stderr: "usage error\n"},
		{name: "invalid api path", args: []string{"--api-url", "http://127.0.0.1:8765/a/../b", "--token", "tok", "status"}, stderr: "usage error\n"},
		{name: "invalid token", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "bad\ntok", "status"}, stderr: "invalid command line token\n"},
		{name: "unknown command", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "unknown"}, stderr: "unknown command: unknown\n"},
		{name: "cancel missing arg", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "cancel-queue"}, stderr: "usage error\n"},
		{name: "cancel invalid arg", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "cancel-queue", "a/b"}, stderr: "usage error\n"},
		{name: "cancel extra arg", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "cancel-queue", "q1", "extra"}, stderr: "usage error\n"},
		{name: "label invalid", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "config-snapshot", ""}, stderr: "usage error\n"},
		{name: "label extra", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "config-snapshot", "a", "b"}, stderr: "usage error\n"},
		{name: "no arg command extra", args: []string{"--api-url", "http://127.0.0.1:8765", "--token", "tok", "status", "extra"}, stderr: "usage error\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code, stdout, stderr := runAdminForTest(tc.args...)
			if code != 2 || stdout != "" || stderr != tc.stderr {
				t.Fatalf("unexpected parse result: code=%d stdout=%q stderr=%q", code, stdout, stderr)
			}
		})
	}
}

func TestRunAdminTransportAndHumanOutput(t *testing.T) {
	type requestRecord struct {
		method string
		path   string
		auth   string
		accept string
		agent  string
		ctype  string
		body   string
	}
	var records []requestRecord
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := new(bytes.Buffer)
		body.ReadFrom(r.Body)
		records = append(records, requestRecord{
			method: r.Method,
			path:   r.URL.EscapedPath(),
			auth:   r.Header.Get("Authorization"),
			accept: r.Header.Get("Accept"),
			agent:  r.Header.Get("User-Agent"),
			ctype:  r.Header.Get("Content-Type"),
			body:   body.String(),
		})
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		switch r.URL.EscapedPath() {
		case "/api/status", "/adlaire/api/status":
			fmt.Fprint(w, `{"last_build_status":"success","running":false}`)
		case "/api/queue":
			fmt.Fprint(w, `{"active":null,"queued":[{"id":"q1"},{"id":"q2"}]}`)
		case "/api/history":
			fmt.Fprint(w, `{"total":7,"history":[]}`)
		case "/api/build":
			fmt.Fprint(w, `{"queued":true,"queue_id":"q2026093001","dispatch":"requested"}`)
		case "/api/queue/q%201%252":
			fmt.Fprint(w, `{"message":"cancelled"}`)
		case "/api/config-snapshots":
			fmt.Fprint(w, `{"id":"snap001"}`)
		case "/api/events":
			fmt.Fprint(w, `{"total":5,"events":[{"id":"e1"}]}`)
		default:
			t.Errorf("unexpected path: %s", r.URL.EscapedPath())
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error":"not found"}`)
		}
	}))
	defer server.Close()

	cases := []struct {
		args   []string
		stdout string
	}{
		{args: []string{"--api-url", server.URL + "/", "--token", "tok", "status"}, stdout: "status=success running=false\n"},
		{args: []string{"--api-url", server.URL + "/adlaire/", "--token", "tok", "status"}, stdout: "status=success running=false\n"},
		{args: []string{"--api-url", server.URL, "--token", "tok", "queue"}, stdout: "active=none queued=2\n"},
		{args: []string{"--api-url", server.URL, "--token", "tok", "history"}, stdout: "total=7 latest=none\n"},
		{args: []string{"--api-url", server.URL, "--token", "tok", "trigger-build"}, stdout: "queued=q2026093001\n"},
		{args: []string{"--api-url", server.URL, "--token", "tok", "cancel-queue", "q 1%2"}, stdout: "queue cancelled\n"},
		{args: []string{"--api-url", server.URL, "--token", "tok", "config-snapshot"}, stdout: "snapshot=snap001\n"},
		{args: []string{"--api-url", server.URL, "--token", "tok", "config-snapshot", "nightly"}, stdout: "snapshot=snap001\n"},
		{args: []string{"--api-url", server.URL, "--token", "tok", "events"}, stdout: "events=5\n"},
	}

	for _, tc := range cases {
		code, stdout, stderr := runAdminForTest(tc.args...)
		if code != 0 || stdout != tc.stdout || stderr != "" {
			t.Fatalf("args %v: code=%d stdout=%q stderr=%q", tc.args, code, stdout, stderr)
		}
	}

	want := []struct {
		method string
		path   string
		ctype  string
		body   string
	}{
		{method: http.MethodGet, path: "/api/status"},
		{method: http.MethodGet, path: "/adlaire/api/status"},
		{method: http.MethodGet, path: "/api/queue"},
		{method: http.MethodGet, path: "/api/history"},
		{method: http.MethodPost, path: "/api/build", ctype: "application/json", body: "{}"},
		{method: http.MethodDelete, path: "/api/queue/q%201%252"},
		{method: http.MethodPost, path: "/api/config-snapshots", ctype: "application/json", body: `{"label":null}`},
		{method: http.MethodPost, path: "/api/config-snapshots", ctype: "application/json", body: `{"label":"nightly"}`},
		{method: http.MethodGet, path: "/api/events"},
	}
	if len(records) != len(want) {
		t.Fatalf("expected %d requests, got %d", len(want), len(records))
	}
	for i, expected := range want {
		got := records[i]
		if got.method != expected.method || got.path != expected.path || got.body != expected.body {
			t.Fatalf("request %d mismatch: got %+v want %+v", i, got, expected)
		}
		if got.auth != "Bearer tok" || got.accept != "application/json" || got.agent != "adlaire-ci-admin/V.0.0-dev" {
			t.Fatalf("request %d headers mismatch: %+v", i, got)
		}
		if got.ctype != expected.ctype {
			t.Fatalf("request %d content type mismatch: got %q want %q", i, got.ctype, expected.ctype)
		}
	}
}

func TestRunAdminJSONModePreservesWireBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, " \n {\"running\":false,\"last_build_status\":\"success\"} \t ")
	}))
	defer server.Close()

	code, stdout, stderr := runAdminForTest("--api-url", server.URL, "--token", "tok", "--json", "status")
	if code != 0 || stdout != "{\"running\":false,\"last_build_status\":\"success\"}\n" || stderr != "" {
		t.Fatalf("json mode mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestRunAdminOutputAndAPIErrors(t *testing.T) {
	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error":"secret-token"}`)
		}))
		defer server.Close()
		code, stdout, stderr := runAdminForTest("--api-url", server.URL, "--token", "secret-token", "status")
		if code != 1 || stdout != "" || stderr != "api error: 401\n" || strings.Contains(stderr, "secret-token") {
			t.Fatalf("http error mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})

	t.Run("redirect not followed", func(t *testing.T) {
		count := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count++
			http.Redirect(w, r, "/leak?token=secret-token", http.StatusFound)
		}))
		defer server.Close()
		code, stdout, stderr := runAdminForTest("--api-url", server.URL, "--token", "secret-token", "status")
		if code != 1 || stdout != "" || stderr != "api error: 302\n" || count != 1 || strings.Contains(stderr, "secret-token") {
			t.Fatalf("redirect mismatch: code=%d stdout=%q stderr=%q count=%d", code, stdout, stderr, count)
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, `{"last_build_status":"success","running":false}`)
		}))
		defer server.Close()
		code, stdout, stderr := runAdminForTest("--api-url", server.URL, "--token", "tok", "status")
		if code != 1 || stdout != "" || stderr != "api error: invalid response\n" {
			t.Fatalf("content-type mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"last_build_status":"success","running":false} {}`)
		}))
		defer server.Close()
		code, stdout, stderr := runAdminForTest("--api-url", server.URL, "--token", "tok", "status")
		if code != 1 || stdout != "" || stderr != "api error: invalid response\n" {
			t.Fatalf("invalid json mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})

	t.Run("over limit", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, strings.Repeat("x", adminMaxResponseBodyBytes+1))
		}))
		defer server.Close()
		code, stdout, stderr := runAdminForTest("--api-url", server.URL, "--token", "tok", "status")
		if code != 1 || stdout != "" || stderr != "api error: invalid response\n" {
			t.Fatalf("over-limit mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})

	t.Run("invalid shape", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"last_build_status":"success","running":"false"}`)
		}))
		defer server.Close()
		code, stdout, stderr := runAdminForTest("--api-url", server.URL, "--token", "tok", "status")
		if code != 1 || stdout != "" || stderr != "api error: invalid response\n" {
			t.Fatalf("invalid shape mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})
}

func TestAdminCLIFixturesExist(t *testing.T) {
	fixtures := map[string][]string{
		"partial-admin-cli-lifecycle": {
			"manifest.json",
			"input/cli.json",
			"expected/stdout.txt",
			"expected/stderr.txt",
			"expected/effects.json",
			"expected/security.json",
		},
		"success-admin-cli-transport": {
			"manifest.json",
			"input/cli.json",
			"input/fakes.json",
			"expected/request.json",
			"expected/stdout.txt",
			"expected/stderr.txt",
			"expected/effects.json",
			"expected/security.json",
		},
		"failure-admin-cli-output-errors": {
			"manifest.json",
			"input/cli.json",
			"input/fakes.json",
			"expected/response.json",
			"expected/stdout.txt",
			"expected/stderr.txt",
			"expected/effects.json",
			"expected/security.json",
		},
		"security-admin-cli-secret-redaction": {
			"manifest.json",
			"input/cli.json",
			"input/fakes.json",
			"expected/request.json",
			"expected/response.json",
			"expected/stdout.txt",
			"expected/stderr.txt",
			"expected/effects.json",
			"expected/security.json",
		},
	}

	for name, required := range fixtures {
		t.Run(name, func(t *testing.T) {
			root := filepath.Join("..", "testdata", "admin", "cli", name)
			for _, rel := range required {
				path := filepath.Join(root, rel)
				info, err := os.Stat(path)
				if err != nil {
					t.Fatalf("missing fixture file %s: %v", path, err)
				}
				if info.IsDir() {
					t.Fatalf("fixture path is a directory: %s", path)
				}
				if filepath.Ext(path) == ".json" {
					data, err := os.ReadFile(path)
					if err != nil {
						t.Fatalf("read fixture json %s: %v", path, err)
					}
					var value any
					if err := json.Unmarshal(data, &value); err != nil {
						t.Fatalf("invalid fixture json %s: %v", path, err)
					}
				}
			}
		})
	}
}
