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

	"github.com/fqwink/build-scripts/components/security"
)

func cicdLogTestValue() map[string]any {
	return map[string]any{"schema_version": json.Number("1"), "project_id": "p", "run_id": "r-" + strings.Repeat("11", 16), "attempt": json.Number("1"), "job_id": "build", "stream": "stdout", "first_offset": json.Number("0"), "next_offset": json.Number("0"), "data": "", "truncated": false, "gap": false, "terminal": false}
}

func TestCICDLogSchema(t *testing.T) {
	for _, name := range []string{"initial", "nul", "utf8", "window", "length", "reverse", "truncated", "run-id", "unknown", "missing", "invalid-utf8", "duplicate", "unsafe-offset"} {
		t.Run(name, func(t *testing.T) {
			value := cicdLogTestValue()
			valid := false
			switch name {
			case "initial":
				valid = true
			case "nul":
				value["data"], value["next_offset"], valid = strings.Repeat("\x00", 1048576), json.Number("1048576"), true
			case "utf8":
				value["data"], value["next_offset"], valid = "\xe6\xac\xa1", json.Number("3"), true
			case "window":
				value["first_offset"], value["next_offset"], value["truncated"], valid = json.Number("5"), json.Number("5"), true, true
			case "length":
				value["data"] = "x"
			case "reverse":
				value["first_offset"] = json.Number("1")
				value["truncated"] = true
			case "truncated":
				value["truncated"] = true
			case "run-id":
				value["run_id"] = "r-bad"
			case "unknown":
				value["other"] = false
			case "missing":
				delete(value, "terminal")
			case "unsafe-offset":
				value["first_offset"], value["next_offset"], value["truncated"] = json.Number("9007199254740992"), json.Number("9007199254740992"), true
			}
			data := cicdValueBytes(value)
			if name == "invalid-utf8" {
				data = bytes.Replace(data, []byte(`"data":""`), []byte{'"', 'd', 'a', 't', 'a', '"', ':', '"', 0xff, '"'}, 1)
			}
			if name == "duplicate" {
				data = bytes.Replace(data, []byte(`"gap":false`), []byte(`"gap":false,"gap":false`), 1)
			}
			err := ValidateCICDLog(data)
			if valid && err != nil || !valid && err != ErrCICDSchema {
				t.Fatal("statefile.phase19-log.schema.reject", name, err)
			}
		})
	}
}

func cicdLogStoreFixture(t *testing.T) (*CICDLogStore, map[string]any) {
	t.Helper()
	s, request, plan := cicdRunFixture(t)
	if _, err := s.StartRun(request, 1, func(CICDPlanSnapshot) (json.RawMessage, error) { return plan, nil }); err != nil {
		t.Fatal("statefile.phase19-log.store.fixture", err)
	}
	logs, err := OpenCICDLogStore(s.store.path)
	if err != nil {
		t.Fatal("statefile.phase19-log.store.open", err)
	}
	return logs, cicdLogTestValue()
}

func TestCICDLogStore(t *testing.T) {
	s, value := cicdLogStoreFixture(t)
	id := cicdLogIdentity(value)
	before, head, _ := s.store.read()
	if data, err := s.ReadLog(id); data != nil || err != ErrCICDNotFound {
		t.Fatal("statefile.phase19-log.store.missing", err)
	}
	if _, err := os.Lstat(filepath.Join(s.store.path, "cicd", "logs")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("statefile.phase19-log.store.read-no-create")
	}
	if err := s.WriteLog(cicdValueBytes(value)); err != nil {
		t.Fatal("statefile.phase19-log.store.initial", err)
	}
	value["data"], value["next_offset"] = "hello\x00\xe6\xac\xa1", json.Number("9")
	if err := s.WriteLog(cicdValueBytes(value)); err != nil {
		t.Fatal("statefile.phase19-log.store.append", err)
	}
	value["data"], value["first_offset"], value["next_offset"], value["truncated"], value["gap"] = "lo\x00\xe6\xac\xa1!", json.Number("3"), json.Number("10"), true, true
	if err := s.WriteLog(cicdValueBytes(value)); err != nil {
		t.Fatal("statefile.phase19-log.store.window", err)
	}
	value["terminal"] = true
	if err := s.WriteLog(cicdValueBytes(value)); err != nil {
		t.Fatal("statefile.phase19-log.store.seal", err)
	}
	s.store.fault = func(stage string) error {
		if strings.HasPrefix(stage, "log.") {
			return errors.New("no IO on replay")
		}
		return nil
	}
	if err := s.WriteLog(cicdValueBytes(value)); err != nil {
		t.Fatal("statefile.phase19-log.store.noop", err)
	}
	s.store.fault = nil
	fresh, err := OpenCICDLogStore(s.store.path)
	if err != nil {
		t.Fatal("statefile.phase19-log.store.reopen", err)
	}
	data, err := fresh.ReadLog(id)
	if err != nil || !bytes.Equal(data, cicdValueBytes(value)) {
		t.Fatal("statefile.phase19-log.store.read", err)
	}
	data[0] = '!'
	if data, err := fresh.ReadLog(id); err != nil || data[0] != '{' {
		t.Fatal("statefile.phase19-log.store.copy", err)
	}
	after, latest, _ := s.store.read()
	if latest != head || !bytes.Equal(cicdValueBytes(before), cicdValueBytes(after)) {
		t.Fatal("statefile.phase19-log.store.no-root-write")
	}
	path := filepath.Join(s.store.path, "cicd", "logs", id.ProjectID, id.RunID, "1", id.JobID, id.Stream+".json")
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatal("statefile.phase19-log.store.mode")
	}
}

func TestCICDLogRejection(t *testing.T) {
	for _, name := range []string{"first-content", "first-terminal", "first-gap", "not-found", "unselected", "overlap", "reverse-first", "reverse-next", "gap-clear", "terminal-change", "terminal-run", "terminal-new"} {
		t.Run(name, func(t *testing.T) {
			s, value := cicdLogStoreFixture(t)
			var expected error = ErrCICDStale
			if strings.HasPrefix(name, "first-") {
				if name == "first-content" {
					value["data"], value["next_offset"] = "x", json.Number("1")
				}
				if name == "first-terminal" {
					value["terminal"] = true
				}
				if name == "first-gap" {
					value["gap"] = true
				}
			} else if name == "not-found" {
				value["attempt"], expected = json.Number("2"), ErrCICDNotFound
			} else if name == "unselected" {
				value["job_id"] = "ignored"
			} else {
				if err := s.WriteLog(cicdValueBytes(value)); err != nil {
					t.Fatal("statefile.phase19-log.reject.initial", err)
				}
				value["data"], value["next_offset"] = "abc", json.Number("3")
				if name == "gap-clear" {
					value["gap"] = true
				}
				if name == "terminal-change" {
					value["terminal"] = true
				}
				if err := s.WriteLog(cicdValueBytes(value)); err != nil {
					t.Fatal("statefile.phase19-log.reject.base", err)
				}
				switch name {
				case "overlap":
					value["data"] = "axc"
				case "reverse-first":
					value["first_offset"], value["data"], value["truncated"] = json.Number("1"), "bc", true
					if err := s.WriteLog(cicdValueBytes(value)); err != nil {
						t.Fatal("statefile.phase19-log.reject.window", err)
					}
					value["first_offset"], value["data"], value["truncated"] = json.Number("0"), "abc", false
				case "reverse-next":
					value["data"], value["next_offset"] = "ab", json.Number("2")
				case "gap-clear":
					value["gap"] = false
				case "terminal-change":
					value["terminal"] = false
				case "terminal-run", "terminal-new":
					_, head, _ := s.store.read()
					if _, err := s.store.commit(head.Generation, cicdTestOperation, func(root map[string]any) error {
						run := cicdObject(cicdList(root["runs"])[0])
						run["state"] = "failed"
						run["finished_at"] = "2026-10-22T00:00:00Z"
						return nil
					}); err != nil {
						t.Fatal("statefile.phase19-log.reject.root", err)
					}
					if name == "terminal-new" {
						value["stream"], value["data"], value["next_offset"] = "stderr", "", json.Number("0")
					}
				}
			}
			before, head, _ := s.store.read()
			if err := s.WriteLog(cicdValueBytes(value)); err != expected {
				t.Fatal("statefile.phase19-log.reject.result", name, err)
			}
			after, latest, _ := s.store.read()
			if head != latest || !bytes.Equal(cicdValueBytes(before), cicdValueBytes(after)) {
				t.Fatal("statefile.phase19-log.reject.no-root-write")
			}
		})
	}
}

func TestCICDLogMaximum(t *testing.T) {
	s, value := cicdLogStoreFixture(t)
	if err := s.WriteLog(cicdValueBytes(value)); err != nil {
		t.Fatal("statefile.phase19-log.maximum.initial", err)
	}
	value["data"], value["next_offset"] = strings.Repeat("\x00", 1048576), json.Number("1048576")
	data := cicdValueBytes(value)
	if len(data) <= 6291456 || len(data) > cicdLogLimit {
		t.Fatal("statefile.phase19-log.maximum.envelope")
	}
	if err := s.WriteLog(data); err != nil {
		t.Fatal("statefile.phase19-log.maximum.write", err)
	}
	actual, err := s.ReadLog(cicdLogIdentity(value))
	if err != nil || !bytes.Equal(actual, data) {
		t.Fatal("statefile.phase19-log.maximum.read", err)
	}
	if err := s.WriteLog(make([]byte, cicdLogLimit+1)); err != ErrCICDStorage {
		t.Fatal("statefile.phase19-log.maximum.over-cap", err)
	}
	var absent *CICDLogStore
	if _, err := absent.ReadLog(cicdLogIdentity(value)); err != ErrCICDSchema {
		t.Fatal("statefile.phase19-log.maximum.nil-read")
	}
	if absent.WriteLog(nil) != ErrCICDSchema {
		t.Fatal("statefile.phase19-log.maximum.nil-write")
	}
	bad := cicdLogIdentity(value)
	bad.JobID = "../build"
	if _, err := s.ReadLog(bad); err != ErrCICDSchema {
		t.Fatal("statefile.phase19-log.maximum.bad-identity")
	}
}

