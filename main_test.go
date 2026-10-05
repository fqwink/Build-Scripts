package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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

		if manifest.Name != filepath.Base(fixtureDir) && !phase13AllowsScopeName(manifestPath, manifest.Name) && !phase14AllowsScopeName(manifestPath, manifest.Name) && !phase15AllowsScopeName(manifestPath, manifest.Name) && !phase16AllowsScopeName(manifestPath, manifest.Name) && !phase17AllowsScopeName(manifestPath, manifest.Name) && !phase18AllowsScopeName(manifestPath, manifest.Name) {
			t.Fatalf("%s name must match fixture directory, got %q", manifestPath, manifest.Name)
		}
		if prior := seenNames[manifest.Name]; prior != "" {
			t.Fatalf("duplicate fixture manifest name %q in %s and %s", manifest.Name, prior, manifestPath)
		}
		seenNames[manifest.Name] = manifestPath

		if phase17AllowsScopeName(manifestPath, manifest.Name) || phase18AllowsScopeName(manifestPath, manifest.Name) {
			continue
		}

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
	if !strings.Contains(roadmap, "現在の active Phase は Phase 18 とする。") {
		t.Fatalf("docs/ROADMAP.md must keep Phase 18 active after Phase 17 closure until real VPS Pull Bootstrap completes")
	}
	if !strings.Contains(roadmap, "初期実装 Phase 1 から Phase 17 まではすべて `実装済み`") {
		t.Fatalf("docs/ROADMAP.md must state that Phase 1 through Phase 17 are all implemented after Phase 17 closure")
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

func TestPhase16QualityEvidenceClosure(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase16/quality-evidence-closure"
	phase16RequireFixtureFiles(t, root)

	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.Name != phase16Scope {
		t.Fatalf("%s name must be %s, got %q", manifestPath, phase16Scope, manifest.Name)
	}
	if manifest.Scope != phase16Scope {
		t.Fatalf("%s scope must be %s, got %q", manifestPath, phase16Scope, manifest.Scope)
	}
	if !phase11StringSlicesEqual(manifest.Owners, phase16Owners) {
		t.Fatalf("%s owners mismatch\nwant: %v\n got: %v", manifestPath, phase16Owners, manifest.Owners)
	}
	if !phase11StringSlicesEqual(manifest.Components, phase16Components) {
		t.Fatalf("%s components mismatch\nwant: %v\n got: %v", manifestPath, phase16Components, manifest.Components)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase16RequiredChecks) {
		t.Fatalf("%s required_checks mismatch\nwant: %v\n got: %v", manifestPath, phase16RequiredChecks, manifest.RequiredChecks)
	}
	if !phase11StringSlicesEqual(manifest.ClosureRecords, phase16ClosureRecords) {
		t.Fatalf("%s closure_records mismatch\nwant: %v\n got: %v", manifestPath, phase16ClosureRecords, manifest.ClosureRecords)
	}
	for _, sourceAnchor := range manifest.SourceAnchors {
		phase11RequireMarkdownReference(t, manifestPath, sourceAnchor)
	}

	counters := phase16ReadCounters(t, filepath.Join(root, "expected", "counters.json"))
	phase16RequireAllCountersZero(t, counters)

	inventory := phase16ReadInventoryFiles(t, root)
	negativeControls := phase16ReadNegativeControls(t, filepath.Join(root, "input", "negative_controls.json"))
	for _, control := range negativeControls {
		if _, exists := inventory[control.ID]; exists {
			t.Fatalf("negative control inventory id is duplicated: %s", control.ID)
		}
		inventory[control.ID] = phase16InventoryRecord{
			ID:             control.ID,
			WorkUnit:       control.WorkUnit,
			OwnerComponent: control.OwnerComponent,
			SourceRef:      control.PositiveControlRef,
			Classification: "negative_control",
			Status:         "closed",
			CounterKey:     control.CounterKey,
			ClosureRef:     control.PositiveControlRef,
			EvidenceRecords: []string{
				"records/execution.jsonl",
			},
		}
	}

	actions := phase16ReadActions(t, filepath.Join(root, "expected", "actions.json"))
	records := phase16ReadEvidenceRecords(t, root)

	phase16RequireSourceCoverage(t, filepath.Join(root, "input", "source_coverage.json"), inventory)
	phase16RequireInventoryActionRecordClosure(t, inventory, actions, records, counters)
	phase16RequireClosureRecordSet(t, records["records/closure.jsonl"])
	phase16RequireRequiredChecks(t, manifest, records)
	phase16RequireNegativeControls(t, negativeControls, actions, records)
	phase16RequireWorkflow(t, ".github/workflows/phase16-quality-evidence-closure.yml")
	phase16RequireDocumentDriftClosed(t)
}

func TestPhase17ProductionValidationEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase17/production-validation"
	phase17RequireFixtureFiles(t, root)

	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.SchemaVersion != 1 || manifest.Name != phase17Scope || manifest.Scope != phase17Scope || manifest.FixtureRoot != root+"/" {
		t.Fatalf("%s must use schema_version=1, name/scope=%s, fixture_root=%s/", manifestPath, phase17Scope, root)
	}
	if !phase11StringSlicesEqual(manifest.ProviderTargets, phase17ProviderTargets) {
		t.Fatalf("%s provider_targets mismatch\nwant: %v\n got: %v", manifestPath, phase17ProviderTargets, manifest.ProviderTargets)
	}
	if !phase11StringSlicesEqual(manifest.ValidationModes, phase17ValidationModes) {
		t.Fatalf("%s validation_modes mismatch\nwant: %v\n got: %v", manifestPath, phase17ValidationModes, manifest.ValidationModes)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase17RequiredChecks) {
		t.Fatalf("%s required_checks mismatch\nwant: %v\n got: %v", manifestPath, phase17RequiredChecks, manifest.RequiredChecks)
	}

	counters := phase17ReadCounters(t, filepath.Join(root, "expected", "counters.json"))
	phase17RequireCountersClosed(t, counters)
	phase17RequireInputsAndExpected(t, root)
	records := phase17ReadRecords(t, root)
	phase17RequireDerivedCounters(t, counters, records)
	phase17RequireRequiredCheckConnections(t, records)
	phase17RequireBoundaryRecords(t, counters, records)
	phase17RequireClosureEvidenceGates(t, counters, records)
	phase17RequireWorkflow(t, ".github/workflows/phase17-production-validation.yml")
	phase17RequireDocumentState(t)
}

func TestPhase18VPSPullBootstrapEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase18/vps-pull-bootstrap"
	phase18RequireFixtureFiles(t, root)

	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.SchemaVersion != 1 || manifest.Name != phase18Scope || manifest.Scope != phase18Scope || manifest.FixtureRoot != root+"/" {
		t.Fatalf("%s must use schema_version=1, name/scope=%s, fixture_root=%s/", manifestPath, phase18Scope, root)
	}
	if !phase11StringSlicesEqual(manifest.ProviderTargets, phase18ProviderTargets) {
		t.Fatalf("%s provider_targets mismatch\nwant: %v\n got: %v", manifestPath, phase18ProviderTargets, manifest.ProviderTargets)
	}
	if !phase11StringSlicesEqual(manifest.ValidationModes, phase18ValidationModes) {
		t.Fatalf("%s validation_modes mismatch\nwant: %v\n got: %v", manifestPath, phase18ValidationModes, manifest.ValidationModes)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase18RequiredChecks) {
		t.Fatalf("%s required_checks mismatch\nwant: %v\n got: %v", manifestPath, phase18RequiredChecks, manifest.RequiredChecks)
	}

	counters := phase18ReadCounters(t, filepath.Join(root, "expected", "counters.json"))
	phase18RequireCountersClosed(t, counters)
	inputs := phase18RequireInputsAndExpected(t, root)
	records := phase18ReadRecords(t, root)
	phase18RequireRequiredCheckConnections(t, records)
	phase18RequireSessionRecords(t, records["records/session.jsonl"], inputs)
	phase18RequireHealthRecords(t, records["records/health.jsonl"])
	phase18RequireRuntimeRecords(t, records["records/runtime.jsonl"])
	phase18RequireUpdateRollbackRecords(t, records["records/update_rollback.jsonl"])
	phase18RequireBugfixRecords(t, records["records/bugfix.jsonl"], inputs.Bugfix)
	phase18RequireSecurityRecords(t, records["records/security.jsonl"], inputs)
	phase18RequireDocumentDriftRecords(t, records["records/document_drift.jsonl"])
	phase18RequireDerivedCounters(t, counters, records)
	phase18RequireWorkflow(t, ".github/workflows/phase18-vps-pull-bootstrap.yml")
	phase18RequireDocumentState(t)
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
	SchemaVersion          int                     `json:"schema_version,omitempty"`
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
	ProviderTargets        []string                `json:"provider_targets,omitempty"`
	ValidationModes        []string                `json:"validation_modes,omitempty"`
	RequiredChecks         []string                `json:"required_checks,omitempty"`
	ClosureRecords         []string                `json:"closure_records,omitempty"`
	NegativeControls       []string                `json:"negative_controls,omitempty"`
	SourceAnchors          []string                `json:"source_anchors,omitempty"`
	FixtureRoot            string                  `json:"fixture_root,omitempty"`
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

const phase16Scope = "phase-16-quality-evidence-closure"

var phase16Owners = []string{
	"admin",
	"api",
	"archive",
	"builder",
	"commitstatus",
	"mcp",
	"obsidian",
	"release",
	"runner",
	"sdk",
	"security",
	"setup",
	"statefile",
	"ui",
}

var phase16Components = []string{
	"admin",
	"api",
	"archive",
	"builder",
	"commitstatus",
	"mcp",
	"obsidian",
	"release",
	"runner",
	"sdk",
	"security",
	"setup",
	"statefile",
	"ui",
}

var phase16RequiredChecks = []string{
	"phase16-go-format",
	"phase16-go-test",
	"phase16-deno-check",
	"phase16-race",
	"phase16-quality-evidence-fixture",
	"phase16-mutation",
	"phase16-fault-injection",
	"phase16-workflow-hardening",
	"phase16-document-drift",
}

var phase16ClosureRecords = []string{
	"records/closure.jsonl",
	"records/execution.jsonl",
	"records/mutation.jsonl",
	"records/fault.jsonl",
	"records/workflow.jsonl",
}

var phase16ClosureCounters = []string{
	"phase16_declarative_evidence_open_count",
	"phase16_execution_evidence_gap_count",
	"phase16_mutation_survived_count",
	"phase16_fault_injection_open_count",
	"phase16_align_unclassified_count",
	"phase16_align_open_count",
	"phase16_ignored_error_unclassified_count",
	"phase16_required_write_failure_open_count",
	"phase16_skip_open_count",
	"phase16_panic_unclassified_count",
	"phase16_determinism_open_count",
	"phase16_http_boundary_open_count",
	"phase16_filesystem_durability_open_count",
	"phase16_workflow_hardening_open_count",
	"phase16_validation_portability_open_count",
	"phase16_release_rehearsal_open_count",
	"phase16_large_owner_risk_open_count",
	"phase16_document_drift_open_count",
	"final_open_item_count",
}

var phase16RequiredFixtureFiles = []string{
	"manifest.json",
	"input/source_coverage.json",
	"input/evidence_inventory.json",
	"input/align_inventory.json",
	"input/error_inventory.json",
	"input/determinism_inventory.json",
	"input/filesystem_inventory.json",
	"input/workflow_inventory.json",
	"input/large_owner_inventory.json",
	"input/negative_controls.json",
	"expected/counters.json",
	"expected/actions.json",
	"records/closure.jsonl",
	"records/execution.jsonl",
	"records/mutation.jsonl",
	"records/fault.jsonl",
	"records/workflow.jsonl",
}

var phase16InventoryFiles = []string{
	"input/evidence_inventory.json",
	"input/align_inventory.json",
	"input/error_inventory.json",
	"input/determinism_inventory.json",
	"input/filesystem_inventory.json",
	"input/workflow_inventory.json",
	"input/large_owner_inventory.json",
}

type phase16CounterFile struct {
	SchemaVersion int            `json:"schema_version"`
	Scope         string         `json:"scope"`
	Counters      map[string]int `json:"counters"`
}

type phase16InventoryFile struct {
	SchemaVersion int                      `json:"schema_version"`
	Scope         string                   `json:"scope"`
	WorkUnit      string                   `json:"work_unit"`
	Records       []phase16InventoryRecord `json:"records"`
}

type phase16InventoryRecord struct {
	ID               string   `json:"id"`
	WorkUnit         string   `json:"work_unit"`
	OwnerComponent   string   `json:"owner_component"`
	SourceRef        string   `json:"source_ref"`
	Classification   string   `json:"classification"`
	Status           string   `json:"status"`
	CounterKey       string   `json:"counter_key"`
	ClosureRef       string   `json:"closure_ref"`
	NotApplicableRef *string  `json:"not_applicable_ref"`
	FutureRef        *string  `json:"future_ref"`
	EvidenceRecords  []string `json:"evidence_records"`
}

type phase16ActionFile struct {
	SchemaVersion int             `json:"schema_version"`
	Scope         string          `json:"scope"`
	Actions       []phase16Action `json:"actions"`
}

type phase16Action struct {
	ID                 string `json:"id"`
	InventoryID        string `json:"inventory_id"`
	Action             string `json:"action"`
	OwnerComponent     string `json:"owner_component"`
	ExpectedCounterKey string `json:"expected_counter_key"`
	ExpectedStatus     string `json:"expected_status"`
	EvidenceRecord     string `json:"evidence_record"`
}

type phase16NegativeControlFile struct {
	SchemaVersion    int                      `json:"schema_version"`
	Scope            string                   `json:"scope"`
	NegativeControls []phase16NegativeControl `json:"negative_controls"`
}

type phase16NegativeControl struct {
	ID                  string `json:"id"`
	WorkUnit            string `json:"work_unit"`
	OwnerComponent      string `json:"owner_component"`
	TargetContract      string `json:"target_contract"`
	ExpectedFailureCode string `json:"expected_failure_code"`
	IsolationMode       string `json:"isolation_mode"`
	RequiredCheck       string `json:"required_check"`
	CounterKey          string `json:"counter_key"`
	PositiveControlRef  string `json:"positive_control_ref"`
}

type phase16SourceCoverageFile struct {
	SchemaVersion     int                   `json:"schema_version"`
	Scope             string                `json:"scope"`
	GeneratedFrom     string                `json:"generated_from"`
	SourceGroupSHA256 string                `json:"source_group_sha256"`
	GoSum             phase16GoSumRecord    `json:"go_sum"`
	Sources           []phase16SourceRecord `json:"sources"`
}

type phase16GoSumRecord struct {
	Path         string  `json:"path"`
	SourceExists bool    `json:"source_exists"`
	SHA256       *string `json:"sha256"`
	SizeBytes    int64   `json:"size_bytes"`
	Reason       string  `json:"reason"`
	InventoryID  string  `json:"inventory_id"`
}

type phase16SourceRecord struct {
	Path           string                `json:"path"`
	SourceExists   bool                  `json:"source_exists"`
	SHA256         string                `json:"sha256"`
	SizeBytes      int64                 `json:"size_bytes"`
	OwnerComponent string                `json:"owner_component"`
	Area           string                `json:"area"`
	DetectedTerms  []phase16DetectedTerm `json:"detected_terms"`
	InventoryRefs  []string              `json:"inventory_refs"`
}

type phase16DetectedTerm struct {
	TermID             string `json:"term_id"`
	Match              string `json:"match"`
	MatchMode          string `json:"match_mode"`
	Count              int    `json:"count"`
	FirstSourceLocator string `json:"first_source_locator"`
}

type phase16EvidenceRecord struct {
	RecordID               string `json:"record_id"`
	Scope                  string `json:"scope"`
	WorkUnit               string `json:"work_unit"`
	InventoryID            string `json:"inventory_id"`
	OwnerComponent         string `json:"owner_component"`
	SourceRef              string `json:"source_ref"`
	Result                 string `json:"result"`
	CounterKey             string `json:"counter_key"`
	CheckName              string `json:"check_name,omitempty"`
	Command                string `json:"command,omitempty"`
	Runtime                string `json:"runtime,omitempty"`
	InputSHA256            string `json:"input_sha256,omitempty"`
	ExpectedSHA256         string `json:"expected_sha256,omitempty"`
	ActualSHA256           string `json:"actual_sha256,omitempty"`
	DiffResult             string `json:"diff_result,omitempty"`
	ExitCode               *int   `json:"exit_code,omitempty"`
	StdoutSHA256           string `json:"stdout_sha256,omitempty"`
	StderrSHA256           string `json:"stderr_sha256,omitempty"`
	MutationTarget         string `json:"mutation_target,omitempty"`
	MutationClass          string `json:"mutation_class,omitempty"`
	BeforeSHA256           string `json:"before_sha256,omitempty"`
	AfterSHA256            string `json:"after_sha256,omitempty"`
	ExpectedFailure        string `json:"expected_failure,omitempty"`
	ActualFailure          string `json:"actual_failure,omitempty"`
	MutationResult         string `json:"mutation_result,omitempty"`
	FaultClass             string `json:"fault_class,omitempty"`
	InjectionResult        string `json:"injection_result,omitempty"`
	RecoveryResult         string `json:"recovery_result,omitempty"`
	SilentSuccess          *bool  `json:"silent_success,omitempty"`
	WorkflowPath           string `json:"workflow_path,omitempty"`
	JobName                string `json:"job_name,omitempty"`
	ActionPinResult        string `json:"action_pin_result,omitempty"`
	PermissionsResult      string `json:"permissions_result,omitempty"`
	TimeoutResult          string `json:"timeout_result,omitempty"`
	RequiredCheckResult    string `json:"required_check_result,omitempty"`
	DockerFallbackResult   string `json:"docker_fallback_result,omitempty"`
	RuntimeVersionResult   string `json:"runtime_version_result,omitempty"`
	DenoCheckResult        string `json:"deno_check_result,omitempty"`
	ReleaseRehearsalResult string `json:"release_rehearsal_result,omitempty"`
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

func phase16AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase16/quality-evidence-closure/manifest.json" &&
		name == phase16Scope
}

func phase16RequireFixtureFiles(t *testing.T, root string) {
	t.Helper()

	info, err := os.Stat(root)
	if err != nil {
		t.Fatalf("Phase 16 fixture root must exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("Phase 16 fixture root must be a directory: %s", root)
	}
	for _, rel := range phase16RequiredFixtureFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Phase 16 fixture path missing %s: %v", rel, err)
		}
		if info.IsDir() {
			t.Fatalf("Phase 16 fixture path must be a file: %s", rel)
		}
		if strings.HasSuffix(rel, ".json") || strings.HasSuffix(rel, ".jsonl") {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			if len(data) == 0 || data[len(data)-1] != '\n' {
				t.Fatalf("%s must end with LF", path)
			}
			if strings.HasSuffix(rel, ".json") {
				phase16RequireNoDuplicateJSONKeys(t, path, data)
			}
		}
	}
}

func phase16ReadCounters(t *testing.T, path string) map[string]int {
	t.Helper()

	var file phase16CounterFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase16Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase16Scope)
	}
	return file.Counters
}

func phase16RequireAllCountersZero(t *testing.T, counters map[string]int) {
	t.Helper()

	if len(counters) != len(phase16ClosureCounters) {
		t.Fatalf("Phase 16 counters must include exactly %d keys, got %d", len(phase16ClosureCounters), len(counters))
	}
	for _, key := range phase16ClosureCounters {
		value, ok := counters[key]
		if !ok {
			t.Fatalf("Phase 16 counters missing %s", key)
		}
		if value != 0 {
			t.Fatalf("Phase 16 counter %s must be 0, got %d", key, value)
		}
	}
}

func phase16ReadInventoryFiles(t *testing.T, root string) map[string]phase16InventoryRecord {
	t.Helper()

	records := map[string]phase16InventoryRecord{}
	for _, rel := range phase16InventoryFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		var file phase16InventoryFile
		phase16ReadStrictJSON(t, path, &file)
		if file.SchemaVersion != 1 || file.Scope != phase16Scope || strings.TrimSpace(file.WorkUnit) == "" {
			t.Fatalf("%s must use schema_version=1, scope=%s, and non-empty work_unit", path, phase16Scope)
		}
		if len(file.Records) == 0 {
			t.Fatalf("%s must include at least one inventory record", path)
		}
		for _, record := range file.Records {
			phase16RequireInventoryRecord(t, path, record, file.WorkUnit)
			if _, exists := records[record.ID]; exists {
				t.Fatalf("duplicate Phase 16 inventory id %s", record.ID)
			}
			records[record.ID] = record
		}
	}
	return records
}

func phase16RequireInventoryRecord(t *testing.T, path string, record phase16InventoryRecord, fileWorkUnit string) {
	t.Helper()

	if record.ID == "" || record.WorkUnit == "" || record.OwnerComponent == "" || record.SourceRef == "" || record.Classification == "" || record.Status == "" || record.CounterKey == "" || record.ClosureRef == "" {
		t.Fatalf("%s inventory record must not contain empty required fields: %+v", path, record)
	}
	if record.WorkUnit != fileWorkUnit {
		t.Fatalf("%s record %s work_unit mismatch: file=%s record=%s", path, record.ID, fileWorkUnit, record.WorkUnit)
	}
	if record.Status != "closed" && record.Status != "not_applicable" && record.Status != "future_plan" {
		t.Fatalf("%s record %s has invalid status %q", path, record.ID, record.Status)
	}
	if !phase16CounterAllowed(record.CounterKey) {
		t.Fatalf("%s record %s uses unknown counter %s", path, record.ID, record.CounterKey)
	}
	if len(record.EvidenceRecords) == 0 {
		t.Fatalf("%s record %s must connect to at least one evidence record file", path, record.ID)
	}
	for _, rel := range record.EvidenceRecords {
		if !phase16ClosureRecordAllowed(rel) {
			t.Fatalf("%s record %s uses unknown evidence record %s", path, record.ID, rel)
		}
	}
}

