package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var runnerGitHubAPIBase = "https://api.github.com"
var runnerNow = time.Now
var runnerSleep = time.Sleep

type RunnerConfig struct {
	StateDir                   string
	PendingFile                string
	APIRetryMax                int
	APIRetryBaseSeconds        int
	BuildCooldownSeconds       int
	HistoryKeepN               int
	ForceBuildIntervalHours    int
	LogKeepN                   int
	APICircuitBreakerThreshold int
	OutputSizeWarnMB           int
	WeeklySummaryEnabled       bool
	WeeklySummaryDay           int
	WeeklySummaryHour          int
	BranchTargets              []BranchTarget
}

type BranchTarget struct {
	Branch        string         `json:"branch"`
	TargetFile    string         `json:"target_file"`
	SHAFile       string         `json:"sha_file"`
	Src           string         `json:"src"`
	Out           string         `json:"out"`
	DeployTargets []DeployTarget `json:"deploy_targets"`
}

type DeployTarget struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	DestDir string `json:"dest_dir"`
}

type branchConfigFile struct {
	BranchTargets []BranchTarget `json:"branch_targets"`
}

type shaCache struct {
	SHA string `json:"sha"`
}

type buildState struct {
	Running                 bool             `json:"running"`
	CurrentBuildID          *string          `json:"current_build_id"`
	Queued                  []map[string]any `json:"queued"`
	LastStartedAt           *string          `json:"last_started_at"`
	LastFinishedAt          *string          `json:"last_finished_at"`
	WeeklySummaryLastSentAt *string          `json:"weekly_summary_last_sent_at"`
	WeeklySummarySentDate   *string          `json:"weekly_summary_sent_date"`
}

type commitInfo struct {
	SHA     *string `json:"sha"`
	Message *string `json:"message"`
	Author  *string `json:"author"`
	Date    *string `json:"date"`
}

type pipelineLog struct {
	ExitCode        *int   `json:"exit_code"`
	Stdout          string `json:"stdout"`
	Stderr          string `json:"stderr"`
	StdoutTruncated bool   `json:"stdout_truncated"`
	StderrTruncated bool   `json:"stderr_truncated"`
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
}

type deployLog struct {
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
}

type buildLog struct {
	ID              string        `json:"id"`
	Branch          string        `json:"branch"`
	TargetFile      string        `json:"target_file"`
	TargetStatus    string        `json:"target_status"`
	StartedAt       string        `json:"started_at"`
	FinishedAt      string        `json:"finished_at"`
	DurationSeconds int64         `json:"duration_seconds"`
	Commit          commitInfo    `json:"commit"`
	BlobSHA         *string       `json:"blob_sha"`
	PreviousBlobSHA string        `json:"previous_blob_sha"`
	Pipeline        pipelineLog   `json:"pipeline"`
	Report          *runnerReport `json:"report"`
	Warnings        []string      `json:"warnings"`
	Deploy          []deployLog   `json:"deploy"`
	SnapshotID      *string       `json:"snapshot_id"`
	Error           *string       `json:"error"`
}

type historyRecord struct {
	ID              string  `json:"id"`
	Branch          string  `json:"branch"`
	TargetFile      string  `json:"target_file"`
	Status          string  `json:"status"`
	StartedAt       string  `json:"started_at"`
	FinishedAt      string  `json:"finished_at"`
	DurationSeconds int64   `json:"duration_seconds"`
	CommitSHA       *string `json:"commit_sha"`
	BlobSHA         *string `json:"blob_sha"`
	Pages           *int    `json:"pages"`
	Warnings        int     `json:"warnings"`
	SizeWarn        bool    `json:"size_warn"`
	OutputSHA256    *string `json:"output_sha256"`
	RollbackFrom    *string `json:"rollback_from"`
}

type pendingTransfer struct {
	BranchIdx  int    `json:"branch_idx"`
	DeployIdx  int    `json:"deploy_idx"`
	Out        string `json:"out"`
	Host       string `json:"host"`
	User       string `json:"user"`
	DestDir    string `json:"dest_dir"`
	FailedAt   string `json:"failed_at"`
	RetryCount int    `json:"retry_count"`
}

