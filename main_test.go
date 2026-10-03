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

	if code := dispatchMain("adlaire-ci-obsidian", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected obsidian help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-obsidian") {
		t.Fatalf("expected obsidian help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for obsidian help, got %q", stderr.String())
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

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-obsidian-linux-amd64", []string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected release asset obsidian version exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.HasPrefix(got, "adlaire-ci-obsidian V.0.0-dev go=") {
		t.Fatalf("expected canonical obsidian version, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for release asset obsidian version, got %q", stderr.String())
	}
}

func TestPhase12StandardArtifactInventory(t *testing.T) {
	t.Parallel()

	ownerPackages := []string{
		"admin",
		"api",
		"archive",
		"builder",
		"commitstatus",
		"mcp",
		"obsidian",
		"release",
		"runner",
		"security",
		"setup",
		"statefile",
	}
	allowedOwnerFiles := map[string][]string{}
	for _, owner := range ownerPackages {
		allowedOwnerFiles[owner] = []string{
			owner + ".go",
			"model.go",
			"validate.go",
			"execute.go",
			owner + "_test.go",
		}
	}

	requiredFiles := []string{
		"main.go",
		"main_test.go",
		"sdk_contract_test.go",
		"ui_contract_test.go",
		"admin/adlaire-ci-sdk.js",
		"admin/index.html",
		".github/workflows/phase12-quality-gate.yml",
	}
	for _, owner := range ownerPackages {
		for _, file := range allowedOwnerFiles[owner] {
			requiredFiles = append(requiredFiles, filepath.ToSlash(filepath.Join("components", owner, file)))
		}
	}
	for _, path := range requiredFiles {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("required Phase 12 artifact %s is missing: %v", path, err)
		}
		if info.IsDir() {
			t.Fatalf("required Phase 12 artifact %s must be a file", path)
		}
	}

	if _, err := os.Stat("cmd"); !os.IsNotExist(err) {
		t.Fatalf("cmd directory must not exist in the standard Phase 12 layout: %v", err)
	}

	entries, err := os.ReadDir("components")
	if err != nil {
		t.Fatalf("read components directory: %v", err)
	}
	var ownerDirs []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			t.Fatalf("components root must not contain Go file after Phase 12 owner package split: %s", entry.Name())
		}
		if entry.IsDir() {
			ownerDirs = append(ownerDirs, entry.Name())
		}
	}
	sort.Strings(ownerDirs)
	if !phase11StringSlicesEqual(ownerDirs, ownerPackages) {
		t.Fatalf("unexpected components owner directories\nwant: %v\n got: %v", ownerPackages, ownerDirs)
	}

	for _, owner := range ownerPackages {
		entries, err := os.ReadDir(filepath.Join("components", owner))
		if err != nil {
			t.Fatalf("read owner package %s: %v", owner, err)
		}
		var goFiles []string
		for _, entry := range entries {
			if entry.IsDir() {
				t.Fatalf("owner package %s must not contain subdirectory %s", owner, entry.Name())
			}
			if strings.HasSuffix(entry.Name(), ".go") {
				goFiles = append(goFiles, entry.Name())
			}
		}
		sort.Strings(goFiles)
		want := append([]string(nil), allowedOwnerFiles[owner]...)
		sort.Strings(want)
		if !phase11StringSlicesEqual(goFiles, want) {
			t.Fatalf("owner package %s must use exactly five fixed files\nwant: %v\n got: %v", owner, want, goFiles)
		}
	}

	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, path := range requiredFiles {
		if !strings.Contains(documentIndex, path) {
			t.Fatalf("docs/DOCUMENT_INDEX.md does not index required artifact %s", path)
		}
	}
}

func TestPhase12QualityGateEvidence(t *testing.T) {
	t.Parallel()

	workflow := phase11MustReadText(t, ".github/workflows/phase12-quality-gate.yml")
	for _, job := range []string{
		"phase12-go-format:",
		"phase12-go-test:",
		"phase12-go-race:",
		"phase12-go-vet:",
		"phase12-deno-check-sdk:",
		"phase12-fixture-harness:",
		"phase12-mutation:",
	} {
		if !strings.Contains(workflow, job) {
			t.Fatalf("Phase 12 workflow must include required job %s", strings.TrimSuffix(job, ":"))
		}
	}

	effectsPath := "testdata/phase12/quality-gate-reconstruction/expected/effects.json"
	var effects map[string]int
	data, err := os.ReadFile(effectsPath)
	if err != nil {
		t.Fatalf("read Phase 12 effects: %v", err)
	}
	if err := json.Unmarshal(data, &effects); err != nil {
		t.Fatalf("decode Phase 12 effects: %v", err)
	}
	for _, key := range []string{
		"phase12_contract_mismatch_count",
		"phase12_state_safety_open_count",
		"phase12_fixture_harness_open_count",
		"phase12_mutation_survived_count",
		"phase12_race_trigger_open_count",
		"phase12_owner_package_violation_count",
		"phase12_ci_required_check_open_count",
		"phase12_release_reproducibility_open_count",
		"final_open_item_count",
	} {
		value, ok := effects[key]
		if !ok {
			t.Fatalf("Phase 12 effects must include %s", key)
		}
		if value != 0 {
			t.Fatalf("Phase 12 effects %s must be 0, got %d", key, value)
		}
	}
}

