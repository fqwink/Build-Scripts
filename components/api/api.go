package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/fqwink/build-scripts/components/security"
	"github.com/fqwink/build-scripts/components/statefile"
)

const (
	apiBinaryName       = "adlaire-ci-api"
	defaultAPIAddr      = "127.0.0.1:8765"
	defaultQueueMaxSize = 3
	apiTimeLayout       = "2006-01-02T15:04:05Z"
	passwordIterations  = security.CredentialIterations
	maxLoginCount       = int64(9223372036854775807)
)

var apiBinaryVersion = "V.0.0-dev"

func SetBinaryVersion(version string) {
	apiBinaryVersion = version
}

var errAPIQueueFull = errors.New("api queue full")
var errAPIQueueConflict = errors.New("api queue conflict")
var errAPIOutputTargetUnavailable = errors.New("api output target unavailable")
var errAPIStateLockConflict = errors.New("api state file lock conflict")
var errAPIPasswordMismatch = errors.New("api password mismatch")
var errAPIPasswordUnchanged = errors.New("api password unchanged")
var errAPITOTPInvalid = errors.New("api totp invalid")
var errAPITokenCorrupt = errors.New("api token state corrupted")
var errAPITokenLimit = errors.New("api token limit exceeded")
var errAPITokenCollision = errors.New("api token hash collision")

type APIConfig struct {
	Addr                   string
	StateDir               string
	Now                    func() time.Time
	RequestID              func() (string, error)
	HTTPClient             *http.Client
	GitHubAPIBaseURL       string
	CommandRunner          apiCommandRunner
	SystemdDir             string
	AuthTransactionTimeout time.Duration
}

type apiCommandRunner func(context.Context, string, ...string) (apiCommandResult, error)

type apiCommandResult struct {
	Stdout   string
	ExitCode int
}

type APIServer struct {
	cfg              APIConfig
	startedAt        time.Time
	mu               sync.Mutex
	sessions         map[string]apiSession
	loginTickets     map[string]apiLoginTicket
	totpSetup        *apiPendingTOTP
	loginFailures    map[string]apiLoginFailure
	authTransactions chan struct{}
}

type apiSession struct {
	TokenHash              string
	CreatedAt              string
	ExpiresAt              string
	LastUsedAt             string
	PasswordChangeRequired bool
}

type apiLoginTicket struct {
	CredentialsFingerprint string
	ExpiresAt              time.Time
}

type apiPendingTOTP struct {
	SecretBase32 string
	ExpiresAt    time.Time
}

type apiLoginFailure struct {
	Count         int
	LastFailureAt time.Time
	LockedUntil   time.Time
}

type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.wroteHeader {
		return
	}
	r.status = status
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if !r.wroteHeader {
		r.status = http.StatusOK
		r.wroteHeader = true
	}
	return r.ResponseWriter.Write(data)
}

type flushStatusRecorder struct {
	*statusRecorder
	flusher http.Flusher
}

func (r *flushStatusRecorder) Flush() {
	r.flusher.Flush()
}

func statusResponseWriter(w http.ResponseWriter) (http.ResponseWriter, *statusRecorder) {
	rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
	if flusher, ok := w.(http.Flusher); ok {
		return &flushStatusRecorder{statusRecorder: rec, flusher: flusher}, rec
	}
	return rec, rec
}

type apiCredentials struct {
	PasswordHash string  `json:"password_hash"`
	Salt         string  `json:"salt"`
	Algorithm    string  `json:"algorithm"`
	Iterations   int     `json:"iterations"`
	MustChange   bool    `json:"must_change"`
	LoginCount   int64   `json:"login_count"`
	LastLoginAt  *string `json:"last_login_at"`
	UpdatedAt    string  `json:"updated_at"`
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
	ActiveQueueEntry        map[string]any   `json:"active_queue_entry"`
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

type apiOutputTarget struct {
	Branch     string
	TargetFile string
	Out        string
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
	HistoryRetention        json.RawMessage  `json:"history_retention"`
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
	DurationAnomaly         json.RawMessage  `json:"duration_anomaly"`
	SessionTimeoutSeconds   int              `json:"session_timeout_seconds"`
	ForceBuildIntervalHours int              `json:"force_build_interval_hours"`
	BuildCooldownSeconds    int              `json:"build_cooldown_seconds"`
	ScheduleIntervalSeconds int              `json:"schedule_interval_seconds"`
	SchedulePaused          bool             `json:"schedule_paused"`
	AllowedHours            *apiAllowedHours `json:"allowed_hours"`
	WatchMode               string           `json:"watch_mode"`
	TagFilter               json.RawMessage  `json:"tag_filter"`
	BuildCacheEnabled       bool             `json:"build_cache_enabled"`
	DeployParallelism       int              `json:"deploy_parallelism"`
	RemoteBuild             json.RawMessage  `json:"remote_build"`
	ApprovalTimeoutSeconds  int              `json:"approval_timeout_seconds"`
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

type apiAccessLogRecord struct {
	At         string            `json:"at"`
	RequestID  string            `json:"request_id"`
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Query      map[string]string `json:"query"`
	Status     int               `json:"status"`
	DurationMS int64             `json:"duration_ms"`
	AuthType   string            `json:"auth_type"`
	Actor      *string           `json:"actor"`
	RemoteAddr *string           `json:"remote_addr"`
	UserAgent  *string           `json:"user_agent"`
	Error      *string           `json:"error"`
}

type apiAuthContextKey struct{}

type apiAuthCaptureContextKey struct{}

type apiRequestIDContextKey struct{}

type apiAuthInfo struct {
	AuthType               string
	Actor                  *string
	Scopes                 []string
	PasswordChangeRequired bool
}

type apiRateLimitPolicy struct {
	Enabled bool                         `json:"enabled"`
	Groups  map[string]apiRateLimitGroup `json:"groups"`
}

type apiRateLimitGroup struct {
	WindowSeconds int `json:"window_seconds"`
	MaxRequests   int `json:"max_requests"`
}

type apiRateLimitState struct {
	Windows map[string]apiRateLimitWindow `json:"windows"`
}

type apiRateLimitWindow struct {
	WindowStart string `json:"window_start"`
	Count       int    `json:"count"`
}

var apiExactRouteMethods = map[string][]string{
	"/api/access-control":          {http.MethodGet, http.MethodPost},
	"/api/access-log":              {http.MethodGet},
	"/api/alert-rules":             {http.MethodGet, http.MethodPost},
	"/api/api-access-log":          {http.MethodGet},
	"/api/api-rate-limit":          {http.MethodGet, http.MethodPost},
	"/api/approvals":               {http.MethodGet},
	"/api/audit-log":               {http.MethodGet},
	"/api/auth/totp":               {http.MethodDelete},
	"/api/auth/totp-confirm":       {http.MethodPost},
	"/api/auth/totp-setup":         {http.MethodPost},
	"/api/auth/totp-status":        {http.MethodGet},
	"/api/backup":                  {http.MethodGet},
	"/api/branch-config":           {http.MethodGet, http.MethodPost},
	"/api/build":                   {http.MethodPost},
	"/api/build/cancel":            {http.MethodPost},
	"/api/build/force":             {http.MethodPost},
	"/api/build/stream":            {http.MethodGet},
	"/api/build-chain-config":      {http.MethodGet, http.MethodPost},
	"/api/change-password":         {http.MethodPost},
	"/api/circuit-breaker/reset":   {http.MethodPost},
	"/api/config":                  {http.MethodGet, http.MethodPost},
	"/api/config-log":              {http.MethodGet},
	"/api/config/validate":         {http.MethodPost},
	"/api/dashboard":               {http.MethodGet},
	"/api/dashboard-layout":        {http.MethodGet, http.MethodPost},
	"/api/diagnostics":             {http.MethodGet},
	"/api/disk-usage":              {http.MethodGet},
	"/api/health":                  {http.MethodGet},
	"/api/history":                 {http.MethodGet},
	"/api/history/export":          {http.MethodGet},
	"/api/hooks":                   {http.MethodGet, http.MethodPost},
	"/api/log-level":               {http.MethodPost},
	"/api/login":                   {http.MethodPost},
	"/api/login/totp":              {http.MethodPost},
	"/api/logout":                  {http.MethodPost},
	"/api/logs":                    {http.MethodGet},
	"/api/logs/archive":            {http.MethodPost},
	"/api/logs/cleanup":            {http.MethodPost},
	"/api/logs/export":             {http.MethodGet},
	"/api/logs/search":             {http.MethodGet},
	"/api/maintenance":             {http.MethodGet},
	"/api/maintenance/disable":     {http.MethodPost},
	"/api/maintenance/enable":      {http.MethodPost},
	"/api/notes":                   {http.MethodGet, http.MethodPost},
	"/api/notify-config":           {http.MethodGet, http.MethodPost},
	"/api/notify-log":              {http.MethodGet},
	"/api/notify-test":             {http.MethodPost},
	"/api/notify/weekly-summary":   {http.MethodPost},
	"/api/output-meta":             {http.MethodGet},
	"/api/pat-status":              {http.MethodGet},
	"/api/pat-update":              {http.MethodPost},
	"/api/pat-verify":              {http.MethodPost},
	"/api/pipeline-config":         {http.MethodGet, http.MethodPost},
	"/api/queue":                   {http.MethodGet, http.MethodDelete},
	"/api/rate-limit":              {http.MethodGet},
	"/api/repo-config":             {http.MethodPost},
	"/api/repo-info":               {http.MethodGet},
	"/api/restore":                 {http.MethodPost},
	"/api/schedule":                {http.MethodGet},
	"/api/schedule/allowed-hours":  {http.MethodPost},
	"/api/schedule/cooldown":       {http.MethodPost},
	"/api/schedule/force-interval": {http.MethodPost},
	"/api/schedule/interval":       {http.MethodPost},
	"/api/schedule/pause":          {http.MethodPost},
	"/api/schedule/resume":         {http.MethodPost},
	"/api/sessions":                {http.MethodGet},
	"/api/sessions/revoke-all":     {http.MethodPost},
	"/api/smtp-config":             {http.MethodGet, http.MethodPost},
	"/api/smtp-test":               {http.MethodPost},
	"/api/snapshots":               {http.MethodGet},
	"/api/stats":                   {http.MethodGet},
	"/api/stats/build-duration":    {http.MethodGet},
	"/api/stats/build-trends":      {http.MethodGet},
	"/api/stats/timeline":          {http.MethodGet},
	"/api/status":                  {http.MethodGet},
	"/api/sysinfo":                 {http.MethodGet},
	"/api/tag-rules":               {http.MethodGet, http.MethodPost},
	"/api/tokens":                  {http.MethodGet, http.MethodPost},
	"/api/verify-output":           {http.MethodPost},
	"/api/webhook":                 {http.MethodPost},
	"/api/webhook-config":          {http.MethodGet, http.MethodPost},
	"/api/webhook-events":          {http.MethodGet},
}

var apiExactRouteQueryKeys = map[string]map[string]bool{
	"/api/access-log":           {"limit": true, "offset": true},
	"/api/api-access-log":       {"limit": true, "offset": true, "method": true, "path": true, "status": true},
	"/api/audit-log":            {"limit": true, "offset": true, "actor": true, "action": true, "result": true},
	"/api/config-log":           {"limit": true, "offset": true},
	"/api/history":              {"page": true, "per_page": true, "trigger": true, "tag": true, "flagged": true, "failure_category": true},
	"/api/logs":                 {"n": true, "q": true},
	"/api/logs/search":          {"q": true, "from": true, "to": true, "level": true},
	"/api/notify-log":           {"limit": true, "offset": true},
	"/api/stats":                {"days": true},
	"/api/stats/build-duration": {"n": true},
	"/api/stats/build-trends":   {"n": true},
	"/api/stats/timeline":       {"days": true},
	"/api/webhook-events":       {"limit": true, "offset": true},
}

type apiTokenRecord struct {
	ID         string   `json:"id"`
	Label      string   `json:"label"`
	Scopes     []string `json:"scopes"`
	TokenHash  string   `json:"token_hash"`
	CreatedAt  string   `json:"created_at"`
	LastUsedAt *string  `json:"last_used_at"`
	ExpiresAt  *string  `json:"expires_at"`
	RevokedAt  *string  `json:"revoked_at"`
}

type apiTokenFile struct {
	Tokens []apiTokenRecord `json:"tokens"`
}

type apiConfigLogRecord struct {
	At        string           `json:"at"`
	Type      string           `json:"type"`
	Action    string           `json:"action,omitempty"`
	Actor     string           `json:"actor,omitempty"`
	RequestID string           `json:"request_id,omitempty"`
	Endpoint  string           `json:"endpoint,omitempty"`
	Result    string           `json:"result,omitempty"`
	Error     *string          `json:"error,omitempty"`
	Diff      map[string][]any `json:"diff,omitempty"`
	DiffText  string           `json:"diff_text,omitempty"`
	Changes   map[string]any   `json:"changes,omitempty"`
}

func (record apiConfigLogRecord) MarshalJSON() ([]byte, error) {
	type wireRecord struct {
		At        string           `json:"at"`
		Type      string           `json:"type"`
		Action    string           `json:"action"`
		Actor     string           `json:"actor"`
		RequestID string           `json:"request_id"`
		Endpoint  string           `json:"endpoint"`
		Result    string           `json:"result"`
		Error     *string          `json:"error"`
		Diff      map[string][]any `json:"diff"`
		DiffText  string           `json:"diff_text"`
		Changes   map[string]any   `json:"changes,omitempty"`
	}
	action := firstNonEmpty(record.Action, apiConfigLogAction(record.Type))
	actor := firstNonEmpty(record.Actor, "admin")
	requestID := firstNonEmpty(record.RequestID, "00000000000000000000000000000000")
	endpoint := firstNonEmpty(record.Endpoint, apiConfigLogEndpoint(record.Type))
	result := firstNonEmpty(record.Result, "success")
	diff := record.Diff
	if diff == nil {
		diff = apiConfigLogDiff(record.Changes)
	}
	diffText := record.DiffText
	if diffText == "" {
		diffText = apiConfigLogDiffText(diff)
	}
	return json.Marshal(wireRecord{
		At:        record.At,
		Type:      record.Type,
		Action:    action,
		Actor:     actor,
		RequestID: requestID,
		Endpoint:  endpoint,
		Result:    result,
		Error:     record.Error,
		Diff:      diff,
		DiffText:  diffText,
		Changes:   record.Changes,
	})
}

func apiConfigLogAction(logType string) string {
	if strings.HasSuffix(logType, "_delete") || logType == "snapshot_delete" {
		return "delete"
	}
	return "update"
}

func apiConfigLogEndpoint(logType string) string {
	switch logType {
	case "server_config":
		return "POST /api/config"
	case "log_level":
		return "POST /api/log-level"
	case "api_rate_limit":
		return "POST /api/api-rate-limit"
	case "access_control":
		return "POST /api/access-control"
	case "repo_config":
		return "POST /api/repo-config"
	case "branch_config":
		return "POST /api/branch-config"
	case "notify_config":
		return "POST /api/notify-config"
	case "schedule_interval":
		return "POST /api/schedule/interval"
	case "schedule_pause":
		return "POST /api/schedule/pause"
	case "schedule_resume":
		return "POST /api/schedule/resume"
	case "schedule_allowed_hours":
		return "POST /api/schedule/allowed-hours"
	case "schedule_force_interval":
		return "POST /api/schedule/force-interval"
	case "schedule_cooldown":
		return "POST /api/schedule/cooldown"
	case "history_comment":
		return "POST /api/history/{id}/comment"
	case "history_flag":
		return "POST /api/history/{id}/flag"
	case "history_tags":
		return "POST /api/history/{id}/tags"
	case "rollback":
		return "POST /api/history/{id}/rollback"
	case "circuit_breaker_reset":
		return "POST /api/circuit-breaker/reset"
	case "maintenance":
		return "POST /api/maintenance"
	case "pat_update":
		return "POST /api/pat-update"
	case "restore":
		return "POST /api/restore"
	case "webhook_config":
		return "POST /api/webhook-config"
	case "snapshot_delete":
		return "DELETE /api/snapshots/{id}"
	case "pipeline_config":
		return "POST /api/pipeline-config"
	case "build_chain_config":
		return "POST /api/build-chain-config"
	case "notes":
		return "POST /api/notes"
	case "smtp_config":
		return "POST /api/smtp-config"
	case "dashboard_layout":
		return "POST /api/dashboard-layout"
	case "hook":
		return "POST /api/hooks"
	case "hook_delete":
		return "DELETE /api/hooks/{id}"
	case "alert_rule":
		return "POST /api/alert-rules"
	case "alert_rule_delete":
		return "DELETE /api/alert-rules/{id}"
	case "tag_rule":
		return "POST /api/tag-rules"
	case "tag_rule_delete":
		return "DELETE /api/tag-rules/{id}"
	default:
		return "POST /api/" + strings.ReplaceAll(strings.Trim(logType, "_"), "_", "-")
	}
}

func apiConfigLogDiff(changes map[string]any) map[string][]any {
	diff := map[string][]any{}
	for _, key := range sortedMapKeys(changes) {
		diff[key] = []any{nil, apiConfigLogMaskValue(key, changes[key])}
	}
	if len(diff) == 0 {
		diff["change"] = []any{nil, true}
	}
	return diff
}

func apiConfigLogDiffText(diff map[string][]any) string {
	keys := sortedMapKeysAnySlice(diff)
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		pair := diff[key]
		if len(pair) != 2 {
			pair = []any{nil, pair}
		}
		before, _ := json.Marshal(pair[0])
		after, _ := json.Marshal(pair[1])
		lines = append(lines, key+": "+string(before)+" -> "+string(after))
	}
	if len(lines) == 0 {
		return "change: null -> true"
	}
	return strings.Join(lines, "\n")
}

func sortedMapKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedMapKeysAnySlice(values map[string][]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func apiConfigLogMaskValue(key string, value any) any {
	lowerKey := strings.ToLower(key)
	if strings.Contains(lowerKey, "password") || strings.Contains(lowerKey, "token") || strings.Contains(lowerKey, "secret") || strings.Contains(lowerKey, "pat") {
		return "***"
	}
	switch typed := value.(type) {
	case map[string]any:
		return maskSecrets(typed)
	default:
		return value
	}
}

type apiConfigValidationIssue struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
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

type apiBuildLog struct {
	ID              string          `json:"id"`
	Branch          string          `json:"branch,omitempty"`
	TargetFile      string          `json:"target_file,omitempty"`
	Status          string          `json:"status"`
	StartedAt       string          `json:"started_at"`
	FinishedAt      string          `json:"finished_at"`
	TargetStatus    string          `json:"target_status"`
	Commit          map[string]any  `json:"commit"`
	Pipeline        apiPipelineLog  `json:"pipeline"`
	Warnings        []string        `json:"warnings"`
	DurationSeconds int64           `json:"duration_seconds"`
	Report          *apiBuildReport `json:"report,omitempty"`
	Error           *string         `json:"error,omitempty"`
	Comment         *string         `json:"comment"`
	Flagged         bool            `json:"flagged"`
	Tags            []string        `json:"tags"`
	OutputSHA256    string          `json:"output_sha256,omitempty"`
	OutputSizeBytes *int64          `json:"output_size_bytes,omitempty"`
	SHA256          string          `json:"sha256,omitempty"`
	BuildMeta       *apiBuildMeta   `json:"build_meta,omitempty"`
}

type apiBuildReport struct {
	Headings        int  `json:"headings"`
	TablesCount     int  `json:"tables_count"`
	CodeBlocksCount int  `json:"code_blocks_count"`
	SizeWarn        bool `json:"size_warn"`
}

type apiBuildMeta struct {
	BuildID   string `json:"build_id"`
	CommitSHA string `json:"commit_sha"`
	BuildAt   string `json:"build_at"`
}

type apiBuildStreamLogFrame struct {
	Type string `json:"type"`
	Line string `json:"line"`
	At   string `json:"at"`
}

type apiBuildStreamEndFrame struct {
	Type            string `json:"type"`
	Status          string `json:"status"`
	DurationSeconds *int64 `json:"duration_seconds"`
}

type apiBuildTrendFile struct {
	SchemaVersion int                   `json:"schema_version"`
	Samples       []apiBuildTrendSample `json:"samples"`
	Summary       map[string]any        `json:"summary"`
}

type apiBuildTrendSample struct {
	BuildID         string  `json:"build_id"`
	FinishedAt      string  `json:"finished_at"`
	Branch          string  `json:"branch"`
	Trigger         string  `json:"trigger"`
	DurationSeconds float64 `json:"duration_seconds"`
	Status          string  `json:"status"`
	TargetStatus    string  `json:"target_status"`
	Anomaly         bool    `json:"anomaly"`
}

type apiBuildChainConfig struct {
	Chains []apiBuildChainJob `json:"chains"`
}

type apiBuildChainJob struct {
	ID         string   `json:"id"`
	Branch     string   `json:"branch"`
	TargetFile string   `json:"target_file"`
	DependsOn  []string `json:"depends_on"`
	Required   bool     `json:"required"`
	Enabled    bool     `json:"enabled"`
}

type apiApprovalRecord struct {
	ID               string  `json:"id"`
	Status           string  `json:"status"`
	Branch           string  `json:"branch"`
	SHA              string  `json:"sha"`
	Target           string  `json:"target"`
	RequestedTrigger string  `json:"requested_trigger"`
	RequestedForce   bool    `json:"requested_force"`
	RequestedBy      string  `json:"requested_by"`
	DeliveryID       *string `json:"delivery_id"`
	CreatedAt        string  `json:"created_at"`
	ExpiresAt        string  `json:"expires_at"`
	DecidedAt        *string `json:"decided_at"`
	DecidedBy        *string `json:"decided_by"`
	QueueID          *string `json:"queue_id"`
	Reason           *string `json:"reason"`
}

type apiPipelineLog struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type apiCLIConfig struct {
	Addr            string
	StateDir        string
	InitCredentials bool
	AddrProvided    bool
}

func RunAPI(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	cfg, handled, err := parseAPIArgs(args, stdout)
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
	if cfg.InitCredentials {
		return executeAPIInitCredentials(cfg, stdin, stdout, stderr)
	}
	return executeAPIListener(cfg, stderr)
}

func parseAPIArgs(args []string, stdout io.Writer) (apiCLIConfig, bool, error) {
	cfg := apiCLIConfig{Addr: defaultAPIAddr}
	for _, arg := range args {
		if arg == "--help" {
			fmt.Fprintln(stdout, "Usage: adlaire-ci-api --state-dir path [--addr 127.0.0.1:port] [--init-credentials] [--version] [--help]")
			return cfg, true, nil
		}
	}
	for _, arg := range args {
		if arg == "--version" {
			fmt.Fprintf(stdout, "%s %s go=%s\n", apiBinaryName, apiBinaryVersion, runtime.Version())
			return cfg, true, nil
		}
	}
	for _, arg := range args {
		if !validCLIArgToken(arg) {
			return cfg, false, exitError{Code: 2, Msg: "invalid command line token"}
		}
	}
	stateDirProvided := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--init-credentials" {
			cfg.InitCredentials = true
			continue
		}
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
		case "--state-dir", "--addr":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
				return cfg, false, exitError{Code: 2, Msg: "missing value: " + arg}
			}
			i++
			if arg == "--state-dir" {
				cfg.StateDir = args[i]
				stateDirProvided = true
			} else {
				cfg.Addr = args[i]
				cfg.AddrProvided = true
			}
		default:
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
	}
	if cfg.InitCredentials && cfg.AddrProvided {
		return cfg, false, exitError{Code: 2, Msg: "--addr is not allowed with --init-credentials"}
	}
	if !stateDirProvided {
		return cfg, false, exitError{Code: 2, Msg: "state directory is required"}
	}
	if err := validateAPIStateDir(cfg.StateDir); err != nil {
		return cfg, false, err
	}
	if !cfg.InitCredentials {
		if err := validateAPIListenAddress(cfg.Addr); err != nil {
			return cfg, false, err
		}
	}
	return cfg, false, nil
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

func validateAPIListenAddress(addr string) error {
	host, portText, err := net.SplitHostPort(addr)
	if err != nil || host != "127.0.0.1" || portText == "" || (len(portText) > 1 && portText[0] == '0') {
		return exitError{Code: 2, Msg: "invalid listen address: " + addr}
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return exitError{Code: 2, Msg: "invalid listen address: " + addr}
	}
	return nil
}

func executeAPIInitCredentials(cfg apiCLIConfig, stdin io.Reader, stdout, stderr io.Writer) int {
	credentialsPath := filepath.Join(cfg.StateDir, ".admin_credentials")
	if _, err := os.Stat(credentialsPath); err == nil {
		fmt.Fprintln(stderr, "credentials already exist")
		return 2
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(stderr, "credentials write failed")
		return 1
	}
	if apiInputIsTerminal(stdin) {
		fmt.Fprintln(stderr, "initial password stdin must not be terminal")
		return 2
	}
	password, ok, err := readInitialPassword(stdin)
	if err != nil {
		fmt.Fprintln(stderr, "password input failed")
		return 1
	}
	if !ok {
		fmt.Fprintln(stderr, "invalid initial password")
		return 2
	}
	if err := InitCredentials(cfg.StateDir, password, time.Now().UTC()); err != nil {
		if err.Error() == "credentials already exist" {
			fmt.Fprintln(stderr, "credentials already exist")
			return 2
		}
		fmt.Fprintln(stderr, "credentials write failed")
		return 1
	}
	fmt.Fprintln(stdout, "credentials initialized")
	return 0
}

