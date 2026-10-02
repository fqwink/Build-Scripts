package runner

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"
)

var runnerGitHubAPIBase = "https://api.github.com"
var runnerNow = time.Now
var runnerSleep = time.Sleep
var runnerBinaryVersion = "V.0.0-dev"
var runnerBuilderVersionTimeout = 2 * time.Second

func SetBinaryVersion(version string) {
	runnerBinaryVersion = version
}

type RunnerConfig struct {
	DryRun                     bool
	StateDir                   string
	RepositoryOwner            string
	RepositoryName             string
	PendingFile                string
	APIRetryMax                int
	APIRetryBaseSeconds        int
	BuildTimeoutSeconds        int
	BuildCooldownSeconds       int
	BuildRetryMax              int
	BuildRetryBaseSeconds      int
	HistoryKeepN               int
	SnapshotsKeep              int
	ForceBuildIntervalHours    int
	LogKeepN                   int
	APICircuitBreakerThreshold int
	OutputSizeWarnMB           int
	LogArchiveAfterDays        int
	CommitStatusEnabled        bool
	CommitStatusContext        string
	CommitStatusTargetURL      *string
	BuildTrendKeepCount        int
	WeeklySummaryEnabled       bool
	WeeklySummaryDay           int
	WeeklySummaryHour          int
	TagFilter                  TagFilterConfig
	RemoteBuild                RemoteBuildConfig
	DurationAnomaly            DurationAnomalyConfig
	SchedulePaused             bool
	BuildCacheEnabled          bool
	WatchMode                  string
	DeployParallelism          int
	AllowedHours               *AllowedHoursConfig
	BranchTargets              []BranchTarget
}

type BranchTarget struct {
	Branch           string            `json:"branch"`
	TargetFile       string            `json:"target_file"`
	TargetFiles      []string          `json:"target_files"`
	SHAFile          string            `json:"sha_file"`
	Src              string            `json:"src"`
	Out              string            `json:"out"`
	ApprovalRequired bool              `json:"approval_required"`
	Env              map[string]string `json:"env"`
	DeployTargets    []DeployTarget    `json:"deploy_targets"`
}

type DeployTarget struct {
	ID      string `json:"id"`
	Host    string `json:"host"`
	User    string `json:"user"`
	DestDir string `json:"dest_dir"`
}

type branchConfigFile struct {
	BranchTargets []BranchTarget `json:"branch_targets"`
}

type repoConfigFile struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type serverConfigFile struct {
	LogMaxLines             *int                `json:"log_max_lines"`
	HistoryMaxCount         *int                `json:"history_max_count"`
	HistoryRetention        json.RawMessage     `json:"history_retention"`
	BuildTimeoutSeconds     *int                `json:"build_timeout_seconds"`
	LogRetentionDays        *int                `json:"log_retention_days"`
	LogLevel                *string             `json:"log_level"`
	PATExpiresAt            *string             `json:"pat_expires_at"`
	BuildCooldownSeconds    *int                `json:"build_cooldown_seconds"`
	SnapshotsKeep           *int                `json:"snapshots_keep"`
	QueueMaxSize            *int                `json:"queue_max_size"`
	BuildRetryMax           *int                `json:"build_retry_max"`
	BuildRetryBaseSeconds   *int                `json:"build_retry_base_seconds"`
	CommitStatusEnabled     *bool               `json:"commit_status_enabled"`
	CommitStatusContext     *string             `json:"commit_status_context"`
	CommitStatusTargetURL   *string             `json:"commit_status_target_url"`
	LogArchiveAfterDays     *int                `json:"log_archive_after_days"`
	BuildTrendKeepCount     *int                `json:"build_trend_keep_count"`
	DurationAnomaly         json.RawMessage     `json:"duration_anomaly"`
	WatchMode               *string             `json:"watch_mode"`
	TagFilter               json.RawMessage     `json:"tag_filter"`
	DeployParallelism       *int                `json:"deploy_parallelism"`
	RemoteBuild             json.RawMessage     `json:"remote_build"`
	ApprovalTimeoutSeconds  *int                `json:"approval_timeout_seconds"`
	ForceBuildIntervalHours *int                `json:"force_build_interval_hours"`
	ScheduleIntervalSeconds *int                `json:"schedule_interval_seconds"`
	SchedulePaused          *bool               `json:"schedule_paused"`
	AllowedHours            *AllowedHoursConfig `json:"allowed_hours"`
	SessionTimeoutSeconds   *int                `json:"session_timeout_seconds"`
	APIRateLimit            json.RawMessage     `json:"api_rate_limit"`
	BuildCacheEnabled       *bool               `json:"build_cache_enabled"`
}

type AllowedHoursConfig struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type TagFilterConfig struct {
	Enabled  bool     `json:"enabled"`
	Patterns []string `json:"patterns"`
}

type RemoteBuildConfig struct {
	Enabled      bool     `json:"enabled"`
	Host         *string  `json:"host"`
	User         *string  `json:"user"`
	WorkDir      *string  `json:"work_dir"`
	CommandArgs  []string `json:"command_args"`
	ArtifactPath *string  `json:"artifact_path"`
}

type DurationAnomalyConfig struct {
	Enabled       bool    `json:"enabled"`
	MinSamples    int     `json:"min_samples"`
	AvgMultiplier float64 `json:"avg_multiplier"`
	P95Multiplier float64 `json:"p95_multiplier"`
}

type shaCache struct {
	SHA string `json:"sha"`
}

type buildState struct {
	Running                 bool             `json:"running"`
	CurrentBuildID          *string          `json:"current_build_id"`
	ActiveQueueEntry        map[string]any   `json:"active_queue_entry"`
	Queued                  []map[string]any `json:"queued"`
	LastStartedAt           *string          `json:"last_started_at"`
	LastFinishedAt          *string          `json:"last_finished_at"`
	WeeklySummaryLastSentAt *string          `json:"weekly_summary_last_sent_at"`
	WeeklySummarySentDate   *string          `json:"weekly_summary_sent_date"`
}

type queueEntry struct {
	ID           string         `json:"id"`
	Trigger      string         `json:"trigger"`
	QueuedAt     string         `json:"queued_at"`
	RequestedBy  string         `json:"requested_by"`
	Priority     string         `json:"priority"`
	CreatedSeq   int64          `json:"created_seq"`
	Payload      map[string]any `json:"payload"`
	OriginalData map[string]any `json:"-"`
}

type queueRunContext struct {
	Entry      *queueEntry
	Trigger    string
	Force      bool
	Branch     string
	TargetFile string
	CommitSHA  string
}

type buildChainConfigFile struct {
	Chains []buildChainJob `json:"chains"`
}

type buildChainJob struct {
	ID         string   `json:"id"`
	Branch     string   `json:"branch"`
	TargetFile string   `json:"target_file"`
	DependsOn  []string `json:"depends_on"`
	Required   bool     `json:"required"`
	Enabled    bool     `json:"enabled"`
}

type chainLog struct {
	ChainRunID string   `json:"chain_run_id"`
	JobID      string   `json:"job_id"`
	DependsOn  []string `json:"depends_on"`
	ChainIndex int      `json:"chain_index"`
}

type chainSummaryLog struct {
	ChainRunID   string `json:"chain_run_id"`
	TotalJobs    int    `json:"total_jobs"`
	SuccessCount int    `json:"success_count"`
	FailureCount int    `json:"failure_count"`
	SkippedCount int    `json:"skipped_count"`
}

type commitInfo struct {
	SHA     *string `json:"sha"`
	Message *string `json:"message"`
	Author  *string `json:"author"`
	Date    *string `json:"date"`
}

type pipelineLog struct {
	ExitCode        *int            `json:"exit_code"`
	Stdout          string          `json:"stdout"`
	Stderr          string          `json:"stderr"`
	StdoutTruncated bool            `json:"stdout_truncated"`
	StderrTruncated bool            `json:"stderr_truncated"`
	TargetStatus    string          `json:"-"`
	ErrorMessage    string          `json:"-"`
	EnvKeys         []string        `json:"-"`
	RemoteBuild     *remoteBuildLog `json:"-"`
}

type pipelineConfigFile struct {
	ExtraArgs []string          `json:"extra_args"`
	Env       map[string]string `json:"env"`
}

type runnerDependencyManifest struct {
	Version          int                             `json:"version"`
	Builder          runnerDependencyManifestBuilder `json:"builder"`
	Pages            map[string]runnerDependencyPage `json:"pages"`
	GeneratedOutputs map[string]string               `json:"generated_outputs"`
}

type runnerDependencyManifestBuilder struct {
	Name        string `json:"name"`
	SpecSection string `json:"spec_section"`
	ConfigHash  string `json:"config_hash"`
}

type runnerDependencyPage struct {
	OutputPath       string            `json:"output_path"`
	Dependencies     []string          `json:"dependencies"`
	DependencySHA256 map[string]string `json:"dependency_sha256"`
	ConfigHash       string            `json:"config_hash"`
}

type hooksFile struct {
	Hooks []runnerHook `json:"hooks"`
}

type runnerHook struct {
	ID             string   `json:"id"`
	Phase          string   `json:"phase"`
	CommandArgs    []string `json:"command_args"`
	Enabled        bool     `json:"enabled"`
	AbortOnFailure bool     `json:"abort_on_failure"`
	TimeoutSeconds int      `json:"timeout_seconds"`
	CreatedAt      string   `json:"created_at"`
}

type hookRunLog struct {
	HookID          string `json:"hook_id"`
	BuildID         string `json:"build_id"`
	Phase           string `json:"phase"`
	Status          string `json:"status"`
	StartedAt       string `json:"started_at"`
	FinishedAt      string `json:"finished_at"`
	DurationSeconds int64  `json:"duration_seconds"`
	Stdout          string `json:"stdout"`
	Stderr          string `json:"stderr"`
	ExitCode        *int   `json:"exit_code"`
	TimedOut        bool   `json:"timed_out"`
	Truncated       bool   `json:"truncated"`
}

type runnerReport struct {
	Pages           int    `json:"pages"`
	Headings        int    `json:"headings"`
	TablesCount     int    `json:"tables_count"`
	CodeBlocksCount int    `json:"code_blocks_count"`
	WarningsCount   int    `json:"warnings_count"`
	SizeWarn        bool   `json:"size_warn"`
	BrokenLinks     int    `json:"broken_links"`
	HeadingSkips    int    `json:"heading_skips"`
	ReadingTime     int    `json:"reading_time"`
	Theme           string `json:"theme"`
	BuildID         string `json:"build_id"`
	CommitSHA       string `json:"commit_sha"`
	BuildAt         string `json:"build_at"`
}

type deployLog struct {
	TargetID         string  `json:"target_id"`
	Host             string  `json:"host"`
	User             string  `json:"user"`
	DestDir          string  `json:"dest_dir"`
	Status           string  `json:"status"`
	TransferVerified bool    `json:"transfer_verified"`
	FilesTotal       int     `json:"files_total"`
	FilesUploaded    int     `json:"files_uploaded"`
	FilesSkipped     int     `json:"files_skipped"`
	BytesUploaded    int64   `json:"bytes_uploaded"`
	Error            *string `json:"error"`
	StartedAt        string  `json:"-"`
	FinishedAt       string  `json:"-"`
}

type buildLog struct {
	ID               string            `json:"id"`
	Status           string            `json:"status"`
	Branch           string            `json:"branch"`
	TargetFile       string            `json:"target_file"`
	TargetFiles      []string          `json:"target_files"`
	ChangedTargets   []string          `json:"changed_targets"`
	MatchedTags      []string          `json:"matched_tags"`
	TargetStatus     string            `json:"target_status"`
	Trigger          string            `json:"trigger"`
	TriggerActor     *string           `json:"trigger_actor"`
	StartedAt        string            `json:"started_at"`
	FinishedAt       string            `json:"finished_at"`
	DurationSeconds  int64             `json:"duration_seconds"`
	Commit           commitInfo        `json:"commit"`
	CommitSHA        *string           `json:"commit_sha"`
	CommitMessage    *string           `json:"commit_message"`
	CommitAuthor     *string           `json:"commit_author"`
	CommitAt         *string           `json:"commit_at"`
	BlobSHA          *string           `json:"blob_sha"`
	PreviousBlobSHA  string            `json:"previous_blob_sha"`
	Pipeline         pipelineLog       `json:"pipeline"`
	SkippedHookIDs   []string          `json:"skipped_hook_ids"`
	Attempts         []pipelineLog     `json:"attempts"`
	RetryCount       int               `json:"retry_count"`
	RemoteBuild      *remoteBuildLog   `json:"remote_build"`
	Chain            any               `json:"chain"`
	ChainSummary     any               `json:"chain_summary"`
	Report           *runnerReport     `json:"report"`
	Warnings         []string          `json:"warnings"`
	Deploy           []deployLog       `json:"deploy"`
	TargetResults    []targetResultLog `json:"target_results"`
	SnapshotID       *string           `json:"snapshot_id"`
	RollbackFrom     *string           `json:"rollback_from"`
	TransferVerified *bool             `json:"transfer_verified"`
	OutputSHA256     *string           `json:"output_sha256"`
	OutputSizeBytes  *int64            `json:"output_size_bytes"`
	SizeWarn         bool              `json:"size_warn"`
	FailureCategory  *string           `json:"failure_category"`
	FailureEvidence  []failureEvidence `json:"failure_evidence"`
	CommitStatus     *string           `json:"commit_status_state"`
	CommitStatusLog  *commitStatusLog  `json:"commit_status"`
	BuildMeta        *buildMetaLog     `json:"build_meta"`
	Environment      buildEnvLog       `json:"environment"`
	Error            *string           `json:"error"`
}

type remoteBuildLog struct {
	Host              string  `json:"host"`
	User              string  `json:"user"`
	WorkDirBasename   string  `json:"work_dir_basename"`
	CommandName       string  `json:"command_name"`
	ExitCode          *int    `json:"exit_code"`
	DurationSeconds   int64   `json:"duration_seconds"`
	ArtifactSizeBytes *int64  `json:"artifact_size_bytes"`
	ManifestFileCount *int    `json:"manifest_file_count"`
	Status            string  `json:"status"`
	Error             *string `json:"error"`
}

type commitStatusLog struct {
	Enabled    bool    `json:"enabled"`
	State      *string `json:"state"`
	Context    string  `json:"context"`
	TargetURL  *string `json:"target_url"`
	SentAt     *string `json:"sent_at"`
	HTTPStatus *int    `json:"http_status"`
	Error      *string `json:"error"`
}

type buildMetaLog struct {
	BuildID   string `json:"build_id"`
	CommitSHA string `json:"commit_sha"`
	BuildAt   string `json:"build_at"`
}

type targetResultLog struct {
	TargetID   string  `json:"target_id"`
	Status     string  `json:"status"`
	StartedAt  string  `json:"started_at"`
	FinishedAt string  `json:"finished_at"`
	ErrorCode  *string `json:"error_code"`
	Error      *string `json:"error"`
}

type failureEvidence struct {
	Source  string `json:"source"`
	Code    string `json:"code"`
	Message string `json:"message"`
	At      string `json:"at"`
}

type buildEnvLog struct {
	OS             string   `json:"os"`
	Arch           string   `json:"arch"`
	GoVersion      string   `json:"go_version"`
	RunnerVersion  string   `json:"runner_version"`
	BuilderVersion string   `json:"builder_version"`
	Hostname       string   `json:"hostname"`
	PID            int      `json:"pid"`
	WatchMode      string   `json:"watch_mode"`
	CacheEnabled   bool     `json:"cache_enabled"`
	RemoteBuild    bool     `json:"remote_build"`
	StateDir       string   `json:"state_dir"`
	DiskFreeBytes  *int64   `json:"disk_free_bytes"`
	CapturedAt     string   `json:"captured_at"`
	EnvKeys        []string `json:"env_keys"`
}

type historyRecord struct {
	ID              string  `json:"id"`
	Branch          string  `json:"branch"`
	TargetFile      string  `json:"target_file"`
	Status          string  `json:"status"`
	Trigger         string  `json:"trigger"`
	StartedAt       string  `json:"started_at"`
	FinishedAt      string  `json:"finished_at"`
	DurationSeconds int64   `json:"duration_seconds"`
	CommitSHA       *string `json:"commit_sha"`
	BlobSHA         *string `json:"blob_sha"`
	Pages           *int    `json:"pages"`
	Warnings        int     `json:"warnings"`
	SizeWarn        bool    `json:"size_warn"`
	OutputSHA256    *string `json:"output_sha256"`
	OutputSizeBytes *int64  `json:"output_size_bytes"`
	RetryCount      int     `json:"retry_count"`
	FailureCategory *string `json:"failure_category"`
	CommitStatus    *string `json:"commit_status_state"`
	SnapshotID      *string `json:"snapshot_id"`
	RollbackFrom    *string `json:"rollback_from"`
	ChainRunID      *string `json:"chain_run_id"`
}

type pendingTransfer struct {
	BuildID      string  `json:"build_id"`
	Trigger      string  `json:"trigger"`
	SourceKind   string  `json:"source_kind"`
	RollbackFrom *string `json:"rollback_from"`
	SnapshotID   *string `json:"snapshot_id"`
	Branch       string  `json:"branch"`
	TargetID     string  `json:"target_id"`
	Out          *string `json:"out"`
	Host         string  `json:"host"`
	User         string  `json:"user"`
	DestDir      string  `json:"dest_dir"`
	OutputSHA256 string  `json:"output_sha256"`
	FailedAt     string  `json:"failed_at"`
	RetryCount   int     `json:"retry_count"`
	LastError    string  `json:"last_error"`
}

type notifyPendingEntry struct {
	ID          string         `json:"id,omitempty"`
	Event       string         `json:"event"`
	ChannelID   string         `json:"channel_id,omitempty"`
	ChannelType string         `json:"channel_type,omitempty"`
	URL         string         `json:"url"`
	Payload     map[string]any `json:"payload"`
	QueuedAt    string         `json:"queued_at"`
	Attempts    int            `json:"attempts,omitempty"`
	NextAttempt string         `json:"next_attempt_at,omitempty"`
	CreatedAt   string         `json:"created_at,omitempty"`
	PayloadSHA  string         `json:"payload_sha256,omitempty"`
	RetryCount  int            `json:"retry_count"`
	LastError   string         `json:"last_error"`
}

type notifyConfigFile struct {
	Webhooks []notifyWebhook `json:"webhooks"`
	Channels []notifyChannel `json:"channels"`
	On       []string        `json:"on"`
	Summary  notifySummary   `json:"summary"`
}

type notifyWebhook struct {
	URL     string   `json:"url"`
	Label   string   `json:"label"`
	Enabled bool     `json:"enabled"`
	On      []string `json:"on"`
}

type notifyChannel struct {
	ID                   string         `json:"id"`
	Type                 string         `json:"type"`
	Enabled              bool           `json:"enabled"`
	On                   []string       `json:"on"`
	Config               map[string]any `json:"config"`
	RetryCount           *int           `json:"retry_count,omitempty"`
	RetryIntervalSeconds *int           `json:"retry_interval_seconds,omitempty"`
}

type notifySummary struct {
	Enabled   bool   `json:"enabled"`
	Interval  string `json:"interval"`
	Hour      int    `json:"hour"`
	DayOfWeek int    `json:"day_of_week"`
}

type smtpConfigFile struct {
	Host    *string  `json:"host"`
	Port    int      `json:"port"`
	User    *string  `json:"user"`
	TLS     bool     `json:"tls"`
	From    *string  `json:"from"`
	To      []string `json:"to"`
	On      []string `json:"on"`
	Enabled bool     `json:"enabled"`
}

type notifyDeliveryResult struct {
	Result        string
	ErrorCode     *string
	Error         *string
	HTTPStatus    *int
	Retryable     bool
	PendingError  string
	AttemptNumber int
}

type buildCircuitState struct {
	Open                bool    `json:"open"`
	ConsecutiveFailures int     `json:"consecutive_failures"`
	OpenedAt            *string `json:"opened_at"`
	LastFailureAt       *string `json:"last_failure_at"`
	LastError           *string `json:"last_error"`
}

type buildStatusSummary struct {
	SchemaVersion              int     `json:"schema_version"`
	UpdatedAt                  *string `json:"updated_at"`
	Status                     string  `json:"status"`
	Running                    bool    `json:"running"`
	CurrentBuildID             *string `json:"current_build_id"`
	LastBuildID                *string `json:"last_build_id"`
	LastTrigger                *string `json:"last_trigger"`
	LastTargetStatus           *string `json:"last_target_status"`
	LastBranch                 *string `json:"last_branch"`
	LastTargetFile             *string `json:"last_target_file"`
	LastBlobSHA                *string `json:"last_blob_sha"`
	LastCommitSHA              *string `json:"last_commit_sha"`
	LastStartedAt              *string `json:"last_started_at"`
	LastFinishedAt             *string `json:"last_finished_at"`
	LastDurationSeconds        *int64  `json:"last_duration_seconds"`
	LastError                  *string `json:"last_error"`
	LastDeployAt               *string `json:"last_deploy_at"`
	LastDeployStatus           *string `json:"last_deploy_status"`
	PendingTransfersCount      int     `json:"pending_transfers_count"`
	NotifyPendingCount         int     `json:"notify_pending_count"`
	CircuitOpen                bool    `json:"circuit_open"`
	CircuitConsecutiveFailures int     `json:"circuit_consecutive_failures"`
	OutputSHA256               *string `json:"output_sha256"`
	SizeWarn                   bool    `json:"size_warn"`
}

type buildTrendsFile struct {
	SchemaVersion int                `json:"schema_version"`
	Samples       []buildTrendSample `json:"samples"`
	Summary       buildTrendsSummary `json:"summary"`
}

type buildTrendSample struct {
	BuildID         string `json:"build_id"`
	FinishedAt      string `json:"finished_at"`
	Branch          string `json:"branch"`
	Trigger         string `json:"trigger"`
	DurationSeconds int64  `json:"duration_seconds"`
	Status          string `json:"status"`
	TargetStatus    string `json:"target_status"`
	Anomaly         bool   `json:"anomaly"`
}

type buildTrendsSummary struct {
	Count         int      `json:"count"`
	AvgSeconds    *float64 `json:"avg_seconds"`
	MedianSeconds *float64 `json:"median_seconds"`
	P95Seconds    *float64 `json:"p95_seconds"`
	AnomalyCount  int      `json:"anomaly_count"`
}

type gitTreeResponse struct {
	Tree []struct {
		Path string `json:"path"`
		Type string `json:"type"`
		SHA  string `json:"sha"`
	} `json:"tree"`
}

type gitBlobResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

type resolvedRunnerTarget struct {
	Digest         string
	ChangedTargets []string
	Contents       map[string][]byte
	Local          bool
}

type gitCommitResponse []struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Name string `json:"name"`
			Date string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

type gitMatchingRefsResponse []struct {
	Ref    string `json:"ref"`
	Object struct {
		SHA  string `json:"sha"`
		Type string `json:"type"`
	} `json:"object"`
}

func RunRunner(args []string, stdout, stderr io.Writer) int {
	cfg, handled, err := parseRunnerArgs(args, stdout)
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
	return executeRunner(cfg, stdout, stderr)
}

