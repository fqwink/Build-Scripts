package main

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
	code := run([]string{"--src", "testdata/build_spec/single/source.md", "--out", out, "--title", "Fixture Site"}, &stdout, &stderr)
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
	code := run([]string{"--src", "testdata/build_spec/site/docs", "--out", out, "--title", "Docs"}, &stdout, &stderr)
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
	cases := []struct {
		args []string
		err  string
	}{
		{[]string{"--theme", "unknown"}, "unknown theme: unknown"},
		{[]string{"--src", "/path/not-found.md"}, "source not found: /path/not-found.md"},
		{[]string{"--src", "testdata/build_spec/empty-dir"}, "no markdown files found:"},
		{[]string{"--title", ""}, "title must not be empty"},
	}
	for _, tc := range cases {
		var stdout, stderr bytes.Buffer
		code := run(tc.args, &stdout, &stderr)
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

func TestFixtureIdempotency(t *testing.T) {
	out1 := t.TempDir()
	out2 := t.TempDir()
	var stdout1, stderr1, stdout2, stderr2 bytes.Buffer
	args1 := []string{"--src", "testdata/build_spec/single/source.md", "--out", out1, "--title", "Fixture Site"}
	args2 := []string{"--src", "testdata/build_spec/single/source.md", "--out", out2, "--title", "Fixture Site"}
	if code := run(args1, &stdout1, &stderr1); code != 0 {
		t.Fatalf("first exit=%d stderr=%s", code, stderr1.String())
	}
	if code := run(args2, &stdout2, &stderr2); code != 0 {
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

func stableRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := stableContent(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
