package runner

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/fqwink/build-scripts/components/security"
	"github.com/fqwink/build-scripts/components/statefile"
)

type cicdRunTestStore struct {
	snapshot statefile.CICDPlanSnapshot
	check    func() (statefile.CICDRunReply, error)
	commit   func(statefile.CICDRunRequest, int64, func(statefile.CICDPlanSnapshot) (json.RawMessage, error)) (statefile.CICDRunReply, error)
}

func (s *cicdRunTestStore) CheckRunStart(statefile.CICDRunRequest) (statefile.CICDPlanSnapshot, statefile.CICDRunReply, error) {
	if s.check != nil {
		reply, err := s.check()
		return s.snapshot, reply, err
	}
	return s.snapshot, statefile.CICDRunReply{}, nil
}

func (s *cicdRunTestStore) StartRun(request statefile.CICDRunRequest, revision int64, resolve func(statefile.CICDPlanSnapshot) (json.RawMessage, error)) (statefile.CICDRunReply, error) {
	return s.commit(request, revision, resolve)
}

func TestCICDRunStart(t *testing.T) {
	for _, name := range []string{"accepted", "replay", "denied", "source-error", "nil-source", "invalid-values", "stale-preview", "latest-worker", "latest-template", "values-copy"} {
		t.Run(name, func(t *testing.T) {
			p, d, source, context := cicdPlanTestFixture(t)
			snapshot := statefile.CICDPlanSnapshot{Generation: 1, Project: cicdPlanTestRaw(t, p), Definition: cicdPlanTestRaw(t, d), Workers: context.Workers}
			draft := cicdPlanTestDraft(t, p, d, `{}`)
			plan, err := draft.Resolve(cicdPlanTestRaw(t, source), nil, context)
			if err != nil {
				t.Fatal("runner.phase19-run.fixture.plan")
			}
			object := cicdPlanTestObject(t, string(plan))
			request := statefile.CICDRunRequest{ProjectID: "p", DefinitionID: "d", DefinitionRevision: 1, PlanDigest: object["digest"].(string)}
			store := &cicdRunTestStore{snapshot: snapshot}
			fetches, commits, resolutions := 0, 0, 0
			values := json.RawMessage(`{}`)
			var wantError error
			switch name {
			case "replay":
				store.check = func() (statefile.CICDRunReply, error) {
					return statefile.CICDRunReply{Replayed: true, Operation: json.RawMessage(`{"state":"accepted"}`)}, nil
				}
			case "denied":
				wantError = security.ErrCICDForbidden
				store.check = func() (statefile.CICDRunReply, error) { return statefile.CICDRunReply{}, wantError }
			case "source-error":
				wantError = ErrCICDSource
			case "nil-source", "invalid-values":
				wantError = ErrCICDSchema
				if name == "invalid-values" {
					values = json.RawMessage(`{"unknown":true}`)
				}
			case "stale-preview":
				request.PlanDigest, wantError = strings.Repeat("f", 64), statefile.ErrCICDStale
			case "latest-worker":
				wantError = ErrCICDCapability
			case "latest-template":
				wantError = ErrCICDSchema
			}
			store.commit = func(actual statefile.CICDRunRequest, revision int64, resolve func(statefile.CICDPlanSnapshot) (json.RawMessage, error)) (statefile.CICDRunReply, error) {
				commits++
				if actual.PlanDigest != request.PlanDigest || revision != 1 {
					t.Fatal("runner.phase19-run.commit.request")
				}
				current := snapshot
				if name == "latest-worker" {
					worker := cicdPlanTestObject(t, string(context.Workers[0]))
					worker["enabled"] = false
					current.Workers = []json.RawMessage{cicdPlanTestRaw(t, worker)}
				} else if name == "latest-template" {
					current.Definition = json.RawMessage(`{}`)
				}
				resolutions++
				got, err := resolve(current)
				if err != nil {
					return statefile.CICDRunReply{}, err
				}
				if !bytes.Equal(got, plan) {
					t.Fatal("runner.phase19-run.commit.frozen-plan")
				}
				return statefile.CICDRunReply{Generation: 2, Operation: json.RawMessage(`{"state":"accepted"}`)}, nil
			}
			got, err := StartCICDRun(store, request, values, func(actual statefile.CICDPlanSnapshot) (*CICDSourceSnapshot, error) {
				fetches++
				if commits != 0 || !bytes.Equal(actual.Project, snapshot.Project) {
					t.Fatal("runner.phase19-run.source.order")
				}
				if name == "source-error" {
					return nil, wantError
				}
				if name == "nil-source" {
					return nil, nil
				}
				if name == "values-copy" {
					values[0] = '!'
				}
				return &CICDSourceSnapshot{metadata: cicdPlanTestRaw(t, source), files: map[string][]byte{}}, nil
			})
			if err != wantError || name == "replay" && !got.Replayed || wantError != nil && got.Operation != nil {
				t.Fatal("runner.phase19-run.start.result", err)
			}
			if name == "replay" || name == "denied" {
				if fetches != 0 || commits != 0 || resolutions != 0 {
					t.Fatal("runner.phase19-run.replay.no-effects")
				}
			} else if commits != 1 || resolutions != 1 || name == "invalid-values" && fetches != 0 || name != "invalid-values" && fetches != 1 {
				t.Fatal("runner.phase19-run.start.call-count")
			}
		})
	}
}

func TestCICDRunRealStore(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skip("P19_CAPABILITY: run admission filesystem requires Linux amd64")
	}
	p, d, source, context := cicdPlanTestFixture(t)
	root := cicdPlanTestObject(t, `{"schema_version":1,"generation":1,"projects":[],"definitions":[],"templates":[],"environments":[],"users":[{"id":"admin","name":"Administrator","revision":1,"enabled":true,"credential":null,"totp":null,"created_at":"2026-10-22T00:00:00Z"}],"services":[],"secrets":[],"grants":[{"subject_kind":"user","subject_id":"admin","project_id":null,"environment_id":null,"role":"platform-admin","enabled":true,"expires_at":null},{"subject_kind":"user","subject_id":"admin","project_id":"p","environment_id":null,"role":"operator","enabled":true,"expires_at":null}],"workers":[],"runs":[],"deployments":[],"artifacts":[],"caches":[],"operations":[],"audits":[],"targets":[],"environment_states":[],"notifications":[],"trigger_cursors":[],"status_deliveries":[]}`)
	root["projects"], root["definitions"], root["workers"] = []any{p}, []any{d}, []any{cicdPlanTestObject(t, string(context.Workers[0]))}
	state := t.TempDir()
	if os.Chmod(state, 0700) != nil || os.MkdirAll(filepath.Join(state, "cicd", "snapshots"), 0700) != nil {
		t.Fatal("runner.phase19-run.real.directory")
	}
	rootBytes := cicdPlanTestRaw(t, root)
	head := cicdPlanTestRaw(t, map[string]any{"schema_version": 1, "generation": 1, "sha256": fmt.Sprintf("%x", sha256.Sum256(rootBytes))})
	for name, data := range map[string][]byte{"snapshots/g-1.json": rootBytes, "head.json": head, "state.lock": {}} {
		if os.WriteFile(filepath.Join(state, "cicd", name), data, 0600) != nil {
			t.Fatal("runner.phase19-run.real.fixture")
		}
	}
	store, err := statefile.OpenCICDRunStore(state)
	if err != nil {
		t.Fatal("runner.phase19-run.real.open", err)
	}
	plan, err := cicdPlanTestDraft(t, p, d, `{}`).Resolve(cicdPlanTestRaw(t, source), nil, context)
	if err != nil {
		t.Fatal("runner.phase19-run.real.preview")
	}
	request := statefile.CICDRunRequest{Principal: security.CICDPrincipal{Kind: "user", ID: "admin", ExpiresAt: time.Now().UTC().Truncate(time.Second).Add(time.Hour), MustChange: "none"}, ProjectID: "p", DefinitionID: "d", DefinitionRevision: 1, OperationID: "o-" + strings.Repeat("a", 32), RequestID: "o-" + strings.Repeat("b", 32), RequestDigest: strings.Repeat("c", 64), PlanDigest: cicdPlanTestObject(t, string(plan))["digest"].(string), Check: func(current statefile.CICDIdentity) error {
		if current.ID != "admin" || current.Revision != 1 {
			return security.ErrCICDAuth
		}
		return nil
	}}
	fetches := 0
	acquire := func(statefile.CICDPlanSnapshot) (*CICDSourceSnapshot, error) {
		fetches++
		return &CICDSourceSnapshot{metadata: cicdPlanTestRaw(t, source), files: map[string][]byte{}}, nil
	}
	first, err := StartCICDRun(store, request, json.RawMessage(`{}`), acquire)
	if err != nil || first.Replayed || first.Generation != 2 || fetches != 1 {
		t.Fatal("runner.phase19-run.real.accept", err)
	}
	reopened, err := statefile.OpenCICDRunStore(state)
	if err != nil {
		t.Fatal("runner.phase19-run.real.reopen", err)
	}
	replay, err := StartCICDRun(reopened, request, json.RawMessage(`{}`), acquire)
	if err != nil || !replay.Replayed || !bytes.Equal(first.Operation, replay.Operation) || replay.Generation != 2 || fetches != 1 {
		t.Fatal("runner.phase19-run.real.replay", err)
	}
	request.Check = func(statefile.CICDIdentity) error { return security.ErrCICDAuth }
	denied, err := StartCICDRun(reopened, request, json.RawMessage(`{}`), acquire)
	if err != security.ErrCICDAuth || denied.Operation != nil || fetches != 1 {
		t.Fatal("runner.phase19-run.real.revoked", err)
	}
}

type cicdSourceTestTransport func(*http.Request) (*http.Response, error)

func (transport cicdSourceTestTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

type cicdSourceTestBody struct {
	io.Reader
	read   int
	closed bool
}

func (body *cicdSourceTestBody) Read(data []byte) (int, error) {
	n, err := body.Reader.Read(data)
	body.read += n
	return n, err
}

func (body *cicdSourceTestBody) Close() error { body.closed = true; return nil }

func cicdSourceTestResponse(data []byte) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(bytes.NewReader(data)), ContentLength: int64(len(data))}
}

func cicdSourceTestRequest() CICDSourceRequest {
	return CICDSourceRequest{Owner: "owner", Repo: "repo", RefKind: "branch", Ref: "feature/next", SHA: strings.Repeat("b", 40)}
}

func cicdSourceTestFixture(t *testing.T) (map[string]any, map[string]any) {
	t.Helper()
	tree := cicdPlanTestObject(t, `{"sha":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","truncated":false,"tree":[{"path":"hello.txt","mode":"100644","type":"blob","size":6,"sha":"ce013625030ba8dba906f756967f9e9ca394464a"},{"path":"empty.txt","mode":"100755","type":"blob","size":0,"sha":"e69de29bb2d1d6434b8b29ae775ad8c2e48c5391"},{"path":"alpha.bin","mode":"100644","type":"blob","size":3,"sha":"506cd141ad4a679eee22d6a21dd267cca5734b92"}]}`)
	blobs := map[string]any{
		"ce013625030ba8dba906f756967f9e9ca394464a": cicdPlanTestObject(t, `{"sha":"ce013625030ba8dba906f756967f9e9ca394464a","size":6,"encoding":"base64","content":"aGVs\nbG8K\r\n","url":"https://ignored.invalid/secret"}`),
		"e69de29bb2d1d6434b8b29ae775ad8c2e48c5391": cicdPlanTestObject(t, `{"sha":"e69de29bb2d1d6434b8b29ae775ad8c2e48c5391","size":0,"encoding":"base64","content":""}`),
		"506cd141ad4a679eee22d6a21dd267cca5734b92": cicdPlanTestObject(t, `{"sha":"506cd141ad4a679eee22d6a21dd267cca5734b92","size":3,"encoding":"base64","content":"AP8K"}`),
	}
	return tree, blobs
}

func cicdSourceTestClient(t *testing.T, tree map[string]any, blobs map[string]any, calls *[]string) *http.Client {
	t.Helper()
	return cicdSourceClient(cicdSourceTestTransport(func(request *http.Request) (*http.Response, error) {
		*calls = append(*calls, request.URL.RequestURI())
		deadline, ok := request.Context().Deadline()
		if request.Method != http.MethodGet || request.URL.Scheme != "https" || request.URL.Host != "api.github.com" || request.URL.User != nil ||
			request.Header.Get("Authorization") != "Bearer fixture-token" || request.Header.Get("Accept") != "application/vnd.github+json" ||
			request.Header.Get("X-GitHub-Api-Version") != "2022-11-28" || request.Header.Get("Accept-Encoding") != "identity" || !ok || time.Until(deadline) > 10*time.Second {
			t.Fatal("runner.phase19-source.http.request")
		}
		if request.URL.RequestURI() == "/repos/owner/repo/git/trees/"+strings.Repeat("b", 40)+"?recursive=1" {
			return cicdSourceTestResponse(cicdPlanTestRaw(t, tree)), nil
		}
		value, exists := blobs[strings.TrimPrefix(request.URL.Path, "/repos/owner/repo/git/blobs/")]
		if !exists {
			t.Fatal("runner.phase19-source.http.unexpected")
		}
		return cicdSourceTestResponse(cicdPlanTestRaw(t, value)), nil
	}))
}

func TestCICDSource(t *testing.T) {
	tree, blobs := cicdSourceTestFixture(t)
	before := cicdPlanTestRaw(t, tree)
	var calls []string
	snapshot, err := cicdFetchSource(context.Background(), cicdSourceTestRequest(), "fixture-token", cicdSourceTestClient(t, tree, blobs, &calls))
	if err != nil || snapshot == nil {
		t.Fatal("runner.phase19-source.fetch.accept", err)
	}
	wantFiles := cicdPlanTestObject(t, `{"files":[{"path":"alpha.bin","blob_sha":"506cd141ad4a679eee22d6a21dd267cca5734b92","sha256":"712450d3c4a79eea9509e75dc1dacdeff58034df538536cfae2da882bd8a0c50","size":3,"executable":false},{"path":"empty.txt","blob_sha":"e69de29bb2d1d6434b8b29ae775ad8c2e48c5391","sha256":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","size":0,"executable":true},{"path":"hello.txt","blob_sha":"ce013625030ba8dba906f756967f9e9ca394464a","sha256":"5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03","size":6,"executable":false}]}`)["files"]
	digest := sha256.Sum256(cicdPlanTestRaw(t, wantFiles))
	want := cicdPlanTestRaw(t, map[string]any{"owner": "owner", "repo": "repo", "ref_kind": "branch", "ref": "feature/next", "sha": strings.Repeat("b", 40), "files": wantFiles, "tree_digest": hex.EncodeToString(digest[:])})
	if !bytes.Equal(snapshot.Metadata(), want) || !bytes.Equal(before, cicdPlanTestRaw(t, tree)) || !reflect.DeepEqual(calls, []string{
		"/repos/owner/repo/git/trees/" + strings.Repeat("b", 40) + "?recursive=1",
		"/repos/owner/repo/git/blobs/506cd141ad4a679eee22d6a21dd267cca5734b92",
		"/repos/owner/repo/git/blobs/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
		"/repos/owner/repo/git/blobs/ce013625030ba8dba906f756967f9e9ca394464a",
	}) {
		t.Fatal("runner.phase19-source.snapshot.value")
	}
	for name, expected := range map[string][]byte{"alpha.bin": {0, 255, 10}, "empty.txt": {}, "hello.txt": []byte("hello\n")} {
		got, err := snapshot.ReadFile(name)
		if err != nil || !bytes.Equal(got, expected) {
			t.Fatal("runner.phase19-source.file.value")
		}
		if len(got) > 0 {
			got[0] ^= 255
		}
		fresh, _ := snapshot.ReadFile(name)
		if !bytes.Equal(fresh, expected) {
			t.Fatal("runner.phase19-source.file.alias")
		}
	}
	metadata := snapshot.Metadata()
	metadata[0] = '!'
	if !bytes.Equal(snapshot.Metadata(), want) {
		t.Fatal("runner.phase19-source.metadata.alias")
	}
	if data, err := snapshot.ReadFile("unknown"); data != nil || err != ErrCICDNotFound {
		t.Fatal("runner.phase19-source.file.missing")
	}
	for _, empty := range []*CICDSourceSnapshot{nil, {}} {
		if empty.Metadata() != nil {
			t.Fatal("runner.phase19-source.zero.metadata")
		}
		if data, err := empty.ReadFile("hello.txt"); data != nil || err != ErrCICDSource {
			t.Fatal("runner.phase19-source.zero.read")
		}
	}
	p, d, _, projection := cicdPlanTestFixture(t)
	plan, err := cicdPlanTestDraft(t, p, d, `{}`).Resolve(snapshot.Metadata(), nil, projection)
	if err != nil || plan == nil {
		t.Fatal("runner.phase19-source.plan.connect")
	}
}

func TestCICDSourceInput(t *testing.T) {
	for _, change := range []func(*CICDSourceRequest, *string){
		func(r *CICDSourceRequest, _ *string) { r.Owner = "../owner" },
		func(r *CICDSourceRequest, _ *string) { r.Owner = "a_b" },
		func(r *CICDSourceRequest, _ *string) { r.Repo = "repo?key=x" },
		func(r *CICDSourceRequest, _ *string) { r.Repo = strings.Repeat("a", 101) },
		func(r *CICDSourceRequest, _ *string) { r.RefKind = "commit" },
		func(r *CICDSourceRequest, _ *string) { r.Ref = "../main" },
		func(r *CICDSourceRequest, _ *string) { r.SHA = "main" },
		func(r *CICDSourceRequest, _ *string) { r.SHA = strings.Repeat("B", 40) },
		func(_ *CICDSourceRequest, token *string) { *token = "" },
		func(_ *CICDSourceRequest, token *string) { *token = "space token" },
		func(_ *CICDSourceRequest, token *string) { *token = "token\r\n" },
		func(_ *CICDSourceRequest, token *string) { *token = "日本語" },
		func(_ *CICDSourceRequest, token *string) { *token = strings.Repeat("a", 4097) },
	} {
		request, token := cicdSourceTestRequest(), "fixture-token"
		change(&request, &token)
		got, err := FetchCICDSource(context.Background(), request, token)
		if got != nil || err != ErrCICDSchema {
			t.Fatal("runner.phase19-source.input.reject")
		}
	}
	if got, err := FetchCICDSource(nil, cicdSourceTestRequest(), "fixture-token"); got != nil || err != ErrCICDSchema {
		t.Fatal("runner.phase19-source.input.nil-context")
	}
	request := cicdSourceTestRequest()
	request.Owner, request.Repo, request.RefKind, request.Ref = "a-1", "repo_1.git", "tag", "v1.0"
	if !cicdSourceRequestValid(request, strings.Repeat("a", 4096)) {
		t.Fatal("runner.phase19-source.input.boundary")
	}
}

func TestCICDSourceTree(t *testing.T) {
	for _, test := range []struct {
		name string
		edit func(map[string]any)
		err  error
	}{
		{"truncated", func(v map[string]any) { v["truncated"] = true }, ErrCICDSource},
		{"missing-truncated", func(v map[string]any) { delete(v, "truncated") }, ErrCICDTransport},
		{"null-truncated", func(v map[string]any) { v["truncated"] = nil }, ErrCICDTransport},
		{"missing-tree", func(v map[string]any) { delete(v, "tree") }, ErrCICDTransport},
		{"null-tree", func(v map[string]any) { v["tree"] = nil }, ErrCICDTransport},
		{"missing-sha", func(v map[string]any) { delete(v, "sha") }, ErrCICDTransport},
		{"entry-null", func(v map[string]any) { v["tree"].([]any)[0] = nil }, ErrCICDSource},
		{"path-traversal", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["path"] = "../file" }, ErrCICDSource},
		{"path-absolute", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["path"] = "/file" }, ErrCICDSource},
		{"path-duplicate", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["path"] = "empty.txt" }, ErrCICDSource},
		{"missing-parent", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["path"] = "missing/file" }, ErrCICDSource},
		{"file-parent", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["path"] = "empty.txt/file" }, ErrCICDSource},
		{"symlink", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["mode"] = "120000" }, ErrCICDSource},
		{"submodule", func(v map[string]any) {
			e := v["tree"].([]any)[0].(map[string]any)
			e["mode"], e["type"] = "160000", "commit"
		}, ErrCICDSource},
		{"unknown-mode", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["mode"] = "100664" }, ErrCICDSource},
		{"unknown-type", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["type"] = "tree" }, ErrCICDSource},
		{"missing-size", func(v map[string]any) { delete(v["tree"].([]any)[0].(map[string]any), "size") }, ErrCICDSource},
		{"negative-size", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["size"] = -1 }, ErrCICDSource},
		{"large-file", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["size"] = 10485761 }, ErrCICDSource},
		{"bad-sha", func(v map[string]any) { v["tree"].([]any)[0].(map[string]any)["sha"] = "bad" }, ErrCICDSource},
	} {
		t.Run(test.name, func(t *testing.T) {
			tree, blobs := cicdSourceTestFixture(t)
			test.edit(tree)
			var calls []string
			got, err := cicdFetchSource(context.Background(), cicdSourceTestRequest(), "fixture-token", cicdSourceTestClient(t, tree, blobs, &calls))
			if got != nil || err != test.err || len(calls) != 1 {
				t.Fatal("runner.phase19-source.tree.reject", err)
			}
		})
	}
	for _, size := range []int{10000, 10001} {
		tree, _ := cicdSourceTestFixture(t)
		entry := tree["tree"].([]any)[0].(map[string]any)
		entries := make([]any, size)
		for i := range entries {
			copy := cicdPlanClone(entry)
			copy["path"] = fmt.Sprintf("file-%05d", i)
			entries[i] = copy
		}
		tree["tree"] = entries
		files, err := cicdSourceTree(tree)
		if size == 10000 && (err != nil || len(files) != size) || size == 10001 && (files != nil || err != ErrCICDSource) {
			t.Fatal("runner.phase19-source.tree.file-limit")
		}
	}
	for _, total := range []int{104857600, 104857601} {
		tree, _ := cicdSourceTestFixture(t)
		entry := tree["tree"].([]any)[0].(map[string]any)
		entries := make([]any, 11)
		for i := range entries {
			copy := cicdPlanClone(entry)
			copy["path"], copy["size"] = fmt.Sprintf("file-%02d", i), json.Number("10485760")
			if i == 10 {
				copy["size"] = json.Number(fmt.Sprint(total - 104857600))
			}
			entries[i] = copy
		}
		tree["tree"] = entries
		got, err := cicdSourceTree(tree)
		if total == 104857600 && (err != nil || len(got) != 11) || total > 104857600 && (got != nil || err != ErrCICDSource) {
			t.Fatal("runner.phase19-source.tree.total-limit")
		}
	}
	tree, blobs := cicdSourceTestFixture(t)
	entries := tree["tree"].([]any)
	entries[0].(map[string]any)["path"] = "dir/hello.txt"
	tree["tree"] = append(entries, map[string]any{"path": "dir", "mode": "040000", "type": "tree", "sha": strings.Repeat("a", 40), "size": nil})
	var calls []string
	if got, err := cicdFetchSource(context.Background(), cicdSourceTestRequest(), "fixture-token", cicdSourceTestClient(t, tree, blobs, &calls)); err != nil || got == nil {
		t.Fatal("runner.phase19-source.tree.directory")
	}
	tree["tree"] = []any{}
	if got, err := cicdFetchSource(context.Background(), cicdSourceTestRequest(), "fixture-token", cicdSourceTestClient(t, tree, nil, &calls)); err != nil || got == nil || len(got.files) != 0 {
		t.Fatal("runner.phase19-source.tree.empty")
	}
}

