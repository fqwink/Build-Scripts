package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	adminBinaryName           = "adlaire-ci-admin"
	adminHTTPTimeout          = 30 * time.Second
	adminMaxResponseBodyBytes = 1024 * 1024
)

var adminBinaryVersion = "V.0.0-dev"

type adminCommandSpec struct {
	Method  string
	Path    string
	HasBody bool
}

var adminCommandSpecs = map[string]adminCommandSpec{
	"status":          {Method: http.MethodGet, Path: "/api/status"},
	"queue":           {Method: http.MethodGet, Path: "/api/queue"},
	"history":         {Method: http.MethodGet, Path: "/api/history"},
	"trigger-build":   {Method: http.MethodPost, Path: "/api/build", HasBody: true},
	"cancel-queue":    {Method: http.MethodDelete, Path: "/api/queue/"},
	"config-snapshot": {Method: http.MethodPost, Path: "/api/config-snapshots", HasBody: true},
	"events":          {Method: http.MethodGet, Path: "/api/events"},
}

type adminCLIConfig struct {
	APIURL string
	Token  string
	JSON   bool
}

type adminCLIRequest struct {
	Command string
	Method  string
	Path    string
	Body    []byte
}

func RunAdmin(args []string, stdout, stderr io.Writer) int {
	return runAdmin(args, stdout, stderr, newAdminHTTPClient())
}

func runAdmin(args []string, stdout, stderr io.Writer, client *http.Client) int {
	if hasExactArg(args, "--help") {
		fmt.Fprintln(stdout, "Usage: adlaire-ci-admin --api-url url --token token [--json] command [command-args]")
		return 0
	}
	if hasExactArg(args, "--version") {
		fmt.Fprintf(stdout, "%s %s go=%s\n", adminBinaryName, adminBinaryVersion, runtime.Version())
		return 0
	}
	if !safeArgvTokens(args) {
		fmt.Fprintln(stderr, "invalid command line token")
		return 2
	}

	cfg, command, commandArgs, parseErr := parseAdminArgs(args)
	if parseErr != "" {
		fmt.Fprintln(stderr, parseErr)
		return 2
	}

	reqSpec, reqErr := buildAdminRequest(command, commandArgs)
	if reqErr != "" {
		fmt.Fprintln(stderr, reqErr)
		return 2
	}

	status, body, contentType, callErr := callAdminAPI(client, cfg, reqSpec)
	if callErr != "" {
		fmt.Fprintln(stderr, callErr)
		return 1
	}
	if status < 200 || status > 299 {
		if status >= 300 && status <= 599 {
			fmt.Fprintf(stderr, "api error: %d\n", status)
		} else {
			fmt.Fprintln(stderr, "api error: connection failed")
		}
		return 1
	}
	if !isValidAdminJSONContentType(contentType) {
		fmt.Fprintln(stderr, "api error: invalid response")
		return 1
	}

	trimmed := trimASCIIWhitespace(body)
	response, ok := decodeAdminJSONObject(trimmed)
	if !ok {
		fmt.Fprintln(stderr, "api error: invalid response")
		return 1
	}
	if !validateAdminResponseShape(command, response) {
		fmt.Fprintln(stderr, "api error: invalid response")
		return 1
	}

	if cfg.JSON {
		stdout.Write(trimmed)
		fmt.Fprintln(stdout)
		return 0
	}

	line, ok := adminHumanOutput(command, response)
	if !ok {
		fmt.Fprintln(stderr, "api error: invalid response")
		return 1
	}
	fmt.Fprintln(stdout, line)
	return 0
}