func TestCICDLogConcurrency(t *testing.T) {
	s, value := cicdLogStoreFixture(t)
	if err := s.WriteLog(cicdValueBytes(value)); err != nil {
		t.Fatal("statefile.phase19-log.concurrent.initial", err)
	}
	start := make(chan struct{})
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			fresh, err := OpenCICDLogStore(s.store.path)
			if err != nil {
				t.Error("statefile.phase19-log.concurrent.open", err)
				return
			}
			local := cicdLogTestValue()
			local["data"], local["next_offset"] = "same", json.Number("4")
			if err := fresh.WriteLog(cicdValueBytes(local)); err != nil {
				t.Error("statefile.phase19-log.concurrent.write", err)
				return
			}
			data, err := fresh.ReadLog(cicdLogIdentity(local))
			if err != nil || !bytes.Equal(data, cicdValueBytes(local)) {
				t.Error("statefile.phase19-log.concurrent.read", err)
			}
		}()
	}
	close(start)
	group.Wait()
	_, head, err := s.store.read()
	if err != nil || head.Generation != 3 {
		t.Fatal("statefile.phase19-log.concurrent.no-root-write", err)
	}
}

func TestCICDLogFaults(t *testing.T) {
	for _, stage := range []string{"log.create", "log.write", "log.chmod", "log.sync", "log.close", "log.publish", "log.directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			s, value := cicdLogStoreFixture(t)
			if err := s.WriteLog(cicdValueBytes(value)); err != nil {
				t.Fatal("statefile.phase19-log.fault.initial", err)
			}
			value["data"], value["next_offset"] = "x", json.Number("1")
			s.store.fault = func(current string) error {
				if current == stage {
					return errors.New("injected")
				}
				return nil
			}
			want := ErrCICDStorage
			if stage == "log.directory-sync" {
				want = ErrCICDUncertain
			}
			if err := s.WriteLog(cicdValueBytes(value)); err != want {
				t.Fatal("statefile.phase19-log.fault.result", err)
			}
			s.store.fault = nil
			data, err := s.ReadLog(cicdLogIdentity(value))
			if err != nil {
				t.Fatal("statefile.phase19-log.fault.read", err)
			}
			log, _ := decodeCICDLog(data)
			wantText := ""
			if stage == "log.directory-sync" {
				wantText = "x"
			}
			if log["data"] != wantText {
				t.Fatal("statefile.phase19-log.fault.atomic")
			}
			path := filepath.Join(s.store.path, "cicd", "logs", "p", value["run_id"].(string), "1", "build")
			entries, err := os.ReadDir(path)
			if err != nil || len(entries) != 1 || entries[0].Name() != "stdout.json" {
				t.Fatal("statefile.phase19-log.fault.no-temp")
			}
		})
	}
}

func TestCICDLogCorruption(t *testing.T) {
	for _, name := range []string{"identity", "malformed", "symlink", "hardlink", "file-mode", "directory-mode", "oversized"} {
		t.Run(name, func(t *testing.T) {
			s, value := cicdLogStoreFixture(t)
			if err := s.WriteLog(cicdValueBytes(value)); err != nil {
				t.Fatal("statefile.phase19-log.corrupt.initial", err)
			}
			id := cicdLogIdentity(value)
			path := filepath.Join(s.store.path, "cicd", "logs", "p", id.RunID, "1", "build", "stdout.json")
			var err error
			switch name {
			case "identity":
				value["stream"] = "stderr"
				err = os.WriteFile(path, cicdValueBytes(value), 0600)
			case "malformed":
				err = os.WriteFile(path, []byte(`{}`), 0600)
			case "symlink":
				other := filepath.Join(t.TempDir(), "other")
				if os.WriteFile(other, cicdValueBytes(value), 0600) != nil || os.Remove(path) != nil {
					t.Fatal("statefile.phase19-log.corrupt.fixture")
				}
				err = os.Symlink(other, path)
			case "hardlink":
				err = os.Link(path, filepath.Join(t.TempDir(), "other"))
			case "file-mode":
				err = os.Chmod(path, 0644)
			case "directory-mode":
				err = os.Chmod(filepath.Dir(path), 0755)
			case "oversized":
				err = os.Truncate(path, cicdLogLimit+1)
			}
			if err != nil {
				t.Fatal("statefile.phase19-log.corrupt.fixture", err)
			}
			if data, err := s.ReadLog(id); err != ErrCICDCorrupt || data != nil {
				t.Fatal("statefile.phase19-log.corrupt.read", name, err)
			}
			value["stream"] = "stdout"
			if err := s.WriteLog(cicdValueBytes(value)); err != ErrCICDCorrupt {
				t.Fatal("statefile.phase19-log.corrupt.write", name, err)
			}
		})
	}
}

func cicdRunFixture(t *testing.T) (*CICDRunStore, CICDRunRequest, json.RawMessage) {
	t.Helper()
	s, operation := cicdOperationFixture(t)
	project := `{"id":"p","name":"Project","revision":1,"enabled":true,"source_owner":"owner","source_repo":"repo","definition_id":"d","definition_revision":1,"quota":{"parallel":1,"queue":64,"artifact_mb":128,"cache_mb":0,"retention_days":7},"toolchain_digest":"` + strings.Repeat("a", 64) + `"}`
	job := `{"id":"build","needs":[],"condition":{"branches":[],"tags":[],"paths":[],"input_equals":{}},"task":{"kind":"builder","source":".","output":"site"},"matrix":{"axes":{},"exclude":[],"include":[]},"resources":{"cpu_percent":100,"memory_mb":256,"disk_mb":128,"tasks":64,"timeout_seconds":900,"output_mb":16},"retry":{"max_attempts":1,"backoff_seconds":1},"cache":{"enabled":false,"paths":[],"ttl_seconds":604800},"env":{},"secrets":{},"optional":false}`
	definition := `{"project_id":"p","id":"d","revision":1,"inputs":{},"manual":true,"watches":[],"schedules":[],"templates":[],"jobs":[` + job + `],"fail_fast":false}`
	p, _ := decodeCICDValue([]byte(project))
	d, _ := decodeCICDValue([]byte(definition))
	j, _ := decodeCICDValue([]byte(job))
	ignoredDefinition := cloneCICDObject(cicdObject(j))
	ignoredDefinition["id"] = "ignored"
	cicdObject(ignoredDefinition["condition"])["branches"] = []any{"release"}
	cicdObject(d)["jobs"] = append(cicdList(cicdObject(d)["jobs"]), ignoredDefinition)
	planned := cicdObject(j)
	delete(planned, "condition")
	delete(planned, "matrix")
	planned["base_job_id"], planned["matrix_values"], planned["selected"] = "build", map[string]any{}, true
	skipped := cloneCICDObject(planned)
	skipped["id"], skipped["base_job_id"], skipped["selected"] = "ignored", "ignored", false
	data := map[string]any{"project_id": "p", "definition_id": "d", "definition_revision": json.Number("1"), "source": map[string]any{"owner": "owner", "repo": "repo", "ref_kind": "branch", "ref": "main", "sha": strings.Repeat("b", 40), "tree_digest": "4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945", "files": []any{}}, "values": map[string]any{}, "jobs": []any{planned, skipped}, "fail_fast": false, "toolchain_digest": strings.Repeat("a", 64)}
	encoded, _ := json.Marshal(data)
	digest := fmt.Sprintf("%x", sha256.Sum256(encoded))
	plan, _ := json.Marshal(map[string]any{"data": data, "digest": digest})
	if _, err := s.commit(1, operation.id, func(root map[string]any) error {
		root["projects"], root["definitions"] = []any{p}, []any{d}
		root["grants"] = append(cicdList(root["grants"]), map[string]any{"subject_kind": "user", "subject_id": "admin", "project_id": "p", "environment_id": nil, "role": "operator", "enabled": true, "expires_at": nil})
		return nil
	}); err != nil {
		t.Fatal("statefile.phase19-run.fixture.root", err)
	}
	request := CICDRunRequest{Principal: operation.principal, ProjectID: "p", DefinitionID: "d", DefinitionRevision: 1, OperationID: operation.id, RequestID: operation.requestID, RequestDigest: operation.digest, PlanDigest: digest, Check: func(current CICDIdentity) error {
		if current.Revision != 1 {
			return security.ErrCICDAuth
		}
		return nil
	}}
	return &CICDRunStore{store: s, entropy: bytes.NewReader(bytes.Repeat([]byte{0x11}, 16))}, request, plan
}