func TestCICDSourceBlob(t *testing.T) {
	for _, edit := range []func(map[string]any){
		func(v map[string]any) { delete(v, "content") },
		func(v map[string]any) { v["content"] = nil },
		func(v map[string]any) { v["content"] = "aGVs bG8K" },
		func(v map[string]any) { v["content"] = "aGVs\rbG8K" },
		func(v map[string]any) { v["content"] = "aGVsbG8=" },
		func(v map[string]any) { v["content"] = "aGVsbG8J" },
		func(v map[string]any) { v["content"] = "aGVsbG8_" },
		func(v map[string]any) { v["encoding"] = "utf-8" },
		func(v map[string]any) { v["sha"] = strings.Repeat("a", 40) },
		func(v map[string]any) { v["size"] = json.Number("5") },
		func(v map[string]any) { v["size"] = nil },
	} {
		tree, blobs := cicdSourceTestFixture(t)
		edit(blobs["ce013625030ba8dba906f756967f9e9ca394464a"].(map[string]any))
		var calls []string
		got, err := cicdFetchSource(context.Background(), cicdSourceTestRequest(), "fixture-token", cicdSourceTestClient(t, tree, blobs, &calls))
		if got != nil || err != ErrCICDDigest || len(calls) != 4 {
			t.Fatal("runner.phase19-source.blob.reject")
		}
	}
	file := map[string]any{"blob_sha": "4d1ae35ba2c8ec712fa2a379db44ad639ca277bd", "size": json.Number("1")}
	blob := map[string]any{"sha": file["blob_sha"], "size": file["size"], "encoding": "base64", "content": "Zh=="}
	if got, err := cicdSourceBlob(file, blob); got != nil || err != ErrCICDDigest {
		t.Fatal("runner.phase19-source.blob.pad-bits")
	}
	blob["content"] = "Zg=="
	if got, err := cicdSourceBlob(file, blob); err != nil || string(got) != "f" {
		t.Fatal("runner.phase19-source.blob.padding")
	}
}

func TestCICDSourceHTTP(t *testing.T) {
	transport := cicdSourceTransport()
	defer transport.CloseIdleConnections()
	client := cicdSourceClient(transport)
	if transport.Proxy != nil || !transport.DisableKeepAlives || !transport.DisableCompression || transport.TLSNextProto == nil || len(transport.TLSNextProto) != 0 ||
		transport.ForceAttemptHTTP2 || transport.DialContext == nil || transport.TLSHandshakeTimeout != 10*time.Second || transport.ResponseHeaderTimeout != 10*time.Second || client.Timeout != 10*time.Second {
		t.Fatal("runner.phase19-source.http.transport-profile")
	}
	var mu sync.Mutex
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		requests++
		if !r.Close {
			t.Error("runner.phase19-source.http.connection-reuse")
		}
		if requests == 1 {
			_, _ = io.WriteString(w, `{}`)
			return
		}
		connection, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error("runner.phase19-source.http.hijack")
			return
		}
		_ = connection.Close()
	}))
	defer server.Close()
	response, err := client.Get(server.URL)
	if err != nil {
		t.Fatal("runner.phase19-source.http.initial")
	}
	_, _ = io.Copy(io.Discard, response.Body)
	_ = response.Body.Close()
	response, err = client.Get(server.URL)
	if response != nil {
		_ = response.Body.Close()
	}
	mu.Lock()
	count := requests
	mu.Unlock()
	if err == nil || response != nil || count != 2 {
		t.Fatal("runner.phase19-source.http.implicit-retry")
	}
	for _, test := range []struct {
		name, body      string
		status          int
		media, encoding string
		length          int64
		failure         bool
		want            error
	}{
		{"redirect", `{}`, 302, "application/json", "", 2, false, ErrCICDTransport},
		{"unauthorized", `secret-response`, 401, "application/json", "", 15, false, ErrCICDTransport},
		{"rate", `{}`, 429, "application/json", "", 2, false, ErrCICDTransport},
		{"media", `{}`, 200, "text/plain", "", 2, false, ErrCICDTransport},
		{"gzip", `{}`, 200, "application/json", "gzip", 2, false, ErrCICDTransport},
		{"json", `{`, 200, "application/json", "", 1, false, ErrCICDTransport},
		{"duplicate", `{"tree":[],"tree":[]}`, 200, "application/json", "", -1, false, ErrCICDTransport},
		{"trailing", `{} {}`, 200, "application/json", "", -1, false, ErrCICDTransport},
		{"null", `null`, 200, "application/json", "", -1, false, ErrCICDTransport},
		{"surrogate", `{"x":"\ud800"}`, 200, "application/json", "", -1, false, ErrCICDTransport},
		{"utf8", "{\"x\":\"\xff\"}", 200, "application/json", "", -1, false, ErrCICDTransport},
		{"fraction", `{"x":1.0}`, 200, "application/json", "", -1, false, ErrCICDTransport},
		{"exponent", `{"x":1e0}`, 200, "application/json", "", -1, false, ErrCICDTransport},
		{"declared-large", `{}`, 200, "application/json", "", 65, false, ErrCICDBody},
		{"actual-large", strings.Repeat("x", 100), 200, "application/json", "", -1, false, ErrCICDBody},
		{"transport", `secret-response`, 200, "application/json", "", -1, true, ErrCICDTransport},
		{"exact", "{}" + strings.Repeat(" ", 62), 200, "application/json; charset=utf-8", "", 64, false, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			body := &cicdSourceTestBody{Reader: strings.NewReader(test.body)}
			calls := 0
			client := cicdSourceClient(cicdSourceTestTransport(func(request *http.Request) (*http.Response, error) {
				calls++
				if test.failure {
					return nil, errors.New("private-token-and-url")
				}
				return &http.Response{StatusCode: test.status, Header: http.Header{"Content-Type": {test.media}, "Content-Encoding": {test.encoding}, "Location": {"https://other.invalid"}}, Body: body, ContentLength: test.length}, nil
			}))
			got, err := cicdSourceGet(context.Background(), client, "/repos/owner/repo/git/trees/"+strings.Repeat("b", 40), "fixture-token", 64)
			if err != test.want || (test.want != nil && got != nil) || calls != 1 || body.read > 65 || (!test.failure && !body.closed) || (test.name == "declared-large" && body.read != 0) {
				t.Fatal("runner.phase19-source.http.result")
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	started, released := make(chan struct{}), make(chan struct{})
	cancelClient := cicdSourceClient(cicdSourceTestTransport(func(request *http.Request) (*http.Response, error) {
		close(started)
		<-request.Context().Done()
		close(released)
		return nil, errors.New("private-cancel-detail")
	}))
	done := make(chan error, 1)
	go func() {
		got, err := cicdFetchSource(ctx, cicdSourceTestRequest(), "fixture-token", cancelClient)
		if got != nil {
			done <- errors.New("partial")
			return
		}
		done <- err
	}()
	<-started
	cancel()
	if err := <-done; err != ErrCICDTransport {
		t.Fatal("runner.phase19-source.http.cancel")
	}
	<-released
	for _, cancelled := range []bool{false, true} {
		body := &cicdSourceTestBody{Reader: iotest.ErrReader(errors.New("private-read-detail"))}
		calls := 0
		client := cicdSourceClient(cicdSourceTestTransport(func(*http.Request) (*http.Response, error) {
			calls++
			response := cicdSourceTestResponse(nil)
			response.Body = body
			return response, nil
		}))
		ctx, cancel := context.WithCancel(context.Background())
		if cancelled {
			cancel()
		}
		got, err := cicdFetchSource(ctx, cicdSourceTestRequest(), "fixture-token", client)
		cancel()
		if got != nil || err != ErrCICDTransport || (!cancelled && (calls != 1 || !body.closed)) || (cancelled && calls != 0) {
			t.Fatal("runner.phase19-source.http.read-failure")
		}
	}
}

func TestCICDSourceConcurrent(t *testing.T) {
	tree, blobs := cicdSourceTestFixture(t)
	var calls []string
	snapshot, err := cicdFetchSource(context.Background(), cicdSourceTestRequest(), "fixture-token", cicdSourceTestClient(t, tree, blobs, &calls))
	if err != nil {
		t.Fatal("runner.phase19-source.concurrent.prepare")
	}
	metadata := snapshot.Metadata()
	start := make(chan struct{})
	var wait sync.WaitGroup
	for i := 0; i < 8; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			for j := 0; j < 16; j++ {
				data, err := snapshot.ReadFile("hello.txt")
				if err != nil || string(data) != "hello\n" || !bytes.Equal(snapshot.Metadata(), metadata) {
					t.Error("runner.phase19-source.concurrent.value")
					return
				}
				data[0] = '!'
			}
		}()
	}
	close(start)
	wait.Wait()
}

func cicdPlanTestObject(t *testing.T, raw string) map[string]any {
	t.Helper()
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	var result map[string]any
	if decoder.Decode(&result) != nil {
		t.Fatal("runner.phase19-plan-unit.fixture.decode")
	}
	return result
}

func cicdPlanTestRaw(t *testing.T, value any) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal("runner.phase19-plan-unit.fixture.encode")
	}
	return data
}

func cicdPlanTestFixture(t *testing.T) (map[string]any, map[string]any, map[string]any, CICDPlanContext) {
	t.Helper()
	project := cicdPlanTestObject(t, `{"id":"p","name":"Project","revision":1,"enabled":true,"source_owner":"owner","source_repo":"repo","definition_id":"d","definition_revision":1,"quota":{"parallel":1,"queue":64,"artifact_mb":1024,"cache_mb":256,"retention_days":7},"toolchain_digest":"`+strings.Repeat("a", 64)+`"}`)
	definition := cicdPlanTestObject(t, `{"project_id":"p","id":"d","revision":1,"inputs":{},"manual":true,"watches":[],"schedules":[],"templates":[],"jobs":[{"id":"build","needs":[],"condition":{"branches":[],"tags":[],"paths":[],"input_equals":{}},"task":{"kind":"builder","source":".","output":"site"},"matrix":{"axes":{},"exclude":[],"include":[]},"resources":{"cpu_percent":100,"memory_mb":256,"disk_mb":128,"tasks":64,"timeout_seconds":900,"output_mb":16},"retry":{"max_attempts":1,"backoff_seconds":1},"cache":{"enabled":false,"paths":[],"ttl_seconds":604800},"env":{},"secrets":{},"optional":false}],"fail_fast":false}`)
	source := cicdPlanTestObject(t, `{"owner":"owner","repo":"repo","ref_kind":"branch","ref":"main","sha":"`+strings.Repeat("b", 40)+`","tree_digest":"4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945","files":[]}`)
	worker := json.RawMessage(`{"id":"w","revision":1,"enabled":true,"draining":false,"host":"worker.invalid","port":22,"user":"worker","allowed_projects":["p"],"capabilities":["builder","go","deno"],"slots":1,"last_seen":null,"leases":[],"toolchain_digest":"` + strings.Repeat("a", 64) + `","total_memory_mb":896}`)
	return project, definition, source, CICDPlanContext{Workers: []json.RawMessage{worker}}
}

func cicdPlanTestJob(t *testing.T, id string) map[string]any {
	t.Helper()
	_, d, _, _ := cicdPlanTestFixture(t)
	job := d["jobs"].([]any)[0].(map[string]any)
	job["id"] = id
	return job
}

func cicdPlanTestDraft(t *testing.T, p, d map[string]any, values string, templates ...json.RawMessage) *CICDPlanDraft {
	t.Helper()
	draft, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, d), json.RawMessage(values), templates)
	if err != nil || draft == nil {
		t.Fatal("runner.phase19-plan-unit.prepare.accept", err)
	}
	return draft
}

func TestCICDPlan(t *testing.T) {
	p, d, s, context := cicdPlanTestFixture(t)
	projectRaw, definitionRaw, sourceRaw := cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, d), cicdPlanTestRaw(t, s)
	draft, err := PrepareCICDPlan(projectRaw, definitionRaw, json.RawMessage(`{}`), nil)
	if err != nil {
		t.Fatal("runner.phase19-plan-unit.normal.prepare", err)
	}
	got, err := draft.Resolve(sourceRaw, nil, context)
	if err != nil {
		t.Fatal("runner.phase19-plan-unit.normal.resolve", err)
	}
	expectedJob := cicdPlanTestObject(t, `{"id":"build","base_job_id":"build","selected":true,"matrix_values":{},"needs":[],"task":{"kind":"builder","source":".","output":"site"},"resources":{"cpu_percent":100,"memory_mb":256,"disk_mb":128,"tasks":64,"timeout_seconds":900,"output_mb":16},"retry":{"max_attempts":1,"backoff_seconds":1},"cache":{"enabled":false,"paths":[],"ttl_seconds":604800},"env":{},"secrets":{},"optional":false}`)
	expectedData := map[string]any{"project_id": "p", "definition_id": "d", "definition_revision": json.Number("1"), "source": s, "values": map[string]any{}, "jobs": []any{expectedJob}, "fail_fast": false, "toolchain_digest": strings.Repeat("a", 64)}
	data := cicdPlanTestRaw(t, expectedData)
	digest := sha256.Sum256(data)
	want := cicdPlanTestRaw(t, map[string]any{"digest": hex.EncodeToString(digest[:]), "data": expectedData})
	if !bytes.Equal(got, want) || statefile.ValidateCICDValue("Plan", got) != nil {
		t.Fatal("runner.phase19-plan-unit.normal.plan: bytes or digest differ")
	}
	if !bytes.Equal(projectRaw, cicdPlanTestRaw(t, p)) || !bytes.Equal(definitionRaw, cicdPlanTestRaw(t, d)) || !bytes.Equal(sourceRaw, cicdPlanTestRaw(t, s)) {
		t.Fatal("runner.phase19-plan-unit.normal.immutable")
	}
	for _, data := range [][]byte{projectRaw, definitionRaw} {
		for i := range data {
			data[i] = 'x'
		}
	}
	got[0] = 'x'
	again, err := draft.Resolve(sourceRaw, nil, context)
	if err != nil || !bytes.Equal(again, want) {
		t.Fatal("runner.phase19-plan-unit.normal.alias")
	}
	for _, empty := range []*CICDPlanDraft{nil, {}} {
		if got, err := empty.Resolve(sourceRaw, nil, context); got != nil || err != ErrCICDSchema {
			t.Fatal("runner.phase19-plan-unit.zero.reject")
		}
	}
}

func TestCICDPlanInputs(t *testing.T) {
	for _, test := range []struct {
		id, declaration, supplied, want string
		bad                             bool
	}{
		{"required-zero", `{"kind":"integer","required":true,"default":null,"choices":[],"min":0,"max":4}`, `{"x":0}`, `{"x":0}`, false},
		{"required-false", `{"kind":"boolean","required":true,"default":null,"choices":[],"min":0,"max":1}`, `{"x":false}`, `{"x":false}`, false},
		{"required-empty", `{"kind":"string","required":true,"default":null,"choices":[],"min":0,"max":4}`, `{"x":""}`, `{"x":""}`, false},
		{"utf8-limit", `{"kind":"string","required":true,"default":null,"choices":[],"min":2,"max":2}`, `{"x":"é"}`, `{"x":"é"}`, false},
		{"utf8-excess", `{"kind":"string","required":true,"default":null,"choices":[],"min":0,"max":1}`, `{"x":"é"}`, ``, true},
		{"optional-omit", `{"kind":"integer","required":false,"default":null,"choices":[],"min":0,"max":4}`, `{}`, `{}`, false},
		{"default-zero", `{"kind":"integer","required":false,"default":0,"choices":[0,2],"min":0,"max":4}`, `{}`, `{"x":0}`, false},
		{"explicit-wins", `{"kind":"integer","required":false,"default":0,"choices":[0,2],"min":0,"max":4}`, `{"x":2}`, `{"x":2}`, false},
		{"required-missing", `{"kind":"integer","required":true,"default":null,"choices":[],"min":0,"max":4}`, `{}`, ``, true},
		{"required-default", `{"kind":"integer","required":true,"default":0,"choices":[],"min":0,"max":4}`, `{"x":0}`, ``, true},
		{"explicit-null", `{"kind":"integer","required":false,"default":0,"choices":[],"min":0,"max":4}`, `{"x":null}`, ``, true},
		{"numeric-string", `{"kind":"integer","required":true,"default":null,"choices":[],"min":0,"max":4}`, `{"x":"0"}`, ``, true},
		{"boolean-number", `{"kind":"boolean","required":true,"default":null,"choices":[],"min":0,"max":1}`, `{"x":0}`, ``, true},
		{"negative-zero", `{"kind":"integer","required":true,"default":null,"choices":[],"min":-4,"max":4}`, `{"x":-0}`, ``, true},
		{"exponent", `{"kind":"integer","required":true,"default":null,"choices":[],"min":0,"max":4}`, `{"x":1e0}`, ``, true},
		{"fraction", `{"kind":"integer","required":true,"default":null,"choices":[],"min":0,"max":4}`, `{"x":1.0}`, ``, true},
		{"below", `{"kind":"integer","required":true,"default":null,"choices":[],"min":0,"max":4}`, `{"x":-1}`, ``, true},
		{"above", `{"kind":"integer","required":true,"default":null,"choices":[],"min":0,"max":4}`, `{"x":5}`, ``, true},
		{"reverse-range", `{"kind":"integer","required":true,"default":null,"choices":[],"min":4,"max":0}`, `{"x":2}`, ``, true},
		{"boolean-range", `{"kind":"boolean","required":true,"default":null,"choices":[],"min":1,"max":1}`, `{"x":true}`, ``, true},
		{"string-range", `{"kind":"string","required":true,"default":null,"choices":[],"min":-1,"max":4097}`, `{"x":""}`, ``, true},
		{"duplicate-choice", `{"kind":"integer","required":true,"default":null,"choices":[1,1],"min":0,"max":4}`, `{"x":1}`, ``, true},
		{"choice-type", `{"kind":"integer","required":true,"default":null,"choices":[false],"min":0,"max":4}`, `{"x":0}`, ``, true},
		{"choice-range", `{"kind":"integer","required":true,"default":null,"choices":[5],"min":0,"max":4}`, `{"x":0}`, ``, true},
		{"choice-missing", `{"kind":"integer","required":true,"default":null,"choices":[1],"min":0,"max":4}`, `{"x":0}`, ``, true},
		{"default-choice", `{"kind":"integer","required":false,"default":0,"choices":[1],"min":0,"max":4}`, `{}`, ``, true},
		{"unknown-input", `{"kind":"integer","required":false,"default":null,"choices":[],"min":0,"max":4}`, `{"other":0}`, ``, true},
		{"duplicate-input", `{"kind":"integer","required":false,"default":null,"choices":[],"min":0,"max":4}`, `{"x":0,"x":1}`, ``, true},
	} {
		t.Run(test.id, func(t *testing.T) {
			p, d, s, context := cicdPlanTestFixture(t)
			d["inputs"] = map[string]any{"x": cicdPlanTestObject(t, test.declaration)}
			draft, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, d), json.RawMessage(test.supplied), nil)
			if test.bad {
				if err != ErrCICDSchema || draft != nil {
					t.Fatal("runner.phase19-plan-unit.inputs.reject")
				}
				return
			}
			if err != nil {
				t.Fatal("runner.phase19-plan-unit.inputs.prepare", err)
			}
			got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, context)
			if err != nil || got == nil {
				t.Fatal("runner.phase19-plan-unit.inputs.resolve", err)
			}
			plan := cicdPlanTestObject(t, string(got))
			if !reflect.DeepEqual(plan["data"].(map[string]any)["values"], cicdPlanTestObject(t, test.want)) {
				t.Fatal("runner.phase19-plan-unit.inputs.value")
			}
		})
	}
}