type notifyPendingEntry struct {
	Event      string         `json:"event"`
	URL        string         `json:"url"`
	Payload    map[string]any `json:"payload"`
	QueuedAt   string         `json:"queued_at"`
	RetryCount int            `json:"retry_count"`
	LastError  string         `json:"last_error"`
}

type notifyConfigFile struct {
	Webhooks []notifyWebhook `json:"webhooks"`
}

type notifyWebhook struct {
	URL     string   `json:"url"`
	Label   string   `json:"label"`
	Enabled bool     `json:"enabled"`
	On      []string `json:"on"`
}

type buildCircuitState struct {
	Open                bool    `json:"open"`
	ConsecutiveFailures int     `json:"consecutive_failures"`
	OpenedAt            *string `json:"opened_at"`
	LastFailureAt       *string `json:"last_failure_at"`
	LastError           *string `json:"last_error"`
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

func runRunner(args []string, stdout, stderr io.Writer) int {
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
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--help":
			fmt.Fprintln(stdout, "Usage: adlaire-ci-runner [--state-dir path] [--once] [--version] [--help]")
			return cfg, true, nil
		case "--version":
			fmt.Fprintf(stdout, "adlaire-ci-runner ADLAIRE_CI_SPEC go=%s\n", runtime.Version())
			return cfg, true, nil
		case "--once":
			continue
		}
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
		if !strings.HasPrefix(arg, "--") {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
		name, value, hasValue := strings.Cut(arg, "=")
		if name != "--state-dir" {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + name}
		}
		if !hasValue {
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
				return cfg, false, exitError{Code: 2, Msg: "missing value: " + name}
			}
			i++
			value = args[i]
		}
		cfg = defaultRunnerConfig(value)
	}
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

func defaultRunnerConfig(stateDir string) RunnerConfig {
	return RunnerConfig{
		StateDir:                   stateDir,
		PendingFile:                filepath.Join(stateDir, ".pending_transfers"),
		APIRetryMax:                5,
		APIRetryBaseSeconds:        1,
		BuildCooldownSeconds:       60,
		HistoryKeepN:               10,
		LogKeepN:                   50,
		APICircuitBreakerThreshold: 3,
		OutputSizeWarnMB:           5,
		WeeklySummaryEnabled:       true,
		WeeklySummaryDay:           0,
		WeeklySummaryHour:          9,
		BranchTargets: []BranchTarget{{
			Branch:     "main",
			TargetFile: "docs",
			SHAFile:    filepath.Join(stateDir, ".last_sha"),
			Src:        filepath.Join(stateDir, "repo", "docs"),
			Out:        filepath.Join(stateDir, "dist", "site"),
		}},
	}
}