func parseRunnerArgs(args []string, stdout io.Writer) (RunnerConfig, bool, error) {
	cfg := defaultRunnerConfig("/opt/adlaire-builder")
	for _, arg := range args {
		if arg == "--help" {
			fmt.Fprintln(stdout, "Usage: adlaire-ci-runner [--state-dir path] [--once] [--dry-run] [--version] [--help]")
			return cfg, true, nil
		}
	}
	for _, arg := range args {
		if arg == "--version" {
			fmt.Fprintf(stdout, "adlaire-ci-runner %s go=%s\n", runnerBinaryVersion, runtime.Version())
			return cfg, true, nil
		}
	}
	for _, arg := range args {
		if !validCLIArgToken(arg) {
			return cfg, false, exitError{Code: 2, Msg: "invalid command line token"}
		}
	}
	stateDir := cfg.StateDir
	dryRun := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--once":
			continue
		case "--dry-run":
			dryRun = true
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
		name := arg
		if name != "--state-dir" {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + name}
		}
		if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
			return cfg, false, exitError{Code: 2, Msg: "missing value: " + name}
		}
		i++
		stateDir = args[i]
	}
	cfg = defaultRunnerConfig(stateDir)
	cfg.DryRun = dryRun
	if cfg.StateDir == "" {
		return cfg, false, exitError{Code: 2, Msg: "state directory must not be empty"}
	}
	if !filepath.IsAbs(cfg.StateDir) {
		return cfg, false, exitError{Code: 2, Msg: "state directory must be absolute: " + cfg.StateDir}
	}
	info, err := os.Stat(cfg.StateDir)
	if err != nil {
		return cfg, false, exitError{Code: 2, Msg: "state directory not found: " + cfg.StateDir}
	}
	if !info.IsDir() {
		return cfg, false, exitError{Code: 2, Msg: "state path is not directory: " + cfg.StateDir}
	}
	return cfg, false, nil
}

func validCLIArgToken(arg string) bool {
	if !utf8.ValidString(arg) {
		return false
	}
	for _, r := range arg {
		if r == 0 || r == '\n' || r == '\r' || r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

func defaultRunnerConfig(stateDir string) RunnerConfig {
	return RunnerConfig{
		DryRun:                     false,
		StateDir:                   stateDir,
		RepositoryOwner:            "fqwink",
		RepositoryName:             "Build-Scripts",
		PendingFile:                filepath.Join(stateDir, ".pending_transfers"),
		APIRetryMax:                5,
		APIRetryBaseSeconds:        1,
		BuildTimeoutSeconds:        300,
		BuildCooldownSeconds:       60,
		BuildRetryMax:              0,
		BuildRetryBaseSeconds:      5,
		HistoryKeepN:               10,
		SnapshotsKeep:              5,
		LogKeepN:                   50,
		APICircuitBreakerThreshold: 3,
		OutputSizeWarnMB:           5,
		LogArchiveAfterDays:        0,
		CommitStatusEnabled:        false,
		CommitStatusContext:        "Adlaire CI",
		BuildTrendKeepCount:        1000,
		WeeklySummaryEnabled:       true,
		WeeklySummaryDay:           0,
		WeeklySummaryHour:          9,
		TagFilter:                  TagFilterConfig{Enabled: false, Patterns: []string{}},
		RemoteBuild:                RemoteBuildConfig{Enabled: false, CommandArgs: []string{}},
		DurationAnomaly:            DurationAnomalyConfig{Enabled: false, MinSamples: 20, AvgMultiplier: 2.0, P95Multiplier: 1.5},
		WatchMode:                  "github",
		DeployParallelism:          1,
		BranchTargets: []BranchTarget{{
			Branch:     "main",
			TargetFile: "docs",
			TargetFiles: []string{
				"docs",
			},
			SHAFile: filepath.Join(stateDir, ".last_sha"),
			Src:     filepath.Join(stateDir, "repo", "docs"),
			Out:     filepath.Join(stateDir, "dist", "site"),
			Env:     map[string]string{},
		}},
	}
}

func executeRunner(cfg RunnerConfig, stdout, stderr io.Writer) int {
	logger := slog.New(slog.NewTextHandler(stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if cfg.DryRun {
		code := executeRunnerDryRun(cfg, stdout)
		_ = stderr
		return code
	}
	lockPath := filepath.Join(cfg.StateDir, ".build_lock")
	release, locked, code := acquireRunnerLock(lockPath, logger)
	if !locked {
		return code
	}
	defer release()
	if err := loadRunnerConfig(&cfg, logger); err != nil {
		logger.Error(err.Error())
		_ = writeBuildStatusError(cfg, "config_error", "startup_config_integrity", err.Error())
		return 2
	}
	token := ""
	if cfg.WatchMode == "github" {
		var err error
		token, err = readRunnerToken(filepath.Join(cfg.StateDir, ".github_token"))
		if err != nil {
			logger.Error(err.Error())
			_ = writeBuildStatusError(cfg, "config_error", "startup_config_integrity", err.Error())
			return 2
		}
	}
	if err := ensureRunnerStateDirectories(cfg.StateDir); err != nil {
		logger.Error("STATE_DIR_INVALID: path=" + cfg.StateDir)
		_ = writeBuildStatusError(cfg, "config_error", "startup_config_integrity", err.Error())
		return 2
	}
	if err := repairCorruptJSONArray(filepath.Join(cfg.StateDir, ".notify_pending"), logger); err != nil {
		logger.Error(err.Error())
		return 1
	}
	if err := retryPendingTransfers(&cfg, logger); err != nil {
		logger.Error("PENDING_TRANSFER_RETRY_FAILED: " + err.Error())
	}
	if err := retryNotifyPending(filepath.Join(cfg.StateDir, ".notify_pending"), logger); err != nil {
		logger.Error("NOTIFY_PENDING_RETRY_FAILED: " + err.Error())
	}
	if err := sendWeeklySummaryIfDue(cfg, logger); err != nil {
		logger.Warn("WEEKLY_SUMMARY_FAILED: " + err.Error())
	}
	if circuitOpen(cfg.StateDir) {
		logger.Error("CIRCUIT_OPEN: polling skipped")
		_ = writeBuildStatusSkip(cfg, "circuit_open", "polling", nil)
		_ = clearRunningBuildState(cfg.StateDir)
		return 0
	}
	queueCtx, err := readQueueRunContext(cfg.StateDir)
	if err != nil {
		logger.Error("QUEUE_READ_FAILED: " + err.Error())
		_ = writeBuildStatusError(cfg, "failure_decode", "queue", err.Error())
		return 1
	}
	hasQueue := queueCtx.Entry != nil
	if cfg.SchedulePaused && !hasQueue {
		_ = writeBuildStatusSkip(cfg, "skipped_schedule_paused", "polling", nil)
		return 0
	}
	if cooldownActive(cfg, logger) && (!hasQueue || !queueCtx.Force) {
		_ = writeBuildStatusSkip(cfg, "skipped_cooldown", "polling", nil)
		return 0
	}
	exit := 0
	chainHandled := false
	if !hasQueue {
		var chainExit int
		chainHandled, chainExit = executeRunnerChain(cfg, token, logger)
		if chainHandled {
			exit = chainExit
		}
	}
	matchedQueueTarget := false
	for i, target := range cfg.BranchTargets {
		if chainHandled {
			break
		}
		if hasQueue && !queueMatchesTarget(queueCtx, target) {
			continue
		}
		matchedQueueTarget = true
		buildID, err := runnerNextBuildID(cfg.StateDir, runnerNow().UTC())
		if err != nil {
			logger.Error("BUILD_ID_ALLOC_FAILED: " + err.Error())
			if exit < 1 {
				exit = 1
			}
			continue
		}
		if status := processRunnerTarget(cfg, i, target, token, buildID, logger, queueCtx); status > exit {
			exit = status
		}
		if hasQueue {
			break
		}
	}
	if !chainHandled && hasQueue && !matchedQueueTarget {
		logger.Error("QUEUE_TARGET_NOT_FOUND: id=" + queueCtx.Entry.ID)
		_ = writeBuildStatusError(cfg, "failure_api", queueCtx.Trigger, "queue target not found")
		exit = 1
	}
	if !chainHandled && hasQueue && matchedQueueTarget && exit == 0 {
		if err := clearActiveQueueEntry(cfg.StateDir, queueCtx.Entry.ID); err != nil {
			logger.Error("QUEUE_FINALIZE_FAILED: " + err.Error())
			exit = 1
		}
	}
	finished := runnerNow().UTC().Format(time.RFC3339)
	state, err := runnerReadBuildState(cfg.StateDir)
	if err != nil {
		state = runnerDefaultBuildState()
	}
	state.Running = false
	state.CurrentBuildID = nil
	state.LastFinishedAt = &finished
	if err := runnerAtomicWriteJSON(filepath.Join(cfg.StateDir, ".build_state"), state, 0600); err != nil {
		logger.Error("STATE_FINISH_FAILED: " + err.Error())
		if exit < 1 {
			exit = 1
		}
	}
	cleanupBuildLogs(filepath.Join(cfg.StateDir, ".build_logs"), cfg.LogKeepN, cfg.LogArchiveAfterDays, logger)
	_ = stderr
	return exit
}

func executeRunnerChain(cfg RunnerConfig, token string, logger *slog.Logger) (bool, int) {
	chainCfg, order, err := readBuildChainConfig(filepath.Join(cfg.StateDir, ".build_chain_config"))
	if errors.Is(err, os.ErrNotExist) || len(order) == 0 {
		return false, 0
	}
	if err != nil {
		logger.Warn("BUILD_CHAIN_DISABLED: " + err.Error())
		return false, 0
	}
	chainRunID, err := nextChainRunID(cfg.StateDir, runnerNow().UTC())
	if err != nil {
		logger.Error("BUILD_CHAIN_ID_FAILED: " + err.Error())
		return true, 1
	}
	targetByJob := map[string]struct {
		index  int
		target BranchTarget
	}{}
	for _, job := range chainCfg.Chains {
		if !job.Enabled {
			continue
		}
		for i, target := range cfg.BranchTargets {
			if target.Branch == job.Branch && target.TargetFile == job.TargetFile {
				targetByJob[job.ID] = struct {
					index  int
					target BranchTarget
				}{index: i, target: target}
				break
			}
		}
	}
	results := map[string]string{}
	exit := 0
	lastBuildLogPath := ""
	for chainIndex, job := range order {
		if shouldSkipChainJob(job, results) {
			buildID, err := runnerNextBuildID(cfg.StateDir, runnerNow().UTC())
			if err != nil {
				logger.Error("BUILD_ID_ALLOC_FAILED: " + err.Error())
				exit = 1
				break
			}
			if err := appendSkippedDependencyHistory(cfg.StateDir, buildID, job, chainRunID); err != nil {
				logger.Error("BUILD_CHAIN_SKIP_HISTORY_FAILED: " + err.Error())
				exit = 1
				break
			}
			results[job.ID] = "skipped_dependency_failed"
			continue
		}
		targetRef, ok := targetByJob[job.ID]
		if !ok {
			results[job.ID] = "failure_api"
			logger.Error("BUILD_CHAIN_TARGET_NOT_FOUND: job_id=" + job.ID)
			exit = 1
			continue
		}
		buildID, err := runnerNextBuildID(cfg.StateDir, runnerNow().UTC())
		if err != nil {
			logger.Error("BUILD_ID_ALLOC_FAILED: " + err.Error())
			exit = 1
			continue
		}
		meta := &chainLog{
			ChainRunID: chainRunID,
			JobID:      job.ID,
			DependsOn:  append([]string(nil), job.DependsOn...),
			ChainIndex: chainIndex,
		}
		status := processRunnerTargetWithChain(cfg, targetRef.index, targetRef.target, token, buildID, logger, queueRunContext{}, meta, nil)
		if status > exit {
			exit = status
		}
		targetStatus := readBuildLogTargetStatus(cfg.StateDir, buildID)
		if targetStatus == "" && status == 0 {
			targetStatus = "success"
		}
		if targetStatus == "" {
			targetStatus = "failure_build"
		}
		results[job.ID] = targetStatus
		if _, err := os.Stat(filepath.Join(cfg.StateDir, ".build_logs", buildID+".json")); err == nil {
			lastBuildLogPath = filepath.Join(cfg.StateDir, ".build_logs", buildID+".json")
		}
	}
	if lastBuildLogPath != "" {
		summary := chainSummaryFromResults(chainRunID, order, results)
		var log buildLog
		if err := runnerReadJSONFile(lastBuildLogPath, &log); err == nil {
			log.ChainSummary = summary
			if err := writeBuildLog(cfg.StateDir, log); err != nil {
				logger.Error("BUILD_CHAIN_SUMMARY_FAILED: " + err.Error())
				exit = 1
			}
		}
	}
	return true, exit
}

func readBuildChainConfig(path string) (buildChainConfigFile, []buildChainJob, error) {
	var config buildChainConfigFile
	if err := runnerReadJSONFile(path, &config); err != nil {
		return config, nil, err
	}
	jobs := []buildChainJob{}
	for _, job := range config.Chains {
		if job.Enabled {
			jobs = append(jobs, job)
		}
	}
	order, err := topologicalBuildChain(jobs)
	if err != nil {
		return config, nil, err
	}
	return config, order, nil
}

func topologicalBuildChain(jobs []buildChainJob) ([]buildChainJob, error) {
	byID := map[string]buildChainJob{}
	dependents := map[string][]string{}
	remaining := map[string]int{}
	for _, job := range jobs {
		if job.ID == "" {
			return nil, errors.New("chain job id missing")
		}
		if _, exists := byID[job.ID]; exists {
			return nil, fmt.Errorf("chain duplicate id: %s", job.ID)
		}
		byID[job.ID] = job
		remaining[job.ID] = len(job.DependsOn)
		for _, dep := range job.DependsOn {
			dependents[dep] = append(dependents[dep], job.ID)
		}
	}
	for _, job := range jobs {
		for _, dep := range job.DependsOn {
			if _, ok := byID[dep]; !ok {
				return nil, fmt.Errorf("chain dependency missing: %s", dep)
			}
			if dep == job.ID {
				return nil, fmt.Errorf("chain self dependency: %s", dep)
			}
		}
	}
	ready := []string{}
	for _, job := range jobs {
		if remaining[job.ID] == 0 {
			ready = append(ready, job.ID)
		}
	}
	out := []buildChainJob{}
	for len(ready) > 0 {
		id := ready[0]
		ready = ready[1:]
		out = append(out, byID[id])
		for _, depID := range dependents[id] {
			remaining[depID]--
			if remaining[depID] == 0 {
				ready = append(ready, depID)
			}
		}
	}
	if len(out) != len(jobs) {
		return nil, errors.New("chain cycle detected")
	}
	return out, nil
}

func shouldSkipChainJob(job buildChainJob, results map[string]string) bool {
	if !job.Required {
		return false
	}
	for _, dep := range job.DependsOn {
		if !chainStatusAllowsContinuation(results[dep]) {
			return true
		}
	}
	return false
}

func chainStatusAllowsContinuation(status string) bool {
	if status == "skipped_dependency_failed" {
		return false
	}
	return status == "" || status == "success" || status == "success_deploy_pending" || strings.HasPrefix(status, "skipped_")
}

func nextChainRunID(stateDir string, now time.Time) (string, error) {
	base := "chain" + now.UTC().Format("20060102150405")
	for suffix := 0; suffix < 1000; suffix++ {
		id := base
		if suffix > 0 {
			id = fmt.Sprintf("%s-%03d", base, suffix)
		}
		if !chainRunIDExists(stateDir, id) {
			return id, nil
		}
	}
	return "", errors.New("chain id suffix exhausted")
}

func chainRunIDExists(stateDir, id string) bool {
	records, err := readHistoryRecords(stateDir)
	if err != nil {
		return false
	}
	for _, rec := range records {
		if rec.ChainRunID != nil && *rec.ChainRunID == id {
			return true
		}
	}
	return false
}

func readBuildLogTargetStatus(stateDir, buildID string) string {
	var log buildLog
	if err := runnerReadJSONFile(filepath.Join(stateDir, ".build_logs", buildID+".json"), &log); err != nil {
		return ""
	}
	return log.TargetStatus
}

func appendSkippedDependencyHistory(stateDir, buildID string, job buildChainJob, chainRunID string) error {
	now := runnerNow().UTC().Format(time.RFC3339)
	rec := historyRecord{
		ID:              buildID,
		Branch:          job.Branch,
		TargetFile:      job.TargetFile,
		Status:          "skipped_dependency_failed",
		Trigger:         "polling",
		StartedAt:       "",
		FinishedAt:      now,
		DurationSeconds: 0,
		Warnings:        0,
		SizeWarn:        false,
		RetryCount:      0,
		ChainRunID:      &chainRunID,
	}
	return appendRunnerJSONLine(filepath.Join(stateDir, ".build_history"), rec, 0600)
}

func chainSummaryFromResults(chainRunID string, order []buildChainJob, results map[string]string) chainSummaryLog {
	summary := chainSummaryLog{ChainRunID: chainRunID, TotalJobs: len(order)}
	for _, job := range order {
		status, ok := results[job.ID]
		if !ok && shouldSkipChainJob(job, results) {
			status = "skipped_dependency_failed"
		} else if !ok {
			status = "failure_api"
		}
		switch {
		case status == "success" || status == "success_deploy_pending":
			summary.SuccessCount++
		case status == "skipped_dependency_failed":
			summary.SkippedCount++
		default:
			summary.FailureCount++
		}
	}
	return summary
}

func ensureRunnerStateDirectories(stateDir string) error {
	for _, dir := range []string{filepath.Join(stateDir, "repo"), filepath.Join(stateDir, "dist"), filepath.Join(stateDir, ".build_logs")} {
		info, err := os.Stat(dir)
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, 0700); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("state path is not directory: %s", dir)
		}
	}
	return nil
}

func loadRunnerConfig(cfg *RunnerConfig, logger *slog.Logger) error {
	if err := loadRepoConfig(cfg); err != nil {
		return err
	}
	if err := loadServerConfig(cfg); err != nil {
		return err
	}
	path := filepath.Join(cfg.StateDir, ".branch_config")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return validateRunnerConfig(*cfg)
	}
	if err != nil {
		return err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key := range raw {
		if key != "branch_targets" {
			return fmt.Errorf("BRANCH_CONFIG_INVALID: unknown key=%s", key)
		}
	}
	var bc branchConfigFile
	if err := json.Unmarshal(data, &bc); err != nil {
		return err
	}
	cfg.BranchTargets = normalizeBranchTargets(bc.BranchTargets)
	return validateRunnerConfig(*cfg)
}

func loadRepoConfig(cfg *RunnerConfig) error {
	var rc repoConfigFile
	err := runnerReadJSONFile(filepath.Join(cfg.StateDir, ".repo_config"), &rc)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil || rc.Owner == "" || rc.Repo == "" {
		return errors.New("REPO_CONFIG_INVALID")
	}
	cfg.RepositoryOwner = rc.Owner
	cfg.RepositoryName = rc.Repo
	return nil
}

func loadServerConfig(cfg *RunnerConfig) error {
	var raw map[string]json.RawMessage
	path := filepath.Join(cfg.StateDir, ".server_config")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	allowed := map[string]bool{
		"log_max_lines": true, "history_max_count": true, "history_retention": true,
		"build_timeout_seconds": true, "log_retention_days": true, "log_level": true,
		"pat_expires_at": true, "snapshots_keep": true, "queue_max_size": true,
		"build_retry_max": true, "build_retry_base_seconds": true,
		"commit_status_enabled": true, "commit_status_context": true, "commit_status_target_url": true,
		"log_archive_after_days": true, "build_trend_keep_count": true, "duration_anomaly": true,
		"watch_mode": true, "tag_filter": true, "build_cache_enabled": true,
		"deploy_parallelism": true, "remote_build": true, "approval_timeout_seconds": true,
		"force_build_interval_hours": true, "build_cooldown_seconds": true,
		"schedule_interval_seconds": true, "schedule_paused": true, "allowed_hours": true,
		"session_timeout_seconds": true, "api_rate_limit": true,
	}
	for key := range raw {
		if !allowed[key] {
			return fmt.Errorf("SERVER_CONFIG_INVALID: unknown key=%s", key)
		}
	}
	var sc serverConfigFile
	if err := json.Unmarshal(data, &sc); err != nil {
		return err
	}
	if sc.BuildTimeoutSeconds != nil {
		cfg.BuildTimeoutSeconds = *sc.BuildTimeoutSeconds
	}
	if sc.BuildCooldownSeconds != nil {
		cfg.BuildCooldownSeconds = *sc.BuildCooldownSeconds
	}
	if sc.BuildRetryMax != nil {
		cfg.BuildRetryMax = *sc.BuildRetryMax
	}
	if sc.BuildRetryBaseSeconds != nil {
		cfg.BuildRetryBaseSeconds = *sc.BuildRetryBaseSeconds
	}
	if sc.SnapshotsKeep != nil {
		cfg.SnapshotsKeep = *sc.SnapshotsKeep
		cfg.HistoryKeepN = *sc.SnapshotsKeep
	}
	if sc.CommitStatusEnabled != nil {
		cfg.CommitStatusEnabled = *sc.CommitStatusEnabled
	}
	if sc.CommitStatusContext != nil {
		cfg.CommitStatusContext = *sc.CommitStatusContext
	}
	if sc.CommitStatusTargetURL != nil {
		cfg.CommitStatusTargetURL = sc.CommitStatusTargetURL
	}
	if sc.LogArchiveAfterDays != nil {
		cfg.LogArchiveAfterDays = *sc.LogArchiveAfterDays
	}
	if sc.BuildTrendKeepCount != nil {
		cfg.BuildTrendKeepCount = *sc.BuildTrendKeepCount
	}
	if rawJSONPresent(sc.TagFilter) {
		var value TagFilterConfig
		if err := json.Unmarshal(sc.TagFilter, &value); err != nil {
			return fmt.Errorf("SERVER_CONFIG_INVALID: tag_filter")
		}
		if value.Patterns == nil {
			value.Patterns = []string{}
		}
		cfg.TagFilter = value
	}
	if rawJSONPresent(sc.RemoteBuild) {
		var value RemoteBuildConfig
		if err := json.Unmarshal(sc.RemoteBuild, &value); err != nil {
			return fmt.Errorf("SERVER_CONFIG_INVALID: remote_build")
		}
		if value.CommandArgs == nil {
			value.CommandArgs = []string{}
		}
		cfg.RemoteBuild = value
	}
	if rawJSONPresent(sc.DurationAnomaly) {
		var value DurationAnomalyConfig
		if err := json.Unmarshal(sc.DurationAnomaly, &value); err != nil {
			return fmt.Errorf("SERVER_CONFIG_INVALID: duration_anomaly")
		}
		cfg.DurationAnomaly = value
	}
	if sc.WatchMode != nil {
		cfg.WatchMode = *sc.WatchMode
	}
	if sc.DeployParallelism != nil {
		cfg.DeployParallelism = *sc.DeployParallelism
	}
	if sc.ForceBuildIntervalHours != nil {
		cfg.ForceBuildIntervalHours = *sc.ForceBuildIntervalHours
	}
	if sc.SchedulePaused != nil {
		cfg.SchedulePaused = *sc.SchedulePaused
	}
	if sc.AllowedHours != nil {
		cfg.AllowedHours = sc.AllowedHours
	}
	if sc.BuildCacheEnabled != nil {
		cfg.BuildCacheEnabled = *sc.BuildCacheEnabled
	}
	return nil
}

func rawJSONPresent(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	return strings.TrimSpace(string(raw)) != "null"
}

func normalizeBranchTargets(targets []BranchTarget) []BranchTarget {
	out := make([]BranchTarget, 0, len(targets))
	for i, t := range targets {
		if len(t.TargetFiles) == 0 && t.TargetFile != "" {
			t.TargetFiles = []string{t.TargetFile}
		}
		sort.Strings(t.TargetFiles)
		if t.Env == nil {
			t.Env = map[string]string{}
		}
		for j := range t.DeployTargets {
			if t.DeployTargets[j].ID == "" {
				t.DeployTargets[j].ID = fmt.Sprintf("%d-%d", i, j)
			}
		}
		if t.Branch == "" {
			t.Branch = fmt.Sprintf("target-%d", i+1)
		}
		out = append(out, t)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Branch == out[j].Branch {
			return out[i].TargetFile < out[j].TargetFile
		}
		return out[i].Branch < out[j].Branch
	})
	return out
}

func validateRunnerConfig(cfg RunnerConfig) error {
	if cfg.RepositoryOwner == "" || cfg.RepositoryName == "" {
		return errors.New("repository identity is required")
	}
	if cfg.BuildTimeoutSeconds < 1 || cfg.BuildTimeoutSeconds > 86400 {
		return errors.New("invalid build timeout seconds")
	}
	if cfg.BuildRetryMax < 0 || cfg.BuildRetryMax > 10 {
		return errors.New("invalid build retry max")
	}
	if cfg.BuildRetryBaseSeconds < 1 || cfg.BuildRetryBaseSeconds > 3600 {
		return errors.New("invalid build retry base seconds")
	}
	if cfg.SnapshotsKeep < 0 || cfg.SnapshotsKeep > 100 {
		return errors.New("invalid snapshots keep")
	}
	if cfg.CommitStatusContext == "" || len(cfg.CommitStatusContext) > 100 || hasUnsafeControl(cfg.CommitStatusContext) {
		return errors.New("invalid commit status context")
	}
	if cfg.CommitStatusTargetURL != nil && *cfg.CommitStatusTargetURL != "" && !strings.HasPrefix(*cfg.CommitStatusTargetURL, "http://") && !strings.HasPrefix(*cfg.CommitStatusTargetURL, "https://") {
		return errors.New("invalid commit status target url")
	}
	if cfg.BuildTrendKeepCount < 10 || cfg.BuildTrendKeepCount > 10000 {
		return errors.New("invalid build trend keep count")
	}
	if cfg.WatchMode != "github" && cfg.WatchMode != "local" {
		return errors.New("invalid watch mode")
	}
	if cfg.WatchMode == "local" && cfg.TagFilter.Enabled {
		return errors.New("invalid tag filter for local watch")
	}
	if err := validateTagFilter(cfg.TagFilter); err != nil {
		return err
	}
	if err := validateRemoteBuild(cfg.RemoteBuild); err != nil {
		return err
	}
	if err := validateDurationAnomaly(cfg.DurationAnomaly); err != nil {
		return err
	}
	if cfg.DeployParallelism < 1 || cfg.DeployParallelism > 16 {
		return errors.New("invalid deploy parallelism")
	}
	if cfg.LogArchiveAfterDays < 0 || cfg.LogArchiveAfterDays > 3650 {
		return errors.New("invalid log archive after days")
	}
	if cfg.AllowedHours != nil && (cfg.AllowedHours.From < 0 || cfg.AllowedHours.From > 23 || cfg.AllowedHours.To < 0 || cfg.AllowedHours.To > 23 || cfg.AllowedHours.From >= cfg.AllowedHours.To) {
		return errors.New("CONFIG_ALLOWED_HOURS_INVALID")
	}
	if len(cfg.BranchTargets) == 0 {
		return errors.New("branch targets must not be empty")
	}
	for _, t := range cfg.BranchTargets {
		if t.Branch == "" || strings.Contains(t.Branch, "..") || strings.Contains(t.Branch, "~") || hasUnsafeControl(t.Branch) {
			return fmt.Errorf("invalid branch target: %s", t.Branch)
		}
		if t.TargetFile == "" || filepath.IsAbs(t.TargetFile) || strings.Contains(t.TargetFile, "..") {
			return fmt.Errorf("invalid target file: %s", t.TargetFile)
		}
		if len(t.TargetFiles) == 0 || len(t.TargetFiles) > 100 {
			return errors.New("target files must not be empty")
		}
		hasPrimaryTarget := false
		for _, targetFile := range t.TargetFiles {
			if targetFile == t.TargetFile {
				hasPrimaryTarget = true
			}
			if targetFile == "" || filepath.IsAbs(targetFile) || strings.Contains(targetFile, "..") || hasUnsafeControl(targetFile) {
				return fmt.Errorf("invalid target file: %s", targetFile)
			}
		}
		if !hasPrimaryTarget {
			return errors.New("target_files must include target_file")
		}
		for _, p := range []string{t.SHAFile, t.Src, t.Out} {
			if p == "" || !filepath.IsAbs(p) {
				return fmt.Errorf("path must be absolute: %s", p)
			}
		}
		if sameOrNestedPath(t.Src, t.Out) {
			return fmt.Errorf("CONFIG_PATH_CONFLICT: src=%s out=%s", t.Src, t.Out)
		}
		if err := validateRunnerEnv(t.Env); err != nil {
			return err
		}
		for _, d := range t.DeployTargets {
			if d.ID == "" || d.Host == "" || d.User == "" || d.DestDir == "" || !strings.HasPrefix(d.DestDir, "/") {
				return errors.New("invalid deploy target")
			}
		}
	}
	return nil
}

func hasUnsafeControl(s string) bool {
	for _, r := range s {
		if r == 0 || r == '\n' || r == '\r' || r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}

func validateTagFilter(filter TagFilterConfig) error {
	if !filter.Enabled {
		return nil
	}
	if len(filter.Patterns) > 100 {
		return errors.New("invalid tag filter patterns")
	}
	for _, pattern := range filter.Patterns {
		if pattern == "" || hasUnsafeControl(pattern) {
			return errors.New("invalid tag filter pattern")
		}
		if strings.Count(pattern, "*") > 1 {
			return errors.New("invalid tag filter pattern")
		}
		if strings.Contains(pattern, "*") && !strings.HasSuffix(pattern, "*") {
			return errors.New("invalid tag filter pattern")
		}
	}
	return nil
}

func validateRemoteBuild(remote RemoteBuildConfig) error {
	if !remote.Enabled {
		return nil
	}
	for _, value := range []*string{remote.Host, remote.User, remote.WorkDir} {
		if value == nil || *value == "" || hasUnsafeControl(*value) {
			return errors.New("invalid remote build")
		}
	}
	if remote.WorkDir != nil && !strings.HasPrefix(*remote.WorkDir, "/") {
		return errors.New("invalid remote build")
	}
	if len(remote.CommandArgs) == 0 || len(remote.CommandArgs) > 64 {
		return errors.New("invalid remote build")
	}
	for _, arg := range remote.CommandArgs {
		if arg == "" || hasUnsafeControl(arg) {
			return errors.New("invalid remote build")
		}
	}
	if remote.ArtifactPath != nil && (*remote.ArtifactPath == "" || hasUnsafeControl(*remote.ArtifactPath)) {
		return errors.New("invalid remote build")
	}
	return nil
}

func validateDurationAnomaly(config DurationAnomalyConfig) error {
	if !config.Enabled {
		return nil
	}
	if config.MinSamples < 1 || config.MinSamples > 10000 || config.AvgMultiplier <= 0 || config.P95Multiplier <= 0 {
		return errors.New("invalid duration anomaly")
	}
	return nil
}

func sameOrNestedPath(a, b string) bool {
	cleanA := filepath.Clean(a)
	cleanB := filepath.Clean(b)
	if cleanA == cleanB {
		return true
	}
	relAB, errAB := filepath.Rel(cleanA, cleanB)
	relBA, errBA := filepath.Rel(cleanB, cleanA)
	return (errAB == nil && relAB != "." && !strings.HasPrefix(relAB, ".."+string(os.PathSeparator)) && relAB != "..") ||
		(errBA == nil && relBA != "." && !strings.HasPrefix(relBA, ".."+string(os.PathSeparator)) && relBA != "..")
}

func validateRunnerEnv(env map[string]string) error {
	if len(env) > 100 {
		return errors.New("invalid branch env")
	}
	for key, value := range env {
		if !validRunnerEnvKey(key) || strings.HasPrefix(key, "ADLAIRE_CI_") {
			return fmt.Errorf("invalid branch env key: %s", key)
		}
		switch key {
		case "PATH", "HOME", "SHELL", "USER", "GITHUB_TOKEN", "ADLAIRE_TOKEN", "ADLAIRE_CHANGED_TARGETS":
			return fmt.Errorf("reserved branch env key: %s", key)
		}
		if !utf8.ValidString(value) || len(value) > 4096 || strings.ContainsAny(value, "\x00\n\r") {
			return fmt.Errorf("invalid branch env value: %s", key)
		}
	}
	return nil
}

func validRunnerEnvKey(key string) bool {
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

func readRunnerPipelineConfig(stateDir string) (pipelineConfigFile, error) {
	path := filepath.Join(stateDir, ".pipeline_config")
	data, err := os.ReadFile(path)
	if err != nil {
		return pipelineConfigFile{}, err
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	var config pipelineConfigFile
	if err := dec.Decode(&config); err != nil {
		return pipelineConfigFile{}, err
	}
	var extra struct{}
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return pipelineConfigFile{}, errors.New("pipeline config must contain one JSON object")
	}
	if err := validateRunnerPipelineConfig(config); err != nil {
		return pipelineConfigFile{}, err
	}
	return config, nil
}

func validateRunnerPipelineConfig(config pipelineConfigFile) error {
	if config.ExtraArgs == nil || config.Env == nil {
		return errors.New("pipeline config requires extra_args and env")
	}
	if len(config.ExtraArgs) > 50 {
		return errors.New("too many pipeline extra args")
	}
	for _, arg := range config.ExtraArgs {
		if !utf8.ValidString(arg) || arg == "" || strings.ContainsAny(arg, "\x00\n\r") {
			return errors.New("invalid pipeline extra arg")
		}
		if reservedPipelineExtraArg(arg) {
			return errors.New("reserved pipeline extra arg")
		}
	}
	if err := validateRunnerEnv(config.Env); err != nil {
		return fmt.Errorf("invalid pipeline env: %w", err)
	}
	return nil
}

func reservedPipelineExtraArg(arg string) bool {
	reserved := []string{"--src", "--out", "--build-id", "--commit-sha", "--build-at", "--cache-dir", "--version", "--help"}
	for _, option := range reserved {
		if arg == option || strings.HasPrefix(arg, option+"=") {
			return true
		}
	}
	return false
}

func executeRunnerDryRun(cfg RunnerConfig, stdout io.Writer) int {
	result := map[string]any{
		"mode":           "dry-run",
		"dry_run":        true,
		"state_dir":      cfg.StateDir,
		"targets":        []map[string]any{},
		"would_write":    []string{},
		"would_call":     []string{},
		"secrets_masked": true,
		"errors":         []map[string]any{},
		"warnings":       []map[string]any{},
	}
	if err := loadRunnerConfig(&cfg, slog.New(slog.NewTextHandler(io.Discard, nil))); err != nil {
		result["errors"] = []map[string]any{{"code": "config_error", "message": err.Error(), "target": nil}}
		data, _ := json.Marshal(result)
		fmt.Fprintln(stdout, string(data))
		return 2
	}
	targets := []map[string]any{}
	wouldWrite := map[string]bool{}
	for _, target := range cfg.BranchTargets {
		prev, prevErr := readSHACache(target.SHAFile)
		prevValue := any(nil)
		if prevErr == nil && prev != "" {
			prevValue = prev
		}
		reason := "sha_changed"
		wouldBuild := true
		if prevErr != nil {
			reason = "config_error"
			wouldBuild = false
		} else if cfg.SchedulePaused {
			reason = "config_error"
			wouldBuild = false
		} else if cooldownActive(cfg, slog.New(slog.NewTextHandler(io.Discard, nil))) {
			reason = "cooldown"
			wouldBuild = false
		}
		if wouldBuild {
			for _, key := range []string{"lock", "build_log", "history", "status", "sha_cache"} {
				wouldWrite[key] = true
			}
		}
		targets = append(targets, map[string]any{
			"branch":             target.Branch,
			"target_file":        target.TargetFile,
			"previous_sha":       prevValue,
			"current_blob_sha":   nil,
			"current_commit_sha": nil,
			"would_build":        wouldBuild,
			"trigger":            "polling",
			"reason":             reason,
			"warnings":           []map[string]any{},
		})
	}
	keys := make([]string, 0, len(wouldWrite))
	for key := range wouldWrite {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result["targets"] = targets
	result["would_write"] = keys
	data, _ := json.Marshal(result)
	fmt.Fprintln(stdout, string(data))
	return 0
}

func cooldownActive(cfg RunnerConfig, logger *slog.Logger) bool {
	if cfg.BuildCooldownSeconds <= 0 {
		return false
	}
	state, err := runnerReadBuildState(cfg.StateDir)
	if err != nil || state.LastFinishedAt == nil || *state.LastFinishedAt == "" {
		return false
	}
	last, err := time.Parse(time.RFC3339, *state.LastFinishedAt)
	if err != nil {
		return false
	}
	if runnerNow().Sub(last) < time.Duration(cfg.BuildCooldownSeconds)*time.Second {
		logger.Info("COOLDOWN: skip")
		return true
	}
	return false
}

func forceIntervalDue(cfg RunnerConfig) bool {
	if cfg.ForceBuildIntervalHours <= 0 {
		return false
	}
	state, err := runnerReadBuildState(cfg.StateDir)
	if err != nil || state.LastFinishedAt == nil || *state.LastFinishedAt == "" {
		return true
	}
	last, err := time.Parse(time.RFC3339, *state.LastFinishedAt)
	if err != nil {
		return true
	}
	return runnerNow().Sub(last) >= time.Duration(cfg.ForceBuildIntervalHours)*time.Hour
}

func readQueueRunContext(stateDir string) (queueRunContext, error) {
	state, err := runnerReadBuildState(stateDir)
	if err != nil {
		return queueRunContext{}, err
	}
	if len(state.ActiveQueueEntry) > 0 {
		entry, err := parseQueueEntry(state.ActiveQueueEntry)
		if err != nil {
			return queueRunContext{}, err
		}
		return queueRunContextFromEntry(entry), nil
	}
	if len(state.Queued) == 0 {
		return queueRunContext{}, nil
	}
	entries := make([]queueEntry, 0, len(state.Queued))
	for _, raw := range state.Queued {
		entry, err := parseQueueEntry(raw)
		if err != nil {
			return queueRunContext{}, err
		}
		entries = append(entries, *entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		left, right := entries[i], entries[j]
		if priorityWeight(left.Priority) != priorityWeight(right.Priority) {
			return priorityWeight(left.Priority) < priorityWeight(right.Priority)
		}
		if left.CreatedSeq != right.CreatedSeq {
			return left.CreatedSeq < right.CreatedSeq
		}
		return left.ID < right.ID
	})
	return queueRunContextFromEntry(&entries[0]), nil
}

func parseQueueEntry(raw map[string]any) (*queueEntry, error) {
	id, ok := raw["id"].(string)
	if !ok || id == "" {
		return nil, errors.New("invalid queue entry id")
	}
	trigger, ok := raw["trigger"].(string)
	if !ok || !validQueueTrigger(trigger) {
		return nil, errors.New("invalid queue entry trigger")
	}
	queuedAt, ok := raw["queued_at"].(string)
	if !ok || queuedAt == "" {
		return nil, errors.New("invalid queue entry queued_at")
	}
	requestedBy, ok := raw["requested_by"].(string)
	if !ok || requestedBy == "" {
		return nil, errors.New("invalid queue entry requested_by")
	}
	priority, ok := raw["priority"].(string)
	if !ok || priorityWeight(priority) < 0 {
		return nil, errors.New("invalid queue entry priority")
	}
	createdSeq, ok := jsonNumberToInt64(raw["created_seq"])
	if !ok || createdSeq < 1 {
		return nil, errors.New("invalid queue entry created_seq")
	}
	payload, ok := raw["payload"].(map[string]any)
	if !ok {
		return nil, errors.New("invalid queue entry payload")
	}
	return &queueEntry{
		ID: id, Trigger: trigger, QueuedAt: queuedAt, RequestedBy: requestedBy,
		Priority: priority, CreatedSeq: createdSeq, Payload: payload, OriginalData: copyStringAnyMap(raw),
	}, nil
}

func jsonNumberToInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case float64:
		if v != float64(int64(v)) {
			return 0, false
		}
		return int64(v), true
	case int:
		return int64(v), true
	case int64:
		return v, true
	default:
		return 0, false
	}
}

func validQueueTrigger(trigger string) bool {
	return trigger == "manual" || trigger == "webhook" || trigger == "approval"
}

func priorityWeight(priority string) int {
	switch priority {
	case "urgent":
		return 0
	case "high":
		return 10
	case "normal":
		return 20
	case "low":
		return 30
	default:
		return -1
	}
}

func queueRunContextFromEntry(entry *queueEntry) queueRunContext {
	ctx := queueRunContext{Entry: entry, Trigger: entry.Trigger}
	switch entry.Trigger {
	case "manual":
		ctx.Force, _ = entry.Payload["force"].(bool)
	case "webhook":
		ctx.Branch, _ = entry.Payload["branch"].(string)
		ctx.CommitSHA, _ = entry.Payload["sha"].(string)
	case "approval":
		ctx.Branch, _ = entry.Payload["branch"].(string)
		ctx.TargetFile, _ = entry.Payload["target"].(string)
		ctx.CommitSHA, _ = entry.Payload["sha"].(string)
		ctx.Force, _ = entry.Payload["requested_force"].(bool)
	}
	return ctx
}

func queueMatchesTarget(ctx queueRunContext, target BranchTarget) bool {
	if ctx.Entry == nil {
		return true
	}
	if ctx.Branch != "" && ctx.Branch != target.Branch {
		return false
	}
	if ctx.TargetFile == "" {
		return true
	}
	for _, targetFile := range normalizedTargetFiles(target) {
		if targetFile == ctx.TargetFile {
			return true
		}
	}
	return false
}

func copyStringAnyMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func runnerReadBuildState(stateDir string) (buildState, error) {
	var state buildState
	err := runnerReadJSONFile(filepath.Join(stateDir, ".build_state"), &state)
	if errors.Is(err, os.ErrNotExist) {
		return runnerDefaultBuildState(), nil
	}
	return state, err
}

func readRunnerToken(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", errors.New("GITHUB_TOKEN_MISSING: path=" + path)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return "", errors.New("GITHUB_TOKEN_INVALID: reason=filetype")
	}
	if info.Mode().Perm()&0077 != 0 {
		return "", fmt.Errorf("GITHUB_TOKEN_INSECURE_MODE: path=%s mode=%04o", path, info.Mode().Perm())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", errors.New("GITHUB_TOKEN_PERMISSION: path=" + path)
	}
	if !utf8.Valid(data) {
		return "", errors.New("GITHUB_TOKEN_INVALID: reason=utf8")
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", errors.New("GITHUB_TOKEN_INVALID: reason=empty")
	}
	for _, r := range token {
		if r == 0 || r == '\n' || r == '\r' || r == '\t' || r < 0x20 || r == 0x7f || r == ' ' {
			return "", errors.New("GITHUB_TOKEN_INVALID: reason=character")
		}
	}
	return token, nil
}

func acquireRunnerLock(path string, logger *slog.Logger) (func(), bool, int) {
	now := runnerNow().UTC().Format(time.RFC3339)
	data := []byte(fmt.Sprintf("pid=%d\nstarted_at=%s\n", os.Getpid(), now))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err == nil {
		_, writeErr := f.Write(data)
		closeErr := f.Close()
		if writeErr != nil || closeErr != nil {
			logger.Error("LOCK_CREATE_FAILED: write failed")
			_ = os.Remove(path)
			return func() {}, false, 4
		}
		return func() { releaseRunnerLock(path, data, logger) }, true, 0
	}
	existing, readErr := os.ReadFile(path)
	if readErr != nil {
		logger.Error("LOCK_CREATE_FAILED: " + readErr.Error())
		return func() {}, false, 4
	}
	pid, parseErr := parseLockPID(string(existing))
	if parseErr != nil {
		logger.Error("LOCK_CORRUPT: path=" + path)
		return func() {}, false, 4
	}
	if pidRunning(pid) {
		logger.Info(fmt.Sprintf("BUILD_SKIP: already running (PID %d)", pid))
		return func() {}, false, 0
	}
	logger.Warn(fmt.Sprintf("STALE_LOCK: pid=%d", pid))
	if err := os.Remove(path); err != nil {
		logger.Error("LOCK_CREATE_FAILED: " + err.Error())
		return func() {}, false, 4
	}
	return acquireRunnerLock(path, logger)
}

func parseLockPID(text string) (int, error) {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "pid=") {
			return strconv.Atoi(strings.TrimPrefix(line, "pid="))
		}
	}
	return 0, errors.New("pid not found")
}

func pidRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	if _, err := os.Stat("/proc"); errors.Is(err, os.ErrNotExist) {
		err = syscall.Kill(pid, 0)
		return err == nil || errors.Is(err, syscall.EPERM)
	}
	if _, err := os.Stat(filepath.Join("/proc", strconv.Itoa(pid))); err == nil {
		return true
	} else if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, os.ErrPermission) {
		err = syscall.Kill(pid, 0)
		return err == nil || errors.Is(err, syscall.EPERM)
	}
	return false
}

func releaseRunnerLock(path string, expected []byte, logger *slog.Logger) {
	current, err := os.ReadFile(path)
	if err != nil {
		if logger != nil {
			logger.Warn("LOCK_RELEASE_FAILED: " + err.Error())
		}
		return
	}
	if !bytes.Equal(current, expected) {
		if logger != nil {
			logger.Warn("LOCK_RELEASE_SKIPPED: owner mismatch")
		}
		return
	}
	if err := os.Remove(path); err != nil && logger != nil {
		logger.Warn("LOCK_RELEASE_FAILED: " + err.Error())
	}
}

func processRunnerTarget(cfg RunnerConfig, idx int, target BranchTarget, token, buildID string, logger *slog.Logger, queueCtx queueRunContext) int {
	return processRunnerTargetWithChain(cfg, idx, target, token, buildID, logger, queueCtx, nil, nil)
}