func apiInputIsTerminal(stdin io.Reader) bool {
	statReader, ok := stdin.(interface {
		Stat() (os.FileInfo, error)
	})
	if !ok {
		return false
	}
	info, err := statReader.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func readInitialPassword(stdin io.Reader) (string, bool, error) {
	data, err := io.ReadAll(io.LimitReader(stdin, 514))
	if err != nil {
		return "", false, err
	}
	defer zeroBytes(data)
	if len(data) == 0 || len(data) >= 514 || data[len(data)-1] != '\n' || bytes.Count(data, []byte{'\n'}) != 1 {
		return "", false, nil
	}
	passwordBytes := data[:len(data)-1]
	if !validInitialPasswordBytes(passwordBytes) {
		return "", false, nil
	}
	return string(passwordBytes), true, nil
}

func validInitialPasswordBytes(password []byte) bool {
	if bytes.ContainsAny(password, "\n\r\x00") || !utf8.Valid(password) {
		return false
	}
	runeCount := utf8.RuneCount(password)
	return runeCount >= 8 && runeCount <= 128
}

func validInitialPasswordString(password string) bool {
	if strings.ContainsAny(password, "\n\r\x00") || !utf8.ValidString(password) {
		return false
	}
	runeCount := utf8.RuneCountInString(password)
	return runeCount >= 8 && runeCount <= 128
}

func zeroBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

func executeAPIListener(cfg apiCLIConfig, stderr io.Writer) int {
	credentialsPath := filepath.Join(cfg.StateDir, ".admin_credentials")
	if err := validateCredentials(credentialsPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(stderr, "credentials are not initialized")
			return 2
		}
		fmt.Fprintln(stderr, "credentials are invalid")
		return 2
	}
	server, err := NewAPIServer(APIConfig{Addr: cfg.Addr, StateDir: cfg.StateDir})
	if err != nil {
		fmt.Fprintln(stderr, "credentials are invalid")
		return 2
	}
	httpServer := server.HTTPServer()
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(signals)
	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return 0
		}
		fmt.Fprintln(stderr, "listen failed")
		return 1
	case <-signals:
		return shutdownAPIHTTPServer(httpServer, errCh, signals, stderr)
	}
}

func shutdownAPIHTTPServer(server *http.Server, errCh <-chan error, signals <-chan os.Signal, stderr io.Writer) int {
	done := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			done <- err
			return
		}
		done <- <-errCh
	}()
	select {
	case <-signals:
		_ = server.Close()
		fmt.Fprintln(stderr, "shutdown failed")
		return 1
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return 0
		}
		_ = server.Close()
		fmt.Fprintln(stderr, "shutdown failed")
		return 1
	}
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
	if cfg.RequestID == nil {
		cfg.RequestID = newAPIRequestID
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}
	if cfg.GitHubAPIBaseURL == "" {
		cfg.GitHubAPIBaseURL = "https://api.github.com"
	}
	if cfg.CommandRunner == nil {
		cfg.CommandRunner = runAPICommand
	}
	if cfg.SystemdDir == "" {
		cfg.SystemdDir = "/etc/systemd/system"
	}
	if cfg.AuthTransactionTimeout <= 0 {
		cfg.AuthTransactionTimeout = 60 * time.Second
	}
	if err := validateCredentials(filepath.Join(cfg.StateDir, ".admin_credentials")); err != nil {
		return nil, err
	}
	return &APIServer{cfg: cfg, startedAt: cfg.Now().UTC(), sessions: map[string]apiSession{}, loginTickets: map[string]apiLoginTicket{}, loginFailures: map[string]apiLoginFailure{}, authTransactions: make(chan struct{}, 1)}, nil
}

func apiRouteMethods(path string) ([]string, bool) {
	if methods, ok := apiExactRouteMethods[path]; ok {
		return methods, true
	}
	switch {
	case strings.HasPrefix(path, "/api/history/"):
		return apiHistoryRouteMethods(path)
	case strings.HasPrefix(path, "/api/snapshots/"):
		return apiSnapshotRouteMethods(path)
	case strings.HasPrefix(path, "/api/tokens/"):
		return []string{http.MethodDelete}, true
	case strings.HasPrefix(path, "/api/alert-rules/"), strings.HasPrefix(path, "/api/tag-rules/"):
		return []string{http.MethodDelete}, true
	case strings.HasPrefix(path, "/api/hooks/"):
		return apiHookRouteMethods(path)
	case strings.HasPrefix(path, "/api/approvals/"):
		return apiApprovalRouteMethods(path)
	default:
		return nil, false
	}
}

func apiRouteQueryKeys(path string) map[string]bool {
	if allowed, ok := apiExactRouteQueryKeys[path]; ok {
		return allowed
	}
	return nil
}

func apiHistoryRouteMethods(path string) ([]string, bool) {
	rest := strings.TrimPrefix(path, "/api/history/")
	id, suffix, ok := strings.Cut(rest, "/")
	if !ok || id == "" {
		return nil, false
	}
	switch suffix {
	case "log":
		return []string{http.MethodGet}, true
	case "comment":
		return []string{http.MethodGet, http.MethodPost}, true
	case "flag", "tags", "rollback":
		return []string{http.MethodPost}, true
	default:
		return nil, false
	}
}

func apiSnapshotRouteMethods(path string) ([]string, bool) {
	rest := strings.TrimPrefix(path, "/api/snapshots/")
	id, suffix, _ := strings.Cut(rest, "/")
	if id == "" {
		return nil, false
	}
	switch suffix {
	case "":
		return []string{http.MethodDelete}, true
	case "download":
		return []string{http.MethodGet}, true
	default:
		return nil, false
	}
}

func apiHookRouteMethods(path string) ([]string, bool) {
	rest := strings.Trim(strings.TrimPrefix(path, "/api/hooks/"), "/")
	id, suffix, hasSuffix := strings.Cut(rest, "/")
	if id == "" {
		return nil, false
	}
	if hasSuffix {
		if suffix == "log" {
			return []string{http.MethodGet}, true
		}
		return nil, false
	}
	return []string{http.MethodDelete}, true
}

func apiApprovalRouteMethods(path string) ([]string, bool) {
	rest := strings.Trim(strings.TrimPrefix(path, "/api/approvals/"), "/")
	id, action, ok := strings.Cut(rest, "/")
	if !ok || id == "" {
		return nil, false
	}
	switch action {
	case "approve", "reject":
		return []string{http.MethodPost}, true
	default:
		return nil, false
	}
}

func apiMethodAllowed(method string, methods []string) bool {
	for _, allowed := range methods {
		if method == allowed {
			return true
		}
	}
	return false
}

func InitCredentials(stateDir, password string, now time.Time) error {
	if !filepath.IsAbs(stateDir) {
		return fmt.Errorf("state directory must be absolute: %s", stateDir)
	}
	if !validInitialPasswordString(password) {
		return errors.New("invalid initial password")
	}
	path := filepath.Join(stateDir, ".admin_credentials")
	if _, err := os.Stat(path); err == nil {
		return errors.New("credentials already exist")
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	record, err := security.NewPasswordRecord(password)
	if err != nil {
		return err
	}
	at := now.UTC().Format(apiTimeLayout)
	cred := apiCredentials{
		PasswordHash: record.PasswordHash,
		Salt:         record.Salt,
		Algorithm:    record.Algorithm,
		Iterations:   record.Iterations,
		MustChange:   true,
		LoginCount:   0,
		LastLoginAt:  nil,
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
	mux.HandleFunc("/api/stats/build-trends", s.withAuth(s.handleStatsBuildTrends))
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
	mux.HandleFunc("/api/alert-rules", s.withAuth(s.handleRuleFile(".alert_rules", "alert_rule", http.StatusCreated)))
	mux.HandleFunc("/api/alert-rules/", s.withAuth(s.handleRulePath(".alert_rules", "alert_rule")))
	mux.HandleFunc("/api/tag-rules", s.withAuth(s.handleRuleFile(".tag_rules", "tag_rule", http.StatusCreated)))
	mux.HandleFunc("/api/tag-rules/", s.withAuth(s.handleRulePath(".tag_rules", "tag_rule")))
	mux.HandleFunc("/api/hooks", s.withAuth(s.handleRuleFile(".hooks", "hook", http.StatusCreated)))
	mux.HandleFunc("/api/hooks/", s.withAuth(s.handleHookPath))
	mux.HandleFunc("/api/verify-output", s.withAuth(s.handleVerifyOutput))
	mux.HandleFunc("/api/pipeline-config", s.withAuth(s.handlePipelineConfig))
	mux.HandleFunc("/api/build-chain-config", s.withAuth(s.handleBuildChainConfig))
	mux.HandleFunc("/api/notes", s.withAuth(s.handleNotes))
	mux.HandleFunc("/api/dashboard-layout", s.withAuth(s.handleDashboardLayout))
	mux.HandleFunc("/api/smtp-config", s.withAuth(s.handleSMTPConfig))
	mux.HandleFunc("/api/smtp-test", s.withAuth(s.handleSMTPTest))
	mux.HandleFunc("/api/circuit-breaker/reset", s.withAuth(s.handleCircuitReset))
	mux.HandleFunc("/api/maintenance", s.withAuth(s.handleMaintenance))
	mux.HandleFunc("/api/maintenance/enable", s.withAuth(s.handleMaintenanceEnable))
	mux.HandleFunc("/api/maintenance/disable", s.withAuth(s.handleMaintenanceDisable))
	mux.HandleFunc("/api/queue", s.withAuth(s.handleQueue))
	mux.HandleFunc("/api/approvals", s.withAuth(s.handleApprovals))
	mux.HandleFunc("/api/approvals/", s.withAuth(s.handleApprovalPath))
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Not found")
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			s.handleAdminStatic(w, r)
			return
		}
		start := time.Now()
		requestID, err := s.cfg.RequestID()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		w.Header().Set("X-Request-Id", requestID)
		authCapture := &apiAuthInfo{AuthType: "none"}
		r = r.WithContext(context.WithValue(r.Context(), apiAuthCaptureContextKey{}, authCapture))
		r = r.WithContext(context.WithValue(r.Context(), apiRequestIDContextKey{}, requestID))
		rw, rec := statusResponseWriter(w)
		if methods, ok := apiRouteMethods(r.URL.Path); !ok {
			writeError(rw, http.StatusNotFound, "Not found")
		} else if !apiMethodAllowed(r.Method, methods) {
			writeError(rw, http.StatusMethodNotAllowed, "Method not allowed")
		} else if allowed, err := s.accessAllowed(r); err != nil {
			writeError(rw, http.StatusServiceUnavailable, "Access control unavailable")
		} else if !allowed {
			writeError(rw, http.StatusForbidden, "Forbidden")
		} else {
			mux.ServeHTTP(rw, r)
		}
		_ = s.appendAPIAccessLog(r, requestID, rec.status, start)
	})
}

func (s *APIServer) handleAdminStatic(w http.ResponseWriter, r *http.Request) {
	var name, contentType, cacheControl string
	switch r.URL.Path {
	case "/", "/admin/", "/admin/index.html":
		name = "index.html"
		contentType = "text/html; charset=utf-8"
		cacheControl = "no-store"
	case "/admin/adlaire-ci-sdk.js":
		name = "adlaire-ci-sdk.js"
		contentType = "text/javascript; charset=utf-8"
		cacheControl = "no-cache"
	default:
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	adminDir := filepath.Join(s.cfg.StateDir, "admin")
	adminInfo, err := os.Lstat(adminDir)
	if err != nil || adminInfo.Mode()&os.ModeSymlink != 0 || !adminInfo.IsDir() {
		http.NotFound(w, r)
		return
	}
	path := filepath.Join(adminDir, name)
	info, err := os.Lstat(path)
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", cacheControl)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodGet {
		_, _ = w.Write(data)
	}
}

func (s *APIServer) HTTPServer() *http.Server {
	return &http.Server{
		Addr:              s.cfg.Addr,
		Handler:           s.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    32768,
		ErrorLog:          log.New(io.Discard, "", 0),
	}
}

