package statefile

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

const cicdTestRoot = `{"schema_version":1,"generation":1,"projects":[],"definitions":[],"templates":[],"environments":[],"users":[{"id":"admin","name":"Administrator","revision":1,"enabled":true,"credential":null,"totp":null,"created_at":"2026-10-22T00:00:00Z"}],"services":[],"secrets":[],"grants":[{"subject_kind":"user","subject_id":"admin","project_id":null,"environment_id":null,"role":"platform-admin","enabled":true,"expires_at":null}],"workers":[],"runs":[],"deployments":[],"artifacts":[],"caches":[],"operations":[],"audits":[],"targets":[],"environment_states":[],"notifications":[],"trigger_cursors":[],"status_deliveries":[]}`
const cicdTestOperation = "o-0123456789abcdef0123456789abcdef"

func TestCICDStorageRegistry(t *testing.T) {
	data, err := os.ReadFile("../../docs/details/statefile.md")
	if err != nil {
		t.Fatal("statefile.phase19-schema.registry.fixture")
	}
	_, after, ok := strings.Cut(string(data), "```json\n")
	input, _, end := strings.Cut(after, "\n```")
	var registry struct {
		Types map[string]any `json:"types"`
	}
	if !ok || !end || json.Unmarshal([]byte(input), &registry) != nil || len(registry.Types) != 85 {
		t.Fatal("statefile.phase19-schema.registry.fixture")
	}
	actual, err := json.Marshal(cicdStorageTypes)
	var normalized map[string]any
	if err != nil || json.Unmarshal(actual, &normalized) != nil || !reflect.DeepEqual(registry.Types, normalized) {
		t.Fatal("statefile.phase19-schema.registry.drift")
	}
}

func TestCICDSchema(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		id, kind, input string
		valid           bool
	}{
		{"root", "Root", cicdTestRoot, true},
		{"head", "Head", `{"schema_version":1,"generation":1,"sha256":"` + strings.Repeat("a", 64) + `"}`, true},
		{"unknown-type", "Unknown", `{}`, false},
		{"primitive-name", "id", `"alpha"`, false},
		{"id", "Id", `"a_1-b"`, true},
		{"id-upper", "Id", `"Alpha"`, false},
		{"id-newline", "Id", `"a\n"`, false},
		{"job-id-limit", "JobId", `"` + strings.Repeat("a", 49) + `"`, false},
		{"digest-upper", "Digest", `"` + strings.Repeat("A", 64) + `"`, false},
		{"utc-leap", "Time", `"2024-02-29T12:34:56Z"`, true},
		{"utc-invalid-day", "Time", `"2023-02-29T12:34:56Z"`, false},
		{"utc-fraction", "Time", `"2024-02-29T12:34:56.0Z"`, false},
		{"utc-offset", "Time", `"2024-02-29T12:34:56+00:00"`, false},
		{"utc-before-epoch", "Time", `"1969-12-31T23:59:59Z"`, false},
		{"relative", "RelPath", `"src/main.go"`, true},
		{"dot-outside-task", "RelPath", `"."`, false},
		{"absolute-relative", "RelPath", `"/src"`, false},
		{"traversal", "RelPath", `"a/../b"`, false},
		{"colon", "RelPath", `"a:b"`, false},
		{"ref", "Ref", `"feature/topic"`, true},
		{"ref-dotdot", "Ref", `"a..b"`, false},
		{"ref-lock", "Ref", `"a.lock/b"`, false},
		{"ref-control", "Ref", `"a\u007f"`, false},
		{"glob", "Glob", `"src/**/test?.go"`, true},
		{"glob-embedded-star", "Glob", `"src/a**b"`, false},
		{"glob-class", "Glob", `"src/[ab]"`, false},
		{"url", "URL", `"https://example.test:8443/"`, true},
		{"url-user", "URL", `"https://user@example.test/"`, false},
		{"url-fragment", "URL", `"https://example.test/#"`, false},
		{"base64", "Bytes", `"Zg=="`, true},
		{"base64-bits", "Bytes", `"Zh=="`, false},
		{"base64-unpadded", "Bytes", `"Zg"`, false},
		{"base64-url", "Bytes", `"_w=="`, false},
		{"scalar-false", "Scalar", `false`, true},
		{"scalar-int", "Scalar", `-2147483648`, true},
		{"scalar-large", "Scalar", `2147483648`, false},
		{"scalar-null", "Scalar", `null`, false},
		{"scalar-exponent", "Scalar", `1e0`, false},
		{"scalar-negative-zero", "Scalar", `-0`, false},
		{"scalar-nul", "Scalar", `"a\u0000b"`, false},
		{"scalar-surrogate", "Scalar", `"\ud800"`, false},
		{"head-version", "Head", `{"schema_version":2,"generation":1,"sha256":"` + strings.Repeat("a", 64) + `"}`, false},
		{"head-unknown", "Head", `{"schema_version":1,"generation":1,"sha256":"` + strings.Repeat("a", 64) + `","x":0}`, false},
		{"head-missing", "Head", `{"schema_version":1,"generation":1}`, false},
		{"head-missing-substituted", "Head", `{"schema_version":1,"generation":1,"other":"` + strings.Repeat("a", 64) + `"}`, false},
		{"head-null", "Head", `{"schema_version":1,"generation":null,"sha256":"` + strings.Repeat("a", 64) + `"}`, false},
		{"head-duplicate", "Head", `{"schema_version":1,"generation":1,"generation":2,"sha256":"` + strings.Repeat("a", 64) + `"}`, false},
		{"task-dot", "Task", `{"kind":"builder","source":".","output":"out"}`, true},
		{"task-output-dot", "Task", `{"kind":"builder","source":"src","output":"."}`, false},
		{"go-packages", "Task", `{"kind":"go-test","packages":["./...","pkg/...","."],"race":false}`, true},
		{"go-race", "Task", `{"kind":"go-test","packages":["./..."],"race":true}`, false},
		{"task-unknown", "Task", `{"kind":"shell","command":"true"}`, false},
		{"resources-minimum", "Resources", `{"cpu_percent":100,"memory_mb":256,"disk_mb":128,"tasks":64,"timeout_seconds":900,"output_mb":16}`, true},
		{"resources-disk-headroom", "Resources", `{"cpu_percent":100,"memory_mb":256,"disk_mb":129,"tasks":64,"timeout_seconds":900,"output_mb":16}`, false},
		{"legacy-totp", "TOTP", `{"enabled":false,"secret_base32":null,"confirmed_at":null,"last_accepted_step":null}`, true},
		{"legacy-totp-lower", "TOTP", `{"enabled":false,"secret_base32":"jbswy3dpehpk3pxp","confirmed_at":null,"last_accepted_step":null}`, false},
	} {
		t.Run(test.id, func(t *testing.T) {
			input := []byte(test.input)
			before := bytes.Clone(input)
			err := ValidateCICDValue(test.kind, input)
			if (err == nil) != test.valid || err != nil && err != ErrCICDSchema || !bytes.Equal(input, before) {
				t.Fatalf("statefile.phase19-schema.%s.result", test.id)
			}
		})
	}
	for _, key := range []string{"PATH", "HOME", "TMPDIR", "LANG", "TZ", "LD_PRELOAD", "DYLD_LIBRARY_PATH", "GOCACHE", "DENO_DIR", "NODE_OPTIONS", "NPM_CONFIG", "lower", "A-B", "A\n"} {
		if validCICDEnvKey(key) {
			t.Fatal("statefile.phase19-schema.env-key.reject")
		}
	}
	if !validCICDEnvKey("BUILD_TARGET") {
		t.Fatal("statefile.phase19-schema.env-key.accept")
	}
}

