package obsidian

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOwnerFileContract(t *testing.T) {
	if Owner() != "obsidian" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executeOwnerFileContract(newOwnerFileContract()) {
		t.Fatalf("obsidian owner file contract failed")
	}
}

func TestRunBuildObsidianVault(t *testing.T) {
	root := t.TempDir()
	vault := filepath.Join(root, "vault")
	writeObsidianTestFile(t, filepath.Join(vault, "Home.md"), "# Home\n\nSee [[Projects/Build#Plan|Build Plan]].\n\n![[images/logo.png|Logo]]\n\n#phase14 #docs\n")
	writeObsidianTestFile(t, filepath.Join(vault, "Projects", "Build.md"), "# Build\n\n## Plan\n\nBack to [[Home]].\n")
	writeObsidianTestFile(t, filepath.Join(vault, "images", "logo.png"), "png")
	writeObsidianTestFile(t, filepath.Join(vault, ".obsidian", "workspace.json"), "{}")
	out := filepath.Join(root, "site")
	report := filepath.Join(out, "obsidian_map.json")
	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{
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
	if strings.Contains(stderr.String(), "OBSIDIAN_") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
	for _, path := range []string{"index.html", "assets/search-index.json", "assets/images/logo.png", "obsidian_map.json"} {
		if _, err := os.Stat(filepath.Join(out, path)); err != nil {
			t.Fatalf("missing output %s: %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(vault, "obsidian_map.json")); !os.IsNotExist(err) {
		t.Fatalf("input vault must remain read-only")
	}
	if !strings.Contains(stdout.String(), "obsidian_notes=2") || !strings.Contains(stdout.String(), "obsidian_assets=1") {
		t.Fatalf("missing obsidian report: %s", stdout.String())
	}
	var got obsidianMap
	data, err := os.ReadFile(report)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("report JSON invalid: %v", err)
	}
	if got.SchemaVersion != obsidianMapSchemaVersion || got.EntrySourcePath != "Home.md" || got.EntryNormalizedPath != "index.md" {
		t.Fatalf("unexpected report identity: %+v", got)
	}
	if len(got.Notes) != 2 || len(got.Assets) != 1 || len(got.Diagnostics) != 0 {
		t.Fatalf("unexpected report counts: %+v", got)
	}
}

func TestObsidianRejectsYAMLFrontmatter(t *testing.T) {
	root := t.TempDir()
	vault := filepath.Join(root, "vault")
	writeObsidianTestFile(t, filepath.Join(vault, "Home.md"), "---\ntitle: Home\n---\n# Home\n")
	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--input-mode", "obsidian-vault", "--obsidian-vault", vault, "--obsidian-entry", "Home.md", "--out", filepath.Join(root, "site")}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.HasPrefix(stderr.String(), "OBSIDIAN_YAML_FRONTMATTER_UNSUPPORTED\n") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout must stay empty on yaml rejection: %s", stdout.String())
	}
}

func TestObsidianInputModeRequired(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--obsidian-vault", "vault"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("exit=%d", code)
	}
	if strings.TrimSpace(stderr.String()) != "OBSIDIAN_INPUT_MODE_REQUIRED" {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestObsidianAmbiguousBasenameFailsWithoutStrict(t *testing.T) {
	root := t.TempDir()
	vault := filepath.Join(root, "vault")
	writeObsidianTestFile(t, filepath.Join(vault, "Home.md"), "# Home\n\nSee [[Guide]].\n")
	writeObsidianTestFile(t, filepath.Join(vault, "A", "Guide.md"), "# Guide A\n")
	writeObsidianTestFile(t, filepath.Join(vault, "B", "Guide.md"), "# Guide B\n")
	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{
		"--input-mode", "obsidian-vault",
		"--obsidian-vault", vault,
		"--obsidian-entry", "Home.md",
		"--out", filepath.Join(root, "site"),
	}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.HasPrefix(stderr.String(), "OBSIDIAN_AMBIGUOUS_NOTE_LINK\n") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout must stay empty on ambiguous basename: %s", stdout.String())
	}
}

