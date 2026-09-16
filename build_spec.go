package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type BuildConfig struct {
	Src     string
	Out     string
	Title   string
	Theme   string
	BaseDir string
	Strict  bool
}

var DefaultBuildConfig = BuildConfig{
	Src:     "/opt/adlaire-builder/repo/docs",
	Out:     "/opt/adlaire-builder/dist/site",
	Title:   "Adlaire Documentation",
	Theme:   "adlaire-default",
	BaseDir: "",
	Strict:  false,
}

type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string { return e.Msg }

type PageInput struct {
	SourcePath   string
	RelativePath string
	RawText      string
}

type Heading struct {
	Level int
	Text  string
	Slug  string
	Line  int
}

type PageData struct {
	Title              string
	Slug               string
	SourcePath         string
	RelativePath       string
	OutputPath         string
	HTML               string
	TocHTML            string
	Headings           []Heading
	Warnings           []string
	ReadingTimeMinutes int
	PlainBlocks        []plainBlock
}

type plainBlock struct {
	Anchor string
	Title  string
	Body   string
}

type RenderContext struct {
	FootnoteDefs     map[string]string
	FootnoteOrder    []string
	FootnoteSeen     map[string]bool
	InternalLinkRefs map[string]string
	BrokenLinks      []string
	HeadingSkipCount int
	CharCount        int
	LinkResolver     func(string) (string, string)
	KnownAnchors     map[string]bool
	Warnings         []string
}

type report struct {
	Pages          int
	Headings       int
	Tables         int
	CodeBlocks     int
	Warnings       int
	SizeWarn       bool
	BrokenLinks    int
	HeadingSkips   int
	ReadingTime    int
	Theme          string
	OutputFiles    int
	OutputBytes    int64
	OutputDir      string
}

type siteFile struct {
	Path string
	Data []byte
}

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(code)
}

func run(args []string, stdout, stderr io.Writer) int {
	cfg, handled, err := parseArgs(args, stdout)
	if err != nil {
		var ee exitError
		if errors.As(err, &ee) {
			fmt.Fprintln(stderr, ee.Msg)
			return ee.Code
		}
		fmt.Fprintln(stderr, err)
		return 1
	}
	if handled {
		return 0
	}
	rep, warnings, err := build(cfg, stdout)
	if err != nil {
		var ee exitError
		if errors.As(err, &ee) {
			fmt.Fprintln(stderr, ee.Msg)
			return ee.Code
		}
		fmt.Fprintln(stderr, err)
		return 1
	}
	for _, w := range warnings {
		fmt.Fprintf(stdout, "[WARN] %s\n", w)
	}
	fmt.Fprintf(stdout, "[REPORT] pages=%d headings=%d tables=%d code_blocks=%d warnings=%d size_warn=%t broken_links=%d heading_skips=%d reading_time=%d theme=%s\n",
		rep.Pages, rep.Headings, rep.Tables, rep.CodeBlocks, rep.Warnings, rep.SizeWarn, rep.BrokenLinks, rep.HeadingSkips, rep.ReadingTime, rep.Theme)
	return 0
}

func parseArgs(args []string, stdout io.Writer) (BuildConfig, bool, error) {
	cfg := DefaultBuildConfig
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--help" {
			fmt.Fprintln(stdout, "Usage: adlaire-ci-build [--src path] [--out path] [--title text] [--theme name] [--base-dir path] [--strict] [--version] [--help]")
			return cfg, true, nil
		}
		if arg == "--version" {
			fmt.Fprintf(stdout, "adlaire-ci-build ADLAIRE_CI_SPEC go=%s\n", runtime.Version())
			return cfg, true, nil
		}
		if arg == "--strict" {
			cfg.Strict = true
			continue
		}
		if !strings.HasPrefix(arg, "--") {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
		name, value, hasValue := strings.Cut(arg, "=")
		if !hasValue {
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
				return cfg, false, exitError{Code: 2, Msg: "missing value: " + name}
			}
			i++
			value = args[i]
		}
		switch name {
		case "--src":
			cfg.Src = value
		case "--out":
			cfg.Out = value
		case "--title":
			cfg.Title = value
		case "--theme":
			cfg.Theme = value
		case "--base-dir":
			cfg.BaseDir = value
		default:
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + name}
		}
	}
	if cfg.Title == "" {
		return cfg, false, exitError{Code: 2, Msg: "title must not be empty"}
	}
	if cfg.Theme != "adlaire-default" {
		return cfg, false, exitError{Code: 2, Msg: "unknown theme: " + cfg.Theme}
	}
	var err error
	cfg.Src, err = absPath(cfg.Src)
	if err != nil {
		return cfg, false, err
	}
	cfg.Out, err = absPath(cfg.Out)
	if err != nil {
		return cfg, false, err
	}
	if cfg.BaseDir != "" {
		cfg.BaseDir, err = absPath(cfg.BaseDir)
		if err != nil {
			return cfg, false, err
		}
		info, statErr := os.Stat(cfg.BaseDir)
		if statErr != nil {
			return cfg, false, exitError{Code: 2, Msg: "base directory not found: " + cfg.BaseDir}
		}
		if !info.IsDir() {
			return cfg, false, exitError{Code: 2, Msg: "base path is not directory: " + cfg.BaseDir}
		}
	}
	return cfg, false, nil
}

func absPath(p string) (string, error) {
	if filepath.IsAbs(p) {
		return filepath.Clean(p), nil
	}
	return filepath.Abs(p)
}