func newAdminHTTPClient() *http.Client {
	return &http.Client{
		Timeout: adminHTTPTimeout,
		Transport: &http.Transport{
			Proxy: nil,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func hasExactArg(args []string, target string) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}
	return false
}

func safeArgvTokens(args []string) bool {
	for _, arg := range args {
		if !utf8.ValidString(arg) {
			return false
		}
		for _, r := range arg {
			if r < 0x20 || r == 0x7f {
				return false
			}
		}
	}
	return true
}

func parseAdminArgs(args []string) (adminCLIConfig, string, []string, string) {
	cfg := adminCLIConfig{}
	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			switch arg {
			case "--api-url":
				if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
					return cfg, "", nil, "missing value: --api-url"
				}
				if cfg.APIURL != "" {
					return cfg, "", nil, "usage error"
				}
				cfg.APIURL = args[i+1]
				i += 2
			case "--token":
				if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
					return cfg, "", nil, "missing value: --token"
				}
				if cfg.Token != "" {
					return cfg, "", nil, "usage error"
				}
				cfg.Token = args[i+1]
				i += 2
			case "--json":
				if cfg.JSON {
					return cfg, "", nil, "usage error"
				}
				cfg.JSON = true
				i++
			default:
				return cfg, "", nil, "unknown option: " + arg
			}
			continue
		}
		if strings.HasPrefix(arg, "-") {
			return cfg, "", nil, "unknown option: " + arg
		}
		break
	}

	if cfg.APIURL == "" || cfg.Token == "" || i >= len(args) {
		return cfg, "", nil, "usage error"
	}
	normalizedURL, ok := normalizeAdminAPIURL(cfg.APIURL)
	if !ok {
		return cfg, "", nil, "usage error"
	}
	cfg.APIURL = normalizedURL
	if !validAdminToken(cfg.Token) {
		return cfg, "", nil, "usage error"
	}

	command := args[i]
	commandArgs := args[i+1:]
	if _, ok := adminCommandSpecs[command]; !ok {
		return cfg, "", nil, "unknown command: " + command
	}
	if !validAdminCommandArgs(command, commandArgs) {
		return cfg, "", nil, "usage error"
	}
	return cfg, command, commandArgs, ""
}

func normalizeAdminAPIURL(input string) (string, bool) {
	u, err := url.Parse(input)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", false
	}
	if u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", false
	}
	p := u.EscapedPath()
	decodedPath := u.Path
	if p == "" {
		decodedPath = ""
	} else if !strings.HasPrefix(decodedPath, "/") {
		return "", false
	}
	if strings.Contains(decodedPath, "\x00") || strings.Contains(decodedPath, "\\") || strings.Contains(decodedPath, "//") {
		return "", false
	}
	for _, segment := range strings.Split(decodedPath, "/") {
		if segment == "." || segment == ".." {
			return "", false
		}
	}
	if decodedPath == "/" {
		decodedPath = ""
	} else if strings.HasSuffix(decodedPath, "/") {
		decodedPath = strings.TrimSuffix(decodedPath, "/")
	}
	return u.Scheme + "://" + u.Host + decodedPath, true
}

func validAdminToken(token string) bool {
	if token == "" || len([]byte(token)) > 4096 || !utf8.ValidString(token) {
		return false
	}
	return !strings.ContainsAny(token, "\x00\r\n")
}

func validAdminCommandArgs(command string, args []string) bool {
	switch command {
	case "status", "queue", "history", "trigger-build", "events":
		return len(args) == 0
	case "cancel-queue":
		if len(args) != 1 {
			return false
		}
		id := args[0]
		return id != "" && !strings.Contains(id, "/") && !strings.Contains(id, "..") && !strings.Contains(id, "\x00")
	case "config-snapshot":
		if len(args) > 1 {
			return false
		}
		if len(args) == 0 {
			return true
		}
		return validAdminConfigSnapshotLabel(args[0])
	default:
		return false
	}
}

func validAdminConfigSnapshotLabel(label string) bool {
	if label == "" || !utf8.ValidString(label) {
		return false
	}
	count := 0
	for _, r := range label {
		if r == '\n' || r == '\r' || r == '\x00' || r == '\ufeff' {
			return false
		}
		count++
	}
	return count <= 128
}

func buildAdminRequest(command string, args []string) (adminCLIRequest, string) {
	spec := adminCommandSpecs[command]
	req := adminCLIRequest{
		Command: command,
		Method:  spec.Method,
		Path:    spec.Path,
	}
	switch command {
	case "trigger-build":
		req.Body = []byte(`{}`)
	case "cancel-queue":
		req.Path += url.PathEscape(args[0])
	case "config-snapshot":
		if len(args) == 0 {
			req.Body = []byte(`{"label":null}`)
		} else {
			body, err := json.Marshal(map[string]string{"label": args[0]})
			if err != nil {
				return req, "usage error"
			}
			req.Body = body
		}
	}
	if !spec.HasBody && len(req.Body) > 0 {
		return req, "usage error"
	}
	return req, ""
}

func callAdminAPI(client *http.Client, cfg adminCLIConfig, reqSpec adminCLIRequest) (int, []byte, []string, string) {
	httpReq, err := http.NewRequest(reqSpec.Method, cfg.APIURL+reqSpec.Path, bytes.NewReader(reqSpec.Body))
	if err != nil {
		return 0, nil, nil, "api error: connection failed"
	}
	httpReq.Header.Set("Authorization", "Bearer "+cfg.Token)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", adminBinaryName+"/"+adminBinaryVersion)
	if len(reqSpec.Body) > 0 {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, nil, nil, "api error: connection failed"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, adminMaxResponseBodyBytes+1))
	if err != nil {
		return resp.StatusCode, nil, resp.Header.Values("Content-Type"), "api error: connection failed"
	}
	if len(body) > adminMaxResponseBodyBytes {
		return resp.StatusCode, nil, resp.Header.Values("Content-Type"), "api error: invalid response"
	}
	return resp.StatusCode, body, resp.Header.Values("Content-Type"), ""
}