func processRunnerTargetWithChain(cfg RunnerConfig, idx int, target BranchTarget, token, buildID string, logger *slog.Logger, queueCtx queueRunContext, chain *chainLog, chainSummary *chainSummaryLog) int {
	started := runnerNow().UTC()
	startedText := started.Format(time.RFC3339)
	trigger := runnerInitialTrigger(cfg)
	if queueCtx.Entry != nil {
		trigger = queueCtx.Trigger
	}
	if err := writeBuildStatusRunning(cfg, buildID, trigger, target.Branch, target.TargetFile, startedText); err != nil {
		logger.Error("BUILD_STATUS_WRITE_FAILED: " + err.Error())
		return 1
	}
	if err := writeBuildStateStart(cfg.StateDir, true, &buildID, queueCtx.Entry); err != nil {
		logger.Error("STATE_START_FAILED: " + err.Error())
		_ = writeBuildStatusError(cfg, "failure_state_write", trigger, err.Error())
		return 1
	}
	if err := writeBuildLog(cfg.StateDir, runningBuildLog(buildID, target, trigger, startedText)); err != nil {
		logger.Error("BUILD_LOG_START_FAILED: " + err.Error())
		_ = writeBuildStatusError(cfg, "failure_state_write", trigger, err.Error())
		return 1
	}
	prevSHA, prevErr := readSHACache(target.SHAFile)
	if prevErr != nil {
		logFailure(cfg, target, buildID, started, nil, "", trigger, "failure_decode", "sha cache decode failed", nil, nil, logger)
		recordCircuitFailure(cfg, "sha cache decode failed")
		return 1
	}
	resolved, err := resolveRunnerTarget(cfg, logger, token, target)
	if err != nil {
		status := "failure_api"
		exit := 3
		if cfg.WatchMode == "local" {
			status = "failure_decode"
			exit = 1
		}
		logFailure(cfg, target, buildID, started, nil, prevSHA, trigger, status, "source resolve failed", nil, nil, logger)
		recordCircuitFailure(cfg, "source resolve failed")
		return exit
	}
	forceBuild := queueCtx.Force || forceIntervalDue(cfg)
	if resolved.Digest == prevSHA && !forceBuild {
		logger.Info(fmt.Sprintf("NO_CHANGE: branch=%s target=%s sha=%s", target.Branch, target.TargetFile, resolved.Digest))
		_ = skipRunnerTarget(cfg, buildID, "skipped_no_change", trigger, &target)
		return 0
	}
	if resolved.Digest == prevSHA && forceBuild && queueCtx.Entry == nil {
		trigger = "force_interval"
	}
	if cfg.AllowedHours != nil && !allowedNow(*cfg.AllowedHours, runnerNow().UTC()) {
		_ = skipRunnerTarget(cfg, buildID, "skipped_allowed_hours", trigger, &target)
		return 0
	}
	commit := commitInfo{}
	matchedTags := []string{}
	if cfg.WatchMode == "github" {
		commit = fetchCommitInfo(cfg, logger, token, target.Branch, target.TargetFile)
		if cfg.TagFilter.Enabled {
			if commit.SHA == nil || *commit.SHA == "" {
				logFailure(cfg, target, buildID, started, &resolved.Digest, prevSHA, trigger, "failure_api", "github api failed", nil, nil, logger)
				recordCircuitFailure(cfg, "github api failed")
				return 3
			}
			tags, err := fetchMatchingTags(cfg, logger, token, *commit.SHA)
			if err != nil {
				logFailure(cfg, target, buildID, started, &resolved.Digest, prevSHA, trigger, "failure_api", "github api failed", nil, nil, logger)
				recordCircuitFailure(cfg, "github api failed")
				return 3
			}
			matchedTags = tags
			if len(matchedTags) == 0 {
				_ = skipRunnerTarget(cfg, buildID, "skipped_tag_filter", trigger, &target)
				return 0
			}
		}
	}
	if target.ApprovalRequired {
		if err := appendApprovalPending(cfg, target, buildID, resolved.Digest, resolved.ChangedTargets); err != nil {
			logFailure(cfg, target, buildID, started, &resolved.Digest, prevSHA, trigger, "failure_state_write", "approval pending write failed", nil, nil, logger)
			return 1
		}
		_ = os.Remove(filepath.Join(cfg.StateDir, ".build_logs", buildID+".json"))
		_ = clearRunningBuildState(cfg.StateDir)
		_ = writeBuildStatusSkip(cfg, "pending_approval", trigger, &target)
		return 0
	}
	if err := materializeResolvedSource(target.Src, resolved); err != nil {
		logFailure(cfg, target, buildID, started, &resolved.Digest, prevSHA, trigger, "failure_decode", "source write failed", nil, nil, logger)
		recordCircuitFailure(cfg, "source write failed")
		return 1
	}
	builderVersion, err := precheckRunnerTarget(cfg, target)
	if err != nil {
		logFailure(cfg, target, buildID, started, &resolved.Digest, prevSHA, trigger, "failure_precheck", "precheck failed", nil, nil, logger)
		recordCircuitFailure(cfg, "precheck failed")
		return 1
	}
	commitSHA := ""
	if commit.SHA != nil {
		commitSHA = *commit.SHA
	}
	hooks, skippedHookIDs, err := readRunnerHooks(cfg.StateDir)
	if err != nil {
		logHookAbortFailure(cfg, target, buildID, started, &resolved.Digest, prevSHA, trigger, "hook config invalid", skippedHookIDs, nil, logger)
		recordCircuitFailure(cfg, "hook config invalid")
		return 1
	}
	hookWarnings, preAbortReason, _ := executeRunnerHooks(cfg, target, buildID, "pre", commitSHA, startedText, resolved.ChangedTargets, hooks, logger)
	if preAbortReason != nil {
		logHookAbortFailure(cfg, target, buildID, started, &resolved.Digest, prevSHA, trigger, *preAbortReason, skippedHookIDs, hookWarnings, logger)
		recordCircuitFailure(cfg, *preAbortReason)
		return 1
	}
	commitStatusState := (*string)(nil)
	if cfg.CommitStatusEnabled {
		state := "error"
		if commit.SHA != nil {
			if err := postCommitStatus(cfg, token, *commit.SHA, "pending", buildID, "Build started"); err == nil {
				state = "pending"
			} else {
				logger.Warn("COMMIT_STATUS_FAILED: " + err.Error())
			}
		} else {
			logger.Warn("COMMIT_STATUS_FAILED: commit sha unavailable")
		}
		commitStatusState = &state
	}
	attempts := []pipelineLog{}
	pl := runPipeline(cfg, target, buildID, commitSHA, startedText, resolved.ChangedTargets)
	attempts = append(attempts, pl)
	for retry := 1; retry <= cfg.BuildRetryMax && shouldRetryPipeline(pl); retry++ {
		runnerSleep(time.Duration(cfg.BuildRetryBaseSeconds*retry) * time.Second)
		pl = runPipeline(cfg, target, buildID, commitSHA, startedText, resolved.ChangedTargets)
		attempts = append(attempts, pl)
	}
	retryCount := len(attempts) - 1
	rep, warns := parseRunnerReport(pl.Stdout)
	if rep != nil {
		rep.BuildID = buildID
		rep.CommitSHA = commitSHA
		rep.BuildAt = startedText
	}
	if pl.ExitCode != nil && *pl.ExitCode == 0 && rep == nil {
		warns = append(warns, "REPORT_MISSING")
	}
	sizeWarn := rep != nil && rep.SizeWarn
	if cfg.OutputSizeWarnMB > 0 {
		if size, ok := outputSizeBytes(target.Out); ok && size > int64(cfg.OutputSizeWarnMB)*1024*1024 {
			sizeWarn = true
			if rep != nil {
				rep.SizeWarn = true
			}
			warns = append(warns, "OUTPUT_SIZE_WARN")
		}
	}
	deploys := []deployLog{}
	status := "success"
	var errText *string
	exitCode := 0
	if pl.TargetStatus != "" {
		status = pl.TargetStatus
		msg := pl.ErrorMessage
		if msg == "" {
			msg = "pipeline failed"
		}
		errText = &msg
		exitCode = 1
		if pl.TargetStatus == "failure_pipeline_config" {
			exitCode = 2
		}
	} else if pl.ExitCode == nil {
		status = "failure_timeout"
		msg := "pipeline timeout"
		errText = &msg
		exitCode = 1
	} else if *pl.ExitCode != 0 {
		status = "failure_build"
		msg := "pipeline failed"
		errText = &msg
		exitCode = 1
	} else if err := validateOutputSite(target.Out); err != nil {
		status = "failure_build"
		msg := "output validation failed"
		errText = &msg
		exitCode = 1
	} else if err := runnerAtomicWriteJSON(target.SHAFile, shaCache{SHA: resolved.Digest}, 0600); err != nil {
		status = "failure_state_write"
		msg := "state write failed"
		errText = &msg
		exitCode = 1
	} else {
		deploys = executeDeploys(cfg, idx, target, buildID)
		for _, dl := range deploys {
			if dl.Status == "pending" {
				status = "success_deploy_pending"
				msg := "deploy pending"
				errText = &msg
				exitCode = 1
			}
		}
	}
	var snapshotID *string
	if status == "success" || status == "success_deploy_pending" {
		if err := writeDependencyManifest(target, pl.EnvKeys); err != nil {
			warns = append(warns, "DEPENDENCY_MANIFEST_FAILED")
			logger.Warn("DEPENDENCY_MANIFEST_FAILED: " + err.Error())
		}
	}
	if status == "success" {
		if id, err := writeSnapshot(cfg, target, buildID); err == nil {
			snapshotID = &id
		} else {
			logger.Warn("SNAPSHOT_FAILED: " + err.Error())
		}
	}
	finished := runnerNow().UTC()
	outputSHA, _ := outputManifestSHA(target.Out)
	var outputSize *int64
	if size, ok := outputSizeBytes(target.Out); ok {
		outputSize = &size
	}
	targetResults := buildTargetResults(deploys)
	transferVerified := buildTransferVerified(deploys)
	failureCategory := buildFailureCategory(status)
	failureEvidence := buildFailureEvidence(status, failureCategory, finished)
	normalizedStatus := buildStatusFromTargetStatus(status)
	if cfg.CommitStatusEnabled {
		state := "error"
		if commit.SHA != nil {
			state = commitStatusStateForTargetStatus(status)
			if err := postCommitStatus(cfg, token, *commit.SHA, state, buildID, "Build "+state); err != nil {
				logger.Warn("COMMIT_STATUS_FAILED: " + err.Error())
				state = "error"
			}
		}
		commitStatusState = &state
	}
	blog := buildLog{
		ID:               buildID,
		Status:           normalizedStatus,
		Branch:           target.Branch,
		TargetFile:       target.TargetFile,
		TargetFiles:      append([]string(nil), target.TargetFiles...),
		ChangedTargets:   append([]string(nil), resolved.ChangedTargets...),
		MatchedTags:      append([]string(nil), matchedTags...),
		TargetStatus:     status,
		Trigger:          trigger,
		TriggerActor:     nil,
		StartedAt:        started.Format(time.RFC3339),
		FinishedAt:       finished.Format(time.RFC3339),
		DurationSeconds:  int64(finished.Sub(started).Seconds()),
		Commit:           commit,
		CommitSHA:        commit.SHA,
		CommitMessage:    commit.Message,
		CommitAuthor:     commit.Author,
		CommitAt:         commit.Date,
		BlobSHA:          &resolved.Digest,
		PreviousBlobSHA:  prevSHA,
		Pipeline:         pl,
		SkippedHookIDs:   skippedHookIDs,
		Attempts:         attempts,
		RetryCount:       retryCount,
		RemoteBuild:      pl.RemoteBuild,
		Chain:            chain,
		ChainSummary:     chainSummary,
		Report:           rep,
		Warnings:         append(warns, hookWarnings...),
		Deploy:           deploys,
		TargetResults:    targetResults,
		SnapshotID:       snapshotID,
		RollbackFrom:     nil,
		TransferVerified: transferVerified,
		OutputSHA256:     outputSHA,
		OutputSizeBytes:  outputSize,
		SizeWarn:         sizeWarn,
		FailureCategory:  failureCategory,
		FailureEvidence:  failureEvidence,
		CommitStatus:     commitStatusState,
		CommitStatusLog:  buildCommitStatusLog(cfg, commitStatusState, finished),
		BuildMeta:        buildBuildMeta(rep),
		Environment:      currentBuildEnv(cfg, builderVersion, pl.EnvKeys),
		Error:            errText,
	}
	if err := writeBuildLog(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_LOG_WRITE_FAILED: " + err.Error())
		_ = writeBuildStatusError(cfg, "failure_state_write", trigger, err.Error())
		return 1
	}
	if err := appendHistory(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_HISTORY_WRITE_FAILED: " + err.Error())
		_ = writeBuildStatusError(cfg, "failure_state_write", trigger, err.Error())
		return 1
	}
	if err := updateBuildTrends(cfg, blog); err != nil {
		logger.Warn("BUILD_TRENDS_WRITE_FAILED: " + err.Error())
	}
	if err := writeBuildStatusFinal(cfg, blog); err != nil {
		logger.Error("BUILD_STATUS_WRITE_FAILED: " + err.Error())
		exitCode = 1
	}
	postWarnings, _, postExitBump := executeRunnerHooks(cfg, target, buildID, "post", commitSHA, startedText, resolved.ChangedTargets, hooks, logger)
	if len(postWarnings) > 0 || len(skippedHookIDs) > 0 {
		blog.Warnings = append(blog.Warnings, postWarnings...)
		if err := writeBuildLog(cfg.StateDir, blog); err != nil {
			logger.Error("BUILD_LOG_WRITE_FAILED: " + err.Error())
			exitCode = 1
		}
	}
	if postExitBump {
		exitCode = 1
	}
	sendBuildNotifications(cfg, blog, logger)
	if exitCode == 0 {
		resetCircuitState(cfg)
	} else if strings.HasPrefix(status, "failure_") {
		recordCircuitFailure(cfg, derefString(errText, "runner target failed"))
	}
	return exitCode
}

func runningBuildLog(buildID string, target BranchTarget, trigger, started string) buildLog {
	return buildLog{
		ID:               buildID,
		Status:           "running",
		Branch:           target.Branch,
		TargetFile:       target.TargetFile,
		TargetFiles:      append([]string(nil), target.TargetFiles...),
		ChangedTargets:   []string{},
		MatchedTags:      []string{},
		TargetStatus:     "running",
		Trigger:          trigger,
		TriggerActor:     nil,
		StartedAt:        started,
		FinishedAt:       "",
		DurationSeconds:  0,
		Commit:           commitInfo{},
		CommitSHA:        nil,
		CommitMessage:    nil,
		CommitAuthor:     nil,
		CommitAt:         nil,
		BlobSHA:          nil,
		PreviousBlobSHA:  "",
		Pipeline:         pipelineLog{},
		SkippedHookIDs:   []string{},
		Attempts:         []pipelineLog{},
		RetryCount:       0,
		RemoteBuild:      nil,
		Chain:            nil,
		ChainSummary:     nil,
		Report:           nil,
		Warnings:         []string{},
		Deploy:           []deployLog{},
		TargetResults:    []targetResultLog{},
		SnapshotID:       nil,
		RollbackFrom:     nil,
		TransferVerified: nil,
		OutputSHA256:     nil,
		OutputSizeBytes:  nil,
		SizeWarn:         false,
		FailureCategory:  nil,
		FailureEvidence:  []failureEvidence{},
		CommitStatus:     nil,
		CommitStatusLog:  nil,
		BuildMeta:        nil,
		Environment:      buildEnvLog{},
		Error:            nil,
	}
}

func runnerInitialTrigger(cfg RunnerConfig) string {
	if cfg.WatchMode == "local" {
		return "local_watch"
	}
	return "polling"
}

func skipRunnerTarget(cfg RunnerConfig, buildID, targetStatus, trigger string, target *BranchTarget) error {
	_ = os.Remove(filepath.Join(cfg.StateDir, ".build_logs", buildID+".json"))
	if err := clearRunningBuildState(cfg.StateDir); err != nil {
		return err
	}
	return writeBuildStatusSkip(cfg, targetStatus, trigger, target)
}

func allowedNow(hours AllowedHoursConfig, now time.Time) bool {
	hour := now.UTC().Hour()
	return hours.From <= hour && hour < hours.To
}

func resolveRunnerTarget(cfg RunnerConfig, logger *slog.Logger, token string, target BranchTarget) (resolvedRunnerTarget, error) {
	if cfg.WatchMode == "local" {
		digest, err := localTargetDigest(target)
		if err != nil {
			return resolvedRunnerTarget{}, err
		}
		return resolvedRunnerTarget{
			Digest:         digest,
			ChangedTargets: normalizedTargetFiles(target),
			Local:          true,
		}, nil
	}
	shas, err := fetchTargetSHAs(cfg, logger, token, target.Branch, normalizedTargetFiles(target))
	if err != nil {
		return resolvedRunnerTarget{}, err
	}
	digest := digestTargetSHAs(target, shas)
	contents := map[string][]byte{}
	for _, path := range normalizedTargetFiles(target) {
		content, err := fetchBlobContent(cfg, logger, token, shas[path])
		if err != nil {
			return resolvedRunnerTarget{}, err
		}
		contents[path] = content
	}
	return resolvedRunnerTarget{
		Digest:         digest,
		ChangedTargets: normalizedTargetFiles(target),
		Contents:       contents,
	}, nil
}

func normalizedTargetFiles(target BranchTarget) []string {
	files := append([]string(nil), target.TargetFiles...)
	if len(files) == 0 && target.TargetFile != "" {
		files = []string{target.TargetFile}
	}
	sort.Strings(files)
	return files
}

func fetchTargetSHAs(cfg RunnerConfig, logger *slog.Logger, token, branch string, targetFiles []string) (map[string]string, error) {
	owner, repo := cfg.RepositoryOwner, cfg.RepositoryName
	apiURL := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=1", strings.TrimRight(runnerGitHubAPIBase, "/"), url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(branch))
	var tree gitTreeResponse
	if err := runnerGetJSON(cfg, logger, token, apiURL, &tree); err != nil {
		return nil, err
	}
	needed := map[string]bool{}
	for _, targetFile := range targetFiles {
		needed[targetFile] = true
	}
	found := map[string]string{}
	for _, item := range tree.Tree {
		if needed[item.Path] {
			found[item.Path] = item.SHA
		}
	}
	for _, targetFile := range targetFiles {
		if found[targetFile] == "" {
			return nil, fmt.Errorf("target not found: %s", targetFile)
		}
	}
	return found, nil
}

func digestTargetSHAs(target BranchTarget, shas map[string]string) string {
	files := normalizedTargetFiles(target)
	if len(files) == 1 {
		return shas[files[0]]
	}
	parts := make([]string, 0, len(files))
	for _, path := range files {
		parts = append(parts, path+"\n"+shas[path]+"\n")
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "")))
	return hex.EncodeToString(sum[:])
}

func localTargetDigest(target BranchTarget) (string, error) {
	files := normalizedTargetFiles(target)
	parts := make([]string, 0, len(files))
	for _, targetFile := range files {
		path := filepath.Join(target.Src, filepath.FromSlash(targetFile))
		if info, err := os.Lstat(path); err == nil {
			if info.IsDir() {
				sum, err := outputManifestSHA(path)
				if err != nil {
					return "", err
				}
				parts = append(parts, targetFile+"\n"+derefString(sum, "")+"\n")
				continue
			}
			if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
				return "", fmt.Errorf("invalid local target: %s", targetFile)
			}
			fileSum, err := runnerFileSHA256(path)
			if err != nil {
				return "", err
			}
			parts = append(parts, targetFile+"\n"+fileSum+"\n")
			continue
		}
		if len(files) == 1 {
			sum, err := outputManifestSHA(target.Src)
			if err != nil {
				return "", err
			}
			return derefString(sum, ""), nil
		}
		return "", fmt.Errorf("local target not found: %s", targetFile)
	}
	if len(files) == 1 {
		fields := strings.Split(parts[0], "\n")
		if len(fields) >= 2 {
			return fields[1], nil
		}
	}
	sort.Strings(parts)
	sum := sha256.Sum256([]byte(strings.Join(parts, "")))
	return hex.EncodeToString(sum[:]), nil
}

func fetchBlobContent(cfg RunnerConfig, logger *slog.Logger, token, sha string) ([]byte, error) {
	owner, repo := cfg.RepositoryOwner, cfg.RepositoryName
	apiURL := fmt.Sprintf("%s/repos/%s/%s/git/blobs/%s", strings.TrimRight(runnerGitHubAPIBase, "/"), url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(sha))
	var blob gitBlobResponse
	if err := runnerGetJSON(cfg, logger, token, apiURL, &blob); err != nil {
		return nil, err
	}
	if blob.Encoding != "base64" {
		return nil, errors.New("unsupported blob encoding")
	}
	content := strings.ReplaceAll(blob.Content, "\n", "")
	return base64.StdEncoding.DecodeString(content)
}

func fetchCommitInfo(cfg RunnerConfig, logger *slog.Logger, token, branch, targetFile string) commitInfo {
	owner, repo := cfg.RepositoryOwner, cfg.RepositoryName
	apiURL := fmt.Sprintf("%s/repos/%s/%s/commits?path=%s&sha=%s&per_page=1", strings.TrimRight(runnerGitHubAPIBase, "/"), url.PathEscape(owner), url.PathEscape(repo), url.QueryEscape(targetFile), url.QueryEscape(branch))
	var commits gitCommitResponse
	if err := runnerGetJSON(cfg, logger, token, apiURL, &commits); err != nil || len(commits) == 0 {
		return commitInfo{}
	}
	msg := strings.SplitN(commits[0].Commit.Message, "\n", 2)[0]
	sha := commits[0].SHA
	author := commits[0].Commit.Author.Name
	date := commits[0].Commit.Author.Date
	return commitInfo{SHA: &sha, Message: &msg, Author: &author, Date: &date}
}

func fetchMatchingTags(cfg RunnerConfig, logger *slog.Logger, token, commitSHA string) ([]string, error) {
	owner, repo := cfg.RepositoryOwner, cfg.RepositoryName
	apiURL := fmt.Sprintf("%s/repos/%s/%s/git/matching-refs/tags", strings.TrimRight(runnerGitHubAPIBase, "/"), url.PathEscape(owner), url.PathEscape(repo))
	var refs gitMatchingRefsResponse
	if err := runnerGetJSON(cfg, logger, token, apiURL, &refs); err != nil {
		return nil, err
	}
	matches := []string{}
	for _, ref := range refs {
		if ref.Object.SHA != commitSHA {
			continue
		}
		name := strings.TrimPrefix(ref.Ref, "refs/tags/")
		if tagMatchesFilter(name, cfg.TagFilter) {
			matches = append(matches, name)
			if len(matches) >= 100 {
				break
			}
		}
	}
	return matches, nil
}

func tagMatchesFilter(tag string, filter TagFilterConfig) bool {
	if !filter.Enabled {
		return true
	}
	if len(filter.Patterns) == 0 {
		return tag != ""
	}
	for _, pattern := range filter.Patterns {
		if pattern == "*" {
			return true
		}
		if strings.HasSuffix(pattern, "*") {
			if strings.HasPrefix(tag, strings.TrimSuffix(pattern, "*")) {
				return true
			}
			continue
		}
		if tag == pattern {
			return true
		}
	}
	return false
}

func runnerGetJSON(cfg RunnerConfig, logger *slog.Logger, token, url string, out any) error {
	client := &http.Client{Timeout: 30 * time.Second}
	var lastErr error
	for attempt := 0; attempt <= cfg.APIRetryMax; attempt++ {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", "adlaire-ci-runner")
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		resp, err := client.Do(req)
		if err == nil && resp != nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				defer resp.Body.Close()
				logPATExpiryWarning(resp, logger)
				return json.NewDecoder(resp.Body).Decode(out)
			}
			lastErr = fmt.Errorf("github api status %d", resp.StatusCode)
			if !runnerRetryable(resp) || attempt == cfg.APIRetryMax {
				resp.Body.Close()
				return lastErr
			}
			wait := retryDelay(cfg, attempt, resp)
			resp.Body.Close()
			runnerSleep(wait)
			continue
		}
		lastErr = err
		if attempt == cfg.APIRetryMax {
			break
		}
		runnerSleep(retryDelay(cfg, attempt, nil))
	}
	return lastErr
}