func cicdRootFixture(t *testing.T) map[string]any {
	t.Helper()
	value, err := decodeCICDValue([]byte(cicdTestRoot))
	if err != nil {
		t.Fatal("statefile.phase19-snapshot.fixture")
	}
	return cicdObject(value)
}

func TestCICDRootIntegrity(t *testing.T) {
	for _, test := range []struct {
		id     string
		mutate func(map[string]any)
	}{
		{"duplicate-user", func(root map[string]any) { root["users"] = append(cicdList(root["users"]), cicdList(root["users"])[0]) }},
		{"missing-user", func(root map[string]any) { cicdObject(cicdList(root["grants"])[0])["subject_id"] = "missing" }},
		{"platform-scope", func(root map[string]any) { cicdObject(cicdList(root["grants"])[0])["environment_id"] = "prod" }},
		{"last-admin", func(root map[string]any) { cicdObject(cicdList(root["users"])[0])["enabled"] = false }},
		{"unknown-root", func(root map[string]any) { root["legacy"] = true }},
		{"null-list", func(root map[string]any) { root["runs"] = nil }},
		{"orphan-delivery", func(root map[string]any) {
			root["status_deliveries"] = []any{map[string]any{"project_id": "missing", "run_id": "r-" + strings.Repeat("a", 32), "attempt": json.Number("1"), "state": "pending", "sent": false, "code": nil, "at": "2026-10-22T00:00:00Z"}}
		}},
	} {
		t.Run(test.id, func(t *testing.T) {
			root := cicdRootFixture(t)
			if !validateCICDRoot(root) {
				t.Fatal("statefile.phase19-root.positive")
			}
			test.mutate(root)
			if validateCICDRoot(root) {
				t.Fatalf("statefile.phase19-root.%s.reject", test.id)
			}
		})
	}
}