func build(cfg BuildConfig, stdout io.Writer) (report, []string, error) {
	fmt.Fprintln(stdout, "Collecting Markdown...")
	inputs, baseDir, err := collectInputs(cfg)
	if err != nil {
		return report{}, nil, err
	}
	fmt.Fprintln(stdout, "Converting MD...")
	srcInfo, _ := os.Stat(cfg.Src)
	pages, rep, warnings, err := renderPages(baseDir, inputs, srcInfo != nil && !srcInfo.IsDir())
	if err != nil {
		return report{}, nil, err
	}
	if cfg.Strict && len(warnings) > 0 {
		return report{}, warnings, exitError{Code: 2, Msg: "strict mode failed with warnings"}
	}
	fmt.Fprintln(stdout, "Building site...")
	files, err := assembleSite(cfg, pages)
	if err != nil {
		return report{}, warnings, exitError{Code: 1, Msg: err.Error()}
	}
	fmt.Fprintln(stdout, "Writing assets...")
	if err := writeAtomic(cfg.Out, files); err != nil {
		return report{}, warnings, exitError{Code: 1, Msg: err.Error()}
	}
	count, size, _ := outputStats(cfg.Out)
	rep.Pages = htmlPageCount(files)
	rep.Warnings = len(warnings)
	rep.Theme = cfg.Theme
	rep.OutputFiles = count
	rep.OutputBytes = size
	rep.OutputDir = cfg.Out
	fmt.Fprintf(stdout, "Done → %s  (pages=%d files=%d bytes=%d)\n", cfg.Out, rep.Pages, count, size)
	return rep, warnings, nil
}

