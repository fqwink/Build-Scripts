package setup

import (
	"encoding/json"
	"errors"
	"net"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/fqwink/build-scripts/components/statefile"
)

var cicdConfigSchema = errors.New("P19_SCHEMA")
var cicdConfigCapability = errors.New("P19_CAPABILITY")
var cicdToolVersion = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+\-]{0,127}$`)

func ParseCICDConfig(data []byte) (CICDConfig, error) {
	var empty CICDConfig
	if len(data) > 1048576 {
		return empty, cicdConfigSchema
	}
	canonical, err := statefile.CanonicalCICDJSON(data)
	if err != nil {
		return empty, cicdConfigSchema
	}
	var shape any
	var cfg CICDConfig
	if json.Unmarshal(canonical, &shape) != nil || !cicdConfigShape(shape, reflect.TypeOf(cfg)) || json.Unmarshal(canonical, &cfg) != nil || !validCICDConfig(cfg) {
		return empty, cicdConfigSchema
	}
	return cfg, nil
}

// The wire structs define the closed key set, including explicit nullable fields.
func cicdConfigShape(value any, kind reflect.Type) bool {
	if kind.Kind() == reflect.Pointer {
		return value == nil || cicdConfigShape(value, kind.Elem())
	}
	if value == nil {
		return false
	}
	switch kind.Kind() {
	case reflect.Struct:
		object, ok := value.(map[string]any)
		if !ok || len(object) != kind.NumField() {
			return false
		}
		for i := 0; i < kind.NumField(); i++ {
			field := kind.Field(i)
			child, exists := object[field.Tag.Get("json")]
			if !exists || !cicdConfigShape(child, field.Type) {
				return false
			}
		}
	case reflect.Slice:
		array, ok := value.([]any)
		if !ok {
			return false
		}
		for _, child := range array {
			if !cicdConfigShape(child, kind.Elem()) {
				return false
			}
		}
	}
	return true
}

func cicdConfigAtom(kind, value string) bool {
	data, _ := json.Marshal(value)
	return statefile.ValidateCICDValue(kind, data) == nil
}

func cicdConfigPath(value string) bool {
	if !filepath.IsAbs(value) || filepath.Clean(value) != value || len(value) > 4096 || strings.Contains(value, "\\") {
		return false
	}
	for _, c := range value {
		if c < 32 || c == 127 {
			return false
		}
	}
	return true
}

func cicdConfigHost(host string) bool {
	if ip := net.ParseIP(host); ip != nil {
		return ip.String() == host
	}
	if host == "" || len(host) > 253 || strings.ToLower(host) != host || strings.Trim(host, "0123456789.") == "" {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, c := range label {
			if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-') {
				return false
			}
		}
	}
	return true
}

func validCICDConfig(cfg CICDConfig) bool {
	if cfg.SchemaVersion != 1 || cfg.Profile != "cicd-v1" || !cicdConfigAtom("Id", cfg.ID) || !cicdConfigAtom("Id", cfg.User) || !cicdConfigPath(cfg.StateDir) || !setupVersionPattern.MatchString(cfg.ReleaseID) || len(cfg.ReleaseID) > 128 {
		return false
	}
	prefix := "/var/lib/adlaire-ci/" + cfg.Role + "/" + cfg.ID
	if cfg.StateDir != prefix && !strings.HasPrefix(cfg.StateDir, prefix+"/") || cfg.WorkerSlots < 1 || cfg.WorkerSlots > 4 || cfg.Role != "worker" && cfg.WorkerSlots != 1 {
		return false
	}
	for _, path := range []*string{cfg.TLSCertFile, cfg.TLSKeyFile, cfg.SSHAuthorizedKeyFile} {
		if path != nil && !cicdConfigPath(*path) {
			return false
		}
	}
	projects := map[string]bool{}
	if len(cfg.AllowedProjects) > 64 || len(cfg.AllowedEnvironments) > 1024 || len(cfg.Toolchains) > 3 {
		return false
	}
	for _, project := range cfg.AllowedProjects {
		if !cicdConfigAtom("Id", project) || projects[project] {
			return false
		}
		projects[project] = true
	}
	bindings, addresses := map[[2]string]bool{}, map[string]bool{}
	for _, binding := range cfg.AllowedEnvironments {
		key := [2]string{binding.ProjectID, binding.EnvironmentID}
		address := net.JoinHostPort(binding.Host, strconv.Itoa(binding.Port))
		if !projects[binding.ProjectID] || !cicdConfigAtom("Id", binding.EnvironmentID) || !cicdConfigHost(binding.Host) || binding.Port < 1 || binding.Port > 65535 || bindings[key] || addresses[address] {
			return false
		}
		bindings[key], addresses[address] = true, true
	}
	switch cfg.Role {
	case "coordinator":
		if len(projects) != 0 || len(bindings) != 0 || cfg.Listen != "127.0.0.1:9321" || cfg.SSHAuthorizedKeyFile != nil {
			return false
		}
	case "worker":
		if len(projects) == 0 || len(bindings) != 0 || cfg.Listen != "" || cfg.SSHAuthorizedKeyFile == nil {
			return false
		}
	case "site":
		host, port, err := net.SplitHostPort(cfg.Listen)
		number, numberErr := strconv.Atoi(port)
		if len(projects) == 0 || len(bindings) == 0 || cfg.SSHAuthorizedKeyFile == nil || cfg.TLSCertFile == nil || cfg.TLSKeyFile == nil || err != nil || numberErr != nil || strconv.Itoa(number) != port || number < 1024 || number > 65535 || !cicdConfigHost(host) || cfg.AllowHTTP && host != "127.0.0.1" {
			return false
		}
	default:
		return false
	}
	if cfg.Role != "site" && (cfg.TLSCertFile != nil || cfg.TLSKeyFile != nil || cfg.AllowHTTP) {
		return false
	}
	capabilities := map[string]bool{}
	for _, tool := range cfg.Toolchains {
		if tool.Capability != "builder" && tool.Capability != "go" && tool.Capability != "deno" || capabilities[tool.Capability] || !cicdToolVersion.MatchString(tool.Version) || !cicdConfigPath(tool.Root) || len(tool.Files) < 1 || len(tool.Files) > 32768 {
			return false
		}
		capabilities[tool.Capability] = true
		paths := map[string]bool{}
		for _, file := range tool.Files {
			if !cicdConfigAtom("RelPath", file.Path) || !cicdConfigAtom("Digest", file.SHA256) || paths[file.Path] {
				return false
			}
			paths[file.Path] = true
		}
	}
	return true
}

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == "setup" && validateExactOwnerFiles(contract.Files, []string{"setup.go", "model.go", "validate.go", "execute.go", "setup_test.go"})
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
