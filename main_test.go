package main

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestDispatchMainExactBasename(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if code := dispatchMain("prefix-adlaire-ci-runner", []string{"--version"}, &stdout, &stderr); code != 2 {
		t.Fatalf("expected unknown command exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout for unknown command, got %q", stdout.String())
	}
	if got := stderr.String(); !strings.Contains(got, "unknown command: prefix-adlaire-ci-runner") {
		t.Fatalf("expected unknown command error, got %q", got)
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-build", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected builder help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-build") {
		t.Fatalf("expected builder help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for builder help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-runner", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected runner help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-runner") {
		t.Fatalf("expected runner help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for runner help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-api", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected api help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-api") {
		t.Fatalf("expected api help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for api help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-admin", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected admin help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-admin") {
		t.Fatalf("expected admin help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for admin help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-setup", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected setup help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-setup") {
		t.Fatalf("expected setup help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for setup help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-release", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected release help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-release") {
		t.Fatalf("expected release help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for release help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-mcp", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected mcp help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-mcp") {
		t.Fatalf("expected mcp help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for mcp help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-build-linux-amd64", []string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected release asset build version exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.HasPrefix(got, "adlaire-ci-build V.0.0-dev go=") {
		t.Fatalf("expected canonical build version, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for release asset build version, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-mcp-linux-amd64", []string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected release asset mcp version exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.HasPrefix(got, "adlaire-ci-mcp V.0.0-dev go=") {
		t.Fatalf("expected canonical mcp version, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for release asset mcp version, got %q", stderr.String())
	}
}

func TestPhase11StandardArtifactInventory(t *testing.T) {
	t.Parallel()

	requiredFiles := []string{
		"main.go",
		"main_test.go",
		"sdk_contract_test.go",
		"ui_contract_test.go",
		"components/admin.go",
		"components/admin_test.go",
		"components/api.go",
		"components/api_test.go",
		"components/builder.go",
		"components/builder_test.go",
		"components/mcp.go",
		"components/mcp_test.go",
		"components/release.go",
		"components/release_test.go",
		"components/runner.go",
		"components/runner_test.go",
		"components/setup.go",
		"components/setup_test.go",
		"admin/adlaire-ci-sdk.js",
		"admin/index.html",
	}
	for _, path := range requiredFiles {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("required Phase 11 artifact %s is missing: %v", path, err)
		}
		if info.IsDir() {
			t.Fatalf("required Phase 11 artifact %s must be a file", path)
		}
	}

	if _, err := os.Stat("cmd"); !os.IsNotExist(err) {
		t.Fatalf("cmd directory must not exist in the standard Phase 11 layout: %v", err)
	}

	expectedComponentFiles := []string{
		"components/admin.go",
		"components/admin_test.go",
		"components/api.go",
		"components/api_test.go",
		"components/builder.go",
		"components/builder_test.go",
		"components/mcp.go",
		"components/mcp_test.go",
		"components/release.go",
		"components/release_test.go",
		"components/runner.go",
		"components/runner_test.go",
		"components/setup.go",
		"components/setup_test.go",
	}
	var componentFiles []string
	entries, err := os.ReadDir("components")
	if err != nil {
		t.Fatalf("read components directory: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			t.Fatalf("components must remain a one-file-per-artifact directory, unexpected subdirectory %s", entry.Name())
		}
		if strings.HasSuffix(entry.Name(), ".go") {
			componentFiles = append(componentFiles, filepath.ToSlash(filepath.Join("components", entry.Name())))
		}
	}
	sort.Strings(componentFiles)
	if !phase11StringSlicesEqual(componentFiles, expectedComponentFiles) {
		t.Fatalf("unexpected components Go files\nwant: %v\n got: %v", expectedComponentFiles, componentFiles)
	}

	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, path := range requiredFiles {
		if !strings.Contains(documentIndex, path) {
			t.Fatalf("docs/DOCUMENT_INDEX.md does not index required artifact %s", path)
		}
	}
}

func TestPhase11FixtureManifestGate(t *testing.T) {
	t.Parallel()

	manifestPaths, err := phase11FindManifestFiles("testdata")
	if err != nil {
		t.Fatalf("find fixture manifests: %v", err)
	}
	if len(manifestPaths) < 42 {
		t.Fatalf("expected at least 42 formal fixture manifests, got %d", len(manifestPaths))
	}

	seenNames := make(map[string]string)
	for _, manifestPath := range manifestPaths {
		manifest := phase11ReadManifest(t, manifestPath)
		fixtureDir := filepath.Dir(manifestPath)

		if manifest.Name != filepath.Base(fixtureDir) {
			t.Fatalf("%s name must match fixture directory, got %q", manifestPath, manifest.Name)
		}
		if prior := seenNames[manifest.Name]; prior != "" {
			t.Fatalf("duplicate fixture manifest name %q in %s and %s", manifest.Name, prior, manifestPath)
		}
		seenNames[manifest.Name] = manifestPath

		if strings.TrimSpace(manifest.Section) == "" {
			t.Fatalf("%s section must not be empty", manifestPath)
		}
		if strings.TrimSpace(manifest.Feature) == "" {
			t.Fatalf("%s feature must not be empty", manifestPath)
		}
		if !phase11AllowedCategories[manifest.Category] {
			t.Fatalf("%s has unsupported category %q", manifestPath, manifest.Category)
		}
		if strings.HasPrefix(manifest.Name, manifest.Category+"-") == false && phase11NameHasCategoryPrefix(manifest.Name) {
			t.Fatalf("%s category %q must match fixture name prefix %q", manifestPath, manifest.Category, manifest.Name)
		}

		if !phase11AllowedComponents[manifest.OwnerComponent] {
			t.Fatalf("%s has unsupported owner component %q", manifestPath, manifest.OwnerComponent)
		}
		phase11RequireSortedUnique(t, manifestPath, "collaborator_components", manifest.CollaboratorComponents)
		for _, component := range manifest.CollaboratorComponents {
			if !phase11AllowedComponents[component] {
				t.Fatalf("%s has unsupported collaborator component %q", manifestPath, component)
			}
			if component == manifest.OwnerComponent {
				t.Fatalf("%s collaborator_components must not repeat owner component %q", manifestPath, component)
			}
		}
		wantComponents := append([]string{manifest.OwnerComponent}, manifest.CollaboratorComponents...)
		sort.Strings(wantComponents)
		if !phase11StringSlicesEqual(manifest.Components, wantComponents) {
			t.Fatalf("%s components must be sorted owner+collaborators\nwant: %v\n got: %v", manifestPath, wantComponents, manifest.Components)
		}

		phase11RequireSortedUnique(t, manifestPath, "components", manifest.Components)
		if len(manifest.References) == 0 {
			t.Fatalf("%s references must not be empty", manifestPath)
		}
		phase11RequireUnique(t, manifestPath, "references", manifest.References)
		for _, reference := range manifest.References {
			phase11RequireMarkdownReference(t, manifestPath, reference)
		}

		if manifest.FakeClock != nil {
			if _, err := time.Parse("2006-01-02T15:04:05Z", *manifest.FakeClock); err != nil {
				t.Fatalf("%s fake_clock must be UTC seconds precision, got %q: %v", manifestPath, *manifest.FakeClock, err)
			}
			if !strings.HasSuffix(*manifest.FakeClock, "Z") {
				t.Fatalf("%s fake_clock must use UTC Z suffix, got %q", manifestPath, *manifest.FakeClock)
			}
		}

		phase11RequireManifestRecords(t, manifestPath, "not_applicable", manifest.NotApplicable)
		phase11RequireSafeRelativePaths(t, manifestPath, "missing_state", manifest.MissingState)
		phase11RequireSafeRelativePaths(t, manifestPath, "input_files", manifest.InputFiles)
		for _, input := range manifest.InputFiles {
			path := filepath.Join(fixtureDir, filepath.FromSlash(input))
			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("%s input file %s is missing: %v", manifestPath, input, err)
			}
			if info.IsDir() {
				t.Fatalf("%s input file %s must not be a directory", manifestPath, input)
			}
			if strings.HasSuffix(input, ".json") {
				phase11RequireValidJSONFile(t, path)
			}
		}

		if len(manifest.Assertions) == 0 {
			t.Fatalf("%s assertions must not be empty", manifestPath)
		}
		phase11RequireUnique(t, manifestPath, "assertions", manifest.Assertions)
		for _, assertion := range manifest.Assertions {
			if !phase11AllowedAssertions[assertion] {
				t.Fatalf("%s has unsupported assertion %q", manifestPath, assertion)
			}
			phase11RequireAssertionEvidence(t, fixtureDir, manifestPath, assertion)
		}
	}
}

func TestPhase11DocumentReferenceGate(t *testing.T) {
	t.Parallel()

	forbidden := []string{
		"docs/details/phase11.md",
		"phase11.md",
		"Phase 11 仕様詳細ファイル",
		"Part 3",
		"Part3",
		"第3部",
	}

	markdownFiles, err := phase11FindMarkdownFiles()
	if err != nil {
		t.Fatalf("find markdown files: %v", err)
	}
	for _, path := range markdownFiles {
		body := phase11MustReadText(t, path)
		for _, term := range forbidden {
			if strings.Contains(body, term) {
				t.Fatalf("%s contains forbidden stale specification term %q", path, term)
			}
		}
	}
}

func TestPhase11RoadmapStateGate(t *testing.T) {
	t.Parallel()

	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	if !strings.Contains(roadmap, "| Phase 11 | バグ修正ゼロ化。") || !strings.Contains(roadmap, "| Phase 11 | バグ修正ゼロ化。source-code audit、横断 regression、正式 fixture harness、意味のあるテスト、test gap inventory / batch closure、race trigger、mutation selection / mutation zero survivor、contract drift の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) を参照する。 | 実装済み | Phase 10 |") {
		t.Fatalf("docs/ROADMAP.md must mark Phase 11 as 実装済み once Phase 11 quality gates are implemented")
	}
	for _, feature := range []string{
		"バグ修正ゼロ化 / 仕様全般完了 gate / cross-owner regression gate",
		"全標準実装 artifact source-code audit / artifact coverage zero gap",
		"意味のあるテスト / test gap inventory / batch closure / race trigger / mutation selection / mutation zero survivor gate",
		"仕様全般不備一括棚卸し / 重複ゼロ / inspection pass / spec-gap closure gate",
		"main dispatch / binary version / output manifest / file tree consistency gate",
		"source-code audit inventory / direct I/O elimination / statefile common persistence gate",
		"正式 fixture directory harness / fixture identity / fixture manifest / contract drift zero gate",
	} {
		want := "| 実装済み | 検証基盤 | " + feature + " |"
		if !strings.Contains(roadmap, want) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 11 feature as 実装済み: %s", feature)
		}
	}
}

type phase11Manifest struct {
	Name                   string                  `json:"name"`
	Section                string                  `json:"section"`
	Feature                string                  `json:"feature"`
	Category               string                  `json:"category"`
	OwnerComponent         string                  `json:"owner_component"`
	CollaboratorComponents []string                `json:"collaborator_components"`
	Components             []string                `json:"components"`
	References             []string                `json:"references"`
	FakeClock              *string                 `json:"fake_clock"`
	NotApplicable          []phase11ManifestRecord `json:"not_applicable"`
	MissingState           []string                `json:"missing_state"`
	InputFiles             []string                `json:"input_files"`
	Assertions             []string                `json:"assertions"`
}

type phase11ManifestRecord struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

var phase11AllowedCategories = map[string]bool{
	"failure":  true,
	"noop":     true,
	"partial":  true,
	"security": true,
	"success":  true,
}

var phase11AllowedComponents = map[string]bool{
	"admin":        true,
	"api":          true,
	"archive":      true,
	"builder":      true,
	"commitstatus": true,
	"mcp":          true,
	"release":      true,
	"runner":       true,
	"sdk":          true,
	"security":     true,
	"setup":        true,
	"statefile":    true,
	"ui":           true,
}

var phase11AllowedAssertions = map[string]bool{
	"effects":     true,
	"idempotency": true,
	"logs":        true,
	"no-write":    true,
	"order":       true,
	"request":     true,
	"response":    true,
	"secret-mask": true,
	"state":       true,
	"stderr":      true,
	"stdout":      true,
}

func phase11ReadManifest(t *testing.T, path string) phase11Manifest {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var manifest phase11Manifest
	if err := decoder.Decode(&manifest); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	if decoder.More() {
		t.Fatalf("%s must contain one JSON object", path)
	}
	return manifest
}

func phase11FindManifestFiles(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "manifest.json" {
			paths = append(paths, filepath.ToSlash(path))
		}
		return nil
	})
	sort.Strings(paths)
	return paths, err
}

func phase11FindMarkdownFiles() ([]string, error) {
	roots := []string{"AGENTS.md", "README.md", "docs"}
	var paths []string
	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			if strings.HasSuffix(root, ".md") {
				paths = append(paths, root)
			}
			continue
		}
		err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".md") {
				paths = append(paths, filepath.ToSlash(path))
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func phase11NameHasCategoryPrefix(name string) bool {
	for category := range phase11AllowedCategories {
		if strings.HasPrefix(name, category+"-") {
			return true
		}
	}
	return false
}

func phase11RequireMarkdownReference(t *testing.T, manifestPath string, reference string) {
	t.Helper()

	if strings.Contains(reference, "\\") || strings.Contains(reference, "../") || strings.HasPrefix(reference, "/") {
		t.Fatalf("%s reference %q must be a repository-relative Markdown link target", manifestPath, reference)
	}
	parts := strings.Split(reference, "#")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		t.Fatalf("%s reference %q must include exactly one file path and one fixed anchor", manifestPath, reference)
	}
	if !strings.HasPrefix(parts[0], "docs/") || !strings.HasSuffix(parts[0], ".md") {
		t.Fatalf("%s reference %q must point to a docs/*.md file", manifestPath, reference)
	}
	data := phase11MustReadText(t, parts[0])
	anchor := parts[1]
	if !strings.Contains(data, `id="`+anchor+`"`) && !strings.Contains(data, "id='"+anchor+"'") && !strings.Contains(data, "#"+anchor) {
		t.Fatalf("%s reference %q points to a missing anchor", manifestPath, reference)
	}
}

func phase11RequireManifestRecords(t *testing.T, manifestPath string, field string, records []phase11ManifestRecord) {
	t.Helper()

	seen := make(map[string]bool)
	for _, record := range records {
		if !phase11SafeRelativePath(record.Path) {
			t.Fatalf("%s %s path %q must be a safe relative path", manifestPath, field, record.Path)
		}
		if strings.TrimSpace(record.Reason) == "" {
			t.Fatalf("%s %s path %q must include a reason", manifestPath, field, record.Path)
		}
		if seen[record.Path] {
			t.Fatalf("%s %s repeats path %q", manifestPath, field, record.Path)
		}
		seen[record.Path] = true
	}
}

func phase11RequireSafeRelativePaths(t *testing.T, manifestPath string, field string, paths []string) {
	t.Helper()

	phase11RequireUnique(t, manifestPath, field, paths)
	for _, path := range paths {
		if !phase11SafeRelativePath(path) {
			t.Fatalf("%s %s path %q must be a safe relative path", manifestPath, field, path)
		}
	}
}

func phase11SafeRelativePath(path string) bool {
	if path == "" || strings.HasPrefix(path, "/") || strings.Contains(path, "\\") {
		return false
	}
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		return false
	}
	for _, segment := range strings.Split(path, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
		if strings.ContainsAny(segment, "\x00\r\n") {
			return false
		}
	}
	return true
}

func phase11RequireAssertionEvidence(t *testing.T, fixtureDir string, manifestPath string, assertion string) {
	t.Helper()

	assertionFiles := map[string]string{
		"effects":     "expected/effects.json",
		"request":     "expected/request.json",
		"response":    "expected/response.json",
		"secret-mask": "expected/security.json",
		"stderr":      "expected/stderr.txt",
		"stdout":      "expected/stdout.txt",
	}
	if rel, ok := assertionFiles[assertion]; ok {
		path := filepath.Join(fixtureDir, filepath.FromSlash(rel))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("%s assertion %q requires %s: %v", manifestPath, assertion, rel, err)
		}
		if info.IsDir() {
			t.Fatalf("%s assertion %q requires %s to be a file", manifestPath, assertion, rel)
		}
		if strings.HasSuffix(rel, ".json") {
			phase11RequireValidJSONFile(t, path)
		}
		return
	}

	assertionDirs := map[string]string{
		"logs":  "expected/logs",
		"state": "expected/state",
	}
	if rel, ok := assertionDirs[assertion]; ok {
		path := filepath.Join(fixtureDir, filepath.FromSlash(rel))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("%s assertion %q requires %s: %v", manifestPath, assertion, rel, err)
		}
		if !info.IsDir() {
			t.Fatalf("%s assertion %q requires %s to be a directory", manifestPath, assertion, rel)
		}
		return
	}
}

func phase11RequireValidJSONFile(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read JSON evidence %s: %v", path, err)
	}
	if !json.Valid(data) {
		t.Fatalf("JSON evidence %s is not valid JSON", path)
	}
}

func phase11RequireUnique(t *testing.T, context string, field string, values []string) {
	t.Helper()

	seen := make(map[string]bool)
	for _, value := range values {
		if value == "" {
			t.Fatalf("%s %s must not contain an empty value", context, field)
		}
		if seen[value] {
			t.Fatalf("%s %s repeats %q", context, field, value)
		}
		seen[value] = true
	}
}

func phase11RequireSortedUnique(t *testing.T, context string, field string, values []string) {
	t.Helper()

	phase11RequireUnique(t, context, field, values)
	if !sort.StringsAreSorted(values) {
		t.Fatalf("%s %s must be sorted: %v", context, field, values)
	}
}

func phase11StringSlicesEqual(a []string, b []string) bool {
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

func phase11MustReadText(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.HasSuffix(path, ".md") {
		return string(data)
	}
	if strings.Contains(string(data), "\x00") {
		t.Fatalf("%s must not contain NUL bytes", path)
	}
	return string(data)
}
