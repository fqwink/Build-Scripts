package setup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

const (
	setupBinaryName         = "adlaire-ci-setup"
	setupDefaultRepository  = "fqwink/Build-Scripts"
	setupDefaultInstallDir  = "/opt/adlaire-builder"
	setupDefaultBinDir      = "/usr/local/bin"
	setupDefaultServiceUser = "root"
	setupDefaultOSArch      = "linux-amd64"
	setupSystemdDir         = "/etc/systemd/system"
	setupHTTPTimeout        = 5 * time.Minute
	setupHealthTimeout      = 10 * time.Second

	setupBinaryLimitBytes   = 128 * 1024 * 1024
	setupAdminUILimitBytes  = 16 * 1024 * 1024
	setupSHA256LimitBytes   = 64 * 1024
	setupHealthLimitBytes   = 1024 * 1024
	setupSecretInputMaxSize = 64 * 1024
)

var setupBinaryVersion = "V.0.0-dev"

func SetBinaryVersion(version string) {
	setupBinaryVersion = version
}

var (
	setupVersionPattern    = regexp.MustCompile(`^V\.[1-9][0-9]*\.[0-9]+$`)
	setupRepositoryPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)
)

type setupMode string

const (
	setupModeInstall    setupMode = "install"
	setupModeInstallAPI setupMode = "install-api"
	setupModeUpdate     setupMode = "update"
)

type setupConfig struct {
	Mode                setupMode
	TargetVersion       string
	Repository          string
	OSArch              string
	InstallDir          string
	BinDir              string
	ServiceUser         string
	DownloadDir         string
	GitHubTokenFile     string
	InitialPasswordFile string
}

type setupAsset struct {
	Name       string
	BinaryName string
	Limit      int64
}

type setupRunState struct {
	Assets        []setupAsset
	APIInstalled  bool
	MCPInstalled  bool
	GitHubToken   string
	HasCredential bool
	Backup        setupBackup
}

type setupBackup struct {
	Root      string
	Binaries  map[string]string
	AdminRoot string
}

type setupFailure struct {
	Code     string
	Stage    string
	Rollback string
}

func (f setupFailure) exitCode() int {
	switch f.Code {
	case "INVALID_INPUT", "UNSUPPORTED_PLATFORM", "PRECONDITION_FAILED", "SECRET_INPUT_INVALID", "EXISTING_STATE_INVALID":
		return 2
	case "DOWNLOAD_FAILED":
		return 3
	default:
		return 1
	}
}

type setupCommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type setupOperations interface {
	ValidatePath(path string) error
	ReadSecretInput(path string) (string, error)
	ValidateExistingToken(path string) error
	ValidateLastSHA(path string) error
	ValidateCredentials(path string) error
	EnsureDownloadDir(path, version string) error
	DownloadAsset(cfg setupConfig, asset setupAsset) error
	VerifyChecksums(downloadDir string, assets []setupAsset, osArch string) error
	InstallBinary(src, dest, binaryName, version string) error
	VerifyBinary(path, binaryName, version string) error
	BinaryExists(path string) bool
	EnsureDir(path string, mode os.FileMode) error
	CreateFileExclusive(path string, data []byte, mode os.FileMode) (bool, error)
	InstallAdminArchive(archivePath, installDir string) error
	WriteRunnerUnits(cfg setupConfig) error
	WriteAPIUnit(cfg setupConfig) error
	InitCredentials(apiPath, installDir, passwordFile string) error
	NewBackupDir() (string, error)
	BackupFile(src, dest string) error
	RestoreFile(src, dest string) error
	BackupDir(src, dest string) error
	RestoreDir(src, dest string) error
	RemoveAll(path string) error
	Systemctl(args ...string) setupCommandResult
	HealthCheck() error
}

type setupRealOps struct{}

func RunSetup(args []string, stdout, stderr io.Writer) int {
	return runSetup(args, stdout, stderr, setupRealOps{})
}

func runSetup(args []string, stdout, stderr io.Writer, ops setupOperations) int {
	if hasExactArg(args, "--help") {
		fmt.Fprintln(stdout, "Usage: adlaire-ci-setup <install|install-api|update> --target-version V.X.N [options] [--version] [--help]")
		return 0
	}
	if hasExactArg(args, "--version") {
		fmt.Fprintf(stdout, "%s %s go=%s\n", setupBinaryName, setupBinaryVersion, runtime.Version())
		return 0
	}
	if !safeArgvTokens(args) {
		fmt.Fprintln(stderr, "invalid command line token")
		return 2
	}

	cfg, parseErr := parseSetupArgs(args)
	if parseErr != "" {
		fmt.Fprintln(stderr, parseErr)
		return 2
	}

	state := &setupRunState{}
	if failure := runSetupStage(cfg, state, stdout, "validate", func() *setupFailure {
		return validateSetupConfig(cfg, state, ops)
	}); failure != nil {
		return finishSetupFailure(stderr, *failure)
	}

	stages := setupStages(cfg, state, ops)
	for _, stage := range stages {
		failure := runSetupStage(cfg, state, stdout, stage.name, stage.run)
		if failure != nil {
			if cfg.Mode == setupModeUpdate && setupStageNeedsRollback(stage.name) {
				if rollbackUpdate(cfg, state, ops, stage.name) {
					failure.Rollback = "completed"
				} else {
					failure.Code = "ROLLBACK_FAILED"
					failure.Rollback = "failed"
				}
			}
			return finishSetupFailure(stderr, *failure)
		}
	}

	fmt.Fprintf(stdout, "setup: success %s %s\n", cfg.Mode, cfg.TargetVersion)
	return 0
}

