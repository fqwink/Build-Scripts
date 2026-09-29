package components

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFixtureSingle(t *testing.T) {
	out := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--src", "../testdata/builder/single/source.md", "--out", out, "--title", "Fixture Site"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	index := readFile(t, filepath.Join(out, "index.html"))
	for _, want := range []string{
		`<h1 id="title" class="mh h1">Title<button class="hn-link" data-href="#title" aria-label="リンクをコピー">¶</button></h1>`,
		`<a href="#title">self</a>`,
		`<div class="cb-wrap" data-lang="bash">`,
		`<span class="cl">bash</span>`,
		`<li class="ml-task"><input type="checkbox" disabled checked>done</li>`,
		`<li class="ml-task"><input type="checkbox" disabled>todo</li>`,
	} {
		if !strings.Contains(index, want) {
			t.Fatalf("index.html missing %q", want)
		}
	}
	for _, rel := range []string{"index.html", "assets/style.css", "assets/app.js", "assets/search-index.json"} {
		if _, err := os.Stat(filepath.Join(out, rel)); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}
	if _, err := os.Stat(filepath.Join(out, "pages")); !os.IsNotExist(err) {
		t.Fatalf("pages directory must not exist for single input")
	}
	if !strings.Contains(stdout.String(), "[REPORT] pages=1") || !strings.Contains(stdout.String(), "theme=adlaire-default") {
		t.Fatalf("unexpected report: %s", stdout.String())
	}
	if strings.Contains(stdout.String(), "BROKEN_LINK") {
		t.Fatalf("unexpected broken link warning: %s", stdout.String())
	}
}

func TestFixtureDirectory(t *testing.T) {
	out := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--src", "../testdata/builder/site/docs", "--out", out, "--title", "Docs"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	for _, rel := range []string{"index.html", "pages/guide-setup.html", "pages/guide-setup-copy.html", "pages/intro.html", "assets/search-index.json"} {
		if _, err := os.Stat(filepath.Join(out, rel)); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}
	intro := readFile(t, filepath.Join(out, "pages/intro.html"))
	if !strings.Contains(intro, `href="guide-setup.html#setup"`) {
		t.Fatalf("intro link was not rewritten: %s", intro)
	}
	for _, rel := range []string{"pages/guide-setup.html", "pages/guide-setup-copy.html"} {
		body := readFile(t, filepath.Join(out, rel))
		if !strings.Contains(body, `id="setup"`) || strings.Contains(body, `id="setup-2"`) {
			t.Fatalf("unexpected heading slug in %s", rel)
		}
	}
	search := readFile(t, filepath.Join(out, "assets/search-index.json"))
	for _, want := range []string{"index.html#site-index", "pages/guide-setup.html#setup", "pages/guide-setup-copy.html#setup", "pages/intro.html#intro"} {
		if !strings.Contains(search, want) {
			t.Fatalf("search index missing %s: %s", want, search)
		}
	}
}

func TestFixtureErrors(t *testing.T) {
	out := t.TempDir()
	cases := []struct {
		args []string
		err  string
	}{
		{[]string{"--theme", "unknown"}, "unknown theme: unknown"},
		{[]string{"--src", "/path/not-found.md"}, "source not found: /path/not-found.md"},
		{[]string{"--src", "../testdata/builder/empty-dir", "--out", out}, "no markdown files found:"},
		{[]string{"--title", ""}, "title must not be empty"},
		{[]string{"--src=../testdata/builder/single/source.md"}, "unknown option: --src=../testdata/builder/single/source.md"},
		{[]string{"--src", "../testdata/builder/single/source.md", "--out", out, "--commit-sha", "abcdef0"}, "invalid commit sha: abcdef0"},
		{[]string{"--src", "../testdata/builder/single/source.md", "--out", out, "--build-id", "build-1"}, "invalid build id: build-1"},
		{[]string{"--src", "../testdata/builder/single/source.md", "--out", out, "--build-at", "2026-09-26T01:02:03+09:00"}, "invalid build at: 2026-09-26T01:02:03+09:00"},
	}
	for _, tc := range cases {
		var stdout, stderr bytes.Buffer
		code := RunBuild(tc.args, &stdout, &stderr)
		if code != 2 {
			t.Fatalf("%v exit=%d stderr=%s", tc.args, code, stderr.String())
		}
		if !strings.Contains(stderr.String(), tc.err) {
			t.Fatalf("%v stderr missing %q: %s", tc.args, tc.err, stderr.String())
		}
		if strings.Contains(stdout.String(), "[REPORT]") {
			t.Fatalf("%v emitted report on failure: %s", tc.args, stdout.String())
		}
	}
}

