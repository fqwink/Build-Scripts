package components

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
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
	cfg          APIConfig
	startedAt    time.Time
	mu           sync.Mutex
	sessions     map[string]apiSession
	loginTickets map[string]apiLoginTicket
	totpSetup    *apiPendingTOTP
}

type apiSession struct {
	TokenHash  string
	CreatedAt  string
	ExpiresAt  string
	LastUsedAt string
}

type apiLoginTicket struct {
	MustChange string
	ExpiresAt  time.Time
}

type apiPendingTOTP struct {
	SecretBase32 string
	ExpiresAt    time.Time
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

type apiTOTPSecret struct {
	Enabled          bool    `json:"enabled"`
	SecretBase32     *string `json:"secret_base32"`
	ConfirmedAt      *string `json:"confirmed_at"`
	LastAcceptedStep *int64  `json:"last_accepted_step"`
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

type apiMaintenanceState struct {
	Enabled bool    `json:"enabled"`
	Reason  *string `json:"reason"`
	Since   *string `json:"since"`
}

type apiRepoConfig struct {
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	Branch     string `json:"branch"`
	TargetFile string `json:"target_file"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type apiBranchConfigFile struct {
	BranchTargets []apiBranchTarget `json:"branch_targets"`
}

type apiBranchTarget struct {
	Branch        string            `json:"branch"`
	TargetFile    string            `json:"target_file"`
	SHAFile       string            `json:"sha_file"`
	Src           string            `json:"src"`
	Out           string            `json:"out"`
	DeployTargets []apiDeployTarget `json:"deploy_targets"`
}

type apiDeployTarget struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	DestDir string `json:"dest_dir"`
}

type apiAllowedHours struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type apiNotifyConfig struct {
	Webhooks []apiNotifyWebhook `json:"webhooks"`
	Channels []apiNotifyChannel `json:"channels"`
	On       []string           `json:"on"`
	Summary  apiNotifySummary   `json:"summary"`
	Email    apiNotifyEmail     `json:"email"`
}

type apiNotifyWebhook struct {
	URL                  string   `json:"url"`
	Label                string   `json:"label"`
	Enabled              bool     `json:"enabled"`
	On                   []string `json:"on"`
	PayloadTemplate      *string  `json:"payload_template"`
	RetryCount           int      `json:"retry_count"`
	RetryIntervalSeconds int      `json:"retry_interval_seconds"`
	Secret               *string  `json:"secret"`
}

func (w *apiNotifyWebhook) UnmarshalJSON(data []byte) error {
	type alias apiNotifyWebhook
	next := alias{Enabled: true}
	if err := json.Unmarshal(data, &next); err != nil {
		return err
	}
	*w = apiNotifyWebhook(next)
	return nil
}

type apiNotifyChannel struct {
	ID                   string         `json:"id"`
	Type                 string         `json:"type"`
	Label                string         `json:"label"`
	Enabled              bool           `json:"enabled"`
	On                   []string       `json:"on"`
	Config               map[string]any `json:"config"`
	RetryCount           int            `json:"retry_count"`
	RetryIntervalSeconds int            `json:"retry_interval_seconds"`
}

func (c *apiNotifyChannel) UnmarshalJSON(data []byte) error {
	type alias apiNotifyChannel
	next := alias{Enabled: true}
	if err := json.Unmarshal(data, &next); err != nil {
		return err
	}
	*c = apiNotifyChannel(next)
	return nil
}

type apiNotifySummary struct {
	Enabled   bool   `json:"enabled"`
	Interval  string `json:"interval"`
	Hour      int    `json:"hour"`
	DayOfWeek int    `json:"day_of_week"`
}

type apiNotifyEmail struct {
	Enabled bool     `json:"enabled"`
	To      []string `json:"to"`
	On      []string `json:"on"`
}

type apiAccessControl struct {
	Allow []string `json:"allow"`
}

type apiServerConfig struct {
	LogMaxLines             int              `json:"log_max_lines"`
	HistoryMaxCount         int              `json:"history_max_count"`
	BuildTimeoutSeconds     int              `json:"build_timeout_seconds"`
	LogRetentionDays        int              `json:"log_retention_days"`
	LogArchiveAfterDays     int              `json:"log_archive_after_days"`
	LogLevel                string           `json:"log_level"`
	PATExpiresAt            *string          `json:"pat_expires_at"`
	SnapshotsKeep           int              `json:"snapshots_keep"`
	QueueMaxSize            int              `json:"queue_max_size"`
	BuildRetryMax           int              `json:"build_retry_max"`
	BuildRetryBaseSeconds   int              `json:"build_retry_base_seconds"`
	CommitStatusEnabled     bool             `json:"commit_status_enabled"`
	CommitStatusContext     string           `json:"commit_status_context"`
	CommitStatusTargetURL   *string          `json:"commit_status_target_url"`
	BuildTrendKeepCount     int              `json:"build_trend_keep_count"`
	SessionTimeoutSeconds   int              `json:"session_timeout_seconds"`
	ForceBuildIntervalHours int              `json:"force_build_interval_hours"`
	BuildCooldownSeconds    int              `json:"build_cooldown_seconds"`
	ScheduleIntervalSeconds int              `json:"schedule_interval_seconds"`
	SchedulePaused          bool             `json:"schedule_paused"`
	AllowedHours            *apiAllowedHours `json:"allowed_hours"`
	APIRateLimit            map[string]any   `json:"api_rate_limit,omitempty"`
}

type apiLogRecord struct {
	At     string `json:"at"`
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
	Status int    `json:"status,omitempty"`
	Action string `json:"action,omitempty"`
	Result string `json:"result,omitempty"`
}

type apiTokenRecord struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	TokenHash string   `json:"token_hash"`
	Scopes    []string `json:"scopes"`
	CreatedAt string   `json:"created_at"`
	RevokedAt *string  `json:"revoked_at"`
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
	OutputSHA256    string         `json:"output_sha256,omitempty"`
	SHA256          string         `json:"sha256,omitempty"`
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
	return &APIServer{cfg: cfg, startedAt: cfg.Now().UTC(), sessions: map[string]apiSession{}, loginTickets: map[string]apiLoginTicket{}}, nil
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
	mux.HandleFunc("/api/auth/totp-status", s.withAuth(s.handleTOTPStatus))
	mux.HandleFunc("/api/auth/totp-setup", s.withAuth(s.handleTOTPSetup))
	mux.HandleFunc("/api/auth/totp-confirm", s.withAuth(s.handleTOTPConfirm))
	mux.HandleFunc("/api/auth/totp", s.withAuth(s.handleTOTPDisable))
	mux.HandleFunc("/api/config", s.withAuth(s.handleConfig))
	mux.HandleFunc("/api/config/validate", s.withAuth(s.handleConfigValidate))
	mux.HandleFunc("/api/log-level", s.withAuth(s.handleLogLevel))
	mux.HandleFunc("/api/api-rate-limit", s.withAuth(s.handleAPIRateLimit))
	mux.HandleFunc("/api/access-control", s.withAuth(s.handleAccessControl))
	mux.HandleFunc("/api/audit-log", s.withAuth(s.handleAuditLog))
	mux.HandleFunc("/api/access-log", s.withAuth(s.handleJSONLinesLog(".access_log", "log")))
	mux.HandleFunc("/api/api-access-log", s.withAuth(s.handleAPIAccessLog))
	mux.HandleFunc("/api/config-log", s.withAuth(s.handleJSONLinesLog(".config_log", "log")))
	mux.HandleFunc("/api/status", s.withAuth(s.handleStatus))
	mux.HandleFunc("/api/sysinfo", s.withAuth(s.handleSysinfo))
	mux.HandleFunc("/api/stats", s.withAuth(s.handleStats))
	mux.HandleFunc("/api/stats/timeline", s.withAuth(s.handleStatsTimeline))
	mux.HandleFunc("/api/stats/build-duration", s.withAuth(s.handleStatsBuildDuration))
	mux.HandleFunc("/api/output-meta", s.withAuth(s.handleOutputMeta))
	mux.HandleFunc("/api/dashboard", s.withAuth(s.handleDashboard))
	mux.HandleFunc("/api/notify-config", s.withAuth(s.handleNotifyConfig))
	mux.HandleFunc("/api/notify-log", s.withAuth(s.handleJSONLinesLog(".notify_log", "log")))
	mux.HandleFunc("/api/notify-test", s.withAuth(s.handleNotifyTest))
	mux.HandleFunc("/api/notify/weekly-summary", s.withAuth(s.handleNotifyWeeklySummary))
	mux.HandleFunc("/api/repo-info", s.withAuth(s.handleRepoInfo))
	mux.HandleFunc("/api/repo-config", s.withAuth(s.handleRepoConfig))
	mux.HandleFunc("/api/branch-config", s.withAuth(s.handleBranchConfig))
	mux.HandleFunc("/api/schedule", s.withAuth(s.handleSchedule))
	mux.HandleFunc("/api/schedule/interval", s.withAuth(s.handleScheduleInterval))
	mux.HandleFunc("/api/schedule/pause", s.withAuth(s.handleSchedulePause))
	mux.HandleFunc("/api/schedule/resume", s.withAuth(s.handleScheduleResume))
	mux.HandleFunc("/api/schedule/allowed-hours", s.withAuth(s.handleScheduleAllowedHours))
	mux.HandleFunc("/api/schedule/force-interval", s.withAuth(s.handleScheduleForceInterval))
	mux.HandleFunc("/api/schedule/cooldown", s.withAuth(s.handleScheduleCooldown))
	mux.HandleFunc("/api/build", s.withAuth(s.handleBuild(false)))
	mux.HandleFunc("/api/build/force", s.withAuth(s.handleBuild(true)))
	mux.HandleFunc("/api/build/cancel", s.withAuth(s.handleCancel))
	mux.HandleFunc("/api/build/stream", s.withAuth(s.handleBuildStream))
	mux.HandleFunc("/api/logs", s.withAuth(s.handleLogs))
	mux.HandleFunc("/api/logs/search", s.withAuth(s.handleLogSearch))
	mux.HandleFunc("/api/logs/export", s.withAuth(s.handleLogExport))
	mux.HandleFunc("/api/logs/cleanup", s.withAuth(s.handleLogsCleanup))
	mux.HandleFunc("/api/logs/archive", s.withAuth(s.handleLogsArchive))
	mux.HandleFunc("/api/history", s.withAuth(s.handleHistory))
	mux.HandleFunc("/api/history/export", s.withAuth(s.handleHistoryExport))
	mux.HandleFunc("/api/history/", s.withAuth(s.handleHistoryPath))
	mux.HandleFunc("/api/pat-status", s.withAuth(s.handlePATStatus))
	mux.HandleFunc("/api/pat-verify", s.withAuth(s.handlePATVerify))
	mux.HandleFunc("/api/pat-update", s.withAuth(s.handlePATUpdate))
	mux.HandleFunc("/api/backup", s.withAuth(s.handleBackup))
	mux.HandleFunc("/api/restore", s.withAuth(s.handleRestore))
	mux.HandleFunc("/api/diagnostics", s.withAuth(s.handleDiagnostics))
	mux.HandleFunc("/api/rate-limit", s.withAuth(s.handleRateLimit))
	mux.HandleFunc("/api/disk-usage", s.withAuth(s.handleDiskUsage))
	mux.HandleFunc("/api/webhook-events", s.withAuth(s.handleWebhookEvents))
	mux.HandleFunc("/api/webhook-config", s.withAuth(s.handleWebhookConfig))
	mux.HandleFunc("/api/webhook", s.handleWebhook)
	mux.HandleFunc("/api/snapshots", s.withAuth(s.handleSnapshots))
	mux.HandleFunc("/api/snapshots/", s.withAuth(s.handleSnapshotPath))
	mux.HandleFunc("/api/tokens", s.withAuth(s.handleTokens))
	mux.HandleFunc("/api/tokens/", s.withAuth(s.handleTokenPath))
	mux.HandleFunc("/api/alert-rules", s.withAuth(s.handleRuleFile(".alert_rules", "alert_rule")))
	mux.HandleFunc("/api/alert-rules/", s.withAuth(s.handleRulePath(".alert_rules", "alert_rule")))
	mux.HandleFunc("/api/tag-rules", s.withAuth(s.handleRuleFile(".tag_rules", "tag_rule")))
	mux.HandleFunc("/api/tag-rules/", s.withAuth(s.handleRulePath(".tag_rules", "tag_rule")))
	mux.HandleFunc("/api/hooks", s.withAuth(s.handleRuleFile(".hooks", "hook")))
	mux.HandleFunc("/api/hooks/", s.withAuth(s.handleHookPath))
	mux.HandleFunc("/api/verify-output", s.withAuth(s.handleVerifyOutput))
	mux.HandleFunc("/api/pipeline-config", s.withAuth(s.handlePipelineConfig))
	mux.HandleFunc("/api/notes", s.withAuth(s.handleNotes))
	mux.HandleFunc("/api/dashboard-layout", s.withAuth(s.handleDashboardLayout))
	mux.HandleFunc("/api/smtp-config", s.withAuth(s.handleSMTPConfig))
	mux.HandleFunc("/api/smtp-test", s.withAuth(s.handleSMTPTest))
	mux.HandleFunc("/api/circuit-breaker/reset", s.withAuth(s.handleCircuitReset))
	mux.HandleFunc("/api/maintenance", s.withAuth(s.handleMaintenance))
	mux.HandleFunc("/api/maintenance/enable", s.withAuth(s.handleMaintenanceEnable))
	mux.HandleFunc("/api/maintenance/disable", s.withAuth(s.handleMaintenanceDisable))
	mux.HandleFunc("/api/queue", s.withAuth(s.handleQueue))
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Not found")
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		if strings.HasPrefix(r.URL.Path, "/api/") && !s.accessAllowed(r) {
			writeError(rec, http.StatusForbidden, "Forbidden")
		} else {
			mux.ServeHTTP(rec, r)
		}
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
	mustChange := "none"
	if cred.MustChange {
		mustChange = "prompt"
	}
	if totp, err := readTOTPSecret(filepath.Join(s.cfg.StateDir, ".totp_secret")); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	} else if totp.Enabled {
		ticket, err := randomToken()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		s.mu.Lock()
		s.loginTickets[tokenHash(ticket)] = apiLoginTicket{MustChange: mustChange, ExpiresAt: s.cfg.Now().UTC().Add(5 * time.Minute)}
		s.mu.Unlock()
		writeJSON(w, http.StatusOK, map[string]any{"must_change": mustChange, "totp_required": true, "ticket": ticket})
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
	ticketHash := tokenHash(req.Ticket)
	s.mu.Lock()
	ticket, ok := s.loginTickets[ticketHash]
	if ok {
		delete(s.loginTickets, ticketHash)
	}
	s.mu.Unlock()
	if !ok || !s.cfg.Now().UTC().Before(ticket.ExpiresAt) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	totp, err := readTOTPSecret(filepath.Join(s.cfg.StateDir, ".totp_secret"))
	if err != nil || !totp.Enabled || totp.SecretBase32 == nil || !verifyTOTPCode(*totp.SecretBase32, req.Code, s.cfg.Now().UTC()) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
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
	now := s.nowString()
	expires := s.cfg.Now().UTC().Add(time.Duration(timeoutSeconds) * time.Second).Format(apiTimeLayout)
	s.mu.Lock()
	s.sessions[tokenHash(token)] = apiSession{TokenHash: tokenHash(token), CreatedAt: now, ExpiresAt: expires, LastUsedAt: now}
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "must_change": ticket.MustChange})
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
	lastBuildAt := any(nil)
	lastBuildStatus := "none"
	lastDeployStatus := any(nil)
	pendingTransfersCount := 0
	if st, ok, err := readBuildStatus(filepath.Join(s.cfg.StateDir, ".build_status.json")); err != nil {
		status = "degraded"
		items = append(items, map[string]any{"name": "build_status", "status": "error"})
	} else if !ok {
		status = "degraded"
		items = append(items, map[string]any{"name": "build_status", "status": "warn"})
	} else {
		items = append(items, map[string]any{"name": "build_status", "status": "ok"})
		lastBuildAt = coalesceString(st.LastFinishedAt, st.LastStartedAt)
		lastBuildStatus = stringOr(st.LastTargetStatus, "none")
		lastDeployStatus = st.LastDeployStatus
		pendingTransfersCount = st.PendingTransfersCount
	}
	if lastBuildAt == nil {
		history := readHistory(filepath.Join(s.cfg.StateDir, ".build_history"))
		if len(history) > 0 {
			lastBuildAt = firstNonEmpty(history[0].FinishedAt, history[0].StartedAt, history[0].BuildAt)
			lastBuildStatus = history[0].Status
		}
	}
	uptime := int64(s.cfg.Now().UTC().Sub(s.startedAt).Seconds())
	if uptime < 0 {
		uptime = 0
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":                  status,
		"checked_at":              s.nowString(),
		"uptime_seconds":          uptime,
		"last_build_at":           lastBuildAt,
		"last_build_status":       lastBuildStatus,
		"last_deploy_status":      lastDeployStatus,
		"pending_transfers_count": pendingTransfersCount,
		"items":                   items,
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

func (s *APIServer) handleTOTPStatus(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	totp, err := readTOTPSecret(filepath.Join(s.cfg.StateDir, ".totp_secret"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	writeJSON(w, http.StatusOK, totpStatusPayload(totp))
}

func (s *APIServer) handleTOTPSetup(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	totp, err := readTOTPSecret(filepath.Join(s.cfg.StateDir, ".totp_secret"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if totp.Enabled {
		writeError(w, http.StatusConflict, "Conflict")
		return
	}
	secret, err := randomBase32Secret()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	s.mu.Lock()
	s.totpSetup = &apiPendingTOTP{SecretBase32: secret, ExpiresAt: s.cfg.Now().UTC().Add(10 * time.Minute)}
	s.mu.Unlock()
	uri := "otpauth://totp/Adlaire%20CI:admin?issuer=Adlaire%20CI&secret=" + secret + "&algorithm=SHA1&digits=6&period=30"
	writeJSON(w, http.StatusOK, map[string]any{"secret": secret, "otpauth_uri": uri})
}

func (s *APIServer) handleTOTPConfirm(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var body struct {
		Code string `json:"code"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	if body.Code == "" {
		writeValidation(w, "code", "required")
		return
	}
	s.mu.Lock()
	pending := s.totpSetup
	s.mu.Unlock()
	if pending == nil || !s.cfg.Now().UTC().Before(pending.ExpiresAt) {
		writeError(w, http.StatusConflict, "Conflict")
		return
	}
	if !verifyTOTPCode(pending.SecretBase32, body.Code, s.cfg.Now().UTC()) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	now := s.nowString()
	secret := pending.SecretBase32
	totp := apiTOTPSecret{Enabled: true, SecretBase32: &secret, ConfirmedAt: &now}
	if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".totp_secret"), totp, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	s.mu.Lock()
	s.totpSetup = nil
	s.mu.Unlock()
	_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": now, "action": "totp_enable", "result": "success"})
	writeJSON(w, http.StatusOK, totpStatusPayload(totp))
}

func (s *APIServer) handleTOTPDisable(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodDelete) {
		return
	}
	var body struct {
		Code string `json:"code"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	if body.Code == "" {
		writeValidation(w, "code", "required")
		return
	}
	path := filepath.Join(s.cfg.StateDir, ".totp_secret")
	totp, err := readTOTPSecret(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if !totp.Enabled || totp.SecretBase32 == nil {
		writeError(w, http.StatusConflict, "Conflict")
		return
	}
	if !verifyTOTPCode(*totp.SecretBase32, body.Code, s.cfg.Now().UTC()) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	next := apiTOTPSecret{Enabled: false}
	if err := atomicWriteJSON(path, next, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": s.nowString(), "action": "totp_disable", "result": "success"})
	writeJSON(w, http.StatusOK, totpStatusPayload(next))
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

func (s *APIServer) handleLogLevel(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var body struct {
		Level string `json:"level"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	level := strings.ToUpper(strings.TrimSpace(body.Level))
	if level != "INFO" && level != "DEBUG" && level != "WARNING" && level != "ERROR" {
		writeValidation(w, "level", "invalid value")
		return
	}
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if cfg.LogLevel == level {
		writeJSON(w, http.StatusOK, map[string]any{"message": "No changes", "level": level})
		return
	}
	cfg.LogLevel = level
	if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".server_config"), cfg, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "log_level", Changes: map[string]any{"log_level": level}})
	writeJSON(w, http.StatusOK, map[string]any{"message": "Log level updated", "level": level})
}