func finishSetupFailure(stderr io.Writer, failure setupFailure) int {
	if failure.Rollback == "" {
		failure.Rollback = "none"
	}
	fmt.Fprintf(stderr, "setup: error %s stage=%s rollback=%s\n", failure.Code, failure.Stage, failure.Rollback)
	return failure.exitCode()
}

type setupStage struct {
	name string
	run  func() *setupFailure
}

func runSetupStage(cfg setupConfig, state *setupRunState, stdout io.Writer, name string, run func() *setupFailure) *setupFailure {
	fmt.Fprintf(stdout, "setup: start %s %s\n", cfg.Mode, name)
	if failure := run(); failure != nil {
		if failure.Stage == "" {
			failure.Stage = name
		}
		if failure.Rollback == "" {
			failure.Rollback = "none"
		}
		return failure
	}
	fmt.Fprintf(stdout, "setup: ok %s %s\n", cfg.Mode, name)
	return nil
}

func setupStages(cfg setupConfig, state *setupRunState, ops setupOperations) []setupStage {
	switch cfg.Mode {
	case setupModeInstall:
		return []setupStage{
			{name: "download", run: func() *setupFailure { return setupDownload(cfg, state, ops) }},
			{name: "checksum", run: func() *setupFailure { return setupChecksum(cfg, state, ops) }},
			{name: "install-binaries", run: func() *setupFailure { return setupInstallBinaries(cfg, state, ops) }},
			{name: "initialize-secrets", run: func() *setupFailure { return setupInitializeSecrets(cfg, state, ops) }},
			{name: "initialize-state", run: func() *setupFailure { return setupInitializeState(cfg, state, ops) }},
			{name: "configure-systemd", run: func() *setupFailure { return setupConfigureSystemd(cfg, state, ops) }},
			{name: "activate-services", run: func() *setupFailure { return setupActivateServices(cfg, state, ops) }},
			{name: "verify", run: func() *setupFailure { return setupVerify(cfg, state, ops) }},
		}
	case setupModeInstallAPI:
		return []setupStage{
			{name: "download", run: func() *setupFailure { return setupDownload(cfg, state, ops) }},
			{name: "checksum", run: func() *setupFailure { return setupChecksum(cfg, state, ops) }},
			{name: "prepare-runtime", run: func() *setupFailure { return setupPrepareRuntime(cfg, state, ops) }},
			{name: "install-api-binary", run: func() *setupFailure { return setupInstallBinaries(cfg, state, ops) }},
			{name: "install-admin", run: func() *setupFailure { return setupInstallAdmin(cfg, state, ops) }},
			{name: "initialize-credentials", run: func() *setupFailure { return setupInitializeCredentials(cfg, state, ops) }},
			{name: "configure-systemd", run: func() *setupFailure { return setupConfigureSystemd(cfg, state, ops) }},
			{name: "activate-api", run: func() *setupFailure { return setupActivateAPI(cfg, state, ops) }},
			{name: "verify", run: func() *setupFailure { return setupVerify(cfg, state, ops) }},
		}
	case setupModeUpdate:
		stages := []setupStage{
			{name: "backup", run: func() *setupFailure { return setupBackupCurrent(cfg, state, ops) }},
			{name: "download", run: func() *setupFailure { return setupDownload(cfg, state, ops) }},
			{name: "checksum", run: func() *setupFailure { return setupChecksum(cfg, state, ops) }},
			{name: "install-binaries", run: func() *setupFailure { return setupInstallBinaries(cfg, state, ops) }},
		}
		if state.APIInstalled {
			stages = append(stages, setupStage{name: "install-admin", run: func() *setupFailure { return setupInstallAdmin(cfg, state, ops) }})
		}
		stages = append(stages,
			setupStage{name: "restart-services", run: func() *setupFailure { return setupRestartServices(cfg, state, ops) }},
			setupStage{name: "verify", run: func() *setupFailure { return setupVerify(cfg, state, ops) }},
		)
		return stages
	default:
		return nil
	}
}

func setupStageNeedsRollback(stage string) bool {
	switch stage {
	case "install-binaries", "install-admin", "restart-services", "verify":
		return true
	default:
		return false
	}
}

