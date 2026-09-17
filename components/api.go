package components

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultAPIAddr      = "127.0.0.1:8765"
	defaultQueueMaxSize = 3
	apiTimeLayout       = "2006-01-02T15:04:05Z"
	passwordIterations  = 260000
)

type APIConfig struct {
	Addr     string
	StateDir string
	Now      func() time.Time
}

type APIServer struct {
	cfg       APIConfig
	startedAt time.Time
	mu        sync.Mutex
	sessions  map[string]apiSession
}

type apiSession struct {
	TokenHash  string
	CreatedAt  string
	ExpiresAt  string
	LastUsedAt string
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

type apiCredentials struct {
	PasswordHash  string  `json:"password_hash"`
	Salt          string  `json:"salt"`
	Algorithm     string  `json:"algorithm"`
	Iterations    int     `json:"iterations"`
	MustChange    bool    `json:"must_change"`
	FailedCount   int     `json:"failed_count"`
	LockedUntil   *string `json:"locked_until"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	LastLoginAt   *string `json:"last_login_at"`
	LastFailureAt *string `json:"last_failure_at"`
}

type apiBuildState struct {
	Running                 bool             `json:"running"`
	CurrentBuildID          *string          `json:"current_build_id"`
	Queued                  []map[string]any `json:"queued"`
	LastStartedAt           *string          `json:"last_started_at"`
	LastFinishedAt          *string          `json:"last_finished_at"`
	WeeklySummaryLastSentAt *string          `json:"weekly_summary_last_sent_at"`
	WeeklySummarySentDate   *string          `json:"weekly_summary_sent_date"`
}

type apiCircuitState struct {
	Open                bool    `json:"open"`
	ConsecutiveFailures int     `json:"consecutive_failures"`
	OpenedAt            *string `json:"opened_at"`
	LastFailureAt       *string `json:"last_failure_at"`
	LastError           *string `json:"last_error"`
}

type apiServerConfig struct {
	LogMaxLines           int     `json:"log_max_lines"`
	HistoryMaxCount       int     `json:"history_max_count"`
	BuildTimeoutSeconds   int     `json:"build_timeout_seconds"`
	LogRetentionDays      int     `json:"log_retention_days"`
	LogArchiveAfterDays   int     `json:"log_archive_after_days"`
	LogLevel              string  `json:"log_level"`
	PATExpiresAt          *string `json:"pat_expires_at"`
	SnapshotsKeep         int     `json:"snapshots_keep"`
	QueueMaxSize          int     `json:"queue_max_size"`
	BuildRetryMax         int     `json:"build_retry_max"`
	BuildRetryBaseSeconds int     `json:"build_retry_base_seconds"`
	CommitStatusEnabled   bool    `json:"commit_status_enabled"`
	CommitStatusContext   string  `json:"commit_status_context"`
	CommitStatusTargetURL *string `json:"commit_status_target_url"`
	BuildTrendKeepCount   int     `json:"build_trend_keep_count"`
	SessionTimeoutSeconds int     `json:"session_timeout_seconds"`
}

type apiLogRecord struct {
	At     string `json:"at"`
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
	Status int    `json:"status,omitempty"`
	Action string `json:"action,omitempty"`
	Result string `json:"result,omitempty"`
}

type apiConfigLogRecord struct {
	At      string         `json:"at"`
	Type    string         `json:"type"`
	Changes map[string]any `json:"changes"`
}

type apiBuildStatus struct {
	LastBlobSHA           *string `json:"last_blob_sha"`
	LastCommitSHA         *string `json:"last_commit_sha"`
	LastStartedAt         *string `json:"last_started_at"`
	LastFinishedAt        *string `json:"last_finished_at"`
	LastTargetStatus      *string `json:"last_target_status"`
	LastDeployStatus      *string `json:"last_deploy_status"`
	LastTrigger           *string `json:"last_trigger"`
	PendingTransfersCount int     `json:"pending_transfers_count"`
	NotifyPendingCount    int     `json:"notify_pending_count"`
	CircuitOpen           bool    `json:"circuit_open"`
	Running               bool    `json:"running"`
}

type apiHistoryRecord struct {
	ID              string   `json:"id"`
	BuildAt         string   `json:"build_at,omitempty"`
	StartedAt       string   `json:"started_at,omitempty"`
	FinishedAt      string   `json:"finished_at,omitempty"`
	SHA             *string  `json:"sha,omitempty"`
	CommitSHA       *string  `json:"commit_sha,omitempty"`
	BlobSHA         *string  `json:"blob_sha,omitempty"`
	Status          string   `json:"status"`
	Trigger         string   `json:"trigger,omitempty"`
	DurationSeconds int64    `json:"duration_seconds"`
	Flagged         bool     `json:"flagged"`
	Tags            []string `json:"tags"`
	OutputSizeBytes *int64   `json:"output_size_bytes"`
}

type apiBuildLog struct {
	ID              string         `json:"id"`
	StartedAt       string         `json:"started_at"`
	FinishedAt      string         `json:"finished_at"`
	TargetStatus    string         `json:"target_status"`
	Commit          map[string]any `json:"commit"`
	Pipeline        apiPipelineLog `json:"pipeline"`
	Warnings        []string       `json:"warnings"`
	DurationSeconds int64          `json:"duration_seconds"`
	Comment         *string        `json:"comment"`
	Flagged         bool           `json:"flagged"`
	Tags            []string       `json:"tags"`
}

type apiPipelineLog struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func NewAPIServer(cfg APIConfig) (*APIServer, error) {
	if cfg.Addr == "" {
		cfg.Addr = defaultAPIAddr
	}
	if cfg.StateDir == "" {
		return nil, errors.New("state directory is required")
	}
	if !filepath.IsAbs(cfg.StateDir) {
		return nil, fmt.Errorf("state directory must be absolute: %s", cfg.StateDir)
	}
	if cfg.Now == nil {
		cfg.Now = func() time.Time { return time.Now().UTC() }
	}
	if err := validateCredentials(filepath.Join(cfg.StateDir, ".admin_credentials")); err != nil {
		return nil, err
	}
	return &APIServer{cfg: cfg, startedAt: cfg.Now().UTC(), sessions: map[string]apiSession{}}, nil
}

func InitCredentials(stateDir, password string, now time.Time) error {
	if !filepath.IsAbs(stateDir) {
		return fmt.Errorf("state directory must be absolute: %s", stateDir)
	}
	if password == "" {
		return errors.New("password must not be empty")
	}
	path := filepath.Join(stateDir, ".admin_credentials")
	if _, err := os.Stat(path); err == nil {
		return errors.New("credentials already exist")
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	salt, err := randomHex(16)
	if err != nil {
		return err
	}
	at := now.UTC().Format(apiTimeLayout)
	cred := apiCredentials{
		PasswordHash: hashPassword(password, salt),
		Salt:         salt,
		Algorithm:    "sha256_iter_v1",
		Iterations:   passwordIterations,
		MustChange:   true,
		CreatedAt:    at,
		UpdatedAt:    at,
	}
	return atomicWriteJSON(path, cred, 0600)
}

func (s *APIServer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", s.handleLogin)
	mux.HandleFunc("/api/login/totp", s.handleLoginTOTP)
	mux.HandleFunc("/api/logout", s.withAuth(s.handleLogout))
	mux.HandleFunc("/api/change-password", s.withAuth(s.handleChangePassword))
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/sessions", s.withAuth(s.handleSessions))
	mux.HandleFunc("/api/sessions/revoke-all", s.withAuth(s.handleRevokeSessions))
	mux.HandleFunc("/api/config", s.withAuth(s.handleConfig))
	mux.HandleFunc("/api/config/validate", s.withAuth(s.handleConfigValidate))
	mux.HandleFunc("/api/access-log", s.withAuth(s.handleJSONLinesLog(".access_log", "log")))
	mux.HandleFunc("/api/api-access-log", s.withAuth(s.handleJSONLinesLog(".api_access_log", "log")))
	mux.HandleFunc("/api/config-log", s.withAuth(s.handleJSONLinesLog(".config_log", "log")))
	mux.HandleFunc("/api/status", s.withAuth(s.handleStatus))
	mux.HandleFunc("/api/sysinfo", s.withAuth(s.handleSysinfo))
	mux.HandleFunc("/api/stats", s.withAuth(s.handleStats))
	mux.HandleFunc("/api/stats/timeline", s.withAuth(s.handleStatsTimeline))
	mux.HandleFunc("/api/stats/build-duration", s.withAuth(s.handleStatsBuildDuration))
	mux.HandleFunc("/api/output-meta", s.withAuth(s.handleOutputMeta))
	mux.HandleFunc("/api/dashboard", s.withAuth(s.handleDashboard))
	mux.HandleFunc("/api/notify-log", s.withAuth(s.handleJSONLinesLog(".notify_log", "log")))
	mux.HandleFunc("/api/build", s.withAuth(s.handleBuild(false)))
	mux.HandleFunc("/api/build/force", s.withAuth(s.handleBuild(true)))
	mux.HandleFunc("/api/build/cancel", s.withAuth(s.handleCancel))
	mux.HandleFunc("/api/build/stream", s.withAuth(s.handleBuildStream))
	mux.HandleFunc("/api/logs", s.withAuth(s.handleLogs))
	mux.HandleFunc("/api/logs/search", s.withAuth(s.handleLogSearch))
	mux.HandleFunc("/api/logs/export", s.withAuth(s.handleLogExport))
	mux.HandleFunc("/api/history", s.withAuth(s.handleHistory))
	mux.HandleFunc("/api/history/export", s.withAuth(s.handleHistoryExport))
	mux.HandleFunc("/api/history/", s.withAuth(s.handleHistoryPath))
	mux.HandleFunc("/api/circuit-breaker/reset", s.withAuth(s.handleCircuitReset))
	mux.HandleFunc("/api/queue", s.withAuth(s.handleQueue))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		mux.ServeHTTP(rec, r)
		if strings.HasPrefix(r.URL.Path, "/api/") {
			_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".api_access_log"), apiLogRecord{
				At: s.nowString(), Method: r.Method, Path: r.URL.Path, Status: rec.status,
			})
		}
	})
}

func (s *APIServer) ListenAndServe() error {
	return http.ListenAndServe(s.cfg.Addr, s.Handler())
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if !decodeBody(w, r, &req, true) {
		return
	}
	if req.Password == "" {
		writeValidation(w, "password", "required")
		return
	}
	path := filepath.Join(s.cfg.StateDir, ".admin_credentials")
	cred, err := readCredentials(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	now := s.nowString()
	if cred.PasswordHash != hashPassword(req.Password, cred.Salt) {
		cred.FailedCount++
		cred.LastFailureAt = &now
		_ = atomicWriteJSON(path, cred, 0600)
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	cred.FailedCount = 0
	cred.LastLoginAt = &now
	if err := atomicWriteJSON(path, cred, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	token, err := randomToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	timeoutSeconds := 28800
	if cfg, err := s.readMergedConfig(); err == nil {
		timeoutSeconds = cfg.SessionTimeoutSeconds
	}
	expires := s.cfg.Now().UTC().Add(time.Duration(timeoutSeconds) * time.Second).Format(apiTimeLayout)
	s.mu.Lock()
	s.sessions[tokenHash(token)] = apiSession{TokenHash: tokenHash(token), CreatedAt: now, ExpiresAt: expires, LastUsedAt: now}
	s.mu.Unlock()
	mustChange := "none"
	if cred.MustChange {
		mustChange = "prompt"
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "must_change": mustChange})
}

func (s *APIServer) handleLoginTOTP(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var req struct {
		Ticket string `json:"ticket"`
		Code   string `json:"code"`
	}
	if !decodeBody(w, r, &req, true) {
		return
	}
	if req.Ticket == "" {
		writeValidation(w, "ticket", "required")
		return
	}
	if req.Code == "" {
		writeValidation(w, "code", "required")
		return
	}
	writeError(w, http.StatusUnauthorized, "Unauthorized")
}

func (s *APIServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	token := bearerToken(r)
	s.mu.Lock()
	delete(s.sessions, tokenHash(token))
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]string{"message": "Logged out"})
}

func (s *APIServer) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if !decodeBody(w, r, &req, true) {
		return
	}
	if req.CurrentPassword == "" {
		writeValidation(w, "current_password", "required")
		return
	}
	if len(req.NewPassword) < 12 {
		writeValidation(w, "new_password", "must be at least 12 characters")
		return
	}
	path := filepath.Join(s.cfg.StateDir, ".admin_credentials")
	cred, err := readCredentials(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if cred.PasswordHash != hashPassword(req.CurrentPassword, cred.Salt) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	salt, err := randomHex(16)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	now := s.nowString()
	cred.Salt = salt
	cred.PasswordHash = hashPassword(req.NewPassword, salt)
	cred.MustChange = false
	cred.UpdatedAt = now
	cred.FailedCount = 0
	if err := atomicWriteJSON(path, cred, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	currentHash := tokenHash(bearerToken(r))
	s.mu.Lock()
	for key := range s.sessions {
		if key != currentHash {
			delete(s.sessions, key)
		}
	}
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]string{"message": "Password changed"})
}

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	status := "ok"
	items := []map[string]any{
		{"name": "api", "status": "ok"},
	}
	if _, ok, err := readBuildStatus(filepath.Join(s.cfg.StateDir, ".build_status.json")); err != nil {
		status = "degraded"
		items = append(items, map[string]any{"name": "build_status", "status": "error"})
	} else if !ok {
		status = "degraded"
		items = append(items, map[string]any{"name": "build_status", "status": "warn"})
	} else {
		items = append(items, map[string]any{"name": "build_status", "status": "ok"})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":     status,
		"checked_at": s.nowString(),
		"items":      items,
	})
}

func (s *APIServer) handleSessions(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	current := tokenHash(bearerToken(r))
	s.mu.Lock()
	defer s.mu.Unlock()
	sessions := []map[string]any{}
	for key, session := range s.sessions {
		sessions = append(sessions, map[string]any{
			"created_at":   session.CreatedAt,
			"expires_at":   session.ExpiresAt,
			"last_used_at": session.LastUsedAt,
			"current":      key == current,
		})
	}
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i]["created_at"].(string) > sessions[j]["created_at"].(string)
	})
	writeJSON(w, http.StatusOK, map[string]any{"sessions": sessions})
}

func (s *APIServer) handleRevokeSessions(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	current := tokenHash(bearerToken(r))
	revoked := 0
	s.mu.Lock()
	for key := range s.sessions {
		if key != current {
			delete(s.sessions, key)
			revoked++
		}
	}
	s.mu.Unlock()
	_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".access_log"), apiLogRecord{At: s.nowString(), Action: "sessions_revoke_all", Result: "success"})
	writeJSON(w, http.StatusOK, map[string]any{"message": "All other sessions revoked", "revoked_count": revoked})
}

func (s *APIServer) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		cfg, err := s.readMergedConfig()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusOK, cfg)
	case http.MethodPost:
		var patch map[string]any
		if !decodeBody(w, r, &patch, true) {
			return
		}
		cfg, err := s.readMergedConfig()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		next, changed, ok := validateConfigPatch(w, cfg, patch)
		if !ok {
			return
		}
		if changed {
			if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".server_config"), next, 0600); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "server_config", Changes: patch})
			writeJSON(w, http.StatusOK, map[string]any{"message": "Config updated", "config": next})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"message": "No changes", "config": next})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleConfigValidate(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var patch map[string]any
	if !decodeBody(w, r, &patch, true) {
		return
	}
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	next, _, ok := validateConfigPatch(w, cfg, patch)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"valid": true, "config": next, "warnings": []string{}, "errors": []string{}})
}

func (s *APIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	resp, ok := s.statusPayload(w)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *APIServer) statusPayload(w http.ResponseWriter) (map[string]any, bool) {
	state, err := readBuildState(filepath.Join(s.cfg.StateDir, ".build_state"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return nil, false
	}
	circuit, err := readCircuitState(filepath.Join(s.cfg.StateDir, ".build_circuit_state"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return nil, false
	}
	statusPath := filepath.Join(s.cfg.StateDir, ".build_status.json")
	if st, ok, err := readBuildStatus(statusPath); err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return nil, false
	} else if ok {
		running := st.Running || state.Running || lockExists(filepath.Join(s.cfg.StateDir, ".build_lock"))
		lastSHA := st.LastBlobSHA
		if lastSHA == nil {
			lastSHA = st.LastCommitSHA
		}
		return map[string]any{
			"last_sha":                lastSHA,
			"last_build_at":           coalesceString(st.LastFinishedAt, st.LastStartedAt),
			"last_build_status":       stringOr(st.LastTargetStatus, "none"),
			"last_trigger":            st.LastTrigger,
			"last_deploy_status":      st.LastDeployStatus,
			"pending_transfers_count": st.PendingTransfersCount,
			"notify_pending_count":    st.NotifyPendingCount,
			"circuit_open":            st.CircuitOpen,
			"output_url":              nil,
			"running":                 running,
		}, true
	}
	history := readHistory(filepath.Join(s.cfg.StateDir, ".build_history"))
	var latest *apiHistoryRecord
	if len(history) > 0 {
		latest = &history[0]
	}
	resp := map[string]any{
		"last_sha":                nil,
		"last_build_at":           nil,
		"last_build_status":       "none",
		"last_trigger":            nil,
		"last_deploy_status":      nil,
		"pending_transfers_count": 0,
		"notify_pending_count":    0,
		"circuit_open":            circuit.Open,
		"output_url":              nil,
		"running":                 state.Running || lockExists(filepath.Join(s.cfg.StateDir, ".build_lock")),
	}
	if latest != nil {
		resp["last_sha"] = firstNonNil(latest.BlobSHA, latest.CommitSHA, latest.SHA)
		resp["last_build_at"] = firstNonEmpty(latest.FinishedAt, latest.StartedAt, latest.BuildAt)
		resp["last_build_status"] = latest.Status
		if latest.Trigger != "" {
			resp["last_trigger"] = latest.Trigger
		}
	}
	return resp, true
}

func (s *APIServer) handleSysinfo(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	meta, err := inspectOutput(s.outputDir())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	uptime := int64(s.cfg.Now().UTC().Sub(s.startedAt).Seconds())
	if uptime < 0 {
		uptime = 0
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"uptime_seconds":    uptime,
		"state_dir":         s.cfg.StateDir,
		"output_exists":     meta.Exists,
		"output_size_bytes": meta.SizeBytes,
		"output_mtime":      meta.MTime,
	})
}

func (s *APIServer) handleStats(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	days, ok := parseBoundedInt(w, r, "days", 7, 1, 366)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, s.statsSummary(days))
}

func (s *APIServer) handleStatsTimeline(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	days, ok := parseBoundedInt(w, r, "days", 30, 1, 366)
	if !ok {
		return
	}
	records := s.historyWithinDays(days)
	type bucket struct {
		Date    string `json:"date"`
		Success int    `json:"success"`
		Failure int    `json:"failure"`
		Total   int    `json:"total"`
	}
	byDate := map[string]*bucket{}
	for _, record := range records {
		at := historyRecordTime(record)
		if at.IsZero() {
			continue
		}
		date := at.Format("2006-01-02")
		if byDate[date] == nil {
			byDate[date] = &bucket{Date: date}
		}
		byDate[date].Total++
		switch statusCategory(record.Status) {
		case "success":
			byDate[date].Success++
		case "failure":
			byDate[date].Failure++
		}
	}
	keys := make([]string, 0, len(byDate))
	for date := range byDate {
		keys = append(keys, date)
	}
	sort.Strings(keys)
	timeline := []bucket{}
	for _, date := range keys {
		timeline = append(timeline, *byDate[date])
	}
	writeJSON(w, http.StatusOK, map[string]any{"days": days, "timeline": timeline})
}

func (s *APIServer) handleStatsBuildDuration(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	n, ok := parseBoundedInt(w, r, "n", 10, 1, 1000)
	if !ok {
		return
	}
	durations := []int64{}
	recent := []map[string]any{}
	for _, log := range s.readBuildLogsNewest() {
		if log.DurationSeconds <= 0 {
			continue
		}
		if len(durations) >= n {
			break
		}
		durations = append(durations, log.DurationSeconds)
		recent = append(recent, map[string]any{
			"id":               log.ID,
			"build_at":         firstNonEmpty(log.FinishedAt, log.StartedAt),
			"duration_seconds": log.DurationSeconds,
			"status":           log.TargetStatus,
		})
	}
	var sum, min, max int64
	for i, duration := range durations {
		sum += duration
		if i == 0 || duration < min {
			min = duration
		}
		if duration > max {
			max = duration
		}
	}
	avg := 0.0
	if len(durations) > 0 {
		avg = float64(sum) / float64(len(durations))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"n": n, "count": len(durations), "avg_seconds": avg,
		"min_seconds": min, "max_seconds": max, "recent": recent,
	})
}

func (s *APIServer) handleOutputMeta(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	meta, err := inspectOutput(s.outputDir())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if !meta.Exists {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	logs := s.readBuildLogsNewest()
	var latest *apiBuildLog
	if len(logs) > 0 {
		latest = &logs[0]
	}
	resp := map[string]any{
		"size_bytes":        meta.SizeBytes,
		"mtime":             meta.MTime,
		"sha256":            meta.SHA256,
		"heading_count":     nil,
		"tables_count":      nil,
		"code_blocks_count": nil,
		"size_diff_bytes":   nil,
		"size_warn":         false,
		"build_warnings":    []string{},
		"build_id":          "",
		"commit_sha":        nil,
		"build_at":          "",
	}
	if latest != nil {
		resp["build_warnings"] = latest.Warnings
		resp["build_id"] = latest.ID
		resp["commit_sha"] = firstCommitSHA(latest.Commit)
		resp["build_at"] = firstNonEmpty(latest.FinishedAt, latest.StartedAt)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *APIServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	status, ok := s.statusPayload(w)
	if !ok {
		return
	}
	meta, err := inspectOutput(s.outputDir())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	uptime := int64(s.cfg.Now().UTC().Sub(s.startedAt).Seconds())
	if uptime < 0 {
		uptime = 0
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status": status,
		"sysinfo": map[string]any{
			"uptime_seconds":    uptime,
			"output_exists":     meta.Exists,
			"output_size_bytes": meta.SizeBytes,
			"output_mtime":      meta.MTime,
		},
		"stats":    s.statsSummary(7),
		"schedule": map[string]any{"enabled": false, "next_run_at": nil},
		"alerts":   map[string]any{"notify_pending_count": status["notify_pending_count"], "circuit_open": status["circuit_open"]},
	})
}

func (s *APIServer) statsSummary(days int) map[string]any {
	records := s.historyWithinDays(days)
	success := 0
	failure := 0
	for _, record := range records {
		switch statusCategory(record.Status) {
		case "success":
			success++
		case "failure":
			failure++
		}
	}
	latestStatus := "none"
	if len(records) > 0 {
		latestStatus = records[0].Status
	}
	return map[string]any{
		"days": days, "total": len(records), "success": success,
		"failure": failure, "latest_status": latestStatus,
	}
}

func (s *APIServer) historyWithinDays(days int) []apiHistoryRecord {
	records := readHistory(filepath.Join(s.cfg.StateDir, ".build_history"))
	since := s.cfg.Now().UTC().AddDate(0, 0, -days)
	out := []apiHistoryRecord{}
	for _, record := range records {
		at := historyRecordTime(record)
		if at.IsZero() || !at.Before(since) {
			out = append(out, record)
		}
	}
	return out
}

func (s *APIServer) outputDir() string {
	return filepath.Join(s.cfg.StateDir, "site")
}

func (s *APIServer) handleBuild(force bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
			return
		}
		circuit, err := readCircuitState(filepath.Join(s.cfg.StateDir, ".build_circuit_state"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if circuit.Open {
			writeError(w, http.StatusConflict, "circuit_open")
			return
		}
		if maintenanceEnabled(filepath.Join(s.cfg.StateDir, ".maintenance")) {
			writeError(w, http.StatusServiceUnavailable, "maintenance")
			return
		}
		statePath := filepath.Join(s.cfg.StateDir, ".build_state")
		state, err := readBuildState(statePath)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if state.Running || lockExists(filepath.Join(s.cfg.StateDir, ".build_lock")) {
			maxSize := s.queueMaxSize()
			if maxSize == 0 || len(state.Queued) >= maxSize {
				writeError(w, http.StatusTooManyRequests, "queue_full")
				return
			}
			entry := map[string]any{
				"id":           s.newID("q"),
				"trigger":      "manual",
				"queued_at":    s.nowString(),
				"requested_by": "admin",
				"priority":     "normal",
				"created_seq":  nextCreatedSeq(state.Queued),
				"payload":      map[string]any{"force": force},
			}
			state.Queued = append(state.Queued, entry)
			if err := atomicWriteJSON(statePath, state, 0600); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			writeJSON(w, http.StatusAccepted, map[string]any{"message": "Build queued", "queued": true})
			return
		}
		buildID := s.newID("b")
		state.Running = true
		state.CurrentBuildID = &buildID
		now := s.nowString()
		state.LastStartedAt = &now
		if force {
			_ = atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".last_sha"), map[string]string{"sha": ""}, 0600)
		}
		if err := atomicWriteJSON(statePath, state, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]any{"message": "Build started", "build_id": buildID, "queued": false})
	}
}

func (s *APIServer) handleCancel(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	statePath := filepath.Join(s.cfg.StateDir, ".build_state")
	state, err := readBuildState(statePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if !state.Running && !lockExists(filepath.Join(s.cfg.StateDir, ".build_lock")) {
		writeError(w, http.StatusConflict, "No build is running")
		return
	}
	state.Running = false
	state.CurrentBuildID = nil
	now := s.nowString()
	state.LastFinishedAt = &now
	if err := atomicWriteJSON(statePath, state, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Build cancelled"})
}

func (s *APIServer) handleBuildStream(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	for _, line := range s.allLogLines() {
		payload, err := json.Marshal(map[string]string{"type": "log", "line": line})
		if err != nil {
			continue
		}
		fmt.Fprintf(w, "data: %s\n\n", payload)
	}
}

func (s *APIServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	n, ok := parseBoundedInt(w, r, "n", 100, 1, 1000)
	if !ok {
		return
	}
	q := r.URL.Query().Get("q")
	lines := s.allLogLines()
	if q != "" {
		lines = filterContains(lines, q)
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	writeJSON(w, http.StatusOK, map[string]any{"lines": lines})
}

func (s *APIServer) handleLogExport(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"exported_at": s.nowString(), "lines": s.allLogLines()})
}

func (s *APIServer) handleLogSearch(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	q := r.URL.Query().Get("q")
	if len(q) > 500 {
		writeValidation(w, "q", "must be 500 characters or less")
		return
	}
	level := strings.ToUpper(r.URL.Query().Get("level"))
	if level == "WARN" {
		level = "WARNING"
	}
	if level != "" && level != "INFO" && level != "WARNING" && level != "ERROR" && level != "DEBUG" {
		writeValidation(w, "level", "invalid value")
		return
	}
	results := []map[string]any{}
	for _, log := range s.readBuildLogsNewest() {
		lines := append(splitLines(log.Pipeline.Stdout), splitLines(log.Pipeline.Stderr)...)
		lines = append(lines, log.Warnings...)
		lines = filterContains(lines, q)
		if level != "" {
			lines = filterContains(lines, "["+level+"]")
		}
		if len(lines) == 0 {
			continue
		}
		results = append(results, map[string]any{"id": log.ID, "build_at": firstNonEmpty(log.FinishedAt, log.StartedAt), "lines": lines})
	}
	writeJSON(w, http.StatusOK, map[string]any{"query": q, "from": r.URL.Query().Get("from"), "to": r.URL.Query().Get("to"), "results": results})
}

func (s *APIServer) handleHistory(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	page, ok := parseBoundedInt(w, r, "page", 1, 1, 1_000_000)
	if !ok {
		return
	}
	perPage, ok := parseBoundedInt(w, r, "per_page", 20, 1, 100)
	if !ok {
		return
	}
	history := readHistory(filepath.Join(s.cfg.StateDir, ".build_history"))
	total := len(history)
	pages := 0
	if total > 0 {
		pages = (total + perPage - 1) / perPage
	}
	start := (page - 1) * perPage
	end := start + perPage
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	writeJSON(w, http.StatusOK, map[string]any{"total": total, "page": page, "per_page": perPage, "pages": pages, "history": history[start:end]})
}

func (s *APIServer) handleHistoryExport(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"exported_at": s.nowString(), "history": readHistory(filepath.Join(s.cfg.StateDir, ".build_history"))})
}

func (s *APIServer) handleHistoryPath(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/api/history/")
	id, suffix, ok := strings.Cut(rest, "/")
	if !ok || id == "" {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	log, err := readBuildLog(filepath.Join(s.cfg.StateDir, ".build_logs", id+".json"))
	if errors.Is(err, os.ErrNotExist) {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	switch suffix {
	case "log":
		lines := append(splitLines(log.Pipeline.Stdout), splitLines(log.Pipeline.Stderr)...)
		lines = append(lines, log.Warnings...)
		writeJSON(w, http.StatusOK, map[string]any{
			"id": log.ID, "build_at": firstNonEmpty(log.FinishedAt, log.StartedAt),
			"sha": firstCommitSHA(log.Commit), "status": log.TargetStatus, "comment": log.Comment,
			"flagged": log.Flagged, "tags": log.Tags, "lines": lines,
		})
	case "comment":
		writeJSON(w, http.StatusOK, map[string]any{"id": log.ID, "comment": log.Comment})
	default:
		writeError(w, http.StatusNotFound, "Not found")
	}
}

func (s *APIServer) handleCircuitReset(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	state := apiCircuitState{Open: false, ConsecutiveFailures: 0}
	if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".build_circuit_state"), state, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Circuit breaker reset", "open": false, "consecutive_failures": 0})
}

func (s *APIServer) handleQueue(w http.ResponseWriter, r *http.Request) {
	statePath := filepath.Join(s.cfg.StateDir, ".build_state")
	state, err := readBuildState(statePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"queued": state.Queued, "max_size": s.queueMaxSize()})
	case http.MethodDelete:
		if !rejectBody(w, r) {
			return
		}
		cleared := len(state.Queued)
		state.Queued = []map[string]any{}
		if err := atomicWriteJSON(statePath, state, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"message": "Queue cleared", "cleared_count": cleared})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleJSONLinesLog(filename, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
			return
		}
		limit, ok := parseBoundedInt(w, r, "limit", 100, 1, 1000)
		if !ok {
			return
		}
		offset, ok := parseBoundedInt(w, r, "offset", 0, 0, 1_000_000)
		if !ok {
			return
		}
		records := readJSONLines(filepath.Join(s.cfg.StateDir, filename))
		sort.Slice(records, func(i, j int) bool {
			return fmt.Sprint(records[i]["at"]) > fmt.Sprint(records[j]["at"])
		})
		total := len(records)
		start := offset
		if start > total {
			start = total
		}
		end := start + limit
		if end > total {
			end = total
		}
		writeJSON(w, http.StatusOK, map[string]any{key: records[start:end], "total": total})
	}
}

func (s *APIServer) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if !s.validSession(token) {
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		s.touchSession(token)
		next(w, r)
	}
}

func bearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}

func (s *APIServer) validSession(token string) bool {
	if token == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[tokenHash(token)]
	if !ok {
		return false
	}
	expires, err := time.Parse(apiTimeLayout, session.ExpiresAt)
	if err != nil || !s.cfg.Now().UTC().Before(expires) {
		delete(s.sessions, tokenHash(token))
		return false
	}
	return true
}

func (s *APIServer) touchSession(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := tokenHash(token)
	session, ok := s.sessions[key]
	if !ok {
		return
	}
	session.LastUsedAt = s.nowString()
	s.sessions[key] = session
}

func (s *APIServer) queueMaxSize() int {
	cfg, err := s.readMergedConfig()
	if err != nil {
		return defaultQueueMaxSize
	}
	return cfg.QueueMaxSize
}

func (s *APIServer) readMergedConfig() (apiServerConfig, error) {
	cfg := defaultServerConfig()
	path := filepath.Join(s.cfg.StateDir, ".server_config")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	var patch map[string]json.RawMessage
	if err := json.Unmarshal(data, &patch); err != nil {
		return cfg, err
	}
	var fileCfg apiServerConfig
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return cfg, err
	}
	fileCfg = normalizeServerConfig(fileCfg)
	if _, ok := patch["log_retention_days"]; ok && fileCfg.LogRetentionDays == defaultServerConfig().LogRetentionDays {
		var v int
		if json.Unmarshal(patch["log_retention_days"], &v) == nil {
			fileCfg.LogRetentionDays = v
		}
	}
	if _, ok := patch["snapshots_keep"]; ok && fileCfg.SnapshotsKeep == defaultServerConfig().SnapshotsKeep {
		var v int
		if json.Unmarshal(patch["snapshots_keep"], &v) == nil {
			fileCfg.SnapshotsKeep = v
		}
	}
	if _, ok := patch["queue_max_size"]; ok && fileCfg.QueueMaxSize == defaultServerConfig().QueueMaxSize {
		var v int
		if json.Unmarshal(patch["queue_max_size"], &v) == nil {
			fileCfg.QueueMaxSize = v
		}
	}
	return fileCfg, nil
}

func (s *APIServer) nowString() string {
	return s.cfg.Now().UTC().Format(apiTimeLayout)
}

func (s *APIServer) newID(prefix string) string {
	return prefix + s.cfg.Now().UTC().Format("20060102150405")
}

func (s *APIServer) readBuildLogsNewest() []apiBuildLog {
	dir := filepath.Join(s.cfg.StateDir, ".build_logs")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	logs := []apiBuildLog{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		log, err := readBuildLog(filepath.Join(dir, entry.Name()))
		if err == nil {
			logs = append(logs, log)
		}
	}
	sort.Slice(logs, func(i, j int) bool {
		return firstNonEmpty(logs[i].FinishedAt, logs[i].StartedAt, logs[i].ID) > firstNonEmpty(logs[j].FinishedAt, logs[j].StartedAt, logs[j].ID)
	})
	return logs
}

type outputInspection struct {
	Exists    bool
	SizeBytes int64
	MTime     any
	SHA256    string
}

func inspectOutput(root string) (outputInspection, error) {
	info, err := os.Stat(root)
	if errors.Is(err, os.ErrNotExist) {
		return outputInspection{Exists: false, MTime: nil}, nil
	}
	if err != nil {
		return outputInspection{}, err
	}
	files := []string{}
	if !info.IsDir() {
		files = append(files, root)
	} else if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return outputInspection{}, err
	}
	sort.Strings(files)
	var total int64
	var latest time.Time
	manifest := sha256.New()
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return outputInspection{}, err
		}
		total += info.Size()
		if info.ModTime().After(latest) {
			latest = info.ModTime()
		}
		rel := filepath.Base(path)
		if info.IsDir() {
			continue
		}
		if r, err := filepath.Rel(root, path); err == nil {
			rel = filepath.ToSlash(r)
		}
		fileHash, err := fileSHA256(path)
		if err != nil {
			return outputInspection{}, err
		}
		_, _ = io.WriteString(manifest, rel)
		_, _ = io.WriteString(manifest, "\n")
		_, _ = io.WriteString(manifest, fileHash)
		_, _ = io.WriteString(manifest, "\n")
	}
	var mtime any
	if !latest.IsZero() {
		mtime = latest.UTC().Format(apiTimeLayout)
	}
	return outputInspection{
		Exists: true, SizeBytes: total, MTime: mtime,
		SHA256: hex.EncodeToString(manifest.Sum(nil)),
	}, nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s *APIServer) allLogLines() []string {
	lines := []string{}
	logs := s.readBuildLogsNewest()
	for i := len(logs) - 1; i >= 0; i-- {
		log := logs[i]
		lines = append(lines, splitLines(log.Pipeline.Stdout)...)
		lines = append(lines, splitLines(log.Pipeline.Stderr)...)
		lines = append(lines, log.Warnings...)
	}
	return lines
}

func validateCredentials(path string) error {
	_, err := readCredentials(path)
	return err
}

func readCredentials(path string) (apiCredentials, error) {
	var cred apiCredentials
	if err := readStrictJSONFile(path, &cred); err != nil {
		return cred, err
	}
	if cred.Algorithm != "sha256_iter_v1" || cred.Iterations != passwordIterations || !isLowerHex(cred.PasswordHash, 64) || !isLowerHex(cred.Salt, 32) {
		return cred, errors.New("invalid credentials")
	}
	return cred, nil
}

func hashPassword(password, saltHex string) string {
	salt, _ := hex.DecodeString(saltHex)
	digest := sha256.Sum256(append(salt, []byte(password)...))
	current := digest[:]
	for i := 1; i < passwordIterations; i++ {
		buf := make([]byte, 0, len(current)+len(salt)+len(password))
		buf = append(buf, current...)
		buf = append(buf, salt...)
		buf = append(buf, []byte(password)...)
		next := sha256.Sum256(buf)
		current = next[:]
	}
	return hex.EncodeToString(current)
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func randomHex(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func randomToken() (string, error) {
	value, err := randomHex(32)
	if err != nil {
		return "", err
	}
	return "acs_" + value, nil
}

func readBuildState(path string) (apiBuildState, error) {
	state := apiBuildState{Queued: []map[string]any{}}
	if err := readJSONIfExists(path, &state); err != nil {
		return state, err
	}
	if state.Queued == nil {
		state.Queued = []map[string]any{}
	}
	return state, nil
}

func readCircuitState(path string) (apiCircuitState, error) {
	var state apiCircuitState
	return state, readJSONIfExists(path, &state)
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
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record apiHistoryRecord
		if json.Unmarshal([]byte(line), &record) == nil {
			records = append(records, record)
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return firstNonEmpty(records[i].FinishedAt, records[i].StartedAt, records[i].BuildAt, records[i].ID) > firstNonEmpty(records[j].FinishedAt, records[j].StartedAt, records[j].BuildAt, records[j].ID)
	})
	return records
}

func readBuildLog(path string) (apiBuildLog, error) {
	var log apiBuildLog
	err := readJSONFile(path, &log)
	return log, err
}

func defaultServerConfig() apiServerConfig {
	return apiServerConfig{
		LogMaxLines: 500, HistoryMaxCount: 100, BuildTimeoutSeconds: 300,
		LogRetentionDays: 30, LogArchiveAfterDays: 0, LogLevel: "INFO",
		SnapshotsKeep: 5, QueueMaxSize: defaultQueueMaxSize,
		BuildRetryMax: 0, BuildRetryBaseSeconds: 5,
		CommitStatusEnabled: false, CommitStatusContext: "Adlaire CI",
		BuildTrendKeepCount: 1000, SessionTimeoutSeconds: 28800,
	}
}

func normalizeServerConfig(cfg apiServerConfig) apiServerConfig {
	def := defaultServerConfig()
	if cfg.LogMaxLines == 0 {
		cfg.LogMaxLines = def.LogMaxLines
	}
	if cfg.HistoryMaxCount == 0 {
		cfg.HistoryMaxCount = def.HistoryMaxCount
	}
	if cfg.BuildTimeoutSeconds == 0 {
		cfg.BuildTimeoutSeconds = def.BuildTimeoutSeconds
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = def.LogLevel
	}
	if cfg.SnapshotsKeep == 0 {
		cfg.SnapshotsKeep = def.SnapshotsKeep
	}
	if cfg.QueueMaxSize == 0 {
		cfg.QueueMaxSize = def.QueueMaxSize
	}
	if cfg.BuildRetryBaseSeconds == 0 {
		cfg.BuildRetryBaseSeconds = def.BuildRetryBaseSeconds
	}
	if cfg.CommitStatusContext == "" {
		cfg.CommitStatusContext = def.CommitStatusContext
	}
	if cfg.BuildTrendKeepCount == 0 {
		cfg.BuildTrendKeepCount = def.BuildTrendKeepCount
	}
	if cfg.SessionTimeoutSeconds == 0 {
		cfg.SessionTimeoutSeconds = def.SessionTimeoutSeconds
	}
	return cfg
}

func validateConfigPatch(w http.ResponseWriter, cfg apiServerConfig, patch map[string]any) (apiServerConfig, bool, bool) {
	data, err := json.Marshal(patch)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return cfg, false, false
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return cfg, false, false
	}
	changed := false
	for key, value := range raw {
		switch key {
		case "log_max_lines":
			var v int
			if !decodeIntField(w, key, value, 1, 10000, &v) {
				return cfg, false, false
			}
			changed = changed || cfg.LogMaxLines != v
			cfg.LogMaxLines = v
		case "history_max_count":
			var v int
			if !decodeIntField(w, key, value, 1, 10000, &v) {
				return cfg, false, false
			}
			changed = changed || cfg.HistoryMaxCount != v
			cfg.HistoryMaxCount = v
		case "build_timeout_seconds":
			var v int
			if !decodeIntField(w, key, value, 1, 86400, &v) {
				return cfg, false, false
			}
			changed = changed || cfg.BuildTimeoutSeconds != v
			cfg.BuildTimeoutSeconds = v
		case "log_retention_days", "log_archive_after_days":
			var v int
			if !decodeIntField(w, key, value, 0, 3650, &v) {
				return cfg, false, false
			}
			if key == "log_retention_days" {
				changed = changed || cfg.LogRetentionDays != v
				cfg.LogRetentionDays = v
			} else {
				changed = changed || cfg.LogArchiveAfterDays != v
				cfg.LogArchiveAfterDays = v
			}
		case "log_level":
			var v string
			if json.Unmarshal(value, &v) != nil || (v != "INFO" && v != "DEBUG" && v != "WARNING" && v != "ERROR") {
				writeValidation(w, key, "invalid value")
				return cfg, false, false
			}
			changed = changed || cfg.LogLevel != v
			cfg.LogLevel = v
		case "pat_expires_at", "commit_status_target_url":
			var v *string
			if json.Unmarshal(value, &v) != nil {
				writeValidation(w, key, "invalid value")
				return cfg, false, false
			}
			if key == "pat_expires_at" {
				changed = changed || stringPtrValue(cfg.PATExpiresAt) != stringPtrValue(v)
				cfg.PATExpiresAt = v
			} else {
				changed = changed || stringPtrValue(cfg.CommitStatusTargetURL) != stringPtrValue(v)
				cfg.CommitStatusTargetURL = v
			}
		case "snapshots_keep", "queue_max_size":
			var v int
			if !decodeIntField(w, key, value, 0, 100, &v) {
				return cfg, false, false
			}
			if key == "snapshots_keep" {
				changed = changed || cfg.SnapshotsKeep != v
				cfg.SnapshotsKeep = v
			} else {
				changed = changed || cfg.QueueMaxSize != v
				cfg.QueueMaxSize = v
			}
		case "build_retry_max":
			var v int
			if !decodeIntField(w, key, value, 0, 10, &v) {
				return cfg, false, false
			}
			changed = changed || cfg.BuildRetryMax != v
			cfg.BuildRetryMax = v
		case "build_retry_base_seconds":
			var v int
			if !decodeIntField(w, key, value, 1, 3600, &v) {
				return cfg, false, false
			}
			changed = changed || cfg.BuildRetryBaseSeconds != v
			cfg.BuildRetryBaseSeconds = v
		case "commit_status_enabled":
			var v bool
			if json.Unmarshal(value, &v) != nil {
				writeValidation(w, key, "invalid value")
				return cfg, false, false
			}
			changed = changed || cfg.CommitStatusEnabled != v
			cfg.CommitStatusEnabled = v
		case "commit_status_context":
			var v string
			if json.Unmarshal(value, &v) != nil || len(v) < 1 || len(v) > 100 {
				writeValidation(w, key, "invalid value")
				return cfg, false, false
			}
			changed = changed || cfg.CommitStatusContext != v
			cfg.CommitStatusContext = v
		case "build_trend_keep_count":
			var v int
			if !decodeIntField(w, key, value, 10, 10000, &v) {
				return cfg, false, false
			}
			changed = changed || cfg.BuildTrendKeepCount != v
			cfg.BuildTrendKeepCount = v
		case "session_timeout_seconds":
			var v int
			if !decodeIntField(w, key, value, 300, 2592000, &v) {
				return cfg, false, false
			}
			changed = changed || cfg.SessionTimeoutSeconds != v
			cfg.SessionTimeoutSeconds = v
		default:
			writeValidation(w, key, "unknown key")
			return cfg, false, false
		}
	}
	return cfg, changed, true
}

func decodeIntField(w http.ResponseWriter, key string, raw json.RawMessage, min, max int, out *int) bool {
	if err := json.Unmarshal(raw, out); err != nil || *out < min || *out > max {
		writeValidation(w, key, "out of range")
		return false
	}
	return true
}

func appendJSONLine(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func readJSONLines(path string) []map[string]any {
	data, err := os.ReadFile(path)
	if err != nil {
		return []map[string]any{}
	}
	records := []map[string]any{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record map[string]any
		if json.Unmarshal([]byte(line), &record) == nil {
			records = append(records, record)
		}
	}
	return records
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func readJSONIfExists(path string, out any) error {
	err := readJSONFile(path, out)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func readJSONFile(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func readStrictJSONFile(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	return dec.Decode(out)
}

func atomicWriteJSON(path string, value any, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func method(w http.ResponseWriter, r *http.Request, want string) bool {
	if r.Method != want {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return false
	}
	return true
}

func rejectBody(w http.ResponseWriter, r *http.Request) bool {
	if r.Body == nil || r.Body == http.NoBody {
		return true
	}
	data, _ := io.ReadAll(r.Body)
	if len(strings.TrimSpace(string(data))) > 0 {
		writeError(w, http.StatusBadRequest, "Request body is not allowed")
		return false
	}
	return true
}

func decodeBody(w http.ResponseWriter, r *http.Request, out any, required bool) bool {
	if r.Body == nil || r.Body == http.NoBody {
		if required {
			writeError(w, http.StatusBadRequest, "Invalid JSON")
			return false
		}
		return true
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(out); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeValidation(w http.ResponseWriter, field, message string) {
	writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
		"error":   "Validation failed",
		"details": []map[string]string{{"field": field, "message": message}},
	})
}

func parseBoundedInt(w http.ResponseWriter, r *http.Request, key string, def, min, max int) (int, bool) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return def, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < min || value > max {
		writeValidation(w, key, "out of range")
		return 0, false
	}
	return value, true
}

func lockExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func maintenanceEnabled(path string) bool {
	var value struct {
		Enabled bool `json:"enabled"`
	}
	return readJSONIfExists(path, &value) == nil && value.Enabled
}

func isLowerHex(value string, length int) bool {
	if len(value) != length {
		return false
	}
	for _, r := range value {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}

func coalesceString(values ...*string) any {
	for _, value := range values {
		if value != nil && *value != "" {
			return *value
		}
	}
	return nil
}

func stringOr(value *string, fallback string) string {
	if value == nil || *value == "" {
		return fallback
	}
	return *value
}

func firstNonNil(values ...*string) any {
	for _, value := range values {
		if value != nil && *value != "" {
			return *value
		}
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func historyRecordTime(record apiHistoryRecord) time.Time {
	for _, value := range []string{record.FinishedAt, record.StartedAt, record.BuildAt, record.ID} {
		if value == "" {
			continue
		}
		if t, err := time.Parse(apiTimeLayout, value); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

func statusCategory(status string) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	switch {
	case normalized == "success" || normalized == "succeeded" || normalized == "passed":
		return "success"
	case normalized == "failure" || normalized == "failed" || normalized == "error" || strings.Contains(normalized, "fail"):
		return "failure"
	default:
		return "other"
	}
}

func splitLines(value string) []string {
	lines := []string{}
	for _, line := range strings.Split(value, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func filterContains(lines []string, q string) []string {
	if q == "" {
		return lines
	}
	out := []string{}
	for _, line := range lines {
		if strings.Contains(line, q) {
			out = append(out, line)
		}
	}
	return out
}

func nextCreatedSeq(queue []map[string]any) int {
	maxSeq := 0
	for _, entry := range queue {
		switch v := entry["created_seq"].(type) {
		case int:
			if v > maxSeq {
				maxSeq = v
			}
		case float64:
			if int(v) > maxSeq {
				maxSeq = int(v)
			}
		}
	}
	return maxSeq + 1
}

func firstCommitSHA(commit map[string]any) any {
	for _, key := range []string{"sha", "commit_sha"} {
		if value, ok := commit[key].(string); ok && value != "" {
			return value
		}
	}
	return nil
}