func (s *APIServer) handleAPIRateLimit(w http.ResponseWriter, r *http.Request) {
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
		policy := cfg.APIRateLimit
		if policy == nil {
			policy = map[string]any{"enabled": false, "groups": []any{}}
		}
		state, _, _ := readOptionalJSONMap(filepath.Join(s.cfg.StateDir, ".api_rate_state"))
		writeJSON(w, http.StatusOK, map[string]any{"policy": policy, "state_summary": state})
	case http.MethodPost:
		var body map[string]any
		if !decodeBody(w, r, &body, true) {
			return
		}
		if _, ok := body["enabled"].(bool); !ok {
			writeValidation(w, "enabled", "required")
			return
		}
		if groups, ok := body["groups"]; ok {
			if _, ok := groups.([]any); !ok {
				writeValidation(w, "groups", "invalid type")
				return
			}
		} else {
			body["groups"] = []any{}
		}
		cfg, err := s.readMergedConfig()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		cfg.APIRateLimit = body
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".server_config"), cfg, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".api_rate_state"), map[string]any{"windows": map[string]any{}}, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "api_rate_limit", Changes: map[string]any{"policy": body}})
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": s.nowString(), "action": "api_rate_limit_update", "result": "success"})
		writeJSON(w, http.StatusOK, map[string]any{"policy": body, "state_summary": map[string]any{"windows": map[string]any{}}})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleAuditLog(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	limit, ok := parseBoundedInt(w, r, "limit", 100, 1, 200)
	if !ok {
		return
	}
	offset, ok := parseBoundedInt(w, r, "offset", 0, 0, 1_000_000)
	if !ok {
		return
	}
	actor := strings.TrimSpace(r.URL.Query().Get("actor"))
	action := strings.TrimSpace(r.URL.Query().Get("action"))
	result := strings.TrimSpace(r.URL.Query().Get("result"))
	for _, pair := range []struct{ field, value string }{{"actor", actor}, {"action", action}, {"result", result}} {
		if len(pair.value) > 200 || containsControl(pair.value) {
			writeValidation(w, pair.field, "invalid value")
			return
		}
	}
	records := readJSONLines(filepath.Join(s.cfg.StateDir, ".audit_log"))
	filtered := []map[string]any{}
	for _, record := range records {
		if actor != "" && fmt.Sprint(record["actor"]) != actor {
			continue
		}
		if action != "" && fmt.Sprint(record["action"]) != action {
			continue
		}
		if result != "" && fmt.Sprint(record["result"]) != result {
			continue
		}
		filtered = append(filtered, record)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return fmt.Sprint(filtered[i]["at"]) > fmt.Sprint(filtered[j]["at"])
	})
	total := len(filtered)
	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}
	writeJSON(w, http.StatusOK, map[string]any{"log": filtered[start:end], "total": total})
}

