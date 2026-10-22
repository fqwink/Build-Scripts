package setup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func cicdConfigTestInput(role string) map[string]any {
	input := map[string]any{"schema_version": 1, "profile": "cicd-v1", "role": role, "id": "test", "allowed_projects": []any{}, "allowed_environments": []any{}, "user": "adlaire-cicd", "state_dir": "/var/lib/adlaire-ci/" + role + "/test", "release_id": "V.1.0", "listen": "127.0.0.1:9321", "tls_cert_file": nil, "tls_key_file": nil, "allow_http": false, "worker_slots": 1, "toolchains": []any{}, "ssh_authorized_key_file": nil}
	if role != "coordinator" {
		input["allowed_projects"] = []any{"project"}
		input["ssh_authorized_key_file"] = "/etc/adlaire-ci/authorized_keys"
		input["listen"] = ""
	}
	if role == "site" {
		input["listen"] = "127.0.0.1:8443"
		input["tls_cert_file"], input["tls_key_file"] = "/etc/adlaire-ci/cert", "/etc/adlaire-ci/key"
		input["allowed_environments"] = []any{map[string]any{"project_id": "project", "environment_id": "prod", "host": "example.test", "port": 443}}
	}
	return input
}

func cicdConfigTestJSON(input any) []byte {
	data, _ := json.Marshal(input)
	return data
}

func TestCICDConfig(t *testing.T) {
	for _, role := range []string{"coordinator", "worker", "site"} {
		t.Run(role, func(t *testing.T) {
			input := cicdConfigTestInput(role)
			input["toolchains"] = []any{map[string]any{"capability": "go", "version": "1.22.12", "root": "/opt/toolchains/go", "files": []any{map[string]any{"path": "bin/go", "sha256": strings.Repeat("a", 64), "executable": true}}}}
			data := cicdConfigTestJSON(input)
			before := append([]byte(nil), data...)
			cfg, err := ParseCICDConfig(data)
			if err != nil || cfg.Role != role || cfg.StateDir != input["state_dir"] || cfg.Toolchains[0].Files[0].Path != "bin/go" || !bytes.Equal(data, before) {
				t.Fatal("setup.phase19-config.positive.value")
			}
		})
	}
	for key := range cicdConfigTestInput("coordinator") {
		t.Run("missing-"+key, func(t *testing.T) {
			input := cicdConfigTestInput("coordinator")
			delete(input, key)
			cfg, err := ParseCICDConfig(cicdConfigTestJSON(input))
			if err == nil || err.Error() != "P19_SCHEMA" || !reflect.DeepEqual(cfg, CICDConfig{}) {
				t.Fatal("setup.phase19-config.missing.reject")
			}
		})
		if key == "tls_cert_file" || key == "tls_key_file" || key == "ssh_authorized_key_file" {
			continue
		}
		t.Run("null-"+key, func(t *testing.T) {
			input := cicdConfigTestInput("coordinator")
			input[key] = nil
			if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err == nil {
				t.Fatal("setup.phase19-config.null.reject")
			}
		})
	}
	for _, tc := range []struct {
		id, role, key string
		value         any
	}{
		{"unknown", "coordinator", "unknown", true}, {"schema", "coordinator", "schema_version", 2}, {"profile", "coordinator", "profile", "legacy-v1"}, {"role", "coordinator", "role", "unknown"}, {"id", "coordinator", "id", "../test"}, {"user", "coordinator", "user", "ROOT"}, {"relative", "coordinator", "state_dir", "relative"}, {"dot", "coordinator", "state_dir", "/var/lib/adlaire-ci/coordinator/test/../test"}, {"foreign-state", "coordinator", "state_dir", "/var/lib/adlaire-ci/worker/test"}, {"prefix", "coordinator", "state_dir", "/var/lib/adlaire-ci/coordinator/test-other"}, {"release", "coordinator", "release_id", "latest"}, {"public", "coordinator", "listen", "0.0.0.0:9321"}, {"port", "coordinator", "listen", "127.0.0.1:8765"}, {"slots", "coordinator", "worker_slots", 2}, {"slot-zero", "worker", "worker_slots", 0}, {"slot-max", "worker", "worker_slots", 5}, {"projects", "coordinator", "allowed_projects", []any{"p"}}, {"empty-worker", "worker", "allowed_projects", []any{}}, {"duplicate", "worker", "allowed_projects", []any{"p", "p"}}, {"worker-key", "worker", "ssh_authorized_key_file", nil}, {"coordinator-key", "coordinator", "ssh_authorized_key_file", "/etc/key"}, {"http", "coordinator", "allow_http", true}, {"tls", "coordinator", "tls_key_file", "/etc/key"}, {"site-cert", "site", "tls_cert_file", nil}, {"site-low-port", "site", "listen", "127.0.0.1:443"}, {"site-port-leading-zero", "site", "listen", "127.0.0.1:08443"}, {"site-host", "site", "listen", "Example.test:8443"}, {"site-binding", "site", "allowed_environments", []any{}},
	} {
		t.Run(tc.id, func(t *testing.T) {
			input := cicdConfigTestInput(tc.role)
			input[tc.key] = tc.value
			cfg, err := ParseCICDConfig(cicdConfigTestJSON(input))
			if err == nil || err.Error() != "P19_SCHEMA" || !reflect.DeepEqual(cfg, CICDConfig{}) {
				t.Fatal("setup.phase19-config.boundary.reject")
			}
		})
	}
	base := string(cicdConfigTestJSON(cicdConfigTestInput("coordinator")))
	for id, data := range map[string]string{"duplicate": strings.Replace(base, `"schema_version":1`, `"schema_version":1,"schema_version":1`, 1), "escaped-duplicate": strings.Replace(base, `"schema_version":1`, `"schema_version":1,"schema_vers\u0069on":1`, 1), "exponent": strings.Replace(base, `"schema_version":1`, `"schema_version":1e0`, 1), "trailing": base + `{}`, "invalid-utf8": strings.Replace(base, "test", "\xff", 1), "surrogate": strings.Replace(base, "test", `\ud800`, 1), "oversize": strings.Repeat(" ", 1048577)} {
		t.Run(id, func(t *testing.T) {
			if _, err := ParseCICDConfig([]byte(data)); err == nil {
				t.Fatal("setup.phase19-config.json.reject")
			}
		})
	}
}