func TestPhase11FixtureManifestGate(t *testing.T) {
	t.Parallel()

	phase11RequireNoUnexpectedFixtureDirectories(t)

	manifestPaths, err := phase11FindManifestFiles("testdata")
	if err != nil {
		t.Fatalf("find fixture manifests: %v", err)
	}
	if len(manifestPaths) < 56 {
		t.Fatalf("expected at least 56 formal fixture manifests, got %d", len(manifestPaths))
	}

	seenNames := make(map[string]string)
	for _, manifestPath := range manifestPaths {
		manifest := phase11ReadManifest(t, manifestPath)
		fixtureDir := filepath.Dir(manifestPath)

		if manifest.Name != filepath.Base(fixtureDir) && !phase13AllowsScopeName(manifestPath, manifest.Name) && !phase14AllowsScopeName(manifestPath, manifest.Name) && !phase15AllowsScopeName(manifestPath, manifest.Name) {
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
		"後続証跡として残る",
		"後続 Phase または横断 owner",
		"後続 Phase の現在状態",
		"formal fixture harness の後続証跡",
		"残証跡",
		"未完了である",
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
	fixtureSpec := phase11MustReadText(t, "docs/details/fixture.md")
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	phase11CompleteRow := "| Phase 11 | バグ修正ゼロ化。source-code audit、横断 regression、正式 fixture harness、意味のあるテスト、test gap inventory / batch closure、race trigger、mutation selection / mutation zero survivor、contract drift の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) を参照する。 | 実装済み | Phase 10 |"
	phase11IncompleteRow := strings.Replace(phase11CompleteRow, "| 実装済み |", "| 実装中・検証未完了 |", 1)
	phase12CompleteRow := "| Phase 12 | 実装品質ゲート再構築。契約不整合、状態安全性、実行型 fixture harness、mutation / race、owner package 5 ファイル固定、release 再現性の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](DETAIL_INDEX.md#phase-12-quality-gate-entry) を参照する。 | 実装済み | Phase 11 |"
	phase12IncompleteRow := strings.Replace(phase12CompleteRow, "| 実装済み |", "| 実装中・検証未完了 |", 1)
	phase13CompleteRow := "| Phase 13 | 実装整合・品質改善。owner package 5 ファイル責務純度、状態安全性、security / archive / commitstatus 責務集約、queue / finalizer / recovery、MCP 実動作化、外部境界 hardening、fixture / mutation / fault / E2E / CI / release governance の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) を参照する。 | 実装済み | Phase 12 |"

	if !strings.Contains(roadmap, phase11CompleteRow) {
		t.Fatalf("docs/ROADMAP.md must mark Phase 11 as 実装済み after Phase 11 closure records reach final_open_item_count=0")
	}
	if strings.Contains(roadmap, phase11IncompleteRow) {
		t.Fatalf("docs/ROADMAP.md must not keep Phase 11 as 実装中・検証未完了 after closure")
	}
	if !strings.Contains(roadmap, phase12CompleteRow) {
		t.Fatalf("docs/ROADMAP.md must mark Phase 12 as 実装済み after Phase 12 closure records reach final_open_item_count=0")
	}
	if strings.Contains(roadmap, phase12IncompleteRow) {
		t.Fatalf("docs/ROADMAP.md must not keep Phase 12 as 実装中・検証未完了 after closure")
	}
	if !strings.Contains(roadmap, phase13CompleteRow) {
		t.Fatalf("docs/ROADMAP.md must define Phase 13 as 実装済み after Phase 13 closure")
	}
	if !strings.Contains(roadmap, "現在の active Phase は Phase 16 とする。") {
		t.Fatalf("docs/ROADMAP.md must state that Phase 16 is the active Phase after Phase 15 closure")
	}
	if !strings.Contains(roadmap, "初期実装 Phase 1 から Phase 15 まではすべて `実装済み`") {
		t.Fatalf("docs/ROADMAP.md must state that Phase 1 through Phase 15 are all implemented after Phase 15 closure")
	}

	for _, feature := range phase11QualityGateFeatures {
		want := "| 実装済み | 検証基盤 | " + feature + " |"
		if !strings.Contains(roadmap, want) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 11 feature as 実装済み: %s", feature)
		}
		forbidden := "| 実装中・検証未完了 | 検証基盤 | " + feature + " |"
		if strings.Contains(roadmap, forbidden) {
			t.Fatalf("docs/ROADMAP.md must not keep Phase 11 feature as 実装中・検証未完了 after closure: %s", feature)
		}
	}
	for _, feature := range phase12QualityGateFeatures {
		want := "| 実装済み | 検証基盤 | " + feature + " |"
		if strings.Contains(feature, "owner package") {
			want = "| 実装済み | 実装構造 | " + feature + " |"
		}
		if strings.Contains(feature, "version tag") {
			want = "| 実装済み | 配布・リリース | " + feature + " |"
		}
		if !strings.Contains(roadmap, want) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 12 feature as 実装済み: %s", feature)
		}
		forbidden := strings.Replace(want, "| 実装済み |", "| 実装中・検証未完了 |", 1)
		if strings.Contains(roadmap, forbidden) {
			t.Fatalf("docs/ROADMAP.md must not keep Phase 12 feature as 実装中・検証未完了 after closure: %s", feature)
		}
	}
	for _, feature := range phase13QualityGateFeatures {
		want := "| 実装済み | 検証基盤 | " + feature + " |"
		if strings.Contains(feature, "owner package") {
			want = "| 実装済み | 実装構造 | " + feature + " |"
		}
		if strings.Contains(feature, "statefile") {
			want = "| 実装済み | 状態管理 | " + feature + " |"
		}
		if strings.Contains(feature, "credential") {
			want = "| 実装済み | セキュリティ | " + feature + " |"
		}
		if strings.Contains(feature, "API route") {
			want = "| 実装済み | 管理ツール・API | " + feature + " |"
		}
		if strings.Contains(feature, "queue state machine") {
			want = "| 実装済み | CI ランナー | " + feature + " |"
		}
		if strings.Contains(feature, "MCP resendWebhook") {
			want = "| 実装済み | MCP サーバー | " + feature + " |"
		}
		if strings.Contains(feature, "HTTP timeout") {
			want = "| 実装済み | 外部境界 | " + feature + " |"
		}
		if strings.Contains(feature, "Git tag") {
			want = "| 実装済み | 配布・リリース | " + feature + " |"
		}
		if strings.Contains(feature, "旧 path") {
			want = "| 実装済み | 文書整合 | " + feature + " |"
		}
		if !strings.Contains(roadmap, want) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 13 feature as 実装済み: %s", feature)
		}
		forbidden := strings.Replace(want, "| 実装済み |", "| 実装中・検証未完了 |", 1)
		if strings.Contains(roadmap, forbidden) {
			t.Fatalf("docs/ROADMAP.md must not keep Phase 13 feature as 実装中・検証未完了 after closure: %s", feature)
		}
	}

	for _, fixture := range phase11BuilderFormalFixtures {
		indexRow := "| [`" + fixture + "/`](../" + fixture + "/) | `builder` fixture | 実在 |"
		if !strings.Contains(documentIndex, indexRow) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index builder formal fixture directory as 実在: %s", fixture)
		}
		if _, err := os.Stat(filepath.Join(fixture, "manifest.json")); err != nil {
			t.Fatalf("builder formal fixture %s must include manifest.json: %v", fixture, err)
		}
		if _, err := os.Stat(filepath.Join(fixture, "expected")); err != nil {
			t.Fatalf("builder formal fixture %s must include expected/: %v", fixture, err)
		}
	}

	for _, coverage := range phase11RootCoverageFixtures {
		rootIndexRow := "| [`" + coverage.Root + "`](../" + coverage.Root + ") | `" + coverage.Owner + "` fixture root | 実在 |"
		if !strings.Contains(documentIndex, rootIndexRow) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index Phase 11 root coverage fixture root as 実在: %s", coverage.Root)
		}
		fixtureIndexRow := "| [`" + coverage.FixtureDir + "/`](../" + coverage.FixtureDir + "/) | `" + coverage.Owner + "` Phase 11 root coverage fixture | 実在 |"
		if !strings.Contains(documentIndex, fixtureIndexRow) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index Phase 11 root coverage fixture directory as 実在: %s", coverage.FixtureDir)
		}
		if _, err := os.Stat(filepath.Join(coverage.FixtureDir, "manifest.json")); err != nil {
			t.Fatalf("Phase 11 root coverage fixture %s must include manifest.json: %v", coverage.FixtureDir, err)
		}
		if !strings.Contains(fixtureSpec, "`"+coverage.Root+"` | Phase 11 root coverage formal fixture root。") {
			t.Fatalf("docs/details/fixture.md must record Phase 11 root coverage closure for %s", coverage.Root)
		}
		if strings.Contains(fixtureSpec, "`"+coverage.Root+"` | Phase 11 対象の"+"未作成 root。") {
			t.Fatalf("docs/details/fixture.md must not keep %s as 未作成 root after root coverage closure", coverage.Root)
		}
	}
	if strings.Contains(fixtureSpec, "`p11-"+"open-"+"fixture-root-coverage`") {
		t.Fatalf("docs/details/fixture.md must not keep closed root coverage as an open Phase 11 blocker")
	}
	if strings.Contains(fixtureSpec, "`p11-"+"open-") {
		t.Fatalf("docs/details/fixture.md must not keep any open Phase 11 blocker ids after closure")
	}
	if !strings.Contains(fixtureSpec, "| `final_open_item_count` | `closed` | `0` | `[]` |") {
		t.Fatalf("docs/details/fixture.md must record Phase 11 final_open_item_count=0")
	}
	if !strings.Contains(fixtureSpec, "`survived=0`") {
		t.Fatalf("docs/details/fixture.md must record Phase 11 mutation zero survivor evidence")
	}
	if !strings.Contains(roadmap, "状態ファイル共通永続化契約 | [`docs/DETAIL_INDEX.md`") || !strings.Contains(roadmap, "| 実装済み | 状態管理 | 状態ファイル共通永続化契約 |") {
		t.Fatalf("docs/ROADMAP.md must mark the statefile common persistence contract implemented after verification")
	}
}

func TestPhase14ObsidianVaultIntegrationEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase14/obsidian-vault-integration"
	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.Name != "phase-14-obsidian-vault-integration" {
		t.Fatalf("%s name mismatch: %q", manifestPath, manifest.Name)
	}
	if !phase11StringSlicesEqual(manifest.Owners, []string{"obsidian"}) {
		t.Fatalf("%s owners mismatch: %v", manifestPath, manifest.Owners)
	}
	if !phase11StringSlicesEqual(manifest.CollaboratorComponents, []string{"builder", "security"}) {
		t.Fatalf("%s collaborators mismatch: %v", manifestPath, manifest.CollaboratorComponents)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, []string{"phase14-go-format", "phase14-go-test", "phase14-obsidian-vault-fixture", "phase14-builder-handoff", "phase14-document-drift"}) {
		t.Fatalf("%s required checks mismatch: %v", manifestPath, manifest.RequiredChecks)
	}
	for _, rel := range []string{
		"input/options.json",
		"input/vault_tree.json",
		"expected/normalized.json",
		"expected/output_tree.json",
		"expected/counters.json",
		"records/closure.jsonl",
	} {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("Phase 14 fixture path missing %s: %v", rel, err)
		}
	}
	counters := phase13ReadCounterFile(t, filepath.Join(root, "expected/counters.json"))
	for _, key := range []string{
		"phase14_obsidian_vault_boundary_open_count",
		"phase14_obsidian_wikilink_open_count",
		"phase14_obsidian_yaml_rejection_open_count",
		"phase14_obsidian_asset_open_count",
		"phase14_obsidian_builder_handoff_open_count",
		"phase14_fixture_execution_gap_count",
		"final_open_item_count",
	} {
		if counters[key] != 0 {
			t.Fatalf("%s must be 0, got %d", key, counters[key])
		}
	}
	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	if !strings.Contains(roadmap, "| Phase 14 | Obsidian Vault 連携。local vault 読取、wikilink / embed / tag / asset 正規化、YAML frontmatter 拒否、builder handoff の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 14 Obsidian Vault 連携参照](DETAIL_INDEX.md#phase-14-obsidian-vault-integration-entry) を参照する。 | 実装済み | Phase 13 |") {
		t.Fatalf("docs/ROADMAP.md must mark Phase 14 implemented")
	}
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, expected := range []string{
		"| [`components/obsidian/`](../components/obsidian/) | `obsidian` owner package directory | 実在 |",
		"| [`testdata/phase14/obsidian-vault-integration/`](../testdata/phase14/obsidian-vault-integration/) | Phase 14 正式 fixture root | 実在 |",
		"| [`.github/workflows/phase14-obsidian-vault-integration.yml`](../.github/workflows/phase14-obsidian-vault-integration.yml) | Phase 14 required check workflow | 実在 |",
	} {
		if !strings.Contains(documentIndex, expected) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must contain Phase 14実在 row: %s", expected)
		}
	}
}