func postCommitStatus(cfg RunnerConfig, token, sha, state, buildID, description string) error {
	if sha == "" {
		return errors.New("commit sha unavailable")
	}
	payload := map[string]any{
		"state":       state,
		"context":     cfg.CommitStatusContext,
		"description": description,
	}
	if cfg.CommitStatusTargetURL != nil && *cfg.CommitStatusTargetURL != "" {
		payload["target_url"] = *cfg.CommitStatusTargetURL
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	apiURL := fmt.Sprintf("%s/repos/%s/%s/statuses/%s", strings.TrimRight(runnerGitHubAPIBase, "/"), url.PathEscape(cfg.RepositoryOwner), url.PathEscape(cfg.RepositoryName), url.PathEscape(sha))
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "adlaire-ci-runner")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Adlaire-Build-ID", buildID)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("github commit status %d", resp.StatusCode)
	}
	return nil
}

func commitStatusStateForTargetStatus(targetStatus string) string {
	if targetStatus == "success" {
		return "success"
	}
	return "failure"
}

func logPATExpiryWarning(resp *http.Response, logger *slog.Logger) {
	if resp == nil || logger == nil {
		return
	}
	raw := strings.TrimSpace(resp.Header.Get("GitHub-Authentication-Token-Expiration"))
	if raw == "" {
		return
	}
	exp, ok := parseGitHubTokenExpiration(raw)
	if !ok {
		logger.Warn("PAT_EXPIRY_PARSE_FAILED: value=" + raw)
		return
	}
	now := runnerNow().UTC()
	remaining := exp.UTC().Sub(now)
	if remaining < 0 {
		logger.Warn("PAT_EXPIRY_WARN: expires_at=" + exp.UTC().Format(time.RFC3339) + " remaining_days=0")
		return
	}
	if remaining <= 7*24*time.Hour {
		days := int(remaining.Hours() / 24)
		logger.Warn("PAT_EXPIRY_WARN: expires_at=" + exp.UTC().Format(time.RFC3339) + " remaining_days=" + strconv.Itoa(days))
	}
}

func parseGitHubTokenExpiration(value string) (time.Time, bool) {
	layouts := []string{time.RFC3339, time.RFC1123, "2006-01-02"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func runnerRetryable(resp *http.Response) bool {
	if resp == nil {
		return true
	}
	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	case http.StatusForbidden:
		return resp.Header.Get("X-RateLimit-Remaining") == "0"
	default:
		return false
	}
}

func retryDelay(cfg RunnerConfig, attempt int, resp *http.Response) time.Duration {
	if resp != nil && resp.Header.Get("X-RateLimit-Remaining") == "0" {
		if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
			if unix, err := strconv.ParseInt(reset, 10, 64); err == nil {
				wait := time.Until(time.Unix(unix, 0))
				if wait > 0 && wait <= time.Hour {
					return wait
				}
			}
		}
	}
	if cfg.APIRetryBaseSeconds <= 0 {
		return 0
	}
	return time.Duration(cfg.APIRetryBaseSeconds) * time.Second * time.Duration(1<<attempt)
}

func materializeResolvedSource(src string, resolved resolvedRunnerTarget) error {
	if resolved.Local {
		return nil
	}
	if err := os.RemoveAll(src); err != nil {
		return err
	}
	if err := os.MkdirAll(src, 0755); err != nil {
		return err
	}
	if len(resolved.Contents) == 1 {
		for _, content := range resolved.Contents {
			return os.WriteFile(filepath.Join(src, "source.md"), content, 0644)
		}
	}
	for rel, content := range resolved.Contents {
		if invalidOutputRelativePath(rel) {
			return fmt.Errorf("invalid source path: %s", rel)
		}
		path := filepath.Join(src, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}
	return nil
}

func writeDependencyManifest(target BranchTarget, envKeys []string) error {
	pageKey := dependencyPageKey(target)
	if pageKey == "" {
		return nil
	}
	deps, hashes, err := dependencyInputs(target.Src)
	if err != nil {
		return err
	}
	configHash := dependencyConfigHash(target, envKeys)
	manifest := runnerDependencyManifest{
		Version: 1,
		Builder: runnerDependencyManifestBuilder{
			Name:        "adlaire-ci-build",
			SpecSection: "28.1",
			ConfigHash:  configHash,
		},
		Pages: map[string]runnerDependencyPage{
			pageKey: {
				OutputPath:       "index.html",
				Dependencies:     deps,
				DependencySHA256: hashes,
				ConfigHash:       configHash,
			},
		},
		GeneratedOutputs: generatedOutputManifest(target.Out, pageKey),
	}
	return runnerAtomicWriteJSON(filepath.Join(target.Out, ".dependency_manifest.json"), manifest, 0600)
}

func dependencyPageKey(target BranchTarget) string {
	files := normalizedTargetFiles(target)
	if len(files) == 0 {
		return ""
	}
	return files[0]
}

func dependencyInputs(src string) ([]string, map[string]string, error) {
	info, err := os.Lstat(src)
	if err != nil {
		return nil, nil, err
	}
	root := src
	if !info.IsDir() {
		root = filepath.Dir(src)
	}
	deps := []string{}
	hashes := map[string]string{}
	err = filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("invalid dependency file: %s", path)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if invalidOutputRelativePath(rel) {
			return fmt.Errorf("invalid dependency path: %s", rel)
		}
		sum, err := runnerFileSHA256(path)
		if err != nil {
			return err
		}
		deps = append(deps, rel)
		hashes[rel] = sum
		return nil
	})
	sort.Strings(deps)
	return deps, hashes, err
}

func dependencyConfigHash(target BranchTarget, envKeys []string) string {
	data, _ := json.Marshal(map[string]any{
		"branch":       target.Branch,
		"target_file":  target.TargetFile,
		"target_files": normalizedTargetFiles(target),
		"env_keys":     envKeys,
	})
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func generatedOutputManifest(root, pageKey string) map[string]string {
	out := map[string]string{}
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if invalidOutputRelativePath(rel) {
			return nil
		}
		switch rel {
		case "index.html":
			out[rel] = pageKey
		case "assets/search-index.json":
			out[rel] = "search-index"
		case ".dependency_manifest.json":
			out[rel] = "manifest"
		default:
			out[rel] = "asset"
		}
		return nil
	})
	out[".dependency_manifest.json"] = "manifest"
	return out
}

func precheckRunnerTarget(cfg RunnerConfig, target BranchTarget) (string, error) {
	if err := ensureDiskAvailable(filepath.Dir(target.Out)); err != nil {
		return "", err
	}
	if cfg.RemoteBuild.Enabled {
		return "", nil
	}
	buildBin := os.Getenv("ADLAIRE_CI_BUILD_BIN")
	if buildBin == "" {
		buildBin = "/usr/local/bin/adlaire-ci-build"
	}
	info, err := os.Stat(buildBin)
	if err != nil {
		return "", err
	}
	if info.IsDir() || info.Mode()&0111 == 0 {
		return "", fmt.Errorf("build binary is not executable: %s", buildBin)
	}
	ctx, cancel := context.WithTimeout(context.Background(), runnerBuilderVersionTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, buildBin, "--version")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("build binary version timeout: %s", buildBin)
	}
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(stderr.String()) != "" {
		return "", fmt.Errorf("build binary version wrote stderr: %s", buildBin)
	}
	fields := strings.Fields(stdout.String())
	if len(fields) != 3 || fields[0] != "adlaire-ci-build" || fields[1] != runnerBinaryVersion || !strings.HasPrefix(fields[2], "go=") || fields[2] == "go=" {
		return "", fmt.Errorf("build binary version mismatch: %s", buildBin)
	}
	return fields[1], nil
}

func ensureDiskAvailable(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return err
	}
	free := stat.Bavail * uint64(stat.Bsize)
	if free < 64*1024*1024 {
		return fmt.Errorf("disk free below 64MiB: %s", path)
	}
	return nil
}

func runPipeline(cfg RunnerConfig, target BranchTarget, buildID, commitSHA, buildAt string, changedTargets []string) pipelineLog {
	if cfg.RemoteBuild.Enabled {
		return runRemoteBuildPipeline(cfg, target, buildID, commitSHA, buildAt, changedTargets)
	}
	pipelineConfig, err := readRunnerPipelineConfig(cfg.StateDir)
	if err != nil {
		code := 2
		stderr, stderrTruncated := trimLog(err.Error())
		return pipelineLog{
			ExitCode:        &code,
			Stdout:          "",
			Stderr:          stderr,
			StdoutTruncated: false,
			StderrTruncated: stderrTruncated,
			TargetStatus:    "failure_pipeline_config",
			ErrorMessage:    "pipeline config invalid",
			EnvKeys:         runnerEnvironmentKeysFrom(target.Env, nil),
		}
	}
	envKeys := runnerEnvironmentKeysFrom(target.Env, pipelineConfig.Env)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.BuildTimeoutSeconds)*time.Second)
	defer cancel()
	buildBin := os.Getenv("ADLAIRE_CI_BUILD_BIN")
	if buildBin == "" {
		buildBin = "/usr/local/bin/adlaire-ci-build"
	}
	args := []string{"--src", target.Src, "--out", target.Out, "--build-id", buildID, "--commit-sha", commitSHA, "--build-at", buildAt}
	if cfg.BuildCacheEnabled {
		args = append(args, "--cache-dir", cfg.StateDir)
	}
	args = append(args, pipelineConfig.ExtraArgs...)
	cmd := exec.CommandContext(ctx, buildBin, args...)
	cmd.Dir = filepath.Dir(target.Src)
	fixedEnv := map[string]string{
		"ADLAIRE_CI_SRC":         target.Src,
		"ADLAIRE_CI_OUT":         target.Out,
		"ADLAIRE_CI_BRANCH":      target.Branch,
		"ADLAIRE_CI_BUILD_ID":    buildID,
		"ADLAIRE_CI_TARGET_FILE": target.TargetFile,
		"ADLAIRE_CI_STATE_DIR":   cfg.StateDir,
	}
	if len(changedTargets) > 0 {
		data, _ := json.Marshal(changedTargets)
		fixedEnv["ADLAIRE_CHANGED_TARGETS"] = string(data)
	}
	cmd.Env = mergeRunnerProcessEnv(os.Environ(), target.Env, pipelineConfig.Env, fixedEnv)
	secrets := runnerSecretValues(target.Env)
	secrets = append(secrets, runnerSecretValues(pipelineConfig.Env)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		trimmedStdout, stdoutTruncated := trimLogWithSecrets(stdout.String(), secrets)
		trimmedStderr, stderrTruncated := trimLogWithSecrets(stderr.String(), secrets)
		return pipelineLog{ExitCode: nil, Stdout: trimmedStdout, Stderr: trimmedStderr, StdoutTruncated: stdoutTruncated, StderrTruncated: stderrTruncated, EnvKeys: envKeys}
	}
	code := 0
	if err != nil {
		code = 1
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			code = ee.ExitCode()
		}
	}
	trimmedStdout, stdoutTruncated := trimLogWithSecrets(stdout.String(), secrets)
	trimmedStderr, stderrTruncated := trimLogWithSecrets(stderr.String(), secrets)
	return pipelineLog{ExitCode: &code, Stdout: trimmedStdout, Stderr: trimmedStderr, StdoutTruncated: stdoutTruncated, StderrTruncated: stderrTruncated, EnvKeys: envKeys}
}

func runRemoteBuildPipeline(cfg RunnerConfig, target BranchTarget, buildID, commitSHA, buildAt string, changedTargets []string) pipelineLog {
	envKeys := runnerEnvironmentKeysFrom(target.Env, nil)
	started := runnerNow().UTC()
	remoteLog := newRemoteBuildLog(cfg.RemoteBuild)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.BuildTimeoutSeconds)*time.Second)
	defer cancel()
	remote := cfg.RemoteBuild
	host := remoteEndpoint(remote)
	env := map[string]string{
		"ADLAIRE_CI_SRC":         target.Src,
		"ADLAIRE_CI_OUT":         target.Out,
		"ADLAIRE_CI_BRANCH":      target.Branch,
		"ADLAIRE_CI_BUILD_ID":    buildID,
		"ADLAIRE_CI_COMMIT_SHA":  commitSHA,
		"ADLAIRE_CI_BUILD_AT":    buildAt,
		"ADLAIRE_CI_TARGET_FILE": target.TargetFile,
		"ADLAIRE_CI_STATE_DIR":   cfg.StateDir,
	}
	if len(changedTargets) > 0 {
		data, _ := json.Marshal(changedTargets)
		env["ADLAIRE_CHANGED_TARGETS"] = string(data)
	}
	for key, value := range target.Env {
		env[key] = value
	}
	cmd := exec.CommandContext(ctx, "ssh", host, remoteCommand(*remote.WorkDir, remote.CommandArgs, env))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		pl := timeoutPipelineLog(stdout.String(), stderr.String())
		pl.EnvKeys = envKeys
		finishRemoteBuildLog(remoteLog, started, nil, "remote build timeout")
		pl.RemoteBuild = remoteLog
		return pl
	}
	code := exitCodeFromError(err)
	if err == nil && remote.ArtifactPath != nil {
		scpCode, scpStdout, scpStderr := copyRemoteArtifact(ctx, host, *remote.ArtifactPath, target.Out)
		stdout.WriteString(scpStdout)
		stderr.WriteString(scpStderr)
		code = scpCode
		if scpCode == 0 {
			if size, ok := outputSizeBytes(target.Out); ok {
				remoteLog.ArtifactSizeBytes = &size
			}
			if count, ok := outputManifestFileCount(target.Out); ok {
				remoteLog.ManifestFileCount = &count
			}
		}
		if ctx.Err() == context.DeadlineExceeded {
			pl := timeoutPipelineLog(stdout.String(), stderr.String())
			pl.EnvKeys = envKeys
			finishRemoteBuildLog(remoteLog, started, nil, "remote artifact timeout")
			pl.RemoteBuild = remoteLog
			return pl
		}
	}
	trimmedStdout, stdoutTruncated := trimLog(stdout.String())
	trimmedStderr, stderrTruncated := trimLog(stderr.String())
	errorText := ""
	if code != 0 {
		errorText = "remote build failed"
	}
	finishRemoteBuildLog(remoteLog, started, &code, errorText)
	return pipelineLog{ExitCode: &code, Stdout: trimmedStdout, Stderr: trimmedStderr, StdoutTruncated: stdoutTruncated, StderrTruncated: stderrTruncated, EnvKeys: envKeys, RemoteBuild: remoteLog}
}

func newRemoteBuildLog(remote RemoteBuildConfig) *remoteBuildLog {
	host := derefString(remote.Host, "")
	user := derefString(remote.User, "")
	workDir := derefString(remote.WorkDir, "")
	commandName := ""
	if len(remote.CommandArgs) > 0 {
		commandName = filepath.Base(remote.CommandArgs[0])
	}
	return &remoteBuildLog{
		Host:              host,
		User:              user,
		WorkDirBasename:   filepath.Base(workDir),
		CommandName:       commandName,
		ExitCode:          nil,
		DurationSeconds:   0,
		ArtifactSizeBytes: nil,
		ManifestFileCount: nil,
		Status:            "failure",
		Error:             nil,
	}
}

func finishRemoteBuildLog(log *remoteBuildLog, started time.Time, exitCode *int, errorText string) {
	if log == nil {
		return
	}
	finished := runnerNow().UTC()
	duration := int64(finished.Sub(started).Seconds())
	if duration < 0 {
		duration = 0
	}
	log.DurationSeconds = duration
	log.ExitCode = exitCode
	if errorText == "" && exitCode != nil && *exitCode == 0 {
		log.Status = "success"
		log.Error = nil
		return
	}
	log.Status = "failure"
	if errorText == "" {
		errorText = "remote build failed"
	}
	log.Error = &errorText
}

func remoteEndpoint(remote RemoteBuildConfig) string {
	if remote.User != nil && *remote.User != "" {
		return *remote.User + "@" + *remote.Host
	}
	return *remote.Host
}

func remoteCommand(workDir string, args []string, env map[string]string) string {
	keys := make([]string, 0, len(env))
	for key := range env {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := []string{"cd -- " + shellQuote(workDir), "&&"}
	for _, key := range keys {
		parts = append(parts, key+"="+shellQuote(env[key]))
	}
	for _, arg := range args {
		parts = append(parts, shellQuote(arg))
	}
	return strings.Join(parts, " ")
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

func copyRemoteArtifact(ctx context.Context, host, artifactPath, out string) (int, string, string) {
	if err := os.RemoveAll(out); err != nil {
		return 1, "", err.Error()
	}
	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		return 1, "", err.Error()
	}
	cmd := exec.CommandContext(ctx, "scp", "-r", host+":"+artifactPath, out)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	return exitCodeFromError(cmd.Run()), stdout.String(), stderr.String()
}

func timeoutPipelineLog(stdout, stderr string) pipelineLog {
	trimmedStdout, stdoutTruncated := trimLog(stdout)
	trimmedStderr, stderrTruncated := trimLog(stderr)
	return pipelineLog{ExitCode: nil, Stdout: trimmedStdout, Stderr: trimmedStderr, StdoutTruncated: stdoutTruncated, StderrTruncated: stderrTruncated}
}

func exitCodeFromError(err error) int {
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return 1
}

func readRunnerHooks(stateDir string) ([]runnerHook, []string, error) {
	path := filepath.Join(stateDir, ".hooks")
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, []string{}, nil
		}
		return nil, nil, err
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	var cfg hooksFile
	if err := dec.Decode(&cfg); err != nil {
		return nil, nil, err
	}
	if len(cfg.Hooks) > 50 {
		return nil, nil, errors.New("hooks count exceeds 50")
	}
	seen := map[string]bool{}
	skipped := []string{}
	for i := range cfg.Hooks {
		hook := &cfg.Hooks[i]
		if !validRunnerHookID(hook.ID) {
			return nil, nil, fmt.Errorf("invalid hook id: %s", hook.ID)
		}
		if seen[hook.ID] {
			return nil, nil, fmt.Errorf("duplicate hook id: %s", hook.ID)
		}
		seen[hook.ID] = true
		if hook.Phase != "pre" && hook.Phase != "post" {
			return nil, nil, fmt.Errorf("invalid hook phase: %s", hook.ID)
		}
		if err := validateRunnerHookCommand(hook.CommandArgs); err != nil {
			return nil, nil, fmt.Errorf("invalid hook command: %s", hook.ID)
		}
		if hook.TimeoutSeconds < 1 || hook.TimeoutSeconds > 3600 {
			return nil, nil, fmt.Errorf("invalid hook timeout: %s", hook.ID)
		}
		if _, err := time.Parse(time.RFC3339, hook.CreatedAt); err != nil {
			return nil, nil, fmt.Errorf("invalid hook created_at: %s", hook.ID)
		}
		if !hook.Enabled {
			skipped = append(skipped, hook.ID)
		}
	}
	sort.Slice(cfg.Hooks, func(i, j int) bool { return cfg.Hooks[i].ID < cfg.Hooks[j].ID })
	sort.Strings(skipped)
	return cfg.Hooks, skipped, nil
}

func validRunnerHookID(id string) bool {
	if len(id) < 1 || len(id) > 64 {
		return false
	}
	for _, r := range id {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func validateRunnerHookCommand(args []string) error {
	if len(args) < 1 || len(args) > 20 {
		return errors.New("command_args length invalid")
	}
	for i, arg := range args {
		limit := 500
		if i == 0 {
			limit = 256
		}
		if arg == "" || len(arg) > limit || strings.ContainsAny(arg, "\x00\n\r") {
			return fmt.Errorf("command_args[%d] invalid", i)
		}
	}
	return nil
}

func executeRunnerHooks(cfg RunnerConfig, target BranchTarget, buildID, phase, commitSHA, buildAt string, changedTargets []string, hooks []runnerHook, logger *slog.Logger) ([]string, *string, bool) {
	warnings := []string{}
	exitBump := false
	for _, hook := range hooks {
		if !hook.Enabled || hook.Phase != phase {
			continue
		}
		result := runRunnerHook(cfg, target, buildID, commitSHA, buildAt, changedTargets, hook)
		if err := writeRunnerHookLog(cfg.StateDir, result); err != nil {
			logger.Error("HOOK_LOG_WRITE_FAILED: " + err.Error())
			if phase == "pre" {
				msg := "hook log write failed"
				return warnings, &msg, true
			}
			warnings = append(warnings, "HOOK_LOG_WRITE_FAILED:"+hook.ID)
			exitBump = true
			continue
		}
		if result.Status == "success" {
			continue
		}
		if phase == "pre" && hook.AbortOnFailure {
			msg := "pre hook failed"
			return warnings, &msg, true
		}
		if phase == "post" {
			warnings = append(warnings, "POST_HOOK_FAILED:"+hook.ID)
			exitBump = true
		} else {
			warnings = append(warnings, "PRE_HOOK_FAILED:"+hook.ID)
		}
	}
	return warnings, nil, exitBump
}

func runRunnerHook(cfg RunnerConfig, target BranchTarget, buildID, commitSHA, buildAt string, changedTargets []string, hook runnerHook) hookRunLog {
	started := runnerNow().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hook.TimeoutSeconds)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, hook.CommandArgs[0], hook.CommandArgs[1:]...)
	cmd.Dir = filepath.Dir(target.Src)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Env = runnerHookEnv(cfg, target, buildID, commitSHA, buildAt, changedTargets, hook)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	timedOut := ctx.Err() == context.DeadlineExceeded
	if timedOut && cmd.Process != nil {
		if killErr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); killErr != nil && !errors.Is(killErr, os.ErrProcessDone) && !errors.Is(killErr, syscall.ESRCH) {
			stderr.WriteString("\nHOOK_PROCESS_GROUP_KILL_FAILED")
		}
	}
	status := "success"
	exitCode := (*int)(nil)
	if timedOut {
		status = "timeout"
	} else {
		code := exitCodeFromError(err)
		exitCode = &code
		if code != 0 {
			status = "failure"
		}
	}
	finished := runnerNow().UTC()
	stdoutText, stdoutTruncated := trimHookOutput(stdout.String(), runnerSecretValues(target.Env))
	stderrText, stderrTruncated := trimHookOutput(stderr.String(), runnerSecretValues(target.Env))
	duration := int64(finished.Sub(started).Seconds())
	if duration < 0 {
		duration = 0
	}
	return hookRunLog{
		HookID: hook.ID, BuildID: buildID, Phase: hook.Phase, Status: status,
		StartedAt: started.Format(time.RFC3339), FinishedAt: finished.Format(time.RFC3339), DurationSeconds: duration,
		Stdout: stdoutText, Stderr: stderrText, ExitCode: exitCode, TimedOut: timedOut,
		Truncated: stdoutTruncated || stderrTruncated,
	}
}

