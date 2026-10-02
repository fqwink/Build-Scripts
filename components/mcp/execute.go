package mcp

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/fqwink/build-scripts/components/statefile"
)

const apiTimeLayout = "2006-01-02T15:04:05Z"

type apiBuildState struct {
	Running          bool             `json:"running"`
	CurrentBuildID   *string          `json:"current_build_id"`
	ActiveQueueEntry map[string]any   `json:"active_queue_entry"`
	Queued           []map[string]any `json:"queued"`
}

type apiBuildStatus struct {
	Status                *string `json:"status"`
	LastBranch            *string `json:"last_branch"`
	LastTargetFile        *string `json:"last_target_file"`
	LastBlobSHA           *string `json:"last_blob_sha"`
	LastCommitSHA         *string `json:"last_commit_sha"`
	LastStartedAt         *string `json:"last_started_at"`
	LastFinishedAt        *string `json:"last_finished_at"`
	LastTargetStatus      *string `json:"last_target_status"`
	LastDeployAt          *string `json:"last_deploy_at"`
	LastDeployStatus      *string `json:"last_deploy_status"`
	LastTrigger           *string `json:"last_trigger"`
	PendingTransfersCount int     `json:"pending_transfers_count"`
	NotifyPendingCount    int     `json:"notify_pending_count"`
	CircuitOpen           bool    `json:"circuit_open"`
	Running               bool    `json:"running"`
}

type apiHistoryRecord struct {
	ID              string   `json:"id"`
	Branch          string   `json:"branch,omitempty"`
	TargetFile      string   `json:"target_file,omitempty"`
	BuildAt         string   `json:"build_at,omitempty"`
	StartedAt       string   `json:"started_at,omitempty"`
	FinishedAt      string   `json:"finished_at,omitempty"`
	SHA             *string  `json:"sha,omitempty"`
	CommitSHA       *string  `json:"commit_sha,omitempty"`
	BlobSHA         *string  `json:"blob_sha,omitempty"`
	Status          string   `json:"status"`
	Trigger         string   `json:"trigger,omitempty"`
	DurationSeconds int64    `json:"duration_seconds"`
	FailureCategory *string  `json:"failure_category,omitempty"`
	Flagged         bool     `json:"flagged"`
	Tags            []string `json:"tags"`
	OutputSizeBytes *int64   `json:"output_size_bytes"`
	OutputSHA256    *string  `json:"output_sha256,omitempty"`
}

type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string { return e.Msg }

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

func validateAPIStateDir(path string) error {
	if path == "" {
		return exitError{Code: 2, Msg: "state directory must not be empty"}
	}
	if !filepath.IsAbs(path) {
		return exitError{Code: 2, Msg: "state directory must be absolute: " + path}
	}
	info, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return exitError{Code: 2, Msg: "state directory not found: " + path}
		}
		return exitError{Code: 2, Msg: "state directory not found: " + path}
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return exitError{Code: 2, Msg: "state directory must not be symlink: " + path}
	}
	if !info.IsDir() {
		return exitError{Code: 2, Msg: "state path is not directory: " + path}
	}
	return nil
}