func TestCICDConfigNested(t *testing.T) {
	for _, tc := range []struct {
		id     string
		change func(map[string]any)
	}{
		{"binding-unknown", func(b map[string]any) { b["extra"] = 1 }}, {"binding-missing", func(b map[string]any) { delete(b, "port") }}, {"binding-null", func(b map[string]any) { b["host"] = nil }}, {"binding-project", func(b map[string]any) { b["project_id"] = "foreign" }}, {"binding-host", func(b map[string]any) { b["host"] = "Example.test" }}, {"binding-port", func(b map[string]any) { b["port"] = 65536 }},
	} {
		t.Run(tc.id, func(t *testing.T) {
			input := cicdConfigTestInput("site")
			tc.change(input["allowed_environments"].([]any)[0].(map[string]any))
			if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err == nil {
				t.Fatal("setup.phase19-config.binding.reject")
			}
		})
	}
	for _, change := range []func(map[string]any){
		func(tool map[string]any) { tool["unknown"] = 1 }, func(tool map[string]any) { delete(tool, "version") }, func(tool map[string]any) { tool["capability"] = "shell" }, func(tool map[string]any) { tool["version"] = "latest version" }, func(tool map[string]any) { tool["root"] = "relative" }, func(tool map[string]any) { tool["files"] = []any{} }, func(tool map[string]any) { tool["files"].([]any)[0].(map[string]any)["path"] = "../go" }, func(tool map[string]any) {
			tool["files"].([]any)[0].(map[string]any)["sha256"] = strings.Repeat("A", 64)
		}, func(tool map[string]any) { delete(tool["files"].([]any)[0].(map[string]any), "executable") }, func(tool map[string]any) { tool["files"] = append(tool["files"].([]any), tool["files"].([]any)[0]) },
	} {
		input := cicdConfigTestInput("worker")
		tool := map[string]any{"capability": "go", "version": "1.22.12", "root": "/opt/go", "files": []any{map[string]any{"path": "bin/go", "sha256": strings.Repeat("a", 64), "executable": true}}}
		change(tool)
		input["toolchains"] = []any{tool}
		if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err == nil {
			t.Fatal("setup.phase19-config.toolchain.reject")
		}
	}
	for _, host := range []string{"example.test", "127.0.0.1", "::1", "2001:db8::1"} {
		if !cicdConfigHost(host) {
			t.Fatal("setup.phase19-config.host.accept")
		}
	}
	for _, host := range []string{"", "Example.test", "example.test.", "-a.test", "a..test", "127.0.0.999", "2001:DB8::1", "[::1]", "::1%lo"} {
		if cicdConfigHost(host) {
			t.Fatal("setup.phase19-config.host.reject")
		}
	}
}

func TestCICDConfigLimits(t *testing.T) {
	data := cicdConfigTestJSON(cicdConfigTestInput("coordinator"))
	data = append(data, bytes.Repeat([]byte{' '}, 1048576-len(data))...)
	if _, err := ParseCICDConfig(data); err != nil {
		t.Fatal("setup.phase19-config.bytes.maximum")
	}
	if _, err := ParseCICDConfig(append(data, ' ')); err == nil {
		t.Fatal("setup.phase19-config.bytes.overflow")
	}
	input := cicdConfigTestInput("worker")
	projects := []any{}
	for i := 0; i < 64; i++ {
		projects = append(projects, fmt.Sprintf("p%d", i))
	}
	input["allowed_projects"], input["worker_slots"] = projects, 4
	if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err != nil {
		t.Fatal("setup.phase19-config.projects.maximum")
	}
	input["allowed_projects"] = append(projects, "overflow")
	if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err == nil {
		t.Fatal("setup.phase19-config.projects.overflow")
	}
	input = cicdConfigTestInput("site")
	bindings := []any{}
	for i := 0; i < 1024; i++ {
		bindings = append(bindings, map[string]any{"project_id": "project", "environment_id": fmt.Sprintf("e%d", i), "host": "example.test", "port": 1024 + i})
	}
	input["allowed_environments"] = bindings
	if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err != nil {
		t.Fatal("setup.phase19-config.environments.maximum")
	}
	for _, binding := range []map[string]any{{"project_id": "project", "environment_id": "overflow", "host": "example.test", "port": 3000}, {"project_id": "project", "environment_id": "other", "host": "example.test", "port": 1024}, {"project_id": "project", "environment_id": "e0", "host": "other.test", "port": 3000}} {
		input["allowed_environments"] = append(append([]any{}, bindings...), binding)
		if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err == nil {
			t.Fatal("setup.phase19-config.environments.overflow")
		}
		input["allowed_environments"] = []any{bindings[0], binding}
		_, err := ParseCICDConfig(cicdConfigTestJSON(input))
		if binding["environment_id"] != "overflow" && err == nil {
			t.Fatal("setup.phase19-config.environments.duplicate")
		}
	}
	input = cicdConfigTestInput("coordinator")
	tools := []any{}
	for _, capability := range []string{"builder", "go", "deno"} {
		tools = append(tools, map[string]any{"capability": capability, "version": strings.Repeat("v", 128), "root": "/opt/tool", "files": []any{map[string]any{"path": "bin/tool", "sha256": strings.Repeat("a", 64), "executable": true}}})
	}
	input["toolchains"] = tools
	if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err != nil {
		t.Fatal("setup.phase19-config.tools.maximum")
	}
	input["toolchains"] = append(tools, tools[0])
	if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err == nil {
		t.Fatal("setup.phase19-config.tools.overflow")
	}
	input["toolchains"] = []any{tools[0], tools[0]}
	if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err == nil {
		t.Fatal("setup.phase19-config.tools.duplicate")
	}
	input["toolchains"] = tools[:1]
	tools[0].(map[string]any)["version"] = strings.Repeat("v", 129)
	if _, err := ParseCICDConfig(cicdConfigTestJSON(input)); err == nil {
		t.Fatal("setup.phase19-config.version.overflow")
	}
}