func runnerHookEnv(cfg RunnerConfig, target BranchTarget, buildID, commitSHA, buildAt string, changedTargets []string, hook runnerHook) []string {
	env := append(os.Environ(),
		"ADLAIRE_CI_SRC="+target.Src,
		"ADLAIRE_CI_OUT="+target.Out,
		"ADLAIRE_CI_BRANCH="+target.Branch,
		"ADLAIRE_CI_BUILD_ID="+buildID,
		"ADLAIRE_CI_COMMIT_SHA="+commitSHA,
		"ADLAIRE_CI_BUILD_AT="+buildAt,
		"ADLAIRE_CI_TARGET_FILE="+target.TargetFile,
		"ADLAIRE_CI_STATE_DIR="+cfg.StateDir,
		"ADLAIRE_CI_HOOK_ID="+hook.ID,
		"ADLAIRE_CI_HOOK_PHASE="+hook.Phase,
	)
	if len(changedTargets) > 0 {
		data, _ := json.Marshal(changedTargets)
		env = append(env, "ADLAIRE_CHANGED_TARGETS="+string(data))
	}
	for key, value := range target.Env {
		env = append(env, key+"="+value)
	}
	return env
}

func runnerSecretValues(env map[string]string) []string {
	secrets := []string{}
	for key, value := range env {
		lower := strings.ToLower(key)
		if value != "" && (strings.Contains(lower, "secret") || strings.Contains(lower, "token") || strings.Contains(lower, "password")) {
			secrets = append(secrets, value)
		}
	}
	return secrets
}

func mergeRunnerProcessEnv(base []string, overlays ...map[string]string) []string {
	order := []string{}
	values := map[string]string{}
	for _, entry := range base {
		key, value, ok := strings.Cut(entry, "=")
		if !ok || key == "" {
			continue
		}
		if _, exists := values[key]; !exists {
			order = append(order, key)
		}
		values[key] = value
	}
	for _, overlay := range overlays {
		keys := make([]string, 0, len(overlay))
		for key := range overlay {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if _, exists := values[key]; !exists {
				order = append(order, key)
			}
			values[key] = overlay[key]
		}
	}
	out := make([]string, 0, len(order))
	for _, key := range order {
		out = append(out, key+"="+values[key])
	}
	return out
}

func trimHookOutput(s string, secrets []string) (string, bool) {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.ReplaceAll(s, "\x00", "\\u0000")
	s = strings.ToValidUTF8(s, "\uFFFD")
	for _, secret := range secrets {
		s = strings.ReplaceAll(s, secret, "***")
	}
	const max = 64 * 1024
	if len(s) <= max {
		return strings.TrimSuffix(s, "\n"), false
	}
	start := len(s) - max
	for start < len(s) && !utf8.RuneStart(s[start]) {
		start++
	}
	return strings.TrimSuffix(s[start:], "\n"), true
}

func writeRunnerHookLog(stateDir string, log hookRunLog) error {
	return runnerAtomicWriteJSON(filepath.Join(stateDir, ".build_logs", log.BuildID+"_hook_"+log.HookID+".json"), log, 0600)
}

func shouldRetryPipeline(pl pipelineLog) bool {
	if pl.TargetStatus == "failure_pipeline_config" {
		return false
	}
	return pl.ExitCode == nil
}

func parseRunnerReport(stdout string) (*runnerReport, []string) {
	var report *runnerReport
	warnings := []string{}
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "[WARN] ") {
			warnings = append(warnings, strings.TrimPrefix(line, "[WARN] "))
		}
		if strings.HasPrefix(line, "[REPORT] ") {
			values := map[string]string{}
			for _, part := range strings.Fields(strings.TrimPrefix(line, "[REPORT] ")) {
				k, v, ok := strings.Cut(part, "=")
				if ok {
					values[k] = v
				}
			}
			if report != nil {
				warnings = append(warnings, "REPORT_DUPLICATE")
				continue
			}
			report = &runnerReport{
				Pages:           atoi(values["pages"]),
				Headings:        atoi(values["headings"]),
				TablesCount:     atoi(values["tables"]),
				CodeBlocksCount: atoi(values["code_blocks"]),
				WarningsCount:   atoi(values["warnings"]),
				SizeWarn:        values["size_warn"] == "true",
				BrokenLinks:     atoi(values["broken_links"]),
				HeadingSkips:    atoi(values["heading_skips"]),
				ReadingTime:     atoi(values["reading_time"]),
				Theme:           values["theme"],
			}
		}
	}
	return report, warnings
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func trimLog(s string) (string, bool) {
	return trimLogWithSecrets(s, nil)
}

func trimLogWithSecrets(s string, secrets []string) (string, bool) {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.ReplaceAll(s, "\x00", "\\u0000")
	s = strings.ToValidUTF8(s, "\uFFFD")
	for _, secret := range secrets {
		s = strings.ReplaceAll(s, secret, "***")
	}
	const max = 1024 * 1024
	if len(s) <= max {
		return strings.TrimSuffix(s, "\n"), false
	}
	start := len(s) - max
	for start < len(s) && !utf8.RuneStart(s[start]) {
		start++
	}
	return strings.TrimSuffix(s[start:], "\n"), true
}

func validateOutputSite(root string) error {
	required := []string{"index.html", "assets/style.css", "assets/app.js", "assets/search-index.json"}
	for _, rel := range required {
		path := filepath.Join(root, rel)
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("invalid output file: %s", rel)
		}
	}
	return nil
}

func executeDeploy(cfg RunnerConfig, branchIdx, deployIdx int, target BranchTarget, d DeployTarget, buildID string) deployLog {
	started := runnerNow().UTC()
	targetID := d.ID
	if targetID == "" {
		targetID = fmt.Sprintf("%d-%d", branchIdx, deployIdx)
	}
	dl := deployLog{TargetID: targetID, Host: d.Host, User: d.User, DestDir: d.DestDir, Status: "success", TransferVerified: true, StartedAt: started.Format(time.RFC3339)}
	total, uploaded, skipped, bytesUploaded, err := deploySite(target.Out, d)
	dl.FilesTotal = total
	dl.FilesUploaded = uploaded
	dl.FilesSkipped = skipped
	dl.BytesUploaded = bytesUploaded
	dl.FinishedAt = runnerNow().UTC().Format(time.RFC3339)
	if err != nil {
		errStr := err.Error()
		dl.Status = "pending"
		dl.TransferVerified = false
		dl.Error = &errStr
		outputSHA, _ := outputManifestSHA(target.Out)
		out := target.Out
		_ = addPendingTransfer(cfg.PendingFile, pendingTransfer{
			BuildID: buildID, Trigger: "deploy", SourceKind: "output",
			RollbackFrom: nil, SnapshotID: nil, Branch: target.Branch,
			TargetID: targetID, Out: &out, Host: d.Host, User: d.User, DestDir: d.DestDir,
			OutputSHA256: derefString(outputSHA, ""), FailedAt: runnerNow().UTC().Format(time.RFC3339),
			RetryCount: 0, LastError: errStr,
		})
	}
	return dl
}

func executeDeploys(cfg RunnerConfig, branchIdx int, target BranchTarget, buildID string) []deployLog {
	if len(target.DeployTargets) == 0 {
		return []deployLog{}
	}
	if cfg.DeployParallelism <= 1 || len(target.DeployTargets) == 1 {
		out := make([]deployLog, 0, len(target.DeployTargets))
		for deployIdx, d := range target.DeployTargets {
			out = append(out, executeDeploy(cfg, branchIdx, deployIdx, target, d, buildID))
		}
		return out
	}
	workers := cfg.DeployParallelism
	if workers > len(target.DeployTargets) {
		workers = len(target.DeployTargets)
	}
	type job struct {
		index  int
		target DeployTarget
	}
	jobs := make(chan job)
	results := make([]deployLog, len(target.DeployTargets))
	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				results[item.index] = executeDeploy(cfg, branchIdx, item.index, target, item.target, buildID)
			}
		}()
	}
	for i, deployTarget := range target.DeployTargets {
		jobs <- job{index: i, target: deployTarget}
	}
	close(jobs)
	wg.Wait()
	return results
}

func buildTargetResults(deploys []deployLog) []targetResultLog {
	results := make([]targetResultLog, 0, len(deploys))
	for _, deploy := range deploys {
		result := targetResultLog{
			TargetID:   deploy.TargetID,
			Status:     "success",
			StartedAt:  deploy.StartedAt,
			FinishedAt: deploy.FinishedAt,
			ErrorCode:  nil,
			Error:      nil,
		}
		if deploy.Status != "success" {
			code := "deploy_ssh_error"
			message := "deploy ssh failed"
			result.Status = "failure"
			result.ErrorCode = &code
			result.Error = &message
		}
		results = append(results, result)
	}
	return results
}

func buildTransferVerified(deploys []deployLog) *bool {
	if len(deploys) == 0 {
		return nil
	}
	verified := true
	for _, deploy := range deploys {
		if deploy.Status != "success" || !deploy.TransferVerified {
			verified = false
			break
		}
	}
	return &verified
}

func buildCommitStatusLog(cfg RunnerConfig, state *string, at time.Time) *commitStatusLog {
	if !cfg.CommitStatusEnabled {
		return nil
	}
	var sentAt *string
	if state != nil {
		value := at.UTC().Format(time.RFC3339)
		sentAt = &value
	}
	return &commitStatusLog{
		Enabled:    cfg.CommitStatusEnabled,
		State:      state,
		Context:    cfg.CommitStatusContext,
		TargetURL:  cfg.CommitStatusTargetURL,
		SentAt:     sentAt,
		HTTPStatus: nil,
		Error:      nil,
	}
}

func buildBuildMeta(report *runnerReport) *buildMetaLog {
	if report == nil {
		return nil
	}
	return &buildMetaLog{
		BuildID:   report.BuildID,
		CommitSHA: report.CommitSHA,
		BuildAt:   report.BuildAt,
	}
}

func deploySite(out string, d DeployTarget) (int, int, int, int64, error) {
	total, uploaded, skipped := 0, 0, 0
	var bytesUploaded int64
	err := filepath.WalkDir(out, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		total++
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("invalid deploy file type: %s", path)
		}
		rel, err := filepath.Rel(out, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." || strings.HasPrefix(rel, "../") || strings.Contains(rel, "/../") {
			return fmt.Errorf("invalid deploy path: %s", rel)
		}
		localSHA, err := runnerFileSHA256(path)
		if err != nil {
			return err
		}
		remotePath := remoteJoin(d.DestDir, rel)
		if remoteSHA256(d, remotePath) == localSHA {
			skipped++
			return nil
		}
		if err := uploadFileSSH(d, path, remotePath); err != nil {
			return err
		}
		if remoteSHA256(d, remotePath) != localSHA {
			return errors.New("checksum mismatch")
		}
		uploaded++
		bytesUploaded += info.Size()
		return nil
	})
	return total, uploaded, skipped, bytesUploaded, err
}

func runnerFileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func remoteSHA256(d DeployTarget, remotePath string) string {
	cmd := exec.Command("ssh", d.User+"@"+d.Host, "sha256sum -- "+shellQuote(remotePath))
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func uploadFileSSH(d DeployTarget, localPath, remotePath string) error {
	data, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}
	remoteDir := filepath.ToSlash(filepath.Dir(remotePath))
	remoteCmd := "mkdir -p -- " + shellQuote(remoteDir) + " && tee -- " + shellQuote(remotePath) + " >/dev/null"
	cmd := exec.Command("ssh", d.User+"@"+d.Host, remoteCmd)
	cmd.Stdin = bytes.NewReader(data)
	return cmd.Run()
}

func remoteJoin(root, rel string) string {
	return strings.TrimRight(root, "/") + "/" + strings.TrimLeft(rel, "/")
}

func writeSnapshot(cfg RunnerConfig, target BranchTarget, buildID string) (string, error) {
	if cfg.SnapshotsKeep == 0 {
		return "", errors.New("snapshot disabled")
	}
	root := filepath.Join(cfg.StateDir, ".snapshots")
	dest := filepath.Join(root, buildID)
	if err := os.MkdirAll(root, 0700); err != nil {
		return "", err
	}
	if _, err := os.Stat(dest); err == nil {
		return buildID, nil
	}
	tmp := filepath.Join(root, "."+buildID+".tmp."+strconv.Itoa(os.Getpid()))
	if err := os.RemoveAll(tmp); err != nil {
		return "", err
	}
	if err := os.Mkdir(tmp, 0700); err != nil {
		return "", err
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	meta, err := createSnapshotArchive(target.Out, filepath.Join(tmp, "site.tar.gz"), buildID)
	if err != nil {
		return "", err
	}
	if err := writeSnapshotMeta(filepath.Join(tmp, "meta.json"), meta); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, dest); err != nil {
		return "", err
	}
	if err := pruneSnapshots(root, cfg.SnapshotsKeep); err != nil {
		return "", err
	}
	return buildID, nil
}

type snapshotMeta struct {
	ID           string `json:"id"`
	BuildID      string `json:"build_id"`
	SavedAt      string `json:"saved_at"`
	SizeBytes    int64  `json:"size_bytes"`
	FileCount    int    `json:"file_count"`
	OutputSHA256 string `json:"output_sha256"`
}

type snapshotFile struct {
	path string
	rel  string
	size int64
}

func createSnapshotArchive(outputRoot, archivePath, buildID string) (snapshotMeta, error) {
	outputSHA, err := outputManifestSHA(outputRoot)
	if err != nil {
		return snapshotMeta{}, err
	}
	files, size, err := snapshotFileList(outputRoot)
	if err != nil {
		return snapshotMeta{}, err
	}
	f, err := os.OpenFile(archivePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return snapshotMeta{}, err
	}
	gz, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		_ = f.Close()
		return snapshotMeta{}, err
	}
	tw := tar.NewWriter(gz)
	for _, item := range files {
		header := &tar.Header{Name: item.rel, Mode: 0644, Size: item.size, Typeflag: tar.TypeReg, Format: tar.FormatUSTAR, ModTime: time.Unix(0, 0).UTC()}
		if err := tw.WriteHeader(header); err != nil {
			_ = f.Close()
			return snapshotMeta{}, err
		}
		in, err := os.Open(item.path)
		if err != nil {
			_ = f.Close()
			return snapshotMeta{}, err
		}
		if _, err := io.Copy(tw, in); err != nil {
			_ = in.Close()
			_ = f.Close()
			return snapshotMeta{}, err
		}
		if err := in.Close(); err != nil {
			_ = f.Close()
			return snapshotMeta{}, err
		}
	}
	if err := tw.Close(); err != nil {
		_ = f.Close()
		return snapshotMeta{}, err
	}
	if err := gz.Close(); err != nil {
		_ = f.Close()
		return snapshotMeta{}, err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return snapshotMeta{}, err
	}
	if err := f.Close(); err != nil {
		return snapshotMeta{}, err
	}
	return snapshotMeta{
		ID:           buildID,
		BuildID:      buildID,
		SavedAt:      runnerNow().UTC().Format(time.RFC3339),
		SizeBytes:    size,
		FileCount:    len(files),
		OutputSHA256: derefString(outputSHA, ""),
	}, nil
}

func snapshotFileList(root string) ([]snapshotFile, int64, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return nil, 0, err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, 0, fmt.Errorf("invalid snapshot root: %s", root)
	}
	files := []snapshotFile{}
	var total int64
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("invalid snapshot file: %s", path)
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok && stat.Nlink > 1 {
			return fmt.Errorf("hardlink snapshot file: %s", path)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if invalidOutputRelativePath(rel) || snapshotSecretPath(rel) {
			return fmt.Errorf("invalid snapshot path: %s", rel)
		}
		files = append(files, snapshotFile{path: path, rel: rel, size: info.Size()})
		total += info.Size()
		return nil
	})
	sort.Slice(files, func(i, j int) bool { return files[i].rel < files[j].rel })
	return files, total, err
}

func snapshotSecretPath(rel string) bool {
	switch filepath.ToSlash(rel) {
	case ".git", ".github_token", ".admin_credentials", ".api_tokens", ".smtp_secret", ".webhook_secret":
		return true
	}
	return strings.HasPrefix(rel, ".git/")
}

func writeSnapshotMeta(path string, meta snapshotMeta) error {
	data := fmt.Sprintf("{\"id\":%q,\"build_id\":%q,\"saved_at\":%q,\"size_bytes\":%d,\"file_count\":%d,\"output_sha256\":%q}\n",
		meta.ID, meta.BuildID, meta.SavedAt, meta.SizeBytes, meta.FileCount, meta.OutputSHA256)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(data)); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		out := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(out, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(out, data, 0644)
	})
}

func pruneSnapshots(root string, keep int) error {
	if keep <= 0 {
		return nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	dirs := []string{}
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)
	for len(dirs) > keep {
		if err := os.RemoveAll(filepath.Join(root, dirs[0])); err != nil {
			return err
		}
		dirs = dirs[1:]
	}
	return nil
}

func addPendingTransfer(path string, entry pendingTransfer) error {
	var entries []pendingTransfer
	_ = readJSONArray(path, &entries)
	for i := range entries {
		if entries[i].Branch == entry.Branch && entries[i].TargetID == entry.TargetID && entries[i].Host == entry.Host && entries[i].User == entry.User && entries[i].DestDir == entry.DestDir {
			entries[i].BuildID = entry.BuildID
			entries[i].Trigger = entry.Trigger
			entries[i].SourceKind = entry.SourceKind
			entries[i].RollbackFrom = entry.RollbackFrom
			entries[i].SnapshotID = entry.SnapshotID
			entries[i].Out = entry.Out
			entries[i].OutputSHA256 = entry.OutputSHA256
			entries[i].RetryCount = 0
			entries[i].FailedAt = entry.FailedAt
			entries[i].LastError = entry.LastError
			return runnerAtomicWriteJSON(path, entries, 0600)
		}
	}
	entries = append(entries, entry)
	return runnerAtomicWriteJSON(path, entries, 0600)
}

func retryPendingTransfers(cfg *RunnerConfig, logger *slog.Logger) error {
	var entries []pendingTransfer
	if err := readJSONArray(cfg.PendingFile, &entries); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err := backupCorruptJSON(cfg.PendingFile, logger); err != nil {
			return err
		}
		return runnerAtomicWriteJSON(cfg.PendingFile, []pendingTransfer{}, 0600)
	}
	remaining := []pendingTransfer{}
	for _, entry := range entries {
		if entry.Out == nil {
			remaining = append(remaining, entry)
			continue
		}
		deploy := DeployTarget{Host: entry.Host, User: entry.User, DestDir: entry.DestDir}
		_, _, _, _, err := deploySite(*entry.Out, deploy)
		if err == nil {
			logger.Info("PENDING RETRY OK site -> " + entry.Host)
			continue
		}
		entry.RetryCount++
		entry.FailedAt = runnerNow().UTC().Format(time.RFC3339)
		entry.LastError = err.Error()
		remaining = append(remaining, entry)
		logger.Error("PENDING RETRY FAILED site: " + entry.Host)
	}
	return runnerAtomicWriteJSON(cfg.PendingFile, remaining, 0600)
}

func retryNotifyPending(path string, logger *slog.Logger) error {
	var entries []notifyPendingEntry
	if err := readJSONArray(path, &entries); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	remaining := []notifyPendingEntry{}
	for _, entry := range entries {
		if err := postNotify(entry.URL, entry.Payload); err == nil {
			logger.Info("NOTIFY_PENDING_RETRY_OK: url=" + entry.URL)
			continue
		} else {
			entry.LastError = err.Error()
		}
		entry.RetryCount++
		remaining = append(remaining, entry)
		logger.Error("NOTIFY_PENDING_RETRY_FAILED: url=" + entry.URL)
	}
	return runnerAtomicWriteJSON(path, remaining, 0600)
}

func sendBuildNotifications(cfg RunnerConfig, log buildLog, logger *slog.Logger) {
	config, err := readNotifyConfig(filepath.Join(cfg.StateDir, ".notify_config"))
	if err != nil {
		return
	}
	event := "failure"
	if log.TargetStatus == "success" {
		event = "success"
	} else if log.TargetStatus == "success_deploy_pending" {
		event = "deploy_failure"
	}
	payload := map[string]any{
		"event": event, "build_id": log.ID, "status": log.TargetStatus,
		"branch": log.Branch, "target_file": log.TargetFile,
		"target_files": log.TargetFiles, "changed_targets": log.ChangedTargets,
		"trigger": log.Trigger, "created_at": runnerNow().UTC().Format(time.RFC3339),
	}
	_ = sendRunnerNotificationEvent(cfg, config, event, payload, logger)
}

func readNotifyConfig(path string) (notifyConfigFile, error) {
	var config notifyConfigFile
	err := runnerReadJSONFile(path, &config)
	if errors.Is(err, os.ErrNotExist) {
		return config, nil
	}
	return config, err
}

func sendWeeklySummaryIfDue(cfg RunnerConfig, logger *slog.Logger) error {
	config, err := readNotifyConfig(filepath.Join(cfg.StateDir, ".notify_config"))
	if err != nil {
		return err
	}
	if len(notifyTargetsForEvent(config, "weekly_summary")) == 0 {
		return nil
	}
	summary := config.Summary
	enabled := cfg.WeeklySummaryEnabled || summary.Enabled
	if !enabled {
		return nil
	}
	if summary.Interval != "" && summary.Interval != "weekly" {
		return nil
	}
	now := runnerNow().UTC()
	day := cfg.WeeklySummaryDay
	hour := cfg.WeeklySummaryHour
	if summary.Enabled {
		day = summary.DayOfWeek
		hour = summary.Hour
	}
	if int(now.Weekday()) != day || now.Hour() != hour {
		return nil
	}
	state, err := runnerReadBuildState(cfg.StateDir)
	if err != nil {
		state = runnerDefaultBuildState()
	}
	sentDate := now.Format("2006-01-02")
	if state.WeeklySummarySentDate != nil && *state.WeeklySummarySentDate == sentDate {
		return nil
	}
	records, err := readHistoryRecords(cfg.StateDir)
	if err != nil {
		return err
	}
	payload := weeklySummaryPayload(records, now)
	if sendRunnerNotificationEvent(cfg, config, "weekly_summary", payload, logger) {
		sentAt := now.Format(time.RFC3339)
		state.WeeklySummaryLastSentAt = &sentAt
		state.WeeklySummarySentDate = &sentDate
		return runnerAtomicWriteJSON(filepath.Join(cfg.StateDir, ".build_state"), state, 0600)
	}
	return nil
}

