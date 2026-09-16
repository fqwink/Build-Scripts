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
	blobSHA, err := fetchTargetSHA(token, target.Branch, target.TargetFile)
	if err != nil {
		logFailure(cfg, target, buildID, started, nil, prevSHA, "failure_api", "github api failed", nil, nil, logger)
		return 3
	}
	if blobSHA == prevSHA {
		logger.Info(fmt.Sprintf("NO_CHANGE: branch=%s target=%s sha=%s", target.Branch, target.TargetFile, blobSHA))
		return 0
	}
	content, err := fetchBlobContent(token, blobSHA)
	if err != nil {
		logFailure(cfg, target, buildID, started, &blobSHA, prevSHA, "failure_decode", "blob decode failed", nil, nil, logger)
		return 1
	}
	if err := materializeSource(target.Src, content); err != nil {
		logFailure(cfg, target, buildID, started, &blobSHA, prevSHA, "failure_decode", "source write failed", nil, nil, logger)
		return 1
	}
	pl := runPipeline(cfg, target, buildID)
	rep, warns := parseRunnerReport(pl.Stdout)
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
	finished := runnerNow().UTC()
	blog := buildLog{
		ID:              buildID,
		Branch:          target.Branch,
		TargetFile:      target.TargetFile,
		TargetStatus:    status,
		StartedAt:       started.Format(time.RFC3339),
		FinishedAt:      finished.Format(time.RFC3339),
		DurationSeconds: int64(finished.Sub(started).Seconds()),
		Commit:          commitInfo{},
		BlobSHA:         &blobSHA,
		PreviousBlobSHA: prevSHA,
		Pipeline:        pl,
		Report:          rep,
		Warnings:        warns,
		Deploy:          deploys,
		SnapshotID:      nil,
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
	return exitCode
}

func fetchTargetSHA(token, branch, targetFile string) (string, error) {
	owner, repo := runnerRepo()
	url := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=1", strings.TrimRight(runnerGitHubAPIBase, "/"), owner, repo, branch)
	var tree gitTreeResponse
	if err := runnerGetJSON(token, url, &tree); err != nil {
		return "", err
	}
	for _, item := range tree.Tree {
		if item.Path == targetFile {
			return item.SHA, nil
		}
	}
	return "", errors.New("target not found")
}

func fetchBlobContent(token, sha string) ([]byte, error) {
	owner, repo := runnerRepo()
	url := fmt.Sprintf("%s/repos/%s/%s/git/blobs/%s", strings.TrimRight(runnerGitHubAPIBase, "/"), owner, repo, sha)
	var blob gitBlobResponse
	if err := runnerGetJSON(token, url, &blob); err != nil {
		return nil, err
	}
	content := strings.ReplaceAll(blob.Content, "\n", "")
	return base64.StdEncoding.DecodeString(content)
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

func runnerGetJSON(token, url string, out any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "adlaire-ci-runner")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("github api status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
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
	errStr := ""
	dl := deployLog{Host: d.Host, User: d.User, DestDir: d.DestDir, Status: "success", TransferVerified: true}
	cmd := exec.Command("ssh", d.User+"@"+d.Host, "mkdir -p "+d.DestDir)
	if err := cmd.Run(); err != nil {
		errStr = err.Error()
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
	state := defaultBuildState()
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
	bak := fmt.Sprintf("%s.corrupt.%s.bak", path, runnerNow().UTC().Format("20060102150405"))
	if err := os.Rename(path, bak); err != nil {
		return err
	}
	logger.Warn("STATE_CORRUPT_BACKUP: path=" + bak)
	return runnerAtomicWriteJSON(path, []any{}, 0600)
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