func collectInputs(cfg BuildConfig) ([]PageInput, string, error) {
	info, err := os.Stat(cfg.Src)
	if err != nil {
		return nil, "", exitError{Code: 2, Msg: "source not found: " + cfg.Src}
	}
	baseDir := cfg.BaseDir
	if baseDir == "" {
		if info.IsDir() {
			baseDir = cfg.Src
		} else {
			baseDir = filepath.Dir(cfg.Src)
		}
	}
	if !info.IsDir() {
		if !isMarkdownPath(cfg.Src) {
			return nil, "", exitError{Code: 2, Msg: "source is not markdown file or directory: " + cfg.Src}
		}
		text, err := readUTF8(cfg.Src)
		if err != nil {
			return nil, "", err
		}
		rel, _ := filepath.Rel(baseDir, cfg.Src)
		return []PageInput{{SourcePath: cfg.Src, RelativePath: slashPath(rel), RawText: text}}, baseDir, nil
	}
	var paths []string
	err = filepath.WalkDir(cfg.Src, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "dist" {
				if path != cfg.Src {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if isMarkdownPath(path) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, "", exitError{Code: 2, Msg: err.Error()}
	}
	sort.Slice(paths, func(i, j int) bool {
		ri, _ := filepath.Rel(baseDir, paths[i])
		rj, _ := filepath.Rel(baseDir, paths[j])
		return slashPath(ri) < slashPath(rj)
	})
	if len(paths) == 0 {
		return nil, "", exitError{Code: 2, Msg: "no markdown files found: " + cfg.Src}
	}
	inputs := make([]PageInput, 0, len(paths))
	for _, path := range paths {
		text, err := readUTF8(path)
		if err != nil {
			return nil, "", err
		}
		rel, _ := filepath.Rel(baseDir, path)
		inputs = append(inputs, PageInput{SourcePath: path, RelativePath: slashPath(rel), RawText: text})
	}
	return inputs, baseDir, nil
}

func isMarkdownPath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

func readUTF8(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", exitError{Code: 2, Msg: "source not found: " + path}
	}
	if !utf8.Valid(data) {
		return "", exitError{Code: 2, Msg: "source is not valid UTF-8: " + path}
	}
	return string(data), nil
}

func renderPages(baseDir string, inputs []PageInput, isSingle bool) ([]PageData, report, []string, error) {
	outputBySource := map[string]string{}
	slugCounts := map[string]int{}
	for _, in := range inputs {
		slug := uniqueSlug(pageSlug(in.RelativePath), slugCounts)
		out := "index.html"
		if !isSingle {
			out = "pages/" + slug + ".html"
		}
		outputBySource[cleanAbs(in.SourcePath)] = out
	}
	var pages []PageData
	var rep report
	var warnings []string
	for _, in := range inputs {
		lines := splitLines(in.RawText)
		headings, slugByLine, headingSkips := collectHeadings(lines)
		ctx := &RenderContext{
			FootnoteDefs:     collectFootnotes(lines),
			FootnoteSeen:     map[string]bool{},
			InternalLinkRefs: map[string]string{},
			KnownAnchors:     map[string]bool{},
		}
		for _, h := range headings {
			ctx.KnownAnchors[h.Slug] = true
		}
		ctx.LinkResolver = func(raw string) (string, string) {
			return resolveLink(raw, in.SourcePath, outputBySource, isSingle)
		}
		body, plain, tables, codeBlocks, err := convert(lines, slugByLine, ctx)
		if err != nil {
			return nil, rep, warnings, err
		}
		body = injectChapterNavigation(body, headings, &ctx.Warnings)
		title := pageTitle(in.RelativePath, headings)
		outPath := outputBySource[cleanAbs(in.SourcePath)]
		pageWarnings := append([]string{}, ctx.Warnings...)
		for _, broken := range ctx.BrokenLinks {
			pageWarnings = append(pageWarnings, broken)
		}
		warnings = append(warnings, pageWarnings...)
		rep.Headings += len(headings)
		rep.Tables += tables
		rep.CodeBlocks += codeBlocks
		rep.BrokenLinks += len(ctx.BrokenLinks)
		rep.HeadingSkips += headingSkips + ctx.HeadingSkipCount
		rep.ReadingTime += readingTime(ctx.CharCount)
		pages = append(pages, PageData{
			Title:              title,
			Slug:               strings.TrimSuffix(filepath.Base(outPath), ".html"),
			SourcePath:         in.SourcePath,
			RelativePath:       in.RelativePath,
			OutputPath:         outPath,
			HTML:               body,
			TocHTML:            buildTOC(headings),
			Headings:           headings,
			Warnings:           pageWarnings,
			ReadingTimeMinutes: readingTime(ctx.CharCount),
			PlainBlocks:        plain,
		})
	}
	return pages, rep, warnings, nil
}

func splitLines(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return strings.Split(text, "\n")
}

func collectHeadings(lines []string) ([]Heading, map[int]string, int) {
	var headings []Heading
	slugByLine := map[int]string{}
	counts := map[string]int{}
	inFence := false
	lastLevel := 0
	skips := 0
	for i, line := range lines {
		if isFenceLine(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if m := headingRe.FindStringSubmatch(line); m != nil {
			level := len(m[1])
			text := strings.TrimSpace(m[2])
			base := slugify(text)
			slug := uniqueSlug(base, counts)
			headings = append(headings, Heading{Level: level, Text: text, Slug: slug, Line: i})
			slugByLine[i] = slug
			if lastLevel > 0 && level-lastLevel >= 2 {
				skips++
			}
			lastLevel = level
		}
	}
	return headings, slugByLine, skips
}

var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
var footnoteDefRe = regexp.MustCompile(`^\[\^([^\]]+)\]:\s*(.*)$`)

func isFenceLine(line string) bool {
	s := strings.TrimSpace(line)
	return strings.HasPrefix(s, "```") || strings.HasPrefix(s, "~~~")
}

func collectFootnotes(lines []string) map[string]string {
	defs := map[string]string{}
	inFence := false
	for _, line := range lines {
		if isFenceLine(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if m := footnoteDefRe.FindStringSubmatch(line); m != nil {
			defs[m[1]] = strings.TrimSpace(m[2])
		}
	}
	return defs
}

func pageTitle(rel string, headings []Heading) string {
	for _, h := range headings {
		if h.Level == 1 {
			return h.Text
		}
	}
	base := strings.TrimSuffix(filepath.Base(rel), filepath.Ext(rel))
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	base = strings.TrimSpace(base)
	if base == "" {
		return "Untitled"
	}
	return base
}

func pageSlug(rel string) string {
	noExt := strings.TrimSuffix(slashPath(rel), filepath.Ext(rel))
	noExt = strings.NewReplacer("/", "-", " ", "-", "_", "-").Replace(noExt)
	return slugify(noExt)
}

func slugify(text string) string {
	text = strings.NewReplacer("`", "", "*", "", "_", "", "~", "", "[", "", "]", "").Replace(text)
	var b strings.Builder
	prevDash := false
	for _, r := range text {
		isDash := unicode.IsSpace(r) || r == '-' || r == '.' || r == '_'
		switch {
		case isDash:
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			b.WriteRune(unicode.ToLower(r))
			prevDash = false
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "section"
	}
	return out
}

func uniqueSlug(base string, counts map[string]int) string {
	n := counts[base]
	counts[base] = n + 1
	if n == 0 {
		return base
	}
	return fmt.Sprintf("%s-%d", base, n+1)
}

func esc(s string) string {
	return html.EscapeString(s)
}

func inline(text string, ctx *RenderContext) string {
	segments := codeSpanSegments(text)
	out := strings.Join(segments, "")
	out = replaceInlineMarkup(out)
	out = replaceImages(out)
	out = replaceLinks(out, ctx)
	out = replaceFootnotes(out, ctx)
	return out
}

func codeSpanSegments(text string) []string {
	var segments []string
	for i := 0; i < len(text); {
		bt := strings.IndexByte(text[i:], '`')
		if bt < 0 {
			segments = append(segments, esc(text[i:]))
			break
		}
		bt += i
		segments = append(segments, esc(text[i:bt]))
		marker := "`"
		if strings.HasPrefix(text[bt:], "``") {
			marker = "``"
		}
		end := strings.Index(text[bt+len(marker):], marker)
		if end < 0 {
			segments = append(segments, esc(text[bt:]))
			break
		}
		start := bt + len(marker)
		end += start
		segments = append(segments, `<code class="ic">`+esc(text[start:end])+`</code>`)
		i = end + len(marker)
	}
	return segments
}

func replaceInlineMarkup(text string) string {
	repls := []struct {
		re   *regexp.Regexp
		repl string
	}{
		{regexp.MustCompile(`\*\*\*(.+?)\*\*\*`), `<strong><em>$1</em></strong>`},
		{regexp.MustCompile(`\*\*(.+?)\*\*`), `<strong>$1</strong>`},
		{regexp.MustCompile(`__(.+?)__`), `<strong>$1</strong>`},
		{regexp.MustCompile(`~~(.+?)~~`), `<del>$1</del>`},
		{regexp.MustCompile(`\*(.+?)\*`), `<em>$1</em>`},
		{regexp.MustCompile(`_(.+?)_`), `<em>$1</em>`},
	}
	for _, r := range repls {
		text = r.re.ReplaceAllString(text, r.repl)
	}
	return text
}

func replaceImages(text string) string {
	re := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		return `<img src="` + esc(parts[2]) + `" alt="` + esc(parts[1]) + `" style="max-width:100%">`
	})
}

func replaceLinks(text string, ctx *RenderContext) string {
	re := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		label := parts[1]
		raw := parts[2]
		href := raw
		if ctx != nil && ctx.LinkResolver != nil {
			resolved, warn := ctx.LinkResolver(raw)
			if warn != "" {
				ctx.Warnings = append(ctx.Warnings, warn)
			}
			href = resolved
		}
		if strings.HasPrefix(raw, "#") && ctx != nil {
			anchor := strings.TrimPrefix(raw, "#")
			ctx.InternalLinkRefs[anchor] = m
			if !ctx.KnownAnchors[anchor] {
				msg := "BROKEN_LINK: " + raw + "  (in: " + m + ")"
				ctx.BrokenLinks = append(ctx.BrokenLinks, msg)
			}
		}
		if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
			return `<a href="` + esc(href) + `" target="_blank" rel="noopener noreferrer">` + label + `</a>`
		}
		return `<a href="` + esc(href) + `">` + label + `</a>`
	})
}