func TestCICDRunAdmission(t *testing.T) {
	s, request, plan := cicdRunFixture(t)
	before, head, _ := s.store.read()
	snapshot, replay, err := s.CheckRunStart(request)
	if err != nil || replay.Replayed || snapshot.Generation != 2 || snapshot.Project == nil || snapshot.Definition == nil {
		t.Fatal("statefile.phase19-run.read.snapshot", err)
	}
	snapshot.Project[0] = '!'
	afterRead, readHead, _ := s.store.read()
	if head != readHead || !bytes.Equal(cicdValueBytes(before), cicdValueBytes(afterRead)) {
		t.Fatal("statefile.phase19-run.read.no-write")
	}
	calls, clockCalls := 0, 0
	now := s.store.clock()
	s.store.clock = func() time.Time { clockCalls++; return now }
	resolve := func(current CICDPlanSnapshot) (json.RawMessage, error) {
		calls++
		if current.Generation != 2 || current.Project[0] != '{' {
			t.Fatal("statefile.phase19-run.resolve.current")
		}
		return plan, nil
	}
	reply, err := s.StartRun(request, 1, resolve)
	if err != nil || reply.Replayed || reply.Generation != 3 || calls != 1 || clockCalls != 1 {
		t.Fatal("statefile.phase19-run.accept.result", err)
	}
	root, committed, _ := s.store.read()
	id := "r-" + strings.Repeat("11", 16)
	operationValue, _ := decodeCICDValue(reply.Operation)
	op := cicdObject(operationValue)
	if committed.Generation != 3 || op["target_id"] != id || op["target_attempt"] != json.Number("1") || op["state"] != "accepted" || len(cicdList(root["runs"])) != 1 || len(cicdList(root["audits"])) != 1 || len(cicdList(root["notifications"])) != 1 {
		t.Fatal("statefile.phase19-run.accept.atomic")
	}
	run := cicdObject(cicdList(root["runs"])[0])
	eventBytes := []byte(`{"definition_id":"d","definition_revision":1,"due":null,"kind":"manual","project_id":"p","source_sha":"` + strings.Repeat("b", 40) + `","trigger_id":"` + request.OperationID + `"}`)
	expectedEvent := fmt.Sprintf("%x", sha256.Sum256(eventBytes))
	if run["id"] != id || run["event_key"] != expectedEvent || run["created_at"] != "2026-10-22T00:00:00Z" || run["requester"] != "admin" || run["trigger"] != "manual" || run["revision"] != json.Number("1") || run["attempt"] != json.Number("1") || run["state"] != "queued" || run["started_at"] != nil || run["finished_at"] != nil || run["cancel_requested"] != false || !bytes.Equal(cicdValueBytes(run["plan"]), plan) {
		t.Fatal("statefile.phase19-run.accept.run-fields")
	}
	wantJobs := `[{"artifact":null,"code":null,"duration_ms":null,"execution_id":null,"finished_at":null,"id":"build","lease_generation":0,"log_gap":false,"reused_from_attempt":null,"started_at":null,"state":"queued","try_number":0,"worker_id":null},{"artifact":null,"code":null,"duration_ms":0,"execution_id":null,"finished_at":"2026-10-22T00:00:00Z","id":"ignored","lease_generation":0,"log_gap":false,"reused_from_attempt":null,"started_at":null,"state":"skipped","try_number":0,"worker_id":null}]`
	if string(cicdValueBytes(run["jobs"])) != wantJobs {
		t.Fatal("statefile.phase19-run.accept.job-fields")
	}
	notification := cicdObject(cicdList(root["notifications"])[0])
	key := []byte(`{"event":"queued","project_id":"p","source_revision":1,"target_attempt":1,"target_id":"` + id + `"}`)
	wantNotificationID := "n" + fmt.Sprintf("%x", sha256.Sum256(key))[:63]
	if notification["id"] != wantNotificationID || notification["target_id"] != id || notification["state"] != "queued" || notification["at"] != run["created_at"] || notification["next_at"] != run["created_at"] {
		t.Fatal("statefile.phase19-run.accept.notification")
	}
	reply.Operation[0] = '!'
	request.DefinitionRevision = 99
	replayed, err := s.StartRun(request, 99, func(CICDPlanSnapshot) (json.RawMessage, error) { calls++; return nil, ErrCICDStale })
	if err != nil || !replayed.Replayed || calls != 1 || replayed.Operation[0] != '{' {
		t.Fatal("statefile.phase19-run.replay.precedence", err)
	}
	_, finalHead, _ := s.store.read()
	if finalHead != committed {
		t.Fatal("statefile.phase19-run.replay.no-write")
	}
	_, checked, err := s.CheckRunStart(request)
	if err != nil || !checked.Replayed || !bytes.Equal(checked.Operation, replayed.Operation) {
		t.Fatal("statefile.phase19-run.replay.read")
	}
}

func TestCICDRunRejection(t *testing.T) {
	for _, name := range []string{"project-revision", "plan-digest", "plan-hash", "scope", "schema", "resolve", "manual", "quota", "entropy", "identity", "grant", "key"} {
		t.Run(name, func(t *testing.T) {
			s, request, plan := cicdRunFixture(t)
			revision, expected := int64(1), ErrCICDStale
			resolveCalls := 0
			var resolutionError error
			switch name {
			case "project-revision":
				revision = 2
			case "plan-digest":
				request.PlanDigest = strings.Repeat("c", 64)
			case "plan-hash", "scope":
				value, _ := decodeCICDValue(plan)
				object := cicdObject(value)
				if name == "plan-hash" {
					object["digest"] = strings.Repeat("f", 64)
				} else {
					cicdObject(object["data"])["project_id"] = "other"
					object["digest"] = cicdValueDigest(object["data"])
				}
				plan, expected = cicdValueBytes(value), ErrCICDSchema
			case "schema":
				plan, expected = json.RawMessage(`{}`), ErrCICDSchema
			case "resolve":
				resolutionError, expected = errors.New("P19_SECRET"), errors.New("P19_SECRET")
			case "entropy":
				s.entropy, expected = bytes.NewReader(nil), errors.New("P19_ENTROPY")
			case "key":
				request.RequestID, expected = request.OperationID, ErrCICDSchema
			case "identity":
				request.Check, expected = func(CICDIdentity) error { return security.ErrCICDAuth }, security.ErrCICDAuth
			case "manual", "quota", "grant":
				if _, err := s.store.commit(2, "o-"+strings.Repeat("e", 32), func(root map[string]any) error {
					switch name {
					case "manual":
						cicdObject(cicdList(root["definitions"])[0])["manual"] = false
					case "quota":
						cicdObject(cicdObject(cicdList(root["projects"])[0])["quota"])["queue"] = json.Number("1")
					case "grant":
						cicdObject(cicdList(root["grants"])[1])["enabled"] = false
					}
					return nil
				}); err != nil {
					t.Fatal("statefile.phase19-run.reject.fixture")
				}
				if name != "quota" {
					expected = security.ErrCICDForbidden
				} else {
					if _, err := s.StartRun(request, 1, func(CICDPlanSnapshot) (json.RawMessage, error) { return plan, nil }); err != nil {
						t.Fatal("statefile.phase19-run.quota.first", err)
					}
					request.OperationID, request.RequestID, expected = "o-"+strings.Repeat("c", 32), "o-"+strings.Repeat("f", 32), errors.New("P19_QUOTA")
				}
			}
			before, beforeHead, _ := s.store.read()
			reply, err := s.StartRun(request, revision, func(CICDPlanSnapshot) (json.RawMessage, error) { resolveCalls++; return plan, resolutionError })
			if err == nil || err.Error() != expected.Error() || reply.Operation != nil {
				t.Fatal("statefile.phase19-run.reject.result", err)
			}
			after, afterHead, _ := s.store.read()
			for _, collection := range []string{"runs", "operations", "notifications", "artifacts"} {
				if !bytes.Equal(cicdValueBytes(before[collection]), cicdValueBytes(after[collection])) {
					t.Fatal("statefile.phase19-run.reject.business-unchanged")
				}
			}
			if name == "key" || name == "identity" || name == "grant" {
				if afterHead != beforeHead || !bytes.Equal(cicdValueBytes(before["audits"]), cicdValueBytes(after["audits"])) || resolveCalls != 0 {
					t.Fatal("statefile.phase19-run.reject.unauthorized-no-write")
				}
			} else if afterHead.Generation != beforeHead.Generation+1 || len(cicdList(after["audits"])) != len(cicdList(before["audits"]))+1 || cicdObject(cicdList(after["audits"])[len(cicdList(after["audits"]))-1])["operation_id"] != request.RequestID {
				t.Fatal("statefile.phase19-run.reject.audit-only")
			}
		})
	}
}