func TestCICDReservedCapacity(t *testing.T) {
	root := cicdRootFixture(t)
	if !validCICDReservedCapacity(root, cicdSnapshotLimit) || validCICDReservedCapacity(root, cicdSnapshotLimit+1) {
		t.Fatal("statefile.phase19-snapshot.capacity.bytes")
	}
	root["operations"] = []any{map[string]any{"state": "accepted"}}
	if !validCICDReservedCapacity(root, cicdSnapshotLimit-262144) || validCICDReservedCapacity(root, cicdSnapshotLimit-262143) {
		t.Fatal("statefile.phase19-snapshot.capacity.accepted")
	}
	cicdObject(cicdList(root["operations"])[0])["state"] = "recovery_required"
	if validCICDReservedCapacity(root, cicdSnapshotLimit) {
		t.Fatal("statefile.phase19-snapshot.capacity.recovery")
	}
	cicdObject(cicdList(root["operations"])[0])["state"] = "succeeded"
	if !validCICDReservedCapacity(root, cicdSnapshotLimit) {
		t.Fatal("statefile.phase19-snapshot.capacity.terminal")
	}
}

func cicdLinuxStore(t *testing.T) *cicdSnapshotStore {
	t.Helper()
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		if _, err := newCICDSnapshotStore("/must-not-access"); err != ErrCICDCapability {
			t.Fatal("statefile.phase19-snapshot.platform.reject")
		}
		t.Skip("P19_CAPABILITY: snapshot filesystem evidence requires Linux amd64; schema evidence is portable")
	}
	path := t.TempDir()
	if os.Chmod(path, 0700) != nil {
		t.Fatal("statefile.phase19-snapshot.directory.fixture")
	}
	store, err := newCICDSnapshotStore(path)
	if err != nil {
		t.Fatal("statefile.phase19-snapshot.store.fixture")
	}
	if _, err := store.initialize([]byte(cicdTestRoot), cicdTestOperation); err != nil {
		t.Fatalf("statefile.phase19-snapshot.initialize: %v", err)
	}
	return store
}

func cicdRenameAdmin(name string) func(map[string]any) error {
	return func(root map[string]any) error { cicdObject(cicdList(root["users"])[0])["name"] = name; return nil }
}

func TestCICDSnapshotStore(t *testing.T) {
	s := cicdLinuxStore(t)
	if _, err := s.initialize([]byte(cicdTestRoot), cicdTestOperation); err != ErrCICDCorrupt {
		t.Fatal("statefile.phase19-snapshot.no-reinitialize")
	}
	for generation := int64(1); generation <= 5; generation++ {
		head, err := s.commit(generation, cicdTestOperation, cicdRenameAdmin(strconv.FormatInt(generation, 10)))
		if err != nil || head.Generation != generation+1 {
			t.Fatalf("statefile.phase19-snapshot.commit: %v", err)
		}
	}
	fresh, err := newCICDSnapshotStore(s.path)
	if err != nil {
		t.Fatal("statefile.phase19-snapshot.restart.fixture")
	}
	root, head, err := fresh.read()
	if err != nil || head.Generation != 6 || cicdObject(cicdList(root["users"])[0])["name"] != "5" {
		t.Fatal("statefile.phase19-snapshot.restart.result")
	}
	calls := 0
	if _, err := fresh.commit(5, cicdTestOperation, func(map[string]any) error { calls++; return nil }); err != ErrCICDRevision || calls != 0 {
		t.Fatal("statefile.phase19-snapshot.stale.callback")
	}
	if noOp, err := fresh.commit(6, cicdTestOperation, func(map[string]any) error { calls++; return nil }); err != nil || noOp != head || calls != 1 {
		t.Fatal("statefile.phase19-snapshot.noop.result")
	}
	before, err := os.ReadFile(filepath.Join(s.path, "cicd", "head.json"))
	if err != nil {
		t.Fatal("statefile.phase19-snapshot.head.fixture")
	}
	if _, err := fresh.commit(6, cicdTestOperation, func(root map[string]any) error { cicdRenameAdmin("not-committed")(root); return ErrCICDBusy }); err != ErrCICDBusy {
		t.Fatal("statefile.phase19-snapshot.callback.error")
	}
	after, err := os.ReadFile(filepath.Join(s.path, "cicd", "head.json"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("statefile.phase19-snapshot.callback.atomic")
	}
	root["generation"] = json.Number("99")
	_, observed, err := fresh.read()
	if err != nil || observed != head {
		t.Fatal("statefile.phase19-snapshot.copy.isolation")
	}
	names, err := os.ReadDir(filepath.Join(s.path, "cicd", "snapshots"))
	if err != nil || len(names) != 3 || names[0].Name() != "g-4.json" || names[1].Name() != "g-5.json" || names[2].Name() != "g-6.json" {
		t.Fatal("statefile.phase19-snapshot.retention")
	}
}

func TestCICDSnapshotFaults(t *testing.T) {
	for _, stage := range []string{"capacity", "snapshot.create", "snapshot.write", "snapshot.chmod", "snapshot.sync", "snapshot.close", "snapshot.publish", "snapshot.directory-sync", "head.create", "head.write", "head.chmod", "head.sync", "head.close", "head.publish", "head.directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			s := cicdLinuxStore(t)
			s.fault = func(actual string) error {
				if actual == stage {
					return errors.New("injected-private-detail")
				}
				return nil
			}
			head, err := s.commit(1, cicdTestOperation, cicdRenameAdmin("next"))
			want := ErrCICDStorage
			if stage == "head.directory-sync" {
				want = ErrCICDUncertain
			}
			if err != want || head != (cicdHead{}) {
				t.Fatalf("statefile.phase19-snapshot.%s.failure: %v", stage, err)
			}
			s.fault = nil
			root, current, err := s.read()
			expected := int64(1)
			if stage == "head.directory-sync" {
				expected = 2
			}
			if err != nil || current.Generation != expected {
				t.Fatalf("statefile.phase19-snapshot.%s.head", stage)
			}
			if expected == 1 && cicdObject(cicdList(root["users"])[0])["name"] != "Administrator" {
				t.Fatal("statefile.phase19-snapshot.failed-mutation.visible")
			}
			if stage == "head.directory-sync" {
				calls := 0
				if _, err := s.commit(expected, cicdTestOperation, func(map[string]any) error { calls++; return nil }); err != ErrCICDUncertain || calls != 0 {
					t.Fatal("statefile.phase19-snapshot.uncertain.stopped")
				}
				var reopenErr error
				s, reopenErr = newCICDSnapshotStore(s.path)
				if reopenErr != nil {
					t.Fatal("statefile.phase19-snapshot.uncertain.reopen")
				}
			}
			if _, err := s.commit(expected, cicdTestOperation, cicdRenameAdmin("retry")); err != nil {
				t.Fatalf("statefile.phase19-snapshot.%s.retry: %v", stage, err)
			}
		})
	}
	t.Run("cleanup", func(t *testing.T) {
		s := cicdLinuxStore(t)
		calls, diagnostics := 0, 0
		s.fault = func(stage string) error {
			if stage == "cleanup" {
				calls++
				if calls == 2 {
					return ErrCICDCleanup
				}
			}
			return nil
		}
		s.diagnostic = func(err error) {
			if err != ErrCICDCleanup {
				t.Fatal("statefile.phase19-snapshot.cleanup.code")
			}
			diagnostics++
		}
		head, err := s.commit(1, cicdTestOperation, cicdRenameAdmin("committed"))
		if err != nil || head.Generation != 2 || diagnostics != 1 {
			t.Fatal("statefile.phase19-snapshot.cleanup.committed")
		}
	})
}