func replaceFootnotes(text string, ctx *RenderContext) string {
	if ctx == nil {
		return text
	}
	re := regexp.MustCompile(`\[\^([^\]]+)\]`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		key := re.FindStringSubmatch(m)[1]
		if !ctx.FootnoteSeen[key] {
			ctx.FootnoteSeen[key] = true
			ctx.FootnoteOrder = append(ctx.FootnoteOrder, key)
		}
		n := 0
		for i, k := range ctx.FootnoteOrder {
			if k == key {
				n = i + 1
				break
			}
		}
		id := esc(key)
		return fmt.Sprintf(`<sup><a href="#fn-%s" id="fnref-%s" class="fn-ref">[%d]</a></sup>`, id, id, n)
	})
}

func resolveLink(raw, sourcePath string, outputBySource map[string]string, isSingle bool) (string, string) {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "mailto:") || strings.HasPrefix(raw, "tel:") || strings.HasPrefix(raw, "#") {
		return raw, ""
	}
	u, err := url.Parse(raw)
	if err != nil || u.Path == "" || !isMarkdownPath(u.Path) {
		return raw, ""
	}
	target := filepath.Clean(filepath.Join(filepath.Dir(sourcePath), filepath.FromSlash(u.Path)))
	out, ok := outputBySource[cleanAbs(target)]
	if !ok {
		return raw, "BROKEN_PAGE_LINK: " + raw + " (in: " + sourcePath + ")"
	}
	if !isSingle && strings.HasPrefix(out, "pages/") {
		out = strings.TrimPrefix(out, "pages/")
	}
	if u.Fragment != "" {
		out += "#" + u.Fragment
	}
	return out, ""
}