func TestCICDRunRequestCollision(t *testing.T) {
	for _, name := range []string{"rejected-audit", "operation", "replay"} {
		t.Run(name, func(t *testing.T) {
			s, request, plan := cicdRunFixture(t)
			resolve := func(CICDPlanSnapshot) (json.RawMessage, error) { return plan, nil }
			if name == "rejected-audit" {
				original := request.PlanDigest
				request.PlanDigest = strings.Repeat("f", 64)
				if _, err := s.StartRun(request, 1, resolve); err == nil || err.Error() != "P19_STALE" {
					t.Fatal("statefile.phase19-run.request-id.fixture")
				}
				request.PlanDigest = original
			} else {
				if _, err := s.StartRun(request, 1, resolve); err != nil {
					t.Fatal("statefile.phase19-run.request-id.fixture", err)
				}
				request.RequestID = request.OperationID
				if name == "operation" {
					request.OperationID = "o-" + strings.Repeat("c", 32)
				} else {
					// A different retained operation supplies the colliding server ID.
					other := request
					other.OperationID, other.RequestID = "o-"+strings.Repeat("c", 32), "o-"+strings.Repeat("d", 32)
					s.entropy = bytes.NewReader(bytes.Repeat([]byte{0x22}, 16))
					if _, err := s.StartRun(other, 1, resolve); err != nil {
						t.Fatal("statefile.phase19-run.request-id.second", err)
					}
					request.RequestID = other.OperationID
				}
			}
			before, head, _ := s.store.read()
			if _, reply, err := s.CheckRunStart(request); err == nil || err.Error() != "P19_ENTROPY" || reply.Operation != nil {
				t.Fatal("statefile.phase19-run.request-id.read-reject", err)
			}
			calls := 0
			if reply, err := s.StartRun(request, 1, func(CICDPlanSnapshot) (json.RawMessage, error) { calls++; return plan, nil }); err == nil || err.Error() != "P19_ENTROPY" || reply.Operation != nil || calls != 0 {
				t.Fatal("statefile.phase19-run.request-id.admission-reject", err)
			}
			after, final, _ := s.store.read()
			if head != final || !bytes.Equal(cicdValueBytes(before), cicdValueBytes(after)) {
				t.Fatal("statefile.phase19-run.request-id.no-write")
			}
		})
	}
}

func TestCICDRunPins(t *testing.T) {
	for _, name := range []string{"selected", "duplicate", "unselected", "missing", "deleting"} {
		t.Run(name, func(t *testing.T) {
			s, request, plan := cicdRunFixture(t)
			value, _ := decodeCICDValue(plan)
			object := cicdObject(value)
			jobs := cicdList(cicdObject(object["data"])["jobs"])
			job := cicdObject(jobs[0])
			job["task"] = map[string]any{"kind": "artifact-verify", "artifact": strings.Repeat("c", 64)}
			if name == "unselected" {
				job["selected"] = false
			}
			if name == "duplicate" {
				cicdObject(jobs[1])["task"], cicdObject(jobs[1])["selected"] = cloneCICDObject(cicdObject(job["task"])), true
			}
			object["digest"] = cicdValueDigest(object["data"])
			request.PlanDigest, plan = object["digest"].(string), cicdValueBytes(object)
			if name != "missing" {
				if _, err := s.store.commit(2, "o-"+strings.Repeat("e", 32), func(root map[string]any) error {
					root["artifacts"] = []any{map[string]any{"project_id": "p", "digest": strings.Repeat("c", 64), "manifest_digest": strings.Repeat("d", 64), "size_bytes": json.Number("1"), "created_at": "2026-10-22T00:00:00Z", "source_sha": strings.Repeat("b", 40), "run_id": "r-" + strings.Repeat("9", 32), "attempt": json.Number("1"), "job_id": "build", "pins": []any{}, "deleting": name == "deleting"}}
					return nil
				}); err != nil {
					t.Fatal("statefile.phase19-run.pin.fixture", err)
				}
			}
			reply, err := s.StartRun(request, 1, func(CICDPlanSnapshot) (json.RawMessage, error) { return plan, nil })
			root, _, _ := s.store.read()
			if name == "missing" || name == "deleting" {
				if err == nil || err.Error() != "P19_NOT_FOUND" || reply.Operation != nil || len(cicdList(root["runs"])) != 0 || len(cicdList(root["notifications"])) != 0 {
					t.Fatal("statefile.phase19-run.pin.reject", err)
				}
			} else {
				pins := cicdList(cicdObject(cicdList(root["artifacts"])[0])["pins"])
				want := 1
				if name == "unselected" {
					want = 0
				}
				if err != nil || len(pins) != want || want == 1 && pins[0] != "r-"+strings.Repeat("11", 16) {
					t.Fatal("statefile.phase19-run.pin.atomic", err)
				}
			}
		})
	}
}

func TestCICDRunReadScope(t *testing.T) {
	s, request, _ := cicdRunFixture(t)
	if _, err := s.store.commit(2, "o-"+strings.Repeat("e", 32), func(root map[string]any) error {
		project := cloneCICDObject(cicdObject(cicdList(root["projects"])[0]))
		project["id"], project["definition_id"], project["definition_revision"] = "other", nil, json.Number("0")
		root["projects"] = append(cicdList(root["projects"]), project)
		cicdObject(cicdList(root["definitions"])[0])["manual"] = false
		root["secrets"] = []any{map[string]any{"project_id": "p", "id": "local", "version": json.Number("1"), "created_at": "2026-10-22T00:00:00Z"}, map[string]any{"project_id": "other", "id": "hidden", "version": json.Number("1"), "created_at": "2026-10-22T00:00:00Z"}}
		return nil
	}); err != nil {
		t.Fatal("statefile.phase19-run.scope.fixture")
	}
	before, head, _ := s.store.read()
	preview, err := s.ReadPlan(request.Principal, "p", "d", 1, request.Check)
	if err != nil || len(preview.Secrets) != 1 || !bytes.Contains(preview.Secrets[0], []byte(`"id":"local"`)) || bytes.Contains(preview.Secrets[0], []byte("hidden")) {
		t.Fatal("statefile.phase19-run.scope.preview", err)
	}
	if _, _, err := s.CheckRunStart(request); err != security.ErrCICDForbidden {
		t.Fatal("statefile.phase19-run.scope.manual-denied")
	}
	if _, err := s.ReadPlan(request.Principal, "other", "d", 1, request.Check); err != security.ErrCICDForbidden {
		t.Fatal("statefile.phase19-run.scope.other-project")
	}
	after, afterHead, _ := s.store.read()
	if head != afterHead || !bytes.Equal(cicdValueBytes(before), cicdValueBytes(after)) {
		t.Fatal("statefile.phase19-run.scope.read-only")
	}
}

func TestCICDRunConcurrency(t *testing.T) {
	s, request, plan := cicdRunFixture(t)
	var calls, accepted, replayed, busy atomic.Int64
	start := make(chan struct{})
	var wait sync.WaitGroup
	for i := 0; i < 8; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			reply, err := s.StartRun(request, 1, func(CICDPlanSnapshot) (json.RawMessage, error) { calls.Add(1); return plan, nil })
			if err == ErrCICDBusy {
				busy.Add(1)
			} else if err != nil {
				t.Error("statefile.phase19-run.concurrent.error", err)
			} else if reply.Replayed {
				replayed.Add(1)
			} else {
				accepted.Add(1)
			}
		}()
	}
	close(start)
	wait.Wait()
	root, head, err := s.store.read()
	if err != nil || calls.Load() != 1 || accepted.Load() != 1 || accepted.Load()+replayed.Load()+busy.Load() != 8 || head.Generation != 3 || len(cicdList(root["runs"])) != 1 || len(cicdList(root["operations"])) != 1 || len(cicdList(root["notifications"])) != 1 || len(cicdList(root["audits"])) != 1 {
		t.Fatal("statefile.phase19-run.concurrent.single-admission", err)
	}
}

func TestCICDRunCapacity(t *testing.T) {
	s, request, plan := cicdRunFixture(t)
	var entropy []byte
	for i := byte(1); i <= 16; i++ {
		entropy = append(entropy, bytes.Repeat([]byte{i}, 16)...)
	}
	s.entropy = bytes.NewReader(entropy)
	for i := 0; i < 17; i++ {
		request.OperationID, request.RequestID = fmt.Sprintf("o-%032x", 100+i), fmt.Sprintf("o-%032x", 200+i)
		reply, err := s.StartRun(request, 1, func(CICDPlanSnapshot) (json.RawMessage, error) { return plan, nil })
		if i < 16 && (err != nil || reply.Replayed) || i == 16 && (err == nil || err.Error() != "P19_QUOTA" || reply.Operation != nil) {
			t.Fatal("statefile.phase19-run.capacity.boundary", err)
		}
	}
	root, head, _ := s.store.read()
	if head.Generation != 19 || len(cicdList(root["runs"])) != 16 || len(cicdList(root["operations"])) != 16 || len(cicdList(root["notifications"])) != 16 || len(cicdList(root["audits"])) != 17 {
		t.Fatal("statefile.phase19-run.capacity.atomic")
	}
}

func TestCICDRunFaults(t *testing.T) {
	for _, stage := range []string{"snapshot.create", "snapshot.write", "snapshot.sync", "head.write", "head.publish", "head.directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			s, request, plan := cicdRunFixture(t)
			s.store.fault = func(current string) error {
				if current == stage {
					return errors.New("fixture failure")
				}
				return nil
			}
			reply, err := s.StartRun(request, 1, func(CICDPlanSnapshot) (json.RawMessage, error) { return plan, nil })
			if err == nil || reply.Operation != nil {
				t.Fatal("statefile.phase19-run.fault.no-success")
			}
			fresh, openErr := OpenCICDRunStore(s.store.path)
			if openErr != nil {
				t.Fatal("statefile.phase19-run.fault.reopen", openErr)
			}
			root, head, readErr := fresh.store.read()
			if stage == "head.directory-sync" {
				if err != ErrCICDUncertain || head.Generation != 3 || len(cicdList(root["runs"])) != 1 || len(cicdList(root["notifications"])) != 1 || len(cicdList(root["audits"])) != 1 {
					t.Fatal("statefile.phase19-run.fault.uncertain-atomic")
				}
			} else if readErr != nil || head.Generation != 2 || len(cicdList(root["runs"])) != 0 || len(cicdList(root["operations"])) != 0 || len(cicdList(root["notifications"])) != 0 || len(cicdList(root["audits"])) != 0 {
				t.Fatal("statefile.phase19-run.fault.before-commit")
			}
		})
	}
}

