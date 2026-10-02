package obsidian

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"unicode"
	"unicode/utf8"

	"github.com/fqwink/build-scripts/components/builder"
)

func executeVaultBuild(cfg obsidianConfig, stdout, stderr io.Writer) int {
	if obsErr := validateVaultRoot(cfg.VaultRoot); obsErr != nil {
		writeObsidianError(stderr, obsErr)
		return obsErr.Exit
	}
	result, obsErr := normalizeVault(cfg)
	if obsErr != nil {
		writeObsidianError(stderr, obsErr)
		return obsErr.Exit
	}
	defer os.RemoveAll(result.intermediateRoot)
	builderArgs := replaceBuilderSource(cfg.BuilderArgs, result.intermediateRoot)
	var builderStdout, builderStderr bytes.Buffer
	code := builder.RunBuild(builderArgs, &builderStdout, &builderStderr)
	_, _ = stdout.Write(builderStdout.Bytes())
	if code != 0 {
		fmt.Fprintln(stderr, "OBSIDIAN_BUILDER_FAILED")
		_, _ = stderr.Write(builderStderr.Bytes())
		return 1
	}
	if obsErr := copyAssetsToOutput(cfg, result.assets); obsErr != nil {
		writeObsidianError(stderr, obsErr)
		return obsErr.Exit
	}
	if cfg.ReportFile != "" {
		if obsErr := copyReport(cfg, result.mapPath); obsErr != nil {
			writeObsidianError(stderr, obsErr)
			return obsErr.Exit
		}
	}
	fmt.Fprintf(stdout, "[REPORT] obsidian_notes=%d obsidian_assets=%d obsidian_diagnostics=%d obsidian_unresolved_links=%d\n", len(result.notes), len(result.assets), len(result.diagnostics), result.unresolvedLinks)
	return 0
}

type normalizeResult struct {
	intermediateRoot string
	mapPath          string
	notes            []noteState
	assets           []obsidianAsset
	diagnostics      []obsidianDiagnostic
	unresolvedLinks  int
}

func normalizeVault(cfg obsidianConfig) (normalizeResult, *obsidianError) {
	filter, obsErr := readFilter(cfg.FilterFile)
	if obsErr != nil {
		return normalizeResult{}, obsErr
	}
	files, obsErr := inventoryVault(cfg.VaultRoot)
	if obsErr != nil {
		return normalizeResult{}, obsErr
	}
	if _, ok := files[cfg.EntryPath]; !ok {
		return normalizeResult{}, newObsError("OBSIDIAN_INVALID_ENTRY", 2, cfg.EntryPath, 0, 0)
	}
	if anyFilterMatch(filter.ExcludeNotes, cfg.EntryPath) {
		return normalizeResult{}, newObsError("OBSIDIAN_ENTRY_EXCLUDED", 2, cfg.EntryPath, 0, 0)
	}
	for _, pattern := range append(filter.IncludeNotes, filter.ExcludeNotes...) {
		if obsErr := validateFilterType(pattern, files, "note"); obsErr != nil {
			return normalizeResult{}, obsErr
		}
	}
	for _, pattern := range append(filter.IncludeAssets, filter.ExcludeAssets...) {
		if obsErr := validateFilterType(pattern, files, "asset"); obsErr != nil {
			return normalizeResult{}, obsErr
		}
	}
	noteSet := map[string]bool{cfg.EntryPath: true}
	assetSet := map[string]bool{}
	states := map[string]*noteState{}
	for changed := true; changed; {
		changed = false
		for rel := range noteSet {
			if states[rel] != nil {
				continue
			}
			state, obsErr := parseNoteState(cfg, files, rel)
			if obsErr != nil {
				return normalizeResult{}, obsErr
			}
			states[rel] = &state
			for _, link := range state.OutgoingLinks {
				if link.Status == "resolved" && link.TargetPath != "" && !noteSet[link.TargetPath] {
					noteSet[link.TargetPath] = true
					changed = true
				}
			}
			for _, asset := range state.AssetEmbeds {
				if asset.Status == "resolved" && asset.SourcePath != "" {
					assetSet[asset.SourcePath] = true
				}
			}
		}
		for _, rel := range filter.IncludeNotes {
			if matchExisting(files, rel, "note") && !noteSet[rel] {
				noteSet[rel] = true
				changed = true
			}
		}
	}
	for rel := range noteSet {
		if anyFilterMatch(filter.ExcludeNotes, rel) {
			delete(noteSet, rel)
			delete(states, rel)
		}
	}
	for _, rel := range filter.IncludeAssets {
		if matchExisting(files, rel, "asset") {
			assetSet[rel] = true
		}
	}
	for rel := range assetSet {
		if anyFilterMatch(filter.ExcludeAssets, rel) {
			delete(assetSet, rel)
		}
	}
	for rel := range assetSet {
		used := false
		for _, state := range states {
			for _, asset := range state.AssetEmbeds {
				if asset.SourcePath == rel && asset.Status == "resolved" {
					used = true
				}
			}
		}
		if !used {
			diag := obsidianDiagnostic{Code: "OBSIDIAN_UNUSED_ASSET", Severity: "warning", SourcePath: rel, Target: rel}
			if cfg.Strict {
				return normalizeResult{}, newObsError("OBSIDIAN_UNUSED_ASSET", 2, rel, 0, 0)
			}
			if states[cfg.EntryPath] != nil {
				states[cfg.EntryPath].Diagnostics = append(states[cfg.EntryPath].Diagnostics, diag)
			}
		}
	}
	root, err := os.MkdirTemp("", "adlaire-ci-obsidian-")
	if err != nil {
		return normalizeResult{}, newObsError("OBSIDIAN_INTERMEDIATE_WRITE_FAILED", 1, "", 0, 0)
	}
	result := normalizeResult{intermediateRoot: root}
	if obsErr := writeIntermediate(root, cfg, files, states, noteSet, assetSet, &result); obsErr != nil {
		_ = os.RemoveAll(root)
		return normalizeResult{}, obsErr
	}
	return result, nil
}