func executeRunner(cfg RunnerConfig, stdout, stderr io.Writer) int {
	logger := slog.New(slog.NewTextHandler(stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if err := loadRunnerConfig(&cfg, logger); err != nil {
		logger.Error(err.Error())
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
	token, err := readRunnerToken(filepath.Join(cfg.StateDir, ".github_token"))
	if err != nil {
		logger.Error(err.Error())
		return 2
	}
	lockPath := filepath.Join(cfg.StateDir, ".build_lock")
	release, locked, code := acquireRunnerLock(lockPath, logger)
	if !locked {
		return code
	}
	defer release()
	if circuitOpen(cfg.StateDir) {
		logger.Error("CIRCUIT_OPEN: polling skipped")
		return 1
	}
	if cooldownActive(cfg, logger) {
		return 0
	}
	buildID := runnerBuildID(runnerNow().UTC(), 1)
	if err := writeBuildState(cfg.StateDir, true, &buildID); err != nil {
		logger.Error("STATE_START_FAILED: " + err.Error())
		return 1
	}
	exit := 0
	for i, target := range cfg.BranchTargets {
		if status := processRunnerTarget(cfg, i, target, token, buildID, logger); status > exit {
			exit = status
		}
	}
	finished := runnerNow().UTC().Format(time.RFC3339)
	state := defaultBuildState()
	state.LastFinishedAt = &finished
	if err := runnerAtomicWriteJSON(filepath.Join(cfg.StateDir, ".build_state"), state, 0600); err != nil {
		logger.Error("STATE_FINISH_FAILED: " + err.Error())
		if exit < 1 {
			exit = 1
		}
	}
	cleanupBuildLogs(filepath.Join(cfg.StateDir, ".build_logs"), cfg.LogKeepN, logger)
	_ = stderr
	return exit
}

func loadRunnerConfig(cfg *RunnerConfig, logger *slog.Logger) error {
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
			logger.Warn("CONFIG_UNKNOWN_KEY: key=" + key)
		}
	}
	var bc branchConfigFile
	if err := json.Unmarshal(data, &bc); err != nil {
		return err
	}
	cfg.BranchTargets = bc.BranchTargets
	return validateRunnerConfig(*cfg)
}

func validateRunnerConfig(cfg RunnerConfig) error {
	if len(cfg.BranchTargets) == 0 {
		return errors.New("branch targets must not be empty")
	}
	for _, t := range cfg.BranchTargets {
		if t.Branch == "" || strings.Contains(t.Branch, "..") || strings.Contains(t.Branch, "~") {
			return fmt.Errorf("invalid branch target: %s", t.Branch)
		}
		if t.TargetFile == "" || filepath.IsAbs(t.TargetFile) || strings.Contains(t.TargetFile, "..") {
			return fmt.Errorf("invalid target file: %s", t.TargetFile)
		}
		for _, p := range []string{t.SHAFile, t.Src, t.Out} {
			if p == "" || !filepath.IsAbs(p) {
				return fmt.Errorf("path must be absolute: %s", p)
			}
		}
		for _, d := range t.DeployTargets {
			if d.Host == "" || d.User == "" || d.DestDir == "" || !filepath.IsAbs(d.DestDir) {
				return errors.New("invalid deploy target")
			}
		}
	}
	return nil
}

func cooldownActive(cfg RunnerConfig, logger *slog.Logger) bool {
	if cfg.BuildCooldownSeconds <= 0 {
		return false
	}
	state, err := readBuildState(cfg.StateDir)
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
	state, err := readBuildState(cfg.StateDir)
	if err != nil || state.LastFinishedAt == nil || *state.LastFinishedAt == "" {
		return true
	}
	last, err := time.Parse(time.RFC3339, *state.LastFinishedAt)
	if err != nil {
		return true
	}
	return runnerNow().Sub(last) >= time.Duration(cfg.ForceBuildIntervalHours)*time.Hour
}

func readBuildState(stateDir string) (buildState, error) {
	var state buildState
	err := readJSONFile(filepath.Join(stateDir, ".build_state"), &state)
	if errors.Is(err, os.ErrNotExist) {
		return defaultBuildState(), nil
	}
	return state, err
}

func readRunnerToken(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", errors.New("github token missing")
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", errors.New("github token missing")
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
		return func() { _ = os.Remove(path) }, true, 0
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
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

func processRunnerTarget(cfg RunnerConfig, idx int, target BranchTarget, token, buildID string, logger *slog.Logger) int {
	started := runnerNow().UTC()
	prevSHA, _ := readSHACache(target.SHAFile)
	blobSHA, err := fetchTargetSHA(cfg, logger, token, target.Branch, target.TargetFile)
	if err != nil {
		logFailure(cfg, target, buildID, started, nil, prevSHA, "failure_api", "github api failed", nil, nil, logger)
		recordCircuitFailure(cfg, "github api failed")
		return 3
	}
	forceBuild := forceIntervalDue(cfg)
	if blobSHA == prevSHA && !forceBuild {
		logger.Info(fmt.Sprintf("NO_CHANGE: branch=%s target=%s sha=%s", target.Branch, target.TargetFile, blobSHA))
		return 0
	}
	content, err := fetchBlobContent(cfg, logger, token, blobSHA)
	if err != nil {
		logFailure(cfg, target, buildID, started, &blobSHA, prevSHA, "failure_decode", "blob decode failed", nil, nil, logger)
		recordCircuitFailure(cfg, "blob decode failed")
		return 1
	}
	if err := materializeSource(target.Src, content); err != nil {
		logFailure(cfg, target, buildID, started, &blobSHA, prevSHA, "failure_decode", "source write failed", nil, nil, logger)
		recordCircuitFailure(cfg, "source write failed")
		return 1
	}
	if err := precheckRunnerTarget(target); err != nil {
		logFailure(cfg, target, buildID, started, &blobSHA, prevSHA, "failure_precheck", "precheck failed", nil, nil, logger)
		recordCircuitFailure(cfg, "precheck failed")
		return 1
	}
	commit := fetchCommitInfo(cfg, logger, token, target.Branch, target.TargetFile)
	pl := runPipeline(cfg, target, buildID)
	rep, warns := parseRunnerReport(pl.Stdout)
	if rep != nil && cfg.OutputSizeWarnMB > 0 {
		if size, ok := outputSizeBytes(target.Out); ok && size > int64(cfg.OutputSizeWarnMB)*1024*1024 {
			rep.SizeWarn = true
			warns = append(warns, "OUTPUT_SIZE_WARN")
		}
	}
	deploys := []deployLog{}
	status := "success"
	var errText *string
	exitCode := 0
	if pl.ExitCode == nil || *pl.ExitCode != 0 {
		status = "failure_build"
		msg := "pipeline failed"
		errText = &msg
		exitCode = 1
	} else if err := runnerAtomicWriteJSON(target.SHAFile, shaCache{SHA: blobSHA}, 0600); err != nil {
		status = "failure_state_write"
		msg := "state write failed"
		errText = &msg
		exitCode = 1
	} else {
		for deployIdx, d := range target.DeployTargets {
			dl := executeDeploy(cfg, idx, deployIdx, target, d)
			deploys = append(deploys, dl)
			if dl.Status == "pending" {
				status = "success_deploy_pending"
				msg := "deploy pending"
				errText = &msg
				exitCode = 1
			}
		}
	}
	var snapshotID *string
	if status == "success" {
		if id, err := writeSnapshot(cfg, target, buildID); err == nil {
			snapshotID = &id
		} else {
			logger.Warn("SNAPSHOT_FAILED: " + err.Error())
		}
	}
	finished := runnerNow().UTC()
	blog := buildLog{
		ID:              buildID,
		Branch:          target.Branch,
		TargetFile:      target.TargetFile,
		TargetStatus:    status,
		StartedAt:       started.Format(time.RFC3339),
		FinishedAt:      finished.Format(time.RFC3339),
		DurationSeconds: int64(finished.Sub(started).Seconds()),
		Commit:          commit,
		BlobSHA:         &blobSHA,
		PreviousBlobSHA: prevSHA,
		Pipeline:        pl,
		Report:          rep,
		Warnings:        warns,
		Deploy:          deploys,
		SnapshotID:      snapshotID,
		Error:           errText,
	}
	if err := writeBuildLog(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_LOG_WRITE_FAILED: " + err.Error())
		return 1
	}
	if err := appendHistory(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_HISTORY_WRITE_FAILED: " + err.Error())
		return 1
	}
	sendBuildNotifications(cfg, blog, logger)
	if exitCode == 0 {
		resetCircuitState(cfg)
	} else if strings.HasPrefix(status, "failure_") {
		recordCircuitFailure(cfg, derefString(errText, "runner target failed"))
	}
	return exitCode
}

func fetchTargetSHA(cfg RunnerConfig, logger *slog.Logger, token, branch, targetFile string) (string, error) {
	owner, repo := runnerRepo()
	url := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=1", strings.TrimRight(runnerGitHubAPIBase, "/"), owner, repo, branch)
	var tree gitTreeResponse
	if err := runnerGetJSON(cfg, logger, token, url, &tree); err != nil {
		return "", err
	}
	for _, item := range tree.Tree {
		if item.Path == targetFile {
			return item.SHA, nil
		}
	}
	return "", errors.New("target not found")
}

func fetchBlobContent(cfg RunnerConfig, logger *slog.Logger, token, sha string) ([]byte, error) {
	owner, repo := runnerRepo()
	url := fmt.Sprintf("%s/repos/%s/%s/git/blobs/%s", strings.TrimRight(runnerGitHubAPIBase, "/"), owner, repo, sha)
	var blob gitBlobResponse
	if err := runnerGetJSON(cfg, logger, token, url, &blob); err != nil {
		return nil, err
	}
	content := strings.ReplaceAll(blob.Content, "\n", "")
	return base64.StdEncoding.DecodeString(content)
}

func fetchCommitInfo(cfg RunnerConfig, logger *slog.Logger, token, branch, targetFile string) commitInfo {
	owner, repo := runnerRepo()
	url := fmt.Sprintf("%s/repos/%s/%s/commits?path=%s&sha=%s&per_page=1", strings.TrimRight(runnerGitHubAPIBase, "/"), owner, repo, targetFile, branch)
	var commits gitCommitResponse
	if err := runnerGetJSON(cfg, logger, token, url, &commits); err != nil || len(commits) == 0 {
		return commitInfo{}
	}
	msg := strings.SplitN(commits[0].Commit.Message, "\n", 2)[0]
	sha := commits[0].SHA
	author := commits[0].Commit.Author.Name
	date := commits[0].Commit.Author.Date
	return commitInfo{SHA: &sha, Message: &msg, Author: &author, Date: &date}
}

func runnerRepo() (string, string) {
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
				if wait > 0 {
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

func materializeSource(src string, content []byte) error {
	if err := os.RemoveAll(src); err != nil {
		return err
	}
	if err := os.MkdirAll(src, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(src, "source.md"), content, 0644)
}

func precheckRunnerTarget(target BranchTarget) error {
	if err := ensureDiskAvailable(filepath.Dir(target.Out)); err != nil {
		return err
	}
	buildBin := os.Getenv("ADLAIRE_CI_BUILD_BIN")
	if buildBin == "" {
		buildBin = "/usr/local/bin/adlaire-ci-build"
	}
	info, err := os.Stat(buildBin)
	if err != nil {
		return err
	}
	if info.IsDir() || info.Mode()&0111 == 0 {
		return fmt.Errorf("build binary is not executable: %s", buildBin)
	}
	cmd := exec.Command(buildBin, "--version")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	if !strings.Contains(string(out), "adlaire-ci-build") || !strings.Contains(string(out), "ADLAIRE_CI_SPEC") {
		return fmt.Errorf("build binary version mismatch: %s", buildBin)
	}
	return nil
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

func runPipeline(cfg RunnerConfig, target BranchTarget, buildID string) pipelineLog {
	script := filepath.Join(filepath.Dir(target.Src), ".ci", "pipeline.sh")
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", script)
	cmd.Dir = filepath.Dir(target.Src)
	cmd.Env = append(os.Environ(),
		"ADLAIRE_CI_SRC="+target.Src,
		"ADLAIRE_CI_OUT="+target.Out,
		"ADLAIRE_CI_BRANCH="+target.Branch,
		"ADLAIRE_CI_BUILD_ID="+buildID,
		"ADLAIRE_CI_TARGET_FILE="+target.TargetFile,
		"ADLAIRE_CI_STATE_DIR="+cfg.StateDir,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		code = 1
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			code = ee.ExitCode()
		}
	}
	return pipelineLog{ExitCode: &code, Stdout: trimLog(stdout.String()), Stderr: trimLog(stderr.String())}
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

func trimLog(s string) string {
	const max = 1024 * 1024
	if len(s) <= max {
		return strings.TrimSuffix(s, "\n")
	}
	return s[len(s)-max:]
}

func executeDeploy(cfg RunnerConfig, branchIdx, deployIdx int, target BranchTarget, d DeployTarget) deployLog {
	dl := deployLog{Host: d.Host, User: d.User, DestDir: d.DestDir, Status: "success", TransferVerified: true}
	total, uploaded, skipped, bytesUploaded, err := deploySite(target.Out, d)
	dl.FilesTotal = total
	dl.FilesUploaded = uploaded
	dl.FilesSkipped = skipped
	dl.BytesUploaded = bytesUploaded
	if err != nil {
		errStr := err.Error()
		dl.Status = "pending"
		dl.TransferVerified = false
		dl.Error = &errStr
		_ = addPendingTransfer(cfg.PendingFile, pendingTransfer{
			BranchIdx: branchIdx, DeployIdx: deployIdx, Out: target.Out,
			Host: d.Host, User: d.User, DestDir: d.DestDir,
			FailedAt: runnerNow().UTC().Format(time.RFC3339), RetryCount: 1,
		})
	}
	return dl
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
		rel, err := filepath.Rel(out, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		localSHA, err := fileSHA256(path)
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

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func remoteSHA256(d DeployTarget, remotePath string) string {
	cmd := exec.Command("ssh", d.User+"@"+d.Host, "sha256sum", remotePath)
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
	cmd := exec.Command("ssh", d.User+"@"+d.Host, "mkdir", "-p", filepath.ToSlash(filepath.Dir(remotePath)), "&&", "tee", remotePath)
	cmd.Stdin = bytes.NewReader(data)
	return cmd.Run()
}

func remoteJoin(root, rel string) string {
	return strings.TrimRight(root, "/") + "/" + strings.TrimLeft(rel, "/")
}

func writeSnapshot(cfg RunnerConfig, target BranchTarget, buildID string) (string, error) {
	if cfg.HistoryKeepN == 0 {
		return "", errors.New("snapshot disabled")
	}
	dest := filepath.Join(cfg.StateDir, ".snapshots", buildID, "site")
	if err := os.RemoveAll(filepath.Dir(dest)); err != nil {
		return "", err
	}
	if err := copyDir(target.Out, dest); err != nil {
		return "", err
	}
	if err := pruneSnapshots(filepath.Join(cfg.StateDir, ".snapshots"), cfg.HistoryKeepN); err != nil {
		return "", err
	}
	return buildID, nil
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
		if e.IsDir() {
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
		if entries[i].Out == entry.Out && entries[i].Host == entry.Host && entries[i].User == entry.User && entries[i].DestDir == entry.DestDir {
			entries[i].RetryCount++
			entries[i].FailedAt = entry.FailedAt
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
		target := BranchTarget{Out: entry.Out}
		deploy := DeployTarget{Host: entry.Host, User: entry.User, DestDir: entry.DestDir}
		result := executeDeploy(*cfg, entry.BranchIdx, entry.DeployIdx, target, deploy)
		if result.Status == "success" {
			logger.Info("PENDING RETRY OK site -> " + entry.Host)
			continue
		}
		entry.RetryCount++
		entry.FailedAt = runnerNow().UTC().Format(time.RFC3339)
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
	client := &http.Client{Timeout: 30 * time.Second}
	for _, entry := range entries {
		body, _ := json.Marshal(entry.Payload)
		resp, err := client.Post(entry.URL, "application/json", bytes.NewReader(body))
		if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			logger.Info("NOTIFY_PENDING_RETRY_OK: url=" + entry.URL)
			continue
		}
		if resp != nil {
			entry.LastError = fmt.Sprintf("http status %d", resp.StatusCode)
			resp.Body.Close()
		} else if err != nil {
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
	if err != nil || len(config.Webhooks) == 0 {
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
	}
	for _, hook := range config.Webhooks {
		if !hook.Enabled || hook.URL == "" || !hookHandlesEvent(hook, event) {
			continue
		}
		if err := postNotify(hook.URL, payload); err != nil {
			_ = addNotifyPending(filepath.Join(cfg.StateDir, ".notify_pending"), notifyPendingEntry{
				Event: event, URL: hook.URL, Payload: payload,
				QueuedAt: runnerNow().UTC().Format(time.RFC3339), RetryCount: 1, LastError: err.Error(),
			})
			logger.Error("NOTIFY_FAILED: url=" + hook.URL)
		}
	}
}

func readNotifyConfig(path string) (notifyConfigFile, error) {
	var config notifyConfigFile
	err := readJSONFile(path, &config)
	if errors.Is(err, os.ErrNotExist) {
		return config, nil
	}
	return config, err
}

func hookHandlesEvent(hook notifyWebhook, event string) bool {
	for _, item := range hook.On {
		if item == event || item == "*" {
			return true
		}
	}
	return false
}

func postNotify(url string, payload map[string]any) error {
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http status %d", resp.StatusCode)
	}
	return nil
}

func addNotifyPending(path string, entry notifyPendingEntry) error {
	var entries []notifyPendingEntry
	_ = readJSONArray(path, &entries)
	entries = append(entries, entry)
	return runnerAtomicWriteJSON(path, entries, 0600)
}

func logFailure(cfg RunnerConfig, target BranchTarget, buildID string, started time.Time, blobSHA *string, prevSHA, status, msg string, pl *pipelineLog, rep *runnerReport, logger *slog.Logger) {
	if pl == nil {
		pl = &pipelineLog{}
	}
	finished := runnerNow().UTC()
	blog := buildLog{
		ID: buildID, Branch: target.Branch, TargetFile: target.TargetFile, TargetStatus: status,
		StartedAt: started.Format(time.RFC3339), FinishedAt: finished.Format(time.RFC3339),
		DurationSeconds: int64(finished.Sub(started).Seconds()), Commit: commitInfo{},
		BlobSHA: blobSHA, PreviousBlobSHA: prevSHA, Pipeline: *pl, Report: rep,
		Warnings: []string{}, Deploy: []deployLog{}, SnapshotID: nil, Error: &msg,
	}
	if err := writeBuildLog(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_LOG_WRITE_FAILED: " + err.Error())
	}
	if err := appendHistory(cfg.StateDir, blog); err != nil {
		logger.Error("BUILD_HISTORY_WRITE_FAILED: " + err.Error())
	}
	sendBuildNotifications(cfg, blog, logger)
}

func circuitOpen(stateDir string) bool {
	state, err := readCircuitState(stateDir)
	return err == nil && state.Open
}

func recordCircuitFailure(cfg RunnerConfig, msg string) {
	if cfg.APICircuitBreakerThreshold <= 0 {
		return
	}
	state, _ := readCircuitState(cfg.StateDir)
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

func readCircuitState(stateDir string) (buildCircuitState, error) {
	var state buildCircuitState
	err := readJSONFile(filepath.Join(stateDir, ".build_circuit_state"), &state)
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
	if err := readJSONFile(path, &cache); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return cache.SHA, nil
}

func writeBuildState(stateDir string, running bool, buildID *string) error {
	now := runnerNow().UTC().Format(time.RFC3339)
	state, err := readBuildState(stateDir)
	if err != nil {
		state = defaultBuildState()
	}
	state.Running = running
	state.CurrentBuildID = buildID
	state.LastStartedAt = &now
	return runnerAtomicWriteJSON(filepath.Join(stateDir, ".build_state"), state, 0600)
}

func defaultBuildState() buildState {
	return buildState{Queued: []map[string]any{}}
}

func writeBuildLog(stateDir string, log buildLog) error {
	return runnerAtomicWriteJSON(filepath.Join(stateDir, ".build_logs", log.ID+".json"), log, 0600)
}

func appendHistory(stateDir string, log buildLog) error {
	path := filepath.Join(stateDir, ".build_history")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	var pages *int
	warnings := 0
	sizeWarn := false
	if log.Report != nil {
		pages = &log.Report.Pages
		warnings = log.Report.WarningsCount
		sizeWarn = log.Report.SizeWarn
	}
	outputSHA, _ := outputManifestSHA(filepath.Join(stateDir, "dist", "site"))
	rec := historyRecord{
		ID: log.ID, Branch: log.Branch, TargetFile: log.TargetFile, Status: log.TargetStatus,
		StartedAt: log.StartedAt, FinishedAt: log.FinishedAt, DurationSeconds: log.DurationSeconds,
		CommitSHA: log.Commit.SHA, BlobSHA: log.BlobSHA, Pages: pages, Warnings: warnings,
		SizeWarn: sizeWarn, OutputSHA256: outputSHA, RollbackFrom: nil,
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	return f.Sync()
}

func outputManifestSHA(root string) (*string, error) {
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}
	parts := []string{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		rel, _ := filepath.Rel(root, path)
		parts = append(parts, filepath.ToSlash(rel)+"\n"+hex.EncodeToString(sum[:])+"\n")
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

func outputSizeBytes(root string) (int64, bool) {
	var size int64
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		size += info.Size()
		return nil
	})
	return size, err == nil
}

func cleanupBuildLogs(dir string, keep int, logger *slog.Logger) {
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
	for len(files) > keep {
		if err := os.Remove(files[0].path); err != nil {
			logger.Warn("LOG_CLEANUP_FAILED: path=" + files[0].path)
		}
		files = files[1:]
	}
}

func runnerBuildID(t time.Time, n int) string {
	base := "b" + t.UTC().Format("20060102150405")
	if n <= 1 {
		return base
	}
	return fmt.Sprintf("%s-%d", base, n)
}

func repairCorruptJSONArray(path string, logger *slog.Logger) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	var raw []json.RawMessage
	if err := readJSONFile(path, &raw); err == nil {
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
	return readJSONFile(path, out)
}

func readJSONFile(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func runnerAtomicWriteJSON(path string, v any, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	base := filepath.Base(path)
	tmp := filepath.Join(filepath.Dir(path), fmt.Sprintf(".%s.tmp.%d", base, os.Getpid()))
	if err := os.WriteFile(tmp, data, mode); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
