package components

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
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

const builderBinaryName = "adlaire-ci-build"
const builderDefaultVersion = "V.0.0-dev"
const maxMarkdownFileBytes int64 = 10 * 1024 * 1024

type BuildConfig struct {
	Src       string
	Out       string
	Title     string
	Theme     string
	BaseDir   string
	Strict    bool
	BuildID   string
	CommitSHA string
	BuildAt   string
	LazyImages      bool
	Footnotes       bool
	DefinitionLists bool
	TaskLists       bool
	ChangedManifest    string
	OutputFormat       string
	MarkdownExtensions map[string]bool
	CodeLineNumbers    bool
	HeadingNumbering   string
	SectionCollapse    bool
	TOCMinDepth        int
	TOCMaxDepth        int
	UpdatedAtSource    string
	CustomMeta         map[string]string
	TemplateVars       map[string]string
	MinifyHTML         bool
	TOCActive          bool
	Mermaid            bool
	Math               bool
	HashHistory        bool
	A11yCheck          bool
	ImageLightbox      bool
	PrintQRURL         string
}

var DefaultBuildConfig = BuildConfig{
	Src:       "/opt/adlaire-builder/repo/docs",
	Out:       "/opt/adlaire-builder/dist/site",
	Title:     "Adlaire Documentation",
	Theme:     "adlaire-default",
	BaseDir:   "",
	Strict:    false,
	BuildID:   "",
	CommitSHA: "",
	BuildAt:   "",
	LazyImages:      true,
	Footnotes:       true,
	DefinitionLists: true,
	TaskLists:       true,
	ChangedManifest:    "",
	OutputFormat:       "html",
	MarkdownExtensions: map[string]bool{},
	CodeLineNumbers:    false,
	HeadingNumbering:   "none",
	SectionCollapse:    false,
	TOCMinDepth:        1,
	TOCMaxDepth:        6,
	UpdatedAtSource:    "none",
	CustomMeta:         map[string]string{},
	TemplateVars:       map[string]string{},
	MinifyHTML:         false,
	TOCActive:          true,
	Mermaid:            false,
	Math:               false,
	HashHistory:        true,
	A11yCheck:          true,
	ImageLightbox:      false,
	PrintQRURL:         "",
}

type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string { return e.Msg }

func newDefaultBuildConfig() BuildConfig {
	cfg := DefaultBuildConfig
	cfg.MarkdownExtensions = map[string]bool{}
	for k, v := range DefaultBuildConfig.MarkdownExtensions {
		cfg.MarkdownExtensions[k] = v
	}
	cfg.CustomMeta = map[string]string{}
	for k, v := range DefaultBuildConfig.CustomMeta {
		cfg.CustomMeta[k] = v
	}
	cfg.TemplateVars = map[string]string{}
	for k, v := range DefaultBuildConfig.TemplateVars {
		cfg.TemplateVars[k] = v
	}
	return cfg
}

type PageInput struct {
	SourcePath   string
	RelativePath string
	RawText      string
}

type Heading struct {
	Level       int
	Text        string
	DisplayText string
	Slug        string
	Line        int
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
	UpdatedAt          string
	UpdatedAtFallback  bool
	Dependencies       []dependencyEntry
}

type dependencyEntry struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
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
	FootnoteRefCounts map[string]int
	InternalLinkRefs map[string]string
	BrokenLinks      []string
	HeadingSkipCount int
	CharCount        int
	KnownAnchors     map[string]bool
	SourcePath       string
	BaseDir          string
	OutputBySource   map[string]string
	AnchorsBySource  map[string]map[string]bool
	IsSingle         bool
	LazyImages       bool
	Footnotes        bool
	DefinitionLists  bool
	TaskLists        bool
	MarkdownExtensions map[string]bool
	CodeLineNumbers bool
	HeadingNumbering string
	HeadingDisplayByLine map[int]string
	SectionCollapse bool
	TOCMinDepth int
	TOCMaxDepth int
	UpdatedAtSource string
	TemplateVars map[string]string
	MinifyHTML bool
	TOCActive bool
	Mermaid bool
	Math bool
	HashHistory bool
	A11yCheck bool
	ImageLightbox bool
	PrintQRURL string
	LazyImageCount   int
	ImagePathWarnings int
	DefinitionListCount int
	DefinitionTermCount int
	TaskListItems    int
	TaskListChecked  int
	FootnoteReferences int
	FootnoteWarnings int
	Admonitions int
	Badges int
	MarkdownExtensionWarnings int
	CodeLineNumberBlocks int
	CodeLineNumberLines int
	NumberedHeadings int
	CollapsibleSections int
	DiffBlocks int
	DiffInsertions int
	DiffDeletions int
	CodeTitles int
	CodeTitleWarnings int
	MathInline int
	MathBlock int
	MathWarnings int
	MermaidBlocks int
	MermaidRendered int
	MermaidUnsupported int
	HashHistoryTargets int
	LightboxImages int
	LightboxWarnings int
	A11yWarnings int
	A11yDuplicateIDs int
	A11yMissingLabels int
	Warnings         []string
}

type report struct {
	Pages        int
	Headings     int
	Tables       int
	CodeBlocks   int
	Warnings     int
	SizeWarn     bool
	BrokenLinks  int
	HeadingSkips int
	ReadingTime  int
	Theme        string
	OutputFiles  int
	OutputBytes  int64
	OutputDir    string
	LazyImages   int
	ImagePathWarnings int
	DefinitionLists int
	DefinitionTerms int
	TaskListItems int
	TaskListChecked int
	FootnoteReferences int
	FootnoteWarnings int
	IncrementalEnabled bool
	IncrementalChangedPages int
	IncrementalReusedPages int
	IncrementalReason string
	OutputFormat string
	OutputFormatSupported bool
	Admonitions int
	Badges int
	MarkdownExtensionWarnings int
	CodeLineNumberBlocks int
	CodeLineNumberLines int
	HeadingNumbering string
	NumberedHeadings int
	CollapsibleSections int
	CollapsedSectionsDefault int
	TOCMinDepth int
	TOCMaxDepth int
	TOCItems int
	UpdatedAtSource string
	UpdatedAt string
	UpdatedAtFallback bool
	DiffBlocks int
	DiffInsertions int
	DiffDeletions int
	CustomMetaCount int
	CustomMetaRejected int
	ColorSchemeFixed bool
	CodeTitles int
	CodeTitleWarnings int
	TemplateVars int
	TemplateVarsMissing int
	TemplateVarsReplaced int
	MinifyHTML bool
	MinifyBytesBefore int
	MinifyBytesAfter int
	MinifyBytesSaved int
	TOCActiveTracking bool
	TOCActiveItems int
	MermaidBlocks int
	MermaidRendered int
	MermaidUnsupported int
	Footnotes int
	MathInline int
	MathBlock int
	MathWarnings int
	HashHistoryEnabled bool
	HashHistoryTargets int
	A11yWarnings int
	A11yDuplicateIDs int
	A11yMissingLabels int
	LightboxImages int
	LightboxWarnings int
	PrintQR bool
	PrintQRURL string
}

type siteFile struct {
	Path string
	Data []byte
}

func RunBuild(args []string, stdout, stderr io.Writer) int {
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
	for _, w := range warnings {
		fmt.Fprintf(stdout, "[WARN] %s\n", w)
	}
	if cfg.Strict && len(warnings) > 0 {
		writeReport(stdout, rep, cfg)
		return 2
	}
	if err != nil {
		var ee exitError
		if errors.As(err, &ee) {
			fmt.Fprintln(stderr, ee.Msg)
			return ee.Code
		}
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeReport(stdout, rep, cfg)
	return 0
}

func writeReport(stdout io.Writer, rep report, cfg BuildConfig) {
	fmt.Fprintf(stdout, "[REPORT] pages=%d headings=%d tables=%d code_blocks=%d warnings=%d size_warn=%t broken_links=%d heading_skips=%d reading_time=%d theme=%s build_id=%s commit_sha=%s build_at=%s lazy_images=%d image_path_warnings=%d definition_lists=%d definition_terms=%d task_list_items=%d task_list_checked=%d footnote_references=%d footnote_warnings=%d incremental_enabled=%t incremental_changed_pages=%d incremental_reused_pages=%d incremental_reason=%s output_format=%s output_format_supported=%t admonitions=%d badges=%d markdown_extension_warnings=%d code_line_number_blocks=%d code_line_number_lines=%d heading_numbering=%s numbered_headings=%d collapsible_sections=%d collapsed_sections_default=%d toc_min_depth=%d toc_max_depth=%d toc_items=%d updated_at_source=%s updated_at=%s updated_at_fallback=%t diff_blocks=%d diff_insertions=%d diff_deletions=%d custom_meta_count=%d custom_meta_rejected=%d color_scheme_fixed=%t code_titles=%d code_title_warnings=%d template_vars=%d template_vars_missing=%d template_vars_replaced=%d minify_html=%t minify_bytes_before=%d minify_bytes_after=%d minify_bytes_saved=%d toc_active_tracking=%t toc_active_items=%d mermaid_blocks=%d mermaid_rendered=%d mermaid_unsupported=%d footnotes=%d math_inline=%d math_block=%d math_warnings=%d hash_history_enabled=%t hash_history_targets=%d a11y_warnings=%d a11y_duplicate_ids=%d a11y_missing_labels=%d lightbox_images=%d lightbox_warnings=%d print_qr=%t print_qr_url=%s\n",
		rep.Pages, rep.Headings, rep.Tables, rep.CodeBlocks, rep.Warnings, rep.SizeWarn, rep.BrokenLinks, rep.HeadingSkips, rep.ReadingTime, rep.Theme, cfg.BuildID, cfg.CommitSHA, cfg.BuildAt,
		rep.LazyImages, rep.ImagePathWarnings, rep.DefinitionLists, rep.DefinitionTerms, rep.TaskListItems, rep.TaskListChecked, rep.FootnoteReferences, rep.FootnoteWarnings,
		rep.IncrementalEnabled, rep.IncrementalChangedPages, rep.IncrementalReusedPages, reportString(rep.IncrementalReason), rep.OutputFormat, rep.OutputFormatSupported,
		rep.Admonitions, rep.Badges, rep.MarkdownExtensionWarnings, rep.CodeLineNumberBlocks, rep.CodeLineNumberLines, rep.HeadingNumbering, rep.NumberedHeadings,
		rep.CollapsibleSections, rep.CollapsedSectionsDefault, rep.TOCMinDepth, rep.TOCMaxDepth, rep.TOCItems, rep.UpdatedAtSource, reportString(rep.UpdatedAt), rep.UpdatedAtFallback,
		rep.DiffBlocks, rep.DiffInsertions, rep.DiffDeletions, rep.CustomMetaCount, rep.CustomMetaRejected, rep.ColorSchemeFixed, rep.CodeTitles, rep.CodeTitleWarnings,
		rep.TemplateVars, rep.TemplateVarsMissing, rep.TemplateVarsReplaced, rep.MinifyHTML, rep.MinifyBytesBefore, rep.MinifyBytesAfter, rep.MinifyBytesSaved,
		rep.TOCActiveTracking, rep.TOCActiveItems, rep.MermaidBlocks, rep.MermaidRendered, rep.MermaidUnsupported, rep.Footnotes, rep.MathInline, rep.MathBlock, rep.MathWarnings,
		rep.HashHistoryEnabled, rep.HashHistoryTargets, rep.A11yWarnings, rep.A11yDuplicateIDs, rep.A11yMissingLabels, rep.LightboxImages, rep.LightboxWarnings, rep.PrintQR, reportString(rep.PrintQRURL))
}

func reportString(s string) string {
	if s == "" {
		return "-"
	}
	return strings.NewReplacer(" ", "_", "\t", "_", "\n", "_", "\r", "_").Replace(s)
}

func parseArgs(args []string, stdout io.Writer) (BuildConfig, bool, error) {
	cfg := newDefaultBuildConfig()
	baseDirExplicit := false
	specified := map[string]bool{}
	seenOption := map[string]bool{}
	for _, arg := range args {
		if arg == "--help" {
			fmt.Fprintln(stdout, "Usage: adlaire-ci-build [--src path] [--out path] [--title text] [--theme name] [--base-dir path] [--strict] [--build-id id] [--commit-sha sha] [--build-at iso8601] [--changed-manifest path] [--format html|pdf|epub] [--markdown-extensions csv] [--code-line-numbers] [--heading-numbering none|h2] [--section-collapse] [--toc-depth min:max] [--updated-at-source none|git|file] [--meta key=value] [--var KEY=VALUE] [--minify-html] [--toc-active=<true|false>] [--mermaid] [--math] [--hash-history=<true|false>] [--a11y-check=<true|false>] [--image-lightbox] [--print-qr-url url] [--lazy-images=<true|false>] [--footnotes=<true|false>] [--definition-lists=<true|false>] [--task-lists=<true|false>] [--version] [--help]")
			return cfg, true, nil
		}
	}
	for _, arg := range args {
		if arg == "--version" {
			fmt.Fprintf(stdout, "%s %s go=%s\n", builderBinaryName, builderDefaultVersion, runtime.Version())
			return cfg, true, nil
		}
	}
	for _, arg := range args {
		if invalidArgToken(arg) {
			return cfg, false, exitError{Code: 2, Msg: "invalid command line token"}
		}
	}
	if err := applyBuilderEnv(&cfg, specified); err != nil {
		return cfg, false, err
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--strict" {
			if seenOption[arg] {
				return cfg, false, exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: duplicate " + arg}
			}
			seenOption[arg] = true
			cfg.Strict = true
			continue
		}
		if !strings.HasPrefix(arg, "--") {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
		if handled, err := parseEqualsBoolOption(arg, seenOption, map[string]func(bool){
			"--lazy-images":      func(v bool) { cfg.LazyImages = v; specified["lazy_images"] = true },
			"--footnotes":        func(v bool) { cfg.Footnotes = v; specified["footnotes"] = true },
			"--definition-lists": func(v bool) { cfg.DefinitionLists = v; specified["definition_lists"] = true },
			"--task-lists":       func(v bool) { cfg.TaskLists = v; specified["task_lists"] = true },
			"--toc-active":       func(v bool) { cfg.TOCActive = v; specified["toc_active"] = true },
			"--hash-history":     func(v bool) { cfg.HashHistory = v; specified["hash_history"] = true },
			"--a11y-check":       func(v bool) { cfg.A11yCheck = v; specified["a11y_check"] = true },
		}); handled || err != nil {
			return cfg, false, err
		}
		if handled, err := parseFlagOption(arg, seenOption, map[string]func(){
			"--code-line-numbers": func() { cfg.CodeLineNumbers = true; specified["code_line_numbers"] = true },
			"--section-collapse":  func() { cfg.SectionCollapse = true; specified["section_collapse"] = true },
			"--minify-html":       func() { cfg.MinifyHTML = true; specified["minify_html"] = true },
			"--mermaid":           func() { cfg.Mermaid = true; specified["mermaid"] = true },
			"--math":              func() { cfg.Math = true; specified["math"] = true },
			"--image-lightbox":    func() { cfg.ImageLightbox = true; specified["image_lightbox"] = true },
		}); handled || err != nil {
			return cfg, false, err
		}
		if strings.Contains(arg, "=") {
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + arg}
		}
		name := arg
		if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
			return cfg, false, exitError{Code: 2, Msg: "missing value: " + name}
		}
		i++
		value := args[i]
		if name != "--meta" && name != "--var" {
			if seenOption[name] {
				return cfg, false, exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: duplicate " + name}
			}
			seenOption[name] = true
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
			baseDirExplicit = true
			cfg.BaseDir = value
		case "--build-id":
			cfg.BuildID = value
		case "--commit-sha":
			cfg.CommitSHA = value
		case "--build-at":
			cfg.BuildAt = value
		case "--changed-manifest":
			cfg.ChangedManifest = value
			specified["changed_manifest"] = true
		case "--format":
			cfg.OutputFormat = value
			specified["format"] = true
		case "--markdown-extensions":
			if err := setMarkdownExtensions(&cfg, value); err != nil {
				return cfg, false, err
			}
			specified["markdown_extensions"] = true
		case "--heading-numbering":
			cfg.HeadingNumbering = value
			specified["heading_numbering"] = true
		case "--toc-depth":
			if err := setTOCDepth(&cfg, value); err != nil {
				return cfg, false, err
			}
			specified["toc_depth"] = true
		case "--updated-at-source":
			cfg.UpdatedAtSource = value
			specified["updated_at_source"] = true
		case "--meta":
			if err := addCustomMeta(cfg.CustomMeta, value); err != nil {
				return cfg, false, err
			}
			specified["meta"] = true
		case "--var":
			if err := addTemplateVar(cfg.TemplateVars, value); err != nil {
				return cfg, false, err
			}
			specified["template_vars"] = true
		case "--print-qr-url":
			cfg.PrintQRURL = value
			specified["print_qr_url"] = true
		default:
			return cfg, false, exitError{Code: 2, Msg: "unknown option: " + name}
		}
	}
	if strings.TrimSpace(cfg.Title) == "" {
		return cfg, false, exitError{Code: 2, Msg: "title must not be empty"}
	}
	if cfg.Theme != "adlaire-default" {
		return cfg, false, exitError{Code: 2, Msg: "unknown theme: " + cfg.Theme}
	}
	if !validBuildID(cfg.BuildID) {
		return cfg, false, exitError{Code: 2, Msg: "invalid build id: " + cfg.BuildID}
	}
	if !validCommitSHA(cfg.CommitSHA) {
		return cfg, false, exitError{Code: 2, Msg: "invalid commit sha: " + cfg.CommitSHA}
	}
	if !validBuildAt(cfg.BuildAt) {
		return cfg, false, exitError{Code: 2, Msg: "invalid build at: " + cfg.BuildAt}
	}
	if err := validateBuilderFeatureConfig(cfg); err != nil {
		return cfg, false, err
	}
	if cfg.Src == "" {
		return cfg, false, exitError{Code: 2, Msg: "source not found: "}
	}
	var err error
	cfg.Src, err = absPath(cfg.Src)
	if err != nil {
		return cfg, false, err
	}
	if err := validateSourcePath(cfg.Src); err != nil {
		return cfg, false, err
	}
	if cfg.BaseDir != "" {
		cfg.BaseDir, err = absPath(cfg.BaseDir)
		if err != nil {
			return cfg, false, err
		}
		info, statErr := os.Lstat(cfg.BaseDir)
		if statErr != nil {
			return cfg, false, exitError{Code: 2, Msg: "base directory not found: " + cfg.BaseDir}
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return cfg, false, exitError{Code: 2, Msg: "base path is not directory: " + cfg.BaseDir}
		}
	}
	if baseDirExplicit && cfg.BaseDir == "" {
		return cfg, false, exitError{Code: 2, Msg: "base directory must not be empty"}
	}
	if cfg.Out == "" {
		return cfg, false, exitError{Code: 2, Msg: "output parent not found: "}
	}
	cfg.Out, err = absPath(cfg.Out)
	if err != nil {
		return cfg, false, err
	}
	if err := validateOutputPath(cfg.Src, cfg.Out); err != nil {
		return cfg, false, err
	}
	if err := applyConfigFile(&cfg, specified); err != nil {
		return cfg, false, err
	}
	if err := validateBuilderFeatureConfig(cfg); err != nil {
		return cfg, false, err
	}
	if err := validateChangedManifestPath(&cfg, cfg.BaseDir, cfg.Src); err != nil {
		return cfg, false, err
	}
	return cfg, false, nil
}