func TestCICDSnapshotCorruption(t *testing.T) {
	for _, test := range []struct {
		id      string
		corrupt func(string) error
	}{
		{"bad-hash", func(path string) error {
			return os.WriteFile(filepath.Join(path, "snapshots", "g-1.json"), []byte(`{}`), 0600)
		}},
		{"valid-shape-bad-hash", func(path string) error {
			return os.WriteFile(filepath.Join(path, "snapshots", "g-1.json"), []byte(strings.Replace(cicdTestRoot, "Administrator", "Tampered", 1)), 0600)
		}},
		{"generation-mismatch", func(path string) error {
			data := []byte(strings.Replace(cicdTestRoot, `"generation":1`, `"generation":2`, 1))
			head, err := json.Marshal(cicdHead{SchemaVersion: 1, Generation: 1, SHA256: fmt.Sprintf("%x", sha256.Sum256(data))})
			if err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(path, "snapshots", "g-1.json"), data, 0600); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(path, "head.json"), head, 0600)
		}},
		{"missing-snapshot", func(path string) error { return os.Remove(filepath.Join(path, "snapshots", "g-1.json")) }},
		{"head-schema", func(path string) error {
			return os.WriteFile(filepath.Join(path, "head.json"), []byte(`{"schema_version":2}`), 0600)
		}},
		{"head-mode", func(path string) error { return os.Chmod(filepath.Join(path, "head.json"), 0644) }},
		{"snapshot-mode", func(path string) error { return os.Chmod(filepath.Join(path, "snapshots", "g-1.json"), 0644) }},
		{"head-symlink", func(path string) error {
			if err := os.Rename(filepath.Join(path, "head.json"), filepath.Join(path, "saved")); err != nil {
				return err
			}
			return os.Symlink("saved", filepath.Join(path, "head.json"))
		}},
		{"snapshot-hardlink", func(path string) error {
			return os.Link(filepath.Join(path, "snapshots", "g-1.json"), filepath.Join(path, "linked"))
		}},
		{"unknown-file", func(path string) error { return os.WriteFile(filepath.Join(path, "unexpected"), nil, 0600) }},
		{"lock-nonempty", func(path string) error { return os.WriteFile(filepath.Join(path, "state.lock"), []byte("x"), 0600) }},
		{"head-oversize", func(path string) error {
			return os.WriteFile(filepath.Join(path, "head.json"), bytes.Repeat([]byte(" "), 4097), 0600)
		}},
		{"head-fifo", func(path string) error {
			if err := os.Remove(filepath.Join(path, "head.json")); err != nil {
				return err
			}
			return syscall.Mkfifo(filepath.Join(path, "head.json"), 0600)
		}},
		{"snapshots-directory-mode", func(path string) error { return os.Chmod(filepath.Join(path, "snapshots"), 0755) }},
	} {
		t.Run(test.id, func(t *testing.T) {
			s := cicdLinuxStore(t)
			path := filepath.Join(s.path, "cicd")
			if test.corrupt(path) != nil {
				t.Fatal("statefile.phase19-snapshot.corruption.fixture")
			}
			var reopenErr error
			s, reopenErr = newCICDSnapshotStore(s.path)
			if reopenErr != nil {
				t.Fatal("statefile.phase19-snapshot.corruption.reopen")
			}
			root, head, err := s.read()
			if err != ErrCICDCorrupt || root != nil || head != (cicdHead{}) {
				t.Fatalf("statefile.phase19-snapshot.%s.reject", test.id)
			}
			calls := 0
			if _, err := s.commit(1, cicdTestOperation, func(map[string]any) error { calls++; return nil }); err != ErrCICDCorrupt || calls != 0 {
				t.Fatal("statefile.phase19-snapshot.corrupt.write.reject")
			}
		})
	}
}