func convert(lines []string, slugByLine map[int]string, ctx *RenderContext) (string, []plainBlock, int, int, error) {
	var out []string
	var plain []plainBlock
	var para []string
	var table []string
	var listOpen string
	var codeBuf []string
	var codeLang string
	var fence string
	inFence := false
	tables := 0
	codeBlocks := 0
	currentAnchor := ""
	currentTitle := ""
	flushPara := func() {
		if len(para) == 0 {
			return
		}
		text := strings.Join(para, " ")
		ctx.CharCount += utf8.RuneCountInString(text)
		out = append(out, `<p class="mp">`+inline(text, ctx)+`</p>`)
		plain = append(plain, plainBlock{Anchor: currentAnchor, Title: currentTitle, Body: normalizePlain(text)})
		para = nil
	}
	flushList := func() {
		if listOpen != "" {
			out = append(out, "</"+listOpen+">")
			listOpen = ""
		}
	}
	flushTable := func() {
		if len(table) == 0 {
			return
		}
		tables++
		out = append(out, `<div class="tw"><table class="mt">`)
		sep := -1
		for i, row := range table {
			cells := splitTableCells(row)
			if len(cells) > 0 {
				allSep := true
				for _, c := range cells {
					c = strings.TrimSpace(c)
					if !regexp.MustCompile(`^:?-+:?$`).MatchString(c) {
						allSep = false
					}
				}
				if allSep {
					sep = i
					break
				}
			}
		}
		for i, row := range table {
			if i == sep {
				continue
			}
			tag := "td"
			if sep >= 0 && i < sep {
				tag = "th"
			}
			cells := splitTableCells(row)
			var b strings.Builder
			b.WriteString("<tr>")
			for idx, c := range cells {
				if tag == "th" {
					b.WriteString(fmt.Sprintf(`<th data-sort="%d" aria-sort="none">%s</th>`, idx, inline(strings.TrimSpace(c), ctx)))
				} else {
					b.WriteString("<td>" + inline(strings.TrimSpace(c), ctx) + "</td>")
				}
			}
			b.WriteString("</tr>")
			out = append(out, b.String())
		}
		out = append(out, `</table></div>`)
		table = nil
	}
	emitCode := func() {
		codeBlocks++
		code := strings.Join(codeBuf, "\n")
		linesCount := 0
		if code != "" {
			linesCount = strings.Count(code, "\n") + 1
		}
		collapsible := ""
		button := ""
		if linesCount > 30 {
			collapsible = fmt.Sprintf(` data-collapsible="true" data-total-lines="%d"`, linesCount)
			button = fmt.Sprintf(`<button class="expand-code">全 %d 行を表示</button>`, linesCount)
		}
		lang := esc(codeLang)
		label := ""
		if lang != "" {
			label = `<span class="cl">` + lang + `</span>`
		}
		out = append(out, fmt.Sprintf(`<div class="cb-wrap" data-lang="%s"><div class="cb-meta">%s<button class="cb-copy" aria-label="コピー">コピー</button></div><pre class="cb"%s><code>%s</code></pre>%s</div>`, lang, label, collapsible, esc(code), button))
	}
	for i, line := range lines {
		raw := strings.TrimRight(line, "\n")
		trim := strings.TrimSpace(raw)
		if !inFence {
			if strings.HasPrefix(trim, "```") || strings.HasPrefix(trim, "~~~") {
				flushPara()
				flushList()
				flushTable()
				inFence = true
				fence = trim[:3]
				codeLang = strings.TrimSpace(trim[3:])
				codeBuf = nil
				continue
			}
		} else {
			if trim == fence {
				emitCode()
				inFence = false
				codeLang = ""
				codeBuf = nil
			} else {
				codeBuf = append(codeBuf, raw)
			}
			continue
		}
		if footnoteDefRe.MatchString(raw) {
			continue
		}
		if m := headingRe.FindStringSubmatch(raw); m != nil {
			flushPara()
			flushList()
			flushTable()
			level := len(m[1])
			if level > 4 {
				level = 4
			}
			text := strings.TrimSpace(m[2])
			slug := slugByLine[i]
			currentAnchor = slug
			currentTitle = text
			out = append(out, fmt.Sprintf(`<h%d id="%s" class="mh h%d">%s<button class="hn-link" data-href="#%s" aria-label="リンクをコピー">¶</button></h%d>`, level, slug, level, inline(text, ctx), slug, level))
			plain = append(plain, plainBlock{Anchor: slug, Title: text, Body: normalizePlain(text)})
			continue
		}
		if trim == "" {
			flushPara()
			flushList()
			flushTable()
			continue
		}
		if strings.Contains(raw, "|") && strings.Count(raw, "|") >= 2 {
			flushPara()
			flushList()
			table = append(table, raw)
			continue
		}
		flushTable()
		if strings.HasPrefix(trim, ">") {
			flushPara()
			flushList()
			out = append(out, `<blockquote class="mbq">`+inline(strings.TrimSpace(strings.TrimPrefix(trim, ">")), ctx)+`</blockquote>`)
			continue
		}
		if trim == "---" || trim == "***" {
			flushPara()
			flushList()
			out = append(out, `<hr class="mr">`)
			continue
		}
		if item, ok := parseList(trim); ok {
			flushPara()
			tag := "ul"
			if item.ordered {
				tag = "ol"
			}
			if listOpen != tag {
				flushList()
				listOpen = tag
				out = append(out, `<`+tag+` class="ml">`)
			}
			if item.task {
				checked := ""
				if item.checked {
					checked = " checked"
				}
				out = append(out, `<li class="ml-task"><input type="checkbox" disabled`+checked+`>`+inline(item.text, ctx)+`</li>`)
			} else {
				out = append(out, `<li>`+inline(item.text, ctx)+`</li>`)
			}
			continue
		}
		para = append(para, trim)
	}
	if inFence {
		return "", nil, tables, codeBlocks, exitError{Code: 2, Msg: "unclosed code fence"}
	}
	flushPara()
	flushList()
	flushTable()
	if len(ctx.FootnoteOrder) > 0 {
		out = append(out, `<section class="fn-section"><ol class="fn-list">`)
		for _, key := range ctx.FootnoteOrder {
			body := ctx.FootnoteDefs[key]
			out = append(out, fmt.Sprintf(`<li id="fn-%s" class="fn-item">%s <a href="#fnref-%s" class="fn-back">↩</a></li>`, esc(key), inline(body, ctx), esc(key)))
		}
		out = append(out, `</ol></section>`)
	}
	return strings.Join(out, "\n"), plain, tables, codeBlocks, nil
}

type listItem struct {
	ordered bool
	task    bool
	checked bool
	text    string
}

func parseList(trim string) (listItem, bool) {
	if strings.HasPrefix(trim, "- ") || strings.HasPrefix(trim, "* ") {
		text := strings.TrimSpace(trim[2:])
		item := listItem{text: text}
		if strings.HasPrefix(text, "[x] ") || strings.HasPrefix(text, "[X] ") || strings.HasPrefix(text, "[ ] ") {
			item.task = true
			item.checked = strings.HasPrefix(text, "[x] ") || strings.HasPrefix(text, "[X] ")
			item.text = strings.TrimSpace(text[4:])
		}
		return item, true
	}
	m := regexp.MustCompile(`^\d+\.\s+(.+)$`).FindStringSubmatch(trim)
	if m != nil {
		return listItem{ordered: true, text: m[1]}, true
	}
	return listItem{}, false
}

func splitTableCells(row string) []string {
	row = strings.TrimSpace(row)
	row = strings.Trim(row, "|")
	return strings.Split(row, "|")
}

func buildTOC(headings []Heading) string {
	var lines []string
	lines = append(lines, `<ul class="tr">`)
	for _, h := range headings {
		if h.Level > 3 {
			continue
		}
		lines = append(lines, fmt.Sprintf(`<li class="ti"><a href="#%s" class="tl lv%d" data-slug="%s">%s</a></li>`, h.Slug, h.Level, h.Slug, esc(h.Text)))
	}
	lines = append(lines, `</ul>`)
	return strings.Join(lines, "\n")
}