func phase16ReadActions(t *testing.T, path string) []phase16Action {
	t.Helper()

	var file phase16ActionFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase16Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase16Scope)
	}
	if len(file.Actions) == 0 {
		t.Fatalf("%s must include actions", path)
	}
	return file.Actions
}

func phase16ReadNegativeControls(t *testing.T, path string) []phase16NegativeControl {
	t.Helper()

	var file phase16NegativeControlFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase16Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase16Scope)
	}
	if len(file.NegativeControls) < len(phase16RequiredChecks) {
		t.Fatalf("%s must include at least one negative control per required check", path)
	}
	seenTargets := map[string]bool{}
	for _, control := range file.NegativeControls {
		if control.ID == "" || control.WorkUnit == "" || control.OwnerComponent == "" || control.TargetContract == "" || control.ExpectedFailureCode == "" || control.RequiredCheck == "" || control.CounterKey == "" || control.PositiveControlRef == "" {
			t.Fatalf("%s negative control must not contain empty required fields: %+v", path, control)
		}
		if control.IsolationMode != "temporary_copy" {
			t.Fatalf("%s negative control %s must use isolation_mode=temporary_copy", path, control.ID)
		}
		if !strings.HasPrefix(control.ExpectedFailureCode, "PHASE16_") {
			t.Fatalf("%s negative control %s failure code must start with PHASE16_", path, control.ID)
		}
		if !phase16RequiredCheckAllowed(control.RequiredCheck) {
			t.Fatalf("%s negative control %s uses unknown required check %s", path, control.ID, control.RequiredCheck)
		}
		if !phase16CounterAllowed(control.CounterKey) {
			t.Fatalf("%s negative control %s uses unknown counter %s", path, control.ID, control.CounterKey)
		}
		seenTargets[control.TargetContract] = true
	}
	for _, target := range []string{
		"json-duplicate-key",
		"source-coverage",
		"counter-reaggregation",
		"mutation-survivor",
		"fault-injection",
		"workflow-hardening",
		"future-phase-boundary",
	} {
		if !seenTargets[target] {
			t.Fatalf("%s must include negative control target_contract=%s", path, target)
		}
	}
	return file.NegativeControls
}

func phase16ReadEvidenceRecords(t *testing.T, root string) map[string][]phase16EvidenceRecord {
	t.Helper()

	records := map[string][]phase16EvidenceRecord{}
	for _, rel := range phase16ClosureRecords {
		path := filepath.Join(root, filepath.FromSlash(rel))
		records[rel] = phase16ReadEvidenceRecordFile(t, path, rel)
		if len(records[rel]) == 0 {
			t.Fatalf("%s must include at least one record", path)
		}
	}
	return records
}

func phase16ReadEvidenceRecordFile(t *testing.T, path string, rel string) []phase16EvidenceRecord {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	text := string(data)
	if !strings.HasSuffix(text, "\n") {
		t.Fatalf("%s must end with LF", path)
	}
	var records []phase16EvidenceRecord
	seen := map[string]bool{}
	for index, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			t.Fatalf("%s line %d must not be empty", path, index+1)
		}
		phase16RequireNoDuplicateJSONKeys(t, fmt.Sprintf("%s line %d", path, index+1), []byte(line))
		var record phase16EvidenceRecord
		decoder := json.NewDecoder(strings.NewReader(line))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&record); err != nil {
			t.Fatalf("decode %s line %d: %v", path, index+1, err)
		}
		if record.RecordID == "" || record.Scope != phase16Scope || record.WorkUnit == "" || record.InventoryID == "" || record.OwnerComponent == "" || record.SourceRef == "" || record.Result == "" || record.CounterKey == "" {
			t.Fatalf("%s line %d common fields are invalid: %+v", path, index+1, record)
		}
		if seen[record.RecordID] {
			t.Fatalf("%s repeats record_id %s", path, record.RecordID)
		}
		seen[record.RecordID] = true
		if !strings.HasPrefix(record.RecordID, "phase16.record."+strings.TrimSuffix(filepath.Base(rel), ".jsonl")+".") {
			t.Fatalf("%s line %d record_id does not match record file: %s", path, index+1, record.RecordID)
		}
		if record.Result == "failed" || record.Result == "rejected" {
			t.Fatalf("%s line %d must not leave failed or rejected result: %s", path, index+1, record.Result)
		}
		if !phase16CounterAllowed(record.CounterKey) {
			t.Fatalf("%s line %d uses unknown counter %s", path, index+1, record.CounterKey)
		}
		phase16RequireRecordSpecificFields(t, path, index+1, rel, record)
		records = append(records, record)
	}
	return records
}

func phase16RequireRecordSpecificFields(t *testing.T, path string, line int, rel string, record phase16EvidenceRecord) {
	t.Helper()

	switch rel {
	case "records/execution.jsonl":
		if record.CheckName == "" || record.Command == "" || record.Runtime == "" || record.InputSHA256 == "" || record.ExpectedSHA256 == "" || record.ActualSHA256 == "" || record.DiffResult == "" || record.ExitCode == nil || record.StdoutSHA256 == "" || record.StderrSHA256 == "" {
			t.Fatalf("%s line %d execution record missing required execution fields: %+v", path, line, record)
		}
		if strings.HasPrefix(record.InventoryID, "phase16.negative.") {
			if !strings.HasPrefix(record.DiffResult, "PHASE16_") {
				t.Fatalf("%s line %d negative execution diff_result must be PHASE16_* failure code, got %q", path, line, record.DiffResult)
			}
		} else if record.DiffResult != "empty" {
			t.Fatalf("%s line %d execution diff_result must be empty, got %q", path, line, record.DiffResult)
		}
	case "records/mutation.jsonl":
		if record.MutationTarget == "" || record.MutationClass == "" || record.BeforeSHA256 == "" || record.AfterSHA256 == "" || record.ExpectedFailure == "" || record.ActualFailure == "" || record.MutationResult == "" {
			t.Fatalf("%s line %d mutation record missing required fields: %+v", path, line, record)
		}
		if record.MutationResult == "survived" {
			t.Fatalf("%s line %d mutation must not survive", path, line)
		}
	case "records/fault.jsonl":
		if record.FaultClass == "" || record.InjectionResult == "" || record.RecoveryResult == "" || record.SilentSuccess == nil {
			t.Fatalf("%s line %d fault record missing required fields: %+v", path, line, record)
		}
		if *record.SilentSuccess {
			t.Fatalf("%s line %d fault record must not be silent success", path, line)
		}
	case "records/workflow.jsonl":
		if record.WorkflowPath == "" || record.JobName == "" || record.ActionPinResult == "" || record.PermissionsResult == "" || record.TimeoutResult == "" || record.RequiredCheckResult == "" || record.DockerFallbackResult == "" || record.RuntimeVersionResult == "" || record.DenoCheckResult == "" || record.ReleaseRehearsalResult == "" {
			t.Fatalf("%s line %d workflow record missing required fields: %+v", path, line, record)
		}
		for _, result := range []string{record.ActionPinResult, record.PermissionsResult, record.TimeoutResult, record.RequiredCheckResult, record.DockerFallbackResult, record.RuntimeVersionResult, record.DenoCheckResult, record.ReleaseRehearsalResult} {
			if result != "passed" && result != "not_applicable" {
				t.Fatalf("%s line %d workflow result must be passed or not_applicable, got %q", path, line, result)
			}
		}
	}
}

func phase16RequireSourceCoverage(t *testing.T, path string, inventory map[string]phase16InventoryRecord) {
	t.Helper()

	var file phase16SourceCoverageFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase16Scope || file.GeneratedFrom == "" {
		t.Fatalf("%s must use schema_version=1, scope=%s, and generated_from", path, phase16Scope)
	}
	if file.GoSum.Path != "go.sum" {
		t.Fatalf("%s go_sum path must be go.sum", path)
	}
	_, goSumErr := os.Stat("go.sum")
	if os.IsNotExist(goSumErr) {
		if file.GoSum.SourceExists || file.GoSum.SHA256 != nil || file.GoSum.SizeBytes != 0 || file.GoSum.InventoryID == "" {
			t.Fatalf("%s must record absent go.sum as non-existent source", path)
		}
	} else if goSumErr == nil {
		if !file.GoSum.SourceExists || file.GoSum.SHA256 == nil || *file.GoSum.SHA256 != phase16FileSHA256(t, "go.sum") {
			t.Fatalf("%s go_sum hash must match live go.sum", path)
		}
	} else {
		t.Fatalf("stat go.sum: %v", goSumErr)
	}

	if len(file.Sources) < 10 {
		t.Fatalf("%s must include at least 10 source coverage records", path)
	}
	seen := map[string]bool{}
	var digestMaterial strings.Builder
	for _, source := range file.Sources {
		if source.Path == "" || source.OwnerComponent == "" || source.Area == "" || len(source.InventoryRefs) == 0 || len(source.DetectedTerms) == 0 {
			t.Fatalf("%s source record must not contain empty required fields: %+v", path, source)
		}
		if seen[source.Path] {
			t.Fatalf("%s repeats source path %s", path, source.Path)
		}
		seen[source.Path] = true
		info, err := os.Stat(filepath.FromSlash(source.Path))
		if err != nil {
			t.Fatalf("%s source path %s must exist: %v", path, source.Path, err)
		}
		if info.IsDir() {
			t.Fatalf("%s source path %s must be a file", path, source.Path)
		}
		if !source.SourceExists {
			t.Fatalf("%s source path %s must use source_exists=true", path, source.Path)
		}
		if source.SizeBytes != info.Size() {
			t.Fatalf("%s source %s size mismatch: want %d got %d", path, source.Path, info.Size(), source.SizeBytes)
		}
		if got := phase16FileSHA256(t, source.Path); got != source.SHA256 {
			t.Fatalf("%s source %s sha256 mismatch: want %s got %s", path, source.Path, source.SHA256, got)
		}
		for _, ref := range source.InventoryRefs {
			if _, ok := inventory[ref]; !ok {
				t.Fatalf("%s source %s references unknown inventory id %s", path, source.Path, ref)
			}
		}
		for _, term := range source.DetectedTerms {
			if term.TermID == "" || term.Match == "" || term.MatchMode == "" || term.Count < 1 || term.FirstSourceLocator == "" {
				t.Fatalf("%s source %s contains invalid detected term: %+v", path, source.Path, term)
			}
		}
		digestMaterial.WriteString(source.Path)
		digestMaterial.WriteByte('\n')
		digestMaterial.WriteString(source.SHA256)
		digestMaterial.WriteByte('\n')
	}
	if file.SourceGroupSHA256 != phase16StringSHA256(digestMaterial.String()) {
		t.Fatalf("%s source_group_sha256 mismatch", path)
	}
}

func phase16RequireInventoryActionRecordClosure(t *testing.T, inventory map[string]phase16InventoryRecord, actions []phase16Action, records map[string][]phase16EvidenceRecord, counters map[string]int) {
	t.Helper()

	actionByInventory := map[string]phase16Action{}
	for _, action := range actions {
		if action.ID == "" || action.InventoryID == "" || action.Action == "" || action.OwnerComponent == "" || action.ExpectedCounterKey == "" || action.ExpectedStatus == "" || action.EvidenceRecord == "" {
			t.Fatalf("Phase 16 action must not contain empty required fields: %+v", action)
		}
		if _, ok := inventory[action.InventoryID]; !ok {
			t.Fatalf("Phase 16 action %s references unknown inventory id %s", action.ID, action.InventoryID)
		}
		if _, duplicate := actionByInventory[action.InventoryID]; duplicate {
			t.Fatalf("Phase 16 inventory id %s has multiple actions", action.InventoryID)
		}
		if _, ok := counters[action.ExpectedCounterKey]; !ok {
			t.Fatalf("Phase 16 action %s references unknown counter %s", action.ID, action.ExpectedCounterKey)
		}
		if !phase16ClosureRecordAllowed(action.EvidenceRecord) {
			t.Fatalf("Phase 16 action %s references unknown evidence record %s", action.ID, action.EvidenceRecord)
		}
		actionByInventory[action.InventoryID] = action
	}
	if len(actionByInventory) != len(inventory) {
		t.Fatalf("Phase 16 actions must cover every inventory record: actions=%d inventory=%d", len(actionByInventory), len(inventory))
	}

	recordByInventory := map[string][]phase16EvidenceRecord{}
	for rel, recordSet := range records {
		for _, record := range recordSet {
			if _, ok := inventory[record.InventoryID]; !ok {
				t.Fatalf("%s record %s references unknown inventory id %s", rel, record.RecordID, record.InventoryID)
			}
			if _, ok := counters[record.CounterKey]; !ok {
				t.Fatalf("%s record %s references unknown counter %s", rel, record.RecordID, record.CounterKey)
			}
			recordByInventory[record.InventoryID] = append(recordByInventory[record.InventoryID], record)
		}
	}
	for id, record := range inventory {
		action := actionByInventory[id]
		if action.OwnerComponent != record.OwnerComponent || action.ExpectedCounterKey != record.CounterKey || action.ExpectedStatus != record.Status {
			t.Fatalf("Phase 16 action %s does not match inventory record %+v", action.ID, record)
		}
		var connected bool
		for _, evidence := range recordByInventory[id] {
			if phase16RecordFileForID(evidence.RecordID) == action.EvidenceRecord {
				connected = true
				break
			}
		}
		if !connected {
			t.Fatalf("Phase 16 inventory %s is not connected to action evidence %s", id, action.EvidenceRecord)
		}
	}
}

func phase16RequireClosureRecordSet(t *testing.T, records []phase16EvidenceRecord) {
	t.Helper()

	if len(records) != 18 {
		t.Fatalf("Phase 16 closure record set must include 18 records, got %d", len(records))
	}
	for index, record := range records {
		want := fmt.Sprintf("phase16.record.closure.all.closure-%02d", index+1)
		if record.RecordID != want {
			t.Fatalf("Phase 16 closure record order mismatch: want %s got %s", want, record.RecordID)
		}
		if record.Result != "passed" || record.CounterKey == "" {
			t.Fatalf("Phase 16 closure record %s must be passed and counter-connected", record.RecordID)
		}
	}
}

func phase16RequireRequiredChecks(t *testing.T, manifest phase11Manifest, records map[string][]phase16EvidenceRecord) {
	t.Helper()

	required := map[string]bool{}
	for _, check := range phase16RequiredChecks {
		required[check] = false
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase16RequiredChecks) {
		t.Fatalf("Phase 16 manifest required checks drift")
	}
	for _, record := range records["records/execution.jsonl"] {
		if _, ok := required[record.CheckName]; ok {
			required[record.CheckName] = true
		}
	}
	for _, record := range records["records/workflow.jsonl"] {
		if _, ok := required[record.JobName]; ok {
			required[record.JobName] = true
		}
	}
	for check, seen := range required {
		if !seen {
			t.Fatalf("Phase 16 required check %s is not connected to execution/workflow records", check)
		}
	}
}

func phase16RequireNegativeControls(t *testing.T, controls []phase16NegativeControl, actions []phase16Action, records map[string][]phase16EvidenceRecord) {
	t.Helper()

	actionByInventory := map[string]phase16Action{}
	for _, action := range actions {
		actionByInventory[action.InventoryID] = action
	}
	executionByInventory := map[string][]phase16EvidenceRecord{}
	for _, record := range records["records/execution.jsonl"] {
		executionByInventory[record.InventoryID] = append(executionByInventory[record.InventoryID], record)
	}
	for _, control := range controls {
		action, ok := actionByInventory[control.ID]
		if !ok {
			t.Fatalf("negative control %s must have action", control.ID)
		}
		if action.Action != "execute_and_close" || action.EvidenceRecord != "records/execution.jsonl" || action.ExpectedStatus != "closed" {
			t.Fatalf("negative control %s action must execute_and_close via records/execution.jsonl", control.ID)
		}
		var detected bool
		for _, record := range executionByInventory[control.ID] {
			if record.CheckName == control.RequiredCheck && record.DiffResult == control.ExpectedFailureCode && record.Result == "passed" {
				detected = true
				break
			}
		}
		if !detected {
			t.Fatalf("negative control %s expected failure %s was not detected", control.ID, control.ExpectedFailureCode)
		}
	}
}

func phase16RequireWorkflow(t *testing.T, path string) {
	t.Helper()

	workflow := phase11MustReadText(t, path)
	if !strings.Contains(workflow, "permissions:\n  contents: read") {
		t.Fatalf("%s must use minimum contents:read permissions", path)
	}
	for _, requiredCheck := range phase16RequiredChecks {
		if !strings.Contains(workflow, "\n  "+requiredCheck+":") {
			t.Fatalf("%s must define required check job %s", path, requiredCheck)
		}
	}
	if strings.Count(workflow, "timeout-minutes:") < len(phase16RequiredChecks) {
		t.Fatalf("%s must set timeout-minutes on every required check", path)
	}
	phase13RequirePinnedActions(t, workflow)
}

func phase16RequireDocumentDriftClosed(t *testing.T) {
	t.Helper()

	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	if !strings.Contains(roadmap, "| Phase 16 | 実装済み品質証跡実体化・追加検証候補 closure。") || !strings.Contains(roadmap, "| 実装済み | Phase 15 |") {
		t.Fatalf("docs/ROADMAP.md must mark Phase 16 as 実装済み")
	}
	if !strings.Contains(roadmap, "現在の active Phase は Phase 18 とする。") {
		t.Fatalf("docs/ROADMAP.md must keep Phase 18 active after Phase 17 closure until real VPS Pull Bootstrap completes")
	}
	for _, feature := range []string{
		"Phase 16 source coverage set / detection registry / evidence package manifest / inventory / negative control / expected / record schema / checker implementation artifact / checker 実行入口 / checker 再導出 gate",
		"Phase 16 宣言型証跡の実行型証跡化 / production entrypoint 実行 / actual expected 比較 gate",
		"Phase 16 `ALIGN-*` 追加検証候補分類・完了 / 未分類 0 gate",
		"Phase 16 ignored error 分類 / required write failure / panic / skip closure gate",
		"Phase 16 clock / sleep / timeout / entropy determinism gate",
		"Phase 16 filesystem durability parity / atomic write / fsync / parent directory fsync / symlink 非追従 gate",
		"Phase 16 Phase 12 workflow hardening / Docker 検証手順固定 / Deno stable JavaScript 検証 / release rehearsal gate",
		"Phase 16 巨大 owner file risk ledger / 5 ファイル原則維持 / 内部責務区画検査 gate",
	} {
		if !strings.Contains(roadmap, "| 実装済み |") || !strings.Contains(roadmap, feature) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 16 feature implemented: %s", feature)
		}
		if strings.Contains(roadmap, "| 仕様化済み・未実装 | 検証基盤 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 実装品質 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 状態管理 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | CI / 配布 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 保守性 | "+feature+" |") {
			t.Fatalf("docs/ROADMAP.md must not keep Phase 16 feature unimplemented: %s", feature)
		}
	}
	for _, expected := range []string{
		"| [`testdata/phase16/quality-evidence-closure/`](../testdata/phase16/quality-evidence-closure/) | Phase 16 quality evidence closure fixture root | 実在 |",
		"| [`.github/workflows/phase16-quality-evidence-closure.yml`](../.github/workflows/phase16-quality-evidence-closure.yml) | Phase 16 required check workflow | 実在 |",
		"| `testdata/phase16/quality-evidence-closure/manifest.json` | Phase 16 evidence package manifest | 実在 |",
		"| `testdata/phase16/quality-evidence-closure/records/workflow.jsonl` | Phase 16 workflow hardening 証跡 | 実在 |",
	} {
		if !strings.Contains(documentIndex, expected) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must contain Phase 16実在 row: %s", expected)
		}
	}
	for _, forbidden := range []string{
		"| `testdata/phase16/quality-evidence-closure/` | Phase 16 quality evidence closure fixture root | 未作成 |",
		"| `.github/workflows/phase16-quality-evidence-closure.yml` | Phase 16 required check workflow | 未作成 |",
		"| `testdata/phase16/quality-evidence-closure/manifest.json` | Phase 16 evidence package manifest | 未作成 |",
		"| `testdata/phase16/quality-evidence-closure/records/workflow.jsonl` | Phase 16 workflow hardening 証跡 | 未作成 |",
	} {
		if strings.Contains(documentIndex, forbidden) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must not keep Phase 16 path as 未作成: %s", forbidden)
		}
	}
}

func phase16ReadStrictJSON(t *testing.T, path string, out any) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Fatalf("%s must end with LF", path)
	}
	phase16RequireNoDuplicateJSONKeys(t, path, data)
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		t.Fatalf("%s must contain exactly one JSON value", path)
	}
}

func phase16RequireNoDuplicateJSONKeys(t *testing.T, path string, data []byte) {
	t.Helper()

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := phase16CheckJSONValueForDuplicateKeys(decoder); err != nil {
		t.Fatalf("%s duplicate-key check failed: %v", path, err)
	}
	if _, err := decoder.Token(); err != io.EOF {
		t.Fatalf("%s must contain exactly one JSON value", path)
	}
}

func phase16CheckJSONValueForDuplicateKeys(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("object key is not string")
			}
			if seen[key] {
				return fmt.Errorf("duplicate key %q", key)
			}
			seen[key] = true
			if err := phase16CheckJSONValueForDuplicateKeys(decoder); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim('}') {
			return fmt.Errorf("object not closed")
		}
	case '[':
		for decoder.More() {
			if err := phase16CheckJSONValueForDuplicateKeys(decoder); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim(']') {
			return fmt.Errorf("array not closed")
		}
	default:
		return fmt.Errorf("unexpected delimiter %q", delim)
	}
	return nil
}

func phase16CounterAllowed(key string) bool {
	for _, allowed := range phase16ClosureCounters {
		if key == allowed {
			return true
		}
	}
	return false
}

func phase16ClosureRecordAllowed(rel string) bool {
	for _, allowed := range phase16ClosureRecords {
		if rel == allowed {
			return true
		}
	}
	return false
}

