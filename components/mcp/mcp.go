package mcp

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	mcpBinaryName        = "adlaire-ci-mcp"
	defaultMCPAddr       = "127.0.0.1:8766"
	mcpProtocolVersion   = "2025-06-18"
	mcpBodyLimit         = 1 << 20
	mcpKeepalivePeriod   = 30 * time.Second
	mcpConfirmationTTL   = 5 * time.Minute
	mcpDefaultToolMS     = 30000
	mcpDefaultSamplingMS = 60000
)

var mcpBinaryVersion = "V.0.0-dev"

func SetBinaryVersion(version string) {
	mcpBinaryVersion = version
}

type mcpConfig struct {
	StateDir         string
	Addr             string
	ReadOnly         bool
	ClientToken      string
	AllowNonLoopback bool
	Now              func() time.Time
}

type mcpServer struct {
	cfg           mcpConfig
	mu            sync.Mutex
	initialized   bool
	clientName    string
	confirmations map[string]mcpConfirmation
	subscribers   map[chan mcpSSEFrame]struct{}
}

type mcpConfirmation struct {
	Tool      string
	ParamsSHA string
	ExpiresAt time.Time
	Summary   string
}

type mcpSSEFrame struct {
	Event string
	Data  map[string]any
}

type mcpRPCResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	ID      any          `json:"id"`
	Result  any          `json:"result,omitempty"`
	Error   *mcpRPCError `json:"error,omitempty"`
}

type mcpRPCError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}

type mcpStateConfig struct {
	ToolTimeoutMS     int              `json:"tool_timeout_ms"`
	SamplingTimeoutMS int              `json:"sampling_timeout_ms"`
	Scopes            []mcpScopeRecord `json:"scopes"`
}