func readHistoryRecords(stateDir string) ([]historyRecord, error) {
	data, err := os.ReadFile(filepath.Join(stateDir, ".build_history"))
	if errors.Is(err, os.ErrNotExist) {
		return []historyRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	records := []historyRecord{}
	for _, line := range bytes.Split(data, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var rec historyRecord
		if err := json.Unmarshal(line, &rec); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func weeklySummaryPayload(records []historyRecord, now time.Time) map[string]any {
	since := now.AddDate(0, 0, -7)
	counts := map[string]int{"success": 0, "failure": 0}
	var totalDuration int64
	durationCount := 0
	var maxDuration *int64
	var maxDurationBuildID *string
	for _, rec := range records {
		finished, err := time.Parse(time.RFC3339, rec.FinishedAt)
		if err != nil || finished.Before(since) || finished.After(now) {
			continue
		}
		switch {
		case rec.Status == "success":
			counts["success"]++
		case rec.Status == "success_deploy_pending":
			counts["success"]++
		default:
			counts["failure"]++
		}
		if rec.DurationSeconds >= 0 {
			totalDuration += rec.DurationSeconds
			durationCount++
			if maxDuration == nil || rec.DurationSeconds > *maxDuration {
				value := rec.DurationSeconds
				maxDuration = &value
				maxDurationBuildID = &rec.ID
			}
		}
	}
	var avg *float64
	if durationCount > 0 {
		value := math.Round((float64(totalDuration)/float64(durationCount))*100) / 100
		avg = &value
	}
	var maxValue any
	if maxDuration != nil {
		maxValue = *maxDuration
	}
	total := counts["success"] + counts["failure"]
	successRate := float64(0)
	if total > 0 {
		successRate = math.Round((float64(counts["success"])/float64(total))*10000) / 100
	}
	return map[string]any{
		"event":                 "weekly_summary",
		"period_days":           7,
		"period_from":           since.Format(time.RFC3339),
		"period_to":             now.Format(time.RFC3339),
		"success_count":         counts["success"],
		"failure_count":         counts["failure"],
		"success_rate":          successRate,
		"avg_duration_seconds":  avg,
		"max_duration_seconds":  maxValue,
		"max_duration_build_id": maxDurationBuildID,
	}
}

func hookHandlesEvent(hook notifyWebhook, event string) bool {
	return eventListHandles(hook.On, nil, event)
}

type notifyTarget struct {
	ChannelID            string
	ChannelType          string
	URL                  string
	Config               map[string]any
	RetryCount           int
	RetryIntervalSeconds int
}

func notifyTargetsForEvent(config notifyConfigFile, event string) []notifyTarget {
	targets := []notifyTarget{}
	for i, hook := range config.Webhooks {
		if !hook.Enabled || hook.URL == "" || !hookHandlesEvent(hook, event) {
			continue
		}
		id := hook.Label
		if id == "" {
			id = fmt.Sprintf("webhook-%d", i+1)
		}
		targets = append(targets, notifyTarget{ChannelID: id, ChannelType: "webhook", URL: hook.URL, Config: map[string]any{"url": hook.URL}, RetryCount: 2, RetryIntervalSeconds: 30})
	}
	for _, channel := range config.Channels {
		if !channel.Enabled || !eventListHandles(channel.On, config.On, event) {
			continue
		}
		target := notifyTarget{
			ChannelID:            channel.ID,
			ChannelType:          channel.Type,
			Config:               channel.Config,
			RetryCount:           notifyRetryCount(channel),
			RetryIntervalSeconds: notifyRetryInterval(channel),
		}
		if target.ChannelID == "" {
			target.ChannelID = "n000000"
		}
		if target.ChannelType == "webhook" {
			rawURL, ok := channel.Config["url"].(string)
			if !ok || rawURL == "" {
				continue
			}
			target.URL = rawURL
		}
		if target.ChannelType != "webhook" && target.ChannelType != "email" && target.ChannelType != "command" {
			continue
		}
		targets = append(targets, target)
	}
	sort.SliceStable(targets, func(i, j int) bool { return targets[i].ChannelID < targets[j].ChannelID })
	return targets
}

func notifyRetryCount(channel notifyChannel) int {
	if channel.RetryCount == nil {
		return 2
	}
	if *channel.RetryCount < 0 {
		return 0
	}
	if *channel.RetryCount > 10 {
		return 10
	}
	return *channel.RetryCount
}

func notifyRetryInterval(channel notifyChannel) int {
	if channel.RetryIntervalSeconds == nil || *channel.RetryIntervalSeconds <= 0 {
		return 30
	}
	return *channel.RetryIntervalSeconds
}

func eventListHandles(primary, fallback []string, event string) bool {
	items := primary
	if len(items) == 0 {
		items = fallback
	}
	for _, item := range items {
		if item == event || item == "*" {
			return true
		}
	}
	return false
}

func sendRunnerNotificationEvent(cfg RunnerConfig, config notifyConfigFile, event string, payload map[string]any, logger *slog.Logger) bool {
	sent := false
	for _, target := range notifyTargetsForEvent(config, event) {
		canonicalPayload := normalizeNotificationPayload(event, payload)
		payloadSHA := notificationPayloadSHA(canonicalPayload)
		result := deliverNotification(cfg, target, canonicalPayload, payloadSHA)
		now := runnerNow().UTC()
		record := map[string]any{
			"id":             nextNotifyRecordID("ntfy", now),
			"at":             now.Format(time.RFC3339),
			"event":          event,
			"channel_id":     target.ChannelID,
			"channel_type":   target.ChannelType,
			"payload_sha256": payloadSHA,
			"result":         result.Result,
			"attempt":        1,
			"http_status":    result.HTTPStatus,
			"error_code":     result.ErrorCode,
			"error":          result.Error,
		}
		if err := appendNotifyLog(cfg.StateDir, record); err != nil {
			logger.Error("NOTIFY_LOG_WRITE_FAILED: channel_id=" + target.ChannelID + " attempt=1")
			continue
		}
		if result.Result == "success" {
			sent = true
			continue
		}
		logger.Error("NOTIFY_FAILED: channel_id=" + target.ChannelID)
		if result.Retryable && target.RetryCount > 0 {
			_ = addNotifyPending(filepath.Join(cfg.StateDir, ".notify_pending"), notifyPendingEntry{
				ID:          nextNotifyRecordID("np", now),
				Event:       event,
				ChannelID:   target.ChannelID,
				ChannelType: "webhook",
				URL:         target.URL,
				Payload:     canonicalPayload,
				QueuedAt:    now.Format(time.RFC3339),
				Attempts:    1,
				NextAttempt: now.Add(time.Duration(target.RetryIntervalSeconds) * time.Second).Format(time.RFC3339),
				CreatedAt:   now.Format(time.RFC3339),
				PayloadSHA:  payloadSHA,
				RetryCount:  1,
				LastError:   result.PendingError,
			})
		}
	}
	return sent
}

func postNotify(url string, payload map[string]any) error {
	result := sendWebhookNotification(notifyTarget{ChannelID: "legacy", ChannelType: "webhook", URL: url}, payload, "")
	if result.Result == "success" {
		return nil
	}
	return errors.New(derefString(result.Error, "notification failed"))
}

func normalizeNotificationPayload(event string, payload map[string]any) map[string]any {
	out := map[string]any{
		"event":      event,
		"build_id":   stringOrNil(payload["build_id"]),
		"status":     stringOrNil(payload["status"]),
		"branch":     stringOrNil(payload["branch"]),
		"trigger":    stringOrNil(payload["trigger"]),
		"created_at": stringOrNil(payload["created_at"]),
	}
	if out["created_at"] == nil {
		out["created_at"] = runnerNow().UTC().Format(time.RFC3339)
	}
	if event == "weekly_summary" {
		for _, key := range []string{"period_days", "period_from", "period_to", "success_count", "failure_count", "success_rate", "avg_duration_seconds", "max_duration_seconds", "max_duration_build_id"} {
			out[key] = payload[key]
		}
	}
	return out
}

func stringOrNil(v any) any {
	if s, ok := v.(string); ok {
		return s
	}
	return nil
}

func notificationPayloadSHA(payload map[string]any) string {
	data, _ := json.Marshal(payload)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func deliverNotification(cfg RunnerConfig, target notifyTarget, payload map[string]any, payloadSHA string) notifyDeliveryResult {
	switch target.ChannelType {
	case "webhook":
		return sendWebhookNotification(target, payload, payloadSHA)
	case "email":
		return sendEmailNotification(cfg, target, payload)
	case "command":
		return sendCommandNotification(target, payload, payloadSHA)
	default:
		code := "command_error"
		msg := "unsupported notification channel"
		return notifyDeliveryResult{Result: "failure", ErrorCode: &code, Error: &msg}
	}
}

func sendWebhookNotification(target notifyTarget, payload map[string]any, payloadSHA string) notifyDeliveryResult {
	body, _ := json.Marshal(payload)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target.URL, bytes.NewReader(body))
	if err != nil {
		code := "network_error"
		msg := "network error"
		return notifyDeliveryResult{Result: "failure", ErrorCode: &code, Error: &msg, PendingError: msg}
	}
	req.Header.Set("Content-Type", "application/json")
	if secret, ok := target.Config["secret"].(string); ok && secret != "" {
		mac := hmac.New(sha256.New, []byte(secret))
		_, _ = mac.Write(body)
		req.Header.Set("X-Adlaire-Signature-256", "sha256="+hex.EncodeToString(mac.Sum(nil)))
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) || isTimeoutError(err) {
			code := "timeout"
			msg := "timeout"
			return notifyDeliveryResult{Result: "failure", ErrorCode: &code, Error: &msg, Retryable: true, PendingError: msg}
		}
		code := "network_error"
		msg := "network error"
		return notifyDeliveryResult{Result: "failure", ErrorCode: &code, Error: &msg, PendingError: msg}
	}
	defer resp.Body.Close()
	status := resp.StatusCode
	if status >= 200 && status <= 299 {
		return notifyDeliveryResult{Result: "success", HTTPStatus: &status}
	}
	code := httpNotifyErrorCode(status)
	msg := fmt.Sprintf("HTTP %03d", status)
	return notifyDeliveryResult{Result: "failure", ErrorCode: &code, Error: &msg, HTTPStatus: &status, Retryable: status >= 500, PendingError: msg}
}

func httpNotifyErrorCode(status int) string {
	switch {
	case status >= 100 && status <= 199:
		return "http_1xx"
	case status >= 300 && status <= 399:
		return "http_3xx"
	case status >= 400 && status <= 499:
		return "http_4xx"
	case status >= 500 && status <= 599:
		return "http_5xx"
	default:
		return "network_error"
	}
}

func isTimeoutError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func sendEmailNotification(cfg RunnerConfig, target notifyTarget, payload map[string]any) notifyDeliveryResult {
	smtpCfg, err := readSMTPConfig(filepath.Join(cfg.StateDir, ".smtp_config"))
	if err != nil || !smtpCfg.Enabled || smtpCfg.Host == nil || *smtpCfg.Host == "" || smtpCfg.From == nil || *smtpCfg.From == "" {
		code := "smtp_not_configured"
		msg := "SMTP not configured"
		return notifyDeliveryResult{Result: "not_configured", ErrorCode: &code, Error: &msg}
	}
	to := smtpCfg.To
	if rawTo, ok := target.Config["to"].([]any); ok {
		to = []string{}
		for _, item := range rawTo {
			if s, ok := item.(string); ok && s != "" {
				to = append(to, s)
			}
		}
	} else if rawTo, ok := target.Config["to"].([]string); ok {
		to = rawTo
	}
	if len(to) == 0 {
		code := "smtp_not_configured"
		msg := "SMTP not configured"
		return notifyDeliveryResult{Result: "not_configured", ErrorCode: &code, Error: &msg}
	}
	port := smtpCfg.Port
	if port == 0 {
		port = 587
	}
	addr := net.JoinHostPort(*smtpCfg.Host, strconv.Itoa(port))
	body, _ := json.Marshal(payload)
	msg := []byte("From: " + *smtpCfg.From + "\r\nTo: " + strings.Join(to, ",") + "\r\nSubject: Adlaire CI notification\r\nContent-Type: application/json\r\n\r\n" + string(body))
	var auth smtp.Auth
	if smtpCfg.User != nil && *smtpCfg.User != "" {
		secret, err := os.ReadFile(filepath.Join(cfg.StateDir, ".smtp_secret"))
		if err != nil {
			code := "smtp_not_configured"
			text := "SMTP not configured"
			return notifyDeliveryResult{Result: "not_configured", ErrorCode: &code, Error: &text}
		}
		auth = smtp.PlainAuth("", *smtpCfg.User, strings.TrimSpace(string(secret)), *smtpCfg.Host)
	}
	if err := smtp.SendMail(addr, auth, *smtpCfg.From, to, msg); err != nil {
		code := "smtp_error"
		text := "SMTP error"
		return notifyDeliveryResult{Result: "failure", ErrorCode: &code, Error: &text}
	}
	return notifyDeliveryResult{Result: "success"}
}

func readSMTPConfig(path string) (smtpConfigFile, error) {
	var config smtpConfigFile
	err := runnerReadJSONFile(path, &config)
	if errors.Is(err, os.ErrNotExist) {
		config.Port = 587
		return config, nil
	}
	return config, err
}

func sendCommandNotification(target notifyTarget, payload map[string]any, payloadSHA string) notifyDeliveryResult {
	args := notifyCommandArgs(target.Config)
	if len(args) == 0 {
		code := "command_error"
		msg := "command start failed"
		return notifyDeliveryResult{Result: "failure", ErrorCode: &code, Error: &msg}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	body, _ := json.Marshal(payload)
	cmd.Stdin = bytes.NewReader(body)
	cmd.Env = append(os.Environ(),
		"ADLAIRE_NOTIFY_EVENT="+fmt.Sprint(payload["event"]),
		"ADLAIRE_NOTIFY_BUILD_ID="+fmt.Sprint(payload["build_id"]),
		"ADLAIRE_NOTIFY_STATUS="+fmt.Sprint(payload["status"]),
		"ADLAIRE_NOTIFY_BRANCH="+fmt.Sprint(payload["branch"]),
		"ADLAIRE_NOTIFY_TRIGGER="+fmt.Sprint(payload["trigger"]),
		"ADLAIRE_NOTIFY_PAYLOAD_SHA256="+payloadSHA,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		code := "timeout"
		msg := "timeout"
		return notifyDeliveryResult{Result: "failure", ErrorCode: &code, Error: &msg}
	}
	if err != nil {
		code := "command_error"
		msg := "command start failed"
		if exit := exitCodeFromError(err); exit != 1 || strings.Contains(fmt.Sprint(err), "exit status") {
			msg = fmt.Sprintf("command exit %d", exit)
		}
		return notifyDeliveryResult{Result: "failure", ErrorCode: &code, Error: &msg}
	}
	return notifyDeliveryResult{Result: "success"}
}

func notifyCommandArgs(config map[string]any) []string {
	raw, ok := config["command_args"]
	if !ok {
		return nil
	}
	if args, ok := raw.([]string); ok {
		return append([]string(nil), args...)
	}
	items, ok := raw.([]any)
	if !ok {
		return nil
	}
	args := []string{}
	for _, item := range items {
		s, ok := item.(string)
		if !ok || s == "" || strings.ContainsAny(s, "\x00\r\n") {
			return nil
		}
		args = append(args, s)
	}
	return args
}

func nextNotifyRecordID(prefix string, now time.Time) string {
	return prefix + now.UTC().Format("20060102150405")
}

func appendNotifyLog(stateDir string, record map[string]any) error {
	return appendRunnerJSONLine(filepath.Join(stateDir, ".notify_log"), record, 0600)
}

func addNotifyPending(path string, entry notifyPendingEntry) error {
	var entries []notifyPendingEntry
	_ = readJSONArray(path, &entries)
	if entry.PayloadSHA == "" {
		entry.PayloadSHA = notificationPayloadSHA(entry.Payload)
	}
	for _, existing := range entries {
		existingSHA := existing.PayloadSHA
		if existingSHA == "" {
			existingSHA = notificationPayloadSHA(existing.Payload)
		}
		if existing.Event == entry.Event && existing.ChannelID == entry.ChannelID && existingSHA == entry.PayloadSHA {
			return nil
		}
	}
	entries = append(entries, entry)
	return runnerAtomicWriteJSON(path, entries, 0600)
}

func logFailure(cfg RunnerConfig, target BranchTarget, buildID string, started time.Time, blobSHA *string, prevSHA, trigger, status, msg string, pl *pipelineLog, rep *runnerReport, logger *slog.Logger) {
	if pl == nil {
		pl = &pipelineLog{}
	}
	finished := runnerNow().UTC()
	failureCategory := buildFailureCategory(status)
	failureEvidence := buildFailureEvidence(status, failureCategory, finished)
	normalizedStatus := buildStatusFromTargetStatus(status)
	blog := buildLog{
		ID: buildID, Status: normalizedStatus, Branch: target.Branch, TargetFile: target.TargetFile,
		TargetFiles: append([]string(nil), target.TargetFiles...), ChangedTargets: []string{}, MatchedTags: []string{}, TargetStatus: status,
		Trigger: trigger, TriggerActor: nil,
		StartedAt: started.Format(time.RFC3339), FinishedAt: finished.Format(time.RFC3339),
		DurationSeconds: int64(finished.Sub(started).Seconds()), Commit: commitInfo{},
		CommitSHA: nil, CommitMessage: nil, CommitAuthor: nil, CommitAt: nil,
		BlobSHA: blobSHA, PreviousBlobSHA: prevSHA, Pipeline: *pl, SkippedHookIDs: []string{}, Attempts: []pipelineLog{}, RetryCount: 0, Report: rep,
		RemoteBuild: pl.RemoteBuild, Chain: nil, ChainSummary: nil,
		Warnings: []string{}, Deploy: []deployLog{}, TargetResults: []targetResultLog{}, SnapshotID: nil,
		RollbackFrom: nil, TransferVerified: nil,
		OutputSHA256: nil, OutputSizeBytes: nil, SizeWarn: false,
		FailureCategory: failureCategory, FailureEvidence: failureEvidence, CommitStatusLog: buildCommitStatusLog(cfg, nil, finished), BuildMeta: buildBuildMeta(rep), Environment: currentBuildEnv(cfg, "", runnerEnvironmentKeysFrom(target.Env, nil)), Error: &msg,
	}
	if err := writeBuildLog(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_LOG_WRITE_FAILED: " + err.Error())
		_ = writeBuildStatusError(cfg, "failure_state_write", trigger, err.Error())
	}
	if err := appendHistory(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_HISTORY_WRITE_FAILED: " + err.Error())
		_ = writeBuildStatusError(cfg, "failure_state_write", trigger, err.Error())
	}
	if err := updateBuildTrends(cfg, blog); err != nil {
		logger.Warn("BUILD_TRENDS_WRITE_FAILED: " + err.Error())
	}
	if err := writeBuildStatusFinal(cfg, blog); err != nil {
		logger.Error("BUILD_STATUS_WRITE_FAILED: " + err.Error())
	}
	sendBuildNotifications(cfg, blog, logger)
}

func logHookAbortFailure(cfg RunnerConfig, target BranchTarget, buildID string, started time.Time, blobSHA *string, prevSHA, trigger, msg string, skippedHookIDs []string, warnings []string, logger *slog.Logger) {
	finished := runnerNow().UTC()
	status := "hook_error"
	failureCategory := buildFailureCategory(status)
	failureEvidence := buildFailureEvidence(status, failureCategory, finished)
	normalizedStatus := buildStatusFromTargetStatus(status)
	if skippedHookIDs == nil {
		skippedHookIDs = []string{}
	}
	if warnings == nil {
		warnings = []string{}
	}
	blog := buildLog{
		ID: buildID, Status: normalizedStatus, Branch: target.Branch, TargetFile: target.TargetFile,
		TargetFiles: append([]string(nil), target.TargetFiles...), ChangedTargets: []string{}, MatchedTags: []string{}, TargetStatus: status,
		Trigger: trigger, TriggerActor: nil,
		StartedAt: started.Format(time.RFC3339), FinishedAt: finished.Format(time.RFC3339),
		DurationSeconds: int64(finished.Sub(started).Seconds()), Commit: commitInfo{},
		CommitSHA: nil, CommitMessage: nil, CommitAuthor: nil, CommitAt: nil,
		BlobSHA: blobSHA, PreviousBlobSHA: prevSHA, Pipeline: pipelineLog{}, SkippedHookIDs: append([]string(nil), skippedHookIDs...),
		Attempts: []pipelineLog{}, RetryCount: 0, Report: nil,
		RemoteBuild: nil, Chain: nil, ChainSummary: nil,
		Warnings: append([]string(nil), warnings...), Deploy: []deployLog{}, TargetResults: []targetResultLog{}, SnapshotID: nil,
		RollbackFrom: nil, TransferVerified: nil,
		OutputSHA256: nil, OutputSizeBytes: nil, SizeWarn: false,
		FailureCategory: failureCategory, FailureEvidence: failureEvidence, CommitStatusLog: buildCommitStatusLog(cfg, nil, finished), BuildMeta: nil, Environment: currentBuildEnv(cfg, "", runnerEnvironmentKeysFrom(target.Env, nil)), Error: &msg,
	}
	if err := writeBuildLog(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_LOG_WRITE_FAILED: " + err.Error())
		_ = writeBuildStatusError(cfg, "failure_state_write", trigger, err.Error())
	}
	if err := appendHistory(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_HISTORY_WRITE_FAILED: " + err.Error())
		_ = writeBuildStatusError(cfg, "failure_state_write", trigger, err.Error())
	}
	if err := updateBuildTrends(cfg, blog); err != nil {
		logger.Warn("BUILD_TRENDS_WRITE_FAILED: " + err.Error())
	}
	if err := writeBuildStatusFinal(cfg, blog); err != nil {
		logger.Error("BUILD_STATUS_WRITE_FAILED: " + err.Error())
	}
	sendBuildNotifications(cfg, blog, logger)
}

func circuitOpen(stateDir string) bool {
	state, err := runnerReadCircuitState(stateDir)
	return err == nil && state.Open
}

func recordCircuitFailure(cfg RunnerConfig, msg string) {
	if cfg.APICircuitBreakerThreshold <= 0 {
		return
	}
	state, _ := runnerReadCircuitState(cfg.StateDir)
	now := runnerNow().UTC().Format(time.RFC3339)
	state.ConsecutiveFailures++
	state.LastFailureAt = &now
	state.LastError = &msg
	if state.ConsecutiveFailures >= cfg.APICircuitBreakerThreshold {
		state.Open = true
		state.OpenedAt = &now
	}
	_ = runnerAtomicWriteJSON(filepath.Join(cfg.StateDir, ".build_circuit_state"), state, 0600)
}

func resetCircuitState(cfg RunnerConfig) {
	state := buildCircuitState{}
	_ = runnerAtomicWriteJSON(filepath.Join(cfg.StateDir, ".build_circuit_state"), state, 0600)
}

func runnerReadCircuitState(stateDir string) (buildCircuitState, error) {
	var state buildCircuitState
	err := runnerReadJSONFile(filepath.Join(stateDir, ".build_circuit_state"), &state)
	if errors.Is(err, os.ErrNotExist) {
		return buildCircuitState{}, nil
	}
	return state, err
}

func derefString(s *string, fallback string) string {
	if s == nil {
		return fallback
	}
	return *s
}

func readSHACache(path string) (string, error) {
	var cache shaCache
	if err := runnerReadJSONFile(path, &cache); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return cache.SHA, nil
}

func writeBuildState(stateDir string, running bool, buildID *string) error {
	return writeBuildStateStart(stateDir, running, buildID, nil)
}

func writeBuildStateStart(stateDir string, running bool, buildID *string, queue *queueEntry) error {
	now := runnerNow().UTC().Format(time.RFC3339)
	state, err := runnerReadBuildState(stateDir)
	if err != nil {
		state = runnerDefaultBuildState()
	}
	if queue != nil {
		activeID := queueEntryID(state.ActiveQueueEntry)
		if activeID == "" {
			nextQueued := make([]map[string]any, 0, len(state.Queued))
			found := false
			for _, entry := range state.Queued {
				if queueEntryID(entry) == queue.ID {
					found = true
					continue
				}
				nextQueued = append(nextQueued, entry)
			}
			if !found {
				return errors.New("queue entry not found")
			}
			state.Queued = nextQueued
			state.ActiveQueueEntry = copyStringAnyMap(queue.OriginalData)
		} else if activeID != queue.ID {
			return errors.New("active queue entry mismatch")
		}
	}
	state.Running = running
	state.CurrentBuildID = buildID
	state.LastStartedAt = &now
	return runnerAtomicWriteJSON(filepath.Join(stateDir, ".build_state"), state, 0600)
}

func clearActiveQueueEntry(stateDir, id string) error {
	state, err := runnerReadBuildState(stateDir)
	if err != nil {
		return err
	}
	if queueEntryID(state.ActiveQueueEntry) == id {
		state.ActiveQueueEntry = nil
	}
	now := runnerNow().UTC().Format(time.RFC3339)
	state.Running = false
	state.CurrentBuildID = nil
	state.LastFinishedAt = &now
	return runnerAtomicWriteJSON(filepath.Join(stateDir, ".build_state"), state, 0600)
}

func queueEntryID(raw map[string]any) string {
	if raw == nil {
		return ""
	}
	id, _ := raw["id"].(string)
	return id
}

func runnerDefaultBuildState() buildState {
	return buildState{Queued: []map[string]any{}}
}

func clearRunningBuildState(stateDir string) error {
	state, err := runnerReadBuildState(stateDir)
	if err != nil {
		state = runnerDefaultBuildState()
	}
	now := runnerNow().UTC().Format(time.RFC3339)
	state.Running = false
	state.CurrentBuildID = nil
	state.LastFinishedAt = &now
	return runnerAtomicWriteJSON(filepath.Join(stateDir, ".build_state"), state, 0600)
}

func runnerNextBuildID(stateDir string, now time.Time) (string, error) {
	for i := 0; i <= 999; i++ {
		id := runnerBuildID(now, i)
		if _, err := os.Stat(filepath.Join(stateDir, ".build_logs", id+".json")); errors.Is(err, os.ErrNotExist) {
			return id, nil
		}
	}
	return "", errors.New("build id suffix exhausted")
}

func buildStatusFromTargetStatus(targetStatus string) string {
	switch {
	case targetStatus == "success" || targetStatus == "success_deploy_pending":
		return "success"
	case targetStatus == "cancelled":
		return "cancelled"
	case strings.HasPrefix(targetStatus, "skipped_") || targetStatus == "circuit_open" || targetStatus == "lock_skipped":
		return "skipped"
	default:
		return "failure"
	}
}

func buildFailureCategory(targetStatus string) *string {
	var category string
	switch targetStatus {
	case "success_deploy_pending":
		category = "deploy_failure"
	case "failure_api":
		category = "github_api"
	case "failure_decode":
		category = "unknown"
	case "failure_precheck":
		category = "resource_error"
	case "failure_build":
		category = "pipeline_exit"
	case "failure_timeout":
		category = "pipeline_timeout"
	case "failure_state_write":
		category = "resource_error"
	case "failure_pipeline_config", "failure_tag_rule", "failure_target_missing":
		category = "config_error"
	case "failure_remote_build":
		category = "pipeline_exit"
	case "hook_error":
		category = "hook_error"
	default:
		return nil
	}
	return &category
}

func buildFailureEvidence(targetStatus string, category *string, at time.Time) []failureEvidence {
	if category == nil {
		return []failureEvidence{}
	}
	source := "runner"
	message := "runner failure"
	switch *category {
	case "github_api":
		source = "github_api"
		message = "github api failure"
	case "pipeline_timeout":
		source = "pipeline"
		message = "pipeline timeout"
	case "pipeline_exit":
		source = "pipeline"
		message = "pipeline exited non-zero"
	case "deploy_failure":
		source = "deploy"
		message = "deploy verification failed"
	case "hook_error":
		source = "hook"
		message = "hook failure"
	case "config_error":
		source = "config"
		message = "configuration invalid"
	case "resource_error":
		source = "resource"
		message = "resource failure"
	case "unknown":
		source = "runner"
		message = "unknown failure"
	}
	return []failureEvidence{{
		Source:  source,
		Code:    *category,
		Message: message,
		At:      at.UTC().Format(time.RFC3339),
	}}
}

func writeBuildStatusRunning(cfg RunnerConfig, buildID, trigger, branch, targetFile, started string) error {
	triggerValue := trigger
	status := buildStatusSummary{
		SchemaVersion:  1,
		Status:         "running",
		Running:        true,
		CurrentBuildID: &buildID,
		LastTrigger:    &triggerValue,
		LastBranch:     &branch,
		LastTargetFile: &targetFile,
		LastStartedAt:  &started,
	}
	now := runnerNow().UTC().Format(time.RFC3339)
	status.UpdatedAt = &now
	fillBuildStatusCounts(cfg.StateDir, &status)
	return runnerAtomicWriteJSON(filepath.Join(cfg.StateDir, ".build_status.json"), status, 0600)
}

func writeBuildStatusFinal(cfg RunnerConfig, log buildLog) error {
	statusText := buildStatusFromTargetStatus(log.TargetStatus)
	trigger := log.Trigger
	branch := log.Branch
	targetFile := log.TargetFile
	status := buildStatusSummary{
		SchemaVersion:       1,
		Status:              statusText,
		Running:             false,
		CurrentBuildID:      nil,
		LastBuildID:         &log.ID,
		LastTrigger:         &trigger,
		LastTargetStatus:    &log.TargetStatus,
		LastBranch:          &branch,
		LastTargetFile:      &targetFile,
		LastBlobSHA:         log.BlobSHA,
		LastCommitSHA:       log.Commit.SHA,
		LastStartedAt:       &log.StartedAt,
		LastFinishedAt:      &log.FinishedAt,
		LastDurationSeconds: &log.DurationSeconds,
		LastError:           log.Error,
		OutputSHA256:        log.OutputSHA256,
		SizeWarn:            log.SizeWarn,
	}
	now := runnerNow().UTC().Format(time.RFC3339)
	status.UpdatedAt = &now
	if len(log.Deploy) == 0 {
		skipped := "skipped"
		status.LastDeployStatus = &skipped
	} else {
		deployAt := log.FinishedAt
		status.LastDeployAt = &deployAt
		if log.TargetStatus == "success_deploy_pending" {
			pending := "pending"
			status.LastDeployStatus = &pending
		} else {
			success := "success"
			status.LastDeployStatus = &success
		}
	}
	fillBuildStatusCounts(cfg.StateDir, &status)
	return runnerAtomicWriteJSON(filepath.Join(cfg.StateDir, ".build_status.json"), status, 0600)
}

func writeBuildStatusSkip(cfg RunnerConfig, targetStatus, trigger string, target *BranchTarget) error {
	status := buildStatusSummary{
		SchemaVersion:    1,
		Status:           "skipped",
		Running:          false,
		LastTargetStatus: &targetStatus,
		LastTrigger:      &trigger,
		CurrentBuildID:   nil,
		LastDeployStatus: stringPtr("skipped"),
		LastFinishedAt:   stringPtr(runnerNow().UTC().Format(time.RFC3339)),
	}
	now := runnerNow().UTC().Format(time.RFC3339)
	status.UpdatedAt = &now
	if target != nil {
		status.LastBranch = &target.Branch
		status.LastTargetFile = &target.TargetFile
	}
	fillBuildStatusCounts(cfg.StateDir, &status)
	return runnerAtomicWriteJSON(filepath.Join(cfg.StateDir, ".build_status.json"), status, 0600)
}

func writeBuildStatusError(cfg RunnerConfig, targetStatus, trigger, message string) error {
	status := buildStatusSummary{
		SchemaVersion:    1,
		Status:           "failure",
		Running:          false,
		LastTargetStatus: &targetStatus,
		LastTrigger:      &trigger,
		LastError:        stringPtr(sanitizeStatusError(message)),
		CurrentBuildID:   nil,
		LastFinishedAt:   stringPtr(runnerNow().UTC().Format(time.RFC3339)),
	}
	now := runnerNow().UTC().Format(time.RFC3339)
	status.UpdatedAt = &now
	fillBuildStatusCounts(cfg.StateDir, &status)
	return runnerAtomicWriteJSON(filepath.Join(cfg.StateDir, ".build_status.json"), status, 0600)
}

func fillBuildStatusCounts(stateDir string, status *buildStatusSummary) {
	status.PendingTransfersCount = runnerJSONArrayLen(filepath.Join(stateDir, ".pending_transfers"))
	status.NotifyPendingCount = runnerJSONArrayLen(filepath.Join(stateDir, ".notify_pending"))
	if circuit, err := runnerReadCircuitState(stateDir); err == nil {
		status.CircuitOpen = circuit.Open
		status.CircuitConsecutiveFailures = circuit.ConsecutiveFailures
	}
}

func runnerJSONArrayLen(path string) int {
	var raw []json.RawMessage
	if err := runnerReadJSONFile(path, &raw); err != nil {
		return 0
	}
	return len(raw)
}

func sanitizeStatusError(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\n")
	if len(s) > 500 {
		s = s[:500]
	}
	return s
}

func stringPtr(s string) *string {
	return &s
}

func currentBuildEnv(cfg RunnerConfig, builderVersion string, envKeys []string) buildEnvLog {
	hostname := "unknown"
	if name, err := os.Hostname(); err == nil && name != "" {
		hostname = name
		if len(hostname) > 255 {
			hostname = hostname[:255]
		}
	}
	return buildEnvLog{
		OS:             runtime.GOOS,
		Arch:           runtime.GOARCH,
		GoVersion:      runtime.Version(),
		RunnerVersion:  runnerBinaryVersion,
		BuilderVersion: builderVersion,
		Hostname:       hostname,
		PID:            os.Getpid(),
		WatchMode:      cfg.WatchMode,
		CacheEnabled:   cfg.BuildCacheEnabled,
		RemoteBuild:    cfg.RemoteBuild.Enabled,
		StateDir:       runnerStateDirForLog(cfg.StateDir),
		DiskFreeBytes:  runnerDiskFreeBytes(cfg.StateDir),
		CapturedAt:     runnerNow().UTC().Format(time.RFC3339),
		EnvKeys:        append([]string(nil), envKeys...),
	}
}

func runnerStateDirForLog(stateDir string) string {
	abs, err := filepath.Abs(stateDir)
	if err != nil {
		return stateDir
	}
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		if rel, err := filepath.Rel(home, abs); err == nil && rel != "." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && rel != ".." {
			return filepath.Base(abs)
		}
	}
	return abs
}

func runnerDiskFreeBytes(path string) *int64 {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil
	}
	free := stat.Bavail * uint64(stat.Bsize)
	maxInt64 := uint64(^uint64(0) >> 1)
	if free > maxInt64 {
		free = maxInt64
	}
	v := int64(free)
	return &v
}

func runnerEnvironmentKeysFrom(branchEnv map[string]string, pipelineEnv map[string]string) []string {
	keys := map[string]bool{}
	for key := range branchEnv {
		keys[key] = true
	}
	for key := range pipelineEnv {
		keys[key] = true
	}
	out := make([]string, 0, len(keys))
	for key := range keys {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func writeBuildLog(stateDir string, log buildLog) error {
	return runnerAtomicWriteJSON(filepath.Join(stateDir, ".build_logs", log.ID+".json"), log, 0600)
}

func appendHistory(stateDir string, log buildLog) error {
	path := filepath.Join(stateDir, ".build_history")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	release, err := acquireStateFileLock(path)
	if err != nil {
		return err
	}
	defer release()
	var pages *int
	warnings := 0
	sizeWarn := false
	if log.Report != nil {
		pages = &log.Report.Pages
		warnings = log.Report.WarningsCount
		sizeWarn = log.Report.SizeWarn
	}
	chainRunID := (*string)(nil)
	if chain, ok := log.Chain.(*chainLog); ok && chain != nil && chain.ChainRunID != "" {
		chainRunID = &chain.ChainRunID
	} else if chain, ok := log.Chain.(chainLog); ok && chain.ChainRunID != "" {
		chainRunID = &chain.ChainRunID
	}
	rec := historyRecord{
		ID: log.ID, Branch: log.Branch, TargetFile: log.TargetFile, Status: log.TargetStatus,
		Trigger:   log.Trigger,
		StartedAt: log.StartedAt, FinishedAt: log.FinishedAt, DurationSeconds: log.DurationSeconds,
		CommitSHA: log.Commit.SHA, BlobSHA: log.BlobSHA, Pages: pages, Warnings: warnings,
		SizeWarn: sizeWarn, OutputSHA256: log.OutputSHA256, OutputSizeBytes: log.OutputSizeBytes,
		RetryCount: log.RetryCount, FailureCategory: log.FailureCategory, CommitStatus: log.CommitStatus,
		SnapshotID: log.SnapshotID, RollbackFrom: nil, ChainRunID: chainRunID,
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	if info, err := f.Stat(); err == nil && info.Size() > 0 {
		var last [1]byte
		if _, err := f.ReadAt(last[:], info.Size()-1); err != nil {
			_ = f.Close()
			return err
		}
		if last[0] != '\n' {
			if _, err := f.Write([]byte("\n")); err != nil {
				_ = f.Close()
				return err
			}
		}
	}
	if _, err := f.Write(append(data, '\n')); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return syncParentDir(path)
}

func appendApprovalPending(cfg RunnerConfig, target BranchTarget, buildID, digest string, changedTargets []string) error {
	record := map[string]any{
		"id":              "approval-" + buildID,
		"build_id":        buildID,
		"status":          "pending",
		"branch":          target.Branch,
		"target_file":     target.TargetFile,
		"target_files":    normalizedTargetFiles(target),
		"target_sha":      digest,
		"changed_targets": changedTargets,
		"trigger":         "polling",
		"created_at":      runnerNow().UTC().Format(time.RFC3339),
	}
	return appendRunnerJSONLine(filepath.Join(cfg.StateDir, ".approval_queue"), record, 0600)
}

func appendRunnerJSONLine(path string, v any, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	release, err := acquireStateFileLock(path)
	if err != nil {
		return err
	}
	defer release()
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, mode)
	if err != nil {
		return err
	}
	if info, err := f.Stat(); err == nil && info.Size() > 0 {
		var last [1]byte
		if _, err := f.ReadAt(last[:], info.Size()-1); err != nil {
			_ = f.Close()
			return err
		}
		if last[0] != '\n' {
			if _, err := f.Write([]byte("\n")); err != nil {
				_ = f.Close()
				return err
			}
		}
	}
	if _, err := f.Write(append(data, '\n')); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Chmod(mode); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return syncParentDir(path)
}

func updateBuildTrends(cfg RunnerConfig, log buildLog) error {
	if cfg.BuildTrendKeepCount <= 0 {
		return nil
	}
	path := filepath.Join(cfg.StateDir, ".build_trends.json")
	trends := buildTrendsFile{SchemaVersion: 1, Samples: []buildTrendSample{}, Summary: summarizeBuildTrends(nil)}
	if err := runnerReadJSONFile(path, &trends); err != nil && !errors.Is(err, os.ErrNotExist) {
		_ = backupCorruptJSON(path, slog.New(slog.NewTextHandler(io.Discard, nil)))
		trends = buildTrendsFile{SchemaVersion: 1, Samples: []buildTrendSample{}}
	}
	trends.SchemaVersion = 1
	anomaly := durationAnomaly(cfg, trends.Samples, log.DurationSeconds)
	trends.Samples = append(trends.Samples, buildTrendSample{
		BuildID:         log.ID,
		FinishedAt:      log.FinishedAt,
		Branch:          log.Branch,
		Trigger:         log.Trigger,
		DurationSeconds: log.DurationSeconds,
		Status:          log.Status,
		TargetStatus:    log.TargetStatus,
		Anomaly:         anomaly,
	})
	if len(trends.Samples) > cfg.BuildTrendKeepCount {
		trends.Samples = trends.Samples[len(trends.Samples)-cfg.BuildTrendKeepCount:]
	}
	trends.Summary = summarizeBuildTrends(trends.Samples)
	return runnerAtomicWriteJSON(path, trends, 0600)
}

func durationAnomaly(cfg RunnerConfig, samples []buildTrendSample, duration int64) bool {
	if !cfg.DurationAnomaly.Enabled || len(samples) < cfg.DurationAnomaly.MinSamples || duration < 0 {
		return false
	}
	values := make([]int64, 0, len(samples))
	var total int64
	for _, sample := range samples {
		if sample.DurationSeconds < 0 {
			continue
		}
		values = append(values, sample.DurationSeconds)
		total += sample.DurationSeconds
	}
	if len(values) < cfg.DurationAnomaly.MinSamples {
		return false
	}
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	avg := float64(total) / float64(len(values))
	p95 := percentileDuration(values, 0.95)
	threshold := math.Max(avg*cfg.DurationAnomaly.AvgMultiplier, p95*cfg.DurationAnomaly.P95Multiplier)
	return float64(duration) > threshold
}

func summarizeBuildTrends(samples []buildTrendSample) buildTrendsSummary {
	summary := buildTrendsSummary{Count: len(samples)}
	if len(samples) == 0 {
		return summary
	}
	values := make([]int64, 0, len(samples))
	var total int64
	for _, sample := range samples {
		if sample.Anomaly {
			summary.AnomalyCount++
		}
		duration := sample.DurationSeconds
		if duration < 0 {
			duration = 0
		}
		total += duration
		values = append(values, duration)
	}
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	avg := float64(total) / float64(len(values))
	median := percentileDuration(values, 0.5)
	p95 := percentileDuration(values, 0.95)
	summary.AvgSeconds = &avg
	summary.MedianSeconds = &median
	summary.P95Seconds = &p95
	return summary
}

func percentileDuration(values []int64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	if len(values) == 1 {
		return float64(values[0])
	}
	idx := int(float64(len(values)-1) * p)
	if idx < 0 {
		idx = 0
	}
	if idx >= len(values) {
		idx = len(values) - 1
	}
	return float64(values[idx])
}

func outputManifestSHA(root string) (*string, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("invalid output root: %s", root)
	}
	parts := []string{}
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Sys() == nil {
			return fmt.Errorf("invalid output file: %s", path)
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok && stat.Nlink > 1 {
			return fmt.Errorf("hardlink output file: %s", path)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if invalidOutputRelativePath(rel) {
			return fmt.Errorf("invalid output path: %s", rel)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		after, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if after.Size() != info.Size() || !after.ModTime().Equal(info.ModTime()) {
			return fmt.Errorf("output changed during read: %s", path)
		}
		sum := sha256.Sum256(data)
		parts = append(parts, rel+"\n"+hex.EncodeToString(sum[:])+"\n")
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(parts)
	sum := sha256.Sum256([]byte(strings.Join(parts, "")))
	out := hex.EncodeToString(sum[:])
	return &out, nil
}

func invalidOutputRelativePath(rel string) bool {
	if rel == "" || strings.HasPrefix(rel, "/") || strings.Contains(rel, "\\") || strings.Contains(rel, "\x00") || strings.Contains(rel, "\n") || strings.Contains(rel, "\r") {
		return true
	}
	for _, part := range strings.Split(rel, "/") {
		if part == "" || part == "." || part == ".." {
			return true
		}
	}
	return false
}

func outputSizeBytes(root string) (int64, bool) {
	info, err := os.Lstat(root)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return 0, false
	}
	var size int64
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("invalid output file: %s", path)
		}
		size += info.Size()
		return nil
	})
	return size, err == nil
}

func outputManifestFileCount(root string) (int, bool) {
	info, err := os.Lstat(root)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return 0, false
	}
	count := 0
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("invalid output file: %s", path)
		}
		count++
		return nil
	})
	return count, err == nil
}

