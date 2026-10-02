package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunMCPHelpVersionAndValidation(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunMCP([]string{"--help", "--bad"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected help success, got %d", code)
	}
	if got := stdout.String(); got != "Usage: adlaire-ci-mcp --state-dir path [--addr host:port] [--read-only] [--client-token token] [--allow-non-loopback] [--version] [--help]\n" {
		t.Fatalf("unexpected help: %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := RunMCP([]string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected version success, got %d", code)
	}
	if got := stdout.String(); !strings.HasPrefix(got, "adlaire-ci-mcp V.0.0-dev go=") {
		t.Fatalf("unexpected version: %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := RunMCP([]string{"--state-dir", t.TempDir(), "--addr", "192.168.1.10:8766"}, &stdout, &stderr); code != 2 {
		t.Fatalf("expected non-loopback validation failure, got %d", code)
	}
	if got := stderr.String(); !strings.Contains(got, "non-loopback address is not allowed") {
		t.Fatalf("unexpected stderr: %q", got)
	}

	stdout.Reset()
	stderr.Reset()
	if code := RunMCP([]string{"--state-dir", t.TempDir(), "--state-dir", t.TempDir()}, &stdout, &stderr); code != 2 {
		t.Fatalf("expected duplicate option failure, got %d", code)
	}
	if got := stderr.String(); got != "duplicate option: --state-dir\n" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}

func TestMCPServerLifecycleToolsAndState(t *testing.T) {
	t.Parallel()

	stateDir := t.TempDir()
	now := time.Date(2026, 9, 30, 1, 2, 3, 0, time.UTC)
	server := newMCPServer(mcpConfig{StateDir: stateDir, Addr: "127.0.0.1:8766", Now: func() time.Time { return now }})
	ts := httptest.NewServer(server)
	defer ts.Close()

	resp := mcpHTTP(t, ts, http.MethodGet, "/health", "", "")
	if resp.StatusCode != http.StatusOK || resp.Body != `{"status":"ok","component":"mcp"}` {
		t.Fatalf("unexpected health: status=%d body=%s", resp.StatusCode, resp.Body)
	}

	initResp := mcpRPC(t, ts, map[string]any{
		"jsonrpc": "2.0",
		"id":      "init-1",
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": mcpProtocolVersion,
			"clientInfo":      map[string]any{"name": "codex", "version": "1.0"},
			"capabilities":    map[string]any{},
		},
	})
	if initResp["error"] != nil {
		t.Fatalf("unexpected initialize error: %#v", initResp["error"])
	}
	if !mcpTestFileExists(filepath.Join(stateDir, ".mcp_client_log")) {
		t.Fatal("expected client log")
	}

	toolsResp := mcpRPC(t, ts, map[string]any{"jsonrpc": "2.0", "id": "tools-1", "method": "tools/list", "params": map[string]any{}})
	result := toolsResp["result"].(map[string]any)
	tools := result["tools"].([]any)
	if len(tools) != len(mcpTools) {
		t.Fatalf("expected %d tools, got %d", len(mcpTools), len(tools))
	}
	first := tools[0].(map[string]any)
	if first["name"] != "adlaire.getStatus" {
		t.Fatalf("unexpected first tool: %#v", first)
	}

	confirmResp := mcpRPC(t, ts, map[string]any{
		"jsonrpc": "2.0",
		"id":      "build-1",
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "adlaire.triggerBuild",
			"arguments": map[string]any{"target": "docs", "source": "git:main", "options": map[string]any{}},
		},
	})
	errObj := confirmResp["error"].(map[string]any)
	if errObj["code"].(float64) != -32004 {
		t.Fatalf("expected confirmation required, got %#v", errObj)
	}
	confirmationID := errObj["data"].(map[string]any)["confirmation_id"].(string)
	successResp := mcpRPC(t, ts, map[string]any{
		"jsonrpc": "2.0",
		"id":      "build-2",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "adlaire.triggerBuild",
			"arguments": map[string]any{
				"target": "docs", "source": "git:main", "options": map[string]any{}, "confirmation_id": confirmationID,
			},
		},
	})
	if successResp["error"] != nil {
		t.Fatalf("unexpected build error: %#v", successResp["error"])
	}
	if !mcpTestFileExists(filepath.Join(stateDir, ".build_state")) {
		t.Fatal("expected build state")
	}
	if !mcpTestFileExists(filepath.Join(stateDir, ".mcp_audit_log")) {
		t.Fatal("expected audit log")
	}
	if !mcpTestFileExists(filepath.Join(stateDir, ".mcp_metrics")) {
		t.Fatal("expected metrics")
	}

	queueResp := mcpRPC(t, ts, map[string]any{
		"jsonrpc": "2.0",
		"id":      "queue-1",
		"method":  "resources/read",
		"params":  map[string]any{"uri": "adlaire://queue"},
	})
	if queueResp["error"] != nil {
		t.Fatalf("unexpected resource error: %#v", queueResp["error"])
	}
}