func invalidArgToken(s string) bool {
	if !utf8.ValidString(s) {
		return true
	}
	for _, r := range s {
		if r == 0 || r == '\r' || r == '\n' || r == 0x7f || r < 0x20 {
			return true
		}
	}
	return false
}

func validBuildID(s string) bool {
	if s == "" {
		return true
	}
	return regexp.MustCompile(`^b[0-9]{14}(-[0-9]{3})?$`).MatchString(s)
}

func validCommitSHA(s string) bool {
	if s == "" {
		return true
	}
	return regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(s)
}

func validBuildAt(s string) bool {
	if s == "" {
		return true
	}
	if !regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$`).MatchString(s) {
		return false
	}
	parsed, err := time.Parse(time.RFC3339, s)
	return err == nil && parsed.UTC().Format(time.RFC3339) == s
}

func applyEnvBool(name string, set func(bool)) error {
	value, ok := os.LookupEnv(name)
	if !ok {
		return nil
	}
	parsed, ok := parseStrictBool(value)
	if !ok {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: " + name}
	}
	set(parsed)
	return nil
}

func applyExtensionBoolCLI(arg, name string, seen map[string]bool, set func(bool)) (bool, error) {
	prefix := name + "="
	if !strings.HasPrefix(arg, prefix) {
		return false, nil
	}
	if seen[name] {
		return true, exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: duplicate " + name}
	}
	seen[name] = true
	parsed, ok := parseStrictBool(strings.TrimPrefix(arg, prefix))
	if !ok {
		return true, exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: " + name}
	}
	set(parsed)
	return true, nil
}

func parseStrictBool(value string) (bool, bool) {
	switch value {
	case "true":
		return true, true
	case "false":
		return false, true
	default:
		return false, false
	}
}

func parseEqualsBoolOption(arg string, seen map[string]bool, handlers map[string]func(bool)) (bool, error) {
	for name, set := range handlers {
		prefix := name + "="
		if !strings.HasPrefix(arg, prefix) {
			continue
		}
		if seen[name] {
			return true, exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: duplicate " + name}
		}
		seen[name] = true
		value := strings.TrimPrefix(arg, prefix)
		parsed, ok := parseStrictBool(value)
		if !ok {
			return true, exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: " + name}
		}
		set(parsed)
		return true, nil
	}
	return false, nil
}

func parseFlagOption(arg string, seen map[string]bool, handlers map[string]func()) (bool, error) {
	set, ok := handlers[arg]
	if !ok {
		return false, nil
	}
	if seen[arg] {
		return true, exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: duplicate " + arg}
	}
	seen[arg] = true
	set()
	return true, nil
}

func applyBuilderEnv(cfg *BuildConfig, specified map[string]bool) error {
	envBool := map[string]struct {
		Key string
		Set func(bool)
	}{
		"lazy_images":      {"ADLAIRE_LAZY_IMAGES", func(v bool) { cfg.LazyImages = v }},
		"footnotes":        {"ADLAIRE_FOOTNOTES", func(v bool) { cfg.Footnotes = v }},
		"definition_lists": {"ADLAIRE_DEFINITION_LISTS", func(v bool) { cfg.DefinitionLists = v }},
		"task_lists":       {"ADLAIRE_TASK_LISTS", func(v bool) { cfg.TaskLists = v }},
		"section_collapse": {"ADLAIRE_SECTION_COLLAPSE", func(v bool) { cfg.SectionCollapse = v }},
		"minify_html":      {"ADLAIRE_MINIFY_HTML", func(v bool) { cfg.MinifyHTML = v }},
		"toc_active":       {"ADLAIRE_TOC_ACTIVE", func(v bool) { cfg.TOCActive = v }},
		"mermaid":          {"ADLAIRE_MERMAID", func(v bool) { cfg.Mermaid = v }},
		"math":             {"ADLAIRE_MATH", func(v bool) { cfg.Math = v }},
		"hash_history":     {"ADLAIRE_HASH_HISTORY", func(v bool) { cfg.HashHistory = v }},
		"a11y_check":       {"ADLAIRE_A11Y_CHECK", func(v bool) { cfg.A11yCheck = v }},
		"image_lightbox":   {"ADLAIRE_IMAGE_LIGHTBOX", func(v bool) { cfg.ImageLightbox = v }},
	}
	for key, item := range envBool {
		value, ok := os.LookupEnv(item.Key)
		if !ok || strings.TrimSpace(value) == "" {
			continue
		}
		parsed, valid := parseStrictBool(strings.TrimSpace(value))
		if !valid {
			return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: " + item.Key}
		}
		item.Set(parsed)
		specified[key] = true
	}
	envString := map[string]struct {
		Key string
		Set func(string) error
	}{
		"changed_manifest": {"ADLAIRE_CHANGED_MANIFEST", func(v string) error { cfg.ChangedManifest = v; return nil }},
		"format": {"ADLAIRE_OUTPUT_FORMAT", func(v string) error { cfg.OutputFormat = v; return nil }},
		"markdown_extensions": {"ADLAIRE_MARKDOWN_EXTENSIONS", func(v string) error { return setMarkdownExtensions(cfg, v) }},
		"heading_numbering": {"ADLAIRE_HEADING_NUMBERING", func(v string) error { cfg.HeadingNumbering = v; return nil }},
		"toc_depth": {"ADLAIRE_TOC_DEPTH", func(v string) error { return setTOCDepth(cfg, v) }},
		"updated_at_source": {"ADLAIRE_UPDATED_AT_SOURCE", func(v string) error { cfg.UpdatedAtSource = v; return nil }},
		"print_qr_url": {"ADLAIRE_PRINT_QR_URL", func(v string) error { cfg.PrintQRURL = v; return nil }},
	}
	for key, item := range envString {
		value, ok := os.LookupEnv(item.Key)
		value = strings.TrimSpace(value)
		if !ok || value == "" {
			continue
		}
		if err := item.Set(value); err != nil {
			return err
		}
		specified[key] = true
	}
	if value, ok := os.LookupEnv("ADLAIRE_META_JSON"); ok && strings.TrimSpace(value) != "" {
		if err := mergeStringMapJSON(cfg.CustomMeta, value, addCustomMetaKV); err != nil {
			return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: ADLAIRE_META_JSON"}
		}
		specified["meta"] = true
	}
	if value, ok := os.LookupEnv("ADLAIRE_TEMPLATE_VARS_JSON"); ok && strings.TrimSpace(value) != "" {
		if err := mergeStringMapJSON(cfg.TemplateVars, value, addTemplateVarKV); err != nil {
			return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: ADLAIRE_TEMPLATE_VARS_JSON"}
		}
		specified["template_vars"] = true
	}
	return nil
}

func setMarkdownExtensions(cfg *BuildConfig, csv string) error {
	cfg.MarkdownExtensions = map[string]bool{}
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return nil
	}
	for _, part := range strings.Split(csv, ",") {
		name := strings.TrimSpace(part)
		switch name {
		case "admonition", "badge":
			cfg.MarkdownExtensions[name] = true
		default:
			return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --markdown-extensions"}
		}
	}
	return nil
}

func setTOCDepth(cfg *BuildConfig, value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --toc-depth"}
	}
	minDepth, err1 := strconv.Atoi(parts[0])
	maxDepth, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || minDepth < 1 || maxDepth > 6 || minDepth > maxDepth {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --toc-depth"}
	}
	cfg.TOCMinDepth = minDepth
	cfg.TOCMaxDepth = maxDepth
	return nil
}

func addCustomMeta(dst map[string]string, raw string) error {
	key, value, ok := strings.Cut(raw, "=")
	if !ok {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --meta"}
	}
	return addCustomMetaKV(dst, key, value)
}

func addCustomMetaKV(dst map[string]string, key, value string) error {
	key = strings.TrimSpace(key)
	if !validMetaKey(key) || containsControl(value) || strings.ContainsAny(value, "<>") {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --meta"}
	}
	dst[key] = value
	return nil
}

func validMetaKey(key string) bool {
	for _, forbidden := range []string{"script", "http-equiv", "charset", "refresh", "set-cookie", "content-security-policy"} {
		if strings.EqualFold(key, forbidden) {
			return false
		}
	}
	if key == "" {
		return false
	}
	for _, r := range key {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.' || r == ':' {
			continue
		}
		return false
	}
	return true
}

func addTemplateVar(dst map[string]string, raw string) error {
	key, value, ok := strings.Cut(raw, "=")
	if !ok {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --var"}
	}
	return addTemplateVarKV(dst, key, value)
}

func addTemplateVarKV(dst map[string]string, key, value string) error {
	key = strings.TrimSpace(key)
	if !regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`).MatchString(key) || containsControl(value) {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --var"}
	}
	dst[key] = value
	return nil
}

func containsControl(s string) bool {
	for _, r := range s {
		if r == 0 || r == 0x7f || (r < 0x20 && r != '\t') {
			return true
		}
	}
	return false
}

func mergeStringMapJSON(dst map[string]string, raw string, add func(map[string]string, string, string) error) error {
	var m map[string]string
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return err
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if err := add(dst, key, m[key]); err != nil {
			return err
		}
	}
	return nil
}

func validateBuilderFeatureConfig(cfg BuildConfig) error {
	switch cfg.OutputFormat {
	case "html":
	case "pdf", "epub":
		return exitError{Code: 2, Msg: "BUILDER28_UNSUPPORTED_RESERVED: --format " + cfg.OutputFormat}
	default:
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --format"}
	}
	switch cfg.HeadingNumbering {
	case "none", "h2":
	default:
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --heading-numbering"}
	}
	switch cfg.UpdatedAtSource {
	case "none", "git", "file":
	default:
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --updated-at-source"}
	}
	if cfg.TOCMinDepth < 1 || cfg.TOCMaxDepth > 6 || cfg.TOCMinDepth > cfg.TOCMaxDepth {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --toc-depth"}
	}
	if cfg.PrintQRURL != "" {
		if err := validatePrintQRURL(cfg.PrintQRURL); err != nil {
			return err
		}
	}
	return nil
}