func TestCICDSnapshotConcurrency(t *testing.T) {
	s := cicdLinuxStore(t)
	const count = 8
	var group sync.WaitGroup
	var ready sync.WaitGroup
	ready.Add(count)
	start := make(chan struct{})
	var conflicts atomic.Int64
	failures := make(chan error, count)
	for i := 0; i < count; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			_, first, err := s.read()
			ready.Done()
			<-start
			if err != nil {
				failures <- err
				return
			}
			generation := first.Generation
			for {
				_, err = s.commit(generation, cicdTestOperation, cicdRenameAdmin(strconv.FormatInt(generation, 10)))
				if err == ErrCICDRevision {
					conflicts.Add(1)
					_, head, readErr := s.read()
					if readErr != nil {
						failures <- readErr
						return
					}
					generation = head.Generation
					continue
				}
				if err != nil {
					failures <- err
				}
				return
			}
		}()
	}
	ready.Wait()
	close(start)
	group.Wait()
	close(failures)
	for err := range failures {
		t.Fatalf("statefile.phase19-snapshot.concurrent: %v", err)
	}
	_, head, err := s.read()
	if err != nil || head.Generation != count+1 || conflicts.Load() < count-1 {
		t.Fatal("statefile.phase19-snapshot.concurrent.no-lost-update")
	}
}

func TestCICDSnapshotCrashHelper(t *testing.T) {
	path := os.Getenv("CICD_TEST_STORE")
	if path == "" {
		return
	}
	s, err := newCICDSnapshotStore(path)
	if err != nil {
		os.Exit(20)
	}
	stage := os.Getenv("CICD_TEST_STAGE")
	s.fault = func(actual string) error {
		if actual == stage {
			_ = syscall.Kill(os.Getpid(), syscall.SIGKILL)
		}
		return nil
	}
	if _, err := s.commit(1, cicdTestOperation, cicdRenameAdmin("after-crash")); err != nil {
		os.Exit(21)
	}
	os.Exit(22)
}

func TestCICDSnapshotReadRaces(t *testing.T) {
	for _, advances := range []int{1, 2, 3} {
		t.Run(strconv.Itoa(advances), func(t *testing.T) {
			s := cicdLinuxStore(t)
			changes := 0
			s.fault = func(stage string) error {
				if stage != "read.snapshot" || changes >= advances {
					return nil
				}
				changes++
				root := cicdRootFixture(t)
				root["generation"] = json.Number(strconv.Itoa(changes + 1))
				data := cicdValueBytes(root)
				head := cicdHead{SchemaVersion: 1, Generation: int64(changes + 1), SHA256: fmt.Sprintf("%x", sha256.Sum256(data))}
				headBytes, err := json.Marshal(head)
				if err != nil || os.WriteFile(filepath.Join(s.path, "cicd", "snapshots", "g-"+strconv.Itoa(changes+1)+".json"), data, 0600) != nil || os.WriteFile(filepath.Join(s.path, "cicd", "head.json"), headBytes, 0600) != nil || os.Remove(filepath.Join(s.path, "cicd", "snapshots", "g-"+strconv.Itoa(changes)+".json")) != nil {
					t.Fatal("statefile.phase19-snapshot.read-race.fixture")
				}
				return nil
			}
			root, head, err := s.read()
			if advances == 3 {
				if err != ErrCICDBusy || root != nil || head != (cicdHead{}) {
					t.Fatal("statefile.phase19-snapshot.read-race.busy")
				}
			} else if err != nil || head.Generation != int64(advances+1) {
				t.Fatal("statefile.phase19-snapshot.read-race.retry")
			}
		})
	}
}