type mcpScopeRecord struct {
	TokenHash string   `json:"token_hash"`
	Scopes    []string `json:"scopes"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type mcpMetricsFile struct {
	Tools map[string]mcpToolMetric `json:"tools"`
}

type mcpToolMetric struct {
	SuccessCount   int    `json:"success_count"`
	ErrorCount     int    `json:"error_count"`
	TimeoutCount   int    `json:"timeout_count"`
	LastStatus     string `json:"last_status"`
	LastDurationMS int64  `json:"last_duration_ms"`
	LastAt         string `json:"last_at"`
}

type mcpToolSpec struct {
	Name        string
	Scope       string
	SideEffect  bool
	Description string
	Required    []string
	Properties  map[string]string
}

var mcpTools = []mcpToolSpec{
	{Name: "adlaire.getStatus", Scope: "read:status", Description: "Read current Adlaire CI status.", Properties: map[string]string{}},
	{Name: "adlaire.getQueue", Scope: "read:queue", Description: "Read current build queue.", Properties: map[string]string{}},
	{Name: "adlaire.triggerBuild", Scope: "write:build", SideEffect: true, Description: "Request a new build through the durable queue.", Required: []string{"target", "source"}, Properties: map[string]string{"target": "string", "source": "string", "options": "object"}},
	{Name: "adlaire.cancelQueueEntry", Scope: "write:queue", SideEffect: true, Description: "Cancel a waiting build queue entry.", Required: []string{"queue_id"}, Properties: map[string]string{"queue_id": "string"}},
	{Name: "adlaire.getHistory", Scope: "read:history", Description: "Read build history summary.", Properties: map[string]string{"limit": "integer", "offset": "integer"}},
	{Name: "adlaire.getBuildLog", Scope: "read:history", Description: "Read one saved build log.", Required: []string{"build_id"}, Properties: map[string]string{"build_id": "string"}},
	{Name: "adlaire.analyzeBuildError", Scope: "read:history", Description: "Request client-side sampling for a saved build error.", Required: []string{"build_id"}, Properties: map[string]string{"build_id": "string"}},
	{Name: "adlaire.getConfig", Scope: "read:config", Description: "Read effective Adlaire CI configuration.", Properties: map[string]string{}},
	{Name: "adlaire.setConfig", Scope: "write:config", SideEffect: true, Description: "Update one allowed configuration path.", Required: []string{"path", "value"}, Properties: map[string]string{"path": "string", "value": "value"}},
	{Name: "adlaire.getMcpConfig", Scope: "read:mcp", Description: "Read MCP configuration without token secrets.", Properties: map[string]string{}},
	{Name: "adlaire.setMcpConfig", Scope: "write:mcp", SideEffect: true, Description: "Update MCP timeout or scope configuration.", Properties: map[string]string{"tool_timeout_ms": "integer", "sampling_timeout_ms": "integer", "scopes": "array"}},
	{Name: "adlaire.createConfigSnapshot", Scope: "write:config", SideEffect: true, Description: "Create a configuration snapshot.", Properties: map[string]string{"label": "string-or-null"}},
	{Name: "adlaire.diffConfigSnapshots", Scope: "read:config", Description: "Diff two configuration snapshots with secrets masked.", Required: []string{"left_id", "right_id"}, Properties: map[string]string{"left_id": "string", "right_id": "string"}},
	{Name: "adlaire.restoreConfigSnapshot", Scope: "write:config", SideEffect: true, Description: "Restore one configuration snapshot.", Required: []string{"snapshot_id"}, Properties: map[string]string{"snapshot_id": "string"}},
	{Name: "adlaire.getMetrics", Scope: "read:metrics", Description: "Read MCP and API metrics snapshot.", Properties: map[string]string{}},
	{Name: "adlaire.getAuditLog", Scope: "read:audit", Description: "Read MCP audit log entries.", Properties: map[string]string{"limit": "integer", "offset": "integer"}},
	{Name: "adlaire.resendWebhook", Scope: "write:notification", SideEffect: true, Description: "Request webhook delivery resend.", Required: []string{"delivery_id"}, Properties: map[string]string{"delivery_id": "string"}},
}

func RunMCP(args []string, stdout, stderr io.Writer) int {
	cfg, handled, err := parseMCPArgs(args, os.Stdin, stdout)
	if err != nil {
		var ee exitError
		if errors.As(err, &ee) {
			fmt.Fprintln(stderr, ee.Msg)
			return ee.Code
		}
		fmt.Fprintln(stderr, err)
		return 1
	}
	if handled {
		return 0
	}
	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		fmt.Fprintln(stderr, "listen failed")
		return 1
	}
	defer listener.Close()
	server := &http.Server{Handler: newMCPServer(cfg)}
	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintln(stderr, "listen failed")
		return 1
	}
	return 0
}

func parseMCPArgs(args []string, stdin io.Reader, stdout io.Writer) (mcpConfig, bool, error) {
	cfg := mcpConfig{Addr: defaultMCPAddr, Now: func() time.Time { return time.Now().UTC() }}
	if hasExactArg(args, "--help") {
		fmt.Fprintln(stdout, "Usage: adlaire-ci-mcp --state-dir path [--addr host:port] [--read-only] [--client-token-file path | --client-token-stdin] [--allow-non-loopback] [--version] [--help]")
		return cfg, true, nil
	}
	if hasExactArg(args, "--version") {
		fmt.Fprintf(stdout, "%s %s go=%s\n", mcpBinaryName, mcpBinaryVersion, runtime.Version())
		return cfg, true, nil
	}
	if !safeArgvTokens(args) {
		return cfg, false, exitError{Code: 2, Msg: "invalid command line token"}
	}
	stateDirProvided := false
	seen := map[string]bool{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
		if !strings.HasPrefix(arg, "--") {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
		if strings.Contains(arg, "=") {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
		switch arg {
		case "--state-dir", "--addr", "--client-token-file":
			if seen[arg] {
				return cfg, false, exitError{Code: 2, Msg: "duplicate option: " + arg}
			}
			seen[arg] = true
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
				return cfg, false, exitError{Code: 2, Msg: "missing value: " + arg}
			}
			i++
			switch arg {
			case "--state-dir":
				cfg.StateDir = args[i]
				stateDirProvided = true
			case "--addr":
				cfg.Addr = args[i]
			case "--client-token-file":
				if seen["--client-token-stdin"] {
					return cfg, false, exitError{Code: 2, Msg: "duplicate option: client token source"}
				}
				token, err := readMCPTokenFile(args[i])
				if err != nil {
					return cfg, false, exitError{Code: 2, Msg: "invalid client token source"}
				}
				cfg.ClientToken = token
			}
		case "--client-token-stdin":
			if seen[arg] || seen["--client-token-file"] {
				return cfg, false, exitError{Code: 2, Msg: "duplicate option: client token source"}
			}
			seen[arg] = true
			token, err := readMCPTokenReader(stdin)
			if err != nil {
				return cfg, false, exitError{Code: 2, Msg: "invalid client token source"}
			}
			cfg.ClientToken = token
		case "--read-only":
			if seen[arg] {
				return cfg, false, exitError{Code: 2, Msg: "duplicate option: " + arg}
			}
			seen[arg] = true
			cfg.ReadOnly = true
		case "--allow-non-loopback":
			if seen[arg] {
				return cfg, false, exitError{Code: 2, Msg: "duplicate option: " + arg}
			}
			seen[arg] = true
			cfg.AllowNonLoopback = true
		default:
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
	}
	if !stateDirProvided {
		return cfg, false, exitError{Code: 2, Msg: "state directory is required"}
	}
	if err := validateAPIStateDir(cfg.StateDir); err != nil {
		return cfg, false, err
	}
	if err := validateMCPListenAddress(cfg.Addr, cfg.AllowNonLoopback); err != nil {
		return cfg, false, err
	}
	if cfg.ClientToken == "" && (seen["--client-token-file"] || seen["--client-token-stdin"]) {
		return cfg, false, exitError{Code: 2, Msg: "client token must not be empty"}
	}
	return cfg, false, nil
}

func validateMCPListenAddress(addr string, allowNonLoopback bool) error {
	host, portText, err := net.SplitHostPort(addr)
	if err != nil || host == "" || portText == "" || (len(portText) > 1 && portText[0] == '0') {
		return exitError{Code: 2, Msg: "invalid listen address: " + addr}
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return exitError{Code: 2, Msg: "invalid listen address: " + addr}
	}
	if host == "localhost" {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil || strings.Contains(host, ":") {
		return exitError{Code: 2, Msg: "invalid listen address: " + addr}
	}
	ip4 := ip.To4()
	if ip4 == nil || host != ip4.String() {
		return exitError{Code: 2, Msg: "invalid listen address: " + addr}
	}
	if ip4[0] == 0 || ip4.Equal(net.IPv4(255, 255, 255, 255)) || ip4[0] >= 224 {
		return exitError{Code: 2, Msg: "invalid listen address: " + addr}
	}
	if ip4[0] == 127 {
		return nil
	}
	if !allowNonLoopback {
		return exitError{Code: 2, Msg: "non-loopback address is not allowed: " + addr}
	}
	return nil
}

func newMCPServer(cfg mcpConfig) *mcpServer {
	if cfg.Addr == "" {
		cfg.Addr = defaultMCPAddr
	}
	if cfg.Now == nil {
		cfg.Now = func() time.Time { return time.Now().UTC() }
	}
	return &mcpServer{
		cfg:           cfg,
		confirmations: map[string]mcpConfirmation{},
		subscribers:   map[chan mcpSSEFrame]struct{}{},
	}
}

func (s *mcpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/health":
		if r.Method != http.MethodGet {
			mcpWriteHTTPError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		mcpWriteJSON(w, http.StatusOK, struct {
			Status    string `json:"status"`
			Component string `json:"component"`
		}{Status: "ok", Component: "mcp"})
	case "/mcp":
		if r.Method != http.MethodPost {
			mcpWriteHTTPError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		s.handleRPC(w, r)
	case "/mcp/events":
		if r.Method != http.MethodGet {
			mcpWriteHTTPError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		s.handleEvents(w, r)
	default:
		mcpWriteHTTPError(w, http.StatusNotFound, "Not found")
	}
}

func (s *mcpServer) handleRPC(w http.ResponseWriter, r *http.Request) {
	if !mcpValidContentType(r.Header.Values("Content-Type")) {
		mcpWriteHTTPError(w, http.StatusBadRequest, "Invalid content type")
		return
	}
	if !s.authorizeHTTP(w, r) {
		return
	}
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, mcpBodyLimit))
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			mcpWriteHTTPError(w, http.StatusRequestEntityTooLarge, "Payload too large")
			return
		}
		s.writeRPCError(w, http.StatusBadRequest, nil, -32700, "Parse error", "invalid_json")
		return
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		s.writeRPCError(w, http.StatusBadRequest, nil, -32700, "Parse error", "empty_body")
		return
	}
	if data[0] == '[' {
		s.writeRPCError(w, http.StatusOK, nil, -32600, "Invalid Request", "batch_not_allowed")
		return
	}
	raw, order, ok := mcpDecodeJSONObject(data)
	if !ok {
		s.writeRPCError(w, http.StatusBadRequest, nil, -32700, "Parse error", "invalid_json")
		return
	}
	id, idOK, idPresent := mcpParseID(raw["id"])
	if !mcpRPCKeysAllowed(order) || !idOK {
		s.writeRPCError(w, http.StatusOK, id, -32600, "Invalid Request", "invalid_request")
		return
	}
	var version string
	if err := json.Unmarshal(raw["jsonrpc"], &version); err != nil || version != "2.0" {
		s.writeRPCError(w, http.StatusOK, id, -32600, "Invalid Request", "invalid_jsonrpc")
		return
	}
	var method string
	if err := json.Unmarshal(raw["method"], &method); err != nil || !mcpValidTinyString(method, 128) {
		s.writeRPCError(w, http.StatusOK, id, -32600, "Invalid Request", "invalid_method")
		return
	}
	paramsRaw, paramsOK := raw["params"]
	if !paramsOK {
		paramsRaw = []byte(`{}`)
	}
	params, paramsValid := mcpDecodeParamsObject(paramsRaw)
	if !paramsValid {
		s.writeRPCError(w, http.StatusOK, id, -32602, "Invalid params", "invalid_params")
		return
	}
	if method == "notifications/initialized" {
		if idPresent {
			s.writeRPCError(w, http.StatusOK, id, -32600, "Invalid Request", "notification_has_id")
			return
		}
		s.mu.Lock()
		s.initialized = true
		s.mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if !idPresent {
		s.writeRPCError(w, http.StatusOK, nil, -32600, "Invalid Request", "missing_id")
		return
	}
	if method != "initialize" && !s.isInitialized() {
		s.writeRPCError(w, http.StatusOK, id, -32600, "Invalid Request", "not_initialized")
		return
	}
	result, rpcErr := s.dispatchRPC(r.Context(), r, method, params)
	if rpcErr != nil {
		mcpWriteJSON(w, http.StatusOK, mcpRPCResponse{JSONRPC: "2.0", ID: id, Error: rpcErr})
		return
	}
	mcpWriteJSON(w, http.StatusOK, mcpRPCResponse{JSONRPC: "2.0", ID: id, Result: result})
}

func (s *mcpServer) handleEvents(w http.ResponseWriter, r *http.Request) {
	if !s.authorizeHTTP(w, r) {
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		mcpWriteHTTPError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	ch := make(chan mcpSSEFrame, 16)
	s.mu.Lock()
	s.subscribers[ch] = struct{}{}
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.subscribers, ch)
		s.mu.Unlock()
	}()
	ticker := time.NewTicker(mcpKeepalivePeriod)
	defer ticker.Stop()
	for {
		select {
		case frame := <-ch:
			if !mcpWriteSSEFrame(w, frame) {
				return
			}
			flusher.Flush()
		case <-ticker.C:
			if _, err := io.WriteString(w, ": keepalive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (s *mcpServer) authorizeHTTP(w http.ResponseWriter, r *http.Request) bool {
	if s.cfg.ClientToken == "" {
		return true
	}
	values := r.Header.Values("Authorization")
	if len(values) != 1 || values[0] != "Bearer "+s.cfg.ClientToken {
		mcpWriteHTTPError(w, http.StatusUnauthorized, "Unauthorized")
		return false
	}
	return true
}

func (s *mcpServer) dispatchRPC(ctx context.Context, r *http.Request, method string, params map[string]any) (any, *mcpRPCError) {
	switch method {
	case "initialize":
		return s.rpcInitialize(r, params)
	case "tools/list":
		return s.rpcToolsList(), nil
	case "tools/call":
		return s.rpcToolsCall(ctx, params)
	case "resources/list":
		return s.rpcResourcesList(), nil
	case "resources/read":
		return s.rpcResourcesRead(params)
	case "resources/subscribe":
		return s.rpcResourceSubscribe(params)
	case "resources/unsubscribe":
		return s.rpcResourceUnsubscribe(params)
	case "prompts/list":
		return s.rpcPromptsList(), nil
	case "prompts/get":
		return s.rpcPromptsGet(params)
	default:
		return nil, mcpError(-32601, "Method not found", "method_not_found")
	}
}

func (s *mcpServer) rpcInitialize(r *http.Request, params map[string]any) (any, *mcpRPCError) {
	protocol, _ := params["protocolVersion"].(string)
	if protocol != mcpProtocolVersion {
		return nil, mcpError(-32602, "Invalid params", "invalid_protocol")
	}
	clientInfo, ok := params["clientInfo"].(map[string]any)
	if !ok {
		return nil, mcpError(-32602, "Invalid params", "invalid_client")
	}
	name, _ := clientInfo["name"].(string)
	version, _ := clientInfo["version"].(string)
	if !mcpValidTinyString(name, 128) || !mcpValidTinyString(version, 128) {
		return nil, mcpError(-32602, "Invalid params", "invalid_client")
	}
	capabilities, ok := params["capabilities"].(map[string]any)
	if !ok {
		return nil, mcpError(-32602, "Invalid params", "invalid_capabilities")
	}
	now := s.now()
	rec := map[string]any{
		"id":               s.nextLineID(filepath.Join(s.cfg.StateDir, ".mcp_client_log"), "mcpcli_", now),
		"connected_at":     mcpTime(now),
		"client_name":      name,
		"client_version":   version,
		"protocol_version": protocol,
		"remote_addr":      mcpRemoteAddr(r.RemoteAddr),
		"capabilities":     maskSecrets(capabilities),
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".mcp_client_log"), rec); err != nil {
		return nil, mcpError(-32603, "Internal error", "client_log_failed")
	}
	s.mu.Lock()
	s.clientName = name
	s.initialized = true
	s.mu.Unlock()
	return map[string]any{
		"protocolVersion": mcpProtocolVersion,
		"serverInfo": map[string]any{
			"name":    mcpBinaryName,
			"version": mcpBinaryVersion,
		},
		"capabilities": map[string]any{
			"tools":       map[string]any{},
			"resources":   map[string]any{"subscribe": true},
			"prompts":     map[string]any{},
			"logging":     map[string]any{},
			"sampling":    map[string]any{},
			"elicitation": map[string]any{},
		},
	}, nil
}

func (s *mcpServer) rpcToolsList() any {
	tools := []map[string]any{}
	for _, tool := range mcpTools {
		if s.cfg.ReadOnly && tool.SideEffect {
			continue
		}
		tools = append(tools, map[string]any{
			"name":        tool.Name,
			"description": tool.Description,
			"inputSchema": mcpToolInputSchema(tool),
			"annotations": map[string]any{"readOnlyHint": !tool.SideEffect},
		})
	}
	return map[string]any{"tools": tools}
}

func (s *mcpServer) rpcToolsCall(ctx context.Context, params map[string]any) (any, *mcpRPCError) {
	name, _ := params["name"].(string)
	args, ok := params["arguments"].(map[string]any)
	if !ok || !mcpValidTinyString(name, 128) || len(params) != 2 {
		return nil, mcpError(-32602, "Invalid params", "invalid_tool_call")
	}
	tool, ok := mcpFindTool(name)
	if !ok {
		return nil, mcpError(-32601, "Method not found", "tool_not_found")
	}
	if err := mcpValidateToolArgs(tool, args); err != nil {
		return nil, mcpError(-32602, "Invalid params", err.Error())
	}
	if s.cfg.ReadOnly && tool.SideEffect {
		return nil, mcpError(-32002, "Forbidden", "read_only")
	}
	if !s.scopeAllowed(tool.Scope) {
		return nil, mcpError(-32002, "Forbidden", "scope_forbidden")
	}
	var confirmationID any
	if tool.SideEffect {
		id, ok := args["confirmation_id"].(string)
		if !ok || id == "" {
			return nil, s.createConfirmation(tool, args)
		}
		confirmationID = id
		if err := s.verifyConfirmation(tool, args, id); err != nil {
			return nil, mcpError(-32004, "Confirmation required", err.Error())
		}
		if err := s.appendAudit(tool, args, id, "accepted", 0, nil, nil, nil); err != nil {
			return nil, mcpError(-32603, "Internal error", "audit_failed")
		}
	}
	start := s.now()
	result, status, rpcErr := s.executeToolWithTimeout(ctx, tool, args)
	duration := s.now().Sub(start).Milliseconds()
	if name == "adlaire.analyzeBuildError" {
		buildID, _ := args["build_id"].(string)
		promptText := "Analyze Adlaire CI build log " + buildID + " and return a concise failure category, likely cause, and next action. Do not call external services."
		promptHash := mcpSHA256Hex([]byte(promptText))
		auditStatus := status
		code := any(nil)
		if rpcErr != nil {
			code = rpcErr.Code
		}
		if err := s.appendAudit(tool, args, "", auditStatus, duration, code, &buildID, &promptHash); err != nil {
			return nil, mcpError(-32603, "Internal error", "audit_failed")
		}
	}
	if err := s.updateMetrics(tool.Name, status, duration); err != nil {
		return nil, mcpError(-32603, "Internal error", "metrics_failed")
	}
	if rpcErr != nil {
		return nil, rpcErr
	}
	if confirmationID != nil {
		s.publish(mcpSSEFrame{Event: "resource-updated", Data: map[string]any{"uri": "adlaire://queue", "updated_at": mcpTime(s.now())}})
	}
	return map[string]any{"ok": true, "data": result, "warnings": []any{}}, nil
}

func (s *mcpServer) executeToolWithTimeout(ctx context.Context, tool mcpToolSpec, args map[string]any) (map[string]any, string, *mcpRPCError) {
	cfg, err := s.readMCPConfig()
	if err != nil {
		return nil, "error", mcpError(-32603, "Internal error", "mcp_config_failed")
	}
	timeout := time.Duration(cfg.ToolTimeoutMS) * time.Millisecond
	if timeout <= 0 {
		timeout = time.Duration(mcpDefaultToolMS) * time.Millisecond
	}
	toolCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := toolCtx.Err(); err != nil {
		return nil, "timeout", mcpError(-32003, "Timeout", "tool_timeout")
	}
	type toolResult struct {
		data   map[string]any
		status string
		err    *mcpRPCError
	}
	done := make(chan toolResult, 1)
	go func() {
		data, status, rpcErr := s.executeTool(toolCtx, tool, args)
		done <- toolResult{data: data, status: status, err: rpcErr}
	}()
	select {
	case result := <-done:
		if toolCtx.Err() != nil {
			return nil, "timeout", mcpError(-32003, "Timeout", "tool_timeout")
		}
		return result.data, result.status, result.err
	case <-toolCtx.Done():
		return nil, "timeout", mcpError(-32003, "Timeout", "tool_timeout")
	}
}

func (s *mcpServer) executeTool(ctx context.Context, tool mcpToolSpec, args map[string]any) (map[string]any, string, *mcpRPCError) {
	if err := ctx.Err(); err != nil {
		return nil, "timeout", mcpError(-32003, "Timeout", "tool_timeout")
	}
	switch tool.Name {
	case "adlaire.getStatus":
		return s.toolStatus(), "success", nil
	case "adlaire.getQueue":
		return s.toolQueue(), "success", nil
	case "adlaire.triggerBuild":
		data, err := s.toolTriggerBuild(args)
		if err != nil {
			return nil, "error", mcpError(-32603, "Internal error", "queue_failed")
		}
		return data, "success", nil
	case "adlaire.cancelQueueEntry":
		data, err := s.toolCancelQueue(args)
		if err != nil {
			return nil, "error", mcpError(-32602, "Invalid params", err.Error())
		}
		return data, "success", nil
	case "adlaire.getHistory":
		return s.toolHistory(args), "success", nil
	case "adlaire.getBuildLog":
		data, err := s.toolBuildLog(args)
		if err != nil {
			return nil, "error", mcpError(-32602, "Invalid params", "build_log_not_found")
		}
		return data, "success", nil
	case "adlaire.analyzeBuildError":
		data, err := s.toolAnalyzeBuildError(args)
		if err != nil {
			return nil, "error", mcpError(-32602, "Invalid params", "build_log_not_found")
		}
		return data, "success", nil
	case "adlaire.getConfig":
		return s.toolConfig(), "success", nil
	case "adlaire.setConfig":
		data, err := s.toolSetConfig(args)
		if err != nil {
			return nil, "error", mcpError(-32602, "Invalid params", err.Error())
		}
		return data, "success", nil
	case "adlaire.getMcpConfig":
		cfg, err := s.readMCPConfig()
		if err != nil {
			return nil, "error", mcpError(-32603, "Internal error", "mcp_config_failed")
		}
		return map[string]any{"config": mcpConfigPublic(cfg)}, "success", nil
	case "adlaire.setMcpConfig":
		data, err := s.toolSetMCPConfig(args)
		if err != nil {
			return nil, "error", mcpError(-32602, "Invalid params", err.Error())
		}
		return data, "success", nil
	case "adlaire.createConfigSnapshot":
		data, err := s.toolCreateConfigSnapshot(args)
		if err != nil {
			return nil, "error", mcpError(-32603, "Internal error", "snapshot_failed")
		}
		return data, "success", nil
	case "adlaire.diffConfigSnapshots":
		data, err := s.toolDiffConfigSnapshots(args)
		if err != nil {
			return nil, "error", mcpError(-32602, "Invalid params", err.Error())
		}
		return data, "success", nil
	case "adlaire.restoreConfigSnapshot":
		data, err := s.toolRestoreConfigSnapshot(args)
		if err != nil {
			return nil, "error", mcpError(-32602, "Invalid params", err.Error())
		}
		return data, "success", nil
	case "adlaire.getMetrics":
		return s.toolMetrics(), "success", nil
	case "adlaire.getAuditLog":
		return s.toolAuditLog(args), "success", nil
	case "adlaire.resendWebhook":
		data, err := s.toolResendWebhook(args)
		if err != nil {
			return nil, "error", mcpError(-32603, "Internal error", "webhook_resend_failed")
		}
		return data, "success", nil
	default:
		return nil, "error", mcpError(-32601, "Method not found", "tool_not_found")
	}
}

func (s *mcpServer) rpcResourcesList() any {
	resources := []map[string]any{}
	for _, item := range mcpResourceSpecs() {
		resources = append(resources, map[string]any{
			"uri":         item[0],
			"name":        item[1],
			"description": item[2],
			"mimeType":    item[3],
		})
	}
	return map[string]any{"resources": resources}
}

func (s *mcpServer) rpcResourcesRead(params map[string]any) (any, *mcpRPCError) {
	uri, _ := params["uri"].(string)
	if !mcpKnownResource(uri) {
		return nil, mcpError(-32602, "Invalid params", "unknown_resource")
	}
	text, mime, err := s.resourceText(uri)
	if err != nil {
		return nil, mcpError(-32603, "Internal error", "resource_read_failed")
	}
	return map[string]any{"contents": []map[string]any{{"uri": uri, "mimeType": mime, "text": text}}}, nil
}

func (s *mcpServer) rpcResourceSubscribe(params map[string]any) (any, *mcpRPCError) {
	uri, _ := params["uri"].(string)
	if !mcpKnownResource(uri) {
		return nil, mcpError(-32602, "Invalid params", "unknown_resource")
	}
	return map[string]any{}, nil
}

func (s *mcpServer) rpcResourceUnsubscribe(params map[string]any) (any, *mcpRPCError) {
	uri, _ := params["uri"].(string)
	if !mcpKnownResource(uri) {
		return nil, mcpError(-32602, "Invalid params", "unknown_resource")
	}
	return map[string]any{}, nil
}

func (s *mcpServer) rpcPromptsList() any {
	return map[string]any{"prompts": []map[string]any{
		{"name": "adlaire.buildFailureTriage", "description": "Generate build failure triage steps.", "arguments": []map[string]any{{"name": "build_id", "required": true}}},
		{"name": "adlaire.releaseReadiness", "description": "Generate release readiness checks.", "arguments": []map[string]any{{"name": "target", "required": true}}},
		{"name": "adlaire.configReview", "description": "Generate configuration snapshot review steps.", "arguments": []map[string]any{{"name": "snapshot_id", "required": true}}},
	}}
}

func (s *mcpServer) rpcPromptsGet(params map[string]any) (any, *mcpRPCError) {
	name, _ := params["name"].(string)
	args, ok := params["arguments"].(map[string]any)
	if !ok {
		return nil, mcpError(-32602, "Invalid params", "invalid_prompt_arguments")
	}
	var text string
	switch name {
	case "adlaire.buildFailureTriage":
		buildID, _ := args["build_id"].(string)
		if !mcpValidID(buildID) {
			return nil, mcpError(-32602, "Invalid params", "invalid_build_id")
		}
		text = "Review build " + buildID + " logs, status, history, and configuration changes. Do not mutate queue, config, or snapshots."
	case "adlaire.releaseReadiness":
		target, _ := args["target"].(string)
		if !mcpValidTarget(target) {
			return nil, mcpError(-32602, "Invalid params", "invalid_target")
		}
		text = "Review release readiness for target " + target + " using status, queue, history, and release evidence. Do not publish or mutate state."
	case "adlaire.configReview":
		snapshotID, _ := args["snapshot_id"].(string)
		if !strings.HasPrefix(snapshotID, "cfgsnap_") {
			return nil, mcpError(-32602, "Invalid params", "invalid_snapshot_id")
		}
		text = "Review configuration snapshot " + snapshotID + " and compare masked differences. Do not restore or change config."
	default:
		return nil, mcpError(-32602, "Invalid params", "unknown_prompt")
	}
	return map[string]any{"messages": []map[string]any{{"role": "user", "content": map[string]any{"type": "text", "text": text}}}}, nil
}

func (s *mcpServer) toolStatus() map[string]any {
	state, stateErr := readBuildState(filepath.Join(s.cfg.StateDir, ".build_state"))
	status, hasStatus, statusErr := readBuildStatus(filepath.Join(s.cfg.StateDir, ".build_status"))
	records := readHistory(filepath.Join(s.cfg.StateDir, ".build_history"))
	out := map[string]any{"state": "ok", "running": state.Running, "current_build_id": state.CurrentBuildID, "queued_count": len(state.Queued)}
	if hasStatus && statusErr == nil {
		out["build_status"] = status
	}
	if len(records) > 0 {
		out["latest_history"] = decorateHistoryRecords(records[:1])[0]
	}
	if stateErr != nil || statusErr != nil {
		out["state"] = "degraded"
	}
	return out
}

func (s *mcpServer) toolQueue() map[string]any {
	state, err := readBuildState(filepath.Join(s.cfg.StateDir, ".build_state"))
	if err != nil {
		return map[string]any{"active": nil, "queued": []any{}, "corrupted": true}
	}
	return map[string]any{"active": state.ActiveQueueEntry, "queued": sortedQueueEntries(state.Queued), "running": state.Running}
}

func (s *mcpServer) toolTriggerBuild(args map[string]any) (map[string]any, error) {
	target, _ := args["target"].(string)
	source, _ := args["source"].(string)
	id, duplicate, err := s.enqueueMCPBuildRequest(filepath.Join(s.cfg.StateDir, ".build_state"), map[string]any{"target": target, "source": source})
	if err != nil {
		return nil, err
	}
	dispatch := "timer_fallback"
	if duplicate {
		dispatch = "duplicate"
	}
	return map[string]any{"queue_id": id, "dispatch": dispatch, "duplicate": duplicate}, nil
}

func (s *mcpServer) toolCancelQueue(args map[string]any) (map[string]any, error) {
	id, _ := args["queue_id"].(string)
	path := filepath.Join(s.cfg.StateDir, ".build_state")
	state, err := readBuildState(path)
	if err != nil {
		return nil, err
	}
	if fmt.Sprint(state.ActiveQueueEntry["id"]) == id {
		return nil, errors.New("queue_entry_running")
	}
	next := []map[string]any{}
	found := false
	for _, entry := range state.Queued {
		if fmt.Sprint(entry["id"]) == id {
			found = true
			continue
		}
		next = append(next, entry)
	}
	if !found {
		return nil, errors.New("queue_entry_not_found")
	}
	state.Queued = next
	if err := atomicWriteJSON(path, state, 0600); err != nil {
		return nil, err
	}
	return map[string]any{"queue_id": id, "cancelled": true}, nil
}

func (s *mcpServer) toolHistory(args map[string]any) map[string]any {
	limit := mcpOptionalInt(args, "limit", 50)
	offset := mcpOptionalInt(args, "offset", 0)
	records := decorateHistoryRecords(readHistory(filepath.Join(s.cfg.StateDir, ".build_history")))
	if offset > len(records) {
		offset = len(records)
	}
	end := offset + limit
	if end > len(records) {
		end = len(records)
	}
	return map[string]any{"history": records[offset:end], "total": len(records), "limit": limit, "offset": offset}
}

func (s *mcpServer) toolBuildLog(args map[string]any) (map[string]any, error) {
	id, _ := args["build_id"].(string)
	logRecord, err := readBuildLogByID(s.cfg.StateDir, id)
	if err != nil {
		return nil, err
	}
	return map[string]any{"build_id": id, "log": logRecord}, nil
}

func (s *mcpServer) toolAnalyzeBuildError(args map[string]any) (map[string]any, error) {
	id, _ := args["build_id"].(string)
	if _, err := readBuildLogByID(s.cfg.StateDir, id); err != nil {
		return nil, err
	}
	promptText := "Analyze Adlaire CI build log " + id + " and return a concise failure category, likely cause, and next action. Do not call external services."
	promptHash := mcpSHA256Hex([]byte(promptText))
	cfg, err := s.readMCPConfig()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"build_id":               id,
		"prompt_hash":            promptHash,
		"external_ai_api_called": false,
		"sampling_request": map[string]any{
			"jsonrpc": "2.0",
			"method":  "sampling/createMessage",
			"params": map[string]any{
				"timeout_ms": cfg.SamplingTimeoutMS,
				"messages": []map[string]any{{
					"role": "user",
					"content": map[string]any{
						"type": "text",
						"text": promptText,
					},
				}},
			},
		},
	}, nil
}

func (s *mcpServer) toolConfig() map[string]any {
	cfg, _, err := readOptionalJSONMap(filepath.Join(s.cfg.StateDir, ".server_config"))
	if err != nil {
		return map[string]any{"config": map[string]any{}, "corrupted": true}
	}
	if cfg == nil {
		cfg = map[string]any{}
	}
	return map[string]any{"config": cfg}
}

func (s *mcpServer) toolSetConfig(args map[string]any) (map[string]any, error) {
	key, _ := args["path"].(string)
	if !mcpAllowedServerConfigPath(key) {
		return nil, errors.New("invalid_config_path")
	}
	path := filepath.Join(s.cfg.StateDir, ".server_config")
	current, _, err := readOptionalJSONMap(path)
	if err != nil {
		return nil, err
	}
	if current == nil {
		current = map[string]any{}
	}
	current[key] = args["value"]
	if err := atomicWriteJSON(path, current, 0600); err != nil {
		return nil, err
	}
	return map[string]any{"path": key, "updated": true}, nil
}

func (s *mcpServer) toolSetMCPConfig(args map[string]any) (map[string]any, error) {
	if len(args) == 0 {
		return nil, errors.New("empty_mcp_config_patch")
	}
	current, err := s.readMCPConfig()
	if err != nil {
		return nil, err
	}
	if v, ok := args["tool_timeout_ms"]; ok {
		n, ok := mcpJSONInt(v)
		if !ok || n < 1000 || n > 300000 {
			return nil, errors.New("invalid_tool_timeout_ms")
		}
		current.ToolTimeoutMS = n
	}
	if v, ok := args["sampling_timeout_ms"]; ok {
		n, ok := mcpJSONInt(v)
		if !ok || n < 1000 || n > 300000 {
			return nil, errors.New("invalid_sampling_timeout_ms")
		}
		current.SamplingTimeoutMS = n
	}
	if v, ok := args["scopes"]; ok {
		list, ok := v.([]any)
		if !ok {
			return nil, errors.New("invalid_scopes")
		}
		scopes := []mcpScopeRecord{}
		seenToken := map[string]bool{}
		for _, item := range list {
			raw, ok := item.(map[string]any)
			if !ok {
				return nil, errors.New("invalid_scope_record")
			}
			tokenHash, _ := raw["token_hash"].(string)
			if !mcpValidSHA256(tokenHash) || seenToken[tokenHash] {
				return nil, errors.New("invalid_token_hash")
			}
			seenToken[tokenHash] = true
			rawScopes, ok := raw["scopes"].([]any)
			if !ok || len(rawScopes) == 0 {
				return nil, errors.New("invalid_scope_values")
			}
			scopeValues := []string{}
			seen := map[string]bool{}
			for _, rawScope := range rawScopes {
				scope, _ := rawScope.(string)
				if !mcpKnownScope(scope) || seen[scope] {
					return nil, errors.New("invalid_scope_values")
				}
				seen[scope] = true
				scopeValues = append(scopeValues, scope)
			}
			sort.Strings(scopeValues)
			now := mcpTime(s.now())
			scopes = append(scopes, mcpScopeRecord{TokenHash: tokenHash, Scopes: scopeValues, CreatedAt: now, UpdatedAt: now})
		}
		current.Scopes = scopes
	}
	if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".mcp_config"), current, 0600); err != nil {
		return nil, err
	}
	return map[string]any{"config": mcpConfigPublic(current)}, nil
}

func (s *mcpServer) toolCreateConfigSnapshot(args map[string]any) (map[string]any, error) {
	label, _ := args["label"].(string)
	now := s.now()
	id := s.nextConfigSnapshotID(now)
	files := map[string]any{}
	for _, name := range []string{".server_config", ".notify_config", ".repo_config", ".branch_config", ".dashboard_layout"} {
		value, exists, err := readOptionalJSONMap(filepath.Join(s.cfg.StateDir, name))
		if err != nil {
			return nil, err
		}
		if exists {
			files[name] = maskSecrets(value)
		}
	}
	snapshot := map[string]any{"id": id, "label": nil, "created_at": mcpTime(now), "created_by": "mcp", "files": files}
	if label != "" {
		snapshot["label"] = label
	}
	path := filepath.Join(s.cfg.StateDir, ".config_snapshots", id+".json")
	if err := atomicWriteJSON(path, snapshot, 0600); err != nil {
		return nil, err
	}
	return map[string]any{"snapshot_id": id}, nil
}

func (s *mcpServer) toolDiffConfigSnapshots(args map[string]any) (map[string]any, error) {
	left, _ := args["left_id"].(string)
	right, _ := args["right_id"].(string)
	if left == right {
		return nil, errors.New("same_snapshot_id")
	}
	leftFiles, err := s.readSnapshotFiles(left)
	if err != nil {
		return nil, err
	}
	rightFiles, err := s.readSnapshotFiles(right)
	if err != nil {
		return nil, err
	}
	changed := []string{}
	keys := map[string]bool{}
	for k := range leftFiles {
		keys[k] = true
	}
	for k := range rightFiles {
		keys[k] = true
	}
	for k := range keys {
		if mcpSHA256Any(leftFiles[k]) != mcpSHA256Any(rightFiles[k]) {
			changed = append(changed, k)
		}
	}
	sort.Strings(changed)
	return map[string]any{"left_id": left, "right_id": right, "changed_paths": changed}, nil
}

func (s *mcpServer) toolRestoreConfigSnapshot(args map[string]any) (map[string]any, error) {
	id, _ := args["snapshot_id"].(string)
	files, err := s.readSnapshotFiles(id)
	if err != nil {
		return nil, err
	}
	restored := []string{}
	for _, name := range mcpSortedKeys(files) {
		if !strings.HasPrefix(name, ".") || strings.Contains(name, "/") {
			continue
		}
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, name), files[name], 0600); err != nil {
			return nil, err
		}
		restored = append(restored, name)
	}
	return map[string]any{"snapshot_id": id, "restored_paths": restored}, nil
}

func (s *mcpServer) toolMetrics() map[string]any {
	metrics, err := s.readMCPMetrics()
	if err != nil {
		return map[string]any{"mcp": map[string]any{"tools": map[string]any{}}, "corrupted": true}
	}
	return map[string]any{"mcp": metrics}
}

func (s *mcpServer) toolAuditLog(args map[string]any) map[string]any {
	limit := mcpOptionalInt(args, "limit", 50)
	offset := mcpOptionalInt(args, "offset", 0)
	records := readJSONLines(filepath.Join(s.cfg.StateDir, ".mcp_audit_log"))
	if offset > len(records) {
		offset = len(records)
	}
	end := offset + limit
	if end > len(records) {
		end = len(records)
	}
	return map[string]any{"audit": records[offset:end], "total": len(records), "limit": limit, "offset": offset}
}

func (s *mcpServer) toolResendWebhook(args map[string]any) (map[string]any, error) {
	deliveryID, _ := args["delivery_id"].(string)
	now := s.now()
	path := filepath.Join(s.cfg.StateDir, ".webhook_resend_requests")
	requestID := s.nextLineID(path, "mcprsnd_", now)
	if err := appendJSONLine(path, map[string]any{
		"id":           requestID,
		"delivery_id":  deliveryID,
		"requested_at": mcpTime(now),
		"requested_by": s.currentClientName(),
		"source":       "mcp",
		"status":       "requested",
	}); err != nil {
		return nil, err
	}
	s.publish(mcpSSEFrame{Event: "resource-updated", Data: map[string]any{"uri": "adlaire://events", "updated_at": mcpTime(s.now())}})
	return map[string]any{"delivery_id": deliveryID, "request_id": requestID, "resend": "requested"}, nil
}

func (s *mcpServer) resourceText(uri string) (string, string, error) {
	var value any
	mime := "application/json"
	switch {
	case uri == "adlaire://status":
		value = s.toolStatus()
	case uri == "adlaire://queue":
		value = s.toolQueue()
	case uri == "adlaire://history":
		value = s.toolHistory(map[string]any{})
	case strings.HasPrefix(uri, "adlaire://history/"):
		id := strings.TrimPrefix(uri, "adlaire://history/")
		value = map[string]any{"build_id": id}
	case strings.HasPrefix(uri, "adlaire://logs/"):
		id := strings.TrimPrefix(uri, "adlaire://logs/")
		logRecord, err := readBuildLogByID(s.cfg.StateDir, id)
		if err != nil {
			return "", "", err
		}
		mime = "text/plain"
		data, _ := json.Marshal(logRecord)
		return string(data), mime, nil
	case uri == "adlaire://config":
		value = s.toolConfig()
	case uri == "adlaire://metrics":
		value = s.toolMetrics()
	case uri == "adlaire://audit":
		value = s.toolAuditLog(map[string]any{})
	default:
		return "", "", errors.New("unknown resource")
	}
	data, err := json.Marshal(value)
	return string(data), mime, err
}

func (s *mcpServer) readMCPConfig() (mcpStateConfig, error) {
	cfg := mcpStateConfig{ToolTimeoutMS: mcpDefaultToolMS, SamplingTimeoutMS: mcpDefaultSamplingMS, Scopes: []mcpScopeRecord{}}
	path := filepath.Join(s.cfg.StateDir, ".mcp_config")
	if err := readJSONIfExists(path, &cfg); err != nil {
		return cfg, err
	}
	if cfg.ToolTimeoutMS == 0 {
		cfg.ToolTimeoutMS = mcpDefaultToolMS
	}
	if cfg.SamplingTimeoutMS == 0 {
		cfg.SamplingTimeoutMS = mcpDefaultSamplingMS
	}
	if cfg.Scopes == nil {
		cfg.Scopes = []mcpScopeRecord{}
	}
	if err := mcpValidateStateConfig(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func mcpValidateStateConfig(cfg mcpStateConfig) error {
	if cfg.ToolTimeoutMS < 1000 || cfg.ToolTimeoutMS > 300000 {
		return errors.New("invalid_tool_timeout_ms")
	}
	if cfg.SamplingTimeoutMS < 1000 || cfg.SamplingTimeoutMS > 300000 {
		return errors.New("invalid_sampling_timeout_ms")
	}
	seenToken := map[string]bool{}
	for _, rec := range cfg.Scopes {
		if !mcpValidSHA256(rec.TokenHash) || seenToken[rec.TokenHash] {
			return errors.New("invalid_token_hash")
		}
		seenToken[rec.TokenHash] = true
		if len(rec.Scopes) == 0 {
			return errors.New("invalid_scope_values")
		}
		seenScope := map[string]bool{}
		for i, scope := range rec.Scopes {
			if !mcpKnownScope(scope) || seenScope[scope] {
				return errors.New("invalid_scope_values")
			}
			if i > 0 && rec.Scopes[i-1] > scope {
				return errors.New("invalid_scope_values")
			}
			seenScope[scope] = true
		}
		if _, err := time.Parse(apiTimeLayout, rec.CreatedAt); err != nil {
			return errors.New("invalid_scope_timestamp")
		}
		if _, err := time.Parse(apiTimeLayout, rec.UpdatedAt); err != nil {
			return errors.New("invalid_scope_timestamp")
		}
	}
	return nil
}

func (s *mcpServer) readMCPMetrics() (mcpMetricsFile, error) {
	metrics := mcpMetricsFile{Tools: map[string]mcpToolMetric{}}
	if err := readJSONIfExists(filepath.Join(s.cfg.StateDir, ".mcp_metrics"), &metrics); err != nil {
		return metrics, err
	}
	if metrics.Tools == nil {
		metrics.Tools = map[string]mcpToolMetric{}
	}
	return metrics, nil
}

func (s *mcpServer) updateMetrics(toolName, status string, durationMS int64) error {
	if status != "success" && status != "timeout" {
		status = "error"
	}
	metrics, err := s.readMCPMetrics()
	if err != nil {
		return err
	}
	item := metrics.Tools[toolName]
	switch status {
	case "success":
		item.SuccessCount++
	case "timeout":
		item.TimeoutCount++
	default:
		item.ErrorCount++
	}
	item.LastStatus = status
	item.LastDurationMS = durationMS
	item.LastAt = mcpTime(s.now())
	metrics.Tools[toolName] = item
	return atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".mcp_metrics"), metrics, 0600)
}

func (s *mcpServer) appendAudit(tool mcpToolSpec, args map[string]any, confirmationID string, status string, durationMS int64, code any, buildID, promptHash *string) error {
	now := s.now()
	paramsHash := mcpSHA256Any(mcpArgsWithoutConfirmation(args))
	confirmationValue := any(nil)
	if confirmationID != "" {
		confirmationValue = confirmationID
	}
	return appendJSONLine(filepath.Join(s.cfg.StateDir, ".mcp_audit_log"), map[string]any{
		"id":                 s.nextLineID(filepath.Join(s.cfg.StateDir, ".mcp_audit_log"), "mcpaud_", now),
		"timestamp":          mcpTime(now),
		"client_name":        s.currentClientName(),
		"tool":               tool.Name,
		"params_hash":        paramsHash,
		"scope":              tool.Scope,
		"confirmation_id":    confirmationValue,
		"status":             status,
		"duration_ms":        durationMS,
		"jsonrpc_error_code": code,
		"build_id":           mcpNullableString(buildID),
		"prompt_hash":        mcpNullableString(promptHash),
	})
}

func (s *mcpServer) createConfirmation(tool mcpToolSpec, args map[string]any) *mcpRPCError {
	id := "mcpconf_" + mcpRandomCrockford26()
	paramsHash := mcpSHA256Any(mcpArgsWithoutConfirmation(args))
	s.mu.Lock()
	s.confirmations[id] = mcpConfirmation{Tool: tool.Name, ParamsSHA: paramsHash, ExpiresAt: s.now().Add(mcpConfirmationTTL), Summary: "Confirm " + tool.Name}
	s.mu.Unlock()
	return &mcpRPCError{Code: -32004, Message: "Confirmation required", Data: map[string]any{"reason": "confirmation_required", "confirmation_id": id, "summary": "Confirm " + tool.Name}}
}

func (s *mcpServer) verifyConfirmation(tool mcpToolSpec, args map[string]any, id string) error {
	hash := mcpSHA256Any(mcpArgsWithoutConfirmation(args))
	s.mu.Lock()
	defer s.mu.Unlock()
	confirmation, ok := s.confirmations[id]
	if !ok {
		return errors.New("confirmation_unknown")
	}
	if s.now().After(confirmation.ExpiresAt) {
		delete(s.confirmations, id)
		return errors.New("confirmation_expired")
	}
	if confirmation.Tool != tool.Name || confirmation.ParamsSHA != hash {
		return errors.New("confirmation_mismatch")
	}
	delete(s.confirmations, id)
	return nil
}

func (s *mcpServer) scopeAllowed(scope string) bool {
	if s.cfg.ClientToken == "" {
		return true
	}
	cfg, err := s.readMCPConfig()
	if err != nil {
		return false
	}
	tokenHash := mcpSHA256Hex([]byte(s.cfg.ClientToken))
	for _, rec := range cfg.Scopes {
		if rec.TokenHash != tokenHash {
			continue
		}
		for _, allowed := range rec.Scopes {
			if allowed == scope {
				return true
			}
		}
		return false
	}
	return false
}

func (s *mcpServer) isInitialized() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.initialized
}

func (s *mcpServer) currentClientName() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clientName == "" {
		return "unknown"
	}
	return s.clientName
}

func (s *mcpServer) now() time.Time {
	return s.cfg.Now().UTC()
}

func (s *mcpServer) publish(frame mcpSSEFrame) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.subscribers {
		select {
		case ch <- frame:
		default:
		}
	}
}

func (s *mcpServer) nextLineID(path, prefix string, now time.Time) string {
	base := prefix + now.UTC().Format("20060102150405")
	records := readJSONLines(path)
	seen := map[string]bool{}
	for _, record := range records {
		if id, ok := record["id"].(string); ok {
			seen[id] = true
		}
	}
	if !seen[base] {
		return base
	}
	for i := 1; i <= 999; i++ {
		id := fmt.Sprintf("%s-%03d", base, i)
		if !seen[id] {
			return id
		}
	}
	return base + "-999"
}

func (s *mcpServer) nextConfigSnapshotID(now time.Time) string {
	return "cfgsnap_" + now.UTC().Format("20060102150405")
}

func (s *mcpServer) readSnapshotFiles(id string) (map[string]any, error) {
	if !strings.HasPrefix(id, "cfgsnap_") {
		return nil, errors.New("invalid_snapshot_id")
	}
	var snapshot struct {
		Files map[string]any `json:"files"`
	}
	if err := readJSONFile(filepath.Join(s.cfg.StateDir, ".config_snapshots", id+".json"), &snapshot); err != nil {
		return nil, err
	}
	if snapshot.Files == nil {
		return map[string]any{}, nil
	}
	return snapshot.Files, nil
}

func mcpToolInputSchema(tool mcpToolSpec) map[string]any {
	props := map[string]any{}
	for _, key := range mcpSortedKeysString(tool.Properties) {
		props[key] = mcpSchemaForKind(tool.Properties[key])
	}
	required := []string{}
	required = append(required, tool.Required...)
	return map[string]any{"type": "object", "required": required, "properties": props, "additionalProperties": false}
}

func mcpSchemaForKind(kind string) any {
	switch kind {
	case "string":
		return map[string]any{"type": "string"}
	case "integer":
		return map[string]any{"type": "integer"}
	case "boolean":
		return map[string]any{"type": "boolean"}
	case "object":
		return map[string]any{"type": "object", "required": []any{}, "properties": map[string]any{}, "additionalProperties": false}
	case "array":
		return map[string]any{"type": "array"}
	case "value":
		return map[string]any{"type": []string{"object", "array", "string", "number", "integer", "boolean", "null"}}
	case "string-or-null":
		return map[string]any{"type": []string{"string", "null"}}
	default:
		return map[string]any{"type": "string"}
	}
}

func mcpFindTool(name string) (mcpToolSpec, bool) {
	for _, tool := range mcpTools {
		if tool.Name == name {
			return tool, true
		}
	}
	return mcpToolSpec{}, false
}

func mcpValidateToolArgs(tool mcpToolSpec, args map[string]any) error {
	for _, required := range tool.Required {
		if _, ok := args[required]; !ok {
			return errors.New("missing_" + required)
		}
	}
	allowed := map[string]bool{"confirmation_id": tool.SideEffect}
	for key := range tool.Properties {
		allowed[key] = true
	}
	for key := range args {
		if !allowed[key] {
			return errors.New("unknown_" + key)
		}
	}
	switch tool.Name {
	case "adlaire.getStatus", "adlaire.getQueue", "adlaire.getConfig", "adlaire.getMcpConfig", "adlaire.getMetrics":
		if len(args) != 0 {
			return errors.New("params_must_be_empty")
		}
	case "adlaire.triggerBuild":
		target, _ := args["target"].(string)
		source, _ := args["source"].(string)
		if !mcpValidTarget(target) || !mcpValidSource(source) {
			return errors.New("invalid_build_request")
		}
		if options, ok := args["options"]; ok {
			optMap, ok := options.(map[string]any)
			if !ok || len(optMap) != 0 {
				return errors.New("invalid_options")
			}
		}
	case "adlaire.cancelQueueEntry":
		if !mcpValidID(fmt.Sprint(args["queue_id"])) {
			return errors.New("invalid_queue_id")
		}
	case "adlaire.getHistory", "adlaire.getAuditLog":
		if !mcpValidLimitOffset(args) {
			return errors.New("invalid_paging")
		}
	case "adlaire.getBuildLog", "adlaire.analyzeBuildError":
		if !mcpValidID(fmt.Sprint(args["build_id"])) {
			return errors.New("invalid_build_id")
		}
	case "adlaire.setConfig":
		path, _ := args["path"].(string)
		if !mcpAllowedServerConfigPath(path) {
			return errors.New("invalid_config_path")
		}
	case "adlaire.setMcpConfig":
		if len(mcpArgsWithoutConfirmation(args)) == 0 {
			return errors.New("empty_mcp_config_patch")
		}
	case "adlaire.createConfigSnapshot":
		if label, ok := args["label"]; ok && label != nil {
			if text, ok := label.(string); !ok || !mcpValidTinyString(text, 128) {
				return errors.New("invalid_label")
			}
		}
	case "adlaire.diffConfigSnapshots":
		left, _ := args["left_id"].(string)
		right, _ := args["right_id"].(string)
		if !strings.HasPrefix(left, "cfgsnap_") || !strings.HasPrefix(right, "cfgsnap_") || left == right {
			return errors.New("invalid_snapshot_id")
		}
	case "adlaire.restoreConfigSnapshot":
		id, _ := args["snapshot_id"].(string)
		if !strings.HasPrefix(id, "cfgsnap_") {
			return errors.New("invalid_snapshot_id")
		}
	case "adlaire.resendWebhook":
		if !mcpValidID(fmt.Sprint(args["delivery_id"])) {
			return errors.New("invalid_delivery_id")
		}
	}
	return nil
}

func mcpAllowedServerConfigPath(key string) bool {
	switch key {
	case "polling_interval_minutes", "github_pat_expires_at", "cooldown_seconds", "maintenance_mode", "notification", "dashboard_layout", "api_rate_limit":
		return true
	default:
		return false
	}
}

func mcpValidLimitOffset(args map[string]any) bool {
	if v, ok := args["limit"]; ok {
		n, ok := mcpJSONInt(v)
		if !ok || n < 1 || n > 100 {
			return false
		}
	}
	if v, ok := args["offset"]; ok {
		n, ok := mcpJSONInt(v)
		if !ok || n < 0 {
			return false
		}
	}
	return true
}

func mcpOptionalInt(args map[string]any, key string, fallback int) int {
	if v, ok := args[key]; ok {
		if n, ok := mcpJSONInt(v); ok {
			return n
		}
	}
	return fallback
}

func mcpValidID(id string) bool {
	return mcpValidTinyString(id, 128) && !strings.ContainsAny(id, "/\x00\r\n")
}

func mcpValidTarget(target string) bool {
	return mcpValidTinyString(target, 128) && target != "/" && !strings.Contains(target, "..") && !strings.ContainsAny(target, "\x00\r\n")
}

func mcpValidSource(source string) bool {
	if strings.HasPrefix(source, "local:") {
		path := strings.TrimPrefix(source, "local:")
		return filepath.IsAbs(path) && !strings.Contains(path, "..") && !strings.ContainsAny(path, "\x00\r\n")
	}
	if strings.HasPrefix(source, "git:") {
		ref := strings.TrimPrefix(source, "git:")
		return mcpValidTinyString(ref, 128) && !strings.ContainsAny(ref, " \x00\r\n") && !strings.Contains(ref, "..") && !strings.Contains(ref, "@{")
	}
	return false
}

func mcpValidTinyString(value string, maxBytes int) bool {
	if value == "" || len(value) > maxBytes || !utf8.ValidString(value) {
		return false
	}
	return !strings.ContainsAny(value, "\x00\r\n")
}

func mcpValidSHA256(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, r := range value {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return false
		}
	}
	return true
}

func mcpKnownScope(scope string) bool {
	for _, tool := range mcpTools {
		if tool.Scope == scope {
			return true
		}
	}
	return false
}

func mcpResourceSpecs() [][4]string {
	return [][4]string{
		{"adlaire://status", "status", "Current status", "application/json"},
		{"adlaire://queue", "queue", "Build queue", "application/json"},
		{"adlaire://history", "history", "Build history", "application/json"},
		{"adlaire://history/{build_id}", "history detail", "Build detail", "application/json"},
		{"adlaire://logs/{build_id}", "build log", "Build log", "text/plain"},
		{"adlaire://config", "config", "Effective config", "application/json"},
		{"adlaire://metrics", "metrics", "Tool and build metrics", "application/json"},
		{"adlaire://audit", "audit", "MCP audit log tail", "application/json"},
	}
}

func mcpKnownResource(uri string) bool {
	for _, item := range mcpResourceSpecs() {
		if uri == item[0] || (strings.HasSuffix(item[0], "{build_id}") && strings.HasPrefix(uri, strings.TrimSuffix(item[0], "{build_id}")) && mcpValidID(strings.TrimPrefix(uri, strings.TrimSuffix(item[0], "{build_id}")))) {
			return true
		}
	}
	return false
}

func mcpDecodeJSONObject(data []byte) (map[string]json.RawMessage, []string, bool) {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return nil, nil, false
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return nil, nil, false
	}
	keys := []string{}
	seen := map[string]bool{}
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return nil, nil, false
		}
		key, ok := tok.(string)
		if !ok || seen[key] {
			return nil, nil, false
		}
		seen[key] = true
		keys = append(keys, key)
		var skip any
		if err := dec.Decode(&skip); err != nil {
			return nil, nil, false
		}
	}
	tok, err = dec.Token()
	if err != nil {
		return nil, nil, false
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '}' {
		return nil, nil, false
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, nil, false
	}
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, false
	}
	return raw, keys, true
}

func mcpRPCKeysAllowed(order []string) bool {
	allowed := map[string]bool{"jsonrpc": false, "id": false, "method": false, "params": true}
	for _, key := range order {
		if _, ok := allowed[key]; !ok {
			return false
		}
		allowed[key] = true
	}
	return allowed["jsonrpc"] && allowed["method"]
}

func mcpParseID(raw json.RawMessage) (any, bool, bool) {
	if raw == nil {
		return nil, true, false
	}
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil, true, true
	}
	var text string
	if json.Unmarshal(raw, &text) == nil {
		return text, true, true
	}
	var number float64
	if json.Unmarshal(raw, &number) == nil {
		return number, true, true
	}
	return nil, false, true
}

func mcpDecodeParamsObject(raw json.RawMessage) (map[string]any, bool) {
	var params map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&params); err != nil || params == nil {
		return nil, false
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, false
	}
	return mcpNormalizeJSONNumbers(params).(map[string]any), true
}

func mcpNormalizeJSONNumbers(value any) any {
	switch v := value.(type) {
	case map[string]any:
		out := map[string]any{}
		for key, item := range v {
			out[key] = mcpNormalizeJSONNumbers(item)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = mcpNormalizeJSONNumbers(item)
		}
		return out
	case json.Number:
		if i, err := strconv.ParseInt(v.String(), 10, 64); err == nil {
			return int(i)
		}
		if f, err := strconv.ParseFloat(v.String(), 64); err == nil {
			return f
		}
		return v.String()
	default:
		return v
	}
}

func mcpJSONInt(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		if v == float64(int(v)) {
			return int(v), true
		}
	}
	return 0, false
}

func mcpValidContentType(values []string) bool {
	if len(values) != 1 {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(values[0])) {
	case "application/json", "application/json; charset=utf-8":
		return true
	default:
		return false
	}
}

func mcpWriteHTTPError(w http.ResponseWriter, status int, message string) {
	mcpWriteJSON(w, status, map[string]string{"error": message})
}

func mcpWriteJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		status = http.StatusInternalServerError
		data = []byte(`{"error":"Internal server error"}`)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

func (s *mcpServer) writeRPCError(w http.ResponseWriter, status int, id any, code int, message, reason string) {
	mcpWriteJSON(w, status, mcpRPCResponse{JSONRPC: "2.0", ID: id, Error: mcpError(code, message, reason)})
}

func mcpError(code int, message, reason string) *mcpRPCError {
	return &mcpRPCError{Code: code, Message: message, Data: map[string]any{"reason": reason}}
}

func mcpWriteSSEFrame(w io.Writer, frame mcpSSEFrame) bool {
	data, err := json.Marshal(frame.Data)
	if err != nil {
		return false
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", frame.Event, data); err != nil {
		return false
	}
	return true
}

func mcpArgsWithoutConfirmation(args map[string]any) map[string]any {
	out := map[string]any{}
	for key, value := range args {
		if key != "confirmation_id" {
			out[key] = value
		}
	}
	return out
}

func mcpConfigPublic(cfg mcpStateConfig) map[string]any {
	scopes := []map[string]any{}
	for _, rec := range cfg.Scopes {
		scopes = append(scopes, map[string]any{"token_hash": rec.TokenHash, "scopes": rec.Scopes, "created_at": rec.CreatedAt, "updated_at": rec.UpdatedAt})
	}
	return map[string]any{"tool_timeout_ms": cfg.ToolTimeoutMS, "sampling_timeout_ms": cfg.SamplingTimeoutMS, "scopes": scopes}
}

func mcpSHA256Any(value any) string {
	data, _ := json.Marshal(value)
	return mcpSHA256Hex(data)
}

func mcpSHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func mcpRandomCrockford26() string {
	const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return strings.Repeat("0", 26)
	}
	var b strings.Builder
	for i := 0; i < 26; i++ {
		bit := i * 5
		value := 0
		for j := 0; j < 5; j++ {
			pos := bit + j
			if pos >= 128 {
				value <<= 1
				continue
			}
			if raw[pos/8]&(1<<uint(7-(pos%8))) != 0 {
				value = (value << 1) | 1
			} else {
				value <<= 1
			}
		}
		b.WriteByte(alphabet[value&31])
	}
	return b.String()
}

func mcpRemoteAddr(value string) string {
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return "127.0.0.1:0"
	}
	if host == "::1" {
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, port)
}

func mcpTime(t time.Time) string {
	return t.UTC().Format(apiTimeLayout)
}

func mcpNullableString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func mcpSortedKeys(value map[string]any) []string {
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func mcpSortedKeysString(value map[string]string) []string {
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