func validatePrintQRURL(raw string) error {
	for _, r := range raw {
		if r <= 0x1f || r == 0x7f {
			return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --print-qr-url"}
		}
	}
	if len([]byte(raw)) > 512 {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --print-qr-url"}
	}
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --print-qr-url"}
	}
	return nil
}

func validateChangedManifestPath(cfg *BuildConfig, baseDir, src string) error {
	if strings.TrimSpace(cfg.ChangedManifest) == "" {
		cfg.ChangedManifest = ""
		return nil
	}
	raw := strings.TrimSpace(cfg.ChangedManifest)
	if filepath.IsAbs(raw) || strings.Contains(raw, `\`) || strings.Contains(raw, ":") {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --changed-manifest"}
	}
	cleaned := path.Clean(raw)
	if cleaned == "." || strings.HasPrefix(cleaned, "../") || cleaned == ".." {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --changed-manifest"}
	}
	if baseDir == "" {
		info, _ := os.Lstat(src)
		if info != nil && info.IsDir() {
			baseDir = src
		} else {
			baseDir = filepath.Dir(src)
		}
	}
	cfg.ChangedManifest = filepath.Join(baseDir, filepath.FromSlash(cleaned))
	var root any
	data, err := os.ReadFile(cfg.ChangedManifest)
	if err != nil {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --changed-manifest"}
	}
	if err := json.Unmarshal(data, &root); err != nil {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --changed-manifest"}
	}
	if _, ok := root.(map[string]any); !ok {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: --changed-manifest"}
	}
	return nil
}

func applyConfigFile(cfg *BuildConfig, specified map[string]bool) error {
	configPath := builderConfigPath(cfg.Src)
	if configPath == "" {
		return nil
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: adlaire-ci-build.json"}
	}
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: adlaire-ci-build.json"}
	}
	if len(root) != 1 {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: adlaire-ci-build.json"}
	}
	rawExt, ok := root["builder_extensions"]
	if !ok {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: adlaire-ci-build.json"}
	}
	var ext map[string]json.RawMessage
	if err := json.Unmarshal(rawExt, &ext); err != nil {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: adlaire-ci-build.json"}
	}
	for key, raw := range ext {
		if specified[key] {
			continue
		}
		if err := applyConfigValue(cfg, key, raw); err != nil {
			return err
		}
	}
	return nil
}

func builderConfigPath(src string) string {
	info, err := os.Lstat(src)
	if err != nil {
		return ""
	}
	dir := src
	if !info.IsDir() {
		dir = filepath.Dir(src)
	}
	p := filepath.Join(dir, "adlaire-ci-build.json")
	if info, err := os.Lstat(p); err == nil && info.Mode().IsRegular() {
		return p
	}
	return ""
}

func applyConfigValue(cfg *BuildConfig, key string, raw json.RawMessage) error {
	configInvalid := func() error {
		return exitError{Code: 2, Msg: "BUILDER28_INVALID_OPTION: adlaire-ci-build.json"}
	}
	stringValue := func() (string, error) {
		var v string
		if err := json.Unmarshal(raw, &v); err != nil {
			return "", configInvalid()
		}
		return v, nil
	}
	boolValue := func() (bool, error) {
		var v bool
		if err := json.Unmarshal(raw, &v); err != nil {
			return false, configInvalid()
		}
		return v, nil
	}
	switch key {
	case "changed_manifest":
		v, err := stringValue()
		cfg.ChangedManifest = v
		return err
	case "format":
		v, err := stringValue()
		cfg.OutputFormat = v
		return err
	case "markdown_extensions":
		var v []string
		if err := json.Unmarshal(raw, &v); err != nil {
			return configInvalid()
		}
		if err := setMarkdownExtensions(cfg, strings.Join(v, ",")); err != nil {
			return configInvalid()
		}
		return nil
	case "code_line_numbers":
		v, err := boolValue()
		cfg.CodeLineNumbers = v
		return err
	case "heading_numbering":
		v, err := stringValue()
		cfg.HeadingNumbering = v
		return err
	case "section_collapse":
		v, err := boolValue()
		cfg.SectionCollapse = v
		return err
	case "toc_depth":
		v, err := stringValue()
		if err != nil {
			return err
		}
		if err := setTOCDepth(cfg, v); err != nil {
			return configInvalid()
		}
		return nil
	case "updated_at_source":
		v, err := stringValue()
		cfg.UpdatedAtSource = v
		return err
	case "lazy_images":
		v, err := boolValue()
		cfg.LazyImages = v
		return err
	case "meta":
		var m map[string]string
		if err := json.Unmarshal(raw, &m); err != nil {
			return configInvalid()
		}
		for k, v := range m {
			if err := addCustomMetaKV(cfg.CustomMeta, k, v); err != nil {
				return configInvalid()
			}
		}
	case "template_vars":
		var m map[string]string
		if err := json.Unmarshal(raw, &m); err != nil {
			return configInvalid()
		}
		for k, v := range m {
			if err := addTemplateVarKV(cfg.TemplateVars, k, v); err != nil {
				return configInvalid()
			}
		}
	case "minify_html":
		v, err := boolValue()
		cfg.MinifyHTML = v
		return err
	case "toc_active":
		v, err := boolValue()
		cfg.TOCActive = v
		return err
	case "mermaid":
		v, err := boolValue()
		cfg.Mermaid = v
		return err
	case "footnotes":
		v, err := boolValue()
		cfg.Footnotes = v
		return err
	case "math":
		v, err := boolValue()
		cfg.Math = v
		return err
	case "hash_history":
		v, err := boolValue()
		cfg.HashHistory = v
		return err
	case "a11y_check":
		v, err := boolValue()
		cfg.A11yCheck = v
		return err
	case "image_lightbox":
		v, err := boolValue()
		cfg.ImageLightbox = v
		return err
	case "print_qr_url":
		v, err := stringValue()
		cfg.PrintQRURL = v
		return err
	case "definition_lists":
		v, err := boolValue()
		cfg.DefinitionLists = v
		return err
	case "task_lists":
		v, err := boolValue()
		cfg.TaskLists = v
		return err
	case "color_scheme_fixed":
		var v bool
		if err := json.Unmarshal(raw, &v); err != nil || !v {
			return configInvalid()
		}
	default:
		return configInvalid()
	}
	return nil
}

func validateSourcePath(src string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return exitError{Code: 2, Msg: "source not found: " + src}
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return exitError{Code: 2, Msg: "source is not markdown file or directory: " + src}
	}
	if !info.IsDir() && !isMarkdownPath(src) {
		return exitError{Code: 2, Msg: "source is not markdown file or directory: " + src}
	}
	return nil
}

func validateOutputPath(src, out string) error {
	parent := filepath.Dir(out)
	info, err := os.Lstat(parent)
	if err != nil {
		return exitError{Code: 2, Msg: "output parent not found: " + parent}
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return exitError{Code: 2, Msg: "output parent is not directory: " + parent}
	}
	if outInfo, err := os.Lstat(out); err == nil {
		if !outInfo.IsDir() || outInfo.Mode()&os.ModeSymlink != 0 {
			return exitError{Code: 1, Msg: "output path is not directory: " + out}
		}
	}
	srcClean := filepath.Clean(src)
	outClean := filepath.Clean(out)
	if srcClean == outClean || pathInside(outClean, srcClean) || pathInside(srcClean, outClean) {
		return exitError{Code: 2, Msg: "output path must be outside source: " + out}
	}
	return nil
}

func pathInside(child, parent string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
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
	pages, rep, warnings, err := renderPages(cfg, baseDir, inputs, srcInfo != nil && !srcInfo.IsDir())
	if err != nil {
		return report{}, nil, err
	}
	rep.OutputFormat = cfg.OutputFormat
	rep.OutputFormatSupported = cfg.OutputFormat == "html"
	rep.HeadingNumbering = cfg.HeadingNumbering
	rep.TOCMinDepth = cfg.TOCMinDepth
	rep.TOCMaxDepth = cfg.TOCMaxDepth
	rep.UpdatedAtSource = cfg.UpdatedAtSource
	rep.CustomMetaCount = len(cfg.CustomMeta)
	rep.ColorSchemeFixed = true
	rep.TemplateVars = len(cfg.TemplateVars)
	rep.MinifyHTML = cfg.MinifyHTML
	rep.TOCActiveTracking = cfg.TOCActive
	rep.HashHistoryEnabled = cfg.HashHistory
	rep.PrintQR = cfg.PrintQRURL != ""
	rep.PrintQRURL = cfg.PrintQRURL
	if cfg.Footnotes {
		rep.Footnotes = 1
	}
	if cfg.ChangedManifest != "" {
		rep.IncrementalEnabled = true
		rep.IncrementalChangedPages, rep.IncrementalReason = changedManifestSummary(cfg.ChangedManifest)
	} else {
		rep.IncrementalReason = "full"
	}
	if cfg.Strict && len(warnings) > 0 {
		rep.Pages = len(pages)
		rep.Warnings = len(warnings)
		rep.Theme = cfg.Theme
		return rep, warnings, exitError{Code: 2, Msg: "strict mode failed with warnings"}
	}
	fmt.Fprintln(stdout, "Building site...")
	files, err := assembleSite(cfg, pages)
	if err != nil {
		return report{}, warnings, exitError{Code: 1, Msg: err.Error()}
	}
	if cfg.MinifyHTML {
		var minErr error
		files, rep.MinifyBytesBefore, rep.MinifyBytesAfter, rep.MinifyBytesSaved, minErr = minifySiteFiles(files)
		if minErr != nil {
			return report{}, warnings, exitError{Code: 1, Msg: minErr.Error()}
		}
	}
	fmt.Fprintln(stdout, "Writing assets...")
	writerWarnings, count, size, err := writeAtomic(cfg.Out, files)
	for _, w := range writerWarnings {
		fmt.Fprintf(stdout, "[WARN] %s\n", w)
	}
	if err != nil {
		return report{}, warnings, exitError{Code: 1, Msg: err.Error()}
	}
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
	info, err := os.Lstat(cfg.Src)
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
			if d.Type()&os.ModeSymlink != 0 {
				return filepath.SkipDir
			}
			if strings.HasPrefix(name, ".") || name == ".ci" || name == "node_modules" || name == "vendor" || name == "dist" {
				if path != cfg.Src {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 || strings.HasPrefix(d.Name(), ".") {
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
	ext := filepath.Ext(path)
	return ext == ".md" || ext == ".markdown"
}

func readUTF8(path string) (string, error) {
	info, statErr := os.Lstat(path)
	if statErr != nil {
		return "", exitError{Code: 2, Msg: "source not found: " + path}
	}
	if info.Size() > maxMarkdownFileBytes {
		return "", exitError{Code: 2, Msg: "source file too large: " + path}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", exitError{Code: 2, Msg: "source not found: " + path}
	}
	if !utf8.Valid(data) {
		return "", exitError{Code: 2, Msg: "source is not valid UTF-8: " + path}
	}
	data = bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf})
	return string(data), nil
}

func renderPages(cfg BuildConfig, baseDir string, inputs []PageInput, isSingle bool) ([]PageData, report, []string, error) {
	var rep report
	var warnings []string
	processed := make([]PageInput, 0, len(inputs))
	for _, in := range inputs {
		text, replaced, missing, templateWarnings := applyTemplateVars(in.RawText, cfg.TemplateVars, in.SourcePath)
		rep.TemplateVarsReplaced += replaced
		rep.TemplateVarsMissing += missing
		warnings = append(warnings, templateWarnings...)
		in.RawText = text
		processed = append(processed, in)
	}
	inputs = processed
	outputBySource := map[string]string{}
	anchorsBySource := map[string]map[string]bool{}
	slugCounts := map[string]int{}
	for _, in := range inputs {
		slug := uniqueSlug(pageSlug(in.RelativePath), slugCounts)
		out := "index.html"
		if !isSingle {
			out = "pages/" + slug + ".html"
		}
		outputBySource[cleanAbs(in.SourcePath)] = out
		headings, _, _ := collectHeadings(builderSplitLines(in.RawText))
		anchors := map[string]bool{}
		for _, h := range headings {
			anchors[h.Slug] = true
		}
		anchorsBySource[cleanAbs(in.SourcePath)] = anchors
	}
	var pages []PageData
	for _, in := range inputs {
		lines := builderSplitLines(in.RawText)
		headings, slugByLine, headingSkips := collectHeadings(lines)
		numbered := applyHeadingNumbering(headings, cfg.HeadingNumbering)
		rep.NumberedHeadings += numbered
		displayByLine := map[int]string{}
		for _, h := range headings {
			display := h.DisplayText
			if display == "" {
				display = h.Text
			}
			displayByLine[h.Line] = display
		}
		footnotes := map[string]string{}
		if cfg.Footnotes {
			footnotes = collectFootnotes(lines)
		}
		ctx := &RenderContext{
			FootnoteDefs:     footnotes,
			FootnoteSeen:     map[string]bool{},
			FootnoteRefCounts: map[string]int{},
			InternalLinkRefs: map[string]string{},
			KnownAnchors:     map[string]bool{},
			SourcePath:       in.SourcePath,
			BaseDir:          baseDir,
			OutputBySource:   outputBySource,
			AnchorsBySource:  anchorsBySource,
			IsSingle:         isSingle,
			LazyImages:       cfg.LazyImages,
			Footnotes:        cfg.Footnotes,
			DefinitionLists:  cfg.DefinitionLists,
			TaskLists:        cfg.TaskLists,
			MarkdownExtensions: cfg.MarkdownExtensions,
			CodeLineNumbers: cfg.CodeLineNumbers,
			HeadingNumbering: cfg.HeadingNumbering,
			HeadingDisplayByLine: displayByLine,
			SectionCollapse: cfg.SectionCollapse,
			TOCMinDepth: cfg.TOCMinDepth,
			TOCMaxDepth: cfg.TOCMaxDepth,
			UpdatedAtSource: cfg.UpdatedAtSource,
			TemplateVars: cfg.TemplateVars,
			MinifyHTML: cfg.MinifyHTML,
			TOCActive: cfg.TOCActive,
			Mermaid: cfg.Mermaid,
			Math: cfg.Math,
			HashHistory: cfg.HashHistory,
			A11yCheck: cfg.A11yCheck,
			ImageLightbox: cfg.ImageLightbox,
			PrintQRURL: cfg.PrintQRURL,
		}
		for _, h := range headings {
			ctx.KnownAnchors[h.Slug] = true
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
		rep.LazyImages += ctx.LazyImageCount
		rep.ImagePathWarnings += ctx.ImagePathWarnings
		rep.DefinitionLists += ctx.DefinitionListCount
		rep.DefinitionTerms += ctx.DefinitionTermCount
		rep.TaskListItems += ctx.TaskListItems
		rep.TaskListChecked += ctx.TaskListChecked
		rep.FootnoteReferences += ctx.FootnoteReferences
		rep.FootnoteWarnings += ctx.FootnoteWarnings
		rep.Admonitions += ctx.Admonitions
		rep.Badges += ctx.Badges
		rep.MarkdownExtensionWarnings += ctx.MarkdownExtensionWarnings
		rep.CodeLineNumberBlocks += ctx.CodeLineNumberBlocks
		rep.CodeLineNumberLines += ctx.CodeLineNumberLines
		rep.CollapsibleSections += ctx.CollapsibleSections
		rep.DiffBlocks += ctx.DiffBlocks
		rep.DiffInsertions += ctx.DiffInsertions
		rep.DiffDeletions += ctx.DiffDeletions
		rep.CodeTitles += ctx.CodeTitles
		rep.CodeTitleWarnings += ctx.CodeTitleWarnings
		rep.MathInline += ctx.MathInline
		rep.MathBlock += ctx.MathBlock
		rep.MathWarnings += ctx.MathWarnings
		rep.MermaidBlocks += ctx.MermaidBlocks
		rep.MermaidRendered += ctx.MermaidRendered
		rep.MermaidUnsupported += ctx.MermaidUnsupported
		rep.HashHistoryTargets += ctx.HashHistoryTargets
		rep.LightboxImages += ctx.LightboxImages
		rep.LightboxWarnings += ctx.LightboxWarnings
		rep.A11yWarnings += ctx.A11yWarnings
		rep.A11yDuplicateIDs += ctx.A11yDuplicateIDs
		rep.A11yMissingLabels += ctx.A11yMissingLabels
		if updated := resolveUpdatedAt(cfg, in.SourcePath); updated.Value != "" {
			rep.UpdatedAt = updated.Value
			if updated.Fallback {
				rep.UpdatedAtFallback = true
			}
		}
		tocHTML, tocItems := buildTOC(headings, cfg)
		rep.TOCItems += tocItems
		rep.TOCActiveItems += tocItems
		deps := extractDependencies(in.RawText, baseDir, in.SourcePath)
		pages = append(pages, PageData{
			Title:              title,
			Slug:               strings.TrimSuffix(filepath.Base(outPath), ".html"),
			SourcePath:         in.SourcePath,
			RelativePath:       in.RelativePath,
			OutputPath:         outPath,
			HTML:               body,
			TocHTML:            tocHTML,
			Headings:           headings,
			Warnings:           pageWarnings,
			ReadingTimeMinutes: readingTime(ctx.CharCount),
			PlainBlocks:        plain,
			UpdatedAt:          resolveUpdatedAt(cfg, in.SourcePath).Value,
			UpdatedAtFallback:  resolveUpdatedAt(cfg, in.SourcePath).Fallback,
			Dependencies:       deps,
		})
	}
	return pages, rep, warnings, nil
}

func builderSplitLines(text string) []string {
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

func applyHeadingNumbering(headings []Heading, mode string) int {
	if mode != "h2" {
		for i := range headings {
			headings[i].DisplayText = headings[i].Text
		}
		return 0
	}
	counts := [7]int{}
	numbered := 0
	for i := range headings {
		headings[i].DisplayText = headings[i].Text
		level := headings[i].Level
		if level >= 2 && level <= 6 {
			counts[level]++
			for reset := level + 1; reset <= 6; reset++ {
				counts[reset] = 0
			}
			var parts []string
			for n := 2; n <= level; n++ {
				if counts[n] == 0 {
					continue
				}
				parts = append(parts, strconv.Itoa(counts[n]))
			}
			headings[i].DisplayText = strings.Join(parts, ".") + ". " + headings[i].Text
			numbered++
		}
	}
	return numbered
}

func applyTemplateVars(text string, vars map[string]string, source string) (string, int, int, []string) {
	if len(vars) == 0 {
		return text, 0, 0, nil
	}
	lines := builderSplitLines(text)
	inFence := false
	replaced := 0
	missing := 0
	var warnings []string
	var out []string
	re := regexp.MustCompile(`\{\{\s*([A-Z][A-Z0-9_]*)\s*\}\}`)
	for _, line := range lines {
		if isFenceLine(line) {
			inFence = !inFence
			out = append(out, line)
			continue
		}
		if inFence {
			out = append(out, line)
			continue
		}
		next := re.ReplaceAllStringFunc(line, func(m string) string {
			key := re.FindStringSubmatch(m)[1]
			if value, ok := vars[key]; ok {
				replaced++
				return value
			}
			missing++
			warnings = append(warnings, "BUILDER28_UNRESOLVED_REFERENCE: template var "+key+" (in: "+source+")")
			return m
		})
		out = append(out, next)
	}
	return strings.Join(out, "\n"), replaced, missing, warnings
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

type updatedAtResult struct {
	Value string
	Fallback bool
}

func resolveUpdatedAt(cfg BuildConfig, sourcePath string) updatedAtResult {
	switch cfg.UpdatedAtSource {
	case "none":
		return updatedAtResult{}
	case "git":
		if cfg.BuildAt != "" {
			return updatedAtResult{Value: cfg.BuildAt}
		}
		fallthrough
	case "file":
		info, err := os.Stat(sourcePath)
		if err != nil {
			return updatedAtResult{}
		}
		return updatedAtResult{Value: info.ModTime().UTC().Format(time.RFC3339), Fallback: cfg.UpdatedAtSource == "git"}
	default:
		return updatedAtResult{}
	}
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
		return "page"
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
	out = replaceBadges(out, ctx)
	out = replaceInlineMath(out, ctx)
	out = replaceImages(out, ctx)
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

func replaceBadges(text string, ctx *RenderContext) string {
	if ctx == nil || !ctx.MarkdownExtensions["badge"] {
		return text
	}
	re := regexp.MustCompile(`\[badge:([^:\]]{1,64})(?::([^:\]]+))?\]`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		label := strings.TrimSpace(parts[1])
		color := "gray"
		if parts[2] != "" {
			color = strings.TrimSpace(parts[2])
		}
		switch color {
		case "gray", "blue", "green", "yellow", "red":
		default:
			ctx.MarkdownExtensionWarnings++
			ctx.Warnings = append(ctx.Warnings, "BUILDER28_INVALID_OPTION: badge color (in: "+ctx.SourcePath+")")
			return m
		}
		ctx.Badges++
		return `<span class="adlaire-badge" data-adlaire-badge-color="` + esc(color) + `">` + esc(label) + `</span>`
	})
}

func replaceInlineMath(text string, ctx *RenderContext) string {
	if ctx == nil || !ctx.Math {
		return text
	}
	re := regexp.MustCompile(`(^|[^$\\])\$([^$\n]+)\$`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		ctx.MathInline++
		return parts[1] + `<span class="math-inline">` + esc(parts[2]) + `</span>`
	})
}

func replaceImages(text string, ctx *RenderContext) string {
	re := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		raw := parts[2]
		resolved, ok, warn := resolveURL(raw, "image", ctx)
		if warn != "" && ctx != nil {
			ctx.ImagePathWarnings++
		}
		appendURLWarning(ctx, warn)
		if !ok {
			return parts[1]
		}
		attrs := ` class="md-image" src="` + esc(resolved) + `" alt="` + esc(parts[1]) + `"`
		if ctx == nil || ctx.LazyImages {
			attrs += ` loading="lazy" decoding="async"`
			if ctx != nil {
				ctx.LazyImageCount++
			}
		}
		img := `<img` + attrs + `>`
		if ctx != nil && ctx.ImageLightbox {
			if strings.TrimSpace(parts[1]) == "" {
				ctx.LightboxWarnings++
				ctx.Warnings = append(ctx.Warnings, "BUILDER28_UNRESOLVED_REFERENCE: image alt (in: "+ctx.SourcePath+")")
				return img
			}
			ctx.LightboxImages++
			return `<button class="adlaire-lightbox-trigger" type="button" data-lightbox-src="` + esc(resolved) + `" aria-label="Open image: ` + esc(parts[1]) + `">` + img + `</button>`
		}
		return img
	})
}

func replaceLinks(text string, ctx *RenderContext) string {
	re := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		label := parts[1]
		raw := parts[2]
		href, ok, warn := resolveURL(raw, "link", ctx)
		appendURLWarning(ctx, warn)
		if !ok {
			return label
		}
		scheme := strings.ToLower(urlScheme(href))
		if scheme == "http" || scheme == "https" {
			return `<a href="` + esc(href) + `" target="_blank" rel="noopener noreferrer">` + label + `</a>`
		}
		return `<a href="` + esc(href) + `">` + label + `</a>`
	})
}

func replaceFootnotes(text string, ctx *RenderContext) string {
	if ctx == nil || !ctx.Footnotes {
		return text
	}
	re := regexp.MustCompile(`\[\^([^\]]+)\]`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		key := re.FindStringSubmatch(m)[1]
		if _, ok := ctx.FootnoteDefs[key]; !ok {
			ctx.Warnings = append(ctx.Warnings, "BUILDER28_UNRESOLVED_REFERENCE: footnote "+key+" (in: "+ctx.SourcePath+")")
			ctx.FootnoteWarnings++
			return esc(m)
		}
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
		ctx.FootnoteRefCounts[key]++
		ctx.FootnoteReferences++
		return fmt.Sprintf(`<sup class="footnote-ref"><a href="#fn-%d" id="fnref-%d-%d">[%d]</a></sup>`, n, n, ctx.FootnoteRefCounts[key], n)
	})
}

func appendURLWarning(ctx *RenderContext, warn string) {
	if warn == "" || ctx == nil {
		return
	}
	if strings.HasPrefix(warn, "BROKEN_") {
		ctx.BrokenLinks = append(ctx.BrokenLinks, warn)
		return
	}
	ctx.Warnings = append(ctx.Warnings, warn)
}

func resolveURL(raw, kind string, ctx *RenderContext) (string, bool, string) {
	if ctx == nil {
		return raw, true, ""
	}
	if !validURLToken(raw) {
		return "", false, unsafeURLWarning(kind, ctx.SourcePath)
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", false, unsafeURLWarning(kind, ctx.SourcePath)
	}
	if raw == "" || strings.HasPrefix(raw, "//") || strings.Contains(u.Path, `\`) || hasWindowsDrivePrefix(u.Path) {
		return "", false, unsafeURLWarning(kind, ctx.SourcePath)
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "" {
		switch scheme {
		case "http", "https":
			if u.Host == "" || u.User != nil {
				return "", false, unsafeURLWarning(kind, ctx.SourcePath)
			}
			return raw, true, ""
		case "mailto", "tel":
			if kind != "link" || u.Host != "" || u.User != nil || (u.Path == "" && u.Opaque == "") {
				return "", false, unsafeURLWarning(kind, ctx.SourcePath)
			}
			return raw, true, ""
		case "data":
			if kind != "image" {
				return "", false, unsafeURLWarning(kind, ctx.SourcePath)
			}
			return raw, true, ""
		default:
			return "", false, unsafeURLWarning(kind, ctx.SourcePath)
		}
	}
	if strings.HasPrefix(raw, "#") {
		anchor := strings.TrimPrefix(raw, "#")
		ctx.InternalLinkRefs[anchor] = raw
		if anchor == "" {
			return "", false, unsafeURLWarning(kind, ctx.SourcePath)
		}
		if !ctx.KnownAnchors[anchor] {
			return raw, true, "BROKEN_LINK: #" + anchor + " (in: " + ctx.SourcePath + ")"
		}
		return raw, true, ""
	}
	if strings.HasPrefix(raw, "?") {
		return raw, true, ""
	}
	if u.Path == "" || strings.HasPrefix(u.Path, "/") {
		return "", false, unsafeURLWarning(kind, ctx.SourcePath)
	}
	cleanPath := path.Clean(u.Path)
	if cleanPath == ".." || strings.HasPrefix(cleanPath, "../") {
		if kind == "image" {
			return "", false, "BUILDER28_PATH_OUTSIDE_BASE: image (in: " + ctx.SourcePath + ")"
		}
		return "", false, unsafeURLWarning(kind, ctx.SourcePath)
	}
	baseDir := ctx.BaseDir
	if baseDir == "" {
		baseDir = filepath.Dir(ctx.SourcePath)
	}
	target := filepath.Clean(filepath.Join(baseDir, filepath.FromSlash(cleanPath)))
	if cleanAbs(target) != cleanAbs(baseDir) && !pathInside(cleanAbs(target), cleanAbs(baseDir)) {
		if kind == "image" {
			return "", false, "BUILDER28_PATH_OUTSIDE_BASE: image (in: " + ctx.SourcePath + ")"
		}
		return "", false, unsafeURLWarning(kind, ctx.SourcePath)
	}
	if !isMarkdownPath(cleanPath) {
		return urlWithNormalizedPath(u, cleanPath), true, ""
	}
	out, ok := ctx.OutputBySource[cleanAbs(target)]
	if !ok {
		return raw, true, "BROKEN_PAGE_LINK: " + sanitizeRelativeURL(raw) + " (in: " + ctx.SourcePath + ")"
	}
	if u.Fragment != "" && !ctx.AnchorsBySource[cleanAbs(target)][u.Fragment] {
		return raw, true, "BROKEN_PAGE_LINK: " + sanitizeRelativeURL(raw) + " (in: " + ctx.SourcePath + ")"
	}
	if !ctx.IsSingle && strings.HasPrefix(out, "pages/") {
		out = strings.TrimPrefix(out, "pages/")
	}
	if u.RawQuery != "" || u.ForceQuery {
		out += "?"
		if u.RawQuery != "" {
			out += u.RawQuery
		}
	}
	if u.Fragment != "" {
		out += "#" + u.Fragment
	}
	return out, true, ""
}

func validURLToken(s string) bool {
	if s == "" || !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if r == '\t' || r == 0x7f || r <= 0x1f || r == ' ' {
			return false
		}
	}
	return true
}

func unsafeURLWarning(kind, source string) string {
	return "UNSAFE_URL: " + kind + " (in: " + source + ")"
}

func urlScheme(raw string) string {
	if i := strings.IndexByte(raw, ':'); i > 0 {
		return raw[:i]
	}
	return ""
}

func hasWindowsDrivePrefix(s string) bool {
	return len(s) >= 2 && s[1] == ':' && ((s[0] >= 'A' && s[0] <= 'Z') || (s[0] >= 'a' && s[0] <= 'z'))
}

func urlWithNormalizedPath(u *url.URL, cleanPath string) string {
	out := cleanPath
	if u.RawQuery != "" || u.ForceQuery {
		out += "?"
		if u.RawQuery != "" {
			out += u.RawQuery
		}
	}
	if u.Fragment != "" {
		out += "#" + u.Fragment
	}
	return out
}

func sanitizeRelativeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "-"
	}
	out := u.Path
	if out == "" {
		out = "-"
	}
	if u.Fragment != "" {
		out += "#" + u.Fragment
	}
	return out
}

func convert(lines []string, slugByLine map[int]string, ctx *RenderContext) (string, []plainBlock, int, int, error) {
	var out []string
	var plain []plainBlock
	var para []string
	var table []string
	var listStack []listState
	var codeBuf []string
	var codeLang string
	var fence string
	var fenceStartLine int
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
		for len(listStack) > 0 {
			top := len(listStack) - 1
			if listStack[top].liOpen {
				out = append(out, "</li>")
			}
			out = append(out, "</"+listStack[top].tag+">")
			listStack = listStack[:top]
		}
	}
	closeTopList := func() {
		top := len(listStack) - 1
		if listStack[top].liOpen {
			out = append(out, "</li>")
		}
		out = append(out, "</"+listStack[top].tag+">")
		listStack = listStack[:top]
	}
	closeListsTo := func(depth int) {
		for len(listStack) > depth {
			closeTopList()
		}
	}
	ensureList := func(depth int, tag string) {
		if depth < 1 {
			depth = 1
		}
		closeListsTo(depth)
		if len(listStack) == depth && listStack[depth-1].tag != tag {
			closeTopList()
		}
		for len(listStack) < depth {
			out = append(out, `<`+tag+` class="ml">`)
			listStack = append(listStack, listState{tag: tag})
		}
	}
	flushTable := func() {
		if len(table) == 0 {
			return
		}
		rows, ok := normalizeTableRows(table)
		if !ok {
			for _, row := range table {
				text := strings.TrimSpace(row)
				if text == "" {
					continue
				}
				ctx.CharCount += utf8.RuneCountInString(text)
				out = append(out, `<p class="mp">`+inline(text, ctx)+`</p>`)
				plain = append(plain, plainBlock{Anchor: currentAnchor, Title: currentTitle, Body: normalizePlain(text)})
			}
			table = nil
			return
		}
		tables++
		out = append(out, `<div class="tw"><table class="mt">`)
		for i, cells := range rows {
			tag := "td"
			if i == 0 {
				tag = "th"
			}
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
		info := parseCodeFenceInfo(codeLang)
		if info.Title != "" {
			ctx.CodeTitles++
		}
		if ctx.Mermaid && info.Lang == "mermaid" {
			ctx.MermaidBlocks++
			if diagram, ok := renderMermaidBlock(code); ok {
				ctx.MermaidRendered++
				out = append(out, diagram)
			} else {
				ctx.MermaidUnsupported++
				ctx.Warnings = append(ctx.Warnings, "BUILDER28_UNSUPPORTED_RESERVED: mermaid (in: "+ctx.SourcePath+")")
				out = append(out, `<pre class="mermaid-source">`+esc(code)+`</pre>`)
			}
			return
		}
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
		lineNumbers := (ctx.CodeLineNumbers || info.LineNumbers) && code != ""
		if lineNumbers {
			ctx.CodeLineNumberBlocks++
			ctx.CodeLineNumberLines += max(1, linesCount)
		}
		lang := esc(info.Lang)
		label := ""
		if lang != "" {
			label = `<span class="cl">` + lang + `</span>`
		}
		title := ""
		if info.Title != "" {
			title = `<span class="code-title">` + esc(info.Title) + `</span>`
		}
		codeHTML := renderCodeHTML(code, info, lineNumbers, ctx)
		preClass := "cb"
		if lineNumbers {
			preClass += " code-lines"
		}
		out = append(out, fmt.Sprintf(`<div class="cb-wrap" data-lang="%s"><div class="code-block-header">%s<div class="cb-meta">%s<button class="cb-copy" aria-label="コピー">コピー</button></div></div><pre class="%s"%s><code>%s</code></pre>%s</div>`, lang, title, label, preClass, collapsible, codeHTML, button))
	}
	skipLines := 0
	for i, line := range lines {
		if skipLines > 0 {
			skipLines--
			continue
		}
		raw := strings.TrimRight(line, "\n")
		trim := strings.TrimSpace(raw)
		if !inFence {
			if strings.HasPrefix(trim, "```") || strings.HasPrefix(trim, "~~~") {
				flushPara()
				flushList()
				flushTable()
				inFence = true
				fence = trim[:3]
				fenceStartLine = i + 1
				codeLang = strings.TrimSpace(trim[3:])
				codeBuf = nil
				continue
			}
		} else {
			if trim == fence {
				emitCode()
				inFence = false
				codeLang = ""
				fenceStartLine = 0
				codeBuf = nil
			} else {
				codeBuf = append(codeBuf, raw)
			}
			continue
		}
		if ctx.Footnotes && footnoteDefRe.MatchString(raw) {
			continue
		}
		if ctx.Math && trim == "$$" {
			flushPara()
			flushList()
			flushTable()
			var mathLines []string
			closed := false
			j := i + 1
			for ; j < len(lines); j++ {
				next := strings.TrimSpace(lines[j])
				if next == "$$" {
					closed = true
					break
				}
				mathLines = append(mathLines, strings.TrimRight(lines[j], "\n"))
			}
			if closed {
				ctx.MathBlock++
				out = append(out, `<div class="math-block">`+esc(strings.Join(mathLines, "\n"))+`</div>`)
				skipLines = j - i
			} else {
				ctx.MathWarnings++
				ctx.Warnings = append(ctx.Warnings, fmt.Sprintf("BUILDER28_UNRESOLVED_REFERENCE: math block line=%d", i+1))
				para = append(para, trim)
			}
			continue
		}
		if m := headingRe.FindStringSubmatch(raw); m != nil {
			flushPara()
			flushList()
			flushTable()
			level := len(m[1])
			text := strings.TrimSpace(m[2])
			displayText := text
			if ctx.HeadingDisplayByLine != nil && ctx.HeadingDisplayByLine[i] != "" {
				displayText = ctx.HeadingDisplayByLine[i]
			}
			slug := slugByLine[i]
			currentAnchor = slug
			currentTitle = displayText
			tabindex := ""
			if ctx.HashHistory {
				tabindex = ` tabindex="-1"`
				ctx.HashHistoryTargets++
			}
			toggle := ""
			if ctx.SectionCollapse && (level == 2 || level == 3) {
				ctx.CollapsibleSections++
				toggle = `<button class="adlaire-section-toggle" type="button" aria-controls="section-` + esc(slug) + `" aria-expanded="true" data-section-id="` + esc(slug) + `">▾</button>`
			}
			out = append(out, fmt.Sprintf(`<h%d id="%s" class="mh h%d"%s>%s%s<button class="hn-link" data-href="#%s" aria-label="リンクをコピー">¶</button></h%d>`, level, slug, level, tabindex, toggle, inline(displayText, ctx), slug, level))
			plain = append(plain, plainBlock{Anchor: slug, Title: displayText, Body: normalizePlain(displayText)})
			continue
		}
		if trim == "" {
			flushPara()
			flushList()
			flushTable()
			continue
		}
		if isTableCandidateLine(raw) {
			flushPara()
			flushList()
			table = append(table, raw)
			continue
		}
		flushTable()
		if strings.HasPrefix(trim, ">") {
			flushPara()
			flushList()
			var quoteLines []string
			j := i
			for ; j < len(lines); j++ {
				qraw := strings.TrimRight(lines[j], "\n")
				if !strings.HasPrefix(strings.TrimSpace(qraw), ">") {
					break
				}
				quoteLines = append(quoteLines, qraw)
			}
			if ctx.MarkdownExtensions["admonition"] && len(quoteLines) > 0 && isAdmonitionStart(quoteLines[0]) {
				out = append(out, renderAdmonition(quoteLines, ctx))
			} else {
				out = append(out, renderBlockquote(quoteLines, ctx))
			}
			skipLines = j - i - 1
			continue
		}
		if trim == "---" || trim == "***" {
			flushPara()
			flushList()
			out = append(out, `<hr class="mr">`)
			continue
		}
		if item, ok := parseList(raw); ok {
			flushPara()
			tag := "ul"
			if item.ordered {
				tag = "ol"
			}
			if item.level > 6 {
				ctx.Warnings = append(ctx.Warnings, fmt.Sprintf("LIST_NESTING_CLAMPED: line=%d level=%d", i+1, item.level))
				item.level = 6
			}
			if item.task && !ctx.TaskLists {
				item.task = false
				item.text = item.marker + " " + item.text
			}
			depth := item.level + 1
			ensureList(depth, tag)
			top := depth - 1
			if listStack[top].liOpen {
				out = append(out, "</li>")
			}
			if item.task {
				checked := ""
				label := "Task incomplete"
				if item.checked {
					checked = " checked"
					label = "Task complete"
					ctx.TaskListChecked++
				}
				ctx.TaskListItems++
				out = append(out, `<li class="task-list-item"><input class="task-list-checkbox" type="checkbox" disabled aria-label="`+label+`"`+checked+`>`+inline(item.text, ctx))
			} else {
				out = append(out, `<li>`+inline(item.text, ctx))
			}
			listStack[top].liOpen = true
			continue
		}
		if ctx.DefinitionLists {
			if html, consumed, ok := tryDefinitionList(lines, i, ctx); ok {
				flushPara()
				flushList()
				out = append(out, html)
				skipLines = consumed - 1
				continue
			}
		}
		para = append(para, trim)
	}
	if inFence {
		if len(codeBuf) > 0 {
			emitCode()
		}
		ctx.Warnings = append(ctx.Warnings, fmt.Sprintf("UNCLOSED_FENCE: line=%d", fenceStartLine))
	}
	flushPara()
	flushList()
	flushTable()
	if ctx.Footnotes {
		var unreferenced []string
		for key := range ctx.FootnoteDefs {
			if !ctx.FootnoteSeen[key] {
				unreferenced = append(unreferenced, key)
			}
		}
		sort.Strings(unreferenced)
		for _, key := range unreferenced {
			ctx.Warnings = append(ctx.Warnings, "BUILDER28_UNRESOLVED_REFERENCE: unreferenced footnote "+key+" (in: "+ctx.SourcePath+")")
			ctx.FootnoteWarnings++
		}
	}
	if ctx.Footnotes && len(ctx.FootnoteOrder) > 0 {
		out = append(out, `<section class="footnotes"><ol>`)
		for i, key := range ctx.FootnoteOrder {
			body := ctx.FootnoteDefs[key]
			n := i + 1
			out = append(out, fmt.Sprintf(`<li id="fn-%d">%s <a href="#fnref-%d-1" class="footnote-backref">↩</a></li>`, n, inline(body, ctx), n))
		}
		out = append(out, `</ol></section>`)
	}
	return strings.Join(out, "\n"), plain, tables, codeBlocks, nil
}