func parseSetupArgs(args []string) (setupConfig, string) {
	cfg := setupConfig{
		Repository:  setupDefaultRepository,
		OSArch:      setupDefaultOSArch,
		InstallDir:  setupDefaultInstallDir,
		BinDir:      setupDefaultBinDir,
		ServiceUser: setupDefaultServiceUser,
	}
	if len(args) == 0 {
		return cfg, "usage error"
	}
	if strings.HasPrefix(args[0], "--") || strings.HasPrefix(args[0], "-") {
		return cfg, "unknown option: " + args[0]
	}
	switch setupMode(args[0]) {
	case setupModeInstall, setupModeInstallAPI, setupModeUpdate:
		cfg.Mode = setupMode(args[0])
	default:
		return cfg, "unknown command: " + args[0]
	}

	seen := map[string]bool{}
	for i := 1; i < len(args); {
		name := args[i]
		if !strings.HasPrefix(name, "-") {
			return cfg, "usage error"
		}
		if !strings.HasPrefix(name, "--") || strings.Contains(name, "=") {
			return cfg, "unknown option: " + name
		}
		if !setupKnownOption(name) {
			return cfg, "unknown option: " + name
		}
		if seen[name] {
			return cfg, "usage error"
		}
		seen[name] = true
		if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
			return cfg, "missing value: " + name
		}
		value := args[i+1]
		if !setupOptionAllowed(cfg.Mode, name) {
			return cfg, "usage error"
		}
		switch name {
		case "--target-version":
			cfg.TargetVersion = value
		case "--repository":
			cfg.Repository = value
		case "--os-arch":
			cfg.OSArch = value
		case "--install-dir":
			cfg.InstallDir = value
		case "--bin-dir":
			cfg.BinDir = value
		case "--service-user":
			cfg.ServiceUser = value
		case "--download-dir":
			cfg.DownloadDir = value
		case "--github-token-file":
			cfg.GitHubTokenFile = value
		case "--initial-password-file":
			cfg.InitialPasswordFile = value
		}
		i += 2
	}
	if cfg.TargetVersion == "" {
		return cfg, "usage error"
	}
	if cfg.DownloadDir == "" {
		cfg.DownloadDir = filepath.Join("/tmp", "adlaire-ci-release-"+cfg.TargetVersion)
	}
	if cfg.Mode == setupModeInstall && cfg.GitHubTokenFile == "" {
		return cfg, "usage error"
	}
	return cfg, ""
}

func setupKnownOption(name string) bool {
	switch name {
	case "--target-version", "--repository", "--os-arch", "--install-dir", "--bin-dir", "--service-user", "--download-dir", "--github-token-file", "--initial-password-file":
		return true
	default:
		return false
	}
}

func setupOptionAllowed(mode setupMode, name string) bool {
	switch name {
	case "--github-token-file":
		return mode == setupModeInstall
	case "--initial-password-file":
		return mode == setupModeInstallAPI
	case "--service-user":
		return mode == setupModeInstall || mode == setupModeInstallAPI
	default:
		return true
	}
}

func validateSetupConfig(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	if !setupVersionPattern.MatchString(cfg.TargetVersion) || !setupRepositoryPattern.MatchString(cfg.Repository) {
		return setupFail("INVALID_INPUT", "validate")
	}
	if cfg.OSArch != setupDefaultOSArch {
		return setupFail("UNSUPPORTED_PLATFORM", "validate")
	}
	if !validSetupServiceUser(cfg.ServiceUser) {
		return setupFail("INVALID_INPUT", "validate")
	}
	for _, path := range []string{cfg.InstallDir, cfg.BinDir, cfg.DownloadDir} {
		if err := ops.ValidatePath(path); err != nil {
			return setupFail("INVALID_INPUT", "validate")
		}
	}
	if setupDangerousPathRelation(cfg.InstallDir, cfg.BinDir) || setupDangerousPathRelation(cfg.InstallDir, cfg.DownloadDir) || setupDangerousPathRelation(cfg.BinDir, cfg.DownloadDir) {
		return setupFail("INVALID_INPUT", "validate")
	}

	switch cfg.Mode {
	case setupModeInstall:
		token, err := ops.ReadSecretInput(cfg.GitHubTokenFile)
		if err != nil {
			return setupFail("SECRET_INPUT_INVALID", "validate")
		}
		state.GitHubToken = token
		state.Assets = []setupAsset{
			setupBinaryAsset("adlaire-ci-build", cfg.OSArch),
			setupBinaryAsset("adlaire-ci-runner", cfg.OSArch),
			setupBinaryAsset("adlaire-ci-setup", cfg.OSArch),
		}
	case setupModeInstallAPI:
		if !setupRunnerCohortReady(cfg, ops) {
			return setupFail("PRECONDITION_FAILED", "validate")
		}
		credPath := filepath.Join(cfg.InstallDir, ".admin_credentials")
		if ops.BinaryExists(credPath) {
			if err := ops.ValidateCredentials(credPath); err != nil {
				return setupFail("EXISTING_STATE_INVALID", "validate")
			}
			state.HasCredential = true
		} else {
			if cfg.InitialPasswordFile == "" {
				return setupFail("SECRET_INPUT_INVALID", "validate")
			}
			if _, err := ops.ReadSecretInput(cfg.InitialPasswordFile); err != nil {
				return setupFail("SECRET_INPUT_INVALID", "validate")
			}
		}
		state.Assets = []setupAsset{
			setupBinaryAsset("adlaire-ci-api", cfg.OSArch),
			setupBinaryAsset("adlaire-ci-admin", cfg.OSArch),
			setupNamedAsset("admin-ui.tar.gz", setupAdminUILimitBytes),
		}
	case setupModeUpdate:
		for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup"} {
			if !ops.BinaryExists(filepath.Join(cfg.BinDir, name)) {
				return setupFail("PRECONDITION_FAILED", "validate")
			}
		}
		apiEnabled := ops.Systemctl("is-enabled", "adlaire-ci-api")
		state.APIInstalled = ops.BinaryExists(filepath.Join(cfg.BinDir, "adlaire-ci-api")) && apiEnabled.ExitCode == 0 && (strings.TrimSpace(apiEnabled.Stdout) == "enabled" || strings.TrimSpace(apiEnabled.Stdout) == "static")
		state.MCPInstalled = ops.BinaryExists(filepath.Join(cfg.BinDir, "adlaire-ci-mcp"))
		state.Assets = []setupAsset{
			setupBinaryAsset("adlaire-ci-build", cfg.OSArch),
			setupBinaryAsset("adlaire-ci-runner", cfg.OSArch),
			setupBinaryAsset("adlaire-ci-setup", cfg.OSArch),
		}
		if state.APIInstalled {
			state.Assets = append(state.Assets,
				setupBinaryAsset("adlaire-ci-api", cfg.OSArch),
				setupBinaryAsset("adlaire-ci-admin", cfg.OSArch),
				setupNamedAsset("admin-ui.tar.gz", setupAdminUILimitBytes),
			)
		}
		if state.MCPInstalled {
			state.Assets = append(state.Assets, setupBinaryAsset("adlaire-ci-mcp", cfg.OSArch))
		}
	}
	return nil
}