func TestPhase14ObsidianVaultProductionEntrypoint(t *testing.T) {
	root := t.TempDir()
	vault := filepath.Join(root, "vault")
	phase14WriteFile(t, filepath.Join(vault, "Home.md"), "# Home\n\nSee [[Guide.md|Guide]].\n\n![[images/logo.png|Logo]]\n")
	phase14WriteFile(t, filepath.Join(vault, "Guide.md"), "# Guide\n")
	phase14WriteFile(t, filepath.Join(vault, "images", "logo.png"), "png")
	out := filepath.Join(root, "site")
	report := filepath.Join(out, "obsidian_map.json")
	var stdout, stderr bytes.Buffer
	code := dispatchMain("adlaire-ci-build", []string{
		"--input-mode", "obsidian-vault",
		"--obsidian-vault", vault,
		"--obsidian-entry", "Home.md",
		"--obsidian-report-file", report,
		"--out", out,
		"--title", "Vault",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr must be empty: %s", stderr.String())
	}
	for _, rel := range []string{"index.html", "assets/images/logo.png", "obsidian_map.json"} {
		if _, err := os.Stat(filepath.Join(out, rel)); err != nil {
			t.Fatalf("missing production output %s: %v", rel, err)
		}
	}
	if !strings.Contains(stdout.String(), "obsidian_notes=2") || !strings.Contains(stdout.String(), "obsidian_assets=1") {
		t.Fatalf("production entrypoint missing obsidian report: %s", stdout.String())
	}
}

func TestPhase15ObsidianLocalSyncProductionEntrypoint(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project")
	vault := filepath.Join(root, "vault")
	state := filepath.Join(root, "state")
	phase14WriteFile(t, filepath.Join(project, "docs", "project.md"), "# Project\n")
	if err := os.MkdirAll(vault, 0o755); err != nil {
		t.Fatalf("mkdir vault: %v", err)
	}
	if err := os.MkdirAll(state, 0o755); err != nil {
		t.Fatalf("mkdir state: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := dispatchMain("adlaire-ci-obsidian", []string{
		"sync", "plan",
		"--project-root", project,
		"--vault", vault,
		"--state-dir", state,
		"--plan-file", "plans/sync.json",
		"--direction", "export-only",
		"--delete-policy", "tombstone",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr must be empty: %s", stderr.String())
	}
	var response struct {
		Command    string `json:"command"`
		PlanFile   string `json:"plan_file"`
		PlanHash   string `json:"plan_hash"`
		Operations int    `json:"operations"`
		Conflicts  int    `json:"conflicts"`
		Tombstones int    `json:"tombstones"`
		Applied    bool   `json:"applied"`
	}
	decoder := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("decode production sync response: %v\n%s", err, stdout.String())
	}
	if response.Command != "sync plan" || response.PlanFile != "plans/sync.json" || response.PlanHash == "" || response.Operations != 1 || response.Applied {
		t.Fatalf("unexpected production sync response: %+v", response)
	}
	if _, err := os.Stat(filepath.Join(state, "plans", "sync.json")); err != nil {
		t.Fatalf("production entrypoint must write plan file: %v", err)
	}
}

func phase14WriteFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestPhase13ImplementationAlignmentQualityEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase13/implementation-alignment-quality"
	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)

	if manifest.Name != "phase-13-implementation-alignment-quality" {
		t.Fatalf("%s name must be phase-13-implementation-alignment-quality, got %q", manifestPath, manifest.Name)
	}
	if manifest.Scope != "phase-13-implementation-alignment-quality" {
		t.Fatalf("%s scope must be phase-13-implementation-alignment-quality, got %q", manifestPath, manifest.Scope)
	}
	if !phase11StringSlicesEqual(manifest.Owners, phase13Owners) {
		t.Fatalf("%s owners mismatch\nwant: %v\n got: %v", manifestPath, phase13Owners, manifest.Owners)
	}
	if !phase11StringSlicesEqual(manifest.Entrypoints, phase13Entrypoints) {
		t.Fatalf("%s entrypoints mismatch\nwant: %v\n got: %v", manifestPath, phase13Entrypoints, manifest.Entrypoints)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase13RequiredChecks) {
		t.Fatalf("%s required_checks mismatch\nwant: %v\n got: %v", manifestPath, phase13RequiredChecks, manifest.RequiredChecks)
	}
	if !phase11StringSlicesEqual(manifest.ClosureRecords, phase13ClosureRecords) {
		t.Fatalf("%s closure_records mismatch\nwant: %v\n got: %v", manifestPath, phase13ClosureRecords, manifest.ClosureRecords)
	}
	if !phase11StringSlicesEqual(manifest.NegativeControls, phase13NegativeControls) {
		t.Fatalf("%s negative_controls mismatch\nwant: %v\n got: %v", manifestPath, phase13NegativeControls, manifest.NegativeControls)
	}
	for _, sourceAnchor := range manifest.SourceAnchors {
		phase11RequireMarkdownReference(t, manifestPath, sourceAnchor)
	}

	phase13RequireOwnerPackages(t)
	phase13RequireDocumentIndexRows(t)
	phase13RequireGovernanceFiles(t)
	phase13RequireZeroExternalDependencies(t)

	effects := phase13ReadCounterFile(t, filepath.Join(root, "expected", "effects.json"))
	counters := phase13ReadCounterFile(t, filepath.Join(root, "expected", "counters.json"))
	if len(counters) != len(phase13ClosureCounters) {
		t.Fatalf("expected/counters.json must include exactly %d counters, got %d", len(phase13ClosureCounters), len(counters))
	}
	for _, key := range phase13ClosureCounters {
		counter, ok := counters[key]
		if !ok {
			t.Fatalf("expected/counters.json must include %s", key)
		}
		if counter != 0 {
			t.Fatalf("expected/counters.json %s must be 0, got %d", key, counter)
		}
		effect, ok := effects[key]
		if !ok {
			t.Fatalf("expected/effects.json must include %s", key)
		}
		if effect != counter {
			t.Fatalf("expected/effects.json %s must equal counters value %d, got %d", key, counter, effect)
		}
	}

	closureRecords := phase13ReadClosureRecords(t, filepath.Join(root, "records", "closure.jsonl"))
	if len(closureRecords) != 18 {
		t.Fatalf("records/closure.jsonl must include 18 closure records, got %d", len(closureRecords))
	}
	for _, record := range closureRecords {
		if record.Scope != "phase-13-implementation-alignment-quality" {
			t.Fatalf("records/closure.jsonl record %s has unexpected scope %q", record.RecordID, record.Scope)
		}
		if record.Status != "closed" {
			t.Fatalf("records/closure.jsonl record %s must be closed, got %q", record.RecordID, record.Status)
		}
		if len(record.OpenItems) != 0 {
			t.Fatalf("records/closure.jsonl record %s must have no open_items, got %v", record.RecordID, record.OpenItems)
		}
	}
	for _, path := range []string{"records/mutation.jsonl", "records/race.jsonl", "records/fault.jsonl", "records/e2e.jsonl", "records/release.jsonl"} {
		if len(phase13ReadJSONLines(t, filepath.Join(root, filepath.FromSlash(path)))) == 0 {
			t.Fatalf("%s must contain at least one JSONL record", path)
		}
	}

	workflow := phase11MustReadText(t, ".github/workflows/phase13-implementation-alignment-quality.yml")
	if !strings.Contains(workflow, "permissions:\n  contents: read") {
		t.Fatalf("Phase 13 workflow must use minimum contents:read permissions")
	}
	for _, requiredCheck := range phase13RequiredChecks {
		if !strings.Contains(workflow, "\n  "+requiredCheck+":") {
			t.Fatalf("Phase 13 workflow must define required check job %s", requiredCheck)
		}
	}
	if strings.Count(workflow, "timeout-minutes:") < len(phase13RequiredChecks) {
		t.Fatalf("Phase 13 workflow must set timeout-minutes on every required check")
	}
	phase13RequirePinnedActions(t, workflow)
}

