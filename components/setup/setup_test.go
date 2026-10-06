package setup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type fakeSetupOps struct {
	real setupRealOps

	versions       map[string]string
	downloaded     map[string][]byte
	systemctlCalls []string
	apiEnabled     bool
	unsafeArchive  bool

	failDownloadName   string
	failChecksum       bool
	failAPIRestartOnce bool
	apiRestartFailures int
}

func newFakeSetupOps() *fakeSetupOps {
	return &fakeSetupOps{
		versions:   map[string]string{},
		downloaded: map[string][]byte{},
	}
}

func runSetupForTest(ops *fakeSetupOps, args ...string) (int, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runSetup(args, &stdout, &stderr, ops)
	return code, stdout.String(), stderr.String()
}

func (ops *fakeSetupOps) ValidatePath(path string) error {
	return ops.real.ValidatePath(path)
}

func (ops *fakeSetupOps) ReadSecretInput(path string) (string, error) {
	return ops.real.ReadSecretInput(path)
}

func (ops *fakeSetupOps) ValidateExistingToken(path string) error {
	return ops.real.ValidateExistingToken(path)
}

func (ops *fakeSetupOps) ValidateLastSHA(path string) error {
	return ops.real.ValidateLastSHA(path)
}

func (ops *fakeSetupOps) ValidateCredentials(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0600 {
		return errors.New("invalid credentials")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if string(data) != "credential\n" {
		return errors.New("invalid credentials")
	}
	return nil
}

func (ops *fakeSetupOps) EnsureDownloadDir(path, version string) error {
	return ops.real.EnsureDownloadDir(path, version)
}

func (ops *fakeSetupOps) DownloadAsset(cfg setupConfig, asset setupAsset) error {
	if ops.failDownloadName == asset.Name {
		return errors.New("download failed")
	}
	if err := os.MkdirAll(cfg.DownloadDir, 0755); err != nil {
		return err
	}
	var data []byte
	if asset.Name == "SHA256SUMS" {
		data = ops.checksumBytes()
	} else if asset.Name == "admin-ui.tar.gz" {
		if ops.unsafeArchive {
			data = setupTestTarGZ(tarFile{name: "../index.html", body: "bad"})
		} else {
			data = setupTestTarGZ(
				tarFile{name: "index.html", body: "<!doctype html><title>Adlaire CI</title>"},
				tarFile{name: "adlaire-ci-sdk.js", body: "globalThis.AdlaireCI = {};"},
			)
		}
	} else {
		data = []byte("binary " + asset.Name + " " + cfg.TargetVersion + "\n")
	}
	ops.downloaded[asset.Name] = data
	return os.WriteFile(filepath.Join(cfg.DownloadDir, asset.Name), data, 0644)
}

func (ops *fakeSetupOps) checksumBytes() []byte {
	names := make([]string, 0, len(ops.downloaded))
	for name := range ops.downloaded {
		names = append(names, name)
	}
	sort.Strings(names)
	var out strings.Builder
	for _, name := range names {
		sum := sha256.Sum256(ops.downloaded[name])
		fmt.Fprintf(&out, "%s  %s\n", hex.EncodeToString(sum[:]), name)
	}
	return []byte(out.String())
}

func (ops *fakeSetupOps) VerifyChecksums(downloadDir string, assets []setupAsset, osArch string) error {
	if ops.failChecksum {
		return errors.New("checksum failed")
	}
	return ops.real.VerifyChecksums(downloadDir, assets, osArch)
}

func (ops *fakeSetupOps) InstallBinary(src, dest, binaryName, version string) error {
	if err := ops.real.InstallBinaryForTest(src, dest); err != nil {
		return err
	}
	ops.versions[dest] = version
	return nil
}

func (ops *fakeSetupOps) VerifyBinary(path, binaryName, version string) error {
	if !ops.BinaryExists(path) {
		return errors.New("missing binary")
	}
	if ops.versions[path] != version {
		return fmt.Errorf("version mismatch: got %q want %q", ops.versions[path], version)
	}
	return nil
}

func (ops *fakeSetupOps) BinaryExists(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0
}

func (ops *fakeSetupOps) EnsureDir(path string, mode os.FileMode) error {
	return ops.real.EnsureDir(path, mode)
}

func (ops *fakeSetupOps) CreateFileExclusive(path string, data []byte, mode os.FileMode) (bool, error) {
	return ops.real.CreateFileExclusive(path, data, mode)
}

func (ops *fakeSetupOps) InstallAdminArchive(archivePath, installDir string) error {
	return ops.real.InstallAdminArchive(archivePath, installDir)
}

func (ops *fakeSetupOps) WriteRunnerUnits(cfg setupConfig) error {
	return nil
}

func (ops *fakeSetupOps) WriteAPIUnit(cfg setupConfig) error {
	return nil
}

func (ops *fakeSetupOps) InitCredentials(apiPath, installDir, passwordFile string) error {
	return os.WriteFile(filepath.Join(installDir, ".admin_credentials"), []byte("credential\n"), 0600)
}

func (ops *fakeSetupOps) NewBackupDir() (string, error) {
	return os.MkdirTemp("", "setup-test-backup-")
}

func (ops *fakeSetupOps) BackupFile(src, dest string) error {
	if err := ops.real.BackupFile(src, dest); err != nil {
		return err
	}
	ops.versions[dest] = ops.versions[src]
	return nil
}

func (ops *fakeSetupOps) RestoreFile(src, dest string) error {
	if err := ops.real.RestoreFile(src, dest); err != nil {
		return err
	}
	ops.versions[dest] = ops.versions[src]
	return nil
}

func (ops *fakeSetupOps) BackupDir(src, dest string) error {
	return ops.real.BackupDir(src, dest)
}

func (ops *fakeSetupOps) RestoreDir(src, dest string) error {
	return ops.real.RestoreDir(src, dest)
}

func (ops *fakeSetupOps) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (ops *fakeSetupOps) Systemctl(args ...string) setupCommandResult {
	call := strings.Join(args, " ")
	ops.systemctlCalls = append(ops.systemctlCalls, call)
	switch call {
	case "is-enabled adlaire-ci-api":
		if ops.apiEnabled {
			return setupCommandResult{Stdout: "enabled\n", ExitCode: 0}
		}
		return setupCommandResult{Stderr: "disabled\n", ExitCode: 1}
	case "is-active adlaire-ci.timer", "is-active adlaire-ci-api":
		return setupCommandResult{Stdout: "active\n", ExitCode: 0}
	case "cat adlaire-ci.service", "cat adlaire-ci.timer", "cat adlaire-ci-api.service":
		return setupCommandResult{Stdout: "unit\n", ExitCode: 0}
	case "restart adlaire-ci-api":
		if ops.failAPIRestartOnce && ops.apiRestartFailures == 0 {
			ops.apiRestartFailures++
			return setupCommandResult{Stderr: "failed\n", ExitCode: 1}
		}
		return setupCommandResult{ExitCode: 0}
	default:
		return setupCommandResult{ExitCode: 0}
	}
}

func (ops *fakeSetupOps) HealthCheck() error {
	return nil
}

func (setupRealOps) InstallBinaryForTest(src, dest string) error {
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
	if err := os.WriteFile(dest, data, 0755); err != nil {
		return err
	}
	return os.Chmod(dest, 0755)
}

type tarFile struct {
	name string
	body string
}

func setupTestTarGZ(files ...tarFile) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, file := range files {
		body := []byte(file.body)
		_ = tw.WriteHeader(&tar.Header{Name: file.name, Mode: 0644, Size: int64(len(body)), Typeflag: tar.TypeReg})
		_, _ = tw.Write(body)
	}
	_ = tw.Close()
	_ = gz.Close()
	return buf.Bytes()
}

