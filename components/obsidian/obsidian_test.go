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

func writeObsidianTestFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