func TestPhase15ObsidianLocalSyncEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase15/obsidian-local-sync"
	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.Name != "phase-15-obsidian-local-sync" {
		t.Fatalf("%s name mismatch: %q", manifestPath, manifest.Name)
	}
	if !phase11StringSlicesEqual(manifest.Owners, []string{"obsidian"}) {
		t.Fatalf("%s owners mismatch: %v", manifestPath, manifest.Owners)
	}
	if !phase11StringSlicesEqual(manifest.CollaboratorComponents, []string{"release", "security", "setup", "statefile"}) {
		t.Fatalf("%s collaborators mismatch: %v", manifestPath, manifest.CollaboratorComponents)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, []string{"phase15-go-format", "phase15-go-test", "phase15-obsidian-sync-fixture", "phase15-sync-atomicity", "phase15-release-setup-distribution", "phase15-document-drift"}) {
		t.Fatalf("%s required checks mismatch: %v", manifestPath, manifest.RequiredChecks)
	}
	for _, rel := range []string{
		"input/options.json",
		"input/sync_state.json",
		"input/project_tree.json",
		"input/vault_tree.json",
		"expected/plan.json",
		"expected/apply_state.json",
		"expected/distribution.json",
		"expected/counters.json",
		"records/closure.jsonl",
	} {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("Phase 15 fixture path missing %s: %v", rel, err)
		}
	}
	counters := phase13ReadCounterFile(t, filepath.Join(root, "expected/counters.json"))
	for _, key := range []string{
		"phase15_sync_plan_mismatch_count",
		"phase15_sync_apply_atomicity_open_count",
		"phase15_sync_conflict_open_count",
		"phase15_sync_tombstone_open_count",
		"phase15_sync_rollback_open_count",
		"phase15_sync_service_dependency_open_count",
		"phase15_distribution_open_count",
		"phase15_fixture_execution_gap_count",
		"final_open_item_count",
	} {
		if counters[key] != 0 {
			t.Fatalf("%s must be 0, got %d", key, counters[key])
		}
	}
	distribution := phase11MustReadText(t, filepath.Join(root, "expected/distribution.json"))
	for _, required := range []string{`"binary_asset":"adlaire-ci-obsidian-linux-amd64"`, `"release_asset_count":9`, `"checksum_line_count":8`, `"setup_mode":"install-obsidian"`} {
		if !strings.Contains(distribution, required) {
			t.Fatalf("expected distribution fixture to contain %s", required)
		}
	}
	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	if !strings.Contains(roadmap, "| Phase 15 | Obsidian local vault 同期。plan / apply / rollback、conflict、tombstone、state schema、atomic write、`adlaire-ci-obsidian` 配布連携の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 15 Obsidian local vault 同期参照](DETAIL_INDEX.md#phase-15-obsidian-local-sync-entry) を参照する。 | 実装済み | Phase 14 |") {
		t.Fatalf("docs/ROADMAP.md must mark Phase 15 implemented")
	}
	if !strings.Contains(roadmap, "| 実装済み | Obsidian 同期 | Phase 15 Obsidian local vault 同期 / plan / apply / rollback / conflict / tombstone / atomic write / `adlaire-ci-obsidian` 配布連携 gate |") {
		t.Fatalf("docs/ROADMAP.md must mark the Phase 15 feature implemented")
	}
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, expected := range []string{
		"| [`testdata/phase15/obsidian-local-sync/`](../testdata/phase15/obsidian-local-sync/) | Phase 15 Obsidian local vault 同期 fixture root | 実在 |",
		"| [`.github/workflows/phase15-obsidian-local-sync.yml`](../.github/workflows/phase15-obsidian-local-sync.yml) | Phase 15 required check workflow | 実在 |",
	} {
		if !strings.Contains(documentIndex, expected) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must contain Phase 15 実在 row: %s", expected)
		}
	}
	workflow := phase11MustReadText(t, ".github/workflows/phase15-obsidian-local-sync.yml")
	for _, requiredCheck := range manifest.RequiredChecks {
		if !strings.Contains(workflow, "\n  "+requiredCheck+":") {
			t.Fatalf("Phase 15 workflow must define required check job %s", requiredCheck)
		}
	}
	phase13RequirePinnedActions(t, workflow)
}