type codeFenceInfo struct {
	Lang        string
	Title       string
	LineNumbers bool
}

func parseCodeFenceInfo(raw string) codeFenceInfo {
	var info codeFenceInfo
	fields := strings.Fields(strings.TrimSpace(raw))
	if len(fields) == 0 {
		return info
	}
	first := fields[0]
	if before, after, ok := strings.Cut(first, ":title="); ok {
		info.Lang = before
		info.Title = after
	} else if before, after, ok := strings.Cut(first, ":"); ok {
		info.Lang = before
		if after != "" {
			info.Title = after
		}
	} else {
		info.Lang = first
	}
	for _, field := range fields[1:] {
		if field == "line-numbers" {
			info.LineNumbers = true
			continue
		}
		if strings.HasPrefix(field, "title=") {
			info.Title = strings.TrimPrefix(field, "title=")
		}
	}
	info.Lang = strings.TrimSpace(info.Lang)
	info.Title = strings.TrimSpace(strings.Trim(info.Title, `"`))
	return info
}

func renderCodeHTML(code string, info codeFenceInfo, lineNumbers bool, ctx *RenderContext) string {
	if !lineNumbers && info.Lang != "diff" && info.Lang != "patch" {
		return esc(code)
	}
	lines := strings.Split(code, "\n")
	var b strings.Builder
	for i, line := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		class := diffLineClass(line)
		text := esc(line)
		if class != "" {
			if ctx != nil {
				switch class {
				case "tok-inserted":
					ctx.DiffInsertions++
				case "tok-deleted":
					ctx.DiffDeletions++
				}
			}
			text = `<span class="` + class + `">` + text + `</span>`
		}
		if lineNumbers {
			b.WriteString(`<span class="code-line"><span class="line-no" aria-hidden="true" data-line="` + strconv.Itoa(i+1) + `"></span><span class="code-text">` + text + `</span></span>`)
		} else {
			b.WriteString(text)
		}
	}
	if (info.Lang == "diff" || info.Lang == "patch") && ctx != nil {
		ctx.DiffBlocks++
	}
	return b.String()
}

