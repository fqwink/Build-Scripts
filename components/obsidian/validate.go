package obsidian

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

func SetBinaryVersion(version string) {
	obsidianBinaryVersion = version
}

func Owner() string {
	return ownerName
}

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == ownerName && validateExactOwnerFiles(contract.Files, []string{"obsidian.go", "model.go", "validate.go", "execute.go", "obsidian_test.go"})
}

func validateExactOwnerFiles(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	seen := map[string]bool{}
	for i, file := range got {
		if file != want[i] || file == "" || seen[file] {
			return false
		}
		seen[file] = true
	}
	return true
}

func parseBuildArgs(args []string) (obsidianConfig, bool, *obsidianError) {
	cfg := obsidianConfig{BuilderArgs: []string{}}
	seen := map[string]bool{}
	hasObsidianOption := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		name, value, hasValue := splitOption(arg)
		switch name {
		case "--input-mode", "--obsidian-vault", "--obsidian-entry", "--obsidian-filter-file", "--obsidian-report-file":
			if seen[name] {
				return cfg, true, newObsError("OBSIDIAN_INVALID_FILTER", 2, name, 0, 0)
			}
			seen[name] = true
			if !hasValue {
				if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
					return cfg, true, newObsError("OBSIDIAN_INVALID_FILTER", 2, name, 0, 0)
				}
				i++
				value = args[i]
			}
			if name != "--input-mode" {
				hasObsidianOption = true
			}
			switch name {
			case "--input-mode":
				cfg.InputMode = value
			case "--obsidian-vault":
				cfg.VaultRoot = value
			case "--obsidian-entry":
				cfg.EntryPath = filepath.ToSlash(value)
			case "--obsidian-filter-file":
				cfg.FilterFile = value
			case "--obsidian-report-file":
				cfg.ReportFile = value
			}
		case "--strict":
			cfg.Strict = true
			cfg.BuilderArgs = append(cfg.BuilderArgs, arg)
		default:
			if strings.HasPrefix(name, "--obsidian-") {
				hasObsidianOption = true
				return cfg, true, newObsError("OBSIDIAN_INVALID_FILTER", 2, name, 0, 0)
			}
			cfg.BuilderArgs = append(cfg.BuilderArgs, arg)
			if !hasValue && optionConsumesValue(name) && i+1 < len(args) {
				i++
				cfg.BuilderArgs = append(cfg.BuilderArgs, args[i])
			}
		}
	}
	if cfg.InputMode == "" {
		if hasObsidianOption {
			return cfg, true, newObsError("OBSIDIAN_INPUT_MODE_REQUIRED", 2, "", 0, 0)
		}
		return cfg, false, nil
	}
	if cfg.InputMode != obsidianInputMode {
		if hasObsidianOption {
			return cfg, true, newObsError("OBSIDIAN_INPUT_MODE_REQUIRED", 2, "", 0, 0)
		}
		return cfg, false, nil
	}
	if cfg.VaultRoot == "" {
		return cfg, true, newObsError("OBSIDIAN_INVALID_VAULT", 2, "", 0, 0)
	}
	if cfg.EntryPath == "" || !strings.HasSuffix(cfg.EntryPath, ".md") {
		return cfg, true, newObsError("OBSIDIAN_INVALID_ENTRY", 2, cfg.EntryPath, 0, 0)
	}
	if err := validateRelativePath(cfg.EntryPath, false); err != nil {
		return cfg, true, newObsError("OBSIDIAN_INVALID_ENTRY", 2, cfg.EntryPath, 0, 0)
	}
	cfg.OutputRoot = findBuilderOutputRoot(cfg.BuilderArgs)
	return cfg, true, nil
}

func splitOption(arg string) (string, string, bool) {
	if eq := strings.IndexByte(arg, '='); eq >= 0 {
		return arg[:eq], arg[eq+1:], true
	}
	return arg, "", false
}

func optionConsumesValue(name string) bool {
	switch name {
	case "--src", "--out", "--title", "--theme", "--base-dir", "--build-id", "--commit-sha", "--build-at", "--changed-manifest", "--format", "--markdown-extensions", "--heading-numbering", "--toc-depth", "--updated-at-source", "--meta", "--var", "--print-qr-url":
		return true
	default:
		return false
	}
}

func findBuilderOutputRoot(args []string) string {
	for i := 0; i < len(args); i++ {
		name, value, hasValue := splitOption(args[i])
		if name == "--out" {
			if hasValue {
				return value
			}
			if i+1 < len(args) {
				return args[i+1]
			}
		}
		if !hasValue && optionConsumesValue(name) {
			i++
		}
	}
	return ""
}