func TestFixtureHelpPriority(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--src", "bad\npath", "--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "Usage: adlaire-ci-build") {
		t.Fatalf("help output missing usage: %s", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for help, got %q", stderr.String())
	}
}

func TestFixtureBuildMetadata(t *testing.T) {
	out := t.TempDir()
	var stdout, stderr bytes.Buffer
	args := []string{
		"--src", "../testdata/builder/single/source.md",
		"--out", out,
		"--title", "Fixture Site",
		"--build-id", "b20260926010203-001",
		"--commit-sha", "abcdef0123456789abcdef0123456789abcdef01",
		"--build-at", "2026-09-26T01:02:03Z",
	}
	code := RunBuild(args, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	index := readFile(t, filepath.Join(out, "index.html"))
	for _, want := range []string{
		`<meta name="adlaire-build-id" content="b20260926010203-001">`,
		`<meta name="adlaire-commit-sha" content="abcdef0123456789abcdef0123456789abcdef01">`,
		`<meta name="adlaire-build-at" content="2026-09-26T01:02:03Z">`,
		`Generated at 2026-09-26T01:02:03Z`,
	} {
		if !strings.Contains(index, want) {
			t.Fatalf("index.html missing %q", want)
		}
	}
	for _, want := range []string{
		"build_id=b20260926010203-001",
		"commit_sha=abcdef0123456789abcdef0123456789abcdef01",
		"build_at=2026-09-26T01:02:03Z",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("report missing %q: %s", want, stdout.String())
		}
	}
}

func TestFixtureURLSafety(t *testing.T) {
	root := t.TempDir()
	docs := filepath.Join(root, "docs")
	writeFile(t, filepath.Join(docs, "intro.md"), `# Intro

## Safe

[ok](guide/setup.md#setup)
[external](https://example.com/path?q=1)
[email](mailto:team@example.com)
[missing](guide/setup.md?token=secret#missing)
[bad](javascript:alert)
[protocol](//example.com/x)
![escape](../secret.png)
![remote](https://example.com/img.png)
`)
	writeFile(t, filepath.Join(docs, "guide", "setup.md"), `# Setup
`)
	out := filepath.Join(root, "dist")

	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--src", docs, "--out", out, "--title", "Docs"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	for _, want := range []string{
		"UNSAFE_URL: link",
		"BUILDER28_PATH_OUTSIDE_BASE: image",
		"BROKEN_PAGE_LINK: guide/setup.md#missing",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q: %s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "token=secret") {
		t.Fatalf("warning leaked query string: %s", stdout.String())
	}
	intro := readFile(t, filepath.Join(out, "pages", "intro.html"))
	for _, want := range []string{
		`href="guide-setup.html#setup"`,
		`href="https://example.com/path?q=1" target="_blank" rel="noopener noreferrer"`,
		`href="mailto:team@example.com"`,
		`src="https://example.com/img.png"`,
	} {
		if !strings.Contains(intro, want) {
			t.Fatalf("intro.html missing %q: %s", want, intro)
		}
	}
	for _, forbidden := range []string{"javascript:alert", `href="//example.com/x"`, `src="../secret.png"`} {
		if strings.Contains(intro, forbidden) {
			t.Fatalf("intro.html contains forbidden URL %q: %s", forbidden, intro)
		}
	}
}

func TestFixtureURLSafetyStrictKeepsExistingOutput(t *testing.T) {
	root := t.TempDir()
	docs := filepath.Join(root, "docs")
	writeFile(t, filepath.Join(docs, "index.md"), `# Intro

[bad](javascript:alert)
`)
	out := filepath.Join(root, "dist")
	writeFile(t, filepath.Join(out, "index.html"), "old")

	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--src", docs, "--out", out, "--title", "Docs", "--strict"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for strict warning, got %q", stderr.String())
	}
	for _, want := range []string{"[WARN] UNSAFE_URL: link", "[REPORT] pages=1", "warnings=1"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q: %s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "Done") {
		t.Fatalf("strict warning must not emit Done: %s", stdout.String())
	}
	if got := readFile(t, filepath.Join(out, "index.html")); got != "old" {
		t.Fatalf("strict warning replaced existing output: %q", got)
	}
}

