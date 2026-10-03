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

func parseSyncArgs(args []string) (obsidianSyncConfig, *obsidianError) {
	cfg := obsidianSyncConfig{
		DeletePolicy: "reject",
		ConflictDir:  "adlaire-ci-conflicts",
		TombstoneDir: "adlaire-ci-tombstones",
	}
	if len(args) < 2 || args[0] != "sync" {
		return cfg, newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
	}
	switch args[1] {
	case "plan", "apply", "rollback":
		cfg.Action = args[1]
	default:
		return cfg, newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
	}
	seen := map[string]bool{}
	for i := 2; i < len(args); {
		name := args[i]
		if !syncKnownOption(cfg.Action, name) || strings.Contains(name, "=") {
			return cfg, newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
		}
		if seen[name] {
			return cfg, newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
		}
		seen[name] = true
		if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
			return cfg, newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
		}
		value := args[i+1]
		switch name {
		case "--vault":
			cfg.VaultRoot = value
		case "--project-root":
			cfg.ProjectRoot = value
		case "--state-dir":
			cfg.StateDir = value
		case "--plan-file":
			cfg.PlanFile = filepath.ToSlash(value)
		case "--plan-hash":
			cfg.PlanHash = value
		case "--rollback-file":
			cfg.RollbackFile = filepath.ToSlash(value)
		case "--direction":
			cfg.Direction = value
		case "--delete-policy":
			cfg.DeletePolicy = value
		case "--conflict-dir":
			cfg.ConflictDir = normalizeSyncDirOption(value)
		case "--tombstone-dir":
			cfg.TombstoneDir = normalizeSyncDirOption(value)
		case "--open-uri":
			cfg.OpenURISet = true
			switch value {
			case "true":
				cfg.OpenURI = true
			case "false":
				cfg.OpenURI = false
			default:
				return cfg, newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
			}
		}
		i += 2
	}
	if err := validateSyncRequired(cfg); err != nil {
		return cfg, err
	}
	if err := validateSyncEnums(cfg); err != nil {
		return cfg, err
	}
	if err := validateSyncPaths(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func syncKnownOption(action, name string) bool {
	switch action {
	case "plan":
		switch name {
		case "--vault", "--project-root", "--state-dir", "--plan-file", "--direction", "--delete-policy", "--conflict-dir", "--tombstone-dir":
			return true
		}
	case "apply":
		switch name {
		case "--vault", "--project-root", "--state-dir", "--plan-file", "--plan-hash", "--rollback-file", "--open-uri":
			return true
		}
	case "rollback":
		switch name {
		case "--vault", "--project-root", "--state-dir", "--rollback-file", "--open-uri":
			return true
		}
	}
	return false
}

func validateSyncRequired(cfg obsidianSyncConfig) *obsidianError {
	if cfg.VaultRoot == "" || cfg.ProjectRoot == "" || cfg.StateDir == "" {
		return newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
	}
	switch cfg.Action {
	case "plan":
		if cfg.PlanFile == "" || cfg.Direction == "" {
			return newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
		}
	case "apply":
		if cfg.PlanFile == "" || cfg.PlanHash == "" {
			return newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
		}
	case "rollback":
		if cfg.RollbackFile == "" {
			return newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
		}
	}
	return nil
}

func validateSyncEnums(cfg obsidianSyncConfig) *obsidianError {
	if cfg.Action == "plan" {
		switch cfg.Direction {
		case "import-only", "export-only", "bidirectional":
		default:
			return newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
		}
		switch cfg.DeletePolicy {
		case "reject", "tombstone":
		default:
			return newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
		}
	}
	if cfg.PlanHash != "" && !isLowerHexDigest(cfg.PlanHash) {
		return newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0)
	}
	return nil
}

func validateSyncPaths(cfg obsidianSyncConfig) *obsidianError {
	for _, root := range []string{cfg.VaultRoot, cfg.ProjectRoot, cfg.StateDir} {
		if root == "" || !filepath.IsAbs(root) || !filepath.IsLocal(filepath.Base(root)) {
			return newObsError("OBSIDIAN_SYNC_PATH_INVALID", 2, "", 0, 0)
		}
	}
	if validateSyncRoot(cfg.VaultRoot) != nil || validateSyncRoot(cfg.ProjectRoot) != nil || validateSyncRoot(cfg.StateDir) != nil {
		return newObsError("OBSIDIAN_SYNC_PATH_INVALID", 2, "", 0, 0)
	}
	if cfg.PlanFile != "" && validateSyncStateFilePath(cfg.StateDir, cfg.PlanFile) != nil {
		return newObsError("OBSIDIAN_SYNC_PATH_INVALID", 2, "", 0, 0)
	}
	if cfg.RollbackFile != "" && validateSyncStateFilePath(cfg.StateDir, cfg.RollbackFile) != nil {
		return newObsError("OBSIDIAN_SYNC_PATH_INVALID", 2, "", 0, 0)
	}
	if cfg.Action == "plan" {
		if validateSyncProjectDirPath(cfg.ProjectRoot, cfg.ConflictDir) != nil || validateSyncProjectDirPath(cfg.ProjectRoot, cfg.TombstoneDir) != nil {
			return newObsError("OBSIDIAN_SYNC_PATH_INVALID", 2, "", 0, 0)
		}
	}
	return nil
}

func validateSyncRoot(root string) error {
	info, err := os.Lstat(root)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("invalid root")
	}
	return nil
}

func validateSyncStateFilePath(stateDir, rel string) error {
	rel = filepath.ToSlash(rel)
	if validateRelativePath(rel, false) != nil {
		return fmt.Errorf("invalid state path")
	}
	abs := filepath.Join(stateDir, filepath.FromSlash(rel))
	cleanAbs := filepath.Clean(abs)
	cleanRoot := filepath.Clean(stateDir)
	r, err := filepath.Rel(cleanRoot, cleanAbs)
	if err != nil || r == "." || strings.HasPrefix(filepath.ToSlash(r), "../") || filepath.IsAbs(r) {
		return fmt.Errorf("state escape")
	}
	if info, err := os.Lstat(cleanAbs); err == nil {
		if info.IsDir() || info.Mode()&os.ModeSymlink != 0 || hasMultipleHardlinks(info) || !info.Mode().IsRegular() {
			return fmt.Errorf("invalid existing state path")
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func validateSyncProjectDirPath(projectRoot, rel string) error {
	rel = normalizeSyncDirOption(rel)
	if validateRelativePath(rel, false) != nil {
		return fmt.Errorf("invalid project dir")
	}
	abs := filepath.Join(projectRoot, filepath.FromSlash(rel))
	cleanAbs := filepath.Clean(abs)
	cleanRoot := filepath.Clean(projectRoot)
	r, err := filepath.Rel(cleanRoot, cleanAbs)
	if err != nil || r == "." || strings.HasPrefix(filepath.ToSlash(r), "../") || filepath.IsAbs(r) {
		return fmt.Errorf("project escape")
	}
	if info, err := os.Lstat(cleanAbs); err == nil {
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return fmt.Errorf("invalid existing project dir")
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func normalizeSyncDirOption(value string) string {
	return strings.TrimSuffix(filepath.ToSlash(value), "/")
}

func isLowerHexDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, r := range value {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}

func exactArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
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