func TestCICDPlanMatrix(t *testing.T) {
	for _, test := range []struct {
		id, matrix, want string
		bad              bool
	}{
		{"empty", `{"axes":{},"exclude":[],"include":[]}`, `[{}]`, false},
		{"axis-order", `{"axes":{"z":[true,false],"a":[2,1]},"exclude":[],"include":[]}`, `[{"a":2,"z":true},{"a":2,"z":false},{"a":1,"z":true},{"a":1,"z":false}]`, false},
		{"exclude-reinclude", `{"axes":{"a":[1,2],"b":["x","y"]},"exclude":[{"a":1}],"include":[{"a":1,"b":"x"},{"a":3,"b":"z"}]}`, `[{"a":2,"b":"x"},{"a":2,"b":"y"},{"a":1,"b":"x"},{"a":3,"b":"z"}]`, false},
		{"empty-exclude", `{"axes":{"a":[1]},"exclude":[{}],"include":[]}`, ``, true},
		{"unknown-exclude", `{"axes":{"a":[1]},"exclude":[{"b":1}],"include":[]}`, ``, true},
		{"type-exclude", `{"axes":{"a":[1]},"exclude":[{"a":"1"}],"include":[]}`, ``, true},
		{"partial-include", `{"axes":{"a":[1],"b":[2]},"exclude":[],"include":[{"a":3}]}`, ``, true},
		{"type-include", `{"axes":{"a":[1]},"exclude":[],"include":[{"a":true}]}`, ``, true},
		{"duplicate-include", `{"axes":{"a":[1]},"exclude":[],"include":[{"a":2},{"a":2}]}`, ``, true},
		{"existing-include", `{"axes":{"a":[1]},"exclude":[],"include":[{"a":1}]}`, ``, true},
		{"empty-axis-include", `{"axes":{},"exclude":[],"include":[{}]}`, ``, true},
		{"empty-result", `{"axes":{"a":[1]},"exclude":[{"a":1}],"include":[]}`, ``, true},
		{"mixed-axis", `{"axes":{"a":[1,true]},"exclude":[],"include":[]}`, ``, true},
		{"duplicate-axis", `{"axes":{"a":[1,1]},"exclude":[],"include":[]}`, ``, true},
	} {
		t.Run(test.id, func(t *testing.T) {
			p, d, _, _ := cicdPlanTestFixture(t)
			matrix := cicdPlanTestObject(t, test.matrix)
			d["jobs"].([]any)[0].(map[string]any)["matrix"] = matrix
			draft, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, d), json.RawMessage(`{}`), nil)
			if test.bad {
				if err != ErrCICDMatrix || draft != nil {
					t.Fatal("runner.phase19-plan-unit.matrix.reject", err)
				}
				return
			}
			if err != nil {
				t.Fatal("runner.phase19-plan-unit.matrix.prepare", err)
			}
			got := make([]any, len(draft.jobs))
			for i, job := range draft.jobs {
				got[i] = job.job["matrix_values"]
				id := "build"
				if len(matrix["axes"].(map[string]any)) != 0 {
					id += fmt.Sprintf("_m%d", i+1)
				}
				if job.job["id"] != id {
					t.Fatal("runner.phase19-plan-unit.matrix.id")
				}
			}
			if string(cicdPlanTestRaw(t, got)) != test.want {
				t.Fatal("runner.phase19-plan-unit.matrix.order")
			}
		})
	}
	for _, count := range []int{256, 257} {
		p, d, _, _ := cicdPlanTestFixture(t)
		axis := make([]any, 16)
		for i := range axis {
			axis[i] = i
		}
		matrix := map[string]any{"axes": map[string]any{"a": axis, "b": axis}, "exclude": []any{}, "include": []any{}}
		if count == 257 {
			matrix["include"] = []any{map[string]any{"a": 17, "b": 17}}
		}
		d["jobs"].([]any)[0].(map[string]any)["matrix"] = matrix
		tuples, matrixErr := cicdMatrix(cicdPlanTestObject(t, string(cicdPlanTestRaw(t, matrix))))
		if count == 256 && (matrixErr != nil || len(tuples) != 256) || count == 257 && (matrixErr != ErrCICDMatrix || tuples != nil) {
			t.Fatal("runner.phase19-plan-unit.matrix.limit")
		}
		draft, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, d), json.RawMessage(`{}`), nil)
		if count == 256 && (err != nil || len(draft.jobs) != count) || count == 257 && (err != ErrCICDMatrix || draft != nil) {
			t.Fatal("runner.phase19-plan-unit.matrix.limit", err)
		}
	}
}

func TestCICDPlanDAG(t *testing.T) {
	p, d, _, _ := cicdPlanTestFixture(t)
	z, x, y := cicdPlanTestJob(t, "z"), cicdPlanTestJob(t, "x"), cicdPlanTestJob(t, "y")
	z["needs"] = []any{"y", "x"}
	d["jobs"] = []any{z, x, y}
	draft := cicdPlanTestDraft(t, p, d, `{}`)
	ids := make([]string, 0)
	for _, item := range draft.jobs {
		ids = append(ids, item.job["id"].(string))
	}
	if !reflect.DeepEqual(ids, []string{"x", "y", "z"}) || !reflect.DeepEqual(draft.jobs[2].job["needs"], []any{"x", "y"}) {
		t.Fatal("runner.phase19-plan-unit.dag.order")
	}
	for _, test := range []struct {
		id    string
		alter func(map[string]any)
		want  error
	}{
		{"cycle", func(d map[string]any) { d["jobs"].([]any)[1].(map[string]any)["needs"] = []any{"z"} }, ErrCICDDAG},
		{"self", func(d map[string]any) { d["jobs"].([]any)[0].(map[string]any)["needs"] = []any{"z"} }, ErrCICDDAG},
		{"duplicate", func(d map[string]any) { d["jobs"].([]any)[0].(map[string]any)["needs"] = []any{"x", "x"} }, ErrCICDDAG},
		{"unknown", func(d map[string]any) { d["jobs"].([]any)[0].(map[string]any)["needs"] = []any{"missing"} }, ErrCICDDAG},
		{"same-base", func(d map[string]any) { d["jobs"].([]any)[1].(map[string]any)["id"] = "z" }, ErrCICDDAG},
		{"generated-collision", func(d map[string]any) {
			d["jobs"].([]any)[1].(map[string]any)["id"] = "y_m1"
			d["jobs"].([]any)[0].(map[string]any)["needs"] = []any{}
			d["jobs"].([]any)[2].(map[string]any)["matrix"] = cicdPlanTestObject(t, `{"axes":{"a":[1]},"exclude":[],"include":[]}`)
		}, ErrCICDDAG},
	} {
		t.Run(test.id, func(t *testing.T) {
			copy := cicdPlanTestObject(t, string(cicdPlanTestRaw(t, d)))
			test.alter(copy)
			got, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, copy), json.RawMessage(`{}`), nil)
			if got != nil || err != test.want {
				t.Fatal("runner.phase19-plan-unit.dag.reject", err)
			}
		})
	}
}

func TestCICDPlanGlob(t *testing.T) {
	for _, test := range []struct {
		pattern, target string
		want, invalid   bool
	}{
		{"**/a", "a", true, false}, {"**/a", "x/y/a", true, false}, {"a/**", "a", true, false},
		{"a/**/b", "a/b", true, false}, {"a/**/b", "a/x/y/b", true, false}, {"a/*/b", "a/b", false, false},
		{"a/*/b", "a/x/y/b", false, false}, {"a/*/b", "a/x/b", true, false}, {"?.js", "é.js", true, false},
		{"?.js", "ab.js", false, false}, {"A", "a", false, false}, {"**", "a/b", true, false},
		{"a/**b", "a/b", false, true}, {"[ab]", "a", false, true}, {"!a", "a", false, true},
		{"a\\*", "a", false, true}, {"{a,b}", "a", false, true}, {"a//b", "a/b", false, true},
		{strings.Repeat("a", 241), "a", false, true}, {strings.Repeat("a/", 32) + "a", "a", false, true},
		{strings.Repeat("?", 240), strings.Repeat("a", 1000), false, true},
	} {
		got, err := cicdMatchGlob(test.pattern, test.target)
		if test.invalid && err != ErrCICDSchema || !test.invalid && (err != nil || got != test.want) {
			t.Fatal("runner.phase19-plan-unit.glob.match")
		}
	}
}

func TestCICDPlanTemplateAndSubstitution(t *testing.T) {
	p, d, s, context := cicdPlanTestFixture(t)
	templateJob := cicdPlanTestJob(t, "test")
	templateJob["task"] = map[string]any{"kind": "go-test", "packages": []any{"${input.package}"}, "race": "${matrix.race}"}
	templateJob["env"] = map[string]any{"SCOPE": "${input.package}"}
	templateJob["matrix"] = cicdPlanTestObject(t, `{"axes":{"race":[false]},"exclude":[],"include":[]}`)
	template := map[string]any{"project_id": "p", "id": "t", "revision": 7, "inputs": cicdPlanTestObject(t, `{"package":{"kind":"string","required":true,"default":null,"choices":[],"min":1,"max":64}}`), "jobs": []any{templateJob}}
	d["templates"] = []any{map[string]any{"id": "t", "revision": 7, "prefix": "pre", "values": map[string]any{"package": "."}}}
	d["jobs"].([]any)[0].(map[string]any)["needs"] = []any{"pre_test"}
	draft := cicdPlanTestDraft(t, p, d, `{}`, cicdPlanTestRaw(t, template))
	got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, context)
	if err != nil {
		t.Fatal("runner.phase19-plan-unit.template.resolve", err)
	}
	jobs := cicdPlanTestObject(t, string(got))["data"].(map[string]any)["jobs"].([]any)
	first := jobs[0].(map[string]any)
	if first["id"] != "pre_test_m1" || first["base_job_id"] != "pre_test" || first["task"].(map[string]any)["race"] != false || first["env"].(map[string]any)["SCOPE"] != "." || !reflect.DeepEqual(jobs[1].(map[string]any)["needs"], []any{"pre_test_m1"}) {
		t.Fatal("runner.phase19-plan-unit.template.expansion")
	}
	for _, test := range []struct {
		id    string
		alter func(map[string]any, map[string]any)
		want  error
	}{
		{"wrong-revision", func(d, _ map[string]any) { d["templates"].([]any)[0].(map[string]any)["revision"] = 8 }, ErrCICDSchema},
		{"cross-input", func(_, template map[string]any) {
			template["jobs"].([]any)[0].(map[string]any)["env"] = map[string]any{"SCOPE": "${input.parent}"}
		}, ErrCICDSchema},
		{"partial-reference", func(_, template map[string]any) {
			template["jobs"].([]any)[0].(map[string]any)["env"] = map[string]any{"SCOPE": "prefix-${input.package}"}
		}, ErrCICDSchema},
		{"external-need", func(_, template map[string]any) {
			template["jobs"].([]any)[0].(map[string]any)["needs"] = []any{"build"}
		}, ErrCICDDAG},
		{"long-prefix", func(d, _ map[string]any) {
			d["templates"].([]any)[0].(map[string]any)["prefix"] = strings.Repeat("a", 44)
		}, ErrCICDSchema},
		{"parent-value", func(d, _ map[string]any) {
			d["templates"].([]any)[0].(map[string]any)["values"] = map[string]any{"package": "${input.package}"}
		}, ErrCICDSchema},
		{"kind-change", func(_, template map[string]any) {
			template["jobs"].([]any)[0].(map[string]any)["task"].(map[string]any)["race"] = true
		}, ErrCICDSchema},
	} {
		t.Run(test.id, func(t *testing.T) {
			definition := cicdPlanTestObject(t, string(cicdPlanTestRaw(t, d)))
			copy := cicdPlanTestObject(t, string(cicdPlanTestRaw(t, template)))
			test.alter(definition, copy)
			got, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, definition), json.RawMessage(`{}`), []json.RawMessage{cicdPlanTestRaw(t, copy)})
			if got != nil || err != test.want {
				t.Fatal("runner.phase19-plan-unit.template.reject", err)
			}
		})
	}
	for _, value := range []string{"${matrix.race}", "literal"} {
		got, err := cicdSubstitute("${input.x}", map[string]any{"x": value}, map[string]any{"race": false})
		if err != nil || got != value {
			t.Fatal("runner.phase19-plan-unit.substitute.no-rescan")
		}
	}
}

func TestCICDPlanSelection(t *testing.T) {
	for _, test := range []struct {
		id                    string
		branches, tags, paths []any
		kind, ref             string
		changed               []string
		want                  bool
	}{
		{"unconditional", []any{}, []any{}, []any{}, "branch", "main", nil, true},
		{"tag-only", []any{}, []any{"v*"}, []any{}, "branch", "main", nil, false},
		{"branch-only", []any{"*"}, []any{}, []any{}, "tag", "v1", nil, false},
		{"tag-match", []any{}, []any{"v?"}, []any{}, "tag", "v1", nil, true},
		{"branch-or", []any{"dev", "main"}, []any{}, []any{}, "branch", "main", nil, true},
		{"empty-source", []any{}, []any{}, []any{"**/*.go"}, "branch", "main", nil, false},
		{"removed-path", []any{}, []any{}, []any{"**/*.go"}, "branch", "main", []string{"old/deleted.go"}, true},
		{"nonmatch-path", []any{}, []any{}, []any{"**/*.go"}, "branch", "main", []string{"README.md"}, false},
		{"ref-and-path", []any{"dev"}, []any{}, []any{"**/*.go"}, "branch", "main", []string{"main.go"}, false},
	} {
		t.Run(test.id, func(t *testing.T) {
			p, d, s, context := cicdPlanTestFixture(t)
			parent := d["jobs"].([]any)[0].(map[string]any)
			parent["condition"] = map[string]any{"branches": test.branches, "tags": test.tags, "paths": test.paths, "input_equals": map[string]any{}}
			parent["optional"] = true
			child := cicdPlanTestJob(t, "child")
			child["needs"] = []any{"build"}
			d["jobs"] = []any{parent, child}
			s["ref_kind"], s["ref"] = test.kind, test.ref
			if !test.want {
				context.Workers = nil
			}
			draft := cicdPlanTestDraft(t, p, d, `{}`)
			got, err := draft.Resolve(cicdPlanTestRaw(t, s), test.changed, context)
			if err != nil {
				t.Fatal("runner.phase19-plan-unit.selection.resolve", err)
			}
			jobs := cicdPlanTestObject(t, string(got))["data"].(map[string]any)["jobs"].([]any)
			if len(jobs) != 2 || jobs[0].(map[string]any)["selected"] != test.want || jobs[1].(map[string]any)["selected"] != test.want {
				t.Fatal("runner.phase19-plan-unit.selection.cascade")
			}
		})
	}
	for _, match := range []bool{false, true} {
		p, d, s, context := cicdPlanTestFixture(t)
		d["inputs"] = cicdPlanTestObject(t, `{"x":{"kind":"boolean","required":false,"default":false,"choices":[],"min":0,"max":1}}`)
		d["jobs"].([]any)[0].(map[string]any)["condition"].(map[string]any)["input_equals"] = map[string]any{"x": match}
		draft := cicdPlanTestDraft(t, p, d, `{}`)
		got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, context)
		if err != nil {
			t.Fatal("runner.phase19-plan-unit.input-equals.resolve", err)
		}
		job := cicdPlanTestObject(t, string(got))["data"].(map[string]any)["jobs"].([]any)[0].(map[string]any)
		if job["selected"] != !match {
			t.Fatal("runner.phase19-plan-unit.input-equals.typed")
		}
	}
}

func TestCICDPlanPrerequisites(t *testing.T) {
	for _, test := range []struct {
		id    string
		alter func(map[string]any, map[string]any, *CICDPlanContext)
		want  error
	}{
		{"disabled", func(p, _ map[string]any, _ *CICDPlanContext) { p["enabled"] = false }, ErrCICDForbidden},
		{"no-toolchain", func(p, _ map[string]any, _ *CICDPlanContext) { p["toolchain_digest"] = nil }, ErrCICDCapability},
		{"no-worker", func(_, _ map[string]any, c *CICDPlanContext) { c.Workers = nil }, ErrCICDCapability},
		{"wrong-toolchain", func(p, _ map[string]any, _ *CICDPlanContext) { p["toolchain_digest"] = strings.Repeat("b", 64) }, ErrCICDCapability},
		{"memory-incapable", func(_, d map[string]any, _ *CICDPlanContext) {
			d["jobs"].([]any)[0].(map[string]any)["resources"].(map[string]any)["memory_mb"] = 385
		}, ErrCICDCapability},
		{"secret-absent", func(_, d map[string]any, _ *CICDPlanContext) {
			d["jobs"].([]any)[0].(map[string]any)["secrets"] = map[string]any{"TOKEN": map[string]any{"id": "secret", "version": 1}}
		}, ErrCICDSecret},
		{"cache-quota", func(p, d map[string]any, _ *CICDPlanContext) {
			p["quota"].(map[string]any)["cache_mb"] = 0
			d["jobs"].([]any)[0].(map[string]any)["cache"] = map[string]any{"enabled": true, "paths": []any{"site"}, "ttl_seconds": 60}
		}, ErrCICDQuota},
		{"duplicate-worker", func(_, _ map[string]any, c *CICDPlanContext) { c.Workers = append(c.Workers, c.Workers[0]) }, ErrCICDSchema},
		{"wrong-project", func(_, _ map[string]any, c *CICDPlanContext) {
			c.Workers[0] = bytes.ReplaceAll(c.Workers[0], []byte(`"p"`), []byte(`"other"`))
		}, ErrCICDCapability},
		{"disabled-worker", func(_, _ map[string]any, c *CICDPlanContext) {
			c.Workers[0] = bytes.ReplaceAll(c.Workers[0], []byte(`"enabled":true`), []byte(`"enabled":false`))
		}, ErrCICDCapability},
		{"draining-not-incapable", func(_, _ map[string]any, c *CICDPlanContext) {
			c.Workers[0] = bytes.ReplaceAll(c.Workers[0], []byte(`"draining":false`), []byte(`"draining":true`))
		}, nil},
	} {
		t.Run(test.id, func(t *testing.T) {
			p, d, s, context := cicdPlanTestFixture(t)
			test.alter(p, d, &context)
			draft := cicdPlanTestDraft(t, p, d, `{}`)
			got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, context)
			if err != test.want || test.want != nil && got != nil || test.want == nil && len(got) == 0 {
				t.Fatal("runner.phase19-plan-unit.prerequisite.result", err)
			}
		})
	}
	for _, test := range []struct {
		id    string
		alter func(map[string]any)
	}{
		{"cache-empty", func(j map[string]any) { j["cache"].(map[string]any)["enabled"] = true }},
		{"cache-secret", func(j map[string]any) {
			j["cache"] = map[string]any{"enabled": true, "paths": []any{"out"}, "ttl_seconds": 60}
			j["secrets"] = map[string]any{"TOKEN": map[string]any{"id": "secret", "version": 1}}
		}},
		{"cache-overlap", func(j map[string]any) { j["cache"].(map[string]any)["paths"] = []any{"out", "out/x"} }},
		{"cache-duplicate", func(j map[string]any) { j["cache"].(map[string]any)["paths"] = []any{"out", "out"} }},
		{"cache-glob", func(j map[string]any) { j["cache"].(map[string]any)["paths"] = []any{"out/*"} }},
		{"missing-equality", func(j map[string]any) {
			j["condition"].(map[string]any)["input_equals"] = map[string]any{"missing": false}
		}},
		{"reserved-env", func(j map[string]any) { j["env"] = map[string]any{"GOPROXY": "off"} }},
		{"bad-substituted-path", func(j map[string]any) { j["task"].(map[string]any)["source"] = "${input.missing}" }},
	} {
		t.Run(test.id, func(t *testing.T) {
			p, d, _, _ := cicdPlanTestFixture(t)
			test.alter(d["jobs"].([]any)[0].(map[string]any))
			got, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, d), json.RawMessage(`{}`), nil)
			if err != ErrCICDSchema || got != nil {
				t.Fatal("runner.phase19-plan-unit.preflight.reject", err)
			}
		})
	}
}