func TestCICDSnapshotLocks(t *testing.T) {
	s := cicdLinuxStore(t)
	path := filepath.Join(s.path, "cicd", "state.lock")
	lock, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatal("statefile.phase19-snapshot.flock.fixture")
	}
	defer lock.Close()
	if syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB) != nil {
		t.Fatal("statefile.phase19-snapshot.flock.acquire")
	}
	before, err := lock.Stat()
	if err != nil {
		t.Fatal("statefile.phase19-snapshot.flock.stat")
	}
	start := time.Now()
	calls := 0
	_, err = s.commit(1, cicdTestOperation, func(map[string]any) error { calls++; return nil })
	if err != ErrCICDBusy || calls != 0 || time.Since(start) < 1900*time.Millisecond {
		t.Fatal("statefile.phase19-snapshot.flock.deadline")
	}
	after, err := os.Stat(path)
	if err != nil || !os.SameFile(before, after) || after.Size() != 0 {
		t.Fatal("statefile.phase19-snapshot.flock.never-unlink")
	}
	if syscall.Flock(int(lock.Fd()), syscall.LOCK_UN) != nil {
		t.Fatal("statefile.phase19-snapshot.flock.release")
	}
	cicdSnapshotMutex.Lock()
	_, err = s.commit(1, cicdTestOperation, func(map[string]any) error { calls++; return nil })
	cicdSnapshotMutex.Unlock()
	if err != ErrCICDBusy || calls != 0 {
		t.Fatal("statefile.phase19-snapshot.mutex.deadline")
	}
	if _, err := s.commit(1, cicdTestOperation, cicdRenameAdmin("unlocked")); err != nil {
		t.Fatal("statefile.phase19-snapshot.lock.reusable")
	}
}

func TestCICDSnapshotNoReplace(t *testing.T) {
	s := cicdLinuxStore(t)
	dir, snapshots, unlock, err := s.lock(false)
	if err != nil {
		t.Fatal("statefile.phase19-snapshot.noreplace.fixture")
	}
	defer unlock()
	_ = dir
	before, err := readCICDBounded(snapshots, "g-1.json", cicdSnapshotLimit)
	if err != nil {
		t.Fatal("statefile.phase19-snapshot.noreplace.before")
	}
	if err := s.writeCICDFile(snapshots, "."+cicdTestOperation+".tmp", "g-1.json", []byte("not-the-snapshot"), "snapshot", false); err != ErrCICDStorage {
		t.Fatal("statefile.phase19-snapshot.noreplace.reject")
	}
	after, err := readCICDBounded(snapshots, "g-1.json", cicdSnapshotLimit)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("statefile.phase19-snapshot.noreplace.immutable")
	}
}