func TestCICDConfigStamp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")
	if os.WriteFile(path, []byte(`{}`), 0640) != nil || os.Chtimes(path, time.Unix(100, 0), time.Unix(100, 0)) != nil {
		t.Fatal("setup.phase19-config.stamp.fixture")
	}
	before, _ := os.Stat(path)
	if !sameCICDConfigFile(before, before, true) {
		t.Fatal("setup.phase19-config.stamp.stable")
	}
	if os.Chtimes(path, time.Unix(200, 0), time.Unix(200, 0)) != nil {
		t.Fatal("setup.phase19-config.stamp.change")
	}
	after, _ := os.Stat(path)
	if sameCICDConfigFile(before, after, true) || !sameCICDConfigFile(before, after, false) {
		t.Fatal("setup.phase19-config.stamp.content-change")
	}
	if os.Chmod(path, 0600) != nil {
		t.Fatal("setup.phase19-config.stamp.chmod")
	}
	after, _ = os.Stat(path)
	if sameCICDConfigFile(before, after, false) {
		t.Fatal("setup.phase19-config.stamp.permission-change")
	}
	if os.Rename(path, path+".old") != nil || os.WriteFile(path, []byte(`{}`), 0640) != nil {
		t.Fatal("setup.phase19-config.stamp.replace")
	}
	after, _ = os.Stat(path)
	if sameCICDConfigFile(before, after, true) {
		t.Fatal("setup.phase19-config.stamp.inode-change")
	}
}

func TestCICDConfigAccount(t *testing.T) {
	cfg, _ := ParseCICDConfig(cicdConfigTestJSON(cicdConfigTestInput("coordinator")))
	line := "adlaire-cicd:x:1001:1001::" + cfg.StateDir + ":/usr/sbin/nologin\n"
	if !cicdConfigAccount([]byte(line), cfg, 1001, 1001) {
		t.Fatal("setup.phase19-config.account.accept")
	}
	for _, value := range []string{line + line, strings.Replace(line, "1001:1001", "0:1001", 1), strings.Replace(line, "1001:1001", "1001:1002", 1), strings.Replace(line, cfg.StateDir, "/other", 1), strings.Replace(line, "/usr/sbin/nologin", "/bin/sh", 1), strings.Replace(line, "adlaire-cicd", "other", 1), line + strings.Replace(line, "adlaire-cicd", "alias", 1), "bad\n"} {
		if cicdConfigAccount([]byte(value), cfg, 1001, 1001) {
			t.Fatal("setup.phase19-config.account.reject")
		}
	}
	if cicdConfigAccount([]byte(line), cfg, 0, 1001) {
		t.Fatal("setup.phase19-config.account.root-reject")
	}
	if cicdConfigAccount([]byte(strings.Replace(line, ":1001::", ":0::", 1)), cfg, 1001, 0) {
		t.Fatal("setup.phase19-config.account.root-group-reject")
	}
	if !cicdConfigAccount([]byte(line+"other:x:1002:1001::/other:/bin/sh\n"), cfg, 1001, 1001) {
		t.Fatal("setup.phase19-config.account.read-group-share")
	}
}