func setupFail(code, stage string) *setupFailure {
	return &setupFailure{Code: code, Stage: stage, Rollback: "none"}
}

func validSetupServiceUser(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r == 0 || r == ':' || r == '\r' || r == '\n' || r == '\t' || r == ' ' {
			return false
		}
	}
	return true
}

func setupDangerousPathRelation(a, b string) bool {
	a = filepath.Clean(a)
	b = filepath.Clean(b)
	if a == b {
		return true
	}
	rel, err := filepath.Rel(a, b)
	if err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
		return true
	}
	rel, err = filepath.Rel(b, a)
	if err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
		return true
	}
	return false
}

func setupRunnerCohortReady(cfg setupConfig, ops setupOperations) bool {
	for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup"} {
		path := filepath.Join(cfg.BinDir, name)
		if !ops.BinaryExists(path) || ops.VerifyBinary(path, name, cfg.TargetVersion) != nil {
			return false
		}
	}
	return setupRunnerUnitActive(ops)
}

func setupRunnerUnitActive(ops setupOperations) bool {
	active := ops.Systemctl("is-active", "adlaire-ci.timer")
	if active.ExitCode != 0 || strings.TrimSpace(active.Stdout) != "active" {
		return false
	}
	for _, unit := range []string{"adlaire-ci.service", "adlaire-ci.timer"} {
		if result := ops.Systemctl("cat", unit); result.ExitCode != 0 {
			return false
		}
	}
	return true
}

func setupDownload(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	if err := ops.EnsureDownloadDir(cfg.DownloadDir, cfg.TargetVersion); err != nil {
		return setupFail("INSTALL_FAILED", "download")
	}
	for _, asset := range append([]setupAsset{}, append(state.Assets, setupNamedAsset("SHA256SUMS", setupSHA256LimitBytes))...) {
		if err := ops.DownloadAsset(cfg, asset); err != nil {
			return setupFail("DOWNLOAD_FAILED", "download")
		}
	}
	return nil
}

func setupChecksum(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	if err := ops.VerifyChecksums(cfg.DownloadDir, state.Assets, cfg.OSArch); err != nil {
		return setupFail("CHECKSUM_FAILED", "checksum")
	}
	return nil
}

func setupInstallBinaries(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	for _, asset := range state.Assets {
		if asset.BinaryName == "" {
			continue
		}
		src := filepath.Join(cfg.DownloadDir, asset.Name)
		dest := filepath.Join(cfg.BinDir, asset.BinaryName)
		if err := ops.InstallBinary(src, dest, asset.BinaryName, cfg.TargetVersion); err != nil {
			return setupFail("INSTALL_FAILED", currentInstallStage(cfg))
		}
	}
	return nil
}

func currentInstallStage(cfg setupConfig) string {
	if cfg.Mode == setupModeInstallAPI {
		return "install-api-binary"
	}
	return "install-binaries"
}

func setupInitializeSecrets(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	target := filepath.Join(cfg.InstallDir, ".github_token")
	created, err := ops.CreateFileExclusive(target, []byte(state.GitHubToken+"\n"), 0600)
	if err != nil {
		return setupFail("SECRET_INIT_FAILED", "initialize-secrets")
	}
	if !created {
		if err := ops.ValidateExistingToken(target); err != nil {
			return setupFail("EXISTING_STATE_INVALID", "initialize-secrets")
		}
	}
	return nil
}

func setupInitializeState(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	target := filepath.Join(cfg.InstallDir, ".last_sha")
	created, err := ops.CreateFileExclusive(target, []byte("{\"sha\":\"\"}\n"), 0600)
	if err != nil {
		return setupFail("STATE_INIT_FAILED", "initialize-state")
	}
	if !created {
		if err := ops.ValidateLastSHA(target); err != nil {
			return setupFail("EXISTING_STATE_INVALID", "initialize-state")
		}
	}
	return nil
}

func setupPrepareRuntime(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	for _, dir := range []string{filepath.Join(cfg.InstallDir, ".build_logs"), filepath.Join(cfg.InstallDir, ".snapshots")} {
		if err := ops.EnsureDir(dir, 0700); err != nil {
			return setupFail("INSTALL_FAILED", "prepare-runtime")
		}
	}
	return nil
}

func setupInstallAdmin(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	if err := ops.InstallAdminArchive(filepath.Join(cfg.DownloadDir, "admin-ui.tar.gz"), cfg.InstallDir); err != nil {
		return setupFail("ADMIN_INSTALL_FAILED", "install-admin")
	}
	return nil
}