func cicdOperationFixture(t *testing.T) (*cicdSnapshotStore, cicdOperationRequest) {
	t.Helper()
	s := cicdLinuxStore(t)
	now := time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC)
	s.clock = func() time.Time { return now }
	request := cicdOperationRequest{id: cicdTestOperation, digest: strings.Repeat("d", 64), requestID: "o-" + strings.Repeat("b", 32), principal: security.CICDPrincipal{Kind: "user", ID: "admin", ExpiresAt: now.Add(time.Hour), MustChange: "none"}, access: security.CICDAccess{Operation: "projectPut"}, check: func() error { return nil }}
	return s, request
}

func cicdIdentityFixture(t *testing.T) (*CICDIdentityStore, time.Time) {
	t.Helper()
	s, request := cicdOperationFixture(t)
	record, err := security.NewPasswordRecord("Fixture-password-2026")
	if err != nil {
		t.Fatal("statefile.phase19-identity.credential.fixture")
	}
	credential := security.CICDCredential{PasswordHash: record.PasswordHash, Salt: record.Salt, Algorithm: record.Algorithm, Iterations: record.Iterations, MustChange: true, UpdatedAt: "2026-10-22T00:00:00Z"}
	encoded, _ := json.Marshal(credential)
	value, _ := decodeCICDValue(encoded)
	if _, err := s.commit(1, request.id, func(root map[string]any) error {
		root["users"] = append(cicdList(root["users"]), map[string]any{"id": "alice", "name": "Alice", "revision": json.Number("1"), "enabled": true, "credential": value, "totp": nil, "created_at": "2026-10-22T00:00:00Z"})
		root["services"] = []any{map[string]any{"id": "build", "revision": json.Number("1"), "enabled": true, "token_sha256": strings.Repeat("a", 64), "expires_at": "2026-10-23T00:00:00Z"}}
		return nil
	}); err != nil {
		t.Fatal("statefile.phase19-identity.root.fixture")
	}
	return &CICDIdentityStore{store: s}, s.clock()
}

func TestCICDIdentityStore(t *testing.T) {
	s, now := cicdIdentityFixture(t)
	before, head, _ := s.store.read()
	if _, err := s.Login("alice", "wrong-password", cicdTestOperation, now, nil); err != security.ErrCICDAuth {
		t.Fatal("statefile.phase19-identity.password.reject")
	}
	_, after, _ := s.store.read()
	if head != after {
		t.Fatal("statefile.phase19-identity.reject.no-write")
	}
	result, err := s.Login("alice", "Fixture-password-2026", cicdTestOperation, now, nil)
	if err != nil || !result.PasswordVerified || result.MustChange != "prompt" || result.Identity.Credential.LoginCount != 1 {
		t.Fatal("statefile.phase19-identity.login.count")
	}
	result.Identity.Credential.PasswordHash = strings.Repeat("c", 64)
	identity, err := s.ReadIdentity("user", "alice")
	if err != nil || identity.Credential.LoginCount != 1 || identity.Credential.PasswordHash == strings.Repeat("c", 64) {
		t.Fatal("statefile.phase19-identity.read.copy")
	}
	p := security.CICDPrincipal{Kind: "user", ID: "alice", ExpiresAt: now.Add(time.Hour), MustChange: "prompt"}
	if _, err := s.Authorize(p, security.CICDAccess{Operation: "logout"}, now, nil); err != security.ErrCICDAuth {
		t.Fatal("statefile.phase19-identity.auth.check-required")
	}
	if _, err := s.Authorize(p, security.CICDAccess{Operation: "logout"}, now, func(CICDIdentity) error { return security.ErrCICDAuth }); err != security.ErrCICDAuth {
		t.Fatal("statefile.phase19-identity.auth.revoked")
	}
	_, head, _ = s.store.read()
	auth, err := s.Authorize(p, security.CICDAccess{Operation: "logout"}, now, func(current CICDIdentity) error {
		if current.ID != "alice" {
			t.Fatal("statefile.phase19-identity.auth.subject")
		}
		return nil
	})
	if err != nil || auth.Generation != head.Generation || len(auth.Grants) != 0 {
		t.Fatal("statefile.phase19-identity.auth.current")
	}
	service, err := s.ReadService(strings.Repeat("a", 64))
	if err != nil || service.ID != "build" {
		t.Fatal("statefile.phase19-identity.service.match")
	}
	if _, err := s.ReadService(strings.Repeat("b", 64)); err != security.ErrCICDAuth {
		t.Fatal("statefile.phase19-identity.service.unknown")
	}
	_, after, _ = s.store.read()
	if head != after || len(cicdList(before["audits"])) != 0 {
		t.Fatal("statefile.phase19-identity.read.no-write")
	}
	if _, err := s.store.commit(head.Generation, cicdTestOperation, func(root map[string]any) error {
		duplicate := cloneCICDObject(cicdObject(cicdList(root["services"])[0]))
		duplicate["id"] = "other"
		root["services"] = append(cicdList(root["services"]), duplicate)
		return nil
	}); err != nil {
		t.Fatal("statefile.phase19-identity.service.duplicate.fixture")
	}
	if _, err := s.ReadService(strings.Repeat("a", 64)); err != ErrCICDCorrupt {
		t.Fatal("statefile.phase19-identity.service.duplicate")
	}
}

func TestCICDIdentityFaults(t *testing.T) {
	for _, stage := range []string{"snapshot.write", "head.directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			s, now := cicdIdentityFixture(t)
			s.store.fault = func(actual string) error {
				if actual == stage {
					return errors.New("PRIVATE_AUTH_IO")
				}
				return nil
			}
			result, err := s.Login("alice", "Fixture-password-2026", cicdTestOperation, now, nil)
			want := ErrCICDStorage
			if stage == "head.directory-sync" {
				want = ErrCICDUncertain
			}
			if err != want || result.Identity.ID != "" || result.MustChange != "" {
				t.Fatal("statefile.phase19-identity.fault.no-identity")
			}
		})
	}
}

func TestCICDIdentityConcurrency(t *testing.T) {
	s, now := cicdIdentityFixture(t)
	start := make(chan struct{})
	errs := make(chan error, 4)
	var group sync.WaitGroup
	for i := 0; i < 4; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			_, err := s.Login("alice", "Fixture-password-2026", cicdTestOperation, now, nil)
			errs <- err
		}()
	}
	close(start)
	group.Wait()
	close(errs)
	succeeded, busy := 0, 0
	for err := range errs {
		switch err {
		case nil:
			succeeded++
		case ErrCICDBusy:
			busy++
		default:
			t.Fatal("statefile.phase19-identity.concurrent.result")
		}
	}
	identity, readErr := s.ReadIdentity("user", "alice")
	_, concurrentHead, headErr := s.store.read()
	if readErr != nil || headErr != nil || succeeded < 1 || succeeded+busy != 4 || identity.Credential.LoginCount != int64(succeeded) || concurrentHead.Generation != int64(2+succeeded) {
		t.Fatal("statefile.phase19-identity.concurrent.busy-no-write")
	}
	for i := 0; i < busy; i++ {
		if _, err := s.Login("alice", "Fixture-password-2026", cicdTestOperation, now, nil); err != nil {
			t.Fatal("statefile.phase19-identity.concurrent.busy-retry")
		}
	}
	identity, _ = s.ReadIdentity("user", "alice")
	if identity.Credential.LoginCount != 4 {
		t.Fatal("statefile.phase19-identity.concurrent.no-lost-update")
	}
	result, err := s.Login("alice", "Fixture-password-2026", cicdTestOperation, now, nil)
	if err != nil || result.MustChange != "forced" {
		t.Fatal("statefile.phase19-identity.concurrent.forced")
	}
	_, head, _ := s.store.read()
	if _, err := s.store.commit(head.Generation, cicdTestOperation, func(root map[string]any) error {
		for _, value := range cicdList(root["users"]) {
			user := cicdObject(value)
			if user["id"] == "alice" {
				cicdObject(user["credential"])["login_count"] = json.Number("9007199254740991")
			}
		}
		return nil
	}); err != nil {
		t.Fatal("statefile.phase19-identity.saturation.fixture")
	}
	result, err = s.Login("alice", "Fixture-password-2026", cicdTestOperation, now, nil)
	if err != nil || result.Identity.Credential.LoginCount != 9007199254740991 {
		t.Fatal("statefile.phase19-identity.saturation.value")
	}
}

func cicdProjectFixture() map[string]any {
	return map[string]any{"id": "p", "name": "Project", "revision": json.Number("1"), "enabled": false, "source_owner": "owner", "source_repo": "repo", "definition_id": nil, "definition_revision": json.Number("0"), "toolchain_digest": nil, "quota": map[string]any{"parallel": json.Number("1"), "queue": json.Number("64"), "artifact_mb": json.Number("1024"), "cache_mb": json.Number("256"), "retention_days": json.Number("7")}}
}

func cicdCreateProject(root map[string]any) (cicdOperationMutation, error) {
	root["projects"] = append(cicdList(root["projects"]), cicdProjectFixture())
	return cicdOperationMutation{targetID: "p", result: json.RawMessage(`{"id":"p"}`)}, nil
}