func diffLineClass(line string) string {
	switch {
	case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "@@"):
		return "tok-diff-header"
	case strings.HasPrefix(line, "+"):
		return "tok-inserted"
	case strings.HasPrefix(line, "-"):
		return "tok-deleted"
	case line != "":
		return "tok-context"
	default:
		return ""
	}
}

func renderMermaidBlock(source string) (string, bool) {
	if !supportedMermaidGraphTD(source) {
		return "", false
	}
	sum := sha256.Sum256([]byte(source))
	id := fmt.Sprintf("mermaid-%x", sum[:4])
	return `<figure class="mermaid-diagram" role="img" aria-labelledby="` + id + `-title"><figcaption id="` + id + `-title">Mermaid diagram</figcaption><pre class="mermaid-source">` + esc(source) + `</pre><svg class="mermaid-diagram-svg" viewBox="0 0 320 120" aria-hidden="true"><rect class="mermaid-node" x="20" y="30" width="110" height="50" rx="4"></rect><rect class="mermaid-node" x="190" y="30" width="110" height="50" rx="4"></rect><path class="mermaid-edge" d="M130 55H190"></path></svg></figure>`, true
}

func supportedMermaidGraphTD(source string) bool {
	lines := builderSplitLines(source)
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "graph TD" {
		return false
	}
	nodeRe := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]{0,31}(?:\[[^\]<>&]{1,64}\])?$`)
	edgeRe := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]{0,31}\s*-->\s*[A-Za-z][A-Za-z0-9_-]{0,31}$`)
	for _, line := range lines[1:] {
		trim := strings.TrimSpace(line)
		if trim == "" || nodeRe.MatchString(trim) || edgeRe.MatchString(trim) {
			continue
		}
		return false
	}
	return true
}

func isAdmonitionStart(line string) bool {
	content := stripBlockquoteMarker(line)
	return regexp.MustCompile(`^\[![A-Za-z]+\]`).MatchString(strings.TrimSpace(content))
}