func setupInitializeCredentials(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	credPath := filepath.Join(cfg.InstallDir, ".admin_credentials")
	if state.HasCredential {
		if err := ops.ValidateCredentials(credPath); err != nil {
			return setupFail("EXISTING_STATE_INVALID", "initialize-credentials")
		}
		return nil
	}
	if err := ops.InitCredentials(filepath.Join(cfg.BinDir, "adlaire-ci-api"), cfg.InstallDir, cfg.InitialPasswordFile); err != nil {
		return setupFail("CREDENTIALS_INIT_FAILED", "initialize-credentials")
	}
	if err := ops.ValidateCredentials(credPath); err != nil {
		return setupFail("CREDENTIALS_INIT_FAILED", "initialize-credentials")
	}
	return nil
}

func setupConfigureSystemd(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	if err := ops.WriteRunnerUnits(cfg); err != nil {
		return setupFail("SYSTEMD_FAILED", "configure-systemd")
	}
	if cfg.Mode == setupModeInstallAPI {
		if err := ops.WriteAPIUnit(cfg); err != nil {
			return setupFail("SYSTEMD_FAILED", "configure-systemd")
		}
	}
	if result := ops.Systemctl("daemon-reload"); result.ExitCode != 0 {
		return setupFail("SYSTEMD_FAILED", "configure-systemd")
	}
	return nil
}

func setupActivateServices(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	if result := ops.Systemctl("enable", "--now", "adlaire-ci.timer"); result.ExitCode != 0 {
		return setupFail("SYSTEMD_FAILED", "activate-services")
	}
	if !setupRunnerUnitActive(ops) {
		return setupFail("SYSTEMD_FAILED", "activate-services")
	}
	return nil
}

func setupActivateAPI(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	if result := ops.Systemctl("enable", "--now", "adlaire-ci-api"); result.ExitCode != 0 {
		return setupFail("SYSTEMD_FAILED", "activate-api")
	}
	if !setupRunnerUnitActive(ops) {
		return setupFail("SYSTEMD_FAILED", "activate-api")
	}
	if active := ops.Systemctl("is-active", "adlaire-ci-api"); active.ExitCode != 0 || strings.TrimSpace(active.Stdout) != "active" {
		return setupFail("SYSTEMD_FAILED", "activate-api")
	}
	if cat := ops.Systemctl("cat", "adlaire-ci-api.service"); cat.ExitCode != 0 {
		return setupFail("SYSTEMD_FAILED", "activate-api")
	}
	if err := ops.HealthCheck(); err != nil {
		return setupFail("HEALTHCHECK_FAILED", "activate-api")
	}
	return nil
}

func setupBackupCurrent(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	root, err := ops.NewBackupDir()
	if err != nil {
		return setupFail("INSTALL_FAILED", "backup")
	}
	state.Backup = setupBackup{Root: root, Binaries: map[string]string{}}
	for _, asset := range state.Assets {
		if asset.BinaryName == "" {
			continue
		}
		src := filepath.Join(cfg.BinDir, asset.BinaryName)
		dst := filepath.Join(root, asset.BinaryName)
		if err := ops.BackupFile(src, dst); err != nil {
			return setupFail("PRECONDITION_FAILED", "backup")
		}
		state.Backup.Binaries[asset.BinaryName] = dst
	}
	if state.APIInstalled {
		src := filepath.Join(cfg.InstallDir, "admin")
		dst := filepath.Join(root, "admin")
		if err := ops.BackupDir(src, dst); err != nil {
			return setupFail("PRECONDITION_FAILED", "backup")
		}
		state.Backup.AdminRoot = dst
	}
	return nil
}

func setupRestartServices(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	if result := ops.Systemctl("restart", "adlaire-ci.timer"); result.ExitCode != 0 {
		return setupFail("SYSTEMD_FAILED", "restart-services")
	}
	if !setupRunnerUnitActive(ops) {
		return setupFail("SYSTEMD_FAILED", "restart-services")
	}
	if state.APIInstalled {
		if result := ops.Systemctl("restart", "adlaire-ci-api"); result.ExitCode != 0 {
			return setupFail("SYSTEMD_FAILED", "restart-services")
		}
		if active := ops.Systemctl("is-active", "adlaire-ci-api"); active.ExitCode != 0 || strings.TrimSpace(active.Stdout) != "active" {
			return setupFail("SYSTEMD_FAILED", "restart-services")
		}
		if cat := ops.Systemctl("cat", "adlaire-ci-api.service"); cat.ExitCode != 0 {
			return setupFail("SYSTEMD_FAILED", "restart-services")
		}
		if err := ops.HealthCheck(); err != nil {
			return setupFail("HEALTHCHECK_FAILED", "restart-services")
		}
	}
	return nil
}