func (s *APIServer) ListenAndServe() error {
	return s.HTTPServer().ListenAndServe()
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	if !s.enforceAPIRateLimit(w, r, apiAuthInfo{AuthType: "anonymous"}, "login") {
		return
	}
	if !validateRouteQuery(w, r) {
		return
	}
	releaseAuthTransaction, ok := s.acquireAuthTransaction(w, r)
	if !ok {
		return
	}
	defer releaseAuthTransaction()
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
	if s.loginLockActive(r) {
		if err := s.appendAuthAccessLog(r, "login_locked", "denied", "too_many_attempts"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := s.appendSecurityAudit(r, "permission_denied", "anonymous", nil, "auth", "login", "denied"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeError(w, http.StatusTooManyRequests, "Too many attempts")
		return
	}
	path := filepath.Join(s.cfg.StateDir, ".admin_credentials")
	cred, err := readCredentialsLocked(path)
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if !credentialsPasswordEqual(cred, req.Password) {
		s.recordLoginFailure(r)
		if err := s.appendAuthAccessLog(r, "login_failure", "failure", "password_mismatch"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := s.appendSecurityAudit(r, "login_failure", "anonymous", nil, "auth", "login", "failure"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	s.clearLoginFailure(r)
	totp, err := readTOTPSecretLocked(filepath.Join(s.cfg.StateDir, ".totp_secret"))
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if totp.Enabled {
		ticket, err := randomToken()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := s.appendAuthAccessLog(r, "totp_required", "success", ""); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := s.appendSecurityAudit(r, "totp_required", "anonymous", nil, "auth", "login", "success"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		s.mu.Lock()
		s.loginTickets[tokenHash(ticket)] = apiLoginTicket{CredentialsFingerprint: credentialsFingerprint(cred), ExpiresAt: s.cfg.Now().UTC().Add(5 * time.Minute)}
		s.mu.Unlock()
		writeJSON(w, http.StatusOK, map[string]any{"totp_required": true, "ticket": ticket})
		return
	}
	token, err := randomToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	now := s.nowString()
	cred, err = updateCredentialsLocked(path, func(current apiCredentials) (apiCredentials, bool, error) {
		if !credentialsPasswordEqual(current, req.Password) {
			return current, false, errAPIPasswordMismatch
		}
		advanceCredentialsLogin(&current, now)
		return current, true, nil
	})
	if errors.Is(err, errAPIPasswordMismatch) {
		s.recordLoginFailure(r)
		if logErr := s.appendAuthAccessLog(r, "login_failure", "failure", "password_mismatch"); logErr != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if logErr := s.appendSecurityAudit(r, "login_failure", "anonymous", nil, "auth", "login", "failure"); logErr != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	mustChange := credentialsMustChange(&cred)
	if err := s.appendAuthAccessLog(r, "login_success", "success", ""); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := s.appendSecurityAudit(r, "login_success", "anonymous", nil, "auth", "login", "success"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	timeoutSeconds := 28800
	if cfg, err := s.readMergedConfig(); err == nil {
		timeoutSeconds = cfg.SessionTimeoutSeconds
	}
	expires := s.cfg.Now().UTC().Add(time.Duration(timeoutSeconds) * time.Second).Format(apiTimeLayout)
	s.mu.Lock()
	s.sessions[tokenHash(token)] = apiSession{TokenHash: tokenHash(token), CreatedAt: now, ExpiresAt: expires, LastUsedAt: now, PasswordChangeRequired: mustChange == "forced"}
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "must_change": mustChange})
}

func (s *APIServer) handleLoginTOTP(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	if !s.enforceAPIRateLimit(w, r, apiAuthInfo{AuthType: "anonymous"}, "login") {
		return
	}
	if !validateRouteQuery(w, r) {
		return
	}
	releaseAuthTransaction, ok := s.acquireAuthTransaction(w, r)
	if !ok {
		return
	}
	defer releaseAuthTransaction()
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
		s.writeTOTPLoginFailure(w, r)
		return
	}
	totpPath := filepath.Join(s.cfg.StateDir, ".totp_secret")
	totp, err := readTOTPSecretLocked(totpPath)
	if err != nil || !totp.Enabled || totp.SecretBase32 == nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		s.writeTOTPLoginFailure(w, r)
		return
	}
	step, ok := verifyTOTPCodeStep(*totp.SecretBase32, req.Code, s.cfg.Now().UTC())
	if !ok || (totp.LastAcceptedStep != nil && step <= *totp.LastAcceptedStep) {
		s.writeTOTPLoginFailure(w, r)
		return
	}
	credPath := filepath.Join(s.cfg.StateDir, ".admin_credentials")
	cred, err := readCredentialsLocked(credPath)
	if err != nil || !passwordHashEqual(credentialsFingerprint(cred), ticket.CredentialsFingerprint) {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		s.writeTOTPLoginFailure(w, r)
		return
	}
	token, err := randomToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	totp, err = updateTOTPSecretLocked(totpPath, func(current apiTOTPSecret) (apiTOTPSecret, bool, error) {
		if !current.Enabled || current.SecretBase32 == nil || !passwordHashEqual(*current.SecretBase32, *totp.SecretBase32) {
			return current, false, errAPITOTPInvalid
		}
		if current.LastAcceptedStep != nil && step <= *current.LastAcceptedStep {
			return current, false, errAPITOTPInvalid
		}
		current.LastAcceptedStep = &step
		return current, true, nil
	})
	if errors.Is(err, errAPITOTPInvalid) {
		s.writeTOTPLoginFailure(w, r)
		return
	}
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	timeoutSeconds := 28800
	if cfg, err := s.readMergedConfig(); err == nil {
		timeoutSeconds = cfg.SessionTimeoutSeconds
	}
	now := s.nowString()
	cred, err = updateCredentialsLocked(credPath, func(current apiCredentials) (apiCredentials, bool, error) {
		if !passwordHashEqual(credentialsFingerprint(current), ticket.CredentialsFingerprint) {
			return current, false, errAPIPasswordMismatch
		}
		advanceCredentialsLogin(&current, now)
		return current, true, nil
	})
	if errors.Is(err, errAPIPasswordMismatch) {
		s.writeTOTPLoginFailure(w, r)
		return
	}
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	mustChange := credentialsMustChange(&cred)
	if err := s.appendAuthAccessLog(r, "login_success", "success", ""); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := s.appendSecurityAudit(r, "login_success", "anonymous", nil, "auth", "login", "success"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	expires := s.cfg.Now().UTC().Add(time.Duration(timeoutSeconds) * time.Second).Format(apiTimeLayout)
	s.mu.Lock()
	s.sessions[tokenHash(token)] = apiSession{TokenHash: tokenHash(token), CreatedAt: now, ExpiresAt: expires, LastUsedAt: now, PasswordChangeRequired: mustChange == "forced"}
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "must_change": mustChange})
}

func (s *APIServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	token := bearerToken(r)
	s.mu.Lock()
	delete(s.sessions, tokenHash(token))
	s.mu.Unlock()
	if err := s.appendAuthAccessLog(r, "logout", "success", ""); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	actorType, actorID := apiAuditActor(r)
	if err := s.appendSecurityAudit(r, "logout", actorType, actorID, "auth", "logout", "success"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
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
	if !validInitialPasswordString(req.NewPassword) {
		writeValidation(w, "new_password", "invalid")
		return
	}
	path := filepath.Join(s.cfg.StateDir, ".admin_credentials")
	now := s.nowString()
	_, err := updateCredentialsLocked(path, func(current apiCredentials) (apiCredentials, bool, error) {
		if !credentialsPasswordEqual(current, req.CurrentPassword) {
			return current, false, errAPIPasswordMismatch
		}
		if credentialsPasswordEqual(current, req.NewPassword) {
			return current, false, errAPIPasswordUnchanged
		}
		record, err := security.NewPasswordRecord(req.NewPassword)
		if err != nil {
			return current, false, err
		}
		current.PasswordHash = record.PasswordHash
		current.Salt = record.Salt
		current.Algorithm = record.Algorithm
		current.Iterations = record.Iterations
		current.MustChange = false
		current.UpdatedAt = now
		return current, true, nil
	})
	if errors.Is(err, errAPIPasswordMismatch) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if errors.Is(err, errAPIPasswordUnchanged) {
		writeValidation(w, "new_password", "must be different from current password")
		return
	}
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	currentSessionHash := tokenHash(bearerToken(r))
	s.mu.Lock()
	for key, session := range s.sessions {
		if key != currentSessionHash {
			delete(s.sessions, key)
			continue
		}
		session.PasswordChangeRequired = false
		s.sessions[key] = session
	}
	s.mu.Unlock()
	if err := s.appendAuthAccessLog(r, "password_change", "success", ""); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	actorType, actorID := apiAuditActor(r)
	if err := s.appendSecurityAudit(r, "password_change", actorType, actorID, "auth", "password", "success"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Password changed"})
}

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	if !validateRouteQuery(w, r) {
		return
	}
	lastBuildAt := any(nil)
	lastBuildStatus := "none"
	lastDeployAt := any(nil)
	lastDeployStatus := any(nil)
	checks := []string{}
	statusPath := filepath.Join(s.cfg.StateDir, ".build_status.json")
	if st, ok, check := readBuildStatusForHealth(statusPath); ok {
		lastBuildAt = st.LastFinishedAt
		lastBuildStatus = stringOr(st.Status, "none")
		lastDeployAt = st.LastDeployAt
		lastDeployStatus = st.LastDeployStatus
	} else {
		if check != "" {
			checks = append(checks, check)
		}
		history, err := readHealthHistory(filepath.Join(s.cfg.StateDir, ".build_history"))
		if err != nil {
			checks = append(checks, "build_history_read_error")
		} else if latest := latestStatusSummaryHistory(history); latest != nil {
			lastBuildAt = nullableString(firstNonEmpty(latest.FinishedAt, latest.StartedAt, latest.BuildAt))
			lastBuildStatus = normalizeHealthBuildStatus(latest.Status)
		}
	}
	pendingTransfers, err := readJSONArrayCount(filepath.Join(s.cfg.StateDir, ".pending_transfers"))
	if err != nil {
		pendingTransfers = 0
		checks = append(checks, "pending_transfers_read_error")
	} else if pendingTransfers > 0 {
		checks = append(checks, "pending_transfers_present")
	}
	if _, err := readJSONArrayCount(filepath.Join(s.cfg.StateDir, ".notify_pending")); err != nil {
		checks = append(checks, "notify_pending_read_error")
	}
	if healthRunnerStale(filepath.Join(s.cfg.StateDir, ".build_status.json"), s.cfg.Now().UTC()) {
		checks = append(checks, "runner_stale")
	}
	uptime := int64(s.cfg.Now().UTC().Sub(s.startedAt).Seconds())
	if uptime < 0 {
		uptime = 0
	}
	status := "ok"
	if len(checks) > 0 {
		status = "degraded"
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":             status,
		"last_build_at":      lastBuildAt,
		"last_build_status":  lastBuildStatus,
		"last_deploy_at":     lastDeployAt,
		"last_deploy_status": lastDeployStatus,
		"pending_transfers":  pendingTransfers,
		"uptime_seconds":     uptime,
		"checks":             checks,
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
	if err := s.appendAuthAccessLog(r, "session_revoke_all", "success", ""); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	actorType, actorID := apiAuditActor(r)
	if err := s.appendSecurityAudit(r, "session_revoke_all", actorType, actorID, "auth", "sessions", "success"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "All other sessions revoked", "revoked_count": revoked})
}

func (s *APIServer) handleTOTPStatus(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	totp, err := readTOTPSecretLocked(filepath.Join(s.cfg.StateDir, ".totp_secret"))
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	writeJSON(w, http.StatusOK, totpStatusPayload(totp))
}

func (s *APIServer) handleTOTPSetup(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	totp, err := readTOTPSecretLocked(filepath.Join(s.cfg.StateDir, ".totp_secret"))
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
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
	if err := s.appendAuthAccessLog(r, "totp_setup", "success", ""); err != nil {
		s.mu.Lock()
		if s.totpSetup != nil && s.totpSetup.SecretBase32 == secret {
			s.totpSetup = nil
		}
		s.mu.Unlock()
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	actorType, actorID := apiAuditActor(r)
	if err := s.appendSecurityAudit(r, "totp_setup", actorType, actorID, "auth", "totp_setup", "success"); err != nil {
		s.mu.Lock()
		if s.totpSetup != nil && s.totpSetup.SecretBase32 == secret {
			s.totpSetup = nil
		}
		s.mu.Unlock()
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
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
		if err := s.appendAuthAccessLog(r, "totp_failure", "failure", "invalid_totp"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		actorType, actorID := apiAuditActor(r)
		if err := s.appendSecurityAudit(r, "totp_failure", actorType, actorID, "auth", "totp_confirm", "failure"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	now := s.nowString()
	secret := pending.SecretBase32
	totp := apiTOTPSecret{Enabled: true, SecretBase32: &secret, ConfirmedAt: &now}
	totp, err := updateTOTPSecretLocked(filepath.Join(s.cfg.StateDir, ".totp_secret"), func(current apiTOTPSecret) (apiTOTPSecret, bool, error) {
		if current.Enabled {
			return current, false, errAPITOTPInvalid
		}
		return apiTOTPSecret{Enabled: true, SecretBase32: &secret, ConfirmedAt: &now}, true, nil
	})
	if errors.Is(err, errAPITOTPInvalid) {
		writeError(w, http.StatusConflict, "Conflict")
		return
	}
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	s.mu.Lock()
	s.totpSetup = nil
	s.mu.Unlock()
	if err := s.appendAuthAccessLog(r, "totp_enabled", "success", ""); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	actorType, actorID := apiAuditActor(r)
	if err := s.appendSecurityAudit(r, "totp_enabled", actorType, actorID, "auth", "totp", "success"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
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
	totp, err := readTOTPSecretLocked(path)
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if !totp.Enabled || totp.SecretBase32 == nil {
		writeError(w, http.StatusConflict, "Conflict")
		return
	}
	step, ok := verifyTOTPCodeStep(*totp.SecretBase32, body.Code, s.cfg.Now().UTC())
	if !ok || (totp.LastAcceptedStep != nil && step <= *totp.LastAcceptedStep) {
		if err := s.appendAuthAccessLog(r, "totp_failure", "failure", "invalid_totp"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		actorType, actorID := apiAuditActor(r)
		if err := s.appendSecurityAudit(r, "totp_failure", actorType, actorID, "auth", "totp_disable", "failure"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	next := apiTOTPSecret{Enabled: false}
	next, err = updateTOTPSecretLocked(path, func(current apiTOTPSecret) (apiTOTPSecret, bool, error) {
		if !current.Enabled || current.SecretBase32 == nil || !passwordHashEqual(*current.SecretBase32, *totp.SecretBase32) {
			return current, false, errAPITOTPInvalid
		}
		if current.LastAcceptedStep != nil && step <= *current.LastAcceptedStep {
			return current, false, errAPITOTPInvalid
		}
		return apiTOTPSecret{Enabled: false}, true, nil
	})
	if errors.Is(err, errAPITOTPInvalid) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	s.mu.Lock()
	s.loginTickets = map[string]apiLoginTicket{}
	s.totpSetup = nil
	s.mu.Unlock()
	if err := s.appendAuthAccessLog(r, "totp_disabled", "success", ""); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	actorType, actorID := apiAuditActor(r)
	if err := s.appendSecurityAudit(r, "totp_disabled", actorType, actorID, "auth", "totp", "success"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
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
	raw, order, ok := decodeRawObjectBody(w, r, true)
	if !ok {
		return
	}
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	next, errorsList, warningsList, valid := validateConfigDryRun(cfg, raw, order)
	for _, item := range errorsList {
		if item.Code == "unknown_key" {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
				"error":   "Validation failed",
				"details": []map[string]string{{"field": item.Field, "message": "Unknown config key"}},
			})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"valid": valid, "config": next, "warnings": warningsList, "errors": errorsList})
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
		policy, ok := apiRateLimitPolicyFromConfig(cfg.APIRateLimit)
		if !ok {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		state, err := s.readAPIRateLimitState()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusOK, s.apiRateLimitResponse(policy, state))
	case http.MethodPost:
		var body map[string]any
		if !decodeBody(w, r, &body, true) {
			return
		}
		policy, ok := validateAPIRateLimitPolicyBody(w, body)
		if !ok {
			return
		}
		cfg, err := s.readMergedConfig()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		currentPolicy, ok := apiRateLimitPolicyFromConfig(cfg.APIRateLimit)
		if !ok {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		state, err := s.readAPIRateLimitState()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if apiRateLimitPoliciesEqual(currentPolicy, policy) {
			writeJSON(w, http.StatusOK, s.apiRateLimitResponse(policy, state))
			return
		}
		policyMap := apiRateLimitPolicyToMap(policy)
		cfg.APIRateLimit = policyMap
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".server_config"), cfg, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		emptyState := apiRateLimitState{Windows: map[string]apiRateLimitWindow{}}
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".api_rate_state"), emptyState, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "api_rate_limit", Changes: map[string]any{"policy": policyMap}}); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		actorType, actorID := apiAuditActor(r)
		if err := s.appendSecurityAudit(r, "rate_limit_update", actorType, actorID, "config", "api_rate_limit", "success"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusOK, s.apiRateLimitResponse(policy, emptyState))
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) enforceAPIRateLimit(w http.ResponseWriter, r *http.Request, auth apiAuthInfo, group string) bool {
	if group == "" {
		return true
	}
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return false
	}
	policy, ok := apiRateLimitPolicyFromConfig(cfg.APIRateLimit)
	if !ok {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return false
	}
	if !policy.Enabled {
		return true
	}
	groupPolicy, ok := policy.Groups[group]
	if !ok {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return false
	}
	keys := apiRateLimitKeys(auth, r, group)
	if len(keys) == 0 {
		return true
	}
	statePath := filepath.Join(s.cfg.StateDir, ".api_rate_state")
	release, err := acquireStateFileLock(statePath)
	if err != nil {
		writeError(w, http.StatusConflict, "Conflict")
		return false
	}
	defer release()
	state, err := s.readAPIRateLimitStateLocked(statePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return false
	}
	now := s.cfg.Now().UTC()
	limited := false
	for _, key := range keys {
		window := apiRateLimitCurrentWindow(state.Windows[key], groupPolicy, now)
		if window.Count >= groupPolicy.MaxRequests {
			limited = true
		}
		state.Windows[key] = window
	}
	if limited {
		actorType, actorID := apiRateLimitAuditActor(auth)
		if err := s.appendSecurityAudit(r, "permission_denied", actorType, actorID, "endpoint", strings.ToUpper(r.Method)+" "+r.URL.Path, "denied"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return false
		}
		writeError(w, http.StatusTooManyRequests, "Too many requests")
		return false
	}
	for _, key := range keys {
		window := state.Windows[key]
		window.Count++
		state.Windows[key] = window
	}
	if err := atomicWriteJSONLocked(statePath, state, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return false
	}
	return true
}

func apiRateLimitAuditActor(auth apiAuthInfo) (string, any) {
	switch auth.AuthType {
	case "session":
		return "admin", "admin"
	case "api_token":
		if auth.Actor != nil {
			return "api_token", *auth.Actor
		}
		return "api_token", nil
	case "webhook":
		return "webhook", "webhook"
	default:
		return "anonymous", nil
	}
}

func apiRateLimitKeys(auth apiAuthInfo, r *http.Request, group string) []string {
	ipKey := "ip:" + apiRateLimitRemoteAddr(r) + ":" + group
	switch auth.AuthType {
	case "session":
		return []string{"session:admin:" + group, ipKey}
	case "api_token":
		if auth.Actor != nil && *auth.Actor != "" {
			return []string{"token:" + *auth.Actor + ":" + group, ipKey}
		}
		return []string{ipKey}
	case "webhook", "anonymous":
		return []string{ipKey}
	default:
		return []string{ipKey}
	}
}

func apiRateLimitRemoteAddr(r *http.Request) string {
	value := apiRemoteAddr(r)
	if value == nil || *value == "" {
		return "unknown"
	}
	return *value
}

func loginFailureKey(r *http.Request) string {
	return "ip:" + apiRateLimitRemoteAddr(r) + ":login"
}

func (s *APIServer) loginLockActive(r *http.Request) bool {
	key := loginFailureKey(r)
	now := s.cfg.Now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.loginFailures[key]
	if !ok {
		return false
	}
	if !entry.LockedUntil.IsZero() && now.Before(entry.LockedUntil) {
		return true
	}
	if !entry.LockedUntil.IsZero() && !now.Before(entry.LockedUntil) {
		delete(s.loginFailures, key)
	}
	return false
}

func (s *APIServer) recordLoginFailure(r *http.Request) {
	key := loginFailureKey(r)
	now := s.cfg.Now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := s.loginFailures[key]
	entry.Count++
	entry.LastFailureAt = now
	if entry.Count >= 10 {
		entry.LockedUntil = now.Add(10 * time.Minute)
	}
	s.loginFailures[key] = entry
}

func (s *APIServer) clearLoginFailure(r *http.Request) {
	key := loginFailureKey(r)
	s.mu.Lock()
	delete(s.loginFailures, key)
	s.mu.Unlock()
}

func apiRateLimitCurrentWindow(window apiRateLimitWindow, group apiRateLimitGroup, now time.Time) apiRateLimitWindow {
	start, err := time.Parse(apiTimeLayout, window.WindowStart)
	if err != nil || !now.Before(start.Add(time.Duration(group.WindowSeconds)*time.Second)) {
		return apiRateLimitWindow{WindowStart: now.Format(apiTimeLayout), Count: 0}
	}
	if window.Count < 0 {
		window.Count = 0
	}
	return window
}

func (s *APIServer) readAPIRateLimitState() (apiRateLimitState, error) {
	path := filepath.Join(s.cfg.StateDir, ".api_rate_state")
	release, err := acquireStateFileLock(path)
	if err != nil {
		return apiRateLimitState{}, err
	}
	defer release()
	return s.readAPIRateLimitStateLocked(path)
}

func (s *APIServer) readAPIRateLimitStateLocked(path string) (apiRateLimitState, error) {
	state := apiRateLimitState{Windows: map[string]apiRateLimitWindow{}}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return state, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&state); err != nil || state.Windows == nil || !apiRateLimitStateValid(state) {
		if backupErr := s.backupCorruptAPIRateLimitState(path); backupErr != nil {
			return apiRateLimitState{}, backupErr
		}
		return apiRateLimitState{Windows: map[string]apiRateLimitWindow{}}, nil
	}
	return state, nil
}

func (s *APIServer) backupCorruptAPIRateLimitState(path string) error {
	backup := fmt.Sprintf("%s.corrupt.%s.bak", path, s.cfg.Now().UTC().Format("20060102150405"))
	if err := os.Rename(path, backup); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return atomicWriteJSONLocked(path, apiRateLimitState{Windows: map[string]apiRateLimitWindow{}}, 0600)
}

func apiRateLimitStateValid(state apiRateLimitState) bool {
	for key, window := range state.Windows {
		if key == "" || window.Count < 0 || !validAPITime(window.WindowStart) {
			return false
		}
		group := apiRateLimitGroupFromKey(key)
		if !apiRateLimitGroupNameAllowed(group) {
			return false
		}
	}
	return true
}

func apiRateLimitGroupFromKey(key string) string {
	if idx := strings.LastIndex(key, ":"); idx >= 0 && idx < len(key)-1 {
		return key[idx+1:]
	}
	return ""
}

func (s *APIServer) apiRateLimitResponse(policy apiRateLimitPolicy, state apiRateLimitState) map[string]any {
	return map[string]any{
		"policy":        apiRateLimitPolicyToMap(policy),
		"state_summary": s.apiRateLimitStateSummary(policy, state),
	}
}

func (s *APIServer) apiRateLimitStateSummary(policy apiRateLimitPolicy, state apiRateLimitState) []map[string]any {
	now := s.cfg.Now().UTC()
	summary := []map[string]any{}
	for key, window := range state.Windows {
		groupName := apiRateLimitGroupFromKey(key)
		group, ok := policy.Groups[groupName]
		if !ok {
			continue
		}
		start, err := time.Parse(apiTimeLayout, window.WindowStart)
		if err != nil {
			continue
		}
		resetAt := start.Add(time.Duration(group.WindowSeconds) * time.Second)
		if !resetAt.After(now) {
			continue
		}
		summary = append(summary, map[string]any{
			"key":          key,
			"group":        groupName,
			"window_start": window.WindowStart,
			"count":        window.Count,
			"reset_at":     resetAt.UTC().Format(apiTimeLayout),
		})
	}
	sort.Slice(summary, func(i, j int) bool {
		left := fmt.Sprint(summary[i]["reset_at"])
		right := fmt.Sprint(summary[j]["reset_at"])
		if left == right {
			return fmt.Sprint(summary[i]["key"]) < fmt.Sprint(summary[j]["key"])
		}
		return left > right
	})
	if len(summary) > 100 {
		summary = summary[:100]
	}
	return summary
}

func validateAPIRateLimitPolicyBody(w http.ResponseWriter, body map[string]any) (apiRateLimitPolicy, bool) {
	policy, ok := apiRateLimitPolicyFromConfig(body)
	if !ok {
		writeValidation(w, "body", "invalid value")
		return apiRateLimitPolicy{}, false
	}
	return policy, true
}

func apiRateLimitPolicyFromConfig(raw map[string]any) (apiRateLimitPolicy, bool) {
	if raw == nil {
		return defaultAPIRateLimitPolicy(), true
	}
	if len(raw) != 2 {
		return apiRateLimitPolicy{}, false
	}
	enabled, ok := raw["enabled"].(bool)
	if !ok {
		return apiRateLimitPolicy{}, false
	}
	groupsRaw, ok := raw["groups"].(map[string]any)
	if !ok || len(groupsRaw) != len(apiRateLimitGroupNames()) {
		return apiRateLimitPolicy{}, false
	}
	policy := apiRateLimitPolicy{Enabled: enabled, Groups: map[string]apiRateLimitGroup{}}
	for _, name := range apiRateLimitGroupNames() {
		groupRaw, ok := groupsRaw[name].(map[string]any)
		if !ok || len(groupRaw) != 2 {
			return apiRateLimitPolicy{}, false
		}
		window, ok := apiRateLimitInt(groupRaw["window_seconds"], 1, 86400)
		if !ok {
			return apiRateLimitPolicy{}, false
		}
		maxRequests, ok := apiRateLimitInt(groupRaw["max_requests"], 1, 100000)
		if !ok {
			return apiRateLimitPolicy{}, false
		}
		policy.Groups[name] = apiRateLimitGroup{WindowSeconds: window, MaxRequests: maxRequests}
	}
	for name := range groupsRaw {
		if !apiRateLimitGroupNameAllowed(name) {
			return apiRateLimitPolicy{}, false
		}
	}
	return policy, true
}

func apiRateLimitInt(value any, min, max int) (int, bool) {
	var v int
	switch typed := value.(type) {
	case int:
		v = typed
	case int64:
		v = int(typed)
	case float64:
		if typed != math.Trunc(typed) {
			return 0, false
		}
		v = int(typed)
	default:
		return 0, false
	}
	return v, v >= min && v <= max
}

func apiRateLimitPolicyToMap(policy apiRateLimitPolicy) map[string]any {
	groups := map[string]any{}
	for _, name := range apiRateLimitGroupNames() {
		group := policy.Groups[name]
		groups[name] = map[string]any{"window_seconds": group.WindowSeconds, "max_requests": group.MaxRequests}
	}
	return map[string]any{"enabled": policy.Enabled, "groups": groups}
}

func apiRateLimitPoliciesEqual(a, b apiRateLimitPolicy) bool {
	if a.Enabled != b.Enabled {
		return false
	}
	for _, name := range apiRateLimitGroupNames() {
		if a.Groups[name] != b.Groups[name] {
			return false
		}
	}
	return true
}

func defaultAPIRateLimitPolicy() apiRateLimitPolicy {
	return apiRateLimitPolicy{Enabled: true, Groups: map[string]apiRateLimitGroup{
		"login":   {WindowSeconds: 60, MaxRequests: 10},
		"read":    {WindowSeconds: 60, MaxRequests: 600},
		"trigger": {WindowSeconds: 60, MaxRequests: 60},
		"operate": {WindowSeconds: 60, MaxRequests: 120},
		"config":  {WindowSeconds: 60, MaxRequests: 60},
		"admin":   {WindowSeconds: 60, MaxRequests: 60},
	}}
}

func apiRateLimitGroupNames() []string {
	return []string{"login", "read", "trigger", "operate", "config", "admin"}
}

func apiRateLimitGroupNameAllowed(name string) bool {
	for _, candidate := range apiRateLimitGroupNames() {
		if name == candidate {
			return true
		}
	}
	return false
}

func apiRateLimitGroupForRequest(method, path string) (string, bool) {
	if strings.ToUpper(method) == http.MethodPost && (path == "/api/login" || path == "/api/login/totp") {
		return "login", true
	}
	if strings.ToUpper(method) == http.MethodPost && path == "/api/webhook" {
		return "trigger", true
	}
	return apiTokenRequiredScope(strings.ToUpper(method), path)
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
	if action != "" && !validAuditAction(action) {
		writeValidation(w, "action", "invalid value")
		return
	}
	if result != "" && result != "success" && result != "failure" && result != "denied" {
		writeValidation(w, "result", "invalid value")
		return
	}
	records := readJSONLines(filepath.Join(s.cfg.StateDir, ".audit_log"))
	type indexedAuditRecord struct {
		index  int
		record map[string]any
	}
	filtered := []indexedAuditRecord{}
	for i, record := range records {
		if actor != "" && !auditActorMatches(record, actor) {
			continue
		}
		if action != "" && fmt.Sprint(record["action"]) != action {
			continue
		}
		if result != "" && fmt.Sprint(record["result"]) != result {
			continue
		}
		filtered = append(filtered, indexedAuditRecord{index: i, record: record})
	}
	sort.Slice(filtered, func(i, j int) bool {
		left := auditRecordTimestamp(filtered[i].record)
		right := auditRecordTimestamp(filtered[j].record)
		if left != right {
			return left > right
		}
		return filtered[i].index > filtered[j].index
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
	out := []map[string]any{}
	for _, entry := range filtered[start:end] {
		out = append(out, entry.record)
	}
	writeJSON(w, http.StatusOK, map[string]any{"log": out, "total": total})
}

func validAuditAction(action string) bool {
	switch action {
	case "login_success", "login_failure", "permission_denied", "totp_required", "logout", "session_revoke_all", "password_change", "token_create", "token_auth", "token_expired", "token_revoked_reject", "token_revoke", "totp_setup", "totp_failure", "totp_enabled", "totp_disabled", "build_trigger", "build_force_trigger", "config_update", "rate_limit_update", "approval_pending", "approval_approved", "approval_rejected", "approval_expired", "share_link_create", "share_link_revoke":
		return true
	default:
		return false
	}
}

func auditActorMatches(record map[string]any, actor string) bool {
	if actor == "anonymous" {
		return fmt.Sprint(record["actor_type"]) == "anonymous" && record["actor_id"] == nil
	}
	if value, ok := record["actor_id"].(string); ok {
		return value == actor
	}
	if value, ok := record["actor"].(string); ok {
		return value == actor
	}
	return false
}

func auditRecordTimestamp(record map[string]any) string {
	if value, ok := record["timestamp"].(string); ok {
		return value
	}
	if value, ok := record["at"].(string); ok {
		return value
	}
	return ""
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
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if cfg.ScheduleIntervalSeconds == body.IntervalSeconds {
		writeJSON(w, http.StatusOK, map[string]any{"message": "No changes", "interval_seconds": body.IntervalSeconds})
		return
	}
	cfg.ScheduleIntervalSeconds = body.IntervalSeconds
	if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".server_config"), cfg, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	changes := map[string]any{"interval_seconds": body.IntervalSeconds}
	if err := s.applyScheduleInterval(r.Context(), body.IntervalSeconds); err != nil {
		changes["result"] = "partial_failure"
		changes["error"] = "systemd_update_failed"
		now := s.nowString()
		if logErr := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: now, Type: "schedule_interval", Changes: changes}); logErr != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if auditErr := appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": now, "action": "config_update", "result": "failure", "target_id": "schedule_interval"}); auditErr != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		log.Printf("SCHEDULE_INTERVAL_APPLY_FAILED")
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	changes["result"] = "success"
	changes["error"] = nil
	now := s.nowString()
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: now, Type: "schedule_interval", Changes: changes}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": now, "action": "config_update", "result": "success", "target_id": "schedule_interval"}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Schedule interval updated", "interval_seconds": body.IntervalSeconds})
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

func (s *APIServer) applyScheduleInterval(ctx context.Context, seconds int) error {
	dropInDir := filepath.Join(s.cfg.SystemdDir, "adlaire-ci.timer.d")
	content := fmt.Sprintf("[Timer]\nOnUnitActiveSec=%ds\nPersistent=true\n", seconds)
	if err := atomicWriteText(filepath.Join(dropInDir, "override.conf"), content, 0644); err != nil {
		return err
	}
	cmdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if result, err := s.cfg.CommandRunner(cmdCtx, "systemctl", "daemon-reload"); err != nil || result.ExitCode != 0 {
		return errors.New("systemd daemon-reload failed")
	}
	if result, err := s.cfg.CommandRunner(cmdCtx, "systemctl", "restart", "adlaire-ci.timer"); err != nil || result.ExitCode != 0 {
		return errors.New("systemd restart failed")
	}
	result, err := s.cfg.CommandRunner(cmdCtx, "systemctl", "show", "adlaire-ci.timer", "-p", "OnUnitActiveSec")
	if err != nil || result.ExitCode != 0 {
		return errors.New("systemd show failed")
	}
	if parseSystemdSeconds(result.Stdout) != seconds {
		return errors.New("systemd interval mismatch")
	}
	return nil
}

func parseSystemdSeconds(value string) int {
	value = strings.TrimSpace(value)
	if before, after, ok := strings.Cut(value, "="); ok && before != "" {
		value = strings.TrimSpace(after)
	}
	value = strings.TrimSuffix(value, "s")
	if seconds, err := strconv.Atoi(value); err == nil {
		return seconds
	}
	if strings.HasSuffix(value, "min") {
		minutes, err := strconv.Atoi(strings.TrimSuffix(value, "min"))
		if err == nil {
			return minutes * 60
		}
	}
	return -1
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
	statusPath := filepath.Join(s.cfg.StateDir, ".build_status.json")
	if st, ok, err := readBuildStatus(statusPath); err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return nil, false
	} else if ok {
		running := st.Running || lockExists(filepath.Join(s.cfg.StateDir, ".build_lock"))
		lastSHA := st.LastBlobSHA
		if lastSHA == nil {
			lastSHA = st.LastCommitSHA
		}
		return map[string]any{
			"last_sha":                lastSHA,
			"last_build_at":           coalesceString(st.LastFinishedAt, st.LastStartedAt),
			"last_build_status":       stringOr(st.Status, "none"),
			"last_target_status":      st.LastTargetStatus,
			"last_trigger":            st.LastTrigger,
			"last_deploy_status":      st.LastDeployStatus,
			"pending_transfers_count": st.PendingTransfersCount,
			"notify_pending_count":    st.NotifyPendingCount,
			"circuit_open":            st.CircuitOpen,
			"running":                 running,
		}, true
	}
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
	pendingTransfers, err := readJSONArrayCount(filepath.Join(s.cfg.StateDir, ".pending_transfers"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return nil, false
	}
	notifyPending, err := readJSONArrayCount(filepath.Join(s.cfg.StateDir, ".notify_pending"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return nil, false
	}
	history := readHistory(filepath.Join(s.cfg.StateDir, ".build_history"))
	latest := latestStatusSummaryHistory(history)
	resp := map[string]any{
		"last_sha":                nil,
		"last_build_at":           nil,
		"last_build_status":       "none",
		"last_target_status":      nil,
		"last_trigger":            nil,
		"last_deploy_status":      nil,
		"pending_transfers_count": pendingTransfers,
		"notify_pending_count":    notifyPending,
		"circuit_open":            circuit.Open,
		"running":                 state.Running || lockExists(filepath.Join(s.cfg.StateDir, ".build_lock")),
	}
	if latest != nil {
		resp["last_sha"] = firstNonNil(latest.BlobSHA, latest.CommitSHA, latest.SHA)
		resp["last_build_at"] = firstNonEmpty(latest.FinishedAt, latest.StartedAt, latest.BuildAt)
		resp["last_build_status"] = normalizeHealthBuildStatus(latest.Status)
		resp["last_target_status"] = latest.Status
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
	target, err := s.selectedOutputTarget()
	if err != nil {
		writeOutputTargetError(w, err)
		return
	}
	meta, err := inspectOutput(target.Out)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	uptime := int64(s.cfg.Now().UTC().Sub(s.startedAt).Seconds())
	if uptime < 0 {
		uptime = 0
	}
	writeJSON(w, http.StatusOK, outputSysinfo(meta, uptime))
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
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
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
		if log.DurationSeconds < 0 {
			continue
		}
		if len(durations) >= n {
			break
		}
		durations = append(durations, log.DurationSeconds)
		status := log.Status
		if status == "" {
			status = normalizeHealthBuildStatus(log.TargetStatus)
		}
		recent = append(recent, map[string]any{
			"id":               log.ID,
			"build_at":         firstNonEmpty(log.FinishedAt, log.StartedAt),
			"duration_seconds": log.DurationSeconds,
			"status":           status,
			"target_status":    log.TargetStatus,
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
	var avg any
	var minValue any
	var maxValue any
	if len(durations) > 0 {
		avg = math.Round((float64(sum)/float64(len(durations)))*100) / 100
		minValue = min
		maxValue = max
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"n": n, "count": len(durations), "avg_seconds": avg,
		"min_seconds": minValue, "max_seconds": maxValue, "recent": recent,
	})
}

func (s *APIServer) handleStatsBuildTrends(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	n, ok := parseBoundedInt(w, r, "n", 100, 1, 1000)
	if !ok {
		return
	}
	var trend apiBuildTrendFile
	if err := readJSONIfExists(filepath.Join(s.cfg.StateDir, ".build_trends.json"), &trend); err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if trend.Samples == nil {
		trend.Samples = []apiBuildTrendSample{}
	}
	samples := trend.Samples
	if len(samples) > n {
		samples = samples[len(samples)-n:]
	}
	samples = append([]apiBuildTrendSample(nil), samples...)
	sort.Slice(samples, func(i, j int) bool {
		if samples[i].FinishedAt == samples[j].FinishedAt {
			return samples[i].BuildID < samples[j].BuildID
		}
		return samples[i].FinishedAt < samples[j].FinishedAt
	})
	writeJSON(w, http.StatusOK, map[string]any{
		"n":       n,
		"count":   len(samples),
		"samples": samples,
		"summary": buildTrendSummary(samples),
	})
}

func (s *APIServer) handleOutputMeta(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	target, err := s.selectedOutputTarget()
	if err != nil {
		writeOutputTargetError(w, err)
		return
	}
	meta, err := inspectOutput(target.Out)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if !meta.Exists {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	histories := outputMetaSuccessHistory(readHistory(filepath.Join(s.cfg.StateDir, ".build_history")), target)
	var history *apiHistoryRecord
	if len(histories) > 0 {
		history = &histories[0]
	}
	var log *apiBuildLog
	if history != nil {
		if found, err := readBuildLogByID(s.cfg.StateDir, history.ID); err == nil {
			log = &found
		}
	}
	resp := map[string]any{
		"size_bytes":        meta.SizeBytes,
		"mtime":             meta.MTime,
		"sha256":            "",
		"heading_count":     nil,
		"tables_count":      nil,
		"code_blocks_count": nil,
		"size_diff_bytes":   nil,
		"size_warn":         false,
		"build_warnings":    []string{},
		"build_id":          "",
		"commit_sha":        "",
		"build_at":          "",
	}
	if history != nil {
		resp["sha256"] = stringPtrValue(history.OutputSHA256)
		compare := outputMetaCompareHistory(histories, log)
		if compare != nil && compare.OutputSizeBytes != nil {
			resp["size_diff_bytes"] = meta.SizeBytes - *compare.OutputSizeBytes
		}
	}
	if log != nil {
		resp["build_warnings"] = log.Warnings
		commitSHA := ""
		if value, ok := firstCommitSHA(log.Commit).(string); ok {
			commitSHA = value
		}
		if log.Report != nil {
			resp["heading_count"] = log.Report.Headings
			resp["tables_count"] = log.Report.TablesCount
			resp["code_blocks_count"] = log.Report.CodeBlocksCount
			resp["size_warn"] = log.Report.SizeWarn
		}
		resp["build_id"] = log.ID
		resp["commit_sha"] = commitSHA
		resp["build_at"] = firstNonEmpty(log.FinishedAt, log.StartedAt)
		if log.BuildMeta != nil {
			resp["build_id"] = firstNonEmpty(log.BuildMeta.BuildID, log.ID)
			resp["commit_sha"] = firstNonEmpty(log.BuildMeta.CommitSHA, commitSHA)
			resp["build_at"] = firstNonEmpty(log.BuildMeta.BuildAt, log.FinishedAt, log.StartedAt)
		}
	}
	writeJSON(w, http.StatusOK, resp)
}

func outputMetaSuccessHistory(records []apiHistoryRecord, target apiOutputTarget) []apiHistoryRecord {
	out := []apiHistoryRecord{}
	for _, record := range records {
		if record.Branch != target.Branch || record.TargetFile != target.TargetFile {
			continue
		}
		switch record.Status {
		case "success", "success_deploy_pending":
			out = append(out, record)
		}
	}
	return out
}

func outputMetaCompareHistory(records []apiHistoryRecord, log *apiBuildLog) *apiHistoryRecord {
	if len(records) == 0 {
		return nil
	}
	currentID := ""
	if log != nil {
		currentID = log.ID
		if log.BuildMeta != nil && log.BuildMeta.BuildID != "" {
			currentID = log.BuildMeta.BuildID
		}
	}
	if currentID != "" && currentID == records[0].ID {
		if len(records) < 2 {
			return nil
		}
		return &records[1]
	}
	return &records[0]
}

func (s *APIServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	status, ok := s.statusPayload(w)
	if !ok {
		return
	}
	target, err := s.selectedOutputTarget()
	if err != nil {
		writeOutputTargetError(w, err)
		return
	}
	meta, err := inspectOutput(target.Out)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	uptime := int64(s.cfg.Now().UTC().Sub(s.startedAt).Seconds())
	if uptime < 0 {
		uptime = 0
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":   status,
		"sysinfo":  outputSysinfo(meta, uptime),
		"stats":    s.statsSummary(7),
		"schedule": map[string]any{"enabled": false, "next_run_at": nil},
		"alerts":   []any{},
	})
}

func (s *APIServer) selectedOutputTarget() (apiOutputTarget, error) {
	cfg, _, err := s.readBranchConfig()
	if err != nil {
		return apiOutputTarget{}, err
	}
	if len(cfg.BranchTargets) == 0 {
		return apiOutputTarget{}, errAPIOutputTargetUnavailable
	}
	targets := append([]apiBranchTarget(nil), cfg.BranchTargets...)
	sort.Slice(targets, func(i, j int) bool {
		if targets[i].Branch == targets[j].Branch {
			return targets[i].TargetFile < targets[j].TargetFile
		}
		return targets[i].Branch < targets[j].Branch
	})
	if status, ok, err := readBuildStatus(filepath.Join(s.cfg.StateDir, ".build_status.json")); err == nil && ok && status.LastBranch != nil && status.LastTargetFile != nil {
		for _, target := range targets {
			if target.Branch == *status.LastBranch && target.TargetFile == *status.LastTargetFile {
				return apiOutputTarget{Branch: target.Branch, TargetFile: target.TargetFile, Out: target.Out}, nil
			}
		}
	}
	target := targets[0]
	return apiOutputTarget{Branch: target.Branch, TargetFile: target.TargetFile, Out: target.Out}, nil
}

func writeOutputTargetError(w http.ResponseWriter, err error) {
	if errors.Is(err, errAPIOutputTargetUnavailable) {
		writeError(w, http.StatusInternalServerError, "Output target unavailable")
		return
	}
	writeError(w, http.StatusInternalServerError, "State file is corrupted")
}

func outputSysinfo(meta outputInspection, uptime int64) map[string]any {
	size := meta.SizeBytes
	mtime := meta.MTime
	if !meta.Exists {
		size = 0
		mtime = nil
	}
	return map[string]any{
		"output_size_bytes": size,
		"output_mtime":      mtime,
		"uptime_seconds":    uptime,
	}
}

func (s *APIServer) statsSummary(days int) map[string]any {
	records := s.historyWithinDays(days)
	success := 0
	failure := 0
	durations := []int64{}
	times := []time.Time{}
	for _, record := range records {
		switch statusCategory(record.Status) {
		case "success":
			success++
		case "failure":
			failure++
		}
		if record.DurationSeconds >= 0 {
			durations = append(durations, record.DurationSeconds)
		}
		if at := historyRecordTime(record); !at.IsZero() {
			times = append(times, at)
		}
	}
	total := len(records)
	var successRate any
	if total > 0 {
		successRate = math.Round((float64(success)/float64(total))*1000) / 1000
	}
	var avgIntervalMinutes any
	if len(times) >= 2 {
		sort.Slice(times, func(i, j int) bool { return times[i].After(times[j]) })
		var totalMinutes float64
		for i := 0; i < len(times)-1; i++ {
			totalMinutes += times[i].Sub(times[i+1]).Minutes()
		}
		avgIntervalMinutes = math.Round((totalMinutes/float64(len(times)-1))*100) / 100
	}
	var avgDuration any
	var maxDuration any
	if len(durations) > 0 {
		var sum int64
		var max int64
		for i, duration := range durations {
			sum += duration
			if i == 0 || duration > max {
				max = duration
			}
		}
		avgDuration = math.Round((float64(sum)/float64(len(durations)))*100) / 100
		maxDuration = max
	}
	return map[string]any{
		"days":                 days,
		"total_builds":         total,
		"success_count":        success,
		"failure_count":        failure,
		"success_rate":         successRate,
		"avg_interval_minutes": avgIntervalMinutes,
		"avg_duration_seconds": avgDuration,
		"max_duration_seconds": maxDuration,
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

func (s *APIServer) handleBuild(force bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
			return
		}
		maintenance, err := maintenanceEnabled(filepath.Join(s.cfg.StateDir, ".maintenance"))
		if err != nil {
			writeError(w, http.StatusServiceUnavailable, "maintenance_unavailable")
			return
		}
		if maintenance {
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
		if !state.Running && state.CurrentBuildID == nil && lockExists(filepath.Join(s.cfg.StateDir, ".build_lock")) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		actorType, actorID := apiAuditActor(r)
		requestedBy := "admin"
		if actorType == "api_token" {
			if id, ok := actorID.(string); ok && id != "" {
				requestedBy = id
			}
		}
		queueID, duplicate, err := s.enqueueBuildRequest(statePath, "manual", requestedBy, map[string]any{"force": force})
		if err != nil {
			if errors.Is(err, errAPIQueueFull) {
				writeError(w, http.StatusTooManyRequests, "queue_full")
				return
			}
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		action := "build_trigger"
		if force {
			action = "build_force_trigger"
		}
		if err := s.appendSecurityAudit(r, action, actorType, actorID, "build", queueID, "success"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		status := http.StatusAccepted
		message := "Build queued"
		if duplicate {
			status = http.StatusOK
			message = "Already queued"
		}
		writeJSON(w, status, map[string]any{"message": message, "queue_id": queueID, "queued": true, "dispatch": s.dispatchRunner(queueID)})
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
	logRecord, running, found, err := s.selectBuildStreamLog()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	at := firstNonEmpty(logRecord.FinishedAt, logRecord.StartedAt)
	if at == "" {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	status := buildStreamStatus(logRecord, running)
	var duration *int64
	if status != "running" {
		value := logRecord.DurationSeconds
		if value < 0 {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		duration = &value
	}
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("X-Accel-Buffering", "no")
	flusher, _ := w.(http.Flusher)
	for _, line := range buildStreamLines(logRecord) {
		if !writeSSEFrame(w, flusher, apiBuildStreamLogFrame{Type: "log", Line: truncateStreamRunes(line, 4000), At: at}) {
			return
		}
	}
	_ = writeSSEFrame(w, flusher, apiBuildStreamEndFrame{Type: "end", Status: status, DurationSeconds: duration})
}

func (s *APIServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	if !validateQueryKeys(w, r, map[string]bool{"n": true, "q": true}) {
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
	if !validateLogSearchQuery(w, r) {
		return
	}
	q := r.URL.Query().Get("q")
	from, fromValue, ok := parseLogSearchDate(w, r.URL.Query().Get("from"), "from")
	if !ok {
		return
	}
	to, toValue, ok := parseLogSearchDate(w, r.URL.Query().Get("to"), "to")
	if !ok {
		return
	}
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		writeValidation(w, "from", "must be before or equal to to")
		return
	}
	levelRaw := strings.TrimSpace(r.URL.Query().Get("level"))
	level := strings.ToUpper(levelRaw)
	if level == "WARN" || level == "WARNING" {
		level = "WARNING"
	}
	var levelValue any
	if level != "" {
		levelValue = level
	}
	results := []map[string]any{}
	for _, log := range s.readBuildLogsNewest() {
		buildAt := firstNonEmpty(log.FinishedAt, log.StartedAt)
		buildTime, _ := time.Parse(apiTimeLayout, buildAt)
		if !from.IsZero() && buildTime.Before(from) {
			continue
		}
		if !to.IsZero() && buildTime.After(to) {
			continue
		}
		lines := logSearchLines(log)
		lines = filterContainsFold(lines, q)
		if level != "" {
			lines = filterContains(lines, "["+level+"]")
		}
		if len(lines) == 0 {
			continue
		}
		results = append(results, map[string]any{"id": log.ID, "build_at": buildAt, "lines": lines})
	}
	writeJSON(w, http.StatusOK, map[string]any{"query": q, "from": fromValue, "to": toValue, "level": levelValue, "results": results})
}

func validateLogSearchQuery(w http.ResponseWriter, r *http.Request) bool {
	allowed := map[string]bool{"q": true, "from": true, "to": true, "level": true}
	if !validateQueryKeys(w, r, allowed) {
		return false
	}
	q := r.URL.Query().Get("q")
	if len(q) > 500 {
		writeValidation(w, "q", "must be 500 characters or less")
		return false
	}
	if strings.ContainsAny(q, "\x00\r\n") {
		writeValidation(w, "q", "invalid value")
		return false
	}
	level := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("level")))
	if level == "WARN" {
		level = "WARNING"
	}
	if level != "" && level != "INFO" && level != "WARNING" && level != "ERROR" && level != "DEBUG" {
		writeValidation(w, "level", "invalid value")
		return false
	}
	return true
}

func parseLogSearchDate(w http.ResponseWriter, value, field string) (time.Time, any, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil, true
	}
	day, err := time.ParseInLocation("2006-01-02", value, time.UTC)
	if err != nil {
		writeValidation(w, field, "invalid value")
		return time.Time{}, nil, false
	}
	if field == "to" {
		day = day.Add(24*time.Hour - time.Second)
	}
	return day, value, true
}

func logSearchLines(log apiBuildLog) []string {
	lines := append([]string{}, splitLines(log.Pipeline.Stdout)...)
	lines = append(lines, splitLines(log.Pipeline.Stderr)...)
	lines = append(lines, log.Warnings...)
	if log.Error != nil && *log.Error != "" {
		lines = append(lines, *log.Error)
	}
	return lines
}

func (s *APIServer) handleHistory(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	if !validateHistoryQuery(w, r) {
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
	history = decorateHistoryRecords(history)
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
	writeJSON(w, http.StatusOK, map[string]any{"exported_at": s.nowString(), "history": decorateHistoryRecords(readHistory(filepath.Join(s.cfg.StateDir, ".build_history")))})
}

func (s *APIServer) handleHistoryPath(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/history/")
	id, suffix, ok := strings.Cut(rest, "/")
	if !ok || id == "" {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if !validSimpleID(id) {
		writeValidation(w, "id", "invalid value")
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
			var nextComment *string
			if comment == "" {
				nextComment = nil
			} else {
				nextComment = &comment
			}
			if nullableStringEqual(log.Comment, nextComment) {
				writeJSON(w, http.StatusOK, map[string]any{"message": "Comment saved"})
				return
			}
			log.Comment = nextComment
			if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".build_logs", id+".json"), log, 0600); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "history_comment", Changes: map[string]any{"id": id}})
			writeJSON(w, http.StatusOK, map[string]any{"message": "Comment saved"})
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
		if log.Flagged == *body.Flagged {
			writeJSON(w, http.StatusOK, map[string]any{"message": "Flag updated"})
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
		writeJSON(w, http.StatusOK, map[string]any{"message": "Flag updated"})
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
		if stringSlicesEqual(log.Tags, tags) {
			writeJSON(w, http.StatusOK, map[string]any{"message": "Tags updated"})
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
		writeJSON(w, http.StatusOK, map[string]any{"message": "Tags updated"})
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
	current, err := readCircuitState(filepath.Join(s.cfg.StateDir, ".build_circuit_state"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if circuitStateInitial(current) {
		writeJSON(w, http.StatusOK, map[string]any{"message": "Circuit breaker reset", "open": false, "consecutive_failures": 0})
		return
	}
	state := apiCircuitState{Open: false, ConsecutiveFailures: 0}
	if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, ".build_circuit_state"), state, 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := appendConfigAndAuditLog(s.cfg.StateDir, s.nowString(), "circuit_breaker_reset", map[string]any{"open": false, "consecutive_failures": 0}, "circuit_breaker_reset"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Circuit breaker reset", "open": false, "consecutive_failures": 0})
}

func circuitStateInitial(state apiCircuitState) bool {
	return !state.Open && state.ConsecutiveFailures == 0 && state.OpenedAt == nil && state.LastFailureAt == nil && state.LastError == nil
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
	if err := appendConfigAndAuditLog(s.cfg.StateDir, s.nowString(), "maintenance", map[string]any{"enabled": true, "reason": reason}, "maintenance"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Maintenance mode enabled", "since": since})
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
	if err := appendConfigAndAuditLog(s.cfg.StateDir, s.nowString(), "maintenance", map[string]any{"enabled": false}, "maintenance"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Maintenance mode disabled"})
}

func appendConfigAndAuditLog(stateDir, at, logType string, changes map[string]any, targetID string) error {
	if err := appendJSONLine(filepath.Join(stateDir, ".config_log"), apiConfigLogRecord{At: at, Type: logType, Changes: changes}); err != nil {
		return err
	}
	return appendJSONLine(filepath.Join(stateDir, ".audit_log"), map[string]any{"at": at, "action": "config_update", "result": "success", "target_id": targetID})
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
		writeJSON(w, http.StatusOK, map[string]any{"active": state.ActiveQueueEntry, "queued": sortedQueueEntries(state.Queued), "max_size": s.queueMaxSize()})
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

func (s *APIServer) handleApprovals(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	records, err := readLatestApprovals(filepath.Join(s.cfg.StateDir, ".approval_queue"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].CreatedAt == records[j].CreatedAt {
			return records[i].ID < records[j].ID
		}
		return records[i].CreatedAt > records[j].CreatedAt
	})
	writeJSON(w, http.StatusOK, map[string]any{"approvals": records})
}

func (s *APIServer) handleApprovalPath(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/approvals/"), "/")
	id, action, ok := strings.Cut(rest, "/")
	if !ok || !validSimpleID(id) {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if !method(w, r, http.MethodPost) || !rejectBody(w, r) {
		return
	}
	switch action {
	case "approve":
		s.approveBuild(w, r, id)
	case "reject":
		s.rejectBuildApproval(w, r, id)
	default:
		writeError(w, http.StatusNotFound, "Not found")
	}
}

func (s *APIServer) approveBuild(w http.ResponseWriter, r *http.Request, id string) {
	approvalPath := filepath.Join(s.cfg.StateDir, ".approval_queue")
	record, found, err := latestApprovalByID(approvalPath, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if record.Status != "pending" {
		writeError(w, http.StatusConflict, "Conflict")
		return
	}
	actor := apiActorID(r)
	payload := map[string]any{
		"approval_id":       record.ID,
		"branch":            record.Branch,
		"sha":               record.SHA,
		"target":            record.Target,
		"requested_trigger": record.RequestedTrigger,
		"requested_force":   record.RequestedForce,
		"delivery_id":       record.DeliveryID,
	}
	queueID, duplicate, err := s.enqueueApprovalRequest(filepath.Join(s.cfg.StateDir, ".build_state"), actor, payload)
	if err != nil {
		if errors.Is(err, errAPIQueueFull) {
			writeError(w, http.StatusTooManyRequests, "queue_full")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	now := s.nowString()
	record.Status = "approved"
	record.DecidedAt = &now
	record.DecidedBy = &actor
	record.QueueID = &queueID
	record.Reason = nil
	if err := appendJSONLine(approvalPath, record); err != nil {
		log.Printf("APPROVAL_APPROVED_RECORD_FAILED: path=%s", approvalPath)
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": now, "actor": actor, "action": "approval_approved", "result": "success", "target_type": "approval", "target_id": id}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	message := "Build approved"
	if duplicate {
		message = "Build approved"
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": message, "queued": true, "queue_id": queueID, "dispatch": s.dispatchRunner(queueID)})
}

func (s *APIServer) rejectBuildApproval(w http.ResponseWriter, r *http.Request, id string) {
	approvalPath := filepath.Join(s.cfg.StateDir, ".approval_queue")
	record, found, err := latestApprovalByID(approvalPath, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if record.Status != "pending" {
		writeError(w, http.StatusConflict, "Conflict")
		return
	}
	actor := apiActorID(r)
	now := s.nowString()
	reason := "rejected"
	record.Status = "rejected"
	record.DecidedAt = &now
	record.DecidedBy = &actor
	record.QueueID = nil
	record.Reason = &reason
	if err := appendJSONLine(approvalPath, record); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	history := apiHistoryRecord{ID: id, Status: "approval_rejected", Trigger: "approval", FinishedAt: now, BlobSHA: &record.SHA, DurationSeconds: 0, Tags: []string{}}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".build_history"), history); err != nil {
		log.Printf("APPROVAL_REJECT_HISTORY_FAILED: path=%s", filepath.Join(s.cfg.StateDir, ".build_history"))
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": now, "actor": actor, "action": "approval_rejected", "result": "success", "target_type": "approval", "target_id": id}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Build rejected"})
}

func (s *APIServer) handlePATStatus(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	configured, err := patConfigured(filepath.Join(s.cfg.StateDir, ".github_token"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	cfg, err := s.readMergedConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	expiresInDays, err := expiresInDays(cfg.PATExpiresAt, s.cfg.Now().UTC())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"configured":      configured,
		"expires_at":      cfg.PATExpiresAt,
		"expires_in_days": expiresInDays,
	})
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
	result, err := s.verifyGitHubPAT(r.Context(), token)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "PAT verification failed")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *APIServer) verifyGitHubPAT(ctx context.Context, token string) (map[string]any, error) {
	base := strings.TrimRight(s.cfg.GitHubAPIBaseURL, "/")
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, base+"/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "adlaire-ci")
	resp, err := s.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	checkedAt := s.nowString()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return map[string]any{"valid": true, "checked_at": checkedAt, "scopes": githubScopes(resp.Header.Get("X-OAuth-Scopes"))}, nil
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return map[string]any{"valid": false, "checked_at": checkedAt, "scopes": []string{}}, nil
	}
	return nil, fmt.Errorf("github pat verify status %d", resp.StatusCode)
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
	backup, err := s.backupObject()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, backup)
}

func (s *APIServer) handleRestore(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	rawBody := map[string]json.RawMessage{}
	if !decodeBody(w, r, &rawBody, true) {
		return
	}
	plan, changed, ok := s.restorePlan(w, rawBody)
	if !ok {
		return
	}
	if !changed {
		writeJSON(w, http.StatusOK, map[string]string{"message": "No changes"})
		return
	}
	for _, write := range plan.JSONWrites {
		if write.Delete {
			if err := removeIfExists(filepath.Join(s.cfg.StateDir, write.File)); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			continue
		}
		if err := atomicWriteJSON(filepath.Join(s.cfg.StateDir, write.File), write.Value, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}
	for _, write := range plan.SecretWrites {
		path := filepath.Join(s.cfg.StateDir, write.File)
		if write.Delete {
			if err := removeIfExists(path); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			continue
		}
		if err := atomicWriteText(path, write.Value, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}
	now := s.nowString()
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: now, Type: "restore", Changes: plan.ChangeSummary}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{"at": now, "action": "config_update", "result": "success", "target_type": "config", "target_id": "restore"}); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Config restored"})
}

func (s *APIServer) handleDiagnostics(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !rejectBody(w, r) {
		return
	}
	items, err := s.diagnosticsItems(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"checked_at": s.nowString(), "items": items})
}

func (s *APIServer) diagnosticsItems(ctx context.Context) ([]map[string]any, error) {
	token, tokenErr := readSecretText(filepath.Join(s.cfg.StateDir, ".github_token"))
	tokenConfigured := tokenErr == nil && strings.TrimSpace(token) != ""
	if tokenErr != nil && !errors.Is(tokenErr, os.ErrNotExist) {
		return nil, tokenErr
	}
	cfg, err := s.readMergedConfig()
	if err != nil {
		return nil, err
	}
	target, err := s.selectedOutputTarget()
	if err != nil {
		return nil, err
	}
	meta, err := inspectOutput(target.Out)
	if err != nil {
		return nil, err
	}
	return []map[string]any{
		s.diagnosticsPATItem(tokenConfigured, cfg.PATExpiresAt),
		s.diagnosticsGitHubAPIItem(ctx, token, tokenConfigured),
		diagnosticsOutputItem(meta),
		s.diagnosticsSystemdItem(ctx),
		diagnosticsWebhookItem(filepath.Join(s.cfg.StateDir, ".webhook_secret")),
	}, nil
}

func (s *APIServer) diagnosticsPATItem(configured bool, expiresAt *string) map[string]any {
	if !configured {
		return diagnosticsItem("pat", "error", "PAT is not configured")
	}
	days, err := expiresInDays(expiresAt, s.cfg.Now().UTC())
	if err != nil {
		return diagnosticsItem("pat", "error", "PAT expiration is invalid")
	}
	if days != nil && *days <= 0 {
		return diagnosticsItem("pat", "error", "PAT is expired")
	}
	if days != nil && *days <= 7 {
		return diagnosticsItem("pat", "warn", "PAT expires within 7 days")
	}
	return diagnosticsItem("pat", "ok", "PAT is valid")
}

func (s *APIServer) diagnosticsGitHubAPIItem(ctx context.Context, token string, configured bool) map[string]any {
	if !configured {
		return diagnosticsItem("github_api", "error", "GitHub API not checked because PAT is not configured")
	}
	if _, err := s.fetchGitHubRateLimit(ctx, token); err != nil {
		return diagnosticsItem("github_api", "error", "GitHub API unreachable")
	}
	return diagnosticsItem("github_api", "ok", "GitHub API reachable")
}

func diagnosticsOutputItem(meta outputInspection) map[string]any {
	if !meta.Exists {
		return diagnosticsItem("output_file", "error", "Output site is missing")
	}
	return diagnosticsItem("output_file", "ok", fmt.Sprintf("Output site exists (%s)", humanBytes(meta.SizeBytes)))
}

func humanBytes(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	value := float64(size)
	for _, unit := range []string{"KB", "MB", "GB", "TB"} {
		value = value / 1024
		if value < 1024 || unit == "TB" {
			return fmt.Sprintf("%.1f %s", value, unit)
		}
	}
	return fmt.Sprintf("%d B", size)
}

func (s *APIServer) diagnosticsSystemdItem(ctx context.Context) map[string]any {
	cmdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	timer, err := s.cfg.CommandRunner(cmdCtx, "systemctl", "is-active", "adlaire-ci.timer")
	if err != nil || timer.ExitCode != 0 || strings.TrimSpace(timer.Stdout) != "active" {
		return diagnosticsItem("systemd", "error", "systemd timer or service is not ready")
	}
	cmdCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	service, err := s.cfg.CommandRunner(cmdCtx, "systemctl", "cat", "adlaire-ci.service")
	if err != nil || service.ExitCode != 0 {
		return diagnosticsItem("systemd", "error", "systemd timer or service is not ready")
	}
	return diagnosticsItem("systemd", "ok", "adlaire-ci.timer is active and adlaire-ci.service is installed")
}

func diagnosticsWebhookItem(path string) map[string]any {
	configured, err := patConfigured(path)
	if err == nil && configured {
		return diagnosticsItem("webhook", "ok", "Webhook secret configured")
	}
	return diagnosticsItem("webhook", "warn", "Webhook URL not configured")
}

func diagnosticsItem(name, status, message string) map[string]any {
	return map[string]any{"name": name, "status": status, "message": message}
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
	result, err := s.fetchGitHubRateLimit(r.Context(), token)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "GitHub rate limit check failed")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *APIServer) fetchGitHubRateLimit(ctx context.Context, token string) (map[string]any, error) {
	base := strings.TrimRight(s.cfg.GitHubAPIBaseURL, "/")
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, base+"/rate_limit", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "adlaire-ci")
	resp, err := s.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
		return nil, fmt.Errorf("github rate limit status %d", resp.StatusCode)
	}
	var payload struct {
		Resources struct {
			Core struct {
				Limit     int   `json:"limit"`
				Remaining int   `json:"remaining"`
				Reset     int64 `json:"reset"`
				Used      int   `json:"used"`
			} `json:"core"`
		} `json:"resources"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Resources.Core.Reset <= 0 {
		return nil, errors.New("missing core rate limit")
	}
	return map[string]any{
		"limit":     payload.Resources.Core.Limit,
		"remaining": payload.Resources.Core.Remaining,
		"reset_at":  time.Unix(payload.Resources.Core.Reset, 0).UTC().Format(apiTimeLayout),
		"used":      payload.Resources.Core.Used,
	}, nil
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
	target, err := s.selectedOutputTarget()
	if err != nil {
		writeOutputTargetError(w, err)
		return
	}
	outputBytes, _ := fileTreeStats(target.Out, func(path string) bool { return true })
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
	sort.Slice(events, func(i, j int) bool {
		left := firstNonEmpty(fmt.Sprint(events[i]["timestamp"]), fmt.Sprint(events[i]["at"]))
		right := firstNonEmpty(fmt.Sprint(events[j]["timestamp"]), fmt.Sprint(events[j]["at"]))
		return left > right
	})
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
	webhookActor := "webhook"
	webhookAuth := apiAuthInfo{AuthType: "webhook", Actor: &webhookActor}
	if capture, ok := r.Context().Value(apiAuthCaptureContextKey{}).(*apiAuthInfo); ok {
		*capture = webhookAuth
	}
	if !s.enforceAPIRateLimit(w, r, webhookAuth, "trigger") {
		return
	}
	if !validateRouteQuery(w, r) {
		return
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	eventID := s.newID("wh")
	deliveryID := strings.TrimSpace(r.Header.Get("X-GitHub-Delivery"))
	eventName := r.Header.Get("X-GitHub-Event")
	if deliveryID == "" || eventName == "" {
		writeError(w, http.StatusUnprocessableEntity, "Validation failed")
		return
	}
	maintenance, err := maintenanceEnabled(filepath.Join(s.cfg.StateDir, ".maintenance"))
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "maintenance_unavailable")
		return
	}
	if maintenance {
		writeError(w, http.StatusServiceUnavailable, "maintenance")
		return
	}
	ref, _ := payload["ref"].(string)
	branch := ""
	if strings.HasPrefix(ref, "refs/heads/") {
		branch = strings.TrimPrefix(ref, "refs/heads/")
	}
	sha, _ := payload["after"].(string)
	repository := webhookRepository(payload)
	if eventName != "push" {
		if err := s.appendWebhookEvent(eventID, deliveryID, eventName, nullableString(ref), nullableString(branch), nullableString(sha), repository, false, nil, "ignored_event"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]any{"message": "Webhook ignored", "queued": false, "event_id": eventID})
		return
	}
	if branch == "" {
		if err := s.appendWebhookEvent(eventID, deliveryID, eventName, nullableString(ref), nil, nullableString(sha), repository, false, nil, "ignored_branch"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]any{"message": "Webhook ignored", "queued": false, "event_id": eventID})
		return
	}
	if !isLowerHex(sha, 40) {
		writeValidation(w, "after", "invalid value")
		return
	}
	allowed, err := s.webhookBranchAllowed(branch)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	if !allowed {
		if err := s.appendWebhookEvent(eventID, deliveryID, eventName, &ref, &branch, &sha, repository, false, nil, "ignored_branch"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]any{"message": "Webhook ignored", "queued": false, "event_id": eventID})
		return
	}
	id, duplicate, err := s.enqueueBuildRequest(filepath.Join(s.cfg.StateDir, ".build_state"), "webhook", "webhook", map[string]any{"delivery_id": deliveryID, "branch": branch, "sha": sha})
	if err != nil {
		if errors.Is(err, errAPIQueueFull) {
			_ = s.appendWebhookEvent(eventID, deliveryID, eventName, &ref, &branch, &sha, repository, false, nil, "queue_full")
			writeError(w, http.StatusTooManyRequests, "queue_full")
			return
		}
		if errors.Is(err, errAPIQueueConflict) {
			writeError(w, http.StatusConflict, "Conflicting delivery")
			return
		}
		writeError(w, http.StatusInternalServerError, "State file is corrupted")
		return
	}
	result := "queued"
	message := "Webhook accepted"
	if duplicate {
		result = "duplicate"
		message = "Webhook already queued"
	}
	eventLogFailed := false
	if err := s.appendWebhookEvent(eventID, deliveryID, eventName, &ref, &branch, &sha, repository, true, &id, result); err != nil {
		log.Printf("WEBHOOK_EVENT_LOG_WRITE_FAILED")
		eventLogFailed = true
	}
	if err := s.appendSecurityAudit(r, "build_trigger", "webhook", "webhook", "build", id, "success"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	resp := map[string]any{"message": message, "queued": true, "event_id": eventID, "queue_id": id, "dispatch": s.dispatchRunner(id)}
	if eventLogFailed {
		resp["event_log_failed"] = true
	}
	writeJSON(w, http.StatusAccepted, resp)
}

func (s *APIServer) webhookBranchAllowed(branch string) (bool, error) {
	cfg, _, err := s.readBranchConfig()
	if err != nil {
		return false, err
	}
	for _, target := range cfg.BranchTargets {
		if target.Branch == branch {
			return true, nil
		}
	}
	return false, nil
}

func (s *APIServer) appendWebhookEvent(id, deliveryID, eventName string, ref, branch, sha, repository *string, buildTriggered bool, queueID *string, result string) error {
	var errorCode *string
	if result == "queue_full" {
		value := "queue_full"
		errorCode = &value
	}
	record := map[string]any{
		"id":              id,
		"timestamp":       s.nowString(),
		"delivery_id":     deliveryID,
		"event":           eventName,
		"ref":             ref,
		"branch":          branch,
		"sha":             sha,
		"repository":      repository,
		"build_triggered": buildTriggered,
		"queued_id":       queueID,
		"result":          result,
		"error_code":      errorCode,
	}
	return appendJSONLine(filepath.Join(s.cfg.StateDir, ".webhook_events.json"), record)
}

func webhookRepository(payload map[string]any) *string {
	repository, ok := payload["repository"].(map[string]any)
	if !ok {
		return nil
	}
	name, _ := repository["name"].(string)
	ownerMap, _ := repository["owner"].(map[string]any)
	owner, _ := ownerMap["login"].(string)
	if owner == "" || name == "" {
		return nil
	}
	value := owner + "/" + name
	return &value
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
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
		now := s.nowString()
		if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: now, Type: "snapshot_delete", Changes: map[string]any{"id": id}}); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		actorType, actorID := apiAuditActor(r)
		if err := s.appendSecurityAudit(r, "snapshot_delete", actorType, actorID, "snapshot", id, "success"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
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
			out = append(out, map[string]any{
				"id":           token.ID,
				"label":        token.Label,
				"scopes":       token.Scopes,
				"created_at":   token.CreatedAt,
				"last_used_at": token.LastUsedAt,
				"expires_at":   token.ExpiresAt,
				"revoked_at":   token.RevokedAt,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"tokens": out})
	case http.MethodPost:
		var body struct {
			Label     string   `json:"label"`
			Scopes    []string `json:"scopes"`
			ExpiresAt *string  `json:"expires_at"`
		}
		if !decodeBody(w, r, &body, true) {
			return
		}
		label := strings.TrimSpace(body.Label)
		if label == "" || len([]byte(label)) > 64 {
			writeValidation(w, "label", "invalid value")
			return
		}
		scopes, ok := normalizeTokenScopes(body.Scopes)
		if !ok {
			writeValidation(w, "scopes", "invalid value")
			return
		}
		expiresAt, ok := normalizeTokenExpiresAt(body.ExpiresAt, s.cfg.Now().UTC())
		if !ok {
			writeValidation(w, "expires_at", "invalid value")
			return
		}
		token, err := newAPIToken()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		hash := tokenHash(token)
		var record apiTokenRecord
		_, err = updateAPITokensLocked(path, func(tokens []apiTokenRecord) ([]apiTokenRecord, bool, error) {
			if len(tokens) >= 100 {
				return nil, false, errAPITokenLimit
			}
			for _, existing := range tokens {
				if tokenHashEqual(existing.TokenHash, hash) {
					return nil, false, errAPITokenCollision
				}
			}
			id, err := nextTokenID(tokens)
			if err != nil {
				return nil, false, err
			}
			record = apiTokenRecord{ID: id, Label: label, TokenHash: hash, Scopes: scopes, CreatedAt: s.nowString(), ExpiresAt: expiresAt}
			return append(tokens, record), true, nil
		})
		if err != nil {
			if errors.Is(err, errAPIStateLockConflict) {
				writeError(w, http.StatusConflict, "Conflict")
				return
			}
			if errors.Is(err, errAPITokenLimit) {
				writeValidation(w, "tokens", "limit exceeded")
				return
			}
			if errors.Is(err, errAPITokenCorrupt) {
				writeError(w, http.StatusInternalServerError, "State file is corrupted")
				return
			}
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := s.appendTokenAccessLog(r, "token_create", record.ID, "success", ""); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		actorType, actorID := apiAuditActor(r)
		if err := s.appendSecurityAudit(r, "token_create", actorType, actorID, "api_token", record.ID, "success"); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"id":           record.ID,
			"label":        record.Label,
			"scopes":       record.Scopes,
			"created_at":   record.CreatedAt,
			"expires_at":   record.ExpiresAt,
			"revoked_at":   record.RevokedAt,
			"last_used_at": record.LastUsedAt,
			"token":        token,
		})
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
	var revoked bool
	_, err := updateAPITokensLocked(path, func(tokens []apiTokenRecord) ([]apiTokenRecord, bool, error) {
		for i := range tokens {
			if tokens[i].ID == id && tokens[i].RevokedAt == nil {
				now := s.nowString()
				tokens[i].RevokedAt = &now
				revoked = true
				return tokens, true, nil
			}
		}
		return tokens, false, nil
	})
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			writeError(w, http.StatusConflict, "Conflict")
			return
		}
		if errors.Is(err, errAPITokenCorrupt) {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if !revoked {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if err := s.appendTokenAccessLog(r, "token_revoke", id, "success", ""); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	actorType, actorID := apiAuditActor(r)
	if err := s.appendSecurityAudit(r, "token_revoke", actorType, actorID, "api_token", id, "success"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Token revoked"})
}

func (s *APIServer) handleRuleFile(filename, logType string, createStatus int) http.HandlerFunc {
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
			id, err := s.nextTimeID(ruleIDPrefix(logType), ruleIDSet(rules))
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			rule["id"] = id
			if logType == "hook" {
				if _, ok := rule["enabled"]; !ok {
					rule["enabled"] = true
				}
				if _, ok := rule["timeout_seconds"]; !ok {
					rule["timeout_seconds"] = 300
				}
			}
			rules = append(rules, rule)
			if err := atomicWriteJSON(path, ruleListWrapper(path, rules), 0600); err != nil {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: logType, Changes: map[string]any{"id": rule["id"]}})
			writeJSON(w, createStatus, rule)
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
		if err := atomicWriteJSON(path, ruleListWrapper(path, next), 0600); err != nil {
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
	target, err := s.selectedOutputTarget()
	if err != nil {
		writeOutputTargetError(w, err)
		return
	}
	histories := outputMetaSuccessHistory(readHistory(filepath.Join(s.cfg.StateDir, ".build_history")), target)
	if len(histories) == 0 {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	expected := stringPtrValue(histories[0].OutputSHA256)
	if expected == "" {
		if log, err := readBuildLogByID(s.cfg.StateDir, histories[0].ID); err == nil {
			expected = firstNonEmpty(log.OutputSHA256, log.SHA256)
		}
	}
	if expected == "" {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	meta, err := inspectOutput(target.Out)
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
	path := filepath.Join(s.cfg.StateDir, ".pipeline_config")
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
			value = defaultPipelineConfig()
		}
		writeJSON(w, http.StatusOK, value)
	case http.MethodPost:
		var value map[string]any
		if !decodeBody(w, r, &value, true) {
			return
		}
		if !validatePipelineConfig(w, value) {
			return
		}
		current, _, err := readOptionalJSONMap(path)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if current == nil {
			current = defaultPipelineConfig()
		}
		if mapsEqual(current, value) {
			writeJSON(w, http.StatusOK, map[string]string{"message": "No changes"})
			return
		}
		if err := atomicWriteJSON(path, value, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		_ = appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "pipeline_config", Changes: maskSecrets(value)})
		writeJSON(w, http.StatusOK, map[string]string{"message": "Pipeline config updated"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *APIServer) handleBuildChainConfig(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.cfg.StateDir, ".build_chain_config")
	switch r.Method {
	case http.MethodGet:
		if !rejectBody(w, r) {
			return
		}
		var cfg apiBuildChainConfig
		if err := readJSONIfExists(path, &cfg); err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if cfg.Chains == nil {
			cfg.Chains = []apiBuildChainJob{}
		}
		writeJSON(w, http.StatusOK, cfg)
	case http.MethodPost:
		var cfg apiBuildChainConfig
		if !decodeBody(w, r, &cfg, true) {
			return
		}
		if cfg.Chains == nil {
			cfg.Chains = []apiBuildChainJob{}
		}
		if !validateBuildChainConfig(w, cfg) {
			return
		}
		var current apiBuildChainConfig
		if err := readJSONIfExists(path, &current); err != nil {
			writeError(w, http.StatusInternalServerError, "State file is corrupted")
			return
		}
		if current.Chains == nil {
			current.Chains = []apiBuildChainJob{}
		}
		if buildChainConfigsEqual(current, cfg) {
			writeJSON(w, http.StatusOK, map[string]any{"message": "No changes", "chains_count": len(cfg.Chains)})
			return
		}
		if err := atomicWriteJSON(path, cfg, 0600); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if err := appendJSONLine(filepath.Join(s.cfg.StateDir, ".config_log"), apiConfigLogRecord{At: s.nowString(), Type: "build_chain_config", Changes: map[string]any{"chains_count": len(cfg.Chains)}}); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"message": "Build chain config updated", "chains_count": len(cfg.Chains)})
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
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
		authTransactionPath := apiAuthTransactionPath(r.URL.Path)
		releaseAuthTransaction, ok := s.acquireAuthTransaction(w, r)
		if !ok {
			return
		}
		authTransactionReleased := false
		releaseAuth := func() {
			if authTransactionReleased {
				return
			}
			releaseAuthTransaction()
			authTransactionReleased = true
		}
		defer releaseAuth()

		token := bearerToken(r)
		auth, status := s.authenticateRequest(r, token)
		if status != http.StatusOK {
			if status == http.StatusForbidden {
				writeError(w, http.StatusForbidden, "Forbidden")
				return
			}
			if status == http.StatusConflict {
				writeError(w, http.StatusConflict, "Conflict")
				return
			}
			if status == http.StatusInternalServerError {
				writeError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if auth.AuthType == "session" {
			if auth.PasswordChangeRequired && r.URL.Path != "/api/change-password" && r.URL.Path != "/api/logout" {
				writeError(w, http.StatusForbidden, "Password change required")
				return
			}
			s.touchSession(token)
		}
		if !authTransactionPath {
			releaseAuth()
		}
		if capture, ok := r.Context().Value(apiAuthCaptureContextKey{}).(*apiAuthInfo); ok {
			*capture = auth
		}
		nextReq := r.WithContext(context.WithValue(r.Context(), apiAuthContextKey{}, auth))
		group, _ := apiRateLimitGroupForRequest(r.Method, r.URL.Path)
		if !s.enforceAPIRateLimit(w, nextReq, auth, group) {
			return
		}
		if !validateRouteQuery(w, nextReq) {
			return
		}
		next(w, nextReq)
	}
}

func (s *APIServer) acquireAuthTransaction(w http.ResponseWriter, r *http.Request) (func(), bool) {
	timer := time.NewTimer(s.cfg.AuthTransactionTimeout)
	defer timer.Stop()
	select {
	case s.authTransactions <- struct{}{}:
		return func() {
			<-s.authTransactions
		}, true
	case <-timer.C:
		writeError(w, http.StatusConflict, "Conflict")
		return nil, false
	case <-r.Context().Done():
		writeError(w, http.StatusConflict, "Conflict")
		return nil, false
	}
}

func apiAuthTransactionPath(path string) bool {
	switch path {
	case "/api/login",
		"/api/login/totp",
		"/api/logout",
		"/api/change-password",
		"/api/sessions",
		"/api/sessions/revoke-all",
		"/api/auth/totp-status",
		"/api/auth/totp-setup",
		"/api/auth/totp-confirm",
		"/api/auth/totp":
		return true
	default:
		return false
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
	_, ok := s.authenticateToken(token)
	return ok
}

func (s *APIServer) authenticateRequest(r *http.Request, token string) (apiAuthInfo, int) {
	if token == "" {
		return apiAuthInfo{}, http.StatusUnauthorized
	}
	hash := tokenHash(token)
	if auth, ok := s.authenticateSession(hash); ok {
		return auth, http.StatusOK
	}
	return s.validAPITokenForRequest(hash, r)
}

func (s *APIServer) authenticateToken(token string) (apiAuthInfo, bool) {
	if token == "" {
		return apiAuthInfo{}, false
	}
	hash := tokenHash(token)
	if auth, ok := s.authenticateSession(hash); ok {
		return auth, true
	}
	auth, status := s.validAPITokenForReadOnlyAuth(hash)
	return auth, status == http.StatusOK
}

func (s *APIServer) authenticateSession(hash string) (apiAuthInfo, bool) {
	s.mu.Lock()
	session, ok := s.sessions[hash]
	if !ok {
		s.mu.Unlock()
		return apiAuthInfo{}, false
	}
	expires, err := time.Parse(apiTimeLayout, session.ExpiresAt)
	if err != nil || !s.cfg.Now().UTC().Before(expires) {
		delete(s.sessions, hash)
		s.mu.Unlock()
		return apiAuthInfo{}, false
	}
	s.mu.Unlock()
	actor := "admin"
	return apiAuthInfo{AuthType: "session", Actor: &actor, PasswordChangeRequired: session.PasswordChangeRequired}, true
}

func (s *APIServer) validAPIToken(token string) bool {
	if token == "" {
		return false
	}
	_, status := s.validAPITokenForReadOnlyAuth(tokenHash(token))
	return status == http.StatusOK
}

func (s *APIServer) validAPITokenForReadOnlyAuth(hash string) (apiAuthInfo, int) {
	tokens, err := readAPITokens(filepath.Join(s.cfg.StateDir, ".api_tokens"))
	if err != nil {
		return apiAuthInfo{}, http.StatusInternalServerError
	}
	for _, record := range tokens {
		if !tokenHashEqual(record.TokenHash, hash) {
			continue
		}
		if record.RevokedAt != nil || apiTokenExpired(record, s.cfg.Now().UTC()) {
			return apiAuthInfo{}, http.StatusUnauthorized
		}
		actor := record.ID
		return apiAuthInfo{AuthType: "api_token", Actor: &actor, Scopes: append([]string{}, record.Scopes...)}, http.StatusOK
	}
	return apiAuthInfo{}, http.StatusUnauthorized
}

func (s *APIServer) validAPITokenForRequest(hash string, r *http.Request) (apiAuthInfo, int) {
	method := strings.ToUpper(r.Method)
	path := r.URL.Path
	tokensPath := filepath.Join(s.cfg.StateDir, ".api_tokens")
	var matched *apiTokenRecord
	var failureAction string
	var failureReason string
	var failureStatus int
	_, err := updateAPITokensLocked(tokensPath, func(tokens []apiTokenRecord) ([]apiTokenRecord, bool, error) {
		for i := range tokens {
			if !tokenHashEqual(tokens[i].TokenHash, hash) {
				continue
			}
			record := tokens[i]
			matched = &record
			if tokens[i].RevokedAt != nil {
				failureAction = "token_revoked_reject"
				failureReason = "revoked"
				failureStatus = http.StatusUnauthorized
				return tokens, false, nil
			}
			if apiTokenExpired(tokens[i], s.cfg.Now().UTC()) {
				failureAction = "token_expired"
				failureReason = "expired"
				failureStatus = http.StatusUnauthorized
				return tokens, false, nil
			}
			if !apiTokenScopesAllow(tokens[i].Scopes, method, path) {
				failureAction = "permission_denied"
				failureReason = "scope_denied"
				failureStatus = http.StatusForbidden
				return tokens, false, nil
			}
			now := s.nowString()
			tokens[i].LastUsedAt = &now
			record = tokens[i]
			matched = &record
			return tokens, true, nil
		}
		return tokens, false, nil
	})
	if err != nil {
		if errors.Is(err, errAPIStateLockConflict) {
			return apiAuthInfo{}, http.StatusConflict
		}
		return apiAuthInfo{}, http.StatusInternalServerError
	}
	if matched == nil {
		return apiAuthInfo{}, http.StatusUnauthorized
	}
	if failureAction != "" {
		if err := s.appendTokenAccessLog(r, failureAction, matched.ID, "failure", failureReason); err != nil {
			return apiAuthInfo{}, http.StatusInternalServerError
		}
		result := "failure"
		targetType := "api_token"
		targetID := any(matched.ID)
		if failureAction == "permission_denied" {
			result = "denied"
			targetType = "endpoint"
			targetID = method + " " + path
		}
		if err := s.appendSecurityAudit(r, failureAction, "api_token", matched.ID, targetType, targetID, result); err != nil {
			return apiAuthInfo{}, http.StatusInternalServerError
		}
		return apiAuthInfo{}, failureStatus
	}
	if err := s.appendTokenAccessLog(r, "token_auth", matched.ID, "success", ""); err != nil {
		return apiAuthInfo{}, http.StatusInternalServerError
	}
	if err := s.appendSecurityAudit(r, "token_auth", "api_token", matched.ID, "api_token", matched.ID, "success"); err != nil {
		return apiAuthInfo{}, http.StatusInternalServerError
	}
	actor := matched.ID
	return apiAuthInfo{AuthType: "api_token", Actor: &actor, Scopes: append([]string{}, matched.Scopes...)}, http.StatusOK
}

func (s *APIServer) appendTokenAccessLog(r *http.Request, action, tokenID, result, reason string) error {
	var reasonValue any
	if reason != "" {
		reasonValue = reason
	}
	return appendJSONLine(filepath.Join(s.cfg.StateDir, ".access_log"), map[string]any{
		"at":          s.nowString(),
		"action":      action,
		"result":      result,
		"session_id":  nil,
		"token_id":    tokenID,
		"remote_addr": apiRemoteAddr(r),
		"reason":      reasonValue,
	})
}

func (s *APIServer) appendAuthAccessLog(r *http.Request, action, result, reason string) error {
	var reasonValue any
	if reason != "" {
		reasonValue = reason
	}
	return appendJSONLine(filepath.Join(s.cfg.StateDir, ".access_log"), map[string]any{
		"at":          s.nowString(),
		"action":      action,
		"result":      result,
		"session_id":  nil,
		"token_id":    nil,
		"remote_addr": apiRemoteAddr(r),
		"reason":      reasonValue,
	})
}

func (s *APIServer) writeTOTPLoginFailure(w http.ResponseWriter, r *http.Request) {
	if err := s.appendAuthAccessLog(r, "totp_failure", "failure", "invalid_totp"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if err := s.appendSecurityAudit(r, "totp_failure", "anonymous", nil, "auth", "login_totp", "failure"); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeError(w, http.StatusUnauthorized, "Unauthorized")
}

func (s *APIServer) appendSecurityAudit(r *http.Request, action, actorType string, actorID any, targetType string, targetID any, result string) error {
	return appendJSONLine(filepath.Join(s.cfg.StateDir, ".audit_log"), map[string]any{
		"timestamp":   s.nowString(),
		"request_id":  apiRequestID(r),
		"action":      action,
		"actor_type":  actorType,
		"actor_id":    actorID,
		"target_type": targetType,
		"target_id":   targetID,
		"result":      result,
		"remote_addr": apiRemoteAddr(r),
		"message":     nil,
	})
}

func apiAuditActor(r *http.Request) (string, any) {
	auth, ok := r.Context().Value(apiAuthContextKey{}).(apiAuthInfo)
	if !ok {
		return "anonymous", nil
	}
	switch auth.AuthType {
	case "session":
		return "admin", "admin"
	case "api_token":
		if auth.Actor != nil {
			return "api_token", *auth.Actor
		}
		return "api_token", nil
	default:
		return "anonymous", nil
	}
}

func apiRequestID(r *http.Request) any {
	requestID, ok := r.Context().Value(apiRequestIDContextKey{}).(string)
	if !ok || requestID == "" {
		return nil
	}
	return requestID
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

func (s *APIServer) enqueueBuildRequest(statePath, trigger, requestedBy string, payload map[string]any) (string, bool, error) {
	state, err := readBuildState(statePath)
	if err != nil {
		return "", false, err
	}
	if trigger == "webhook" {
		id, duplicate, conflict := webhookQueueMatch(state, payload)
		if duplicate {
			return id, true, nil
		}
		if conflict {
			return "", false, errAPIQueueConflict
		}
	} else if queueEntryMatches(state.ActiveQueueEntry, trigger, payload) {
		return fmt.Sprint(state.ActiveQueueEntry["id"]), true, nil
	}
	for _, entry := range sortedQueueEntries(state.Queued) {
		if queueEntryMatches(entry, trigger, payload) {
			return fmt.Sprint(entry["id"]), true, nil
		}
	}
	limit := s.queueMaxSize()
	active := len(state.ActiveQueueEntry) > 0
	if !state.Running && !active && !lockExists(filepath.Join(s.cfg.StateDir, ".build_lock")) && state.CurrentBuildID == nil && limit < 1 {
		limit = 1
	}
	if limit == 0 || len(state.Queued) >= limit {
		return "", false, errAPIQueueFull
	}
	id, err := s.nextTimeID("q", queueIDSet(state.ActiveQueueEntry, state.Queued))
	if err != nil {
		return "", false, err
	}
	state.Queued = append(state.Queued, map[string]any{
		"id":           id,
		"trigger":      trigger,
		"queued_at":    s.nowString(),
		"requested_by": requestedBy,
		"priority":     "normal",
		"created_seq":  nextCreatedSeq(state.Queued),
		"payload":      payload,
	})
	return id, false, atomicWriteJSON(statePath, state, 0600)
}

func (s *APIServer) enqueueApprovalRequest(statePath, requestedBy string, payload map[string]any) (string, bool, error) {
	state, err := readBuildState(statePath)
	if err != nil {
		return "", false, err
	}
	matches := []map[string]any{}
	if fmt.Sprint(state.ActiveQueueEntry["trigger"]) == "approval" && queuePayloadApprovalID(state.ActiveQueueEntry["payload"]) == fmt.Sprint(payload["approval_id"]) {
		matches = append(matches, state.ActiveQueueEntry)
	}
	for _, entry := range state.Queued {
		if fmt.Sprint(entry["trigger"]) != "approval" {
			continue
		}
		if queuePayloadApprovalID(entry["payload"]) == fmt.Sprint(payload["approval_id"]) {
			matches = append(matches, entry)
		}
	}
	if len(matches) == 1 {
		if !queuePayloadEqual(matches[0]["payload"], payload) {
			log.Printf("APPROVAL_QUEUE_INCONSISTENT: approval_id=%s", payload["approval_id"])
			return "", false, errors.New("approval queue inconsistent")
		}
		return fmt.Sprint(matches[0]["id"]), true, nil
	}
	if len(matches) > 1 {
		log.Printf("APPROVAL_QUEUE_INCONSISTENT: approval_id=%s", payload["approval_id"])
		return "", false, errors.New("approval queue inconsistent")
	}
	limit := s.queueMaxSize()
	active := len(state.ActiveQueueEntry) > 0
	if !state.Running && !active && !lockExists(filepath.Join(s.cfg.StateDir, ".build_lock")) && state.CurrentBuildID == nil && limit < 1 {
		limit = 1
	}
	if limit == 0 || len(state.Queued) >= limit {
		return "", false, errAPIQueueFull
	}
	id, err := s.nextTimeID("q", queueIDSet(state.ActiveQueueEntry, state.Queued))
	if err != nil {
		return "", false, err
	}
	state.Queued = append(state.Queued, map[string]any{
		"id":           id,
		"trigger":      "approval",
		"queued_at":    s.nowString(),
		"requested_by": requestedBy,
		"priority":     "normal",
		"created_seq":  nextCreatedSeq(state.Queued),
		"payload":      payload,
	})
	return id, false, atomicWriteJSON(statePath, state, 0600)
}

func (s *APIServer) dispatchRunner(queueID string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	commandRunner := s.cfg.CommandRunner
	if commandRunner == nil {
		commandRunner = runAPICommand
	}
	result, err := commandRunner(ctx, "systemctl", "start", "--no-block", "adlaire-ci.service")
	if err == nil && result.ExitCode == 0 {
		return "requested"
	}
	log.Printf("RUNNER_ACTIVATION_DEFERRED: queue_id=%s", queueID)
	return "timer_fallback"
}

func runAPICommand(ctx context.Context, name string, args ...string) (apiCommandResult, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = io.Discard
	err := cmd.Run()
	result := apiCommandResult{Stdout: stdout.String(), ExitCode: 0}
	if err != nil {
		result.ExitCode = -1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
		}
	}
	return result, err
}

func queuePayloadApprovalID(value any) string {
	if payload, ok := value.(map[string]any); ok {
		return fmt.Sprint(payload["approval_id"])
	}
	return ""
}

func queuePayloadEqual(a any, b map[string]any) bool {
	aj, errA := json.Marshal(a)
	bj, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aj) == string(bj)
}

func queueEntryMatches(entry map[string]any, trigger string, payload map[string]any) bool {
	return len(entry) > 0 && fmt.Sprint(entry["trigger"]) == trigger && queuePayloadEqual(entry["payload"], payload)
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

func webhookQueueMatch(state apiBuildState, payload map[string]any) (string, bool, bool) {
	deliveryID := fmt.Sprint(payload["delivery_id"])
	branch := fmt.Sprint(payload["branch"])
	sha := fmt.Sprint(payload["sha"])
	activeID, activeDuplicate, activeConflict := webhookEntryMatch(state.ActiveQueueEntry, deliveryID, branch, sha)
	if activeDuplicate || activeConflict {
		if activeDuplicate {
			for _, entry := range sortedQueueEntries(state.Queued) {
				if waitingID, waitingDuplicate, _ := webhookEntryMatch(entry, deliveryID, branch, sha); waitingDuplicate {
					log.Printf("QUEUE_DUPLICATE_STATE: queue_id=%s", waitingID)
					break
				}
			}
		}
		return activeID, activeDuplicate, activeConflict
	}
	for _, entry := range sortedQueueEntries(state.Queued) {
		id, duplicate, conflict := webhookEntryMatch(entry, deliveryID, branch, sha)
		if duplicate || conflict {
			return id, duplicate, conflict
		}
	}
	return "", false, false
}

func webhookEntryMatch(entry map[string]any, deliveryID, branch, sha string) (string, bool, bool) {
	if len(entry) == 0 || fmt.Sprint(entry["trigger"]) != "webhook" {
		return "", false, false
	}
	payload, ok := entry["payload"].(map[string]any)
	if !ok || fmt.Sprint(payload["delivery_id"]) != deliveryID {
		return "", false, false
	}
	id := fmt.Sprint(entry["id"])
	if fmt.Sprint(payload["branch"]) == branch && fmt.Sprint(payload["sha"]) == sha {
		return id, true, false
	}
	return id, false, true
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
	for key := range patch {
		if !serverConfigStorageKeyAllowed(key) {
			return cfg, fmt.Errorf("unknown server config key")
		}
	}
	defaultData, err := json.Marshal(cfg)
	if err != nil {
		return cfg, err
	}
	var merged map[string]json.RawMessage
	if err := json.Unmarshal(defaultData, &merged); err != nil {
		return cfg, err
	}
	for key, value := range patch {
		merged[key] = value
	}
	mergedData, err := json.Marshal(merged)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(mergedData, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
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

func (s *APIServer) nextTimeID(prefix string, existing map[string]bool) (string, error) {
	base := s.newID(prefix)
	if !existing[base] {
		return base, nil
	}
	for i := 1; i <= 999; i++ {
		id := fmt.Sprintf("%s-%03d", base, i)
		if !existing[id] {
			return id, nil
		}
	}
	return "", errors.New("id suffix exhausted")
}

func queueIDSet(active map[string]any, entries []map[string]any) map[string]bool {
	ids := map[string]bool{}
	if id := fmt.Sprint(active["id"]); id != "" {
		ids[id] = true
	}
	for _, entry := range entries {
		id := fmt.Sprint(entry["id"])
		if id != "" {
			ids[id] = true
		}
	}
	return ids
}

func ruleIDSet(rules []map[string]any) map[string]bool {
	ids := map[string]bool{}
	for _, rule := range rules {
		id := fmt.Sprint(rule["id"])
		if id != "" {
			ids[id] = true
		}
	}
	return ids
}

func ruleIDPrefix(logType string) string {
	switch logType {
	case "hook":
		return "h"
	case "alert_rule":
		return "r"
	case "tag_rule":
		return "t"
	default:
		return "r"
	}
}

func nextTokenID(tokens []apiTokenRecord) (string, error) {
	maxID := 0
	for _, token := range tokens {
		if !strings.HasPrefix(token.ID, "tok") || len(token.ID) != 9 {
			continue
		}
		n, err := strconv.Atoi(token.ID[3:])
		if err != nil {
			continue
		}
		if n > maxID {
			maxID = n
		}
	}
	if maxID >= 999999 {
		return "", errors.New("token id exhausted")
	}
	return fmt.Sprintf("tok%06d", maxID+1), nil
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
	mtimeSource := info.ModTime()
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
	manifest := sha256.New()
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return outputInspection{}, err
		}
		total += info.Size()
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
	if !mtimeSource.IsZero() {
		mtime = mtimeSource.UTC().Format(apiTimeLayout)
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

func (s *APIServer) selectBuildStreamLog() (apiBuildLog, bool, bool, error) {
	state, err := readBuildState(filepath.Join(s.cfg.StateDir, ".build_state"))
	if err != nil {
		return apiBuildLog{}, false, false, err
	}
	if state.Running && state.CurrentBuildID != nil && *state.CurrentBuildID != "" {
		logRecord, err := readBuildLog(filepath.Join(s.cfg.StateDir, ".build_logs", *state.CurrentBuildID+".json"))
		if errors.Is(err, os.ErrNotExist) {
			return apiBuildLog{}, true, false, nil
		}
		if err != nil {
			return apiBuildLog{}, true, false, err
		}
		if logRecord.ID == "" {
			logRecord.ID = *state.CurrentBuildID
		}
		return logRecord, true, true, nil
	}
	logs := s.readBuildStreamLogsNewest()
	if len(logs) == 0 {
		return apiBuildLog{}, false, false, nil
	}
	return logs[0], false, true, nil
}

func (s *APIServer) readBuildStreamLogsNewest() []apiBuildLog {
	dir := filepath.Join(s.cfg.StateDir, ".build_logs")
	logs := []apiBuildLog{}
	normalIDs := map[string]bool{}
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}
			logRecord, err := readBuildLog(filepath.Join(dir, entry.Name()))
			if err != nil || logRecord.FinishedAt == "" {
				if err != nil {
					log.Printf("BUILD_STREAM_SKIP_CORRUPT: build_id=%s", strings.TrimSuffix(entry.Name(), ".json"))
				}
				continue
			}
			if logRecord.ID == "" {
				logRecord.ID = strings.TrimSuffix(entry.Name(), ".json")
			}
			normalIDs[logRecord.ID] = true
			logs = append(logs, logRecord)
		}
	}
	archiveDir := filepath.Join(dir, "archive")
	if entries, err := os.ReadDir(archiveDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json.gz") {
				continue
			}
			id := strings.TrimSuffix(entry.Name(), ".json.gz")
			if normalIDs[id] {
				continue
			}
			logRecord, err := readArchivedBuildLog(filepath.Join(archiveDir, entry.Name()))
			if err != nil || logRecord.FinishedAt == "" {
				log.Printf("BUILD_STREAM_SKIP_CORRUPT: build_id=%s", id)
				continue
			}
			if logRecord.ID == "" {
				logRecord.ID = id
			}
			logs = append(logs, logRecord)
		}
	}
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].FinishedAt == logs[j].FinishedAt {
			return logs[i].ID > logs[j].ID
		}
		return logs[i].FinishedAt > logs[j].FinishedAt
	})
	return logs
}

func readArchivedBuildLog(path string) (apiBuildLog, error) {
	file, err := os.Open(path)
	if err != nil {
		return apiBuildLog{}, err
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		return apiBuildLog{}, err
	}
	defer reader.Close()
	var logRecord apiBuildLog
	if err := json.NewDecoder(reader).Decode(&logRecord); err != nil {
		return apiBuildLog{}, err
	}
	return logRecord, nil
}

func buildStreamLines(logRecord apiBuildLog) []string {
	lines := []string{}
	appendTextLines := func(value string) {
		if value == "" {
			return
		}
		parts := strings.Split(value, "\n")
		if len(parts) > 0 && parts[len(parts)-1] == "" {
			parts = parts[:len(parts)-1]
		}
		for _, line := range parts {
			lines = append(lines, strings.TrimRight(line, "\r"))
		}
	}
	appendTextLines(logRecord.Pipeline.Stdout)
	appendTextLines(logRecord.Pipeline.Stderr)
	lines = append(lines, logRecord.Warnings...)
	if logRecord.Error != nil && *logRecord.Error != "" {
		lines = append(lines, *logRecord.Error)
	}
	return lines
}

func buildStreamStatus(logRecord apiBuildLog, running bool) string {
	if running {
		return "running"
	}
	switch logRecord.Status {
	case "success", "failure", "cancelled":
		return logRecord.Status
	}
	switch {
	case logRecord.TargetStatus == "success" || logRecord.TargetStatus == "success_deploy_pending":
		return "success"
	case logRecord.TargetStatus == "cancelled":
		return "cancelled"
	default:
		return "failure"
	}
}

func writeSSEFrame(w http.ResponseWriter, flusher http.Flusher, value any) bool {
	payload, err := json.Marshal(value)
	if err != nil {
		return false
	}
	frame := fmt.Sprintf("data: %s\n\n", payload)
	n, err := io.WriteString(w, frame)
	if err != nil || n != len(frame) {
		return false
	}
	if flusher != nil {
		flusher.Flush()
	}
	return true
}

func truncateStreamRunes(value string, limit int) string {
	if limit < 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}

func validateCredentials(path string) error {
	_, err := readCredentials(path)
	return err
}

func acquireAPIStateFileLock(path string) (func(), error) {
	release, err := acquireStateFileLock(path)
	if err != nil {
		return func() {}, fmt.Errorf("%w: %v", errAPIStateLockConflict, err)
	}
	return release, nil
}

func readCredentialsLocked(path string) (apiCredentials, error) {
	release, err := acquireAPIStateFileLock(path)
	if err != nil {
		return apiCredentials{}, err
	}
	defer release()
	return readCredentials(path)
}

func updateCredentialsLocked(path string, mutate func(apiCredentials) (apiCredentials, bool, error)) (apiCredentials, error) {
	release, err := acquireAPIStateFileLock(path)
	if err != nil {
		return apiCredentials{}, err
	}
	defer release()
	current, err := readCredentials(path)
	if err != nil {
		return apiCredentials{}, err
	}
	next, write, err := mutate(current)
	if err != nil {
		return apiCredentials{}, err
	}
	if !write {
		return next, nil
	}
	if err := validateCredentialsValue(next); err != nil {
		return apiCredentials{}, err
	}
	if err := atomicWriteJSONLocked(path, next, 0600); err != nil {
		return apiCredentials{}, err
	}
	return next, nil
}

func readCredentials(path string) (apiCredentials, error) {
	var cred apiCredentials
	if err := readStrictJSONFile(path, &cred); err != nil {
		return cred, err
	}
	return cred, validateCredentialsValue(cred)
}

func validateCredentialsValue(cred apiCredentials) error {
	if err := security.ValidatePasswordRecord(credentialsPasswordRecord(cred)); err != nil {
		return errors.New("invalid credentials")
	}
	if cred.LoginCount < 0 || !validAPITime(cred.UpdatedAt) {
		return errors.New("invalid credentials")
	}
	if cred.LastLoginAt != nil && !validAPITime(*cred.LastLoginAt) {
		return errors.New("invalid credentials")
	}
	return nil
}

func credentialsPasswordRecord(cred apiCredentials) security.PasswordRecord {
	return security.PasswordRecord{
		PasswordHash: cred.PasswordHash,
		Salt:         cred.Salt,
		Algorithm:    cred.Algorithm,
		Iterations:   cred.Iterations,
	}
}

func credentialsPasswordEqual(cred apiCredentials, password string) bool {
	ok, err := security.VerifyPassword(credentialsPasswordRecord(cred), password)
	return err == nil && ok
}

func passwordHashEqual(left, right string) bool {
	return hmac.Equal([]byte(left), []byte(right))
}

func advanceCredentialsLogin(cred *apiCredentials, now string) string {
	if cred.LoginCount < maxLoginCount {
		cred.LoginCount++
	}
	cred.LastLoginAt = &now
	return credentialsMustChange(cred)
}

func credentialsMustChange(cred *apiCredentials) string {
	if !cred.MustChange {
		return "none"
	}
	if cred.LoginCount >= 5 {
		return "forced"
	}
	return "prompt"
}

func credentialsFingerprint(cred apiCredentials) string {
	payload := strings.Join([]string{
		cred.PasswordHash,
		cred.Salt,
		cred.Algorithm,
		strconv.Itoa(cred.Iterations),
		cred.UpdatedAt,
	}, "\n")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func tokenHashEqual(left, right string) bool {
	return hmac.Equal([]byte(left), []byte(right))
}

func newAPIToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}
	return "act_" + base64.RawURLEncoding.EncodeToString(buf), nil
}

func randomHex(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func newAPIRequestID() (string, error) {
	return randomHex(16)
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

func readTOTPSecretLocked(path string) (apiTOTPSecret, error) {
	release, err := acquireAPIStateFileLock(path)
	if err != nil {
		return apiTOTPSecret{}, err
	}
	defer release()
	return readTOTPSecret(path)
}

func updateTOTPSecretLocked(path string, mutate func(apiTOTPSecret) (apiTOTPSecret, bool, error)) (apiTOTPSecret, error) {
	release, err := acquireAPIStateFileLock(path)
	if err != nil {
		return apiTOTPSecret{}, err
	}
	defer release()
	current, err := readTOTPSecret(path)
	if err != nil {
		return apiTOTPSecret{}, err
	}
	next, write, err := mutate(current)
	if err != nil {
		return apiTOTPSecret{}, err
	}
	if !write {
		return next, nil
	}
	if err := validateTOTPSecretValue(next); err != nil {
		return apiTOTPSecret{}, err
	}
	if err := atomicWriteJSONLocked(path, next, 0600); err != nil {
		return apiTOTPSecret{}, err
	}
	return next, nil
}

func validateTOTPSecretValue(secret apiTOTPSecret) error {
	if !secret.Enabled {
		if secret.SecretBase32 != nil || secret.ConfirmedAt != nil || secret.LastAcceptedStep != nil {
			return errors.New("invalid totp secret")
		}
		return nil
	}
	if secret.SecretBase32 == nil || strings.TrimSpace(*secret.SecretBase32) == "" {
		return errors.New("invalid totp secret")
	}
	if _, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(strings.TrimSpace(*secret.SecretBase32))); err != nil {
		return errors.New("invalid totp secret")
	}
	if secret.ConfirmedAt != nil && !validAPITime(*secret.ConfirmedAt) {
		return errors.New("invalid totp secret")
	}
	if secret.LastAcceptedStep != nil && *secret.LastAcceptedStep < 0 {
		return errors.New("invalid totp secret")
	}
	return nil
}

func totpStatusPayload(secret apiTOTPSecret) map[string]any {
	return map[string]any{"enabled": secret.Enabled, "confirmed_at": secret.ConfirmedAt}
}

func verifyTOTPCode(secretBase32, code string, at time.Time) bool {
	_, ok := verifyTOTPCodeStep(secretBase32, code, at)
	return ok
}

func verifyTOTPCodeStep(secretBase32, code string, at time.Time) (int64, bool) {
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return 0, false
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	step := at.UTC().Unix() / 30
	for _, offset := range []int64{-1, 0, 1} {
		if totpCode(secretBase32, step+offset) == code {
			return step + offset, true
		}
	}
	return 0, false
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

func readBuildStatusForHealth(path string) (apiBuildStatus, bool, string) {
	var status apiBuildStatus
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return status, false, "build_status_missing"
	}
	if err != nil {
		return status, false, "build_status_read_error"
	}
	if err := json.Unmarshal(data, &status); err != nil {
		return status, false, "build_status_corrupt"
	}
	return status, true, ""
}

func readHistory(path string) []apiHistoryRecord {
	data, err := os.ReadFile(path)
	if err != nil {
		return []apiHistoryRecord{}
	}
	records := []apiHistoryRecord{}
	seen := map[string]bool{}
	for lineNo, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record apiHistoryRecord
		if json.Unmarshal([]byte(line), &record) == nil && record.ID != "" && record.Status != "" {
			if seen[record.ID] {
				fmt.Fprintf(os.Stderr, "BUILD_HISTORY_DUPLICATE_ID: path=%s line=%d id=%s\n", path, lineNo+1, record.ID)
				continue
			}
			seen[record.ID] = true
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

func readHealthHistory(path string) ([]apiHistoryRecord, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []apiHistoryRecord{}, nil
	}
	if err != nil {
		return nil, err
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
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return firstNonEmpty(records[i].FinishedAt, records[i].StartedAt, records[i].BuildAt, records[i].ID) > firstNonEmpty(records[j].FinishedAt, records[j].StartedAt, records[j].BuildAt, records[j].ID)
	})
	return records, nil
}

func latestStatusSummaryHistory(records []apiHistoryRecord) *apiHistoryRecord {
	for i := range records {
		if statusSummaryEligible(records[i].Status) {
			return &records[i]
		}
	}
	return nil
}

func statusSummaryEligible(status string) bool {
	switch status {
	case "success", "success_deploy_pending", "failure", "failure_api", "failure_decode", "failure_precheck",
		"failure_build", "failure_state_write", "failure_target_missing", "failure_pipeline_config",
		"failure_timeout", "failure_remote_build", "failure_tag_rule", "hook_error", "cancelled",
		"skipped_no_change", "skipped_cooldown", "skipped_schedule_paused", "skipped_allowed_hours",
		"skipped_tag_filter", "skipped_maintenance", "circuit_open", "config_recovered", "config_error",
		"lock_skipped":
		return true
	default:
		return false
	}
}

func filterHistory(w http.ResponseWriter, r *http.Request, records []apiHistoryRecord) []apiHistoryRecord {
	trigger := strings.TrimSpace(r.URL.Query().Get("trigger"))
	tag := strings.TrimSpace(r.URL.Query().Get("tag"))
	flaggedRaw := strings.TrimSpace(r.URL.Query().Get("flagged"))
	failureCategory := strings.TrimSpace(r.URL.Query().Get("failure_category"))
	if trigger != "" && !validHistoryTrigger(trigger) {
		writeValidation(w, "trigger", "invalid value")
		return nil
	}
	if failureCategory != "" && !validFailureCategory(failureCategory) {
		writeValidation(w, "failure_category", "invalid value")
		return nil
	}
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
		if failureCategory != "" && stringPtrValue(record.FailureCategory) != failureCategory {
			continue
		}
		out = append(out, record)
	}
	return out
}

func validateHistoryQuery(w http.ResponseWriter, r *http.Request) bool {
	allowed := map[string]bool{"page": true, "per_page": true, "trigger": true, "tag": true, "flagged": true, "failure_category": true}
	for key := range r.URL.Query() {
		if !allowed[key] {
			writeValidation(w, key, "unknown query")
			return false
		}
	}
	return true
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

func validHistoryTrigger(value string) bool {
	switch value {
	case "polling", "manual", "webhook", "approval", "rollback", "retry":
		return true
	default:
		return false
	}
}

func validFailureCategory(value string) bool {
	switch value {
	case "github_api", "pipeline_timeout", "pipeline_exit", "deploy_failure", "hook_error", "config_error", "resource_error", "unknown":
		return true
	default:
		return false
	}
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
		HistoryRetention: json.RawMessage(`{"enabled":false,"max_count":1000,"max_age_days":null,"updated_at":null}`),
		LogRetentionDays: 30, LogArchiveAfterDays: 0, LogLevel: "INFO",
		SnapshotsKeep: 5, QueueMaxSize: defaultQueueMaxSize,
		BuildRetryMax: 0, BuildRetryBaseSeconds: 5,
		CommitStatusEnabled: false, CommitStatusContext: "Adlaire CI",
		BuildTrendKeepCount:     1000,
		DurationAnomaly:         json.RawMessage(`{"enabled":false,"min_samples":20,"avg_multiplier":2.0,"p95_multiplier":1.5}`),
		SessionTimeoutSeconds:   28800,
		ScheduleIntervalSeconds: 300,
		WatchMode:               "github",
		TagFilter:               json.RawMessage(`{"enabled":false,"patterns":[]}`),
		BuildCacheEnabled:       false,
		DeployParallelism:       1,
		RemoteBuild:             json.RawMessage(`{"enabled":false,"host":null,"user":null,"work_dir":null,"command_args":[],"artifact_path":null}`),
		ApprovalTimeoutSeconds:  86400,
		APIRateLimit:            apiRateLimitPolicyToMap(defaultAPIRateLimitPolicy()),
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
	if cfg.WatchMode == "" {
		cfg.WatchMode = def.WatchMode
	}
	if len(cfg.HistoryRetention) == 0 {
		cfg.HistoryRetention = append(json.RawMessage(nil), def.HistoryRetention...)
	}
	if len(cfg.DurationAnomaly) == 0 {
		cfg.DurationAnomaly = append(json.RawMessage(nil), def.DurationAnomaly...)
	}
	if len(cfg.TagFilter) == 0 {
		cfg.TagFilter = append(json.RawMessage(nil), def.TagFilter...)
	}
	if cfg.DeployParallelism == 0 {
		cfg.DeployParallelism = def.DeployParallelism
	}
	if len(cfg.RemoteBuild) == 0 {
		cfg.RemoteBuild = append(json.RawMessage(nil), def.RemoteBuild...)
	}
	if cfg.ApprovalTimeoutSeconds == 0 {
		cfg.ApprovalTimeoutSeconds = def.ApprovalTimeoutSeconds
	}
	if cfg.APIRateLimit == nil {
		cfg.APIRateLimit = apiRateLimitPolicyToMap(defaultAPIRateLimitPolicy())
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

func validateConfigDryRun(cfg apiServerConfig, raw map[string]json.RawMessage, order []string) (apiServerConfig, []apiConfigValidationIssue, []apiConfigValidationIssue, bool) {
	errorsList := []apiConfigValidationIssue{}
	warningsList := []apiConfigValidationIssue{}
	seenWarnings := map[string]bool{}
	for _, key := range order {
		value := raw[key]
		if secretLikeKey(key) || !configPatchKeyAllowed(key) {
			errorsList = append(errorsList, apiConfigValidationIssue{Field: key, Code: "unknown_key", Message: "Unknown config key"})
			continue
		}
		switch key {
		case "log_max_lines":
			if v, ok := rawInt(value, 1, 10000); ok {
				cfg.LogMaxLines = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "log_max_lines must be 1..10000"))
			}
		case "history_max_count":
			if v, ok := rawInt(value, 1, 10000); ok {
				cfg.HistoryMaxCount = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "history_max_count must be 1..10000"))
			}
		case "build_timeout_seconds":
			if v, ok := rawInt(value, 1, 86400); ok {
				cfg.BuildTimeoutSeconds = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "build_timeout_seconds must be 1..86400"))
			}
		case "log_retention_days":
			if v, ok := rawInt(value, 0, 3650); ok {
				cfg.LogRetentionDays = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "log_retention_days must be 0..3650"))
			}
		case "log_archive_after_days":
			if v, ok := rawInt(value, 0, 3650); ok {
				cfg.LogArchiveAfterDays = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "log_archive_after_days must be 0..3650"))
			}
		case "log_level":
			var v string
			if json.Unmarshal(value, &v) == nil && (v == "INFO" || v == "DEBUG" || v == "WARNING" || v == "ERROR") {
				cfg.LogLevel = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "invalid_value", "log_level must be INFO, DEBUG, WARNING, or ERROR"))
			}
		case "pat_expires_at":
			if v, ok := rawNullableString(value); ok {
				cfg.PATExpiresAt = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "invalid_value", "pat_expires_at must be string or null"))
			}
		case "commit_status_target_url":
			if v, ok := rawNullableString(value); ok {
				cfg.CommitStatusTargetURL = v
				if v != nil && strings.HasPrefix(*v, "http://") && !seenWarnings[key] {
					warningsList = append(warningsList, apiConfigValidationIssue{Field: key, Code: "http_url", Message: "https is recommended"})
					seenWarnings[key] = true
				}
			} else {
				errorsList = append(errorsList, configValidationError(key, "invalid_value", "commit_status_target_url must be string or null"))
			}
		case "snapshots_keep":
			if v, ok := rawInt(value, 0, 100); ok {
				cfg.SnapshotsKeep = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "snapshots_keep must be 0..100"))
			}
		case "queue_max_size":
			if v, ok := rawInt(value, 0, 100); ok {
				cfg.QueueMaxSize = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "queue_max_size must be 0..100"))
			}
		case "build_retry_max":
			if v, ok := rawInt(value, 0, 10); ok {
				cfg.BuildRetryMax = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "build_retry_max must be 0..10"))
			}
		case "build_retry_base_seconds":
			if v, ok := rawInt(value, 1, 3600); ok {
				cfg.BuildRetryBaseSeconds = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "build_retry_base_seconds must be 1..3600"))
			}
		case "commit_status_enabled":
			var v bool
			if json.Unmarshal(value, &v) == nil {
				cfg.CommitStatusEnabled = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "invalid_value", "commit_status_enabled must be boolean"))
			}
		case "commit_status_context":
			var v string
			if json.Unmarshal(value, &v) == nil && len(v) >= 1 && len(v) <= 100 && !containsControl(v) {
				cfg.CommitStatusContext = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "invalid_value", "commit_status_context must be 1..100 characters"))
			}
		case "build_trend_keep_count":
			if v, ok := rawInt(value, 10, 10000); ok {
				cfg.BuildTrendKeepCount = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "build_trend_keep_count must be 10..10000"))
			}
		case "session_timeout_seconds":
			if v, ok := rawInt(value, 300, 2592000); ok {
				cfg.SessionTimeoutSeconds = v
			} else {
				errorsList = append(errorsList, configValidationError(key, "out_of_range", "session_timeout_seconds must be 300..2592000"))
			}
		}
	}
	sort.Slice(warningsList, func(i, j int) bool {
		return warningsList[i].Field < warningsList[j].Field
	})
	return cfg, errorsList, warningsList, len(errorsList) == 0
}

func configPatchKeyAllowed(key string) bool {
	switch key {
	case "log_max_lines", "history_max_count", "build_timeout_seconds", "log_retention_days", "log_archive_after_days", "log_level", "pat_expires_at", "snapshots_keep", "queue_max_size", "build_retry_max", "build_retry_base_seconds", "commit_status_enabled", "commit_status_context", "commit_status_target_url", "build_trend_keep_count", "session_timeout_seconds":
		return true
	default:
		return false
	}
}

func serverConfigStorageKeyAllowed(key string) bool {
	switch key {
	case "log_max_lines", "history_max_count", "history_retention", "build_timeout_seconds",
		"log_retention_days", "log_archive_after_days", "log_level", "pat_expires_at",
		"snapshots_keep", "queue_max_size", "build_retry_max", "build_retry_base_seconds",
		"commit_status_enabled", "commit_status_context", "commit_status_target_url",
		"build_trend_keep_count", "duration_anomaly", "watch_mode", "tag_filter",
		"build_cache_enabled", "deploy_parallelism", "remote_build", "approval_timeout_seconds",
		"force_build_interval_hours", "build_cooldown_seconds", "schedule_interval_seconds",
		"schedule_paused", "allowed_hours", "session_timeout_seconds", "api_rate_limit":
		return true
	default:
		return false
	}
}

func configValidationError(field, code, message string) apiConfigValidationIssue {
	return apiConfigValidationIssue{Field: field, Code: code, Message: message}
}

func rawInt(raw json.RawMessage, min, max int) (int, bool) {
	var value int
	if err := json.Unmarshal(raw, &value); err != nil || value < min || value > max {
		return 0, false
	}
	return value, true
}

func rawNullableString(raw json.RawMessage) (*string, bool) {
	var value *string
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, false
	}
	if value != nil && containsControl(*value) {
		return nil, false
	}
	return value, true
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

func (s *APIServer) accessAllowed(r *http.Request) (bool, error) {
	if r.Method == http.MethodGet && r.URL.Path == "/api/health" {
		return true, nil
	}
	cfg, err := s.readAccessControl()
	if err != nil {
		return false, err
	}
	if len(cfg.Allow) == 0 {
		return true, nil
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(strings.TrimSpace(host))
	if ip == nil || ip.To4() == nil {
		return false, nil
	}
	ip = ip.To4()
	for _, entry := range cfg.Allow {
		if strings.Contains(entry, "/") {
			_, network, err := net.ParseCIDR(entry)
			if err == nil && network.Contains(ip) {
				return true, nil
			}
			continue
		}
		allowed := net.ParseIP(entry)
		if allowed != nil && allowed.To4() != nil && allowed.To4().Equal(ip) {
			return true, nil
		}
	}
	return false, nil
}

func (s *APIServer) appendAPIAccessLog(r *http.Request, requestID string, status int, start time.Time) error {
	duration := time.Since(start).Milliseconds()
	if duration < 0 {
		duration = 0
	}
	auth := apiAuthInfo{AuthType: "none"}
	if capture, ok := r.Context().Value(apiAuthCaptureContextKey{}).(*apiAuthInfo); ok && capture.AuthType != "" {
		auth = *capture
	} else if value, ok := r.Context().Value(apiAuthContextKey{}).(apiAuthInfo); ok {
		auth = value
	} else if r.URL.Path == "/api/webhook" && status < 400 {
		actor := "webhook"
		auth = apiAuthInfo{AuthType: "webhook", Actor: &actor}
	}
	record := apiAccessLogRecord{
		At:         s.nowString(),
		RequestID:  requestID,
		Method:     strings.ToUpper(r.Method),
		Path:       r.URL.Path,
		Query:      apiAccessLogQuery(r, status),
		Status:     status,
		DurationMS: duration,
		AuthType:   auth.AuthType,
		Actor:      auth.Actor,
		RemoteAddr: apiRemoteAddr(r),
		UserAgent:  apiUserAgent(r),
		Error:      apiErrorField(status),
	}
	return appendJSONLine(filepath.Join(s.cfg.StateDir, ".api_access_log"), record)
}

func apiAccessLogQuery(r *http.Request, status int) map[string]string {
	query := map[string]string{}
	if status >= 400 {
		return query
	}
	for key, values := range r.URL.Query() {
		if len(values) == 0 || !apiAccessLogQueryKeyAllowed(key) {
			continue
		}
		value := strings.TrimSpace(values[0])
		if value == "" {
			continue
		}
		if apiAccessLogQueryValuePlain(key, value) {
			query[key] = value
		} else {
			query[key] = "***"
		}
	}
	return query
}

func apiAccessLogQueryKeyAllowed(key string) bool {
	switch key {
	case "action", "days", "enabled", "failure_category", "flagged", "from", "limit", "method", "n", "offset", "page", "path", "per_page", "q", "result", "status", "tag", "to", "trigger", "type":
		return true
	default:
		return false
	}
}

func apiAccessLogQueryValuePlain(key, value string) bool {
	switch key {
	case "days", "from", "limit", "n", "offset", "page", "per_page", "status", "to":
		_, err := strconv.Atoi(value)
		return err == nil
	case "enabled", "flagged":
		return value == "true" || value == "false"
	case "method":
		return value == http.MethodGet || value == http.MethodPost || value == http.MethodDelete || value == http.MethodPut || value == http.MethodPatch
	default:
		return false
	}
}

func apiRemoteAddr(r *http.Request) *string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(strings.TrimSpace(host))
	if ip == nil {
		return nil
	}
	value := ip.String()
	return &value
}

func apiUserAgent(r *http.Request) *string {
	raw := strings.TrimSpace(r.UserAgent())
	if raw == "" {
		return nil
	}
	var b strings.Builder
	count := 0
	for _, rr := range raw {
		if rr < 0x20 || rr == 0x7f {
			continue
		}
		if count >= 512 {
			break
		}
		b.WriteRune(rr)
		count++
	}
	if b.Len() == 0 {
		return nil
	}
	value := b.String()
	return &value
}

func apiErrorField(status int) *string {
	if status < 400 {
		return nil
	}
	message := http.StatusText(status)
	if message == "" {
		message = "HTTP error"
	}
	return &message
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

type apiRestoreJSONWrite struct {
	File   string
	Value  map[string]any
	Delete bool
}

type apiRestoreSecretWrite struct {
	File   string
	Value  string
	Delete bool
}

type apiRestorePlan struct {
	JSONWrites    []apiRestoreJSONWrite
	SecretWrites  []apiRestoreSecretWrite
	ChangeSummary map[string]any
}

func (s *APIServer) backupObject() (map[string]any, error) {
	out := map[string]any{"exported_at": s.nowString()}
	for _, key := range backupConfigKeys() {
		value, err := s.backupJSONValue(key, true)
		if err != nil {
			return nil, err
		}
		out[key] = value
	}
	out["webhook_secret_set"] = fileExists(filepath.Join(s.cfg.StateDir, ".webhook_secret"))
	out["smtp_password_set"] = fileExists(filepath.Join(s.cfg.StateDir, ".smtp_secret"))
	return out, nil
}

func (s *APIServer) restorePlan(w http.ResponseWriter, body map[string]json.RawMessage) (apiRestorePlan, bool, bool) {
	if len(body) == 0 {
		writeValidation(w, "body", "required")
		return apiRestorePlan{}, false, false
	}
	if raw, ok := body["exported_at"]; ok {
		var exportedAt string
		if err := json.Unmarshal(raw, &exportedAt); err != nil || !validAPITime(exportedAt) {
			writeValidation(w, "exported_at", "invalid value")
			return apiRestorePlan{}, false, false
		}
	}
	allowed := restoreAllowedKeys()
	for key := range body {
		if !allowed[key] {
			writeValidation(w, key, "unknown key")
			return apiRestorePlan{}, false, false
		}
	}
	plan := apiRestorePlan{ChangeSummary: map[string]any{"files": []string{}, "secrets": []string{}}}
	changed := false
	for _, key := range backupConfigKeys() {
		raw, ok := body[key]
		if !ok {
			writeValidation(w, key, "required")
			return apiRestorePlan{}, false, false
		}
		value, deleteFile, valid := s.restoreJSONValue(w, key, raw)
		if !valid {
			return apiRestorePlan{}, false, false
		}
		current, err := s.backupJSONValue(key, false)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
			return apiRestorePlan{}, false, false
		}
		if mapsEqual(current, value) && !deleteFile {
			continue
		}
		changed = true
		file := "." + key
		plan.JSONWrites = append(plan.JSONWrites, apiRestoreJSONWrite{File: file, Value: value, Delete: deleteFile})
		plan.ChangeSummary["files"] = appendStringAny(plan.ChangeSummary["files"], key)
	}
	secretChanged, ok := s.addRestoreSecretPlan(w, body, "webhook_secret", ".webhook_secret", &plan)
	if !ok {
		return apiRestorePlan{}, false, false
	}
	changed = changed || secretChanged
	secretChanged, ok = s.addRestoreSecretPlan(w, body, "smtp_password", ".smtp_secret", &plan)
	if !ok {
		return apiRestorePlan{}, false, false
	}
	changed = changed || secretChanged
	return plan, changed, true
}

func backupConfigKeys() []string {
	out := make([]string, 0, len(backupConfigFiles()))
	for _, file := range backupConfigFiles() {
		out = append(out, strings.TrimPrefix(file, "."))
	}
	return out
}

func restoreAllowedKeys() map[string]bool {
	allowed := map[string]bool{
		"exported_at":        true,
		"webhook_secret_set": true,
		"smtp_password_set":  true,
		"webhook_secret":     true,
		"smtp_password":      true,
	}
	for _, key := range backupConfigKeys() {
		allowed[key] = true
	}
	return allowed
}

func (s *APIServer) backupJSONValue(key string, mask bool) (map[string]any, error) {
	path := filepath.Join(s.cfg.StateDir, "."+key)
	var value map[string]any
	switch key {
	case "server_config":
		raw, ok, err := readOptionalJSONMap(path)
		if err != nil {
			return nil, err
		}
		if ok {
			value = raw
		} else {
			value = map[string]any{}
		}
	case "notify_config":
		cfg, err := s.readNotifyConfig()
		if err != nil {
			return nil, err
		}
		if mask {
			cfg = maskNotifyConfig(cfg)
		}
		return mapFromJSONValue(cfg)
	case "repo_config":
		cfg, err := s.readRepoConfig()
		if err != nil {
			return nil, err
		}
		value = map[string]any{"owner": cfg.Owner, "repo": cfg.Repo, "updated_at": nil}
		if cfg.UpdatedAt != "" {
			value["updated_at"] = cfg.UpdatedAt
		}
	case "branch_config":
		cfg, _, err := s.readBranchConfig()
		if err != nil {
			return nil, err
		}
		return mapFromJSONValue(cfg)
	case "access_control":
		cfg, err := s.readAccessControl()
		if err != nil {
			return nil, err
		}
		return mapFromJSONValue(cfg)
	case "hooks":
		rules, err := readRuleList(path)
		if err != nil {
			return nil, err
		}
		value = map[string]any{"hooks": rules}
	case "alert_rules", "tag_rules":
		rules, err := readRuleList(path)
		if err != nil {
			return nil, err
		}
		value = map[string]any{"rules": rules}
	case "pipeline_config":
		raw, ok, err := readOptionalJSONMap(path)
		if err != nil {
			return nil, err
		}
		if ok {
			value = raw
		} else {
			value = defaultPipelineConfig()
		}
	case "dashboard_layout":
		raw, ok, err := readOptionalJSONMap(path)
		if err != nil {
			return nil, err
		}
		if ok {
			value = raw
		} else {
			value = defaultDashboardLayout()
		}
	case "smtp_config":
		raw, ok, err := readOptionalJSONMap(path)
		if err != nil {
			return nil, err
		}
		if ok {
			value = raw
		} else {
			value = defaultSMTPConfig()
		}
	default:
		return nil, fmt.Errorf("unknown backup key %s", key)
	}
	if mask {
		return maskSecrets(value), nil
	}
	return value, nil
}

func (s *APIServer) restoreJSONValue(w http.ResponseWriter, key string, raw json.RawMessage) (map[string]any, bool, bool) {
	value, ok := decodeRestoreObject(w, key, raw)
	if !ok {
		return nil, false, false
	}
	current, _ := s.backupJSONValue(key, false)
	if containsMaskedSecret(value) {
		value = mergeMaskedSecrets(value, current)
	}
	switch key {
	case "server_config":
		if !validateRestoreServerConfig(w, raw) {
			return nil, false, false
		}
	case "notify_config":
		var cfg apiNotifyConfig
		if err := json.Unmarshal(mustJSON(value), &cfg); err != nil {
			writeValidation(w, key, "invalid value")
			return nil, false, false
		}
		normalized, valid := normalizeNotifyConfig(w, cfg)
		if !valid {
			return nil, false, false
		}
		next, err := mapFromJSONValue(normalized)
		if err != nil {
			writeValidation(w, key, "invalid value")
			return nil, false, false
		}
		value = next
	case "repo_config":
		next, deleteFile, valid := s.restoreRepoConfigValue(w, value)
		return next, deleteFile, valid
	case "branch_config":
		var cfg apiBranchConfigFile
		if err := json.Unmarshal(mustJSON(value), &cfg); err != nil {
			writeValidation(w, key, "invalid value")
			return nil, false, false
		}
		if cfg.BranchTargets == nil {
			cfg.BranchTargets = []apiBranchTarget{}
		}
		if !validateBranchTargets(w, cfg.BranchTargets) {
			return nil, false, false
		}
		next, err := mapFromJSONValue(cfg)
		if err != nil {
			writeValidation(w, key, "invalid value")
			return nil, false, false
		}
		return next, len(cfg.BranchTargets) == 0, true
	case "access_control":
		var cfg apiAccessControl
		if err := json.Unmarshal(mustJSON(value), &cfg); err != nil {
			writeValidation(w, key, "invalid value")
			return nil, false, false
		}
		normalized, valid := normalizeAccessControl(w, cfg)
		if !valid {
			return nil, false, false
		}
		next, err := mapFromJSONValue(normalized)
		if err != nil {
			writeValidation(w, key, "invalid value")
			return nil, false, false
		}
		value = next
	case "hooks":
		if !validateRestoreRuleWrapper(w, key, value, "hooks", 50) {
			return nil, false, false
		}
	case "alert_rules", "tag_rules":
		if !validateRestoreRuleWrapper(w, key, value, "rules", 100) {
			return nil, false, false
		}
	case "pipeline_config":
		if !validatePipelineConfig(w, value) {
			return nil, false, false
		}
	case "dashboard_layout":
		if !validateDashboardLayout(w, value) {
			return nil, false, false
		}
	case "smtp_config":
		if !validateSMTPConfig(w, value) {
			return nil, false, false
		}
	}
	return value, false, true
}

func (s *APIServer) restoreRepoConfigValue(w http.ResponseWriter, value map[string]any) (map[string]any, bool, bool) {
	for key := range value {
		if key != "owner" && key != "repo" && key != "updated_at" {
			writeValidation(w, key, "unknown key")
			return nil, false, false
		}
	}
	owner, ok := value["owner"].(string)
	if !ok || strings.TrimSpace(owner) == "" || containsControl(owner) {
		writeValidation(w, "owner", "invalid value")
		return nil, false, false
	}
	repo, ok := value["repo"].(string)
	if !ok || strings.TrimSpace(repo) == "" || containsControl(repo) {
		writeValidation(w, "repo", "invalid value")
		return nil, false, false
	}
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	var updatedAt any
	if raw, ok := value["updated_at"]; ok && raw != nil {
		text, ok := raw.(string)
		if !ok || !validAPITime(text) {
			writeValidation(w, "updated_at", "invalid value")
			return nil, false, false
		}
		updatedAt = text
	}
	def := defaultRepoConfig()
	if owner == def.Owner && repo == def.Repo && updatedAt == nil {
		return map[string]any{"owner": owner, "repo": repo, "updated_at": nil}, fileExists(filepath.Join(s.cfg.StateDir, ".repo_config")), true
	}
	if updatedAt == nil {
		updatedAt = s.nowString()
	}
	return map[string]any{"owner": owner, "repo": repo, "updated_at": updatedAt}, false, true
}

func (s *APIServer) addRestoreSecretPlan(w http.ResponseWriter, body map[string]json.RawMessage, field, file string, plan *apiRestorePlan) (bool, bool) {
	raw, exists := body[field]
	if !exists {
		return false, true
	}
	deleteSecret, value, keep, ok := parseRestoreSecret(w, field, raw)
	if !ok || keep {
		return false, ok
	}
	path := filepath.Join(s.cfg.StateDir, file)
	current, err := os.ReadFile(path)
	currentExists := err == nil
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return false, false
	}
	if deleteSecret {
		if !currentExists {
			return false, true
		}
		plan.SecretWrites = append(plan.SecretWrites, apiRestoreSecretWrite{File: file, Delete: true})
		plan.ChangeSummary["secrets"] = appendStringAny(plan.ChangeSummary["secrets"], field)
		return true, true
	}
	if currentExists && string(current) == value {
		return false, true
	}
	plan.SecretWrites = append(plan.SecretWrites, apiRestoreSecretWrite{File: file, Value: value})
	plan.ChangeSummary["secrets"] = appendStringAny(plan.ChangeSummary["secrets"], field)
	return true, true
}

func decodeRestoreObject(w http.ResponseWriter, field string, raw json.RawMessage) (map[string]any, bool) {
	var value map[string]any
	if err := json.Unmarshal(raw, &value); err != nil || value == nil {
		writeValidation(w, field, "object required")
		return nil, false
	}
	return value, true
}

func parseRestoreSecret(w http.ResponseWriter, field string, raw json.RawMessage) (bool, string, bool, bool) {
	if string(raw) == "null" {
		return true, "", false, true
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		writeValidation(w, field, "invalid value")
		return false, "", false, false
	}
	if value == "***" {
		return false, "", true, true
	}
	if !validRestoreSecret(value) {
		writeValidation(w, field, "invalid value")
		return false, "", false, false
	}
	return false, value, false, true
}

func validRestoreSecret(value string) bool {
	return value != "" && len([]byte(value)) <= 4096 && strings.ToValidUTF8(value, "") == value && !strings.HasPrefix(value, "\ufeff") && !containsControl(value)
}

func validateRestoreServerConfig(w http.ResponseWriter, raw json.RawMessage) bool {
	fields := map[string]json.RawMessage{}
	if err := json.Unmarshal(raw, &fields); err != nil {
		writeValidation(w, "server_config", "invalid value")
		return false
	}
	for key, value := range fields {
		if secretLikeKey(key) {
			writeValidation(w, key, "unknown key")
			return false
		}
		switch key {
		case "log_max_lines":
			if !validRawInt(value, 1, 10000) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "history_max_count":
			if !validRawInt(value, 1, 10000) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "build_timeout_seconds":
			if !validRawInt(value, 1, 86400) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "log_retention_days", "log_archive_after_days":
			if !validRawInt(value, 0, 3650) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "log_level":
			var v string
			if json.Unmarshal(value, &v) != nil || (v != "INFO" && v != "DEBUG" && v != "WARNING" && v != "ERROR") {
				writeValidation(w, key, "invalid value")
				return false
			}
		case "pat_expires_at", "commit_status_target_url":
			var v *string
			if json.Unmarshal(value, &v) != nil || (v != nil && containsControl(*v)) {
				writeValidation(w, key, "invalid value")
				return false
			}
		case "snapshots_keep", "queue_max_size":
			if !validRawInt(value, 0, 100) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "build_retry_max":
			if !validRawInt(value, 0, 10) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "build_retry_base_seconds":
			if !validRawInt(value, 1, 3600) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "commit_status_enabled", "schedule_paused":
			var v bool
			if json.Unmarshal(value, &v) != nil {
				writeValidation(w, key, "invalid value")
				return false
			}
		case "commit_status_context":
			var v string
			if json.Unmarshal(value, &v) != nil || len(v) < 1 || len(v) > 100 || containsControl(v) {
				writeValidation(w, key, "invalid value")
				return false
			}
		case "build_trend_keep_count":
			if !validRawInt(value, 10, 10000) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "session_timeout_seconds":
			if !validRawInt(value, 300, 2592000) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "force_build_interval_hours":
			if !validRawInt(value, 0, 8760) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "build_cooldown_seconds":
			if !validRawInt(value, 0, 86400) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "schedule_interval_seconds":
			if !validRawInt(value, 30, 86400) {
				writeValidation(w, key, "out of range")
				return false
			}
		case "allowed_hours":
			if string(value) == "null" {
				continue
			}
			var hours apiAllowedHours
			if json.Unmarshal(value, &hours) != nil || hours.From < 0 || hours.From > 23 || hours.To < 0 || hours.To > 23 || hours.From > hours.To {
				writeValidation(w, key, "invalid value")
				return false
			}
		case "api_rate_limit":
			var rate map[string]any
			if json.Unmarshal(value, &rate) != nil || rate == nil {
				writeValidation(w, key, "invalid value")
				return false
			}
			if _, ok := apiRateLimitPolicyFromConfig(rate); !ok {
				writeValidation(w, key, "invalid value")
				return false
			}
		default:
			writeValidation(w, key, "unknown key")
			return false
		}
	}
	return true
}

func validateRestoreRuleWrapper(w http.ResponseWriter, field string, value map[string]any, listKey string, max int) bool {
	if len(value) != 1 {
		writeValidation(w, field, "invalid value")
		return false
	}
	raw, ok := value[listKey].([]any)
	if !ok || len(raw) > max {
		writeValidation(w, listKey, "invalid value")
		return false
	}
	for i, item := range raw {
		child, ok := item.(map[string]any)
		if !ok || fmt.Sprint(child["id"]) == "" {
			writeValidation(w, fmt.Sprintf("%s[%d]", listKey, i), "invalid value")
			return false
		}
	}
	return true
}

func validateSMTPConfig(w http.ResponseWriter, value map[string]any) bool {
	required := map[string]bool{"host": true, "port": true, "user": true, "tls": true, "from": true, "to": true, "on": true, "enabled": true}
	if len(value) != len(required) {
		writeValidation(w, "smtp_config", "invalid value")
		return false
	}
	for key := range value {
		if !required[key] {
			writeValidation(w, key, "unknown key")
			return false
		}
	}
	if raw, ok := value["port"].(float64); !ok || raw < 1 || raw > 65535 || math.Trunc(raw) != raw {
		writeValidation(w, "port", "out of range")
		return false
	}
	if _, ok := value["tls"].(bool); !ok {
		writeValidation(w, "tls", "invalid value")
		return false
	}
	if _, ok := value["enabled"].(bool); !ok {
		writeValidation(w, "enabled", "invalid value")
		return false
	}
	for _, key := range []string{"host", "user", "from"} {
		if value[key] == nil {
			continue
		}
		text, ok := value[key].(string)
		if !ok || containsControl(text) || len(text) > 255 {
			writeValidation(w, key, "invalid value")
			return false
		}
	}
	for _, key := range []string{"to", "on"} {
		list, ok := value[key].([]any)
		if !ok || len(list) > 50 {
			writeValidation(w, key, "invalid value")
			return false
		}
		for _, item := range list {
			text, ok := item.(string)
			if !ok || text == "" || containsControl(text) {
				writeValidation(w, key, "invalid value")
				return false
			}
		}
	}
	return true
}

func defaultDashboardLayout() map[string]any {
	return map[string]any{"widgets": []any{"status", "stats", "schedule", "alerts", "disk", "rate_limit", "snapshots", "maintenance", "queue"}}
}

func defaultSMTPConfig() map[string]any {
	return map[string]any{"host": nil, "port": 587, "user": nil, "tls": true, "from": nil, "to": []any{}, "on": []any{}, "enabled": false}
}

func validRawInt(raw json.RawMessage, min, max int) bool {
	var value int
	return json.Unmarshal(raw, &value) == nil && value >= min && value <= max
}

func validAPITime(value string) bool {
	parsed, err := time.Parse(apiTimeLayout, value)
	return err == nil && parsed.UTC().Format(apiTimeLayout) == value
}

func secretLikeKey(key string) bool {
	lower := strings.ToLower(key)
	return strings.Contains(lower, "secret") || strings.Contains(lower, "password") || strings.Contains(lower, "token")
}

func mapFromJSONValue(value any) (map[string]any, error) {
	var out map[string]any
	if err := json.Unmarshal(mustJSON(value), &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, nil
}

func mustJSON(value any) []byte {
	data, _ := json.Marshal(value)
	return data
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func appendStringAny(value any, item string) []string {
	if list, ok := value.([]string); ok {
		return append(list, item)
	}
	return []string{item}
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
		if list, ok := value.([]any); ok {
			curList, _ := current[key].([]any)
			out[key] = mergeMaskedSecretLists(list, curList)
			continue
		}
		out[key] = value
	}
	return out
}

func mergeMaskedSecretLists(next, current []any) []any {
	out := make([]any, len(next))
	for i, value := range next {
		if child, ok := value.(map[string]any); ok {
			var curChild map[string]any
			if i < len(current) {
				curChild, _ = current[i].(map[string]any)
			}
			out[i] = mergeMaskedSecrets(child, curChild)
			continue
		}
		out[i] = value
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
	if id == "" || len(id) > 64 {
		return false
	}
	for _, ch := range id {
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-' {
			continue
		}
		return false
	}
	return true
}

func readAPITokens(path string) ([]apiTokenRecord, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []apiTokenRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	if len(raw) != 1 {
		return nil, errors.New("invalid api token file")
	}
	listRaw, ok := raw["tokens"]
	if !ok {
		return nil, errors.New("invalid api token file")
	}
	var rows []map[string]json.RawMessage
	if err := json.Unmarshal(listRaw, &rows); err != nil {
		return nil, err
	}
	if len(rows) > 100 {
		return nil, errors.New("invalid api token file")
	}
	tokens := make([]apiTokenRecord, 0, len(rows))
	for _, row := range rows {
		record, err := decodeAPITokenRecord(row)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, record)
	}
	return tokens, nil
}

func writeAPITokens(path string, tokens []apiTokenRecord) error {
	sort.SliceStable(tokens, func(i, j int) bool {
		if tokens[i].CreatedAt == tokens[j].CreatedAt {
			return tokens[i].ID < tokens[j].ID
		}
		return tokens[i].CreatedAt > tokens[j].CreatedAt
	})
	return atomicWriteJSONLocked(path, apiTokenFile{Tokens: tokens}, 0600)
}

func updateAPITokensLocked(path string, mutate func([]apiTokenRecord) ([]apiTokenRecord, bool, error)) ([]apiTokenRecord, error) {
	release, err := acquireAPIStateFileLock(path)
	if err != nil {
		return nil, err
	}
	defer release()
	current, err := readAPITokens(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errAPITokenCorrupt, err)
	}
	next, write, err := mutate(current)
	if err != nil {
		return nil, err
	}
	if !write {
		return next, nil
	}
	if err := validateAPITokenRecords(next); err != nil {
		return nil, err
	}
	if err := writeAPITokens(path, next); err != nil {
		return nil, err
	}
	return next, nil
}

func validateAPITokenRecords(tokens []apiTokenRecord) error {
	if len(tokens) > 100 {
		return errors.New("invalid api token file")
	}
	encoded, err := json.Marshal(apiTokenFile{Tokens: tokens})
	if err != nil {
		return err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &raw); err != nil {
		return err
	}
	listRaw, ok := raw["tokens"]
	if !ok || len(raw) != 1 {
		return errors.New("invalid api token file")
	}
	var rows []map[string]json.RawMessage
	if err := json.Unmarshal(listRaw, &rows); err != nil {
		return err
	}
	for _, row := range rows {
		if _, err := decodeAPITokenRecord(row); err != nil {
			return err
		}
	}
	return nil
}

func decodeAPITokenRecord(row map[string]json.RawMessage) (apiTokenRecord, error) {
	allowed := map[string]bool{"id": true, "label": true, "scopes": true, "token_hash": true, "created_at": true, "last_used_at": true, "expires_at": true, "revoked_at": true}
	if len(row) != len(allowed) {
		return apiTokenRecord{}, errors.New("invalid api token record")
	}
	for key := range row {
		if !allowed[key] {
			return apiTokenRecord{}, errors.New("invalid api token record")
		}
	}
	id, ok := rawJSONString(row["id"])
	if !ok || !validTokenID(id) {
		return apiTokenRecord{}, errors.New("invalid api token id")
	}
	label, ok := rawJSONString(row["label"])
	if !ok || strings.TrimSpace(label) == "" || len([]byte(label)) > 64 {
		return apiTokenRecord{}, errors.New("invalid api token label")
	}
	scopes, ok := rawJSONStringList(row["scopes"])
	if !ok {
		return apiTokenRecord{}, errors.New("invalid api token scopes")
	}
	scopes, ok = normalizeTokenScopes(scopes)
	if !ok {
		return apiTokenRecord{}, errors.New("invalid api token scopes")
	}
	hash, ok := rawJSONString(row["token_hash"])
	if !ok || !validLowerHex(hash, 64) {
		return apiTokenRecord{}, errors.New("invalid api token hash")
	}
	createdAt, ok := rawJSONString(row["created_at"])
	if !ok || !validAPITime(createdAt) {
		return apiTokenRecord{}, errors.New("invalid api token created_at")
	}
	lastUsedAt, ok := rawNullableAPITime(row["last_used_at"])
	if !ok {
		return apiTokenRecord{}, errors.New("invalid api token last_used_at")
	}
	expiresAt, ok := rawNullableAPITime(row["expires_at"])
	if !ok {
		return apiTokenRecord{}, errors.New("invalid api token expires_at")
	}
	revokedAt, ok := rawNullableAPITime(row["revoked_at"])
	if !ok {
		return apiTokenRecord{}, errors.New("invalid api token revoked_at")
	}
	return apiTokenRecord{ID: id, Label: label, Scopes: scopes, TokenHash: hash, CreatedAt: createdAt, LastUsedAt: lastUsedAt, ExpiresAt: expiresAt, RevokedAt: revokedAt}, nil
}

func normalizeTokenScopes(scopes []string) ([]string, bool) {
	seen := map[string]bool{}
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if !validTokenScope(scope) || seen[scope] {
			return nil, false
		}
		seen[scope] = true
	}
	if len(seen) == 0 || len(seen) > 5 {
		return nil, false
	}
	out := []string{}
	for _, scope := range []string{"trigger", "read", "operate", "config", "admin"} {
		if seen[scope] {
			out = append(out, scope)
		}
	}
	return out, true
}

func normalizeTokenExpiresAt(value *string, now time.Time) (*string, bool) {
	if value == nil {
		return nil, true
	}
	trimmed := strings.TrimSpace(*value)
	parsed, err := time.Parse(apiTimeLayout, trimmed)
	if err != nil || !now.Before(parsed) {
		return nil, false
	}
	return &trimmed, true
}

func apiTokenExpired(record apiTokenRecord, now time.Time) bool {
	if record.ExpiresAt == nil {
		return false
	}
	expires, err := time.Parse(apiTimeLayout, *record.ExpiresAt)
	return err != nil || !now.Before(expires)
}

func validTokenScope(scope string) bool {
	switch scope {
	case "trigger", "read", "operate", "config", "admin":
		return true
	default:
		return false
	}
}

func validTokenID(id string) bool {
	if !strings.HasPrefix(id, "tok") || len(id) < 9 {
		return false
	}
	for _, ch := range id[3:] {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func validLowerHex(value string, size int) bool {
	if len(value) != size {
		return false
	}
	for _, ch := range value {
		if (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') {
			continue
		}
		return false
	}
	return true
}

func rawJSONString(raw json.RawMessage) (string, bool) {
	var value string
	if len(raw) == 0 || json.Unmarshal(raw, &value) != nil {
		return "", false
	}
	return value, true
}

func rawJSONStringList(raw json.RawMessage) ([]string, bool) {
	var value []string
	if len(raw) == 0 || json.Unmarshal(raw, &value) != nil {
		return nil, false
	}
	return value, true
}

func rawNullableAPITime(raw json.RawMessage) (*string, bool) {
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil, true
	}
	value, ok := rawJSONString(raw)
	if !ok || !validAPITime(value) {
		return nil, false
	}
	return &value, true
}

func apiTokenScopesAllow(scopes []string, method, path string) bool {
	required, ok := apiTokenRequiredScope(method, path)
	if !ok {
		return false
	}
	for _, scope := range scopes {
		if scope == required {
			return true
		}
	}
	return false
}

func apiTokenRequiredScope(method, path string) (string, bool) {
	if method == http.MethodPost && (path == "/api/build" || path == "/api/build/force") {
		return "trigger", true
	}
	if apiTokenReadEndpoint(method, path) {
		return "read", true
	}
	if apiTokenOperateEndpoint(method, path) {
		return "operate", true
	}
	if apiTokenConfigEndpoint(method, path) {
		return "config", true
	}
	if apiTokenAdminEndpoint(method, path) {
		return "admin", true
	}
	return "", false
}

func apiTokenReadEndpoint(method, path string) bool {
	if method != http.MethodGet {
		return false
	}
	if apiTokenExactPath(path, []string{
		"/api/status", "/api/logs", "/api/logs/search", "/api/logs/export", "/api/history", "/api/history/export",
		"/api/sysinfo", "/api/health", "/api/schedule", "/api/notify-config", "/api/notify-log", "/api/config",
		"/api/config-log", "/api/pat-status", "/api/stats", "/api/stats/timeline", "/api/stats/build-duration",
		"/api/stats/build-trends", "/api/output-meta", "/api/repo-info", "/api/branch-config", "/api/backup",
		"/api/dashboard", "/api/diagnostics", "/api/rate-limit", "/api/disk-usage", "/api/webhook-events",
		"/api/webhook-config", "/api/snapshots", "/api/maintenance", "/api/access-control", "/api/hooks",
		"/api/alert-rules", "/api/tag-rules", "/api/pipeline-config", "/api/build-chain-config", "/api/notes",
		"/api/smtp-config", "/api/queue", "/api/approvals", "/api/dashboard-layout",
	}) {
		return true
	}
	return apiTokenSuffixEndpoint(path, "/api/history/", []string{"/log", "/comment"}) ||
		apiTokenSuffixEndpoint(path, "/api/snapshots/", []string{"/download"}) ||
		apiTokenSuffixEndpoint(path, "/api/hooks/", []string{"/log"})
}

func apiTokenOperateEndpoint(method, path string) bool {
	if method == http.MethodGet && path == "/api/build/stream" {
		return true
	}
	if method == http.MethodPost && apiTokenExactPath(path, []string{"/api/build/cancel", "/api/notify-test", "/api/notify/weekly-summary", "/api/pat-verify", "/api/circuit-breaker/reset", "/api/smtp-test", "/api/verify-output"}) {
		return true
	}
	if method == http.MethodDelete && path == "/api/queue" {
		return true
	}
	if method == http.MethodPost && apiTokenSuffixEndpoint(path, "/api/history/", []string{"/rollback"}) {
		return true
	}
	if method == http.MethodPost && apiTokenSuffixEndpoint(path, "/api/approvals/", []string{"/approve", "/reject"}) {
		return true
	}
	return false
}

func apiTokenConfigEndpoint(method, path string) bool {
	if method == http.MethodPost && apiTokenExactPath(path, []string{
		"/api/schedule/interval", "/api/schedule/pause", "/api/schedule/resume", "/api/schedule/allowed-hours",
		"/api/schedule/force-interval", "/api/schedule/cooldown", "/api/notify-config", "/api/config/validate",
		"/api/config", "/api/log-level", "/api/pat-update", "/api/repo-config", "/api/branch-config", "/api/restore",
		"/api/webhook-config", "/api/maintenance/enable", "/api/maintenance/disable", "/api/access-control",
		"/api/hooks", "/api/alert-rules", "/api/tag-rules", "/api/pipeline-config", "/api/build-chain-config",
		"/api/notes", "/api/smtp-config", "/api/dashboard-layout", "/api/logs/cleanup", "/api/logs/archive",
	}) {
		return true
	}
	if method == http.MethodDelete && apiTokenSuffixEndpoint(path, "/api/snapshots/", []string{""}) {
		return true
	}
	if method == http.MethodDelete && (apiTokenSuffixEndpoint(path, "/api/hooks/", []string{""}) || apiTokenSuffixEndpoint(path, "/api/alert-rules/", []string{""}) || apiTokenSuffixEndpoint(path, "/api/tag-rules/", []string{""})) {
		return true
	}
	if method == http.MethodPost && apiTokenSuffixEndpoint(path, "/api/history/", []string{"/comment", "/flag", "/tags"}) {
		return true
	}
	return false
}

func apiTokenAdminEndpoint(method, path string) bool {
	if method == http.MethodGet && apiTokenExactPath(path, []string{"/api/access-log", "/api/api-access-log", "/api/audit-log", "/api/api-rate-limit", "/api/sessions", "/api/auth/totp-status", "/api/tokens"}) {
		return true
	}
	if method == http.MethodPost && apiTokenExactPath(path, []string{"/api/api-rate-limit", "/api/sessions/revoke-all", "/api/auth/totp-setup", "/api/auth/totp-confirm", "/api/tokens"}) {
		return true
	}
	if method == http.MethodDelete && (path == "/api/auth/totp" || apiTokenSuffixEndpoint(path, "/api/tokens/", []string{""})) {
		return true
	}
	return false
}

func apiTokenExactPath(path string, allowed []string) bool {
	for _, item := range allowed {
		if path == item {
			return true
		}
	}
	return false
}

func apiTokenSuffixEndpoint(path, prefix string, suffixes []string) bool {
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	rest := strings.TrimPrefix(path, prefix)
	if rest == "" || strings.Contains(rest, "/") && !apiTokenSuffixRestAllowed(rest, suffixes) {
		return false
	}
	if len(suffixes) == 1 && suffixes[0] == "" {
		return !strings.Contains(rest, "/")
	}
	for _, suffix := range suffixes {
		if strings.HasSuffix(rest, suffix) {
			id := strings.TrimSuffix(rest, suffix)
			return id != "" && !strings.Contains(id, "/")
		}
	}
	return false
}

func apiTokenSuffixRestAllowed(rest string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if suffix != "" && strings.HasSuffix(rest, suffix) {
			return true
		}
	}
	return false
}

func readRuleList(path string) ([]map[string]any, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []map[string]any{}, nil
	}
	if err != nil {
		return nil, err
	}
	var wrapper map[string][]map[string]any
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	rules, ok := wrapper[ruleListKey(path)]
	if !ok || len(wrapper) != 1 {
		return nil, errors.New("invalid rule list")
	}
	if rules == nil {
		rules = []map[string]any{}
	}
	return rules, nil
}

func ruleListWrapper(path string, rules []map[string]any) map[string]any {
	return map[string]any{ruleListKey(path): rules}
}

func ruleListKey(path string) string {
	if filepath.Base(path) == ".hooks" {
		return "hooks"
	}
	return "rules"
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

func defaultPipelineConfig() map[string]any {
	return map[string]any{"extra_args": []any{}, "env": map[string]any{}}
}

func validatePipelineConfig(w http.ResponseWriter, value map[string]any) bool {
	if len(value) != 2 {
		writeValidation(w, "body", "extra_args and env required")
		return false
	}
	for key := range value {
		if key != "extra_args" && key != "env" {
			writeValidation(w, key, "unknown field")
			return false
		}
	}
	raw, ok := value["extra_args"].([]any)
	if !ok {
		writeValidation(w, "extra_args", "required")
		return false
	}
	if len(raw) > 50 {
		writeValidation(w, "extra_args", "too many")
		return false
	}
	reserved := []string{"--src", "--out", "--build-id", "--commit-sha", "--build-at", "--cache-dir", "--version", "--help"}
	for _, item := range raw {
		arg, ok := item.(string)
		if !ok || arg == "" || strings.ContainsAny(arg, "\x00\n\r") {
			writeValidation(w, "extra_args", "invalid value")
			return false
		}
		for _, option := range reserved {
			if arg == option || strings.HasPrefix(arg, option+"=") {
				writeValidation(w, "extra_args", "reserved arg")
				return false
			}
		}
	}
	env, ok := value["env"].(map[string]any)
	if !ok {
		writeValidation(w, "env", "required")
		return false
	}
	if len(env) > 100 {
		writeValidation(w, "env", "too many")
		return false
	}
	for key, rawValue := range env {
		stringValue, ok := rawValue.(string)
		if !ok || !validPipelineEnvKey(key) || strings.HasPrefix(key, "ADLAIRE_CI_") || strings.ContainsAny(stringValue, "\x00\n\r") || len(stringValue) > 4096 {
			writeValidation(w, "env", "invalid value")
			return false
		}
		switch key {
		case "PATH", "HOME", "SHELL", "USER", "GITHUB_TOKEN", "ADLAIRE_TOKEN", "ADLAIRE_CHANGED_TARGETS":
			writeValidation(w, "env", "reserved key")
			return false
		}
	}
	return true
}

func validPipelineEnvKey(key string) bool {
	if key == "" || len(key) > 64 {
		return false
	}
	for i, r := range key {
		if i == 0 {
			if r != '_' && (r < 'A' || r > 'Z') {
				return false
			}
			continue
		}
		if r != '_' && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
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
	if len(raw) > 9 {
		writeValidation(w, "widgets", "too many")
		return false
	}
	allowed := map[string]bool{"status": true, "stats": true, "schedule": true, "alerts": true, "disk": true, "rate_limit": true, "snapshots": true, "maintenance": true, "queue": true}
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

func buildTrendSummary(samples []apiBuildTrendSample) map[string]any {
	if len(samples) == 0 {
		return map[string]any{"count": 0, "avg_seconds": nil, "median_seconds": nil, "p95_seconds": nil, "anomaly_count": 0}
	}
	values := make([]float64, 0, len(samples))
	anomalies := 0
	total := 0.0
	for _, sample := range samples {
		values = append(values, sample.DurationSeconds)
		total += sample.DurationSeconds
		if sample.Anomaly {
			anomalies++
		}
	}
	sort.Float64s(values)
	return map[string]any{
		"count":          len(samples),
		"avg_seconds":    math.Round((total/float64(len(samples)))*100) / 100,
		"median_seconds": percentile(values, 0.5),
		"p95_seconds":    percentile(values, 0.95),
		"anomaly_count":  anomalies,
	}
}

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	if p == 0.5 && len(values)%2 == 0 {
		mid := len(values) / 2
		return math.Round(((values[mid-1]+values[mid])/2)*100) / 100
	}
	index := int(math.Ceil(float64(len(values))*p)) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(values) {
		index = len(values) - 1
	}
	return math.Round(values[index]*100) / 100
}

func buildChainConfigsEqual(a, b apiBuildChainConfig) bool {
	aj, errA := json.Marshal(a)
	bj, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aj) == string(bj)
}

func validateBuildChainConfig(w http.ResponseWriter, cfg apiBuildChainConfig) bool {
	if len(cfg.Chains) > 100 {
		writeValidation(w, "chains", "too many items")
		return false
	}
	ids := map[string]apiBuildChainJob{}
	for _, job := range cfg.Chains {
		if !validBuildChainID(job.ID) {
			writeValidation(w, "id", "invalid value")
			return false
		}
		if _, exists := ids[job.ID]; exists {
			writeValidation(w, "id", "duplicate value")
			return false
		}
		if !validateBranchName(w, job.Branch) || !validateTargetFile(w, job.TargetFile) {
			return false
		}
		if job.DependsOn == nil {
			writeValidation(w, "depends_on", "required")
			return false
		}
		if len(job.DependsOn) > 20 {
			writeValidation(w, "depends_on", "too many items")
			return false
		}
		ids[job.ID] = job
	}
	for _, job := range cfg.Chains {
		seen := map[string]bool{}
		for _, dep := range job.DependsOn {
			if seen[dep] || dep == job.ID {
				writeValidation(w, "depends_on", "invalid value")
				return false
			}
			target, ok := ids[dep]
			if !ok || (job.Enabled && !target.Enabled) {
				writeValidation(w, "depends_on", "invalid value")
				return false
			}
			seen[dep] = true
		}
	}
	visiting := map[string]bool{}
	visited := map[string]bool{}
	var visit func(string) bool
	visit = func(id string) bool {
		if visiting[id] {
			return false
		}
		if visited[id] {
			return true
		}
		visiting[id] = true
		for _, dep := range ids[id].DependsOn {
			if !visit(dep) {
				return false
			}
		}
		visiting[id] = false
		visited[id] = true
		return true
	}
	for id := range ids {
		if !visit(id) {
			writeValidation(w, "depends_on", "cycle detected")
			return false
		}
	}
	return true
}

func validBuildChainID(id string) bool {
	if id == "" || len(id) > 64 {
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

func apiActorID(r *http.Request) string {
	if auth, ok := r.Context().Value(apiAuthContextKey{}).(apiAuthInfo); ok && auth.Actor != nil && *auth.Actor != "" {
		return *auth.Actor
	}
	return "admin"
}

func latestApprovalByID(path, id string) (apiApprovalRecord, bool, error) {
	records, err := readLatestApprovals(path)
	if err != nil {
		return apiApprovalRecord{}, false, err
	}
	for _, record := range records {
		if record.ID == id {
			return record, true, nil
		}
	}
	return apiApprovalRecord{}, false, nil
}

func readLatestApprovals(path string) ([]apiApprovalRecord, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []apiApprovalRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	latest := map[string]apiApprovalRecord{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record apiApprovalRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil || !validApprovalRecord(record) {
			log.Printf("APPROVAL_QUEUE_CORRUPT_LINE: path=%s", path)
			continue
		}
		latest[record.ID] = record
	}
	records := make([]apiApprovalRecord, 0, len(latest))
	for _, record := range latest {
		records = append(records, record)
	}
	return records, nil
}

func validApprovalRecord(record apiApprovalRecord) bool {
	switch record.Status {
	case "pending", "approved", "rejected", "expired":
	default:
		return false
	}
	return record.ID != "" && record.Branch != "" && record.SHA != "" && record.Target != "" && record.RequestedTrigger != "" && record.RequestedBy != "" && record.CreatedAt != "" && record.ExpiresAt != ""
}

func appendJSONLine(path string, value any) error {
	return statefile.AppendJSONLine(path, value, 0600)
}

func atomicWriteHistory(path string, records []apiHistoryRecord, mode os.FileMode) error {
	values := make([]any, 0, len(records))
	for _, item := range records {
		values = append(values, item)
	}
	return statefile.WriteJSONLinesAtomic(path, values, mode)
}

func readJSONLines(path string) []map[string]any {
	records, err := statefile.ReadJSONLines(path)
	if err != nil {
		return []map[string]any{}
	}
	return records
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func normalizeHealthBuildStatus(status string) string {
	switch {
	case status == "":
		return "none"
	case status == "success" || status == "success_deploy_pending":
		return "success"
	case status == "cancelled":
		return "cancelled"
	case strings.HasPrefix(status, "skipped_") || status == "circuit_open" || status == "config_recovered" || status == "config_error" || status == "lock_skipped":
		return "skipped"
	case strings.HasPrefix(status, "failure_") || status == "hook_error":
		return "failure"
	default:
		return status
	}
}

func readJSONArrayCount(path string) (int, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	var values []json.RawMessage
	if err := json.Unmarshal(data, &values); err != nil {
		return 0, err
	}
	return len(values), nil
}

func healthRunnerStale(path string, now time.Time) bool {
	status, ok, check := readBuildStatusForHealth(path)
	if !ok || check != "" || !status.Running || status.LastStartedAt == nil || *status.LastStartedAt == "" {
		return false
	}
	startedAt, err := time.Parse(apiTimeLayout, *status.LastStartedAt)
	if err != nil {
		return false
	}
	return now.Sub(startedAt.UTC()) > 24*time.Hour
}

func patConfigured(path string) (bool, error) {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !info.Mode().IsRegular() {
		return false, nil
	}
	token, err := readSecretText(path)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(token) != "", nil
}

func expiresInDays(expiresAt *string, now time.Time) (*int, error) {
	if expiresAt == nil {
		return nil, nil
	}
	expiresDate, err := time.ParseInLocation("2006-01-02", *expiresAt, time.UTC)
	if err != nil {
		return nil, err
	}
	nowDate := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), 0, 0, 0, 0, time.UTC)
	days := int(expiresDate.Sub(nowDate).Hours() / 24)
	return &days, nil
}

func githubScopes(header string) []string {
	seen := map[string]bool{}
	scopes := []string{}
	for _, raw := range strings.Split(header, ",") {
		scope := strings.TrimSpace(raw)
		if scope == "" || seen[scope] {
			continue
		}
		seen[scope] = true
		scopes = append(scopes, scope)
	}
	sort.Strings(scopes)
	return scopes
}

func readJSONIfExists(path string, out any) error {
	return statefile.ReadJSONIfExists(path, out)
}

func readJSONFile(path string, out any) error {
	return statefile.ReadJSON(path, out)
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
	return statefile.WriteJSONAtomic(path, value, mode)
}

func atomicWriteJSONLocked(path string, value any, mode os.FileMode) error {
	return statefile.WriteJSONAtomicLocked(path, value, mode)
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
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "Payload too large")
			return false
		}
		writeError(w, http.StatusBadRequest, "Request body is not allowed")
		return false
	}
	if len(bytes.TrimSpace(data)) > 0 {
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
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "Payload too large")
			return false
		}
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return false
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return false
	}
	return true
}

func decodeRawObjectBody(w http.ResponseWriter, r *http.Request, required bool) (map[string]json.RawMessage, []string, bool) {
	if r.Body == nil || r.Body == http.NoBody {
		if required {
			writeError(w, http.StatusBadRequest, "Invalid JSON")
			return nil, nil, false
		}
		return map[string]json.RawMessage{}, nil, true
	}
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "Payload too large")
			return nil, nil, false
		}
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return nil, nil, false
	}
	var raw map[string]json.RawMessage
	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&raw); err != nil || raw == nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return nil, nil, false
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return nil, nil, false
	}
	order, ok := jsonObjectKeyOrder(data)
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return nil, nil, false
	}
	return raw, order, true
}

func jsonObjectKeyOrder(data []byte) ([]string, bool) {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return nil, false
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return nil, false
	}
	keys := []string{}
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return nil, false
		}
		key, ok := tok.(string)
		if !ok {
			return nil, false
		}
		keys = append(keys, key)
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return nil, false
		}
	}
	tok, err = dec.Token()
	if err != nil {
		return nil, false
	}
	delim, ok = tok.(json.Delim)
	return keys, ok && delim == '}'
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		status = http.StatusInternalServerError
		data = []byte(`{"error":"Internal server error"}`)
	}
	data = append(data, '\n')
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(data)
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

func validateQueryKeys(w http.ResponseWriter, r *http.Request, allowed map[string]bool) bool {
	for key, values := range r.URL.Query() {
		if !allowed[key] {
			writeValidation(w, key, "unknown query")
			return false
		}
		if len(values) > 1 {
			writeValidation(w, key, "duplicate value")
			return false
		}
	}
	return true
}

func validateRouteQuery(w http.ResponseWriter, r *http.Request) bool {
	return validateQueryKeys(w, r, apiRouteQueryKeys(r.URL.Path))
}

func parseBoundedInt(w http.ResponseWriter, r *http.Request, key string, def, min, max int) (int, bool) {
	values, ok := r.URL.Query()[key]
	if !ok || len(values) == 0 || values[0] == "" {
		return def, true
	}
	if len(values) != 1 {
		writeValidation(w, key, "duplicate value")
		return 0, false
	}
	raw := values[0]
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

func maintenanceEnabled(path string) (bool, error) {
	state, err := readMaintenanceState(path)
	if err != nil {
		return false, err
	}
	return state.Enabled, nil
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

func nullableStringEqual(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
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

func filterContainsFold(lines []string, q string) []string {
	if q == "" {
		return lines
	}
	q = strings.ToLower(q)
	out := []string{}
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), q) {
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