func cleanupBuildLogs(dir string, keep int, archiveAfterDays int, logger *slog.Logger) {
	if keep <= 0 {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	type fileInfo struct {
		path string
		mod  time.Time
	}
	files := []fileInfo{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		info, err := e.Info()
		if err == nil {
			files = append(files, fileInfo{filepath.Join(dir, e.Name()), info.ModTime()})
		}
	}
	sort.Slice(files, func(i, j int) bool { return files[i].mod.Before(files[j].mod) })
	cutoff := time.Time{}
	if archiveAfterDays > 0 {
		cutoff = runnerNow().UTC().Add(-time.Duration(archiveAfterDays) * 24 * time.Hour)
	}
	for len(files) > keep {
		if !cutoff.IsZero() && files[0].mod.Before(cutoff) {
			if err := archiveBuildLog(dir, files[0].path); err != nil {
				logger.Warn("LOG_ARCHIVE_FAILED: path=" + files[0].path)
				files = files[1:]
				continue
			}
		}
		if err := os.Remove(files[0].path); err != nil {
			logger.Warn("LOG_CLEANUP_FAILED: path=" + files[0].path)
		}
		files = files[1:]
	}
}

func archiveBuildLog(dir, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	archiveDir := filepath.Join(dir, "archive")
	if err := os.MkdirAll(archiveDir, 0700); err != nil {
		return err
	}
	outPath := filepath.Join(archiveDir, filepath.Base(path)+".gz")
	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	gz, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		_ = f.Close()
		_ = os.Remove(outPath)
		return err
	}
	if _, err := gz.Write(data); err != nil {
		_ = gz.Close()
		_ = f.Close()
		_ = os.Remove(outPath)
		return err
	}
	if err := gz.Close(); err != nil {
		_ = f.Close()
		_ = os.Remove(outPath)
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(outPath)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(outPath)
		return err
	}
	return syncParentDir(outPath)
}

func runnerBuildID(t time.Time, n int) string {
	base := "b" + t.UTC().Format("20060102150405")
	if n <= 0 {
		return base
	}
	return fmt.Sprintf("%s-%03d", base, n)
}

func repairCorruptJSONArray(path string, logger *slog.Logger) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	var raw []json.RawMessage
	if err := runnerReadJSONFile(path, &raw); err == nil {
		return nil
	}
	if err := backupCorruptJSON(path, logger); err != nil {
		return err
	}
	return runnerAtomicWriteJSON(path, []any{}, 0600)
}

func backupCorruptJSON(path string, logger *slog.Logger) error {
	bak := fmt.Sprintf("%s.corrupt.%s.bak", path, runnerNow().UTC().Format("20060102150405"))
	if err := os.Rename(path, bak); err != nil {
		return err
	}
	logger.Warn("STATE_CORRUPT_BACKUP: path=" + bak)
	return nil
}

func readJSONArray(path string, out any) error {
	return runnerReadJSONFile(path, out)
}

func runnerReadJSONFile(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func runnerAtomicWriteJSON(path string, v any, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	release, err := acquireStateFileLock(path)
	if err != nil {
		return err
	}
	defer release()
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	base := filepath.Base(path)
	tmp := filepath.Join(filepath.Dir(path), fmt.Sprintf(".%s.tmp.%d", base, os.Getpid()))
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Chmod(mode); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return syncParentDir(path)
}

func acquireStateFileLock(path string) (func(), error) {
	lockPath := path + ".lock"
	var lastErr error
	for attempt := 0; attempt <= 100; attempt++ {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			_, writeErr := fmt.Fprintf(f, "pid=%d\nstarted_at=%s\n", os.Getpid(), runnerNow().UTC().Format(time.RFC3339))
			closeErr := f.Close()
			if writeErr != nil || closeErr != nil {
				_ = os.Remove(lockPath)
				if writeErr != nil {
					return func() {}, writeErr
				}
				return func() {}, closeErr
			}
			return func() { _ = os.Remove(lockPath) }, nil
		}
		lastErr = err
		if !errors.Is(err, os.ErrExist) || attempt == 100 {
			break
		}
		runnerSleep(100 * time.Millisecond)
	}
	if lastErr == nil {
		lastErr = errors.New("state file lock failed")
	}
	return func() {}, lastErr
}

func syncParentDir(path string) error {
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}