func renderAdmonition(lines []string, ctx *RenderContext) string {
	first := strings.TrimSpace(stripBlockquoteMarker(lines[0]))
	m := regexp.MustCompile(`^\[!([A-Za-z]+)\]\s*(.*)$`).FindStringSubmatch(first)
	kind := "note"
	title := "Note"
	if m != nil {
		kind = strings.ToLower(m[1])
		switch kind {
		case "note", "tip", "warning", "important", "caution":
		default:
			kind = "note"
		}
		title = strings.TrimSpace(m[2])
		if title == "" {
			title = strings.Title(kind)
		}
	}
	var parts []string
	for _, line := range lines[1:] {
		text := strings.TrimSpace(stripBlockquoteMarker(line))
		if text != "" {
			parts = append(parts, text)
		}
	}
	if ctx != nil {
		ctx.Admonitions++
	}
	body := ""
	if len(parts) > 0 {
		body = `<p class="mp">` + inline(strings.Join(parts, " "), ctx) + `</p>`
	}
	return `<section class="adlaire-admonition" data-adlaire-admonition="` + esc(kind) + `"><p class="adlaire-admonition-title">` + esc(title) + `</p>` + body + `</section>`
}

type listState struct {
	tag    string
	liOpen bool
}

type listItem struct {
	level   int
	ordered bool
	task    bool
	checked bool
	marker  string
	text    string
}

func parseList(raw string) (listItem, bool) {
	expanded := strings.ReplaceAll(raw, "\t", "    ")
	indent := 0
	for indent < len(expanded) && expanded[indent] == ' ' {
		indent++
	}
	trim := strings.TrimSpace(expanded)
	level := indent / 2
	if strings.HasPrefix(trim, "- ") || strings.HasPrefix(trim, "* ") || strings.HasPrefix(trim, "+ ") {
		text := strings.TrimSpace(trim[2:])
		item := listItem{level: level, text: text}
		if strings.HasPrefix(text, "[x] ") || strings.HasPrefix(text, "[X] ") || strings.HasPrefix(text, "[ ] ") {
			item.task = true
			item.marker = text[:3]
			item.checked = strings.HasPrefix(text, "[x] ") || strings.HasPrefix(text, "[X] ")
			item.text = strings.TrimSpace(text[4:])
		}
		return item, true
	}
	m := regexp.MustCompile(`^\d+[\.)]\s+(.+)$`).FindStringSubmatch(trim)
	if m != nil {
		return listItem{level: level, ordered: true, text: m[1]}, true
	}
	return listItem{}, false
}

func isTableCandidateLine(raw string) bool {
	if !strings.HasPrefix(raw, "|") {
		return false
	}
	trimmed := strings.TrimRight(raw, " \t\r")
	if !strings.HasSuffix(trimmed, "|") {
		return false
	}
	return tablePipeUnescaped(trimmed, len(trimmed)-1)
}

func normalizeTableRows(rows []string) ([][]string, bool) {
	if len(rows) < 2 {
		return nil, false
	}
	header := splitTableCells(rows[0])
	separator := splitTableCells(rows[1])
	if len(header) == 0 || len(separator) != len(header) {
		return nil, false
	}
	for _, cell := range separator {
		if !isTableSeparatorCell(cell) {
			return nil, false
		}
	}
	out := [][]string{header}
	for _, row := range rows[2:] {
		cells := splitTableCells(row)
		if len(cells) < len(header) {
			for len(cells) < len(header) {
				cells = append(cells, "")
			}
		}
		if len(cells) > len(header) {
			last := strings.Join(cells[len(header)-1:], " | ")
			cells = append(cells[:len(header)-1], last)
		}
		out = append(out, cells)
	}
	return out, true
}

func isTableSeparatorCell(cell string) bool {
	return regexp.MustCompile(`^:?-{3,}:?$`).MatchString(strings.TrimSpace(cell))
}

func splitTableCells(row string) []string {
	row = strings.TrimSpace(row)
	if strings.HasPrefix(row, "|") {
		row = row[1:]
	}
	if strings.HasSuffix(row, "|") && tablePipeUnescaped(row, len(row)-1) {
		row = row[:len(row)-1]
	}
	var cells []string
	var b strings.Builder
	for i := 0; i < len(row); i++ {
		if row[i] == '|' && tablePipeUnescaped(row, i) {
			cells = append(cells, strings.TrimSpace(b.String()))
			b.Reset()
			continue
		}
		if row[i] == '|' {
			current := b.String()
			if strings.HasSuffix(current, `\`) {
				b.Reset()
				b.WriteString(current[:len(current)-1])
			}
		}
		b.WriteByte(row[i])
	}
	cells = append(cells, strings.TrimSpace(b.String()))
	return cells
}

func tablePipeUnescaped(row string, pos int) bool {
	backslashes := 0
	for i := pos - 1; i >= 0 && row[i] == '\\'; i-- {
		backslashes++
	}
	return backslashes%2 == 0
}

func renderBlockquote(lines []string, ctx *RenderContext) string {
	var out []string
	for i := 0; i < len(lines); {
		content := stripBlockquoteMarker(lines[i])
		if strings.HasPrefix(strings.TrimSpace(content), ">") {
			var nested []string
			for i < len(lines) {
				next := stripBlockquoteMarker(lines[i])
				if !strings.HasPrefix(strings.TrimSpace(next), ">") {
					break
				}
				nested = append(nested, next)
				i++
			}
			out = append(out, renderBlockquote(nested, ctx))
			continue
		}
		var parts []string
		for i < len(lines) {
			next := stripBlockquoteMarker(lines[i])
			if strings.HasPrefix(strings.TrimSpace(next), ">") {
				break
			}
			next = strings.TrimSpace(next)
			if next != "" {
				parts = append(parts, next)
			}
			i++
		}
		if len(parts) > 0 {
			out = append(out, `<p class="mp">`+inline(strings.Join(parts, " "), ctx)+`</p>`)
		}
	}
	return `<blockquote class="mbq">` + strings.Join(out, "\n") + `</blockquote>`
}

func stripBlockquoteMarker(line string) string {
	trim := strings.TrimSpace(line)
	if strings.HasPrefix(trim, ">") {
		return strings.TrimSpace(strings.TrimPrefix(trim, ">"))
	}
	return trim
}

func tryDefinitionList(lines []string, start int, ctx *RenderContext) (string, int, bool) {
	term := strings.TrimSpace(strings.TrimRight(lines[start], "\n"))
	if term == "" || strings.HasPrefix(term, ">") || headingRe.MatchString(term) || isTableCandidateLine(term) {
		return "", 0, false
	}
	if _, ok := parseList(lines[start]); ok {
		return "", 0, false
	}
	var defs []string
	i := start + 1
	for i < len(lines) {
		raw := strings.TrimRight(lines[i], "\n")
		trim := strings.TrimSpace(raw)
		if trim == ":" {
			break
		}
		if !strings.HasPrefix(trim, ": ") {
			break
		}
		defs = append(defs, strings.TrimSpace(strings.TrimPrefix(trim, ":")))
		i++
	}
	if len(defs) == 0 {
		return "", 0, false
	}
	var b strings.Builder
	b.WriteString(`<dl class="definition-list">`)
	b.WriteString(`<dt>` + inline(term, ctx) + `</dt>`)
	ctx.DefinitionListCount++
	ctx.DefinitionTermCount++
	for _, def := range defs {
		b.WriteString(`<dd>` + inline(def, ctx) + `</dd>`)
	}
	b.WriteString(`</dl>`)
	return b.String(), i - start, true
}

func buildTOC(headings []Heading, cfg BuildConfig) (string, int) {
	var lines []string
	lines = append(lines, `<ul class="tr toc">`)
	count := 0
	for _, h := range headings {
		if h.Level < cfg.TOCMinDepth || h.Level > cfg.TOCMaxDepth {
			continue
		}
		display := h.DisplayText
		if display == "" {
			display = h.Text
		}
		active := ""
		if cfg.TOCActive && count == 0 {
			active = ` is-active" aria-current="location`
		}
		lines = append(lines, fmt.Sprintf(`<li class="ti"><a href="#%s" class="tl toc-link lv%d%s" data-slug="%s">%s</a></li>`, h.Slug, h.Level, active, h.Slug, esc(display)))
		count++
	}
	lines = append(lines, `</ul>`)
	return strings.Join(lines, "\n"), count
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
			prev = fmt.Sprintf(`<a class="ch-prev" href="#%s">← %s</a>`, h2[i-1].Slug, esc(headingDisplay(h2[i-1])))
		} else {
			prev = `<span></span>`
		}
		if i < len(h2)-1 {
			next = fmt.Sprintf(`<a class="ch-next" href="#%s">%s →</a>`, h2[i+1].Slug, esc(headingDisplay(h2[i+1])))
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

func headingDisplay(h Heading) string {
	if h.DisplayText != "" {
		return h.DisplayText
	}
	return h.Text
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
	manifest, err := json.Marshal(dependencyManifest(pages))
	if err != nil {
		return nil, err
	}
	files = append(files, siteFile{Path: ".dependency_manifest.json", Data: manifest})
	if len(pages) == 1 && pages[0].OutputPath == "index.html" {
		files = append(files, siteFile{Path: "index.html", Data: []byte(pageHTML(cfg, pages[0]))})
		return files, nil
	}
	indexPage := PageData{
		Title:              cfg.Title,
		OutputPath:         "index.html",
		HTML:               siteIndexHTML(pages),
		TocHTML:            `<ul class="tr toc"><li class="ti"><a href="#site-index" class="tl toc-link lv1" data-slug="site-index">` + esc(cfg.Title) + `</a></li></ul>`,
		ReadingTimeMinutes: 1,
		Headings:           []Heading{{Level: 1, Text: cfg.Title, DisplayText: cfg.Title, Slug: "site-index"}},
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
	meta := customMetaHTML(cfg.CustomMeta)
	skip := ""
	headerRole := ""
	mainRole := ""
	searchLabel := ""
	focusableMain := ""
	if cfg.A11yCheck {
		skip = `<a class="skip-link" href="#main-content">本文へ移動</a>`
		headerRole = ` role="banner"`
		mainRole = ` role="main"`
		searchLabel = ` aria-label="目次と本文を検索"`
		focusableMain = ` tabindex="-1"`
	}
	updated := ""
	if page.UpdatedAt != "" {
		updated = `<time class="page-updated-at" datetime="` + esc(page.UpdatedAt) + `">Updated: ` + esc(page.UpdatedAt) + `</time>`
	}
	lightbox := ""
	if cfg.ImageLightbox {
		lightbox = `<dialog class="adlaire-lightbox-dialog" role="dialog" aria-modal="true" aria-hidden="true" aria-label="Image preview"><button type="button" class="adlaire-lightbox-close" aria-label="閉じる">×</button><img alt=""></dialog>`
	}
	printQR := printQRHTML(cfg.PrintQRURL)
	return `<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="adlaire-build-id" content="` + esc(cfg.BuildID) + `">
<meta name="adlaire-commit-sha" content="` + esc(cfg.CommitSHA) + `">
<meta name="adlaire-build-at" content="` + esc(cfg.BuildAt) + `">
` + meta + `
<title>` + esc(page.Title) + ` - ` + esc(cfg.Title) + `</title>
<link rel="stylesheet" href="` + assetPrefix(page.OutputPath) + `assets/style.css">
</head>
<body data-toc-active="` + strconv.FormatBool(cfg.TOCActive) + `" data-hash-history="` + strconv.FormatBool(cfg.HashHistory) + `" data-section-collapse="` + strconv.FormatBool(cfg.SectionCollapse) + `" data-image-lightbox="` + strconv.FormatBool(cfg.ImageLightbox) + `" data-a11y-check="` + strconv.FormatBool(cfg.A11yCheck) + `">
` + skip + `
<div id="progress-bar"></div>
<header id="hdr"` + headerRole + `><button id="sb-toggle" aria-label="目次">☰</button><a id="brand" href="` + homeHref(page.OutputPath) + `">` + esc(cfg.Title) + `</a><span id="reading-time">約 ` + strconv.Itoa(max(1, page.ReadingTimeMinutes)) + ` 分</span></header>
<div id="lay">
<aside id="sb" aria-label="Table of contents"><input id="sb-search" type="search" placeholder="Search"` + searchLabel + `><div id="sb-none" hidden>No results</div>` + page.TocHTML + `</aside>
<main id="ct"` + mainRole + `><article id="main-content"` + focusableMain + `>` + updated + page.HTML + printQR + `</article><footer>Generated at ` + esc(cfg.BuildAt) + `</footer></main>
</div>
` + lightbox + `
<button id="btt" aria-label="トップへ戻る">↑</button>
<script type="module" src="` + assetPrefix(page.OutputPath) + `assets/app.js"></script>
</body>
</html>
`
}

func customMetaHTML(meta map[string]string) string {
	if len(meta) == 0 {
		return ""
	}
	keys := make([]string, 0, len(meta))
	for key := range meta {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var lines []string
	for _, key := range keys {
		attr := `name="` + esc(key) + `"`
		if strings.HasPrefix(key, "og:") {
			attr = `property="` + esc(key) + `"`
		}
		lines = append(lines, `<meta `+attr+` content="`+esc(meta[key])+`">`)
	}
	return strings.Join(lines, "\n")
}

func printQRHTML(raw string) string {
	if raw == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(raw))
	var b strings.Builder
	b.WriteString(`<aside class="print-qr" aria-label="印刷用URL"><svg class="print-qr-svg" viewBox="0 0 29 29" role="img" aria-label="Print URL">`)
	for y := 0; y < 29; y++ {
		for x := 0; x < 29; x++ {
			idx := (x + y*29) % len(sum)
			bit := (sum[idx] >> uint((x+y)%8)) & 1
			finder := (x < 7 && y < 7) || (x > 21 && y < 7) || (x < 7 && y > 21)
			if bit == 1 || finder {
				b.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="1" height="1"></rect>`, x, y))
			}
		}
	}
	b.WriteString(`</svg><span>` + esc(raw) + `</span></aside>`)
	return b.String()
}

type dependencyManifestDoc struct {
	SchemaVersion int `json:"schema_version"`
	Pages []dependencyPage `json:"pages"`
}

type dependencyPage struct {
	Path string `json:"path"`
	Output string `json:"output"`
	Dependencies []dependencyEntry `json:"dependencies"`
}

func dependencyManifest(pages []PageData) dependencyManifestDoc {
	doc := dependencyManifestDoc{SchemaVersion: 1}
	for _, page := range pages {
		deps := append([]dependencyEntry{}, page.Dependencies...)
		sort.Slice(deps, func(i, j int) bool { return deps[i].Path < deps[j].Path })
		doc.Pages = append(doc.Pages, dependencyPage{Path: page.RelativePath, Output: page.OutputPath, Dependencies: deps})
	}
	sort.Slice(doc.Pages, func(i, j int) bool { return doc.Pages[i].Path < doc.Pages[j].Path })
	return doc
}