func inventoryVault(root string) (map[string]vaultFile, *obsidianError) {
	files := map[string]vaultFile{}
	err := filepath.WalkDir(root, func(abs string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if abs == root {
			return nil
		}
		rel, err := filepath.Rel(root, abs)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, ".obsidian/") || rel == ".obsidian" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := os.Lstat(abs)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return newWalkError("OBSIDIAN_PATH_ESCAPE", rel)
		}
		if entry.IsDir() {
			if validateRelativePath(rel, false) != nil {
				return newWalkError("OBSIDIAN_PATH_ESCAPE", rel)
			}
			return nil
		}
		if !info.Mode().IsRegular() {
			return newWalkError("OBSIDIAN_INVALID_VAULT", rel)
		}
		if hasMultipleHardlinks(info) {
			return newWalkError("OBSIDIAN_HARDLINK_FORBIDDEN", rel)
		}
		if validateRelativePath(rel, false) != nil {
			return newWalkError("OBSIDIAN_PATH_ESCAPE", rel)
		}
		ext := strings.ToLower(filepath.Ext(rel))
		kind := "asset"
		if ext == ".md" {
			kind = "note"
		}
		if ext == ".canvas" {
			return newWalkError("OBSIDIAN_ASSET_UNSUPPORTED", rel)
		}
		sum, readErr := fileSHA256(abs)
		if readErr != nil {
			return readErr
		}
		files[rel] = vaultFile{SourcePath: rel, AbsPath: abs, Size: info.Size(), SHA256: sum, Kind: kind}
		return nil
	})
	if err != nil {
		if wrapped, ok := err.(walkError); ok {
			return nil, newObsError(wrapped.Code, 2, wrapped.Path, 0, 0)
		}
		return nil, newObsError("OBSIDIAN_INVALID_VAULT", 2, "", 0, 0)
	}
	return files, nil
}

type walkError struct {
	Code string
	Path string
}

func (e walkError) Error() string {
	return e.Code + ": " + e.Path
}

func newWalkError(code, rel string) error {
	return walkError{Code: code, Path: rel}
}

func hasMultipleHardlinks(info os.FileInfo) bool {
	stat, ok := info.Sys().(*syscall.Stat_t)
	return ok && stat.Nlink >= 2
}

func validateFilterType(pattern string, files map[string]vaultFile, want string) *obsidianError {
	for rel, file := range files {
		if filterMatch(pattern, rel) && file.Kind != want {
			return newObsError("OBSIDIAN_FILTER_TYPE_MISMATCH", 2, pattern, 0, 0)
		}
	}
	return nil
}

func matchExisting(files map[string]vaultFile, rel string, kind string) bool {
	file, ok := files[rel]
	return ok && file.Kind == kind
}