func TestCICDSnapshotMutations(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skip("P19_CAPABILITY: storage mutation execution requires Linux amd64")
	}
	const selection = "^TestCICD(Schema|RootIntegrity|Snapshot(Store|Faults|Corruption|Locks|NoReplace))$"
	for _, mutation := range []struct{ id, path, old, new, reason string }{
		{"unknown-key", "validate.go", "!ok || len(object) != len(schema)", "!ok", "statefile.phase19-schema.head-unknown.result"},
		{"missing-key", "validate.go", "!exists || !validateCICDType(childType, child, expression+\".\"+key, depth+1)", "exists && !validateCICDType(childType, child, expression+\".\"+key, depth+1)", "statefile.phase19-schema.head-missing-substituted.result"},
		{"duplicate-key", "execute.go", "if _, exists := object[key]; exists {", "if _, exists := object[key]; exists && false {", "statefile.phase19-schema.head-duplicate.result"},
		{"digest", "execute.go", "fmt.Sprintf(\"%x\", sha256.Sum256(data)) != head.SHA256", "false", "statefile.phase19-snapshot.valid-shape-bad-hash.reject"},
		{"generation", "execute.go", "if err != nil || !ok || !validateCICDRoot(root) || cicdInt(root[\"generation\"]) != head.Generation {", "if err != nil || !ok || !validateCICDRoot(root) {", "statefile.phase19-snapshot.generation-mismatch.reject"},
		{"cas", "execute.go", "head.Generation != expectedGeneration || head.Generation >= 9007199254740991", "head.Generation >= 9007199254740991", "statefile.phase19-snapshot.stale.callback"},
		{"noop", "execute.go", "if bytes.Equal(before, cicdValueBytes(root)) {", "if false && bytes.Equal(before, cicdValueBytes(root)) {", "statefile.phase19-snapshot.noop.result"},
		{"noreplace", "execute.go", "err = renameCICDNoReplace(dir, tmp, name)", "err = os.Rename(cicdChildPath(dir, tmp), cicdChildPath(dir, name))", "statefile.phase19-snapshot.noreplace.reject"},
		{"directory-sync", "execute.go", "s.inject(stage+\".directory-sync\") != nil || dir.Sync() != nil", "false", "statefile.phase19-snapshot.snapshot.directory-sync.failure"},
		{"uncertain", "execute.go", "return ErrCICDUncertain", "return ErrCICDStorage", "statefile.phase19-snapshot.head.directory-sync.failure"},
		{"retention", "execute.go", "n < generation-2", "n < generation-1", "statefile.phase19-snapshot.retention"},
		{"flock", "execute.go", "err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)", "err = nil", "statefile.phase19-snapshot.flock.deadline"},
	} {
		t.Run(mutation.id, func(t *testing.T) {
			copyDir := t.TempDir()
			for _, name := range []string{"statefile.go", "model.go", "validate.go", "execute.go", "statefile_test.go"} {
				data, err := os.ReadFile(name)
				if err != nil {
					t.Fatal("statefile.phase19-mutation.source.fixture")
				}
				if name == mutation.path {
					if strings.Count(string(data), mutation.old) != 1 {
						t.Fatal("statefile.phase19-mutation.target.fixture")
					}
					data = []byte(strings.Replace(string(data), mutation.old, mutation.new, 1))
				}
				if os.WriteFile(filepath.Join(copyDir, name), data, 0600) != nil {
					t.Fatal("statefile.phase19-mutation.copy.fixture")
				}
			}
			if os.WriteFile(filepath.Join(copyDir, "go.mod"), []byte("module storage-mutation\n\ngo 1.22\n"), 0600) != nil {
				t.Fatal("statefile.phase19-mutation.module.fixture")
			}
			compile := exec.Command("go", "test", "-c", "-o", "mutation.test", ".")
			compile.Dir = copyDir
			compile.Env = append(os.Environ(), "GOENV=off", "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off", "GOWORK=off")
			if _, err := compile.CombinedOutput(); err != nil {
				t.Fatal("statefile.phase19-mutation.compile.invalid")
			}
			command := exec.Command("./mutation.test", "-test.run="+selection)
			command.Dir = copyDir
			output, err := command.CombinedOutput()
			if err == nil || !bytes.Contains(output, []byte(mutation.reason)) {
				t.Fatal("statefile.phase19-mutation.expected-assertion.not-killed")
			}
		})
	}
}

func TestCICDSnapshotProcessCrash(t *testing.T) {
	for _, stage := range []string{"snapshot.sync", "snapshot.directory-sync", "head.sync", "head.directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			s := cicdLinuxStore(t)
			// This is the current native Go test binary, not a shell or runtime lookup.
			command := exec.Command("/proc/self/exe", "-test.run=^TestCICDSnapshotCrashHelper$")
			command.Env = []string{"CICD_TEST_STORE=" + s.path, "CICD_TEST_STAGE=" + stage}
			output, processErr := command.CombinedOutput()
			exitErr, exited := processErr.(*exec.ExitError)
			var status syscall.WaitStatus
			if exited {
				status, _ = exitErr.Sys().(syscall.WaitStatus)
			}
			if !exited || !status.Signaled() || status.Signal() != syscall.SIGKILL || len(output) != 0 {
				t.Fatal("statefile.phase19-snapshot.crash.kill")
			}
			fresh, err := newCICDSnapshotStore(s.path)
			if err != nil {
				t.Fatal("statefile.phase19-snapshot.crash.reopen")
			}
			_, head, err := fresh.read()
			expected := int64(1)
			if stage == "head.directory-sync" {
				expected = 2
			}
			if err != nil || head.Generation != expected {
				t.Fatal("statefile.phase19-snapshot.crash.reconcile")
			}
			if _, err := fresh.commit(expected, cicdTestOperation, cicdRenameAdmin("reconciled")); err != nil {
				t.Fatalf("statefile.phase19-snapshot.crash.retry: %v", err)
			}
		})
	}
}