func TestCICDConfigRead(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		if _, err := ReadCICDConfig("/etc/adlaire-ci/cicd/coordinator-test.json"); err == nil || err.Error() != "P19_CAPABILITY" {
			t.Fatal("setup.phase19-config.platform.reject")
		}
		t.Skip("P19_CAPABILITY: FD config fixture requires Linux amd64")
	}
	for _, fault := range []string{"none", "supplementary", "shared", "no-membership", "file-mode", "parent-mode", "parent-readable", "symlink", "parent-symlink", "hardlink", "directory", "oversize", "account", "filename", "owner", "group"} {
		t.Run(fault, func(t *testing.T) {
			root := t.TempDir()
			for _, path := range []string{"etc", "etc/adlaire-ci", "etc/adlaire-ci/cicd"} {
				mode := os.FileMode(0750)
				if path == "etc" {
					mode = 0755
				}
				if os.Mkdir(filepath.Join(root, path), mode) != nil {
					t.Fatal("setup.phase19-config.read.fixture")
				}
			}
			input := cicdConfigTestInput("coordinator")
			path := "/etc/adlaire-ci/cicd/coordinator-test.json"
			file := filepath.Join(root, path)
			data := cicdConfigTestJSON(input)
			uid, gid := uint32(os.Geteuid()), uint32(os.Getegid())
			if gid == 0 {
				gid = 1001
			}
			configGroup := gid
			var memberships []uint32
			if fault == "supplementary" {
				memberships = []uint32{gid}
				gid++
			}
			if fault == "no-membership" {
				gid++
			}
			serviceUID := uid
			if serviceUID == 0 {
				serviceUID = 1001
			}
			passwd := "adlaire-cicd:x:" + strconv.FormatUint(uint64(serviceUID), 10) + ":" + strconv.FormatUint(uint64(gid), 10) + "::" + input["state_dir"].(string) + ":/usr/sbin/nologin\n"
			if fault == "shared" {
				passwd += "other:x:" + strconv.FormatUint(uint64(serviceUID+1), 10) + ":" + strconv.FormatUint(uint64(gid), 10) + "::/other:/bin/sh\n"
			}
			if fault == "account" {
				passwd = strings.Replace(passwd, "nologin", "shell", 1)
			}
			if fault == "oversize" {
				data = bytes.Repeat([]byte{' '}, 1048577)
			}
			if os.WriteFile(filepath.Join(root, "etc/passwd"), []byte(passwd), 0644) != nil || os.WriteFile(file, data, 0640) != nil {
				t.Fatal("setup.phase19-config.read.fixture-write")
			}
			if os.Getegid() == 0 {
				for _, name := range []string{"etc/adlaire-ci", "etc/adlaire-ci/cicd", "etc/adlaire-ci/cicd/coordinator-test.json"} {
					if os.Chown(filepath.Join(root, name), int(uid), int(configGroup)) != nil {
						t.Fatal("setup.phase19-config.read.fixture-group")
					}
				}
			}
			switch fault {
			case "file-mode":
				_ = os.Chmod(file, 0644)
			case "parent-mode":
				_ = os.Chmod(filepath.Dir(file), 0770)
			case "parent-readable":
				_ = os.Chmod(filepath.Dir(file), 0755)
			case "symlink":
				_ = os.Rename(file, file+".real")
				_ = os.Symlink(file+".real", file)
			case "parent-symlink":
				_ = os.Rename(filepath.Dir(file), filepath.Dir(file)+".real")
				_ = os.Symlink(filepath.Dir(file)+".real", filepath.Dir(file))
			case "hardlink":
				_ = os.Link(file, file+".link")
			case "directory":
				_ = os.Remove(file)
				_ = os.Mkdir(file, 0640)
			case "filename":
				path = "/etc/adlaire-ci/cicd/coordinator-other.json"
				_ = os.Rename(file, filepath.Join(root, path))
			case "owner":
				uid++
			case "group":
				gid++
			}
			cfg, err := readCICDConfigAt(root, path, uid, serviceUID, gid, memberships...)
			if fault == "none" || fault == "supplementary" || fault == "shared" {
				if err != nil || cfg.ID != "test" {
					t.Fatalf("setup.phase19-config.read.accept: %v", err)
				}
			} else if err == nil || !reflect.DeepEqual(cfg, CICDConfig{}) {
				t.Fatal("setup.phase19-config.read.reject")
			}
		})
	}
}

type fakeSetupOps struct {
	real setupRealOps

	versions       map[string]string
	downloaded     map[string][]byte
	systemctlCalls []string
	apiEnabled     bool
	unsafeArchive  bool

	failDownloadName   string
	failChecksum       bool
	failAPIRestartOnce bool
	apiRestartFailures int
}