func isValidAdminJSONContentType(values []string) bool {
	if len(values) != 1 {
		return false
	}
	ct := strings.TrimSpace(values[0])
	return ct == "application/json" || ct == "application/json; charset=utf-8"
}

func trimASCIIWhitespace(data []byte) []byte {
	return bytes.Trim(data, " \t\r\n")
}

func decodeAdminJSONObject(data []byte) (map[string]any, bool) {
	if len(data) == 0 {
		return nil, false
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var value any
	if err := dec.Decode(&value); err != nil {
		return nil, false
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, false
	}
	obj, ok := value.(map[string]any)
	return obj, ok
}

func validateAdminResponseShape(command string, response map[string]any) bool {
	switch command {
	case "status":
		_, okStatus := response["last_build_status"].(string)
		_, okRunning := response["running"].(bool)
		return okStatus && okRunning
	case "queue":
		if active := response["active"]; active != nil {
			obj, ok := active.(map[string]any)
			if !ok {
				return false
			}
			if _, ok := obj["id"].(string); !ok {
				return false
			}
		}
		_, okQueued := response["queued"].([]any)
		return okQueued
	case "history":
		if _, ok := integerJSONNumber(response["total"]); !ok {
			return false
		}
		history, ok := response["history"].([]any)
		if !ok {
			return false
		}
		if len(history) == 0 {
			return true
		}
		first, ok := history[0].(map[string]any)
		if !ok {
			return false
		}
		_, ok = first["id"].(string)
		return ok
	case "trigger-build":
		queueID, okID := response["queue_id"].(string)
		queued, okQueued := response["queued"].(bool)
		dispatch, okDispatch := response["dispatch"].(string)
		return okID && queueID != "" && okQueued && queued && okDispatch && (dispatch == "requested" || dispatch == "timer_fallback")
	case "cancel-queue":
		return response != nil
	case "config-snapshot":
		id, ok := response["id"].(string)
		return ok && id != ""
	case "events":
		if _, ok := integerJSONNumber(response["total"]); !ok {
			return false
		}
		_, ok := response["events"].([]any)
		return ok
	default:
		return false
	}
}

func adminHumanOutput(command string, response map[string]any) (string, bool) {
	switch command {
	case "status":
		status, okStatus := response["last_build_status"].(string)
		running, okRunning := response["running"].(bool)
		if !okStatus || !okRunning {
			return "", false
		}
		return fmt.Sprintf("status=%s running=%t", status, running), true
	case "queue":
		activeText := "none"
		if active := response["active"]; active != nil {
			obj, ok := active.(map[string]any)
			if !ok {
				return "", false
			}
			id, ok := obj["id"].(string)
			if !ok {
				return "", false
			}
			activeText = id
		}
		queued, ok := response["queued"].([]any)
		if !ok {
			return "", false
		}
		return fmt.Sprintf("active=%s queued=%d", activeText, len(queued)), true
	case "history":
		total, ok := integerJSONNumber(response["total"])
		if !ok {
			return "", false
		}
		history, ok := response["history"].([]any)
		if !ok {
			return "", false
		}
		latest := "none"
		if len(history) > 0 {
			first, ok := history[0].(map[string]any)
			if !ok {
				return "", false
			}
			id, ok := first["id"].(string)
			if !ok {
				return "", false
			}
			latest = id
		}
		return fmt.Sprintf("total=%d latest=%s", total, latest), true
	case "trigger-build":
		queueID, ok := response["queue_id"].(string)
		if !ok || queueID == "" {
			return "", false
		}
		return "queued=" + queueID, true
	case "cancel-queue":
		return "queue cancelled", true
	case "config-snapshot":
		id, ok := response["id"].(string)
		if !ok || id == "" {
			return "", false
		}
		return "snapshot=" + id, true
	case "events":
		total, ok := integerJSONNumber(response["total"])
		if !ok {
			return "", false
		}
		return fmt.Sprintf("events=%d", total), true
	default:
		return "", false
	}
}

func integerJSONNumber(value any) (int64, bool) {
	number, ok := value.(json.Number)
	if !ok {
		return 0, false
	}
	text := number.String()
	if strings.ContainsAny(text, ".eE") {
		return 0, false
	}
	parsed, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}
