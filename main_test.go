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

func TestPhase12StandardArtifactInventory(t *testing.T) {
	t.Parallel()

	ownerPackages := []string{
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
	if !strings.Contains(roadmap, "現在の active Phase はなしとする。") {
		t.Fatalf("docs/ROADMAP.md must state that no active Phase remains after Phase 12 completion")
	}
	if !strings.Contains(roadmap, "初期実装 Phase 1 から Phase 12 まではすべて `実装済み`") {
		t.Fatalf("docs/ROADMAP.md must state that Phase 1 through Phase 12 are all implemented after Phase 12 closure")
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
		if strings.Contains(roadmap, "| 実装中・検証未完了 |") && strings.Contains(roadmap, feature) {
			t.Fatalf("docs/ROADMAP.md must not keep Phase 12 feature as 実装中・検証未完了 after closure: %s", feature)
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