func writeSetupSecret(t *testing.T, dir, name, value string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(value+"\n"), 0600); err != nil {
		t.Fatalf("write secret: %v", err)
	}
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatalf("chmod secret: %v", err)
	}
	return path
}

func seedSetupBinary(t *testing.T, ops *fakeSetupOps, dir, name, version string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	if err := os.WriteFile(path, []byte("old "+name+" "+version+"\n"), 0755); err != nil {
		t.Fatalf("write binary: %v", err)
	}
	if err := os.Chmod(path, 0755); err != nil {
		t.Fatalf("chmod binary: %v", err)
	}
	ops.versions[path] = version
}

func TestRunSetupHelpVersionPrecedence(t *testing.T) {
	ops := newFakeSetupOps()
	code, stdout, stderr := runSetupForTest(ops, "--help", "bad\narg")
	if code != 0 || stdout != "Usage: adlaire-ci-setup <install|install-api|install-obsidian|update> --target-version V.X.N [options] [--version] [--help]\n" || stderr != "" {
		t.Fatalf("help mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	code, stdout, stderr = runSetupForTest(ops, "--version", "install", "--target-version", "bad")
	if code != 0 || !strings.HasPrefix(stdout, "adlaire-ci-setup V.0.0-dev go=") || stderr != "" {
		t.Fatalf("version mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestRunSetupInstallSuccess(t *testing.T) {
	root := t.TempDir()
	installDir := filepath.Join(root, "install")
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	tokenFile := writeSetupSecret(t, root, "github_token", "github-token")
	ops := newFakeSetupOps()

	code, stdout, stderr := runSetupForTest(ops,
		"install",
		"--target-version", "V.1.1",
		"--github-token-file", tokenFile,
		"--install-dir", installDir,
		"--bin-dir", binDir,
		"--download-dir", downloadDir,
	)
	if code != 0 || stderr != "" {
		t.Fatalf("install failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	wantStdout := strings.Join([]string{
		"setup: start install validate",
		"setup: ok install validate",
		"setup: start install download",
		"setup: ok install download",
		"setup: start install checksum",
		"setup: ok install checksum",
		"setup: start install install-binaries",
		"setup: ok install install-binaries",
		"setup: start install initialize-secrets",
		"setup: ok install initialize-secrets",
		"setup: start install initialize-state",
		"setup: ok install initialize-state",
		"setup: start install configure-systemd",
		"setup: ok install configure-systemd",
		"setup: start install activate-services",
		"setup: ok install activate-services",
		"setup: start install verify",
		"setup: ok install verify",
		"setup: success install V.1.1",
		"",
	}, "\n")
	if stdout != wantStdout {
		t.Fatalf("stdout mismatch:\n got %q\nwant %q", stdout, wantStdout)
	}
	for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup"} {
		if ops.versions[filepath.Join(binDir, name)] != "V.1.1" {
			t.Fatalf("missing installed version for %s", name)
		}
	}
	if data, err := os.ReadFile(filepath.Join(installDir, ".github_token")); err != nil || string(data) != "github-token\n" {
		t.Fatalf("token mismatch: data=%q err=%v", string(data), err)
	}
	if data, err := os.ReadFile(filepath.Join(installDir, ".last_sha")); err != nil || string(data) != "{\"sha\":\"\"}\n" {
		t.Fatalf("last sha mismatch: data=%q err=%v", string(data), err)
	}
	if !containsSetupCall(ops.systemctlCalls, "enable --now adlaire-ci.timer") {
		t.Fatalf("missing timer enable call: %v", ops.systemctlCalls)
	}
}

func TestRunSetupInstallObsidianSuccess(t *testing.T) {
	root := t.TempDir()
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	ops := newFakeSetupOps()

	code, stdout, stderr := runSetupForTest(ops,
		"install-obsidian",
		"--target-version", "V.1.1",
		"--bin-dir", binDir,
		"--download-dir", downloadDir,
	)
	if code != 0 || stderr != "" {
		t.Fatalf("install-obsidian failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	wantStdout := strings.Join([]string{
		"setup: start install-obsidian validate",
		"setup: ok install-obsidian validate",
		"setup: start install-obsidian download",
		"setup: ok install-obsidian download",
		"setup: start install-obsidian checksum",
		"setup: ok install-obsidian checksum",
		"setup: start install-obsidian install-binaries",
		"setup: ok install-obsidian install-binaries",
		"setup: start install-obsidian verify",
		"setup: ok install-obsidian verify",
		"setup: success install-obsidian V.1.1",
		"",
	}, "\n")
	if stdout != wantStdout {
		t.Fatalf("stdout mismatch:\n got %q\nwant %q", stdout, wantStdout)
	}
	if ops.versions[filepath.Join(binDir, "adlaire-ci-obsidian")] != "V.1.1" {
		t.Fatalf("missing installed obsidian version")
	}
	for _, forbidden := range []string{"daemon-reload", "enable --now adlaire-ci.timer", "is-active adlaire-ci.timer"} {
		if containsSetupCall(ops.systemctlCalls, forbidden) {
			t.Fatalf("install-obsidian must not call systemd %q: %v", forbidden, ops.systemctlCalls)
		}
	}
}

func TestRunSetupFailureContracts(t *testing.T) {
	root := t.TempDir()
	installDir := filepath.Join(root, "install")
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	tokenFile := writeSetupSecret(t, root, "github_token", "github-token")

	t.Run("install api precondition", func(t *testing.T) {
		ops := newFakeSetupOps()
		passFile := writeSetupSecret(t, root, "admin_password_precondition", "Initial-password-123")
		code, stdout, stderr := runSetupForTest(ops,
			"install-api",
			"--target-version", "V.1.1",
			"--initial-password-file", passFile,
			"--install-dir", installDir,
			"--bin-dir", binDir,
			"--download-dir", downloadDir,
		)
		if code != 2 || stdout != "setup: start install-api validate\n" || stderr != "setup: error PRECONDITION_FAILED stage=validate rollback=none\n" {
			t.Fatalf("precondition mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})

	t.Run("download", func(t *testing.T) {
		ops := newFakeSetupOps()
		ops.failDownloadName = "adlaire-ci-runner-linux-amd64"
		code, stdout, stderr := runSetupForTest(ops,
			"install",
			"--target-version", "V.1.1",
			"--github-token-file", tokenFile,
			"--install-dir", filepath.Join(root, "install-download"),
			"--bin-dir", filepath.Join(root, "bin-download"),
			"--download-dir", filepath.Join(root, "download-download"),
		)
		if code != 3 || !strings.HasSuffix(stdout, "setup: start install download\n") || stderr != "setup: error DOWNLOAD_FAILED stage=download rollback=none\n" {
			t.Fatalf("download mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})

	t.Run("checksum", func(t *testing.T) {
		ops := newFakeSetupOps()
		ops.failChecksum = true
		code, stdout, stderr := runSetupForTest(ops,
			"install",
			"--target-version", "V.1.1",
			"--github-token-file", tokenFile,
			"--install-dir", filepath.Join(root, "install-checksum"),
			"--bin-dir", filepath.Join(root, "bin-checksum"),
			"--download-dir", filepath.Join(root, "download-checksum"),
		)
		if code != 1 || !strings.HasSuffix(stdout, "setup: start install checksum\n") || stderr != "setup: error CHECKSUM_FAILED stage=checksum rollback=none\n" {
			t.Fatalf("checksum mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})
}

func TestRunSetupInstallAPIUnsafeAdminArchivePreservesExistingAdmin(t *testing.T) {
	root := t.TempDir()
	installDir := filepath.Join(root, "install")
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	ops := newFakeSetupOps()
	for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup"} {
		seedSetupBinary(t, ops, binDir, name, "V.1.1")
	}
	adminDir := filepath.Join(installDir, "admin")
	if err := os.MkdirAll(adminDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "index.html"), []byte("old-index"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "adlaire-ci-sdk.js"), []byte("old-sdk"), 0644); err != nil {
		t.Fatal(err)
	}
	passFile := writeSetupSecret(t, root, "admin_password", "Initial-password-123")
	ops.unsafeArchive = true

	code, stdout, stderr := runSetupForTest(ops,
		"install-api",
		"--target-version", "V.1.1",
		"--initial-password-file", passFile,
		"--install-dir", installDir,
		"--bin-dir", binDir,
		"--download-dir", downloadDir,
	)
	if code != 1 || !strings.HasSuffix(stdout, "setup: start install-api install-admin\n") || stderr != "setup: error ADMIN_INSTALL_FAILED stage=install-admin rollback=none\n" {
		t.Fatalf("unsafe admin mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if data, err := os.ReadFile(filepath.Join(adminDir, "index.html")); err != nil || string(data) != "old-index" {
		t.Fatalf("admin index was not preserved: data=%q err=%v", string(data), err)
	}
}

func TestRunSetupUpdateAPIRollback(t *testing.T) {
	root := t.TempDir()
	installDir := filepath.Join(root, "install")
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	ops := newFakeSetupOps()
	ops.apiEnabled = true
	ops.failAPIRestartOnce = true
	for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup", "adlaire-ci-api", "adlaire-ci-admin"} {
		seedSetupBinary(t, ops, binDir, name, "V.1.0")
	}
	adminDir := filepath.Join(installDir, "admin")
	if err := os.MkdirAll(adminDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "index.html"), []byte("old-index"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "adlaire-ci-sdk.js"), []byte("old-sdk"), 0644); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := runSetupForTest(ops,
		"update",
		"--target-version", "V.1.1",
		"--install-dir", installDir,
		"--bin-dir", binDir,
		"--download-dir", downloadDir,
	)
	if code != 1 || !strings.HasSuffix(stdout, "setup: start update restart-services\n") || stderr != "setup: error SYSTEMD_FAILED stage=restart-services rollback=completed\n" {
		t.Fatalf("rollback mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup", "adlaire-ci-api", "adlaire-ci-admin"} {
		if got := ops.versions[filepath.Join(binDir, name)]; got != "V.1.0" {
			t.Fatalf("binary %s not rolled back: %s", name, got)
		}
	}
	if data, err := os.ReadFile(filepath.Join(adminDir, "index.html")); err != nil || string(data) != "old-index" {
		t.Fatalf("admin index not rolled back: data=%q err=%v", string(data), err)
	}
}

func TestPhase18MaterializeAdminAssets(t *testing.T) {
	originalIndex := append([]byte(nil), setupEmbeddedAdminIndex...)
	originalSDK := append([]byte(nil), setupEmbeddedAdminSDK...)
	t.Cleanup(func() { SetEmbeddedAdminAssets(originalIndex, originalSDK) })
	SetEmbeddedAdminAssets([]byte("<main>phase18</main>\n"), []byte("export const phase18 = true;\n"))

	installDir := filepath.Join(t.TempDir(), "runtime")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunSetup([]string{"phase18-materialize-admin", "--install-dir", installDir}, &stdout, &stderr)
	if code != 0 || stdout.String() != "setup: success phase18-materialize-admin\n" || stderr.Len() != 0 {
		t.Fatalf("materialize result code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for name, want := range map[string]string{
		"index.html":        "<main>phase18</main>\n",
		"adlaire-ci-sdk.js": "export const phase18 = true;\n",
	} {
		path := filepath.Join(installDir, "admin", name)
		data, err := os.ReadFile(path)
		if err != nil || string(data) != want {
			t.Fatalf("asset %s mismatch data=%q err=%v", name, string(data), err)
		}
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != 0644 {
			t.Fatalf("asset %s mode mismatch info=%v err=%v", name, info, err)
		}
	}
}

func TestPhase18MaterializeAdminRejectsUnsafeTarget(t *testing.T) {
	originalIndex := append([]byte(nil), setupEmbeddedAdminIndex...)
	originalSDK := append([]byte(nil), setupEmbeddedAdminSDK...)
	t.Cleanup(func() { SetEmbeddedAdminAssets(originalIndex, originalSDK) })
	SetEmbeddedAdminAssets([]byte("index"), []byte("sdk"))

	root := t.TempDir()
	target := filepath.Join(root, "target")
	if err := os.Symlink(filepath.Join(root, "missing"), target); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunSetup([]string{"phase18-materialize-admin", "--install-dir", target}, &stdout, &stderr)
	if code != 1 || stdout.Len() != 0 || stderr.String() != "admin asset materialization failed\n" {
		t.Fatalf("unsafe target result code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestSetupAtomicWriteTextSyncsFileAndParent(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "asset.txt")
	ops := defaultSetupAtomicWriteOps()
	syncTargets := []string{}
	ops.sync = func(file *os.File) error {
		syncTargets = append(syncTargets, file.Name())
		return file.Sync()
	}
	if err := atomicWriteTextWithOps(target, "new-value", 0644, ops); err != nil {
		t.Fatalf("atomic write failed: %v", err)
	}
	if len(syncTargets) != 2 || filepath.Dir(syncTargets[0]) != dir || syncTargets[1] != dir {
		t.Fatalf("file and parent directory must be synced in order: %#v", syncTargets)
	}
	data, err := os.ReadFile(target)
	if err != nil || string(data) != "new-value" {
		t.Fatalf("unexpected target data=%q err=%v", string(data), err)
	}
	info, err := os.Stat(target)
	if err != nil || info.Mode().Perm() != 0644 {
		t.Fatalf("unexpected target mode info=%v err=%v", info, err)
	}
}

func TestSetupAtomicWriteTextFailureKeepsExistingTarget(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*setupAtomicWriteOps)
	}{
		{
			name: "file-sync",
			mutate: func(ops *setupAtomicWriteOps) {
				ops.sync = func(*os.File) error { return errors.New("sync failed") }
			},
		},
		{
			name: "rename",
			mutate: func(ops *setupAtomicWriteOps) {
				ops.rename = func(string, string) error { return errors.New("rename failed") }
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			target := filepath.Join(dir, "asset.txt")
			if err := os.WriteFile(target, []byte("old-value"), 0644); err != nil {
				t.Fatal(err)
			}
			ops := defaultSetupAtomicWriteOps()
			test.mutate(&ops)
			if err := atomicWriteTextWithOps(target, "new-value", 0644, ops); err == nil {
				t.Fatal("expected atomic write failure")
			}
			data, err := os.ReadFile(target)
			if err != nil || string(data) != "old-value" {
				t.Fatalf("existing target changed data=%q err=%v", string(data), err)
			}
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), ".asset.txt.tmp-") {
					t.Fatalf("temporary file was not removed: %s", entry.Name())
				}
			}
		})
	}
}

func TestPhase18MaterializeAdminRejectsAssetSymlink(t *testing.T) {
	originalIndex := append([]byte(nil), setupEmbeddedAdminIndex...)
	originalSDK := append([]byte(nil), setupEmbeddedAdminSDK...)
	t.Cleanup(func() { SetEmbeddedAdminAssets(originalIndex, originalSDK) })
	SetEmbeddedAdminAssets([]byte("index"), []byte("sdk"))

	installDir := filepath.Join(t.TempDir(), "runtime")
	adminDir := filepath.Join(installDir, "admin")
	if err := os.MkdirAll(adminDir, 0755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("unchanged"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(adminDir, "index.html")); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunSetup([]string{"phase18-materialize-admin", "--install-dir", installDir}, &stdout, &stderr)
	if code != 1 || stdout.Len() != 0 || stderr.String() != "admin asset materialization failed\n" {
		t.Fatalf("asset symlink result code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	data, err := os.ReadFile(outside)
	if err != nil || string(data) != "unchanged" {
		t.Fatalf("symlink target changed data=%q err=%v", string(data), err)
	}
}

func TestSetupInitCredentialsAcceptsSpecifiedChildOutput(t *testing.T) {
	root := t.TempDir()
	apiPath := filepath.Join(root, "api")
	passwordPath := filepath.Join(root, "password")
	if err := os.WriteFile(apiPath, []byte("#!/bin/sh\nprintf 'credentials initialized\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(passwordPath, []byte("password123\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := (setupRealOps{}).InitCredentials(apiPath, root, passwordPath); err != nil {
		t.Fatalf("specified init output rejected: %v", err)
	}
	if err := os.WriteFile(apiPath, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := (setupRealOps{}).InitCredentials(apiPath, root, passwordPath); err == nil || err.Error() != "credentials command output mismatch" {
		t.Fatalf("unexpected mismatch result: %v", err)
	}
}

func TestSetupFixturesExist(t *testing.T) {
	fixtures := []string{
		"success-setup-admin-release-asset-layout",
		"failure-setup-download-boundary",
		"failure-setup-api-version-cohort",
		"security-setup-admin-archive-boundary",
		"partial-setup-systemd-rollback-boundary",
		"partial-setup-api-runner-dispatch",
		"security-setup-secret-preservation",
	}
	required := []string{
		"manifest.json",
		"input/cli.json",
		"expected/stdout.txt",
		"expected/stderr.txt",
		"expected/effects.json",
		"expected/security.json",
	}
	for _, fixture := range fixtures {
		root := filepath.Join("..", "..", "testdata", "setup", fixture)
		for _, rel := range required {
			path := filepath.Join(root, rel)
			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("missing fixture file %s: %v", path, err)
			}
			if info.IsDir() {
				t.Fatalf("fixture path is a directory: %s", path)
			}
			if strings.HasSuffix(rel, ".json") {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read fixture json %s: %v", path, err)
				}
				var decoded any
				if err := json.Unmarshal(data, &decoded); err != nil {
					t.Fatalf("invalid fixture json %s: %v", path, err)
				}
			}
		}
	}
}

func containsSetupCall(calls []string, want string) bool {
	for _, call := range calls {
		if call == want {
			return true
		}
	}
	return false
}