func setupTempDir(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func TestSetupValidatePathRejectsSymlinkSegments(t *testing.T) {
	root := setupTempDir(t)
	target, link := filepath.Join(root, "target"), filepath.Join(root, "link")
	if err := os.Mkdir(target, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	ops := setupRealOps{}
	if err := ops.ValidatePath(filepath.Join(target, "missing")); err != nil {
		t.Fatalf("real temp path rejected: %v", err)
	}
	for _, path := range []string{link, filepath.Join(link, "missing")} {
		if err := ops.ValidatePath(path); err == nil || err.Error() != "symlink path segment" {
			t.Fatalf("symlink segment accepted: %v", err)
		}
	}
}

func newFakeSetupOps() *fakeSetupOps {
	return &fakeSetupOps{
		versions:   map[string]string{},
		downloaded: map[string][]byte{},
	}
}

func runSetupForTest(ops *fakeSetupOps, args ...string) (int, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runSetup(args, &stdout, &stderr, ops)
	return code, stdout.String(), stderr.String()
}

func (ops *fakeSetupOps) ValidatePath(path string) error {
	return ops.real.ValidatePath(path)
}

func (ops *fakeSetupOps) ReadSecretInput(path string) (string, error) {
	return ops.real.ReadSecretInput(path)
}

func (ops *fakeSetupOps) ValidateExistingToken(path string) error {
	return ops.real.ValidateExistingToken(path)
}

func (ops *fakeSetupOps) ValidateLastSHA(path string) error {
	return ops.real.ValidateLastSHA(path)
}

func (ops *fakeSetupOps) ValidateCredentials(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0600 {
		return errors.New("invalid credentials")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if string(data) != "credential\n" {
		return errors.New("invalid credentials")
	}
	return nil
}

func (ops *fakeSetupOps) EnsureDownloadDir(path, version string) error {
	return ops.real.EnsureDownloadDir(path, version)
}

func (ops *fakeSetupOps) DownloadAsset(cfg setupConfig, asset setupAsset) error {
	if ops.failDownloadName == asset.Name {
		return errors.New("download failed")
	}
	if err := os.MkdirAll(cfg.DownloadDir, 0755); err != nil {
		return err
	}
	var data []byte
	if asset.Name == "SHA256SUMS" {
		data = ops.checksumBytes()
	} else if asset.Name == "admin-ui.tar.gz" {
		if ops.unsafeArchive {
			data = setupTestTarGZ(tarFile{name: "../index.html", body: "bad"})
		} else {
			data = setupTestTarGZ(
				tarFile{name: "index.html", body: "<!doctype html><title>Adlaire CI</title>"},
				tarFile{name: "adlaire-ci-sdk.js", body: "globalThis.AdlaireCI = {};"},
			)
		}
	} else {
		data = []byte("binary " + asset.Name + " " + cfg.TargetVersion + "\n")
	}
	ops.downloaded[asset.Name] = data
	return os.WriteFile(filepath.Join(cfg.DownloadDir, asset.Name), data, 0644)
}

func (ops *fakeSetupOps) checksumBytes() []byte {
	names := make([]string, 0, len(ops.downloaded))
	for name := range ops.downloaded {
		names = append(names, name)
	}
	sort.Strings(names)
	var out strings.Builder
	for _, name := range names {
		sum := sha256.Sum256(ops.downloaded[name])
		fmt.Fprintf(&out, "%s  %s\n", hex.EncodeToString(sum[:]), name)
	}
	return []byte(out.String())
}

func (ops *fakeSetupOps) VerifyChecksums(downloadDir string, assets []setupAsset, osArch string) error {
	if ops.failChecksum {
		return errors.New("checksum failed")
	}
	return ops.real.VerifyChecksums(downloadDir, assets, osArch)
}

func (ops *fakeSetupOps) InstallBinary(src, dest, binaryName, version string) error {
	if err := ops.real.InstallBinaryForTest(src, dest); err != nil {
		return err
	}
	ops.versions[dest] = version
	return nil
}

func (ops *fakeSetupOps) VerifyBinary(path, binaryName, version string) error {
	if !ops.BinaryExists(path) {
		return errors.New("missing binary")
	}
	if ops.versions[path] != version {
		return fmt.Errorf("version mismatch: got %q want %q", ops.versions[path], version)
	}
	return nil
}

func (ops *fakeSetupOps) BinaryExists(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0
}

func (ops *fakeSetupOps) EnsureDir(path string, mode os.FileMode) error {
	return ops.real.EnsureDir(path, mode)
}

func (ops *fakeSetupOps) CreateFileExclusive(path string, data []byte, mode os.FileMode) (bool, error) {
	return ops.real.CreateFileExclusive(path, data, mode)
}

func (ops *fakeSetupOps) InstallAdminArchive(archivePath, installDir string) error {
	return ops.real.InstallAdminArchive(archivePath, installDir)
}

func (ops *fakeSetupOps) WriteRunnerUnits(cfg setupConfig) error {
	return nil
}

func (ops *fakeSetupOps) WriteAPIUnit(cfg setupConfig) error {
	return nil
}

func (ops *fakeSetupOps) InitCredentials(apiPath, installDir, passwordFile string) error {
	return os.WriteFile(filepath.Join(installDir, ".admin_credentials"), []byte("credential\n"), 0600)
}

func (ops *fakeSetupOps) NewBackupDir() (string, error) {
	return os.MkdirTemp("", "setup-test-backup-")
}

func (ops *fakeSetupOps) BackupFile(src, dest string) error {
	if err := ops.real.BackupFile(src, dest); err != nil {
		return err
	}
	ops.versions[dest] = ops.versions[src]
	return nil
}

func (ops *fakeSetupOps) RestoreFile(src, dest string) error {
	if err := ops.real.RestoreFile(src, dest); err != nil {
		return err
	}
	ops.versions[dest] = ops.versions[src]
	return nil
}

func (ops *fakeSetupOps) BackupDir(src, dest string) error {
	return ops.real.BackupDir(src, dest)
}

func (ops *fakeSetupOps) RestoreDir(src, dest string) error {
	return ops.real.RestoreDir(src, dest)
}

func (ops *fakeSetupOps) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (ops *fakeSetupOps) Systemctl(args ...string) setupCommandResult {
	call := strings.Join(args, " ")
	ops.systemctlCalls = append(ops.systemctlCalls, call)
	switch call {
	case "is-enabled adlaire-ci-api":
		if ops.apiEnabled {
			return setupCommandResult{Stdout: "enabled\n", ExitCode: 0}
		}
		return setupCommandResult{Stderr: "disabled\n", ExitCode: 1}
	case "is-active adlaire-ci.timer", "is-active adlaire-ci-api":
		return setupCommandResult{Stdout: "active\n", ExitCode: 0}
	case "cat adlaire-ci.service", "cat adlaire-ci.timer", "cat adlaire-ci-api.service":
		return setupCommandResult{Stdout: "unit\n", ExitCode: 0}
	case "restart adlaire-ci-api":
		if ops.failAPIRestartOnce && ops.apiRestartFailures == 0 {
			ops.apiRestartFailures++
			return setupCommandResult{Stderr: "failed\n", ExitCode: 1}
		}
		return setupCommandResult{ExitCode: 0}
	default:
		return setupCommandResult{ExitCode: 0}
	}
}

func (ops *fakeSetupOps) HealthCheck() error {
	return nil
}

func (setupRealOps) InstallBinaryForTest(src, dest string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("empty binary")
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(dest, data, 0755); err != nil {
		return err
	}
	return os.Chmod(dest, 0755)
}

type tarFile struct {
	name string
	body string
}

func setupTestTarGZ(files ...tarFile) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, file := range files {
		body := []byte(file.body)
		_ = tw.WriteHeader(&tar.Header{Name: file.name, Mode: 0644, Size: int64(len(body)), Typeflag: tar.TypeReg})
		_, _ = tw.Write(body)
	}
	_ = tw.Close()
	_ = gz.Close()
	return buf.Bytes()
}

func writeSetupSecret(t *testing.T, dir, name, value string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(value+"\n"), 0600); err != nil {
		t.Fatalf("write secret: %v", err)
	}
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatalf("chmod secret: %v", err)
	}
	return path
}

func seedSetupBinary(t *testing.T, ops *fakeSetupOps, dir, name, version string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	if err := os.WriteFile(path, []byte("old "+name+" "+version+"\n"), 0755); err != nil {
		t.Fatalf("write binary: %v", err)
	}
	if err := os.Chmod(path, 0755); err != nil {
		t.Fatalf("chmod binary: %v", err)
	}
	ops.versions[path] = version
}