func TestMCPReadOnlyAndTokenScope(t *testing.T) {
	t.Parallel()

	stateDir := t.TempDir()
	tokenHash := mcpSHA256Hex([]byte("secret-token"))
	cfg := mcpStateConfig{ToolTimeoutMS: mcpDefaultToolMS, SamplingTimeoutMS: mcpDefaultSamplingMS, Scopes: []mcpScopeRecord{{
		TokenHash: tokenHash,
		Scopes:    []string{"read:status"},
		CreatedAt: "2026-09-30T00:00:00Z",
		UpdatedAt: "2026-09-30T00:00:00Z",
	}}}
	if err := atomicWriteJSON(filepath.Join(stateDir, ".mcp_config"), cfg, 0600); err != nil {
		t.Fatal(err)
	}
	server := newMCPServer(mcpConfig{StateDir: stateDir, ClientToken: "secret-token", ReadOnly: true, Now: func() time.Time {
		return time.Date(2026, 9, 30, 1, 2, 3, 0, time.UTC)
	}})
	ts := httptest.NewServer(server)
	defer ts.Close()

	unauthorized := mcpHTTP(t, ts, http.MethodPost, "/mcp", `{"jsonrpc":"2.0","id":"x","method":"initialize","params":{}}`, "")
	if unauthorized.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized, got %d", unauthorized.StatusCode)
	}

	mcpRPCWithToken(t, ts, "secret-token", map[string]any{
		"jsonrpc": "2.0",
		"id":      "init",
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": mcpProtocolVersion,
			"clientInfo":      map[string]any{"name": "codex", "version": "1.0"},
			"capabilities":    map[string]any{},
		},
	})
	toolsResp := mcpRPCWithToken(t, ts, "secret-token", map[string]any{"jsonrpc": "2.0", "id": "tools", "method": "tools/list", "params": map[string]any{}})
	for _, raw := range toolsResp["result"].(map[string]any)["tools"].([]any) {
		if raw.(map[string]any)["name"] == "adlaire.triggerBuild" {
			t.Fatal("read-only tools/list must exclude side-effect tools")
		}
	}
	forbidden := mcpRPCWithToken(t, ts, "secret-token", map[string]any{
		"jsonrpc": "2.0",
		"id":      "queue",
		"method":  "tools/call",
		"params":  map[string]any{"name": "adlaire.getQueue", "arguments": map[string]any{}},
	})
	if forbidden["error"].(map[string]any)["code"].(float64) != -32002 {
		t.Fatalf("expected scope forbidden, got %#v", forbidden)
	}

	_, err := server.toolSetMCPConfig(map[string]any{"scopes": []any{
		map[string]any{"token_hash": tokenHash, "scopes": []any{"read:status"}},
		map[string]any{"token_hash": tokenHash, "scopes": []any{"read:queue"}},
	}})
	if err == nil || !strings.Contains(err.Error(), "invalid_token_hash") {
		t.Fatalf("expected duplicate token hash rejection, got %v", err)
	}
}

func TestMCPSamplingAndTimeout(t *testing.T) {
	t.Parallel()

	stateDir := t.TempDir()
	logDir := filepath.Join(stateDir, ".build_logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatal(err)
	}
	buildID := "b20260930010203"
	if err := atomicWriteJSON(filepath.Join(logDir, buildID+".json"), map[string]any{
		"id":            buildID,
		"target_status": "failure",
		"lines":         []string{"build failed"},
	}, 0600); err != nil {
		t.Fatal(err)
	}
	server := newMCPServer(mcpConfig{StateDir: stateDir, Now: func() time.Time {
		return time.Date(2026, 9, 30, 1, 2, 3, 0, time.UTC)
	}})
	ts := httptest.NewServer(server)
	defer ts.Close()

	mcpRPC(t, ts, map[string]any{
		"jsonrpc": "2.0",
		"id":      "init",
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": mcpProtocolVersion,
			"clientInfo":      map[string]any{"name": "codex", "version": "1.0"},
			"capabilities":    map[string]any{},
		},
	})
	resp := mcpRPC(t, ts, map[string]any{
		"jsonrpc": "2.0",
		"id":      "sample",
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "adlaire.analyzeBuildError",
			"arguments": map[string]any{"build_id": buildID},
		},
	})
	if resp["error"] != nil {
		t.Fatalf("unexpected sampling error: %#v", resp["error"])
	}
	data := resp["result"].(map[string]any)["data"].(map[string]any)
	request := data["sampling_request"].(map[string]any)
	if request["method"] != "sampling/createMessage" || data["external_ai_api_called"] != false {
		t.Fatalf("unexpected sampling response: %#v", data)
	}
	if !mcpTestFileExists(filepath.Join(stateDir, ".mcp_audit_log")) {
		t.Fatal("expected sampling audit log")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	tool, ok := mcpFindTool("adlaire.getStatus")
	if !ok {
		t.Fatal("missing status tool")
	}
	_, status, rpcErr := server.executeToolWithTimeout(ctx, tool, map[string]any{})
	if status != "timeout" || rpcErr == nil || rpcErr.Code != -32003 {
		t.Fatalf("expected timeout result, got status=%s err=%#v", status, rpcErr)
	}
}