func parseNoteState(cfg obsidianConfig, files map[string]vaultFile, rel string) (noteState, *obsidianError) {
	file := files[rel]
	data, err := os.ReadFile(file.AbsPath)
	if err != nil || !utf8.Valid(data) || bytes.IndexByte(data, 0) >= 0 {
		return noteState{}, newObsError("OBSIDIAN_INVALID_ENTRY", 2, rel, 0, 0)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if containsForbiddenControl(text) {
		return noteState{}, newObsError("OBSIDIAN_INVALID_ENTRY", 2, rel, 0, 0)
	}
	if strings.HasPrefix(text, "---\n") {
		return noteState{}, newObsError("OBSIDIAN_YAML_FRONTMATTER_UNSUPPORTED", 2, rel, 1, 1)
	}
	if strings.Contains(text, "<%") || strings.Contains(text, "fetch(") {
		return noteState{}, newObsError("OBSIDIAN_WIKILINK_MALFORMED", 2, rel, 1, 1)
	}
	protected, obsErr := protectedMarkdownRanges(text, rel)
	if obsErr != nil {
		return noteState{}, obsErr
	}
	headings := collectHeadingSlugs(text)
	state := noteState{
		SourcePath:     rel,
		NormalizedPath: normalizedNotePath(cfg.EntryPath, rel),
		RawText:        text,
		Title:          firstTitle(text, rel),
		SHA256:         file.SHA256,
		TrailingLF:     strings.HasSuffix(text, "\n"),
		Headings:       headings,
		OutgoingLinks:  []obsidianLink{},
		AssetEmbeds:    []obsidianAsset{},
		Tags:           []string{},
		Diagnostics:    []obsidianDiagnostic{},
	}
	tokens, obsErr := scanTokens(text, protected, rel)
	if obsErr != nil {
		return noteState{}, obsErr
	}
	replacements := []replacement{}
	for _, token := range tokens {
		if token.Embed {
			asset, replacementText, obsErr := resolveAssetToken(cfg, files, state, token)
			if obsErr != nil {
				return noteState{}, obsErr
			}
			state.AssetEmbeds = append(state.AssetEmbeds, asset)
			replacements = append(replacements, replacement{Start: token.Start, End: token.End, Text: replacementText})
			if asset.Status != "resolved" {
				state.Diagnostics = append(state.Diagnostics, obsidianDiagnostic{Code: asset.ErrorCode, Severity: "warning", SourcePath: rel, Line: token.Line, Column: token.Column, Target: asset.SourcePath})
			}
			continue
		}
		link, replacementText, obsErr := resolveLinkToken(cfg, files, state, token)
		if obsErr != nil {
			return noteState{}, obsErr
		}
		state.OutgoingLinks = append(state.OutgoingLinks, link)
		replacements = append(replacements, replacement{Start: token.Start, End: token.End, Text: replacementText})
		if link.Status != "resolved" {
			state.Diagnostics = append(state.Diagnostics, obsidianDiagnostic{Code: link.ErrorCode, Severity: "warning", SourcePath: rel, Line: token.Line, Column: token.Column, Target: link.Target})
		}
	}
	tags, tagDiagnostics, obsErr := collectTags(cfg, text, protected, rel)
	if obsErr != nil {
		return noteState{}, obsErr
	}
	state.Tags = tags
	state.Diagnostics = append(state.Diagnostics, tagDiagnostics...)
	sort.Slice(replacements, func(i, j int) bool { return replacements[i].Start > replacements[j].Start })
	normalized := text
	for _, repl := range replacements {
		normalized = normalized[:repl.Start] + repl.Text + normalized[repl.End:]
	}
	state.NormalizedText = normalized
	sortLinks(state.OutgoingLinks)
	sortAssets(state.AssetEmbeds)
	sortDiagnostics(state.Diagnostics)
	return state, nil
}

type token struct {
	Raw    string
	Body   string
	Embed  bool
	Start  int
	End    int
	Line   int
	Column int
}

type replacement struct {
	Start int
	End   int
	Text  string
}

func containsForbiddenControl(text string) bool {
	for _, r := range text {
		if r < 0x20 && r != '\n' && r != '\r' && r != '\t' {
			return true
		}
	}
	return false
}

func protectedMarkdownRanges(text, rel string) ([]bool, *obsidianError) {
	protected := make([]bool, len(text))
	inFence := false
	offset := 0
	for _, line := range strings.SplitAfter(text, "\n") {
		trimmed := strings.TrimSpace(line)
		isFence := strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
		if isFence {
			info := strings.ToLower(strings.TrimSpace(strings.TrimLeft(trimmed, "`~")))
			if !inFence && (info == "yaml" || info == "yml") {
				ln, col := lineColumn(text, offset)
				return nil, newObsError("OBSIDIAN_YAML_BLOCK_UNSUPPORTED", 2, rel, ln, col)
			}
			if !inFence && info == "dataview" {
				ln, col := lineColumn(text, offset)
				return nil, newObsError("OBSIDIAN_WIKILINK_MALFORMED", 2, rel, ln, col)
			}
			inFence = !inFence
		}
		if inFence || isFence {
			for i := offset; i < offset+len(line) && i < len(protected); i++ {
				protected[i] = true
			}
		}
		markInlineCode(protected, line, offset)
		offset += len(line)
	}
	return protected, nil
}

func markInlineCode(protected []bool, line string, offset int) {
	start := -1
	for i := 0; i < len(line); i++ {
		if line[i] != '`' {
			continue
		}
		if start < 0 {
			start = i
			continue
		}
		for j := offset + start; j <= offset+i && j < len(protected); j++ {
			protected[j] = true
		}
		start = -1
	}
}

func scanTokens(text string, protected []bool, rel string) ([]token, *obsidianError) {
	var tokens []token
	for i := 0; i < len(text); i++ {
		if protected[i] || (i > 0 && text[i-1] == '\\') {
			continue
		}
		embed := false
		start := i
		if strings.HasPrefix(text[i:], "![[") {
			embed = true
		} else if strings.HasPrefix(text[i:], "[[") {
			start = i
		} else {
			continue
		}
		prefix := 2
		if embed {
			prefix = 3
		}
		endRel := strings.Index(text[i+prefix:], "]]")
		if endRel < 0 {
			ln, col := lineColumn(text, i)
			return nil, newObsError("OBSIDIAN_WIKILINK_MALFORMED", 2, rel, ln, col)
		}
		end := i + prefix + endRel + 2
		body := text[i+prefix : end-2]
		if strings.Contains(body, "[[") || strings.Contains(body, "![[") {
			ln, col := lineColumn(text, i)
			return nil, newObsError("OBSIDIAN_WIKILINK_MALFORMED", 2, rel, ln, col)
		}
		ln, col := lineColumn(text, start)
		tokens = append(tokens, token{Raw: text[start:end], Body: body, Embed: embed, Start: start, End: end, Line: ln, Column: col})
		i = end - 1
	}
	return tokens, nil
}

func resolveLinkToken(cfg obsidianConfig, files map[string]vaultFile, state noteState, token token) (obsidianLink, string, *obsidianError) {
	target, heading, alias, ok := splitWikilinkBody(token.Body)
	link := obsidianLink{Raw: token.Raw, Target: target, Heading: heading, Alias: alias, Status: "unresolved", ErrorCode: "OBSIDIAN_UNRESOLVED_LINK", Line: token.Line, Column: token.Column}
	if !ok || invalidDisplay(alias) || invalidDisplay(heading) {
		return link, "", newObsError("OBSIDIAN_WIKILINK_MALFORMED", 2, state.SourcePath, token.Line, token.Column)
	}
	if validateRelativePath(noteTargetPath(target), false) != nil {
		return link, "", newObsError("OBSIDIAN_WIKILINK_MALFORMED", 2, state.SourcePath, token.Line, token.Column)
	}
	resolved, status, code := resolveNoteTarget(files, target)
	link.TargetPath = resolved
	link.Status = status
	link.ErrorCode = code
	if link.Status == "ambiguous" {
		return link, "", newObsError(link.ErrorCode, 2, state.SourcePath, token.Line, token.Column)
	}
	if status == "resolved" && heading != "" {
		targetState, obsErr := parseNoteHeadingOnly(files, resolved)
		if obsErr != nil {
			return link, "", obsErr
		}
		if !targetState.Headings[slugify(heading)] {
			link.Status = "unresolved"
			link.ErrorCode = "OBSIDIAN_UNRESOLVED_LINK"
		}
	}
	if link.Status != "resolved" {
		if cfg.Strict {
			return link, "", newObsError(link.ErrorCode, 2, state.SourcePath, token.Line, token.Column)
		}
		return link, token.Raw, nil
	}
	display := alias
	if display == "" {
		if heading != "" {
			display = heading
		} else {
			display = strings.TrimSuffix(path.Base(target), path.Ext(target))
		}
	}
	if invalidDisplay(display) {
		return link, "", newObsError("OBSIDIAN_WIKILINK_MALFORMED", 2, state.SourcePath, token.Line, token.Column)
	}
	dest := relativeNormalizedPath(state.NormalizedPath, normalizedNotePath(cfg.EntryPath, resolved))
	if heading != "" {
		dest += "#" + slugify(heading)
	}
	return link, "[" + display + "](" + dest + ")", nil
}

func resolveAssetToken(cfg obsidianConfig, files map[string]vaultFile, state noteState, token token) (obsidianAsset, string, *obsidianError) {
	target, _, alias, ok := splitWikilinkBody(token.Body)
	asset := obsidianAsset{Raw: token.Raw, SourcePath: target, Status: "unresolved", ErrorCode: "OBSIDIAN_UNRESOLVED_LINK", Line: token.Line, Column: token.Column}
	if !ok || invalidDisplay(alias) {
		return asset, "", newObsError("OBSIDIAN_WIKILINK_MALFORMED", 2, state.SourcePath, token.Line, token.Column)
	}
	ext := strings.ToLower(path.Ext(target))
	if ext == "" || ext == ".md" {
		return asset, "", newObsError("OBSIDIAN_NOTE_EMBED_UNSUPPORTED", 2, state.SourcePath, token.Line, token.Column)
	}
	if !allowedAssetExt(ext) || validateRelativePath(target, false) != nil {
		return asset, "", newObsError("OBSIDIAN_ASSET_UNSUPPORTED", 2, state.SourcePath, token.Line, token.Column)
	}
	resolved, status, code := resolveAssetTarget(files, target)
	asset.SourcePath = resolved
	asset.Status = status
	asset.ErrorCode = code
	if status != "resolved" {
		if cfg.Strict {
			return asset, "", newObsError(code, 2, state.SourcePath, token.Line, token.Column)
		}
		return asset, token.Raw, nil
	}
	file := files[resolved]
	asset.NormalizedPath = "assets/" + resolved
	asset.SHA256 = file.SHA256
	asset.Size = file.Size
	asset.MediaType = mediaTypeForExt(ext)
	asset.ErrorCode = ""
	display := alias
	if display == "" {
		display = path.Base(resolved)
	}
	if invalidDisplay(display) {
		return asset, "", newObsError("OBSIDIAN_WIKILINK_MALFORMED", 2, state.SourcePath, token.Line, token.Column)
	}
	dest := relativeNormalizedPath(state.NormalizedPath, asset.NormalizedPath)
	if ext == ".pdf" {
		return asset, "[" + display + "](" + dest + ")", nil
	}
	return asset, "![" + display + "](" + dest + ")", nil
}

func splitWikilinkBody(body string) (string, string, string, bool) {
	parts := strings.Split(body, "|")
	if len(parts) > 2 {
		return "", "", "", false
	}
	targetPart := strings.TrimSpace(parts[0])
	alias := ""
	if len(parts) == 2 {
		alias = strings.TrimSpace(parts[1])
	}
	heading := ""
	if hash := strings.IndexByte(targetPart, '#'); hash >= 0 {
		heading = strings.TrimSpace(targetPart[hash+1:])
		targetPart = strings.TrimSpace(targetPart[:hash])
	}
	if targetPart == "" || strings.ContainsAny(targetPart+heading+alias, "\x00\r\n\\[]") {
		return "", "", "", false
	}
	return targetPart, heading, alias, true
}

func invalidDisplay(value string) bool {
	return strings.ContainsAny(value, "\x00\r\n[]")
}

func noteTargetPath(target string) string {
	if strings.HasSuffix(target, ".md") {
		return target
	}
	return target + ".md"
}

func resolveNoteTarget(files map[string]vaultFile, target string) (string, string, string) {
	candidate := noteTargetPath(target)
	if file, ok := files[candidate]; ok && file.Kind == "note" {
		return candidate, "resolved", ""
	}
	if strings.Contains(candidate, "/") {
		return candidate, "unresolved", "OBSIDIAN_UNRESOLVED_LINK"
	}
	var matches []string
	for rel, file := range files {
		if file.Kind == "note" && path.Base(rel) == candidate {
			matches = append(matches, rel)
		}
	}
	sort.Strings(matches)
	if len(matches) == 1 {
		return matches[0], "resolved", ""
	}
	if len(matches) > 1 {
		return candidate, "ambiguous", "OBSIDIAN_AMBIGUOUS_NOTE_LINK"
	}
	return candidate, "unresolved", "OBSIDIAN_UNRESOLVED_LINK"
}

func resolveAssetTarget(files map[string]vaultFile, target string) (string, string, string) {
	if file, ok := files[target]; ok && file.Kind == "asset" {
		return target, "resolved", ""
	}
	if strings.Contains(target, "/") {
		return target, "unresolved", "OBSIDIAN_UNRESOLVED_LINK"
	}
	var matches []string
	for rel, file := range files {
		if file.Kind == "asset" && path.Base(rel) == target {
			matches = append(matches, rel)
		}
	}
	sort.Strings(matches)
	if len(matches) == 1 {
		return matches[0], "resolved", ""
	}
	return target, "unresolved", "OBSIDIAN_UNRESOLVED_LINK"
}

func parseNoteHeadingOnly(files map[string]vaultFile, rel string) (noteState, *obsidianError) {
	file := files[rel]
	data, err := os.ReadFile(file.AbsPath)
	if err != nil || !utf8.Valid(data) {
		return noteState{}, newObsError("OBSIDIAN_INVALID_ENTRY", 2, rel, 0, 0)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	return noteState{Headings: collectHeadingSlugs(text)}, nil
}

func collectHeadingSlugs(text string) map[string]bool {
	counts := map[string]int{}
	slugs := map[string]bool{}
	for _, line := range strings.Split(text, "\n") {
		level := 0
		for level < len(line) && level < 6 && line[level] == '#' {
			level++
		}
		if level == 0 || level >= len(line) || line[level] != ' ' {
			continue
		}
		base := slugify(strings.TrimSpace(line[level+1:]))
		slug := uniqueSlug(base, counts)
		slugs[slug] = true
	}
	return slugs
}

func firstTitle(text, rel string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return strings.TrimSuffix(path.Base(rel), path.Ext(rel))
}

func collectTags(cfg obsidianConfig, text string, protected []bool, rel string) ([]string, []obsidianDiagnostic, *obsidianError) {
	seen := map[string]bool{}
	dupes := map[string]bool{}
	for i := 0; i < len(text); i++ {
		if protected[i] || text[i] != '#' {
			continue
		}
		if i > 0 && !strings.ContainsRune(" \t([{", rune(text[i-1])) {
			continue
		}
		j := i + 1
		if j >= len(text) || !isTagStart(text[j]) {
			continue
		}
		for j < len(text) && isTagBody(text[j]) {
			j++
		}
		tag := text[i+1 : j]
		if seen[tag] {
			dupes[tag] = true
		}
		seen[tag] = true
		i = j - 1
	}
	tags := make([]string, 0, len(seen))
	for tag := range seen {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	var diagnostics []obsidianDiagnostic
	for tag := range dupes {
		if cfg.Strict {
			return nil, nil, newObsError("OBSIDIAN_DUPLICATE_TAG", 2, rel, 0, 0)
		}
		diagnostics = append(diagnostics, obsidianDiagnostic{Code: "OBSIDIAN_DUPLICATE_TAG", Severity: "warning", SourcePath: rel, Target: tag})
	}
	sortDiagnostics(diagnostics)
	return tags, diagnostics, nil
}

func isTagStart(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9')
}

func isTagBody(b byte) bool {
	return isTagStart(b) || b == '_' || b == '/' || b == '-'
}

func writeIntermediate(root string, cfg obsidianConfig, files map[string]vaultFile, states map[string]*noteState, noteSet map[string]bool, assetSet map[string]bool, result *normalizeResult) *obsidianError {
	var notePaths []string
	for rel := range noteSet {
		notePaths = append(notePaths, rel)
	}
	sort.Strings(notePaths)
	for _, rel := range notePaths {
		state := states[rel]
		target := filepath.Join(root, filepath.FromSlash(state.NormalizedPath))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return newObsError("OBSIDIAN_INTERMEDIATE_WRITE_FAILED", 1, state.SourcePath, 0, 0)
		}
		if err := os.WriteFile(target, []byte(state.NormalizedText), 0o644); err != nil {
			return newObsError("OBSIDIAN_INTERMEDIATE_WRITE_FAILED", 1, state.SourcePath, 0, 0)
		}
		result.notes = append(result.notes, *state)
	}
	var assetPaths []string
	for rel := range assetSet {
		assetPaths = append(assetPaths, rel)
	}
	sort.Strings(assetPaths)
	for _, rel := range assetPaths {
		file := files[rel]
		if file.Kind != "asset" {
			continue
		}
		asset := obsidianAsset{SourcePath: rel, NormalizedPath: "assets/" + rel, SHA256: file.SHA256, Size: file.Size, MediaType: mediaTypeForExt(strings.ToLower(path.Ext(rel))), Status: "resolved", ErrorCode: ""}
		result.assets = append(result.assets, asset)
		target := filepath.Join(root, filepath.FromSlash(asset.NormalizedPath))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return newObsError("OBSIDIAN_INTERMEDIATE_WRITE_FAILED", 1, rel, 0, 0)
		}
		if err := copyFile(file.AbsPath, target); err != nil {
			return newObsError("OBSIDIAN_INTERMEDIATE_WRITE_FAILED", 1, rel, 0, 0)
		}
	}
	obsMap := buildMap(cfg, files, result.notes, result.assets)
	data, err := json.Marshal(obsMap)
	if err != nil {
		return newObsError("OBSIDIAN_INTERMEDIATE_WRITE_FAILED", 1, "", 0, 0)
	}
	data = append(data, '\n')
	mapPath := filepath.Join(root, "obsidian_map.json")
	if err := os.WriteFile(mapPath, data, 0o644); err != nil {
		return newObsError("OBSIDIAN_INTERMEDIATE_WRITE_FAILED", 1, "", 0, 0)
	}
	result.mapPath = mapPath
	result.diagnostics = obsMap.Diagnostics
	for _, diag := range obsMap.Diagnostics {
		if diag.Code == "OBSIDIAN_UNRESOLVED_LINK" {
			result.unresolvedLinks++
		}
	}
	return nil
}

func buildMap(cfg obsidianConfig, files map[string]vaultFile, notes []noteState, assets []obsidianAsset) obsidianMap {
	sort.Slice(notes, func(i, j int) bool { return notes[i].SourcePath < notes[j].SourcePath })
	sortAssets(assets)
	mapNotes := make([]obsidianMapNote, 0, len(notes))
	var diagnostics []obsidianDiagnostic
	for _, state := range notes {
		sortLinks(state.OutgoingLinks)
		sortAssets(state.AssetEmbeds)
		sort.Strings(state.Tags)
		sortDiagnostics(state.Diagnostics)
		diagnostics = append(diagnostics, state.Diagnostics...)
		mapNotes = append(mapNotes, obsidianMapNote{
			SourcePath:       state.SourcePath,
			NormalizedPath:   state.NormalizedPath,
			Title:            state.Title,
			SourceSHA256:     state.SHA256,
			SourceTrailingLF: state.TrailingLF,
			OutgoingLinks:    nonNilLinks(state.OutgoingLinks),
			AssetEmbeds:      nonNilAssets(state.AssetEmbeds),
			Tags:             nonNilStrings(state.Tags),
			Diagnostics:      nonNilDiagnostics(state.Diagnostics),
		})
	}
	sortDiagnostics(diagnostics)
	return obsidianMap{
		SchemaVersion:       obsidianMapSchemaVersion,
		VaultDigest:         vaultDigest(files, notes, assets),
		EntrySourcePath:     cfg.EntryPath,
		EntryNormalizedPath: "index.md",
		Notes:               mapNotes,
		Assets:              nonNilAssets(assets),
		Diagnostics:         nonNilDiagnostics(diagnostics),
	}
}

func vaultDigest(files map[string]vaultFile, notes []noteState, assets []obsidianAsset) string {
	seen := map[string]bool{}
	for _, note := range notes {
		seen[note.SourcePath] = true
	}
	for _, asset := range assets {
		seen[asset.SourcePath] = true
	}
	var paths []string
	for rel := range seen {
		paths = append(paths, rel)
	}
	sort.Strings(paths)
	var b strings.Builder
	for _, rel := range paths {
		file := files[rel]
		b.WriteString(rel)
		b.WriteByte('\n')
		b.WriteString(file.SHA256)
		b.WriteByte('\n')
		b.WriteString(fmt.Sprintf("%d", file.Size))
		b.WriteByte('\n')
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

func copyAssetsToOutput(cfg obsidianConfig, assets []obsidianAsset) *obsidianError {
	if cfg.OutputRoot == "" || len(assets) == 0 {
		return nil
	}
	for _, asset := range assets {
		src := filepath.Join(cfg.VaultRoot, filepath.FromSlash(asset.SourcePath))
		dst := filepath.Join(cfg.OutputRoot, filepath.FromSlash(asset.NormalizedPath))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return newObsError("OBSIDIAN_INTERMEDIATE_WRITE_FAILED", 1, asset.SourcePath, 0, 0)
		}
		if err := copyFile(src, dst); err != nil {
			return newObsError("OBSIDIAN_INTERMEDIATE_WRITE_FAILED", 1, asset.SourcePath, 0, 0)
		}
	}
	return nil
}

func copyReport(cfg obsidianConfig, mapPath string) *obsidianError {
	if cfg.OutputRoot == "" {
		return newObsError("OBSIDIAN_REPORT_PATH_INVALID", 2, cfg.ReportFile, 0, 0)
	}
	outAbs, err := filepath.Abs(cfg.OutputRoot)
	if err != nil {
		return newObsError("OBSIDIAN_REPORT_PATH_INVALID", 2, cfg.ReportFile, 0, 0)
	}
	reportAbs, err := filepath.Abs(cfg.ReportFile)
	if err != nil {
		return newObsError("OBSIDIAN_REPORT_PATH_INVALID", 2, cfg.ReportFile, 0, 0)
	}
	rel, err := filepath.Rel(outAbs, reportAbs)
	if err != nil || strings.HasPrefix(filepath.ToSlash(rel), "../") || filepath.IsAbs(rel) || rel == ".." {
		return newObsError("OBSIDIAN_REPORT_PATH_INVALID", 2, cfg.ReportFile, 0, 0)
	}
	if err := os.MkdirAll(filepath.Dir(reportAbs), 0o755); err != nil {
		return newObsError("OBSIDIAN_REPORT_COPY_FAILED", 1, cfg.ReportFile, 0, 0)
	}
	if err := copyFile(mapPath, reportAbs); err != nil {
		return newObsError("OBSIDIAN_REPORT_COPY_FAILED", 1, cfg.ReportFile, 0, 0)
	}
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func normalizedNotePath(entry, rel string) string {
	if rel == entry {
		return "index.md"
	}
	return "notes/" + rel
}

func relativeNormalizedPath(from, to string) string {
	base := path.Dir(from)
	if base == "." {
		base = ""
	}
	rel, err := filepath.Rel(filepath.FromSlash(base), filepath.FromSlash(to))
	if err != nil {
		return to
	}
	return filepath.ToSlash(rel)
}

func allowedAssetExt(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".pdf":
		return true
	default:
		return false
	}
}

func mediaTypeForExt(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	default:
		if got := mime.TypeByExtension(ext); got != "" {
			return got
		}
		return "application/octet-stream"
	}
}

func lineColumn(text string, offset int) (int, int) {
	line := 1
	lineStart := 0
	for i := 0; i < offset && i < len(text); i++ {
		if text[i] == '\n' {
			line++
			lineStart = i + 1
		}
	}
	return line, offset - lineStart + 1
}

func sortLinks(items []obsidianLink) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Line != items[j].Line {
			return items[i].Line < items[j].Line
		}
		if items[i].Column != items[j].Column {
			return items[i].Column < items[j].Column
		}
		return items[i].Raw < items[j].Raw
	})
}

func sortAssets(items []obsidianAsset) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].SourcePath != items[j].SourcePath {
			return items[i].SourcePath < items[j].SourcePath
		}
		if items[i].Line != items[j].Line {
			return items[i].Line < items[j].Line
		}
		return items[i].Column < items[j].Column
	})
}

func sortDiagnostics(items []obsidianDiagnostic) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].SourcePath != items[j].SourcePath {
			return items[i].SourcePath < items[j].SourcePath
		}
		if items[i].Line != items[j].Line {
			return items[i].Line < items[j].Line
		}
		if items[i].Column != items[j].Column {
			return items[i].Column < items[j].Column
		}
		return items[i].Code < items[j].Code
	})
}

func nonNilLinks(items []obsidianLink) []obsidianLink {
	if items == nil {
		return []obsidianLink{}
	}
	return items
}

func nonNilAssets(items []obsidianAsset) []obsidianAsset {
	if items == nil {
		return []obsidianAsset{}
	}
	return items
}

func nonNilDiagnostics(items []obsidianDiagnostic) []obsidianDiagnostic {
	if items == nil {
		return []obsidianDiagnostic{}
	}
	return items
}

func nonNilStrings(items []string) []string {
	if items == nil {
		return []string{}
	}
	return items
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