func injectChapterNavigation(body string, headings []Heading, warnings *[]string) string {
	var h2 []Heading
	for _, h := range headings {
		if h.Level == 2 {
			h2 = append(h2, h)
		}
	}
	if len(h2) <= 1 {
		return body
	}
	for i, h := range h2 {
		next := ""
		prev := ""
		if i > 0 {
			prev = fmt.Sprintf(`<a class="ch-prev" href="#%s">← %s</a>`, h2[i-1].Slug, esc(h2[i-1].Text))
		} else {
			prev = `<span></span>`
		}
		if i < len(h2)-1 {
			next = fmt.Sprintf(`<a class="ch-next" href="#%s">%s →</a>`, h2[i+1].Slug, esc(h2[i+1].Text))
		} else {
			next = `<span></span>`
		}
		nav := `<nav class="ch-nav">` + prev + next + `</nav>`
		if i == len(h2)-1 {
			body += "\n" + nav
			continue
		}
		marker := `id="` + h2[i+1].Slug + `"`
		pos := strings.Index(body, marker)
		if pos < 0 {
			*warnings = append(*warnings, "CHAPTER_NAV_SKIPPED: slug="+h.Slug)
			continue
		}
		tagStart := strings.LastIndex(body[:pos], "<h")
		if tagStart < 0 {
			*warnings = append(*warnings, "CHAPTER_NAV_SKIPPED: slug="+h.Slug)
			continue
		}
		body = body[:tagStart] + nav + "\n" + body[tagStart:]
	}
	return body
}

func assembleSite(cfg BuildConfig, pages []PageData) ([]siteFile, error) {
	var files []siteFile
	style := []byte(defaultCSS())
	app := []byte(defaultJS())
	search, err := json.Marshal(searchIndex(pages))
	if err != nil {
		return nil, err
	}
	files = append(files, siteFile{Path: "assets/style.css", Data: style})
	files = append(files, siteFile{Path: "assets/app.js", Data: app})
	files = append(files, siteFile{Path: "assets/search-index.json", Data: search})
	if len(pages) == 1 && pages[0].OutputPath == "index.html" {
		files = append(files, siteFile{Path: "index.html", Data: []byte(pageHTML(cfg, pages[0]))})
		return files, nil
	}
	indexPage := PageData{
		Title:              cfg.Title,
		OutputPath:         "index.html",
		HTML:               siteIndexHTML(pages),
		TocHTML:            `<ul class="tr"><li class="ti"><a href="#site-index" class="tl lv1" data-slug="site-index">` + esc(cfg.Title) + `</a></li></ul>`,
		ReadingTimeMinutes: 1,
		Headings:           []Heading{{Level: 1, Text: cfg.Title, Slug: "site-index"}},
	}
	files = append(files, siteFile{Path: "index.html", Data: []byte(pageHTML(cfg, indexPage))})
	for _, p := range pages {
		files = append(files, siteFile{Path: p.OutputPath, Data: []byte(pageHTML(cfg, p))})
	}
	return files, nil
}