type phase11Manifest struct {
	Name                   string                  `json:"name"`
	Scope                  string                  `json:"scope,omitempty"`
	Section                string                  `json:"section"`
	Feature                string                  `json:"feature"`
	Category               string                  `json:"category"`
	OwnerComponent         string                  `json:"owner_component"`
	CollaboratorComponents []string                `json:"collaborator_components"`
	Components             []string                `json:"components"`
	Owners                 []string                `json:"owners,omitempty"`
	Entrypoints            []string                `json:"entrypoints,omitempty"`
	RequiredChecks         []string                `json:"required_checks,omitempty"`
	ClosureRecords         []string                `json:"closure_records,omitempty"`
	NegativeControls       []string                `json:"negative_controls,omitempty"`
	SourceAnchors          []string                `json:"source_anchors,omitempty"`
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
	"obsidian": true,
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
	"obsidian":     true,
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

var phase11QualityGateFeatures = []string{
	"バグ修正ゼロ化 / 仕様全般完了 gate / cross-owner regression gate",
	"全標準実装 artifact source-code audit / artifact coverage zero gap",
	"意味のあるテスト / test gap inventory / batch closure / race trigger / mutation selection / mutation zero survivor gate",
	"仕様全般不備一括棚卸し / 重複ゼロ / inspection pass / spec-gap closure gate",
	"main dispatch / binary version / output manifest / file tree consistency gate",
	"source-code audit inventory / direct I/O elimination / statefile common persistence gate",
	"builder parser / config / output / generated site regression gate",
	"runner source / pipeline / deploy / notification / corrupt state regression gate",
	"API auth / config / backup / webhook / read model / filesystem regression gate",
	"SDK / UI / admin CLI / MCP bridge / route parity regression gate",
	"setup / release distribution / filesystem / external boundary regression gate",
	"正式 fixture directory harness / fixture identity / fixture manifest / contract drift zero gate",
}

var phase12QualityGateFeatures = []string{
	"Phase 12 契約不整合・状態安全性 gate",
	"Phase 12 実行型 fixture harness / dependency injection / production entrypoint / expected 比較 gate",
	"Phase 12 mutation / race / queue state machine / CI required checks gate",
	"Phase 12 owner package 5 ファイル固定 / 巨大コンポーネント分割 gate",
	"Phase 12 version tag / release notes / release reproducibility gate",
}

var phase13QualityGateFeatures = []string{
	"Phase 13 owner package 5 ファイル責務純度 / 6 ファイル目・サブディレクトリ・空ファイル・ダミー実装ゼロ gate",
	"Phase 13 statefile 直接更新排除 / process lock / atomic write / JSON Lines 破損検出・隔離・復旧 gate",
	"Phase 13 credential / token / Go 標準ライブラリ内製 KDF / secret mask / file safety / required log write gate",
	"Phase 13 API route / Admin / SDK / UI client binding / MCP bridge / setup stdout / credential 初期化契約統一 gate",
	"Phase 13 queue state machine / finalizer / active recovery / at-least-once / backup-restore transaction gate",
	"Phase 13 MCP resendWebhook / subscribe / unsubscribe / sampling 実動作化と未実装 error gate",
	"Phase 13 HTTP timeout / body 上限 / graceful shutdown / Webhook SSRF 防止 / SSH strict / systemd 最小権限 gate",
	"Phase 13 fixture execution / mutation / race / fault injection / integration / E2E / CI required check gate",
	"Phase 13 Git tag / GitHub Release / SHA256SUMS / signature / SBOM / reproducible build evidence / recovery procedure gate",
	"Phase 13 旧 path / 重複仕様 / 実装済み表記 drift ゼロ gate",
}

var phase13Owners = []string{
	"admin",
	"api",
	"archive",
	"builder",
	"commitstatus",
	"mcp",
	"release",
	"runner",
	"sdk",
	"security",
	"setup",
	"statefile",
	"ui",
}

var phase13GoOwners = []string{
	"admin",
	"api",
	"archive",
	"builder",
	"commitstatus",
	"mcp",
	"release",
	"runner",
	"security",
	"setup",
	"statefile",
}

var phase13Entrypoints = []string{
	"adlaire-ci-build",
	"adlaire-ci-runner",
	"adlaire-ci-api",
	"adlaire-ci-admin",
	"adlaire-ci-setup",
	"adlaire-ci-release",
	"adlaire-ci-mcp",
	"admin-ui-browser",
	"phase13-fixture-harness",
}

var phase13RequiredChecks = []string{
	"phase13-go-format",
	"phase13-go-test",
	"phase13-go-race",
	"phase13-go-vet",
	"phase13-stdlib-static-analysis",
	"phase13-dependency-inventory",
	"phase13-secret-boundary-scan",
	"phase13-deno-check-sdk",
	"phase13-owner-shape",
	"phase13-contract-parity",
	"phase13-state-safety",
	"phase13-security-boundary",
	"phase13-archive-commitstatus",
	"phase13-runner-recovery",
	"phase13-mcp-real-behavior",
	"phase13-external-boundary",
	"phase13-executable-fixture",
	"phase13-mutation",
	"phase13-concurrency",
	"phase13-fault-injection",
	"phase13-browser",
	"phase13-integration-e2e",
	"phase13-setup-e2e",
	"phase13-release-e2e",
	"phase13-actions-pinning",
	"phase13-recovery-procedure",
	"phase13-document-drift",
}

var phase13ClosureRecords = []string{
	"records/closure.jsonl",
	"records/mutation.jsonl",
	"records/race.jsonl",
	"records/fault.jsonl",
	"records/e2e.jsonl",
	"records/release.jsonl",
}

var phase13NegativeControls = []string{
	"contract-mismatch",
	"direct-state-mutation",
	"github-actions-unpinned",
	"jsonl-corruption",
	"mutation-survivor",
	"race-trigger",
	"fault-injection-failure",
	"token-argv",
}

var phase13ClosureCounters = []string{
	"phase13_owner_file_violation_count",
	"phase13_dummy_or_empty_file_count",
	"phase13_direct_state_mutation_count",
	"phase13_contract_mismatch_count",
	"phase13_state_safety_open_count",
	"phase13_jsonl_corruption_open_count",
	"phase13_security_kdf_open_count",
	"phase13_token_arg_open_count",
	"phase13_archive_commitstatus_ownership_open_count",
	"phase13_queue_recovery_open_count",
	"phase13_mcp_unimplemented_success_count",
	"phase13_external_boundary_open_count",
	"phase13_required_log_write_ignore_count",
	"phase13_fixture_execution_gap_count",
	"phase13_mutation_survived_count",
	"phase13_race_or_concurrency_open_count",
	"phase13_fault_injection_open_count",
	"phase13_e2e_open_count",
	"phase13_ci_required_check_open_count",
	"phase13_action_pin_open_count",
	"phase13_release_evidence_open_count",
	"phase13_recovery_procedure_open_count",
	"phase13_document_drift_open_count",
	"final_open_item_count",
}

var phase11BuilderFormalFixtures = []string{
	"testdata/builder/empty-dir",
	"testdata/builder/safe",
	"testdata/builder/single",
	"testdata/builder/site",
	"testdata/builder/strict",
	"testdata/builder/url-safety",
}

type phase11RootCoverageFixture struct {
	Root       string
	Owner      string
	FixtureDir string
}

var phase11RootCoverageFixtures = []phase11RootCoverageFixture{
	{Root: "testdata/runner/", Owner: "runner", FixtureDir: "testdata/runner/success-runner-phase11-root-coverage"},
	{Root: "testdata/api/", Owner: "api", FixtureDir: "testdata/api/success-api-phase11-root-coverage"},
	{Root: "testdata/sdk/", Owner: "sdk", FixtureDir: "testdata/sdk/success-sdk-phase11-root-coverage"},
	{Root: "testdata/ui/", Owner: "ui", FixtureDir: "testdata/ui/success-ui-phase11-root-coverage"},
	{Root: "testdata/statefile/", Owner: "statefile", FixtureDir: "testdata/statefile/success-statefile-phase11-root-coverage"},
	{Root: "testdata/archive/", Owner: "archive", FixtureDir: "testdata/archive/success-archive-phase11-root-coverage"},
	{Root: "testdata/commitstatus/", Owner: "commitstatus", FixtureDir: "testdata/commitstatus/success-commitstatus-phase11-root-coverage"},
	{Root: "testdata/security/", Owner: "security", FixtureDir: "testdata/security/success-security-phase11-root-coverage"},
}

type phase13ClosureRecord struct {
	RecordID  string   `json:"record_id"`
	Scope     string   `json:"scope"`
	Status    string   `json:"status"`
	OpenItems []string `json:"open_items"`
}

func phase13AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase13/implementation-alignment-quality/manifest.json" &&
		name == "phase-13-implementation-alignment-quality"
}

