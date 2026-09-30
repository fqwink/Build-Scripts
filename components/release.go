package components

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
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	releaseBinaryName   = "adlaire-ci-release"
	releaseHTTPTimeout  = 30 * time.Second
	releaseAssetTimeout = 5 * time.Minute
	releaseMaxOutput    = 1024 * 1024
	releaseMaxArchive   = 128 * 1024 * 1024
	releaseMaxNotes     = 262144
	releaseMaxToken     = 4096
)

var (
	releaseBinaryVersion     = "V.0.0-dev"
	releaseRepositoryPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,100}/[A-Za-z0-9_.-]{1,100}$`)
	releaseTagPattern        = regexp.MustCompile(`^V\.[1-9][0-9]*\.[0-9]+$`)
)

type releaseConfig struct {
	Repository string
	Tag        string
	Commit     string
	NotesFile  string
	TokenFile  string
	Out        string
}

type releaseState struct {
	CheckoutRoot string
	Notes        string
	Token        string
	Assets       []releaseAsset
	ReleaseID    int64
	ReleaseURL   string
}

type releaseAsset struct {
	Name   string
	Data   []byte
	Mode   os.FileMode
	SHA256 string
}

type releaseCommandResult struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}

type releaseGitHubRelease struct {
	ID              int64
	UploadURL       string
	HTMLURL         string
	TagName         string
	TargetCommitish string
	Name            string
	Body            string
	Draft           bool
	Prerelease      bool
}

type releaseGitHubOps interface {
	CheckRemote(cfg releaseConfig, token string) (string, error)
	CheckReleaseAbsent(cfg releaseConfig, token string) error
	CreateDraft(cfg releaseConfig, token, notes string) (releaseGitHubRelease, error)
	UploadAsset(releaseGitHubRelease, string, releaseAsset, string) error
	VerifyAssets(releaseGitHubRelease, releaseConfig, []releaseAsset, string) error
	Publish(releaseGitHubRelease, releaseConfig, string, string, []releaseAsset) error
	DeleteDraft(releaseGitHubRelease, releaseConfig, string) error
}

type releaseOperations interface {
	CheckoutRoot() (string, error)
	VerifyClean(root string) error
	VerifyCommitAndTag(root, commit, tag string) error
	ReadNotes(path, checkoutRoot string) (string, error)
	ReadToken(path string) (string, error)
	PrepareOutput(path, checkoutRoot string) error
	BuildArtifacts(cfg releaseConfig, state releaseState) ([]releaseAsset, error)
	WriteOutput(path string, assets []releaseAsset) error
	GitHub() releaseGitHubOps
}

type releaseRealOps struct{}

type releaseRealGitHubOps struct{}

func SetBinaryVersion(version string) {
	if version == "" {
		return
	}
	builderDefaultVersion = version
	runnerBinaryVersion = version
	apiBinaryVersion = version
	adminBinaryVersion = version
	setupBinaryVersion = version
	releaseBinaryVersion = version
}

func RunRelease(args []string, stdout, stderr io.Writer) int {
	return runRelease(args, stdout, stderr, releaseRealOps{})
}

func runRelease(args []string, stdout, stderr io.Writer, ops releaseOperations) int {
	if hasExactArg(args, "--help") {
		fmt.Fprintln(stdout, "Usage: adlaire-ci-release --repository owner/repository --tag V.X.N --commit 40-hex-sha --notes-file path --token-file path --out path [--version] [--help]")
		return 0
	}
	if hasExactArg(args, "--version") {
		fmt.Fprintf(stdout, "%s %s go=%s\n", releaseBinaryName, releaseBinaryVersion, runtime.Version())
		return 0
	}
	if !safeArgvTokens(args) {
		fmt.Fprintln(stderr, "invalid command line token")
		return 2
	}
	cfg, parseErr := parseReleaseArgs(args)
	if parseErr != "" {
		fmt.Fprintln(stderr, parseErr)
		return 2
	}
	state := releaseState{}
	if code, errCode := validateRelease(cfg, &state, ops); errCode != "" {
		return finishReleaseError(stderr, errCode, code)
	}
	if code, errCode := buildRelease(cfg, &state, ops); errCode != "" {
		return finishReleaseError(stderr, errCode, code)
	}
	if code, errCode := publishRelease(cfg, &state, ops); errCode != "" {
		return finishReleaseError(stderr, errCode, code)
	}
	writeReleaseSuccess(stdout, cfg, state)
	return 0
}

func finishReleaseError(stderr io.Writer, code string, exit int) int {
	fmt.Fprintf(stderr, "release: %s\n", code)
	return exit
}

func parseReleaseArgs(args []string) (releaseConfig, string) {
	cfg := releaseConfig{}
	seen := map[string]bool{}
	for i := 0; i < len(args); {
		name := args[i]
		if !strings.HasPrefix(name, "-") {
			return cfg, "usage error"
		}
		if !strings.HasPrefix(name, "--") || strings.Contains(name, "=") {
			return cfg, "unknown option: " + name
		}
		if !releaseKnownOption(name) {
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
		switch name {
		case "--repository":
			cfg.Repository = value
		case "--tag":
			cfg.Tag = value
		case "--commit":
			cfg.Commit = value
		case "--notes-file":
			cfg.NotesFile = value
		case "--token-file":
			cfg.TokenFile = value
		case "--out":
			cfg.Out = value
		}
		i += 2
	}
	if cfg.Repository == "" || cfg.Tag == "" || cfg.Commit == "" || cfg.NotesFile == "" || cfg.TokenFile == "" || cfg.Out == "" {
		return cfg, "usage error"
	}
	return cfg, ""
}

func releaseKnownOption(name string) bool {
	switch name {
	case "--repository", "--tag", "--commit", "--notes-file", "--token-file", "--out":
		return true
	default:
		return false
	}
}

func validateRelease(cfg releaseConfig, state *releaseState, ops releaseOperations) (int, string) {
	if !releaseRepositoryPattern.MatchString(cfg.Repository) || strings.Contains(cfg.Repository, "/./") || strings.Contains(cfg.Repository, "/../") {
		return 2, "INVALID_INPUT"
	}
	if !releaseTagPattern.MatchString(cfg.Tag) || cfg.Tag == releaseBinaryVersion {
		return 2, "INVALID_INPUT"
	}
	if !isLowerHex(cfg.Commit, 40) {
		return 2, "INVALID_INPUT"
	}
	root, err := ops.CheckoutRoot()
	if err != nil {
		return 2, "INVALID_INPUT"
	}
	state.CheckoutRoot = root
	notes, err := ops.ReadNotes(cfg.NotesFile, root)
	if err != nil {
		return 2, "INVALID_INPUT"
	}
	token, err := ops.ReadToken(cfg.TokenFile)
	if err != nil {
		return 2, "INVALID_INPUT"
	}
	state.Notes = notes
	state.Token = token
	if err := ops.PrepareOutput(cfg.Out, root); err != nil {
		return 2, "INVALID_INPUT"
	}
	if err := ops.VerifyClean(root); err != nil {
		return 2, "DIRTY_WORKTREE"
	}
	if err := ops.VerifyCommitAndTag(root, cfg.Commit, cfg.Tag); err != nil {
		return 2, "TAG_MISMATCH"
	}
	if _, err := ops.GitHub().CheckRemote(cfg, token); err != nil {
		return mapReleaseGitHubReadError(err)
	}
	if err := ops.GitHub().CheckReleaseAbsent(cfg, token); err != nil {
		if errors.Is(err, errReleaseExists) {
			return 2, "RELEASE_EXISTS"
		}
		return mapReleaseGitHubReadError(err)
	}
	return 0, ""
}

func buildRelease(cfg releaseConfig, state *releaseState, ops releaseOperations) (int, string) {
	assets, err := ops.BuildArtifacts(cfg, *state)
	if err != nil {
		return releaseBuildError(err)
	}
	if err := verifyReleaseAssetSet(assets); err != nil {
		return 1, "CHECKSUM_FAILED"
	}
	if err := ops.VerifyClean(state.CheckoutRoot); err != nil {
		return 2, "DIRTY_WORKTREE"
	}
	if err := ops.VerifyCommitAndTag(state.CheckoutRoot, cfg.Commit, cfg.Tag); err != nil {
		return 2, "TAG_MISMATCH"
	}
	if err := ops.WriteOutput(cfg.Out, assets); err != nil {
		return 1, "OUTPUT_FAILED"
	}
	state.Assets = assets
	return 0, ""
}

func publishRelease(cfg releaseConfig, state *releaseState, ops releaseOperations) (int, string) {
	gh := ops.GitHub()
	if _, err := gh.CheckRemote(cfg, state.Token); err != nil {
		return mapReleaseGitHubReadError(err)
	}
	if err := gh.CheckReleaseAbsent(cfg, state.Token); err != nil {
		if errors.Is(err, errReleaseExists) {
			return 2, "RELEASE_EXISTS"
		}
		return mapReleaseGitHubReadError(err)
	}
	draft, err := gh.CreateDraft(cfg, state.Token, state.Notes)
	if err != nil {
		return 3, "DRAFT_CREATE_FAILED"
	}
	state.ReleaseID = draft.ID
	state.ReleaseURL = draft.HTMLURL
	createdDraft := draft.ID > 0
	cleanup := func(originalExit int, originalCode string) (int, string) {
		if !createdDraft {
			return originalExit, originalCode
		}
		if err := gh.DeleteDraft(draft, cfg, state.Token); err != nil {
			return 1, "DRAFT_CLEANUP_FAILED"
		}
		return originalExit, originalCode
	}
	if _, err := gh.CheckRemote(cfg, state.Token); err != nil {
		return cleanup(3, "PUBLISH_FAILED")
	}
	for _, asset := range state.Assets {
		if err := gh.UploadAsset(draft, state.Token, asset, cfg.Repository); err != nil {
			return cleanup(3, "ASSET_UPLOAD_FAILED")
		}
	}
	if err := gh.VerifyAssets(draft, cfg, state.Assets, state.Token); err != nil {
		return cleanup(3, "ASSET_VERIFY_FAILED")
	}
	if _, err := gh.CheckRemote(cfg, state.Token); err != nil {
		return cleanup(3, "PUBLISH_FAILED")
	}
	if err := gh.Publish(draft, cfg, state.Token, state.Notes, state.Assets); err != nil {
		return cleanup(3, "PUBLISH_FAILED")
	}
	return 0, ""
}

func writeReleaseSuccess(stdout io.Writer, cfg releaseConfig, state releaseState) {
	assets := releaseAssetNames()
	obj := map[string]any{
		"tag":         cfg.Tag,
		"commit":      cfg.Commit,
		"release_url": state.ReleaseURL,
		"assets":      assets,
		"published":   true,
	}
	data, _ := json.Marshal(obj)
	// json.Marshal orders map keys lexicographically; the contract requires a stable
	// public line, so write it explicitly instead of relying on map ordering.
	_ = data
	fmt.Fprintf(stdout, `{"tag":%q,"commit":%q,"release_url":%q,"assets":[`, cfg.Tag, cfg.Commit, state.ReleaseURL)
	for i, name := range assets {
		if i > 0 {
			fmt.Fprint(stdout, ",")
		}
		fmt.Fprintf(stdout, "%q", name)
	}
	fmt.Fprint(stdout, `],"published":true}`+"\n")
}

func releaseAssetNames() []string {
	return []string{
		"adlaire-ci-build-linux-amd64",
		"adlaire-ci-runner-linux-amd64",
		"adlaire-ci-api-linux-amd64",
		"adlaire-ci-setup-linux-amd64",
		"adlaire-ci-admin-linux-amd64",
		"adlaire-ci-mcp-linux-amd64",
		"admin-ui.tar.gz",
		"SHA256SUMS",
	}
}

func releaseBinaryAssetNames() []string {
	return releaseAssetNames()[:6]
}

func releaseBuildError(err error) (int, string) {
	switch {
	case errors.Is(err, errReleaseFormatFailed):
		return 1, "FORMAT_FAILED"
	case errors.Is(err, errReleaseTestFailed):
		return 1, "TEST_FAILED"
	case errors.Is(err, errReleaseBuildFailed):
		return 1, "BUILD_FAILED"
	case errors.Is(err, errReleaseVersionMismatch):
		return 1, "VERSION_MISMATCH"
	case errors.Is(err, errReleaseNonReproducible):
		return 1, "NON_REPRODUCIBLE"
	case errors.Is(err, errReleaseArchiveFailed):
		return 1, "ARCHIVE_FAILED"
	case errors.Is(err, errReleaseChecksumFailed):
		return 1, "CHECKSUM_FAILED"
	default:
		return 1, "BUILD_FAILED"
	}
}

func mapReleaseGitHubReadError(err error) (int, string) {
	if errors.Is(err, errReleaseTagMismatch) {
		return 2, "TAG_MISMATCH"
	}
	return 3, "GITHUB_READ_FAILED"
}

var (
	errReleaseExists          = errors.New("release exists")
	errReleaseTagMismatch     = errors.New("tag mismatch")
	errReleaseFormatFailed    = errors.New("format failed")
	errReleaseTestFailed      = errors.New("test failed")
	errReleaseBuildFailed     = errors.New("build failed")
	errReleaseVersionMismatch = errors.New("version mismatch")
	errReleaseNonReproducible = errors.New("non reproducible")
	errReleaseArchiveFailed   = errors.New("archive failed")
	errReleaseChecksumFailed  = errors.New("checksum failed")
)

func verifyReleaseAssetSet(assets []releaseAsset) error {
	want := releaseAssetNames()
	if len(assets) != len(want) {
		return errors.New("asset count")
	}
	seen := map[string]releaseAsset{}
	for _, asset := range assets {
		if asset.Name == "" || len(asset.Data) == 0 || asset.Mode == 0 {
			return errors.New("invalid asset")
		}
		sum := sha256.Sum256(asset.Data)
		if asset.SHA256 != "" && asset.SHA256 != hex.EncodeToString(sum[:]) {
			return errors.New("asset checksum")
		}
		asset.SHA256 = hex.EncodeToString(sum[:])
		seen[asset.Name] = asset
	}
	for _, name := range want {
		if _, ok := seen[name]; !ok {
			return errors.New("missing asset")
		}
	}
	checksum, ok := seen["SHA256SUMS"]
	if !ok {
		return errors.New("missing checksum")
	}
	return validateReleaseChecksumManifest(checksum.Data, seen)
}

func validateReleaseChecksumManifest(data []byte, assets map[string]releaseAsset) error {
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	wantNames := append([]string{}, releaseAssetNames()[:7]...)
	sort.Strings(wantNames)
	if len(lines) != len(wantNames) {
		return errors.New("checksum line count")
	}
	for i, name := range wantNames {
		fields := strings.Split(lines[i], "  ")
		if len(fields) != 2 || fields[1] != name || len(fields[0]) != 64 {
			return errors.New("checksum line")
		}
		sum := sha256.Sum256(assets[name].Data)
		if fields[0] != hex.EncodeToString(sum[:]) {
			return errors.New("checksum mismatch")
		}
	}
	return nil
}

func (releaseRealOps) CheckoutRoot() (string, error) {
	result := releaseRunCommand("", 30*time.Second, "git", "rev-parse", "--show-toplevel")
	if result.ExitCode != 0 || len(result.Stderr) != 0 {
		return "", errors.New("checkout root")
	}
	root := strings.TrimSpace(string(result.Stdout))
	if root == "" || !filepath.IsAbs(root) {
		return "", errors.New("checkout root")
	}
	return root, nil
}

func (releaseRealOps) VerifyClean(root string) error {
	result := releaseRunCommand(root, 30*time.Second, "git", "status", "--porcelain=v1", "--untracked-files=all")
	if result.ExitCode != 0 || len(result.Stdout) != 0 {
		return errors.New("dirty")
	}
	return nil
}

func (releaseRealOps) VerifyCommitAndTag(root, commit, tag string) error {
	head := releaseRunCommand(root, 30*time.Second, "git", "rev-parse", "HEAD")
	if head.ExitCode != 0 || strings.TrimSpace(string(head.Stdout)) != commit {
		return errors.New("head mismatch")
	}
	tagCommit := releaseRunCommand(root, 30*time.Second, "git", "rev-parse", "refs/tags/"+tag+"^{commit}")
	if tagCommit.ExitCode != 0 || strings.TrimSpace(string(tagCommit.Stdout)) != commit {
		return errors.New("tag mismatch")
	}
	return nil
}

func (releaseRealOps) ReadNotes(path, checkoutRoot string) (string, error) {
	if path == "" || !filepath.IsAbs(path) {
		return "", errors.New("invalid notes path")
	}
	rel, err := filepath.Rel(checkoutRoot, path)
	if err != nil || strings.HasPrefix(rel, "..") || rel == "." || filepath.IsAbs(rel) {
		return "", errors.New("notes outside checkout")
	}
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() <= 0 || info.Size() > releaseMaxNotes {
		return "", errors.New("invalid notes")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if len(data) != int(info.Size()) || !utf8.Valid(data) || bytes.ContainsAny(data, "\x00\r") || !bytes.HasSuffix(data, []byte("\n")) || bytes.HasSuffix(data, []byte("\n\n")) {
		return "", errors.New("invalid notes data")
	}
	return string(data), nil
}

func (releaseRealOps) ReadToken(path string) (string, error) {
	if path == "" || !filepath.IsAbs(path) {
		return "", errors.New("invalid token path")
	}
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Mode().Perm() != 0600 || info.Size() <= 0 || info.Size() > releaseMaxToken+1 || !setupOwnedByCurrentUser(info) {
		return "", errors.New("invalid token")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	token := strings.TrimSpace(string(data))
	if token == "" || len(token) > releaseMaxToken || !utf8.ValidString(token) {
		return "", errors.New("invalid token data")
	}
	for _, r := range token {
		if r <= 0x20 || r == 0x7f {
			return "", errors.New("invalid token byte")
		}
	}
	return token, nil
}

func (releaseRealOps) PrepareOutput(path, checkoutRoot string) error {
	if path == "" || !filepath.IsAbs(path) || filepath.Clean(path) != path || path == string(filepath.Separator) {
		return errors.New("invalid output path")
	}
	if rel, err := filepath.Rel(checkoutRoot, path); err == nil && (rel == "." || !strings.HasPrefix(rel, "..")) {
		return errors.New("output inside checkout")
	}
	if _, err := os.Lstat(path); err == nil {
		return errors.New("output exists")
	} else if !os.IsNotExist(err) {
		return err
	}
	if _, err := os.Lstat(path + ".tmp"); err == nil {
		return errors.New("staging exists")
	} else if !os.IsNotExist(err) {
		return err
	}
	parent := filepath.Dir(path)
	info, err := os.Lstat(parent)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("invalid output parent")
	}
	return nil
}

func (releaseRealOps) BuildArtifacts(cfg releaseConfig, state releaseState) ([]releaseAsset, error) {
	a, err := os.MkdirTemp("", "adlaire-ci-release-a-")
	if err != nil {
		return nil, errReleaseBuildFailed
	}
	b, err := os.MkdirTemp("", "adlaire-ci-release-b-")
	if err != nil {
		os.RemoveAll(a)
		return nil, errReleaseBuildFailed
	}
	defer os.RemoveAll(a)
	defer os.RemoveAll(b)
	if err := extractGitArchive(state.CheckoutRoot, cfg.Commit, a); err != nil {
		return nil, errReleaseBuildFailed
	}
	if err := extractGitArchive(state.CheckoutRoot, cfg.Commit, b); err != nil {
		return nil, errReleaseBuildFailed
	}
	if err := compareSnapshotTrees(a, b); err != nil {
		return nil, errReleaseNonReproducible
	}
	modTime, err := releaseCommitTimestamp(state.CheckoutRoot, cfg.Commit)
	if err != nil {
		return nil, err
	}
	if err := gofmtSnapshot(a); err != nil {
		return nil, err
	}
	if result := releaseRunCommand(a, 15*time.Minute, "go", "test", "./..."); result.ExitCode != 0 || len(result.Stderr) > releaseMaxOutput || len(result.Stdout) > releaseMaxOutput {
		return nil, errReleaseTestFailed
	}
	assetsA, err := buildSnapshotAssets(a, cfg.Tag, modTime)
	if err != nil {
		return nil, err
	}
	assetsB, err := buildSnapshotAssets(b, cfg.Tag, modTime)
	if err != nil {
		return nil, err
	}
	if !sameReleaseAssets(assetsA, assetsB) {
		return nil, errReleaseNonReproducible
	}
	return withChecksumAsset(assetsA), nil
}

func (releaseRealOps) WriteOutput(path string, assets []releaseAsset) error {
	staging := path + ".tmp"
	if err := os.Mkdir(staging, 0700); err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			os.RemoveAll(staging)
		}
	}()
	for _, asset := range assets {
		target := filepath.Join(staging, asset.Name)
		if err := os.WriteFile(target, asset.Data, asset.Mode); err != nil {
			return err
		}
		if err := os.Chmod(target, asset.Mode); err != nil {
			return err
		}
	}
	if err := os.Rename(staging, path); err != nil {
		return err
	}
	ok = true
	return nil
}

func (releaseRealOps) GitHub() releaseGitHubOps {
	return releaseRealGitHubOps{}
}

func extractGitArchive(checkoutRoot, commit, dest string) error {
	cmd := exec.Command("git", "archive", "--format=tar", commit)
	cmd.Dir = checkoutRoot
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	tr := tar.NewReader(io.LimitReader(stdout, releaseMaxArchive+1))
	seen := map[string]bool{}
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			cmd.Wait()
			return err
		}
		if !safeArchivePath(hdr.Name) || seen[hdr.Name] {
			cmd.Wait()
			return errors.New("unsafe archive")
		}
		seen[hdr.Name] = true
		target := filepath.Join(dest, filepath.FromSlash(hdr.Name))
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				cmd.Wait()
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				cmd.Wait()
				return err
			}
			mode := os.FileMode(0644)
			if hdr.FileInfo().Mode()&0111 != 0 {
				mode = 0755
			}
			data, err := io.ReadAll(io.LimitReader(tr, hdr.Size+1))
			if err != nil || int64(len(data)) != hdr.Size {
				cmd.Wait()
				return errors.New("archive size")
			}
			if err := os.WriteFile(target, data, mode); err != nil {
				cmd.Wait()
				return err
			}
		default:
			cmd.Wait()
			return errors.New("unsafe archive type")
		}
	}
	if err := cmd.Wait(); err != nil || stderr.Len() != 0 {
		return errors.New("git archive failed")
	}
	return nil
}

func safeArchivePath(name string) bool {
	if name == "" || strings.HasPrefix(name, "/") || strings.Contains(name, "\\") || strings.ContainsAny(name, "\x00\r\n") {
		return false
	}
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return utf8.ValidString(name)
}

func compareSnapshotTrees(a, b string) error {
	left, err := snapshotDigests(a)
	if err != nil {
		return err
	}
	right, err := snapshotDigests(b)
	if err != nil {
		return err
	}
	if len(left) != len(right) {
		return errors.New("tree mismatch")
	}
	for path, l := range left {
		if right[path] != l {
			return errors.New("tree mismatch")
		}
	}
	return nil
}

func snapshotDigests(root string) (map[string]string, error) {
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return errors.New("unsafe file")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		rel, _ := filepath.Rel(root, path)
		out[filepath.ToSlash(rel)] = fmt.Sprintf("%o:%d:%s", info.Mode().Perm(), info.Size(), hex.EncodeToString(sum[:]))
		return nil
	})
	return out, err
}

func gofmtSnapshot(root string) error {
	files := []string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil || len(files) == 0 {
		return errReleaseFormatFailed
	}
	sort.Strings(files)
	for _, file := range files {
		result := releaseRunCommand(root, 2*time.Minute, "gofmt", "-d", file)
		if result.ExitCode != 0 || len(result.Stdout) != 0 || len(result.Stderr) != 0 {
			return errReleaseFormatFailed
		}
	}
	return nil
}

func releaseCommitTimestamp(root, commit string) (time.Time, error) {
	result := releaseRunCommand(root, 30*time.Second, "git", "show", "-s", "--format=%ct", commit)
	text := strings.TrimSuffix(string(result.Stdout), "\n")
	if result.ExitCode != 0 || len(result.Stderr) != 0 || text == "" || strings.Trim(text, "0123456789") != "" {
		return time.Time{}, errReleaseBuildFailed
	}
	seconds, err := strconv.ParseInt(text, 10, 64)
	if err != nil || seconds <= 0 {
		return time.Time{}, errReleaseBuildFailed
	}
	return time.Unix(seconds, 0).UTC(), nil
}

func buildSnapshotAssets(root, tag string, modTime time.Time) ([]releaseAsset, error) {
	assets := []releaseAsset{}
	for _, name := range releaseBinaryAssetNames() {
		out := filepath.Join(root, ".release-out", name)
		if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
			return nil, errReleaseBuildFailed
		}
		ldflags := "-s -w -X main.binaryVersion=" + tag
		result := releaseRunCommand(root, 10*time.Minute, "go", "build", "-trimpath", "-buildvcs=false", "-ldflags", ldflags, "-o", out, ".")
		if result.ExitCode != 0 {
			return nil, errReleaseBuildFailed
		}
		if err := verifyReleaseBinary(out, strings.TrimSuffix(name, "-linux-amd64"), tag); err != nil {
			return nil, err
		}
		data, err := os.ReadFile(out)
		if err != nil || len(data) == 0 {
			return nil, errReleaseBuildFailed
		}
		assets = append(assets, releaseAsset{Name: name, Data: data, Mode: 0755})
	}
	admin, err := buildAdminArchive(root, modTime)
	if err != nil {
		return nil, err
	}
	assets = append(assets, releaseAsset{Name: "admin-ui.tar.gz", Data: admin, Mode: 0644})
	return assets, nil
}

func verifyReleaseBinary(path, binary, tag string) error {
	result := releaseRunCommand("", 10*time.Second, path, "--version")
	if result.ExitCode != 0 || len(result.Stderr) != 0 {
		return errReleaseBuildFailed
	}
	wantPrefix := binary + " " + tag + " go="
	if !strings.HasPrefix(string(result.Stdout), wantPrefix) || !strings.HasSuffix(string(result.Stdout), "\n") || strings.Contains(string(result.Stdout), releaseBinaryVersion) {
		return errReleaseVersionMismatch
	}
	return nil
}

func buildAdminArchive(root string, modTime time.Time) ([]byte, error) {
	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, errReleaseArchiveFailed
	}
	gz.Name = ""
	gz.Comment = ""
	gz.ModTime = modTime
	tw := tar.NewWriter(gz)
	for _, name := range []string{"adlaire-ci-sdk.js", "index.html"} {
		data, err := os.ReadFile(filepath.Join(root, "admin", name))
		if err != nil || len(data) == 0 {
			tw.Close()
			gz.Close()
			return nil, errReleaseArchiveFailed
		}
		hdr := &tar.Header{Name: name, Typeflag: tar.TypeReg, Format: tar.FormatUSTAR, Mode: 0644, Size: int64(len(data)), ModTime: modTime}
		if err := tw.WriteHeader(hdr); err != nil {
			tw.Close()
			gz.Close()
			return nil, errReleaseArchiveFailed
		}
		if _, err := tw.Write(data); err != nil {
			tw.Close()
			gz.Close()
			return nil, errReleaseArchiveFailed
		}
	}
	if err := tw.Close(); err != nil {
		gz.Close()
		return nil, errReleaseArchiveFailed
	}
	if err := gz.Close(); err != nil {
		return nil, errReleaseArchiveFailed
	}
	return buf.Bytes(), nil
}

func sameReleaseAssets(a, b []releaseAsset) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name || a[i].Mode != b[i].Mode || !bytes.Equal(a[i].Data, b[i].Data) {
			return false
		}
	}
	return true
}

func withChecksumAsset(assets []releaseAsset) []releaseAsset {
	var builder strings.Builder
	names := make([]string, 0, len(assets))
	byName := map[string]releaseAsset{}
	for _, asset := range assets {
		names = append(names, asset.Name)
		byName[asset.Name] = asset
	}
	sort.Strings(names)
	for _, name := range names {
		sum := sha256.Sum256(byName[name].Data)
		fmt.Fprintf(&builder, "%s  %s\n", hex.EncodeToString(sum[:]), name)
	}
	return append(assets, releaseAsset{Name: "SHA256SUMS", Data: []byte(builder.String()), Mode: 0644})
}

func releaseRunCommand(dir string, timeout time.Duration, name string, args ...string) releaseCommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = releaseCommandEnv(os.Environ())
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
	return releaseCommandResult{Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), ExitCode: exitCode}
}

func releaseCommandEnv(base []string) []string {
	keep := map[string]bool{"PATH": true, "HOME": true, "TMPDIR": true, "GOPATH": true, "GOMODCACHE": true}
	out := []string{}
	for _, entry := range base {
		key, _, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		if strings.HasPrefix(key, "GO") && !keep[key] {
			continue
		}
		if keep[key] {
			out = append(out, entry)
		}
	}
	fixed := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        "linux",
		"GOARCH":      "amd64",
		"GOFLAGS":     "",
		"GOENV":       "off",
		"GOWORK":      "off",
		"GOTOOLCHAIN": "local",
		"GOPROXY":     "off",
		"GOSUMDB":     "off",
	}
	for key, value := range fixed {
		out = append(out, key+"="+value)
	}
	sort.Strings(out)
	return out
}

func (releaseRealGitHubOps) CheckRemote(cfg releaseConfig, token string) (string, error) {
	branch, err := releaseVerifyRemoteRefs(cfg, token)
	if err != nil {
		return "", err
	}
	return branch, nil
}

func (releaseRealGitHubOps) CheckReleaseAbsent(cfg releaseConfig, token string) error {
	owner, repo, _ := strings.Cut(cfg.Repository, "/")
	req, err := releaseGitHubRequest(http.MethodGet, "https://api.github.com/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/releases/tags/"+url.PathEscape(cfg.Tag), token, nil)
	if err != nil {
		return err
	}
	resp, _, err := releaseDoJSON(req, releaseHTTPTimeout, 1024*1024)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode == http.StatusOK {
		return errReleaseExists
	}
	return errors.New("release read")
}

func (releaseRealGitHubOps) CreateDraft(cfg releaseConfig, token, notes string) (releaseGitHubRelease, error) {
	owner, repo, _ := releaseSplitRepository(cfg.Repository)
	body, _ := json.Marshal(map[string]any{"tag_name": cfg.Tag, "target_commitish": cfg.Commit, "name": cfg.Tag, "body": notes, "draft": true, "prerelease": false, "generate_release_notes": false})
	req, err := releaseGitHubRequest(http.MethodPost, "https://api.github.com/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/releases", token, bytes.NewReader(body))
	if err != nil {
		return releaseGitHubRelease{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, raw, err := releaseDoJSON(req, releaseHTTPTimeout, 4*1024*1024)
	if err != nil || resp.StatusCode != http.StatusCreated {
		return releaseGitHubRelease{}, errors.New("draft create")
	}
	var parsed struct {
		ID         int64  `json:"id"`
		UploadURL  string `json:"upload_url"`
		HTMLURL    string `json:"html_url"`
		TagName    string `json:"tag_name"`
		Target     string `json:"target_commitish"`
		Name       string `json:"name"`
		Body       string `json:"body"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil || parsed.ID <= 0 || parsed.TagName != cfg.Tag || parsed.Target != cfg.Commit || parsed.Name != cfg.Tag || parsed.Body != notes || !parsed.Draft || parsed.Prerelease || !validReleaseURL(parsed.HTMLURL, cfg.Repository, cfg.Tag) {
		return releaseGitHubRelease{}, errors.New("draft body")
	}
	if _, err := releaseUploadURL(parsed.UploadURL, "probe"); err != nil {
		return releaseGitHubRelease{}, err
	}
	return releaseGitHubRelease{ID: parsed.ID, UploadURL: parsed.UploadURL, HTMLURL: parsed.HTMLURL, TagName: parsed.TagName, TargetCommitish: parsed.Target, Name: parsed.Name, Body: parsed.Body, Draft: parsed.Draft, Prerelease: parsed.Prerelease}, nil
}