func (s *APIServer) handleAccessControl(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		cfg, err := s.readAccessControl()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		writeJSON(w, http.StatusOK, cfg)
	case http.MethodPost:
		var body apiAccessControl
		if !decodeBody(w, r, &body, true) {
			return
		}
		next, ok := normalizeAccessControl(w, body)
		if !ok {
			return
		}
		current, err := s.readAccessControl()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if stringSlicesEqual(current.Allow, next.Allow) {
			writeJSON(w, http.StatusOK, map[string]any{"message": "No changes", "allow": next.Allow})
			return
		}
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".access_control"), next, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "access_control", Changes: map[string]any{"allow": next.Allow}}); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"message": "Access control updated", "allow": next.Allow})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleRepoInfo(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	cfg, err := s.readRepoConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (s *APIServer) handleRepoConfig(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var patch map[string]json.RawMessage
	if !decodeBody(w, r, &patch, true) {
		return
	}
	if len(patch) == 0 {
		writeValidation(w, "body", "must include at least one field")
		return
	}
	current, err := s.readRepoConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	next := current
	changed := false
	seen := false
	for key, raw := range patch {
		switch key {
		case "owner":
			value, ok := decodeTrimmedString(w, key, raw, 1, 100)
			if !ok {
				return
			}
			seen = true
			changed = changed || next.Owner != value
			next.Owner = value
		case "repo":
			value, ok := decodeTrimmedString(w, key, raw, 1, 100)
			if !ok {
				return
			}
			seen = true
			changed = changed || next.Repo != value
			next.Repo = value
		case "branch":
			value, ok := decodeTrimmedString(w, key, raw, 1, 128)
			if !ok || !validateBranchName(w, value) {
				return
			}
			seen = true
			changed = changed || next.Branch != value
			next.Branch = value
		case "target_file":
			value, ok := decodeTrimmedString(w, key, raw, 1, 500)
			if !ok || !validateTargetFile(w, value) {
				return
			}
			seen = true
			changed = changed || next.TargetFile != value
			next.TargetFile = value
		default:
			writeValidation(w, key, "unknown key")
			return
		}
	}
	if !seen {
		writeValidation(w, "body", "must include at least one field")
		return
	}
	if !changed {
		writeJSON(w, http.StatusOK, map[string]string{"message": "No changes"})
		return
	}
	next.UpdatedAt = s.nowString()
	if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".repo_config"), next, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "repo_config", Changes: maskRepoConfigChanges(patch)}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Repo config updated"})
}

func (s *APIServer) handleBranchConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		cfg, source, err := s.readBranchConfig()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"source": source, "branches": cfg.BranchTargets})
	case http.MethodPost:
		var body struct {
			Branches      []apiBranchTarget `json:"branches"`
			BranchTargets []apiBranchTarget `json:"branch_targets"`
		}
		if !decodeBody(w, r, &body, true) {
			return
		}
		targets := body.Branches
		if targets == nil {
			targets = body.BranchTargets
		}
		if targets == nil {
			writeValidation(w, "branches", "required")
			return
		}
		if len(targets) == 0 {
			if _, err := os.Stat(filepath.Join(s.cfg.StateDir, ".branch_config")); errors.Is(err, os.ErrNotExist) {
				writeJSON(w, http.StatusOK, map[string]any{"message": "No changes", "branches_count": 0})
				return
			}
			if err := removeIfExists(filepath.Join(s.cfg.StateDir, ".branch_config")); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "branch_config", Changes: map[string]any{"branches_count": 0}}); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"message": "Branch config updated", "branches_count": 0})
			return
		}
		if !validateBranchTargets(w, targets) {
			return
		}
		next := apiBranchConfigFile{BranchTargets: targets}
		current, source, err := s.readBranchConfig()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if source == "file" && branchTargetsEqual(current.BranchTargets, targets) {
			writeJSON(w, http.StatusOK, map[string]any{"message": "No changes", "branches_count": len(targets)})
			return
		}
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".branch_config"), next, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "branch_config", Changes: map[string]any{"branches_count": len(targets)}}); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"message": "Branch config updated", "branches_count": len(targets)})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleNotifyConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		cfg, err := s.readNotifyConfig()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		writeJSON(w, http.StatusOK, maskNotifyConfig(cfg))
	case http.MethodPost:
		var cfg apiNotifyConfig
		if !decodeBody(w, r, &cfg, true) {
			return
		}
		next, ok := normalizeNotifyConfig(w, cfg)
		if !ok {
			return
		}
		current, err := s.readNotifyConfig()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if notifyConfigsEqual(current, next) {
			writeJSON(w, http.StatusOK, map[string]string{"message": "No changes"})
			return
		}
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".notify_config"), next, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "notify_config", Changes: notifyConfigChangeSummary(next)}); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Notify config updated"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleNotifyTest(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	cfg, err := s.readNotifyConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	targets := notifyWebhookTargets(cfg, "failure")
	if len(targets) == 0 {
		writeValidation(w, "webhook", "not configured")
		return
	}
	target := targets[0]
	payload := map[string]any{"event": "test", "sent_at": s.nowString()}
	status, sendErr := sendWebhook(target.URL, target.Secret, payload)
	result := "success"
	errText := any(nil)
	if sendErr != nil {
		result = "failure"
		errText = "Webhook delivery failed"
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".notify_log"), map[string]any{
		"at":          s.nowString(),
		"event":       "test",
		"result":      result,
		"http_status": status,
		"attempt":     1,
		"error":       errText,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if sendErr != nil {
		writeError(w, http.StatusUnprocessableEntity, "Webhook delivery failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Test notification sent", "webhook_url": target.URL})
}

func (s *APIServer) handleNotifyWeeklySummary(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	cfg, err := s.readNotifyConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	targets := notifyWebhookTargets(cfg, "weekly_summary")
	if len(targets) == 0 {
		writeValidation(w, "webhook", "not configured")
		return
	}
	summary := s.weeklySummary()
	payload := map[string]any{
		"event":                "weekly_summary",
		"period_days":          7,
		"success_count":        summary["success_count"],
		"failure_count":        summary["failure_count"],
		"success_rate":         summary["success_rate"],
		"avg_duration_seconds": summary["avg_duration_seconds"],
	}
	status, sendErr := sendWebhook(targets[0].URL, targets[0].Secret, payload)
	result := "success"
	errText := any(nil)
	if sendErr != nil {
		result = "failure"
		errText = "Webhook delivery failed"
	}
	logRecord := map[string]any{
		"at":          s.nowString(),
		"event":       "weekly_summary",
		"result":      result,
		"http_status": status,
		"attempt":     1,
		"error":       errText,
	}
	for key, value := range summary {
		logRecord[key] = value
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".notify_log"), logRecord); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if sendErr != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message":       "Weekly summary sent",
		"period":        summary["period"],
		"success_count": summary["success_count"],
		"failure_count": summary["failure_count"],
		"success_rate":  summary["success_rate"],
	})
}

func (s *APIServer) handleSchedule(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, schedulePayload(cfg))
}

func (s *APIServer) handleScheduleInterval(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var body struct {
		IntervalSeconds int `json:"interval_seconds"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	if body.IntervalSeconds < 30 || body.IntervalSeconds > 86400 {
		writeValidation(w, "interval_seconds", "out of range")
		return
	}
	if !s.updateScheduleConfig(w, "schedule_interval", map[string]any{"interval_seconds": body.IntervalSeconds}, func(cfg *apiServerConfig) (bool, map[string]any) {
		changed := cfg.ScheduleIntervalSeconds != body.IntervalSeconds
		cfg.ScheduleIntervalSeconds = body.IntervalSeconds
		return changed, map[string]any{"message": "Interval updated", "interval_seconds": body.IntervalSeconds}
	}) {
		return
	}
}

func (s *APIServer) handleSchedulePause(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	if !s.updateScheduleConfig(w, "schedule_pause", map[string]any{"schedule_paused": true}, func(cfg *apiServerConfig) (bool, map[string]any) {
		if cfg.SchedulePaused {
			writeError(w, http.StatusConflict, "Schedule already paused")
			return false, nil
		}
		cfg.SchedulePaused = true
		return true, map[string]any{"message": "Schedule paused"}
	}) {
		return
	}
}

func (s *APIServer) handleScheduleResume(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	if !s.updateScheduleConfig(w, "schedule_resume", map[string]any{"schedule_paused": false}, func(cfg *apiServerConfig) (bool, map[string]any) {
		if !cfg.SchedulePaused {
			writeError(w, http.StatusConflict, "Schedule already running")
			return false, nil
		}
		cfg.SchedulePaused = false
		return true, map[string]any{"message": "Schedule resumed"}
	}) {
		return
	}
}

func (s *APIServer) handleScheduleAllowedHours(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var body struct {
		From *int `json:"from"`
		To   *int `json:"to"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	var next *apiAllowedHours
	if body.From != nil || body.To != nil {
		if body.From == nil || body.To == nil || *body.From < 0 || *body.From > 23 || *body.To < 0 || *body.To > 23 || *body.From == *body.To {
			writeValidation(w, "allowed_hours", "invalid value")
			return
		}
		next = &apiAllowedHours{From: *body.From, To: *body.To}
	}
	if !s.updateScheduleConfig(w, "schedule_allowed_hours", map[string]any{"allowed_hours": next}, func(cfg *apiServerConfig) (bool, map[string]any) {
		changed := !allowedHoursEqual(cfg.AllowedHours, next)
		cfg.AllowedHours = next
		return changed, map[string]any{"message": "Allowed hours updated", "allowed_hours": next}
	}) {
		return
	}
}

func (s *APIServer) handleScheduleForceInterval(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var body struct {
		Hours int `json:"hours"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	if body.Hours < 0 || body.Hours > 8760 {
		writeValidation(w, "hours", "out of range")
		return
	}
	if !s.updateScheduleConfig(w, "schedule_force_interval", map[string]any{"hours": body.Hours}, func(cfg *apiServerConfig) (bool, map[string]any) {
		changed := cfg.ForceBuildIntervalHours != body.Hours
		cfg.ForceBuildIntervalHours = body.Hours
		return changed, map[string]any{"message": "Force build interval updated", "hours": body.Hours}
	}) {
		return
	}
}

func (s *APIServer) handleScheduleCooldown(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var body struct {
		Seconds int `json:"seconds"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	if body.Seconds < 0 || body.Seconds > 86400 {
		writeValidation(w, "seconds", "out of range")
		return
	}
	if !s.updateScheduleConfig(w, "schedule_cooldown", map[string]any{"seconds": body.Seconds}, func(cfg *apiServerConfig) (bool, map[string]any) {
		changed := cfg.BuildCooldownSeconds != body.Seconds
		cfg.BuildCooldownSeconds = body.Seconds
		return changed, map[string]any{"message": "Build cooldown updated", "seconds": body.Seconds}
	}) {
		return
	}
}

func (s *APIServer) updateScheduleConfig(w http.ResponseWriter, logType string, changes map[string]any, update func(*apiServerConfig) (bool, map[string]any)) bool {
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return false
	}
	changed, resp := update(&cfg)
	if resp == nil {
		return false
	}
	if !changed {
		writeJSON(w, http.StatusOK, map[string]any{"message": "No changes"})
		return true
	}
	if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".server_config"), cfg, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return false
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: logType, Changes: changes}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return false
	}
	writeJSON(w, http.StatusOK, resp)
	return true
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
			"queued":                  state.Queued,
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
		"queued":                  state.Queued,
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
		if maintenanceEnabled(filepath.Join(s.cfg.StateDir, ".maintenance")) {
			writeError(w, http.StatusServiceUnavailable, "maintenance")
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
		statePath := filepath.Join(s.cfg.StateDir, ".build_state")
		state, err := readBuildState(statePath)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if !state.Running && lockExists(filepath.Join(s.cfg.StateDir, ".build_lock")) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		if state.Running {
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

func (s *APIServer) handleLogsCleanup(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if cfg.LogRetentionDays <= 0 {
		writeJSON(w, http.StatusOK, map[string]any{"message": "No logs deleted", "deleted_count": 0, "failed_count": 0})
		return
	}
	cutoff := s.cfg.Now().UTC().AddDate(0, 0, -cfg.LogRetentionDays)
	deleted, failed := cleanupLogFiles(filepath.Join(s.cfg.StateDir, ".build_logs"), cutoff)
	writeJSON(w, http.StatusOK, map[string]any{"message": "Logs cleaned up", "deleted_count": deleted, "failed_count": failed})
}

func (s *APIServer) handleLogsArchive(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if cfg.LogArchiveAfterDays <= 0 {
		writeJSON(w, http.StatusOK, map[string]any{"message": "No logs archived", "archived_count": 0})
		return
	}
	cutoff := s.cfg.Now().UTC().AddDate(0, 0, -cfg.LogArchiveAfterDays)
	count, err := archiveLogFiles(filepath.Join(s.cfg.StateDir, ".build_logs"), cutoff)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Logs archived", "archived_count": count})
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
	history = filterHistory(w, r, history)
	if history == nil {
		return
	}
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
	rest := strings.TrimPrefix(r.URL.Path, "/api/history/")
	id, suffix, ok := strings.Cut(rest, "/")
	if !ok || id == "" {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	log, err := readBuildLogByID(s.cfg.StateDir, id)
	if errors.Is(err, os.ErrNotExist) {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	switch suffix {
	case "log":
		if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
			return
		}
		lines := append(splitLines(log.Pipeline.Stdout), splitLines(log.Pipeline.Stderr)...)
		lines = append(lines, log.Warnings...)
		writeJSON(w, http.StatusOK, map[string]any{
			"id": log.ID, "build_at": firstNonEmpty(log.FinishedAt, log.StartedAt),
			"sha": firstCommitSHA(log.Commit), "status": log.TargetStatus, "comment": log.Comment,
			"flagged": log.Flagged, "tags": log.Tags, "lines": lines,
		})
	case "comment":
		switch r.Method {
		case http.MethodGet:
			if !rejectBody(w, r) {
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"id": log.ID, "comment": log.Comment})
		case http.MethodPost:
			var body struct {
				Comment *string `json:"comment"`
			}
			if !decodeBody(w, r, &body, true) {
				return
			}
			if body.Comment == nil {
				writeValidation(w, "comment", "required")
				return
			}
			comment := strings.TrimSpace(*body.Comment)
			if len(comment) > 2000 {
				writeValidation(w, "comment", "must be 2000 characters or less")
				return
			}
			if comment == "" {
				log.Comment = nil
			} else {
				log.Comment = &comment
			}
			if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".build_logs", id+".json"), log, 0600); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "history_comment", Changes: map[string]any{"id": id}})
			writeJSON(w, http.StatusOK, map[string]any{"id": log.ID, "comment": log.Comment})
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	case "flag":
		if !method(w, r, http.MethodPost) {
			return
		}
		var body struct {
			Flagged *bool `json:"flagged"`
		}
		if !decodeBody(w, r, &body, true) {
			return
		}
		if body.Flagged == nil {
			writeValidation(w, "flagged", "required")
			return
		}
		log.Flagged = *body.Flagged
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".build_logs", id+".json"), log, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := updateHistoryRecord(filepath.Join(s.cfg.StateDir, ".build_history"), id, func(record *apiHistoryRecord) {
			record.Flagged = *body.Flagged
		}); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "history_flag", Changes: map[string]any{"id": id, "flagged": *body.Flagged}})
		writeJSON(w, http.StatusOK, map[string]any{"id": log.ID, "flagged": log.Flagged})
	case "tags":
		if !method(w, r, http.MethodPost) {
			return
		}
		var body struct {
			Tags []string `json:"tags"`
		}
		if !decodeBody(w, r, &body, true) {
			return
		}
		tags, ok := normalizeTags(w, body.Tags)
		if !ok {
			return
		}
		log.Tags = tags
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".build_logs", id+".json"), log, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := updateHistoryRecord(filepath.Join(s.cfg.StateDir, ".build_history"), id, func(record *apiHistoryRecord) {
			record.Tags = tags
		}); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "history_tags", Changes: map[string]any{"id": id, "tags": tags}})
		writeJSON(w, http.StatusOK, map[string]any{"id": log.ID, "tags": log.Tags})
	case "rollback":
		if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
			return
		}
		state, err := readBuildState(filepath.Join(s.cfg.StateDir, ".build_state"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if state.Running || lockExists(filepath.Join(s.cfg.StateDir, ".build_lock")) {
			writeError(w, http.StatusConflict, "Build is running")
			return
		}
		snapshotDir := filepath.Join(s.cfg.StateDir, ".snapshots", id)
		if _, err := os.Stat(snapshotDir); errors.Is(err, os.ErrNotExist) {
			writeError(w, http.StatusNotFound, "Not found")
			return
		} else if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		buildID := s.newID("rollback")
		now := s.nowString()
		if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".build_history"), apiHistoryRecord{ID: buildID, BuildAt: now, Status: "success", Trigger: "rollback"}); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: now, Type: "rollback", Changes: map[string]any{"snapshot_id": id, "build_id": buildID}})
		writeJSON(w, http.StatusAccepted, map[string]any{"message": "Rollback queued", "build_id": buildID})
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