func setupVerify(cfg setupConfig, state *setupRunState, ops setupOperations) *setupFailure {
	for _, asset := range state.Assets {
		if asset.BinaryName == "" {
			continue
		}
		if err := ops.VerifyBinary(filepath.Join(cfg.BinDir, asset.BinaryName), asset.BinaryName, cfg.TargetVersion); err != nil {
			return setupFail("VERIFY_FAILED", "verify")
		}
	}
	if cfg.Mode == setupModeInstall {
		if err := ops.ValidateExistingToken(filepath.Join(cfg.InstallDir, ".github_token")); err != nil {
			return setupFail("VERIFY_FAILED", "verify")
		}
		if err := ops.ValidateLastSHA(filepath.Join(cfg.InstallDir, ".last_sha")); err != nil {
			return setupFail("VERIFY_FAILED", "verify")
		}
	}
	if cfg.Mode == setupModeInstall || cfg.Mode == setupModeUpdate {
		if !setupRunnerUnitActive(ops) {
			return setupFail("VERIFY_FAILED", "verify")
		}
	}
	if cfg.Mode == setupModeInstallAPI || (cfg.Mode == setupModeUpdate && state.APIInstalled) {
		if err := ops.ValidateCredentials(filepath.Join(cfg.InstallDir, ".admin_credentials")); err != nil {
			return setupFail("VERIFY_FAILED", "verify")
		}
		if !ops.BinaryExists(filepath.Join(cfg.InstallDir, "admin", "index.html")) || !ops.BinaryExists(filepath.Join(cfg.InstallDir, "admin", "adlaire-ci-sdk.js")) {
			return setupFail("VERIFY_FAILED", "verify")
		}
		if active := ops.Systemctl("is-active", "adlaire-ci-api"); active.ExitCode != 0 || strings.TrimSpace(active.Stdout) != "active" {
			return setupFail("VERIFY_FAILED", "verify")
		}
		if cat := ops.Systemctl("cat", "adlaire-ci-api.service"); cat.ExitCode != 0 {
			return setupFail("VERIFY_FAILED", "verify")
		}
		if !setupRunnerUnitActive(ops) {
			return setupFail("VERIFY_FAILED", "verify")
		}
		if err := ops.HealthCheck(); err != nil {
			return setupFail("VERIFY_FAILED", "verify")
		}
	}
	return nil
}

func rollbackUpdate(cfg setupConfig, state *setupRunState, ops setupOperations, failedStage string) bool {
	if state.Backup.Root == "" {
		return false
	}
	for name, backup := range state.Backup.Binaries {
		if err := ops.RestoreFile(backup, filepath.Join(cfg.BinDir, name)); err != nil {
			return false
		}
	}
	if state.APIInstalled && state.Backup.AdminRoot != "" {
		if err := ops.RestoreDir(state.Backup.AdminRoot, filepath.Join(cfg.InstallDir, "admin")); err != nil {
			return false
		}
	}
	if failedStage == "restart-services" || failedStage == "install-admin" || failedStage == "verify" {
		if result := ops.Systemctl("restart", "adlaire-ci.timer"); result.ExitCode != 0 {
			return false
		}
		if state.APIInstalled {
			if result := ops.Systemctl("restart", "adlaire-ci-api"); result.ExitCode != 0 {
				return false
			}
		}
	}
	if !setupRunnerUnitActive(ops) {
		return false
	}
	if state.APIInstalled {
		if active := ops.Systemctl("is-active", "adlaire-ci-api"); active.ExitCode != 0 || strings.TrimSpace(active.Stdout) != "active" {
			return false
		}
		if cat := ops.Systemctl("cat", "adlaire-ci-api.service"); cat.ExitCode != 0 {
			return false
		}
		if err := ops.HealthCheck(); err != nil {
			return false
		}
	}
	return true
}

func setupBinaryAsset(binaryName, osArch string) setupAsset {
	return setupAsset{Name: binaryName + "-" + osArch, BinaryName: binaryName, Limit: setupBinaryLimitBytes}
}

func setupNamedAsset(name string, limit int64) setupAsset {
	return setupAsset{Name: name, Limit: limit}
}

func setupReleaseAssetNames(osArch string) map[string]bool {
	names := map[string]bool{
		"admin-ui.tar.gz": true,
	}
	for _, binary := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-api", "adlaire-ci-setup", "adlaire-ci-admin", "adlaire-ci-mcp"} {
		names[binary+"-"+osArch] = true
	}
	return names
}

func (setupRealOps) ValidatePath(path string) error {
	if path == "" || !filepath.IsAbs(path) || filepath.Clean(path) != path || path == string(filepath.Separator) {
		return errors.New("invalid path")
	}
	parts := strings.Split(strings.TrimPrefix(path, string(filepath.Separator)), string(filepath.Separator))
	current := string(filepath.Separator)
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return errors.New("invalid path")
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return errors.New("symlink path segment")
		}
	}
	return nil
}

func (setupRealOps) ReadSecretInput(path string) (string, error) {
	value, err := readSetupSecretFile(path)
	if err != nil {
		return "", err
	}
	if !validSetupSecretValue(value) {
		return "", errors.New("invalid secret")
	}
	return value, nil
}

func (setupRealOps) ValidateExistingToken(path string) error {
	value, err := readSetupSecretFile(path)
	if err != nil {
		return err
	}
	if !validSetupSecretValue(value) {
		return errors.New("invalid existing token")
	}
	return nil
}