func phase16RequiredCheckAllowed(check string) bool {
	for _, allowed := range phase16RequiredChecks {
		if check == allowed {
			return true
		}
	}
	return false
}

func phase16RecordFileForID(recordID string) string {
	parts := strings.Split(recordID, ".")
	if len(parts) < 3 {
		return ""
	}
	return "records/" + parts[2] + ".jsonl"
}

func phase16FileSHA256(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.FromSlash(path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func phase16StringSHA256(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

const phase17Scope = "phase-17-production-validation"

const phase17OperationWindow = "initial-production-validation"

const phase17MonitorIntervalSeconds = 0

const phase17RequiredMonitorSampleCount = 1

const phase17MaxConsecutiveMonitorFailures = 0

var phase17ProviderTargets = []string{"conoha-vps-primary", "xserver-vps-future"}

var phase17ValidationModes = []string{"trial-production-conoha-vps", "production-equivalent-simulation"}

var phase17RequiredChecks = []string{
	"phase17-conoha-vps-preflight",
	"phase17-trial-production-buildout",
	"phase17-conoha-install-update-rollback",
	"phase17-systemd-lifecycle",
	"phase17-runtime-flow",
	"phase17-trial-operation-loop",
	"phase17-production-simulation",
	"phase17-failure-injection",
	"phase17-security-boundary",
	"phase17-bugfix-spec-first",
	"phase17-document-drift",
}

var phase17FixtureFiles = []string{
	"manifest.json",
	"input/provider_targets.json",
	"input/conoha_vps_environment.json",
	"input/buildout_plan.json",
	"input/trial_operation.json",
	"input/simulation_matrix.json",
	"input/security_boundary.json",
	"expected/counters.json",
	"expected/runtime_flow.json",
	"expected/buildout.json",
	"expected/bugfix_loop.json",
	"records/trial_conoha.jsonl",
	"records/buildout.jsonl",
	"records/operation.jsonl",
	"records/bugfix.jsonl",
	"records/simulation.jsonl",
	"records/failure.jsonl",
	"records/security.jsonl",
	"records/document_drift.jsonl",
}

var phase17RecordFiles = []string{
	"records/trial_conoha.jsonl",
	"records/buildout.jsonl",
	"records/operation.jsonl",
	"records/bugfix.jsonl",
	"records/simulation.jsonl",
	"records/failure.jsonl",
	"records/security.jsonl",
	"records/document_drift.jsonl",
}

var phase17Counters = []string{
	"phase17_conoha_trial_operation_open_count",
	"phase17_buildout_open_count",
	"phase17_production_simulation_open_count",
	"phase17_install_update_rollback_open_count",
	"phase17_systemd_lifecycle_open_count",
	"phase17_runtime_flow_open_count",
	"phase17_failure_injection_open_count",
	"phase17_secret_leak_open_count",
	"phase17_destructive_operation_open_count",
	"phase17_xserver_future_misclassified_count",
	"phase17_document_drift_open_count",
	"phase17_trial_operation_open_count",
	"phase17_bugfix_spec_gap_count",
	"phase17_known_bug_open_count",
	"final_open_item_count",
}

var phase17BuildoutWorkUnits = []string{
	"provider-prerequisite",
	"access-baseline",
	"system-baseline",
	"release-acquisition",
	"install-bootstrap",
	"runtime-baseline",
	"operation-baseline",
	"update-rollback-baseline",
	"simulation-baseline",
	"monitoring-baseline",
}

var phase17ForbiddenProviderOperations = []string{
	"vps-create",
	"vps-delete",
	"plan-change",
	"disk-rebuild",
	"volume-create",
	"volume-delete",
	"volume-attach",
	"volume-detach",
	"firewall-lockout",
	"ssh-lockout",
}

var phase17TrialSteps = []string{
	"preflight",
	"install",
	"runtime-flow",
	"update",
	"rollback",
	"reboot-recovery",
	"cleanup",
}

var phase17OperationSteps = []string{
	"operation-start",
	"monitor",
	"issue-detect",
	"spec-general-update",
	"implementation-fix",
	"redeploy",
	"revalidate",
	"evidence-record",
	"known-bug-zero-check",
}

var phase17FailureClasses = []string{
	"disk full",
	"permission denied",
	"short write",
	"fsync failure",
	"rename failure",
	"process kill",
	"network refusal",
	"DNS failure",
	"SSH failure",
	"log write failure",
	"reboot recovery",
	"interrupted recovery",
}

var phase17AbortConditions = []string{
	"secret-boundary-failure",
	"destructive-operation-boundary-failure",
	"state-corruption",
	"service-unrecoverable",
	"known-critical-bug",
}

var phase17SecretClasses = []string{
	"provider-token",
	"admin-token",
	"ssh-private-key",
	"ip-specific-secret",
	"host-specific-secret",
}

var phase17AllowedSecretInputChannels = []string{
	"operator-local-file",
	"operator-stdin",
	"ssh-agent",
	"systemd-credential",
	"not-required",
}

var phase17ForbiddenSecretInputChannels = []string{
	"repository-file",
	"fixture-file",
	"command-argument",
	"environment-variable",
	"pr-body",
	"issue-comment",
	"chat-message",
	"stdout",
	"stderr",
	"log",
	"screen-shot",
}

var phase17RequiredCheckRecordFiles = map[string][]string{
	"phase17-conoha-vps-preflight":           {"records/trial_conoha.jsonl"},
	"phase17-trial-production-buildout":      {"records/buildout.jsonl"},
	"phase17-conoha-install-update-rollback": {"records/trial_conoha.jsonl"},
	"phase17-systemd-lifecycle":              {"records/trial_conoha.jsonl"},
	"phase17-runtime-flow":                   {"records/trial_conoha.jsonl", "records/operation.jsonl"},
	"phase17-trial-operation-loop":           {"records/operation.jsonl"},
	"phase17-production-simulation":          {"records/simulation.jsonl"},
	"phase17-failure-injection":              {"records/failure.jsonl"},
	"phase17-security-boundary":              {"records/security.jsonl"},
	"phase17-bugfix-spec-first":              {"records/bugfix.jsonl"},
	"phase17-document-drift":                 {"records/document_drift.jsonl"},
}

var phase17RequiredCheckCounterKeys = map[string][]string{
	"phase17-conoha-vps-preflight":           {"phase17_conoha_trial_operation_open_count"},
	"phase17-trial-production-buildout":      {"phase17_buildout_open_count"},
	"phase17-conoha-install-update-rollback": {"phase17_install_update_rollback_open_count"},
	"phase17-systemd-lifecycle":              {"phase17_systemd_lifecycle_open_count"},
	"phase17-runtime-flow":                   {"phase17_runtime_flow_open_count"},
	"phase17-trial-operation-loop":           {"phase17_trial_operation_open_count"},
	"phase17-production-simulation":          {"phase17_production_simulation_open_count"},
	"phase17-failure-injection":              {"phase17_failure_injection_open_count"},
	"phase17-security-boundary":              {"phase17_secret_leak_open_count", "phase17_destructive_operation_open_count"},
	"phase17-bugfix-spec-first":              {"phase17_bugfix_spec_gap_count", "phase17_known_bug_open_count"},
	"phase17-document-drift":                 {"phase17_document_drift_open_count", "phase17_xserver_future_misclassified_count"},
}

type phase17CountersFile struct {
	SchemaVersion int            `json:"schema_version"`
	Scope         string         `json:"scope"`
	Counters      map[string]int `json:"counters"`
}

type phase17ProviderTargetsFile struct {
	SchemaVersion  int                     `json:"schema_version"`
	Scope          string                  `json:"scope"`
	ProviderTarget []phase17ProviderTarget `json:"provider_targets"`
}

type phase17ProviderTarget struct {
	ProviderTarget    string  `json:"provider_target"`
	Classification    string  `json:"classification"`
	CompletionBlocker bool    `json:"completion_blocker"`
	OS                *string `json:"os"`
	MinimumPlanClass  *string `json:"minimum_plan_class"`
	MinimumRAMMB      *int    `json:"minimum_ram_mb"`
	ExcludedPlan      *string `json:"excluded_plan"`
	FutureRef         *string `json:"future_ref"`
}

type phase17EnvironmentFile struct {
	SchemaVersion        int                   `json:"schema_version"`
	Scope                string                `json:"scope"`
	ProviderTarget       string                `json:"provider_target"`
	ValidationMode       string                `json:"validation_mode"`
	OS                   string                `json:"os"`
	MinimumPlanClass     string                `json:"minimum_plan_class"`
	MinimumRAMMB         int                   `json:"minimum_ram_mb"`
	ExcludedPlan         string                `json:"excluded_plan"`
	EnvironmentIdentity  phase17EnvironmentID  `json:"environment_identity"`
	Capabilities         []phase17ResultRecord `json:"capabilities"`
	StateDirectoryParent phase17StateParent    `json:"state_directory_parent"`
	Network              []phase17ResultRecord `json:"network"`
	RequiredCommands     []phase17ResultRecord `json:"required_commands"`
}

type phase17EnvironmentID struct {
	ProviderPlanLabel   string `json:"provider_plan_label"`
	ProviderRegionLabel string `json:"provider_region_label"`
	VPSInstanceLabel    string `json:"vps_instance_label"`
	PublicEndpointLabel string `json:"public_endpoint_label"`
	MetadataPolicy      string `json:"metadata_policy"`
}

type phase17ResultRecord struct {
	Name     string `json:"name"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Result   string `json:"result"`
}

type phase17StateParent struct {
	Path          string `json:"path"`
	OwnerExpected string `json:"owner_expected"`
	ModeExpected  string `json:"mode_expected"`
	Result        string `json:"result"`
}

type phase17BuildoutPlanFile struct {
	SchemaVersion              int      `json:"schema_version"`
	Scope                      string   `json:"scope"`
	ProviderTarget             string   `json:"provider_target"`
	ValidationMode             string   `json:"validation_mode"`
	RequiredCheck              string   `json:"required_check"`
	WorkUnits                  []string `json:"work_units"`
	ForbiddenProviderOperation []string `json:"forbidden_provider_operations"`
}

type phase17TrialOperationFile struct {
	SchemaVersion                 int      `json:"schema_version"`
	Scope                         string   `json:"scope"`
	ProviderTarget                string   `json:"provider_target"`
	ValidationMode                string   `json:"validation_mode"`
	OperationSteps                []string `json:"operation_steps"`
	OperationWindow               string   `json:"operation_window"`
	MonitorIntervalSeconds        int      `json:"monitor_interval_seconds"`
	RequiredMonitorSampleCount    int      `json:"required_monitor_sample_count"`
	MaxConsecutiveMonitorFailures int      `json:"max_consecutive_monitor_failures"`
	AbortConditions               []string `json:"abort_conditions"`
	BugfixSpecFirstRequired       bool     `json:"bugfix_spec_first_required"`
	RevalidationRequired          bool     `json:"revalidation_required"`
	KnownBugZeroRequired          bool     `json:"known_bug_zero_required"`
}

type phase17SimulationMatrixFile struct {
	SchemaVersion                int                   `json:"schema_version"`
	Scope                        string                `json:"scope"`
	ProviderTarget               string                `json:"provider_target"`
	ValidationMode               string                `json:"validation_mode"`
	FailureClasses               []phase17FailureClass `json:"failure_classes"`
	DestructiveOperationBoundary []string              `json:"destructive_operation_boundaries"`
	ProviderExecutionAllowed     bool                  `json:"provider_execution_allowed"`
}

type phase17FailureClass struct {
	FailureClass      string `json:"failure_class"`
	InjectionLocation string `json:"injection_location"`
	ExpectedResult    string `json:"expected_result"`
	RecoveryRequired  bool   `json:"recovery_required"`
}

type phase17SecurityBoundaryFile struct {
	SchemaVersion                  int      `json:"schema_version"`
	Scope                          string   `json:"scope"`
	ProviderTarget                 string   `json:"provider_target"`
	ValidationMode                 string   `json:"validation_mode"`
	RequiredCheck                  string   `json:"required_check"`
	SecretClasses                  []string `json:"secret_classes"`
	AllowedSecretInputChannels     []string `json:"allowed_secret_input_channels"`
	ForbiddenSecretInputChannels   []string `json:"forbidden_secret_input_channels"`
	CredentialStoragePolicy        string   `json:"credential_storage_policy"`
	SecretReferencePolicy          string   `json:"secret_reference_policy"`
	DestructiveOperationBoundaries []string `json:"destructive_operation_boundaries"`
	MaskRequired                   bool     `json:"mask_required"`
}

type phase17RuntimeFlowFile struct {
	SchemaVersion      int      `json:"schema_version"`
	Scope              string   `json:"scope"`
	ProviderTarget     string   `json:"provider_target"`
	ValidationMode     string   `json:"validation_mode"`
	ExpectedSteps      []string `json:"expected_steps"`
	CounterConnections []string `json:"counter_connections"`
}

type phase17BuildoutExpectedFile struct {
	SchemaVersion                   int      `json:"schema_version"`
	Scope                           string   `json:"scope"`
	RequiredCheck                   string   `json:"required_check"`
	ExpectedWorkUnitOrder           []string `json:"expected_work_unit_order"`
	AllowedProviderOperationClasses []string `json:"allowed_provider_operation_classes"`
	ForbiddenProviderOperations     []string `json:"forbidden_provider_operations"`
	CounterKey                      string   `json:"counter_key"`
	CounterValue                    int      `json:"counter_value"`
}

type phase17BugfixLoopFile struct {
	SchemaVersion                   int      `json:"schema_version"`
	Scope                           string   `json:"scope"`
	ExpectedOperationOrder          []string `json:"expected_operation_order"`
	OperationWindow                 string   `json:"operation_window"`
	MonitorIntervalSeconds          int      `json:"monitor_interval_seconds"`
	RequiredMonitorSampleCount      int      `json:"required_monitor_sample_count"`
	MaxConsecutiveMonitorFailures   int      `json:"max_consecutive_monitor_failures"`
	AbortConditions                 []string `json:"abort_conditions"`
	SpecUpdateRequired              bool     `json:"spec_update_required"`
	ResponsibilitySourceRefRequired bool     `json:"responsibility_source_ref_required"`
	RevalidationRequired            bool     `json:"revalidation_required"`
	AllowedKnownBugStatuses         []string `json:"allowed_known_bug_statuses"`
	ForbiddenKnownBugStatuses       []string `json:"forbidden_known_bug_statuses"`
	CounterConnections              []string `json:"counter_connections"`
}

type phase17EvidenceRecord struct {
	RecordID                           string  `json:"record_id"`
	Scope                              string  `json:"scope"`
	ProviderTarget                     string  `json:"provider_target"`
	ValidationMode                     string  `json:"validation_mode"`
	CheckName                          string  `json:"check_name"`
	SourceRef                          string  `json:"source_ref"`
	Result                             string  `json:"result"`
	CounterKey                         string  `json:"counter_key"`
	StartedAt                          *string `json:"started_at"`
	EndedAt                            *string `json:"ended_at"`
	StdoutSHA256                       *string `json:"stdout_sha256"`
	StderrSHA256                       *string `json:"stderr_sha256"`
	SecretScanResult                   string  `json:"secret_scan_result"`
	DestructiveOperationResult         string  `json:"destructive_operation_result"`
	EvidenceRef                        string  `json:"evidence_ref"`
	TrialStep                          string  `json:"trial_step,omitempty"`
	EnvironmentRef                     string  `json:"environment_ref,omitempty"`
	EnvironmentIdentityPolicy          string  `json:"environment_identity_policy,omitempty"`
	RuntimeFlowRef                     string  `json:"runtime_flow_ref,omitempty"`
	ServiceState                       string  `json:"service_state,omitempty"`
	StateDigest                        string  `json:"state_digest,omitempty"`
	LogDigest                          string  `json:"log_digest,omitempty"`
	RollbackState                      string  `json:"rollback_state,omitempty"`
	CleanupResult                      string  `json:"cleanup_result,omitempty"`
	BuildoutStep                       string  `json:"buildout_step,omitempty"`
	Sequence                           *int    `json:"sequence,omitempty"`
	WorkUnitResult                     string  `json:"work_unit_result,omitempty"`
	ProviderOperationClass             string  `json:"provider_operation_class,omitempty"`
	ForbiddenProviderOperationDetected *bool   `json:"forbidden_provider_operation_detected,omitempty"`
	PreconditionRef                    string  `json:"precondition_ref,omitempty"`
	ActualRef                          string  `json:"actual_ref,omitempty"`
	ExpectedRef                        string  `json:"expected_ref,omitempty"`
	OperationStep                      string  `json:"operation_step,omitempty"`
	OperationWindow                    string  `json:"operation_window,omitempty"`
	MonitorIntervalSeconds             *int    `json:"monitor_interval_seconds,omitempty"`
	MonitorSampleIndex                 *int    `json:"monitor_sample_index,omitempty"`
	MonitorObservedAt                  *string `json:"monitor_observed_at,omitempty"`
	AbortCondition                     string  `json:"abort_condition,omitempty"`
	AbortResult                        string  `json:"abort_result,omitempty"`
	HealthResult                       string  `json:"health_result,omitempty"`
	StateWriteResult                   string  `json:"state_write_result,omitempty"`
	LogWriteResult                     string  `json:"log_write_result,omitempty"`
	RevalidationRef                    string  `json:"revalidation_ref,omitempty"`
	BugID                              string  `json:"bug_id,omitempty"`
	SpecUpdateRef                      string  `json:"spec_update_ref,omitempty"`
	ResponsibilitySourceRef            string  `json:"responsibility_source_ref,omitempty"`
	ImplementationRef                  string  `json:"implementation_ref,omitempty"`
	KnownBugStatus                     string  `json:"known_bug_status,omitempty"`
	SimulationCase                     string  `json:"simulation_case,omitempty"`
	FailureClass                       string  `json:"failure_class,omitempty"`
	DestructiveOperationToken          string  `json:"destructive_operation_token,omitempty"`
	ProviderExecutionResult            string  `json:"provider_execution_result,omitempty"`
	RecoveryResult                     string  `json:"recovery_result,omitempty"`
	SimulationEnvironmentRef           string  `json:"simulation_environment_ref,omitempty"`
	InjectionMethod                    string  `json:"injection_method,omitempty"`
	ExpectedFailure                    string  `json:"expected_failure,omitempty"`
	ActualFailure                      string  `json:"actual_failure,omitempty"`
	SilentSuccessDetected              *bool   `json:"silent_success_detected,omitempty"`
	SecretClass                        string  `json:"secret_class,omitempty"`
	InputChannel                       string  `json:"input_channel,omitempty"`
	RedactedReference                  string  `json:"redacted_reference,omitempty"`
	SecretReferencePolicy              string  `json:"secret_reference_policy,omitempty"`
	SecretValuePresent                 *bool   `json:"secret_value_present,omitempty"`
	ScanTarget                         string  `json:"scan_target,omitempty"`
	MaskResult                         string  `json:"mask_result,omitempty"`
	CredentialStorageResult            string  `json:"credential_storage_result,omitempty"`
	BoundaryResult                     string  `json:"boundary_result,omitempty"`
	DocumentRef                        string  `json:"document_ref,omitempty"`
	AnchorRef                          string  `json:"anchor_ref,omitempty"`
	FixturePathRef                     string  `json:"fixture_path_ref,omitempty"`
	WorkflowPathRef                    string  `json:"workflow_path_ref,omitempty"`
	ExpectedState                      string  `json:"expected_state,omitempty"`
	ActualState                        string  `json:"actual_state,omitempty"`
	DriftResult                        string  `json:"drift_result,omitempty"`
}

func phase17AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase17/production-validation/manifest.json" &&
		name == phase17Scope
}

func phase18AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase18/vps-pull-bootstrap/manifest.json" &&
		name == phase18Scope
}

func phase17RequireFixtureFiles(t *testing.T, root string) {
	t.Helper()

	allowed := map[string]bool{}
	for _, rel := range phase17FixtureFiles {
		allowed[rel] = true
		path := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Lstat(path)
		if err != nil {
			t.Fatalf("Phase 17 fixture path missing %s: %v", rel, err)
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			t.Fatalf("Phase 17 fixture path must be a regular file: %s", rel)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if len(data) == 0 || data[len(data)-1] != '\n' {
			t.Fatalf("%s must end with LF", path)
		}
		if strings.HasSuffix(rel, ".json") {
			phase16RequireNoDuplicateJSONKeys(t, path, data)
		}
		phase17RequireNoSecretLikeFixtureContent(t, path, data)
	}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root || entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !allowed[rel] {
			t.Fatalf("Phase 17 fixture contains unknown file %s", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk Phase 17 fixture: %v", err)
	}
}

func phase17ReadCounters(t *testing.T, path string) map[string]int {
	t.Helper()

	var file phase17CountersFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase17Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase17Scope)
	}
	return file.Counters
}

func phase17RequireCountersClosed(t *testing.T, counters map[string]int) {
	t.Helper()

	if len(counters) != len(phase17Counters) {
		t.Fatalf("Phase 17 counters must include exactly %d keys, got %d", len(phase17Counters), len(counters))
	}
	var sum int
	for _, key := range phase17Counters {
		value, ok := counters[key]
		if !ok {
			t.Fatalf("Phase 17 counters missing %s", key)
		}
		if value < 0 {
			t.Fatalf("Phase 17 counter %s must not be negative", key)
		}
		if key != "final_open_item_count" {
			sum += value
		}
	}
	if counters["final_open_item_count"] != sum {
		t.Fatalf("Phase 17 final_open_item_count must equal sum of open counters: want %d got %d", sum, counters["final_open_item_count"])
	}
	if counters["final_open_item_count"] != 0 {
		t.Fatalf("Phase 17 closure requires final_open_item_count=0, got %d", counters["final_open_item_count"])
	}
}

func phase17RequireInputsAndExpected(t *testing.T, root string) {
	t.Helper()

	var targets phase17ProviderTargetsFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/provider_targets.json"), &targets)
	if targets.SchemaVersion != 1 || targets.Scope != phase17Scope || len(targets.ProviderTarget) != 2 {
		t.Fatalf("Phase 17 provider_targets.json schema mismatch")
	}
	for i, target := range targets.ProviderTarget {
		if target.ProviderTarget != phase17ProviderTargets[i] {
			t.Fatalf("Phase 17 provider target order mismatch")
		}
		if target.ProviderTarget == "conoha-vps-primary" {
			if !target.CompletionBlocker || target.OS == nil || *target.OS != "ubuntu-server-24.04-lts-64bit" || target.MinimumPlanClass == nil || *target.MinimumPlanClass != "conoha-vps-1gb-memory-class" || target.MinimumRAMMB == nil || *target.MinimumRAMMB != 1024 {
				t.Fatalf("ConoHa provider target must fix OS and minimum plan class: %+v", target)
			}
		}
		if target.ProviderTarget == "xserver-vps-future" && (target.CompletionBlocker || target.FutureRef == nil || *target.FutureRef == "") {
			t.Fatalf("xserver-vps-future must be future_plan and non-blocking: %+v", target)
		}
	}

	var env phase17EnvironmentFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/conoha_vps_environment.json"), &env)
	if env.SchemaVersion != 1 || env.Scope != phase17Scope || env.ProviderTarget != "conoha-vps-primary" || env.ValidationMode != "trial-production-conoha-vps" || env.OS != "ubuntu-server-24.04-lts-64bit" || env.MinimumPlanClass != "conoha-vps-1gb-memory-class" || env.MinimumRAMMB != 1024 || env.ExcludedPlan != "512mb" {
		t.Fatalf("Phase 17 ConoHa environment fixed values mismatch: %+v", env)
	}
	if env.EnvironmentIdentity.MetadataPolicy != "opaque-non-secret-labels" || phase17ContainsSecretLikeValue(env.EnvironmentIdentity.ProviderPlanLabel, env.EnvironmentIdentity.ProviderRegionLabel, env.EnvironmentIdentity.VPSInstanceLabel, env.EnvironmentIdentity.PublicEndpointLabel) {
		t.Fatalf("Phase 17 environment identity must be opaque and non-secret: %+v", env.EnvironmentIdentity)
	}
	if env.StateDirectoryParent.ModeExpected != "0700" || len(env.Capabilities) == 0 || len(env.Network) == 0 || len(env.RequiredCommands) == 0 {
		t.Fatalf("Phase 17 environment must include state parent, capabilities, network, and commands")
	}

	var buildout phase17BuildoutPlanFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/buildout_plan.json"), &buildout)
	if buildout.SchemaVersion != 1 || buildout.Scope != phase17Scope || buildout.RequiredCheck != "phase17-trial-production-buildout" || !phase11StringSlicesEqual(buildout.WorkUnits, phase17BuildoutWorkUnits) || !phase11StringSlicesEqual(buildout.ForbiddenProviderOperation, phase17ForbiddenProviderOperations) {
		t.Fatalf("Phase 17 buildout plan mismatch: %+v", buildout)
	}

	var trial phase17TrialOperationFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/trial_operation.json"), &trial)
	if trial.SchemaVersion != 1 || trial.Scope != phase17Scope || !phase11StringSlicesEqual(trial.OperationSteps, phase17OperationSteps) || trial.OperationWindow != phase17OperationWindow || trial.MonitorIntervalSeconds != phase17MonitorIntervalSeconds || trial.RequiredMonitorSampleCount != phase17RequiredMonitorSampleCount || trial.MaxConsecutiveMonitorFailures != phase17MaxConsecutiveMonitorFailures || !phase11StringSlicesEqual(trial.AbortConditions, phase17AbortConditions) || !trial.BugfixSpecFirstRequired || !trial.RevalidationRequired || !trial.KnownBugZeroRequired {
		t.Fatalf("Phase 17 trial operation input mismatch: %+v", trial)
	}

	var simulation phase17SimulationMatrixFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/simulation_matrix.json"), &simulation)
	if simulation.SchemaVersion != 1 || simulation.Scope != phase17Scope || simulation.ProviderExecutionAllowed || !phase11StringSlicesEqual(simulation.DestructiveOperationBoundary, phase17ForbiddenProviderOperations) || len(simulation.FailureClasses) != len(phase17FailureClasses) {
		t.Fatalf("Phase 17 simulation matrix mismatch: %+v", simulation)
	}
	for i, failure := range simulation.FailureClasses {
		if failure.FailureClass != phase17FailureClasses[i] || failure.ExpectedResult == "" || !failure.RecoveryRequired {
			t.Fatalf("Phase 17 failure class mismatch at %d: %+v", i, failure)
		}
	}

	var security phase17SecurityBoundaryFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/security_boundary.json"), &security)
	if security.SchemaVersion != 1 || security.Scope != phase17Scope || security.RequiredCheck != "phase17-security-boundary" || !phase11StringSlicesEqual(security.SecretClasses, phase17SecretClasses) || !phase11StringSlicesEqual(security.AllowedSecretInputChannels, phase17AllowedSecretInputChannels) || !phase11StringSlicesEqual(security.ForbiddenSecretInputChannels, phase17ForbiddenSecretInputChannels) || security.CredentialStoragePolicy != "repository-forbidden" || security.SecretReferencePolicy != "metadata-only" || !phase11StringSlicesEqual(security.DestructiveOperationBoundaries, phase17ForbiddenProviderOperations) || !security.MaskRequired {
		t.Fatalf("Phase 17 security boundary mismatch: %+v", security)
	}

	var runtime phase17RuntimeFlowFile
	phase16ReadStrictJSON(t, filepath.Join(root, "expected/runtime_flow.json"), &runtime)
	if runtime.SchemaVersion != 1 || runtime.Scope != phase17Scope || !phase11StringSlicesEqual(runtime.ExpectedSteps, phase17TrialSteps) {
		t.Fatalf("Phase 17 runtime flow expected steps mismatch: %+v", runtime)
	}

	var expectedBuildout phase17BuildoutExpectedFile
	phase16ReadStrictJSON(t, filepath.Join(root, "expected/buildout.json"), &expectedBuildout)
	if expectedBuildout.SchemaVersion != 1 || expectedBuildout.Scope != phase17Scope || expectedBuildout.RequiredCheck != "phase17-trial-production-buildout" || !phase11StringSlicesEqual(expectedBuildout.ExpectedWorkUnitOrder, phase17BuildoutWorkUnits) || !phase11StringSlicesEqual(expectedBuildout.ForbiddenProviderOperations, phase17ForbiddenProviderOperations) || expectedBuildout.CounterKey != "phase17_buildout_open_count" {
		t.Fatalf("Phase 17 buildout expected mismatch: %+v", expectedBuildout)
	}

	var bugfix phase17BugfixLoopFile
	phase16ReadStrictJSON(t, filepath.Join(root, "expected/bugfix_loop.json"), &bugfix)
	if bugfix.SchemaVersion != 1 || bugfix.Scope != phase17Scope || !phase11StringSlicesEqual(bugfix.ExpectedOperationOrder, phase17OperationSteps) || bugfix.OperationWindow != trial.OperationWindow || bugfix.MonitorIntervalSeconds != trial.MonitorIntervalSeconds || bugfix.RequiredMonitorSampleCount != trial.RequiredMonitorSampleCount || bugfix.MaxConsecutiveMonitorFailures != trial.MaxConsecutiveMonitorFailures || !phase11StringSlicesEqual(bugfix.AbortConditions, trial.AbortConditions) || !bugfix.SpecUpdateRequired || !bugfix.ResponsibilitySourceRefRequired || !bugfix.RevalidationRequired {
		t.Fatalf("Phase 17 bugfix loop expected mismatch: %+v", bugfix)
	}
}

func phase17ReadRecords(t *testing.T, root string) map[string][]phase17EvidenceRecord {
	t.Helper()

	records := map[string][]phase17EvidenceRecord{}
	for _, rel := range phase17RecordFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
		if len(lines) == 0 {
			t.Fatalf("%s must include at least one JSONL record", path)
		}
		seen := map[string]bool{}
		for index, line := range lines {
			if strings.TrimSpace(line) == "" {
				t.Fatalf("%s line %d must not be empty", path, index+1)
			}
			phase16RequireNoDuplicateJSONKeys(t, fmt.Sprintf("%s line %d", path, index+1), []byte(line))
			var record phase17EvidenceRecord
			decoder := json.NewDecoder(strings.NewReader(line))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&record); err != nil {
				t.Fatalf("decode %s line %d: %v", path, index+1, err)
			}
			phase17RequireRecordCommon(t, rel, index+1, record)
			if seen[record.RecordID] {
				t.Fatalf("%s repeats record_id %s", path, record.RecordID)
			}
			seen[record.RecordID] = true
			phase17RequireRecordSpecific(t, rel, index+1, record)
			records[rel] = append(records[rel], record)
		}
	}
	return records
}

func phase17RequireRecordCommon(t *testing.T, rel string, line int, record phase17EvidenceRecord) {
	t.Helper()

	prefix := "phase17.record." + strings.TrimSuffix(filepath.Base(rel), ".jsonl") + "."
	if record.RecordID == "" || !strings.HasPrefix(record.RecordID, prefix) || record.Scope != phase17Scope || record.ProviderTarget == "" || record.ValidationMode == "" || record.CheckName == "" || record.SourceRef == "" || record.Result == "" || record.CounterKey == "" || record.SecretScanResult == "" || record.DestructiveOperationResult == "" || record.EvidenceRef == "" {
		t.Fatalf("%s line %d common fields invalid: %+v", rel, line, record)
	}
	if !phase17RequiredCheckAllowed(record.CheckName) {
		t.Fatalf("%s line %d unknown check_name %s", rel, line, record.CheckName)
	}
	if !phase17CounterAllowed(record.CounterKey) {
		t.Fatalf("%s line %d unknown counter %s", rel, line, record.CounterKey)
	}
	switch record.Result {
	case "passed", "failed", "not_applicable", "future_plan", "rejected":
	default:
		t.Fatalf("%s line %d invalid result %s", rel, line, record.Result)
	}
	if record.Result == "failed" || record.Result == "rejected" {
		t.Fatalf("%s line %d must not leave failed or rejected result", rel, line)
	}
	if record.Result == "passed" && (record.StartedAt == nil || *record.StartedAt == "" || record.EndedAt == nil || *record.EndedAt == "") {
		t.Fatalf("%s line %d passed record must include started_at and ended_at", rel, line)
	}
	for label, value := range map[string]*string{"stdout_sha256": record.StdoutSHA256, "stderr_sha256": record.StderrSHA256} {
		if value != nil && !phase17IsLowerHexSHA256(*value) {
			t.Fatalf("%s line %d %s must be lowercase SHA-256 hex or null: %v", rel, line, label, value)
		}
	}
	switch record.SecretScanResult {
	case "passed", "not_applicable":
	default:
		t.Fatalf("%s line %d invalid secret_scan_result %s", rel, line, record.SecretScanResult)
	}
	switch record.DestructiveOperationResult {
	case "not_applicable", "not_executed", "blocked":
	default:
		t.Fatalf("%s line %d invalid destructive_operation_result %s", rel, line, record.DestructiveOperationResult)
	}
	if phase17ContainsSecretLikeValue(record.EvidenceRef) {
		t.Fatalf("%s line %d evidence_ref must not contain secret-like or endpoint-like value: %s", rel, line, record.EvidenceRef)
	}
}

func phase17RequireRecordSpecific(t *testing.T, rel string, line int, record phase17EvidenceRecord) {
	t.Helper()

	switch rel {
	case "records/trial_conoha.jsonl":
		if record.TrialStep == "" || record.EnvironmentRef == "" || record.EnvironmentIdentityPolicy != "opaque-non-secret-labels" || record.RuntimeFlowRef == "" || record.ServiceState == "" || record.StateDigest == "" || record.LogDigest == "" || record.RollbackState == "" || record.CleanupResult == "" {
			t.Fatalf("%s line %d missing trial ConoHa fields: %+v", rel, line, record)
		}
	case "records/buildout.jsonl":
		if record.BuildoutStep == "" || record.Sequence == nil || record.WorkUnitResult == "" || record.ProviderOperationClass == "" || record.ForbiddenProviderOperationDetected == nil || *record.ForbiddenProviderOperationDetected || record.PreconditionRef == "" || record.ActualRef == "" || record.ExpectedRef == "" {
			t.Fatalf("%s line %d missing buildout fields or detected forbidden operation: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.BuildoutStep, phase17BuildoutWorkUnits) {
			t.Fatalf("%s line %d uses unknown buildout_step %s", rel, line, record.BuildoutStep)
		}
	case "records/operation.jsonl":
		if record.OperationStep == "" || record.OperationWindow != phase17OperationWindow || record.MonitorIntervalSeconds == nil || *record.MonitorIntervalSeconds != phase17MonitorIntervalSeconds || record.AbortCondition == "" || record.AbortResult == "" || record.ServiceState == "" || record.HealthResult == "" || record.StateWriteResult == "" || record.LogWriteResult == "" || record.RevalidationRef == "" {
			t.Fatalf("%s line %d missing operation fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.OperationStep, phase17OperationSteps) {
			t.Fatalf("%s line %d uses unknown operation_step %s", rel, line, record.OperationStep)
		}
		if record.OperationStep == "monitor" {
			if record.MonitorSampleIndex == nil || *record.MonitorSampleIndex != 1 || record.MonitorObservedAt == nil || *record.MonitorObservedAt == "" {
				t.Fatalf("%s line %d monitor baseline sample must include sample index 1 and observed_at: %+v", rel, line, record)
			}
		}
		if record.OperationStep != "monitor" && (record.MonitorSampleIndex != nil || record.MonitorObservedAt != nil) {
			t.Fatalf("%s line %d non-monitor operation record must not claim monitor sample fields: %+v", rel, line, record)
		}
	case "records/bugfix.jsonl":
		if record.BugID == "" || record.SpecUpdateRef == "" || record.ResponsibilitySourceRef == "" || record.ImplementationRef == "" || record.RevalidationRef == "" || record.KnownBugStatus == "" {
			t.Fatalf("%s line %d missing bugfix fields: %+v", rel, line, record)
		}
		switch record.KnownBugStatus {
		case "closed", "future_plan":
		default:
			t.Fatalf("%s line %d invalid known_bug_status %s", rel, line, record.KnownBugStatus)
		}
	case "records/simulation.jsonl":
		if record.SimulationCase == "" || record.FailureClass == "" || record.DestructiveOperationToken == "" || record.ProviderExecutionResult != "not_executed_on_provider" || record.RecoveryResult == "" || record.SimulationEnvironmentRef == "" {
			t.Fatalf("%s line %d missing simulation fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.FailureClass, phase17FailureClasses) {
			t.Fatalf("%s line %d uses unknown simulation failure_class %s", rel, line, record.FailureClass)
		}
		if record.DestructiveOperationToken != "not_applicable" && !phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) {
			t.Fatalf("%s line %d uses unknown destructive_operation_token %s", rel, line, record.DestructiveOperationToken)
		}
	case "records/failure.jsonl":
		if record.FailureClass == "" || record.InjectionMethod == "" || record.ExpectedFailure == "" || record.ActualFailure == "" || record.RecoveryResult == "" || record.SilentSuccessDetected == nil || *record.SilentSuccessDetected {
			t.Fatalf("%s line %d missing failure fields or silent success: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.FailureClass, phase17FailureClasses) {
			t.Fatalf("%s line %d uses unknown failure_class %s", rel, line, record.FailureClass)
		}
	case "records/security.jsonl":
		if record.SecretClass == "" || record.InputChannel == "" || record.RedactedReference == "" || record.SecretReferencePolicy != "metadata-only" || record.SecretValuePresent == nil || *record.SecretValuePresent || record.ScanTarget == "" || record.MaskResult == "" || record.CredentialStorageResult != "not_stored" || record.DestructiveOperationToken == "" || record.BoundaryResult == "" {
			t.Fatalf("%s line %d missing security fields or leaked secret: %+v", rel, line, record)
		}
		phase17RequireSecurityRecordBoundary(t, rel, line, record)
	case "records/document_drift.jsonl":
		if record.DocumentRef == "" || record.AnchorRef == "" || record.FixturePathRef == "" || record.WorkflowPathRef == "" || record.ExpectedState == "" || record.ActualState == "" || record.DriftResult == "" {
			t.Fatalf("%s line %d missing document drift fields: %+v", rel, line, record)
		}
		switch record.DriftResult {
		case "no_drift", "future_plan_confirmed":
		default:
			t.Fatalf("%s line %d invalid drift_result %s", rel, line, record.DriftResult)
		}
	}
}

func phase17RequireRequiredCheckConnections(t *testing.T, records map[string][]phase17EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for rel, recordSet := range records {
		for _, record := range recordSet {
			if !phase17StringInSlice(rel, phase17RequiredCheckRecordFiles[record.CheckName]) {
				t.Fatalf("Phase 17 record %s places check %s in invalid record file %s", record.RecordID, record.CheckName, rel)
			}
			if !phase17StringInSlice(record.CounterKey, phase17RequiredCheckCounterKeys[record.CheckName]) {
				t.Fatalf("Phase 17 record %s connects check %s to invalid counter %s", record.RecordID, record.CheckName, record.CounterKey)
			}
			seen[record.CheckName] = true
		}
	}
	for _, check := range phase17RequiredChecks {
		if !seen[check] {
			t.Fatalf("Phase 17 required check %s is not connected to records", check)
		}
	}
}

func phase17RequireBoundaryRecords(t *testing.T, counters map[string]int, records map[string][]phase17EvidenceRecord) {
	t.Helper()

	phase17RequireBuildoutOrder(t, records["records/buildout.jsonl"], counters["phase17_buildout_open_count"])
	phase17RequireOperationStepCoverage(t, records["records/operation.jsonl"])

	simulationFailureClasses := map[string]bool{}
	for _, record := range records["records/simulation.jsonl"] {
		simulationFailureClasses[record.FailureClass] = true
		if record.ProviderExecutionResult != "not_executed_on_provider" || record.DestructiveOperationResult != "not_executed" {
			t.Fatalf("Phase 17 simulation must not execute destructive operation on provider: %+v", record)
		}
	}
	for _, failureClass := range phase17FailureClasses {
		if !simulationFailureClasses[failureClass] {
			t.Fatalf("Phase 17 failure class missing from simulation records: %s", failureClass)
		}
	}

	failureClasses := map[string]bool{}
	for _, record := range records["records/failure.jsonl"] {
		failureClasses[record.FailureClass] = true
	}
	for _, failureClass := range phase17FailureClasses {
		if !failureClasses[failureClass] {
			t.Fatalf("Phase 17 failure class missing from failure records: %s", failureClass)
		}
	}

	var blockedForbiddenSecretInput bool
	var blockedDestructiveOperation bool
	for _, record := range records["records/security.jsonl"] {
		if phase17StringInSlice(record.InputChannel, phase17ForbiddenSecretInputChannels) && record.BoundaryResult == "blocked" && !*record.SecretValuePresent {
			blockedForbiddenSecretInput = true
		}
		if phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) && record.DestructiveOperationResult == "blocked" && record.BoundaryResult == "blocked" {
			blockedDestructiveOperation = true
		}
	}
	if !blockedForbiddenSecretInput {
		t.Fatalf("Phase 17 security records must include a blocked forbidden secret input channel")
	}
	if !blockedDestructiveOperation {
		t.Fatalf("Phase 17 security records must include a blocked destructive provider operation")
	}

	var futurePlan bool
	for _, record := range records["records/document_drift.jsonl"] {
		if record.ProviderTarget == "xserver-vps-future" && record.Result == "future_plan" && record.DriftResult == "future_plan_confirmed" {
			futurePlan = true
		}
	}
	if !futurePlan {
		t.Fatalf("Phase 17 document drift records must keep xserver-vps-future as future_plan")
	}

	if counters["phase17_xserver_future_misclassified_count"] != 0 {
		t.Fatalf("Phase 17 xserver-vps-future must remain non-blocking with counter 0")
	}
}

func phase17RequireBuildoutOrder(t *testing.T, records []phase17EvidenceRecord, _ int) {
	t.Helper()

	if len(records) != len(phase17BuildoutWorkUnits) {
		t.Fatalf("Phase 17 buildout records must contain exactly %d work units, got %d", len(phase17BuildoutWorkUnits), len(records))
	}
	seen := map[string]bool{}
	for index, record := range records {
		wantStep := phase17BuildoutWorkUnits[index]
		wantSequence := index + 1
		if record.BuildoutStep != wantStep || record.Sequence == nil || *record.Sequence != wantSequence {
			t.Fatalf("Phase 17 buildout record %d order mismatch: want %s/%d got %+v", index+1, wantStep, wantSequence, record)
		}
		if seen[record.BuildoutStep] {
			t.Fatalf("Phase 17 buildout repeats work unit %s", record.BuildoutStep)
		}
		seen[record.BuildoutStep] = true
		if record.Result != "passed" || record.WorkUnitResult != "passed" || record.CounterKey != "phase17_buildout_open_count" {
			t.Fatalf("Phase 17 buildout work unit must be closed before buildout counter can be 0: %+v", record)
		}
	}
}

func phase17RequireOperationStepCoverage(t *testing.T, records []phase17EvidenceRecord) {
	t.Helper()

	positions := map[string]int{}
	for index, record := range records {
		if record.CheckName != "phase17-trial-operation-loop" {
			continue
		}
		if _, exists := positions[record.OperationStep]; !exists {
			positions[record.OperationStep] = index
		}
	}
	previous := -1
	for _, step := range phase17OperationSteps {
		position, exists := positions[step]
		if !exists {
			t.Fatalf("Phase 17 operation records must include placeholder or evidence for step %s", step)
		}
		if position <= previous {
			t.Fatalf("Phase 17 operation records must keep operation step order; step %s is out of order", step)
		}
		previous = position
	}
}

func phase17RequireSecurityRecordBoundary(t *testing.T, rel string, line int, record phase17EvidenceRecord) {
	t.Helper()

	if !phase17StringInSlice(record.SecretClass, phase17SecretClasses) {
		t.Fatalf("%s line %d unknown secret_class %s", rel, line, record.SecretClass)
	}
	allowedInput := phase17StringInSlice(record.InputChannel, phase17AllowedSecretInputChannels)
	forbiddenInput := phase17StringInSlice(record.InputChannel, phase17ForbiddenSecretInputChannels)
	if allowedInput == forbiddenInput {
		t.Fatalf("%s line %d input_channel must be in exactly one allow/deny set: %s", rel, line, record.InputChannel)
	}
	if phase17ContainsSecretLikeValue(record.RedactedReference, record.ScanTarget) {
		t.Fatalf("%s line %d security metadata must remain opaque and non-secret: %+v", rel, line, record)
	}
	switch record.MaskResult {
	case "masked", "not_present", "not_applicable":
	default:
		t.Fatalf("%s line %d invalid mask_result %s", rel, line, record.MaskResult)
	}
	if record.DestructiveOperationToken != "not_applicable" && !phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) {
		t.Fatalf("%s line %d unknown destructive operation token %s", rel, line, record.DestructiveOperationToken)
	}
	switch record.BoundaryResult {
	case "passed", "blocked", "not_applicable":
	default:
		t.Fatalf("%s line %d invalid boundary_result %s", rel, line, record.BoundaryResult)
	}
	if forbiddenInput && (record.Result != "passed" || record.SecretScanResult != "not_applicable" || record.BoundaryResult != "blocked") {
		t.Fatalf("%s line %d forbidden input channel must be blocked as a passing boundary check: %+v", rel, line, record)
	}
	if phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) && (record.DestructiveOperationResult != "blocked" || record.BoundaryResult != "blocked") {
		t.Fatalf("%s line %d forbidden provider operation must be blocked: %+v", rel, line, record)
	}
}

func phase17RequireClosureEvidenceGates(t *testing.T, counters map[string]int, records map[string][]phase17EvidenceRecord) {
	t.Helper()

	if phase17CountMonitorSamples(records["records/operation.jsonl"]) < phase17RequiredMonitorSampleCount {
		t.Fatalf("Phase 17 trial operation closure requires %d baseline monitor sample", phase17RequiredMonitorSampleCount)
	}
	if phase17AnyPendingExternalEvidence(records) {
		t.Fatalf("Phase 17 closure must not leave pending external evidence")
	}
	if phase17RecordWithResult(records["records/operation.jsonl"], "not_applicable") {
		t.Fatalf("Phase 17 operation records must not leave not_applicable placeholders")
	}
	if phase17RecordWithResult(records["records/trial_conoha.jsonl"], "not_applicable") {
		t.Fatalf("Phase 17 trial ConoHa records must not leave not_applicable placeholders")
	}
}

func phase17RequireDerivedCounters(t *testing.T, counters map[string]int, records map[string][]phase17EvidenceRecord) {
	t.Helper()

	derived := phase17DeriveCounters(records)
	for _, key := range phase17Counters {
		if counters[key] != derived[key] {
			t.Fatalf("Phase 17 counter %s must be derived from records: want %d got %d", key, derived[key], counters[key])
		}
	}
}

func phase17DeriveCounters(records map[string][]phase17EvidenceRecord) map[string]int {
	derived := map[string]int{}
	for _, key := range phase17Counters {
		derived[key] = 0
	}

	for _, record := range records["records/trial_conoha.jsonl"] {
		if phase17TrialRecordOpen(record) {
			switch record.CheckName {
			case "phase17-conoha-vps-preflight":
				derived["phase17_conoha_trial_operation_open_count"] = 1
			case "phase17-conoha-install-update-rollback":
				derived["phase17_install_update_rollback_open_count"] = 1
			case "phase17-systemd-lifecycle":
				derived["phase17_systemd_lifecycle_open_count"] = 1
			case "phase17-runtime-flow":
				derived["phase17_runtime_flow_open_count"] = 1
			}
		}
	}
	if !phase17BuildoutClosed(records["records/buildout.jsonl"]) {
		derived["phase17_buildout_open_count"] = 1
	}
	if !phase17SimulationClosed(records["records/simulation.jsonl"]) {
		derived["phase17_production_simulation_open_count"] = 1
	}
	if !phase17FailureInjectionClosed(records["records/failure.jsonl"]) {
		derived["phase17_failure_injection_open_count"] = 1
	}
	if !phase17DocumentDriftClosed(records["records/document_drift.jsonl"]) {
		derived["phase17_document_drift_open_count"] = 1
	}
	if phase17XserverFutureMisclassified(records["records/document_drift.jsonl"]) {
		derived["phase17_xserver_future_misclassified_count"] = 1
	}
	if !phase17TrialOperationClosed(records["records/operation.jsonl"]) {
		derived["phase17_trial_operation_open_count"] = 1
	}
	for _, record := range records["records/security.jsonl"] {
		if phase17SecurityLeakOpen(record) {
			derived["phase17_secret_leak_open_count"] = 1
		}
		if phase17DestructiveOperationOpen(record) {
			derived["phase17_destructive_operation_open_count"] = 1
		}
	}
	for _, record := range records["records/bugfix.jsonl"] {
		if phase17BugfixSpecGapOpen(record) {
			derived["phase17_bugfix_spec_gap_count"] = 1
		}
		if phase17KnownBugOpen(record) {
			derived["phase17_known_bug_open_count"] = 1
		}
	}
	for _, key := range phase17Counters {
		if key != "final_open_item_count" {
			derived["final_open_item_count"] += derived[key]
		}
	}
	return derived
}

func phase17TrialRecordOpen(record phase17EvidenceRecord) bool {
	if record.Result != "passed" {
		return true
	}
	for _, value := range []string{
		record.ServiceState,
		record.StateDigest,
		record.LogDigest,
		record.RollbackState,
		record.CleanupResult,
	} {
		if value == "pending-external-evidence" || value == "not_started" || value == "" {
			return true
		}
	}
	return false
}

func phase17BuildoutClosed(records []phase17EvidenceRecord) bool {
	if len(records) != len(phase17BuildoutWorkUnits) {
		return false
	}
	for _, record := range records {
		if record.Result != "passed" || record.WorkUnitResult != "passed" || record.ForbiddenProviderOperationDetected == nil || *record.ForbiddenProviderOperationDetected {
			return false
		}
	}
	return true
}

func phase17SimulationClosed(records []phase17EvidenceRecord) bool {
	if len(records) != len(phase17FailureClasses) {
		return false
	}
	for _, record := range records {
		if record.Result != "passed" || record.ProviderExecutionResult != "not_executed_on_provider" || record.RecoveryResult == "" {
			return false
		}
	}
	return true
}

func phase17FailureInjectionClosed(records []phase17EvidenceRecord) bool {
	if len(records) != len(phase17FailureClasses) {
		return false
	}
	for _, record := range records {
		if record.Result != "passed" || record.ActualFailure == "" || record.RecoveryResult == "" || record.SilentSuccessDetected == nil || *record.SilentSuccessDetected {
			return false
		}
	}
	return true
}

func phase17DocumentDriftClosed(records []phase17EvidenceRecord) bool {
	if len(records) == 0 {
		return false
	}
	for _, record := range records {
		switch record.DriftResult {
		case "no_drift", "future_plan_confirmed":
		default:
			return false
		}
	}
	return true
}

func phase17XserverFutureMisclassified(records []phase17EvidenceRecord) bool {
	for _, record := range records {
		if record.ProviderTarget == "xserver-vps-future" {
			return record.Result != "future_plan" || record.ValidationMode != "future_plan" || record.DriftResult != "future_plan_confirmed"
		}
	}
	return true
}

func phase17TrialOperationClosed(records []phase17EvidenceRecord) bool {
	if phase17CountMonitorSamples(records) < phase17RequiredMonitorSampleCount {
		return false
	}
	for _, record := range records {
		if record.CheckName != "phase17-trial-operation-loop" {
			continue
		}
		if record.Result != "passed" || phase17OperationRecordPending(record) || record.AbortResult == "triggered" {
			return false
		}
	}
	return true
}

func phase17SecurityLeakOpen(record phase17EvidenceRecord) bool {
	if record.SecretScanResult == "failed" || record.SecretValuePresent == nil || *record.SecretValuePresent || record.CredentialStorageResult != "not_stored" {
		return true
	}
	return record.BoundaryResult == "failed" || record.MaskResult == ""
}

func phase17DestructiveOperationOpen(record phase17EvidenceRecord) bool {
	if record.DestructiveOperationResult == "boundary_failed" || record.BoundaryResult == "failed" {
		return true
	}
	if phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) {
		return record.DestructiveOperationResult != "blocked" || record.BoundaryResult != "blocked"
	}
	return false
}

func phase17BugfixSpecGapOpen(record phase17EvidenceRecord) bool {
	if record.Result == "failed" || record.Result == "rejected" {
		return true
	}
	if record.KnownBugStatus == "future_plan" {
		return false
	}
	return record.SpecUpdateRef == "" || record.ResponsibilitySourceRef == "" || record.ImplementationRef == "" || record.ImplementationRef == "not_started" || record.RevalidationRef == ""
}

func phase17KnownBugOpen(record phase17EvidenceRecord) bool {
	switch record.KnownBugStatus {
	case "closed", "future_plan":
		return false
	default:
		return true
	}
}

func phase17CountMonitorSamples(records []phase17EvidenceRecord) int {
	count := 0
	seen := map[int]bool{}
	for _, record := range records {
		if record.OperationStep != "monitor" || record.MonitorSampleIndex == nil {
			continue
		}
		if !seen[*record.MonitorSampleIndex] {
			seen[*record.MonitorSampleIndex] = true
			count++
		}
	}
	return count
}

func phase17RecordWithResult(records []phase17EvidenceRecord, result string) bool {
	for _, record := range records {
		if record.Result == result {
			return true
		}
	}
	return false
}

func phase17AnyPendingExternalEvidence(records map[string][]phase17EvidenceRecord) bool {
	for _, recordSet := range records {
		for _, record := range recordSet {
			for _, value := range []string{
				record.ServiceState,
				record.StateDigest,
				record.LogDigest,
				record.RollbackState,
				record.CleanupResult,
				record.HealthResult,
				record.StateWriteResult,
				record.LogWriteResult,
				record.ImplementationRef,
			} {
				if value == "pending-external-evidence" || value == "not_started" {
					return true
				}
			}
		}
	}
	return false
}

func phase17OperationRecordPending(record phase17EvidenceRecord) bool {
	for _, value := range []string{
		record.ServiceState,
		record.HealthResult,
		record.StateWriteResult,
		record.LogWriteResult,
	} {
		if value == "pending-external-evidence" {
			return true
		}
	}
	return false
}

func phase17RequireWorkflow(t *testing.T, path string) {
	t.Helper()

	workflow := phase11MustReadText(t, path)
	if !strings.Contains(workflow, "name: phase17-production-validation") || !strings.Contains(workflow, "permissions:\n  contents: read") {
		t.Fatalf("%s must use the Phase 17 workflow name and contents:read permissions", path)
	}
	for _, requiredCheck := range phase17RequiredChecks {
		if !strings.Contains(workflow, "\n  "+requiredCheck+":") {
			t.Fatalf("%s must define required check job %s", path, requiredCheck)
		}
	}
	if strings.Count(workflow, "timeout-minutes:") < len(phase17RequiredChecks) {
		t.Fatalf("%s must set timeout-minutes on every Phase 17 job", path)
	}
	phase13RequirePinnedActions(t, workflow)
}

func phase17RequireDocumentState(t *testing.T) {
	t.Helper()

	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	if !strings.Contains(roadmap, "| Phase 17 | ConoHa VPS 試験本番運用。") || !strings.Contains(roadmap, "| 実装済み | Phase 16 |") {
		t.Fatalf("docs/ROADMAP.md must mark Phase 17 as 実装済み after closure evidence reaches zero")
	}
	if strings.Contains(roadmap, "現在の active Phase は Phase 17") {
		t.Fatalf("docs/ROADMAP.md must not keep Phase 17 as active after closure evidence reaches zero")
	}
	expectedPaths := []string{"testdata/phase17/production-validation/"}
	for _, rel := range phase17FixtureFiles {
		expectedPaths = append(expectedPaths, "testdata/phase17/production-validation/"+rel)
	}
	expectedPaths = append(expectedPaths, ".github/workflows/phase17-production-validation.yml")
	for _, rel := range expectedPaths {
		if !strings.Contains(documentIndex, "| `"+rel+"` |") && !strings.Contains(documentIndex, "| [`"+rel+"`]") {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index Phase 17 path %s", rel)
		}
	}
	for _, forbidden := range []string{
		"| `testdata/phase17/production-validation/` | Phase 17 正式 fixture root | 未作成 |",
		"| `.github/workflows/phase17-production-validation.yml` | Phase 17 required check workflow | 未作成 |",
	} {
		if strings.Contains(documentIndex, forbidden) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must not keep Phase 17 path as 未作成: %s", forbidden)
		}
	}
}

func phase17CounterAllowed(key string) bool {
	return phase17StringInSlice(key, phase17Counters)
}

func phase17RequiredCheckAllowed(check string) bool {
	return phase17StringInSlice(check, phase17RequiredChecks)
}

func phase17RequireNoSecretLikeFixtureContent(t *testing.T, path string, data []byte) {
	t.Helper()

	text := string(data)
	for _, forbidden := range []string{
		"-----BEGIN OPENSSH PRIVATE KEY-----",
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN PRIVATE KEY-----",
		"ghp_",
		"github_pat_",
		"xoxb-",
		"sk-proj-",
		"sk-svcacct-",
		"password=",
		"token=",
		"secret=",
		"http://",
		"https://",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("%s must not contain secret-like literal or endpoint-like value %q", path, forbidden)
		}
	}
	for _, token := range strings.FieldsFunc(text, func(r rune) bool {
		return !(r == '.' || r == '-' || r == '_' || r == ':' || r == '/' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z')
	}) {
		if phase17LooksLikeIPv4(token) {
			t.Fatalf("%s must not contain IPv4-like value %q", path, token)
		}
	}
}

func phase17LooksLikeIPv4(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		if part == "" || len(part) > 3 {
			return false
		}
		n := 0
		for _, r := range part {
			if r < '0' || r > '9' {
				return false
			}
			n = n*10 + int(r-'0')
		}
		if n > 255 {
			return false
		}
	}
	return true
}

func phase17IsLowerHexSHA256(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, r := range value {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}

func phase17StringInSlice(value string, values []string) bool {
	for _, candidate := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func phase17ContainsSecretLikeValue(values ...string) bool {
	for _, value := range values {
		lower := strings.ToLower(value)
		for _, forbidden := range []string{"token", "secret", "key", "password", "credential", "http://", "https://", "@"} {
			if strings.Contains(lower, forbidden) {
				return true
			}
		}
	}
	return false
}

const phase18Scope = "phase-18-vps-pull-bootstrap"

var phase18ProviderTargets = []string{"conoha-vps-primary", "xserver-vps-future"}

var phase18ValidationModes = []string{"vps-pull-bootstrap", "future_plan"}

var phase18RequiredChecks = []string{
	"phase18-bootstrap-prerequisite",
	"phase18-bootstrap-session",
	"phase18-bootstrap-artifact-health",
	"phase18-bootstrap-runtime-flow",
	"phase18-pull-runner-foundation",
	"phase18-bootstrap-update-rollback",
	"phase18-bugfix-spec-first",
	"phase18-security-boundary",
	"phase18-document-drift",
}

var phase18FixtureFiles = []string{
	"manifest.json",
	"input/environment_handover.json",
	"input/session_plan.json",
	"input/operation_scope.json",
	"input/security_boundary.json",
	"expected/counters.json",
	"expected/session_flow.json",
	"expected/bugfix_closure.json",
	"expected/negative_boundaries.json",
	"records/session.jsonl",
	"records/health.jsonl",
	"records/runtime.jsonl",
	"records/update_rollback.jsonl",
	"records/bugfix.jsonl",
	"records/security.jsonl",
	"records/document_drift.jsonl",
}

var phase18RecordFiles = []string{
	"records/session.jsonl",
	"records/health.jsonl",
	"records/runtime.jsonl",
	"records/update_rollback.jsonl",
	"records/bugfix.jsonl",
	"records/security.jsonl",
	"records/document_drift.jsonl",
}

var phase18Counters = []string{
	"phase18_bootstrap_prerequisite_open_count",
	"phase18_bootstrap_session_open_count",
	"phase18_bootstrap_runtime_open_count",
	"phase18_pull_runner_foundation_open_count",
	"phase18_bootstrap_update_rollback_open_count",
	"phase18_bugfix_spec_gap_count",
	"phase18_known_bug_open_count",
	"phase18_secret_leak_open_count",
	"phase18_destructive_operation_open_count",
	"phase18_document_drift_open_count",
	"final_open_item_count",
}

var phase18Phase17RequiredChecks = []string{
	"phase17-conoha-vps-preflight",
	"phase17-trial-production-buildout",
	"phase17-conoha-install-update-rollback",
	"phase17-systemd-lifecycle",
	"phase17-runtime-flow",
	"phase17-trial-operation-loop",
	"phase17-production-simulation",
	"phase17-failure-injection",
	"phase17-security-boundary",
	"phase17-bugfix-spec-first",
	"phase17-document-drift",
}

var phase18WorkUnits = []string{
	"phase17-handover",
	"bootstrap-session-open",
	"bootstrap-artifact-health",
	"bootstrap-runtime-flow",
	"pull-runner-foundation",
	"bootstrap-update-rollback-drill",
	"bootstrap-issue-triage",
	"spec-first-fix",
	"rebootstrap-and-revalidate",
	"bootstrap-session-close",
}

var phase18AbortConditions = []string{
	"secret-boundary-failure",
	"destructive-operation-boundary-failure",
	"state-corruption",
	"service-unrecoverable",
	"known-critical-bug",
	"document-drift-blocker",
}

var phase18OperationTargets = []string{
	"bootstrap-artifact-health",
	"bootstrap-runtime-flow",
	"pull-runner-foundation",
	"bootstrap-update-rollback-drill",
	"bootstrap-issue-triage",
	"spec-first-fix",
	"rebootstrap-and-revalidate",
}

var phase18SourceChannels = []string{"integration-head", "stable-release"}

var phase18TopologyRoles = []string{"ci-cd", "site", "single-node"}

var phase18BootstrapInvocationTokens = []string{
	"sh",
	"adlaire-ci-vps-pull-bootstrap.sh",
	"--source-channel",
	"integration-head|stable-release",
	"--topology-role",
	"ci-cd|site|single-node",
	"--source",
	"absolute-https-url-or-absolute-file-path",
	"--sha256",
	"64-lowercase-hex",
	"--install-dir",
	"/opt/adlaire-builder",
	"--bin-dir",
	"/usr/local/bin",
	"--state-dir",
	"/opt/adlaire-builder",
	"--service-user",
	"root",
}

var phase18BootstrapRequiredCommands = []string{"sh", "mktemp", "mkdir", "chmod", "install", "mv", "rm", "sha256sum", "systemctl"}

var phase18BootstrapHTTPClientOrder = []string{"curl", "wget"}

var phase18BootstrapFailureCodes = []string{
	"bootstrap-script-http-client-missing",
	"bootstrap-script-source-unreachable",
	"bootstrap-script-digest-mismatch",
	"bootstrap-script-version-mismatch",
	"bootstrap-script-executable-permission-invalid",
	"bootstrap-script-systemd-unit-invalid",
	"bootstrap-script-state-directory-invalid",
	"bootstrap-script-health-failed",
	"bootstrap-script-rollback-failed",
}

var phase18ForbiddenProviderOperations = []string{
	"vps-create",
	"vps-delete",
	"plan-change",
	"disk-rebuild",
	"volume-create",
	"volume-delete",
	"volume-attach",
	"volume-detach",
	"firewall-lockout",
	"ssh-lockout",
}

var phase18SecretClasses = []string{
	"provider-token",
	"admin-token",
	"ssh-private-key",
	"ip-specific-secret",
	"host-specific-secret",
}

var phase18ForbiddenEvidenceValues = []string{
	"provider-account-id",
	"provider-resource-id",
	"ip-address",
	"hostname",
	"fqdn",
	"credential-file-path",
	"secret-value",
	"secret-hash",
	"token-prefix-suffix",
}

var phase18OperationSamples = []string{"open", "runtime", "close"}

var phase18SampleObservations = []string{
	"service_state",
	"api_health",
	"admin_cli_health",
	"runner_state",
	"statefile_digest",
	"audit_log_write",
	"access_log_write",
	"config_log_write",
	"secret_leakage_boundary",
	"known_bug_status",
	"document_drift_status",
}

var phase18HealthTargets = []string{
	"bootstrap-script",
	"release-asset",
	"binary-version",
	"systemd-service",
	"api-health",
	"admin-cli",
	"sdk-client",
	"ui-runtime",
	"state-directory",
	"log-write",
}

var phase18RuntimeTargets = []string{
	"admin-api-sdk-ui",
	"runner-queue",
	"build-execution",
	"statefile-update",
	"audit-access-config-log",
	"commit-status",
	"deploy-simulation",
}

var phase18UpdateRollbackOperations = []string{"update", "rollback", "cleanup"}

var phase18IssueFlow = []string{
	"issue-detect",
	"spec-general-update",
	"implementation-fix",
	"redeploy",
	"revalidate",
	"evidence-record",
	"known-bug-zero-check",
}

var phase18TriageResults = []string{"fix_required", "environment_only", "not_required", "no_issue_detected"}

var phase18FixModes = []string{"spec-first-fix", "environment-only-revalidate", "not-required-revalidate"}

var phase18RequiredBugfixRefs = []string{
	"issue_ref",
	"spec_general_update_ref",
	"responsibility_source_ref",
	"implementation_fix_ref",
	"revalidation_ref",
	"evidence_ref",
}

var phase18DisallowedFixModes = []string{"hotfix-without-spec", "implementation-only-fix", "counter-only-closure"}

var phase18DocumentDriftTargets = []string{
	"state",
	"path",
	"anchor",
	"fixture-root",
	"workflow",
	"required-check",
	"checker-artifact",
	"closure-counter",
}

var phase18RequiredCheckRecordFiles = map[string][]string{
	"phase18-bootstrap-prerequisite":    {"records/session.jsonl"},
	"phase18-bootstrap-session":         {"records/session.jsonl"},
	"phase18-bootstrap-artifact-health": {"records/health.jsonl"},
	"phase18-bootstrap-runtime-flow":    {"records/health.jsonl", "records/runtime.jsonl"},
	"phase18-pull-runner-foundation":    {"records/runtime.jsonl"},
	"phase18-bootstrap-update-rollback": {"records/update_rollback.jsonl"},
	"phase18-bugfix-spec-first":         {"records/bugfix.jsonl"},
	"phase18-security-boundary":         {"records/security.jsonl"},
	"phase18-document-drift":            {"records/document_drift.jsonl"},
}

var phase18RequiredCheckCounterKeys = map[string][]string{
	"phase18-bootstrap-prerequisite":    {"phase18_bootstrap_prerequisite_open_count"},
	"phase18-bootstrap-session":         {"phase18_bootstrap_session_open_count"},
	"phase18-bootstrap-artifact-health": {"phase18_bootstrap_runtime_open_count"},
	"phase18-bootstrap-runtime-flow":    {"phase18_bootstrap_runtime_open_count"},
	"phase18-pull-runner-foundation":    {"phase18_pull_runner_foundation_open_count"},
	"phase18-bootstrap-update-rollback": {"phase18_bootstrap_update_rollback_open_count"},
	"phase18-bugfix-spec-first":         {"phase18_bugfix_spec_gap_count", "phase18_known_bug_open_count"},
	"phase18-security-boundary":         {"phase18_secret_leak_open_count", "phase18_destructive_operation_open_count"},
	"phase18-document-drift":            {"phase18_document_drift_open_count"},
}

type phase18NegativeBoundarySpec struct {
	Condition      string
	Group          string
	DiagnosticCode string
	CounterKey     string
}

var phase18NegativeBoundarySpecs = []phase18NegativeBoundarySpec{
	{"phase17-handover-missing-treated-as-success", "bootstrap-prerequisite", "record-link-missing", "phase18_bootstrap_prerequisite_open_count"},
	{"bootstrap-session-three-samples-missing-treated-as-success", "bootstrap-session", "session-mismatch", "phase18_bootstrap_session_open_count"},
	{"handover-open-sample-counted-as-session-open", "bootstrap-session", "session-mismatch", "phase18_bootstrap_session_open_count"},
	{"twenty-four-hour-monitoring-absence-treated-as-blocker", "bootstrap-session", "schema-mismatch", "phase18_bootstrap_session_open_count"},
	{"bootstrap-script-command-contract-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"bootstrap-script-contract-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"bootstrap-script-failure-code-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"bootstrap-script-output-contract-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"release-asset-or-version-evidence-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"runner-evidence-missing-counter-zero", "pull-runner-foundation", "counter-mismatch", "phase18_pull_runner_foundation_open_count"},
	{"update-rollback-missing-treated-as-success", "update-rollback", "record-link-missing", "phase18_bootstrap_update_rollback_open_count"},
	{"hotfix-without-spec-treated-as-success", "bugfix-spec-first", "record-link-missing", "phase18_bugfix_spec_gap_count"},
	{"known-bug-open-treated-as-complete", "known-bug", "schema-mismatch", "phase18_known_bug_open_count"},
	{"provider-identifier-written-to-evidence", "secret-boundary", "secret-boundary-failure", "phase18_secret_leak_open_count"},
	{"secret-written-to-output-or-record", "secret-boundary", "secret-boundary-failure", "phase18_secret_leak_open_count"},
	{"forbidden-provider-operation-treated-as-success", "destructive-operation", "destructive-operation-boundary-failure", "phase18_destructive_operation_open_count"},
	{"checker-artifact-drift-ignored", "document-drift", "document-drift", "phase18_document_drift_open_count"},
	{"closure-counter-drift-ignored", "document-drift", "counter-mismatch", "phase18_document_drift_open_count"},
	{"document-drift-ignored", "document-drift", "document-drift", "phase18_document_drift_open_count"},
	{"screenshot-only-treated-as-success", "document-drift", "record-link-missing", "phase18_document_drift_open_count"},
	{"xserver-future-treated-as-blocker", "document-drift", "document-drift", "phase18_document_drift_open_count"},
	{"xserver-future-used-as-required-check-success", "document-drift", "required-check-mismatch", "phase18_document_drift_open_count"},
}

type phase18CounterFile struct {
	SchemaVersion int            `json:"schema_version"`
	Scope         string         `json:"scope"`
	Counters      map[string]int `json:"counters"`
}

type phase18EnvironmentHandoverFile struct {
	SchemaVersion             int                    `json:"schema_version"`
	Scope                     string                 `json:"scope"`
	Phase17Ref                string                 `json:"phase17_ref"`
	ProviderTarget            string                 `json:"provider_target"`
	OS                        string                 `json:"os"`
	MinimumPlanClass          string                 `json:"minimum_plan_class"`
	MinimumRAMMB              int                    `json:"minimum_ram_mb"`
	EnvironmentIdentityPolicy string                 `json:"environment_identity_policy"`
	HandoverResults           phase18HandoverResults `json:"handover_results"`
}

type phase18HandoverResults struct {
	Phase17Status         string                         `json:"phase17_status"`
	Phase17RequiredChecks phase18Phase17RequiredCheckSet `json:"phase17_required_checks"`
	Phase17OpenItems      int                            `json:"phase17_open_items"`
	BootstrapPrerequisite string                         `json:"bootstrap_prerequisite"`
	Result                string                         `json:"result"`
}

type phase18Phase17RequiredCheckSet struct {
	PassedChecks  []string `json:"passed_checks"`
	FailedChecks  []string `json:"failed_checks"`
	MissingChecks []string `json:"missing_checks"`
}

type phase18SessionPlanFile struct {
	SchemaVersion                int      `json:"schema_version"`
	Scope                        string   `json:"scope"`
	SessionKind                  string   `json:"session_kind"`
	RequiredSessionCount         int      `json:"required_session_count"`
	RequiredOperationSampleCount int      `json:"required_operation_sample_count"`
	WorkUnits                    []string `json:"work_units"`
	AbortConditions              []string `json:"abort_conditions"`
	BugfixSpecFirstRequired      bool     `json:"bugfix_spec_first_required"`
	KnownBugZeroRequired         bool     `json:"known_bug_zero_required"`
}

type phase18OperationScopeFile struct {
	SchemaVersion               int      `json:"schema_version"`
	Scope                       string   `json:"scope"`
	ProviderTarget              string   `json:"provider_target"`
	ValidationMode              string   `json:"validation_mode"`
	BootstrapScriptArtifact     string   `json:"bootstrap_script_artifact"`
	BootstrapScriptRuntime      string   `json:"bootstrap_script_runtime"`
	BootstrapScriptInvocation   string   `json:"bootstrap_script_invocation"`
	SourceChannels              []string `json:"source_channels"`
	ActiveSourceChannel         string   `json:"active_source_channel"`
	TopologyRoles               []string `json:"topology_roles"`
	MinimumTopologyRole         string   `json:"minimum_topology_role"`
	BootstrapSourceIDFormat     string   `json:"bootstrap_source_identifier_format"`
	BootstrapDigestAlgorithm    string   `json:"bootstrap_digest_algorithm"`
	BootstrapInstallDirectory   string   `json:"bootstrap_install_directory"`
	BootstrapBinDirectory       string   `json:"bootstrap_bin_directory"`
	BootstrapStateDirectory     string   `json:"bootstrap_state_directory"`
	BootstrapServiceUser        string   `json:"bootstrap_service_user"`
	BootstrapSecretInputPolicy  string   `json:"bootstrap_secret_input_policy"`
	BootstrapStateModel         string   `json:"bootstrap_state_model"`
	BootstrapCLIContract        string   `json:"bootstrap_cli_contract"`
	BootstrapInvocationTokens   []string `json:"bootstrap_invocation_tokens"`
	BootstrapRequiredCommands   []string `json:"bootstrap_required_commands"`
	BootstrapHTTPClientOrder    []string `json:"bootstrap_http_client_order"`
	BootstrapStdoutContract     string   `json:"bootstrap_stdout_contract"`
	BootstrapStderrContract     string   `json:"bootstrap_stderr_contract"`
	BootstrapExitCodeContract   string   `json:"bootstrap_exit_code_contract"`
	BootstrapFailureCodes       []string `json:"bootstrap_failure_codes"`
	OperationTargets            []string `json:"operation_targets"`
	ForbiddenProviderOperations []string `json:"forbidden_provider_operations"`
	XserverFuturePolicy         string   `json:"xserver_future_policy"`
}

type phase18SecurityBoundaryFile struct {
	SchemaVersion                  int      `json:"schema_version"`
	Scope                          string   `json:"scope"`
	ProviderTarget                 string   `json:"provider_target"`
	ValidationMode                 string   `json:"validation_mode"`
	SecretReferencePolicy          string   `json:"secret_reference_policy"`
	MetadataPolicy                 string   `json:"metadata_policy"`
	SecretClasses                  []string `json:"secret_classes"`
	ForbiddenEvidenceValues        []string `json:"forbidden_evidence_values"`
	DestructiveOperationBoundaries []string `json:"destructive_operation_boundaries"`
	MaskRequired                   bool     `json:"mask_required"`
}

type phase18SessionFlowFile struct {
	SchemaVersion      int      `json:"schema_version"`
	Scope              string   `json:"scope"`
	SessionKind        string   `json:"session_kind"`
	OperationSamples   []string `json:"operation_samples"`
	SampleObservations []string `json:"sample_observations"`
	WorkUnits          []string `json:"work_units"`
	NoOpPolicy         string   `json:"no_op_policy"`
	ClosurePolicy      string   `json:"closure_policy"`
}

type phase18BugfixClosureFile struct {
	SchemaVersion          int      `json:"schema_version"`
	Scope                  string   `json:"scope"`
	IssueFlow              []string `json:"issue_flow"`
	TriageResults          []string `json:"triage_results"`
	FixModes               []string `json:"fix_modes"`
	RequiredRefs           []string `json:"required_refs"`
	TerminalKnownBugStatus string   `json:"terminal_known_bug_status"`
	DisallowedFixModes     []string `json:"disallowed_fix_modes"`
}

type phase18NegativeBoundariesFile struct {
	SchemaVersion int                       `json:"schema_version"`
	Scope         string                    `json:"scope"`
	Cases         []phase18NegativeBoundary `json:"cases"`
}

type phase18NegativeBoundary struct {
	CaseID                    string `json:"case_id"`
	BoundaryGroup             string `json:"boundary_group"`
	ForbiddenSuccessCondition string `json:"forbidden_success_condition"`
	ExpectedDiagnosticCode    string `json:"expected_diagnostic_code"`
	ExpectedCounterKey        string `json:"expected_counter_key"`
	SourceRef                 string `json:"source_ref"`
}

type phase18InputSet struct {
	Handover phase18EnvironmentHandoverFile
	Session  phase18SessionPlanFile
	Scope    phase18OperationScopeFile
	Security phase18SecurityBoundaryFile
	Flow     phase18SessionFlowFile
	Bugfix   phase18BugfixClosureFile
}

type phase18EvidenceRecord struct {
	RecordID                           string   `json:"record_id"`
	Scope                              string   `json:"scope"`
	EvidenceOrigin                     string   `json:"evidence_origin"`
	ExecutionEnvironment               string   `json:"execution_environment"`
	ProviderTarget                     string   `json:"provider_target"`
	ValidationMode                     string   `json:"validation_mode"`
	CheckName                          string   `json:"check_name"`
	SourceRef                          string   `json:"source_ref"`
	Result                             string   `json:"result"`
	CounterKey                         string   `json:"counter_key"`
	StartedAt                          *string  `json:"started_at"`
	EndedAt                            *string  `json:"ended_at"`
	StdoutSHA256                       *string  `json:"stdout_sha256"`
	StderrSHA256                       *string  `json:"stderr_sha256"`
	SecretScanResult                   string   `json:"secret_scan_result"`
	DestructiveOperationResult         string   `json:"destructive_operation_result"`
	EvidenceRef                        string   `json:"evidence_ref"`
	SessionID                          string   `json:"session_id,omitempty"`
	SessionKind                        string   `json:"session_kind,omitempty"`
	OperatorApprovalRef                string   `json:"operator_approval_ref,omitempty"`
	SampleIndex                        *int     `json:"sample_index,omitempty"`
	OperationSample                    string   `json:"operation_sample,omitempty"`
	WorkUnit                           string   `json:"work_unit,omitempty"`
	SampleObservations                 []string `json:"sample_observations,omitempty"`
	ServiceState                       string   `json:"service_state,omitempty"`
	APIHealth                          string   `json:"api_health,omitempty"`
	AdminCLIHealth                     string   `json:"admin_cli_health,omitempty"`
	RunnerState                        string   `json:"runner_state,omitempty"`
	StatefileDigest                    string   `json:"statefile_digest,omitempty"`
	AuditLogWrite                      string   `json:"audit_log_write,omitempty"`
	AccessLogWrite                     string   `json:"access_log_write,omitempty"`
	ConfigLogWrite                     string   `json:"config_log_write,omitempty"`
	SecretLeakageBoundary              string   `json:"secret_leakage_boundary,omitempty"`
	KnownBugStatus                     string   `json:"known_bug_status,omitempty"`
	DocumentDriftStatus                string   `json:"document_drift_status,omitempty"`
	AbortCondition                     string   `json:"abort_condition,omitempty"`
	CloseResult                        string   `json:"close_result,omitempty"`
	HealthTarget                       string   `json:"health_target,omitempty"`
	OwnerRef                           string   `json:"owner_ref,omitempty"`
	EndpointOrCommandRef               string   `json:"endpoint_or_command_ref,omitempty"`
	ExpectedState                      string   `json:"expected_state,omitempty"`
	ActualState                        string   `json:"actual_state,omitempty"`
	HealthResult                       string   `json:"health_result,omitempty"`
	LatencyMS                          *int     `json:"latency_ms,omitempty"`
	LogWriteResult                     string   `json:"log_write_result,omitempty"`
	StateReadResult                    string   `json:"state_read_result,omitempty"`
	SecretMaskResult                   string   `json:"secret_mask_result,omitempty"`
	Diagnostic                         string   `json:"diagnostic,omitempty"`
	ClosureAction                      string   `json:"closure_action,omitempty"`
	RuntimeTarget                      string   `json:"runtime_target,omitempty"`
	OperationRef                       string   `json:"operation_ref,omitempty"`
	QueueState                         string   `json:"queue_state,omitempty"`
	BuildResult                        string   `json:"build_result,omitempty"`
	StateUpdateResult                  string   `json:"state_update_result,omitempty"`
	AuditResult                        string   `json:"audit_result,omitempty"`
	CommitStatusResult                 string   `json:"commit_status_result,omitempty"`
	DeployResult                       string   `json:"deploy_result,omitempty"`
	NoCustomerDataResult               string   `json:"no_customer_data_result,omitempty"`
	OperationType                      string   `json:"operation_type,omitempty"`
	FromVersion                        string   `json:"from_version,omitempty"`
	ToVersion                          string   `json:"to_version,omitempty"`
	ServiceRestartResult               string   `json:"service_restart_result,omitempty"`
	StateMigrationResult               string   `json:"state_migration_result,omitempty"`
	RollbackState                      string   `json:"rollback_state,omitempty"`
	FailureCleanupResult               string   `json:"failure_cleanup_result,omitempty"`
	VersionAfterOperation              string   `json:"version_after_operation,omitempty"`
	IssueID                            string   `json:"issue_id,omitempty"`
	IssueClassification                string   `json:"issue_classification,omitempty"`
	TriageResult                       string   `json:"triage_result,omitempty"`
	FixMode                            string   `json:"fix_mode,omitempty"`
	IssueRef                           string   `json:"issue_ref,omitempty"`
	SpecGeneralUpdateRef               string   `json:"spec_general_update_ref,omitempty"`
	ResponsibilitySourceRef            string   `json:"responsibility_source_ref,omitempty"`
	ImplementationFixRef               string   `json:"implementation_fix_ref,omitempty"`
	RevalidationRef                    string   `json:"revalidation_ref,omitempty"`
	SecretClass                        string   `json:"secret_class,omitempty"`
	ScanTarget                         string   `json:"scan_target,omitempty"`
	MetadataPolicyResult               string   `json:"metadata_policy_result,omitempty"`
	ForbiddenValueDetected             *bool    `json:"forbidden_value_detected,omitempty"`
	ForbiddenProviderOperationToken    string   `json:"forbidden_provider_operation_token,omitempty"`
	ForbiddenProviderOperationDetected *bool    `json:"forbidden_provider_operation_detected,omitempty"`
	BoundaryResult                     string   `json:"boundary_result,omitempty"`
	MaskResult                         string   `json:"mask_result,omitempty"`
	StorageResult                      string   `json:"storage_result,omitempty"`
	DriftTarget                        string   `json:"drift_target,omitempty"`
	DocumentRef                        string   `json:"document_ref,omitempty"`
	AnchorRef                          string   `json:"anchor_ref,omitempty"`
	FixturePathRef                     string   `json:"fixture_path_ref,omitempty"`
	WorkflowPathRef                    string   `json:"workflow_path_ref,omitempty"`
	ExpectedRef                        string   `json:"expected_ref,omitempty"`
	ActualRef                          string   `json:"actual_ref,omitempty"`
	DriftResult                        string   `json:"drift_result,omitempty"`
	CorrectiveActionRef                string   `json:"corrective_action_ref,omitempty"`
	ClosureResult                      string   `json:"closure_result,omitempty"`
}

func phase18RequireFixtureFiles(t *testing.T, root string) {
	t.Helper()

	info, err := os.Lstat(root)
	if err != nil {
		t.Fatalf("Phase 18 fixture root must exist: %v", err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("Phase 18 fixture root must be a directory and not a symlink: %s", root)
	}
	allowed := map[string]bool{}
	for _, rel := range phase18FixtureFiles {
		allowed[rel] = true
		path := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Lstat(path)
		if err != nil {
			t.Fatalf("Phase 18 fixture path missing %s: %v", rel, err)
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			t.Fatalf("Phase 18 fixture path must be a regular file: %s", rel)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if len(data) == 0 || data[len(data)-1] != '\n' {
			t.Fatalf("%s must be non-empty and end with LF", path)
		}
		if strings.HasPrefix(string(data), "\ufeff") {
			t.Fatalf("%s must not include UTF-8 BOM", path)
		}
		if phase18ContainsForbiddenEvidenceLiteral(string(data)) {
			t.Fatalf("%s contains forbidden secret-like, endpoint-like, or provider-identity literal", path)
		}
		if strings.HasSuffix(rel, ".json") {
			phase16RequireNoDuplicateJSONKeys(t, path, data)
		}
		if strings.HasSuffix(rel, ".jsonl") {
			for index, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
				if strings.TrimSpace(line) == "" {
					t.Fatalf("%s line %d must not be empty", path, index+1)
				}
				phase16RequireNoDuplicateJSONKeys(t, fmt.Sprintf("%s line %d", path, index+1), []byte(line))
			}
		}
	}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root || entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !allowed[rel] {
			t.Fatalf("Phase 18 fixture contains unknown file %s", rel)
		}
		return nil
	}); err != nil {
		t.Fatalf("walk Phase 18 fixture: %v", err)
	}
}

func phase18ReadCounters(t *testing.T, path string) map[string]int {
	t.Helper()

	var file phase18CounterFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase18Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase18Scope)
	}
	return file.Counters
}

func phase18RequireCountersClosed(t *testing.T, counters map[string]int) {
	t.Helper()

	if len(counters) != len(phase18Counters) {
		t.Fatalf("Phase 18 counters must include exactly %d keys, got %d", len(phase18Counters), len(counters))
	}
	var sum int
	for _, key := range phase18Counters {
		value, ok := counters[key]
		if !ok {
			t.Fatalf("Phase 18 counters missing %s", key)
		}
		if value < 0 {
			t.Fatalf("Phase 18 counter %s must not be negative", key)
		}
		if key != "final_open_item_count" {
			sum += value
		}
	}
	if counters["final_open_item_count"] != sum {
		t.Fatalf("Phase 18 final_open_item_count must equal open counter sum: want %d got %d", sum, counters["final_open_item_count"])
	}
	if counters["final_open_item_count"] != 0 {
		t.Fatalf("Phase 18 checker acceptance fixture requires final_open_item_count=0, got %d", counters["final_open_item_count"])
	}
}

func phase18RequireInputsAndExpected(t *testing.T, root string) phase18InputSet {
	t.Helper()

	var inputs phase18InputSet
	phase16ReadStrictJSON(t, filepath.Join(root, "input/environment_handover.json"), &inputs.Handover)
	if inputs.Handover.SchemaVersion != 1 || inputs.Handover.Scope != phase18Scope || inputs.Handover.ProviderTarget != "conoha-vps-primary" || inputs.Handover.OS != "ubuntu-server-24.04-lts-64bit" || inputs.Handover.MinimumPlanClass != "conoha-vps-1gb-memory-class" || inputs.Handover.MinimumRAMMB != 1024 || inputs.Handover.EnvironmentIdentityPolicy != "opaque-non-secret-labels" {
		t.Fatalf("Phase 18 environment handover fixed values mismatch: %+v", inputs.Handover)
	}
	phase18RequireSafeReference(t, "input/environment_handover.json phase17_ref", inputs.Handover.Phase17Ref)
	if inputs.Handover.HandoverResults.Phase17Status != "implemented" || inputs.Handover.HandoverResults.Phase17OpenItems != 0 || inputs.Handover.HandoverResults.BootstrapPrerequisite != "passed" || inputs.Handover.HandoverResults.Result != "passed" {
		t.Fatalf("Phase 18 handover must close Phase 17 prerequisites: %+v", inputs.Handover.HandoverResults)
	}
	if !phase11StringSlicesEqual(inputs.Handover.HandoverResults.Phase17RequiredChecks.PassedChecks, phase18Phase17RequiredChecks) || len(inputs.Handover.HandoverResults.Phase17RequiredChecks.FailedChecks) != 0 || len(inputs.Handover.HandoverResults.Phase17RequiredChecks.MissingChecks) != 0 {
		t.Fatalf("Phase 18 handover Phase 17 required checks mismatch: %+v", inputs.Handover.HandoverResults.Phase17RequiredChecks)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "input/session_plan.json"), &inputs.Session)
	if inputs.Session.SchemaVersion != 1 || inputs.Session.Scope != phase18Scope || inputs.Session.SessionKind != "operator-approved-bootstrap-session" || inputs.Session.RequiredSessionCount != 1 || inputs.Session.RequiredOperationSampleCount != 3 || !phase11StringSlicesEqual(inputs.Session.WorkUnits, phase18WorkUnits) || !phase11StringSlicesEqual(inputs.Session.AbortConditions, phase18AbortConditions) || !inputs.Session.BugfixSpecFirstRequired || !inputs.Session.KnownBugZeroRequired {
		t.Fatalf("Phase 18 session plan mismatch: %+v", inputs.Session)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "input/operation_scope.json"), &inputs.Scope)
	if inputs.Scope.SchemaVersion != 1 || inputs.Scope.Scope != phase18Scope || inputs.Scope.ProviderTarget != "conoha-vps-primary" || inputs.Scope.ValidationMode != "vps-pull-bootstrap" || inputs.Scope.BootstrapScriptArtifact != "adlaire-ci-vps-pull-bootstrap.sh" || inputs.Scope.BootstrapScriptRuntime != "/bin/sh" || inputs.Scope.BootstrapScriptInvocation != "operator-runs-sh-script-on-vps" || !phase11StringSlicesEqual(inputs.Scope.SourceChannels, phase18SourceChannels) || inputs.Scope.ActiveSourceChannel != "integration-head" || !phase11StringSlicesEqual(inputs.Scope.TopologyRoles, phase18TopologyRoles) || inputs.Scope.MinimumTopologyRole != "single-node" || inputs.Scope.BootstrapSourceIDFormat != "absolute-https-url-or-absolute-file-path" || inputs.Scope.BootstrapDigestAlgorithm != "sha256" || inputs.Scope.BootstrapInstallDirectory != "/opt/adlaire-builder" || inputs.Scope.BootstrapBinDirectory != "/usr/local/bin" || inputs.Scope.BootstrapStateDirectory != "/opt/adlaire-builder" || inputs.Scope.BootstrapServiceUser != "root" || inputs.Scope.BootstrapSecretInputPolicy != "no-command-argument-secret" || inputs.Scope.BootstrapStateModel != "staging-verify-commit-rollback" || inputs.Scope.BootstrapCLIContract != "fixed-token-flag-interface" || !phase11StringSlicesEqual(inputs.Scope.BootstrapInvocationTokens, phase18BootstrapInvocationTokens) || !phase11StringSlicesEqual(inputs.Scope.BootstrapRequiredCommands, phase18BootstrapRequiredCommands) || !phase11StringSlicesEqual(inputs.Scope.BootstrapHTTPClientOrder, phase18BootstrapHTTPClientOrder) || inputs.Scope.BootstrapStdoutContract != "single-json-object-secret-safe" || inputs.Scope.BootstrapStderrContract != "diagnostic-json-lines-on-failure" || inputs.Scope.BootstrapExitCodeContract != "common-cli-exit-code-contract" || !phase11StringSlicesEqual(inputs.Scope.BootstrapFailureCodes, phase18BootstrapFailureCodes) || !phase11StringSlicesEqual(inputs.Scope.OperationTargets, phase18OperationTargets) || !phase11StringSlicesEqual(inputs.Scope.ForbiddenProviderOperations, phase18ForbiddenProviderOperations) || inputs.Scope.XserverFuturePolicy != "future_plan_not_blocker" {
		t.Fatalf("Phase 18 operation scope mismatch: %+v", inputs.Scope)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "input/security_boundary.json"), &inputs.Security)
	if inputs.Security.SchemaVersion != 1 || inputs.Security.Scope != phase18Scope || inputs.Security.ProviderTarget != "conoha-vps-primary" || inputs.Security.ValidationMode != "vps-pull-bootstrap" || inputs.Security.SecretReferencePolicy != "metadata-only" || inputs.Security.MetadataPolicy != "opaque-non-secret-labels" || !phase11StringSlicesEqual(inputs.Security.SecretClasses, phase18SecretClasses) || !phase11StringSlicesEqual(inputs.Security.ForbiddenEvidenceValues, phase18ForbiddenEvidenceValues) || !phase11StringSlicesEqual(inputs.Security.DestructiveOperationBoundaries, phase18ForbiddenProviderOperations) || !inputs.Security.MaskRequired {
		t.Fatalf("Phase 18 security boundary mismatch: %+v", inputs.Security)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "expected/session_flow.json"), &inputs.Flow)
	if inputs.Flow.SchemaVersion != 1 || inputs.Flow.Scope != phase18Scope || inputs.Flow.SessionKind != inputs.Session.SessionKind || !phase11StringSlicesEqual(inputs.Flow.OperationSamples, phase18OperationSamples) || !phase11StringSlicesEqual(inputs.Flow.SampleObservations, phase18SampleObservations) || !phase11StringSlicesEqual(inputs.Flow.WorkUnits, inputs.Session.WorkUnits) || inputs.Flow.NoOpPolicy != "explicit-not-required-record-required" || inputs.Flow.ClosurePolicy != "single-session-id-required" {
		t.Fatalf("Phase 18 session flow mismatch: %+v", inputs.Flow)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "expected/bugfix_closure.json"), &inputs.Bugfix)
	if inputs.Bugfix.SchemaVersion != 1 || inputs.Bugfix.Scope != phase18Scope || !phase11StringSlicesEqual(inputs.Bugfix.IssueFlow, phase18IssueFlow) || !phase11StringSlicesEqual(inputs.Bugfix.TriageResults, phase18TriageResults) || !phase11StringSlicesEqual(inputs.Bugfix.FixModes, phase18FixModes) || !phase11StringSlicesEqual(inputs.Bugfix.RequiredRefs, phase18RequiredBugfixRefs) || inputs.Bugfix.TerminalKnownBugStatus != "none" || !phase11StringSlicesEqual(inputs.Bugfix.DisallowedFixModes, phase18DisallowedFixModes) {
		t.Fatalf("Phase 18 bugfix closure mismatch: %+v", inputs.Bugfix)
	}

	var negative phase18NegativeBoundariesFile
	phase16ReadStrictJSON(t, filepath.Join(root, "expected/negative_boundaries.json"), &negative)
	if negative.SchemaVersion != 1 || negative.Scope != phase18Scope || len(negative.Cases) != len(phase18NegativeBoundarySpecs) {
		t.Fatalf("Phase 18 negative boundary file mismatch: %+v", negative)
	}
	seenCondition := map[string]bool{}
	for index, want := range phase18NegativeBoundarySpecs {
		got := negative.Cases[index]
		wantID := "phase18.negative." + want.Group + "." + want.Condition
		if got.CaseID != wantID || got.BoundaryGroup != want.Group || got.ForbiddenSuccessCondition != want.Condition || got.ExpectedDiagnosticCode != want.DiagnosticCode || got.ExpectedCounterKey != want.CounterKey {
			t.Fatalf("Phase 18 negative boundary %d mismatch\nwant: %+v / %s\n got: %+v", index+1, want, wantID, got)
		}
		if seenCondition[got.ForbiddenSuccessCondition] {
			t.Fatalf("Phase 18 negative boundary repeats forbidden_success_condition %s", got.ForbiddenSuccessCondition)
		}
		seenCondition[got.ForbiddenSuccessCondition] = true
		phase18RequireSafeReference(t, "negative boundary source_ref", got.SourceRef)
	}
	return inputs
}

func phase18ReadRecords(t *testing.T, root string) map[string][]phase18EvidenceRecord {
	t.Helper()

	records := map[string][]phase18EvidenceRecord{}
	for _, rel := range phase18RecordFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
		if len(lines) == 0 {
			t.Fatalf("%s must include at least one JSONL record", path)
		}
		seen := map[string]bool{}
		for index, line := range lines {
			if strings.TrimSpace(line) == "" {
				t.Fatalf("%s line %d must not be empty", path, index+1)
			}
			phase16RequireNoDuplicateJSONKeys(t, fmt.Sprintf("%s line %d", path, index+1), []byte(line))
			var record phase18EvidenceRecord
			decoder := json.NewDecoder(strings.NewReader(line))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&record); err != nil {
				t.Fatalf("decode %s line %d: %v", path, index+1, err)
			}
			var trailing any
			if err := decoder.Decode(&trailing); err != io.EOF {
				t.Fatalf("%s line %d must contain exactly one JSON object", path, index+1)
			}
			phase18RequireRecordCommon(t, rel, index+1, record)
			if seen[record.RecordID] {
				t.Fatalf("%s repeats record_id %s", path, record.RecordID)
			}
			seen[record.RecordID] = true
			phase18RequireRecordSpecific(t, rel, index+1, record)
			records[rel] = append(records[rel], record)
		}
	}
	return records
}

func phase18RequireRecordCommon(t *testing.T, rel string, line int, record phase18EvidenceRecord) {
	t.Helper()

	prefix := "phase18.record." + strings.TrimSuffix(filepath.Base(rel), ".jsonl") + "."
	if record.RecordID == "" || !strings.HasPrefix(record.RecordID, prefix) || record.Scope != phase18Scope || record.EvidenceOrigin != "checker-acceptance-fixture" || record.ExecutionEnvironment != "repository-fixture" || record.ProviderTarget == "" || record.ValidationMode == "" || record.CheckName == "" || record.SourceRef == "" || record.Result == "" || record.CounterKey == "" || record.SecretScanResult == "" || record.DestructiveOperationResult == "" || record.EvidenceRef == "" {
		t.Fatalf("%s line %d common fields invalid: %+v", rel, line, record)
	}
	if !phase17StringInSlice(record.ProviderTarget, phase18ProviderTargets) {
		t.Fatalf("%s line %d unknown provider_target %s", rel, line, record.ProviderTarget)
	}
	if !phase17StringInSlice(record.ValidationMode, phase18ValidationModes) {
		t.Fatalf("%s line %d unknown validation_mode %s", rel, line, record.ValidationMode)
	}
	if !phase17StringInSlice(record.CheckName, phase18RequiredChecks) {
		t.Fatalf("%s line %d unknown check_name %s", rel, line, record.CheckName)
	}
	if !phase17StringInSlice(record.CounterKey, phase18Counters) {
		t.Fatalf("%s line %d unknown counter_key %s", rel, line, record.CounterKey)
	}
	switch record.Result {
	case "passed", "not_applicable", "future_plan":
	case "failed", "rejected":
		t.Fatalf("%s line %d must not leave failed or rejected result", rel, line)
	default:
		t.Fatalf("%s line %d invalid result %s", rel, line, record.Result)
	}
	if record.Result == "passed" {
		if record.StartedAt == nil || record.EndedAt == nil || *record.StartedAt == "" || *record.EndedAt == "" {
			t.Fatalf("%s line %d passed record must include started_at and ended_at", rel, line)
		}
		if _, err := time.Parse(time.RFC3339, *record.StartedAt); err != nil {
			t.Fatalf("%s line %d invalid started_at: %v", rel, line, err)
		}
		if _, err := time.Parse(time.RFC3339, *record.EndedAt); err != nil {
			t.Fatalf("%s line %d invalid ended_at: %v", rel, line, err)
		}
	}
	for label, value := range map[string]*string{"stdout_sha256": record.StdoutSHA256, "stderr_sha256": record.StderrSHA256} {
		if value != nil && !phase17IsLowerHexSHA256(*value) {
			t.Fatalf("%s line %d %s must be lowercase SHA-256 hex or null", rel, line, label)
		}
	}
	switch record.SecretScanResult {
	case "passed", "not_applicable":
	default:
		t.Fatalf("%s line %d invalid secret_scan_result %s", rel, line, record.SecretScanResult)
	}
	switch record.DestructiveOperationResult {
	case "not_applicable", "not_executed", "blocked":
	default:
		t.Fatalf("%s line %d invalid destructive_operation_result %s", rel, line, record.DestructiveOperationResult)
	}
	phase18RequireSafeReference(t, fmt.Sprintf("%s line %d source_ref", rel, line), record.SourceRef)
	phase18RequireSafeReference(t, fmt.Sprintf("%s line %d evidence_ref", rel, line), record.EvidenceRef)
	phase18RequireRecordReferences(t, rel, line, record)
	if record.ProviderTarget == "xserver-vps-future" && (rel != "records/document_drift.jsonl" || record.ValidationMode != "future_plan" || record.Result != "future_plan" || record.CheckName != "phase18-document-drift") {
		t.Fatalf("%s line %d xserver-vps-future is allowed only as a document_drift future plan classification record: %+v", rel, line, record)
	}
}

func phase18RequireRecordSpecific(t *testing.T, rel string, line int, record phase18EvidenceRecord) {
	t.Helper()

	switch rel {
	case "records/session.jsonl":
		if record.SessionID == "" || record.SessionKind != "operator-approved-bootstrap-session" || record.OperatorApprovalRef == "" || record.SampleIndex == nil || record.OperationSample == "" || record.WorkUnit == "" || len(record.SampleObservations) == 0 || record.ServiceState == "" || record.APIHealth == "" || record.AdminCLIHealth == "" || record.RunnerState == "" || record.StatefileDigest == "" || record.AuditLogWrite == "" || record.AccessLogWrite == "" || record.ConfigLogWrite == "" || record.SecretLeakageBoundary == "" || record.KnownBugStatus == "" || record.DocumentDriftStatus == "" || record.AbortCondition == "" || record.CloseResult == "" {
			t.Fatalf("%s line %d missing session fields: %+v", rel, line, record)
		}
		if !phase17IsLowerHexSHA256(record.StatefileDigest) || !phase11StringSlicesEqual(record.SampleObservations, phase18SampleObservations) || !phase17StringInSlice(record.OperationSample, phase18OperationSamples) || (!phase17StringInSlice(record.AbortCondition, phase18AbortConditions) && record.AbortCondition != "not_triggered") {
			t.Fatalf("%s line %d session fixed fields mismatch: %+v", rel, line, record)
		}
	case "records/health.jsonl":
		if record.SessionID == "" || record.HealthTarget == "" || record.OwnerRef == "" || record.EndpointOrCommandRef == "" || record.ExpectedState == "" || record.ActualState == "" || record.HealthResult == "" || record.LogWriteResult == "" || record.StateReadResult == "" || record.SecretMaskResult == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing health fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.HealthTarget, phase18HealthTargets) || record.CounterKey != "phase18_bootstrap_runtime_open_count" || record.HealthResult != "passed" || record.StateReadResult != "passed" || record.SecretMaskResult != "passed" {
			t.Fatalf("%s line %d health result mismatch: %+v", rel, line, record)
		}
		if record.HealthTarget == "log-write" && record.LogWriteResult != "passed" {
			t.Fatalf("%s line %d log-write target must close log_write_result: %+v", rel, line, record)
		}
		if phase17StringInSlice(record.HealthTarget, []string{"api-health", "admin-cli", "sdk-client", "ui-runtime"}) && record.LatencyMS == nil {
			t.Fatalf("%s line %d interactive health target must include latency_ms", rel, line)
		}
	case "records/runtime.jsonl":
		if record.SessionID == "" || record.RuntimeTarget == "" || record.OwnerRef == "" || record.OperationRef == "" || record.QueueState == "" || record.BuildResult == "" || record.StateUpdateResult == "" || record.AuditResult == "" || record.CommitStatusResult == "" || record.DeployResult == "" || record.NoCustomerDataResult == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing runtime fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.RuntimeTarget, phase18RuntimeTargets) || record.NoCustomerDataResult != "passed" {
			t.Fatalf("%s line %d runtime target mismatch: %+v", rel, line, record)
		}
	case "records/update_rollback.jsonl":
		if record.SessionID == "" || record.OperationType == "" || record.FromVersion == "" || record.ToVersion == "" || record.ServiceRestartResult == "" || record.StateMigrationResult == "" || record.RollbackState == "" || record.FailureCleanupResult == "" || record.VersionAfterOperation == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing update/rollback fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.OperationType, phase18UpdateRollbackOperations) || record.CounterKey != "phase18_bootstrap_update_rollback_open_count" {
			t.Fatalf("%s line %d update/rollback operation mismatch: %+v", rel, line, record)
		}
	case "records/bugfix.jsonl":
		if record.SessionID == "" || record.IssueID == "" || record.IssueClassification == "" || record.TriageResult == "" || record.FixMode == "" || record.IssueRef == "" || record.SpecGeneralUpdateRef == "" || record.ResponsibilitySourceRef == "" || record.ImplementationFixRef == "" || record.RevalidationRef == "" || record.KnownBugStatus == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing bugfix fields: %+v", rel, line, record)
		}
	case "records/security.jsonl":
		if record.SessionID == "" || record.SecretClass == "" || record.ScanTarget == "" || record.MetadataPolicyResult == "" || record.ForbiddenValueDetected == nil || record.ForbiddenProviderOperationToken == "" || record.ForbiddenProviderOperationDetected == nil || record.BoundaryResult == "" || record.MaskResult == "" || record.StorageResult == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing security fields: %+v", rel, line, record)
		}
	case "records/document_drift.jsonl":
		if record.SessionID == "" || record.DriftTarget == "" || record.DocumentRef == "" || record.AnchorRef == "" || record.FixturePathRef == "" || record.WorkflowPathRef == "" || record.ExpectedRef == "" || record.ActualRef == "" || record.DriftResult == "" || record.CorrectiveActionRef == "" || record.ClosureResult == "" {
			t.Fatalf("%s line %d missing document drift fields: %+v", rel, line, record)
		}
	}
}

func phase18RequireRequiredCheckConnections(t *testing.T, records map[string][]phase18EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for rel, recordSet := range records {
		for _, record := range recordSet {
			files := phase18RequiredCheckRecordFiles[record.CheckName]
			if !phase17StringInSlice(rel, files) {
				t.Fatalf("Phase 18 record %s places check %s in invalid record file %s", record.RecordID, record.CheckName, rel)
			}
			if !phase17StringInSlice(record.CounterKey, phase18RequiredCheckCounterKeys[record.CheckName]) {
				t.Fatalf("Phase 18 record %s connects check %s to invalid counter %s", record.RecordID, record.CheckName, record.CounterKey)
			}
			if record.ProviderTarget != "xserver-vps-future" {
				seen[record.CheckName] = true
			}
		}
	}
	for _, check := range phase18RequiredChecks {
		if !seen[check] {
			t.Fatalf("Phase 18 required check %s is not connected to evidence records", check)
		}
	}
}

func phase18RequireSessionRecords(t *testing.T, records []phase18EvidenceRecord, inputs phase18InputSet) {
	t.Helper()

	var sessionID string
	operationSamples := map[string]phase18EvidenceRecord{}
	var handover bool
	for _, record := range records {
		if record.ProviderTarget != "conoha-vps-primary" {
			continue
		}
		if !strings.HasPrefix(record.SessionID, "phase18.session.") {
			t.Fatalf("Phase 18 session_id must use phase18.session.<slug>: %s", record.SessionID)
		}
		if sessionID == "" {
			sessionID = record.SessionID
		}
		if record.SessionID != sessionID {
			t.Fatalf("Phase 18 bootstrap records must reuse a single session_id: %s vs %s", record.SessionID, sessionID)
		}
		if record.CheckName == "phase18-bootstrap-prerequisite" {
			if record.WorkUnit != "phase17-handover" || record.OperationSample != "open" || record.CounterKey != "phase18_bootstrap_prerequisite_open_count" {
				t.Fatalf("Phase 18 handover session record mismatch: %+v", record)
			}
			handover = true
			continue
		}
		if record.CheckName == "phase18-bootstrap-session" {
			operationSamples[record.OperationSample] = record
		}
	}
	if !handover {
		t.Fatalf("Phase 18 session records must include environment handover record")
	}
	if len(operationSamples) != inputs.Session.RequiredOperationSampleCount {
		t.Fatalf("Phase 18 bootstrap session must include exactly %d samples, got %d", inputs.Session.RequiredOperationSampleCount, len(operationSamples))
	}
	for index, sample := range phase18OperationSamples {
		record, ok := operationSamples[sample]
		if !ok {
			t.Fatalf("Phase 18 bootstrap session missing sample %s", sample)
		}
		if record.SampleIndex == nil || *record.SampleIndex != index+1 {
			t.Fatalf("Phase 18 bootstrap session sample %s index mismatch: %+v", sample, record)
		}
	}
	closeRecord := operationSamples["close"]
	if closeRecord.CloseResult != "closed" || closeRecord.SecretLeakageBoundary != "passed" || closeRecord.KnownBugStatus != "none" || closeRecord.DocumentDriftStatus != "no_drift" || closeRecord.AbortCondition != "not_triggered" {
		t.Fatalf("Phase 18 close sample must close session without leaks, known bugs, drift, or abort: %+v", closeRecord)
	}
}

func phase18RequireHealthRecords(t *testing.T, records []phase18EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for _, record := range records {
		if record.Result != "passed" || record.HealthResult != "passed" || record.StateReadResult != "passed" || record.SecretMaskResult != "passed" {
			t.Fatalf("Phase 18 health target must pass: %+v", record)
		}
		if record.HealthTarget != "log-write" && record.LogWriteResult != "not_applicable" {
			t.Fatalf("Phase 18 non-log health target must use log_write_result=not_applicable: %+v", record)
		}
		if record.HealthTarget == "log-write" && record.LogWriteResult != "passed" {
			t.Fatalf("Phase 18 log-write health target must pass log write: %+v", record)
		}
		if record.HealthTarget == "bootstrap-script" {
			expectedState := "adlaire-ci-vps-pull-bootstrap.sh /bin/sh integration-head single-node absolute-https-url-or-absolute-file-path sha256 /opt/adlaire-builder /usr/local/bin /opt/adlaire-builder root no-command-argument-secret staging-verify-commit-rollback fixed-token-flag-interface sh adlaire-ci-vps-pull-bootstrap.sh --source-channel integration-head|stable-release --topology-role ci-cd|site|single-node --source absolute-https-url-or-absolute-file-path --sha256 64-lowercase-hex --install-dir /opt/adlaire-builder --bin-dir /usr/local/bin --state-dir /opt/adlaire-builder --service-user root sh mktemp mkdir chmod install mv rm sha256sum systemctl curl wget single-json-object-secret-safe diagnostic-json-lines-on-failure common-cli-exit-code-contract bootstrap-script-http-client-missing bootstrap-script-source-unreachable bootstrap-script-digest-mismatch bootstrap-script-version-mismatch bootstrap-script-executable-permission-invalid bootstrap-script-systemd-unit-invalid bootstrap-script-state-directory-invalid bootstrap-script-health-failed bootstrap-script-rollback-failed"
			if record.SourceRef != "input/operation_scope.json" || record.OwnerRef != "docs/details/production-validation.md" || record.EndpointOrCommandRef != "input/operation_scope.json" || record.ExpectedState != expectedState || record.ActualState != expectedState || !strings.Contains(record.Diagnostic, "source identifier") || !strings.Contains(record.Diagnostic, "secret policy") || !strings.Contains(record.Diagnostic, "cli token") || !strings.Contains(record.Diagnostic, "required command") || !strings.Contains(record.Diagnostic, "http client") || !strings.Contains(record.Diagnostic, "exit code") || !strings.Contains(record.Diagnostic, "failure code") {
				t.Fatalf("Phase 18 bootstrap script health target must prove exact bootstrap contract: %+v", record)
			}
		}
		seen[record.HealthTarget] = true
	}
	for _, target := range phase18HealthTargets {
		if !seen[target] {
			t.Fatalf("Phase 18 health records missing target %s", target)
		}
	}
}

func phase18RequireRuntimeRecords(t *testing.T, records []phase18EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for _, record := range records {
		seen[record.RuntimeTarget] = true
		if record.RuntimeTarget == "admin-api-sdk-ui" {
			if record.CheckName != "phase18-bootstrap-runtime-flow" || record.CounterKey != "phase18_bootstrap_runtime_open_count" {
				t.Fatalf("Phase 18 admin/api/sdk/ui runtime target must connect to runtime-flow: %+v", record)
			}
		} else if record.CheckName != "phase18-pull-runner-foundation" || record.CounterKey != "phase18_pull_runner_foundation_open_count" {
			t.Fatalf("Phase 18 runner runtime target must connect to pull-runner-foundation: %+v", record)
		}
		if record.RuntimeTarget == "runner-queue" && record.QueueState != "succeeded" {
			t.Fatalf("Phase 18 runner queue must reach succeeded: %+v", record)
		}
		if record.RuntimeTarget == "build-execution" && record.BuildResult != "passed" {
			t.Fatalf("Phase 18 build execution must pass: %+v", record)
		}
		if record.RuntimeTarget == "statefile-update" && record.StateUpdateResult != "passed" {
			t.Fatalf("Phase 18 statefile update must pass: %+v", record)
		}
		if record.RuntimeTarget == "audit-access-config-log" && record.AuditResult != "passed" {
			t.Fatalf("Phase 18 audit/access/config log runtime must pass: %+v", record)
		}
		if record.RuntimeTarget == "commit-status" && record.CommitStatusResult != "passed" && record.CommitStatusResult != "not_applicable" {
			t.Fatalf("Phase 18 commit status result invalid: %+v", record)
		}
		if record.RuntimeTarget == "deploy-simulation" && record.DeployResult != "passed" {
			t.Fatalf("Phase 18 deploy simulation must pass: %+v", record)
		}
		if record.NoCustomerDataResult != "passed" {
			t.Fatalf("Phase 18 runtime records must prove no customer data use: %+v", record)
		}
	}
	for _, target := range phase18RuntimeTargets {
		if !seen[target] {
			t.Fatalf("Phase 18 runtime records missing target %s", target)
		}
	}
}

func phase18RequireUpdateRollbackRecords(t *testing.T, records []phase18EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for _, record := range records {
		seen[record.OperationType] = true
		if record.Result != "passed" || record.FromVersion == "unknown" || record.ToVersion == "unknown" || record.VersionAfterOperation == "unknown" {
			t.Fatalf("Phase 18 update/rollback records must pass with concrete versions: %+v", record)
		}
		switch record.OperationType {
		case "update":
			if record.FromVersion == record.ToVersion || record.VersionAfterOperation != record.ToVersion || record.ServiceRestartResult != "passed" || record.StateMigrationResult != "passed" || record.RollbackState != "committed" {
				t.Fatalf("Phase 18 update record mismatch: %+v", record)
			}
		case "rollback":
			if record.VersionAfterOperation != record.ToVersion || record.ServiceRestartResult != "passed" || record.StateMigrationResult != "passed" || record.RollbackState != "rolled_back" {
				t.Fatalf("Phase 18 rollback record mismatch: %+v", record)
			}
		case "cleanup":
			if record.FailureCleanupResult != "passed" || record.RollbackState != "not_applicable" {
				t.Fatalf("Phase 18 cleanup record mismatch: %+v", record)
			}
		default:
			t.Fatalf("Phase 18 unknown update/rollback operation %s", record.OperationType)
		}
	}
	for _, operation := range phase18UpdateRollbackOperations {
		if !seen[operation] {
			t.Fatalf("Phase 18 update/rollback records missing operation %s", operation)
		}
	}
}

func phase18RequireBugfixRecords(t *testing.T, records []phase18EvidenceRecord, expected phase18BugfixClosureFile) {
	t.Helper()

	var specFirst bool
	var noIssue bool
	for _, record := range records {
		if record.CounterKey != "phase18_bugfix_spec_gap_count" && record.CounterKey != "phase18_known_bug_open_count" {
			t.Fatalf("Phase 18 bugfix record uses invalid counter: %+v", record)
		}
		if !phase17StringInSlice(record.TriageResult, expected.TriageResults) || !phase17StringInSlice(record.FixMode, expected.FixModes) {
			t.Fatalf("Phase 18 bugfix record uses unknown triage or fix mode: %+v", record)
		}
		if record.KnownBugStatus != expected.TerminalKnownBugStatus {
			t.Fatalf("Phase 18 bugfix records must close known bugs to %s: %+v", expected.TerminalKnownBugStatus, record)
		}
		if record.FixMode == "spec-first-fix" {
			specFirst = true
			for label, ref := range map[string]string{
				"spec_general_update_ref":   record.SpecGeneralUpdateRef,
				"responsibility_source_ref": record.ResponsibilitySourceRef,
				"implementation_fix_ref":    record.ImplementationFixRef,
				"revalidation_ref":          record.RevalidationRef,
				"evidence_ref":              record.EvidenceRef,
			} {
				phase18RequireSafeReference(t, "bugfix "+label, ref)
			}
		}
		if record.IssueClassification == "no_issue" {
			noIssue = true
			if record.IssueID != "not_applicable" || record.TriageResult != "no_issue_detected" || record.FixMode != "not-required-revalidate" || record.ImplementationFixRef != "not_applicable" {
				t.Fatalf("Phase 18 no_issue bugfix record mismatch: %+v", record)
			}
		}
	}
	if !specFirst || !noIssue {
		t.Fatalf("Phase 18 bugfix records must include spec-first-fix and no_issue closure records")
	}
}

func phase18RequireSecurityRecords(t *testing.T, records []phase18EvidenceRecord, inputs phase18InputSet) {
	t.Helper()

	secretCoverage := map[string]bool{}
	operationCoverage := map[string]bool{}
	for _, record := range records {
		if record.SecretClass != "not_applicable" {
			if !phase17StringInSlice(record.SecretClass, inputs.Security.SecretClasses) {
				t.Fatalf("Phase 18 security record unknown secret_class: %+v", record)
			}
			if *record.ForbiddenValueDetected || record.MetadataPolicyResult != "opaque-only" || record.StorageResult != "not_stored" || (record.MaskResult != "masked" && record.MaskResult != "not_present") || record.Result != "passed" {
				t.Fatalf("Phase 18 secret coverage record must prove opaque non-storage and masking: %+v", record)
			}
			secretCoverage[record.SecretClass] = true
		}
		if record.ForbiddenProviderOperationToken != "not_applicable" {
			if !phase17StringInSlice(record.ForbiddenProviderOperationToken, inputs.Scope.ForbiddenProviderOperations) {
				t.Fatalf("Phase 18 security record unknown forbidden operation: %+v", record)
			}
			if *record.ForbiddenProviderOperationDetected || record.BoundaryResult != "blocked" || record.DestructiveOperationResult != "blocked" || record.Result != "passed" {
				t.Fatalf("Phase 18 destructive operation must be blocked, not executed: %+v", record)
			}
			operationCoverage[record.ForbiddenProviderOperationToken] = true
		}
	}
	for _, secretClass := range inputs.Security.SecretClasses {
		if !secretCoverage[secretClass] {
			t.Fatalf("Phase 18 security records missing secret class coverage %s", secretClass)
		}
	}
	for _, operation := range inputs.Scope.ForbiddenProviderOperations {
		if !operationCoverage[operation] {
			t.Fatalf("Phase 18 security records missing forbidden provider operation coverage %s", operation)
		}
	}
}

func phase18RequireDocumentDriftRecords(t *testing.T, records []phase18EvidenceRecord) {
	t.Helper()

	targetCoverage := map[string]bool{}
	var futurePlan bool
	futurePlanSourceRef := "[docs/details/production-validation.md 本番検証詳細本文責務 Phase 18 VPS Pull Bootstrap 契約](docs/details/production-validation.md#phase-18-vps-pull-bootstrap-contract)"
	futurePlanEvidenceRef := "[docs/ROADMAP.md 状態・計画責務 §4.1](docs/ROADMAP.md#roadmap-initial-phase-plan)"
	for _, record := range records {
		if record.ProviderTarget == "xserver-vps-future" {
			if record.SessionID != "not_applicable" || record.ValidationMode != "future_plan" || record.Result != "future_plan" || record.DriftResult != "future_plan_confirmed" || record.ClosureResult != "not_applicable" || record.CounterKey != "phase18_document_drift_open_count" || record.SourceRef != futurePlanSourceRef || record.EvidenceRef != futurePlanEvidenceRef || record.DocumentRef != futurePlanSourceRef || record.AnchorRef != futurePlanSourceRef || record.FixturePathRef != "not_applicable" || record.WorkflowPathRef != "not_applicable" || record.ExpectedRef != "not_applicable" || record.ActualRef != "not_applicable" || record.CorrectiveActionRef != "not_applicable" {
				t.Fatalf("Phase 18 xserver future plan classification mismatch: %+v", record)
			}
			futurePlan = true
			continue
		}
		if !phase17StringInSlice(record.DriftTarget, phase18DocumentDriftTargets) || record.DriftResult != "no_drift" || record.ClosureResult != "closed" || record.CounterKey != "phase18_document_drift_open_count" {
			t.Fatalf("Phase 18 document drift record mismatch: %+v", record)
		}
		targetCoverage[record.DriftTarget] = true
	}
	for _, target := range phase18DocumentDriftTargets {
		if !targetCoverage[target] {
			t.Fatalf("Phase 18 document drift records missing target %s", target)
		}
	}
	if !futurePlan {
		t.Fatalf("Phase 18 document drift records must keep xserver-vps-future as non-blocking future_plan")
	}
}

func phase18RequireDerivedCounters(t *testing.T, counters map[string]int, records map[string][]phase18EvidenceRecord) {
	t.Helper()

	derived := phase18DeriveCounters(records)
	for _, key := range phase18Counters {
		if counters[key] != derived[key] {
			t.Fatalf("Phase 18 counter %s must be derived from records: want %d got %d", key, derived[key], counters[key])
		}
	}
}

func phase18DeriveCounters(records map[string][]phase18EvidenceRecord) map[string]int {
	derived := map[string]int{}
	for _, key := range phase18Counters {
		derived[key] = 0
	}
	for _, recordSet := range records {
		for _, record := range recordSet {
			if record.Result == "failed" || record.Result == "rejected" {
				derived[record.CounterKey] = 1
			}
		}
	}
	for _, key := range phase18Counters {
		if key != "final_open_item_count" {
			derived["final_open_item_count"] += derived[key]
		}
	}
	return derived
}

func phase18RequireWorkflow(t *testing.T, path string) {
	t.Helper()

	workflow := phase11MustReadText(t, path)
	if !strings.Contains(workflow, "name: phase18-vps-pull-bootstrap") || !strings.Contains(workflow, "permissions:\n  contents: read") {
		t.Fatalf("%s must use the Phase 18 workflow name and contents:read permissions", path)
	}
	if !strings.Contains(workflow, "\n  phase18-vps-pull-bootstrap:") {
		t.Fatalf("%s must define primary Phase 18 job", path)
	}
	if !strings.Contains(workflow, "timeout-minutes: 30") {
		t.Fatalf("%s must set timeout-minutes <= 30 for Phase 18 job", path)
	}
	if !strings.Contains(workflow, "go test ./... -run TestPhase18VPSPullBootstrapEvidence -count=1") {
		t.Fatalf("%s must run the Phase 18 evidence checker directly", path)
	}
	phase13RequirePinnedActions(t, workflow)
}

func phase18RequireDocumentState(t *testing.T) {
	t.Helper()

	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	phase18Row := "| Phase 18 | VPS Pull Bootstrap。"
	phase18InProgressStatus := "| 実装中・bootstrap 検証未完了 | Phase 17 |"
	phase18ImplementedStatus := "| 実装済み | Phase 17 |"
	if !strings.Contains(roadmap, phase18Row) || !strings.Contains(roadmap, phase18InProgressStatus) {
		t.Fatalf("docs/ROADMAP.md must keep Phase 18 as 実装中・bootstrap 検証未完了 until real VPS Pull Bootstrap evidence exists")
	}
	if !strings.Contains(roadmap, "現在の active Phase は Phase 18") {
		t.Fatalf("docs/ROADMAP.md must keep Phase 18 as the active Phase until real VPS Pull Bootstrap is complete")
	}
	if strings.Contains(roadmap, phase18ImplementedStatus) || strings.Contains(roadmap, "Phase 18 も `実装済み`") || strings.Contains(roadmap, "Phase 18 は `仕様化済み・未実装`") {
		t.Fatalf("docs/ROADMAP.md must not mark Phase 18 implemented or unimplemented while real VPS Pull Bootstrap remains incomplete")
	}
	for _, feature := range []string{
		"Phase 18 VPS Pull Bootstrap / Phase 17 handover / bootstrap script artifact / topology role / source channel / bootstrap session / bootstrap artifact health / bootstrap runtime flow gate",
		"Phase 18 pull runner foundation / bootstrap update / rollback drill / bootstrap 中 issue / 仕様先行バグ修正 / 再検証 / known bug 0 gate",
		"Phase 18 secret / destructive operation / opaque metadata / document drift gate",
	} {
		if !strings.Contains(roadmap, "| 実装中・bootstrap 検証未完了 |") || !strings.Contains(roadmap, feature) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 18 feature in progress until real bootstrap completes: %s", feature)
		}
		if strings.Contains(roadmap, "| 実装済み | CI/CD Bootstrap | "+feature+" |") ||
			strings.Contains(roadmap, "| 実装済み | 試験本番運用 | "+feature+" |") ||
			strings.Contains(roadmap, "| 実装済み | 本番検証 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | CI/CD Bootstrap | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 試験本番運用 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 本番検証 | "+feature+" |") {
			t.Fatalf("docs/ROADMAP.md must not mark Phase 18 feature implemented or unimplemented while real bootstrap remains incomplete: %s", feature)
		}
	}
	expectedPaths := []string{"testdata/phase18/vps-pull-bootstrap/"}
	for _, rel := range phase18FixtureFiles {
		expectedPaths = append(expectedPaths, "testdata/phase18/vps-pull-bootstrap/"+rel)
	}
	expectedPaths = append(expectedPaths, ".github/workflows/phase18-vps-pull-bootstrap.yml")
	for _, rel := range expectedPaths {
		if !strings.Contains(documentIndex, "| `"+rel+"` |") && !strings.Contains(documentIndex, "| [`"+rel+"`]") {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index Phase 18 path %s", rel)
		}
	}
	phase18IndexStart := strings.Index(documentIndex, "<a id=\"phase-18-target-paths\"></a>")
	if phase18IndexStart < 0 {
		t.Fatalf("docs/DOCUMENT_INDEX.md must define phase-18-target-paths anchor")
	}
	phase18IndexSection := documentIndex[phase18IndexStart:]
	if strings.Contains(phase18IndexSection, "未作成") {
		t.Fatalf("docs/DOCUMENT_INDEX.md must not keep Phase 18 paths as 未作成")
	}
}

func phase18RequireSafeReference(t *testing.T, label string, value string) {
	t.Helper()

	if strings.TrimSpace(value) == "" {
		t.Fatalf("%s must not be empty", label)
	}
	if value == "not_applicable" || value == "not_required_for_non_open_sample" {
		return
	}
	if filepath.IsAbs(value) || strings.Contains(value, "../") || strings.Contains(value, `..\`) {
		t.Fatalf("%s must not use absolute or parent path reference: %s", label, value)
	}
	if phase18ContainsForbiddenEvidenceLiteral(value) || phase17LooksLikeIPv4(value) {
		t.Fatalf("%s contains forbidden literal: %s", label, value)
	}
}

func phase18RequireRecordReferences(t *testing.T, rel string, line int, record phase18EvidenceRecord) {
	t.Helper()

	references := map[string]string{
		"source_ref":                record.SourceRef,
		"evidence_ref":              record.EvidenceRef,
		"operator_approval_ref":     record.OperatorApprovalRef,
		"owner_ref":                 record.OwnerRef,
		"endpoint_or_command_ref":   record.EndpointOrCommandRef,
		"operation_ref":             record.OperationRef,
		"issue_ref":                 record.IssueRef,
		"spec_general_update_ref":   record.SpecGeneralUpdateRef,
		"responsibility_source_ref": record.ResponsibilitySourceRef,
		"implementation_fix_ref":    record.ImplementationFixRef,
		"revalidation_ref":          record.RevalidationRef,
		"document_ref":              record.DocumentRef,
		"fixture_path_ref":          record.FixturePathRef,
		"workflow_path_ref":         record.WorkflowPathRef,
		"expected_ref":              record.ExpectedRef,
		"actual_ref":                record.ActualRef,
		"corrective_action_ref":     record.CorrectiveActionRef,
	}
	for field, value := range references {
		if value == "" {
			continue
		}
		label := fmt.Sprintf("%s line %d %s", rel, line, field)
		phase18RequireSafeReference(t, label, value)
		phase18RequireResolvableReference(t, label, value)
	}
}

func phase18RequireResolvableReference(t *testing.T, label string, value string) {
	t.Helper()

	if value == "not_applicable" || value == "not_required_for_non_open_sample" {
		return
	}
	target := phase18ReferenceTarget(value)
	if target == "" {
		t.Fatalf("%s must include a non-empty reference target: %s", label, value)
	}
	target = strings.SplitN(target, "#", 2)[0]
	if target == "" {
		t.Fatalf("%s must not use an anchor-only reference: %s", label, value)
	}
	candidates := []string{target}
	if phase18LooksFixtureLocalReference(target) {
		candidates = append(candidates, filepath.Join("testdata/phase18/vps-pull-bootstrap", target))
	}
	for _, candidate := range candidates {
		if _, err := os.Lstat(candidate); err == nil {
			return
		}
	}
	t.Fatalf("%s must resolve to an existing repository or Phase 18 fixture path: %s", label, value)
}

func phase18ReferenceTarget(value string) string {
	if start := strings.LastIndex(value, "]("); strings.HasPrefix(value, "[") && start >= 0 && strings.HasSuffix(value, ")") {
		return strings.TrimSpace(value[start+2 : len(value)-1])
	}
	return strings.TrimSpace(value)
}

func phase18LooksFixtureLocalReference(target string) bool {
	if strings.HasPrefix(target, "docs/") ||
		strings.HasPrefix(target, ".github/") ||
		strings.HasPrefix(target, "testdata/") ||
		strings.HasPrefix(target, "components/") ||
		target == "main.go" ||
		target == "main_test.go" ||
		target == "README.md" ||
		target == "AGENTS.md" {
		return false
	}
	return true
}

func phase18ContainsForbiddenEvidenceLiteral(value string) bool {
	lower := strings.ToLower(value)
	for _, forbidden := range []string{
		"-----begin openssh private key-----",
		"-----begin rsa private key-----",
		"-----begin private key-----",
		"ghp_",
		"github_pat_",
		"xoxb-",
		"sk-proj-",
		"sk-svcacct-",
		"password=",
		"token=",
		"secret=",
		"http://",
		"https://",
	} {
		if strings.Contains(lower, forbidden) {
			return true
		}
	}
	for _, token := range strings.FieldsFunc(value, func(r rune) bool {
		return !(r == '.' || r == '-' || r == '_' || r == ':' || r == '/' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z')
	}) {
		if phase17LooksLikeIPv4(token) {
			return true
		}
	}
	return false
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