func (s *APIServer) handleMaintenance(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	state, err := readMaintenanceState(filepath.Join(s.cfg.StateDir, ".maintenance"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (s *APIServer) handleMaintenanceEnable(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	reason := strings.TrimSpace(body.Reason)
	if reason == "" || len(reason) > 500 {
		writeValidation(w, "reason", "must be 1 to 500 characters")
		return
	}
	path := filepath.Join(s.cfg.StateDir, ".maintenance")
	current, err := readMaintenanceState(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if current.Enabled && stringPtrValue(current.Reason) == reason {
		writeJSON(w, http.StatusOK, map[string]any{"message": "No changes", "since": current.Since})
		return
	}
	since := s.nowString()
	next := apiMaintenanceState{Enabled: true, Reason: &reason, Since: &since}
	if err := atomicWriteJSON(path, next, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "maintenance", Changes: map[string]any{"enabled": true, "reason": reason}}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Maintenance enabled", "since": since})
}

func (s *APIServer) handleMaintenanceDisable(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	path := filepath.Join(s.cfg.StateDir, ".maintenance")
	current, err := readMaintenanceState(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if !current.Enabled {
		writeJSON(w, http.StatusOK, map[string]any{"message": "No changes"})
		return
	}
	next := apiMaintenanceState{Enabled: false}
	if err := atomicWriteJSON(path, next, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "maintenance", Changes: map[string]any{"enabled": false}}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Maintenance disabled"})
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

func (s *APIServer) handlePATStatus(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	info, err := os.Stat(filepath.Join(s.cfg.StateDir, ".github_token"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"configured": err == nil, "mode": fileModeString(info)})
}

func (s *APIServer) handlePATVerify(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	token, err := readSecretText(filepath.Join(s.cfg.StateDir, ".github_token"))
	if errors.Is(err, os.ErrNotExist) || strings.TrimSpace(token) == "" {
		writeError(w, http.StatusNotImplemented, "Not configured")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"configured": true, "valid": true, "checked_at": s.nowString(), "remote_checked": false})
}

func (s *APIServer) handlePATUpdate(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var body struct {
		Token string `json:"token"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	token := strings.TrimSpace(body.Token)
	if len(token) < 8 || len(token) > 512 || containsControl(token) {
		writeValidation(w, "token", "invalid value")
		return
	}
	if err := atomicWriteText(filepath.Join(s.cfg.StateDir, ".github_token"), token+"\n", 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "pat_update", Changes: map[string]any{"token": "***"}})
	writeJSON(w, http.StatusOK, map[string]string{"message": "PAT updated"})
}

func (s *APIServer) handleBackup(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	backup := map[string]any{}
	for _, file := range backupConfigFiles() {
		if value, ok, err := readOptionalJSONMap(filepath.Join(s.cfg.StateDir, file)); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		} else if ok {
			backup[strings.TrimPrefix(file, ".")] = maskSecrets(value)
		}
	}
	if _, err := os.Stat(filepath.Join(s.cfg.StateDir, ".github_token")); err == nil {
		backup["github_token_set"] = true
	}
	if _, err := os.Stat(filepath.Join(s.cfg.StateDir, ".webhook_secret")); err == nil {
		backup["webhook_secret_set"] = true
	}
	writeJSON(w, http.StatusOK, map[string]any{"exported_at": s.nowString(), "config": backup})
}

func (s *APIServer) handleRestore(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var body struct {
		Config map[string]map[string]any `json:"config"`
	}
	if !decodeBody(w, r, &body, true) {
		return
	}
	if body.Config == nil {
		writeValidation(w, "config", "required")
		return
	}
	allowed := map[string]string{}
	for _, file := range backupConfigFiles() {
		allowed[strings.TrimPrefix(file, ".")] = file
	}
	writes := map[string]map[string]any{}
	for key, value := range body.Config {
		file, ok := allowed[key]
		if !ok {
			writeValidation(w, key, "unknown key")
			return
		}
		if containsMaskedSecret(value) {
			current, _, _ := readOptionalJSONMap(filepath.Join(s.cfg.StateDir, file))
			value = mergeMaskedSecrets(value, current)
		}
		writes[file] = value
	}
	for _, file := range backupConfigFiles() {
		value, ok := writes[file]
		if !ok {
			continue
		}
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, file), value, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}
	_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "restore", Changes: map[string]any{"files": len(writes)}})
	writeJSON(w, http.StatusOK, map[string]string{"message": "Config restored"})
}

func (s *APIServer) handleDiagnostics(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	meta, err := inspectOutput(s.outputDir())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	_, tokenErr := os.Stat(filepath.Join(s.cfg.StateDir, ".github_token"))
	writeJSON(w, http.StatusOK, map[string]any{"items": []map[string]any{
		{"name": "github_token", "status": statusFromExists(tokenErr)},
		{"name": "output", "status": map[bool]string{true: "ok", false: "warn"}[meta.Exists]},
	}})
}

func (s *APIServer) handleRateLimit(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	token, err := readSecretText(filepath.Join(s.cfg.StateDir, ".github_token"))
	if errors.Is(err, os.ErrNotExist) || strings.TrimSpace(token) == "" {
		writeError(w, http.StatusNotImplemented, "Not configured")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"checked_at":       s.nowString(),
		"remote_checked":   false,
		"limit":            nil,
		"remaining":        nil,
		"reset_at":         nil,
		"resource":         "core",
		"token_configured": true,
	})
}

func (s *APIServer) handleDiskUsage(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	logBytes, logCount := fileTreeStats(filepath.Join(s.cfg.StateDir, ".build_logs"), func(path string) bool {
		return strings.HasSuffix(path, ".json") && !strings.Contains(filepath.ToSlash(path), "/archive/")
	})
	archiveBytes, archiveCount := fileTreeStats(filepath.Join(s.cfg.StateDir, ".build_logs", "archive"), func(path string) bool {
		return strings.HasSuffix(path, ".json.gz")
	})
	outputBytes, _ := fileTreeStats(s.outputDir(), func(path string) bool { return true })
	writeJSON(w, http.StatusOK, map[string]any{
		"build_logs_bytes":         logBytes,
		"build_logs_count":         logCount,
		"build_logs_archive_bytes": archiveBytes,
		"build_logs_archive_count": archiveCount,
		"output_file_bytes":        outputBytes,
		"total_bytes":              logBytes + archiveBytes + outputBytes,
	})
}

func (s *APIServer) handleWebhookConfig(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.cfg.StateDir, ".webhook_secret")
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		_, err := os.Stat(path)
		writeJSON(w, http.StatusOK, map[string]any{"configured": err == nil})
	case http.MethodPost:
		var body struct {
			Secret string `json:"secret"`
		}
		if !decodeBody(w, r, &body, true) {
			return
		}
		secret := strings.TrimSpace(body.Secret)
		if len(secret) < 8 || len(secret) > 256 || containsControl(secret) {
			writeValidation(w, "secret", "invalid value")
			return
		}
		if err := atomicWriteText(path, secret+"\n", 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "webhook_config", Changes: map[string]any{"secret": "***"}})
		writeJSON(w, http.StatusOK, map[string]string{"message": "Webhook config updated"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleWebhookEvents(w http.ResponseWriter, r *http.Request) {
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
	events := readJSONLines(filepath.Join(s.cfg.StateDir, ".webhook_events.json"))
	sort.Slice(events, func(i, j int) bool { return fmt.Sprint(events[i]["at"]) > fmt.Sprint(events[j]["at"]) })
	start := offset
	if start > len(events) {
		start = len(events)
	}
	end := start + limit
	if end > len(events) {
		end = len(events)
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": events[start:end], "total": len(events)})
}

func (s *APIServer) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusRequestEntityTooLarge, "Payload too large")
		return
	}
	secretData, err := os.ReadFile(filepath.Join(s.cfg.StateDir, ".webhook_secret"))
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if !validWebhookSignature(strings.TrimSpace(string(secretData)), raw, r.Header.Get("X-Hub-Signature-256")) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if maintenanceEnabled(filepath.Join(s.cfg.StateDir, ".maintenance")) {
		writeError(w, http.StatusServiceUnavailable, "maintenance")
		return
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Validation failed")
		return
	}
	eventID := firstNonEmpty(r.Header.Get("X-GitHub-Delivery"), s.newID("evt"))
	eventName := r.Header.Get("X-GitHub-Event")
	result := "ignored_event"
	queued := false
	var queueID any
	if eventName == "push" || eventName == "" {
		statePath := filepath.Join(s.cfg.StateDir, ".build_state")
		state, err := readBuildState(statePath)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		maxSize := s.queueMaxSize()
		if maxSize == 0 || len(state.Queued) >= maxSize {
			writeError(w, http.StatusTooManyRequests, "queue_full")
			return
		}
		id := s.newID("q")
		state.Queued = append(state.Queued, map[string]any{"id": id, "trigger": "webhook", "queued_at": s.nowString(), "requested_by": "webhook", "priority": "normal", "created_seq": nextCreatedSeq(state.Queued), "payload": map[string]any{"delivery_id": eventID, "ref": payload["ref"], "after": payload["after"]}})
		if err := atomicWriteJSON(statePath, state, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		result, queued, queueID = "queued", true, id
	}
	_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".webhook_events.json"), map[string]any{"at": s.nowString(), "event_id": eventID, "event": eventName, "result": result, "queued_id": queueID})
	writeJSON(w, http.StatusAccepted, map[string]any{"message": "Webhook accepted", "queued": queued, "event_id": eventID, "queue_id": queueID})
}

func (s *APIServer) handleSnapshots(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	root := filepath.Join(s.cfg.StateDir, ".snapshots")
	entries, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		writeJSON(w, http.StatusOK, map[string]any{"snapshots": []map[string]any{}})
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	snapshots := []map[string]any{}
	for _, entry := range entries {
		if !entry.IsDir() || !validSimpleID(entry.Name()) {
			continue
		}
		meta, _, _ := readOptionalJSONMap(filepath.Join(root, entry.Name(), "meta.json"))
		if meta == nil {
			meta = map[string]any{}
		}
		meta["id"] = entry.Name()
		snapshots = append(snapshots, meta)
	}
	sort.Slice(snapshots, func(i, j int) bool { return fmt.Sprint(snapshots[i]["id"]) > fmt.Sprint(snapshots[j]["id"]) })
	writeJSON(w, http.StatusOK, map[string]any{"snapshots": snapshots})
}

func (s *APIServer) handleSnapshotPath(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/snapshots/")
	id, suffix, _ := strings.Cut(rest, "/")
	if !validSimpleID(id) {
		writeValidation(w, "id", "invalid value")
		return
	}
	dir := filepath.Join(s.cfg.StateDir, ".snapshots", id)
	switch {
	case suffix == "download" && r.Method == http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		path := filepath.Join(dir, "site.tar.gz")
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			writeError(w, http.StatusNotFound, "Not found")
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", id+".tar.gz"))
		http.ServeFile(w, r, path)
	case suffix == "" && r.Method == http.MethodDelete:
		if !rejectBody(w, r) {
			return
		}
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			writeError(w, http.StatusNotFound, "Not found")
			return
		}
		if err := os.RemoveAll(dir); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "snapshot_delete", Changes: map[string]any{"id": id}})
		writeJSON(w, http.StatusOK, map[string]string{"message": "Snapshot deleted"})
	default:
		writeError(w, http.StatusNotFound, "Not found")
	}
}

func (s *APIServer) handleTokens(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.cfg.StateDir, ".api_tokens")
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		tokens, err := readAPITokens(path)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		out := []map[string]any{}
		for _, token := range tokens {
			out = append(out, map[string]any{"id": token.ID, "name": token.Name, "scopes": token.Scopes, "created_at": token.CreatedAt, "revoked_at": token.RevokedAt})
		}
		writeJSON(w, http.StatusOK, map[string]any{"tokens": out})
	case http.MethodPost:
		var body struct {
			Name   string   `json:"name"`
			Scopes []string `json:"scopes"`
		}
		if !decodeBody(w, r, &body, true) {
			return
		}
		name := strings.TrimSpace(body.Name)
		if name == "" || len(name) > 100 {
			writeValidation(w, "name", "invalid value")
			return
		}
		tokens, err := readAPITokens(path)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		raw, err := randomHex(32)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		token := "act_" + raw
		record := apiTokenRecord{ID: s.newID("tok"), Name: name, TokenHash: tokenHash(token), Scopes: normalizeTokenScopes(body.Scopes), CreatedAt: s.nowString()}
		tokens = append(tokens, record)
		if err := atomicWriteJSON(path, tokens, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": s.nowString(), "action": "token_issue", "id": record.ID})
		writeJSON(w, http.StatusOK, map[string]any{"id": record.ID, "token": token, "name": record.Name, "scopes": record.Scopes})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleTokenPath(w http.ResponseWriter, r *http.Request) {
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/tokens/"), "/")
	if !validSimpleID(id) {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if !method(w, r, http.MethodDelete) || !rejectBody(w, r) {
		return
	}
	path := filepath.Join(s.cfg.StateDir, ".api_tokens")
	tokens, err := readAPITokens(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	for i := range tokens {
		if tokens[i].ID == id {
			now := s.nowString()
			tokens[i].RevokedAt = &now
			if err := atomicWriteJSON(path, tokens, 0600); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": now, "action": "token_revoke", "id": id})
			writeJSON(w, http.StatusOK, map[string]string{"message": "Token revoked"})
			return
		}
	}
	writeError(w, http.StatusNotFound, "Not found")
}

func (s *APIServer) handleRuleFile(filename, logType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(s.cfg.StateDir, filename)
		switch r.Method {
		case http.MethodGet:
			if !rejectBody(w, r) {
				return
			}
			rules, err := readRuleList(path)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "State file is corrupted")
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"rules": rules})
		case http.MethodPost:
			var rule map[string]any
			if !decodeBody(w, r, &rule, true) {
				return
			}
			if len(rule) == 0 {
				writeValidation(w, "body", "required")
				return
			}
			rules, err := readRuleList(path)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "State file is corrupted")
				return
			}
			if duplicateRule(rules, rule) {
				writeError(w, http.StatusConflict, "Conflict")
				return
			}
			rule["id"] = s.newID("rule")
			rules = append(rules, rule)
			if err := atomicWriteJSON(path, rules, 0600); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: logType, Changes: map[string]any{"id": rule["id"]}})
			writeJSON(w, http.StatusOK, map[string]any{"message": "Rule created", "rule": rule})
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}

func (s *APIServer) handleRulePath(filename, logType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.Trim(strings.TrimPrefix(r.URL.Path, strings.TrimSuffix(r.URL.Path, filepath.Base(r.URL.Path))), "/")
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) == 0 {
			writeError(w, http.StatusNotFound, "Not found")
			return
		}
		id = parts[len(parts)-1]
		if !validSimpleID(id) {
			writeError(w, http.StatusNotFound, "Not found")
			return
		}
		if !method(w, r, http.MethodDelete) || !rejectBody(w, r) {
			return
		}
		path := filepath.Join(s.cfg.StateDir, filename)
		rules, err := readRuleList(path)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		next := []map[string]any{}
		found := false
		for _, rule := range rules {
			if fmt.Sprint(rule["id"]) == id {
				found = true
				continue
			}
			next = append(next, rule)
		}
		if !found {
			writeError(w, http.StatusNotFound, "Not found")
			return
		}
		if err := atomicWriteJSON(path, next, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: logType + "_delete", Changes: map[string]any{"id": id}})
		writeJSON(w, http.StatusOK, map[string]string{"message": "Rule deleted"})
	}
}

func (s *APIServer) handleHookPath(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/hooks/"), "/")
	id, suffix, hasSuffix := strings.Cut(rest, "/")
	if !validSimpleID(id) {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if hasSuffix && suffix == "log" {
		if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
			return
		}
		runs := []map[string]any{}
		root := filepath.Join(s.cfg.StateDir, ".build_logs")
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), "_hook_"+id+".json") {
				return nil
			}
			if value, ok, err := readOptionalJSONMap(path); err == nil && ok {
				runs = append(runs, value)
			}
			return nil
		})
		sort.Slice(runs, func(i, j int) bool { return fmt.Sprint(runs[i]["at"]) > fmt.Sprint(runs[j]["at"]) })
		if len(runs) > 20 {
			runs = runs[:20]
		}
		writeJSON(w, http.StatusOK, map[string]any{"id": id, "runs": runs})
		return
	}
	if hasSuffix {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	s.handleRulePath(".hooks", "hook")(w, r)
}

func (s *APIServer) handleVerifyOutput(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	logs := s.readBuildLogsNewest()
	if len(logs) == 0 {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	expected := firstNonEmpty(logs[0].OutputSHA256, logs[0].SHA256)
	if expected == "" {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	meta, err := inspectOutput(s.outputDir())
	if err == nil && !meta.Exists {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	actual := meta.SHA256
	writeJSON(w, http.StatusOK, map[string]any{"match": expected == actual, "expected": expected, "actual": actual})
}

func (s *APIServer) handlePipelineConfig(w http.ResponseWriter, r *http.Request) {
	s.handleGenericConfigFile(w, r, ".pipeline_config", "pipeline_config", validatePipelineConfig)
}

func (s *APIServer) handleNotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		value, _, _ := readOptionalJSONMap(filepath.Join(s.cfg.StateDir, ".notes"))
		if value == nil {
			value = map[string]any{"content": ""}
		}
		writeJSON(w, http.StatusOK, value)
	case http.MethodPost:
		var body map[string]any
		if !decodeBody(w, r, &body, true) {
			return
		}
		content, _ := body["content"].(string)
		if len(content) > 20000 {
			writeValidation(w, "content", "too long")
			return
		}
		current, _, _ := readOptionalJSONMap(filepath.Join(s.cfg.StateDir, ".notes"))
		if current != nil && fmt.Sprint(current["content"]) == content {
			writeJSON(w, http.StatusOK, map[string]string{"message": "No changes"})
			return
		}
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".notes"), map[string]any{"content": content, "updated_at": s.nowString()}, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "notes", Changes: map[string]any{"content_changed": true}})
		writeJSON(w, http.StatusOK, map[string]string{"message": "Notes updated"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleDashboardLayout(w http.ResponseWriter, r *http.Request) {
	s.handleGenericConfigFile(w, r, ".dashboard_layout", "dashboard_layout", validateDashboardLayout)
}

func (s *APIServer) handleSMTPConfig(w http.ResponseWriter, r *http.Request) {
	s.handleGenericConfigFile(w, r, ".smtp_config", "smtp_config", nil)
}

func (s *APIServer) handleSMTPTest(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	cfg, _, _ := readOptionalJSONMap(filepath.Join(s.cfg.StateDir, ".smtp_config"))
	if cfg == nil || !boolFromAny(cfg["enabled"]) {
		writeValidation(w, "smtp", "disabled")
		return
	}
	_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".notify_log"), map[string]any{"at": s.nowString(), "event": "smtp_test", "result": "success"})
	writeJSON(w, http.StatusOK, map[string]string{"message": "SMTP test sent"})
}

func (s *APIServer) handleGenericConfigFile(w http.ResponseWriter, r *http.Request, filename, logType string, validate func(http.ResponseWriter, map[string]any) bool) {
	path := filepath.Join(s.cfg.StateDir, filename)
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		value, _, err := readOptionalJSONMap(path)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if value == nil {
			value = map[string]any{}
		}
		writeJSON(w, http.StatusOK, maskSecrets(value))
	case http.MethodPost:
		var value map[string]any
		if !decodeBody(w, r, &value, true) {
			return
		}
		if validate != nil && !validate(w, value) {
			return
		}
		current, _, _ := readOptionalJSONMap(path)
		if mapsEqual(current, value) {
			writeJSON(w, http.StatusOK, map[string]string{"message": "No changes"})
			return
		}
		if err := atomicWriteJSON(path, value, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: logType, Changes: maskSecrets(value)})
		writeJSON(w, http.StatusOK, map[string]string{"message": "Config updated"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleAPIAccessLog(w http.ResponseWriter, r *http.Request) {
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
	methodFilter := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("method")))
	if methodFilter != "" && !validHTTPMethod(methodFilter) {
		writeValidation(w, "method", "invalid value")
		return
	}
	pathFilter := r.URL.Query().Get("path")
	statusFilter := 0
	if raw := r.URL.Query().Get("status"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 100 || value > 599 {
			writeValidation(w, "status", "out of range")
			return
		}
		statusFilter = value
	}
	records := readJSONLines(filepath.Join(s.cfg.StateDir, ".api_access_log"))
	filtered := []map[string]any{}
	for _, record := range records {
		if methodFilter != "" && strings.ToUpper(fmt.Sprint(record["method"])) != methodFilter {
			continue
		}
		if pathFilter != "" && !strings.HasPrefix(fmt.Sprint(record["path"]), pathFilter) {
			continue
		}
		if statusFilter != 0 && intFromAny(record["status"]) != statusFilter {
			continue
		}
		filtered = append(filtered, record)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return fmt.Sprint(filtered[i]["at"]) > fmt.Sprint(filtered[j]["at"])
	})
	total := len(filtered)
	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}
	writeJSON(w, http.StatusOK, map[string]any{"log": filtered[start:end], "total": total})
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
	session, ok := s.sessions[tokenHash(token)]
	if !ok {
		s.mu.Unlock()
		return s.validAPIToken(token)
	}
	expires, err := time.Parse(apiTimeLayout, session.ExpiresAt)
	if err != nil || !s.cfg.Now().UTC().Before(expires) {
		delete(s.sessions, tokenHash(token))
		s.mu.Unlock()
		return false
	}
	s.mu.Unlock()
	return true
}

func (s *APIServer) validAPIToken(token string) bool {
	if token == "" {
		return false
	}
	tokens, err := readAPITokens(filepath.Join(s.cfg.StateDir, ".api_tokens"))
	if err != nil {
		return false
	}
	hash := tokenHash(token)
	for _, record := range tokens {
		if record.TokenHash == hash && record.RevokedAt == nil {
			return true
		}
	}
	return false
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

func (s *APIServer) readRepoConfig() (apiRepoConfig, error) {
	cfg := defaultRepoConfig()
	path := filepath.Join(s.cfg.StateDir, ".repo_config")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	var fileCfg apiRepoConfig
	if err := readJSONFile(path, &fileCfg); err != nil {
		return cfg, err
	}
	if fileCfg.Owner != "" {
		cfg.Owner = fileCfg.Owner
	}
	if fileCfg.Repo != "" {
		cfg.Repo = fileCfg.Repo
	}
	if fileCfg.Branch != "" {
		cfg.Branch = fileCfg.Branch
	}
	if fileCfg.TargetFile != "" {
		cfg.TargetFile = fileCfg.TargetFile
	}
	cfg.UpdatedAt = fileCfg.UpdatedAt
	return cfg, nil
}

func (s *APIServer) readBranchConfig() (apiBranchConfigFile, string, error) {
	path := filepath.Join(s.cfg.StateDir, ".branch_config")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return apiBranchConfigFile{BranchTargets: []apiBranchTarget{defaultBranchTarget(s.cfg.StateDir)}}, "default", nil
	}
	var cfg apiBranchConfigFile
	if err := readJSONFile(path, &cfg); err != nil {
		return cfg, "file", err
	}
	if cfg.BranchTargets == nil {
		cfg.BranchTargets = []apiBranchTarget{}
	}
	return cfg, "file", nil
}

func (s *APIServer) readNotifyConfig() (apiNotifyConfig, error) {
	path := filepath.Join(s.cfg.StateDir, ".notify_config")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return defaultNotifyConfig(), nil
	}
	var cfg apiNotifyConfig
	if err := readJSONFile(path, &cfg); err != nil {
		return cfg, err
	}
	normalized, ok := normalizeNotifyConfig(nil, cfg)
	if !ok {
		return cfg, errors.New("invalid notify config")
	}
	return normalized, nil
}

func (s *APIServer) readAccessControl() (apiAccessControl, error) {
	path := filepath.Join(s.cfg.StateDir, ".access_control")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return apiAccessControl{Allow: []string{}}, nil
	}
	var cfg apiAccessControl
	if err := readJSONFile(path, &cfg); err != nil {
		return cfg, err
	}
	normalized, ok := normalizeAccessControl(nil, cfg)
	if !ok {
		return cfg, errors.New("invalid access control")
	}
	return normalized, nil
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

func randomBase32Secret() (string, error) {
	raw := make([]byte, 20)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw), nil
}

func readTOTPSecret(path string) (apiTOTPSecret, error) {
	var secret apiTOTPSecret
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return secret, nil
	}
	if err != nil {
		return secret, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return secret, nil
	}
	if err := json.Unmarshal(data, &secret); err != nil {
		return secret, err
	}
	return secret, nil
}

func totpStatusPayload(secret apiTOTPSecret) map[string]any {
	return map[string]any{"enabled": secret.Enabled, "confirmed_at": secret.ConfirmedAt}
}

func verifyTOTPCode(secretBase32, code string, at time.Time) bool {
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return false
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return false
		}
	}
	step := at.UTC().Unix() / 30
	for _, offset := range []int64{-1, 0, 1} {
		if totpCode(secretBase32, step+offset) == code {
			return true
		}
	}
	return false
}

func totpCode(secretBase32 string, step int64) string {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(strings.TrimSpace(secretBase32)))
	if err != nil {
		return ""
	}
	var counter [8]byte
	binary.BigEndian.PutUint64(counter[:], uint64(step))
	mac := hmac.New(sha1.New, key)
	_, _ = mac.Write(counter[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	value := (int(sum[offset])&0x7f)<<24 |
		(int(sum[offset+1])&0xff)<<16 |
		(int(sum[offset+2])&0xff)<<8 |
		(int(sum[offset+3]) & 0xff)
	return fmt.Sprintf("%06d", value%1_000_000)
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

func readMaintenanceState(path string) (apiMaintenanceState, error) {
	var state apiMaintenanceState
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
		if json.Unmarshal([]byte(line), &record) == nil && record.ID != "" && record.Status != "" {
			records = append(records, record)
		} else {
			fmt.Fprintf(os.Stderr, "BUILD_HISTORY_SKIP_CORRUPT: path=%s\n", path)
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return firstNonEmpty(records[i].FinishedAt, records[i].StartedAt, records[i].BuildAt, records[i].ID) > firstNonEmpty(records[j].FinishedAt, records[j].StartedAt, records[j].BuildAt, records[j].ID)
	})
	return records
}

func filterHistory(w http.ResponseWriter, r *http.Request, records []apiHistoryRecord) []apiHistoryRecord {
	trigger := strings.TrimSpace(r.URL.Query().Get("trigger"))
	tag := strings.TrimSpace(r.URL.Query().Get("tag"))
	flaggedRaw := strings.TrimSpace(r.URL.Query().Get("flagged"))
	var flagged *bool
	if flaggedRaw != "" {
		switch strings.ToLower(flaggedRaw) {
		case "true":
			value := true
			flagged = &value
		case "false":
			value := false
			flagged = &value
		default:
			writeValidation(w, "flagged", "invalid value")
			return nil
		}
	}
	out := []apiHistoryRecord{}
	for _, record := range records {
		if trigger != "" && record.Trigger != trigger {
			continue
		}
		if tag != "" && !containsString(record.Tags, tag) {
			continue
		}
		if flagged != nil && record.Flagged != *flagged {
			continue
		}
		out = append(out, record)
	}
	return out
}

func updateHistoryRecord(path, id string, update func(*apiHistoryRecord)) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	records := []apiHistoryRecord{}
	found := false
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record apiHistoryRecord
		if json.Unmarshal([]byte(line), &record) != nil {
			continue
		}
		if record.ID == id {
			update(&record)
			found = true
		}
		records = append(records, record)
	}
	if !found {
		return nil
	}
	return atomicWriteHistory(path, records, 0600)
}

func readBuildLog(path string) (apiBuildLog, error) {
	var log apiBuildLog
	err := readJSONFile(path, &log)
	return log, err
}

func readBuildLogByID(stateDir, id string) (apiBuildLog, error) {
	log, err := readBuildLog(filepath.Join(stateDir, ".build_logs", id+".json"))
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		return log, err
	}
	archivePath := filepath.Join(stateDir, ".build_logs", "archive", id+".json.gz")
	file, archiveErr := os.Open(archivePath)
	if archiveErr != nil {
		return apiBuildLog{}, archiveErr
	}
	defer file.Close()
	reader, archiveErr := gzip.NewReader(file)
	if archiveErr != nil {
		return apiBuildLog{}, archiveErr
	}
	defer reader.Close()
	var archived apiBuildLog
	if archiveErr := json.NewDecoder(reader).Decode(&archived); archiveErr != nil {
		return apiBuildLog{}, archiveErr
	}
	return archived, nil
}

func defaultServerConfig() apiServerConfig {
	return apiServerConfig{
		LogMaxLines: 500, HistoryMaxCount: 100, BuildTimeoutSeconds: 300,
		LogRetentionDays: 30, LogArchiveAfterDays: 0, LogLevel: "INFO",
		SnapshotsKeep: 5, QueueMaxSize: defaultQueueMaxSize,
		BuildRetryMax: 0, BuildRetryBaseSeconds: 5,
		CommitStatusEnabled: false, CommitStatusContext: "Adlaire CI",
		BuildTrendKeepCount: 1000, SessionTimeoutSeconds: 28800,
		ScheduleIntervalSeconds: 300,
	}
}

func defaultRepoConfig() apiRepoConfig {
	owner, repo := apiDefaultRepo()
	return apiRepoConfig{Owner: owner, Repo: repo, Branch: "main", TargetFile: "docs"}
}

func apiDefaultRepo() (string, string) {
	repo := os.Getenv("ADLAIRE_CI_REPOSITORY")
	if repo == "" {
		repo = "fqwink/Build-Scripts"
	}
	owner, name, ok := strings.Cut(repo, "/")
	if !ok || owner == "" || name == "" {
		return "fqwink", "Build-Scripts"
	}
	return owner, name
}

func defaultBranchTarget(stateDir string) apiBranchTarget {
	return apiBranchTarget{
		Branch:     "main",
		TargetFile: "docs",
		SHAFile:    filepath.Join(stateDir, ".last_sha"),
		Src:        filepath.Join(stateDir, "repo", "docs"),
		Out:        filepath.Join(stateDir, "dist", "site"),
	}
}

func defaultNotifyConfig() apiNotifyConfig {
	return apiNotifyConfig{
		Webhooks: []apiNotifyWebhook{},
		Channels: []apiNotifyChannel{},
		On:       []string{},
		Summary:  apiNotifySummary{Enabled: false, Interval: "weekly", Hour: 9, DayOfWeek: 1},
		Email:    apiNotifyEmail{Enabled: false, To: []string{}, On: []string{}},
	}
}

func normalizeNotifyConfig(w http.ResponseWriter, cfg apiNotifyConfig) (apiNotifyConfig, bool) {
	if cfg.Webhooks == nil {
		cfg.Webhooks = []apiNotifyWebhook{}
	}
	if cfg.Channels == nil {
		cfg.Channels = []apiNotifyChannel{}
	}
	if cfg.On == nil {
		cfg.On = []string{}
	}
	on, ok := normalizeNotifyEvents(w, "on", cfg.On, false)
	if !ok {
		return cfg, false
	}
	cfg.On = on
	if cfg.Summary.Interval == "" {
		cfg.Summary.Interval = "weekly"
	}
	if cfg.Summary.Interval != "weekly" {
		writeValidationIf(w, "summary.interval", "must be weekly")
		return cfg, false
	}
	if cfg.Summary.Hour < 0 || cfg.Summary.Hour > 23 {
		writeValidationIf(w, "summary.hour", "out of range")
		return cfg, false
	}
	if cfg.Summary.DayOfWeek < 0 || cfg.Summary.DayOfWeek > 6 {
		writeValidationIf(w, "summary.day_of_week", "out of range")
		return cfg, false
	}
	for i := range cfg.Webhooks {
		webhook := &cfg.Webhooks[i]
		webhook.URL = strings.TrimSpace(webhook.URL)
		if !validateHTTPURL(webhook.URL) {
			writeValidationIf(w, fmt.Sprintf("webhooks[%d].url", i), "invalid url")
			return cfg, false
		}
		webhook.Label = strings.TrimSpace(webhook.Label)
		if len(webhook.Label) > 64 {
			writeValidationIf(w, fmt.Sprintf("webhooks[%d].label", i), "too long")
			return cfg, false
		}
		events, ok := normalizeNotifyEvents(w, fmt.Sprintf("webhooks[%d].on", i), webhook.On, true)
		if !ok {
			return cfg, false
		}
		webhook.On = events
		if webhook.PayloadTemplate != nil && len(*webhook.PayloadTemplate) > 10000 {
			writeValidationIf(w, fmt.Sprintf("webhooks[%d].payload_template", i), "too long")
			return cfg, false
		}
		if webhook.RetryCount < 0 || webhook.RetryCount > 10 {
			writeValidationIf(w, fmt.Sprintf("webhooks[%d].retry_count", i), "out of range")
			return cfg, false
		}
		if webhook.RetryIntervalSeconds == 0 {
			webhook.RetryIntervalSeconds = 30
		}
		if webhook.RetryIntervalSeconds < 1 || webhook.RetryIntervalSeconds > 3600 {
			writeValidationIf(w, fmt.Sprintf("webhooks[%d].retry_interval_seconds", i), "out of range")
			return cfg, false
		}
		if webhook.Secret != nil {
			secret := strings.TrimSpace(*webhook.Secret)
			if secret == "" || len(secret) > 256 {
				writeValidationIf(w, fmt.Sprintf("webhooks[%d].secret", i), "invalid value")
				return cfg, false
			}
			webhook.Secret = &secret
		}
	}
	seenChannels := map[string]bool{}
	for i := range cfg.Channels {
		channel := &cfg.Channels[i]
		channel.ID = strings.TrimSpace(channel.ID)
		if channel.ID == "" {
			channel.ID = fmt.Sprintf("n%03d", i+1)
		}
		if !validNotifyChannelID(channel.ID) || seenChannels[channel.ID] {
			writeValidationIf(w, fmt.Sprintf("channels[%d].id", i), "invalid value")
			return cfg, false
		}
		seenChannels[channel.ID] = true
		channel.Type = strings.TrimSpace(channel.Type)
		if channel.Type != "webhook" && channel.Type != "email" && channel.Type != "command" {
			writeValidationIf(w, fmt.Sprintf("channels[%d].type", i), "invalid value")
			return cfg, false
		}
		channel.Label = strings.TrimSpace(channel.Label)
		if len(channel.Label) > 64 {
			writeValidationIf(w, fmt.Sprintf("channels[%d].label", i), "too long")
			return cfg, false
		}
		events, ok := normalizeNotifyEvents(w, fmt.Sprintf("channels[%d].on", i), channel.On, true)
		if !ok {
			return cfg, false
		}
		channel.On = events
		if channel.Config == nil {
			channel.Config = map[string]any{}
		}
		if channel.Type == "webhook" {
			rawURL, _ := channel.Config["url"].(string)
			rawURL = strings.TrimSpace(rawURL)
			if !validateHTTPURL(rawURL) {
				writeValidationIf(w, fmt.Sprintf("channels[%d].config.url", i), "invalid url")
				return cfg, false
			}
			channel.Config["url"] = rawURL
		}
		if channel.RetryCount < 0 || channel.RetryCount > 10 {
			writeValidationIf(w, fmt.Sprintf("channels[%d].retry_count", i), "out of range")
			return cfg, false
		}
		if channel.RetryIntervalSeconds == 0 {
			channel.RetryIntervalSeconds = 30
		}
		if channel.RetryIntervalSeconds < 1 || channel.RetryIntervalSeconds > 3600 {
			writeValidationIf(w, fmt.Sprintf("channels[%d].retry_interval_seconds", i), "out of range")
			return cfg, false
		}
	}
	if cfg.Email.To == nil {
		cfg.Email.To = []string{}
	}
	if len(cfg.Email.To) > 50 {
		writeValidationIf(w, "email.to", "too many")
		return cfg, false
	}
	if cfg.Email.On == nil {
		cfg.Email.On = []string{}
	}
	emailOn, ok := normalizeNotifyEmailEvents(w, cfg.Email.On)
	if !ok {
		return cfg, false
	}
	cfg.Email.On = emailOn
	for i, addr := range cfg.Email.To {
		cfg.Email.To[i] = strings.TrimSpace(addr)
		if cfg.Email.To[i] == "" || !strings.Contains(cfg.Email.To[i], "@") {
			writeValidationIf(w, fmt.Sprintf("email.to[%d]", i), "invalid email")
			return cfg, false
		}
	}
	return cfg, true
}

func writeValidationIf(w http.ResponseWriter, field, message string) {
	if w != nil {
		writeValidation(w, field, message)
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
	if cfg.ScheduleIntervalSeconds == 0 {
		cfg.ScheduleIntervalSeconds = def.ScheduleIntervalSeconds
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

func decodeTrimmedString(w http.ResponseWriter, key string, raw json.RawMessage, min, max int) (string, bool) {
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		writeValidation(w, key, "invalid value")
		return "", false
	}
	value = strings.TrimSpace(value)
	if len(value) < min || len(value) > max {
		writeValidation(w, key, "out of range")
		return "", false
	}
	return value, true
}

func validateBranchName(w http.ResponseWriter, value string) bool {
	if strings.HasPrefix(value, "refs/heads/") || strings.Contains(value, "..") || strings.Contains(value, "~") || containsControl(value) {
		writeValidation(w, "branch", "invalid value")
		return false
	}
	return true
}

func validateTargetFile(w http.ResponseWriter, value string) bool {
	if filepath.IsAbs(value) || strings.Contains(value, "..") || strings.HasPrefix(value, "/") || containsControl(value) {
		writeValidation(w, "target_file", "invalid value")
		return false
	}
	return true
}

func validateBranchTargets(w http.ResponseWriter, targets []apiBranchTarget) bool {
	if len(targets) > 50 {
		writeValidation(w, "branches", "must contain 50 items or less")
		return false
	}
	for i, target := range targets {
		field := fmt.Sprintf("branches[%d]", i)
		if target.Branch == "" || len(target.Branch) > 128 || strings.HasPrefix(target.Branch, "refs/heads/") || strings.Contains(target.Branch, "..") || strings.Contains(target.Branch, "~") || containsControl(target.Branch) {
			writeValidation(w, field+".branch", "invalid value")
			return false
		}
		if target.TargetFile == "" || len(target.TargetFile) > 500 || filepath.IsAbs(target.TargetFile) || strings.Contains(target.TargetFile, "..") || containsControl(target.TargetFile) {
			writeValidation(w, field+".target_file", "invalid value")
			return false
		}
		for _, pair := range []struct {
			name  string
			value string
		}{{"sha_file", target.SHAFile}, {"src", target.Src}, {"out", target.Out}} {
			if pair.value == "" || !filepath.IsAbs(pair.value) || containsControl(pair.value) {
				writeValidation(w, field+"."+pair.name, "invalid value")
				return false
			}
		}
		if len(target.DeployTargets) > 20 {
			writeValidation(w, field+".deploy_targets", "must contain 20 items or less")
			return false
		}
		for j, deploy := range target.DeployTargets {
			deployField := fmt.Sprintf("%s.deploy_targets[%d]", field, j)
			if deploy.Host == "" || len(deploy.Host) > 255 || containsControl(deploy.Host) {
				writeValidation(w, deployField+".host", "invalid value")
				return false
			}
			if deploy.User == "" || len(deploy.User) > 64 || containsControl(deploy.User) {
				writeValidation(w, deployField+".user", "invalid value")
				return false
			}
			if deploy.DestDir == "" || !filepath.IsAbs(deploy.DestDir) || containsControl(deploy.DestDir) {
				writeValidation(w, deployField+".dest_dir", "invalid value")
				return false
			}
		}
	}
	return true
}

func containsControl(value string) bool {
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}

func branchTargetsEqual(a, b []apiBranchTarget) bool {
	aj, errA := json.Marshal(a)
	bj, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aj) == string(bj)
}

func schedulePayload(cfg apiServerConfig) map[string]any {
	return map[string]any{
		"next_run_at":                nil,
		"interval":                   formatInterval(cfg.ScheduleIntervalSeconds),
		"interval_seconds":           cfg.ScheduleIntervalSeconds,
		"paused":                     cfg.SchedulePaused,
		"allowed_hours":              cfg.AllowedHours,
		"force_build_interval_hours": cfg.ForceBuildIntervalHours,
		"build_cooldown_seconds":     cfg.BuildCooldownSeconds,
	}
}

func formatInterval(seconds int) string {
	if seconds%3600 == 0 {
		return fmt.Sprintf("%dh", seconds/3600)
	}
	if seconds%60 == 0 {
		return fmt.Sprintf("%dmin", seconds/60)
	}
	return fmt.Sprintf("%ds", seconds)
}

func allowedHoursEqual(a, b *apiAllowedHours) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.From == b.From && a.To == b.To
}

func normalizeNotifyEvents(w http.ResponseWriter, field string, values []string, allowWildcard bool) ([]string, bool) {
	allowed := map[string]bool{
		"start": true, "success": true, "failure": true, "deploy_failure": true,
		"weekly_summary": true, "approval_required": true, "duration_anomaly": true, "config_corrupt": true,
	}
	if allowWildcard {
		allowed["*"] = true
	}
	out := []string{}
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if !allowed[value] {
			writeValidationIf(w, field, "invalid event")
			return nil, false
		}
		if !seen[value] {
			out = append(out, value)
			seen[value] = true
		}
	}
	return out, true
}

func normalizeNotifyEmailEvents(w http.ResponseWriter, values []string) ([]string, bool) {
	allowed := map[string]bool{"start": true, "success": true, "failure": true, "duration_anomaly": true}
	out := []string{}
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if !allowed[value] {
			writeValidationIf(w, "email.on", "invalid event")
			return nil, false
		}
		if !seen[value] {
			out = append(out, value)
			seen[value] = true
		}
	}
	return out, true
}

func validateHTTPURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func validNotifyChannelID(id string) bool {
	if len(id) < 1 || len(id) > 64 {
		return false
	}
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func maskNotifyConfig(cfg apiNotifyConfig) apiNotifyConfig {
	for i := range cfg.Webhooks {
		if cfg.Webhooks[i].Secret != nil {
			masked := "***"
			cfg.Webhooks[i].Secret = &masked
		}
	}
	return cfg
}

func notifyConfigsEqual(a, b apiNotifyConfig) bool {
	aj, errA := json.Marshal(a)
	bj, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aj) == string(bj)
}

func normalizeAccessControl(w http.ResponseWriter, cfg apiAccessControl) (apiAccessControl, bool) {
	if cfg.Allow == nil {
		cfg.Allow = []string{}
	}
	if len(cfg.Allow) > 100 {
		writeValidationIf(w, "allow", "too many")
		return cfg, false
	}
	seen := map[string]bool{}
	out := []string{}
	for i, raw := range cfg.Allow {
		value := strings.TrimSpace(raw)
		if value == "" {
			writeValidationIf(w, fmt.Sprintf("allow[%d]", i), "empty value")
			return cfg, false
		}
		normalized, ok := normalizeIPv4OrCIDR(value)
		if !ok {
			writeValidationIf(w, fmt.Sprintf("allow[%d]", i), "invalid value")
			return cfg, false
		}
		if !seen[normalized] {
			out = append(out, normalized)
			seen[normalized] = true
		}
	}
	sort.Strings(out)
	return apiAccessControl{Allow: out}, true
}

func normalizeIPv4OrCIDR(value string) (string, bool) {
	if strings.Contains(value, "/") {
		ip, network, err := net.ParseCIDR(value)
		if err != nil || ip == nil || ip.To4() == nil || network == nil {
			return "", false
		}
		ones, bits := network.Mask.Size()
		if bits != 32 || ones < 0 || ones > 32 {
			return "", false
		}
		network.IP = network.IP.To4()
		return network.String(), true
	}
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() == nil {
		return "", false
	}
	return ip.To4().String(), true
}

func (s *APIServer) accessAllowed(r *http.Request) bool {
	if r.Method == http.MethodGet && r.URL.Path == "/api/health" {
		return true
	}
	cfg, err := s.readAccessControl()
	if err != nil || len(cfg.Allow) == 0 {
		return err == nil
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(strings.TrimSpace(host))
	if ip == nil || ip.To4() == nil {
		return false
	}
	ip = ip.To4()
	for _, entry := range cfg.Allow {
		if strings.Contains(entry, "/") {
			_, network, err := net.ParseCIDR(entry)
			if err == nil && network.Contains(ip) {
				return true
			}
			continue
		}
		allowed := net.ParseIP(entry)
		if allowed != nil && allowed.To4() != nil && allowed.To4().Equal(ip) {
			return true
		}
	}
	return false
}

func stringSlicesEqual(a, b []string) bool {
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

func notifyConfigChangeSummary(cfg apiNotifyConfig) map[string]any {
	return map[string]any{
		"webhooks_count": len(cfg.Webhooks),
		"channels_count": len(cfg.Channels),
		"on":             cfg.On,
		"summary":        cfg.Summary,
		"email_enabled":  cfg.Email.Enabled,
	}
}

type apiNotifyTarget struct {
	URL    string
	Secret *string
}

func notifyWebhookTargets(cfg apiNotifyConfig, event string) []apiNotifyTarget {
	targets := []apiNotifyTarget{}
	for _, channel := range cfg.Channels {
		if channel.Type != "webhook" || !channel.Enabled || !notifyEventEnabled(event, cfg.On, channel.On) {
			continue
		}
		rawURL, _ := channel.Config["url"].(string)
		if rawURL != "" {
			targets = append(targets, apiNotifyTarget{URL: rawURL})
		}
	}
	if len(cfg.Channels) > 0 {
		return targets
	}
	for _, webhook := range cfg.Webhooks {
		if !webhook.Enabled || !notifyEventEnabled(event, cfg.On, webhook.On) {
			continue
		}
		targets = append(targets, apiNotifyTarget{URL: webhook.URL, Secret: webhook.Secret})
	}
	return targets
}

func notifyEventEnabled(event string, defaults, specific []string) bool {
	values := specific
	if len(values) == 0 {
		values = defaults
	}
	for _, value := range values {
		if value == "*" || value == event {
			return true
		}
	}
	return false
}

func sendWebhook(rawURL string, secret *string, payload map[string]any) (int, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if secret != nil {
		mac := hmac.New(sha256.New, []byte(*secret))
		_, _ = mac.Write(body)
		req.Header.Set("X-Adlaire-Signature", "sha256="+hex.EncodeToString(mac.Sum(nil)))
	}
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("webhook returned %d", resp.StatusCode)
	}
	return resp.StatusCode, nil
}

func (s *APIServer) weeklySummary() map[string]any {
	now := s.cfg.Now().UTC()
	from := now.AddDate(0, 0, -7)
	successCount := 0
	failureCount := 0
	durationCount := 0
	var durationTotal int64
	for _, record := range readHistory(filepath.Join(s.cfg.StateDir, ".build_history")) {
		at := firstNonEmpty(record.FinishedAt, record.StartedAt, record.BuildAt)
		if at == "" {
			continue
		}
		parsed, err := time.Parse(apiTimeLayout, at)
		if err != nil || parsed.Before(from) || parsed.After(now) {
			continue
		}
		switch record.Status {
		case "success":
			successCount++
		case "failure", "cancelled", "hook_error":
			failureCount++
		default:
			continue
		}
		if record.DurationSeconds != 0 {
			durationTotal += record.DurationSeconds
			durationCount++
		}
	}
	total := successCount + failureCount
	successRate := 0.0
	if total > 0 {
		successRate = math.Round((float64(successCount)/float64(total))*10000) / 100
	}
	var avgDuration any
	if durationCount > 0 {
		avgDuration = math.Round(float64(durationTotal)/float64(durationCount)*100) / 100
	}
	return map[string]any{
		"period":               from.Format("2006-01-02") + "/" + now.Format("2006-01-02"),
		"success_count":        successCount,
		"failure_count":        failureCount,
		"success_rate":         successRate,
		"avg_duration_seconds": avgDuration,
	}
}

func maskRepoConfigChanges(patch map[string]json.RawMessage) map[string]any {
	changes := map[string]any{}
	for key, raw := range patch {
		var value any
		if json.Unmarshal(raw, &value) == nil {
			changes[key] = value
		}
	}
	return changes
}

func removeIfExists(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func atomicWriteText(path, value string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(value), mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func cleanupLogFiles(root string, cutoff time.Time) (int, int) {
	deleted, failed := 0, 0
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry == nil || entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, ".json.gz") {
			return nil
		}
		info, statErr := entry.Info()
		if statErr != nil || info.ModTime().After(cutoff) {
			return nil
		}
		if err := os.Remove(path); err != nil {
			failed++
		} else {
			deleted++
		}
		return nil
	})
	return deleted, failed
}

func archiveLogFiles(root string, cutoff time.Time) (int, error) {
	archiveDir := filepath.Join(root, "archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return 0, err
	}
	count := 0
	entries, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		src := filepath.Join(root, entry.Name())
		info, err := entry.Info()
		if err != nil || info.ModTime().After(cutoff) {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".json")
		if _, err := readBuildLog(src); err != nil {
			fmt.Fprintf(os.Stderr, "LOG_ARCHIVE_SKIP_CORRUPT: id=%s\n", id)
			continue
		}
		dest := filepath.Join(archiveDir, id+".json.gz")
		if _, err := os.Stat(dest); err == nil {
			continue
		}
		if err := gzipFile(src, dest); err != nil {
			return count, err
		}
		if err := os.Remove(src); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func gzipFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	writer := gzip.NewWriter(out)
	if _, err := io.Copy(writer, in); err != nil {
		_ = writer.Close()
		return err
	}
	return writer.Close()
}

func backupConfigFiles() []string {
	return []string{".server_config", ".notify_config", ".repo_config", ".branch_config", ".access_control", ".hooks", ".alert_rules", ".tag_rules", ".pipeline_config", ".dashboard_layout", ".smtp_config"}
}

func readOptionalJSONMap(path string) (map[string]any, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, true, err
	}
	return value, true, nil
}

func readSecretText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func maskSecrets(value map[string]any) map[string]any {
	out := map[string]any{}
	for key, v := range value {
		lower := strings.ToLower(key)
		if strings.Contains(lower, "secret") || strings.Contains(lower, "password") || strings.Contains(lower, "token") {
			if v != nil && fmt.Sprint(v) != "" {
				out[key] = "***"
			} else {
				out[key] = v
			}
			continue
		}
		if child, ok := v.(map[string]any); ok {
			out[key] = maskSecrets(child)
			continue
		}
		if list, ok := v.([]any); ok {
			next := make([]any, len(list))
			for i, item := range list {
				if child, ok := item.(map[string]any); ok {
					next[i] = maskSecrets(child)
				} else {
					next[i] = item
				}
			}
			out[key] = next
			continue
		}
		out[key] = v
	}
	return out
}

func containsMaskedSecret(value map[string]any) bool {
	for _, v := range value {
		if v == "***" {
			return true
		}
		if child, ok := v.(map[string]any); ok && containsMaskedSecret(child) {
			return true
		}
	}
	return false
}

func mergeMaskedSecrets(next, current map[string]any) map[string]any {
	out := map[string]any{}
	for key, value := range next {
		if value == "***" && current != nil {
			out[key] = current[key]
			continue
		}
		if child, ok := value.(map[string]any); ok {
			curChild, _ := current[key].(map[string]any)
			out[key] = mergeMaskedSecrets(child, curChild)
			continue
		}
		out[key] = value
	}
	return out
}

func fileModeString(info os.FileInfo) any {
	if info == nil {
		return nil
	}
	return fmt.Sprintf("%04o", info.Mode().Perm())
}

func statusFromExists(err error) string {
	if err == nil {
		return "ok"
	}
	if errors.Is(err, os.ErrNotExist) {
		return "warn"
	}
	return "error"
}

func fileTreeStats(root string, include func(string) bool) (int64, int) {
	var bytes int64
	count := 0
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry == nil || entry.IsDir() || !include(path) {
			return nil
		}
		if info, err := entry.Info(); err == nil {
			bytes += info.Size()
			count++
		}
		return nil
	})
	return bytes, count
}

func validWebhookSignature(secret string, body []byte, header string) bool {
	if !strings.HasPrefix(header, "sha256=") {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(header))
}

func validSimpleID(id string) bool {
	if id == "" || len(id) > 128 || containsControl(id) || strings.Contains(id, "..") || strings.ContainsAny(id, `/\`) {
		return false
	}
	return true
}

func readAPITokens(path string) ([]apiTokenRecord, error) {
	tokens := []apiTokenRecord{}
	if err := readJSONIfExists(path, &tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

func normalizeTokenScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return []string{"read", "write"}
	}
	out := []string{}
	seen := map[string]bool{}
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" || seen[scope] {
			continue
		}
		seen[scope] = true
		out = append(out, scope)
	}
	return out
}

func readRuleList(path string) ([]map[string]any, error) {
	rules := []map[string]any{}
	if err := readJSONIfExists(path, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func duplicateRule(rules []map[string]any, rule map[string]any) bool {
	next := mapWithoutID(rule)
	for _, existing := range rules {
		if mapsEqual(mapWithoutID(existing), next) {
			return true
		}
	}
	return false
}

func mapWithoutID(value map[string]any) map[string]any {
	out := map[string]any{}
	for key, v := range value {
		if key != "id" {
			out[key] = v
		}
	}
	return out
}

func validatePipelineConfig(w http.ResponseWriter, value map[string]any) bool {
	raw, ok := value["extra_args"].([]any)
	if !ok {
		return true
	}
	reserved := map[string]bool{"--src": true, "--out": true, "--state-dir": true}
	for _, item := range raw {
		if reserved[fmt.Sprint(item)] {
			writeValidation(w, "extra_args", "reserved arg")
			return false
		}
	}
	return true
}

func validateDashboardLayout(w http.ResponseWriter, value map[string]any) bool {
	raw, ok := value["widgets"].([]any)
	if !ok || len(raw) == 0 {
		writeValidation(w, "widgets", "required")
		return false
	}
	allowed := map[string]bool{"status": true, "sysinfo": true, "stats": true, "schedule": true, "alerts": true}
	seen := map[string]bool{}
	for _, item := range raw {
		widget := fmt.Sprint(item)
		if !allowed[widget] || seen[widget] {
			writeValidation(w, "widgets", "invalid value")
			return false
		}
		seen[widget] = true
	}
	return true
}

func boolFromAny(value any) bool {
	v, _ := value.(bool)
	return v
}

func mapsEqual(a, b map[string]any) bool {
	aj, errA := json.Marshal(a)
	bj, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aj) == string(bj)
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

func atomicWriteHistory(path string, records []apiHistoryRecord, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	var b strings.Builder
	for _, item := range records {
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		b.Write(data)
		b.WriteByte('\n')
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(b.String()), mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
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
	state, err := readMaintenanceState(path)
	return err == nil && state.Enabled
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

func normalizeTags(w http.ResponseWriter, values []string) ([]string, bool) {
	if len(values) > 20 {
		writeValidation(w, "tags", "must contain 20 items or less")
		return nil, false
	}
	seen := map[string]bool{}
	tags := []string{}
	for _, value := range values {
		tag := strings.TrimSpace(value)
		if tag == "" || len(tag) > 40 {
			writeValidation(w, "tags", "invalid tag")
			return nil, false
		}
		if seen[tag] {
			continue
		}
		seen[tag] = true
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags, true
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func intFromAny(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		i, _ := v.Int64()
		return int(i)
	default:
		return 0
	}
}

func validHTTPMethod(value string) bool {
	switch value {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
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