func TestCICDPlanDeploy(t *testing.T) {
	p, d, s, context := cicdPlanTestFixture(t)
	deploy := cicdPlanTestJob(t, "publish")
	deploy["needs"] = []any{"build"}
	deploy["task"] = map[string]any{"kind": "deploy", "environment_id": "prod", "artifact_from": "build"}
	d["jobs"] = append(d["jobs"].([]any), deploy)
	context.Environments = []json.RawMessage{json.RawMessage(`{"project_id":"p","id":"prod","revision":1,"production":false,"approvals":0,"allow_self_approval":false,"targets":["site"],"strategy":"rolling","stages":[{"target_ids":["site"],"percent":100,"observe_seconds":10,"max_error_percent":0,"min_requests":0}],"freeze":[],"health_timeout_seconds":1,"retained_successes":2,"health_path":"index.html"}`)}
	draft := cicdPlanTestDraft(t, p, d, `{}`)
	if got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, context); got == nil || err != nil {
		t.Fatal("runner.phase19-plan-unit.deploy.accept", err)
	}
	without := context
	without.Environments = nil
	if got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, without); got != nil || err != ErrCICDNotFound {
		t.Fatal("runner.phase19-plan-unit.deploy.environment")
	}
	for _, kind := range []string{"matrix", "go-build", "not-needed"} {
		copy := cicdPlanTestObject(t, string(cicdPlanTestRaw(t, d)))
		builder := copy["jobs"].([]any)[0].(map[string]any)
		switch kind {
		case "matrix":
			builder["matrix"] = cicdPlanTestObject(t, `{"axes":{"a":[1,2]},"exclude":[],"include":[]}`)
		case "go-build":
			builder["task"] = map[string]any{"kind": "go-build", "package": ".", "output": "bin"}
		case "not-needed":
			copy["jobs"].([]any)[1].(map[string]any)["needs"] = []any{}
		}
		if got, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, copy), json.RawMessage(`{}`), nil); got != nil || err != ErrCICDSchema {
			t.Fatal("runner.phase19-plan-unit.deploy.reference", err)
		}
	}
	verify := cicdPlanTestJob(t, "verify")
	verify["task"] = map[string]any{"kind": "artifact-verify", "artifact": strings.Repeat("c", 64)}
	d["jobs"] = []any{verify}
	context.Workers = nil
	context.Artifacts = []json.RawMessage{json.RawMessage(`{"project_id":"p","digest":"` + strings.Repeat("c", 64) + `","manifest_digest":"` + strings.Repeat("d", 64) + `","size_bytes":1,"created_at":"2026-01-01T00:00:00Z","source_sha":"` + strings.Repeat("b", 40) + `","run_id":"r","attempt":1,"job_id":"build","pins":[],"deleting":false}`)}
	draft = cicdPlanTestDraft(t, p, d, `{}`)
	if got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, context); got == nil || err != nil {
		t.Fatal("runner.phase19-plan-unit.verify.no-worker", err)
	}
	context.Artifacts[0] = bytes.ReplaceAll(context.Artifacts[0], []byte(`"deleting":false`), []byte(`"deleting":true`))
	if got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, context); got != nil || err != ErrCICDNotFound {
		t.Fatal("runner.phase19-plan-unit.verify.deleting")
	}
}

func cicdPlanTestSourceFiles(t *testing.T, source map[string]any, files []any) {
	t.Helper()
	digest := sha256.Sum256(cicdPlanTestRaw(t, files))
	source["files"], source["tree_digest"] = files, hex.EncodeToString(digest[:])
}

func cicdPlanTestFile(path, digest string, size int) map[string]any {
	return map[string]any{"path": path, "blob_sha": strings.Repeat("b", 40), "sha256": strings.Repeat(digest, 64), "size": size, "executable": false}
}

func TestCICDPlanSource(t *testing.T) {
	p, d, s, context := cicdPlanTestFixture(t)
	draft := cicdPlanTestDraft(t, p, d, `{}`)
	before := cicdPlanTestObject(t, string(cicdPlanTestRaw(t, s)))
	cicdPlanTestSourceFiles(t, before, []any{cicdPlanTestFile("deleted.go", "a", 1), cicdPlanTestFile("edited.go", "a", 1), cicdPlanTestFile("unchanged.go", "a", 1)})
	cicdPlanTestSourceFiles(t, s, []any{cicdPlanTestFile("added.go", "a", 1), cicdPlanTestFile("edited.go", "b", 1), cicdPlanTestFile("unchanged.go", "a", 1)})
	if got, err := CICDChangedPaths(cicdPlanTestRaw(t, before), cicdPlanTestRaw(t, s)); err != nil || !reflect.DeepEqual(got, []string{"added.go", "deleted.go", "edited.go"}) {
		t.Fatal("runner.phase19-plan-unit.source.diff", err)
	}
	if got, err := CICDChangedPaths(nil, cicdPlanTestRaw(t, s)); err != nil || !reflect.DeepEqual(got, []string{"added.go", "edited.go", "unchanged.go"}) {
		t.Fatal("runner.phase19-plan-unit.source.first", err)
	}
	if got, err := CICDChangedPaths(cicdPlanTestRaw(t, s), cicdPlanTestRaw(t, s)); err != nil || got == nil || len(got) != 0 {
		t.Fatal("runner.phase19-plan-unit.source.unchanged", err)
	}
	before["ref"] = "different"
	if got, err := CICDChangedPaths(cicdPlanTestRaw(t, before), cicdPlanTestRaw(t, s)); got != nil || err != ErrCICDSchema {
		t.Fatal("runner.phase19-plan-unit.source.identity")
	}
	for _, kind := range []string{"wrong-owner", "digest", "order", "duplicate", "excess-file", "excess-total"} {
		copy := cicdPlanTestObject(t, string(cicdPlanTestRaw(t, s)))
		switch kind {
		case "wrong-owner":
			copy["owner"] = "other"
		case "digest":
			copy["tree_digest"] = strings.Repeat("0", 64)
		case "order":
			cicdPlanTestSourceFiles(t, copy, []any{cicdPlanTestFile("z", "a", 1), cicdPlanTestFile("a", "a", 1)})
		case "duplicate":
			cicdPlanTestSourceFiles(t, copy, []any{cicdPlanTestFile("a", "a", 1), cicdPlanTestFile("a", "a", 1)})
		case "excess-file":
			cicdPlanTestSourceFiles(t, copy, []any{cicdPlanTestFile("a", "a", 10485761)})
		case "excess-total":
			files := make([]any, 11)
			for i := range files {
				files[i] = cicdPlanTestFile(fmt.Sprintf("f%02d", i), "a", 10485760)
			}
			cicdPlanTestSourceFiles(t, copy, files)
		}
		if got, err := draft.Resolve(cicdPlanTestRaw(t, copy), nil, context); got != nil || err != ErrCICDSchema {
			t.Fatal("runner.phase19-plan-unit.source.reject", err)
		}
	}
	for _, paths := range [][]string{{"z", "a"}, {"a", "a"}, {"../a"}} {
		if got, err := draft.Resolve(cicdPlanTestRaw(t, s), paths, context); got != nil || err != ErrCICDSchema {
			t.Fatal("runner.phase19-plan-unit.source.paths")
		}
	}
	files := make([]any, 3000)
	for i := range files {
		files[i] = cicdPlanTestFile(fmt.Sprintf("f%04d", i), "a", 1)
	}
	cicdPlanTestSourceFiles(t, s, files)
	if got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, context); got != nil || err != ErrCICDBody {
		t.Fatal("runner.phase19-plan-unit.source.plan-size", err)
	}
}

func TestCICDPlanConcurrent(t *testing.T) {
	p, d, s, context := cicdPlanTestFixture(t)
	draft := cicdPlanTestDraft(t, p, d, `{}`)
	source := cicdPlanTestRaw(t, s)
	want, err := draft.Resolve(source, nil, context)
	if err != nil {
		t.Fatal("runner.phase19-plan-unit.concurrent.control", err)
	}
	start := make(chan struct{})
	results := make(chan bool, 8)
	var ready sync.WaitGroup
	ready.Add(8)
	for i := 0; i < 8; i++ {
		go func() {
			ready.Done()
			<-start
			for n := 0; n < 8; n++ {
				got, err := draft.Resolve(source, nil, context)
				if err != nil || !bytes.Equal(got, want) {
					results <- false
					return
				}
			}
			results <- true
		}()
	}
	ready.Wait()
	close(start)
	for i := 0; i < 8; i++ {
		if !<-results {
			t.Fatal("runner.phase19-plan-unit.concurrent.immutable")
		}
	}
}

func TestCICDPlanWire(t *testing.T) {
	p, d, s, context := cicdPlanTestFixture(t)
	for _, raw := range []json.RawMessage{nil, json.RawMessage(`null`), json.RawMessage(`[]`), json.RawMessage(`{"unknown":true}`), json.RawMessage(`{} {}`), json.RawMessage(`{"x":true,"x":false}`), {0xff}} {
		if got, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, d), raw, nil); got != nil || err != ErrCICDSchema {
			t.Fatal("runner.phase19-plan-unit.wire.reject")
		}
	}
	for _, field := range []string{"manual", "inputs", "jobs", "fail_fast"} {
		copy := cicdPlanTestObject(t, string(cicdPlanTestRaw(t, d)))
		delete(copy, field)
		if got, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, copy), json.RawMessage(`{}`), nil); got != nil || err != ErrCICDSchema {
			t.Fatal("runner.phase19-plan-unit.wire.required")
		}
	}
	d["manual"] = false
	draft := cicdPlanTestDraft(t, p, d, `{}`)
	if got, err := draft.Resolve(cicdPlanTestRaw(t, s), nil, context); got == nil || err != nil {
		t.Fatal("runner.phase19-plan-unit.wire.manual-false-preview", err)
	}
	for _, automatic := range []bool{false, true} {
		d["inputs"] = cicdPlanTestObject(t, `{"x":{"kind":"boolean","required":true,"default":null,"choices":[],"min":0,"max":1}}`)
		d["schedules"] = []any{map[string]any{"id": "daily", "ref_kind": "branch", "ref": "main", "anchor": "2026-01-01T00:00:00Z", "every_seconds": 60, "catch_up": "latest", "enabled": automatic}}
		got, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, d), json.RawMessage(`{"x":true}`), nil)
		if automatic && (got != nil || err != ErrCICDSchema) || !automatic && (got == nil || err != nil) {
			t.Fatal("runner.phase19-plan-unit.wire.automatic-input", err)
		}
	}
}

func TestCICDPlanDeferredSchema(t *testing.T) {
	job := cicdPlanTestJob(t, "test")
	job["task"] = map[string]any{"kind": "go-test", "packages": []any{"${input.package}"}, "race": "${matrix.race}"}
	job["cache"] = map[string]any{"enabled": false, "paths": []any{"${input." + strings.Repeat("a", 63) + "}"}, "ttl_seconds": 60}
	if statefile.ValidateCICDValue("Job", cicdPlanTestRaw(t, job)) != nil {
		t.Fatal("runner.phase19-plan-unit.deferred.job")
	}
	planned := cicdPlanTestObject(t, string(cicdPlanTestRaw(t, job)))
	delete(planned, "matrix")
	delete(planned, "condition")
	planned["base_job_id"], planned["selected"], planned["matrix_values"] = "test", true, map[string]any{}
	if statefile.ValidateCICDValue("PlannedJob", cicdPlanTestRaw(t, planned)) != statefile.ErrCICDSchema {
		t.Fatal("runner.phase19-plan-unit.deferred.planned")
	}
	for _, task := range []map[string]any{
		{"kind": "go-test", "packages": "${input.package}", "race": false},
		{"kind": "go-test", "packages": []any{"prefix-${input.package}"}, "race": false},
		{"kind": "go-test", "packages": []any{"${other.package}"}, "race": false},
		{"kind": "go-test", "packages": []any{"${input.Bad}"}, "race": false},
		{"kind": "${input.kind}", "packages": []any{"."}, "race": false},
		{"kind": "deploy", "environment_id": "${input.env}", "artifact_from": "build"},
		{"kind": "deploy", "environment_id": "prod", "artifact_from": "${input.builder}"},
	} {
		job["task"] = task
		if statefile.ValidateCICDValue("Job", cicdPlanTestRaw(t, job)) != statefile.ErrCICDSchema {
			t.Fatal("runner.phase19-plan-unit.deferred.reject")
		}
	}
}

func TestCICDPlanTaskSyntax(t *testing.T) {
	for _, test := range []struct {
		task  map[string]any
		valid bool
	}{
		{map[string]any{"kind": "go-test", "packages": []any{"./...", "sub/-race"}, "race": false}, true},
		{map[string]any{"kind": "go-test", "packages": []any{"-race"}, "race": false}, false},
		{map[string]any{"kind": "go-vet", "packages": []any{"-toolexec=tool"}}, false},
		{map[string]any{"kind": "go-build", "package": "-buildmode=plugin", "output": "bin"}, false},
		{map[string]any{"kind": "go-build", "package": ".", "output": "bin"}, true},
		{map[string]any{"kind": "deno-check", "files": []any{"sub/--config=local.js"}}, true},
		{map[string]any{"kind": "deno-check", "files": []any{"--config=local.js"}}, false},
		{map[string]any{"kind": "deno-check", "files": []any{"source.ts"}}, false},
	} {
		p, d, _, _ := cicdPlanTestFixture(t)
		d["jobs"].([]any)[0].(map[string]any)["task"] = test.task
		got, err := PrepareCICDPlan(cicdPlanTestRaw(t, p), cicdPlanTestRaw(t, d), json.RawMessage(`{}`), nil)
		if test.valid && (got == nil || err != nil) || !test.valid && (got != nil || err != ErrCICDSchema) {
			t.Fatal("runner.phase19-plan-unit.task-syntax.result", err)
		}
	}
}

func TestCICDPlanMutations(t *testing.T) {
	for _, mutation := range []struct{ id, path, old, new, reason string }{
		{"default", "runner/validate.go", "value = d[\"default\"]", "value = nil", "runner.phase19-plan-unit.inputs.value"},
		{"utf8-bytes", "runner/validate.go", "return ok && int64(len(v)) >= min && int64(len(v)) <= max", "return ok && int64(utf8.RuneCountInString(v)) >= min && int64(utf8.RuneCountInString(v)) <= max", "runner.phase19-plan-unit.inputs.prepare"},
		{"matrix-order", "runner/validate.go", "keys := cicdPlanKeys(axes)", "keys := cicdPlanKeys(axes); sort.Sort(sort.Reverse(sort.StringSlice(keys)))", "runner.phase19-plan-unit.matrix.order"},
		{"matrix-duplicate", "runner/validate.go", "if reflect.DeepEqual(previous, add) {", "if false && reflect.DeepEqual(previous, add) {", "runner.phase19-plan-unit.matrix.reject"},
		{"matrix-limit", "runner/validate.go", "if len(result) > 256 {", "if len(result) > 257 {", "runner.phase19-plan-unit.matrix.limit"},
		{"dag-cycle", "runner/execute.go", "if !found {\n\t\t\treturn nil, ErrCICDDAG", "if !found {\n\t\t\treturn draft, nil", "runner.phase19-plan-unit.dag.reject"},
		{"dag-unknown", "runner/execute.go", "if _, exists := baseIndex[base]; !exists {", "if _, exists := baseIndex[base]; false && !exists {", "runner.phase19-plan-unit.dag.reject"},
		{"selection", "runner/execute.go", "match = match && selected[need.(string)]", "match = match || selected[need.(string)]", "runner.phase19-plan-unit.selection.resolve"},
		{"project-enabled", "runner/execute.go", "if draft.project[\"enabled\"] != true {", "if false && draft.project[\"enabled\"] != true {", "runner.phase19-plan-unit.prerequisite.result"},
		{"plan-body", "runner/execute.go", "if len(canonical) > 524288 {", "if len(canonical) > 33554432 {", "runner.phase19-plan-unit.source.plan-size"},
		{"artifact-gc", "runner/execute.go", "&& object[\"deleting\"] != true", "&& true", "runner.phase19-plan-unit.verify.deleting"},
		{"secret-version", "runner/execute.go", "if !found {\n\t\t\treturn ErrCICDSecret", "if false && !found {\n\t\t\treturn ErrCICDSecret", "runner.phase19-plan-unit.prerequisite.result"},
		{"deferred-schema", "statefile/validate.go", "if expression == \"Job\" && (key == \"task\" || key == \"cache\") {", "if false && expression == \"Job\" && (key == \"task\" || key == \"cache\") {", "runner.phase19-plan-unit.prepare.accept"},
		{"argv", "runner/execute.go", "|| !cicdPlanTaskValid(job[\"task\"].(map[string]any))", "|| false && !cicdPlanTaskValid(job[\"task\"].(map[string]any))", "runner.phase19-plan-unit.task-syntax.result"},
		{"source-truncated", "runner/validate.go", "if truncated || len(entries) > 100000 {", "if false && truncated || len(entries) > 100000 {", "runner.phase19-source.tree.reject"},
		{"source-file-limit", "runner/validate.go", "if len(files) > 10000 || total > 104857600 {", "if len(files) > 10001 || total > 104857600 {", "runner.phase19-source.tree.file-limit"},
		{"source-total-limit", "runner/validate.go", "if len(files) > 10000 || total > 104857600 {", "if len(files) > 10000 || total > 104857601 {", "runner.phase19-source.tree.total-limit"},
		{"source-parent", "runner/validate.go", "if paths[name] != \"tree\" {", "if false && paths[name] != \"tree\" {", "runner.phase19-source.tree.reject"},
		{"source-sha1", "runner/execute.go", "if hex.EncodeToString(hash.Sum(nil)) != file[\"blob_sha\"] {", "if false && hex.EncodeToString(hash.Sum(nil)) != file[\"blob_sha\"] {", "runner.phase19-source.blob.reject"},
		{"source-blob-size", "runner/execute.go", "blob[\"size\"] != file[\"size\"]", "false && blob[\"size\"] != file[\"size\"]", "runner.phase19-source.blob.reject"},
		{"source-pad-bits", "runner/execute.go", "base64.StdEncoding.Strict().DecodeString(content)", "base64.StdEncoding.DecodeString(content)", "runner.phase19-source.blob.pad-bits"},
		{"source-body-limit", "runner/execute.go", "if int64(len(data)) > limit {", "if false && int64(len(data)) > limit {", "runner.phase19-source.http.result"},
		{"source-metadata-copy", "runner/execute.go", "return append(json.RawMessage(nil), snapshot.metadata...)", "return snapshot.metadata", "runner.phase19-source.metadata.alias"},
		{"source-file-copy", "runner/execute.go", "return append([]byte{}, data...), nil", "return data, nil", "runner.phase19-source.file.alias"},
		{"source-no-reuse", "runner/execute.go", "DisableKeepAlives: true", "DisableKeepAlives: false", "runner.phase19-source.http.transport-profile"},
		{"run-values-copy", "runner/execute.go", "values = bytes.Clone(values)", "values = values[:]", "runner.phase19-run.start.result"},
		{"run-replay", "runner/execute.go", "if err != nil || replay.Replayed {", "if err != nil {", "runner.phase19-run.start.result"},
	} {
		t.Run(mutation.id, func(t *testing.T) {
			root := t.TempDir()
			for _, owner := range []string{"runner", "statefile", "security"} {
				target := filepath.Join(root, "components", owner)
				if os.MkdirAll(target, 0700) != nil {
					t.Fatal("runner.phase19-plan-unit.mutation.directory")
				}
				for _, name := range []string{owner + ".go", "model.go", "validate.go", "execute.go", owner + "_test.go"} {
					data, err := os.ReadFile(filepath.Join("..", owner, name))
					if err != nil {
						t.Fatal("runner.phase19-plan-unit.mutation.source")
					}
					if owner+"/"+name == mutation.path {
						if strings.Count(string(data), mutation.old) != 1 {
							t.Fatal("runner.phase19-plan-unit.mutation.target")
						}
						data = []byte(strings.Replace(string(data), mutation.old, mutation.new, 1))
					}
					if os.WriteFile(filepath.Join(target, name), data, 0600) != nil {
						t.Fatal("runner.phase19-plan-unit.mutation.copy")
					}
				}
			}
			if os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/fqwink/build-scripts\n\ngo 1.22\n"), 0600) != nil {
				t.Fatal("runner.phase19-plan-unit.mutation.module")
			}
			args := []string{"test", "-c", "-o", "mutation.test"}
			if runtime.GOOS == "darwin" {
				args = append(args, "-ldflags=-linkmode=external")
			}
			compile := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), append(args, "./components/runner")...)
			compile.Dir = root
			compile.Env = append(os.Environ(), "GOENV=off", "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off", "GOWORK=off")
			if _, err := compile.CombinedOutput(); err != nil {
				t.Fatal("runner.phase19-plan-unit.mutation.compile-invalid")
			}
			command := exec.Command(filepath.Join(root, "mutation.test"), "-test.run=^TestCICD(RunStart|Plan(Inputs|Matrix|DAG|Selection|Prerequisites|Source|Deploy|TemplateAndSubstitution|TaskSyntax)|Source(Input|Tree|Blob|HTTP)?)$")
			command.Dir = root
			output, err := command.CombinedOutput()
			if err == nil || !bytes.Contains(output, []byte(mutation.reason)) {
				t.Fatal("runner.phase19-plan-unit.mutation.not-killed")
			}
		})
	}
}

func TestCICDResources(t *testing.T) {
	want := CICDResources{CPUPercent: 100, MemoryMB: 256, DiskMB: 128, Tasks: 64, TimeoutSeconds: 900, OutputMB: 16}
	if got := DefaultCICDResources(); got != want || ValidateCICDResources(got) != nil {
		t.Fatal("runner.phase19-memory.defaults.value: unexpected defaults")
	}
	fields := []struct {
		name     string
		min, max int64
		set      func(*CICDResources, int64)
	}{
		{"cpu", 10, 400, func(r *CICDResources, v int64) { r.CPUPercent = v }},
		{"memory", 256, 8192, func(r *CICDResources, v int64) { r.MemoryMB = v }},
		{"disk", 128, 2048, func(r *CICDResources, v int64) { r.DiskMB = v }},
		{"tasks", 8, 256, func(r *CICDResources, v int64) { r.Tasks = v }},
		{"timeout", 1, 3600, func(r *CICDResources, v int64) { r.TimeoutSeconds = v }},
		{"output", 1, 16, func(r *CICDResources, v int64) { r.OutputMB = v }},
	}
	for _, field := range fields {
		for _, value := range []int64{field.min - 1, field.min, field.max, field.max + 1, -1, 1<<63 - 1} {
			t.Run(fmt.Sprintf("%s-%d", field.name, value), func(t *testing.T) {
				r := want
				if field.name == "disk" {
					r.MemoryMB = 8192
				}
				field.set(&r, value)
				before := r
				err := ValidateCICDResources(r)
				valid := value >= field.min && value <= field.max
				if valid && err != nil || !valid && err != ErrCICDSchema || r != before {
					t.Fatal("runner.phase19-memory.resources.boundary: unexpected validation or mutation")
				}
			})
		}
	}
	for _, gap := range []int64{127, 128, 129} {
		r := want
		r.MemoryMB, r.DiskMB = 512, 512-gap
		if err := ValidateCICDResources(r); (gap >= 128 && err != nil) || (gap < 128 && err != ErrCICDSchema) {
			t.Fatal("runner.phase19-memory.resources.gap: unexpected compound boundary")
		}
	}
	if ValidateCICDResources(CICDResources{}) != ErrCICDSchema {
		t.Fatal("runner.phase19-memory.resources.missing: zero input accepted")
	}
}