func TestCICDOperationAdmission(t *testing.T) {
	s, request := cicdOperationFixture(t)
	calls := 0
	first, err := s.updateCICD(request, func(root map[string]any) (cicdOperationMutation, error) { calls++; return cicdCreateProject(root) })
	if err != nil || first.replayed || first.operation["state"] != "succeeded" || calls != 1 || string(first.result) != `{"id":"p"}` || first.generation != 2 {
		t.Fatalf("statefile.phase19-operation.admission.first: %v", err)
	}
	root, head, err := s.read()
	if err != nil || len(cicdList(root["operations"])) != 1 || len(cicdList(root["audits"])) != 2 || first.operation["expires_at"] != "2026-11-21T00:00:00Z" {
		t.Fatal("statefile.phase19-operation.admission.atomic")
	}
	_, err = s.commit(head.Generation, request.id, func(root map[string]any) error {
		project := cicdObject(cicdList(root["projects"])[0])
		project["revision"], project["name"] = json.Number("2"), "Newer"
		return nil
	})
	if err != nil {
		t.Fatal("statefile.phase19-operation.replay.revision.fixture")
	}
	replayed, err := s.updateCICD(request, func(map[string]any) (cicdOperationMutation, error) {
		calls++
		return cicdOperationMutation{}, ErrCICDRevision
	})
	if err != nil || !replayed.replayed || calls != 1 || replayed.result != nil || replayed.oneTime != "" || replayed.generation != 3 {
		t.Fatal("statefile.phase19-operation.replay.before-revision")
	}
	replayed.operation["state"] = "tampered"
	root, after, err := s.read()
	if err != nil || after.Generation != 3 || findCICDOperation(root, request.id)["state"] != "succeeded" || len(cicdList(root["audits"])) != 2 {
		t.Fatal("statefile.phase19-operation.replay.immutable")
	}
	request.digest = strings.Repeat("e", 64)
	denied, err := s.updateCICD(request, nil)
	if err != ErrCICDIdempotency || !reflect.DeepEqual(denied, cicdOperationReply{}) {
		t.Fatal("statefile.phase19-operation.replay.conflict")
	}
}

func TestCICDOperationSchema(t *testing.T) {
	for _, test := range []struct {
		id     string
		mutate func(map[string]any)
	}{
		{"unknown-kind", func(root map[string]any) { cicdObject(cicdList(root["operations"])[0])["kind"] = "arbitrary" }},
		{"unknown-code", func(root map[string]any) { cicdObject(cicdList(root["operations"])[0])["code"] = "PRIVATE_ERROR" }},
		{"null-target", func(root map[string]any) { cicdObject(cicdList(root["operations"])[0])["target_id"] = nil }},
		{"orphan-target", func(root map[string]any) { root["projects"] = []any{} }},
		{"missing-accepted", func(root map[string]any) { root["audits"] = cicdList(root["audits"])[1:] }},
		{"missing-finished", func(root map[string]any) { root["audits"] = cicdList(root["audits"])[:1] }},
		{"duplicate-audit", func(root map[string]any) {
			root["audits"] = append(cicdList(root["audits"]), cicdList(root["audits"])[0])
		}},
		{"audit-action", func(root map[string]any) { cicdObject(cicdList(root["audits"])[0])["action"] = "userPut" }},
		{"audit-identity", func(root map[string]any) { cicdObject(cicdList(root["audits"])[0])["subject"] = "bob" }},
		{"audit-time", func(root map[string]any) { cicdObject(cicdList(root["audits"])[1])["at"] = "2026-10-21T00:00:00Z" }},
		{"expiry", func(root map[string]any) {
			cicdObject(cicdList(root["operations"])[0])["expires_at"] = "2026-11-22T00:00:00Z"
		}},
	} {
		t.Run(test.id, func(t *testing.T) {
			root := cicdRootFixture(t)
			root["projects"] = []any{cicdProjectFixture()}
			operation := map[string]any{"id": cicdTestOperation, "project_id": nil, "subject": "admin", "request_digest": strings.Repeat("d", 64), "kind": "projectPut", "state": "succeeded", "target_id": "p", "target_attempt": nil, "created_at": "2026-10-22T00:00:00Z", "expires_at": "2026-11-21T00:00:00Z", "code": nil}
			root["operations"] = []any{operation}
			root["audits"] = []any{map[string]any{"operation_id": cicdTestOperation, "subject": "admin", "project_id": nil, "environment_id": nil, "action": "projectPut", "phase": "accepted", "at": "2026-10-22T00:00:00Z", "target_id": "p", "code": nil}, map[string]any{"operation_id": cicdTestOperation, "subject": "admin", "project_id": nil, "environment_id": nil, "action": "projectPut", "phase": "finished", "at": "2026-10-22T00:00:00Z", "target_id": "p", "code": nil}}
			if !validateCICDRoot(root) {
				t.Fatal("statefile.phase19-operation.schema.positive")
			}
			test.mutate(root)
			if validateCICDRoot(root) {
				t.Fatal("statefile.phase19-operation.schema.reject")
			}
		})
	}
}

func TestCICDOperationOneTime(t *testing.T) {
	s, request := cicdOperationFixture(t)
	request.access.Operation = "servicePut"
	const secret = "MEMORY_ONLY_TEST_CREDENTIAL_7f1dbd8a"
	calls := 0
	first, err := s.updateCICD(request, func(root map[string]any) (cicdOperationMutation, error) {
		calls++
		root["services"] = append(cicdList(root["services"]), map[string]any{"id": "svc", "revision": json.Number("1"), "token_sha256": strings.Repeat("1", 64), "enabled": true, "expires_at": "2026-11-21T00:00:00Z"})
		return cicdOperationMutation{targetID: "svc", result: json.RawMessage(`{"id":"svc"}`), oneTime: secret}, nil
	})
	if err != nil || first.oneTime != secret {
		t.Fatalf("statefile.phase19-operation.one-time.first: %v", err)
	}
	reply, err := s.updateCICD(request, func(map[string]any) (cicdOperationMutation, error) {
		calls++
		return cicdOperationMutation{}, ErrCICDInternal
	})
	if err != nil || !reply.replayed || reply.oneTime != "" || reply.result != nil || calls != 1 {
		t.Fatal("statefile.phase19-operation.one-time.replay")
	}
	root, _, err := s.read()
	if err != nil || bytes.Contains(cicdValueBytes(root), []byte(secret)) || bytes.Contains(cicdValueBytes(root), []byte(`"result"`)) {
		t.Fatal("statefile.phase19-operation.one-time.not-stored")
	}
	request.id = "o-" + strings.Repeat("c", 32)
	if reply, err := s.updateCICD(request, func(map[string]any) (cicdOperationMutation, error) {
		return cicdOperationMutation{targetID: "svc", oneTime: secret}, nil
	}); err == nil || !reflect.DeepEqual(reply, cicdOperationReply{}) {
		t.Fatal("statefile.phase19-operation.one-time.no-reissue")
	}
}

func TestCICDOperationScope(t *testing.T) {
	for _, mode := range []string{"other-project", "lower-role"} {
		t.Run(mode, func(t *testing.T) {
			s, request := cicdCacheOperationFixture(t)
			_, head, _ := s.read()
			if _, err := s.commit(head.Generation, request.id, func(root map[string]any) error {
				if mode == "other-project" {
					project := cicdProjectFixture()
					project["id"] = "q"
					root["projects"] = append(cicdList(root["projects"]), project)
				} else {
					cicdObject(cicdList(root["grants"])[1])["role"] = "operator"
				}
				return nil
			}); err != nil {
				t.Fatal("statefile.phase19-operation.scope.fixture")
			}
			if mode == "other-project" {
				request.access.ProjectID = "q"
			}
			reply, err := s.updateCICD(request, nil)
			if err != security.ErrCICDForbidden || !reflect.DeepEqual(reply, cicdOperationReply{}) {
				t.Fatal("statefile.phase19-operation.scope.reject")
			}
		})
	}
}

func TestCICDOperationAuthorization(t *testing.T) {
	for _, mode := range []string{"check-missing", "check-revoked", "grant-revoked", "identity-disabled", "expired", "other-subject"} {
		t.Run(mode, func(t *testing.T) {
			s, request := cicdOperationFixture(t)
			if _, err := s.updateCICD(request, cicdCreateProject); err != nil {
				t.Fatal("statefile.phase19-operation.auth.fixture")
			}
			_, head, _ := s.read()
			if _, err := s.commit(head.Generation, request.id, func(root map[string]any) error {
				credential := map[string]any{"password_hash": strings.Repeat("a", 64), "salt": strings.Repeat("b", 64), "algorithm": "pbkdf2_hmac_sha256_v1", "iterations": json.Number("600000"), "must_change": false, "login_count": json.Number("0"), "last_login_at": nil, "updated_at": "2026-10-22T00:00:00Z"}
				root["users"] = append(cicdList(root["users"]), map[string]any{"id": "bob", "name": "Bob", "revision": json.Number("1"), "enabled": true, "credential": credential, "totp": nil, "created_at": "2026-10-22T00:00:00Z"})
				root["grants"] = append(cicdList(root["grants"]), map[string]any{"subject_kind": "user", "subject_id": "bob", "project_id": nil, "environment_id": nil, "role": "platform-admin", "enabled": true, "expires_at": nil})
				if mode == "grant-revoked" {
					cicdObject(cicdList(root["grants"])[0])["enabled"] = false
				}
				if mode == "identity-disabled" {
					cicdObject(cicdList(root["users"])[0])["enabled"] = false
				}
				return nil
			}); err != nil {
				t.Fatal("statefile.phase19-operation.auth.change.fixture")
			}
			want := "P19_AUTH"
			switch mode {
			case "check-missing":
				request.check = nil
			case "check-revoked":
				request.check = func() error { return errors.New("PRIVATE_SESSION_DETAIL") }
			case "grant-revoked":
				want = "P19_FORBIDDEN"
			case "expired":
				request.principal.ExpiresAt = s.clock()
			case "other-subject":
				request.principal.ID = "bob"
				want = "P19_IDEMPOTENCY"
			}
			_, before, _ := s.read()
			calls := 0
			reply, err := s.updateCICD(request, func(map[string]any) (cicdOperationMutation, error) { calls++; return cicdOperationMutation{}, nil })
			_, after, _ := s.read()
			if err == nil || err.Error() != want || !reflect.DeepEqual(reply, cicdOperationReply{}) || calls != 0 || before != after {
				t.Fatal("statefile.phase19-operation.auth.replay-denied")
			}
		})
	}
}