func readMCPTokenFile(path string) (string, error) {
	if path == "" || !utf8.ValidString(path) || strings.ContainsAny(path, "\x00\n\r") {
		return "", errors.New("invalid_client_token_file")
	}
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Size() > 4096 {
		return "", errors.New("invalid_client_token_file")
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return readMCPTokenReader(file)
}

func readMCPTokenReader(reader io.Reader) (string, error) {
	data, err := io.ReadAll(io.LimitReader(reader, 4097))
	if err != nil {
		return "", err
	}
	if len(data) > 4096 {
		return "", errors.New("client_token_too_large")
	}
	token := strings.TrimSpace(string(data))
	if token == "" || !utf8.ValidString(token) || strings.ContainsAny(token, "\x00\r\n") {
		return "", errors.New("invalid_client_token")
	}
	return token, nil
}

func readJSONIfExists(path string, out any) error {
	return statefile.ReadJSONIfExists(path, out)
}

func readJSONFile(path string, out any) error {
	return statefile.ReadJSON(path, out)
}

func atomicWriteJSON(path string, value any, mode os.FileMode) error {
	return statefile.WriteJSONAtomic(path, value, mode)
}

func readOptionalJSONMap(path string) (map[string]any, bool, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	var value map[string]any
	if err := readJSONFile(path, &value); err != nil {
		return nil, true, err
	}
	return value, true, nil
}

func readJSONLines(path string) []map[string]any {
	records, err := statefile.ReadJSONLines(path)
	if err != nil {
		return []map[string]any{}
	}
	return records
}

func appendJSONLine(path string, value any) error {
	return statefile.AppendJSONLine(path, value, 0600)
}

func readBuildLogByID(stateDir, id string) (map[string]any, error) {
	logPath := filepath.Join(stateDir, ".build_logs", id+".json")
	var logRecord map[string]any
	if err := readJSONFile(logPath, &logRecord); err == nil {
		return logRecord, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	archivePath := filepath.Join(stateDir, ".build_logs", "archive", id+".json.gz")
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	if err := json.NewDecoder(reader).Decode(&logRecord); err != nil {
		return nil, err
	}
	return logRecord, nil
}

func maskSecrets(value map[string]any) map[string]any {
	out := map[string]any{}
	for key, raw := range value {
		lower := strings.ToLower(key)
		if strings.Contains(lower, "token") || strings.Contains(lower, "secret") || strings.Contains(lower, "password") {
			out[key] = "***"
			continue
		}
		if nested, ok := raw.(map[string]any); ok {
			out[key] = maskSecrets(nested)
			continue
		}
		out[key] = raw
	}
	return out
}

func readBuildState(path string) (apiBuildState, error) {
	state := apiBuildState{Queued: []map[string]any{}}
	if err := readJSONIfExists(path, &state); err != nil {
		return state, err
	}
	if state.Queued == nil {
		state.Queued = []map[string]any{}
	}
	if state.ActiveQueueEntry == nil {
		state.ActiveQueueEntry = map[string]any{}
	}
	return state, nil
}

func readBuildStatus(path string) (apiBuildStatus, bool, error) {
	var status apiBuildStatus
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return status, false, nil
	}
	if err := readJSONFile(path, &status); err != nil {
		return status, true, err
	}
	return status, true, nil
}

func readHistory(path string) []apiHistoryRecord {
	data, err := os.ReadFile(path)
	if err != nil {
		return []apiHistoryRecord{}
	}
	records := []apiHistoryRecord{}
	seen := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record apiHistoryRecord
		if json.Unmarshal([]byte(line), &record) == nil && record.ID != "" && record.Status != "" && !seen[record.ID] {
			seen[record.ID] = true
			records = append(records, record)
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return firstNonEmpty(records[i].FinishedAt, records[i].StartedAt, records[i].BuildAt, records[i].ID) > firstNonEmpty(records[j].FinishedAt, records[j].StartedAt, records[j].BuildAt, records[j].ID)
	})
	return records
}

func decorateHistoryRecords(records []apiHistoryRecord) []apiHistoryRecord {
	out := make([]apiHistoryRecord, len(records))
	for i, record := range records {
		if record.BuildAt == "" {
			record.BuildAt = firstNonEmpty(record.FinishedAt, record.StartedAt)
		}
		if record.SHA == nil {
			if sha := firstNonEmpty(stringPtrValue(record.CommitSHA), stringPtrValue(record.BlobSHA)); sha != "" {
				record.SHA = &sha
			}
		}
		out[i] = record
	}
	return out
}

func sortedQueueEntries(entries []map[string]any) []map[string]any {
	out := make([]map[string]any, len(entries))
	copy(out, entries)
	sort.SliceStable(out, func(i, j int) bool {
		leftPriority := queuePriorityRank(fmt.Sprint(out[i]["priority"]))
		rightPriority := queuePriorityRank(fmt.Sprint(out[j]["priority"]))
		if leftPriority != rightPriority {
			return leftPriority < rightPriority
		}
		leftSeq := queueCreatedSeq(out[i])
		rightSeq := queueCreatedSeq(out[j])
		if leftSeq != rightSeq {
			return leftSeq < rightSeq
		}
		return fmt.Sprint(out[i]["id"]) < fmt.Sprint(out[j]["id"])
	})
	return out
}

func queuePriorityRank(priority string) int {
	switch priority {
	case "urgent":
		return 0
	case "high":
		return 1
	case "normal":
		return 2
	case "low":
		return 3
	default:
		return 4
	}
}

func queueCreatedSeq(entry map[string]any) int {
	switch v := entry["created_seq"].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		n, _ := strconv.Atoi(v.String())
		return n
	default:
		return 0
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (s *mcpServer) enqueueMCPBuildRequest(path string, payload map[string]any) (string, bool, error) {
	state, err := readBuildState(path)
	if err != nil {
		return "", false, err
	}
	for _, entry := range append([]map[string]any{state.ActiveQueueEntry}, state.Queued...) {
		if fmt.Sprint(entry["trigger"]) == "mcp" && fmt.Sprint(entry["payload"]) == fmt.Sprint(payload) {
			return fmt.Sprint(entry["id"]), true, nil
		}
	}
	now := s.now()
	id := s.nextLineID(path, "queue_", now)
	state.Queued = append(state.Queued, map[string]any{
		"id":           id,
		"trigger":      "mcp",
		"priority":     "normal",
		"requested_by": "mcp",
		"created_at":   mcpTime(now),
		"created_seq":  len(state.Queued) + 1,
		"payload":      payload,
	})
	if err := atomicWriteJSON(path, state, 0600); err != nil {
		return "", false, err
	}
	return id, false, nil
}

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}