func TestCICDMemoryAdmission(t *testing.T) {
	tests := []struct {
		id               string
		total, available int64
		reserved         []int64
		slots, request   int64
		want             error
	}{
		{"minimum", 896, 320, nil, 1, 256, nil},
		{"nominal-one-gib", 1024, 1024, nil, 4, 256, nil},
		{"kernel-reservation", 960, 600, nil, 1, 256, nil},
		{"host-below", 895, 895, nil, 1, 256, ErrCICDCapability},
		{"host-upper", 1048576, 1048576, nil, 4, 8192, nil},
		{"host-above", 1048577, 1048577, nil, 4, 256, ErrCICDCapability},
		{"host-overflow", 1<<63 - 1, 320, nil, 1, 256, ErrCICDCapability},
		{"available-negative", 896, -1, nil, 1, 256, ErrCICDCapability},
		{"available-excess", 896, 897, nil, 1, 256, ErrCICDCapability},
		{"available-zero", 896, 0, nil, 1, 256, ErrCICDBusy},
		{"guard-below", 896, 319, nil, 1, 256, ErrCICDBusy},
		{"budget-exact", 896, 448, nil, 1, 384, nil},
		{"budget-above", 896, 896, nil, 1, 385, ErrCICDCapability},
		{"old-default-rejected", 960, 960, nil, 1, 1024, ErrCICDCapability},
		{"low-host-single-slot", 1024, 1024, []int64{256}, 4, 256, ErrCICDBusy},
		{"slot-threshold-below", 1535, 1535, []int64{256}, 4, 256, ErrCICDBusy},
		{"slot-threshold-exact", 1536, 1536, []int64{256}, 4, 256, nil},
		{"configured-single", 8192, 8192, []int64{256}, 1, 256, ErrCICDBusy},
		{"configured-four", 8192, 8192, []int64{256, 256, 256}, 4, 256, nil},
		{"four-full", 8192, 8192, []int64{256, 256, 256, 256}, 4, 256, ErrCICDBusy},
		{"claim-count-invalid", 8192, 8192, []int64{256, 256, 256, 256, 256}, 4, 256, ErrCICDStateCorrupt},
		{"configured-count-invalid", 8192, 8192, []int64{256, 256}, 1, 256, ErrCICDStateCorrupt},
		{"reserved-negative", 8192, 8192, []int64{-1}, 4, 256, ErrCICDStateCorrupt},
		{"reserved-below", 8192, 8192, []int64{255}, 4, 256, ErrCICDStateCorrupt},
		{"reserved-above", 8192, 8192, []int64{8193}, 4, 256, ErrCICDStateCorrupt},
		{"reserved-overflow", 8192, 8192, []int64{1<<63 - 1}, 4, 256, ErrCICDStateCorrupt},
		{"reserve-total-exact", 2048, 1600, []int64{1024}, 4, 512, nil},
		{"reserve-total-above", 2048, 2048, []int64{1025}, 4, 512, ErrCICDBusy},
		{"pressure-with-claim", 2048, 1599, []int64{1024}, 4, 512, ErrCICDBusy},
		{"request-below", 896, 896, nil, 1, 255, ErrCICDSchema},
		{"request-above", 16384, 16384, nil, 1, 8193, ErrCICDSchema},
		{"slots-below", 896, 896, nil, 0, 256, ErrCICDSchema},
		{"slots-above", 896, 896, nil, 5, 256, ErrCICDSchema},
		{"schema-before-capability", 0, 0, nil, 0, 0, ErrCICDSchema},
		{"capability-before-state", 0, 0, []int64{-1}, 1, 256, ErrCICDCapability},
		{"state-before-budget", 896, 896, []int64{-1}, 1, 512, ErrCICDStateCorrupt},
	}
	for _, tc := range tests {
		t.Run(tc.id, func(t *testing.T) {
			memory := CICDHostMemory{TotalMB: tc.total, AvailableMB: tc.available}
			before := append([]int64(nil), tc.reserved...)
			if got := CheckCICDMemory(memory, tc.reserved, tc.slots, tc.request); got != tc.want || !reflect.DeepEqual(before, tc.reserved) {
				t.Fatal("runner.phase19-memory." + tc.id + ".admit: unexpected code or input mutation")
			}
		})
	}
}

func checkCICDMemorySpecification(data []byte) error {
	parts := bytes.SplitN(data, []byte("```json\n"), 2)
	if len(parts) != 2 {
		return errors.New("memory-spec: missing registry")
	}
	block := bytes.SplitN(parts[1], []byte("\n```"), 2)
	var registry struct {
		Types map[string]json.RawMessage `json:"types"`
	}
	if len(block) != 2 || json.Unmarshal(block[0], &registry) != nil {
		return errors.New("memory-spec: invalid registry")
	}
	var resources map[string]string
	var worker map[string]json.RawMessage
	want := map[string]string{"cpu_percent": "int(10,400)", "memory_mb": "int(256,8192)", "disk_mb": "int(128,2048)", "tasks": "int(8,256)", "timeout_seconds": "int(1,3600)", "output_mb": "int(1,16)"}
	if json.Unmarshal(registry.Types["Resources"], &resources) != nil || !reflect.DeepEqual(resources, want) ||
		json.Unmarshal(registry.Types["Worker"], &worker) != nil || string(worker["total_memory_mb"]) != `"int(896,1048576)"` {
		return errors.New("memory-spec: resource bounds")
	}
	if !bytes.Contains(data, []byte("Resources={cpu_percent:100,memory_mb:256,disk_mb:128,tasks:64,timeout_seconds:900,output_mb:16}")) ||
		!bytes.Contains(data, []byte("memory_mb >= disk_mb+128")) {
		return errors.New("memory-spec: default or compound bound")
	}
	return nil
}

func TestCICDMemorySpecification(t *testing.T) {
	data, err := os.ReadFile("../../docs/details/statefile.md")
	if err != nil || checkCICDMemorySpecification(data) != nil {
		t.Fatal("runner.phase19-memory.specification.accept: contract drift")
	}
	for _, old := range []string{`"int(256,8192)"`, `"int(896,1048576)"`, "memory_mb:256", "disk_mb:128", "memory_mb >= disk_mb+128"} {
		if !bytes.Contains(data, []byte(old)) {
			t.Fatal("runner.phase19-memory.specification.fixture: negative-control target absent")
		}
		changed := bytes.Replace(data, []byte(old), []byte("invalid-control"), 1)
		if checkCICDMemorySpecification(changed) == nil {
			t.Fatal("runner.phase19-memory.specification.reject: modified contract accepted")
		}
	}
}

func TestCICDMemoryConcurrent(t *testing.T) {
	memory := CICDHostMemory{TotalMB: 8192, AvailableMB: 8192}
	claims := []int64{256, 512, 1024}
	errorsSeen := make(chan bool, 8)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for j := 0; j < 100; j++ {
				parsed, err := ParseCICDHostMemory(strings.NewReader("MemTotal: 8388608 kB\nMemAvailable: 8388608 kB\n"))
				if err != nil || parsed != memory || CheckCICDMemory(memory, claims, 4, 256) != nil || ValidateCICDResources(DefaultCICDResources()) != nil {
					errorsSeen <- true
					return
				}
			}
		}()
	}
	group.Wait()
	close(errorsSeen)
	if len(errorsSeen) != 0 || !reflect.DeepEqual(claims, []int64{256, 512, 1024}) {
		t.Fatal("runner.phase19-memory.concurrent.admit: shared input mutated or inconsistent result")
	}
}

type cicdMemoryFaultReader struct{}

func (cicdMemoryFaultReader) Read(p []byte) (int, error) {
	return copy(p, "MemTotal: 1048576 kB\nMemAvailable: 655360 kB\n"), errors.New("private-memory-marker")
}

func TestCICDHostMemoryParser(t *testing.T) {
	base := "MemTotal: 1048576 kB\nMemAvailable: 655360 kB\n"
	tests := []struct {
		id, text string
		want     CICDHostMemory
		valid    bool
	}{
		{"normal", base, CICDHostMemory{1024, 640}, true},
		{"padded", "MemAvailable:\t655360\tkB\nMemTotal:  983040  kB", CICDHostMemory{960, 640}, true},
		{"fraction-floor", "MemTotal: 917505 kB\nMemAvailable: 328703 kB\n", CICDHostMemory{896, 320}, true},
		{"available-zero", "MemTotal: 917504 kB\nMemAvailable: 0 kB\n", CICDHostMemory{896, 0}, true},
		{"unknown-key", base + "SwapTotal: 99999999 kB\n", CICDHostMemory{1024, 640}, true},
		{"maximum", "MemTotal: 1073741824 kB\nMemAvailable: 1073741824 kB\n", CICDHostMemory{1048576, 1048576}, true},
		{"size-exact", base + strings.Repeat(" ", 65536-len(base)), CICDHostMemory{1024, 640}, true},
		{"size-excess", base + strings.Repeat(" ", 65537-len(base)), CICDHostMemory{}, false},
		{"empty", "", CICDHostMemory{}, false},
		{"missing-available", "MemTotal: 1048576 kB\nMemFree: 655360 kB\n", CICDHostMemory{}, false},
		{"missing-total", "MemAvailable: 655360 kB\n", CICDHostMemory{}, false},
		{"duplicate-total", base + "MemTotal: 1048576 kB\n", CICDHostMemory{}, false},
		{"duplicate-available", base + "MemAvailable: 655360 kB\n", CICDHostMemory{}, false},
		{"zero-total", "MemTotal: 0 kB\nMemAvailable: 0 kB\n", CICDHostMemory{}, false},
		{"available-excess-before-floor", "MemTotal: 917504 kB\nMemAvailable: 917505 kB\n", CICDHostMemory{}, false},
		{"value-above", "MemTotal: 1073741825 kB\nMemAvailable: 0 kB\n", CICDHostMemory{}, false},
	}
	for _, value := range []string{"-1", "+1", "1.5", "1e6", "9223372036854775808", "NaN", ""} {
		tests = append(tests, struct {
			id, text string
			want     CICDHostMemory
			valid    bool
		}{"invalid-value-" + value, "MemTotal: " + value + " kB\nMemAvailable: 0 kB\n", CICDHostMemory{}, false})
	}
	for _, suffix := range []string{"MB", "kb", "kB extra", "kB\r", "kB\x00", ""} {
		tests = append(tests, struct {
			id, text string
			want     CICDHostMemory
			valid    bool
		}{"invalid-unit-" + suffix, "MemTotal: 1048576 " + suffix + "\nMemAvailable: 0 kB\n", CICDHostMemory{}, false})
	}
	for _, tc := range tests {
		t.Run(tc.id, func(t *testing.T) {
			got, err := ParseCICDHostMemory(strings.NewReader(tc.text))
			if got != tc.want || tc.valid && err != nil || !tc.valid && err != ErrCICDCapability {
				t.Fatal("runner.phase19-memory." + tc.id + ".parse: unexpected value or error")
			}
		})
	}
	for _, reader := range []io.Reader{nil, cicdMemoryFaultReader{}} {
		if got, err := ParseCICDHostMemory(reader); got != (CICDHostMemory{}) || err != ErrCICDCapability {
			t.Fatal("runner.phase19-memory.read-error.reject: error or partial data leaked")
		}
	}
	long := strings.NewReader(base + strings.Repeat(" ", 131072))
	if _, err := ParseCICDHostMemory(long); err != ErrCICDCapability || long.Len() != len(base)+131072-65537 {
		t.Fatal("runner.phase19-memory.bounded-read.reject: read exceeded detection byte")
	}
	if runtime.GOOS == "linux" {
		got, err := ReadCICDHostMemory()
		if err != nil || got.TotalMB <= 0 || got.AvailableMB < 0 || got.AvailableMB > got.TotalMB {
			t.Fatal("runner.phase19-memory.linux-read.value: physical host observation failed")
		}
	} else {
		if got, err := ReadCICDHostMemory(); got != (CICDHostMemory{}) || err != ErrCICDCapability {
			t.Fatal("runner.phase19-memory.unsupported.reject: unsupported host accepted")
		}
	}
}

type nativeBuilderFixture struct {
	Code         int    `json:"code"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	VersionFault string `json:"version_fault"`
	Trace        string `json:"trace"`
	PipelineArgs string `json:"pipeline_args"`
}

func TestMain(m *testing.M) {
	if raw := os.Getenv("RUNNER_NATIVE_BUILDER_FIXTURE"); raw != "" {
		var fixture nativeBuilderFixture
		if err := json.Unmarshal([]byte(raw), &fixture); err != nil {
			os.Exit(98)
		}
		stage := "build"
		if len(os.Args) == 2 && os.Args[1] == "--version" {
			stage = "version"
		}
		if fixture.Trace != "" {
			file, err := os.OpenFile(fixture.Trace, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				os.Exit(98)
			}
			_, writeErr := fmt.Fprintln(file, stage)
			closeErr := file.Close()
			if writeErr != nil || closeErr != nil {
				os.Exit(98)
			}
		}
		if stage == "version" {
			switch fixture.VersionFault {
			case "timeout":
				time.Sleep(time.Hour)
			case "stderr":
				fmt.Fprintln(os.Stderr, "warn")
			case "empty_go_token":
				fmt.Println("adlaire-ci-build V.0.0-dev go=")
				os.Exit(0)
			case "extra_token":
				fmt.Println("adlaire-ci-build V.0.0-dev go=fake extra")
				os.Exit(0)
			}
			fmt.Println("adlaire-ci-build V.0.0-dev go=fake")
			os.Exit(0)
		}
		if fixture.PipelineArgs != "" {
			for key, want := range map[string]string{"PIPELINE_ENV": "pipeline", "BRANCH_VALUE": "pipeline-wins", "BRANCH_ONLY": "branch-only"} {
				if os.Getenv(key) != want {
					fmt.Fprintln(os.Stderr, "pipeline environment mismatch")
					os.Exit(6)
				}
			}
			if err := os.WriteFile(fixture.PipelineArgs, []byte(strings.Join(os.Args[1:], "\n")+"\n"), 0600); err != nil {
				os.Exit(98)
			}
			if err := os.WriteFile(filepath.Join(os.Getenv("ADLAIRE_CI_STATE_DIR"), ".pipeline_config"), []byte(`{"extra_args":[],"env":{"MUTATED":"yes"}}`), 0600); err != nil {
				os.Exit(98)
			}
			fmt.Fprintln(os.Stdout, "secret="+os.Getenv("PIPELINE_SECRET"))
		}
		if fixture.Code == 0 {
			out := os.Getenv("ADLAIRE_CI_OUT")
			if out == "" || os.MkdirAll(filepath.Join(out, "assets"), 0755) != nil {
				os.Exit(98)
			}
			for name, body := range map[string]string{"index.html": "<html></html>", "assets/style.css": "body{}", "assets/app.js": `console.log("ok")`, "assets/search-index.json": "[]"} {
				if err := os.WriteFile(filepath.Join(out, name), []byte(body), 0644); err != nil {
					os.Exit(98)
				}
			}
		}
		fmt.Fprint(os.Stdout, fixture.Stdout)
		fmt.Fprint(os.Stderr, fixture.Stderr)
		os.Exit(fixture.Code)
	}
	os.Exit(m.Run())
}

func useNativeBuilder(t *testing.T, fixture nativeBuilderFixture) {
	t.Helper()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(fixture)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADLAIRE_CI_BUILD_BIN", executable)
	t.Setenv("RUNNER_NATIVE_BUILDER_FIXTURE", string(data))
	t.Setenv("GORACE", "atexit_sleep_ms=0")
}

func rejectShellStartup(t *testing.T) {
	t.Helper()
	startup := filepath.Join(t.TempDir(), "startup")
	if err := os.WriteFile(startup, []byte("exit 97\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BASH_ENV", startup)
	if err := exec.Command("bash", "--noprofile", "--norc", "-c", "exit 0").Run(); err == nil {
		t.Fatal("shell startup negative control did not fail")
	}
}

func TestRunnerFixtureR1CLI(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := RunRunner([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("help exit=%d", code)
	}
	if strings.TrimSpace(stdout.String()) != "Usage: adlaire-ci-runner [--state-dir path] [--once] [--dry-run] [--version] [--help]" || stderr.Len() != 0 {
		t.Fatalf("unexpected help stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunRunner([]string{"--state-dir", "relative"}, &stdout, &stderr); code != 2 {
		t.Fatalf("relative exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "state directory must be absolute: relative") {
		t.Fatalf("missing relative error: %s", stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunRunner([]string{"--unknown"}, &stdout, &stderr); code != 2 {
		t.Fatalf("unknown exit=%d", code)
	}
	if !strings.Contains(stderr.String(), "unknown option: --unknown") {
		t.Fatalf("missing unknown error: %s", stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunRunner([]string{"--version", "--state-dir", "relative"}, &stdout, &stderr); code != 0 {
		t.Fatalf("version exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout.String()), "adlaire-ci-runner V.0.0-dev go=") || stderr.Len() != 0 {
		t.Fatalf("unexpected version stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunnerFixtureR2NoChange(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"blob-1\"}\n" {
			t.Fatalf("sha changed")
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
			t.Fatalf("history must not exist")
		}
		if !strings.Contains(stdout.String(), "NO_CHANGE: branch=main target=docs sha=blob-1") {
			t.Fatalf("missing no change log: %s", stdout.String())
		}
	})
}

func TestRunnerFixtureR3BuildSuccessNoDeploy(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	rejectShellStartup(t)
	trace := filepath.Join(state, "builder-trace")
	useNativeBuilder(t, nativeBuilderFixture{Stdout: `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, Trace: trace})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 0 {
			log := onlyBuildLog(t, state)
			problem := ""
			if log.Error != nil {
				problem = *log.Error
			}
			t.Fatalf("exit=%d stdout=%s stderr=%s status=%s error=%s", code, stdout.String(), stderr.String(), log.TargetStatus, problem)
		}
		if readFile(t, trace) != "version\nbuild\n" {
			t.Fatal("native builder invocation order mismatch")
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"new-blob\"}\n" {
			t.Fatalf("sha not updated: %s", readFile(t, filepath.Join(state, ".last_sha")))
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "success" || log.Pipeline.ExitCode == nil || *log.Pipeline.ExitCode != 0 || log.Report == nil || log.Report.Pages != 1 || len(log.Deploy) != 0 || log.Error != nil {
			t.Fatalf("unexpected log: %+v", log)
		}
		if log.SnapshotID == nil {
			t.Fatalf("snapshot id must be recorded")
		}
		if log.Status != "success" || log.OutputSHA256 == nil || log.OutputSizeBytes == nil || log.Trigger != "polling" {
			t.Fatalf("normalized log fields missing: %+v", log)
		}
		var manifest runnerDependencyManifest
		readJSON(t, filepath.Join(state, "dist", "site", ".dependency_manifest.json"), &manifest)
		if manifest.Version != 1 || manifest.Builder.Name != "adlaire-ci-build" || len(manifest.Pages) != 1 || manifest.GeneratedOutputs[".dependency_manifest.json"] != "manifest" {
			t.Fatalf("dependency manifest mismatch: %+v", manifest)
		}
		if log.BuildMeta == nil || log.BuildMeta.BuildID != log.ID || log.BuildMeta.BuildAt != log.StartedAt || log.TransferVerified != nil || log.RemoteBuild != nil || log.CommitStatusLog != nil {
			t.Fatalf("schema auxiliary fields mismatch: %+v", log)
		}
		if log.Environment.RunnerVersion != "V.0.0-dev" || log.Environment.GoVersion == "" || log.Environment.OS == "" || log.Environment.Arch == "" {
			t.Fatalf("environment record missing: %+v", log.Environment)
		}
		snapshotDir := filepath.Join(state, ".snapshots", *log.SnapshotID)
		var meta snapshotMeta
		readJSON(t, filepath.Join(snapshotDir, "meta.json"), &meta)
		if meta.ID != *log.SnapshotID || meta.BuildID != *log.SnapshotID || meta.FileCount == 0 || meta.OutputSHA256 == "" {
			t.Fatalf("snapshot meta mismatch: %+v", meta)
		}
		if !tarGzContains(t, filepath.Join(snapshotDir, "site.tar.gz"), "index.html") {
			t.Fatalf("snapshot archive missing index.html")
		}
		if _, err := os.Stat(filepath.Join(snapshotDir, "site")); !os.IsNotExist(err) {
			t.Fatalf("legacy snapshot site directory must not exist")
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".build_history")), `"status":"success"`) {
			t.Fatalf("history missing success")
		}
		var summary buildStatusSummary
		readJSON(t, filepath.Join(state, ".build_status.json"), &summary)
		if summary.Status != "success" || summary.Running || summary.LastBuildID == nil || *summary.LastBuildID != log.ID || summary.OutputSHA256 == nil {
			t.Fatalf("unexpected build status summary: %+v", summary)
		}
		var trends buildTrendsFile
		readJSON(t, filepath.Join(state, ".build_trends.json"), &trends)
		if trends.SchemaVersion != 1 || len(trends.Samples) != 1 || trends.Samples[0].BuildID != log.ID || trends.Summary.Count != 1 || trends.Summary.AvgSeconds == nil {
			t.Fatalf("unexpected build trends: %+v", trends)
		}
		var bs buildState
		readJSON(t, filepath.Join(state, ".build_state"), &bs)
		if bs.Running || bs.CurrentBuildID != nil {
			t.Fatalf("state not finalized: %+v", bs)
		}
		if _, err := os.Stat(filepath.Join(state, ".build_lock")); !os.IsNotExist(err) {
			t.Fatalf("lock remains")
		}
	})
}