func TestCICDOperationRejectedMutation(t *testing.T) {
	for _, mode := range []string{"revision", "unknown-error", "bookkeeping", "generation", "bad-target"} {
		t.Run(mode, func(t *testing.T) {
			s, request := cicdOperationFixture(t)
			reply, err := s.updateCICD(request, func(root map[string]any) (cicdOperationMutation, error) {
				result, _ := cicdCreateProject(root)
				switch mode {
				case "revision":
					return result, ErrCICDRevision
				case "unknown-error":
					return result, errors.New("PRIVATE_ERROR_VALUE")
				case "bookkeeping":
					root["audits"] = []any{}
					root["operations"] = []any{map[string]any{"bad": true}}
				case "generation":
					root["generation"] = json.Number("2")
				case "bad-target":
					result.targetID = "not-created"
				}
				return result, nil
			})
			want := "P19_STATE_CORRUPT"
			if mode == "revision" {
				want = "P19_REVISION"
			}
			if mode == "unknown-error" {
				want = "P19_INTERNAL"
			}
			root, head, readErr := s.read()
			if err == nil || err.Error() != want || !reflect.DeepEqual(reply, cicdOperationReply{}) || readErr != nil || head.Generation != 2 || len(cicdList(root["operations"])) != 0 || len(cicdList(root["projects"])) != 0 || len(cicdList(root["audits"])) != 1 {
				t.Fatalf("statefile.phase19-operation.reject.atomic: %v", err)
			}
			audit := cicdObject(cicdList(root["audits"])[0])
			if audit["phase"] != "rejected" || audit["operation_id"] != request.requestID || audit["code"] != want || audit["target_id"] != nil {
				t.Fatal("statefile.phase19-operation.reject.audit")
			}
			request.requestID = "o-" + strings.Repeat("c", 32)
			if _, err := s.updateCICD(request, cicdCreateProject); err != nil {
				t.Fatal("statefile.phase19-operation.reject.key-not-accepted")
			}
		})
	}
}

func cicdCacheOperationFixture(t *testing.T) (*cicdSnapshotStore, cicdOperationRequest) {
	t.Helper()
	s, request := cicdOperationFixture(t)
	if _, err := s.commit(1, request.id, func(root map[string]any) error {
		root["projects"] = []any{cicdProjectFixture()}
		root["grants"] = append(cicdList(root["grants"]), map[string]any{"subject_kind": "user", "subject_id": "admin", "project_id": "p", "environment_id": nil, "role": "project-admin", "enabled": true, "expires_at": nil})
		root["caches"] = []any{map[string]any{"project_id": "p", "key": strings.Repeat("2", 64), "digest": strings.Repeat("1", 64), "size_bytes": json.Number("1"), "created_at": "2026-10-22T00:00:00Z", "last_used_at": "2026-10-22T00:00:00Z", "expires_at": "2026-10-23T00:00:00Z", "in_use": []any{}, "deleting": false}}
		return nil
	}); err != nil {
		t.Fatal("statefile.phase19-operation.cache.fixture")
	}
	request.access = security.CICDAccess{Operation: "cacheDelete", ProjectID: "p"}
	if _, err := s.updateCICD(request, func(root map[string]any) (cicdOperationMutation, error) {
		cicdObject(cicdList(root["caches"])[0])["deleting"] = true
		return cicdOperationMutation{targetID: strings.Repeat("2", 64)}, nil
	}); err != nil {
		t.Fatalf("statefile.phase19-operation.async.fixture: %v", err)
	}
	return s, request
}

func TestCICDOperationAdvance(t *testing.T) {
	s, request := cicdCacheOperationFixture(t)
	root, head, _ := s.read()
	if findCICDOperation(root, request.id)["state"] != "accepted" || len(cicdList(root["audits"])) != 1 {
		t.Fatal("statefile.phase19-operation.async.accepted")
	}
	calls := 0
	if _, err := s.advanceCICDOperation(request.id, head.Generation-1, func(map[string]any, map[string]any) error { calls++; return nil }); err != ErrCICDRevision || calls != 0 {
		t.Fatal("statefile.phase19-operation.advance.cas")
	}
	if _, err := s.advanceCICDOperation(request.id, head.Generation, func(_ map[string]any, operation map[string]any) error {
		operation["state"] = "succeeded"
		return nil
	}); err != ErrCICDCorrupt {
		t.Fatal("statefile.phase19-operation.advance.cache-not-removed")
	}
	_, unchanged, readErr := s.read()
	if readErr != nil || unchanged != head {
		t.Fatal("statefile.phase19-operation.advance.rejected-no-commit")
	}
	if _, err := s.advanceCICDOperation(request.id, head.Generation, func(_ map[string]any, operation map[string]any) error {
		operation["state"], operation["code"] = "recovery_required", "P19_TRANSPORT"
		return nil
	}); err != nil {
		t.Fatalf("statefile.phase19-operation.advance.recovery: %v", err)
	}
	root, head, _ = s.read()
	if len(cicdList(root["audits"])) != 1 || !validCICDReservedCapacity(root, cicdSnapshotLimit-262144) || validCICDReservedCapacity(root, cicdSnapshotLimit) {
		t.Fatal("statefile.phase19-operation.advance.reservation")
	}
	if _, err := s.advanceCICDOperation(request.id, head.Generation, func(_ map[string]any, operation map[string]any) error {
		operation["state"], operation["code"] = "failed", "P19_CLEANUP"
		return nil
	}); err != nil {
		t.Fatal("statefile.phase19-operation.advance.finished")
	}
	root, head, _ = s.read()
	if len(cicdList(root["audits"])) != 2 || !validCICDReservedCapacity(root, cicdSnapshotLimit) {
		t.Fatal("statefile.phase19-operation.advance.released")
	}
	if _, err := s.advanceCICDOperation(request.id, head.Generation, func(map[string]any, map[string]any) error { return nil }); err != nil {
		t.Fatal("statefile.phase19-operation.advance.identical")
	}
	_, after, _ := s.read()
	if head != after {
		t.Fatal("statefile.phase19-operation.advance.identical-no-generation")
	}
	if _, err := s.advanceCICDOperation(request.id, head.Generation, func(_ map[string]any, operation map[string]any) error {
		operation["state"], operation["code"] = "succeeded", nil
		return nil
	}); err != ErrCICDStale {
		t.Fatal("statefile.phase19-operation.advance.terminal-sealed")
	}
	if _, err := s.advanceCICDOperation(request.id, head.Generation, func(root map[string]any, _ map[string]any) error {
		cicdObject(cicdList(root["projects"])[0])["name"] = "Changed"
		return nil
	}); err != ErrCICDStale {
		t.Fatal("statefile.phase19-operation.advance.result-sealed")
	}
}

func TestCICDOperationGC(t *testing.T) {
	for _, mode := range []string{"before", "equal", "after", "active", "recovery", "protected"} {
		t.Run(mode, func(t *testing.T) {
			s, request := cicdCacheOperationFixture(t)
			_, head, _ := s.read()
			if mode != "active" {
				if _, err := s.advanceCICDOperation(request.id, head.Generation, func(root map[string]any, operation map[string]any) error {
					if mode == "recovery" {
						operation["state"], operation["code"] = "recovery_required", "P19_TRANSPORT"
					} else {
						operation["state"], operation["code"] = "failed", "P19_CLEANUP"
					}
					if mode == "protected" {
						cicdObject(cicdList(root["caches"])[0])["in_use"] = []any{request.id}
					}
					return nil
				}); err != nil {
					t.Fatal("statefile.phase19-operation.gc.terminal.fixture")
				}
			}
			at := time.Date(2026, 11, 21, 0, 0, 0, 0, time.UTC)
			if mode == "before" {
				at = at.Add(-time.Second)
			}
			if mode == "after" {
				at = at.Add(time.Second)
			}
			s.clock = func() time.Time { return at }
			count, err := s.gcCICDOperations("o-" + strings.Repeat("f", 32))
			want := 0
			if mode == "equal" || mode == "after" {
				want = 1
			}
			root, _, readErr := s.read()
			if err != nil || readErr != nil || count != want || len(cicdList(root["operations"])) != 1-want || want == 1 && len(cicdList(root["audits"])) != 0 {
				t.Fatalf("statefile.phase19-operation.gc.result: %v", err)
			}
		})
	}
}