func TestCICDCanonicalJSON(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		id, input, want string
	}{
		{"object-order", ` { "b":2,"a":1 } `, `{"a":1,"b":2}`},
		{"nested-order", `{"z":[{"b":false,"a":null},3,1],"a":{}}`, `{"a":{},"z":[{"a":null,"b":false},3,1]}`},
		{"array-order", `[3,1,2]`, `[3,1,2]`},
		{"empty-array", `[]`, `[]`},
		{"empty-object", `{}`, `{}`},
		{"null", `null`, `null`},
		{"true", `true`, `true`},
		{"false", `false`, `false`},
		{"zero", `0`, `0`},
		{"safe-upper", `9007199254740991`, `9007199254740991`},
		{"safe-lower", `-9007199254740991`, `-9007199254740991`},
		{"quoted-number", `"1e0"`, `"1e0"`},
		{"html-slash", `"\u003c\u0026\u003e\/"`, `"<&>/"`},
		{"quotes-backslash", `"\"\\"`, `"\"\\"`},
		{"line-separators", `"\u2028\u2029"`, "\"\u2028\u2029\""},
		{"surrogate-pair", `"\uD83D\uDE00"`, "\"\U0001f600\""},
		{"escaped-surrogate-text", `"\\ud800"`, `"\\ud800"`},
		{"no-normalization", `{"\u00e9":1,"e\u0301":2}`, "{\"e\u0301\":2,\"\u00e9\":1}"},
		{"utf8-key-order", `{"\uffff":2,"\ud800\udc00":3,"a":1}`, "{\"a\":1,\"\uffff\":2,\"\U00010000\":3}"},
	} {
		t.Run(test.id, func(t *testing.T) {
			input := []byte(test.input)
			before := bytes.Clone(input)
			got, err := CanonicalCICDJSON(input)
			if err != nil || !bytes.Equal(got, []byte(test.want)) || !bytes.Equal(input, before) {
				t.Fatalf("statefile.phase19-canonical.%s.bytes: canonical bytes or input immutability mismatch", test.id)
			}
			again, err := CanonicalCICDJSON(got)
			if err != nil || !bytes.Equal(again, got) {
				t.Fatalf("statefile.phase19-canonical.%s.bytes: canonical form not stable", test.id)
			}
		})
	}
	controls := []string{
		`\u0000`, `\u0001`, `\u0002`, `\u0003`, `\u0004`, `\u0005`, `\u0006`, `\u0007`,
		`\b`, `\t`, `\n`, `\u000b`, `\f`, `\r`, `\u000e`, `\u000f`,
		`\u0010`, `\u0011`, `\u0012`, `\u0013`, `\u0014`, `\u0015`, `\u0016`, `\u0017`,
		`\u0018`, `\u0019`, `\u001a`, `\u001b`, `\u001c`, `\u001d`, `\u001e`, `\u001f`,
	}
	for c, escaped := range controls {
		t.Run(fmt.Sprintf("control-%02x", c), func(t *testing.T) {
			input := []byte(fmt.Sprintf(`"\u%04X"`, c))
			got, err := CanonicalCICDJSON(input)
			if err != nil || string(got) != `"`+escaped+`"` {
				t.Fatalf("statefile.phase19-canonical.control-%02x.bytes: escape mismatch", c)
			}
		})
	}
	for _, test := range []struct{ id, input string }{
		{"empty", ""}, {"whitespace", " \n\t"}, {"truncated", `{"a":`},
		{"trailing-token", `{} []`}, {"trailing-garbage", `{}x`},
		{"duplicate", `{"a":1,"a":2}`}, {"escaped-duplicate", `{"a":1,"\u0061":2}`},
		{"null-duplicate", `{"a":null,"a":1}`},
		{"nested-duplicate", `[{"a":{"b":1,"b":2}}]`},
		{"float", `1.0`}, {"exponent", `1e0`}, {"negative-zero", `-0`},
		{"leading-zero", `01`}, {"positive-sign", `+1`},
		{"unsafe-upper", `9007199254740992`}, {"unsafe-lower", `-9007199254740992`},
		{"overflow", `9223372036854775808`}, {"underflow", `-9223372036854775809`},
		{"nested-float", `{"a":[1.5]}`},
		{"high-surrogate", `"\ud800"`}, {"low-surrogate", `"\udc00"`},
		{"wrong-pair", `"\ud800\u0041"`}, {"double-high", `"\ud800\ud800"`},
		{"truncated-pair", `"\ud800\u00"`}, {"surrogate-key", `{"\ud800":1}`},
		{"invalid-utf8", "\"\xff\""}, {"raw-control", "\"\x01\""},
		{"invalid-escape", `"\x01"`},
	} {
		t.Run(test.id, func(t *testing.T) {
			input := []byte(test.input)
			before := bytes.Clone(input)
			got, err := CanonicalCICDJSON(input)
			digest, digestErr := DigestCICDJSON(input)
			if got != nil || !errors.Is(err, ErrCICDSchema) || err.Error() != "P19_SCHEMA" ||
				digest != "" || !errors.Is(digestErr, ErrCICDSchema) || !bytes.Equal(input, before) {
				t.Fatalf("statefile.phase19-canonical.%s.reject: expected schema rejection without result or input mutation", test.id)
			}
		})
	}
	t.Run("digest", func(t *testing.T) {
		const want = "43258cff783fe7036d8a43033f830adfc60ec037382473548ac742b888292777"
		for _, input := range []string{`{"a":1,"b":2}`, " {\n\"b\":2,\"a\":1}\n"} {
			got, err := DigestCICDJSON([]byte(input))
			if err != nil || got != want {
				t.Fatal("statefile.phase19-canonical.fixed.digest: expected independent SHA256 over canonical bytes")
			}
		}
	})
}

func TestOwnerFileContract(t *testing.T) {
	if Owner() != "statefile" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executeOwnerFileContract(newOwnerFileContract()) {
		t.Fatalf("statefile owner file contract failed")
	}
}