func TestRunnerFixtureR4PipelineFailure(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	rejectShellStartup(t)
	trace := filepath.Join(state, "builder-trace")
	useNativeBuilder(t, nativeBuilderFixture{Code: 7, Stdout: "before fail", Stderr: "failed", Trace: trace})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, trace) != "version\nbuild\n" {
			t.Fatal("native builder invocation order mismatch")
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"old-blob\"}\n" {
			t.Fatalf("sha changed on failure")
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "failure_build" || log.Pipeline.ExitCode == nil || *log.Pipeline.ExitCode != 7 || log.Pipeline.Stdout != "before fail" || log.Pipeline.Stderr != "failed" || log.Error == nil || *log.Error != "pipeline failed" {
			t.Fatalf("unexpected failure log: %+v", log)
		}
		if log.FailureCategory == nil || *log.FailureCategory != "pipeline_exit" || log.Status != "failure" {
			t.Fatalf("failure category not recorded: %+v", log)
		}
		if len(log.FailureEvidence) != 1 || log.FailureEvidence[0].Source != "pipeline" || log.FailureEvidence[0].Code != "pipeline_exit" || log.FailureEvidence[0].Message != "pipeline exited non-zero" {
			t.Fatalf("failure evidence not recorded: %+v", log.FailureEvidence)
		}
	})
}

func TestRunnerPhase2PipelineConfigAppliesArgsEnvAndMasksSecret(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	argsFile := filepath.Join(state, "args.txt")
	trace := filepath.Join(state, "builder-trace")
	useNativeBuilder(t, nativeBuilderFixture{PipelineArgs: argsFile, Trace: trace, Stdout: "[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default"})
	rejectShellStartup(t)
	var branches branchConfigFile
	readJSON(t, filepath.Join(state, ".branch_config"), &branches)
	branches.BranchTargets[0].Env = map[string]string{"BRANCH_VALUE": "branch", "BRANCH_ONLY": "branch-only"}
	writeRunnerTestJSON(t, filepath.Join(state, ".branch_config"), branches)
	writeRunnerTestJSON(t, filepath.Join(state, ".pipeline_config"), pipelineConfigFile{
		ExtraArgs: []string{"--theme", "adlaire-alt", "--flag=value"},
		Env: map[string]string{
			"BRANCH_VALUE":    "pipeline-wins",
			"PIPELINE_ENV":    "pipeline",
			"PIPELINE_SECRET": "super-secret",
		},
	})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("runner.pipeline-env.native.result: expected exit0 actual=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if got := readFile(t, trace); got != "version\nbuild\n" {
			t.Fatalf("runner.pipeline-env.native.order: fixture.md#15a-f-runner-初期受け入れ-fixture-契約 expected version/build actual=%q", got)
		}
		args := strings.Split(strings.TrimSpace(readFile(t, argsFile)), "\n")
		wantSuffix := []string{"--theme", "adlaire-alt", "--flag=value"}
		if len(args) < len(wantSuffix) || strings.Join(args[len(args)-len(wantSuffix):], "\n") != strings.Join(wantSuffix, "\n") {
			t.Fatalf("runner.pipeline-env.native.args: expected pipeline extra args appended actual=%#v", args)
		}
		var changed pipelineConfigFile
		readJSON(t, filepath.Join(state, ".pipeline_config"), &changed)
		if len(changed.Env) != 1 || changed.Env["MUTATED"] != "yes" {
			t.Fatal("runner.pipeline-env.native.snapshot: expected fake config replacement, actual unchanged")
		}
		log := onlyBuildLog(t, state)
		if strings.Contains(log.Pipeline.Stdout, "super-secret") || !strings.Contains(log.Pipeline.Stdout, "secret=***") {
			t.Fatalf("runner.pipeline-env.native.mask: expected secret masked actual=%q", log.Pipeline.Stdout)
		}
		if log.Environment.BuilderVersion != "V.0.0-dev" || log.Environment.Hostname == "" || log.Environment.StateDir == "" || log.Environment.CapturedAt == "" || log.Environment.DiskFreeBytes == nil {
			t.Fatalf("environment record incomplete: %+v", log.Environment)
		}
		if got, want := strings.Join(log.Environment.EnvKeys, ","), "BRANCH_ONLY,BRANCH_VALUE,PIPELINE_ENV,PIPELINE_SECRET"; got != want {
			t.Fatalf("runner.pipeline-env.native.env-keys: expected=%q actual=%q", want, got)
		}
	})

	for _, wrong := range []string{"PIPELINE_ENV", "BRANCH_VALUE", "BRANCH_ONLY"} {
		t.Run("native-reject-"+wrong, func(t *testing.T) {
			caseID := strings.ReplaceAll(strings.ToLower(wrong), "_", "-")
			root := t.TempDir()
			fixture := nativeBuilderFixture{PipelineArgs: filepath.Join(root, "args")}
			data, err := json.Marshal(fixture)
			if err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(os.Getenv("ADLAIRE_CI_BUILD_BIN"), "--out", filepath.Join(root, "out"))
			for _, value := range os.Environ() {
				key, _, _ := strings.Cut(value, "=")
				if key != "RUNNER_NATIVE_BUILDER_FIXTURE" && key != "PIPELINE_ENV" && key != "BRANCH_VALUE" && key != "BRANCH_ONLY" && key != "ADLAIRE_CI_OUT" && key != "ADLAIRE_CI_STATE_DIR" {
					cmd.Env = append(cmd.Env, value)
				}
			}
			cmd.Env = append(cmd.Env, "RUNNER_NATIVE_BUILDER_FIXTURE="+string(data), "ADLAIRE_CI_OUT="+filepath.Join(root, "out"), "ADLAIRE_CI_STATE_DIR="+state)
			for key, value := range map[string]string{"PIPELINE_ENV": "pipeline", "BRANCH_VALUE": "pipeline-wins", "BRANCH_ONLY": "branch-only"} {
				if key == wrong {
					value = "wrong"
				}
				cmd.Env = append(cmd.Env, key+"="+value)
			}
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err = cmd.Run()
			var exit *exec.ExitError
			if !errors.As(err, &exit) || exit.ExitCode() != 6 || stdout.Len() != 0 || stderr.String() != "pipeline environment mismatch\n" {
				t.Fatalf("runner.pipeline-env.native.%s.reject: fixture.md#15a-f-runner-初期受け入れ-fixture-契約 expected exit6/empty stdout/fixed reason actual=%v/%q/%q", caseID, err, stdout.String(), stderr.String())
			}
			for _, name := range []string{fixture.PipelineArgs, filepath.Join(root, "out")} {
				if _, err := os.Lstat(name); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("runner.pipeline-env.native.%s.%s.no-write: expected absent fake output actual=%v", caseID, filepath.Base(name), err)
				}
			}
		})
	}
}

func TestRunnerPhase2BuilderVersionPrecheckStrict(t *testing.T) {
	for _, fault := range []string{"stderr", "empty_go_token", "extra_token", "timeout"} {
		t.Run(fault, func(t *testing.T) {
			if fault == "timeout" {
				old := runnerBuilderVersionTimeout
				runnerBuilderVersionTimeout = 20 * time.Millisecond
				t.Cleanup(func() { runnerBuilderVersionTimeout = old })
			}
			state := newRunnerState(t, "old-blob", nil)
			trace := filepath.Join(state, "builder-trace")
			useNativeBuilder(t, nativeBuilderFixture{VersionFault: fault, Trace: trace})
			prefix := "build binary version mismatch: "
			if fault == "stderr" {
				prefix = "build binary version wrote stderr: "
			} else if fault == "timeout" {
				prefix = "build binary version timeout: "
			}
			if _, err := precheckRunnerTarget(RunnerConfig{}, BranchTarget{Out: filepath.Join(state, "dist", "site")}); err == nil || !strings.HasPrefix(err.Error(), prefix) {
				t.Fatalf("wrong internal precheck cause: expected=%s actual=%v", prefix, err)
			}
			server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
			withRunnerServer(t, server.URL, func() {
				var stdout, stderr bytes.Buffer
				if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 1 {
					t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
				}
				data, err := os.ReadFile(trace)
				if (err != nil && !os.IsNotExist(err)) || (fault != "timeout" && string(data) != "version\nversion\n") || strings.Contains(string(data), "build\n") {
					t.Fatalf("builder must not start after precheck failure: trace=%q err=%v", data, err)
				}
				log := onlyBuildLog(t, state)
				if log.TargetStatus != "failure_precheck" || log.FailureCategory == nil || *log.FailureCategory != "resource_error" || log.Environment.BuilderVersion != "" {
					t.Fatalf("unexpected precheck log: %+v", log)
				}
				if log.Error == nil || *log.Error != "precheck failed" {
					t.Fatal("public precheck error must remain fixed")
				}
			})
		})
	}
}