func (releaseRealGitHubOps) UploadAsset(rel releaseGitHubRelease, token string, asset releaseAsset, repository string) error {
	upload, err := releaseUploadURL(rel.UploadURL, asset.Name)
	if err != nil {
		return err
	}
	req, err := releaseGitHubRequest(http.MethodPost, upload, token, bytes.NewReader(asset.Data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = int64(len(asset.Data))
	resp, raw, err := releaseDoJSON(req, releaseAssetTimeout, 4*1024*1024)
	if err != nil || resp.StatusCode != http.StatusCreated {
		return errors.New("upload")
	}
	var parsed struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Size  int64  `json:"size"`
		State string `json:"state"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil || parsed.ID <= 0 || parsed.Name != asset.Name || parsed.Size != int64(len(asset.Data)) || parsed.State != "uploaded" {
		return errors.New("upload body")
	}
	return nil
}

func (releaseRealGitHubOps) VerifyAssets(rel releaseGitHubRelease, cfg releaseConfig, assets []releaseAsset, token string) error {
	if rel.ID <= 0 {
		return errors.New("release id")
	}
	owner, repo, _ := releaseSplitRepository(cfg.Repository)
	infos, err := releaseListAssets(owner, repo, rel.ID, token)
	if err != nil {
		return err
	}
	if len(infos) != len(assets) {
		return errors.New("asset count")
	}
	want := map[string]releaseAsset{}
	for _, asset := range assets {
		want[asset.Name] = asset
	}
	seenID := map[int64]bool{}
	seenName := map[string]bool{}
	for _, info := range infos {
		if info.ID <= 0 || seenID[info.ID] || seenName[info.Name] {
			return errors.New("asset identity")
		}
		seenID[info.ID] = true
		seenName[info.Name] = true
		asset, ok := want[info.Name]
		if !ok || info.Size != int64(len(asset.Data)) {
			return errors.New("asset metadata")
		}
		data, err := releaseDownloadAsset(owner, repo, info.ID, info.Size, token)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		if len(data) != len(asset.Data) || hex.EncodeToString(sum[:]) != asset.SHA256 && asset.SHA256 != "" {
			return errors.New("asset digest")
		}
		local := sha256.Sum256(asset.Data)
		if !bytes.Equal(sum[:], local[:]) {
			return errors.New("asset digest")
		}
	}
	return nil
}

func (releaseRealGitHubOps) Publish(rel releaseGitHubRelease, cfg releaseConfig, token, notes string, assets []releaseAsset) error {
	owner, repo, _ := releaseSplitRepository(cfg.Repository)
	if err := releaseVerifyReleaseMetadata(owner, repo, rel.ID, cfg, token, notes, true, assets); err != nil {
		return err
	}
	body := strings.NewReader(`{"draft":false,"prerelease":false}`)
	req, err := releaseGitHubRequest(http.MethodPatch, fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d", url.PathEscape(owner), url.PathEscape(repo), rel.ID), token, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, raw, err := releaseDoJSON(req, releaseHTTPTimeout, 4*1024*1024)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("publish")
	}
	if err := releaseValidateReleaseMetadata(raw, rel.ID, cfg, notes, false, assets); err != nil {
		return err
	}
	if err := releaseVerifyReleaseMetadata(owner, repo, rel.ID, cfg, token, notes, false, assets); err != nil {
		return err
	}
	return nil
}

func (releaseRealGitHubOps) DeleteDraft(rel releaseGitHubRelease, cfg releaseConfig, token string) error {
	if rel.ID <= 0 {
		return nil
	}
	owner, repo, _ := releaseSplitRepository(cfg.Repository)
	req, err := releaseGitHubRequest(http.MethodDelete, fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d", url.PathEscape(owner), url.PathEscape(repo), rel.ID), token, nil)
	if err != nil {
		return err
	}
	resp, body, err := releaseDoJSON(req, releaseHTTPTimeout, 1024*1024)
	if err != nil || resp.StatusCode != http.StatusNoContent || len(body) != 0 {
		return errors.New("delete draft")
	}
	return nil
}

type releaseAssetInfo struct {
	ID   int64
	Name string
	Size int64
}

func releaseVerifyRemoteRefs(cfg releaseConfig, token string) (string, error) {
	owner, repo, ok := releaseSplitRepository(cfg.Repository)
	if !ok {
		return "", errors.New("repository")
	}
	branch, err := releaseReadRepository(owner, repo, cfg.Repository, token)
	if err != nil {
		return "", err
	}
	head, err := releaseReadBranchHead(owner, repo, branch, token)
	if err != nil {
		return "", err
	}
	if err := releaseVerifyCompare(owner, repo, cfg.Commit, head, token); err != nil {
		return "", err
	}
	tagCommit, err := releaseReadTagCommit(owner, repo, cfg.Tag, token)
	if err != nil {
		return "", err
	}
	if tagCommit != cfg.Commit {
		return "", errReleaseTagMismatch
	}
	return branch, nil
}

func releaseSplitRepository(repository string) (string, string, bool) {
	owner, repo, ok := strings.Cut(repository, "/")
	return owner, repo, ok && owner != "" && repo != ""
}

func releaseReadRepository(owner, repo, fullName, token string) (string, error) {
	req, err := releaseGitHubRequest(http.MethodGet, "https://api.github.com/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo), token, nil)
	if err != nil {
		return "", err
	}
	resp, body, err := releaseDoJSON(req, releaseHTTPTimeout, 4*1024*1024)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("repository status")
	}
	var parsed struct {
		FullName      string `json:"full_name"`
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || parsed.FullName != fullName || !validGitHubRefName(parsed.DefaultBranch) {
		return "", errors.New("repository body")
	}
	return parsed.DefaultBranch, nil
}

func validGitHubRefName(name string) bool {
	return name != "" && len(name) <= 255 && utf8.ValidString(name) && strings.TrimSpace(name) == name && !strings.ContainsAny(name, "\x00\r\n") && !strings.HasPrefix(name, "refs/") && name != "." && name != ".."
}

func releaseReadBranchHead(owner, repo, branch, token string) (string, error) {
	req, err := releaseGitHubRequest(http.MethodGet, "https://api.github.com/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/git/ref/heads/"+url.PathEscape(branch), token, nil)
	if err != nil {
		return "", err
	}
	resp, body, err := releaseDoJSON(req, releaseHTTPTimeout, 4*1024*1024)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", errReleaseTagMismatch
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("branch status")
	}
	var parsed struct {
		Ref    string `json:"ref"`
		Object struct {
			Type string `json:"type"`
			SHA  string `json:"sha"`
		} `json:"object"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || parsed.Ref != "refs/heads/"+branch || parsed.Object.Type != "commit" || !isLowerHex(parsed.Object.SHA, 40) {
		return "", errors.New("branch body")
	}
	return parsed.Object.SHA, nil
}

func releaseVerifyCompare(owner, repo, commit, defaultHead, token string) error {
	req, err := releaseGitHubRequest(http.MethodGet, "https://api.github.com/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/compare/"+url.PathEscape(commit+"..."+defaultHead), token, nil)
	if err != nil {
		return err
	}
	resp, body, err := releaseDoJSON(req, releaseHTTPTimeout, 4*1024*1024)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusConflict {
		return errReleaseTagMismatch
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("compare status")
	}
	var parsed struct {
		Status     string `json:"status"`
		BaseCommit struct {
			SHA string `json:"sha"`
		} `json:"base_commit"`
		MergeBaseCommit struct {
			SHA string `json:"sha"`
		} `json:"merge_base_commit"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || parsed.BaseCommit.SHA != commit || parsed.MergeBaseCommit.SHA != commit || (parsed.Status != "ahead" && parsed.Status != "identical") {
		return errReleaseTagMismatch
	}
	return nil
}

func releaseReadTagCommit(owner, repo, tag, token string) (string, error) {
	req, err := releaseGitHubRequest(http.MethodGet, "https://api.github.com/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/git/ref/tags/"+url.PathEscape(tag), token, nil)
	if err != nil {
		return "", err
	}
	resp, body, err := releaseDoJSON(req, releaseHTTPTimeout, 4*1024*1024)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", errReleaseTagMismatch
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("tag status")
	}
	var parsed struct {
		Object struct {
			Type string `json:"type"`
			SHA  string `json:"sha"`
		} `json:"object"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || !isLowerHex(parsed.Object.SHA, 40) {
		return "", errors.New("tag body")
	}
	switch parsed.Object.Type {
	case "commit":
		return parsed.Object.SHA, nil
	case "tag":
		return releasePeelAnnotatedTag(owner, repo, parsed.Object.SHA, token)
	default:
		return "", errReleaseTagMismatch
	}
}

func releasePeelAnnotatedTag(owner, repo, sha, token string) (string, error) {
	seen := map[string]bool{}
	for i := 0; i < 5; i++ {
		if seen[sha] || !isLowerHex(sha, 40) {
			return "", errReleaseTagMismatch
		}
		seen[sha] = true
		req, err := releaseGitHubRequest(http.MethodGet, "https://api.github.com/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/git/tags/"+url.PathEscape(sha), token, nil)
		if err != nil {
			return "", err
		}
		resp, body, err := releaseDoJSON(req, releaseHTTPTimeout, 4*1024*1024)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == http.StatusNotFound {
			return "", errReleaseTagMismatch
		}
		if resp.StatusCode != http.StatusOK {
			return "", errors.New("tag object status")
		}
		var parsed struct {
			Object struct {
				Type string `json:"type"`
				SHA  string `json:"sha"`
			} `json:"object"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil || !isLowerHex(parsed.Object.SHA, 40) {
			return "", errors.New("tag object body")
		}
		if parsed.Object.Type == "commit" {
			return parsed.Object.SHA, nil
		}
		if parsed.Object.Type != "tag" {
			return "", errReleaseTagMismatch
		}
		sha = parsed.Object.SHA
	}
	return "", errReleaseTagMismatch
}

func releaseListAssets(owner, repo string, releaseID int64, token string) ([]releaseAssetInfo, error) {
	req, err := releaseGitHubRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d/assets?per_page=100&page=1", url.PathEscape(owner), url.PathEscape(repo), releaseID), token, nil)
	if err != nil {
		return nil, err
	}
	resp, body, err := releaseDoJSON(req, releaseHTTPTimeout, 4*1024*1024)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK || strings.Contains(resp.Header.Get("Link"), `rel="next"`) {
		return nil, errors.New("asset list")
	}
	var parsed []struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Size  int64  `json:"size"`
		State string `json:"state"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	infos := make([]releaseAssetInfo, 0, len(parsed))
	for _, item := range parsed {
		if item.ID <= 0 || item.Name == "" || item.Size <= 0 || (item.State != "" && item.State != "uploaded") {
			return nil, errors.New("asset list body")
		}
		infos = append(infos, releaseAssetInfo{ID: item.ID, Name: item.Name, Size: item.Size})
	}
	return infos, nil
}

func releaseDownloadAsset(owner, repo string, assetID, size int64, token string) ([]byte, error) {
	raw := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/assets/%d", url.PathEscape(owner), url.PathEscape(repo), assetID)
	req, err := releaseGitHubRequest(http.MethodGet, raw, token, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")
	resp, body, err := releaseDoJSON(req, releaseAssetTimeout, size)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusFound {
		location := resp.Header.Get("Location")
		u, err := url.Parse(location)
		if err != nil || u.Scheme != "https" || u.Host != "release-assets.githubusercontent.com" || u.User != nil || u.Fragment != "" {
			return nil, errors.New("asset redirect")
		}
		redirectReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		redirectReq.Header.Set("User-Agent", releaseBinaryName)
		resp, body, err = releaseDoJSON(redirectReq, releaseAssetTimeout, size)
		if err != nil {
			return nil, err
		}
	}
	if resp.StatusCode != http.StatusOK || resp.ContentLength != size || int64(len(body)) != size {
		return nil, errors.New("asset download")
	}
	return body, nil
}

func releaseVerifyReleaseMetadata(owner, repo string, releaseID int64, cfg releaseConfig, token, notes string, draft bool, assets []releaseAsset) error {
	req, err := releaseGitHubRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d", url.PathEscape(owner), url.PathEscape(repo), releaseID), token, nil)
	if err != nil {
		return err
	}
	resp, body, err := releaseDoJSON(req, releaseHTTPTimeout, 4*1024*1024)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("release metadata")
	}
	return releaseValidateReleaseMetadata(body, releaseID, cfg, notes, draft, assets)
}

func releaseValidateReleaseMetadata(body []byte, releaseID int64, cfg releaseConfig, notes string, draft bool, assets []releaseAsset) error {
	var parsed struct {
		ID         int64           `json:"id"`
		TagName    string          `json:"tag_name"`
		Target     string          `json:"target_commitish"`
		Name       string          `json:"name"`
		Body       string          `json:"body"`
		Draft      bool            `json:"draft"`
		Prerelease bool            `json:"prerelease"`
		Assets     json.RawMessage `json:"assets"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || parsed.ID != releaseID || parsed.TagName != cfg.Tag || parsed.Target != cfg.Commit || parsed.Name != cfg.Tag || parsed.Body != notes || parsed.Draft != draft || parsed.Prerelease {
		return errors.New("release metadata body")
	}
	var listed []any
	if len(parsed.Assets) > 0 {
		if err := json.Unmarshal(parsed.Assets, &listed); err != nil {
			return errors.New("release assets body")
		}
		if len(listed) != len(assets) {
			return errors.New("release assets count")
		}
	}
	return nil
}

func releaseGitHubRequest(method, rawurl, token string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, rawurl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", releaseBinaryName)
	return req, nil
}

func releaseDoJSON(req *http.Request, timeout time.Duration, limit int64) (*http.Response, []byte, error) {
	client := &http.Client{Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, nil, err
	}
	if int64(len(body)) > limit {
		return nil, nil, errors.New("body too large")
	}
	return resp, body, nil
}

func validReleaseURL(raw, repository, tag string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host != "github.com" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return false
	}
	return u.Path == "/"+repository+"/releases/tag/"+tag
}

func releaseUploadURL(raw, name string) (string, error) {
	if !strings.HasSuffix(raw, "{?name,label}") {
		return "", errors.New("invalid upload url")
	}
	raw = strings.TrimSuffix(raw, "{?name,label}")
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host != "uploads.github.com" || u.User != nil || u.Fragment != "" || u.RawQuery != "" {
		return "", errors.New("invalid upload url")
	}
	q := url.Values{}
	q.Set("name", name)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