func TestRunObsidianSyncPlanApplyRollback(t *testing.T) {
	withFixedSyncClock(t, 1000)
	root := t.TempDir()
	project := filepath.Join(root, "project")
	vault := filepath.Join(root, "vault")
	state := filepath.Join(root, "state")
	writeObsidianTestFile(t, filepath.Join(project, "docs", "project.md"), "# Project\n")
	writeObsidianTestFile(t, filepath.Join(vault, "docs", "vault.md"), "# Vault\n")
	mkdirObsidianTestDir(t, state)

	plan := runObsidianSyncPlan(t, project, vault, state, "bidirectional")
	if plan.Operations != 2 || plan.Conflicts != 0 || plan.Tombstones != 0 || plan.Applied {
		t.Fatalf("unexpected plan response: %+v", plan)
	}
	if _, err := os.Stat(filepath.Join(state, "sync_state.json")); !os.IsNotExist(err) {
		t.Fatalf("sync plan must not write sync_state.json")
	}

	var stdout, stderr bytes.Buffer
	code := RunObsidian([]string{
		"sync", "apply",
		"--project-root", project,
		"--vault", vault,
		"--state-dir", state,
		"--plan-file", "plans/sync.json",
		"--plan-hash", plan.PlanHash,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("apply exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("apply stderr must be empty: %s", stderr.String())
	}
	var apply obsidianSyncApplyResponse
	decodeObsidianSyncJSON(t, stdout.Bytes(), &apply)
	if apply.OperationsApplied != 2 || apply.RollbackFile == "" || apply.StateDigest == "" {
		t.Fatalf("unexpected apply response: %+v", apply)
	}
	for _, path := range []string{filepath.Join(project, "docs", "vault.md"), filepath.Join(vault, "docs", "project.md"), filepath.Join(state, "sync_state.json"), filepath.Join(state, filepath.FromSlash(apply.RollbackFile))} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing apply path %s: %v", path, err)
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = RunObsidian([]string{
		"sync", "rollback",
		"--project-root", project,
		"--vault", vault,
		"--state-dir", state,
		"--rollback-file", apply.RollbackFile,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("rollback exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	var rollback obsidianSyncRollbackResponse
	decodeObsidianSyncJSON(t, stdout.Bytes(), &rollback)
	if rollback.OperationsRolledBack != 2 || rollback.StateDigest == "" {
		t.Fatalf("unexpected rollback response: %+v", rollback)
	}
	for _, path := range []string{filepath.Join(project, "docs", "vault.md"), filepath.Join(vault, "docs", "project.md")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("rollback must remove created path %s", path)
		}
	}
}

func TestRunObsidianSyncPlanHashMismatch(t *testing.T) {
	withFixedSyncClock(t, 1100)
	root := t.TempDir()
	project := filepath.Join(root, "project")
	vault := filepath.Join(root, "vault")
	state := filepath.Join(root, "state")
	writeObsidianTestFile(t, filepath.Join(project, "docs", "project.md"), "# Project\n")
	mkdirObsidianTestDir(t, vault)
	mkdirObsidianTestDir(t, state)
	_ = runObsidianSyncPlan(t, project, vault, state, "export-only")

	var stdout, stderr bytes.Buffer
	code := RunObsidian([]string{
		"sync", "apply",
		"--project-root", project,
		"--vault", vault,
		"--state-dir", state,
		"--plan-file", "plans/sync.json",
		"--plan-hash", strings.Repeat("0", 64),
	}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected plan hash mismatch exit 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("failure stdout must be empty: %s", stdout.String())
	}
	if stderr.String() != "obsidian: OBSIDIAN_SYNC_PLAN_HASH_MISMATCH\n" {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunObsidianSyncConflictRequiresManualResolution(t *testing.T) {
	withFixedSyncClock(t, 1200)
	root := t.TempDir()
	project := filepath.Join(root, "project")
	vault := filepath.Join(root, "vault")
	state := filepath.Join(root, "state")
	writeObsidianTestFile(t, filepath.Join(project, "shared.md"), "# Project\n")
	writeObsidianTestFile(t, filepath.Join(vault, "shared.md"), "# Vault\n")
	mkdirObsidianTestDir(t, state)
	plan := runObsidianSyncPlan(t, project, vault, state, "bidirectional")
	if plan.Conflicts != 1 {
		t.Fatalf("expected one conflict, got %+v", plan)
	}

	var stdout, stderr bytes.Buffer
	code := RunObsidian([]string{
		"sync", "apply",
		"--project-root", project,
		"--vault", vault,
		"--state-dir", state,
		"--plan-file", "plans/sync.json",
		"--plan-hash", plan.PlanHash,
	}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected conflict exit 1, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("conflict stdout must be empty: %s", stdout.String())
	}
	if stderr.String() != "obsidian: OBSIDIAN_SYNC_CONFLICT\n" {
		t.Fatalf("unexpected conflict stderr: %s", stderr.String())
	}
}

func TestRunObsidianSyncRollbackRestoresPreviousContent(t *testing.T) {
	withFixedSyncClock(t, 1300)
	root := t.TempDir()
	project := filepath.Join(root, "project")
	vault := filepath.Join(root, "vault")
	state := filepath.Join(root, "state")
	writeObsidianTestFile(t, filepath.Join(project, "shared.md"), "# New\n")
	writeObsidianTestFile(t, filepath.Join(vault, "shared.md"), "# Old\n")
	mkdirObsidianTestDir(t, state)
	oldDigest, err := fileSHA256(filepath.Join(vault, "shared.md"))
	if err != nil {
		t.Fatalf("digest old vault file: %v", err)
	}
	initialState := obsidianSyncState{
		SchemaVersion:     obsidianSyncStateSchemaVersion,
		ProjectRootDigest: "",
		VaultRootDigest:   "",
		Entries: []obsidianSyncStateEntry{{
			Path:           "shared.md",
			ProjectDigest:  oldDigest,
			VaultDigest:    oldDigest,
			LastSyncDigest: oldDigest,
			LastSyncUnix:   1,
		}},
		Tombstones:  []obsidianSyncTombstone{},
		Conflicts:   []obsidianSyncConflict{},
		LastApplyID: "seed",
	}
	initialStateBytes, err := syncCanonicalJSON(initialState)
	if err != nil {
		t.Fatalf("encode initial state: %v", err)
	}
	if err := os.WriteFile(filepath.Join(state, "sync_state.json"), initialStateBytes, 0o600); err != nil {
		t.Fatalf("write initial state: %v", err)
	}

	plan := runObsidianSyncPlan(t, project, vault, state, "bidirectional")
	var stdout, stderr bytes.Buffer
	code := RunObsidian([]string{
		"sync", "apply",
		"--project-root", project,
		"--vault", vault,
		"--state-dir", state,
		"--plan-file", "plans/sync.json",
		"--plan-hash", plan.PlanHash,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("apply exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if got := readObsidianTestFile(t, filepath.Join(vault, "shared.md")); got != "# New\n" {
		t.Fatalf("vault should receive project content, got %q", got)
	}
	var apply obsidianSyncApplyResponse
	decodeObsidianSyncJSON(t, stdout.Bytes(), &apply)

	stdout.Reset()
	stderr.Reset()
	code = RunObsidian([]string{
		"sync", "rollback",
		"--project-root", project,
		"--vault", vault,
		"--state-dir", state,
		"--rollback-file", apply.RollbackFile,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("rollback exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if got := readObsidianTestFile(t, filepath.Join(vault, "shared.md")); got != "# Old\n" {
		t.Fatalf("vault rollback content mismatch: %q", got)
	}
}

func runObsidianSyncPlan(t *testing.T, project, vault, state, direction string) obsidianSyncPlanResponse {
	t.Helper()

	var stdout, stderr bytes.Buffer
	code := RunObsidian([]string{
		"sync", "plan",
		"--project-root", project,
		"--vault", vault,
		"--state-dir", state,
		"--plan-file", "plans/sync.json",
		"--direction", direction,
		"--delete-policy", "tombstone",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("plan exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("plan stderr must be empty: %s", stderr.String())
	}
	var response obsidianSyncPlanResponse
	decodeObsidianSyncJSON(t, stdout.Bytes(), &response)
	if response.Command != "sync plan" || response.PlanFile != "plans/sync.json" || response.PlanHash == "" {
		t.Fatalf("unexpected plan response: %+v", response)
	}
	return response
}

func decodeObsidianSyncJSON(t *testing.T, data []byte, target any) {
	t.Helper()
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Fatalf("sync JSON response must end with LF: %q", string(data))
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		t.Fatalf("decode sync JSON: %v\n%s", err, string(data))
	}
}

func withFixedSyncClock(t *testing.T, unix int64) {
	t.Helper()
	previous := obsidianSyncNowUnix
	obsidianSyncNowUnix = func() int64 { return unix }
	t.Cleanup(func() {
		obsidianSyncNowUnix = previous
	})
}

func mkdirObsidianTestDir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func writeObsidianTestFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readObsidianTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