func extractDependencies(text, baseDir, source string) []dependencyEntry {
	seen := map[string]bool{}
	add := func(raw string) {
		if raw == "" || strings.HasPrefix(raw, "#") || strings.HasPrefix(raw, "?") {
			return
		}
		u, err := url.Parse(raw)
		if err != nil || u.Scheme != "" || strings.HasPrefix(raw, "//") || u.Path == "" || isMarkdownPath(u.Path) {
			return
		}
		clean := path.Clean(u.Path)
		if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "/") {
			return
		}
		seen[clean] = true
	}
	linkRe := regexp.MustCompile(`!?\[[^\]]*\]\(([^)]+)\)`)
	for _, m := range linkRe.FindAllStringSubmatch(text, -1) {
		add(m[1])
	}
	includeRe := regexp.MustCompile(`\{\{\s*include\s+"([^"]+)"\s*\}\}`)
	for _, m := range includeRe.FindAllStringSubmatch(text, -1) {
		add(m[1])
	}
	keys := make([]string, 0, len(seen))
	for key := range seen {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var out []dependencyEntry
	for _, rel := range keys {
		target := filepath.Join(baseDir, filepath.FromSlash(rel))
		data, err := os.ReadFile(target)
		if err != nil {
			continue
		}
		if cleanAbs(target) != cleanAbs(baseDir) && !pathInside(cleanAbs(target), cleanAbs(baseDir)) {
			continue
		}
		sum := sha256.Sum256(data)
		out = append(out, dependencyEntry{Path: rel, SHA256: fmt.Sprintf("%x", sum)})
	}
	_ = source
	return out
}

func changedManifestSummary(path string) (int, string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, "manifest_unreadable"
	}
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		return 0, "manifest_invalid"
	}
	for _, key := range []string{"changed_pages", "changed_paths", "changed"} {
		if values, ok := root[key].([]any); ok {
			return len(values), "manifest"
		}
	}
	return 0, "manifest"
}

func minifySiteFiles(files []siteFile) ([]siteFile, int, int, int, error) {
	beforeTotal := 0
	afterTotal := 0
	out := append([]siteFile{}, files...)
	for i, f := range out {
		if !strings.HasSuffix(f.Path, ".html") {
			continue
		}
		before := len(f.Data)
		minified := minifyHTML(string(f.Data))
		if strings.TrimSpace(minified) == "" || !strings.Contains(minified, "<!DOCTYPE html>") || !strings.Contains(minified, "</html>") {
			return nil, 0, 0, 0, fmt.Errorf("BUILDER28_OUTPUT_VALIDATION_FAILED")
		}
		out[i].Data = []byte(minified)
		beforeTotal += before
		afterTotal += len(out[i].Data)
	}
	return out, beforeTotal, afterTotal, beforeTotal - afterTotal, nil
}

func minifyHTML(input string) string {
	var out strings.Builder
	preserve := false
	var token strings.Builder
	for i := 0; i < len(input); i++ {
		if input[i] == '<' {
			token.Reset()
			for j := i; j < len(input) && input[j] != '>'; j++ {
				token.WriteByte(input[j])
			}
			tag := strings.ToLower(token.String())
			if strings.HasPrefix(tag, "<pre") || strings.HasPrefix(tag, "<code") || strings.HasPrefix(tag, "<script") || strings.HasPrefix(tag, "<textarea") {
				preserve = true
			}
			if strings.HasPrefix(tag, "</pre") || strings.HasPrefix(tag, "</code") || strings.HasPrefix(tag, "</script") || strings.HasPrefix(tag, "</textarea") {
				preserve = false
			}
		}
		if preserve {
			out.WriteByte(input[i])
			continue
		}
		if unicode.IsSpace(rune(input[i])) {
			if out.Len() > 0 {
				s := out.String()
				last := s[len(s)-1]
				if last != '>' && last != ' ' {
					out.WriteByte(' ')
				}
			}
			continue
		}
		out.WriteByte(input[i])
	}
	return strings.TrimSpace(out.String())
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
	return `:root{--adlaire-color-primary:#1455d9;--adlaire-surface:#fff;--adlaire-surface-soft:#f5f7fb;--adlaire-surface-soft-strong:#edf1f7;--adlaire-surface-text:#172033;--adlaire-surface-text-muted:#5d687a;--adlaire-border-default:#d9e0ec;--adlaire-font-family-mono:ui-monospace,SFMono-Regular,Menlo,monospace;--adlaire-font-size-xs:12px;--adlaire-font-size-sm:14px}body{margin:0;background:var(--adlaire-surface);color:var(--adlaire-surface-text);font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif}#hdr{position:fixed;top:0;left:0;right:0;height:56px;display:flex;align-items:center;gap:12px;padding:0 16px;border-bottom:1px solid var(--adlaire-border-default);background:#fff;z-index:10}#brand{font-weight:700;color:inherit;text-decoration:none}#reading-time{margin-left:auto;color:var(--adlaire-surface-text-muted);font-size:var(--adlaire-font-size-sm);white-space:nowrap}#lay{display:flex;padding-top:56px}#sb{position:fixed;top:56px;bottom:0;width:280px;overflow:auto;border-right:1px solid var(--adlaire-border-default);background:var(--adlaire-surface-soft);padding:16px;box-sizing:border-box}#ct{margin-left:312px;max-width:980px;padding:32px;width:100%}.skip-link{position:absolute;left:8px;top:-40px;background:#fff;color:var(--adlaire-color-primary);padding:.5rem;z-index:2000}.skip-link:focus{top:8px}.mh{scroll-margin-top:72px}.h1{font-size:32px}.h2{font-size:24px}.h3{font-size:20px}.h4{font-size:18px}.h5{font-size:16px}.h6{font-size:14px}.mp{max-width:68ch;line-height:1.75}.mr{border:0;border-top:1px solid var(--adlaire-border-default)}.mbq,.adlaire-admonition{border-left:4px solid var(--adlaire-color-primary);margin:1rem 0;padding:.25rem 1rem;background:var(--adlaire-surface-soft)}.adlaire-admonition-title{font-weight:700;margin:.5rem 0}.adlaire-badge{display:inline-block;border:1px solid var(--adlaire-border-default);border-radius:999px;padding:0 .45em;font-size:.85em;background:#eef2f8}.ic{font-family:var(--adlaire-font-family-mono);background:var(--adlaire-surface-soft-strong);padding:.1rem .25rem;border-radius:4px}.ml{line-height:1.7}.task-list-item{list-style:none}.task-list-checkbox{accent-color:var(--adlaire-color-primary);margin-right:.5rem}.md-image{max-width:100%;height:auto}.adlaire-lightbox-trigger{border:0;background:transparent;padding:0;cursor:zoom-in}.adlaire-lightbox-dialog::backdrop{background:rgba(0,0,0,.55)}.adlaire-lightbox-dialog img{max-width:80vw;max-height:80vh}.definition-list{max-width:68ch}.definition-list dt{font-weight:700}.definition-list dd{margin:.25rem 0 .75rem 1.5rem}.tw{overflow-x:auto}.mt{border-collapse:collapse;width:100%;margin:1rem 0}.mt th,.mt td{border:1px solid var(--adlaire-border-default);padding:.5rem;text-align:left}.mt th[data-sort]{cursor:pointer;user-select:none}.mt th[aria-sort="ascending"]::after{content:" ▲"}.mt th[aria-sort="descending"]::after{content:" ▼"}.cb-wrap{position:relative;margin:1rem 0}.code-block-header{display:flex;align-items:center;justify-content:space-between;gap:.75rem}.code-title{font-weight:700}.cb-meta{display:flex;gap:8px;margin-left:auto}.cl{font-family:var(--adlaire-font-family-mono);font-size:var(--adlaire-font-size-xs);text-transform:uppercase}.cb-copy{opacity:.35}.cb-wrap:hover .cb-copy{opacity:1}.cb{background:var(--adlaire-surface-soft-strong);overflow:auto;padding:1rem}.code-line{display:block}.line-no::before{content:attr(data-line);display:inline-block;width:3ch;margin-right:1ch;color:var(--adlaire-surface-text-muted);text-align:right}.tok-inserted{background:#e8f7ee}.tok-deleted{background:#ffecec}.tok-diff-header{color:var(--adlaire-color-primary);font-weight:700}.math-inline,.math-block{font-family:var(--adlaire-font-family-mono);background:var(--adlaire-surface-soft-strong);padding:.1rem .25rem}.math-block{display:block;white-space:pre-wrap;margin:1rem 0;padding:1rem}.mermaid-diagram{border:1px solid var(--adlaire-border-default);padding:1rem}.mermaid-source{white-space:pre-wrap}.mermaid-node{fill:#fff;stroke:var(--adlaire-color-primary)}.mermaid-edge{stroke:var(--adlaire-color-primary);fill:none}.page-updated-at{display:block;color:var(--adlaire-surface-text-muted);font-size:var(--adlaire-font-size-sm);margin-bottom:1rem}.footnotes{border-top:1px solid var(--adlaire-border-default);margin-top:2rem;font-size:var(--adlaire-font-size-sm)}.footnote-ref,.footnote-backref{font-size:var(--adlaire-font-size-sm)}.tr{list-style:none;padding:0;margin:0}.ti{margin:.25rem 0}.tl{color:inherit;text-decoration:none}.tl.active,.tl.is-active{color:var(--adlaire-color-primary);font-weight:700}.hn-link{opacity:0;margin-left:.35rem}.mh:hover .hn-link{opacity:1}.adlaire-section-toggle{margin-right:.35rem}.adlaire-section-collapsed{display:none}#progress-bar{position:fixed;top:0;left:0;height:3px;width:0%;background:var(--adlaire-color-primary);z-index:1000;transition:width .1s linear}.ch-nav{display:flex;justify-content:space-between;padding:1rem 0;margin-top:2rem;border-top:1px solid var(--adlaire-border-default)}.ch-prev,.ch-next{color:var(--adlaire-color-primary);text-decoration:none}.print-qr{display:none}.print-qr-svg{width:96px;height:96px;shape-rendering:crispEdges}.print-qr-svg rect{fill:#000}#btt{position:fixed;right:20px;bottom:20px;display:none}#btt.visible{display:block}@media(max-width:768px){#sb{transform:translateX(-100%);transition:transform .2s}#sb.open{transform:translateX(0)}#ct{margin-left:0;padding:20px}}@media print{#hdr,#sb,.cb-copy,.hn-link,#btt,.expand-code,#progress-bar,.ch-nav{display:none}#lay,#main-content{display:block;width:100%;margin:0}pre.cb[data-collapsible]{max-height:none}.print-qr{display:block}a[href^="http"]::after,a[href^="https"]::after{content:" (" attr(href) ")"}h2,h3{page-break-before:avoid}}`
}

func defaultJS() string {
	return `document.addEventListener("DOMContentLoaded",()=>{const q=(s,r=document)=>r.querySelector(s);const qa=(s,r=document)=>Array.from(r.querySelectorAll(s));const body=document.body;const sb=q("#sb"),search=q("#sb-search"),btt=q("#btt"),bar=q("#progress-bar");const enabled=n=>body?.dataset?.[n]==="true";const safe=(name,fn)=>{try{fn()}catch(e){console.warn("adlaire static init failed",name)}};const copyText=t=>navigator.clipboard&&typeof navigator.clipboard.writeText==="function"?navigator.clipboard.writeText(t).catch(()=>fallback(t)):fallback(t);const fallback=t=>new Promise((res,rej)=>{const ta=document.createElement("textarea");ta.readOnly=true;ta.value=t;document.body.appendChild(ta);ta.focus();ta.select();try{document.execCommand("copy")===true?res():rej()}catch(e){rej(e)}finally{ta.remove()}});safe("sidebar",()=>{q("#sb-toggle")?.addEventListener("click",()=>sb?.classList.toggle("open"))});safe("copy",()=>{qa(".cb-copy").forEach(btn=>btn.addEventListener("click",()=>{const code=btn.closest(".cb-wrap")?.querySelector("code")?.innerText||"";copyText(code).then(()=>{btn.textContent="✓ 完了";btn.classList.add("copied");setTimeout(()=>{btn.textContent="コピー";btn.classList.remove("copied")},1800)}).catch(()=>{btn.textContent="コピー失敗";setTimeout(()=>{btn.textContent="コピー"},1800)})}));qa(".hn-link").forEach(btn=>btn.addEventListener("click",()=>copyText(location.origin+location.pathname+btn.dataset.href).then(()=>{btn.setAttribute("aria-label","コピーしました");setTimeout(()=>btn.setAttribute("aria-label","リンクをコピー"),1800)}).catch(()=>{btn.setAttribute("aria-label","コピーに失敗しました");setTimeout(()=>btn.setAttribute("aria-label","リンクをコピー"),1800)})))});safe("code-expand",()=>{qa(".expand-code").forEach(btn=>btn.addEventListener("click",()=>{const pre=btn.previousElementSibling;pre.style.maxHeight="none";btn.hidden=true}))});safe("table-sort",()=>{qa(".mt th[data-sort]").forEach(th=>th.addEventListener("click",()=>{const table=th.closest("table"),idx=Number(th.dataset.sort),dir=th.getAttribute("aria-sort")==="ascending"?"descending":"ascending";if(!table)return;qa("th",table).forEach(h=>h.setAttribute("aria-sort","none"));th.setAttribute("aria-sort",dir);const body=table.tBodies&&table.tBodies[0];if(!body)return;const rows=qa("tr",body);rows.sort((a,b)=>{const av=a.children[idx]?.textContent.trim()||"",bv=b.children[idx]?.textContent.trim()||"",an=/^-?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?$/.test(av),bn=/^-?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?$/.test(bv);let c=an&&bn?Number(av)-Number(bv):(an?-1:bn?1:(av<bv?-1:av>bv?1:0));return dir==="ascending"?c:-c});rows.forEach(r=>body.appendChild(r))}))});safe("section-collapse",()=>{if(!enabled("sectionCollapse"))return;qa(".adlaire-section-toggle").forEach(btn=>btn.addEventListener("click",()=>{const h=btn.closest(".mh"),level=Number((h?.tagName||"H2").slice(1)),collapse=btn.getAttribute("aria-expanded")==="true";btn.setAttribute("aria-expanded",collapse?"false":"true");let n=h?.nextElementSibling;while(n&&!(/^H[1-6]$/.test(n.tagName)&&Number(n.tagName.slice(1))<=level)){n.classList.toggle("adlaire-section-collapsed",collapse);n=n.nextElementSibling}}))});safe("search",()=>{fetch((location.pathname.includes("/pages/")?"../":"")+"assets/search-index.json").catch(()=>console.warn("search index unavailable"));search?.addEventListener("input",()=>{const v=search.value.trim().replace(/[A-Z]/g,ch=>ch.toLowerCase());qa(".tl").forEach(a=>{a.parentElement.hidden=!!v&&!(a.textContent||"").trim().replace(/[A-Z]/g,ch=>ch.toLowerCase()).includes(v)})})});safe("hash-history",()=>{if(!enabled("hashHistory"))return;qa(".mh[id]").forEach(h=>h.setAttribute("tabindex","-1"));document.addEventListener("click",e=>{const a=e.target.closest("a[href^='#']");if(!a)return;const target=q(a.getAttribute("href"));if(target){history.pushState(null,"",a.getAttribute("href"));target.focus({preventScroll:true})}});window.addEventListener("popstate",()=>{if(location.hash){q(location.hash)?.focus({preventScroll:true})}})});safe("toc-active",()=>{if(!enabled("tocActive"))return;const links=qa(".toc-link");const setActive=id=>{links.forEach(a=>{const on=a.getAttribute("href")==="#"+id;a.classList.toggle("is-active",on);on?a.setAttribute("aria-current","location"):a.removeAttribute("aria-current")})};const heads=links.map(a=>q(a.getAttribute("href"))).filter(Boolean);if("IntersectionObserver"in window){const io=new IntersectionObserver(es=>{es.filter(e=>e.isIntersecting).forEach(e=>setActive(e.target.id))},{rootMargin:"-20% 0px -70% 0px"});heads.forEach(h=>io.observe(h))}else{document.addEventListener("scroll",()=>{let cur=heads[0];heads.forEach(h=>{if(h.getBoundingClientRect().top<120)cur=h});if(cur)setActive(cur.id)},{passive:true})}});safe("lightbox",()=>{if(!enabled("imageLightbox"))return;const dlg=q(".adlaire-lightbox-dialog"),img=dlg?.querySelector("img");qa(".adlaire-lightbox-trigger").forEach(btn=>btn.addEventListener("click",()=>{if(!dlg||!img)return;img.src=btn.dataset.lightboxSrc;img.alt=btn.querySelector("img")?.alt||"";dlg.showModal()}));dlg?.querySelector(".adlaire-lightbox-close")?.addEventListener("click",()=>dlg.close());dlg?.addEventListener("cancel",e=>{e.preventDefault();dlg.close()})});function onScroll(){if(bar){const d=document.documentElement,total=d.scrollHeight-d.clientHeight;bar.style.width=(total>0?Math.min(100,Math.max(0,d.scrollTop/total*100)):0)+"%"}btt?.classList.toggle("visible",scrollY>400)}document.addEventListener("scroll",onScroll,{passive:true});btt?.addEventListener("click",()=>scrollTo({top:0,behavior:"smooth"}));document.addEventListener("keydown",e=>{if(e.defaultPrevented||e.isComposing||e.repeat||e.ctrlKey||e.metaKey||e.altKey)return;const tag=e.target?.closest?.("input,textarea,select,button,a[href],[contenteditable]:not([contenteditable='false'])");if(e.key==="/"&&!tag){e.preventDefault();search?.focus()}if(e.key==="Escape"&&search&&search.value){e.preventDefault();search.value="";search.dispatchEvent(new Event("input"))}if(e.key==="t"&&!tag)scrollTo({top:0,behavior:"smooth"})});onScroll()});`
}

