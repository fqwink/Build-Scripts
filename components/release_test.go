package components

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeReleaseOps struct {
	root      string
	notes     string
	token     string
	dirty     bool
	tagBad    bool
	buildErr  error
	outputErr error
	wroteOut  bool
	github    *fakeReleaseGitHub
}

type fakeReleaseGitHub struct {
	checkRemoteCalls int
	releaseExists    bool
	createErr        error
	uploadErrAt      int
	verifyErr        error
	publishErr       error
	deleteErr        error
	uploads          []string
	deletes          int
}

func newFakeReleaseOps(t *testing.T) *fakeReleaseOps {
	t.Helper()
	root := t.TempDir()
	notes := filepath.Join(root, "notes.md")
	tokenDir := t.TempDir()
	token := filepath.Join(tokenDir, "token")
	if err := os.WriteFile(notes, []byte("Release notes\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(token, []byte("ghp_fake_token\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return &fakeReleaseOps{root: root, notes: notes, token: token, github: &fakeReleaseGitHub{}}
}

func runReleaseForTest(ops *fakeReleaseOps, args ...string) (int, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runRelease(args, &stdout, &stderr, ops)
	return code, stdout.String(), stderr.String()
}

func releaseArgs(ops *fakeReleaseOps, out string) []string {
	return []string{
		"--repository", "fqwink/Build-Scripts",
		"--tag", "V.1.1",
		"--commit", strings.Repeat("a", 40),
		"--notes-file", ops.notes,
		"--token-file", ops.token,
		"--out", out,
	}
}

func (ops *fakeReleaseOps) CheckoutRoot() (string, error) {
	return ops.root, nil
}

func (ops *fakeReleaseOps) VerifyClean(root string) error {
	if ops.dirty {
		return errors.New("dirty")
	}
	return nil
}

func (ops *fakeReleaseOps) VerifyCommitAndTag(root, commit, tag string) error {
	if ops.tagBad {
		return errors.New("tag mismatch")
	}
	return nil
}

func (ops *fakeReleaseOps) ReadNotes(path, checkoutRoot string) (string, error) {
	if path != ops.notes || checkoutRoot != ops.root {
		return "", errors.New("bad notes")
	}
	return "Release notes\n", nil
}

func (ops *fakeReleaseOps) ReadToken(path string) (string, error) {
	if path != ops.token {
		return "", errors.New("bad token")
	}
	return "ghp_fake_token", nil
}

func (ops *fakeReleaseOps) PrepareOutput(path, checkoutRoot string) error {
	if path == "" || checkoutRoot != ops.root {
		return errors.New("bad output")
	}
	return nil
}

func (ops *fakeReleaseOps) BuildArtifacts(cfg releaseConfig, state releaseState) ([]releaseAsset, error) {
	if ops.buildErr != nil {
		return nil, ops.buildErr
	}
	return releaseTestAssets(cfg.Tag), nil
}

func (ops *fakeReleaseOps) WriteOutput(path string, assets []releaseAsset) error {
	if ops.outputErr != nil {
		return ops.outputErr
	}
	if err := verifyReleaseAssetSet(assets); err != nil {
		return err
	}
	ops.wroteOut = true
	return nil
}

func (ops *fakeReleaseOps) GitHub() releaseGitHubOps {
	return ops.github
}

func (gh *fakeReleaseGitHub) CheckRemote(cfg releaseConfig, token string) (string, error) {
	gh.checkRemoteCalls++
	if cfg.Commit == strings.Repeat("b", 40) {
		return "", errReleaseTagMismatch
	}
	return "main", nil
}

func (gh *fakeReleaseGitHub) CheckReleaseAbsent(cfg releaseConfig, token string) error {
	if gh.releaseExists {
		return errReleaseExists
	}
	return nil
}

func (gh *fakeReleaseGitHub) CreateDraft(cfg releaseConfig, token, notes string) (releaseGitHubRelease, error) {
	if gh.createErr != nil {
		return releaseGitHubRelease{}, gh.createErr
	}
	return releaseGitHubRelease{ID: 42, UploadURL: "https://uploads.github.com/repos/fqwink/Build-Scripts/releases/42/assets{?name,label}", HTMLURL: "https://github.com/fqwink/Build-Scripts/releases/tag/" + cfg.Tag, Draft: true}, nil
}

func (gh *fakeReleaseGitHub) UploadAsset(rel releaseGitHubRelease, token string, asset releaseAsset, repository string) error {
	if gh.uploadErrAt > 0 && len(gh.uploads)+1 == gh.uploadErrAt {
		return errors.New("upload")
	}
	gh.uploads = append(gh.uploads, asset.Name)
	return nil
}

func (gh *fakeReleaseGitHub) VerifyAssets(rel releaseGitHubRelease, cfg releaseConfig, assets []releaseAsset, token string) error {
	return gh.verifyErr
}

func (gh *fakeReleaseGitHub) Publish(rel releaseGitHubRelease, cfg releaseConfig, token, notes string, assets []releaseAsset) error {
	return gh.publishErr
}

func (gh *fakeReleaseGitHub) DeleteDraft(rel releaseGitHubRelease, cfg releaseConfig, token string) error {
	gh.deletes++
	return gh.deleteErr
}

func releaseTestAssets(tag string) []releaseAsset {
	assets := []releaseAsset{}
	for _, name := range releaseAssetNames()[:7] {
		mode := os.FileMode(0755)
		if name == "admin-ui.tar.gz" {
			mode = 0644
		}
		assets = append(assets, releaseAsset{Name: name, Data: []byte(name + " " + tag + "\n"), Mode: mode})
	}
	return withChecksumAsset(assets)
}

func TestRunReleaseHelpAndVersion(t *testing.T) {
	t.Parallel()

	ops := newFakeReleaseOps(t)
	code, stdout, stderr := runReleaseForTest(ops, "--help", "--bad")
	if code != 0 || !strings.Contains(stdout, "Usage: adlaire-ci-release") || stderr != "" {
		t.Fatalf("unexpected help result code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	code, stdout, stderr = runReleaseForTest(ops, "--version", "--repository")
	if code != 0 || !strings.HasPrefix(stdout, "adlaire-ci-release V.0.0-dev go=") || stderr != "" {
		t.Fatalf("unexpected version result code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestRunReleaseSuccessPublishesAssets(t *testing.T) {
	t.Parallel()

	ops := newFakeReleaseOps(t)
	code, stdout, stderr := runReleaseForTest(ops, releaseArgs(ops, filepath.Join(t.TempDir(), "release-out"))...)
	if code != 0 || stderr != "" {
		t.Fatalf("expected success, code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !ops.wroteOut {
		t.Fatal("expected output artifacts to be written")
	}
	if len(ops.github.uploads) != 8 {
		t.Fatalf("expected 8 uploads, got %v", ops.github.uploads)
	}
	if ops.github.deletes != 0 {
		t.Fatalf("expected no cleanup on success")
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(stdout), &response); err != nil {
		t.Fatalf("stdout is not json: %v", err)
	}
	if response["tag"] != "V.1.1" || response["commit"] != strings.Repeat("a", 40) || response["published"] != true {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestRunReleaseRejectsDirtyWorktreeBeforeBuild(t *testing.T) {
	t.Parallel()

	ops := newFakeReleaseOps(t)
	ops.dirty = true
	code, stdout, stderr := runReleaseForTest(ops, releaseArgs(ops, filepath.Join(t.TempDir(), "release-out"))...)
	if code != 2 || stdout != "" || stderr != "release: DIRTY_WORKTREE\n" {
		t.Fatalf("unexpected dirty result code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if ops.wroteOut || len(ops.github.uploads) != 0 {
		t.Fatal("dirty worktree must not create artifacts or upload assets")
	}
}

func TestRunReleaseVersionMismatchStopsBeforeGitHubWrite(t *testing.T) {
	t.Parallel()

	ops := newFakeReleaseOps(t)
	ops.buildErr = errReleaseVersionMismatch
	code, stdout, stderr := runReleaseForTest(ops, releaseArgs(ops, filepath.Join(t.TempDir(), "release-out"))...)
	if code != 1 || stdout != "" || stderr != "release: VERSION_MISMATCH\n" {
		t.Fatalf("unexpected version mismatch code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if len(ops.github.uploads) != 0 {
		t.Fatal("version mismatch must not upload assets")
	}
}

func TestRunReleaseUploadFailureCleansDraft(t *testing.T) {
	t.Parallel()

	ops := newFakeReleaseOps(t)
	ops.github.uploadErrAt = 3
	code, stdout, stderr := runReleaseForTest(ops, releaseArgs(ops, filepath.Join(t.TempDir(), "release-out"))...)
	if code != 3 || stdout != "" || stderr != "release: ASSET_UPLOAD_FAILED\n" {
		t.Fatalf("unexpected upload failure code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if ops.github.deletes != 1 {
		t.Fatalf("expected one draft cleanup, got %d", ops.github.deletes)
	}
}

func TestRunReleaseCleanupFailureOverridesOriginal(t *testing.T) {
	t.Parallel()

	ops := newFakeReleaseOps(t)
	ops.github.uploadErrAt = 1
	ops.github.deleteErr = errors.New("cleanup")
	code, stdout, stderr := runReleaseForTest(ops, releaseArgs(ops, filepath.Join(t.TempDir(), "release-out"))...)
	if code != 1 || stdout != "" || stderr != "release: DRAFT_CLEANUP_FAILED\n" {
		t.Fatalf("unexpected cleanup failure code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestReleaseChecksumManifest(t *testing.T) {
	t.Parallel()

	assets := releaseTestAssets("V.1.1")
	if err := verifyReleaseAssetSet(assets); err != nil {
		t.Fatalf("expected release asset set to verify: %v", err)
	}
	byName := map[string]releaseAsset{}
	for _, asset := range assets {
		byName[asset.Name] = asset
	}
	lines := strings.Split(strings.TrimSuffix(string(byName["SHA256SUMS"].Data), "\n"), "\n")
	if len(lines) != 7 {
		t.Fatalf("expected 7 checksum lines, got %d", len(lines))
	}
	wantNames := releaseAssetNames()[:7]
	sortStrings(wantNames)
	for i, line := range lines {
		fields := strings.Split(line, "  ")
		if len(fields) != 2 || fields[1] != wantNames[i] {
			t.Fatalf("line %d mismatch: %q", i, line)
		}
		sum := sha256.Sum256(byName[fields[1]].Data)
		if fields[0] != hex.EncodeToString(sum[:]) {
			t.Fatalf("line %d checksum mismatch", i)
		}
	}
}

func TestReleaseFixtureFilesExist(t *testing.T) {
	t.Parallel()

	names := []string{
		"success-release-assets-reproducible",
		"failure-release-output-parent-race",
		"success-release-draft-publish",
		"failure-release-dirty-worktree",
		"failure-release-version-mismatch",
		"failure-release-source-snapshot-boundary",
		"failure-release-source-race",
		"failure-release-notes-identity",
		"failure-release-token-identity",
		"failure-release-go-environment",
		"failure-release-gofmt-order",
		"failure-release-permission",
		"partial-release-remote-race-cleanup",
		"failure-release-non-reproducible",
		"failure-release-output-atomicity",
		"failure-release-existing-release",
		"partial-release-upload-cleanup",
		"partial-release-cleanup-failure",
		"failure-release-timeout-boundaries",
		"security-release-http-boundary",
		"security-release-token-mask",
	}
	for _, name := range names {
		root := filepath.Join("..", "testdata", "release", name)
		for _, rel := range []string{
			"manifest.json",
			"input/fakes.json",
			"expected/stdout.txt",
			"expected/stderr.txt",
			"expected/effects.json",
			"expected/release-assets.json",
			"expected/security.json",
		} {
			path := filepath.Join(root, rel)
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("missing release fixture file %s: %v", path, err)
			}
			if strings.HasSuffix(rel, ".json") {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				var v any
				if err := json.Unmarshal(data, &v); err != nil {
					t.Fatalf("invalid json %s: %v", path, err)
				}
			}
		}
	}
}

func sortStrings(values []string) {
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] < values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
}