func TestRunSetupHelpVersionPrecedence(t *testing.T) {
	ops := newFakeSetupOps()
	code, stdout, stderr := runSetupForTest(ops, "--help", "bad\narg")
	if code != 0 || stdout != "Usage: adlaire-ci-setup <install|install-api|install-obsidian|update> --target-version V.X.N [options] [--version] [--help]\n" || stderr != "" {
		t.Fatalf("help mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	code, stdout, stderr = runSetupForTest(ops, "--version", "install", "--target-version", "bad")
	if code != 0 || !strings.HasPrefix(stdout, "adlaire-ci-setup V.0.0-dev go=") || stderr != "" {
		t.Fatalf("version mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestRunSetupInstallSuccess(t *testing.T) {
	root := setupTempDir(t)
	installDir := filepath.Join(root, "install")
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	tokenFile := writeSetupSecret(t, root, "github_token", "github-token")
	ops := newFakeSetupOps()

	code, stdout, stderr := runSetupForTest(ops,
		"install",
		"--target-version", "V.1.1",
		"--github-token-file", tokenFile,
		"--install-dir", installDir,
		"--bin-dir", binDir,
		"--download-dir", downloadDir,
	)
	if code != 0 || stderr != "" {
		t.Fatalf("install failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	wantStdout := strings.Join([]string{
		"setup: start install validate",
		"setup: ok install validate",
		"setup: start install download",
		"setup: ok install download",
		"setup: start install checksum",
		"setup: ok install checksum",
		"setup: start install install-binaries",
		"setup: ok install install-binaries",
		"setup: start install initialize-secrets",
		"setup: ok install initialize-secrets",
		"setup: start install initialize-state",
		"setup: ok install initialize-state",
		"setup: start install configure-systemd",
		"setup: ok install configure-systemd",
		"setup: start install activate-services",
		"setup: ok install activate-services",
		"setup: start install verify",
		"setup: ok install verify",
		"setup: success install V.1.1",
		"",
	}, "\n")
	if stdout != wantStdout {
		t.Fatalf("stdout mismatch:\n got %q\nwant %q", stdout, wantStdout)
	}
	for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup"} {
		if ops.versions[filepath.Join(binDir, name)] != "V.1.1" {
			t.Fatalf("missing installed version for %s", name)
		}
	}
	if data, err := os.ReadFile(filepath.Join(installDir, ".github_token")); err != nil || string(data) != "github-token\n" {
		t.Fatalf("token mismatch: data=%q err=%v", string(data), err)
	}
	if data, err := os.ReadFile(filepath.Join(installDir, ".last_sha")); err != nil || string(data) != "{\"sha\":\"\"}\n" {
		t.Fatalf("last sha mismatch: data=%q err=%v", string(data), err)
	}
	if !containsSetupCall(ops.systemctlCalls, "enable --now adlaire-ci.timer") {
		t.Fatalf("missing timer enable call: %v", ops.systemctlCalls)
	}
}

func TestRunSetupInstallObsidianSuccess(t *testing.T) {
	root := setupTempDir(t)
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	ops := newFakeSetupOps()

	code, stdout, stderr := runSetupForTest(ops,
		"install-obsidian",
		"--target-version", "V.1.1",
		"--bin-dir", binDir,
		"--download-dir", downloadDir,
	)
	if code != 0 || stderr != "" {
		t.Fatalf("install-obsidian failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	wantStdout := strings.Join([]string{
		"setup: start install-obsidian validate",
		"setup: ok install-obsidian validate",
		"setup: start install-obsidian download",
		"setup: ok install-obsidian download",
		"setup: start install-obsidian checksum",
		"setup: ok install-obsidian checksum",
		"setup: start install-obsidian install-binaries",
		"setup: ok install-obsidian install-binaries",
		"setup: start install-obsidian verify",
		"setup: ok install-obsidian verify",
		"setup: success install-obsidian V.1.1",
		"",
	}, "\n")
	if stdout != wantStdout {
		t.Fatalf("stdout mismatch:\n got %q\nwant %q", stdout, wantStdout)
	}
	if ops.versions[filepath.Join(binDir, "adlaire-ci-obsidian")] != "V.1.1" {
		t.Fatalf("missing installed obsidian version")
	}
	for _, forbidden := range []string{"daemon-reload", "enable --now adlaire-ci.timer", "is-active adlaire-ci.timer"} {
		if containsSetupCall(ops.systemctlCalls, forbidden) {
			t.Fatalf("install-obsidian must not call systemd %q: %v", forbidden, ops.systemctlCalls)
		}
	}
}

func TestRunSetupFailureContracts(t *testing.T) {
	root := setupTempDir(t)
	installDir := filepath.Join(root, "install")
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	tokenFile := writeSetupSecret(t, root, "github_token", "github-token")

	t.Run("install api precondition", func(t *testing.T) {
		ops := newFakeSetupOps()
		passFile := writeSetupSecret(t, root, "admin_password_precondition", "Initial-password-123")
		code, stdout, stderr := runSetupForTest(ops,
			"install-api",
			"--target-version", "V.1.1",
			"--initial-password-file", passFile,
			"--install-dir", installDir,
			"--bin-dir", binDir,
			"--download-dir", downloadDir,
		)
		if code != 2 || stdout != "setup: start install-api validate\n" || stderr != "setup: error PRECONDITION_FAILED stage=validate rollback=none\n" {
			t.Fatalf("precondition mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})

	t.Run("download", func(t *testing.T) {
		ops := newFakeSetupOps()
		ops.failDownloadName = "adlaire-ci-runner-linux-amd64"
		code, stdout, stderr := runSetupForTest(ops,
			"install",
			"--target-version", "V.1.1",
			"--github-token-file", tokenFile,
			"--install-dir", filepath.Join(root, "install-download"),
			"--bin-dir", filepath.Join(root, "bin-download"),
			"--download-dir", filepath.Join(root, "download-download"),
		)
		if code != 3 || !strings.HasSuffix(stdout, "setup: start install download\n") || stderr != "setup: error DOWNLOAD_FAILED stage=download rollback=none\n" {
			t.Fatalf("download mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})

	t.Run("checksum", func(t *testing.T) {
		ops := newFakeSetupOps()
		ops.failChecksum = true
		code, stdout, stderr := runSetupForTest(ops,
			"install",
			"--target-version", "V.1.1",
			"--github-token-file", tokenFile,
			"--install-dir", filepath.Join(root, "install-checksum"),
			"--bin-dir", filepath.Join(root, "bin-checksum"),
			"--download-dir", filepath.Join(root, "download-checksum"),
		)
		if code != 1 || !strings.HasSuffix(stdout, "setup: start install checksum\n") || stderr != "setup: error CHECKSUM_FAILED stage=checksum rollback=none\n" {
			t.Fatalf("checksum mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})
}

func TestRunSetupInstallAPIUnsafeAdminArchivePreservesExistingAdmin(t *testing.T) {
	root := setupTempDir(t)
	installDir := filepath.Join(root, "install")
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	ops := newFakeSetupOps()
	for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup"} {
		seedSetupBinary(t, ops, binDir, name, "V.1.1")
	}
	adminDir := filepath.Join(installDir, "admin")
	if err := os.MkdirAll(adminDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "index.html"), []byte("old-index"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "adlaire-ci-sdk.js"), []byte("old-sdk"), 0644); err != nil {
		t.Fatal(err)
	}
	passFile := writeSetupSecret(t, root, "admin_password", "Initial-password-123")
	ops.unsafeArchive = true

	code, stdout, stderr := runSetupForTest(ops,
		"install-api",
		"--target-version", "V.1.1",
		"--initial-password-file", passFile,
		"--install-dir", installDir,
		"--bin-dir", binDir,
		"--download-dir", downloadDir,
	)
	if code != 1 || !strings.HasSuffix(stdout, "setup: start install-api install-admin\n") || stderr != "setup: error ADMIN_INSTALL_FAILED stage=install-admin rollback=none\n" {
		t.Fatalf("unsafe admin mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if data, err := os.ReadFile(filepath.Join(adminDir, "index.html")); err != nil || string(data) != "old-index" {
		t.Fatalf("admin index was not preserved: data=%q err=%v", string(data), err)
	}
}

func TestRunSetupUpdateAPIRollback(t *testing.T) {
	root := setupTempDir(t)
	installDir := filepath.Join(root, "install")
	binDir := filepath.Join(root, "bin")
	downloadDir := filepath.Join(root, "download")
	ops := newFakeSetupOps()
	ops.apiEnabled = true
	ops.failAPIRestartOnce = true
	for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup", "adlaire-ci-api", "adlaire-ci-admin"} {
		seedSetupBinary(t, ops, binDir, name, "V.1.0")
	}
	adminDir := filepath.Join(installDir, "admin")
	if err := os.MkdirAll(adminDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "index.html"), []byte("old-index"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adminDir, "adlaire-ci-sdk.js"), []byte("old-sdk"), 0644); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := runSetupForTest(ops,
		"update",
		"--target-version", "V.1.1",
		"--install-dir", installDir,
		"--bin-dir", binDir,
		"--download-dir", downloadDir,
	)
	if code != 1 || !strings.HasSuffix(stdout, "setup: start update restart-services\n") || stderr != "setup: error SYSTEMD_FAILED stage=restart-services rollback=completed\n" {
		t.Fatalf("rollback mismatch: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	for _, name := range []string{"adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-setup", "adlaire-ci-api", "adlaire-ci-admin"} {
		if got := ops.versions[filepath.Join(binDir, name)]; got != "V.1.0" {
			t.Fatalf("binary %s not rolled back: %s", name, got)
		}
	}
	if data, err := os.ReadFile(filepath.Join(adminDir, "index.html")); err != nil || string(data) != "old-index" {
		t.Fatalf("admin index not rolled back: data=%q err=%v", string(data), err)
	}
}

func TestPhase18MaterializeAdminAssets(t *testing.T) {
	originalIndex := append([]byte(nil), setupEmbeddedAdminIndex...)
	originalSDK := append([]byte(nil), setupEmbeddedAdminSDK...)
	t.Cleanup(func() { SetEmbeddedAdminAssets(originalIndex, originalSDK) })
	SetEmbeddedAdminAssets([]byte("<main>phase18</main>\n"), []byte("export const phase18 = true;\n"))

	installDir := filepath.Join(setupTempDir(t), "runtime")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunSetup([]string{"phase18-materialize-admin", "--install-dir", installDir}, &stdout, &stderr)
	if code != 0 || stdout.String() != "setup: success phase18-materialize-admin\n" || stderr.Len() != 0 {
		t.Fatalf("materialize result code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for name, want := range map[string]string{
		"index.html":        "<main>phase18</main>\n",
		"adlaire-ci-sdk.js": "export const phase18 = true;\n",
	} {
		path := filepath.Join(installDir, "admin", name)
		data, err := os.ReadFile(path)
		if err != nil || string(data) != want {
			t.Fatalf("asset %s mismatch data=%q err=%v", name, string(data), err)
		}
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != 0644 {
			t.Fatalf("asset %s mode mismatch info=%v err=%v", name, info, err)
		}
	}
}

func TestPhase18MaterializeAdminRejectsUnsafeTarget(t *testing.T) {
	originalIndex := append([]byte(nil), setupEmbeddedAdminIndex...)
	originalSDK := append([]byte(nil), setupEmbeddedAdminSDK...)
	t.Cleanup(func() { SetEmbeddedAdminAssets(originalIndex, originalSDK) })
	SetEmbeddedAdminAssets([]byte("index"), []byte("sdk"))

	root := setupTempDir(t)
	target := filepath.Join(root, "target")
	if err := os.Symlink(filepath.Join(root, "missing"), target); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunSetup([]string{"phase18-materialize-admin", "--install-dir", target}, &stdout, &stderr)
	if code != 1 || stdout.Len() != 0 || stderr.String() != "admin asset materialization failed\n" {
		t.Fatalf("unsafe target result code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestSetupAtomicWriteTextSyncsFileAndParent(t *testing.T) {
	dir := setupTempDir(t)
	target := filepath.Join(dir, "asset.txt")
	ops := defaultSetupAtomicWriteOps()
	syncTargets := []string{}
	ops.sync = func(file *os.File) error {
		syncTargets = append(syncTargets, file.Name())
		return file.Sync()
	}
	if err := atomicWriteTextWithOps(target, "new-value", 0644, ops); err != nil {
		t.Fatalf("atomic write failed: %v", err)
	}
	if len(syncTargets) != 2 || filepath.Dir(syncTargets[0]) != dir || syncTargets[1] != dir {
		t.Fatalf("file and parent directory must be synced in order: %#v", syncTargets)
	}
	data, err := os.ReadFile(target)
	if err != nil || string(data) != "new-value" {
		t.Fatalf("unexpected target data=%q err=%v", string(data), err)
	}
	info, err := os.Stat(target)
	if err != nil || info.Mode().Perm() != 0644 {
		t.Fatalf("unexpected target mode info=%v err=%v", info, err)
	}
}

func TestSetupAtomicWriteTextFailureKeepsExistingTarget(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*setupAtomicWriteOps)
	}{
		{
			name: "file-sync",
			mutate: func(ops *setupAtomicWriteOps) {
				ops.sync = func(*os.File) error { return errors.New("sync failed") }
			},
		},
		{
			name: "rename",
			mutate: func(ops *setupAtomicWriteOps) {
				ops.rename = func(string, string) error { return errors.New("rename failed") }
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir := setupTempDir(t)
			target := filepath.Join(dir, "asset.txt")
			if err := os.WriteFile(target, []byte("old-value"), 0644); err != nil {
				t.Fatal(err)
			}
			ops := defaultSetupAtomicWriteOps()
			test.mutate(&ops)
			if err := atomicWriteTextWithOps(target, "new-value", 0644, ops); err == nil {
				t.Fatal("expected atomic write failure")
			}
			data, err := os.ReadFile(target)
			if err != nil || string(data) != "old-value" {
				t.Fatalf("existing target changed data=%q err=%v", string(data), err)
			}
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), ".asset.txt.tmp-") {
					t.Fatalf("temporary file was not removed: %s", entry.Name())
				}
			}
		})
	}
}

func TestPhase18MaterializeAdminRejectsAssetSymlink(t *testing.T) {
	originalIndex := append([]byte(nil), setupEmbeddedAdminIndex...)
	originalSDK := append([]byte(nil), setupEmbeddedAdminSDK...)
	t.Cleanup(func() { SetEmbeddedAdminAssets(originalIndex, originalSDK) })
	SetEmbeddedAdminAssets([]byte("index"), []byte("sdk"))

	installDir := filepath.Join(setupTempDir(t), "runtime")
	adminDir := filepath.Join(installDir, "admin")
	if err := os.MkdirAll(adminDir, 0755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(setupTempDir(t), "outside")
	if err := os.WriteFile(outside, []byte("unchanged"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(adminDir, "index.html")); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunSetup([]string{"phase18-materialize-admin", "--install-dir", installDir}, &stdout, &stderr)
	if code != 1 || stdout.Len() != 0 || stderr.String() != "admin asset materialization failed\n" {
		t.Fatalf("asset symlink result code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	data, err := os.ReadFile(outside)
	if err != nil || string(data) != "unchanged" {
		t.Fatalf("symlink target changed data=%q err=%v", string(data), err)
	}
}

func TestSetupInitCredentialsAcceptsSpecifiedChildOutput(t *testing.T) {
	root := setupTempDir(t)
	apiPath := filepath.Join(root, "api")
	passwordPath := filepath.Join(root, "password")
	if err := os.WriteFile(apiPath, []byte("#!/bin/sh\nprintf 'credentials initialized\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(passwordPath, []byte("password123\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := (setupRealOps{}).InitCredentials(apiPath, root, passwordPath); err != nil {
		t.Fatalf("specified init output rejected: %v", err)
	}
	if err := os.WriteFile(apiPath, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := (setupRealOps{}).InitCredentials(apiPath, root, passwordPath); err == nil || err.Error() != "credentials command output mismatch" {
		t.Fatalf("unexpected mismatch result: %v", err)
	}
}

func TestSetupFixturesExist(t *testing.T) {
	fixtures := []string{
		"success-setup-admin-release-asset-layout",
		"failure-setup-download-boundary",
		"failure-setup-api-version-cohort",
		"security-setup-admin-archive-boundary",
		"partial-setup-systemd-rollback-boundary",
		"partial-setup-api-runner-dispatch",
		"security-setup-secret-preservation",
	}
	required := []string{
		"manifest.json",
		"input/cli.json",
		"expected/stdout.txt",
		"expected/stderr.txt",
		"expected/effects.json",
		"expected/security.json",
	}
	for _, fixture := range fixtures {
		root := filepath.Join("..", "..", "testdata", "setup", fixture)
		for _, rel := range required {
			path := filepath.Join(root, rel)
			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("missing fixture file %s: %v", path, err)
			}
			if info.IsDir() {
				t.Fatalf("fixture path is a directory: %s", path)
			}
			if strings.HasSuffix(rel, ".json") {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read fixture json %s: %v", path, err)
				}
				var decoded any
				if err := json.Unmarshal(data, &decoded); err != nil {
					t.Fatalf("invalid fixture json %s: %v", path, err)
				}
			}
		}
	}
}

func containsSetupCall(calls []string, want string) bool {
	for _, call := range calls {
		if call == want {
			return true
		}
	}
	return false
}