func readSetupSecretFile(path string) (string, error) {
	if path == "" || !filepath.IsAbs(path) {
		return "", errors.New("invalid secret path")
	}
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Mode().Perm() != 0600 || !setupOwnedByCurrentUser(info) {
		return "", errors.New("invalid secret file")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if len(data) == 0 || len(data) > setupSecretInputMaxSize || !utf8.Valid(data) {
		return "", errors.New("invalid secret data")
	}
	return strings.TrimSpace(string(data)), nil
}

func validSetupSecretValue(value string) bool {
	if value == "" || !utf8.ValidString(value) {
		return false
	}
	for _, r := range value {
		if r <= 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

func setupOwnedByCurrentUser(info os.FileInfo) bool {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return false
	}
	return stat.Uid == uint32(os.Geteuid())
}

func (setupRealOps) ValidateLastSHA(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0600 {
		return errors.New("invalid state file")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var value struct {
		SHA string `json:"sha"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	for _, r := range value.SHA {
		if r <= 0x20 || r == 0x7f {
			return errors.New("invalid sha")
		}
	}
	return nil
}

func (setupRealOps) ValidateCredentials(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0600 {
		return errors.New("invalid credentials file")
	}
	return validateCredentials(path)
}

func (setupRealOps) EnsureDownloadDir(path, version string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	marker := filepath.Join(path, ".adlaire-ci-setup-version")
	if data, err := os.ReadFile(marker); err == nil {
		if strings.TrimSpace(string(data)) != version {
			return errors.New("download dir belongs to another version")
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(marker, []byte(version+"\n"), 0644)
}

func (ops setupRealOps) DownloadAsset(cfg setupConfig, asset setupAsset) error {
	target := filepath.Join(cfg.DownloadDir, asset.Name)
	partial := target + ".partial"
	if err := os.Remove(partial); err != nil && !os.IsNotExist(err) {
		return err
	}
	downloadURL, err := setupReleaseURL(cfg.Repository, cfg.TargetVersion, asset.Name)
	if err != nil {
		return err
	}
	client := &http.Client{
		Timeout: setupHTTPTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return errors.New("too many redirects")
			}
			if !setupAllowedRedirect(req.URL) {
				return errors.New("unsafe redirect")
			}
			req.Header.Del("Authorization")
			req.Header.Del("Cookie")
			req.Header.Del("Referer")
			return nil
		},
	}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("download status %d", resp.StatusCode)
	}
	if resp.ContentLength == 0 || resp.ContentLength > asset.Limit {
		return errors.New("invalid content length")
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, asset.Limit+1))
	if err != nil {
		return err
	}
	if int64(len(data)) == 0 || int64(len(data)) > asset.Limit {
		return errors.New("invalid response size")
	}
	if err := os.WriteFile(partial, data, 0644); err != nil {
		return err
	}
	return os.Rename(partial, target)
}

func setupReleaseURL(repository, version, asset string) (string, error) {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return "", errors.New("invalid repository")
	}
	return "https://github.com/" + url.PathEscape(parts[0]) + "/" + url.PathEscape(parts[1]) + "/releases/download/" + url.PathEscape(version) + "/" + url.PathEscape(asset), nil
}

func setupAllowedRedirect(u *url.URL) bool {
	if u == nil || u.Scheme != "https" || u.User != nil || u.Fragment != "" {
		return false
	}
	if net.ParseIP(u.Hostname()) != nil {
		return false
	}
	switch u.Hostname() {
	case "github.com", "objects.githubusercontent.com", "release-assets.githubusercontent.com":
		return true
	default:
		return false
	}
}

func (setupRealOps) VerifyChecksums(downloadDir string, assets []setupAsset, osArch string) error {
	raw, err := os.ReadFile(filepath.Join(downloadDir, "SHA256SUMS"))
	if err != nil {
		return err
	}
	if len(raw) == 0 || len(raw) > setupSHA256LimitBytes || bytes.Contains(raw, []byte{0}) {
		return errors.New("invalid SHA256SUMS")
	}
	allowed := setupReleaseAssetNames(osArch)
	sums := map[string]string{}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return errors.New("invalid checksum line")
		}
		hash, name := fields[0], strings.TrimPrefix(fields[1], "*")
		if name == "SHA256SUMS" || !allowed[name] || sums[name] != "" || len(hash) != 64 {
			return errors.New("invalid checksum target")
		}
		if _, err := hex.DecodeString(hash); err != nil {
			return err
		}
		sums[name] = hash
	}
	for _, asset := range assets {
		want := sums[asset.Name]
		if want == "" {
			return errors.New("missing checksum")
		}
		data, err := os.ReadFile(filepath.Join(downloadDir, asset.Name))
		if err != nil {
			return err
		}
		got := sha256.Sum256(data)
		if hex.EncodeToString(got[:]) != strings.ToLower(want) {
			return errors.New("checksum mismatch")
		}
	}
	return nil
}

func (ops setupRealOps) InstallBinary(src, dest, binaryName, version string) error {
	info, err := os.Lstat(dest)
	if err == nil && info.Mode()&os.ModeSymlink != 0 {
		return errors.New("destination symlink")
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("empty binary")
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, data, 0755); err != nil {
		return err
	}
	if err := os.Chmod(tmp, 0755); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return err
	}
	return ops.VerifyBinary(dest, binaryName, version)
}

func (setupRealOps) VerifyBinary(path, binaryName, version string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "--version")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil || stderr.Len() != 0 {
		return errors.New("version command failed")
	}
	line := stdout.String()
	if !strings.HasPrefix(line, binaryName+" "+version+" go=") || !strings.HasSuffix(line, "\n") || strings.Contains(line, "V.0.0-dev") {
		return errors.New("version mismatch")
	}
	return nil
}

func (setupRealOps) BinaryExists(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0
}

func (setupRealOps) EnsureDir(path string, mode os.FileMode) error {
	if err := os.MkdirAll(path, mode); err != nil {
		return err
	}
	return os.Chmod(path, mode)
}

func (setupRealOps) CreateFileExclusive(path string, data []byte, mode os.FileMode) (bool, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return false, err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		if os.IsExist(err) {
			return false, nil
		}
		return false, err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return false, err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return false, err
	}
	if err := f.Close(); err != nil {
		return false, err
	}
	return true, os.Chmod(path, mode)
}

func (setupRealOps) InstallAdminArchive(archivePath, installDir string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()
	parent := installDir
	tmp, err := os.MkdirTemp(parent, ".admin-next-")
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			os.RemoveAll(tmp)
		}
	}()
	seen := map[string]bool{}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		name := filepath.Clean(hdr.Name)
		if name != "index.html" && name != "adlaire-ci-sdk.js" {
			return errors.New("unsafe admin archive entry")
		}
		if hdr.Typeflag != tar.TypeReg || hdr.FileInfo().Mode()&os.ModeSymlink != 0 || hdr.Size < 0 || hdr.Size > setupAdminUILimitBytes {
			return errors.New("invalid admin archive entry")
		}
		data, err := io.ReadAll(io.LimitReader(tr, setupAdminUILimitBytes+1))
		if err != nil {
			return err
		}
		if int64(len(data)) != hdr.Size || int64(len(data)) > setupAdminUILimitBytes {
			return errors.New("invalid admin archive size")
		}
		if err := os.WriteFile(filepath.Join(tmp, name), data, 0644); err != nil {
			return err
		}
		seen[name] = true
	}
	if !seen["index.html"] || !seen["adlaire-ci-sdk.js"] {
		return errors.New("missing admin files")
	}
	dest := filepath.Join(installDir, "admin")
	old := filepath.Join(installDir, ".admin-old")
	os.RemoveAll(old)
	if _, err := os.Lstat(dest); err == nil {
		if err := os.Rename(dest, old); err != nil {
			return err
		}
	}
	if err := os.Rename(tmp, dest); err != nil {
		if _, oldErr := os.Lstat(old); oldErr == nil {
			os.Rename(old, dest)
		}
		return err
	}
	os.RemoveAll(old)
	ok = true
	return nil
}

func (setupRealOps) WriteRunnerUnits(cfg setupConfig) error {
	service := fmt.Sprintf(`[Unit]
Description=Adlaire CI Runner

[Service]
Type=oneshot
User=%s
WorkingDirectory=%s
ExecStart=%s --state-dir %s
`, cfg.ServiceUser, cfg.InstallDir, filepath.Join(cfg.BinDir, "adlaire-ci-runner"), cfg.InstallDir)
	timer := `[Unit]
Description=Adlaire CI Runner Timer

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min
Unit=adlaire-ci.service

[Install]
WantedBy=timers.target
`
	if err := atomicWriteText(filepath.Join(setupSystemdDir, "adlaire-ci.service"), service, 0644); err != nil {
		return err
	}
	return atomicWriteText(filepath.Join(setupSystemdDir, "adlaire-ci.timer"), timer, 0644)
}

func (setupRealOps) WriteAPIUnit(cfg setupConfig) error {
	unit := fmt.Sprintf(`[Unit]
Description=Adlaire CI API Server
Wants=adlaire-ci.timer
After=network.target adlaire-ci.timer

[Service]
Type=simple
User=%s
WorkingDirectory=%s
ExecStart=%s --addr 127.0.0.1:8765 --state-dir %s
Restart=always
RestartSec=10
KillSignal=SIGTERM
TimeoutStopSec=15s

[Install]
WantedBy=multi-user.target
`, cfg.ServiceUser, cfg.InstallDir, filepath.Join(cfg.BinDir, "adlaire-ci-api"), cfg.InstallDir)
	return atomicWriteText(filepath.Join(setupSystemdDir, "adlaire-ci-api.service"), unit, 0644)
}

func (setupRealOps) InitCredentials(apiPath, installDir, passwordFile string) error {
	file, err := os.Open(passwordFile)
	if err != nil {
		return err
	}
	defer file.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, apiPath, "--init-credentials", "--state-dir", installDir)
	cmd.Stdin = file
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		return errors.New("credentials command produced output")
	}
	return nil
}

func (setupRealOps) NewBackupDir() (string, error) {
	return os.MkdirTemp("", "adlaire-ci-setup-backup-")
}

func (setupRealOps) BackupFile(src, dest string) error {
	return copySetupFile(src, dest, 0755)
}

func (setupRealOps) RestoreFile(src, dest string) error {
	return copySetupFile(src, dest, 0755)
}

func copySetupFile(src, dest string, mode os.FileMode) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("not regular")
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(dest, data, mode); err != nil {
		return err
	}
	return os.Chmod(dest, mode)
}

func (setupRealOps) BackupDir(src, dest string) error {
	return copySetupDir(src, dest)
}

func (setupRealOps) RestoreDir(src, dest string) error {
	os.RemoveAll(dest)
	return copySetupDir(src, dest)
}

func copySetupDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return errors.New("symlink in directory")
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		return copySetupFile(path, target, info.Mode().Perm())
	})
}

func (setupRealOps) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (setupRealOps) Systemctl(args ...string) setupCommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "systemctl", args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}
	return setupCommandResult{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: exitCode}
}

func (setupRealOps) HealthCheck() error {
	client := &http.Client{Timeout: setupHealthTimeout}
	resp, err := client.Get("http://127.0.0.1:8765/api/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("health status")
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, setupHealthLimitBytes+1))
	if err != nil {
		return err
	}
	if len(data) == 0 || len(data) > setupHealthLimitBytes {
		return errors.New("health body")
	}
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil || obj == nil {
		return errors.New("health json")
	}
	return nil
}

func setupSortedAssetNames(assets []setupAsset) []string {
	names := make([]string, 0, len(assets))
	for _, asset := range assets {
		names = append(names, asset.Name)
	}
	sort.Strings(names)
	return names
}