func TestMCPFixtureManifestFiles(t *testing.T) {
	t.Parallel()

	fixtures := map[string][]string{
		"cli/mcp-cli-lifecycle":                {"manifest.json", "expected/stdout.txt", "expected/stderr.txt", "expected/effects.json", "expected/security.json"},
		"jsonrpc/mcp-jsonrpc-errors":           {"manifest.json", "expected/response.json", "expected/state/state-diff.json", "expected/effects.json", "expected/security.json"},
		"jsonrpc/mcp-initialize-client-log":    {"manifest.json", "expected/response.json", "expected/state/state-diff.json", "expected/effects.json", "expected/security.json"},
		"tools/mcp-tools-list-call":            {"manifest.json", "expected/response.json", "expected/state/state-diff.json", "expected/effects.json", "expected/security.json"},
		"tools/mcp-tools-confirmation":         {"manifest.json", "expected/response.json", "expected/state/state-diff.json", "expected/effects.json", "expected/security.json"},
		"resources/mcp-resources-subscription": {"manifest.json", "expected/response.json", "expected/events.json", "expected/effects.json", "expected/security.json"},
		"prompts/mcp-prompts":                  {"manifest.json", "expected/response.json", "expected/effects.json", "expected/security.json"},
		"sampling/mcp-sampling":                {"manifest.json", "expected/request.json", "expected/response.json", "expected/state/state-diff.json", "expected/effects.json", "expected/security.json"},
		"sse/mcp-sse":                          {"manifest.json", "expected/events.json", "expected/effects.json", "expected/security.json"},
		"state/mcp-state-metrics-audit":        {"manifest.json", "expected/state/state-diff.json", "expected/effects.json", "expected/security.json"},
	}
	for name, required := range fixtures {
		root := filepath.Join("..", "..", "testdata", "mcp", name)
		for _, rel := range required {
			path := filepath.Join(root, rel)
			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("missing MCP fixture file %s: %v", path, err)
			}
			if info.IsDir() {
				t.Fatalf("fixture path is a directory: %s", path)
			}
			if strings.HasSuffix(rel, ".json") {
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
		for _, rel := range []string{"expected/logs"} {
			path := filepath.Join(root, rel)
			_, statErr := os.Stat(path)
			if statErr == nil {
				continue
			}
			if strings.Contains(name, "client-log") || strings.Contains(name, "confirmation") || strings.Contains(name, "state-metrics") {
				t.Fatalf("missing MCP fixture directory %s: %v", path, statErr)
			}
		}
	}
}

type mcpHTTPResponse struct {
	StatusCode int
	Body       string
}

func mcpRPC(t *testing.T, ts *httptest.Server, value map[string]any) map[string]any {
	t.Helper()
	return mcpRPCWithToken(t, ts, "", value)
}

func mcpRPCWithToken(t *testing.T, ts *httptest.Server, token string, value map[string]any) map[string]any {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	resp := mcpHTTP(t, ts, http.MethodPost, "/mcp", string(data), token)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected rpc status=%d body=%s", resp.StatusCode, resp.Body)
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(resp.Body), &out); err != nil {
		t.Fatalf("invalid rpc json: %v body=%s", err, resp.Body)
	}
	return out
}

func mcpHTTP(t *testing.T, ts *httptest.Server, method, path, body, token string) mcpHTTPResponse {
	t.Helper()
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return mcpHTTPResponse{StatusCode: resp.StatusCode, Body: string(data)}
}

func mcpTestFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