func phase14AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase14/obsidian-vault-integration/manifest.json" &&
		name == "phase-14-obsidian-vault-integration"
}

func phase15AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase15/obsidian-local-sync/manifest.json" &&
		name == "phase-15-obsidian-local-sync"
}

func phase13RequireOwnerPackages(t *testing.T) {
	t.Helper()

	for _, owner := range phase13GoOwners {
		dir := filepath.Join("components", owner)
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read owner package %s: %v", dir, err)
		}
		var files []string
		for _, entry := range entries {
			if entry.IsDir() {
				t.Fatalf("%s must not contain subdirectory %s", dir, entry.Name())
			}
			files = append(files, entry.Name())
			info, err := entry.Info()
			if err != nil {
				t.Fatalf("stat %s/%s: %v", dir, entry.Name(), err)
			}
			if info.Size() == 0 {
				t.Fatalf("%s/%s must not be empty", dir, entry.Name())
			}
		}
		sort.Strings(files)
		expected := []string{owner + ".go", owner + "_test.go", "execute.go", "model.go", "validate.go"}
		sort.Strings(expected)
		if !phase11StringSlicesEqual(files, expected) {
			t.Fatalf("%s must contain exactly the five owner files\nwant: %v\n got: %v", dir, expected, files)
		}
	}
}

func phase13RequireDocumentIndexRows(t *testing.T) {
	t.Helper()

	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, expected := range []string{
		"| [`testdata/phase13/implementation-alignment-quality/`](../testdata/phase13/implementation-alignment-quality/) | Phase 13 implementation alignment quality fixture root | 実在 |",
		"| [`.github/workflows/phase13-implementation-alignment-quality.yml`](../.github/workflows/phase13-implementation-alignment-quality.yml) | Phase 13 required check workflow | 実在 |",
		"| [`testdata/phase13/implementation-alignment-quality/manifest.json`](../testdata/phase13/implementation-alignment-quality/manifest.json) | Phase 13 evidence package manifest | 実在 |",
		"| [`testdata/phase13/implementation-alignment-quality/expected/counters.json`](../testdata/phase13/implementation-alignment-quality/expected/counters.json) | Phase 13 closure counter 期待値 | 実在 |",
		"| [`testdata/phase13/implementation-alignment-quality/records/closure.jsonl`](../testdata/phase13/implementation-alignment-quality/records/closure.jsonl) | Phase 13 closure record set | 実在 |",
	} {
		if !strings.Contains(documentIndex, expected) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must contain Phase 13実在 row: %s", expected)
		}
	}
}