func TestCICDOperationConcurrency(t *testing.T) {
	s, request := cicdOperationFixture(t)
	const count = 8
	start := make(chan struct{})
	var group sync.WaitGroup
	var calls atomic.Int64
	results := make(chan error, count)
	for i := 0; i < count; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			_, err := s.updateCICD(request, func(root map[string]any) (cicdOperationMutation, error) { calls.Add(1); return cicdCreateProject(root) })
			results <- err
		}()
	}
	close(start)
	group.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatalf("statefile.phase19-operation.concurrent.result: %v", err)
		}
	}
	root, head, err := s.read()
	if err != nil || calls.Load() != 1 || head.Generation != 2 || len(cicdList(root["operations"])) != 1 || len(cicdList(root["audits"])) != 2 {
		t.Fatal("statefile.phase19-operation.concurrent.once")
	}
}

func TestCICDOperationFaults(t *testing.T) {
	for _, stage := range []string{"snapshot.write", "snapshot.directory-sync", "head.publish", "head.directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			s, request := cicdOperationFixture(t)
			s.fault = func(actual string) error {
				if actual == stage {
					return errors.New("PRIVATE_IO_DETAIL")
				}
				return nil
			}
			reply, err := s.updateCICD(request, cicdCreateProject)
			want := ErrCICDStorage
			if stage == "head.directory-sync" {
				want = ErrCICDUncertain
			}
			if err != want || !reflect.DeepEqual(reply, cicdOperationReply{}) {
				t.Fatal("statefile.phase19-operation.fault.no-result")
			}
			fresh, reopenErr := newCICDSnapshotStore(s.path)
			if reopenErr != nil {
				t.Fatal("statefile.phase19-operation.fault.reopen")
			}
			fresh.clock = s.clock
			root, head, readErr := fresh.read()
			if readErr != nil {
				t.Fatal("statefile.phase19-operation.fault.read")
			}
			if stage == "head.directory-sync" {
				if head.Generation != 2 || findCICDOperation(root, request.id) == nil {
					t.Fatal("statefile.phase19-operation.fault.uncertain-visible")
				}
				calls := 0
				replay, err := fresh.updateCICD(request, func(map[string]any) (cicdOperationMutation, error) { calls++; return cicdOperationMutation{}, nil })
				if err != nil || !replay.replayed || calls != 0 {
					t.Fatal("statefile.phase19-operation.fault.reconcile-once")
				}
			} else if head.Generation != 1 || len(cicdList(root["projects"])) != 0 || len(cicdList(root["operations"])) != 0 || len(cicdList(root["audits"])) != 0 {
				t.Fatal("statefile.phase19-operation.fault.atomic")
			}
		})
	}
}

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
	const selection = "^TestCICD(Log(Schema|Store|Rejection|Faults|Corruption)|Schema|RootIntegrity|Run(Admission|Rejection|RequestCollision|Pins|ReadScope|Capacity)|Operation(Schema|Admission|OneTime|Authorization|RejectedMutation|Scope|Advance|GC|Concurrency|Faults)|Snapshot(Store|Faults|Corruption|Locks|NoReplace))$"
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
		{"session-check", "execute.go", "request.check == nil || request.check() != nil", "false", "statefile.phase19-operation.auth.replay-denied"},
		{"authorization", "execute.go", "_, err = authorizeCICDRoot(root, request.principal, request.access, now)", "err = nil", "statefile.phase19-operation.auth.replay-denied"},
		{"operation-digest", "execute.go", "operation[\"request_digest\"] != request.digest", "false", "statefile.phase19-operation.replay.conflict"},
		{"operation-replay", "execute.go", "if operation := findCICDOperation(root, request.id); operation != nil {", "if operation := findCICDOperation(root, request.id); false && operation != nil {", "statefile.phase19-operation.replay.before-revision"},
		{"one-time-reissue", "execute.go", "result.oneTime != \"\" && previousSubjects[result.targetID]", "false && result.oneTime != \"\" && previousSubjects[result.targetID]", "statefile.phase19-operation.one-time.no-reissue"},
		{"operation-cas", "execute.go", "if head.Generation != expectedGeneration {", "if false && head.Generation != expectedGeneration {", "statefile.phase19-operation.advance.cas"},
		{"operation-finished-audit", "execute.go", "if cicdTerminal(next[\"state\"]) {", "if false && cicdTerminal(next[\"state\"]) {", "statefile.phase19-operation.advance.finished"},
		{"operation-reserve", "validate.go", "state == \"accepted\" || state == \"recovery_required\"", "state == \"accepted\"", "statefile.phase19-operation.advance.reservation"},
		{"operation-gc-expiry", "execute.go", "now.Format(time.RFC3339) >= operation[\"expires_at\"].(string)", "true", "statefile.phase19-operation.gc.result"},
		{"operation-gc-active", "execute.go", "cicdTerminal(operation[\"state\"]) && now.Format(time.RFC3339)", "true && now.Format(time.RFC3339)", "statefile.phase19-operation.gc.result"},
		{"operation-gc-protection", "execute.go", "&& !protected[operation[\"id\"].(string)]", "&& true", "statefile.phase19-operation.gc.result"},
		{"operation-error-redaction", "execute.go", "return ErrCICDInternal", "return err", "statefile.phase19-operation.reject.atomic"},
		{"session-expiry", "security/execute.go", "!now.Before(principal.ExpiresAt)", "false", "statefile.phase19-operation.auth.replay-denied"},
		{"project-scope", "security/execute.go", "grant.ProjectID != access.ProjectID", "false", "statefile.phase19-operation.scope.reject"},
		{"run-project-revision", "execute.go", "if cicdInt(project[\"revision\"]) != expectedProjectRevision {", "if false && cicdInt(project[\"revision\"]) != expectedProjectRevision {", "statefile.phase19-run.reject.result"},
		{"run-plan-digest", "execute.go", "if plan[\"digest\"] != request.PlanDigest {", "if false && plan[\"digest\"] != request.PlanDigest {", "statefile.phase19-run.reject.result"},
		{"run-identity", "execute.go", "request.checkCurrent(identity) != nil", "false && request.checkCurrent(identity) != nil", "statefile.phase19-run.reject.result"},
		{"run-queue", "execute.go", "queued >= cicdInt(cicdObject(project[\"quota\"])[\"queue\"])", "false", "statefile.phase19-run.reject.result"},
		{"run-global-limit", "execute.go", "active >= 16", "active >= 17", "statefile.phase19-run.capacity.boundary"},
		{"run-pin", "execute.go", "artifact[\"pins\"] = append(pins, id)", "artifact[\"pins\"] = pins", "statefile.phase19-run.pin.atomic"},
		{"run-entropy", "execute.go", "if _, err := io.ReadFull(reader, entropy[:]); err != nil {", "if _, err := io.ReadFull(reader, entropy[:]); false && err != nil {", "statefile.phase19-run.reject.result"},
		{"run-request-id", "execute.go", "if !cicdRequestIDAvailable(root, request.requestID) {", "if false && !cicdRequestIDAvailable(root, request.requestID) {", "statefile.phase19-run.request-id.admission-reject"},
		{"log-length", "validate.go", "next-first != int64(len(log[\"data\"].(string)))", "false", "statefile.phase19-log.schema.reject"},
		{"log-terminal", "validate.go", "if old[\"terminal\"] == true {", "if false && old[\"terminal\"] == true {", "statefile.phase19-log.reject.result"},
		{"log-overlap", "validate.go", "return a[start-cicdInt(old[\"first_offset\"]):end-cicdInt(old[\"first_offset\"])] == b[:end-start]", "return len(a)+len(b) >= 0", "statefile.phase19-log.reject.result"},
		{"log-gap", "validate.go", "old[\"gap\"] == true && next[\"gap\"] != true", "false", "statefile.phase19-log.reject.result"},
		{"log-selected", "execute.go", "if !selected {", "if false && !selected {", "statefile.phase19-log.reject.result"},
		{"log-first", "execute.go", "cicdInt(log[\"next_offset\"]) != 0", "false", "statefile.phase19-log.reject.result"},
		{"log-replay", "execute.go", "if bytes.Equal(cicdValueBytes(old), cicdValueBytes(log)) {", "if false && bytes.Equal(cicdValueBytes(old), cicdValueBytes(log)) {", "statefile.phase19-log.store.noop"},
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
			if os.MkdirAll(filepath.Join(copyDir, "components", "security"), 0700) != nil {
				t.Fatal("statefile.phase19-mutation.security.fixture")
			}
			for _, name := range []string{"security.go", "model.go", "validate.go", "execute.go", "security_test.go"} {
				data, err := os.ReadFile(filepath.Join("..", "security", name))
				if mutation.path == "security/"+name {
					if strings.Count(string(data), mutation.old) != 1 {
						t.Fatal("statefile.phase19-mutation.security.target.fixture")
					}
					data = []byte(strings.Replace(string(data), mutation.old, mutation.new, 1))
				}
				if err != nil || os.WriteFile(filepath.Join(copyDir, "components", "security", name), data, 0600) != nil {
					t.Fatal("statefile.phase19-mutation.security.copy")
				}
			}
			if os.WriteFile(filepath.Join(copyDir, "go.mod"), []byte("module github.com/fqwink/build-scripts\n\ngo 1.22\n"), 0600) != nil {
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