func TestRunnerPhase2PipelineConfigInvalidDoesNotStartBuild(t *testing.T) {
	cases := []struct {
		name    string
		content string
	}{
		{
			name:    "reserved_arg",
			content: `{"extra_args":["--src=/tmp/other"],"env":{}}`,
		},
		{
			name:    "reserved_env",
			content: `{"extra_args":[],"env":{"ADLAIRE_CI_BAD":"x"}}`,
		},
		{
			name:    "unknown_key",
			content: `{"extra_args":[],"env":{},"unknown":true}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := newRunnerState(t, "old-blob", nil)
			rejectShellStartup(t)
			trace := filepath.Join(state, "builder-trace")
			useNativeBuilder(t, nativeBuilderFixture{Trace: trace})
			caseID := strings.ReplaceAll(tc.name, "_", "-")
			if err := os.WriteFile(filepath.Join(state, ".pipeline_config"), []byte(tc.content+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
			withRunnerServer(t, server.URL, func() {
				var stdout, stderr bytes.Buffer
				code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
				log := onlyBuildLog(t, state)
				if code != 2 {
					t.Fatalf("runner.pipeline-config-invalid.%s.result: fixture.md#15a-f-runner-初期受け入れ-fixture-契約 expected exit=2 actual=%d stdout=%q stderr=%q status=%s category=%s error=%s; reproduce go test ./components/runner -run TestRunnerPhase2PipelineConfigInvalidDoesNotStartBuild -count=1", caseID, code, stdout.String(), stderr.String(), log.TargetStatus, derefString(log.FailureCategory, ""), derefString(log.Error, ""))
				}
				if actual := readFile(t, trace); actual != "version\n" {
					t.Fatalf("runner.pipeline-config-invalid.%s.native-trace: expected version-only trace actual=%q", caseID, actual)
				}
				if _, err := os.Stat(filepath.Join(state, "dist", "site")); !os.IsNotExist(err) {
					t.Fatalf("runner.pipeline-config-invalid.%s.no-output: expected output absent actual=%v", caseID, err)
				}
				if actual := readFile(t, filepath.Join(state, ".last_sha")); actual != "{\"sha\":\"old-blob\"}\n" {
					t.Fatalf("runner.pipeline-config-invalid.%s.sha-unchanged: expected old SHA actual=%q", caseID, actual)
				}
				if log.TargetStatus != "failure_pipeline_config" || log.Pipeline.ExitCode == nil || *log.Pipeline.ExitCode != 2 || log.RetryCount != 0 || len(log.Attempts) != 1 {
					t.Fatalf("runner.pipeline-config-invalid.%s.log: expected failure_pipeline_config/exit2/retry0/attempt1 actual=%+v", caseID, log)
				}
				if log.FailureCategory == nil || *log.FailureCategory != "config_error" || log.Error == nil || *log.Error != "pipeline config invalid" {
					t.Fatalf("runner.pipeline-config-invalid.%s.category: expected config_error/pipeline config invalid actual=%+v", caseID, log)
				}
				if len(log.FailureEvidence) != 1 || log.FailureEvidence[0].Source != "config" || log.FailureEvidence[0].Code != "config_error" {
					t.Fatalf("runner.pipeline-config-invalid.%s.evidence: expected one config/config_error evidence actual=%+v", caseID, log.FailureEvidence)
				}
			})
		})
	}
}

func TestRunnerCompletionServerConfigSparseSchema(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"log_max_lines":              750,
		"history_max_count":          500,
		"history_retention":          map[string]any{"enabled": false, "max_count": 1000, "max_age_days": nil, "updated_at": nil},
		"build_timeout_seconds":      120,
		"log_retention_days":         30,
		"log_level":                  "INFO",
		"pat_expires_at":             nil,
		"snapshots_keep":             0,
		"queue_max_size":             3,
		"build_retry_max":            1,
		"build_retry_base_seconds":   5,
		"commit_status_enabled":      false,
		"commit_status_context":      "Adlaire CI",
		"commit_status_target_url":   nil,
		"log_archive_after_days":     0,
		"build_trend_keep_count":     1000,
		"duration_anomaly":           map[string]any{"enabled": false, "min_samples": 20, "avg_multiplier": 2.0, "p95_multiplier": 1.5},
		"watch_mode":                 "github",
		"tag_filter":                 map[string]any{"enabled": false, "patterns": []any{}},
		"build_cache_enabled":        false,
		"deploy_parallelism":         1,
		"remote_build":               map[string]any{"enabled": false, "host": nil, "user": nil, "work_dir": nil, "command_args": []any{}, "artifact_path": nil},
		"approval_timeout_seconds":   86400,
		"force_build_interval_hours": 0,
		"build_cooldown_seconds":     0,
		"schedule_interval_seconds":  300,
		"schedule_paused":            false,
		"allowed_hours":              nil,
		"session_timeout_seconds":    28800,
		"api_rate_limit": map[string]any{"enabled": true, "groups": map[string]any{
			"login": map[string]any{"window_seconds": 60, "max_requests": 10},
		}},
	})
	cfg := defaultRunnerConfig(state)
	if err := loadRunnerConfig(&cfg, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))); err != nil {
		t.Fatalf("valid sparse server config rejected: %v", err)
	}
	if cfg.BuildTimeoutSeconds != 120 || cfg.SnapshotsKeep != 0 || cfg.BuildRetryMax != 1 || cfg.WatchMode != "github" || cfg.DeployParallelism != 1 {
		t.Fatalf("server config not applied: %+v", cfg)
	}
}

func TestRunnerCompletionServerConfigRejectsUnknownKey(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{"unknown": true})
	cfg := defaultRunnerConfig(state)
	err := loadRunnerConfig(&cfg, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	if err == nil || !strings.Contains(err.Error(), "SERVER_CONFIG_INVALID") {
		t.Fatalf("unknown server config key accepted: %v", err)
	}
}

func TestRunnerCompletionDryRunNoWrites(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	for _, path := range []string{filepath.Join(state, "repo"), filepath.Join(state, "dist"), filepath.Join(state, ".build_logs"), filepath.Join(state, ".build_lock")} {
		_ = os.RemoveAll(path)
	}
	var stdout, stderr bytes.Buffer
	if code := RunRunner([]string{"--state-dir", state, "--dry-run"}, &stdout, &stderr); code != 0 {
		t.Fatalf("dry-run exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	var result map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &result); err != nil {
		t.Fatalf("invalid dry-run json: %v stdout=%s", err, stdout.String())
	}
	if result["dry_run"] != true {
		t.Fatalf("dry_run flag missing: %+v", result)
	}
	for _, path := range []string{filepath.Join(state, "repo"), filepath.Join(state, "dist"), filepath.Join(state, ".build_logs"), filepath.Join(state, ".build_lock"), filepath.Join(state, ".build_history")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("dry-run wrote path %s", path)
		}
	}
}

func TestRunnerCompletionLocalWatchDoesNotRequireToken(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	if err := os.Remove(filepath.Join(state, ".github_token")); err != nil {
		t.Fatal(err)
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"watch_mode":             "local",
		"build_cooldown_seconds": 0,
	})
	src := filepath.Join(state, "repo", "docs")
	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "source.md"), []byte("# Local\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	log := onlyBuildLog(t, state)
	if log.Environment.WatchMode != "local" || log.BlobSHA == nil || *log.BlobSHA == "old-blob" {
		t.Fatalf("local watch log mismatch: %+v", log)
	}
}

func TestRunnerPhase2HooksPrePostSuccess(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	preMarker := filepath.Join(t.TempDir(), "pre.txt")
	postMarker := filepath.Join(t.TempDir(), "post.txt")
	preHook := filepath.Join(t.TempDir(), "pre-hook")
	postHook := filepath.Join(t.TempDir(), "post-hook")
	writeExecutable(t, preHook, fmt.Sprintf("#!/usr/bin/env bash\nset -u\nprintf '%%s' \"$MY_SECRET\"\nprintf '%%s' \"$ADLAIRE_CI_HOOK_PHASE:$ADLAIRE_CI_BRANCH\" > %q\n", preMarker))
	writeExecutable(t, postHook, fmt.Sprintf("#!/usr/bin/env bash\nset -u\nprintf '%%s' \"$ADLAIRE_CI_HOOK_PHASE:$ADLAIRE_CI_BRANCH\" > %q\n", postMarker))
	target := BranchTarget{
		Branch: "main", TargetFile: "docs", SHAFile: filepath.Join(state, ".last_sha"),
		Src: filepath.Join(state, "repo", "docs"), Out: filepath.Join(state, "dist", "site"),
		Env: map[string]string{"MY_SECRET": "super-secret"},
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".branch_config"), branchConfigFile{BranchTargets: []BranchTarget{target}})
	writeRunnerTestJSON(t, filepath.Join(state, ".hooks"), hooksFile{Hooks: []runnerHook{
		{ID: "disabled", Phase: "pre", CommandArgs: []string{preHook}, Enabled: false, AbortOnFailure: true, TimeoutSeconds: 300, CreatedAt: "2026-09-16T00:00:00Z"},
		{ID: "pre", Phase: "pre", CommandArgs: []string{preHook}, Enabled: true, AbortOnFailure: true, TimeoutSeconds: 300, CreatedAt: "2026-09-16T00:00:00Z"},
		{ID: "post", Phase: "post", CommandArgs: []string{postHook}, Enabled: true, AbortOnFailure: false, TimeoutSeconds: 300, CreatedAt: "2026-09-16T00:00:00Z"},
	}})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.Status != "success" || log.TargetStatus != "success" || len(log.SkippedHookIDs) != 1 || log.SkippedHookIDs[0] != "disabled" {
			t.Fatalf("unexpected hook build log: %+v", log)
		}
		if readFile(t, preMarker) != "pre:main" || readFile(t, postMarker) != "post:main" {
			t.Fatalf("hook markers missing pre=%q post=%q", readFile(t, preMarker), readFile(t, postMarker))
		}
		preLog := readHookLog(t, state, log.ID, "pre")
		if preLog.Status != "success" || preLog.Stdout != "***" || preLog.Truncated {
			t.Fatalf("unexpected pre hook log: %+v", preLog)
		}
		postLog := readHookLog(t, state, log.ID, "post")
		if postLog.Status != "success" || postLog.Phase != "post" {
			t.Fatalf("unexpected post hook log: %+v", postLog)
		}
	})
}

func TestRunnerPhase2HooksPreAbortSkipsBuild(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipelineWithExtra(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "", `printf builder > "$out/builder-ran"`)
	preHook := filepath.Join(t.TempDir(), "pre-hook")
	writeExecutable(t, preHook, "#!/usr/bin/env bash\nprintf pre-failed\nexit 9\n")
	writeRunnerTestJSON(t, filepath.Join(state, ".hooks"), hooksFile{Hooks: []runnerHook{
		{ID: "pre_abort", Phase: "pre", CommandArgs: []string{preHook}, Enabled: true, AbortOnFailure: true, TimeoutSeconds: 300, CreatedAt: "2026-09-16T00:00:00Z"},
	}})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"old-blob\"}\n" {
			t.Fatalf("sha changed after hook abort")
		}
		if _, err := os.Stat(filepath.Join(state, "dist", "site", "builder-ran")); !os.IsNotExist(err) {
			t.Fatalf("builder ran after hook abort")
		}
		log := onlyBuildLog(t, state)
		if log.Status != "failure" || log.TargetStatus != "hook_error" || log.Error == nil || *log.Error != "pre hook failed" {
			t.Fatalf("unexpected abort log: %+v", log)
		}
		hookLog := readHookLog(t, state, log.ID, "pre_abort")
		if hookLog.Status != "failure" || hookLog.ExitCode == nil || *hookLog.ExitCode != 9 || hookLog.Stdout != "pre-failed" {
			t.Fatalf("unexpected hook log: %+v", hookLog)
		}
	})
}

func TestRunnerPhase2HooksPostFailureKeepsBuildStatus(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	postHook := filepath.Join(t.TempDir(), "post-hook")
	writeExecutable(t, postHook, "#!/usr/bin/env bash\nprintf post-failed >&2\nexit 5\n")
	writeRunnerTestJSON(t, filepath.Join(state, ".hooks"), hooksFile{Hooks: []runnerHook{
		{ID: "post_fail", Phase: "post", CommandArgs: []string{postHook}, Enabled: true, AbortOnFailure: false, TimeoutSeconds: 300, CreatedAt: "2026-09-16T00:00:00Z"},
	}})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"new-blob\"}\n" {
			t.Fatalf("sha not updated after successful build")
		}
		log := onlyBuildLog(t, state)
		if log.Status != "success" || log.TargetStatus != "success" || !runnerTestContainsString(log.Warnings, "POST_HOOK_FAILED:post_fail") {
			t.Fatalf("post failure changed build result or warning missing: %+v", log)
		}
		hookLog := readHookLog(t, state, log.ID, "post_fail")
		if hookLog.Status != "failure" || hookLog.ExitCode == nil || *hookLog.ExitCode != 5 || hookLog.Stderr != "post-failed" {
			t.Fatalf("unexpected hook log: %+v", hookLog)
		}
	})
}

func TestRunnerCompletionMultiFileChangedTargets(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipelineWithExtra(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "", `if [ "${ADLAIRE_CHANGED_TARGETS:-}" != '["docs/a.md","docs/b.md"]' ]; then exit 8; fi`)
	target := BranchTarget{
		Branch: "main", TargetFile: "docs/a.md", TargetFiles: []string{"docs/b.md", "docs/a.md"},
		SHAFile: filepath.Join(state, ".last_sha"), Src: filepath.Join(state, "repo", "docs"), Out: filepath.Join(state, "dist", "site"),
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".branch_config"), branchConfigFile{BranchTargets: []BranchTarget{target}})
	encodedA := base64.StdEncoding.EncodeToString([]byte("# A\n"))
	encodedB := base64.StdEncoding.EncodeToString([]byte("# B\n"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{
				{"path": "docs/a.md", "type": "blob", "sha": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
				{"path": "docs/b.md", "type": "blob", "sha": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
			}})
		case strings.Contains(r.URL.Path, "/git/blobs/aaaaaaaa"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encodedA, "encoding": "base64"})
		case strings.Contains(r.URL.Path, "/git/blobs/bbbbbbbb"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encodedB, "encoding": "base64"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if strings.Join(log.TargetFiles, ",") != "docs/a.md,docs/b.md" || strings.Join(log.ChangedTargets, ",") != "docs/a.md,docs/b.md" {
			t.Fatalf("target tracking mismatch: %+v", log)
		}
		if _, err := os.Stat(filepath.Join(state, "repo", "docs", "docs", "a.md")); err != nil {
			t.Fatalf("source a not materialized: %v", err)
		}
		if readFile(t, filepath.Join(state, ".last_sha")) == "{\"sha\":\"old-blob\"}\n" {
			t.Fatalf("multi-file digest not updated")
		}
	})
}

func TestRunnerCompletionApprovalRequiredSkipsBuild(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	target := BranchTarget{
		Branch: "main", TargetFile: "docs", TargetFiles: []string{"docs"},
		SHAFile: filepath.Join(state, ".last_sha"), Src: filepath.Join(state, "repo", "docs"), Out: filepath.Join(state, "dist", "site"),
		ApprovalRequired: true,
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".branch_config"), branchConfigFile{BranchTargets: []BranchTarget{target}})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
			t.Fatalf("approval pending must not create history")
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"old-blob\"}\n" {
			t.Fatalf("approval pending updated sha")
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".approval_queue")), `"status":"pending"`) {
			t.Fatalf("approval queue missing")
		}
		var summary buildStatusSummary
		readJSON(t, filepath.Join(state, ".build_status.json"), &summary)
		if summary.Status != "skipped" || summary.LastTargetStatus == nil || *summary.LastTargetStatus != "pending_approval" {
			t.Fatalf("approval status mismatch: %+v", summary)
		}
	})
}

func TestRunnerCompletionReportDuplicateUsesFirst(t *testing.T) {
	report, warnings := parseRunnerReport("[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=first\n[REPORT] pages=2 headings=2 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=second\n")
	if report == nil || report.Pages != 1 || report.Theme != "first" {
		t.Fatalf("first report not preserved: %+v", report)
	}
	if len(warnings) != 1 || warnings[0] != "REPORT_DUPLICATE" {
		t.Fatalf("missing duplicate warning: %+v", warnings)
	}
}

func TestRunnerCompletionTrimLogUTF8Safe(t *testing.T) {
	raw := strings.Repeat("a", 1024*1024-1) + "あいう"
	trimmed, truncated := trimLog(raw)
	if !truncated {
		t.Fatalf("expected truncation")
	}
	if !strings.Contains(trimmed, "あいう") || strings.Contains(trimmed, "\uFFFD") {
		t.Fatalf("trimmed log is not utf8 safe: suffix=%q", trimmed[len(trimmed)-12:])
	}
}

func TestRunnerCompletionArchivesOldBuildLogs(t *testing.T) {
	state := t.TempDir()
	dir := filepath.Join(state, ".build_logs")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(dir, "old.json")
	newer := filepath.Join(dir, "new.json")
	if err := os.WriteFile(old, []byte(`{"id":"old"}`+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newer, []byte(`{"id":"new"}`+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	oldNow := runnerNow
	runnerNow = func() time.Time { return time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC) }
	defer func() { runnerNow = oldNow }()
	oldTime := runnerNow().Add(-48 * time.Hour)
	if err := os.Chtimes(old, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}
	newTime := runnerNow()
	if err := os.Chtimes(newer, newTime, newTime); err != nil {
		t.Fatal(err)
	}
	cleanupBuildLogs(dir, 1, 1, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatalf("old log not removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "archive", "old.json.gz")); err != nil {
		t.Fatalf("old log not archived: %v", err)
	}
	if _, err := os.Stat(newer); err != nil {
		t.Fatalf("new log removed: %v", err)
	}
}

func TestRunnerHardeningRetriesGitHubAPI(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	hits := 0
	encoded := base64.StdEncoding.EncodeToString([]byte("# Title\n"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/git/trees/") {
			hits++
			if hits == 1 {
				http.Error(w, "temporary", http.StatusInternalServerError)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": "docs", "type": "blob", "sha": "retry-blob"}}})
			return
		}
		if strings.Contains(r.URL.Path, "/git/blobs/") {
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
			return
		}
		http.NotFound(w, r)
	}))
	oldSleep := runnerSleep
	runnerSleep = func(time.Duration) {}
	defer func() { runnerSleep = oldSleep }()
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if hits != 2 {
			t.Fatalf("expected retry hits=2 got %d", hits)
		}
	})
}

func TestRunnerHardeningPrecheckFailure(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	t.Setenv("ADLAIRE_CI_BUILD_BIN", filepath.Join(state, "missing-build-bin"))
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "failure_precheck" || log.Error == nil || *log.Error != "precheck failed" {
			t.Fatalf("unexpected precheck log: %+v", log)
		}
	})
}

func TestRunnerHardeningNotifyPendingRetry(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	received := 0
	hook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received++
		w.WriteHeader(http.StatusNoContent)
	}))
	defer hook.Close()
	writeRunnerTestJSON(t, filepath.Join(state, ".notify_pending"), []notifyPendingEntry{{
		Event: "success", URL: hook.URL, Payload: map[string]any{"ok": true},
		QueuedAt: "2026-09-16T00:00:00Z", RetryCount: 1,
	}})
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if received != 1 {
			t.Fatalf("expected one notify retry, got %d", received)
		}
		if strings.TrimSpace(readFile(t, filepath.Join(state, ".notify_pending"))) != "[]" {
			t.Fatalf("notify pending must be empty")
		}
	})
}

func TestRunnerHardeningCircuitOpenSkipsPolling(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	now := "2026-09-16T00:00:00Z"
	msg := "github api failed"
	writeRunnerTestJSON(t, filepath.Join(state, ".build_circuit_state"), buildCircuitState{Open: true, ConsecutiveFailures: 3, OpenedAt: &now, LastFailureAt: &now, LastError: &msg})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
			t.Fatalf("history must not be touched while circuit is open")
		}
		if !strings.Contains(stdout.String(), "CIRCUIT_OPEN") {
			t.Fatalf("missing circuit log: %s", stdout.String())
		}
	})
}

func TestRunnerCompletionPendingTransferRetrySuccess(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	out := filepath.Join(state, "dist", "site")
	if err := os.MkdirAll(out, 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(out, "index.html")
	if err := os.WriteFile(file, []byte("<html></html>"), 0644); err != nil {
		t.Fatal(err)
	}
	sha, err := runnerFileSHA256(file)
	if err != nil {
		t.Fatal(err)
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".pending_transfers"), []pendingTransfer{{BuildID: "b20260916000000", Trigger: "deploy", SourceKind: "output", Branch: "main", TargetID: "deploy-1", Out: &out, Host: "host", User: "deploy", DestDir: "/var/www/html", OutputSHA256: sha, FailedAt: "2026-09-16T00:00:00Z", RetryCount: 1, LastError: "previous failure"}})
	fakeBin := t.TempDir()
	writeExecutable(t, filepath.Join(fakeBin, "ssh"), "#!/bin/sh\ncase \"$2\" in\n  sha256sum*) printf '"+sha+"  file\\n'; exit 0 ;;\n  mkdir*) cat >/dev/null; exit 0 ;;\nesac\nexit 1\n")
	t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if strings.TrimSpace(readFile(t, filepath.Join(state, ".pending_transfers"))) != "[]" {
			t.Fatalf("pending transfers not cleared")
		}
	})
}

func TestRunnerCompletionDeployUploadsAndVerifies(t *testing.T) {
	out := t.TempDir()
	file := filepath.Join(out, "index.html")
	if err := os.WriteFile(file, []byte("<html></html>"), 0644); err != nil {
		t.Fatal(err)
	}
	sha, err := runnerFileSHA256(file)
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "uploaded")
	fakeBin := t.TempDir()
	script := fmt.Sprintf("#!/bin/sh\ncase \"$2\" in\n  sha256sum*) if [ -f %q ]; then printf '%s  file\\n'; else printf 'old  file\\n'; fi; exit 0 ;;\n  mkdir*) cat >/dev/null; touch %q; exit 0 ;;\nesac\nexit 1\n", marker, sha, marker)
	writeExecutable(t, filepath.Join(fakeBin, "ssh"), script)
	t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	total, uploaded, skipped, bytesUploaded, err := deploySite(out, DeployTarget{Host: "host", User: "deploy", DestDir: "/var/www/html"})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || uploaded != 1 || skipped != 0 || bytesUploaded != int64(len("<html></html>")) {
		t.Fatalf("unexpected deploy stats total=%d uploaded=%d skipped=%d bytes=%d", total, uploaded, skipped, bytesUploaded)
	}
}

func TestRunnerCompletionCommitInfoAndNotification(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	notified := 0
	hook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notified++
		w.WriteHeader(http.StatusNoContent)
	}))
	defer hook.Close()
	writeRunnerTestJSON(t, filepath.Join(state, ".notify_config"), notifyConfigFile{Webhooks: []notifyWebhook{{URL: hook.URL, Enabled: true, On: []string{"success"}}}})
	server := fakeGitHubWithCommit(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.Commit.SHA == nil || *log.Commit.SHA != "commit-sha" || log.Commit.Message == nil || *log.Commit.Message != "Update docs" {
			t.Fatalf("commit info not recorded: %+v", log.Commit)
		}
		if log.CommitSHA == nil || *log.CommitSHA != "commit-sha" || log.CommitMessage == nil || *log.CommitMessage != "Update docs" || log.CommitAuthor == nil || *log.CommitAuthor != "A. Developer" || log.CommitAt == nil {
			t.Fatalf("top-level commit fields not recorded: %+v", log)
		}
		if notified != 1 {
			t.Fatalf("expected one notification, got %d", notified)
		}
	})
}

func TestRunnerPhase2NotificationChannelsClassifyFailures(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://example.invalid/elsewhere", http.StatusFound)
	}))
	defer redirect.Close()
	failing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Adlaire-Signature-256") == "" {
			t.Fatalf("webhook signature missing")
		}
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer failing.Close()
	marker := filepath.Join(state, "command-notified")
	command := filepath.Join(t.TempDir(), "notify-command")
	writeExecutable(t, command, fmt.Sprintf("#!/usr/bin/env bash\ncat >/dev/null\nprintf ok > %q\nexit 0\n", marker))
	retryCount := 2
	retryInterval := 1
	writeRunnerTestJSON(t, filepath.Join(state, ".notify_config"), notifyConfigFile{
		Channels: []notifyChannel{
			{ID: "cmd", Type: "command", Enabled: true, On: []string{"success"}, Config: map[string]any{"command_args": []any{command}}},
			{ID: "mail", Type: "email", Enabled: true, On: []string{"success"}, Config: map[string]any{}},
			{ID: "redirect", Type: "webhook", Enabled: true, On: []string{"success"}, Config: map[string]any{"url": redirect.URL}, RetryCount: &retryCount, RetryIntervalSeconds: &retryInterval},
			{ID: "server", Type: "webhook", Enabled: true, On: []string{"success"}, Config: map[string]any{"url": failing.URL, "secret": "sign-me"}, RetryCount: &retryCount, RetryIntervalSeconds: &retryInterval},
		},
	})
	notifyConfig, err := readNotifyConfig(filepath.Join(state, ".notify_config"))
	if err != nil {
		t.Fatal(err)
	}
	if targets := notifyTargetsForEvent(notifyConfig, "success"); len(targets) != 4 {
		t.Fatalf("notify targets mismatch: %+v", targets)
	}
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, marker) != "ok" {
			t.Fatalf("command notification did not run")
		}
		var pending []notifyPendingEntry
		pendingPath := filepath.Join(state, ".notify_pending")
		if _, err := os.Stat(pendingPath); err != nil {
			t.Fatalf("notify pending missing: %v log=%s", err, readFile(t, filepath.Join(state, ".notify_log")))
		}
		readJSON(t, pendingPath, &pending)
		if len(pending) != 1 || pending[0].ChannelID != "server" || pending[0].Attempts != 1 || pending[0].PayloadSHA == "" {
			t.Fatalf("unexpected pending entries: %+v", pending)
		}
		logText := readFile(t, filepath.Join(state, ".notify_log"))
		for _, needle := range []string{`"channel_id":"cmd"`, `"result":"success"`, `"channel_id":"mail"`, `"error_code":"smtp_not_configured"`, `"channel_id":"redirect"`, `"error_code":"http_3xx"`, `"channel_id":"server"`, `"error_code":"http_5xx"`} {
			if !strings.Contains(logText, needle) {
				t.Fatalf("notify log missing %s: %s", needle, logText)
			}
		}
	})
}

func TestRunnerPhase2TagFilterBuildsMatchingTag(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"tag_filter": map[string]any{"enabled": true, "patterns": []string{"v*"}},
	})
	server := fakeGitHubWithTags(t, "docs", "new-blob", "# Title\n", []map[string]string{{"ref": "refs/tags/v1.0.0", "sha": "commit-sha"}})
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "success" || strings.Join(log.MatchedTags, ",") != "v1.0.0" {
			t.Fatalf("tag filter log mismatch: %+v", log)
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"new-blob\"}\n" {
			t.Fatalf("sha not updated after matching tag")
		}
	})
}

func TestRunnerPhase2TagFilterSkipsNonMatchingTag(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"tag_filter": map[string]any{"enabled": true, "patterns": []string{"release-*"}},
	})
	server := fakeGitHubWithTags(t, "docs", "new-blob", "# Title\n", []map[string]string{{"ref": "refs/tags/v1.0.0", "sha": "commit-sha"}})
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
			t.Fatalf("tag-filter skip must not create history")
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"old-blob\"}\n" {
			t.Fatalf("sha updated despite non-matching tag")
		}
		var summary buildStatusSummary
		readJSON(t, filepath.Join(state, ".build_status.json"), &summary)
		if summary.Status != "skipped" || summary.LastTargetStatus == nil || *summary.LastTargetStatus != "skipped_tag_filter" {
			t.Fatalf("tag-filter skip status mismatch: %+v", summary)
		}
	})
}

func TestRunnerPhase2RemoteBuild(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	workDir := t.TempDir()
	fakeBin := t.TempDir()
	writeExecutable(t, filepath.Join(fakeBin, "ssh"), "#!/usr/bin/env bash\nset -euo pipefail\neval \"$2\"\n")
	t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"remote_build": map[string]any{
			"enabled": true, "host": "host", "user": "deploy", "work_dir": workDir,
			"command_args": []string{"sh", "-c", "mkdir -p \"$ADLAIRE_CI_OUT/assets\" && printf '<html></html>' > \"$ADLAIRE_CI_OUT/index.html\" && printf 'body{}' > \"$ADLAIRE_CI_OUT/assets/style.css\" && printf 'console.log(\"ok\")' > \"$ADLAIRE_CI_OUT/assets/app.js\" && printf '[]' > \"$ADLAIRE_CI_OUT/assets/search-index.json\" && printf '[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=remote'"},
		},
	})
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if !log.Environment.RemoteBuild || log.Report == nil || log.Report.Theme != "remote" || log.TargetStatus != "success" {
			t.Fatalf("remote build log mismatch: %+v", log)
		}
		if log.RemoteBuild == nil || log.RemoteBuild.Status != "success" || log.RemoteBuild.CommandName != "sh" || log.RemoteBuild.WorkDirBasename != filepath.Base(workDir) || log.RemoteBuild.ExitCode == nil || *log.RemoteBuild.ExitCode != 0 {
			t.Fatalf("remote build object mismatch: %+v", log.RemoteBuild)
		}
	})
}

func TestRunnerPhase2WeeklySummaryChannel(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	received := 0
	hook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received++
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("summary payload decode: %v", err)
		}
		if payload["event"] != "weekly_summary" || payload["success_count"] == nil {
			t.Fatalf("summary payload mismatch: %+v", payload)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer hook.Close()
	writeRunnerTestJSON(t, filepath.Join(state, ".notify_config"), notifyConfigFile{
		Channels: []notifyChannel{{ID: "weekly", Type: "webhook", Enabled: true, On: []string{"weekly_summary"}, Config: map[string]any{"url": hook.URL}}},
		Summary:  notifySummary{Enabled: true, Interval: "weekly", Hour: 9, DayOfWeek: 0},
	})
	if err := appendRunnerJSONLine(filepath.Join(state, ".build_history"), historyRecord{ID: "b1", Branch: "main", TargetFile: "docs", Status: "success", FinishedAt: "2026-09-27T08:00:00Z", DurationSeconds: 4}, 0600); err != nil {
		t.Fatal(err)
	}
	oldNow := runnerNow
	runnerNow = func() time.Time { return time.Date(2026, 9, 27, 9, 0, 0, 0, time.UTC) }
	defer func() { runnerNow = oldNow }()
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if received != 1 {
			t.Fatalf("weekly summary not sent: %d", received)
		}
		var bs buildState
		readJSON(t, filepath.Join(state, ".build_state"), &bs)
		if bs.WeeklySummarySentDate == nil || *bs.WeeklySummarySentDate != "2026-09-27" {
			t.Fatalf("summary date not recorded: %+v", bs)
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".notify_log")), `"event":"weekly_summary"`) {
			t.Fatalf("notify log missing weekly summary")
		}
	})
}

func TestRunnerPhase2DurationAnomaly(t *testing.T) {
	state := t.TempDir()
	cfg := defaultRunnerConfig(state)
	cfg.DurationAnomaly = DurationAnomalyConfig{Enabled: true, MinSamples: 2, AvgMultiplier: 2, P95Multiplier: 2}
	writeRunnerTestJSON(t, filepath.Join(state, ".build_trends.json"), buildTrendsFile{SchemaVersion: 1, Samples: []buildTrendSample{
		{BuildID: "a", FinishedAt: "2026-09-16T00:00:00Z", DurationSeconds: 10, Status: "success", TargetStatus: "success"},
		{BuildID: "b", FinishedAt: "2026-09-16T00:01:00Z", DurationSeconds: 10, Status: "success", TargetStatus: "success"},
	}})
	if err := updateBuildTrends(cfg, buildLog{ID: "c", FinishedAt: "2026-09-16T00:02:00Z", DurationSeconds: 100, Status: "success", TargetStatus: "success"}); err != nil {
		t.Fatal(err)
	}
	var trends buildTrendsFile
	readJSON(t, filepath.Join(state, ".build_trends.json"), &trends)
	if len(trends.Samples) != 3 || !trends.Samples[2].Anomaly || trends.Summary.AnomalyCount != 1 {
		t.Fatalf("duration anomaly not recorded: %+v", trends)
	}
}

func TestRunnerPhase2PriorityQueueRunsUrgentBeforeNormal(t *testing.T) {
	state := newRunnerState(t, "same-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"schedule_paused": true,
	})
	writeRunnerTestJSON(t, filepath.Join(state, ".build_state"), buildState{Queued: []map[string]any{
		{"id": "q-normal", "trigger": "manual", "queued_at": "2026-09-16T00:00:00Z", "requested_by": "admin", "priority": "normal", "created_seq": 1, "payload": map[string]any{"force": true}},
		{"id": "q-urgent", "trigger": "manual", "queued_at": "2026-09-16T00:00:01Z", "requested_by": "admin", "priority": "urgent", "created_seq": 2, "payload": map[string]any{"force": true}},
	}})
	server := fakeGitHub(t, "docs", "same-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.Trigger != "manual" || log.TargetStatus != "success" {
			t.Fatalf("queue build log mismatch: %+v", log)
		}
		var bs buildState
		readJSON(t, filepath.Join(state, ".build_state"), &bs)
		if bs.ActiveQueueEntry != nil || len(bs.Queued) != 1 || fmt.Sprint(bs.Queued[0]["id"]) != "q-normal" {
			t.Fatalf("queue state mismatch: %+v", bs)
		}
	})
}

func TestRunnerPhase2BuildChainRequiredDependencySkipsDependent(t *testing.T) {
	state := newRunnerState(t, "old-a", nil)
	writePipeline(t, state, 7, "failed", "boom")
	targetA := BranchTarget{
		Branch: "main", TargetFile: "docs-a", SHAFile: filepath.Join(state, ".last_sha_a"),
		Src: filepath.Join(state, "repo", "docs-a"), Out: filepath.Join(state, "dist", "site-a"),
	}
	targetB := BranchTarget{
		Branch: "main", TargetFile: "docs-b", SHAFile: filepath.Join(state, ".last_sha_b"),
		Src: filepath.Join(state, "repo", "docs-b"), Out: filepath.Join(state, "dist", "site-b"),
	}
	if err := os.WriteFile(targetA.SHAFile, []byte("{\"sha\":\"old-a\"}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetB.SHAFile, []byte("{\"sha\":\"old-b\"}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".branch_config"), branchConfigFile{BranchTargets: []BranchTarget{targetA, targetB}})
	writeRunnerTestJSON(t, filepath.Join(state, ".build_chain_config"), buildChainConfigFile{Chains: []buildChainJob{
		{ID: "a", Branch: "main", TargetFile: "docs-a", DependsOn: []string{}, Required: true, Enabled: true},
		{ID: "b", Branch: "main", TargetFile: "docs-b", DependsOn: []string{"a"}, Required: true, Enabled: true},
	}})
	server := fakeGitHubMulti(t, map[string]struct {
		sha  string
		blob string
	}{
		"docs-a": {sha: "new-a", blob: "# A\n"},
		"docs-b": {sha: "new-b", blob: "# B\n"},
	})
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.Chain == nil || log.ChainSummary == nil || log.TargetStatus != "failure_build" {
			t.Fatalf("chain log mismatch: %+v", log)
		}
		summaryJSON, _ := json.Marshal(log.ChainSummary)
		if !strings.Contains(string(summaryJSON), `"failure_count":1`) || !strings.Contains(string(summaryJSON), `"skipped_count":1`) {
			t.Fatalf("chain summary mismatch: %s", summaryJSON)
		}
		history := readFile(t, filepath.Join(state, ".build_history"))
		if !strings.Contains(history, `"status":"failure_build"`) || !strings.Contains(history, `"status":"skipped_dependency_failed"`) || !strings.Contains(history, `"chain_run_id":"chain`) {
			t.Fatalf("chain history mismatch: %s", history)
		}
	})
}