func phase13RequireGovernanceFiles(t *testing.T) {
	t.Helper()

	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, file := range []string{"LICENSE", "SECURITY.md", "CONTRIBUTING.md", "CODEOWNERS", "CHANGELOG.md"} {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Phase 13 governance file %s must exist: %v", file, err)
		}
		if strings.TrimSpace(string(data)) == "" {
			t.Fatalf("Phase 13 governance file %s must not be empty", file)
		}
		row := "| [`" + file + "`](../" + file + ") | release governance artifact | 実在 |"
		if !strings.Contains(documentIndex, row) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index governance file as 実在: %s", file)
		}
	}
}

func phase13RequireZeroExternalDependencies(t *testing.T) {
	t.Helper()

	goMod := phase11MustReadText(t, "go.mod")
	for _, forbidden := range []string{"\nrequire ", "\nreplace ", "golang.org/x/"} {
		if strings.Contains(goMod, forbidden) {
			t.Fatalf("go.mod must not contain external dependency marker %q", forbidden)
		}
	}
	for _, path := range []string{"docs/SPEC.md", "docs/details/fixture.md", ".github/workflows/phase13-implementation-alignment-quality.yml"} {
		body := phase11MustReadText(t, path)
		for _, forbidden := range []string{"staticcheck", "govulncheck", "gosec", "Argon2id", "argon2id"} {
			if strings.Contains(body, forbidden) {
				t.Fatalf("%s must not contain forbidden external dependency term %q", path, forbidden)
			}
		}
	}
	implementationFiles, err := filepath.Glob("components/*/*.go")
	if err != nil {
		t.Fatalf("glob implementation files: %v", err)
	}
	implementationFiles = append(implementationFiles, "admin/adlaire-ci-sdk.js", "admin/index.html")
	for _, path := range implementationFiles {
		body := phase11MustReadText(t, path)
		for _, forbidden := range []string{"golang.org/x/", "npm install", "cdn.jsdelivr", "unpkg.com"} {
			if strings.Contains(body, forbidden) {
				t.Fatalf("%s must not contain external dependency marker %q", path, forbidden)
			}
		}
	}
}

func phase13ReadCounterFile(t *testing.T, path string) map[string]int {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var counters map[string]int
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&counters); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return counters
}

func phase13ReadClosureRecords(t *testing.T, path string) []phase13ClosureRecord {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	text := string(data)
	if !strings.HasSuffix(text, "\n") {
		t.Fatalf("%s must end with LF", path)
	}
	var records []phase13ClosureRecord
	for index, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			t.Fatalf("%s line %d must not be empty", path, index+1)
		}
		var record phase13ClosureRecord
		decoder := json.NewDecoder(strings.NewReader(line))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&record); err != nil {
			t.Fatalf("decode %s line %d: %v", path, index+1, err)
		}
		if strings.TrimSpace(record.RecordID) == "" {
			t.Fatalf("%s line %d record_id must not be empty", path, index+1)
		}
		records = append(records, record)
	}
	return records
}

func phase13ReadJSONLines(t *testing.T, path string) []map[string]any {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	text := string(data)
	if !strings.HasSuffix(text, "\n") {
		t.Fatalf("%s must end with LF", path)
	}
	var records []map[string]any
	for index, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			t.Fatalf("%s line %d must not be empty", path, index+1)
		}
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("decode %s line %d: %v", path, index+1, err)
		}
		records = append(records, record)
	}
	return records
}

func phase13RequirePinnedActions(t *testing.T, workflow string) {
	t.Helper()

	for _, line := range strings.Split(workflow, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "- uses: ") {
			continue
		}
		parts := strings.Split(line, "@")
		if len(parts) != 2 {
			t.Fatalf("workflow action must include one @ ref: %s", line)
		}
		ref := strings.TrimSpace(parts[1])
		if len(ref) != 40 {
			t.Fatalf("workflow action ref must be a 40 character commit SHA, got %q in %s", ref, line)
		}
		for _, char := range ref {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Fatalf("workflow action ref must be lowercase hex commit SHA, got %q in %s", ref, line)
			}
		}
	}
}

func phase11RequireNoUnexpectedFixtureDirectories(t *testing.T) {
	t.Helper()

	var unexpected []string
	err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if phase11HasNumericCopySuffix(d.Name()) {
			unexpected = append(unexpected, filepath.ToSlash(path))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan fixture directories: %v", err)
	}
	formalRoots := []string{"testdata/admin/cli", "testdata/builder", "testdata/setup", "testdata/release"}
	for _, coverage := range phase11RootCoverageFixtures {
		formalRoots = append(formalRoots, strings.TrimSuffix(coverage.Root, "/"))
	}
	for _, formalRoot := range formalRoots {
		entries, err := os.ReadDir(formalRoot)
		if err != nil {
			t.Fatalf("read formal fixture root %s: %v", formalRoot, err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			manifestPath := filepath.Join(formalRoot, entry.Name(), "manifest.json")
			if _, err := os.Stat(manifestPath); err != nil {
				unexpected = append(unexpected, filepath.ToSlash(filepath.Join(formalRoot, entry.Name())))
			}
		}
	}
	sort.Strings(unexpected)
	if len(unexpected) != 0 {
		t.Fatalf("unexpected or unregistered fixture directories remain: %v", unexpected)
	}
}

func phase11HasNumericCopySuffix(name string) bool {
	index := strings.LastIndex(name, " ")
	if index < 0 || index == len(name)-1 {
		return false
	}
	for _, char := range name[index+1:] {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
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