func TestFixtureAtomicStagingCollision(t *testing.T) {
	root := t.TempDir()
	docs := filepath.Join(root, "docs")
	writeFile(t, filepath.Join(docs, "index.md"), "# Intro\n")
	out := filepath.Join(root, "dist")
	staging := filepath.Join(root, "dist.tmp.1")
	if err := os.Mkdir(staging, 0755); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--src", docs, "--out", out, "--title", "Docs"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "output staging path already exists: "+staging) {
		t.Fatalf("stderr missing staging collision: %s", stderr.String())
	}
	if strings.Contains(stdout.String(), "[REPORT]") {
		t.Fatalf("staging collision emitted report: %s", stdout.String())
	}
	if _, err := os.Stat(staging); err != nil {
		t.Fatalf("staging entry was modified or removed: %v", err)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatalf("output should not be created on staging collision")
	}
}

func TestFixtureAtomicExistingOutputFileIsPreserved(t *testing.T) {
	root := t.TempDir()
	docs := filepath.Join(root, "docs")
	writeFile(t, filepath.Join(docs, "index.md"), "# Intro\n")
	out := filepath.Join(root, "dist")
	if err := os.WriteFile(out, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--src", docs, "--out", out, "--title", "Docs"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "output path is not directory: "+out) {
		t.Fatalf("stderr missing output type error: %s", stderr.String())
	}
	if strings.Contains(stdout.String(), "[REPORT]") {
		t.Fatalf("existing file failure emitted report: %s", stdout.String())
	}
	if got := readFile(t, out); got != "old" {
		t.Fatalf("existing output file was replaced: %q", got)
	}
}

func TestFixtureAtomicRejectsUnsafeSiteFilePath(t *testing.T) {
	out := filepath.Join(t.TempDir(), "dist")
	files := []siteFile{
		{Path: "../escape.html", Data: []byte("escape")},
		{Path: "index.html", Data: []byte("index")},
		{Path: "assets/style.css", Data: []byte("css")},
		{Path: "assets/app.js", Data: []byte("js")},
		{Path: "assets/search-index.json", Data: []byte("[]")},
	}
	warnings, err := writeAtomic(out, files)
	if err == nil {
		t.Fatalf("expected unsafe path error")
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if !strings.Contains(err.Error(), "invalid output path: ../escape.html") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(out); !os.IsNotExist(statErr) {
		t.Fatalf("output should not be created for unsafe path")
	}
}

func TestFixtureAtomicOutputModes(t *testing.T) {
	out := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := RunBuild([]string{"--src", "../testdata/builder/single/source.md", "--out", out, "--title", "Fixture Site"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	modeOf := func(path string) os.FileMode {
		t.Helper()
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		return info.Mode().Perm()
	}
	if got := modeOf(filepath.Join(out, "assets")); got != 0755 {
		t.Fatalf("assets mode=%#o", got)
	}
	for _, rel := range []string{"index.html", "assets/style.css", "assets/app.js", "assets/search-index.json"} {
		if got := modeOf(filepath.Join(out, rel)); got != 0644 {
			t.Fatalf("%s mode=%#o", rel, got)
		}
	}
}

func TestFixtureIdempotency(t *testing.T) {
	out1 := t.TempDir()
	out2 := t.TempDir()
	var stdout1, stderr1, stdout2, stderr2 bytes.Buffer
	args1 := []string{"--src", "../testdata/builder/single/source.md", "--out", out1, "--title", "Fixture Site"}
	args2 := []string{"--src", "../testdata/builder/single/source.md", "--out", out2, "--title", "Fixture Site"}
	if code := RunBuild(args1, &stdout1, &stderr1); code != 0 {
		t.Fatalf("first exit=%d stderr=%s", code, stderr1.String())
	}
	if code := RunBuild(args2, &stdout2, &stderr2); code != 0 {
		t.Fatalf("second exit=%d stderr=%s", code, stderr2.String())
	}
	for _, rel := range []string{"index.html", "assets/style.css", "assets/app.js", "assets/search-index.json"} {
		a := stableRead(t, filepath.Join(out1, rel))
		b := stableRead(t, filepath.Join(out2, rel))
		if !bytes.Equal(a, b) {
			t.Fatalf("%s differs", rel)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func writeFile(t *testing.T, path string, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
}

func stableRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := stableContent(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