func replaceBuilderSource(args []string, src string) []string {
	out := []string{}
	for i := 0; i < len(args); i++ {
		name, _, hasValue := splitOption(args[i])
		if name == "--src" {
			if !hasValue && i+1 < len(args) {
				i++
			}
			continue
		}
		out = append(out, args[i])
		if !hasValue && optionConsumesValue(name) && i+1 < len(args) {
			i++
			out = append(out, args[i])
		}
	}
	return append([]string{"--src", src}, out...)
}

func validateVaultRoot(path string) *obsidianError {
	info, err := os.Lstat(path)
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return newObsError("OBSIDIAN_INVALID_VAULT", 2, "", 0, 0)
	}
	return nil
}

func validateRelativePath(value string, allowStar bool) error {
	if value == "" || !utf8.ValidString(value) || strings.HasPrefix(value, "/") || strings.HasSuffix(value, "/") || strings.HasPrefix(value, "~") || strings.Contains(value, "\\") || strings.ContainsAny(value, "\x00\r\n") {
		return fmt.Errorf("invalid relative path")
	}
	if len([]byte(value)) > 4096 || (len(value) >= 2 && value[1] == ':') {
		return fmt.Errorf("invalid relative path")
	}
	for _, segment := range strings.Split(value, "/") {
		if segment == "" || segment == "." || segment == ".." || len([]byte(segment)) > 255 {
			return fmt.Errorf("invalid relative path")
		}
		if segment == "*" && allowStar {
			continue
		}
		if strings.Contains(segment, "*") || strings.ContainsAny(segment, "?[]{}") {
			return fmt.Errorf("invalid pattern")
		}
	}
	return nil
}

func readFilter(path string) (obsidianFilter, *obsidianError) {
	var filter obsidianFilter
	if path == "" {
		return filter, nil
	}
	data, err := os.ReadFile(path)
	if err != nil || !utf8.Valid(data) {
		return filter, newObsError("OBSIDIAN_INVALID_FILTER", 2, "", 0, 0)
	}
	if hasDuplicateJSONKey(data) {
		return filter, newObsError("OBSIDIAN_INVALID_FILTER", 2, "", 0, 0)
	}
	var raw map[string][]string
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&raw); err != nil {
		return filter, newObsError("OBSIDIAN_INVALID_FILTER", 2, "", 0, 0)
	}
	for key, values := range raw {
		dedup := dedupSorted(values)
		for _, value := range dedup {
			if value == "" {
				return filter, newObsError("OBSIDIAN_INVALID_FILTER", 2, "", 0, 0)
			}
			if strings.Contains(value, "**") || validateRelativePath(value, true) != nil {
				return filter, newObsError("OBSIDIAN_INVALID_FILTER_PATTERN", 2, value, 0, 0)
			}
		}
		switch key {
		case "include_notes":
			filter.IncludeNotes = dedup
		case "exclude_notes":
			filter.ExcludeNotes = dedup
		case "include_assets":
			filter.IncludeAssets = dedup
		case "exclude_assets":
			filter.ExcludeAssets = dedup
		default:
			return filter, newObsError("OBSIDIAN_INVALID_FILTER", 2, key, 0, 0)
		}
	}
	return filter, nil
}

func hasDuplicateJSONKey(data []byte) bool {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil || tok != json.Delim('{') {
		return false
	}
	seen := map[string]bool{}
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return false
		}
		key, ok := keyTok.(string)
		if !ok {
			return false
		}
		if seen[key] {
			return true
		}
		seen[key] = true
		var skip any
		if err := dec.Decode(&skip); err != nil {
			return false
		}
	}
	return false
}

func dedupSorted(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func filterMatch(pattern, rel string) bool {
	pp := strings.Split(pattern, "/")
	rp := strings.Split(rel, "/")
	if len(pp) != len(rp) {
		return false
	}
	for i := range pp {
		if pp[i] != "*" && pp[i] != rp[i] {
			return false
		}
	}
	return true
}

func anyFilterMatch(patterns []string, rel string) bool {
	for _, pattern := range patterns {
		if filterMatch(pattern, rel) {
			return true
		}
	}
	return false
}

func newObsError(code string, exit int, target string, line int, column int) *obsidianError {
	return &obsidianError{Code: code, Exit: exit, Target: target, Line: line, Column: column}
}