func TestRunnerCompletionCommitStatus(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	writeRunnerTestJSON(t, filepath.Join(state, ".server_config"), map[string]any{
		"build_cooldown_seconds": 0,
		"commit_status_enabled":  true,
		"commit_status_context":  "Adlaire CI",
	})
	encoded := base64.StdEncoding.EncodeToString([]byte("# Title\n"))
	statuses := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": "docs", "type": "blob", "sha": "new-blob"}}})
		case strings.Contains(r.URL.Path, "/git/blobs/"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
		case strings.Contains(r.URL.Path, "/commits"):
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"sha": "abcdef0123456789abcdef0123456789abcdef01",
				"commit": map[string]any{
					"message": "Update docs",
					"author":  map[string]string{"name": "A. Developer", "date": "2026-09-16T00:00:00Z"},
				},
			}})
		case strings.Contains(r.URL.Path, "/statuses/"):
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("status payload decode: %v", err)
			}
			statuses = append(statuses, fmt.Sprint(payload["state"]))
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if got := strings.Join(statuses, ","); got != "pending,success" {
			t.Fatalf("unexpected commit statuses: %s", got)
		}
		log := onlyBuildLog(t, state)
		if log.CommitStatus == nil || *log.CommitStatus != "success" {
			t.Fatalf("commit status not recorded in log: %+v", log.CommitStatus)
		}
		if log.CommitStatusLog == nil || !log.CommitStatusLog.Enabled || log.CommitStatusLog.State == nil || *log.CommitStatusLog.State != "success" || log.CommitStatusLog.Context != "Adlaire CI" {
			t.Fatalf("commit status object mismatch: %+v", log.CommitStatusLog)
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".build_history")), `"commit_status_state":"success"`) {
			t.Fatalf("commit status not recorded in history")
		}
	})
}

func TestRunnerCompletionBuildRetryTimeoutSucceeds(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	counter := filepath.Join(state, "retry-count")
	buildBin := filepath.Join(t.TempDir(), "adlaire-ci-build")
	writeExecutable(t, buildBin, fmt.Sprintf(`#!/usr/bin/env bash
set -u
if [ "${1:-}" = "--version" ]; then
  printf 'adlaire-ci-build V.0.0-dev go=fake\n'
  exit 0
fi
out=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --out)
      shift
      out="${1:-}"
      ;;
  esac
  shift || true
done
if [ ! -f %q ]; then
  printf '1' > %q
  printf 'first timeout'
  sleep 2
  exit 0
fi
mkdir -p "$out/assets"
printf '<html></html>' > "$out/index.html"
printf 'body{}' > "$out/assets/style.css"
printf 'console.log("ok")' > "$out/assets/app.js"
printf '[]' > "$out/assets/search-index.json"
printf '[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default'
exit 0
`, counter, counter))
	t.Setenv("ADLAIRE_CI_BUILD_BIN", buildBin)
	oldSleep := runnerSleep
	runnerSleep = func(time.Duration) {}
	defer func() { runnerSleep = oldSleep }()
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		cfg := defaultRunnerConfig(state)
		cfg.BuildTimeoutSeconds = 1
		cfg.BuildRetryMax = 1
		cfg.BuildRetryBaseSeconds = 1
		cfg.BuildCooldownSeconds = 0
		var stdout, stderr bytes.Buffer
		if code := executeRunner(cfg, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.RetryCount != 1 || len(log.Attempts) != 2 || log.Attempts[0].ExitCode != nil || log.Attempts[1].ExitCode == nil || *log.Attempts[1].ExitCode != 0 || log.TargetStatus != "success" {
			t.Fatalf("retry log mismatch: %+v", log)
		}
		if !strings.Contains(readFile(t, filepath.Join(state, ".build_history")), `"retry_count":1`) {
			t.Fatalf("history missing retry count")
		}
	})
}

func TestRunnerCompletionBuildRetrySkipsPipelineExit(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipeline(t, state, 7, "first failure", "failed")
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		cfg := defaultRunnerConfig(state)
		cfg.BuildRetryMax = 1
		cfg.BuildRetryBaseSeconds = 1
		cfg.BuildCooldownSeconds = 0
		var stdout, stderr bytes.Buffer
		if code := executeRunner(cfg, &stdout, &stderr); code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.RetryCount != 0 || len(log.Attempts) != 1 || log.Attempts[0].ExitCode == nil || *log.Attempts[0].ExitCode != 7 || log.TargetStatus != "failure_build" {
			t.Fatalf("pipeline exit should not retry: %+v", log)
		}
	})
}

func TestRunnerCompletionCooldownSkips(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	now := time.Date(2026, 9, 16, 1, 2, 3, 0, time.UTC)
	nowText := now.Format(time.RFC3339)
	writeRunnerTestJSON(t, filepath.Join(state, ".build_state"), buildState{Queued: []map[string]any{}, LastFinishedAt: &nowText})
	oldNow := runnerNow
	runnerNow = func() time.Time { return now.Add(10 * time.Second) }
	defer func() { runnerNow = oldNow }()
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		http.NotFound(w, r)
	}))
	defer server.Close()
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if hits != 0 {
			t.Fatalf("cooldown should skip API calls, hits=%d", hits)
		}
	})
}

func TestRunnerCompletionForceIntervalBuildsSameSHA(t *testing.T) {
	state := newRunnerState(t, "same-blob", nil)
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	old := time.Date(2026, 9, 15, 1, 0, 0, 0, time.UTC).Format(time.RFC3339)
	writeRunnerTestJSON(t, filepath.Join(state, ".build_state"), buildState{Queued: []map[string]any{}, LastFinishedAt: &old})
	server := fakeGitHub(t, "docs", "same-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		cfg := defaultRunnerConfig(state)
		cfg.ForceBuildIntervalHours = 1
		cfg.BuildCooldownSeconds = 0
		var stdout, stderr bytes.Buffer
		if code := executeRunner(cfg, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".build_history")); err != nil {
			t.Fatalf("force interval did not build: %v", err)
		}
	})
}

func TestRunnerCompletionSizeWarning(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipelineWithExtra(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "", `dd if=/dev/zero of="$out/big.bin" bs=1048576 count=2 2>/dev/null`)
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		cfg := defaultRunnerConfig(state)
		cfg.OutputSizeWarnMB = 1
		cfg.BuildCooldownSeconds = 0
		var stdout, stderr bytes.Buffer
		if code := executeRunner(cfg, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.Report == nil || !log.Report.SizeWarn {
			t.Fatalf("size warn should be true when threshold exceeded: %+v", log.Report)
		}
		if len(log.Warnings) != 1 || !strings.Contains(log.Warnings[0], "OUTPUT_SIZE_WARN") {
			t.Fatalf("missing size warning: %+v", log.Warnings)
		}
	})
}

func TestRunnerCompletionSizeWarningWithoutReport(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	writePipelineWithExtra(t, state, 0, "", "", `dd if=/dev/zero of="$out/big.bin" bs=1048576 count=2 2>/dev/null`)
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		cfg := defaultRunnerConfig(state)
		cfg.OutputSizeWarnMB = 1
		cfg.BuildCooldownSeconds = 0
		var stdout, stderr bytes.Buffer
		if code := executeRunner(cfg, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		log := onlyBuildLog(t, state)
		if log.Report != nil {
			t.Fatalf("report must be absent: %+v", log.Report)
		}
		if !log.SizeWarn {
			t.Fatalf("size warn should be independent of report: %+v", log)
		}
		if !strings.Contains(strings.Join(log.Warnings, ","), "REPORT_MISSING") || !strings.Contains(strings.Join(log.Warnings, ","), "OUTPUT_SIZE_WARN") {
			t.Fatalf("missing warnings: %+v", log.Warnings)
		}
	})
}

func TestRunnerCompletionPATExpiryWarning(t *testing.T) {
	oldNow := runnerNow
	runnerNow = func() time.Time { return time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC) }
	defer func() { runnerNow = oldNow }()

	var out bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&out, &slog.HandlerOptions{Level: slog.LevelDebug}))
	resp := &http.Response{Header: make(http.Header)}
	resp.Header.Set("GitHub-Authentication-Token-Expiration", "2026-09-20T00:00:00Z")
	logPATExpiryWarning(resp, logger)
	if !strings.Contains(out.String(), "PAT_EXPIRY_WARN") || !strings.Contains(out.String(), "remaining_days=4") {
		t.Fatalf("missing PAT expiry warning: %s", out.String())
	}
}

func TestRunnerFixtureR5DeployPending(t *testing.T) {
	state := newRunnerState(t, "old-blob", []DeployTarget{{Host: "example.invalid", User: "deploy", DestDir: "/var/www/html"}})
	writePipeline(t, state, 0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "")
	fakeBin := t.TempDir()
	writeExecutable(t, filepath.Join(fakeBin, "ssh"), "#!/bin/sh\nexit 255\n")
	t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	server := fakeGitHub(t, "docs", "new-blob", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr)
		if code != 1 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if readFile(t, filepath.Join(state, ".last_sha")) != "{\"sha\":\"new-blob\"}\n" {
			t.Fatalf("sha not updated before deploy pending")
		}
		var pending []pendingTransfer
		readJSON(t, filepath.Join(state, ".pending_transfers"), &pending)
		if len(pending) != 1 || pending[0].RetryCount != 0 || pending[0].DestDir != "/var/www/html" || pending[0].TargetID != "0-0" || pending[0].Out == nil {
			t.Fatalf("unexpected pending: %+v", pending)
		}
		log := onlyBuildLog(t, state)
		if log.TargetStatus != "success_deploy_pending" || len(log.Deploy) != 1 || log.Deploy[0].Status != "pending" || log.Deploy[0].TransferVerified || log.Error == nil || *log.Error != "deploy pending" {
			t.Fatalf("unexpected deploy log: %+v", log)
		}
		if log.TransferVerified == nil || *log.TransferVerified {
			t.Fatalf("transfer_verified summary mismatch: %+v", log.TransferVerified)
		}
		if len(log.TargetResults) != 1 || log.TargetResults[0].TargetID != "0-0" || log.TargetResults[0].Status != "failure" || log.TargetResults[0].ErrorCode == nil || *log.TargetResults[0].ErrorCode != "deploy_ssh_error" || log.TargetResults[0].Error == nil || *log.TargetResults[0].Error != "deploy ssh failed" {
			t.Fatalf("unexpected target results: %+v", log.TargetResults)
		}
		if _, err := time.Parse(time.RFC3339, log.TargetResults[0].StartedAt); err != nil {
			t.Fatalf("target result started_at invalid: %+v", log.TargetResults[0])
		}
		if _, err := time.Parse(time.RFC3339, log.TargetResults[0].FinishedAt); err != nil {
			t.Fatalf("target result finished_at invalid: %+v", log.TargetResults[0])
		}
	})
}

func TestRunnerFixtureR6LockConflict(t *testing.T) {
	state := newRunnerState(t, "old-blob", nil)
	lock := fmt.Sprintf("pid=%d\nstarted_at=2026-09-16T00:00:00Z\n", os.Getpid())
	if err := os.WriteFile(filepath.Join(state, ".build_lock"), []byte(lock), 0600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d", code)
	}
	if _, err := os.Stat(filepath.Join(state, ".build_history")); !os.IsNotExist(err) {
		t.Fatalf("history must not be touched")
	}
	if !strings.Contains(stdout.String(), "BUILD_SKIP: already running") {
		t.Fatalf("missing lock log: %s", stdout.String())
	}
}

func TestRunnerFixtureR7CorruptNotifyPending(t *testing.T) {
	state := newRunnerState(t, "blob-1", nil)
	if err := os.WriteFile(filepath.Join(state, ".notify_pending"), []byte("{bad json"), 0600); err != nil {
		t.Fatal(err)
	}
	fixedNow := time.Date(2026, 9, 16, 1, 2, 3, 0, time.UTC)
	oldNow := runnerNow
	runnerNow = func() time.Time { return fixedNow }
	defer func() { runnerNow = oldNow }()
	server := fakeGitHub(t, "docs", "blob-1", "# Title\n")
	withRunnerServer(t, server.URL, func() {
		var stdout, stderr bytes.Buffer
		if code := RunRunner([]string{"--state-dir", state}, &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(filepath.Join(state, ".notify_pending.corrupt.20260916010203.bak")); err != nil {
			t.Fatalf("missing corrupt backup: %v", err)
		}
		if strings.TrimSpace(readFile(t, filepath.Join(state, ".notify_pending"))) != "[]" {
			t.Fatalf("notify pending not reset")
		}
	})
}

func newRunnerState(t *testing.T, sha string, deploy []DeployTarget) string {
	t.Helper()
	state := t.TempDir()
	buildBin := filepath.Join(t.TempDir(), "adlaire-ci-build")
	writeExecutable(t, buildBin, fakeBuilderScript(0, `[REPORT] pages=1 headings=1 tables=0 code_blocks=0 warnings=0 size_warn=false broken_links=0 heading_skips=0 reading_time=1 theme=adlaire-default`, "", ""))
	t.Setenv("ADLAIRE_CI_BUILD_BIN", buildBin)
	if err := os.WriteFile(filepath.Join(state, ".github_token"), []byte("token\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(state, ".last_sha"), []byte(fmt.Sprintf("{\"sha\":\"%s\"}\n", sha)), 0600); err != nil {
		t.Fatal(err)
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".pipeline_config"), pipelineConfigFile{ExtraArgs: []string{}, Env: map[string]string{}})
	target := BranchTarget{
		Branch: "main", TargetFile: "docs", SHAFile: filepath.Join(state, ".last_sha"),
		Src: filepath.Join(state, "repo", "docs"), Out: filepath.Join(state, "dist", "site"),
		DeployTargets: deploy,
	}
	writeRunnerTestJSON(t, filepath.Join(state, ".branch_config"), branchConfigFile{BranchTargets: []BranchTarget{target}})
	return state
}

func writePipeline(t *testing.T, state string, code int, stdout string, stderr string) {
	t.Helper()
	writePipelineWithExtra(t, state, code, stdout, stderr, "")
}

func writePipelineWithExtra(t *testing.T, state string, code int, stdout string, stderr string, extra string) {
	t.Helper()
	buildBin := filepath.Join(t.TempDir(), "adlaire-ci-build")
	writeExecutable(t, buildBin, fakeBuilderScript(code, stdout, stderr, extra))
	t.Setenv("ADLAIRE_CI_BUILD_BIN", buildBin)
}

func fakeBuilderScript(code int, stdout string, stderr string, extra string) string {
	return fmt.Sprintf(`#!/usr/bin/env bash
set -u
if [ "${1:-}" = "--version" ]; then
  printf 'adlaire-ci-build V.0.0-dev go=fake\n'
  exit 0
fi
out=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --out)
      shift
      out="${1:-}"
      ;;
  esac
  shift || true
done
printf '%%s' %q
printf '%%s' %q >&2
if [ %d -eq 0 ]; then
  mkdir -p "$out/assets"
  printf '<html></html>' > "$out/index.html"
  printf 'body{}' > "$out/assets/style.css"
  printf 'console.log("ok")' > "$out/assets/app.js"
  printf '[]' > "$out/assets/search-index.json"
  %s
fi
exit %d
`, stdout, stderr, code, extra, code)
}

func writeExecutable(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
}

func fakeGitHub(t *testing.T, target, sha, blob string) *httptest.Server {
	t.Helper()
	encoded := base64.StdEncoding.EncodeToString([]byte(blob))
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": target, "type": "blob", "sha": sha}}})
		case strings.Contains(r.URL.Path, "/git/blobs/"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
		default:
			http.NotFound(w, r)
		}
	}))
}

func fakeGitHubMulti(t *testing.T, blobs map[string]struct {
	sha  string
	blob string
}) *httptest.Server {
	t.Helper()
	encodedBySHA := map[string]string{}
	tree := []map[string]string{}
	for path, blob := range blobs {
		encodedBySHA[blob.sha] = base64.StdEncoding.EncodeToString([]byte(blob.blob))
		tree = append(tree, map[string]string{"path": path, "type": "blob", "sha": blob.sha})
	}
	sort.Slice(tree, func(i, j int) bool { return tree[i]["path"] < tree[j]["path"] })
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": tree})
		case strings.Contains(r.URL.Path, "/git/blobs/"):
			sha := path.Base(r.URL.Path)
			encoded, ok := encodedBySHA[sha]
			if !ok {
				http.NotFound(w, r)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
		default:
			http.NotFound(w, r)
		}
	}))
}

func fakeGitHubWithCommit(t *testing.T, target, sha, blob string) *httptest.Server {
	t.Helper()
	encoded := base64.StdEncoding.EncodeToString([]byte(blob))
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": target, "type": "blob", "sha": sha}}})
		case strings.Contains(r.URL.Path, "/git/blobs/"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
		case strings.Contains(r.URL.Path, "/commits"):
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"sha": "commit-sha",
				"commit": map[string]any{
					"message": "Update docs\n\nBody",
					"author":  map[string]string{"name": "A. Developer", "date": "2026-09-16T00:00:00Z"},
				},
			}})
		default:
			http.NotFound(w, r)
		}
	}))
}

func fakeGitHubWithTags(t *testing.T, target, sha, blob string, refs []map[string]string) *httptest.Server {
	t.Helper()
	encoded := base64.StdEncoding.EncodeToString([]byte(blob))
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/trees/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"tree": []map[string]string{{"path": target, "type": "blob", "sha": sha}}})
		case strings.Contains(r.URL.Path, "/git/blobs/"):
			_ = json.NewEncoder(w).Encode(map[string]string{"content": encoded, "encoding": "base64"})
		case strings.Contains(r.URL.Path, "/commits"):
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"sha": "commit-sha",
				"commit": map[string]any{
					"message": "Release docs",
					"author":  map[string]string{"name": "A. Developer", "date": "2026-09-16T00:00:00Z"},
				},
			}})
		case strings.Contains(r.URL.Path, "/git/matching-refs/tags"):
			out := []map[string]any{}
			for _, ref := range refs {
				out = append(out, map[string]any{"ref": ref["ref"], "object": map[string]string{"sha": ref["sha"], "type": "commit"}})
			}
			_ = json.NewEncoder(w).Encode(out)
		default:
			http.NotFound(w, r)
		}
	}))
}

func withRunnerServer(t *testing.T, url string, fn func()) {
	t.Helper()
	old := runnerGitHubAPIBase
	runnerGitHubAPIBase = url
	t.Setenv("ADLAIRE_CI_REPOSITORY", "owner/repo")
	defer func() { runnerGitHubAPIBase = old }()
	fn()
}

func onlyBuildLog(t *testing.T, state string) buildLog {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(state, ".build_logs"))
	if err != nil {
		t.Fatal(err)
	}
	mainLogs := []os.DirEntry{}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") && !strings.Contains(entry.Name(), "_hook_") {
			mainLogs = append(mainLogs, entry)
		}
	}
	if len(mainLogs) != 1 {
		t.Fatalf("expected one build log, got %d", len(mainLogs))
	}
	var log buildLog
	readJSON(t, filepath.Join(state, ".build_logs", mainLogs[0].Name()), &log)
	return log
}

func readHookLog(t *testing.T, state, buildID, hookID string) hookRunLog {
	t.Helper()
	var log hookRunLog
	readJSON(t, filepath.Join(state, ".build_logs", buildID+"_hook_"+hookID+".json"), &log)
	return log
}

func runnerTestContainsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func writeRunnerTestJSON(t *testing.T, path string, v any) {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
}

func readJSON(t *testing.T, path string, v any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatal(err)
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

func tarGzContains(t *testing.T, archivePath, name string) bool {
	t.Helper()
	f, err := os.Open(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return false
		}
		if err != nil {
			t.Fatal(err)
		}
		if header.Name == name {
			return true
		}
	}
}