type searchEntry struct {
	URL   string `json:"url"`
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func searchIndex(pages []PageData) []searchEntry {
	entries := []searchEntry{{URL: "index.html#site-index", ID: "site-index", Title: "Site Index", Body: "Site Index"}}
	single := len(pages) == 1 && pages[0].OutputPath == "index.html"
	if single {
		entries = entries[:0]
	}
	for _, p := range pages {
		for i, block := range p.PlainBlocks {
			id := block.Anchor
			if id == "" {
				id = p.Slug
			}
			title := block.Title
			if title == "" {
				title = p.Title
			}
			url := p.OutputPath + "#" + id
			if single && i == 0 {
				url = "index.html#" + id
			}
			entries = append(entries, searchEntry{URL: url, ID: id, Title: title, Body: truncateRunes(block.Body, 200)})
		}
	}
	return entries
}

func siteIndexHTML(pages []PageData) string {
	var b strings.Builder
	b.WriteString(`<h1 id="site-index" class="mh h1">Site Index<button class="hn-link" data-href="#site-index" aria-label="リンクをコピー">¶</button></h1>`)
	b.WriteString(`<ul class="ml">`)
	for _, p := range pages {
		b.WriteString(`<li><a href="` + esc(p.OutputPath) + `">` + esc(p.Title) + `</a></li>`)
	}
	b.WriteString(`</ul>`)
	return b.String()
}

func pageHTML(cfg BuildConfig, page PageData) string {
	now := time.Now().UTC().Format(time.RFC3339)
	return `<!doctype html>
<html lang="ja">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="adlaire-generated-at" content="` + now + `">
<title>` + esc(page.Title) + ` - ` + esc(cfg.Title) + `</title>
<link rel="stylesheet" href="` + assetPrefix(page.OutputPath) + `assets/style.css">
</head>
<body>
<div id="progress-bar"></div>
<header id="hdr"><button id="sb-toggle" aria-label="目次">☰</button><a id="brand" href="` + homeHref(page.OutputPath) + `">` + esc(cfg.Title) + `</a><span id="reading-time">約 ` + strconv.Itoa(max(1, page.ReadingTimeMinutes)) + ` 分</span></header>
<div id="lay">
<aside id="sb"><input id="sb-search" type="search" placeholder="Search"><div id="sb-none" hidden>No results</div>` + page.TocHTML + `</aside>
<main id="ct"><article id="main">` + page.HTML + `</article><footer>Generated at ` + now + `</footer></main>
</div>
<button id="btt" aria-label="トップへ戻る">↑</button>
<script type="module" src="` + assetPrefix(page.OutputPath) + `assets/app.js"></script>
</body>
</html>
`
}

func assetPrefix(out string) string {
	if strings.Contains(out, "/") {
		return "../"
	}
	return ""
}

func homeHref(out string) string {
	if strings.Contains(out, "/") {
		return "../index.html"
	}
	return "index.html"
}

func defaultCSS() string {
	return `:root{--adlaire-color-primary:#1455d9;--adlaire-surface:#fff;--adlaire-surface-soft:#f5f7fb;--adlaire-surface-soft-strong:#edf1f7;--adlaire-surface-text:#172033;--adlaire-surface-text-muted:#5d687a;--adlaire-border-default:#d9e0ec;--adlaire-font-family-mono:ui-monospace,SFMono-Regular,Menlo,monospace;--adlaire-font-size-xs:12px;--adlaire-font-size-sm:14px}body{margin:0;background:var(--adlaire-surface);color:var(--adlaire-surface-text);font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif}#hdr{position:fixed;top:0;left:0;right:0;height:56px;display:flex;align-items:center;gap:12px;padding:0 16px;border-bottom:1px solid var(--adlaire-border-default);background:#fff;z-index:10}#brand{font-weight:700;color:inherit;text-decoration:none}#reading-time{margin-left:auto;color:var(--adlaire-surface-text-muted);font-size:var(--adlaire-font-size-sm);white-space:nowrap}#lay{display:flex;padding-top:56px}#sb{position:fixed;top:56px;bottom:0;width:280px;overflow:auto;border-right:1px solid var(--adlaire-border-default);background:var(--adlaire-surface-soft);padding:16px;box-sizing:border-box}#ct{margin-left:312px;max-width:980px;padding:32px;width:100%}.mh{scroll-margin-top:72px}.h1{font-size:32px}.h2{font-size:24px}.h3{font-size:20px}.h4{font-size:18px}.mp{max-width:68ch;line-height:1.75}.mr{border:0;border-top:1px solid var(--adlaire-border-default)}.mbq{border-left:4px solid var(--adlaire-color-primary);margin:1rem 0;padding:.25rem 1rem;background:var(--adlaire-surface-soft)}.ic{font-family:var(--adlaire-font-family-mono);background:var(--adlaire-surface-soft-strong);padding:.1rem .25rem;border-radius:4px}.ml{line-height:1.7}.ml-task{list-style:none}.ml-task input[type="checkbox"]{accent-color:var(--adlaire-color-primary);margin-right:.5rem}.tw{overflow-x:auto}.mt{border-collapse:collapse;width:100%;margin:1rem 0}.mt th,.mt td{border:1px solid var(--adlaire-border-default);padding:.5rem;text-align:left}.mt th[data-sort]{cursor:pointer;user-select:none}.mt th[aria-sort="ascending"]::after{content:" ▲"}.mt th[aria-sort="descending"]::after{content:" ▼"}.cb-wrap{position:relative;margin:1rem 0}.cb-meta{position:absolute;top:8px;right:10px;display:flex;gap:8px}.cl{font-family:var(--adlaire-font-family-mono);font-size:var(--adlaire-font-size-xs);text-transform:uppercase}.cb-copy{opacity:.15}.cb-wrap:hover .cb-copy{opacity:1}.cb{background:var(--adlaire-surface-soft-strong);overflow:auto;padding:2.75rem 1rem 1rem}.fn-section{border-top:1px solid var(--adlaire-border-default);margin-top:2rem}.fn-list{font-size:var(--adlaire-font-size-sm)}.tr{list-style:none;padding:0;margin:0}.ti{margin:.25rem 0}.tl{color:inherit;text-decoration:none}.tl.active{color:var(--adlaire-color-primary);font-weight:700}.hn-link{opacity:0;margin-left:.35rem}.mh:hover .hn-link{opacity:1}#progress-bar{position:fixed;top:0;left:0;height:3px;width:0%;background:var(--adlaire-color-primary);z-index:1000;transition:width .1s linear}.ch-nav{display:flex;justify-content:space-between;padding:1rem 0;margin-top:2rem;border-top:1px solid var(--adlaire-border-default)}.ch-prev,.ch-next{color:var(--adlaire-color-primary);text-decoration:none}#btt{position:fixed;right:20px;bottom:20px;display:none}#btt.visible{display:block}@media(max-width:768px){#sb{transform:translateX(-100%);transition:transform .2s}#sb.open{transform:translateX(0)}#ct{margin-left:0;padding:20px}}@media print{#hdr,#sb,.cb-copy,.hn-link,#btt,.expand-code,#progress-bar,.ch-nav{display:none}#lay,#main{display:block;width:100%;margin:0}pre.cb[data-collapsible]{max-height:none}a[href^="http"]::after,a[href^="https"]::after{content:" (" attr(href) ")"}h2,h3{page-break-before:avoid}}`
}

func defaultJS() string {
	return `document.addEventListener("DOMContentLoaded",()=>{const q=(s,r=document)=>r.querySelector(s);const qa=(s,r=document)=>Array.from(r.querySelectorAll(s));const sb=q("#sb"),search=q("#sb-search"),btt=q("#btt"),bar=q("#progress-bar");q("#sb-toggle")?.addEventListener("click",()=>sb?.classList.toggle("open"));qa(".cb-copy").forEach(btn=>btn.addEventListener("click",()=>{const code=btn.closest(".cb-wrap")?.querySelector("code")?.innerText||"";const done=()=>{btn.textContent="✓ 完了";btn.classList.add("copied");setTimeout(()=>{btn.textContent="コピー";btn.classList.remove("copied")},1800)};navigator.clipboard?.writeText(code).then(done).catch(()=>{const t=document.createElement("textarea");t.value=code;document.body.appendChild(t);t.select();document.execCommand("copy");t.remove();done()})}));qa(".hn-link").forEach(btn=>btn.addEventListener("click",()=>navigator.clipboard?.writeText(location.origin+location.pathname+btn.dataset.href).then(()=>{btn.setAttribute("aria-label","コピーしました");setTimeout(()=>btn.setAttribute("aria-label","リンクをコピー"),1800)})));qa(".expand-code").forEach(btn=>btn.addEventListener("click",()=>{const pre=btn.previousElementSibling;pre.style.maxHeight="none";btn.hidden=true}));qa(".mt th[data-sort]").forEach(th=>th.addEventListener("click",()=>{const table=th.closest("table"),idx=Number(th.dataset.sort),dir=th.getAttribute("aria-sort")==="ascending"?"descending":"ascending";qa("th",table).forEach(h=>h.setAttribute("aria-sort","none"));th.setAttribute("aria-sort",dir);const bodyRows=qa("tbody tr",table),rows=bodyRows.length?bodyRows:qa("tr",table).slice(1);rows.sort((a,b)=>{const av=a.children[idx]?.textContent.trim()||"",bv=b.children[idx]?.textContent.trim()||"",an=Number(av),bn=Number(bv);let c=!Number.isNaN(an)&&!Number.isNaN(bn)?an-bn:av.localeCompare(bv);return dir==="ascending"?c:-c});rows.forEach(r=>table.appendChild(r))}));function onScroll(){if(bar){const d=document.documentElement,total=d.scrollHeight-d.clientHeight;bar.style.width=(total>0?d.scrollTop/total*100:0)+"%"}btt?.classList.toggle("visible",scrollY>400)}document.addEventListener("scroll",onScroll,{passive:true});btt?.addEventListener("click",()=>scrollTo({top:0,behavior:"smooth"}));document.addEventListener("keydown",e=>{const tag=document.activeElement?.tagName;if(e.key==="/"&&!["INPUT","TEXTAREA","SELECT"].includes(tag)){e.preventDefault();search?.focus()}if(e.key==="Escape"&&search){search.value="";search.dispatchEvent(new Event("input"))}if(e.key==="t"&&!["INPUT","TEXTAREA","SELECT"].includes(tag))scrollTo({top:0,behavior:"smooth"})});fetch((location.pathname.includes("/pages/")?"../":"")+"assets/search-index.json").catch(()=>{});search?.addEventListener("input",()=>{const v=search.value.toLowerCase();qa(".tl").forEach(a=>{a.parentElement.hidden=v&&!(a.textContent||"").toLowerCase().includes(v)})});onScroll()});`
}

func writeAtomic(out string, files []siteFile) error {
	info, err := os.Stat(out)
	if err == nil && !info.IsDir() {
		return fmt.Errorf("output path is not directory: %s", out)
	}
	parent := filepath.Dir(out)
	base := filepath.Base(out)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return fmt.Errorf("cannot create output directory: %s", out)
	}
	tmp := filepath.Join(parent, fmt.Sprintf("%s.tmp.%d", base, os.Getpid()))
	prev := filepath.Join(parent, fmt.Sprintf("%s.prev.%d", base, os.Getpid()))
	os.RemoveAll(tmp)
	os.RemoveAll(prev)
	if err := os.MkdirAll(tmp, 0755); err != nil {
		return fmt.Errorf("cannot create output directory: %s", out)
	}
	for _, f := range files {
		path := filepath.Join(tmp, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			os.RemoveAll(tmp)
			return fmt.Errorf("cannot write output: %s", path)
		}
		if err := os.WriteFile(path, f.Data, 0644); err != nil {
			os.RemoveAll(tmp)
			return fmt.Errorf("cannot write output: %s", path)
		}
	}
	if _, err := os.Stat(out); err == nil {
		if err := os.Rename(out, prev); err != nil {
			os.RemoveAll(tmp)
			return fmt.Errorf("cannot write output: %s", out)
		}
	}
	if err := os.Rename(tmp, out); err != nil {
		if _, statErr := os.Stat(prev); statErr == nil {
			if restoreErr := os.Rename(prev, out); restoreErr != nil {
				return fmt.Errorf("cannot restore previous output: %s", out)
			}
		}
		os.RemoveAll(tmp)
		return fmt.Errorf("cannot write output: %s", out)
	}
	os.RemoveAll(prev)
	return nil
}

func outputStats(root string) (int, int64, error) {
	count := 0
	var size int64
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		count++
		size += info.Size()
		return nil
	})
	return count, size, err
}

func htmlPageCount(files []siteFile) int {
	n := 0
	for _, f := range files {
		if strings.HasSuffix(f.Path, ".html") {
			n++
		}
	}
	return n
}

func cleanAbs(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return slashPath(filepath.Clean(p))
	}
	return slashPath(filepath.Clean(abs))
}

func slashPath(p string) string {
	return filepath.ToSlash(filepath.Clean(p))
}

func normalizePlain(s string) string {
	s = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	s = strings.Join(strings.Fields(s), " ")
	return s
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func readingTime(chars int) int {
	if chars <= 0 {
		return 1
	}
	v := chars / 200
	if chars%200 != 0 {
		v++
	}
	if v < 1 {
		return 1
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func stableContent(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(data, []byte("\n"))
	var kept [][]byte
	for _, line := range lines {
		if bytes.Contains(line, []byte(`name="adlaire-generated-at"`)) || bytes.Contains(line, []byte("Generated at ")) {
			continue
		}
		kept = append(kept, line)
	}
	return bytes.Join(kept, []byte("\n")), nil
}
