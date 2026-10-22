package runner

import (
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fqwink/build-scripts/components/statefile"
)

func FetchCICDSource(ctx context.Context, request CICDSourceRequest, token string) (*CICDSourceSnapshot, error) {
	transport := &http.Transport{
		DialContext:         (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second, ResponseHeaderTimeout: 10 * time.Second,
		DisableCompression: true, MaxIdleConnsPerHost: 1,
	}
	defer transport.CloseIdleConnections()
	return cicdFetchSource(ctx, request, token, cicdSourceClient(transport))
}

func cicdSourceClient(transport http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: transport, Timeout: 10 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
}

func cicdFetchSource(ctx context.Context, request CICDSourceRequest, token string, client *http.Client) (*CICDSourceSnapshot, error) {
	if ctx == nil || !cicdSourceRequestValid(request, token) {
		return nil, ErrCICDSchema
	}
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	base := "/repos/" + request.Owner + "/" + request.Repo + "/git/"
	tree, err := cicdSourceGet(ctx, client, base+"trees/"+request.SHA+"?recursive=1", token, 8388608)
	if err != nil {
		return nil, err
	}
	files, err := cicdSourceTree(tree)
	if err != nil {
		return nil, err
	}
	snapshot := &CICDSourceSnapshot{files: make(map[string][]byte, len(files))}
	for _, file := range files {
		blob, err := cicdSourceGet(ctx, client, base+"blobs/"+file["blob_sha"].(string), token, 16777216)
		if err != nil {
			return nil, err
		}
		data, err := cicdSourceBlob(file, blob)
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256(data)
		file["sha256"] = hex.EncodeToString(sum[:])
		snapshot.files[file["path"].(string)] = data
	}
	digest, err := statefile.DigestCICDJSON(cicdPlanBytes(files))
	if err != nil {
		return nil, ErrCICDSource
	}
	metadata := map[string]any{
		"owner": request.Owner, "repo": request.Repo, "ref_kind": request.RefKind,
		"ref": request.Ref, "sha": request.SHA, "tree_digest": digest, "files": files,
	}
	snapshot.metadata, err = statefile.CanonicalCICDJSON(cicdPlanBytes(metadata))
	if err != nil || statefile.ValidateCICDValue("Source", snapshot.metadata) != nil {
		return nil, ErrCICDSource
	}
	if ctx.Err() != nil {
		return nil, ErrCICDTransport
	}
	return snapshot, nil
}

func cicdSourceGet(ctx context.Context, client *http.Client, path, token string, limit int64) (map[string]any, error) {
	if ctx.Err() != nil {
		return nil, ErrCICDTransport
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com"+path, nil)
	if err != nil {
		return nil, ErrCICDTransport
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept-Encoding", "identity")
	response, err := client.Do(request)
	if err != nil {
		return nil, ErrCICDTransport
	}
	defer response.Body.Close()
	media, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if response.StatusCode != http.StatusOK || err != nil ||
		(media != "application/json" && media != "application/vnd.github+json") || response.Header.Get("Content-Encoding") != "" {
		return nil, ErrCICDTransport
	}
	if response.ContentLength > limit {
		return nil, ErrCICDBody
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return nil, ErrCICDTransport
	}
	if int64(len(data)) > limit {
		return nil, ErrCICDBody
	}
	canonical, err := statefile.CanonicalCICDJSON(data)
	if err != nil {
		return nil, ErrCICDTransport
	}
	decoder := json.NewDecoder(bytes.NewReader(canonical))
	decoder.UseNumber()
	var result map[string]any
	if decoder.Decode(&result) != nil || result == nil {
		return nil, ErrCICDTransport
	}
	return result, nil
}

func cicdSourceBlob(file, blob map[string]any) ([]byte, error) {
	content, contentOK := blob["content"].(string)
	if blob["sha"] != file["blob_sha"] || blob["size"] != file["size"] || blob["encoding"] != "base64" || !contentOK {
		return nil, ErrCICDDigest
	}
	// GitHub folds base64; all other whitespace and noncanonical pad bits are rejected.
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\n", "")
	size := cicdPlanInt(file["size"])
	if len(content) != base64.StdEncoding.EncodedLen(int(size)) {
		return nil, ErrCICDDigest
	}
	data, err := base64.StdEncoding.Strict().DecodeString(content)
	if err != nil || int64(len(data)) != size {
		return nil, ErrCICDDigest
	}
	hash := sha1.New()
	hash.Write([]byte("blob " + strconv.FormatInt(size, 10) + "\x00"))
	hash.Write(data)
	if hex.EncodeToString(hash.Sum(nil)) != file["blob_sha"] {
		return nil, ErrCICDDigest
	}
	return data, nil
}

func (snapshot *CICDSourceSnapshot) Metadata() json.RawMessage {
	if snapshot == nil {
		return nil
	}
	return append(json.RawMessage(nil), snapshot.metadata...)
}

func (snapshot *CICDSourceSnapshot) ReadFile(path string) ([]byte, error) {
	if snapshot == nil || snapshot.metadata == nil || snapshot.files == nil {
		return nil, ErrCICDSource
	}
	data, exists := snapshot.files[path]
	if !exists {
		return nil, ErrCICDNotFound
	}
	return append([]byte{}, data...), nil
}

func PrepareCICDPlan(project, definition, values json.RawMessage, templates []json.RawMessage) (*CICDPlanDraft, error) {
	size := int64(len(project)) + int64(len(definition)) + int64(len(values))
	for _, raw := range templates {
		size += int64(len(raw))
		if size > 33554432 {
			return nil, ErrCICDBody
		}
	}
	if size > 33554432 {
		return nil, ErrCICDBody
	}
	p, err := cicdPlanObject("Project", project)
	if err != nil {
		return nil, err
	}
	d, err := cicdPlanObject("Definition", definition)
	if err != nil || d["project_id"] != p["id"] || len(templates) > 1024 {
		return nil, ErrCICDSchema
	}
	automatic := len(d["watches"].([]any)) != 0
	for _, field := range []string{"watches", "schedules"} {
		seen := make(map[string]bool)
		for _, value := range d[field].([]any) {
			trigger := value.(map[string]any)
			id := trigger["id"].(string)
			if seen[id] {
				return nil, ErrCICDSchema
			}
			seen[id] = true
			automatic = automatic || field == "schedules" && trigger["enabled"] == true
		}
	}
	if automatic {
		for _, declaration := range d["inputs"].(map[string]any) {
			if declaration.(map[string]any)["required"] == true {
				return nil, ErrCICDSchema
			}
		}
	}
	v, err := cicdPlanObject("ScalarMap", values)
	if err != nil {
		return nil, err
	}
	v, err = cicdNormalizeInputs(d["inputs"].(map[string]any), v)
	if err != nil {
		return nil, err
	}
	index := make(map[string]map[string]any)
	for _, data := range templates {
		template, err := cicdPlanObject("Template", data)
		if err != nil || template["project_id"] != p["id"] {
			return nil, ErrCICDSchema
		}
		key := template["id"].(string) + "/" + string(template["revision"].(json.Number))
		if index[key] != nil {
			return nil, ErrCICDSchema
		}
		index[key] = template
	}
	draft := &CICDPlanDraft{project: p, definition: d, values: v}
	baseJobs := make([]cicdDraftJob, 0)
	for _, value := range d["templates"].([]any) {
		use := value.(map[string]any)
		template := index[use["id"].(string)+"/"+string(use["revision"].(json.Number))]
		if template == nil {
			return nil, ErrCICDSchema
		}
		for _, scalar := range use["values"].(map[string]any) {
			if text, ok := scalar.(string); ok && strings.Contains(text, "${") {
				return nil, ErrCICDSchema
			}
		}
		inputs, err := cicdNormalizeInputs(template["inputs"].(map[string]any), use["values"].(map[string]any))
		if err != nil {
			return nil, err
		}
		prefix := use["prefix"].(string) + "_"
		local := make(map[string]bool)
		for _, item := range template["jobs"].([]any) {
			id := item.(map[string]any)["id"].(string)
			if local[id] {
				return nil, ErrCICDDAG
			}
			local[id] = true
		}
		for _, item := range template["jobs"].([]any) {
			job := cicdPlanClone(item.(map[string]any))
			job["id"] = prefix + job["id"].(string)
			if statefile.ValidateCICDValue("JobId", cicdPlanBytes(job["id"])) != nil {
				return nil, ErrCICDSchema
			}
			for i, need := range job["needs"].([]any) {
				if !local[need.(string)] {
					return nil, ErrCICDDAG
				}
				job["needs"].([]any)[i] = prefix + need.(string)
			}
			task := job["task"].(map[string]any)
			if task["kind"] == "deploy" {
				if !local[task["artifact_from"].(string)] {
					return nil, ErrCICDSchema
				}
				task["artifact_from"] = prefix + task["artifact_from"].(string)
			}
			baseJobs = append(baseJobs, cicdDraftJob{job: job, values: inputs})
		}
	}
	for _, item := range d["jobs"].([]any) {
		baseJobs = append(baseJobs, cicdDraftJob{job: item.(map[string]any), values: v})
	}
	baseIndex := make(map[string][]int)
	for _, base := range baseJobs {
		id := base.job["id"].(string)
		if _, exists := baseIndex[id]; exists {
			return nil, ErrCICDDAG
		}
		baseIndex[id] = nil
	}
	ids := make(map[string]bool)
	for _, base := range baseJobs {
		id := base.job["id"].(string)
		matrix := base.job["matrix"].(map[string]any)
		tuples, err := cicdMatrix(matrix)
		if err != nil || len(draft.jobs)+len(tuples) > 256 {
			return nil, ErrCICDMatrix
		}
		for n, tuple := range tuples {
			job := cicdPlanClone(base.job)
			job["id"], job["base_job_id"], job["selected"], job["matrix_values"] = id, id, false, tuple
			if len(matrix["axes"].(map[string]any)) != 0 {
				job["id"] = id + "_m" + strconv.Itoa(n+1)
			}
			expanded := job["id"].(string)
			if len(expanded) > 64 {
				return nil, ErrCICDMatrix
			}
			if ids[expanded] {
				return nil, ErrCICDDAG
			}
			ids[expanded] = true
			for _, field := range []string{"task", "env"} {
				job[field], err = cicdSubstitute(job[field], base.values, tuple)
				if err != nil {
					return nil, err
				}
			}
			cache := job["cache"].(map[string]any)
			cache["paths"], err = cicdSubstitute(cache["paths"], base.values, tuple)
			if err != nil {
				return nil, err
			}
			condition := job["condition"].(map[string]any)
			condition["input_equals"], err = cicdSubstitute(condition["input_equals"], base.values, tuple)
			if err != nil || statefile.ValidateCICDValue("Condition", cicdPlanBytes(condition)) != nil {
				return nil, ErrCICDSchema
			}
			for key := range condition["input_equals"].(map[string]any) {
				if _, exists := base.values[key]; !exists {
					return nil, ErrCICDSchema
				}
			}
			for _, field := range []string{"branches", "tags", "paths"} {
				for _, pattern := range condition[field].([]any) {
					if _, err := cicdMatchGlob(pattern.(string), ""); err != nil {
						return nil, err
					}
				}
			}
			delete(job, "condition")
			delete(job, "matrix")
			baseIndex[id] = append(baseIndex[id], len(draft.jobs))
			draft.jobs = append(draft.jobs, cicdDraftJob{job: job, condition: condition, values: base.values})
		}
	}
	for i := range draft.jobs {
		job := draft.jobs[i].job
		seen := make(map[string]bool)
		needs := make([]any, 0)
		for _, value := range job["needs"].([]any) {
			base := value.(string)
			if seen[base] || base == job["base_job_id"] {
				return nil, ErrCICDDAG
			}
			seen[base] = true
			if _, exists := baseIndex[base]; !exists {
				return nil, ErrCICDDAG
			}
		}
		// Dependencies are serialized in origin/matrix order, not request order.
		for _, candidate := range draft.jobs {
			if seen[candidate.job["base_job_id"].(string)] {
				needs = append(needs, candidate.job["id"])
			}
		}
		job["needs"] = needs
		task := job["task"].(map[string]any)
		if task["kind"] == "deploy" {
			base := task["artifact_from"].(string)
			indices := baseIndex[base]
			if !seen[base] || len(indices) != 1 || draft.jobs[indices[0]].job["task"].(map[string]any)["kind"] != "builder" {
				return nil, ErrCICDSchema
			}
			task["artifact_from"] = draft.jobs[indices[0]].job["id"]
		}
		if statefile.ValidateCICDValue("PlannedJob", cicdPlanBytes(job)) != nil || !cicdPlanCacheValid(job) || !cicdPlanTaskValid(job["task"].(map[string]any)) {
			return nil, ErrCICDSchema
		}
	}
	ordered := make([]cicdDraftJob, 0, len(draft.jobs))
	placed := make(map[string]bool)
	for len(ordered) < len(draft.jobs) {
		found := false
		for _, candidate := range draft.jobs {
			id := candidate.job["id"].(string)
			if placed[id] {
				continue
			}
			ready := true
			for _, need := range candidate.job["needs"].([]any) {
				if !placed[need.(string)] {
					ready = false
					break
				}
			}
			if ready {
				ordered = append(ordered, candidate)
				placed[id], found = true, true
				break
			}
		}
		if !found {
			return nil, ErrCICDDAG
		}
	}
	draft.jobs = ordered
	return draft, nil
}

func cicdPlanCacheValid(job map[string]any) bool {
	cache := job["cache"].(map[string]any)
	paths := cache["paths"].([]any)
	if cache["enabled"] == true && (len(paths) == 0 || len(job["secrets"].(map[string]any)) != 0) {
		return false
	}
	for i, path := range paths {
		text := path.(string)
		if strings.ContainsAny(text, "*?[]{}") {
			return false
		}
		for _, previous := range paths[:i] {
			other := previous.(string)
			if text == other || strings.HasPrefix(text, other+"/") || strings.HasPrefix(other, text+"/") {
				return false
			}
		}
	}
	return true
}

func cicdPlanTaskValid(task map[string]any) bool {
	var arguments []any
	switch task["kind"] {
	case "go-test", "go-vet":
		arguments = task["packages"].([]any)
	case "go-build":
		arguments = []any{task["package"]}
	case "deno-check":
		arguments = task["files"].([]any)
	}
	for _, argument := range arguments {
		path := argument.(string)
		if strings.HasPrefix(path, "-") || task["kind"] == "deno-check" && !strings.HasSuffix(path, ".js") {
			return false
		}
	}
	return true
}

func cicdPlanSource(data json.RawMessage) (map[string]any, error) {
	source, err := cicdPlanObject("Source", data)
	if err != nil {
		return nil, err
	}
	last, total := "", int64(0)
	for _, value := range source["files"].([]any) {
		file := value.(map[string]any)
		name := file["path"].(string)
		if last >= name {
			return nil, ErrCICDSchema
		}
		last = name
		total += cicdPlanInt(file["size"])
	}
	digest, err := statefile.DigestCICDJSON(cicdPlanBytes(source["files"]))
	if err != nil || total > 104857600 || digest != source["tree_digest"] {
		return nil, ErrCICDSchema
	}
	return source, nil
}

func CICDChangedPaths(previous, current json.RawMessage) ([]string, error) {
	next, err := cicdPlanSource(current)
	if err != nil {
		return nil, err
	}
	before := make(map[string]any)
	if previous != nil {
		prior, err := cicdPlanSource(previous)
		if err != nil {
			return nil, err
		}
		for _, field := range []string{"owner", "repo", "ref_kind", "ref"} {
			if prior[field] != next[field] {
				return nil, ErrCICDSchema
			}
		}
		for _, item := range prior["files"].([]any) {
			file := item.(map[string]any)
			before[file["path"].(string)] = file["sha256"]
		}
	}
	changed := make(map[string]bool)
	for _, item := range next["files"].([]any) {
		file := item.(map[string]any)
		path := file["path"].(string)
		if before[path] != file["sha256"] {
			changed[path] = true
		}
		delete(before, path)
	}
	for path := range before {
		changed[path] = true
	}
	paths := make([]string, 0, len(changed))
	for path := range changed {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

func (draft *CICDPlanDraft) Resolve(source json.RawMessage, changedPaths []string, context CICDPlanContext) (json.RawMessage, error) {
	if draft == nil || draft.project == nil || draft.definition == nil || len(draft.jobs) == 0 {
		return nil, ErrCICDSchema
	}
	if draft.project["enabled"] != true {
		return nil, ErrCICDForbidden
	}
	if draft.project["toolchain_digest"] == nil {
		return nil, ErrCICDCapability
	}
	s, err := cicdPlanSource(source)
	if err != nil || s["owner"] != draft.project["source_owner"] || s["repo"] != draft.project["source_repo"] {
		return nil, ErrCICDSchema
	}
	if changedPaths == nil {
		changedPaths = make([]string, 0, len(s["files"].([]any)))
		for _, value := range s["files"].([]any) {
			changedPaths = append(changedPaths, value.(map[string]any)["path"].(string))
		}
	}
	if len(changedPaths) > 20000 {
		return nil, ErrCICDSchema
	}
	for i, path := range changedPaths {
		if statefile.ValidateCICDValue("RelPath", cicdPlanBytes(path)) != nil || i > 0 && changedPaths[i-1] >= path {
			return nil, ErrCICDSchema
		}
	}
	projections := make(map[string][]map[string]any)
	var projectionBytes int64
	for _, group := range []struct {
		kind string
		list []json.RawMessage
		max  int
	}{{"Worker", context.Workers, 32}, {"SecretMeta", context.Secrets, 10000}, {"Artifact", context.Artifacts, 10000}, {"Environment", context.Environments, 1024}} {
		if len(group.list) > group.max {
			return nil, ErrCICDSchema
		}
		seen := make(map[string]bool)
		for _, raw := range group.list {
			projectionBytes += int64(len(raw))
			if projectionBytes > 33554432 {
				return nil, ErrCICDBody
			}
			object, err := cicdPlanObject(group.kind, raw)
			if err != nil {
				return nil, err
			}
			key := cicdProjectionKey(group.kind, object)
			if seen[key] {
				return nil, ErrCICDSchema
			}
			seen[key] = true
			projections[group.kind] = append(projections[group.kind], object)
		}
	}
	jobs := make([]any, 0, len(draft.jobs))
	selected := make(map[string]bool)
	for _, candidate := range draft.jobs {
		job := cicdPlanClone(candidate.job)
		match, err := cicdPlanCondition(candidate.condition, candidate.values, s, changedPaths)
		if err != nil {
			return nil, err
		}
		for _, need := range job["needs"].([]any) {
			match = match && selected[need.(string)]
		}
		job["selected"] = match
		selected[job["id"].(string)] = match
		if match {
			if err := cicdPlanPrerequisites(draft.project, job, projections); err != nil {
				return nil, err
			}
		}
		jobs = append(jobs, job)
	}
	d := draft.definition
	data := map[string]any{"project_id": draft.project["id"], "definition_id": d["id"], "definition_revision": d["revision"], "source": s, "values": draft.values, "jobs": jobs, "fail_fast": d["fail_fast"], "toolchain_digest": draft.project["toolchain_digest"]}
	canonical, err := statefile.CanonicalCICDJSON(cicdPlanBytes(data))
	if err != nil {
		return nil, ErrCICDSchema
	}
	if len(canonical) > 524288 {
		return nil, ErrCICDBody
	}
	digest, _ := statefile.DigestCICDJSON(canonical)
	plan, err := statefile.CanonicalCICDJSON(cicdPlanBytes(map[string]any{"digest": digest, "data": data}))
	if err != nil || statefile.ValidateCICDValue("Plan", plan) != nil {
		return nil, ErrCICDSchema
	}
	return plan, nil
}

func cicdProjectionKey(kind string, object map[string]any) string {
	project, _ := object["project_id"].(string)
	id, _ := object["id"].(string)
	if kind == "Artifact" {
		id = object["digest"].(string)
	}
	if kind == "SecretMeta" {
		id += "/" + string(object["version"].(json.Number))
	}
	return project + "/" + id
}

func cicdPlanCondition(condition, values, source map[string]any, paths []string) (bool, error) {
	kind, opposite := "branches", "tags"
	if source["ref_kind"] == "tag" {
		kind, opposite = opposite, kind
	}
	refMatch := len(condition[kind].([]any)) == 0 && len(condition[opposite].([]any)) == 0
	for _, pattern := range condition[kind].([]any) {
		match, err := cicdMatchGlob(pattern.(string), source["ref"].(string))
		if err != nil {
			return false, err
		}
		refMatch = refMatch || match
	}
	pathMatch := len(condition["paths"].([]any)) == 0
	for _, pattern := range condition["paths"].([]any) {
		for _, path := range paths {
			match, err := cicdMatchGlob(pattern.(string), path)
			if err != nil {
				return false, err
			}
			pathMatch = pathMatch || match
		}
	}
	inputMatch := true
	for key, expected := range condition["input_equals"].(map[string]any) {
		inputMatch = inputMatch && reflect.DeepEqual(values[key], expected)
	}
	return refMatch && pathMatch && inputMatch, nil
}

func cicdPlanPrerequisites(project, job map[string]any, projections map[string][]map[string]any) error {
	projectID := project["id"].(string)
	for _, value := range job["secrets"].(map[string]any) {
		ref := value.(map[string]any)
		found := false
		for _, secret := range projections["SecretMeta"] {
			found = found || secret["project_id"] == projectID && secret["id"] == ref["id"] && secret["version"] == ref["version"]
		}
		if !found {
			return ErrCICDSecret
		}
	}
	if job["cache"].(map[string]any)["enabled"] == true && cicdPlanInt(project["quota"].(map[string]any)["cache_mb"]) == 0 {
		return ErrCICDQuota
	}
	task := job["task"].(map[string]any)
	if task["kind"] == "artifact-verify" || task["kind"] == "deploy" {
		kind, field, wanted := "Artifact", "digest", task["artifact"]
		if task["kind"] == "deploy" {
			kind, field, wanted = "Environment", "id", task["environment_id"]
		}
		for _, object := range projections[kind] {
			if object["project_id"] == projectID && object[field] == wanted && object["deleting"] != true {
				return nil
			}
		}
		return ErrCICDNotFound
	}
	capability := task["kind"].(string)
	if strings.HasPrefix(capability, "go-") {
		capability = "go"
	} else if capability == "deno-check" {
		capability = "deno"
	}
	var resources CICDResources
	_ = json.Unmarshal(cicdPlanBytes(job["resources"]), &resources)
	for _, worker := range projections["Worker"] {
		if worker["enabled"] != true || worker["toolchain_digest"] != project["toolchain_digest"] {
			continue
		}
		allowed, capable := false, false
		for _, id := range worker["allowed_projects"].([]any) {
			allowed = allowed || id == projectID
		}
		for _, name := range worker["capabilities"].([]any) {
			capable = capable || name == capability
		}
		total := cicdPlanInt(worker["total_memory_mb"])
		if allowed && capable && CheckCICDMemory(CICDHostMemory{TotalMB: total, AvailableMB: total}, nil, cicdPlanInt(worker["slots"]), resources.MemoryMB) == nil {
			return nil
		}
	}
	return ErrCICDCapability
}

func DefaultCICDResources() CICDResources {
	return CICDResources{CPUPercent: 100, MemoryMB: 256, DiskMB: 128, Tasks: 64, TimeoutSeconds: 900, OutputMB: 16}
}

// ReadCICDHostMemory observes physical memory, not the complete sandbox capability.
func ReadCICDHostMemory() (CICDHostMemory, error) {
	if runtime.GOOS != "linux" {
		return CICDHostMemory{}, ErrCICDCapability
	}
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return CICDHostMemory{}, ErrCICDCapability
	}
	memory, readErr := ParseCICDHostMemory(file)
	closeErr := file.Close()
	if readErr != nil || closeErr != nil {
		return CICDHostMemory{}, ErrCICDCapability
	}
	return memory, nil
}

func ParseCICDHostMemory(reader io.Reader) (CICDHostMemory, error) {
	if reader == nil {
		return CICDHostMemory{}, ErrCICDCapability
	}
	data, err := io.ReadAll(io.LimitReader(reader, cicdMeminfoMaxSize+1))
	if err != nil || len(data) > cicdMeminfoMaxSize {
		return CICDHostMemory{}, ErrCICDCapability
	}
	values := make(map[string]int64, 2)
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.FieldsFunc(line, func(r rune) bool { return r == ' ' || r == '\t' })
		if len(fields) == 0 || (fields[0] != "MemTotal:" && fields[0] != "MemAvailable:") {
			continue
		}
		key := fields[0]
		if _, duplicate := values[key]; duplicate || len(fields) != 3 || fields[2] != "kB" {
			return CICDHostMemory{}, ErrCICDCapability
		}
		for _, digit := range fields[1] {
			if digit < '0' || digit > '9' {
				return CICDHostMemory{}, ErrCICDCapability
			}
		}
		value, parseErr := strconv.ParseInt(fields[1], 10, 64)
		if parseErr != nil || value < 0 || value > cicdMaximumHostMB*1024 {
			return CICDHostMemory{}, ErrCICDCapability
		}
		values[key] = value
	}
	total, hasTotal := values["MemTotal:"]
	available, hasAvailable := values["MemAvailable:"]
	if !hasTotal || !hasAvailable || total == 0 || available > total {
		return CICDHostMemory{}, ErrCICDCapability
	}
	return CICDHostMemory{TotalMB: total / 1024, AvailableMB: available / 1024}, nil
}

type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string { return e.Msg }

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}