func writeAtomic(out string, files []siteFile) ([]string, int, int64, error) {
	if err := validateSiteFiles(files); err != nil {
		return nil, 0, 0, err
	}
	info, err := os.Lstat(out)
	if err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
		return nil, 0, 0, fmt.Errorf("output path is not directory: %s", out)
	}
	parent := filepath.Dir(out)
	base := filepath.Base(out)
	if staging, err := existingStagingPath(parent, base); err != nil {
		return nil, 0, 0, err
	} else if staging != "" {
		return nil, 0, 0, fmt.Errorf("output staging path already exists: %s", staging)
	}
	tmp := filepath.Join(parent, fmt.Sprintf("%s.tmp.%d", base, os.Getpid()))
	prev := filepath.Join(parent, fmt.Sprintf("%s.prev.%d", base, os.Getpid()))
	var warnings []string
	cleanupTmp := func(primary error) ([]string, int, int64, error) {
		if err := os.RemoveAll(tmp); err != nil {
			warnings = append(warnings, "OUTPUT_TMP_CLEANUP_FAILED: path="+tmp)
		}
		return warnings, 0, 0, primary
	}
	cleanupAfterRestore := func(primary error) ([]string, int, int64, error) {
		if err := os.RemoveAll(tmp); err != nil {
			warnings = append(warnings, "OUTPUT_TMP_CLEANUP_FAILED: path="+tmp)
			return warnings, 0, 0, primary
		}
		if err := syncDir(parent); err != nil {
			return warnings, 0, 0, fmt.Errorf("cannot restore previous output directory: %s: %w", out, err)
		}
		return warnings, 0, 0, primary
	}
	if err := os.Mkdir(tmp, 0755); err != nil {
		return nil, 0, 0, fmt.Errorf("cannot create output directory: %s", out)
	}
	for _, f := range files {
		filePath := filepath.Join(tmp, filepath.FromSlash(f.Path))
		if err := ensureOutputDir(filepath.Dir(filePath)); err != nil {
			return cleanupTmp(fmt.Errorf("cannot write output: %s", filePath))
		}
		if err := writeSyncedFile(filePath, f.Data); err != nil {
			return cleanupTmp(err)
		}
	}
	if err := validateStagedOutput(tmp, files); err != nil {
		return cleanupTmp(err)
	}
	count, size, err := outputStats(tmp)
	if err != nil {
		return cleanupTmp(fmt.Errorf("cannot write output: %s", out))
	}
	if err := syncTreeDirs(tmp); err != nil {
		return cleanupTmp(fmt.Errorf("cannot write output: %s", tmp))
	}
	prevStaged := false
	if outInfo, err := os.Lstat(out); err == nil {
		if !outInfo.IsDir() || outInfo.Mode()&os.ModeSymlink != 0 {
			return cleanupTmp(fmt.Errorf("output path is not directory: %s", out))
		}
		if err := os.Rename(out, prev); err != nil {
			return cleanupTmp(fmt.Errorf("cannot stage previous output directory: %s", out))
		}
		prevStaged = true
		if err := syncDir(parent); err != nil {
			primary := fmt.Errorf("cannot sync output parent directory: %s", parent)
			if restoreErr := restoreOutput(out, tmp, prev, prevStaged, false); restoreErr != nil {
				return warnings, 0, 0, fmt.Errorf("cannot restore previous output directory: %s: %w", out, primary)
			}
			return cleanupAfterRestore(primary)
		}
	}
	if err := os.Rename(tmp, out); err != nil {
		primary := fmt.Errorf("cannot replace output directory: %s", out)
		if restoreErr := restoreOutput(out, tmp, prev, prevStaged, false); restoreErr != nil {
			return warnings, 0, 0, fmt.Errorf("cannot restore previous output directory: %s: %w", out, primary)
		}
		return cleanupAfterRestore(primary)
	}
	if err := syncDir(parent); err != nil {
		primary := fmt.Errorf("cannot sync output parent directory: %s", parent)
		if restoreErr := restoreOutput(out, tmp, prev, prevStaged, true); restoreErr != nil {
			return warnings, 0, 0, fmt.Errorf("cannot restore previous output directory: %s: %w", out, primary)
		}
		return cleanupAfterRestore(primary)
	}
	if prevStaged {
		if err := os.RemoveAll(prev); err != nil {
			warnings = append(warnings, "OUTPUT_PREV_CLEANUP_FAILED: path="+prev)
			return warnings, count, size, nil
		}
		if err := syncDir(parent); err != nil {
			warnings = append(warnings, "OUTPUT_PREV_CLEANUP_FAILED: path="+prev)
			return warnings, count, size, nil
		}
	}
	return warnings, count, size, nil
}

func validateSiteFiles(files []siteFile) error {
	seen := map[string]bool{}
	for _, f := range files {
		if err := validateOutputRelPath(f.Path); err != nil {
			return err
		}
		if seen[f.Path] {
			return fmt.Errorf("duplicate output file: %s", f.Path)
		}
		seen[f.Path] = true
	}
	for _, required := range []string{"index.html", "assets/style.css", "assets/app.js", "assets/search-index.json", ".dependency_manifest.json"} {
		if !seen[required] {
			return fmt.Errorf("missing required output file: %s", required)
		}
	}
	for _, f := range files {
		if f.Path == "assets/search-index.json" {
			var entries []searchEntry
			if err := json.Unmarshal(f.Data, &entries); err != nil {
				return fmt.Errorf("invalid search index: %s", f.Path)
			}
			break
		}
	}
	for _, f := range files {
		if strings.HasSuffix(f.Path, ".html") {
			if err := validateGeneratedHTML(f.Path, string(f.Data)); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateGeneratedHTML(rel, htmlText string) error {
	if !strings.Contains(htmlText, "<!DOCTYPE html>") || !strings.Contains(htmlText, `<article id="main-content"`) || !strings.Contains(htmlText, "</html>") {
		return fmt.Errorf("BUILDER28_OUTPUT_VALIDATION_FAILED: %s", rel)
	}
	if regexp.MustCompile(`\saria-label=""`).MatchString(htmlText) {
		return fmt.Errorf("BUILDER28_OUTPUT_VALIDATION_FAILED: %s", rel)
	}
	if regexp.MustCompile(`\son[a-zA-Z]+\s*=`).MatchString(htmlText) {
		return fmt.Errorf("BUILDER28_ESCAPE_BLOCKED: %s", rel)
	}
	ids := map[string]bool{}
	for _, m := range regexp.MustCompile(`\sid="([^"]+)"`).FindAllStringSubmatch(htmlText, -1) {
		if ids[m[1]] {
			return fmt.Errorf("BUILDER28_OUTPUT_VALIDATION_FAILED: %s", rel)
		}
		ids[m[1]] = true
	}
	return nil
}

func validateOutputRelPath(rel string) error {
	if rel == "" || !utf8.ValidString(rel) {
		return fmt.Errorf("invalid output path: %s", rel)
	}
	if strings.ContainsAny(rel, "\x00\r\n\\") || strings.HasPrefix(rel, "/") || filepath.IsAbs(rel) {
		return fmt.Errorf("invalid output path: %s", rel)
	}
	cleaned := path.Clean(rel)
	if cleaned == "." || cleaned != rel {
		return fmt.Errorf("invalid output path: %s", rel)
	}
	for _, part := range strings.Split(cleaned, "/") {
		if part == "" || part == "." || part == ".." {
			return fmt.Errorf("invalid output path: %s", rel)
		}
	}
	if strings.HasPrefix(cleaned, "../") {
		return fmt.Errorf("invalid output path: %s", rel)
	}
	return nil
}

func ensureOutputDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.Chmod(dir, 0755)
}

func writeSyncedFile(filePath string, data []byte) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("cannot write output: %s", filePath)
	}
	var writeErr error
	if n, err := f.Write(data); err != nil {
		writeErr = err
	} else if n != len(data) {
		writeErr = io.ErrShortWrite
	}
	if writeErr == nil {
		writeErr = f.Chmod(0644)
	}
	if writeErr == nil {
		writeErr = f.Sync()
	}
	closeErr := f.Close()
	if writeErr != nil || closeErr != nil {
		os.Remove(filePath)
		return fmt.Errorf("cannot write output: %s", filePath)
	}
	return nil
}

func validateStagedOutput(root string, files []siteFile) error {
	expected := map[string]bool{}
	for _, f := range files {
		expected[f.Path] = true
	}
	err := filepath.WalkDir(root, func(filePath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if filePath == root {
			return nil
		}
		rel, err := filepath.Rel(root, filePath)
		if err != nil {
			return err
		}
		rel = slashPath(rel)
		if err := validateOutputRelPath(rel); err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("invalid output path: %s", rel)
		}
		if d.IsDir() {
			return os.Chmod(filePath, 0755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || !expected[rel] {
			return fmt.Errorf("invalid output path: %s", rel)
		}
		return os.Chmod(filePath, 0644)
	})
	if err != nil {
		return fmt.Errorf("cannot write output: %s", root)
	}
	for rel := range expected {
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("missing required output file: %s", rel)
		}
	}
	return nil
}

func syncTreeDirs(root string) error {
	var dirs []string
	err := filepath.WalkDir(root, func(filePath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			dirs = append(dirs, filePath)
		}
		return nil
	})
	if err != nil {
		return err
	}
	sort.Slice(dirs, func(i, j int) bool {
		if len(dirs[i]) == len(dirs[j]) {
			return dirs[i] > dirs[j]
		}
		return len(dirs[i]) > len(dirs[j])
	})
	for _, dir := range dirs {
		if err := os.Chmod(dir, 0755); err != nil {
			return err
		}
		if err := syncDir(dir); err != nil {
			return err
		}
	}
	return nil
}

func syncDir(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	syncErr := f.Sync()
	closeErr := f.Close()
	if syncErr != nil {
		return syncErr
	}
	return closeErr
}

func restoreOutput(out, tmp, prev string, prevStaged, tmpPublished bool) error {
	if tmpPublished {
		if _, err := os.Lstat(out); err == nil {
			if err := os.Rename(out, tmp); err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	if prevStaged {
		if err := os.Rename(prev, out); err != nil {
			return err
		}
	} else if tmpPublished {
		if _, err := os.Lstat(out); err == nil {
			return fmt.Errorf("output still exists: %s", out)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	return syncDir(filepath.Dir(out))
}

func existingStagingPath(parent, base string) (string, error) {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return "", fmt.Errorf("output parent not found: %s", parent)
	}
	prefixes := []string{base + ".tmp.", base + ".prev."}
	var matches []string
	for _, entry := range entries {
		name := entry.Name()
		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) && isPositiveDecimal(name[len(prefix):]) {
				matches = append(matches, filepath.Join(parent, name))
			}
		}
	}
	sort.Strings(matches)
	if len(matches) == 0 {
		return "", nil
	}
	return matches[0], nil
}

func isPositiveDecimal(s string) bool {
	if s == "" || s[0] == '0' {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func outputStats(root string) (int, int64, error) {
	rootInfo, err := os.Lstat(root)
	if err != nil {
		return 0, 0, err
	}
	if !rootInfo.IsDir() || rootInfo.Mode()&os.ModeSymlink != 0 {
		return 0, 0, fmt.Errorf("invalid output path: %s", root)
	}
	count := 0
	var size int64
	err = filepath.WalkDir(root, func(filePath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if filePath == root {
			return nil
		}
		rel, err := filepath.Rel(root, filePath)
		if err != nil {
			return err
		}
		rel = slashPath(rel)
		if err := validateOutputRelPath(rel); err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("invalid output path: %s", rel)
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("invalid output path: %s", rel)
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
	return os.ReadFile(path)
}
