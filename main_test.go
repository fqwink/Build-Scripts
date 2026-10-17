package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"math/big"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestSpecificationReferenceIntegrity(t *testing.T) {
	t.Parallel()
	documents := map[string]string{}
	paths := []string{"AGENTS.md", "README.md", "CONTRIBUTING.md", "SECURITY.md"}
	if err := filepath.WalkDir("docs", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(path, ".md") {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		documents[filepath.ToSlash(path)] = phase11MustReadText(t, path)
	}
	exists := func(path string) bool {
		_, err := os.Lstat(filepath.FromSlash(path))
		return err == nil
	}
	if err := specificationReferenceIntegrity(documents, exists); err != nil {
		t.Fatalf("statefile.spec-reference.positive.accept: expected valid references, actual %v", err)
	}
	if err := specificationPhaseReferenceIntegrity(documents); err != nil {
		t.Fatalf("statefile.spec-phase.positive.accept: expected registered phases, actual %v", err)
	}
	if reason := specificationGovernanceIntegrity(documents["docs/SPEC.md"]); reason != "" {
		t.Fatalf("statefile.spec-governance.positive.accept: fixture.md#specification-reference-checker expected valid responsibility routing, actual %s", reason)
	}
	if reason := specificationTestRoutingIntegrity(documents); reason != "" {
		t.Fatalf("statefile.spec-routing.positive.accept: fixture.md#test-artifact-traceability-contract expected valid owners and test routing, actual %s", reason)
	}
	for _, control := range []struct{ id, from, to, reason string }{
		{"owner-missing", "`mcp`、`obsidian` の14件", "`mcp` の14件", "assertion-owner"},
		{"owner-extra", "`mcp`、`obsidian` の14件", "`mcp`、`obsidian`、`unknown` の14件", "assertion-owner"},
		{"local-test", "../../components/commitstatus/commitstatus_test.go", "missing-local-test", "test-routing"},
		{"obsidian-test", "../../components/obsidian/obsidian_test.go", "missing-obsidian-test", "test-routing"},
		{"old-glob", "`components/*/*_test.go`", "`components/*_test.go`", "test-glob"},
	} {
		t.Run("routing-"+control.id, func(t *testing.T) {
			changed := make(map[string]string, len(documents))
			for path, content := range documents {
				changed[path] = content
			}
			if !strings.Contains(changed["docs/details/fixture.md"], control.from) {
				t.Fatalf("statefile.spec-routing.%s.fixture: expected mutation target, actual missing", control.id)
			}
			changed["docs/details/fixture.md"] = strings.ReplaceAll(changed["docs/details/fixture.md"], control.from, control.to)
			if reason := specificationTestRoutingIntegrity(changed); reason != control.reason {
				t.Fatalf("statefile.spec-routing.%s.reject: fixture.md#test-artifact-traceability-contract expected=%s actual=%s; reproduce go test . -run TestSpecificationReferenceIntegrity -count=1", control.id, control.reason, reason)
			}
		})
	}
	for _, control := range []struct{ id, anchor, from, to string }{
		{"artifact-absence", "sec-4-2a", "", "単独の専用 Go 実装 artifact を持たない"},
		{"layout-current-state", "sec-4-3", "", "実在化済み"},
		{"freeze-route", "policy-spec-freeze", "DETAIL_INDEX.md#0i-詳細節対応表", "DETAIL_INDEX.md#0i1-builder--静的-web-サイト出力"},
		{"format-source", "sec-5-3", "SPEC.md#direction-technical", "SPEC.md#sec-4-1"},
		{"format-restatement", "sec-5-3", "", "workflow から読み書きする Adlaire CI の入力、出力"},
	} {
		t.Run("governance-"+control.id, func(t *testing.T) {
			marker := `<a id="` + control.anchor + `"></a>`
			before, after, ok := strings.Cut(documents["docs/SPEC.md"], marker)
			if !ok || (control.from != "" && !strings.Contains(after, control.from)) {
				t.Fatalf("statefile.spec-governance.%s.fixture: expected mutation target at SPEC.md#%s, actual missing", control.id, control.anchor)
			}
			if control.from == "" {
				after = "\n" + control.to + "\n" + after
			} else {
				after = strings.Replace(after, control.from, control.to, 1)
			}
			if reason := specificationGovernanceIntegrity(before + marker + after); reason != control.id {
				t.Fatalf("statefile.spec-governance.%s.reject: fixture.md#specification-reference-checker expected=%s actual=%q; reproduce go test . -run TestSpecificationReferenceIntegrity -count=1", control.id, control.id, reason)
			}
		})
	}
	for _, control := range []struct{ id, expression string }{
		{"joined-list", "Phase20/21/" + strconv.Itoa(21+1)},
		{"comma-list", "Phase19,20," + strconv.Itoa(21+1)},
		{"range", "Phase19~" + strconv.Itoa(21+1)},
		{"spaced", "Phase " + strconv.Itoa(21+1)},
		{"path", "phase-" + strconv.Itoa(21+1) + "-synthetic"},
		{"japanese", "フェーズ" + strconv.Itoa(21+1)},
		{"leading-zero", "Phase020"},
		{"zero", "Phase0"},
		{"unregistered-row", "| Phase " + strconv.Itoa(21+1) + " | synthetic |"},
		{"duplicate-row", "| Phase 21 | synthetic |"},
	} {
		t.Run("phase-"+control.id, func(t *testing.T) {
			changed := make(map[string]string, len(documents))
			for path, content := range documents {
				changed[path] = content
			}
			changed["docs/ROADMAP.md"] += "\n" + control.expression + "\n"
			if err := specificationPhaseReferenceIntegrity(changed); err == nil {
				t.Fatalf("statefile.spec-phase.%s.reject: expected unregistered/duplicate Phase rejected, actual accepted; fixture.md#specification-reference-checker", control.id)
			}
		})
	}
	for _, control := range []struct {
		name    string
		content string
		want    string
	}{
		{"missing-path", "[target](missing.md)", "missing path"},
		{"missing-anchor", "[target](target.md#missing)", "missing anchor"},
		{"duplicate-anchor", "<a id=\"same\"></a><a id=\"same\"></a>", "duplicate anchor"},
	} {
		t.Run(control.name, func(t *testing.T) {
			input := map[string]string{"index.md": control.content, "target.md": "<a id=\"valid\"></a>"}
			err := specificationReferenceIntegrity(input, func(path string) bool { _, ok := input[path]; return ok })
			if err == nil || !strings.Contains(err.Error(), control.want) {
				t.Fatalf("statefile.spec-reference.%s.reject: expected %s, actual %v", control.name, control.want, err)
			}
		})
	}
	input := map[string]string{"index.md": "`[example](missing.md)`\n```text\n[example](missing.md)\n```\n[ok](target.md#valid)", "target.md": "<a id=\"valid\"></a>"}
	if err := specificationReferenceIntegrity(input, func(path string) bool { _, ok := input[path]; return ok }); err != nil {
		t.Fatalf("statefile.spec-reference.examples.accept: expected literal examples ignored, actual %v", err)
	}
}

type phase19SpecificationGroup struct {
	Anchor string `json:"anchor"`
	Group  string `json:"group"`
	Cases  int    `json:"case_count"`
}

// This shape is independent of the case-ID registry being checked.
func phase19SpecificationRequiredGroups() []phase19SpecificationGroup {
	return []phase19SpecificationGroup{
		{"phase-19-status-vectors", "cases", 21},
		{"phase-19-plan-vectors", "input_cases", 22},
		{"phase-19-plan-vectors", "matrix_cases", 11},
		{"phase-19-plan-vectors", "dag_cases", 7},
		{"phase-19-plan-vectors", "selector_cases", 7},
		{"phase-19-plan-vectors", "event_cases", 3},
		{"phase-19-worker-vectors", "ack_cases", 12},
		{"phase-19-worker-vectors", "chunk_cases", 16},
		{"phase-19-worker-vectors", "lease_cases", 12},
		{"phase-19-content-vectors", "chunk_cases", 11},
		{"phase-19-content-vectors", "query_cases", 8},
		{"phase-19-content-vectors", "wire_cases", 12},
		{"phase-19-content-vectors", "cache_cases", 9},
		{"phase-19-content-vectors", "key_cases", 3},
		{"phase-19-content-vectors", "envelope_cases", 13},
		{"phase-19-operation-vectors", "admission_cases", 8},
		{"phase-19-operation-vectors", "transition_cases", 12},
		{"phase-19-operation-vectors", "gc_cases", 8},
		{"phase-19-log-boundary-vectors", "mask_cases", 12},
		{"phase-19-log-boundary-vectors", "page_cases", 13},
		{"phase-19-log-boundary-vectors", "budget_cases", 4},
	}
}

func phase19SpecificationRegistryIntegrity(registry map[string]map[string]string) error {
	required := phase19SpecificationRequiredGroups()
	want := map[string]map[string]int{}
	total := 0
	for _, r := range required {
		if want[r.Anchor] == nil {
			want[r.Anchor] = map[string]int{}
		}
		if r.Anchor == "" || r.Group == "" || r.Cases < 2 || want[r.Anchor][r.Group] != 0 {
			return fmt.Errorf("case registry: requirements")
		}
		want[r.Anchor][r.Group] = r.Cases
		total += r.Cases
	}
	if len(required) != 21 || len(want) != 6 || total != 224 {
		return fmt.Errorf("case registry: requirements count")
	}
	if len(registry) != len(want) {
		return fmt.Errorf("case registry: scope count")
	}
	seen := map[string]map[string]bool{}
	validID := regexp.MustCompile("^[a-z][a-z0-9-]*$")
	for _, r := range required {
		groups, exists := registry[r.Anchor]
		if !exists || len(groups) != len(want[r.Anchor]) {
			return fmt.Errorf("case registry: %s group set", r.Anchor)
		}
		ids, exists := groups[r.Group]
		if !exists {
			return fmt.Errorf("case registry: %s/%s missing group", r.Anchor, r.Group)
		}
		fields := strings.Fields(ids)
		if len(fields) != r.Cases {
			return fmt.Errorf("case registry: %s/%s expected %d cases actual %d", r.Anchor, r.Group, r.Cases, len(fields))
		}
		if seen[r.Anchor] == nil {
			seen[r.Anchor] = map[string]bool{}
		}
		for _, id := range fields {
			if !validID.MatchString(id) || seen[r.Anchor][id] {
				return fmt.Errorf("case registry: %s/%s invalid or duplicate identity", r.Anchor, r.Group)
			}
			seen[r.Anchor][id] = true
		}
	}
	return nil
}

func phase19SpecificationExecutionPlan() []string {
	plan := []string{"canonical-number-token", "canonical-array-order", "canonical-object-order", "canonical-string-escape"}
	seen := map[string]bool{}
	for _, r := range phase19SpecificationRequiredGroups() {
		if !seen[r.Anchor] {
			plan = append(plan, r.Anchor+"-baseline")
			seen[r.Anchor] = true
		}
		for _, control := range []string{"unknown-id", "clone-body", "missing-case", "duplicate-id", "reorder"} {
			plan = append(plan, strings.TrimPrefix(r.Anchor, "phase-19-")+"-"+strings.ReplaceAll(r.Group, "_", "-")+"-"+control)
		}
	}
	return plan
}

func phase19SpecificationExecutionIntegrity(planned, completed []string) error {
	want := map[string]bool{}
	if len(planned) == 0 || len(completed) != len(planned) {
		return fmt.Errorf("case execution: count")
	}
	for _, id := range planned {
		if id == "" || want[id] {
			return fmt.Errorf("case execution: plan identity")
		}
		want[id] = true
	}
	seen := map[string]bool{}
	for _, id := range completed {
		if !want[id] || seen[id] {
			return fmt.Errorf("case execution: missing, extra or duplicate identity")
		}
		seen[id] = true
	}
	return nil
}

func TestSpecificationVectorRegistryIntegrity(t *testing.T) {
	t.Parallel()
	doc := phase11MustReadText(t, "docs/details/fixture.md")
	_, block, found := strings.Cut(doc, "<a id=\"specification-vector-registry-integrity\"></a>")
	_, block, jsonFound := strings.Cut(block, "```json\n")
	block, _, endFound := strings.Cut(block, "\n```")
	var declared []phase19SpecificationGroup
	if !found || !jsonFound || !endFound {
		t.Fatal("statefile.spec-registry.document.fixture: fixture.md#specification-vector-registry-integrity missing fixed input")
	}
	err := developmentActionsDecodeJSON([]byte(block), &declared, 65536, true)
	actual, _ := json.Marshal(declared)
	expected, _ := json.Marshal(phase19SpecificationRequiredGroups())
	if err != nil || !bytes.Equal(expected, actual) {
		t.Fatalf("statefile.spec-registry.document.match: fixture.md#specification-vector-registry-integrity expected=%s actual=%s error=%v", expected, actual, err)
	}
	if err := phase19SpecificationRegistryIntegrity(phase19SpecificationCaseGroups()); err != nil {
		t.Fatalf("statefile.spec-registry.original.accept: expected valid registry actual=%v", err)
	}
	controls := []struct{ anchor, group, operation string }{{"", "", "extra-suite"}}
	seen := map[string]bool{}
	for _, r := range phase19SpecificationRequiredGroups() {
		if !seen[r.Anchor] {
			for _, op := range []string{"drop-suite", "rename-suite", "empty-suite", "extra-group"} {
				controls = append(controls, struct{ anchor, group, operation string }{r.Anchor, "", op})
			}
			seen[r.Anchor] = true
		}
		for _, op := range []string{"drop", "rename", "short", "long", "duplicate", "invalid", "empty"} {
			controls = append(controls, struct{ anchor, group, operation string }{r.Anchor, r.Group, op})
		}
	}
	if len(controls) != 172 {
		t.Fatalf("statefile.spec-registry.controls.fixture: expected 172 controls actual=%d", len(controls))
	}
	for _, c := range controls {
		name := strings.TrimPrefix(c.anchor, "phase-19-") + "-" + strings.ReplaceAll(c.group, "_", "-") + "-" + c.operation
		name = strings.Trim(name, "-")
		t.Run(name, func(t *testing.T) {
			registry := phase19SpecificationCaseGroups()
			groups := registry[c.anchor]
			switch c.operation {
			case "extra-suite":
				registry["unknown-suite"] = map[string]string{"cases": "one two"}
			case "drop-suite":
				delete(registry, c.anchor)
			case "rename-suite":
				registry["unknown-suite"] = groups
				delete(registry, c.anchor)
			case "empty-suite":
				registry[c.anchor] = nil
			case "extra-group":
				groups["unknown_cases"] = "one two"
			case "drop":
				delete(groups, c.group)
			case "rename":
				groups["unknown_cases"] = groups[c.group]
				delete(groups, c.group)
			case "empty":
				groups[c.group] = ""
			default:
				ids := strings.Fields(groups[c.group])
				switch c.operation {
				case "short":
					ids = ids[1:]
				case "long":
					ids = append(ids, "extra-case")
				case "duplicate":
					ids[0] = ids[1]
				case "invalid":
					ids[0] = "INVALID_CASE"
				}
				groups[c.group] = strings.Join(ids, " ")
			}
			err := phase19SpecificationRegistryIntegrity(registry)
			if err == nil || !strings.HasPrefix(err.Error(), "case registry:") {
				t.Fatalf("statefile.spec-registry.%s.reject: fixture.md#specification-vector-registry-integrity expected case registry rejection actual=%v; reproduce go test . -run TestSpecificationVectorRegistryIntegrity -count=1", name, err)
			}
		})
	}
}

func TestSpecificationVectorExecutionIntegrity(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		id                 string
		planned, completed []string
		valid              bool
	}{
		{"original", []string{"one", "two"}, []string{"one", "two"}, true},
		{"reorder", []string{"one", "two"}, []string{"two", "one"}, true},
		{"missing", []string{"one", "two"}, []string{"one"}, false},
		{"extra", []string{"one", "two"}, []string{"one", "two", "three"}, false},
		{"duplicate", []string{"one", "two"}, []string{"one", "one"}, false},
		{"substitute", []string{"one", "two"}, []string{"one", "three"}, false},
		{"empty-plan", nil, nil, false},
		{"duplicate-plan", []string{"one", "one"}, []string{"one", "two"}, false},
		{"empty-completed", []string{"one", "two"}, nil, false},
	} {
		t.Run(c.id, func(t *testing.T) {
			err := phase19SpecificationExecutionIntegrity(c.planned, c.completed)
			if (err == nil) != c.valid || err != nil && !strings.HasPrefix(err.Error(), "case execution:") {
				outcome := "reject"
				if c.valid {
					outcome = "accept"
				}
				t.Fatalf("statefile.spec-execution.%s.%s: fixture.md#specification-vector-registry-integrity expected valid=%t actual=%v; reproduce go test . -run TestSpecificationVectorExecutionIntegrity -count=1", c.id, outcome, c.valid, err)
			}
		})
	}
}

func phase19SpecificationCaseGroups() map[string]map[string]string {
	return map[string]map[string]string{
		"phase-19-status-vectors": {
			"cases": "running different-run-attempt new-run new-attempt same-time-id same-time-attempt scope-owner scope-repo scope-sha scope-project-id scope-definition-id queued-only queued-newer cancel-requested succeeded failed cancelled recovery-required delayed-start existing-claim scope-ascii-case",
		},
		"phase-19-plan-vectors": {
			"input_cases":    "default-missing false-present zero-present empty-present optional-absent required-value required-missing required-default explicit-null string-number boolean-number integer-fraction integer-range utf8-equal utf8-over unknown-key duplicate-choice default-not-choice choice-mismatch reversed-range integer-exponent integer-negative-zero",
			"matrix_cases":   "matrix-none matrix-order matrix-exclude matrix-reinclude matrix-duplicate matrix-type matrix-empty matrix-unknown-axis matrix-long-id matrix-limit matrix-over",
			"dag_cases":      "dag-origin dag-transitive-skip dag-own-skip dag-unknown dag-cycle dag-duplicate-edge dag-duplicate-id",
			"selector_cases": "selector-tag-only-branch selector-branch-only-tag selector-both-empty selector-branch-match selector-typed-input selector-deleted-path selector-empty-source",
			"event_cases":    "event-base event-definition event-revision",
		},
		"phase-19-worker-vectors": {
			"ack_cases":   "ack-normal ack-duplex ack-reverse ack-next-seq ack-pending ack-no-pending ack-mismatch ack-duplicate ack-seq-gap ack-seq-repeat ack-error-stop ack-after-ack",
			"chunk_cases": "chunk-empty chunk-one chunk-two chunk-max chunk-tail chunk-over chunk-no-final chunk-gap chunk-duplicate-final chunk-early-final chunk-extra-byte chunk-empty-middle chunk-bad-padding chunk-newline chunk-digest chunk-changed-digest",
			"lease_cases": "lease-before lease-wall-equal lease-mono-equal lease-backward lease-renew lease-old-ack lease-same-ack lease-generation lease-save-failed lease-too-far lease-skew-limit lease-no-query-renew",
		},
		"phase-19-content-vectors": {
			"chunk_cases":    "chunk-one chunk-exact chunk-first chunk-last chunk-unaligned chunk-tail chunk-future chunk-negative chunk-empty-metadata chunk-max chunk-large-metadata",
			"query_cases":    "query-zero query-max query-large query-leading-zero query-sign query-space query-exponent query-overflow",
			"wire_cases":     "receive-one receive-tail receive-padding receive-digest receive-offset receive-short receive-final receive-empty receive-gap receive-zero receive-overlong receive-whitespace",
			"cache_cases":    "cache-hit cache-no-metadata cache-no-payload cache-corrupt cache-other-key cache-deleting cache-disabled-quota cache-at-expiry cache-after-expiry",
			"key_cases":      "key-alpha key-beta key-env",
			"envelope_cases": "envelope-empty-file envelope-one envelope-block envelope-split envelope-extra-member envelope-gzip-zero envelope-gzip-byte envelope-crc envelope-truncated envelope-tar-zero envelope-short-end envelope-padding envelope-pax",
		},
		"phase-19-operation-vectors": {
			"admission_cases":  "new-create new-denied replay-denied replay-stale-revision replay-current digest-conflict subject-conflict new-stale",
			"transition_cases": "accept-success accept-failure accept-uncertain recover-success recover-failure accept-noop recover-noop success-noop failure-noop terminal-reverse terminal-change recover-reaccept",
			"gc_cases":         "gc-before gc-at gc-after gc-failure gc-accepted gc-recovery gc-protected gc-protected-failure",
		},
		"phase-19-log-boundary-vectors": {
			"mask_cases":   "split-secret longest-overlap single-byte partial-at-eof invalid-tail replacement-not-rescanned unicode-split unicode-secret mask-before-decode short-pattern-at-eof invalid-run empty",
			"page_cases":   "empty-page ascii-page tail stale future middle-of-scalar too-small unicode-exact unread-terminal emoji corrupt-utf8 limit-zero limit-large",
			"budget_cases": "below exact combined-over one-stream-over",
		},
	}
}

// Case identities are independent of the document under test.
func phase19SpecificationCaseCoverage(input []byte, groups map[string]string) error {
	var document map[string]any
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.UseNumber()
	if err := decoder.Decode(&document); err != nil || len(groups) == 0 {
		return fmt.Errorf("case coverage: invalid registry or input")
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		cases, ok := document[key].([]any)
		expected := strings.Fields(groups[key])
		if !ok || len(cases) != len(expected) {
			return fmt.Errorf("case coverage: %s count", key)
		}
		seenIDs, seenBodies := map[string]bool{}, map[string]bool{}
		for _, value := range cases {
			record, ok := value.(map[string]any)
			id, idOK := record["id"].(string)
			if !ok || !idOK || !phase17StringInSlice(id, expected) || seenIDs[id] {
				return fmt.Errorf("case coverage: %s identity", key)
			}
			seenIDs[id] = true
			delete(record, "id")
			body := string(iacSpecificationCanonical(record))
			if seenBodies[body] {
				return fmt.Errorf("case coverage: %s duplicate body", key)
			}
			seenBodies[body] = true
		}
	}
	return nil
}

func decodePhase19SpecificationVector(anchor string, input []byte) error {
	switch anchor {
	case "phase-19-status-vectors":
		_, err := decodePhase19StatusSpecification(input)
		return err
	case "phase-19-plan-vectors":
		_, err := decodePhase19PlanSpecification(input)
		return err
	case "phase-19-worker-vectors":
		_, err := decodePhase19WorkerSpecification(input)
		return err
	case "phase-19-content-vectors":
		_, err := decodePhase19ContentSpecification(input)
		return err
	case "phase-19-operation-vectors":
		_, err := decodePhase19OperationSpecification(input)
		return err
	case "phase-19-log-boundary-vectors":
		_, err := decodePhase19LogSpecification(input)
		return err
	default:
		return fmt.Errorf("unknown vector anchor")
	}
}

func TestSpecificationVectorCaseIntegrity(t *testing.T) {
	t.Parallel()
	planned := phase19SpecificationExecutionPlan()
	completed := []string{}
	if len(planned) != 115 {
		t.Fatalf("statefile.spec-execution.plan.fixture: expected 115 identities actual=%d", len(planned))
	}
	if err := phase19SpecificationRegistryIntegrity(phase19SpecificationCaseGroups()); err != nil {
		t.Fatalf("statefile.spec-registry.current.accept: expected complete registry actual=%v", err)
	}
	for _, c := range []struct {
		id    string
		input string
		valid bool
	}{
		{"number-token", `{"cases":[{"id":"one","value":1},{"id":"two","value":1e0}]}`, true},
		{"array-order", `{"cases":[{"id":"one","value":[1,2]},{"id":"two","value":[2,1]}]}`, true},
		{"object-order", `{"cases":[{"id":"one","value":{"a":1,"b":2}},{"id":"two","value":{"b":2,"a":1}}]}`, false},
		{"string-escape", `{"cases":[{"id":"one","value":"a"},{"id":"two","value":"\u0061"}]}`, false},
	} {
		t.Run(c.id, func(t *testing.T) {
			err := phase19SpecificationCaseCoverage([]byte(c.input), map[string]string{"cases": "one two"})
			if (err == nil) != c.valid {
				t.Fatalf("statefile.spec-vector.%s.canonical: fixture.md#specification-vector-case-integrity expected valid=%t actual=%v; reproduce go test . -run TestSpecificationVectorCaseIntegrity -count=1", c.id, c.valid, err)
			}
			completed = append(completed, "canonical-"+c.id)
		})
	}
	document := phase11MustReadText(t, "docs/details/fixture.md")
	registry := map[string]map[string]int{}
	for _, r := range phase19SpecificationRequiredGroups() {
		if registry[r.Anchor] == nil {
			registry[r.Anchor] = map[string]int{}
		}
		registry[r.Anchor][r.Group] = r.Cases
	}
	anchors := make([]string, 0, len(registry))
	for anchor := range registry {
		anchors = append(anchors, anchor)
	}
	sort.Strings(anchors)
	for _, anchor := range anchors {
		_, block, found := strings.Cut(document, "<a id=\""+anchor+"\"></a>")
		_, block, jsonFound := strings.Cut(block, "```json\n")
		block, _, endFound := strings.Cut(block, "\n```")
		if !found || !jsonFound || !endFound {
			t.Fatalf("statefile.spec-vector.%s.fixture: fixture.md#%s fixed input missing", anchor, anchor)
		}
		input := []byte(block)
		if err := decodePhase19SpecificationVector(anchor, input); err != nil {
			t.Fatalf("statefile.spec-vector.%s.accept: fixture.md#%s expected valid input actual=%v", anchor, anchor, err)
		}
		completed = append(completed, anchor+"-baseline")
		keys := make([]string, 0, len(registry[anchor]))
		for key := range registry[anchor] {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			for _, control := range []string{"unknown-id", "clone-body", "missing-case", "duplicate-id", "reorder"} {
				name := strings.TrimPrefix(anchor, "phase-19-") + "-" + strings.ReplaceAll(key, "_", "-") + "-" + control
				t.Run(name, func(t *testing.T) {
					var raw map[string]any
					decoder := json.NewDecoder(bytes.NewReader(input))
					decoder.UseNumber()
					if err := decoder.Decode(&raw); err != nil {
						t.Fatalf("statefile.spec-vector.%s.fixture: %v", name, err)
					}
					cases := raw[key].([]any)
					first, second := cases[0].(map[string]any), cases[1].(map[string]any)
					switch control {
					case "unknown-id":
						first["id"] = "replacement-case"
					case "clone-body":
						id := first["id"]
						for k := range first {
							delete(first, k)
						}
						for k, v := range second {
							first[k] = v
						}
						first["id"] = id
					case "missing-case":
						raw[key] = cases[1:]
					case "duplicate-id":
						first["id"] = second["id"]
					case "reorder":
						for i, j := 0, len(cases)-1; i < j; i, j = i+1, j-1 {
							cases[i], cases[j] = cases[j], cases[i]
						}
					}
					mutated, err := json.Marshal(raw)
					if err != nil {
						t.Fatalf("statefile.spec-vector.%s.fixture: %v", name, err)
					}
					err = decodePhase19SpecificationVector(anchor, mutated)
					if control == "reorder" {
						if err != nil {
							t.Fatalf("statefile.spec-vector.%s.accept: fixture.md#specification-vector-case-integrity expected valid reordered set actual=%v", name, err)
						}
					} else if err == nil || !strings.HasPrefix(err.Error(), "case coverage:") {
						t.Fatalf("statefile.spec-vector.%s.reject: fixture.md#specification-vector-case-integrity expected case coverage rejection actual=%v; reproduce go test . -run TestSpecificationVectorCaseIntegrity -count=1", name, err)
					}
					completed = append(completed, name)
				})
			}
		}
	}
	if err := phase19SpecificationExecutionIntegrity(planned, completed); err != nil {
		t.Fatalf("statefile.spec-execution.complete.match: fixture.md#specification-vector-registry-integrity expected planned/completed set equality actual=%v; reproduce go test . -run TestSpecificationVectorCaseIntegrity -count=1", err)
	}
}

type phase19StatusRun struct {
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	SHA        string `json:"sha"`
	Project    string `json:"project_id"`
	Definition string `json:"definition_id"`
	Run        string `json:"run_id"`
	Attempt    int    `json:"attempt"`
	Created    string `json:"created_at"`
	State      string `json:"state"`
}

type phase19StatusDecision struct {
	Run      string `json:"run_id"`
	Attempt  int    `json:"attempt"`
	State    string `json:"state"`
	Decision string `json:"decision"`
}

type phase19StatusSpecification struct {
	Version  int              `json:"schema_version"`
	Baseline phase19StatusRun `json:"baseline"`
	Cases    []struct {
		ID       string                       `json:"id"`
		Runs     []map[string]json.RawMessage `json:"runs"`
		Event    map[string]json.RawMessage   `json:"event"`
		Existing bool                         `json:"existing"`
		Expected phase19StatusDecision        `json:"expected"`
	} `json:"cases"`
}

func phase19StatusProjection(runs []phase19StatusRun, event phase19StatusRun, existing bool) phase19StatusDecision {
	fold := func(s string) string {
		return strings.Map(func(r rune) rune {
			if r >= 'A' && r <= 'Z' {
				return r + ('a' - 'A')
			}
			return r
		}, s)
	}
	candidates := []phase19StatusRun{}
	for _, r := range runs {
		if fold(r.Owner) == fold(event.Owner) && fold(r.Repo) == fold(event.Repo) && r.SHA == event.SHA && r.Project == event.Project && r.Definition == event.Definition && r.State != "queued" {
			candidates = append(candidates, r)
		}
	}
	if len(candidates) == 0 {
		return phase19StatusDecision{Decision: "idle"}
	}
	sort.Slice(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.Created != b.Created {
			return a.Created < b.Created
		}
		if a.Run != b.Run {
			return a.Run < b.Run
		}
		return a.Attempt < b.Attempt
	})
	latest := candidates[len(candidates)-1]
	states := map[string]string{"running": "pending", "cancel_requested": "pending", "succeeded": "success", "failed": "failure", "cancelled": "error", "recovery_required": "error"}
	result := phase19StatusDecision{Run: latest.Run, Attempt: latest.Attempt, State: states[latest.State], Decision: "send"}
	if event.Run != latest.Run || event.Attempt != latest.Attempt || event.State != latest.State {
		result.Decision = "superseded"
	} else if existing {
		result.Decision = "deduplicated"
	}
	return result
}

func phase19StatusExpand(base phase19StatusRun, patch map[string]json.RawMessage) (phase19StatusRun, error) {
	var r phase19StatusRun
	if patch == nil {
		return r, fmt.Errorf("null projection")
	}
	data, _ := json.Marshal(base)
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(data, &fields)
	for k, v := range patch {
		if _, ok := fields[k]; !ok || bytes.Equal(v, []byte("null")) {
			return r, fmt.Errorf("projection field")
		}
		fields[k] = v
	}
	data, _ = json.Marshal(fields)
	if err := developmentActionsDecodeJSON(data, &r, 65536, true); err != nil {
		return r, err
	}
	clock, err := time.Parse("2006-01-02T15:04:05Z", r.Created)
	if err != nil || clock.Year() < 1970 || clock.Format("2006-01-02T15:04:05Z") != r.Created || r.Attempt < 1 || r.Attempt > 100 || !regexp.MustCompile("^[0-9a-f]{40}$").MatchString(r.SHA) {
		return r, fmt.Errorf("projection value")
	}
	for _, name := range []string{r.Owner, r.Repo} {
		if !utf8.ValidString(name) || len(name) < 1 || len(name) > 100 || strings.ContainsRune(name, 0) {
			return r, fmt.Errorf("projection source name")
		}
	}
	for _, id := range []string{r.Project, r.Definition, r.Run} {
		if !regexp.MustCompile("^[a-z][a-z0-9_-]{0,63}$").MatchString(id) {
			return r, fmt.Errorf("projection identity")
		}
	}
	if !phase17StringInSlice(r.State, []string{"queued", "running", "cancel_requested", "succeeded", "failed", "cancelled", "recovery_required"}) {
		return r, fmt.Errorf("projection state")
	}
	return r, nil
}

func decodePhase19StatusSpecification(input []byte) (phase19StatusSpecification, error) {
	var v phase19StatusSpecification
	if err := developmentActionsDecodeJSON(input, &v, 65536, true); err != nil {
		return v, err
	}
	if err := phase19SpecificationCaseCoverage(input, phase19SpecificationCaseGroups()["phase-19-status-vectors"]); err != nil {
		return v, err
	}
	encoded, _ := json.Marshal(v)
	var a, b any
	if json.Unmarshal(input, &a) != nil || json.Unmarshal(encoded, &b) != nil || !bytes.Equal(iacSpecificationCanonical(a), iacSpecificationCanonical(b)) || v.Version != 1 || len(v.Cases) != 21 {
		return v, fmt.Errorf("required fields / case set")
	}
	if _, err := phase19StatusExpand(v.Baseline, map[string]json.RawMessage{}); err != nil {
		return v, err
	}
	seen := map[string]bool{}
	for _, c := range v.Cases {
		if !regexp.MustCompile("^[a-z][a-z0-9-]*$").MatchString(c.ID) || seen[c.ID] || c.Runs == nil || c.Event == nil {
			return v, fmt.Errorf("case identity")
		}
		seen[c.ID] = true
		for _, p := range append(append([]map[string]json.RawMessage{}, c.Runs...), c.Event) {
			if _, err := phase19StatusExpand(v.Baseline, p); err != nil {
				return v, err
			}
		}
		if !phase17StringInSlice(c.Expected.Decision, []string{"send", "idle", "superseded", "deduplicated"}) {
			return v, fmt.Errorf("expected decision")
		}
	}
	return v, nil
}

func TestPhase19StatusSpecificationVectors(t *testing.T) {
	t.Parallel()
	doc := phase11MustReadText(t, "docs/details/fixture.md")
	_, block, found := strings.Cut(doc, `<a id="phase-19-status-vectors"></a>`)
	_, block, jsonFound := strings.Cut(block, "```json\n")
	block, _, endFound := strings.Cut(block, "\n```")
	if !found || !jsonFound || !endFound {
		t.Fatal("commitstatus.status-spec.presence.reject: fixture.md#phase-19-status-vectors expected fixed JSON input, actual missing")
	}
	input := []byte(block)
	v, err := decodePhase19StatusSpecification(input)
	if err != nil {
		t.Fatalf("commitstatus.status-spec.schema.reject: expected valid fixed input, actual %v", err)
	}
	for _, c := range v.Cases {
		t.Run(c.ID, func(t *testing.T) {
			runs := []phase19StatusRun{}
			for _, p := range c.Runs {
				r, _ := phase19StatusExpand(v.Baseline, p)
				runs = append(runs, r)
			}
			event, _ := phase19StatusExpand(v.Baseline, c.Event)
			actual := phase19StatusProjection(runs, event, c.Existing)
			if actual != c.Expected {
				t.Fatalf("commitstatus.status-spec.%s.result: commitstatus.md#phase-19-status-delivery expected=%+v actual=%+v; reproduce go test . -run TestPhase19StatusSpecificationVectors -count=1", c.ID, c.Expected, actual)
			}
		})
	}
	for _, id := range []string{"unknown", "missing", "null", "duplicate-case", "projection-field", "projection-null"} {
		t.Run("schema-"+id, func(t *testing.T) {
			var raw map[string]any
			_ = json.Unmarshal(input, &raw)
			cases := raw["cases"].([]any)
			switch id {
			case "unknown":
				raw["unknown"] = true
			case "missing":
				delete(raw, "baseline")
			case "null":
				raw["cases"] = nil
			case "duplicate-case":
				cases[1].(map[string]any)["id"] = cases[0].(map[string]any)["id"]
			case "projection-field":
				cases[0].(map[string]any)["event"] = map[string]any{"unknown": true}
			case "projection-null":
				cases[0].(map[string]any)["event"] = map[string]any{"state": nil}
			}
			changed, _ := json.Marshal(raw)
			if _, err := decodePhase19StatusSpecification(changed); err == nil {
				t.Fatalf("commitstatus.status-spec.%s.reject: fixture.md#phase-19-status-vectors expected schema rejection, actual accepted", id)
			}
		})
	}
}

func specificationTestRoutingIntegrity(documents map[string]string) string {
	detail := documents["docs/DETAIL_INDEX.md"]
	_, detail, _ = strings.Cut(detail, `<a id="0b-詳細仕様参照表"></a>`)
	detail, _, _ = strings.Cut(detail, `<a id="owner-detail-verification-route"></a>`)
	owners := regexp.MustCompile("(?m)^\\| `([a-z]+)` \\|").FindAllStringSubmatch(detail, -1)
	fixture := documents["docs/details/fixture.md"]
	_, assertion, _ := strings.Cut(fixture, "owner segment は ")
	assertion, _, _ = strings.Cut(assertion, "必要な場合だけ")
	listed := regexp.MustCompile("`([a-z]+)`").FindAllStringSubmatch(assertion, -1)
	want, got := []string{}, []string{}
	for _, owner := range owners {
		want = append(want, owner[1])
	}
	for _, owner := range listed {
		got = append(got, owner[1])
	}
	sort.Strings(want)
	sort.Strings(got)
	if len(want) == 0 || !phase11StringSlicesEqual(want, got) {
		return "assertion-owner"
	}
	_, routing, _ := strings.Cut(fixture, `<a id="test-artifact-traceability-contract"></a>`)
	routing, _, _ = strings.Cut(routing, "\n以下は実在する実装 artifact")
	paths := []string{"main_test.go"}
	for _, owner := range want {
		path := "components/" + owner + "/" + owner + "_test.go"
		if owner == "sdk" || owner == "ui" {
			path = owner + "_contract_test.go"
		}
		paths = append(paths, path)
	}
	for _, path := range paths {
		if !strings.Contains(routing, "](../../"+path+")") {
			return "test-routing"
		}
	}
	if strings.Contains(fixture, "`components/*_test.go`") {
		return "test-glob"
	}
	return ""
}

func specificationGovernanceIntegrity(spec string) string {
	section := func(start, end string) string {
		_, body, ok := strings.Cut(spec, `<a id="`+start+`"></a>`)
		if !ok {
			return ""
		}
		body, _, ok = strings.Cut(body, `<a id="`+end+`"></a>`)
		if !ok {
			return ""
		}
		return body
	}
	principle := section("sec-4-2a", "direction-change-history-single-source")
	if principle == "" || strings.Contains(principle, "単独の専用 Go 実装 artifact を持たない") {
		return "artifact-absence"
	}
	layout := section("sec-4-3", "sec-4-4")
	if layout == "" || strings.Contains(layout, "実在化済み") {
		return "layout-current-state"
	}
	freeze := section("policy-spec-freeze", "policy-phase-unit")
	row := ""
	for _, line := range strings.Split(freeze, "\n") {
		if strings.HasPrefix(line, "| 対応表 |") {
			row = line
		}
	}
	if !strings.Contains(row, "DETAIL_INDEX.md#0i-詳細節対応表") || !strings.Contains(row, "DETAIL_INDEX.md#phase-19-owner-route") || !strings.Contains(row, "DETAIL_INDEX.md#phase-20-inhouse-iac-entry") {
		return "freeze-route"
	}
	actions := section("sec-5-3", "direction-inhouse-cicd")
	if !strings.Contains(actions, "SPEC.md#direction-technical") {
		return "format-source"
	}
	if strings.Contains(actions, "workflow から読み書きする Adlaire CI の入力、出力") {
		return "format-restatement"
	}
	return ""
}

type phase19PlanInputCase struct {
	ID       string                     `json:"id"`
	Kind     string                     `json:"kind"`
	Min      int64                      `json:"min"`
	Max      int64                      `json:"max"`
	Required bool                       `json:"required"`
	Default  json.RawMessage            `json:"default"`
	Choices  []json.RawMessage          `json:"choices"`
	Values   map[string]json.RawMessage `json:"values"`
	Expected map[string]json.RawMessage `json:"expected"`
	Code     string                     `json:"expected_code"`
}

type phase19PlanMatrixCase struct {
	ID        string                       `json:"id"`
	Base      string                       `json:"base"`
	Axes      map[string][]json.RawMessage `json:"axes"`
	Exclude   []map[string]json.RawMessage `json:"exclude"`
	Include   []map[string]json.RawMessage `json:"include"`
	Count     int                          `json:"expected_count"`
	Digest    string                       `json:"expected_digest"`
	IDsDigest string                       `json:"expected_ids_digest"`
	Code      string                       `json:"expected_code"`
}

type phase19PlanDAGJob struct {
	ID    string   `json:"id"`
	Needs []string `json:"needs"`
	Match bool     `json:"own_match"`
}

type phase19PlanEvent struct {
	Project    string  `json:"project_id"`
	Definition string  `json:"definition_id"`
	Revision   int     `json:"definition_revision"`
	Kind       string  `json:"kind"`
	Trigger    string  `json:"trigger_id"`
	SHA        string  `json:"source_sha"`
	Due        *string `json:"due"`
}

type phase19PlanSpecificationVector struct {
	Version  int                     `json:"schema_version"`
	Inputs   []phase19PlanInputCase  `json:"input_cases"`
	Matrices []phase19PlanMatrixCase `json:"matrix_cases"`
	DAGs     []struct {
		ID       string              `json:"id"`
		Jobs     []phase19PlanDAGJob `json:"jobs"`
		IDs      []string            `json:"expected_ids"`
		Selected []bool              `json:"expected_selected"`
		Code     string              `json:"expected_code"`
	} `json:"dag_cases"`
	Selectors []struct {
		ID       string                     `json:"id"`
		Kind     string                     `json:"kind"`
		Branches []string                   `json:"branches"`
		Tags     []string                   `json:"tags"`
		Paths    []string                   `json:"paths"`
		Changed  []string                   `json:"changed"`
		Equals   map[string]json.RawMessage `json:"equals"`
		Values   map[string]json.RawMessage `json:"values"`
		Expected bool                       `json:"expected"`
	} `json:"selector_cases"`
	Events []struct {
		ID     string           `json:"id"`
		Input  phase19PlanEvent `json:"input"`
		Digest string           `json:"expected_digest"`
	} `json:"event_cases"`
}

func phase19PlanScalar(raw json.RawMessage) (string, []byte, bool) {
	var value any
	if json.Unmarshal(raw, &value) != nil {
		return "", nil, false
	}
	canonical, err := json.Marshal(value)
	if err != nil {
		return "", nil, false
	}
	switch value.(type) {
	case string:
		return "string", canonical, true
	case bool:
		return "boolean", canonical, true
	case float64:
		s := strings.TrimSpace(string(raw))
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil || strconv.FormatInt(n, 10) != s || n < -2147483648 || n > 2147483647 {
			return "", nil, false
		}
		return "integer", canonical, true
	default:
		return "", nil, false
	}
}

func phase19PlanInputOracle(c phase19PlanInputCase) (map[string]json.RawMessage, string) {
	empty := map[string]json.RawMessage{}
	fail := func() (map[string]json.RawMessage, string) { return empty, "P19_SCHEMA" }
	if c.Min > c.Max || (c.Kind == "string" && (c.Min < 0 || c.Max > 4096)) ||
		(c.Kind == "boolean" && (c.Min != 0 || c.Max != 1)) ||
		!phase17StringInSlice(c.Kind, []string{"string", "integer", "boolean"}) {
		return fail()
	}
	valid := func(raw json.RawMessage) bool {
		kind, _, ok := phase19PlanScalar(raw)
		if !ok || kind != c.Kind {
			return false
		}
		if kind == "string" {
			var s string
			if json.Unmarshal(raw, &s) != nil {
				return false
			}
			return int64(len([]byte(s))) >= c.Min && int64(len([]byte(s))) <= c.Max
		}
		if kind == "integer" {
			n, err := strconv.ParseInt(strings.TrimSpace(string(raw)), 10, 64)
			return err == nil && n >= c.Min && n <= c.Max
		}
		return true
	}
	choiceSet := map[string]bool{}
	for _, raw := range c.Choices {
		if !valid(raw) {
			return fail()
		}
		_, key, _ := phase19PlanScalar(raw)
		if choiceSet[string(key)] {
			return fail()
		}
		choiceSet[string(key)] = true
	}
	allowed := func(raw json.RawMessage) bool {
		if !valid(raw) {
			return false
		}
		_, key, _ := phase19PlanScalar(raw)
		return len(choiceSet) == 0 || choiceSet[string(key)]
	}
	hasDefault := string(c.Default) != "null"
	if (c.Required && hasDefault) || (hasDefault && !allowed(c.Default)) {
		return fail()
	}
	for key := range c.Values {
		if key != "x" {
			return fail()
		}
	}
	raw, present := c.Values["x"]
	if !present {
		if c.Required {
			return fail()
		}
		if !hasDefault {
			return empty, "ok"
		}
		raw = c.Default
	}
	if !allowed(raw) {
		return fail()
	}
	return map[string]json.RawMessage{"x": raw}, "ok"
}

func phase19PlanCanonicalHash(value any) string {
	b, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	var decoded any
	if json.Unmarshal(b, &decoded) != nil {
		return ""
	}
	h := sha256.Sum256(iacSpecificationCanonical(decoded))
	return hex.EncodeToString(h[:])
}

func phase19PlanMatrixOracle(c phase19PlanMatrixCase) ([]map[string]json.RawMessage, []string, string) {
	fail := func() ([]map[string]json.RawMessage, []string, string) {
		return []map[string]json.RawMessage{}, []string{}, "P19_MATRIX"
	}
	axes := []string{}
	kinds := map[string]string{}
	if len(c.Axes) > 4 {
		return fail()
	}
	for axis, values := range c.Axes {
		if len(values) == 0 || len(values) > 32 {
			return fail()
		}
		axes = append(axes, axis)
		seen := map[string]bool{}
		for _, raw := range values {
			kind, key, ok := phase19PlanScalar(raw)
			if !ok || (kinds[axis] != "" && kinds[axis] != kind) || seen[string(key)] {
				return fail()
			}
			kinds[axis] = kind
			seen[string(key)] = true
		}
	}
	sort.Strings(axes)
	for _, list := range []struct {
		values   []map[string]json.RawMessage
		complete bool
	}{{c.Exclude, false}, {c.Include, true}} {
		for _, tuple := range list.values {
			if (list.complete && len(tuple) != len(axes)) || (!list.complete && len(tuple) == 0) {
				return fail()
			}
			for key, raw := range tuple {
				kind, _, ok := phase19PlanScalar(raw)
				if !ok || kinds[key] == "" || kinds[key] != kind {
					return fail()
				}
			}
		}
	}
	tuples := []map[string]json.RawMessage{}
	var walk func(int, map[string]json.RawMessage)
	walk = func(index int, tuple map[string]json.RawMessage) {
		if index < len(axes) {
			axis := axes[index]
			for _, raw := range c.Axes[axis] {
				tuple[axis] = raw
				walk(index+1, tuple)
			}
			return
		}
		for _, excluded := range c.Exclude {
			matches := true
			for key, raw := range excluded {
				_, a, _ := phase19PlanScalar(raw)
				_, b, _ := phase19PlanScalar(tuple[key])
				if !bytes.Equal(a, b) {
					matches = false
					break
				}
			}
			if matches {
				return
			}
		}
		copyTuple := map[string]json.RawMessage{}
		for key, raw := range tuple {
			copyTuple[key] = raw
		}
		tuples = append(tuples, copyTuple)
	}
	walk(0, map[string]json.RawMessage{})
	seen := map[string]bool{}
	for _, tuple := range tuples {
		seen[phase19PlanCanonicalHash(tuple)] = true
	}
	for _, tuple := range c.Include {
		h := phase19PlanCanonicalHash(tuple)
		if seen[h] {
			return fail()
		}
		seen[h] = true
		tuples = append(tuples, tuple)
	}
	if len(tuples) == 0 || len(tuples) > 256 {
		return fail()
	}
	ids := []string{}
	for i := range tuples {
		id := c.Base
		if len(axes) > 0 {
			id += "_m" + strconv.Itoa(i+1)
		}
		if len(id) > 64 {
			return fail()
		}
		ids = append(ids, id)
	}
	return tuples, ids, "ok"
}

func phase19PlanDAGOracle(jobs []phase19PlanDAGJob) ([]string, []bool, string) {
	fail := func() ([]string, []bool, string) { return []string{}, []bool{}, "P19_DAG" }
	index := map[string]int{}
	for i, job := range jobs {
		if _, ok := index[job.ID]; ok {
			return fail()
		}
		index[job.ID] = i
	}
	degrees := make([]int, len(jobs))
	done := make([]bool, len(jobs))
	selection := map[string]bool{}
	for i, job := range jobs {
		seen := map[string]bool{}
		for _, need := range job.Needs {
			if _, ok := index[need]; !ok || need == job.ID || seen[need] {
				return fail()
			}
			seen[need] = true
			degrees[i]++
		}
	}
	ids := []string{}
	selected := []bool{}
	for len(ids) < len(jobs) {
		chosen := -1
		for i := range jobs {
			if !done[i] && degrees[i] == 0 {
				chosen = i
				break
			}
		}
		if chosen < 0 {
			return fail()
		}
		job := jobs[chosen]
		value := job.Match
		for _, need := range job.Needs {
			value = value && selection[need]
		}
		done[chosen] = true
		selection[job.ID] = value
		ids = append(ids, job.ID)
		selected = append(selected, value)
		for i, other := range jobs {
			for _, need := range other.Needs {
				if need == job.ID {
					degrees[i]--
				}
			}
		}
	}
	return ids, selected, "ok"
}

func phase19PlanSelectorOracle(kind string, branches, tags, paths, changed []string, equals, values map[string]json.RawMessage) bool {
	list, other, ref := branches, tags, "main"
	if kind == "tag" {
		list, other, ref = tags, branches, "v1"
	}
	if len(list) == 0 && len(other) > 0 {
		return false
	}
	if len(list) > 0 && !phase17StringInSlice(ref, list) {
		return false
	}
	if len(paths) > 0 {
		matched := false
		for _, p := range paths {
			if phase17StringInSlice(p, changed) {
				matched = true
			}
		}
		if !matched {
			return false
		}
	}
	for key, want := range equals {
		got, ok := values[key]
		if !ok {
			return false
		}
		wk, wb, wok := phase19PlanScalar(want)
		gk, gb, gok := phase19PlanScalar(got)
		if !wok || !gok || wk != gk || !bytes.Equal(wb, gb) {
			return false
		}
	}
	return true
}

func decodePhase19PlanSpecification(input []byte) (phase19PlanSpecificationVector, error) {
	var v phase19PlanSpecificationVector
	if err := developmentActionsDecodeJSON(input, &v, 65536, true); err != nil {
		return v, err
	}
	if err := phase19SpecificationCaseCoverage(input, phase19SpecificationCaseGroups()["phase-19-plan-vectors"]); err != nil {
		return v, err
	}
	encoded, err := json.Marshal(v)
	var a, b any
	if err != nil || json.Unmarshal(input, &a) != nil || json.Unmarshal(encoded, &b) != nil || !bytes.Equal(iacSpecificationCanonical(a), iacSpecificationCanonical(b)) {
		return v, fmt.Errorf("required fields")
	}
	if v.Version != 1 || v.Inputs == nil || v.Matrices == nil || v.DAGs == nil || v.Selectors == nil || v.Events == nil {
		return v, fmt.Errorf("vector identity")
	}
	seen := map[string]bool{}
	checkID := func(id string) bool {
		if !regexp.MustCompile("^[a-z][a-z0-9-]*$").MatchString(id) || seen[id] {
			return false
		}
		seen[id] = true
		return true
	}
	code := func(s string) bool {
		return phase17StringInSlice(s, []string{"ok", "P19_SCHEMA", "P19_MATRIX", "P19_DAG"})
	}
	for _, c := range v.Inputs {
		if !checkID(c.ID) || c.Choices == nil || c.Values == nil || c.Expected == nil || len(c.Default) == 0 || !code(c.Code) {
			return v, fmt.Errorf("input identity")
		}
	}
	for _, c := range v.Matrices {
		if !checkID(c.ID) || c.Axes == nil || c.Exclude == nil || c.Include == nil || !phase17IsLowerHexSHA256(c.Digest) || !phase17IsLowerHexSHA256(c.IDsDigest) || !code(c.Code) {
			return v, fmt.Errorf("matrix identity")
		}
	}
	for _, c := range v.DAGs {
		if !checkID(c.ID) || c.Jobs == nil || c.IDs == nil || c.Selected == nil || !code(c.Code) {
			return v, fmt.Errorf("dag identity")
		}
		for _, j := range c.Jobs {
			if j.Needs == nil {
				return v, fmt.Errorf("dag needs")
			}
		}
	}
	for _, c := range v.Selectors {
		if !checkID(c.ID) || !phase17StringInSlice(c.Kind, []string{"branch", "tag"}) || c.Branches == nil || c.Tags == nil || c.Paths == nil || c.Changed == nil || c.Equals == nil || c.Values == nil {
			return v, fmt.Errorf("selector identity")
		}
	}
	for _, c := range v.Events {
		if !checkID(c.ID) || c.Input.Project == "" || c.Input.Definition == "" || c.Input.Revision < 1 || !phase17IsLowerHexSHA256(c.Digest) {
			return v, fmt.Errorf("event identity")
		}
	}
	return v, nil
}

func TestPhase19PlanSpecificationVectors(t *testing.T) {
	t.Parallel()
	doc := phase11MustReadText(t, "docs/details/fixture.md")
	match := regexp.MustCompile("(?s)<a id=\"phase-19-plan-vectors\"></a>.*?\x60\x60\x60json\\n(.*?)\\n\x60\x60\x60").FindStringSubmatch(doc)
	if len(match) != 2 {
		t.Fatal("runner.plan-spec.presence: missing fixed input")
	}
	input := []byte(match[1])
	v, err := decodePhase19PlanSpecification(input)
	if err != nil || len(v.Inputs) != 22 || len(v.Matrices) != 11 || len(v.DAGs) != 7 || len(v.Selectors) != 7 || len(v.Events) != 3 {
		t.Fatalf("runner.plan-spec.schema: expected 22/11/7/7/3 vectors actual %v", err)
	}
	for _, c := range v.Inputs {
		t.Run(c.ID, func(t *testing.T) {
			actual, code := phase19PlanInputOracle(c)
			if code != c.Code || phase19PlanCanonicalHash(actual) != phase19PlanCanonicalHash(c.Expected) {
				t.Fatalf("runner.plan-spec.%s.result: runner.md#phase-19-input-normalization expected=%s/%s actual=%s/%s", c.ID, c.Code, phase19PlanCanonicalHash(c.Expected), code, phase19PlanCanonicalHash(actual))
			}
		})
	}
	for _, c := range v.Matrices {
		t.Run(c.ID, func(t *testing.T) {
			tuples, ids, code := phase19PlanMatrixOracle(c)
			if code != c.Code || len(tuples) != c.Count || phase19PlanCanonicalHash(tuples) != c.Digest || phase19PlanCanonicalHash(ids) != c.IDsDigest {
				t.Fatalf("runner.plan-spec.%s.result: runner.md#phase-19-matrix-order expected=%s/%d/%s/%s actual=%s/%d/%s/%s", c.ID, c.Code, c.Count, c.Digest, c.IDsDigest, code, len(tuples), phase19PlanCanonicalHash(tuples), phase19PlanCanonicalHash(ids))
			}
		})
	}
	for _, c := range v.DAGs {
		t.Run(c.ID, func(t *testing.T) {
			ids, selected, code := phase19PlanDAGOracle(c.Jobs)
			if code != c.Code || phase19PlanCanonicalHash(ids) != phase19PlanCanonicalHash(c.IDs) || phase19PlanCanonicalHash(selected) != phase19PlanCanonicalHash(c.Selected) {
				t.Fatalf("runner.plan-spec.%s.result: runner.md#phase-19-plan-selection expected=%s/%v/%v actual=%s/%v/%v", c.ID, c.Code, c.IDs, c.Selected, code, ids, selected)
			}
		})
	}
	for _, c := range v.Selectors {
		t.Run(c.ID, func(t *testing.T) {
			actual := phase19PlanSelectorOracle(c.Kind, c.Branches, c.Tags, c.Paths, c.Changed, c.Equals, c.Values)
			if actual != c.Expected {
				t.Fatalf("runner.plan-spec.%s.result: runner.md#phase-19-plan-selection expected=%t actual=%t", c.ID, c.Expected, actual)
			}
		})
	}
	for _, c := range v.Events {
		t.Run(c.ID, func(t *testing.T) {
			actual := phase19PlanCanonicalHash(c.Input)
			if actual != c.Digest {
				t.Fatalf("runner.plan-spec.%s.result: runner.md#phase-19-trigger-state expected=%s actual=%s", c.ID, c.Digest, actual)
			}
		})
	}
	var nullValue map[string]json.RawMessage
	if err := json.Unmarshal(input, &nullValue); err != nil {
		t.Fatal(err)
	}
	nullValue["input_cases"] = json.RawMessage("null")
	nullInput, err := json.Marshal(nullValue)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		id    string
		input []byte
	}{
		{"missing", []byte(strings.Replace(string(input), "\"schema_version\": 1,", "", 1))},
		{"unknown", []byte(strings.Replace(string(input), "\"schema_version\": 1,", "\"schema_version\": 1,\"unknown\":true,", 1))},
		{"duplicate", []byte(strings.Replace(string(input), "\"schema_version\": 1,", "\"schema_version\": 1,\"schema_version\": 1,", 1))},
		{"null", nullInput},
		{"trailing", append(append([]byte{}, input...), []byte("{}")...)},
		{"invalid-utf8", append([]byte{0xff}, input...)},
	} {
		t.Run("schema-"+c.id, func(t *testing.T) {
			if _, err := decodePhase19PlanSpecification(c.input); err == nil {
				t.Fatalf("runner.plan-spec.schema-%s.reject: fixture.md#phase-19-plan-vectors expected rejected actual accepted", c.id)
			}
		})
	}
}

type phase19WorkerFrame struct {
	Direction string `json:"direction"`
	Type      string `json:"type"`
	Seq       int64  `json:"seq"`
	Ack       int64  `json:"ack_seq"`
}

type phase19WorkerAckCase struct {
	ID      string               `json:"id"`
	Frames  []phase19WorkerFrame `json:"frames"`
	Pending []int                `json:"expected_pending"`
	Code    string               `json:"expected_code"`
}

type phase19WorkerChunk struct {
	Offset int64  `json:"offset"`
	Data   string `json:"data"`
	Repeat int    `json:"repeat"`
	Final  bool   `json:"final"`
	Digest string `json:"digest"`
}

type phase19WorkerChunkCase struct {
	ID     string               `json:"id"`
	Total  int64                `json:"total_bytes"`
	Digest string               `json:"digest"`
	Chunks []phase19WorkerChunk `json:"chunks"`
	Code   string               `json:"expected_code"`
}

type phase19WorkerLeaseCase struct {
	ID           string `json:"id"`
	Wall         int64  `json:"wall_ms"`
	Mono         int64  `json:"mono_ms"`
	Until        int64  `json:"until_ms"`
	Deadline     int64  `json:"deadline_ms"`
	Proposed     int64  `json:"proposed_ms"`
	Generation   bool   `json:"generation_matches"`
	Durable      bool   `json:"durable"`
	WantUntil    int64  `json:"expected_until_ms"`
	WantDeadline int64  `json:"expected_deadline_ms"`
	Code         string `json:"expected_code"`
}

type phase19WorkerSpecificationVector struct {
	Version int                      `json:"schema_version"`
	Acks    []phase19WorkerAckCase   `json:"ack_cases"`
	Chunks  []phase19WorkerChunkCase `json:"chunk_cases"`
	Leases  []phase19WorkerLeaseCase `json:"lease_cases"`
}

func phase19WorkerAckOracle(frames []phase19WorkerFrame) ([2]int, string) {
	next, waiting := [2]int64{1, 1}, [2]int64{}
	pending := func() [2]int {
		count := [2]int{}
		for i, seq := range waiting {
			if seq != 0 {
				count[i] = 1
			}
		}
		return count
	}
	for _, f := range frames {
		d := 0
		if f.Direction == "w2c" {
			d = 1
		} else if f.Direction != "c2w" {
			return pending(), "P19_INPUT"
		}
		if f.Seq != next[d] || f.Seq > 9007199254740991 {
			return pending(), "P19_INPUT"
		}
		switch f.Type {
		case "ack":
			if waiting[1-d] == 0 || f.Ack != waiting[1-d] {
				return pending(), "P19_INPUT"
			}
			waiting[1-d] = 0
		case "error":
			return pending(), "P19_TRANSPORT"
		default:
			if !phase17StringInSlice(f.Type, []string{"hello", "assign", "chunk", "start", "heartbeat", "log", "result", "cancel", "query"}) || f.Ack != 0 || waiting[d] != 0 {
				return pending(), "P19_INPUT"
			}
			waiting[d] = f.Seq
		}
		next[d]++
	}
	return pending(), "ok"
}

func phase19WorkerChunkOracle(c phase19WorkerChunkCase) string {
	data, final := []byte{}, false
	for _, chunk := range c.Chunks {
		wire := chunk.Data
		if chunk.Repeat > 1 {
			one, err := base64.StdEncoding.Strict().DecodeString(wire)
			if err != nil || len(one) != 1 {
				return "P19_INPUT"
			}
			wire = base64.StdEncoding.EncodeToString(bytes.Repeat(one, chunk.Repeat))
		}
		part, err := base64.StdEncoding.Strict().DecodeString(wire)
		if err != nil || base64.StdEncoding.EncodeToString(part) != wire || final || chunk.Offset != int64(len(data)) || len(part) > 12288 {
			return "P19_INPUT"
		}
		end := int64(len(data)) + int64(len(part))
		if end > c.Total || chunk.Final != (end == c.Total) || (len(part) == 0 && !(c.Total == 0 && chunk.Final && chunk.Offset == 0)) {
			return "P19_INPUT"
		}
		if chunk.Digest != c.Digest {
			return "P19_DIGEST"
		}
		data = append(data, part...)
		final = chunk.Final
	}
	if !final || int64(len(data)) != c.Total {
		return "P19_INPUT"
	}
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != c.Digest {
		return "P19_DIGEST"
	}
	return "ok"
}

func phase19WorkerLeaseOracle(c phase19WorkerLeaseCase) (int64, int64, string) {
	if c.Wall >= c.Until || c.Mono >= c.Deadline || !c.Generation {
		return c.Until, c.Deadline, "P19_LEASE"
	}
	if c.Proposed == 0 || c.Proposed <= c.Until {
		return c.Until, c.Deadline, "ok"
	}
	remaining := c.Proposed - c.Wall
	if remaining <= 0 || remaining > 35000 || !c.Durable {
		return c.Until, c.Deadline, "P19_LEASE"
	}
	if remaining > 30000 {
		remaining = 30000
	}
	return c.Proposed, c.Mono + remaining, "ok"
}

func decodePhase19WorkerSpecification(input []byte) (phase19WorkerSpecificationVector, error) {
	var v phase19WorkerSpecificationVector
	if err := developmentActionsDecodeJSON(input, &v, 65536, true); err != nil {
		return v, err
	}
	if err := phase19SpecificationCaseCoverage(input, phase19SpecificationCaseGroups()["phase-19-worker-vectors"]); err != nil {
		return v, err
	}
	encoded, err := json.Marshal(v)
	var a, b any
	if err != nil || json.Unmarshal(input, &a) != nil || json.Unmarshal(encoded, &b) != nil || !bytes.Equal(iacSpecificationCanonical(a), iacSpecificationCanonical(b)) {
		return v, fmt.Errorf("required fields")
	}
	if v.Version != 1 || v.Acks == nil || v.Chunks == nil || v.Leases == nil {
		return v, fmt.Errorf("vector identity")
	}
	seen := map[string]bool{}
	id := func(s, code string) bool {
		if !regexp.MustCompile("^[a-z][a-z0-9-]*$").MatchString(s) || seen[s] || !phase17StringInSlice(code, []string{"ok", "P19_INPUT", "P19_TRANSPORT", "P19_DIGEST", "P19_LEASE"}) {
			return false
		}
		seen[s] = true
		return true
	}
	for _, c := range v.Acks {
		if !id(c.ID, c.Code) || c.Frames == nil || len(c.Frames) == 0 || len(c.Pending) != 2 || c.Pending[0] < 0 || c.Pending[0] > 1 || c.Pending[1] < 0 || c.Pending[1] > 1 {
			return v, fmt.Errorf("ack case")
		}
	}
	for _, c := range v.Chunks {
		if !id(c.ID, c.Code) || c.Total < 0 || c.Total > 2097152 || !phase17IsLowerHexSHA256(c.Digest) || len(c.Chunks) == 0 {
			return v, fmt.Errorf("chunk case")
		}
		for _, x := range c.Chunks {
			if x.Repeat < 1 || x.Repeat > 12289 || !phase17IsLowerHexSHA256(x.Digest) {
				return v, fmt.Errorf("chunk generator")
			}
		}
	}
	for _, c := range v.Leases {
		if !id(c.ID, c.Code) || c.Mono < 0 || c.Until < 0 || c.Deadline < 0 || c.Proposed < 0 {
			return v, fmt.Errorf("lease case")
		}
	}
	return v, nil
}

func TestPhase19WorkerSpecificationVectors(t *testing.T) {
	t.Parallel()
	doc := phase11MustReadText(t, "docs/details/fixture.md")
	match := regexp.MustCompile("(?s)<a id=\"phase-19-worker-vectors\"></a>.*?\\x60\\x60\\x60json\\n(.*?)\\n\\x60\\x60\\x60").FindStringSubmatch(doc)
	if len(match) != 2 {
		t.Fatal("runner.worker-spec.presence: missing fixed input")
	}
	input := []byte(match[1])
	v, err := decodePhase19WorkerSpecification(input)
	if err != nil || len(v.Acks) != 12 || len(v.Chunks) != 16 || len(v.Leases) != 12 {
		t.Fatalf("runner.worker-spec.schema: fixture.md#phase-19-worker-vectors expected 12/16/12 cases actual %v", err)
	}
	for _, c := range v.Acks {
		t.Run(c.ID, func(t *testing.T) {
			pending, code := phase19WorkerAckOracle(c.Frames)
			if code != c.Code || pending != [2]int{c.Pending[0], c.Pending[1]} {
				t.Fatalf("runner.worker-spec.%s.result: runner.md#phase-19-ipc-framing expected=%s/%v actual=%s/%v; reproduce: go test . -run '^TestPhase19WorkerSpecificationVectors$' -count=1", c.ID, c.Code, c.Pending, code, pending)
			}
		})
	}
	for _, c := range v.Chunks {
		t.Run(c.ID, func(t *testing.T) {
			if code := phase19WorkerChunkOracle(c); code != c.Code {
				t.Fatalf("runner.worker-spec.%s.result: runner.md#phase-19-worker-session expected=%s actual=%s; reproduce: go test . -run '^TestPhase19WorkerSpecificationVectors$' -count=1", c.ID, c.Code, code)
			}
		})
	}
	for _, c := range v.Leases {
		t.Run(c.ID, func(t *testing.T) {
			until, deadline, code := phase19WorkerLeaseOracle(c)
			if code != c.Code || until != c.WantUntil || deadline != c.WantDeadline {
				t.Fatalf("runner.worker-spec.%s.result: runner.md#phase-19-worker-lease-deadlines expected=%s/%d/%d actual=%s/%d/%d; reproduce: go test . -run '^TestPhase19WorkerSpecificationVectors$' -count=1", c.ID, c.Code, c.WantUntil, c.WantDeadline, code, until, deadline)
			}
		})
	}
	var nullValue map[string]json.RawMessage
	if err := json.Unmarshal(input, &nullValue); err != nil {
		t.Fatal(err)
	}
	nullValue["ack_cases"] = json.RawMessage("null")
	nullInput, err := json.Marshal(nullValue)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		id    string
		input []byte
	}{
		{"missing", []byte(strings.Replace(string(input), "\"schema_version\": 1,", "", 1))},
		{"unknown", []byte(strings.Replace(string(input), "\"schema_version\": 1,", "\"schema_version\": 1,\"unknown\":true,", 1))},
		{"duplicate", []byte(strings.Replace(string(input), "\"schema_version\": 1,", "\"schema_version\": 1,\"schema_version\": 1,", 1))},
		{"null", nullInput},
		{"trailing", append(append([]byte{}, input...), []byte("{}")...)},
		{"invalid-utf8", append([]byte{0xff}, input...)},
	} {
		t.Run("schema-"+c.id, func(t *testing.T) {
			if _, err := decodePhase19WorkerSpecification(c.input); err == nil {
				t.Fatalf("runner.worker-spec.schema-%s.reject: fixture.md#phase-19-worker-vectors expected rejected actual accepted", c.id)
			}
		})
	}
}

func phase19LogRingOracle(data []byte, first int64) ([]byte, int64, int64) {
	drop := len(data) - 1048576
	if drop < 0 {
		drop = 0
	}
	for drop < len(data) && !utf8.RuneStart(data[drop]) {
		drop++
	}
	return data[drop:], first + int64(drop), first + int64(len(data))
}

type phase19ContentChunk struct {
	Offset int64  `json:"offset"`
	Data   string `json:"data"`
	Next   int64  `json:"next_offset"`
	Total  int64  `json:"total_bytes"`
	Digest string `json:"digest"`
	Final  bool   `json:"final"`
}

type phase19CacheKeyInput struct {
	Project   string            `json:"project_id"`
	Trust     string            `json:"trust"`
	Plan      string            `json:"plan_digest"`
	Job       string            `json:"job_id"`
	Source    string            `json:"source_tree_digest"`
	Toolchain string            `json:"toolchain_digest"`
	OS        string            `json:"os"`
	Arch      string            `json:"arch"`
	Env       map[string]string `json:"env"`
}

type phase19ContentSpecificationVector struct {
	SchemaVersion int `json:"schema_version"`
	Chunks        []struct {
		ID     string `json:"id"`
		Total  int64  `json:"total_bytes"`
		Offset int64  `json:"offset"`
		Bytes  int64  `json:"expected_bytes"`
		Next   int64  `json:"expected_next"`
		Final  bool   `json:"expected_final"`
		Code   string `json:"expected_code"`
	} `json:"chunk_cases"`
	Queries []struct {
		ID     string `json:"id"`
		Offset string `json:"offset"`
		Code   string `json:"expected_code"`
	} `json:"query_cases"`
	Wires []struct {
		ID     string              `json:"id"`
		Offset int64               `json:"requested_offset"`
		Digest string              `json:"requested_digest"`
		Wire   phase19ContentChunk `json:"wire"`
		Code   string              `json:"expected_code"`
	} `json:"wire_cases"`
	Caches []struct {
		ID          string `json:"id"`
		Metadata    bool   `json:"metadata"`
		MatchingKey bool   `json:"matching_key"`
		Published   bool   `json:"published"`
		Deleting    bool   `json:"deleting"`
		Quota       int    `json:"quota_mb"`
		Now         int64  `json:"now"`
		Expires     int64  `json:"expires_at"`
		Verified    bool   `json:"verified"`
		Expected    string `json:"expected"`
	} `json:"cache_cases"`
	Keys []struct {
		ID       string               `json:"id"`
		Input    phase19CacheKeyInput `json:"input"`
		Expected string               `json:"expected_sha256"`
	} `json:"key_cases"`
	Envelopes []struct {
		ID       string  `json:"id"`
		Sizes    []int64 `json:"sizes"`
		Mutation string  `json:"mutation"`
		Code     string  `json:"expected_code"`
	} `json:"envelope_cases"`
}

func decodePhase19ContentSpecification(input []byte) (phase19ContentSpecificationVector, error) {
	var v phase19ContentSpecificationVector
	if err := developmentActionsDecodeJSON(input, &v, 65536, true); err != nil {
		return v, err
	}
	if err := phase19SpecificationCaseCoverage(input, phase19SpecificationCaseGroups()["phase-19-content-vectors"]); err != nil {
		return v, err
	}
	encoded, err := json.Marshal(v)
	var a, b any
	if err != nil || json.Unmarshal(input, &a) != nil || json.Unmarshal(encoded, &b) != nil ||
		!bytes.Equal(iacSpecificationCanonical(a), iacSpecificationCanonical(b)) {
		return v, fmt.Errorf("required fields and null boundary")
	}
	if v.SchemaVersion != 1 || v.Chunks == nil || v.Queries == nil || v.Wires == nil || v.Caches == nil || v.Keys == nil || v.Envelopes == nil {
		return v, fmt.Errorf("vector identity")
	}
	seen := map[string]bool{}
	id := func(s string) bool {
		if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(s) || seen[s] {
			return false
		}
		seen[s] = true
		return true
	}
	validCode := func(s string) bool {
		return phase17StringInSlice(s, []string{"ok", "P19_INPUT", "P19_SCHEMA", "P19_STATE_CORRUPT", "P19_ARCHIVE", "protocol"})
	}
	for _, c := range v.Chunks {
		if !id(c.ID) || !validCode(c.Code) {
			return v, fmt.Errorf("chunk schema")
		}
	}
	for _, c := range v.Queries {
		if !id(c.ID) || !validCode(c.Code) {
			return v, fmt.Errorf("query schema")
		}
	}
	for _, c := range v.Wires {
		if !id(c.ID) || !validCode(c.Code) || !phase17IsLowerHexSHA256(c.Digest) || !phase17IsLowerHexSHA256(c.Wire.Digest) {
			return v, fmt.Errorf("wire schema")
		}
	}
	for _, c := range v.Caches {
		if !id(c.ID) || c.Quota < 0 || c.Now < 0 || c.Expires < 0 || !phase17StringInSlice(c.Expected, []string{"hit", "miss", "invalidate"}) {
			return v, fmt.Errorf("cache schema")
		}
	}
	for _, c := range v.Keys {
		k := c.Input
		if !id(c.ID) || !phase17IsLowerHexSHA256(c.Expected) || k.Env == nil || k.Project == "" || k.Job == "" || k.Trust != "registered" || k.OS != "linux" || k.Arch != "amd64" ||
			!phase17IsLowerHexSHA256(k.Plan) || !phase17IsLowerHexSHA256(k.Source) || !phase17IsLowerHexSHA256(k.Toolchain) {
			return v, fmt.Errorf("key schema")
		}
	}
	for _, c := range v.Envelopes {
		if !id(c.ID) || c.Sizes == nil || len(c.Sizes) < 1 || len(c.Sizes) > 2 || !validCode(c.Code) ||
			!phase17StringInSlice(c.Mutation, []string{"none", "gzip-empty-member", "gzip-trailing-zero", "gzip-trailing-byte", "gzip-crc", "gzip-truncated", "tar-extra-zero", "tar-short-terminator", "tar-data-padding", "tar-pax-prefix"}) {
			return v, fmt.Errorf("envelope schema")
		}
		for _, size := range c.Sizes {
			if size < 0 || size > 1024 {
				return v, fmt.Errorf("synthetic file size")
			}
		}
	}
	return v, nil
}

func phase19ChunkRangeOracle(total, offset int64) (int64, int64, bool, string) {
	if total < 1 || total > 1073741824 {
		return 0, 0, false, "P19_STATE_CORRUPT"
	}
	if offset < 0 || offset > total {
		return 0, 0, false, "P19_SCHEMA"
	}
	n := total - offset
	if n > 12288 {
		n = 12288
	}
	next := offset + n
	return n, next, next == total, "ok"
}

func phase19ChunkQueryOracle(offset string) string {
	if !regexp.MustCompile(`^(0|[1-9][0-9]*)$`).MatchString(offset) {
		return "P19_INPUT"
	}
	n, err := strconv.ParseUint(offset, 10, 64)
	if err != nil || n > 1073741824 {
		return "P19_SCHEMA"
	}
	return "ok"
}

func phase19ChunkWireOracle(offset int64, digest string, c phase19ContentChunk) string {
	n, next, final, code := phase19ChunkRangeOracle(c.Total, offset)
	data, err := base64.StdEncoding.Strict().DecodeString(c.Data)
	if code != "ok" || c.Offset != offset || c.Digest != digest || err != nil || base64.StdEncoding.EncodeToString(data) != c.Data ||
		int64(len(data)) != n || c.Next != next || c.Final != final {
		return "protocol"
	}
	return "ok"
}

func phase19CacheHitOracle(metadata, matchingKey, published, deleting bool, quota int, now, expires int64, verified bool) string {
	if !metadata || !matchingKey || deleting || quota == 0 || now >= expires {
		return "miss"
	}
	if !published || !verified {
		return "invalidate"
	}
	return "hit"
}

func phase19CacheKeyOracle(input phase19CacheKeyInput) (string, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	var value any
	if err := json.Unmarshal(b, &value); err != nil {
		return "", err
	}
	h := sha256.Sum256(iacSpecificationCanonical(value))
	return hex.EncodeToString(h[:]), nil
}

func phase19ContentEnvelopeOracle(payload []byte, sizes []int64) string {
	expected := int64(1024)
	for _, size := range sizes {
		expected += 512 + ((size+511)/512)*512
	}
	if !bytes.HasPrefix(payload, []byte{0x1f, 0x8b, 8, 0, 0, 0, 0, 0, 2, 255}) {
		return "P19_ARCHIVE"
	}
	compressed := bytes.NewReader(payload)
	gr, err := gzip.NewReader(compressed)
	if err != nil {
		return "P19_ARCHIVE"
	}
	gr.Multistream(false)
	raw, err := io.ReadAll(io.LimitReader(gr, expected+1))
	closeErr := gr.Close()
	if err != nil || closeErr != nil || int64(len(raw)) != expected || compressed.Len() != 0 {
		return "P19_ARCHIVE"
	}
	pos := int64(0)
	for _, size := range sizes {
		header := raw[pos : pos+512]
		if header[156] != '0' || !bytes.Equal(header[257:265], []byte{'u', 's', 't', 'a', 'r', 0, '0', '0'}) {
			return "P19_ARCHIVE"
		}
		end := pos + 512 + size
		padded := pos + 512 + ((size+511)/512)*512
		if !bytes.Equal(raw[end:padded], make([]byte, padded-end)) {
			return "P19_ARCHIVE"
		}
		pos = padded
	}
	if !bytes.Equal(raw[pos:], make([]byte, 1024)) {
		return "P19_ARCHIVE"
	}
	tr := tar.NewReader(bytes.NewReader(raw))
	for i, size := range sizes {
		h, err := tr.Next()
		if err != nil || h.Name != string(rune('a'+i)) || h.Typeflag != tar.TypeReg || h.Size != size || h.Mode != 0644 || h.Uid != 0 || h.Gid != 0 || h.Uname != "" || h.Gname != "" || h.ModTime.Unix() != 0 {
			return "P19_ARCHIVE"
		}
		data, err := io.ReadAll(tr)
		if err != nil || !bytes.Equal(data, bytes.Repeat([]byte{'a'}, int(size))) {
			return "P19_ARCHIVE"
		}
	}
	if _, err := tr.Next(); err != io.EOF {
		return "P19_ARCHIVE"
	}
	return "ok"
}

func TestPhase19ContentSpecificationVectors(t *testing.T) {
	t.Parallel()
	doc := phase11MustReadText(t, "docs/details/fixture.md")
	match := regexp.MustCompile("(?s)<a id=\"phase-19-content-vectors\"></a>.*?\x60\x60\x60json\\n(.*?)\\n\x60\x60\x60").FindStringSubmatch(doc)
	if len(match) != 2 {
		t.Fatal("archive.content-spec.presence: missing fixed input")
	}
	input := []byte(match[1])
	v, err := decodePhase19ContentSpecification(input)
	if err != nil || len(v.Chunks) != 11 || len(v.Queries) != 8 || len(v.Wires) != 12 || len(v.Caches) != 9 || len(v.Keys) != 3 || len(v.Envelopes) != 13 {
		t.Fatalf("archive.content-spec.schema: expected 11/8/12/9/3/13 valid vectors actual %v", err)
	}
	for _, c := range v.Chunks {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			n, next, final, code := phase19ChunkRangeOracle(c.Total, c.Offset)
			if n != c.Bytes || next != c.Next || final != c.Final || code != c.Code {
				t.Fatalf("api.content-spec.%s.result: spec=docs/details/api.md#phase-19-artifact-download expected=%d/%d/%t/%s actual=%d/%d/%t/%s", c.ID, c.Bytes, c.Next, c.Final, c.Code, n, next, final, code)
			}
		})
	}
	for _, c := range v.Queries {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			if got := phase19ChunkQueryOracle(c.Offset); got != c.Code {
				t.Fatalf("api.content-spec.%s.result: spec=docs/details/api.md#phase-19-artifact-download expected=%s actual=%s", c.ID, c.Code, got)
			}
		})
	}
	for _, c := range v.Wires {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			if got := phase19ChunkWireOracle(c.Offset, c.Digest, c.Wire); got != c.Code {
				t.Fatalf("api.content-spec.%s.result: spec=docs/details/api.md#phase-19-artifact-download expected=%s actual=%s", c.ID, c.Code, got)
			}
		})
	}
	for _, c := range v.Caches {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			if got := phase19CacheHitOracle(c.Metadata, c.MatchingKey, c.Published, c.Deleting, c.Quota, c.Now, c.Expires, c.Verified); got != c.Expected {
				t.Fatalf("runner.content-spec.%s.result: spec=docs/details/runner.md#phase-19-cache-lifecycle expected=%s actual=%s", c.ID, c.Expected, got)
			}
		})
	}
	for _, c := range v.Keys {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			if got, err := phase19CacheKeyOracle(c.Input); err != nil || got != c.Expected {
				t.Fatalf("runner.content-spec.%s.result: spec=docs/details/runner.md#phase-19-cache-lifecycle expected=%s actual=%s error=%v", c.ID, c.Expected, got, err)
			}
		})
	}
	gzipBytes := func(raw []byte) []byte {
		var out bytes.Buffer
		gw, err := gzip.NewWriterLevel(&out, gzip.BestCompression)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = gw.Write(raw); err != nil {
			t.Fatal(err)
		}
		if err = gw.Close(); err != nil {
			t.Fatal(err)
		}
		return out.Bytes()
	}
	for _, c := range v.Envelopes {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			var buf bytes.Buffer
			tw := tar.NewWriter(&buf)
			for i, size := range c.Sizes {
				header := &tar.Header{Name: string(rune('a' + i)), Mode: 0644, Size: size, Typeflag: tar.TypeReg, ModTime: time.Unix(0, 0).UTC(), Format: tar.FormatUSTAR}
				if c.Mutation == "tar-pax-prefix" && i == 0 {
					header.Format = tar.FormatPAX
					header.PAXRecords = map[string]string{"comment": "x"}
				}
				if err := tw.WriteHeader(header); err != nil {
					t.Fatal(err)
				}
				if _, err := tw.Write(bytes.Repeat([]byte{'a'}, int(size))); err != nil {
					t.Fatal(err)
				}
			}
			if err := tw.Close(); err != nil {
				t.Fatal(err)
			}
			raw := append([]byte{}, buf.Bytes()...)
			switch c.Mutation {
			case "tar-extra-zero":
				raw = append(raw, make([]byte, 512)...)
			case "tar-short-terminator":
				raw = raw[:len(raw)-512]
			case "tar-data-padding":
				raw[512+c.Sizes[0]] = 1
			}
			payload := gzipBytes(raw)
			switch c.Mutation {
			case "gzip-empty-member":
				payload = append(payload, gzipBytes(nil)...)
			case "gzip-trailing-zero":
				payload = append(payload, 0)
			case "gzip-trailing-byte":
				payload = append(payload, 1)
			case "gzip-crc":
				payload[len(payload)-8] ^= 1
			case "gzip-truncated":
				payload = payload[:len(payload)-1]
			}
			if got := phase19ContentEnvelopeOracle(payload, c.Sizes); got != c.Code {
				t.Fatalf("archive.content-spec.%s.result: spec=docs/details/archive.md#phase-19-archive-envelope expected=%s actual=%s", c.ID, c.Code, got)
			}
		})
	}
	for _, bad := range []struct {
		id    string
		input []byte
	}{
		{"missing", []byte(strings.Replace(string(input), `"schema_version": 1,`, "", 1))},
		{"unknown", []byte(strings.Replace(string(input), `"schema_version": 1,`, `"schema_version": 1, "unknown": true,`, 1))},
		{"duplicate", []byte(strings.Replace(string(input), `"schema_version": 1,`, `"schema_version": 1, "schema_version": 1,`, 1))},
		{"null", []byte(strings.Replace(string(input), `"metadata": true`, `"metadata": null`, 1))},
		{"trailing", append(append([]byte{}, input...), []byte("{}")...)},
		{"utf8", append(append([]byte{}, input...), 0xff)},
	} {
		t.Run("schema-"+bad.id, func(t *testing.T) {
			if _, err := decodePhase19ContentSpecification(bad.input); err == nil {
				t.Fatalf("statefile.content-spec.schema-%s.reject: expected rejected actual accepted", bad.id)
			}
		})
	}
}

type phase19OperationSpecificationVector struct {
	SchemaVersion int `json:"schema_version"`
	Admission     []struct {
		ID              string `json:"id"`
		Authorized      bool   `json:"authorized"`
		Existing        bool   `json:"existing"`
		SameSubject     bool   `json:"same_subject"`
		SameDigest      bool   `json:"same_digest"`
		RevisionMatches bool   `json:"revision_matches"`
		Expected        string `json:"expected"`
	} `json:"admission_cases"`
	Transitions []struct {
		ID         string `json:"id"`
		From       string `json:"from"`
		To         string `json:"to"`
		SameResult bool   `json:"same_result"`
		Expected   string `json:"expected"`
	} `json:"transition_cases"`
	GC []struct {
		ID        string `json:"id"`
		State     string `json:"state"`
		Now       int64  `json:"now"`
		ExpiresAt int64  `json:"expires_at"`
		Protected bool   `json:"protected_reference"`
		Expected  bool   `json:"expected_collect"`
	} `json:"gc_cases"`
}

func decodePhase19OperationSpecification(input []byte) (phase19OperationSpecificationVector, error) {
	var v phase19OperationSpecificationVector
	if err := developmentActionsDecodeJSON(input, &v, 65536, true); err != nil {
		return v, err
	}
	if err := phase19SpecificationCaseCoverage(input, phase19SpecificationCaseGroups()["phase-19-operation-vectors"]); err != nil {
		return v, err
	}
	encoded, err := json.Marshal(v)
	var before, after any
	if err != nil || json.Unmarshal(input, &before) != nil || json.Unmarshal(encoded, &after) != nil ||
		!bytes.Equal(iacSpecificationCanonical(before), iacSpecificationCanonical(after)) {
		return v, fmt.Errorf("required fields and null boundary")
	}
	if v.SchemaVersion != 1 || v.Admission == nil || v.Transitions == nil || v.GC == nil {
		return v, fmt.Errorf("vector identity")
	}
	seen := map[string]bool{}
	validID := func(id string) bool {
		if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(id) || seen[id] {
			return false
		}
		seen[id] = true
		return true
	}
	validState := func(state string) bool {
		return state == "accepted" || state == "succeeded" || state == "failed" || state == "recovery_required"
	}
	for _, c := range v.Admission {
		if !validID(c.ID) || !phase17StringInSlice(c.Expected, []string{"create", "forbidden", "replay", "idempotency", "revision"}) {
			return v, fmt.Errorf("admission schema")
		}
	}
	for _, c := range v.Transitions {
		if !validID(c.ID) || !validState(c.From) || !validState(c.To) ||
			!phase17StringInSlice(c.Expected, []string{"advance", "unchanged", "stale"}) {
			return v, fmt.Errorf("transition schema")
		}
	}
	for _, c := range v.GC {
		if !validID(c.ID) || !validState(c.State) || c.Now < 0 || c.ExpiresAt != 2592000 {
			return v, fmt.Errorf("gc schema")
		}
	}
	return v, nil
}

func phase19OperationAdmissionOracle(authorized, existing, sameSubject, sameDigest, revisionMatches bool) string {
	if !authorized {
		return "forbidden"
	}
	if existing {
		if !sameSubject || !sameDigest {
			return "idempotency"
		}
		return "replay"
	}
	if !revisionMatches {
		return "revision"
	}
	return "create"
}

func phase19OperationTransitionOracle(from, to string, sameResult bool) string {
	if from == to && sameResult {
		return "unchanged"
	}
	if (from == "accepted" && (to == "succeeded" || to == "failed" || to == "recovery_required")) ||
		(from == "recovery_required" && (to == "succeeded" || to == "failed")) {
		return "advance"
	}
	return "stale"
}

func phase19OperationGCOracle(state string, now, expiresAt int64, protected bool) bool {
	return (state == "succeeded" || state == "failed") && now >= expiresAt && !protected
}

func TestPhase19OperationSpecificationVectors(t *testing.T) {
	t.Parallel()
	document := phase11MustReadText(t, "docs/details/fixture.md")
	match := regexp.MustCompile("(?s)<a id=\"phase-19-operation-vectors\"></a>.*?\x60\x60\x60json\\n(.*?)\\n\x60\x60\x60").FindStringSubmatch(document)
	if len(match) != 2 {
		t.Fatal("statefile.operation-spec.presence: missing fixed input")
	}
	input := []byte(match[1])
	v, err := decodePhase19OperationSpecification(input)
	if err != nil || len(v.Admission) != 8 || len(v.Transitions) != 12 || len(v.GC) != 8 {
		t.Fatalf("statefile.operation-spec.schema: expected 8/12/8 valid vectors, actual %v", err)
	}
	for _, c := range v.Admission {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			got := phase19OperationAdmissionOracle(c.Authorized, c.Existing, c.SameSubject, c.SameDigest, c.RevisionMatches)
			if got != c.Expected {
				t.Fatalf("api.operation-spec.%s.result: spec=docs/details/api.md#phase-19-operation-admission expected=%s actual=%s input=%+v", c.ID, c.Expected, got, c)
			}
		})
	}
	for _, c := range v.Transitions {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			got := phase19OperationTransitionOracle(c.From, c.To, c.SameResult)
			if got != c.Expected {
				t.Fatalf("statefile.operation-spec.%s.result: spec=docs/details/statefile.md#phase-19-operation-lifecycle expected=%s actual=%s input=%+v", c.ID, c.Expected, got, c)
			}
		})
	}
	for _, c := range v.GC {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			got := phase19OperationGCOracle(c.State, c.Now, c.ExpiresAt, c.Protected)
			if got != c.Expected {
				t.Fatalf("statefile.operation-spec.%s.result: spec=docs/details/statefile.md#phase-19-operation-lifecycle expected=%t actual=%t input=%+v", c.ID, c.Expected, got, c)
			}
		})
	}
	for _, bad := range []struct {
		id    string
		input []byte
	}{
		{"missing", []byte(strings.Replace(string(input), `"schema_version": 1,`, "", 1))},
		{"unknown", []byte(strings.Replace(string(input), `"schema_version": 1,`, `"schema_version": 1, "unknown": true,`, 1))},
		{"duplicate", []byte(strings.Replace(string(input), `"schema_version": 1,`, `"schema_version": 1, "schema_version": 1,`, 1))},
		{"null", []byte(strings.Replace(string(input), `"authorized": true`, `"authorized": null`, 1))},
		{"trailing", append(append([]byte{}, input...), []byte("{}")...)},
		{"utf8", append(append([]byte{}, input...), 0xff)},
	} {
		t.Run("schema-"+bad.id, func(t *testing.T) {
			if _, err := decodePhase19OperationSpecification(bad.input); err == nil {
				t.Fatalf("statefile.operation-spec.schema-%s.reject: expected rejected actual accepted", bad.id)
			}
		})
	}
	api := phase11MustReadText(t, "docs/details/api.md")
	if err := phase19OperationHTTPPartition(api); err != nil {
		t.Fatalf("api.operation-spec.http-partition.result: %v", err)
	}
	for _, bad := range []struct{ id, from, to string }{
		{"wrong-status", "| 202: workerResume/", "| 200: workerResume/"},
		{"missing-kind", "| 200: passwordChange/projectPut/", "| 200: projectPut/"},
		{"duplicate-kind", "| 202: workerResume/", "| 202: projectPut/workerResume/"},
		{"read-kind", "| 202: workerResume/", "| 202: plan/workerResume/"},
	} {
		t.Run("http-"+bad.id, func(t *testing.T) {
			if !strings.Contains(api, bad.from) {
				t.Fatalf("api.operation-spec.http-%s.fixture: missing control target", bad.id)
			}
			if err := phase19OperationHTTPPartition(strings.Replace(api, bad.from, bad.to, 1)); err == nil {
				t.Fatalf("api.operation-spec.http-%s.reject: expected rejected actual accepted", bad.id)
			}
		})
	}
}

func phase19OperationHTTPPartition(api string) error {
	section := strings.Split(api, "<a id=\"phase-19-error-catalog\"></a>")[0]
	classes := regexp.MustCompile("(?m)^\\| (200|202): ([A-Za-z/]+) \\|").FindAllStringSubmatch(section, -1)
	classified := map[string]string{}
	for _, row := range classes {
		for _, kind := range strings.Split(row[2], "/") {
			if classified[kind] != "" {
				return fmt.Errorf("duplicate kind %s", kind)
			}
			classified[kind] = row[1]
		}
	}
	catalog := regexp.MustCompile("(?m)^\\| ([A-Za-z][A-Za-z0-9]*) \\| (GET|POST|PUT|DELETE) [^|]+ \\| [^|]+ \\| ([^|]+) \\|").FindAllStringSubmatch(section, -1)
	mutationCount := 0
	for _, row := range catalog {
		kind := row[1]
		if row[2] == "GET" || phase17StringInSlice(kind, []string{"login", "logout", "plan", "deploymentPlan", "artifactCompare"}) {
			if classified[kind] != "" {
				return fmt.Errorf("read/auth kind %s", kind)
			}
			continue
		}
		want := "200"
		if strings.TrimSpace(row[3]) == "Operation" {
			want = "202"
		}
		if classified[kind] != want {
			return fmt.Errorf("kind=%s expected=%s actual=%s", kind, want, classified[kind])
		}
		mutationCount++
	}
	if mutationCount != 25 || len(classified) != mutationCount {
		return fmt.Errorf("expected 25 kinds actual catalog=%d classified=%d", mutationCount, len(classified))
	}
	return nil
}

type phase19LogSpecificationVector struct {
	SchemaVersion int `json:"schema_version"`
	MaskCases     []struct {
		ID        string   `json:"id"`
		Patterns  []string `json:"patterns_hex"`
		Chunks    []string `json:"chunks_hex"`
		Committed []string `json:"committed_hex"`
		Expected  string   `json:"expected_hex"`
	} `json:"mask_cases"`
	PageCases []struct {
		ID               string `json:"id"`
		Data             string `json:"data_hex"`
		First            int64  `json:"first_offset"`
		Cursor           *int64 `json:"cursor_offset"`
		Limit            int    `json:"limit_bytes"`
		Terminal         bool   `json:"stored_terminal"`
		Gap              bool   `json:"stored_gap"`
		Code             string `json:"expected_code"`
		Expected         string `json:"expected_hex"`
		ExpectedFirst    int64  `json:"expected_first"`
		ExpectedNext     int64  `json:"expected_next"`
		ExpectedTerminal bool   `json:"expected_terminal"`
		Unchanged        bool   `json:"expected_cursor_unchanged"`
	} `json:"page_cases"`
	BudgetCases []struct {
		ID       string `json:"id"`
		OutputMB int64  `json:"output_mb"`
		Stdout   int64  `json:"stdout_bytes"`
		Stderr   int64  `json:"stderr_bytes"`
		Code     string `json:"expected_code"`
	} `json:"budget_cases"`
}

func decodePhase19LogSpecification(input []byte) (phase19LogSpecificationVector, error) {
	var v phase19LogSpecificationVector
	if err := developmentActionsDecodeJSON(input, &v, 65536, true); err != nil {
		return v, err
	}
	if err := phase19SpecificationCaseCoverage(input, phase19SpecificationCaseGroups()["phase-19-log-boundary-vectors"]); err != nil {
		return v, err
	}
	encoded, err := json.Marshal(v)
	if err != nil {
		return v, err
	}
	var before, after any
	if json.Unmarshal(input, &before) != nil || json.Unmarshal(encoded, &after) != nil ||
		!bytes.Equal(iacSpecificationCanonical(before), iacSpecificationCanonical(after)) {
		return v, fmt.Errorf("required fields and null boundary")
	}
	if v.SchemaVersion != 1 || v.MaskCases == nil || v.PageCases == nil || v.BudgetCases == nil {
		return v, fmt.Errorf("vector identity")
	}
	validHex := func(s string) bool {
		_, err := hex.DecodeString(s)
		return err == nil && strings.ToLower(s) == s
	}
	seen := map[string]bool{}
	id := func(s string) bool {
		if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(s) || seen[s] {
			return false
		}
		seen[s] = true
		return true
	}
	for _, c := range v.MaskCases {
		if !id(c.ID) || c.Patterns == nil || c.Chunks == nil || c.Committed == nil ||
			len(c.Chunks) != len(c.Committed) || !validHex(c.Expected) {
			return v, fmt.Errorf("mask schema")
		}
		for _, p := range c.Patterns {
			if p == "" || !validHex(p) {
				return v, fmt.Errorf("pattern schema")
			}
		}
		for _, s := range append(append([]string{}, c.Chunks...), c.Committed...) {
			if !validHex(s) {
				return v, fmt.Errorf("chunk schema")
			}
		}
	}
	for _, c := range v.PageCases {
		if !id(c.ID) || !validHex(c.Data) || !validHex(c.Expected) {
			return v, fmt.Errorf("page schema")
		}
	}
	for _, c := range v.BudgetCases {
		if !id(c.ID) || c.OutputMB < 1 || c.OutputMB > 16 || c.Stdout < 0 || c.Stderr < 0 {
			return v, fmt.Errorf("budget schema")
		}
	}
	return v, nil
}

// Pure published-input oracles; these do not run a product worker or read its logs.
func phase19LogMaskOracle(patterns, chunks [][]byte) ([]byte, []string) {
	sorted := append([][]byte(nil), patterns...)
	sort.Slice(sorted, func(i, j int) bool {
		if len(sorted[i]) != len(sorted[j]) {
			return len(sorted[i]) > len(sorted[j])
		}
		return bytes.Compare(sorted[i], sorted[j]) < 0
	})
	longest := 1
	for _, p := range sorted {
		if len(p) > longest {
			longest = len(p)
		}
	}
	var pending, carry, output []byte
	commit := func(final bool) {
		for len(pending) > 0 && (final || len(pending) >= longest) {
			count := 0
			for _, p := range sorted {
				if bytes.HasPrefix(pending, p) {
					carry = append(carry, []byte("[REDACTED]")...)
					count = len(p)
					break
				}
			}
			if count == 0 {
				carry = append(carry, pending[0])
				count = 1
			}
			pending = pending[count:]
		}
		for len(carry) > 0 {
			if !final && !utf8.FullRune(carry) {
				break
			}
			r, size := utf8.DecodeRune(carry)
			if r == utf8.RuneError && size == 1 {
				output = append(output, []byte("\uFFFD")...)
			} else {
				output = append(output, carry[:size]...)
			}
			carry = carry[size:]
		}
	}
	var committed []string
	for _, chunk := range chunks {
		pending = append(pending, chunk...)
		commit(false)
		committed = append(committed, hex.EncodeToString(output))
	}
	commit(true)
	return output, committed
}

func phase19LogPageOracle(data []byte, first int64, cursor *int64, limit int, terminal, gap bool) (string, []byte, int64, int64, bool, bool, bool, bool) {
	if limit < 1 || limit > 65536 {
		return "P19_SCHEMA", nil, 0, 0, false, false, false, false
	}
	if !utf8.Valid(data) || first < 0 || first > 9007199254740991-int64(len(data)) {
		return "P19_STATE_CORRUPT", nil, 0, 0, false, false, false, false
	}
	start := first
	if cursor != nil {
		start = *cursor
	}
	last := first + int64(len(data))
	if start < first || start > last {
		return "P19_CURSOR", nil, 0, 0, false, false, false, false
	}
	index := int(start - first)
	if index < len(data) && !utf8.RuneStart(data[index]) {
		return "P19_CURSOR", nil, 0, 0, false, false, false, false
	}
	end := index
	for end < len(data) {
		_, size := utf8.DecodeRune(data[end:])
		if end-index+size > limit {
			break
		}
		end += size
	}
	if end == index && index < len(data) {
		return "P19_SCHEMA", nil, 0, 0, false, false, false, false
	}
	return "ok", data[index:end], start, first + int64(end), terminal && end == len(data), cursor != nil && end == index, first > 0, gap
}

func phase19LogBudgetOracle(outputMB, stdout, stderr int64) string {
	if stdout+stderr > outputMB*1048576 {
		return "P19_OUTPUT"
	}
	return "ok"
}

func TestPhase19LogSpecificationVectors(t *testing.T) {
	t.Parallel()
	document := phase11MustReadText(t, "docs/details/fixture.md")
	match := regexp.MustCompile("(?s)<a id=\"phase-19-log-boundary-vectors\"></a>.*?\x60\x60\x60json\\n(.*?)\\n\x60\x60\x60").FindStringSubmatch(document)
	if len(match) != 2 {
		t.Fatal("runner.log-spec.presence: missing fixed input")
	}
	v, err := decodePhase19LogSpecification([]byte(match[1]))
	if err != nil || len(v.MaskCases) != 12 || len(v.PageCases) != 13 || len(v.BudgetCases) != 4 {
		t.Fatalf("runner.log-spec.schema: expected 12/13/4 valid cases, actual %v", err)
	}
	decode := func(s string) []byte {
		b, err := hex.DecodeString(s)
		if err != nil {
			t.Fatal("fixture hex")
		}
		return b
	}
	for _, c := range v.MaskCases {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			var patterns, chunks [][]byte
			for _, p := range c.Patterns {
				patterns = append(patterns, decode(p))
			}
			for _, b := range c.Chunks {
				chunks = append(chunks, decode(b))
			}
			got, commits := phase19LogMaskOracle(patterns, chunks)
			if hex.EncodeToString(got) != c.Expected || !phase11StringSlicesEqual(commits, c.Committed) {
				t.Fatalf("runner.log-spec.%s.result: expected published committed/final bytes, actual differs; runner.md#phase-19-streaming-log-contract; reproduce: go test . -run '^TestPhase19LogSpecificationVectors$' -count=1", c.ID)
			}
			joined := bytes.Join(chunks, nil)
			for split := 0; split <= len(joined); split++ {
				got, _ := phase19LogMaskOracle(patterns, [][]byte{joined[:split], joined[split:]})
				if hex.EncodeToString(got) != c.Expected {
					t.Fatalf("runner.log-spec.%s.result: split=%d changed final oracle", c.ID, split)
				}
			}
			var singles [][]byte
			for i := range joined {
				singles = append(singles, joined[i:i+1])
			}
			got, _ = phase19LogMaskOracle(patterns, singles)
			if hex.EncodeToString(got) != c.Expected {
				t.Fatalf("runner.log-spec.%s.result: byte chunks changed final oracle", c.ID)
			}
		})
	}
	for _, c := range v.PageCases {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			code, data, first, next, terminal, unchanged, truncated, gap := phase19LogPageOracle(decode(c.Data), c.First, c.Cursor, c.Limit, c.Terminal, c.Gap)
			if code != c.Code || hex.EncodeToString(data) != c.Expected || first != c.ExpectedFirst || next != c.ExpectedNext || terminal != c.ExpectedTerminal || unchanged != c.Unchanged || truncated != (code == "ok" && c.First > 0) || gap != (code == "ok" && c.Gap) {
				t.Fatalf("api.log-spec.%s.result: expected code=%s offsets=%d..%d terminal=%t, actual code=%s offsets=%d..%d terminal=%t; api.md#phase-19-log-page-contract; reproduce: go test . -run '^TestPhase19LogSpecificationVectors$' -count=1", c.ID, c.Code, c.ExpectedFirst, c.ExpectedNext, c.ExpectedTerminal, code, first, next, terminal)
			}
		})
	}
	for _, c := range v.BudgetCases {
		if got := phase19LogBudgetOracle(c.OutputMB, c.Stdout, c.Stderr); got != c.Code {
			t.Fatalf("runner.log-spec.%s.result: expected %s, actual %s; runner.md#phase-19-streaming-log-contract", c.ID, c.Code, got)
		}
	}
	for _, c := range []struct {
		id          string
		data        []byte
		first, next int64
		size        int
	}{
		{"ring-exact", bytes.Repeat([]byte("a"), 1048576), 0, 1048576, 1048576},
		{"ring-scalar-prefix", append([]byte{0xe3, 0x81, 0x82}, bytes.Repeat([]byte("a"), 1048575)...), 3, 1048578, 1048575},
		{"ring-unicode", bytes.Repeat([]byte{0xe3, 0x81, 0x82}, 349526), 3, 1048578, 1048575},
	} {
		data, first, next := phase19LogRingOracle(c.data, 0)
		if !utf8.Valid(data) || len(data) != c.size || first != c.first || next != c.next {
			t.Fatalf("runner.log-spec.%s.result: expected ring bytes=%d offsets=%d..%d, actual bytes=%d offsets=%d..%d; runner.md#phase-19-streaming-log-contract", c.id, c.size, c.first, c.next, len(data), first, next)
		}
	}
	log := map[string]any{"schema_version": 1, "project_id": strings.Repeat("p", 64), "run_id": "r-" + strings.Repeat("a", 32), "attempt": 100, "job_id": strings.Repeat("j", 64), "stream": "stdout", "first_offset": int64(9007199254740991 - 1048576), "next_offset": int64(9007199254740991), "data": strings.Repeat("\x00", 1048576), "truncated": true, "gap": true, "terminal": true}
	encoded, err := json.Marshal(log)
	state := phase11MustReadText(t, "docs/details/statefile.md")
	capacity := regexp.MustCompile("\\| Log / ([0-9]+)").FindStringSubmatch(state)
	if err != nil || len(capacity) != 2 {
		t.Fatal("statefile.log-spec.escape-capacity.result: missing serialized bound")
	}
	maximum, err := strconv.Atoi(capacity[1])
	if err != nil || len(encoded) <= 6*1048576 || len(encoded) > maximum {
		t.Fatalf("statefile.log-spec.escape-capacity.result: expected encoded NUL log above 6MiB and within cap, actual bytes=%d cap=%d; statefile.md#phase-19-log-invariants", len(encoded), maximum)
	}
	for _, c := range []struct{ id, input string }{
		{"missing", strings.Replace(match[1], "\"schema_version\": 1,", "", 1)},
		{"unknown", strings.Replace(match[1], "{", "{\"unknown\":true,", 1)},
		{"duplicate", strings.Replace(match[1], "\"schema_version\": 1", "\"schema_version\":1,\"schema_version\":1", 1)},
		{"null", regexp.MustCompile(`"patterns_hex": \[[^\]]*\]`).ReplaceAllString(match[1], `"patterns_hex": null`)},
		{"trailing", match[1] + " {}"},
		{"utf8", string([]byte{0xff}) + match[1]},
	} {
		if _, err := decodePhase19LogSpecification([]byte(c.input)); err == nil {
			t.Fatalf("runner.log-spec.schema-%s.result: expected malformed input rejected, actual accepted", c.id)
		}
	}
}

func specificationPhaseReferenceIntegrity(documents map[string]string) error {
	rows := regexp.MustCompile(`(?m)^\| Phase ([0-9]+) \|`).FindAllStringSubmatch(documents["docs/ROADMAP.md"], -1)
	known := make(map[string]bool)
	for _, row := range rows {
		phase, err := strconv.Atoi(row[1])
		if err != nil || phase < 1 || phase > 21 || known[row[1]] || strconv.Itoa(phase) != row[1] {
			return fmt.Errorf("invalid or duplicate Phase row %s", row[1])
		}
		known[row[1]] = true
	}
	if len(known) != 21 {
		return fmt.Errorf("expected 21 registered Phase rows, got %d", len(known))
	}
	expressions := regexp.MustCompile(`(?i:phase|フェーズ)[ _\t-]*[0-9]+(?:[ \t]*[/,、〜～~-][ \t]*(?:(?i:phase|フェーズ)[ _\t-]*)?[0-9]+)*`)
	numbers := regexp.MustCompile(`[0-9]+`)
	paths := make([]string, 0, len(documents))
	for path := range documents {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		for number, line := range strings.Split(documents[path], "\n") {
			for _, expression := range expressions.FindAllString(line, -1) {
				for _, phase := range numbers.FindAllString(expression, -1) {
					if !known[phase] {
						return fmt.Errorf("unregistered Phase %s in %s:%d", phase, path, number+1)
					}
				}
			}
		}
	}
	return nil
}

func specificationReferenceIntegrity(documents map[string]string, exists func(string) bool) error {
	anchors := map[string]map[string]bool{}
	anchorPattern := regexp.MustCompile(`<a id="([^"]+)"></a>`)
	for path, content := range documents {
		anchors[path] = map[string]bool{}
		for _, match := range anchorPattern.FindAllStringSubmatch(content, -1) {
			if anchors[path][match[1]] {
				return fmt.Errorf("duplicate anchor %s#%s", path, match[1])
			}
			anchors[path][match[1]] = true
		}
	}
	linkPattern := regexp.MustCompile(`\[[^\]\n]*\]\(([^)\n]+)\)`)
	inlineCode := regexp.MustCompile("(`+)[^`]*`+")
	for path, content := range documents {
		fence := ""
		for number, line := range strings.Split(content, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
				marker := trimmed[:3]
				if fence == "" {
					fence = marker
				} else if fence == marker {
					fence = ""
				}
				continue
			}
			if fence != "" {
				continue
			}
			for _, match := range linkPattern.FindAllStringSubmatch(inlineCode.ReplaceAllString(line, ""), -1) {
				ref, err := url.Parse(strings.Trim(match[1], "<>"))
				if err != nil {
					return fmt.Errorf("invalid reference %s:%d: %s", path, number+1, match[1])
				}
				if ref.IsAbs() || ref.Host != "" {
					continue
				}
				target := path
				if ref.Path != "" {
					target = filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(path), ref.Path)))
				}
				if !exists(target) {
					return fmt.Errorf("missing path %s:%d: %s", path, number+1, target)
				}
				if ref.Fragment != "" && strings.HasSuffix(target, ".md") && !anchors[target][ref.Fragment] {
					return fmt.Errorf("missing anchor %s:%d: %s#%s", path, number+1, target, ref.Fragment)
				}
			}
		}
	}
	return nil
}

type additionalOIDCVector struct {
	Challenge string `json:"challenge"`
	IDToken   string `json:"id_token"`
	JWK       struct {
		Algorithm string `json:"alg"`
		Exponent  string `json:"e"`
		ID        string `json:"kid"`
		Type      string `json:"kty"`
		Modulus   string `json:"n"`
		Use       string `json:"use"`
	} `json:"jwk"`
	Now      string `json:"now"`
	Verifier string `json:"verifier"`
}

func TestAdditionalSpecificationVectors(t *testing.T) {
	t.Parallel()
	document := phase11MustReadText(t, "docs/details/fixture.md")
	pattern := regexp.MustCompile("(?s)<a id=\"oidc-fixed-vector\"></a>.*?```json\\n(.*?)\\n```")
	match := pattern.FindStringSubmatch(document)
	if len(match) != 2 {
		t.Fatal("security.spec-vector.presence: missing fixed JSON vector")
	}
	vector, err := additionalDecodeOIDCVector(match[1])
	if err != nil {
		t.Fatalf("security.spec-vector.schema: %v", err)
	}
	for _, input := range []string{
		strings.Replace(match[1], "{", `{"unknown":true,`, 1),
		strings.Replace(match[1], "{", `{"now":"duplicate",`, 1),
		strings.Replace(match[1], "dBjft", string([]byte{0xff})+"dBjft", 1),
		match[1] + " {}",
	} {
		if _, err := additionalDecodeOIDCVector(input); err == nil {
			t.Fatal("security.spec-vector.schema: invalid input accepted")
		}
	}
	if err := additionalVerifyOIDCVector(vector); err != nil {
		t.Fatalf("security.spec-vector.positive: %v", err)
	}
	for _, control := range []struct {
		name   string
		change func(*additionalOIDCVector)
		want   string
	}{
		{"pkce", func(v *additionalOIDCVector) { v.Challenge += "x" }, "pkce"},
		{"key", func(v *additionalOIDCVector) { v.JWK.Exponent = "Aw" }, "key"},
		{"signature", func(v *additionalOIDCVector) {
			parts := strings.Split(v.IDToken, ".")
			signature, err := base64.RawURLEncoding.Strict().DecodeString(parts[2])
			if err != nil {
				t.Fatal(err)
			}
			signature[0] ^= 1
			parts[2] = base64.RawURLEncoding.EncodeToString(signature)
			v.IDToken = strings.Join(parts, ".")
		}, "signature"},
		{"algorithm", func(v *additionalOIDCVector) {
			parts := strings.Split(v.IDToken, ".")
			parts[0] = base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","kid":"fixture-key-1","typ":"JWT"}`))
			v.IDToken = strings.Join(parts, ".")
		}, "header"},
		{"expiry", func(v *additionalOIDCVector) { v.Now = "2026-01-02T00:05:30Z" }, "claims"},
	} {
		t.Run(control.name, func(t *testing.T) {
			changed := vector
			control.change(&changed)
			if err := additionalVerifyOIDCVector(changed); err == nil || err.Error() != control.want {
				t.Fatalf("security.spec-vector.negative.%s: expected %s, actual %v", control.name, control.want, err)
			}
		})
	}
}

func additionalDecodeOIDCVector(input string) (additionalOIDCVector, error) {
	var vector additionalOIDCVector
	if !utf8.ValidString(input) {
		return vector, fmt.Errorf("invalid UTF-8")
	}
	if err := phase16CheckJSONValueForDuplicateKeys(json.NewDecoder(strings.NewReader(input))); err != nil {
		return vector, err
	}
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&vector); err != nil {
		return vector, err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return vector, fmt.Errorf("trailing input")
	}
	return vector, nil
}

// This checks the published oracle, not the future OIDC HTTP implementation.
func additionalVerifyOIDCVector(v additionalOIDCVector) error {
	digest := sha256.Sum256([]byte(v.Verifier))
	if len(v.Verifier) != 43 || base64.RawURLEncoding.EncodeToString(digest[:]) != v.Challenge {
		return fmt.Errorf("pkce")
	}
	decode := base64.RawURLEncoding.Strict().DecodeString
	n, err := decode(v.JWK.Modulus)
	if err != nil || len(n) < 256 || len(n) > 512 || v.JWK.Exponent != "AQAB" || v.JWK.Type != "RSA" || v.JWK.Algorithm != "RS256" || v.JWK.Use != "sig" {
		return fmt.Errorf("key")
	}
	key := &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: 65537}
	parts := strings.Split(v.IDToken, ".")
	if len(parts) != 3 {
		return fmt.Errorf("header")
	}
	header, err := decode(parts[0])
	if err != nil {
		return fmt.Errorf("header")
	}
	var h struct {
		Algorithm string `json:"alg"`
		ID        string `json:"kid"`
		Type      string `json:"typ"`
	}
	d := json.NewDecoder(bytes.NewReader(header))
	d.DisallowUnknownFields()
	if d.Decode(&h) != nil || h.Algorithm != "RS256" || h.ID != v.JWK.ID || h.Type != "JWT" {
		return fmt.Errorf("header")
	}
	signature, err := decode(parts[2])
	if err != nil {
		return fmt.Errorf("signature")
	}
	signed := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if rsa.VerifyPKCS1v15(key, crypto.SHA256, signed[:], signature) != nil {
		return fmt.Errorf("signature")
	}
	payload, err := decode(parts[1])
	if err != nil {
		return fmt.Errorf("claims")
	}
	var c struct {
		Issuer   string `json:"iss"`
		Audience string `json:"aud"`
		Subject  string `json:"sub"`
		Nonce    string `json:"nonce"`
		Issued   int64  `json:"iat"`
		Expires  int64  `json:"exp"`
	}
	now, err := time.Parse(time.RFC3339, v.Now)
	if err != nil || json.Unmarshal(payload, &c) != nil || c.Issuer != "https://id.example.test" || c.Audience != "adlaire-ci-test" || c.Subject != "subject-1" || c.Nonce != strings.Repeat("A", 43) || c.Issued != 1767312000 || c.Expires != 1767312300 || now.Unix() >= c.Expires+30 {
		return fmt.Errorf("claims")
	}
	return nil
}

type additionalOIDCSignedVectors struct {
	FlowStartedAt string `json:"flow_started_at"`
	JWK           struct {
		Algorithm string `json:"alg"`
		Exponent  string `json:"e"`
		ID        string `json:"kid"`
		Type      string `json:"kty"`
		Modulus   string `json:"n"`
		Use       string `json:"use"`
	} `json:"jwk"`
	Cases []struct {
		ID           string `json:"id"`
		Header       string `json:"header"`
		ClaimsBase64 string `json:"claims_base64"`
		Now          string `json:"now"`
		Want         string `json:"want"`
		IDToken      string `json:"id_token"`
	} `json:"cases"`
}

func TestAdditionalOIDCSignedClaimVectors(t *testing.T) {
	t.Parallel()
	document := phase11MustReadText(t, "docs/details/fixture.md")
	pattern := regexp.MustCompile("(?s)<a id=\"oidc-signed-claim-vectors\"></a>.*?```json\\n(.*?)\\n```")
	match := pattern.FindStringSubmatch(document)
	if len(match) != 2 {
		t.Fatal("security.signed-vector.presence: missing bundle")
	}
	if _, err := additionalOIDCJSONObject([]byte(match[1])); err != nil {
		t.Fatalf("security.signed-vector.schema: %v", err)
	}
	var vectors additionalOIDCSignedVectors
	decoder := json.NewDecoder(strings.NewReader(match[1]))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&vectors); err != nil {
		t.Fatalf("security.signed-vector.schema: %v", err)
	}
	decode := base64.RawURLEncoding.Strict().DecodeString
	n, err := decode(vectors.JWK.Modulus)
	if err != nil || len(n) != 256 || vectors.JWK.Algorithm != "RS256" || vectors.JWK.Exponent != "AQAB" || vectors.JWK.ID != "fixture-signed-1" || vectors.JWK.Type != "RSA" || vectors.JWK.Use != "sig" || vectors.FlowStartedAt != "2026-01-02T00:00:00Z" {
		t.Fatal("security.signed-vector.schema: unexpected key or clock")
	}
	key := &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: 65537}
	wants := map[string]string{
		"baseline": "ok", "issuer-slash": "issuer", "issuer-wrong": "issuer", "issuer-missing": "issuer",
		"audience-wrong": "audience", "audience-empty": "audience", "audience-duplicate": "audience", "audience-non-string": "audience", "audience-missing": "audience",
		"audience-single-array": "ok", "audience-multi-no-azp": "azp", "audience-multi-wrong-azp": "azp", "audience-multi-valid": "ok", "audience-no-client": "audience", "azp-single-wrong": "azp",
		"nonce-wrong": "nonce", "nonce-missing": "nonce", "subject-empty": "subject", "subject-cr": "subject", "subject-255": "ok", "subject-256": "subject", "subject-null": "subject",
		"iat-future-boundary": "ok", "iat-future-outside": "issued", "iat-start-boundary": "ok", "iat-start-outside": "issued", "iat-fraction": "issued", "iat-exponent": "issued", "iat-null": "issued",
		"exp-missing": "expires", "exp-equals-iat": "expires", "exp-before-iat": "expires", "lifetime-boundary": "ok", "lifetime-outside": "expires", "exp-fraction": "expires", "exp-overflow": "expires",
		"nbf-boundary": "ok", "nbf-outside": "not-before", "nbf-null": "not-before", "unknown-claim": "ok", "duplicate-claim": "json", "invalid-utf8": "json", "trailing-json": "json", "expiry-before-boundary": "ok", "expiry-boundary": "expires",
		"algorithm-none": "header", "kid-unknown": "header", "header-jku": "header", "header-crit": "header", "header-duplicate": "header",
	}
	if len(vectors.Cases) != len(wants) {
		t.Fatalf("security.signed-vector.coverage: got %d cases, want %d", len(vectors.Cases), len(wants))
	}
	seen := map[string]bool{}
	for _, c := range vectors.Cases {
		c := c
		want, exists := wants[c.ID]
		if !exists || seen[c.ID] || c.Want != want {
			t.Fatalf("security.signed-vector.coverage: unknown, duplicate or changed expectation %q", c.ID)
		}
		seen[c.ID] = true
		t.Run(c.ID, func(t *testing.T) {
			parts := strings.Split(c.IDToken, ".")
			if len(parts) != 3 || parts[0] != base64.RawURLEncoding.EncodeToString([]byte(c.Header)) || parts[1] != c.ClaimsBase64 || len(c.IDToken) > 16384 {
				t.Fatal("security.signed-vector.input: fixed input mismatch")
			}
			signature, err := decode(parts[2])
			digest := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
			if err != nil || rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], signature) != nil {
				t.Fatal("security.signed-vector.signature: negative case must retain a valid RSA signature")
			}
			payload, err := decode(parts[1])
			if err != nil {
				t.Fatal("security.signed-vector.input: invalid payload encoding")
			}
			got := additionalOIDCClaimOracle([]byte(c.Header), payload, c.Now, vectors.FlowStartedAt)
			if got != want {
				t.Fatalf("security.signed-vector.claim.%s: got %s, want %s", c.ID, got, want)
			}
		})
	}
}

func additionalOIDCJSONObject(input []byte) (map[string]json.RawMessage, error) {
	if !utf8.Valid(input) {
		return nil, fmt.Errorf("invalid UTF-8")
	}
	if err := phase16CheckJSONValueForDuplicateKeys(json.NewDecoder(bytes.NewReader(input))); err != nil {
		return nil, err
	}
	var object map[string]json.RawMessage
	decoder := json.NewDecoder(bytes.NewReader(input))
	if err := decoder.Decode(&object); err != nil || object == nil {
		return nil, fmt.Errorf("object required")
	}
	if decoder.Decode(new(any)) != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	return object, nil
}

// This evaluates specification inputs only; it is not a product authentication path.
func additionalOIDCClaimOracle(header, payload []byte, nowText, startText string) string {
	h, err := additionalOIDCJSONObject(header)
	text := func(object map[string]json.RawMessage, key string) string {
		var s string
		if json.Unmarshal(object[key], &s) != nil {
			return ""
		}
		return s
	}
	if err != nil || text(h, "alg") != "RS256" || text(h, "kid") != "fixture-signed-1" {
		return "header"
	}
	for _, forbidden := range []string{"jku", "x5u", "x5c", "crit"} {
		if _, exists := h[forbidden]; exists {
			return "header"
		}
	}
	c, err := additionalOIDCJSONObject(payload)
	if err != nil {
		return "json"
	}
	if text(c, "iss") != "https://id.example.test" {
		return "issuer"
	}
	var audience []string
	if value := text(c, "aud"); value != "" {
		audience = []string{value}
	} else if json.Unmarshal(c["aud"], &audience) != nil {
		return "audience"
	}
	seen := map[string]bool{}
	for _, value := range audience {
		if value == "" || seen[value] {
			return "audience"
		}
		seen[value] = true
	}
	if !seen["adlaire-ci-test"] {
		return "audience"
	}
	if _, exists := c["azp"]; exists || len(audience) > 1 {
		if text(c, "azp") != "adlaire-ci-test" {
			return "azp"
		}
	}
	if text(c, "nonce") != strings.Repeat("A", 43) {
		return "nonce"
	}
	subject := text(c, "sub")
	if len(subject) < 1 || len(subject) > 255 || strings.ContainsAny(subject, "\x00\r\n") {
		return "subject"
	}
	now, nowErr := time.Parse(time.RFC3339, nowText)
	start, startErr := time.Parse(time.RFC3339, startText)
	integer := func(key string) (int64, bool) {
		raw := string(c[key])
		if !regexp.MustCompile(`^-?(0|[1-9][0-9]*)$`).MatchString(raw) {
			return 0, false
		}
		var value int64
		err := json.Unmarshal(c[key], &value)
		return value, err == nil
	}
	iat, valid := integer("iat")
	if !valid || nowErr != nil || startErr != nil || iat > now.Unix()+30 || iat < start.Unix()-30 {
		return "issued"
	}
	exp, valid := integer("exp")
	if !valid || exp <= iat || exp-iat > 3600 || exp <= now.Unix()-30 {
		return "expires"
	}
	if _, exists := c["nbf"]; exists {
		nbf, valid := integer("nbf")
		if !valid || nbf > now.Unix()+30 {
			return "not-before"
		}
	}
	return "ok"
}

func TestAdditionalIaCSpecificationVectors(t *testing.T) {
	t.Parallel()
	document := phase11MustReadText(t, "docs/details/fixture.md")
	pattern := regexp.MustCompile("(?s)<a id=\"inhouse-iac-concrete-oracles\"></a>.*?```json\\n(.*?)\\n```")
	match := pattern.FindStringSubmatch(document)
	if len(match) != 2 {
		t.Fatal("setup.spec-vector.presence: missing fixed input")
	}
	var input struct {
		SchemaVersion        int            `json:"schema_version"`
		Now                  string         `json:"now"`
		OperationID          string         `json:"operation_id"`
		Identity             map[string]any `json:"identity"`
		ConfigPath           string         `json:"config_path"`
		ConfigMetadata       map[string]any `json:"config_metadata"`
		BeforeConfig         map[string]any `json:"before_config"`
		BeforeSHA256         string         `json:"before_canonical_sha256"`
		AfterSHA256          string         `json:"after_canonical_sha256"`
		Manifest             map[string]any `json:"manifest"`
		Plan                 map[string]any `json:"plan"`
		ApplyConfirmation    map[string]any `json:"apply_confirmation"`
		RollbackConfirmation map[string]any `json:"rollback_confirmation"`
	}
	if _, err := additionalOIDCJSONObject([]byte(match[1])); err != nil {
		t.Fatalf("setup.spec-vector.schema: %v", err)
	}
	decoder := json.NewDecoder(strings.NewReader(match[1]))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(&input); err != nil {
		t.Fatalf("setup.spec-vector.schema: %v", err)
	}
	canonical := func(v any) []byte {
		var output bytes.Buffer
		encoder := json.NewEncoder(&output)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(v); err != nil {
			t.Fatal(err)
		}
		return bytes.TrimSuffix(output.Bytes(), []byte("\n"))
	}
	hash := func(v any) string {
		digest := sha256.Sum256(canonical(v))
		return fmt.Sprintf("%x", digest)
	}
	if input.SchemaVersion != 1 || input.Now != "2026-01-02T00:00:00Z" || input.OperationID != "iac_00000000000000000000000001" || input.ConfigPath != "/etc/adlaire-ci/cicd/coordinator-main.json" || string(canonical(input.ConfigMetadata)) != `{"gid":991,"mode":"0640","nlink":1,"uid":0}` || string(canonical(input.Identity)) != `{"schema_version":1,"target_id":"fixture-vps-1"}` {
		t.Fatal("setup.spec-vector.binding: changed identity, path, clock or metadata")
	}
	before := "a7dbfd1f06913f8a2bac4b1c7e1031e927a28c861196455617534e563449f0d1"
	after := "cb4d35943f79df39c14c5f9b2885a7437d6b189b25abf532be8390064d6ee451"
	manifest := "77591088e3e13b2ccff82c1b77be73cc310b262f9dd17c799a2304aad5886831"
	plan := "a184604a39947060afba533395847318f9d676385f88cc5cdce9a4f1acbf2eb2"
	if hash(input.BeforeConfig) != before || input.BeforeSHA256 != before || input.AfterSHA256 != after || hash(input.Manifest) != manifest || input.Plan["manifest_sha256"] != manifest || input.Plan["sha256"] != plan {
		t.Fatal("setup.spec-vector.digest: fixed config or manifest mismatch")
	}
	delete(input.Plan, "sha256")
	if hash(input.Plan) != plan || input.Plan["created_at"] != input.Now || input.Plan["expires_at"] != "2026-01-02T00:05:00Z" {
		t.Fatal("setup.spec-vector.plan: self-excluded hash or expiry mismatch")
	}
	resources, resourceOK := input.Manifest["resources"].([]any)
	changes, changesOK := input.Plan["changes"].([]any)
	observed, observedOK := input.Plan["observed"].([]any)
	if !resourceOK || !changesOK || !observedOK || len(resources) != 2 || len(changes) != 2 || len(observed) != 2 {
		t.Fatal("setup.spec-vector.order: missing resources")
	}
	if hash(resources[0].(map[string]any)["config"]) != after || hash(changes[0].(map[string]any)["after"]) != after || observed[0].(map[string]any)["id"] != "api" || observed[1].(map[string]any)["id"] != "main" || changes[0].(map[string]any)["id"] != "main" || changes[1].(map[string]any)["id"] != "api" {
		t.Fatal("setup.spec-vector.order: before/after or config-before-service mismatch")
	}
	for action, confirmation := range map[string]map[string]any{"apply": input.ApplyConfirmation, "rollback": input.RollbackConfirmation} {
		if len(confirmation) != 3 || confirmation["action"] != action || confirmation["target_id"] != "fixture-vps-1" || confirmation["plan_sha256"] != plan {
			t.Fatal("setup.spec-vector.confirmation: operation binding mismatch")
		}
	}
}

type iacSpecificationAdmissionVector struct {
	SchemaVersion          int    `json:"schema_version"`
	Mode                   string `json:"mode"`
	CreatedAt              string `json:"created_at"`
	ExpiresAt              string `json:"expires_at"`
	PreflightAt            string `json:"preflight_at"`
	LockAcquiredAt         string `json:"lock_acquired_at"`
	ObservationsFinishedAt string `json:"observations_finished_at"`
	AdmissionAt            string `json:"admission_at"`
	ExecutionAt            string `json:"execution_at"`
}

func decodeIaCSpecificationAdmission(input []byte) (iacSpecificationAdmissionVector, error) {
	var v iacSpecificationAdmissionVector
	object, err := additionalOIDCJSONObject(input)
	if err != nil || len(object) != 9 {
		return v, fmt.Errorf("closed object required")
	}
	for _, key := range []string{"schema_version", "mode", "created_at", "expires_at", "preflight_at", "lock_acquired_at", "observations_finished_at", "admission_at", "execution_at"} {
		value, exists := object[key]
		if !exists || bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			return v, fmt.Errorf("required field %s", key)
		}
	}
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&v); err != nil {
		return v, err
	}
	return v, nil
}

// These oracles evaluate published vectors, not the future privileged IaC implementation.
func iacSpecificationAdmissionOracle(v iacSpecificationAdmissionVector) string {
	if v.SchemaVersion != 1 || (v.Mode != "apply" && v.Mode != "noop" && v.Mode != "rollback" && v.Mode != "recover") {
		return "invalid_input"
	}
	texts := []string{v.CreatedAt, v.ExpiresAt, v.PreflightAt, v.LockAcquiredAt, v.ObservationsFinishedAt, v.AdmissionAt, v.ExecutionAt}
	times := make([]time.Time, len(texts))
	for index, value := range texts {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil || !regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$`).MatchString(value) || parsed.Year() < 1970 {
			return "invalid_input"
		}
		times[index] = parsed
	}
	if !times[1].Equal(times[0].Add(300 * time.Second)) {
		return "invalid_input"
	}
	if v.Mode == "rollback" || v.Mode == "recover" {
		return "ok"
	}
	for _, now := range times[2:6] {
		if now.Before(times[0]) || !now.Before(times[1]) {
			return "stale_plan"
		}
	}
	return "ok"
}

func iacSpecificationCanonical(value any) []byte {
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	if encoder.Encode(value) != nil {
		return nil
	}
	return bytes.TrimSuffix(output.Bytes(), []byte("\n"))
}

func iacSpecificationHash(value any) string {
	digest := sha256.Sum256(iacSpecificationCanonical(value))
	return hex.EncodeToString(digest[:])
}

func iacSpecificationPlanOracle(manifest, plan map[string]any) string {
	if len(plan) != 8 || plan["schema_version"] != json.Number("1") || plan["manifest_sha256"] != iacSpecificationHash(manifest) {
		return "invalid_input"
	}
	if plan["target_id"] != manifest["target_id"] {
		return "target_mismatch"
	}
	unsigned := make(map[string]any, len(plan)-1)
	for key, value := range plan {
		if key != "sha256" {
			unsigned[key] = value
		}
	}
	if plan["sha256"] != iacSpecificationHash(unsigned) {
		return "invalid_input"
	}
	resources, resourcesOK := manifest["resources"].([]any)
	observed, observedOK := plan["observed"].([]any)
	changes, changesOK := plan["changes"].([]any)
	if !resourcesOK || !observedOK || !changesOK || len(resources) == 0 || len(resources) > 64 || len(observed) != len(resources) {
		return "invalid_input"
	}
	byID := map[string]map[string]any{}
	for _, value := range resources {
		resource, ok := value.(map[string]any)
		id, idOK := resource["id"].(string)
		if !ok || !idOK || !regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`).MatchString(id) || byID[id] != nil {
			return "invalid_input"
		}
		byID[id] = resource
	}
	expected := []any{}
	previousID := ""
	for _, value := range observed {
		observation, ok := value.(map[string]any)
		id, idOK := observation["id"].(string)
		before, beforeOK := observation["before"].(map[string]any)
		resource := byID[id]
		if !ok || !idOK || !beforeOK || len(observation) != 3 || resource == nil || id <= previousID || observation["type"] != resource["type"] {
			return "invalid_input"
		}
		previousID = id
		var after any
		action := ""
		different := false
		switch resource["type"] {
		case "config":
			digest, digestOK := before["sha256"].(string)
			if len(before) != 4 || !digestOK || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(digest) || before["uid"] != json.Number("0") || before["mode"] != "0640" {
				return "invalid_input"
			}
			gid, gidOK := before["gid"].(json.Number)
			if group, err := gid.Int64(); !gidOK || err != nil || group < 0 || group > 4294967295 {
				return "invalid_input"
			}
			after, action = resource["config"], "replace_config"
			different = digest != iacSpecificationHash(after)
		case "service":
			_, activeOK := before["active"].(bool)
			_, enabledOK := before["enabled"].(bool)
			if len(before) != 2 || !activeOK || !enabledOK {
				return "invalid_input"
			}
			after, action = map[string]any{"active": resource["active"], "enabled": resource["enabled"]}, "set_service"
			different = !bytes.Equal(iacSpecificationCanonical(before), iacSpecificationCanonical(after))
		default:
			return "invalid_input"
		}
		if different {
			expected = append(expected, map[string]any{"id": id, "action": action, "after": after})
		}
	}
	sort.Slice(expected, func(i, j int) bool {
		left, right := expected[i].(map[string]any), expected[j].(map[string]any)
		if left["action"] != right["action"] {
			return left["action"] == "replace_config"
		}
		return left["id"].(string) < right["id"].(string)
	})
	if !bytes.Equal(iacSpecificationCanonical(changes), iacSpecificationCanonical(expected)) {
		return "invalid_input"
	}
	return "ok"
}

func TestIaCSpecificationAdmissionVectors(t *testing.T) {
	t.Parallel()
	document := phase11MustReadText(t, "docs/details/fixture.md")
	read := func(anchor string) []byte {
		pattern := regexp.MustCompile("(?s)<a id=\"" + regexp.QuoteMeta(anchor) + "\"></a>.*?```json\\n(.*?)\\n```")
		match := pattern.FindStringSubmatch(document)
		if len(match) != 2 {
			t.Fatalf("setup.iac-admission.vector.present: anchor=%s missing fixed input", anchor)
		}
		return []byte(match[1])
	}
	baselineBytes := read("inhouse-iac-admission-vectors")
	baseline, err := decodeIaCSpecificationAdmission(baselineBytes)
	if err != nil {
		t.Fatalf("setup.iac-admission.vector.schema: %v", err)
	}
	assert := func(t *testing.T, id, got, want string) {
		t.Helper()
		if got != want {
			t.Fatalf("setup.iac-admission.%s.result: anchor=docs/details/setup.md#inhouse-iac-admission expected=%s actual=%s diff=code reason=admission-contract command=go test . -run '^TestIaCSpecificationAdmissionVectors$' -count=1", id, want, got)
		}
	}
	check := func(id string, mutate func(*iacSpecificationAdmissionVector), want string) {
		t.Run(id, func(t *testing.T) {
			v := baseline
			mutate(&v)
			assert(t, id, iacSpecificationAdmissionOracle(v), want)
		})
	}
	check("apply-baseline", func(*iacSpecificationAdmissionVector) {}, "ok")
	check("noop-baseline", func(v *iacSpecificationAdmissionVector) { v.Mode = "noop" }, "ok")
	check("lower-inclusive", func(v *iacSpecificationAdmissionVector) {
		v.PreflightAt, v.LockAcquiredAt, v.ObservationsFinishedAt, v.AdmissionAt = v.CreatedAt, v.CreatedAt, v.CreatedAt, v.CreatedAt
	}, "ok")
	check("last-second", func(v *iacSpecificationAdmissionVector) { v.AdmissionAt = "2026-01-02T00:04:59Z" }, "ok")
	for _, mode := range []string{"rollback", "recover"} {
		mode := mode
		check(mode+"-expired", func(v *iacSpecificationAdmissionVector) {
			v.Mode = mode
			v.PreflightAt, v.LockAcquiredAt, v.ObservationsFinishedAt, v.AdmissionAt = v.ExecutionAt, v.ExecutionAt, v.ExecutionAt, v.ExecutionAt
		}, "ok")
	}
	for _, point := range []string{"preflight", "locked", "observed", "admission"} {
		for _, edge := range []struct{ name, value string }{{"before", "2026-01-01T23:59:59Z"}, {"equal", "2026-01-02T00:05:00Z"}, {"after", "2026-01-02T00:05:01Z"}} {
			point, edge := point, edge
			check(point+"-"+edge.name, func(v *iacSpecificationAdmissionVector) {
				switch point {
				case "preflight":
					v.PreflightAt = edge.value
				case "locked":
					v.LockAcquiredAt = edge.value
				case "observed":
					v.ObservationsFinishedAt = edge.value
				case "admission":
					v.AdmissionAt = edge.value
				}
			}, "stale_plan")
		}
	}
	for _, c := range []struct{ id, value string }{{"ttl-short", "2026-01-02T00:04:59Z"}, {"ttl-long", "2026-01-02T00:05:01Z"}} {
		c := c
		check(c.id, func(v *iacSpecificationAdmissionVector) { v.ExpiresAt = c.value }, "invalid_input")
	}
	for _, c := range []struct{ id, value string }{{"invalid-time", "not-time"}, {"offset-time", "2026-01-02T00:00:00+00:00"}, {"fractional-time", "2026-01-02T00:00:00.1Z"}, {"pre-epoch", "1969-01-02T00:00:00Z"}} {
		c := c
		check(c.id, func(v *iacSpecificationAdmissionVector) { v.CreatedAt = c.value }, "invalid_input")
	}
	check("invalid-mode", func(v *iacSpecificationAdmissionVector) { v.Mode = "resume" }, "invalid_input")
	for _, c := range []struct{ id, input string }{
		{"schema-unknown", strings.Replace(string(baselineBytes), `"schema_version": 1`, `"schema_version": 1, "x": true`, 1)},
		{"schema-missing", strings.Replace(string(baselineBytes), ",\n  \"execution_at\": \"2026-01-02T00:05:01Z\"", "", 1)},
		{"schema-null", strings.Replace(string(baselineBytes), `"mode": "apply"`, `"mode": null`, 1)},
		{"schema-duplicate", strings.Replace(string(baselineBytes), `"mode": "apply"`, `"mode": "apply", "mode": "apply"`, 1)},
		{"schema-utf8", string(baselineBytes) + string([]byte{255})},
		{"schema-trailing", string(baselineBytes) + `{}`},
	} {
		c := c
		t.Run(c.id, func(t *testing.T) {
			_, err := decodeIaCSpecificationAdmission([]byte(c.input))
			got := "accepted"
			if err != nil {
				got = "schema"
			}
			assert(t, c.id, got, "schema")
		})
	}
	planBytes := read("inhouse-iac-concrete-oracles")
	planCheck := func(id string, mutate func(map[string]any), rehash bool, want string) {
		t.Run(id, func(t *testing.T) {
			var bundle map[string]any
			decoder := json.NewDecoder(bytes.NewReader(planBytes))
			decoder.UseNumber()
			if err := decoder.Decode(&bundle); err != nil {
				t.Fatal(err)
			}
			plan := bundle["plan"].(map[string]any)
			mutate(plan)
			if rehash {
				delete(plan, "sha256")
				plan["sha256"] = iacSpecificationHash(plan)
			}
			assert(t, id, iacSpecificationPlanOracle(bundle["manifest"].(map[string]any), plan), want)
		})
	}
	planCheck("plan-baseline", func(map[string]any) {}, false, "ok")
	for _, id := range []string{"plan-after", "plan-missing", "plan-duplicate", "plan-order", "plan-unknown", "plan-observed-missing", "plan-observed-order", "plan-observed-type", "plan-before-key", "plan-noop-extra", "plan-noop-valid", "plan-manifest-hash", "plan-target"} {
		id := id
		want := "invalid_input"
		if id == "plan-noop-valid" {
			want = "ok"
		} else if id == "plan-target" {
			want = "target_mismatch"
		}
		planCheck(id, func(plan map[string]any) {
			changes := plan["changes"].([]any)
			observed := plan["observed"].([]any)
			switch id {
			case "plan-after":
				changes[0].(map[string]any)["after"].(map[string]any)["allow_http"] = false
			case "plan-missing":
				plan["changes"] = changes[:1]
			case "plan-duplicate":
				plan["changes"] = append(changes, changes[0])
			case "plan-order":
				changes[0], changes[1] = changes[1], changes[0]
			case "plan-unknown":
				changes[0].(map[string]any)["id"] = "other"
			case "plan-observed-missing":
				plan["observed"] = observed[:1]
			case "plan-observed-order":
				observed[0], observed[1] = observed[1], observed[0]
			case "plan-observed-type":
				observed[0].(map[string]any)["type"] = "config"
			case "plan-before-key":
				observed[0].(map[string]any)["before"].(map[string]any)["extra"] = true
			case "plan-noop-extra", "plan-noop-valid":
				observed[0].(map[string]any)["before"] = map[string]any{"active": true, "enabled": true}
				if id == "plan-noop-valid" {
					plan["changes"] = changes[:1]
				}
			case "plan-manifest-hash":
				plan["manifest_sha256"] = strings.Repeat("0", 64)
			case "plan-target":
				plan["target_id"] = "other"
			}
		}, true, want)
	}
	planCheck("plan-self-hash", func(plan map[string]any) { plan["sha256"] = strings.Repeat("0", 64) }, false, "invalid_input")
}

func TestAdditionalManagementInputVectors(t *testing.T) {
	t.Parallel()
	document := phase11MustReadText(t, "docs/details/fixture.md")
	read := func(anchor string) map[string]any {
		pattern := regexp.MustCompile("(?s)<a id=\"" + regexp.QuoteMeta(anchor) + "\"></a>.*?```json\\n(.*?)\\n```")
		match := pattern.FindStringSubmatch(document)
		if len(match) != 2 {
			t.Fatalf("fixture.concrete.presence.%s", anchor)
		}
		if _, err := additionalOIDCJSONObject([]byte(match[1])); err != nil {
			t.Fatalf("fixture.concrete.schema.%s: %v", anchor, err)
		}
		decoder := json.NewDecoder(strings.NewReader(match[1]))
		decoder.UseNumber()
		var result map[string]any
		if err := decoder.Decode(&result); err != nil {
			t.Fatalf("fixture.concrete.schema.%s: %v", anchor, err)
		}
		return result
	}
	canonical := func(value any) []byte {
		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(value); err != nil {
			t.Fatal(err)
		}
		return bytes.TrimSuffix(buffer.Bytes(), []byte("\n"))
	}
	hash := func(value any) string {
		digest := sha256.Sum256(canonical(value))
		return hex.EncodeToString(digest[:])
	}
	management := read("additional-management-fixed-inputs")
	history := management["history"].([]any)
	if len(history) != 3 || history[0].(map[string]any)["id"] != "b20260101000000" || history[1].(map[string]any)["id"] != "b20260102000000" || history[2].(map[string]any)["id"] != "b20260102000001" {
		t.Fatal("fixture.concrete.history.binding")
	}
	stats := management["stats_response"].(map[string]any)
	if string(canonical(stats["summary"])) != `{"duration_count":1,"duration_sum_seconds":10,"status_counts":{"failure_build":1,"success":1},"total":2}` {
		t.Fatal("fixture.concrete.stats.summary")
	}
	for _, value := range history {
		row := value.(map[string]any)
		if row["status"] == "failure_build" && row["failure_category"] != "pipeline_exit" {
			t.Fatal("fixture.concrete.history.category")
		}
	}
	state := management["build_state"].(map[string]any)
	queued := state["queued"].([]any)
	ids := management["reorder_request"].(map[string]any)["body"].(map[string]any)["queue_ids"].([]any)
	if len(queued) != 2 || len(ids) != 2 || queued[0].(map[string]any)["id"] != ids[0] || queued[1].(map[string]any)["id"] != ids[1] || ids[0] != "q20260102000000" || ids[1] != "q20260102000000-001" {
		t.Fatal("fixture.concrete.queue.binding")
	}
	inputs := read("additional-management-state-inputs")
	snapshot := inputs["snapshot"].(map[string]any)
	if len(snapshot["files"].(map[string]any)) != 13 || hash(snapshot["files"]) != "6d21ac0c8cbd9a2e802717f11ec3a7a1d9fbffb17e91d8042a9f4331a061bcc5" || snapshot["sha256"] != hash(snapshot["files"]) {
		t.Fatal("fixture.concrete.snapshot.digest")
	}
	payload := inputs["notify_payload"].(map[string]any)
	if hash(payload["payload"]) != "3ab1644e8b1fcc545756baa961ec40b34cdeabd048650c420887fbdbbed7eebb" || payload["payload_sha256"] != hash(payload["payload"]) || inputs["notify_delivery"].(map[string]any)["payload_sha256"] != payload["payload_sha256"] {
		t.Fatal("fixture.concrete.resend.payload")
	}
	public := read("additional-management-public-inputs")
	key := "76144c2400bc476b3d4081c4f2ae5b76668b33fa9887ab95f86f97680e73638e"
	entry := public["response_cache"].(map[string]any)["entries"].(map[string]any)[key].(map[string]any)
	if hash(public["cache_key_input"]) != key || hash(entry["vary"]) != key || entry["key"] != key || hash(entry["body"]) != "1f68af6ca60ce02c628ec72cdf1ff65bb2180aecb73dc2091c3072a1d84a70c9" || entry["body_sha256"] != hash(entry["body"]) {
		t.Fatal("fixture.concrete.cache.digest")
	}
	raw, err := hex.DecodeString(public["share_entropy_hex"].(string))
	if err != nil || len(raw) != 32 {
		t.Fatal("fixture.concrete.share.entropy")
	}
	digest := sha256.Sum256(raw)
	link := public["share_links"].(map[string]any)["links"].([]any)[0].(map[string]any)
	if link["token_hash"] != hex.EncodeToString(digest[:]) {
		t.Fatal("fixture.concrete.share.hash: expected raw-byte hash")
	}
	event := public["event_appends"].([]any)[0].(map[string]any)
	wantFrame := "id: evt_20260102000000\nevent: admin-event\ndata: " + string(canonical(event)) + "\n\n: keepalive\n\nevent: error\ndata: {\"error\":\"Admin events read failed\"}\n\n"
	if public["event_stream_expected"] != wantFrame {
		t.Fatal("fixture.concrete.event.frame")
	}
	for status, value := range public["badge_svg"].(map[string]any) {
		data := value.(string)
		decoder := xml.NewDecoder(strings.NewReader(data))
		var names, texts, fills []string
		for {
			token, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("fixture.concrete.badge.xml.%s", status)
			}
			switch element := token.(type) {
			case xml.StartElement:
				names = append(names, element.Name.Local)
				for _, attribute := range element.Attr {
					if strings.HasPrefix(attribute.Name.Local, "on") || strings.Contains(attribute.Value, "url(") {
						t.Fatalf("fixture.concrete.badge.security.%s", status)
					}
					if element.Name.Local == "rect" && attribute.Name.Local == "fill" {
						fills = append(fills, attribute.Value)
					}
				}
			case xml.CharData:
				if strings.TrimSpace(string(element)) != "" {
					texts = append(texts, string(element))
				}
			}
		}
		colors := map[string]string{"success": "#2e7d32", "failure": "#c62828", "running": "#1565c0", "skipped": "#6d6d6d", "unknown": "#455a64"}
		if len(public["badge_svg"].(map[string]any)) != len(colors) || colors[status] == "" || strings.Join(names, ",") != "svg,rect,rect,text,text" || strings.Join(texts, ",") != "Adlaire CI,"+status || strings.Join(fills, ",") != "#555555,"+colors[status] || !strings.HasSuffix(data, "</svg>\n") || !strings.Contains(data, `width="120" height="20" viewBox="0 0 120 20" role="img" aria-label="Adlaire CI: `+status+`"`) {
			t.Fatalf("fixture.concrete.badge.structure.%s", status)
		}
	}
	permissions := []string{"build:write", "config:read", "config:write", "status:read", "system:admin", "user:admin"}
	bootstrap := read("additional-management-bootstrap-inputs")
	for _, name := range []string{"self_admin", "self_oidc", "self_manager"} {
		self := bootstrap[name].(map[string]any)
		capabilities := self["capabilities"].(map[string]any)
		if len(capabilities) != len(permissions) {
			t.Fatal("fixture.concrete.self.capability-closure")
		}
		for _, permission := range permissions {
			want := name == "self_admin" || name == "self_oidc" && (permission == "status:read" || permission == "build:write" || permission == "config:read") || name == "self_manager" && permission == "user:admin"
			if capabilities[permission] != want {
				t.Fatalf("fixture.concrete.self.capability.%s.%s", name, permission)
			}
		}
		user := self["user"].(map[string]any)
		if len(user) != 10 || user["credential"] != nil || user["totp"] != nil || user["password_set"] != (name != "self_oidc") {
			t.Fatal("fixture.concrete.self.projection")
		}
	}
	if string(canonical(bootstrap["options_legacy"])) != `{"profile":"legacy-v1","providers":[]}` || string(canonical(bootstrap["options_multi"])) != `{"profile":"multi-user-v1","providers":[{"id":"fixture","issuer_origin":"https://id.example.test"}]}` {
		t.Fatal("fixture.concrete.bootstrap.options")
	}
	authorization := read("additional-management-authorization-inputs")
	requests := authorization["requests"].([]any)
	sdk := read("additional-management-sdk-inputs")["rows"].([]any)
	if len(requests) != 50 || len(sdk) != 50 {
		t.Fatal("fixture.concrete.transport.count")
	}
	seenRoutes, seenMethods := map[string]bool{}, map[string]bool{}
	protected, publicCount, sessionSelf := 0, 0, 0
	api := phase11MustReadText(t, "docs/details/api.md")
	api = strings.Split(strings.Split(api, "<a id=\"additional-management-api-contract\"></a>")[1], "<a id=\"phase-13-api-alignment-contract\"></a>")[0]
	apiRoutes := regexp.MustCompile("(?m)^\\| `(GET|POST|PATCH|DELETE)` \\| `([^`]+)` \\|").FindAllStringSubmatch(api, -1)
	if len(apiRoutes) != len(requests) {
		t.Fatal("fixture.concrete.transport.api-closure")
	}
	for index, value := range requests {
		request := value.(map[string]any)
		row := sdk[index].(map[string]any)
		key := request["method"].(string) + " " + request["route_template"].(string)
		if seenRoutes[key] || seenMethods[row["sdk_method"].(string)] || request["method"] != row["http_method"] || request["path"] != row["http_path"] {
			t.Fatalf("fixture.concrete.transport.binding.%d", index)
		}
		seenRoutes[key], seenMethods[row["sdk_method"].(string)] = true, true
		if !strings.Contains(api, "| `"+request["method"].(string)+"` | `"+request["route_template"].(string)+"` |") {
			t.Fatalf("fixture.concrete.transport.route.%s", key)
		}
		if request["required_permission"] != nil {
			protected++
			if request["denied_actor_id"] == nil {
				t.Fatal("fixture.concrete.transport.denied-actor")
			}
		} else if request["auth_boundary"] == "session-self" {
			sessionSelf++
		} else {
			publicCount++
		}
	}
	if protected != 45 || publicCount != 4 || sessionSelf != 1 || !seenMethods["getCurrentUser"] || !seenMethods["getAuthOptions"] {
		t.Fatal("fixture.concrete.transport.auth-closure")
	}
	migration := read("additional-management-migration-inputs")
	journal := migration["journal"].(map[string]any)
	for name, value := range migration["source"].(map[string]any) {
		digest := sha256.Sum256(append(canonical(value), '\n'))
		if journal["source_sha256"].(map[string]any)[name] != hex.EncodeToString(digest[:]) {
			t.Fatalf("fixture.concrete.migration.source.%s", name)
		}
	}
	for name, value := range migration["outputs"].(map[string]any) {
		if journal["output_sha256"].(map[string]any)[name] != hash(value) {
			t.Fatalf("fixture.concrete.migration.output.%s", name)
		}
	}
	if journal["output_sha256"].(map[string]any)[".users"] != "ec4dc9f8bd7d0e46ba7a6c9529e7e03386c7740784f93260ae4c615408eea836" {
		t.Fatal("fixture.concrete.migration.user-golden")
	}
	archives := read("additional-management-site-diff-inputs")
	for _, side := range []string{"left", "right"} {
		input := archives[side].(map[string]any)
		data, err := base64.StdEncoding.Strict().DecodeString(input["archive_base64"].(string))
		if err != nil {
			t.Fatalf("fixture.concrete.archive.base64.%s", side)
		}
		digest := sha256.Sum256(data)
		wantArchive := map[string]string{"left": "83f678968b2ab601b3865c90db00cf74da13546adb9ae8cb5067909021c52d8e", "right": "d1027eb4c3c0d48f695b713584e91e0d385fef32637e60ca170ac2b15b87faaa"}
		if input["archive_sha256"] != hex.EncodeToString(digest[:]) || input["archive_sha256"] != wantArchive[side] {
			t.Fatalf("fixture.concrete.archive.digest.%s", side)
		}
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("fixture.concrete.archive.gzip.%s", side)
		}
		tarReader := tar.NewReader(reader)
		files := input["files"].(map[string]any)
		seen := map[string]bool{}
		var total int64
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil || header.Typeflag != tar.TypeReg || seen[header.Name] {
				t.Fatalf("fixture.concrete.archive.header.%s", side)
			}
			body, err := io.ReadAll(io.LimitReader(tarReader, 1025))
			if err != nil || len(body) > 1024 || files[header.Name] != string(body) || header.Size != int64(len(body)) {
				t.Fatalf("fixture.concrete.archive.file.%s", side)
			}
			digest := sha256.Sum256(body)
			if input["file_sha256"].(map[string]any)[header.Name] != hex.EncodeToString(digest[:]) {
				t.Fatalf("fixture.concrete.archive.file-digest.%s", side)
			}
			seen[header.Name] = true
			total += header.Size
		}
		if _, err := io.Copy(io.Discard, reader); err != nil {
			t.Fatalf("fixture.concrete.archive.trailer.%s", side)
		}
		if err := reader.Close(); err != nil {
			t.Fatalf("fixture.concrete.archive.close.%s", side)
		}
		meta := input["meta"].(map[string]any)
		if len(seen) != 3 || len(files) != 3 || meta["file_count"] != json.Number("3") || meta["size_bytes"] != json.Number(strconv.FormatInt(total, 10)) {
			t.Fatalf("fixture.concrete.archive.metadata.%s", side)
		}
	}
}

func plannedPhaseMaturityBoundary(roadmap string) error {
	if !strings.Contains(roadmap, "現在の active Phase は設定しない。") || regexp.MustCompile("現在の active Phase は Phase (19|20|21)").MatchString(roadmap) {
		return fmt.Errorf("planned Phase activation boundary violated")
	}
	for phase := 19; phase <= 21; phase++ {
		pattern := fmt.Sprintf("(?m)^\\| Phase %d \\| [^\\n]+ \\| 仕様化済み・未実装 \\| Phase %d \\|$", phase, phase-1)
		if len(regexp.MustCompile(pattern).FindAllString(roadmap, -1)) != 1 {
			return fmt.Errorf("planned Phase %d must remain specified and unimplemented", phase)
		}
	}
	return nil
}

func TestPhase19SpecificationContracts(t *testing.T) {
	t.Parallel()

	files := []string{
		"docs/SPEC.md", "docs/ROADMAP.md", "docs/DETAIL_INDEX.md",
		"docs/DOCUMENT_INDEX.md", "docs/DESIGN.md",
		"docs/details/statefile.md", "docs/details/runner.md",
		"docs/details/archive.md", "docs/details/security.md",
		"docs/details/api.md", "docs/details/sdk.md", "docs/details/ui.md",
		"docs/details/admin.md", "docs/details/release.md",
		"docs/details/setup.md", "docs/details/commitstatus.md",
		"docs/details/fixture.md", "docs/details/production-validation.md",
	}
	documents := make(map[string]string, len(files))
	for _, path := range files {
		documents[path] = phase11MustReadText(t, path)
	}
	if err := phase19ValidateSpecification(documents); err != nil {
		t.Fatalf("runner.phase19-spec.positive.accept: expected valid specification, actual %v", err)
	}
	for _, mutation := range []struct {
		name string
		path string
		from string
		to   string
	}{
		{"invalid-arity", "docs/details/statefile.md", `"parallel": "int(1,16)"`, `"parallel": "int(1)"`},
		{"reversed-range", "docs/details/statefile.md", `"parallel": "int(1,16)"`, `"parallel": "int(16,1)"`},
		{"unknown-function", "docs/details/statefile.md", `"parallel": "int(1,16)"`, `"parallel": "execute(1,16)"`},
		{"duplicate-json-key", "docs/details/statefile.md", `"parallel": "int(1,16)"`, `"parallel": "int(1,16)", "parallel": "int(1,16)"`},
		{"duplicate-enum", "docs/details/statefile.md", `enum(builder)`, `enum(builder,builder)`},
		{"recursive-type", "docs/details/statefile.md", `"Id": "id"`, `"Id": "Id"`},
		{"environment-key-type", "docs/details/statefile.md", `"env": "envmap(EnvValue,0,64)"`, `"env": "map(EnvValue,0,64)"`},
		{"site-atomic-state", "docs/details/statefile.md", `"counters": "SiteCounters"`, `"wrong_counters": "SiteCounters"`},
		{"site-retention-clock", "docs/details/statefile.md", `"committed_at": "Time"`, `"committed_at": "Id"`},
		{"planned-selection", "docs/details/statefile.md", `"selected": "bool"`, `"wrong_selected": "bool"`},
		{"planned-task", "docs/details/statefile.md", `"task": "PlannedTask"`, `"task": "Task"`},
		{"planned-artifact-id", "docs/details/statefile.md", "\"PlannedDeployTask\": {\n      \"kind\": \"enum(deploy)\",\n      \"environment_id\": \"Id\",\n      \"artifact_from\": \"Id\"", "\"PlannedDeployTask\": {\n      \"kind\": \"enum(deploy)\",\n      \"environment_id\": \"Id\",\n      \"artifact_from\": \"JobId\""},
		{"cursor-revision", "docs/details/statefile.md", "\"TriggerCursor\": {\n      \"project_id\": \"Id\",\n      \"definition_id\": \"Id\",\n      \"definition_revision\": \"int(1,2147483647)\"", "\"TriggerCursor\": {\n      \"project_id\": \"Id\",\n      \"definition_id\": \"Id\",\n      \"definition_revision\": \"Id\""},
		{"operation-attempt", "docs/details/statefile.md", `"target_attempt": "nullable(int(1,100))"`, `"target_attempt": "nullable(Id)"`},
		{"operation-digest-target", "docs/details/statefile.md", `"target_id": "nullable(text(1,64))"`, `"target_id": "nullable(Id)"`},
		{"audit-digest-target", "docs/details/statefile.md", "\"at\": \"Time\",\n      \"target_id\": \"nullable(text(1,64))\"", "\"at\": \"Time\",\n      \"target_id\": \"nullable(Id)\""},
		{"deployment-observation", "docs/details/statefile.md", `"observations": "list(StageObservation,0,512)"`, `"observations": "list(StageObservation,0,32)"`},
		{"notification-attempt", "docs/details/statefile.md", `"target_attempt": "int(1,100)"`, `"target_attempt": "int(1,3)"`},
		{"worker-lease-project", "docs/details/statefile.md", "\"Lease\": {\n      \"project_id\": \"Id\"", "\"Lease\": {\n      \"project_id\": \"JobId\""},
		{"worker-assignment-digest", "docs/details/statefile.md", `"assignment_digest": "Digest"`, `"assignment_digest": "Id"`},
		{"worker-execution-id", "docs/details/statefile.md", `"execution_id": "nullable(Id)"`, `"execution_id": "nullable(Digest)"`},
		{"worker-assignment-size", "docs/details/statefile.md", `"assignment_bytes": "int(1,2097152)"`, `"assignment_bytes": "int(0,2097152)"`},
		{"worker-toolchain-digest", "docs/details/statefile.md", "\"assignment_digest\": \"Digest\",\n      \"toolchain_digest\": \"Digest\"", "\"assignment_digest\": \"Digest\",\n      \"toolchain_digest\": \"Id\""},
		{"roadmap-inventory-state", "docs/ROADMAP.md", "| 仕様化済み・未実装 | CI / CD | Phase 19", "| 改訂予定 | CI / CD | Phase 19"},
		{"roadmap-phase20-implemented", "docs/ROADMAP.md", "| 仕様化済み・未実装 | Phase 19 |", "| 実装済み | Phase 19 |"},
		{"roadmap-phase21-active", "docs/ROADMAP.md", "現在の active Phase は設定しない。", "現在の active Phase は Phase 21。"},
		{"roadmap-phase21-missing", "docs/ROADMAP.md", "| Phase 21 |", "| Phase 99 |"},
		{"unknown-type", "docs/details/statefile.md", `"parallel": "int(1,16)"`, `"parallel": "MissingType"`},
		{"missing-sdk", "docs/details/sdk.md", "| `runStart(params,body,options)` |", "| `missing(params,body,options)` |"},
		{"missing-ui", "docs/details/ui.md", "| runStart |", "| missing |"},
		{"missing-cli", "docs/details/admin.md", "| `run-start` | runStart |", "| `missing` | missing |"},
		{"duplicate-route", "docs/details/api.md", "| runStart | POST", "| plan | POST"},
		{"duplicate-anchor", "docs/details/security.md", "<a id=\"phase-19-permission-matrix\"></a>", "<a id=\"phase-19-permission-matrix\"></a><a id=\"phase-19-permission-matrix\"></a>"},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			changed := make(map[string]string, len(documents))
			for path, content := range documents {
				changed[path] = content
			}
			if !strings.Contains(changed[mutation.path], mutation.from) {
				t.Fatalf("runner.phase19-spec.%s.fixture: expected mutation target, actual missing", mutation.name)
			}
			changed[mutation.path] = strings.Replace(changed[mutation.path], mutation.from, mutation.to, 1)
			if err := phase19ValidateSpecification(changed); err == nil {
				t.Fatalf("runner.phase19-spec.%s.reject: expected rejected, actual accepted; specification mutation survived: %s", mutation.name, mutation.name)
			}
		})
	}
}

func phase19SchemaReferences(expression string, types map[string]json.RawMessage) ([]string, error) {
	tokens := regexp.MustCompile("[A-Za-z][A-Za-z0-9_.-]*|-?[0-9]+|[(),]").FindAllString(expression, -1)
	if strings.Join(tokens, "") != expression || len(tokens) == 0 {
		return nil, fmt.Errorf("invalid schema syntax %q", expression)
	}
	position := 0
	take := func() string {
		if position >= len(tokens) {
			return ""
		}
		value := tokens[position]
		position++
		return value
	}
	require := func(want string) error {
		if got := take(); got != want {
			return fmt.Errorf("expected %q, got %q in %s", want, got, expression)
		}
		return nil
	}
	number := func() (int64, error) {
		token := take()
		value, err := strconv.ParseInt(token, 10, 64)
		if err != nil || value < -9007199254740991 || value > 9007199254740991 || strconv.FormatInt(value, 10) != token {
			return 0, fmt.Errorf("invalid schema integer %q", token)
		}
		return value, nil
	}
	atoms := map[string]bool{"id": true, "jobid": true, "digest": true, "sha": true, "utc": true,
		"relpath": true, "abspath": true, "glob": true, "ref": true, "url": true,
		"scalar": true, "base64": true, "bool": true}
	references := []string{}
	var parse func(int) error
	parse = func(depth int) error {
		if depth > 32 {
			return fmt.Errorf("schema nesting too deep")
		}
		name := take()
		if atoms[name] {
			return nil
		}
		if types[name] != nil {
			references = append(references, name)
			return nil
		}
		if err := require("("); err != nil {
			return fmt.Errorf("unknown schema type/function %q: %w", name, err)
		}
		switch name {
		case "nullable", "list", "map", "envmap":
			if err := parse(depth + 1); err != nil {
				return err
			}
			if name != "nullable" {
				if err := require(","); err != nil {
					return err
				}
				minimum, err := number()
				if err != nil {
					return err
				}
				if err := require(","); err != nil {
					return err
				}
				maximum, err := number()
				if err != nil || minimum < 0 || maximum < minimum {
					return fmt.Errorf("invalid schema collection range")
				}
			}
		case "int", "text":
			minimum, err := number()
			if err != nil {
				return err
			}
			if err := require(","); err != nil {
				return err
			}
			maximum, err := number()
			if err != nil || maximum < minimum || (name == "text" && minimum < 0) {
				return fmt.Errorf("invalid schema scalar range")
			}
		case "oneof", "enum":
			seen := map[string]bool{}
			for {
				start := position
				if name == "oneof" {
					if err := parse(depth + 1); err != nil {
						return err
					}
				} else if !regexp.MustCompile("^[A-Za-z][A-Za-z0-9_.-]*$").MatchString(take()) {
					return fmt.Errorf("invalid enum value")
				}
				value := strings.Join(tokens[start:position], "")
				if seen[value] {
					return fmt.Errorf("duplicate schema alternative %s", value)
				}
				seen[value] = true
				if position >= len(tokens) || tokens[position] != "," {
					break
				}
				position++
			}
			if name == "oneof" && len(seen) < 2 {
				return fmt.Errorf("oneof requires two alternatives")
			}
		case "legacy":
			if name := take(); name != "AdminCredentials" && name != "TOTPSecret" {
				return fmt.Errorf("unknown legacy type")
			}
		case "const":
			if value := take(); value != "1" && value != "false" {
				return fmt.Errorf("unsupported schema constant")
			}
		default:
			return fmt.Errorf("unknown schema function %s", name)
		}
		return require(")")
	}
	if err := parse(0); err != nil {
		return nil, err
	}
	if position != len(tokens) {
		return nil, fmt.Errorf("trailing schema tokens in %s", expression)
	}
	return references, nil
}

func phase19ValidateSpecification(documents map[string]string) error {
	for path, content := range documents {
		anchors := regexp.MustCompile("<a id=\"(phase-19-[^\"]+)\"></a>").FindAllStringSubmatch(content, -1)
		seen := make(map[string]bool)
		for _, anchor := range anchors {
			if seen[anchor[1]] {
				return fmt.Errorf("duplicate Phase 19 anchor in %s: %s", path, anchor[1])
			}
			seen[anchor[1]] = true
		}
	}
	state := documents["docs/details/statefile.md"]
	block := regexp.MustCompile("(?s)```json\\n(.*?)\\n```").FindStringSubmatch(state)
	if len(block) != 2 {
		return fmt.Errorf("Phase 19 schema registry missing")
	}
	var registry struct {
		SchemaVersion int                        `json:"schema_version"`
		Profile       string                     `json:"profile"`
		Types         map[string]json.RawMessage `json:"types"`
	}
	if err := developmentActionsDecodeJSON([]byte(block[1]), &registry, 1048576, true); err != nil {
		return fmt.Errorf("Phase 19 schema registry: %w", err)
	}
	if registry.SchemaVersion != 1 || registry.Profile != "cicd-v1" || len(registry.Types) != 85 {
		return fmt.Errorf("Phase 19 registry identity or coverage invalid")
	}
	dependencies := make(map[string][]string)
	for name, raw := range registry.Types {
		if !regexp.MustCompile("^[A-Za-z][A-Za-z0-9]*$").MatchString(name) {
			return fmt.Errorf("invalid Phase 19 type name %s", name)
		}
		expressions := map[string]string{}
		var alias string
		if json.Unmarshal(raw, &alias) == nil {
			expressions[name] = alias
		} else if err := json.Unmarshal(raw, &expressions); err != nil || len(expressions) == 0 {
			return fmt.Errorf("invalid Phase 19 object schema %s", name)
		}
		for field, expression := range expressions {
			refs, err := phase19SchemaReferences(expression, registry.Types)
			if err != nil {
				return fmt.Errorf("%s.%s: %w", name, field, err)
			}
			dependencies[name] = append(dependencies[name], refs...)
		}
	}
	visiting, complete := map[string]bool{}, map[string]bool{}
	var visit func(string) error
	visit = func(name string) error {
		if visiting[name] {
			return fmt.Errorf("recursive Phase 19 schema %s", name)
		}
		if complete[name] {
			return nil
		}
		visiting[name] = true
		for _, dependency := range dependencies[name] {
			if err := visit(dependency); err != nil {
				return err
			}
		}
		visiting[name], complete[name] = false, true
		return nil
	}
	for name := range registry.Types {
		if err := visit(name); err != nil {
			return err
		}
	}
	for name, fields := range map[string]map[string]string{
		"Job":               {"env": "envmap(EnvValue,0,64)", "secrets": "envmap(SecretRef,0,16)"},
		"PlannedJob":        {"env": "envmap(EnvValue,0,64)", "secrets": "envmap(SecretRef,0,16)", "selected": "bool", "task": "PlannedTask"},
		"PlannedDeployTask": {"artifact_from": "Id"},
		"TriggerCursor":     {"definition_revision": "int(1,2147483647)"},
		"SiteState":         {"route": "Route", "operations": "list(SiteOperation,0,1024)", "counters": "SiteCounters"},
		"Root":              {"secrets": "list(SecretMeta,0,10000)"},
		"Lease":             {"project_id": "Id"},
		"JobResult":         {"execution_id": "nullable(Id)"},
		"WorkerExecution":   {"job": "JobResult", "lease_generation": "int(1,9007199254740991)", "assignment_bytes": "int(1,2097152)", "assignment_digest": "Digest", "toolchain_digest": "Digest"},
		"SiteOperation":     {"committed_at": "Time"},
		"Operation":         {"target_id": "nullable(text(1,64))", "target_attempt": "nullable(int(1,100))"},
		"Audit":             {"target_id": "nullable(text(1,64))"},
		"Deployment":        {"observations": "list(StageObservation,0,512)"},
		"Notification":      {"target_attempt": "int(1,100)", "source_revision": "int(1,2147483647)"},
	} {
		var actual map[string]string
		if err := json.Unmarshal(registry.Types[name], &actual); err != nil {
			return fmt.Errorf("required Phase 19 schema %s: %w", name, err)
		}
		for field, want := range fields {
			if actual[field] != want {
				return fmt.Errorf("Phase 19 %s.%s: want %s, got %s", name, field, want, actual[field])
			}
		}
	}
	if strings.Contains(documents["docs/ROADMAP.md"], "| 改訂予定 | CI / CD | Phase 19") {
		return fmt.Errorf("Phase 19 inventory maturity contradicts phase plan")
	}
	api := documents["docs/details/api.md"]
	section := strings.Split(api, "<a id=\"phase-19-error-catalog\"></a>")[0]
	routes := regexp.MustCompile("(?m)^\\| ([A-Za-z][A-Za-z0-9]*) \\| (GET|POST|PUT|DELETE) `([^\x60]+)` \\|").FindAllStringSubmatch(section, -1)
	if len(routes) != 53 {
		return fmt.Errorf("Phase 19 operation catalog incomplete")
	}
	seenIDs := make(map[string]bool)
	seenRoutes := make(map[string]bool)
	for _, route := range routes {
		id, wire := route[1], route[2]+" "+route[3]
		if seenIDs[id] || seenRoutes[wire] {
			return fmt.Errorf("duplicate Phase 19 operation %s", id)
		}
		seenIDs[id], seenRoutes[wire] = true, true
		if !strings.Contains(documents["docs/details/sdk.md"], "| `"+id+"(params,body,options)` | "+id+" |") {
			return fmt.Errorf("Phase 19 SDK method missing: %s", id)
		}
		if !regexp.MustCompile("(?m)^\\| " + id + " \\| [a-z]+ \\| `cicd-[a-z-]+` \\|$").MatchString(documents["docs/details/ui.md"]) {
			return fmt.Errorf("Phase 19 UI action missing: %s", id)
		}
		if !regexp.MustCompile("(?m)^\\| `[a-z-]+` \\| " + id + " \\|$").MatchString(documents["docs/details/admin.md"]) {
			return fmt.Errorf("Phase 19 CLI command missing: %s", id)
		}
	}
	if err := plannedPhaseMaturityBoundary(documents["docs/ROADMAP.md"]); err != nil {
		return err
	}
	fixture := documents["docs/details/fixture.md"]
	for _, id := range []string{
		"PLAN", "EVENT", "DAG", "RETRY", "CANCEL", "WORKER", "ISOLATION", "CACHE",
		"ARTIFACT", "APPROVAL", "DEPLOY", "DURABILITY", "ACCESS", "LOG-UI", "RELEASE", "STATUS",
	} {
		if !strings.Contains(fixture, "| P19-"+id+" |") {
			return fmt.Errorf("Phase 19 fixture contract missing: %s", id)
		}
	}
	return nil
}

func TestDispatchMainExactBasename(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if code := dispatchMain("prefix-adlaire-ci-runner", []string{"--version"}, &stdout, &stderr); code != 2 {
		t.Fatalf("expected unknown command exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout for unknown command, got %q", stdout.String())
	}
	if got := stderr.String(); !strings.Contains(got, "unknown command: prefix-adlaire-ci-runner") {
		t.Fatalf("expected unknown command error, got %q", got)
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-build", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected builder help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-build") {
		t.Fatalf("expected builder help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for builder help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-runner", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected runner help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-runner") {
		t.Fatalf("expected runner help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for runner help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-api", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected api help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-api") {
		t.Fatalf("expected api help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for api help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-admin", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected admin help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-admin") {
		t.Fatalf("expected admin help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for admin help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-setup", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected setup help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-setup") {
		t.Fatalf("expected setup help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for setup help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-release", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected release help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-release") {
		t.Fatalf("expected release help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for release help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-mcp", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected mcp help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-mcp") {
		t.Fatalf("expected mcp help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for mcp help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-obsidian", []string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected obsidian help exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.Contains(got, "Usage: adlaire-ci-obsidian") {
		t.Fatalf("expected obsidian help usage, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for obsidian help, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-build-linux-amd64", []string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected release asset build version exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.HasPrefix(got, "adlaire-ci-build V.0.0-dev go=") {
		t.Fatalf("expected canonical build version, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for release asset build version, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-mcp-linux-amd64", []string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected release asset mcp version exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.HasPrefix(got, "adlaire-ci-mcp V.0.0-dev go=") {
		t.Fatalf("expected canonical mcp version, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for release asset mcp version, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	if code := dispatchMain("adlaire-ci-obsidian-linux-amd64", []string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected release asset obsidian version exit code 0, got %d", code)
	}
	if got := stdout.String(); !strings.HasPrefix(got, "adlaire-ci-obsidian V.0.0-dev go=") {
		t.Fatalf("expected canonical obsidian version, got %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for release asset obsidian version, got %q", stderr.String())
	}
}

func TestPhase18EmbeddedAdminMaterialization(t *testing.T) {
	installDir := filepath.Join(t.TempDir(), "runtime")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := dispatchMain("adlaire-ci-setup", []string{"phase18-materialize-admin", "--install-dir", installDir}, &stdout, &stderr)
	if code != 0 || stdout.String() != "setup: success phase18-materialize-admin\n" || stderr.Len() != 0 {
		t.Fatalf("materialize code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for _, name := range []string{"index.html", "adlaire-ci-sdk.js"} {
		want, err := os.ReadFile(filepath.Join("admin", name))
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(installDir, "admin", name))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("materialized %s differs from canonical asset", name)
		}
	}
}

type phase18BootstrapHarness struct {
	root       string
	scriptPath string
	sourcePath string
	stateDir   string
	binDir     string
	systemdDir string
	lockDir    string
	fakeBinDir string
	tmpDir     string
	digest     string
}

func phase18WriteExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		t.Fatal(err)
	}
}

func phase18FakeSource(version string) string {
	return `#!/bin/sh
name=${0##*/}
if [ "${1:-}" = "--version" ]; then
  printf '%s ` + version + ` go=fake\n' "$name"
  exit 0
fi
case "$name" in
  adlaire-ci-setup)
    if [ "${1:-}" != "phase18-materialize-admin" ] || [ "${2:-}" != "--install-dir" ] || [ -z "${3:-}" ]; then exit 2; fi
    mkdir -p "$3/admin" || exit 1
    printf '<html>phase18</html>\n' >"$3/admin/index.html" || exit 1
    printf 'export const phase18 = true;\n' >"$3/admin/adlaire-ci-sdk.js" || exit 1
    chmod 0644 "$3/admin/index.html" "$3/admin/adlaire-ci-sdk.js" || exit 1
    printf 'setup: success phase18-materialize-admin\n'
    ;;
  adlaire-ci-api)
    state=""
    while [ "$#" -gt 0 ]; do
      case "$1" in
        --state-dir) state="$2"; shift 2 ;;
        --init-credentials) shift ;;
        *) shift ;;
      esac
    done
    [ -n "$state" ] || exit 2
    IFS= read -r password || exit 1
    [ -n "$password" ] || exit 1
    mkdir -p "$state" || exit 1
    printf '{"password_hash":"fake","salt":"fake","algorithm":"pbkdf2-sha256","iterations":1,"updated_at":"2026-10-06T00:00:00Z","last_login_at":null,"login_count":0}\n' >"$state/.admin_credentials" || exit 1
    chmod 0600 "$state/.admin_credentials" || exit 1
    printf 'credentials initialized\n'
    ;;
  adlaire-ci-admin)
    printf '{"last_build_status":"success","running":false}\n'
    ;;
  adlaire-ci-runner)
    state=""
    while [ "$#" -gt 0 ]; do
      case "$1" in
        --state-dir) state="$2"; shift 2 ;;
        *) shift ;;
      esac
    done
    [ -n "$state" ] || exit 2
    out=$(sed -n 's/.*"out":"\([^"]*\)".*/\1/p' "$state/.branch_config")
    [ -n "$out" ] || exit 1
    mkdir -p "$out" || exit 1
    printf '<html>built</html>\n' >"$out/index.html" || exit 1
    printf '{}\n' >"$out/.dependency_manifest.json" || exit 1
    printf '{"id":"phase18","status":"success"}\n' >"$state/.build_history" || exit 1
    printf '{"status":"success"}\n' >"$state/.build_status.json" || exit 1
    ;;
  *) exit 2 ;;
esac
`
}

func newPhase18BootstrapHarness(t *testing.T) *phase18BootstrapHarness {
	t.Helper()
	root := t.TempDir()
	h := &phase18BootstrapHarness{
		root:       root,
		scriptPath: filepath.Join(root, "adlaire-ci-vps-pull-bootstrap.sh"),
		sourcePath: filepath.Join(root, "adlaire-ci-source"),
		stateDir:   filepath.Join(root, "state"),
		binDir:     filepath.Join(root, "bin"),
		systemdDir: filepath.Join(root, "systemd"),
		lockDir:    filepath.Join(root, "bootstrap.lock"),
		fakeBinDir: filepath.Join(root, "fake-bin"),
		tmpDir:     filepath.Join(root, "tmp"),
	}
	for _, dir := range []string{h.binDir, h.systemdDir, h.fakeBinDir, h.tmpDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	phase18WriteExecutable(t, h.sourcePath, phase18FakeSource("V.18.0"))
	source, err := os.ReadFile(h.sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(source)
	h.digest = hex.EncodeToString(sum[:])

	script, err := os.ReadFile(phase18BootstrapScriptPath)
	if err != nil {
		t.Fatal(err)
	}
	transformed := strings.NewReplacer(
		"/run/lock/adlaire-ci-vps-pull-bootstrap.lock", h.lockDir,
		"/etc/systemd/system", h.systemdDir,
		"/opt/adlaire-builder", h.stateDir,
		"/usr/local/bin", h.binDir,
	).Replace(string(script))
	phase18WriteExecutable(t, h.scriptPath, transformed)

	phase18WriteExecutable(t, filepath.Join(h.fakeBinDir, "systemctl"), `#!/bin/sh
state=$PHASE18_FAKE_SYSTEMD_STATE
mkdir -p "$state" || exit 1
command_name=${1:-}
shift || true
case "$command_name" in
  is-enabled) [ -f "$state/enabled-$1" ] ;;
  is-active) [ -f "$state/active-$1" ] ;;
  enable) for unit in "$@"; do : >"$state/enabled-$unit" || exit 1; done ;;
  disable) for unit in "$@"; do rm -f "$state/enabled-$unit" || exit 1; done ;;
  start) for unit in "$@"; do : >"$state/active-$unit" || exit 1; done ;;
  stop) for unit in "$@"; do rm -f "$state/active-$unit" || exit 1; done ;;
  restart)
    if [ "${PHASE18_FAKE_FAIL_API_RESTART_ONCE:-0}" = "1" ] && [ "${1:-}" = "adlaire-ci-api.service" ] && [ ! -f "$state/restart-failed" ]; then
      : >"$state/restart-failed"
      exit 1
    fi
    for unit in "$@"; do : >"$state/active-$unit" || exit 1; done
    ;;
  cat|daemon-reload) exit 0 ;;
  *) exit 1 ;;
esac
`)
	phase18WriteExecutable(t, filepath.Join(h.fakeBinDir, "cp"), `#!/bin/sh
if [ "${1:-}" = "-a" ] && [ "${2:-}" = "$PHASE18_TEST_STATE" ]; then
  for unit in adlaire-ci-api.service adlaire-ci.timer adlaire-ci.service; do
    [ ! -f "$PHASE18_FAKE_SYSTEMD_STATE/active-$unit" ] || exit 1
  done
fi
exec /bin/cp "$@"
`)
	phase18WriteExecutable(t, filepath.Join(h.fakeBinDir, "rm"), `#!/bin/sh
if [ "${PHASE18_FAKE_FAIL_STAGE_CLEANUP_ONCE:-0}" = "1" ] && [ "${1:-}" = "-rf" ]; then
  case "${2:-}" in
    "$TMPDIR"/adlaire-ci-bootstrap.*)
      if [ ! -f "$PHASE18_FAKE_CLEANUP_MARKER" ]; then
        : >"$PHASE18_FAKE_CLEANUP_MARKER" || exit 1
        exit 1
      fi
      ;;
  esac
fi
exec /bin/rm "$@"
`)
	phase18WriteExecutable(t, filepath.Join(h.fakeBinDir, "sleep"), "#!/bin/sh\nexit 0\n")
	phase18WriteExecutable(t, filepath.Join(h.fakeBinDir, "curl"), `#!/bin/sh
config=${2:-}
[ "$1" = "--config" ] && [ -f "$config" ] || exit 2
output=$(sed -n 's/^output = "\(.*\)"$/\1/p' "$config")
url=$(sed -n 's/^url = "\(.*\)"$/\1/p' "$config")
data=$(sed -n 's/^data-binary = "@\(.*\)"$/\1/p' "$config")
[ -n "$output" ] || exit 2
mkdir -p "$PHASE18_TEST_STATE" || exit 1
for log in .access_log .audit_log .api_access_log .config_log; do printf '{}\n' >>"$PHASE18_TEST_STATE/$log" || exit 1; done
case "$url" in
  */api/health)
    if [ -n "${PHASE18_FAKE_HEALTH_BODY:-}" ]; then
      printf '%s\n' "$PHASE18_FAKE_HEALTH_BODY" >"$output"
    else
      printf '{"checks":[],"last_build_at":null,"last_build_status":"none","last_deploy_at":null,"last_deploy_status":null,"pending_transfers":0,"status":"ok","uptime_seconds":1}\n' >"$output"
    fi
    ;;
  */api/login)
    printf 'login\n' >>"$PHASE18_TEST_SESSION_EVENTS" || exit 1
    printf '{"must_change":"none","token":"acs_0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}\n' >"$output"
    ;;
  */api/change-password) printf '{"message":"Password changed"}\n' >"$output" ;;
  */api/logout)
    printf 'logout\n' >>"$PHASE18_TEST_SESSION_EVENTS" || exit 1
    printf '{"message":"Logged out"}\n' >"$output"
    ;;
  */api/schedule/cooldown)
    if grep -q '"seconds":0' "$data"; then seconds=0; else seconds=2; fi
    printf '{"message":"Build cooldown updated","seconds":%s}\n' "$seconds" >"$output"
    ;;
  */admin/adlaire-ci-sdk.js) cp "$PHASE18_TEST_STATE/admin/adlaire-ci-sdk.js" "$output" ;;
  */) cp "$PHASE18_TEST_STATE/admin/index.html" "$output" ;;
  *) exit 22 ;;
esac
`)
	phase18WriteExecutable(t, filepath.Join(h.fakeBinDir, "sha256sum"), `#!/bin/sh
if [ -x /usr/bin/sha256sum ]; then exec /usr/bin/sha256sum "$@"; fi
if [ -x /usr/bin/openssl ]; then exec /usr/bin/openssl dgst -sha256 -r "$@"; fi
exit 127
`)
	return h
}

func (h *phase18BootstrapHarness) args() []string {
	return []string{
		h.scriptPath,
		"--source-channel", "integration-head",
		"--topology-role", "single-node",
		"--source", h.sourcePath,
		"--sha256", h.digest,
		"--install-dir", h.stateDir,
		"--bin-dir", h.binDir,
		"--state-dir", h.stateDir,
		"--service-user", "root",
	}
}

func (h *phase18BootstrapHarness) run(input string, extraEnv ...string) (int, string, string) {
	cmd := exec.Command("sh", h.args()...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Env = append(os.Environ(),
		"PATH="+h.fakeBinDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"TMPDIR="+h.tmpDir,
		"PHASE18_FAKE_SYSTEMD_STATE="+filepath.Join(h.root, "systemctl-state"),
		"PHASE18_TEST_STATE="+h.stateDir,
		"PHASE18_TEST_SESSION_EVENTS="+filepath.Join(h.root, "session-events"),
		"PHASE18_FAKE_CLEANUP_MARKER="+filepath.Join(h.root, "cleanup-failed-once"),
	)
	cmd.Env = append(cmd.Env, extraEnv...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return 0, stdout.String(), stderr.String()
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), stdout.String(), stderr.String()
	}
	return -1, stdout.String(), stderr.String()
}

func phase18SeedOperationalState(t *testing.T, h *phase18BootstrapHarness) map[string][]byte {
	t.Helper()
	if err := os.MkdirAll(h.stateDir, 0700); err != nil {
		t.Fatal(err)
	}
	values := map[string][]byte{
		".admin_credentials": []byte("existing-credentials\n"),
		".server_config":     []byte("{\"watch_mode\":\"local\",\"snapshots_keep\":0,\"queue_max_size\":0,\"build_cooldown_seconds\":37}\n"),
		".pipeline_config":   []byte("{\"extra_args\":[\"--strict\"],\"env\":{\"MODE\":\"existing\"}}\n"),
		".branch_config":     []byte("{\"branch_targets\":[]}\n"),
		".last_sha":          []byte("{\"sha\":\"existing-sha\"}\n"),
		".pending_transfers": []byte("[{\"id\":\"existing-transfer\"}]\n"),
		".notify_pending":    []byte("[{\"id\":\"existing-notification\"}]\n"),
		"preserve.txt":       []byte("preserve-me\n"),
	}
	for name, data := range values {
		if err := os.WriteFile(filepath.Join(h.stateDir, name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	return values
}

func phase18ReadDiagnostic(t *testing.T, stderr string) map[string]any {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	if len(lines) != 1 {
		t.Fatalf("bootstrap failure must emit one diagnostic JSON line: %q", stderr)
	}
	var diagnostic map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &diagnostic); err != nil {
		t.Fatalf("invalid bootstrap diagnostic: %v output=%q", err, stderr)
	}
	return diagnostic
}

func TestPhase18BootstrapAcceptsExactArgumentVector(t *testing.T) {
	for _, test := range []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "missing-line-feed", input: "ValidPassword123"},
		{name: "trailing-data-without-line-feed", input: "ValidPassword123\ntrailing"},
		{name: "second-line", input: "ValidPassword123\ntrailing\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			h := newPhase18BootstrapHarness(t)
			code, stdout, stderr := h.run(test.input)
			if code != 2 || stdout != "" {
				t.Fatalf("stdin boundary code=%d stdout=%q stderr=%q", code, stdout, stderr)
			}
			diagnostic := phase18ReadDiagnostic(t, stderr)
			if diagnostic["failure_code"] != "bootstrap-script-health-failed" || diagnostic["stage"] != "admin-password-stdin" {
				t.Fatalf("exact argument vector must reach strict stdin boundary: %#v", diagnostic)
			}
			if _, err := os.Lstat(h.lockDir); !os.IsNotExist(err) {
				t.Fatalf("failed bootstrap retained process lock: %v", err)
			}
		})
	}
}

func phase18TreeSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	snapshot := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(root, path)
		if err != nil || rel == "." {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		value := info.Mode().String()
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			value += ":" + target
		} else if info.Mode().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			sum := sha256.Sum256(data)
			value += ":" + hex.EncodeToString(sum[:])
		}
		snapshot[filepath.ToSlash(rel)] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func TestPhase18BootstrapPreservesOperationalStateAndClosesEvidence(t *testing.T) {
	h := newPhase18BootstrapHarness(t)
	preserved := phase18SeedOperationalState(t, h)
	code, stdout, stderr := h.run("ValidPassword123\n")
	if code != 0 || stderr != "" {
		t.Fatalf("bootstrap success code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	var completion map[string]any
	if err := json.Unmarshal(bytes.TrimSpace([]byte(stdout)), &completion); err != nil {
		t.Fatalf("bootstrap completion must be one JSON object: %v output=%q", err, stdout)
	}
	if completion["status"] != "ok" || completion["final_open_item_result"] != float64(0) || completion["cleanup_result"] != "passed" {
		t.Fatalf("unexpected bootstrap completion: %#v", completion)
	}
	for name, want := range preserved {
		if name == "preserve.txt" || strings.HasPrefix(name, ".") {
			got, err := os.ReadFile(filepath.Join(h.stateDir, name))
			if err != nil || !bytes.Equal(got, want) {
				t.Fatalf("operational state %s changed got=%q want=%q err=%v", name, string(got), string(want), err)
			}
		}
	}
	runtimeStateDir := filepath.Join(h.stateDir, ".phase18-runtime")
	var runtimeBranchConfig struct {
		BranchTargets []struct {
			SHAFile string `json:"sha_file"`
			Src     string `json:"src"`
			Out     string `json:"out"`
		} `json:"branch_targets"`
	}
	runtimeBranchData, err := os.ReadFile(filepath.Join(runtimeStateDir, ".branch_config"))
	if err != nil || json.Unmarshal(runtimeBranchData, &runtimeBranchConfig) != nil || len(runtimeBranchConfig.BranchTargets) != 1 {
		t.Fatalf("invalid isolated runtime branch config: %v data=%q", err, string(runtimeBranchData))
	}
	runtimeTarget := runtimeBranchConfig.BranchTargets[0]
	if runtimeTarget.SHAFile != filepath.Join(runtimeStateDir, ".last_sha") || runtimeTarget.Src != filepath.Join(runtimeStateDir, "phase18-source") || runtimeTarget.Out != filepath.Join(runtimeStateDir, "phase18-site") {
		t.Fatalf("runtime branch paths escaped isolated state: %+v", runtimeTarget)
	}
	for _, name := range []string{"open.json", "runtime.json", "update-rollback.json", "issue-triage.json", "spec-first-fix.json", "revalidation.json", "close.json"} {
		if info, err := os.Lstat(filepath.Join(h.stateDir, ".phase18-evidence", name)); err != nil || !info.Mode().IsRegular() || info.Mode().Perm() != 0600 {
			t.Fatalf("missing or unsafe evidence %s info=%v err=%v", name, info, err)
		}
	}
	var drill map[string]any
	drillData, err := os.ReadFile(filepath.Join(h.stateDir, ".phase18-evidence", "update-rollback.json"))
	if err != nil || json.Unmarshal(drillData, &drill) != nil {
		t.Fatalf("invalid update/rollback evidence: %v data=%q", err, string(drillData))
	}
	if drill["start_digest"] == drill["update_digest"] || drill["source_digest"] != h.digest || drill["result"] != "passed" {
		t.Fatalf("update/rollback drill did not use distinct cohorts: %#v", drill)
	}
	if _, err := os.Lstat(h.lockDir); !os.IsNotExist(err) {
		t.Fatalf("successful bootstrap retained process lock: %v", err)
	}
	sessionEvents, err := os.ReadFile(filepath.Join(h.root, "session-events"))
	if err != nil || string(sessionEvents) != "login\nlogout\nlogin\nlogout\n" {
		t.Fatalf("bootstrap session tokens were not balanced events=%q err=%v", string(sessionEvents), err)
	}
	entries, err := os.ReadDir(h.tmpDir)
	if err != nil || len(entries) != 0 {
		t.Fatalf("successful bootstrap retained staging content entries=%v err=%v", entries, err)
	}
}

func TestPhase18BootstrapRollsBackWhenFinalCleanupFails(t *testing.T) {
	h := newPhase18BootstrapHarness(t)
	phase18SeedOperationalState(t, h)
	stateBefore := phase18TreeSnapshot(t, h.stateDir)
	binBefore := phase18TreeSnapshot(t, h.binDir)
	unitBefore := phase18TreeSnapshot(t, h.systemdDir)

	code, stdout, stderr := h.run("ValidPassword123\n", "PHASE18_FAKE_FAIL_STAGE_CLEANUP_ONCE=1")
	if code != 1 || stdout != "" {
		t.Fatalf("cleanup failure path code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	diagnostic := phase18ReadDiagnostic(t, stderr)
	if diagnostic["failure_code"] != "bootstrap-script-health-failed" || diagnostic["stage"] != "cleanup" || diagnostic["rollback_action"] != "restored" || diagnostic["cleanup_result"] != "passed" {
		t.Fatalf("unexpected cleanup rollback diagnostic: %#v", diagnostic)
	}
	if got := phase18TreeSnapshot(t, h.stateDir); !mapsEqual(stateBefore, got) {
		t.Fatalf("cleanup failure did not restore state\nwant=%#v\n got=%#v", stateBefore, got)
	}
	if got := phase18TreeSnapshot(t, h.binDir); !mapsEqual(binBefore, got) {
		t.Fatalf("cleanup failure did not restore binaries\nwant=%#v\n got=%#v", binBefore, got)
	}
	if got := phase18TreeSnapshot(t, h.systemdDir); !mapsEqual(unitBefore, got) {
		t.Fatalf("cleanup failure did not restore units\nwant=%#v\n got=%#v", unitBefore, got)
	}
	if _, err := os.Lstat(h.lockDir); !os.IsNotExist(err) {
		t.Fatalf("cleanup rollback retained process lock: %v", err)
	}
}

func TestPhase18BootstrapRestoresFullStateAfterCommitFailure(t *testing.T) {
	h := newPhase18BootstrapHarness(t)
	phase18SeedOperationalState(t, h)
	oldSource := phase18FakeSource("V.17.0")
	for _, name := range phase18BootstrapInstalledBinaries {
		phase18WriteExecutable(t, filepath.Join(h.binDir, name), oldSource)
	}
	for _, name := range []string{"adlaire-ci.service", "adlaire-ci.timer", "adlaire-ci-api.service"} {
		if err := os.WriteFile(filepath.Join(h.systemdDir, name), []byte("old-"+name+"\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	fakeSystemd := filepath.Join(h.root, "systemctl-state")
	if err := os.MkdirAll(fakeSystemd, 0755); err != nil {
		t.Fatal(err)
	}
	for _, unit := range []string{"adlaire-ci.service", "adlaire-ci.timer", "adlaire-ci-api.service"} {
		for _, prefix := range []string{"enabled-", "active-"} {
			if err := os.WriteFile(filepath.Join(fakeSystemd, prefix+unit), []byte{}, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	stateBefore := phase18TreeSnapshot(t, h.stateDir)
	binBefore := phase18TreeSnapshot(t, h.binDir)
	unitBefore := phase18TreeSnapshot(t, h.systemdDir)

	code, stdout, stderr := h.run("ValidPassword123\n", "PHASE18_FAKE_FAIL_API_RESTART_ONCE=1")
	if code != 1 || stdout != "" {
		t.Fatalf("rollback failure path code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	diagnostic := phase18ReadDiagnostic(t, stderr)
	if diagnostic["failure_code"] != "bootstrap-script-health-failed" || diagnostic["stage"] != "api-service-start" || diagnostic["rollback_action"] != "restored" || diagnostic["cleanup_result"] != "passed" {
		t.Fatalf("unexpected rollback diagnostic: %#v", diagnostic)
	}
	if got := phase18TreeSnapshot(t, h.stateDir); !mapsEqual(stateBefore, got) {
		t.Fatalf("state tree was not fully restored\nwant=%#v\n got=%#v", stateBefore, got)
	}
	if got := phase18TreeSnapshot(t, h.binDir); !mapsEqual(binBefore, got) {
		t.Fatalf("binary cohort was not restored\nwant=%#v\n got=%#v", binBefore, got)
	}
	if got := phase18TreeSnapshot(t, h.systemdDir); !mapsEqual(unitBefore, got) {
		t.Fatalf("systemd units were not restored\nwant=%#v\n got=%#v", unitBefore, got)
	}
	for _, unit := range []string{"adlaire-ci.service", "adlaire-ci.timer", "adlaire-ci-api.service"} {
		for _, prefix := range []string{"enabled-", "active-"} {
			if _, err := os.Stat(filepath.Join(fakeSystemd, prefix+unit)); err != nil {
				t.Fatalf("service state %s%s was not restored: %v", prefix, unit, err)
			}
		}
	}
	if _, err := os.Lstat(h.lockDir); !os.IsNotExist(err) {
		t.Fatalf("rollback retained process lock: %v", err)
	}
}

func TestPhase18BootstrapRejectsMalformedRuntimeJSON(t *testing.T) {
	h := newPhase18BootstrapHarness(t)
	phase18SeedOperationalState(t, h)
	code, stdout, stderr := h.run("ValidPassword123\n", `PHASE18_FAKE_HEALTH_BODY={"note":"status ok"}`)
	if code != 1 || stdout != "" {
		t.Fatalf("malformed health code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	diagnostic := phase18ReadDiagnostic(t, stderr)
	if diagnostic["failure_code"] != "bootstrap-script-health-failed" || diagnostic["stage"] != "api-health-open" || diagnostic["rollback_action"] != "restored" {
		t.Fatalf("unexpected malformed health diagnostic: %#v", diagnostic)
	}
	if _, err := os.Lstat(filepath.Join(h.stateDir, ".phase18-evidence", "close.json")); !os.IsNotExist(err) {
		t.Fatalf("failed runtime validation created close evidence: %v", err)
	}
}

func TestPhase18BootstrapRejectsNonRegularCredentialState(t *testing.T) {
	h := newPhase18BootstrapHarness(t)
	if err := os.MkdirAll(filepath.Join(h.stateDir, ".admin_credentials"), 0700); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := h.run("ValidPassword123\n")
	if code != 1 || stdout != "" {
		t.Fatalf("credential type code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	diagnostic := phase18ReadDiagnostic(t, stderr)
	if diagnostic["failure_code"] != "bootstrap-script-state-directory-invalid" || diagnostic["stage"] != "credentials-type" || diagnostic["rollback_action"] != "not_required" {
		t.Fatalf("unexpected credential type diagnostic: %#v", diagnostic)
	}
	if info, err := os.Stat(filepath.Join(h.stateDir, ".admin_credentials")); err != nil || !info.IsDir() {
		t.Fatalf("invalid credential state changed info=%v err=%v", info, err)
	}
}

func TestPhase18BootstrapProcessLockRejectsLiveOwnerAndRecoversStaleOwner(t *testing.T) {
	t.Run("live-owner", func(t *testing.T) {
		h := newPhase18BootstrapHarness(t)
		if err := os.Mkdir(h.lockDir, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(h.lockDir, "owner"), []byte(fmt.Sprintf("%d\n", os.Getpid())), 0600); err != nil {
			t.Fatal(err)
		}
		code, stdout, stderr := h.run("ValidPassword123\n")
		if code != 1 || stdout != "" {
			t.Fatalf("live lock code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		diagnostic := phase18ReadDiagnostic(t, stderr)
		if diagnostic["stage"] != "process-lock" {
			t.Fatalf("unexpected live lock diagnostic: %#v", diagnostic)
		}
		if _, err := os.Stat(filepath.Join(h.lockDir, "owner")); err != nil {
			t.Fatalf("foreign live lock was modified: %v", err)
		}
	})

	t.Run("malformed-owner", func(t *testing.T) {
		h := newPhase18BootstrapHarness(t)
		if err := os.Mkdir(h.lockDir, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(h.lockDir, "owner"), []byte("99999999\n123\n"), 0600); err != nil {
			t.Fatal(err)
		}
		code, stdout, stderr := h.run("ValidPassword123\n")
		if code != 1 || stdout != "" {
			t.Fatalf("malformed lock code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		diagnostic := phase18ReadDiagnostic(t, stderr)
		if diagnostic["stage"] != "process-lock" {
			t.Fatalf("unexpected malformed lock diagnostic: %#v", diagnostic)
		}
		data, err := os.ReadFile(filepath.Join(h.lockDir, "owner"))
		if err != nil || string(data) != "99999999\n123\n" {
			t.Fatalf("malformed foreign lock was modified data=%q err=%v", string(data), err)
		}
	})

	t.Run("stale-owner", func(t *testing.T) {
		h := newPhase18BootstrapHarness(t)
		if err := os.Mkdir(h.lockDir, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(h.lockDir, "owner"), []byte("99999999\n"), 0600); err != nil {
			t.Fatal(err)
		}
		staleStage := filepath.Join(h.tmpDir, "adlaire-ci-bootstrap.stale")
		if err := os.Mkdir(staleStage, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(staleStage, "partial"), []byte("stale"), 0600); err != nil {
			t.Fatal(err)
		}
		code, stdout, stderr := h.run("")
		if code != 2 || stdout != "" {
			t.Fatalf("stale lock recovery code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		diagnostic := phase18ReadDiagnostic(t, stderr)
		if diagnostic["stage"] != "admin-password-stdin" {
			t.Fatalf("stale lock was not recovered before stdin validation: %#v", diagnostic)
		}
		if _, err := os.Lstat(h.lockDir); !os.IsNotExist(err) {
			t.Fatalf("stale lock recovery retained lock: %v", err)
		}
		if _, err := os.Lstat(staleStage); !os.IsNotExist(err) {
			t.Fatalf("stale staging cleanup retained path: %v", err)
		}
	})
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if b[key] != value {
			return false
		}
	}
	return true
}

func TestPhase12StandardArtifactInventory(t *testing.T) {
	t.Parallel()

	ownerPackages := []string{
		"admin",
		"api",
		"archive",
		"builder",
		"commitstatus",
		"mcp",
		"obsidian",
		"release",
		"runner",
		"security",
		"setup",
		"statefile",
	}
	allowedOwnerFiles := map[string][]string{}
	for _, owner := range ownerPackages {
		allowedOwnerFiles[owner] = []string{
			owner + ".go",
			"model.go",
			"validate.go",
			"execute.go",
			owner + "_test.go",
		}
	}

	requiredFiles := []string{
		"main.go",
		"main_test.go",
		"sdk_contract_test.go",
		"ui_contract_test.go",
		"admin/adlaire-ci-sdk.js",
		"admin/index.html",
		".github/workflows/phase12-quality-gate.yml",
	}
	for _, owner := range ownerPackages {
		for _, file := range allowedOwnerFiles[owner] {
			requiredFiles = append(requiredFiles, filepath.ToSlash(filepath.Join("components", owner, file)))
		}
	}
	for _, path := range requiredFiles {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("required Phase 12 artifact %s is missing: %v", path, err)
		}
		if info.IsDir() {
			t.Fatalf("required Phase 12 artifact %s must be a file", path)
		}
	}

	if _, err := os.Stat("cmd"); !os.IsNotExist(err) {
		t.Fatalf("cmd directory must not exist in the standard Phase 12 layout: %v", err)
	}

	entries, err := os.ReadDir("components")
	if err != nil {
		t.Fatalf("read components directory: %v", err)
	}
	var ownerDirs []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			t.Fatalf("components root must not contain Go file after Phase 12 owner package split: %s", entry.Name())
		}
		if entry.IsDir() {
			ownerDirs = append(ownerDirs, entry.Name())
		}
	}
	sort.Strings(ownerDirs)
	if !phase11StringSlicesEqual(ownerDirs, ownerPackages) {
		t.Fatalf("unexpected components owner directories\nwant: %v\n got: %v", ownerPackages, ownerDirs)
	}

	for _, owner := range ownerPackages {
		entries, err := os.ReadDir(filepath.Join("components", owner))
		if err != nil {
			t.Fatalf("read owner package %s: %v", owner, err)
		}
		var goFiles []string
		for _, entry := range entries {
			if entry.IsDir() {
				t.Fatalf("owner package %s must not contain subdirectory %s", owner, entry.Name())
			}
			if strings.HasSuffix(entry.Name(), ".go") {
				goFiles = append(goFiles, entry.Name())
			}
		}
		sort.Strings(goFiles)
		want := append([]string(nil), allowedOwnerFiles[owner]...)
		sort.Strings(want)
		if !phase11StringSlicesEqual(goFiles, want) {
			t.Fatalf("owner package %s must use exactly five fixed files\nwant: %v\n got: %v", owner, want, goFiles)
		}
	}

	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, path := range requiredFiles {
		if !strings.Contains(documentIndex, path) {
			t.Fatalf("docs/DOCUMENT_INDEX.md does not index required artifact %s", path)
		}
	}
}

func TestPhase12QualityGateEvidence(t *testing.T) {
	t.Parallel()

	workflow := phase11MustReadText(t, ".github/workflows/phase12-quality-gate.yml")
	for _, job := range []string{
		"phase12-go-format:",
		"phase12-go-test:",
		"phase12-go-race:",
		"phase12-go-vet:",
		"phase12-deno-check-sdk:",
		"phase12-fixture-harness:",
		"phase12-mutation:",
	} {
		if !strings.Contains(workflow, job) {
			t.Fatalf("Phase 12 workflow must include required job %s", strings.TrimSuffix(job, ":"))
		}
	}

	effectsPath := "testdata/phase12/quality-gate-reconstruction/expected/effects.json"
	var effects map[string]int
	data, err := os.ReadFile(effectsPath)
	if err != nil {
		t.Fatalf("read Phase 12 effects: %v", err)
	}
	if err := json.Unmarshal(data, &effects); err != nil {
		t.Fatalf("decode Phase 12 effects: %v", err)
	}
	for _, key := range []string{
		"phase12_contract_mismatch_count",
		"phase12_state_safety_open_count",
		"phase12_fixture_harness_open_count",
		"phase12_mutation_survived_count",
		"phase12_race_trigger_open_count",
		"phase12_owner_package_violation_count",
		"phase12_ci_required_check_open_count",
		"phase12_release_reproducibility_open_count",
		"final_open_item_count",
	} {
		value, ok := effects[key]
		if !ok {
			t.Fatalf("Phase 12 effects must include %s", key)
		}
		if value != 0 {
			t.Fatalf("Phase 12 effects %s must be 0, got %d", key, value)
		}
	}
}

func TestPhase11FixtureManifestGate(t *testing.T) {
	t.Parallel()

	phase11RequireNoUnexpectedFixtureDirectories(t)

	manifestPaths, err := phase11FindManifestFiles("testdata")
	if err != nil {
		t.Fatalf("find fixture manifests: %v", err)
	}
	if len(manifestPaths) < 56 {
		t.Fatalf("expected at least 56 formal fixture manifests, got %d", len(manifestPaths))
	}

	seenNames := make(map[string]string)
	for _, manifestPath := range manifestPaths {
		manifest := phase11ReadManifest(t, manifestPath)
		fixtureDir := filepath.Dir(manifestPath)

		if manifest.Name != filepath.Base(fixtureDir) && !phase13AllowsScopeName(manifestPath, manifest.Name) && !phase14AllowsScopeName(manifestPath, manifest.Name) && !phase15AllowsScopeName(manifestPath, manifest.Name) && !phase16AllowsScopeName(manifestPath, manifest.Name) && !phase17AllowsScopeName(manifestPath, manifest.Name) && !phase18AllowsScopeName(manifestPath, manifest.Name) {
			t.Fatalf("%s name must match fixture directory, got %q", manifestPath, manifest.Name)
		}
		if prior := seenNames[manifest.Name]; prior != "" {
			t.Fatalf("duplicate fixture manifest name %q in %s and %s", manifest.Name, prior, manifestPath)
		}
		seenNames[manifest.Name] = manifestPath

		if phase17AllowsScopeName(manifestPath, manifest.Name) || phase18AllowsScopeName(manifestPath, manifest.Name) {
			continue
		}

		if strings.TrimSpace(manifest.Section) == "" {
			t.Fatalf("%s section must not be empty", manifestPath)
		}
		if strings.TrimSpace(manifest.Feature) == "" {
			t.Fatalf("%s feature must not be empty", manifestPath)
		}
		if !phase11AllowedCategories[manifest.Category] {
			t.Fatalf("%s has unsupported category %q", manifestPath, manifest.Category)
		}
		if strings.HasPrefix(manifest.Name, manifest.Category+"-") == false && phase11NameHasCategoryPrefix(manifest.Name) {
			t.Fatalf("%s category %q must match fixture name prefix %q", manifestPath, manifest.Category, manifest.Name)
		}

		if !phase11AllowedComponents[manifest.OwnerComponent] {
			t.Fatalf("%s has unsupported owner component %q", manifestPath, manifest.OwnerComponent)
		}
		phase11RequireSortedUnique(t, manifestPath, "collaborator_components", manifest.CollaboratorComponents)
		for _, component := range manifest.CollaboratorComponents {
			if !phase11AllowedComponents[component] {
				t.Fatalf("%s has unsupported collaborator component %q", manifestPath, component)
			}
			if component == manifest.OwnerComponent {
				t.Fatalf("%s collaborator_components must not repeat owner component %q", manifestPath, component)
			}
		}
		wantComponents := append([]string{manifest.OwnerComponent}, manifest.CollaboratorComponents...)
		sort.Strings(wantComponents)
		if !phase11StringSlicesEqual(manifest.Components, wantComponents) {
			t.Fatalf("%s components must be sorted owner+collaborators\nwant: %v\n got: %v", manifestPath, wantComponents, manifest.Components)
		}

		phase11RequireSortedUnique(t, manifestPath, "components", manifest.Components)
		if len(manifest.References) == 0 {
			t.Fatalf("%s references must not be empty", manifestPath)
		}
		phase11RequireUnique(t, manifestPath, "references", manifest.References)
		for _, reference := range manifest.References {
			phase11RequireMarkdownReference(t, manifestPath, reference)
		}

		if manifest.FakeClock != nil {
			if _, err := time.Parse("2006-01-02T15:04:05Z", *manifest.FakeClock); err != nil {
				t.Fatalf("%s fake_clock must be UTC seconds precision, got %q: %v", manifestPath, *manifest.FakeClock, err)
			}
			if !strings.HasSuffix(*manifest.FakeClock, "Z") {
				t.Fatalf("%s fake_clock must use UTC Z suffix, got %q", manifestPath, *manifest.FakeClock)
			}
		}

		phase11RequireManifestRecords(t, manifestPath, "not_applicable", manifest.NotApplicable)
		phase11RequireSafeRelativePaths(t, manifestPath, "missing_state", manifest.MissingState)
		phase11RequireSafeRelativePaths(t, manifestPath, "input_files", manifest.InputFiles)
		for _, input := range manifest.InputFiles {
			path := filepath.Join(fixtureDir, filepath.FromSlash(input))
			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("%s input file %s is missing: %v", manifestPath, input, err)
			}
			if info.IsDir() {
				t.Fatalf("%s input file %s must not be a directory", manifestPath, input)
			}
			if strings.HasSuffix(input, ".json") {
				phase11RequireValidJSONFile(t, path)
			}
		}

		if len(manifest.Assertions) == 0 {
			t.Fatalf("%s assertions must not be empty", manifestPath)
		}
		phase11RequireUnique(t, manifestPath, "assertions", manifest.Assertions)
		for _, assertion := range manifest.Assertions {
			if !phase11AllowedAssertions[assertion] {
				t.Fatalf("%s has unsupported assertion %q", manifestPath, assertion)
			}
			phase11RequireAssertionEvidence(t, fixtureDir, manifestPath, assertion)
		}
	}
}

func TestPhase11DocumentReferenceGate(t *testing.T) {
	t.Parallel()

	forbidden := []string{
		"docs/details/phase11.md",
		"phase11.md",
		"Phase 11 仕様詳細ファイル",
		"Part 3",
		"Part3",
		"第3部",
		"後続証跡として残る",
		"後続 Phase または横断 owner",
		"後続 Phase の現在状態",
		"formal fixture harness の後続証跡",
		"残証跡",
		"未完了である",
	}

	markdownFiles, err := phase11FindMarkdownFiles()
	if err != nil {
		t.Fatalf("find markdown files: %v", err)
	}
	for _, path := range markdownFiles {
		body := phase11MustReadText(t, path)
		for _, term := range forbidden {
			if strings.Contains(body, term) {
				t.Fatalf("%s contains forbidden stale specification term %q", path, term)
			}
		}
	}
}

func TestPhase11RoadmapStateGate(t *testing.T) {
	t.Parallel()

	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	fixtureSpec := phase11MustReadText(t, "docs/details/fixture.md")
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	phase11CompleteRow := "| Phase 11 | バグ修正ゼロ化。source-code audit、横断 regression、正式 fixture harness、意味のあるテスト、test gap inventory / batch closure、race trigger、mutation selection / mutation zero survivor、contract drift の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) を参照する。 | 実装済み | Phase 10 |"
	phase11IncompleteRow := strings.Replace(phase11CompleteRow, "| 実装済み |", "| 実装中・検証未完了 |", 1)
	phase12CompleteRow := "| Phase 12 | 実装品質ゲート再構築。契約不整合、状態安全性、実行型 fixture harness、mutation / race、owner package 5 ファイル固定、release 再現性の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](DETAIL_INDEX.md#phase-12-quality-gate-entry) を参照する。 | 実装済み | Phase 11 |"
	phase12IncompleteRow := strings.Replace(phase12CompleteRow, "| 実装済み |", "| 実装中・検証未完了 |", 1)
	phase13CompleteRow := "| Phase 13 | 実装整合・品質改善。owner package 5 ファイル責務純度、状態安全性、security / archive / commitstatus 責務集約、queue / finalizer / recovery、MCP 実動作化、外部境界 hardening、fixture / mutation / fault / E2E / CI / release governance の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) を参照する。 | 実装済み | Phase 12 |"

	if !strings.Contains(roadmap, phase11CompleteRow) {
		t.Fatalf("docs/ROADMAP.md must mark Phase 11 as 実装済み after Phase 11 closure records reach final_open_item_count=0")
	}
	if strings.Contains(roadmap, phase11IncompleteRow) {
		t.Fatalf("docs/ROADMAP.md must not keep Phase 11 as 実装中・検証未完了 after closure")
	}
	if !strings.Contains(roadmap, phase12CompleteRow) {
		t.Fatalf("docs/ROADMAP.md must mark Phase 12 as 実装済み after Phase 12 closure records reach final_open_item_count=0")
	}
	if strings.Contains(roadmap, phase12IncompleteRow) {
		t.Fatalf("docs/ROADMAP.md must not keep Phase 12 as 実装中・検証未完了 after closure")
	}
	if !strings.Contains(roadmap, phase13CompleteRow) {
		t.Fatalf("docs/ROADMAP.md must define Phase 13 as 実装済み after Phase 13 closure")
	}
	for _, feature := range phase11QualityGateFeatures {
		want := "| 実装済み | 検証基盤 | " + feature + " |"
		if !strings.Contains(roadmap, want) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 11 feature as 実装済み: %s", feature)
		}
		forbidden := "| 実装中・検証未完了 | 検証基盤 | " + feature + " |"
		if strings.Contains(roadmap, forbidden) {
			t.Fatalf("docs/ROADMAP.md must not keep Phase 11 feature as 実装中・検証未完了 after closure: %s", feature)
		}
	}
	for _, feature := range phase12QualityGateFeatures {
		want := "| 実装済み | 検証基盤 | " + feature + " |"
		if strings.Contains(feature, "owner package") {
			want = "| 実装済み | 実装構造 | " + feature + " |"
		}
		if strings.Contains(feature, "version tag") {
			want = "| 実装済み | 配布・リリース | " + feature + " |"
		}
		if !strings.Contains(roadmap, want) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 12 feature as 実装済み: %s", feature)
		}
		forbidden := strings.Replace(want, "| 実装済み |", "| 実装中・検証未完了 |", 1)
		if strings.Contains(roadmap, forbidden) {
			t.Fatalf("docs/ROADMAP.md must not keep Phase 12 feature as 実装中・検証未完了 after closure: %s", feature)
		}
	}
	for _, feature := range phase13QualityGateFeatures {
		want := "| 実装済み | 検証基盤 | " + feature + " |"
		if strings.Contains(feature, "owner package") {
			want = "| 実装済み | 実装構造 | " + feature + " |"
		}
		if strings.Contains(feature, "statefile") {
			want = "| 実装済み | 状態管理 | " + feature + " |"
		}
		if strings.Contains(feature, "credential") {
			want = "| 実装済み | セキュリティ | " + feature + " |"
		}
		if strings.Contains(feature, "API route") {
			want = "| 実装済み | 管理ツール・API | " + feature + " |"
		}
		if strings.Contains(feature, "queue state machine") {
			want = "| 実装済み | CI ランナー | " + feature + " |"
		}
		if strings.Contains(feature, "MCP resendWebhook") {
			want = "| 実装済み | MCP サーバー | " + feature + " |"
		}
		if strings.Contains(feature, "HTTP timeout") {
			want = "| 実装済み | 外部境界 | " + feature + " |"
		}
		if strings.Contains(feature, "Git tag") {
			want = "| 実装済み | 配布・リリース | " + feature + " |"
		}
		if strings.Contains(feature, "旧 path") {
			want = "| 実装済み | 文書整合 | " + feature + " |"
		}
		if !strings.Contains(roadmap, want) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 13 feature as 実装済み: %s", feature)
		}
		forbidden := strings.Replace(want, "| 実装済み |", "| 実装中・検証未完了 |", 1)
		if strings.Contains(roadmap, forbidden) {
			t.Fatalf("docs/ROADMAP.md must not keep Phase 13 feature as 実装中・検証未完了 after closure: %s", feature)
		}
	}

	for _, fixture := range phase11BuilderFormalFixtures {
		indexRow := "| [`" + fixture + "/`](../" + fixture + "/) | `builder` fixture | 実在 |"
		if !strings.Contains(documentIndex, indexRow) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index builder formal fixture directory as 実在: %s", fixture)
		}
		if _, err := os.Stat(filepath.Join(fixture, "manifest.json")); err != nil {
			t.Fatalf("builder formal fixture %s must include manifest.json: %v", fixture, err)
		}
		if _, err := os.Stat(filepath.Join(fixture, "expected")); err != nil {
			t.Fatalf("builder formal fixture %s must include expected/: %v", fixture, err)
		}
	}

	for _, coverage := range phase11RootCoverageFixtures {
		rootIndexRow := "| [`" + coverage.Root + "`](../" + coverage.Root + ") | `" + coverage.Owner + "` fixture root | 実在 |"
		if !strings.Contains(documentIndex, rootIndexRow) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index Phase 11 root coverage fixture root as 実在: %s", coverage.Root)
		}
		fixtureIndexRow := "| [`" + coverage.FixtureDir + "/`](../" + coverage.FixtureDir + "/) | `" + coverage.Owner + "` Phase 11 root coverage fixture | 実在 |"
		if !strings.Contains(documentIndex, fixtureIndexRow) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index Phase 11 root coverage fixture directory as 実在: %s", coverage.FixtureDir)
		}
		if _, err := os.Stat(filepath.Join(coverage.FixtureDir, "manifest.json")); err != nil {
			t.Fatalf("Phase 11 root coverage fixture %s must include manifest.json: %v", coverage.FixtureDir, err)
		}
		if !strings.Contains(fixtureSpec, "`"+coverage.Root+"` | Phase 11 root coverage formal fixture root。") {
			t.Fatalf("docs/details/fixture.md must record Phase 11 root coverage closure for %s", coverage.Root)
		}
		if strings.Contains(fixtureSpec, "`"+coverage.Root+"` | Phase 11 対象の"+"未作成 root。") {
			t.Fatalf("docs/details/fixture.md must not keep %s as 未作成 root after root coverage closure", coverage.Root)
		}
	}
	if strings.Contains(fixtureSpec, "`p11-"+"open-"+"fixture-root-coverage`") {
		t.Fatalf("docs/details/fixture.md must not keep closed root coverage as an open Phase 11 blocker")
	}
	if strings.Contains(fixtureSpec, "`p11-"+"open-") {
		t.Fatalf("docs/details/fixture.md must not keep any open Phase 11 blocker ids after closure")
	}
	if !strings.Contains(fixtureSpec, "| `final_open_item_count` | `closed` | `0` | `[]` |") {
		t.Fatalf("docs/details/fixture.md must record Phase 11 final_open_item_count=0")
	}
	if !strings.Contains(fixtureSpec, "`survived=0`") {
		t.Fatalf("docs/details/fixture.md must record Phase 11 mutation zero survivor evidence")
	}
	if !strings.Contains(roadmap, "状態ファイル共通永続化契約 | [`docs/DETAIL_INDEX.md`") || !strings.Contains(roadmap, "| 実装済み | 状態管理 | 状態ファイル共通永続化契約 |") {
		t.Fatalf("docs/ROADMAP.md must mark the statefile common persistence contract implemented after verification")
	}
}

func TestPhase16QualityEvidenceClosure(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase16/quality-evidence-closure"
	phase16RequireFixtureFiles(t, root)

	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.Name != phase16Scope {
		t.Fatalf("%s name must be %s, got %q", manifestPath, phase16Scope, manifest.Name)
	}
	if manifest.Scope != phase16Scope {
		t.Fatalf("%s scope must be %s, got %q", manifestPath, phase16Scope, manifest.Scope)
	}
	if !phase11StringSlicesEqual(manifest.Owners, phase16Owners) {
		t.Fatalf("%s owners mismatch\nwant: %v\n got: %v", manifestPath, phase16Owners, manifest.Owners)
	}
	if !phase11StringSlicesEqual(manifest.Components, phase16Components) {
		t.Fatalf("%s components mismatch\nwant: %v\n got: %v", manifestPath, phase16Components, manifest.Components)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase16RequiredChecks) {
		t.Fatalf("%s required_checks mismatch\nwant: %v\n got: %v", manifestPath, phase16RequiredChecks, manifest.RequiredChecks)
	}
	if !phase11StringSlicesEqual(manifest.ClosureRecords, phase16ClosureRecords) {
		t.Fatalf("%s closure_records mismatch\nwant: %v\n got: %v", manifestPath, phase16ClosureRecords, manifest.ClosureRecords)
	}
	for _, sourceAnchor := range manifest.SourceAnchors {
		phase11RequireMarkdownReference(t, manifestPath, sourceAnchor)
	}

	counters := phase16ReadCounters(t, filepath.Join(root, "expected", "counters.json"))
	phase16RequireAllCountersZero(t, counters)

	inventory := phase16ReadInventoryFiles(t, root)
	negativeControls := phase16ReadNegativeControls(t, filepath.Join(root, "input", "negative_controls.json"))
	for _, control := range negativeControls {
		if _, exists := inventory[control.ID]; exists {
			t.Fatalf("negative control inventory id is duplicated: %s", control.ID)
		}
		inventory[control.ID] = phase16InventoryRecord{
			ID:             control.ID,
			WorkUnit:       control.WorkUnit,
			OwnerComponent: control.OwnerComponent,
			SourceRef:      control.PositiveControlRef,
			Classification: "negative_control",
			Status:         "closed",
			CounterKey:     control.CounterKey,
			ClosureRef:     control.PositiveControlRef,
			EvidenceRecords: []string{
				"records/execution.jsonl",
			},
		}
	}

	actions := phase16ReadActions(t, filepath.Join(root, "expected", "actions.json"))
	records := phase16ReadEvidenceRecords(t, root)

	phase16RequireSourceCoverage(t, filepath.Join(root, "input", "source_coverage.json"), inventory)
	phase16RequireInventoryActionRecordClosure(t, inventory, actions, records, counters)
	phase16RequireClosureRecordSet(t, records["records/closure.jsonl"])
	phase16RequireRequiredChecks(t, manifest, records)
	phase16RequireNegativeControls(t, negativeControls, actions, records)
	phase16RequireWorkflow(t, ".github/workflows/phase16-quality-evidence-closure.yml")
	phase16RequireDocumentDriftClosed(t)
}

func TestPhase17ProductionValidationEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase17/production-validation"
	phase17RequireFixtureFiles(t, root)

	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.SchemaVersion != 1 || manifest.Name != phase17Scope || manifest.Scope != phase17Scope || manifest.FixtureRoot != root+"/" {
		t.Fatalf("%s must use schema_version=1, name/scope=%s, fixture_root=%s/", manifestPath, phase17Scope, root)
	}
	if !phase11StringSlicesEqual(manifest.ProviderTargets, phase17ProviderTargets) {
		t.Fatalf("%s provider_targets mismatch\nwant: %v\n got: %v", manifestPath, phase17ProviderTargets, manifest.ProviderTargets)
	}
	if !phase11StringSlicesEqual(manifest.ValidationModes, phase17ValidationModes) {
		t.Fatalf("%s validation_modes mismatch\nwant: %v\n got: %v", manifestPath, phase17ValidationModes, manifest.ValidationModes)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase17RequiredChecks) {
		t.Fatalf("%s required_checks mismatch\nwant: %v\n got: %v", manifestPath, phase17RequiredChecks, manifest.RequiredChecks)
	}

	counters := phase17ReadCounters(t, filepath.Join(root, "expected", "counters.json"))
	phase17RequireCountersClosed(t, counters)
	phase17RequireInputsAndExpected(t, root)
	records := phase17ReadRecords(t, root)
	phase17RequireDerivedCounters(t, counters, records)
	phase17RequireRequiredCheckConnections(t, records)
	phase17RequireBoundaryRecords(t, counters, records)
	phase17RequireClosureEvidenceGates(t, counters, records)
	phase17RequireWorkflow(t, ".github/workflows/phase17-production-validation.yml")
	phase17RequireDocumentState(t)
}

func TestPhase18VPSPullBootstrapEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase18/vps-pull-bootstrap"
	phase18RequireFixtureFiles(t, root)

	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.SchemaVersion != 1 || manifest.Name != phase18Scope || manifest.Scope != phase18Scope || manifest.FixtureRoot != root+"/" {
		t.Fatalf("%s must use schema_version=1, name/scope=%s, fixture_root=%s/", manifestPath, phase18Scope, root)
	}
	if !phase11StringSlicesEqual(manifest.ProviderTargets, phase18ProviderTargets) {
		t.Fatalf("%s provider_targets mismatch\nwant: %v\n got: %v", manifestPath, phase18ProviderTargets, manifest.ProviderTargets)
	}
	if !phase11StringSlicesEqual(manifest.ValidationModes, phase18ValidationModes) {
		t.Fatalf("%s validation_modes mismatch\nwant: %v\n got: %v", manifestPath, phase18ValidationModes, manifest.ValidationModes)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase18RequiredChecks) {
		t.Fatalf("%s required_checks mismatch\nwant: %v\n got: %v", manifestPath, phase18RequiredChecks, manifest.RequiredChecks)
	}

	counters := phase18ReadCounters(t, filepath.Join(root, "expected", "counters.json"))
	phase18RequireCountersClosed(t, counters)
	inputs := phase18RequireInputsAndExpected(t, root)
	phase18RequireBootstrapScriptArtifact(t, phase18BootstrapScriptPath, inputs.Scope)
	records := phase18ReadRecords(t, root)
	phase18RequireRequiredCheckConnections(t, records)
	phase18RequireSessionRecords(t, records["records/session.jsonl"], inputs)
	phase18RequireHealthRecords(t, records["records/health.jsonl"])
	phase18RequireRuntimeRecords(t, records["records/runtime.jsonl"])
	phase18RequireUpdateRollbackRecords(t, records["records/update_rollback.jsonl"])
	phase18RequireBugfixRecords(t, records["records/bugfix.jsonl"], inputs.Bugfix)
	phase18RequireSecurityRecords(t, records["records/security.jsonl"], inputs)
	phase18RequireDocumentDriftRecords(t, records["records/document_drift.jsonl"])
	phase18RequireDerivedCounters(t, counters, records)
	phase18RequireWorkflow(t, ".github/workflows/phase18-vps-pull-bootstrap.yml")
	phase18RequireDocumentState(t)
	t.Run("strict-stdin", TestPhase18BootstrapAcceptsExactArgumentVector)
	t.Run("state-preservation-and-close", TestPhase18BootstrapPreservesOperationalStateAndClosesEvidence)
	t.Run("full-rollback", TestPhase18BootstrapRestoresFullStateAfterCommitFailure)
	t.Run("cleanup-rollback", TestPhase18BootstrapRollsBackWhenFinalCleanupFails)
	t.Run("runtime-json-boundary", TestPhase18BootstrapRejectsMalformedRuntimeJSON)
	t.Run("credential-type-boundary", TestPhase18BootstrapRejectsNonRegularCredentialState)
	t.Run("process-lock", TestPhase18BootstrapProcessLockRejectsLiveOwnerAndRecoversStaleOwner)
}

func TestPhase14ObsidianVaultIntegrationEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase14/obsidian-vault-integration"
	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.Name != "phase-14-obsidian-vault-integration" {
		t.Fatalf("%s name mismatch: %q", manifestPath, manifest.Name)
	}
	if !phase11StringSlicesEqual(manifest.Owners, []string{"obsidian"}) {
		t.Fatalf("%s owners mismatch: %v", manifestPath, manifest.Owners)
	}
	if !phase11StringSlicesEqual(manifest.CollaboratorComponents, []string{"builder", "security"}) {
		t.Fatalf("%s collaborators mismatch: %v", manifestPath, manifest.CollaboratorComponents)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, []string{"phase14-go-format", "phase14-go-test", "phase14-obsidian-vault-fixture", "phase14-builder-handoff", "phase14-document-drift"}) {
		t.Fatalf("%s required checks mismatch: %v", manifestPath, manifest.RequiredChecks)
	}
	for _, rel := range []string{
		"input/options.json",
		"input/vault_tree.json",
		"expected/normalized.json",
		"expected/output_tree.json",
		"expected/counters.json",
		"records/closure.jsonl",
	} {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("Phase 14 fixture path missing %s: %v", rel, err)
		}
	}
	counters := phase13ReadCounterFile(t, filepath.Join(root, "expected/counters.json"))
	for _, key := range []string{
		"phase14_obsidian_vault_boundary_open_count",
		"phase14_obsidian_wikilink_open_count",
		"phase14_obsidian_yaml_rejection_open_count",
		"phase14_obsidian_asset_open_count",
		"phase14_obsidian_builder_handoff_open_count",
		"phase14_fixture_execution_gap_count",
		"final_open_item_count",
	} {
		if counters[key] != 0 {
			t.Fatalf("%s must be 0, got %d", key, counters[key])
		}
	}
	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	if !strings.Contains(roadmap, "| Phase 14 | Obsidian Vault 連携。local vault 読取、wikilink / embed / tag / asset 正規化、YAML frontmatter 拒否、builder handoff の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 14 Obsidian Vault 連携参照](DETAIL_INDEX.md#phase-14-obsidian-vault-integration-entry) を参照する。 | 実装済み | Phase 13 |") {
		t.Fatalf("docs/ROADMAP.md must mark Phase 14 implemented")
	}
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, expected := range []string{
		"| [`components/obsidian/`](../components/obsidian/) | `obsidian` owner package directory | 実在 |",
		"| [`testdata/phase14/obsidian-vault-integration/`](../testdata/phase14/obsidian-vault-integration/) | Phase 14 正式 fixture root | 実在 |",
		"| [`.github/workflows/phase14-obsidian-vault-integration.yml`](../.github/workflows/phase14-obsidian-vault-integration.yml) | Phase 14 required check workflow | 実在 |",
	} {
		if !strings.Contains(documentIndex, expected) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must contain Phase 14実在 row: %s", expected)
		}
	}
}

func TestPhase14ObsidianVaultProductionEntrypoint(t *testing.T) {
	root := t.TempDir()
	vault := filepath.Join(root, "vault")
	phase14WriteFile(t, filepath.Join(vault, "Home.md"), "# Home\n\nSee [[Guide.md|Guide]].\n\n![[images/logo.png|Logo]]\n")
	phase14WriteFile(t, filepath.Join(vault, "Guide.md"), "# Guide\n")
	phase14WriteFile(t, filepath.Join(vault, "images", "logo.png"), "png")
	out := filepath.Join(root, "site")
	report := filepath.Join(out, "obsidian_map.json")
	var stdout, stderr bytes.Buffer
	code := dispatchMain("adlaire-ci-build", []string{
		"--input-mode", "obsidian-vault",
		"--obsidian-vault", vault,
		"--obsidian-entry", "Home.md",
		"--obsidian-report-file", report,
		"--out", out,
		"--title", "Vault",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr must be empty: %s", stderr.String())
	}
	for _, rel := range []string{"index.html", "assets/images/logo.png", "obsidian_map.json"} {
		if _, err := os.Stat(filepath.Join(out, rel)); err != nil {
			t.Fatalf("missing production output %s: %v", rel, err)
		}
	}
	if !strings.Contains(stdout.String(), "obsidian_notes=2") || !strings.Contains(stdout.String(), "obsidian_assets=1") {
		t.Fatalf("production entrypoint missing obsidian report: %s", stdout.String())
	}
}

func TestPhase15ObsidianLocalSyncProductionEntrypoint(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project")
	vault := filepath.Join(root, "vault")
	state := filepath.Join(root, "state")
	phase14WriteFile(t, filepath.Join(project, "docs", "project.md"), "# Project\n")
	if err := os.MkdirAll(vault, 0o755); err != nil {
		t.Fatalf("mkdir vault: %v", err)
	}
	if err := os.MkdirAll(state, 0o755); err != nil {
		t.Fatalf("mkdir state: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := dispatchMain("adlaire-ci-obsidian", []string{
		"sync", "plan",
		"--project-root", project,
		"--vault", vault,
		"--state-dir", state,
		"--plan-file", "plans/sync.json",
		"--direction", "export-only",
		"--delete-policy", "tombstone",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr must be empty: %s", stderr.String())
	}
	var response struct {
		Command    string `json:"command"`
		PlanFile   string `json:"plan_file"`
		PlanHash   string `json:"plan_hash"`
		Operations int    `json:"operations"`
		Conflicts  int    `json:"conflicts"`
		Tombstones int    `json:"tombstones"`
		Applied    bool   `json:"applied"`
	}
	decoder := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("decode production sync response: %v\n%s", err, stdout.String())
	}
	if response.Command != "sync plan" || response.PlanFile != "plans/sync.json" || response.PlanHash == "" || response.Operations != 1 || response.Applied {
		t.Fatalf("unexpected production sync response: %+v", response)
	}
	if _, err := os.Stat(filepath.Join(state, "plans", "sync.json")); err != nil {
		t.Fatalf("production entrypoint must write plan file: %v", err)
	}
}

func phase14WriteFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestPhase13ImplementationAlignmentQualityEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase13/implementation-alignment-quality"
	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)

	if manifest.Name != "phase-13-implementation-alignment-quality" {
		t.Fatalf("%s name must be phase-13-implementation-alignment-quality, got %q", manifestPath, manifest.Name)
	}
	if manifest.Scope != "phase-13-implementation-alignment-quality" {
		t.Fatalf("%s scope must be phase-13-implementation-alignment-quality, got %q", manifestPath, manifest.Scope)
	}
	if !phase11StringSlicesEqual(manifest.Owners, phase13Owners) {
		t.Fatalf("%s owners mismatch\nwant: %v\n got: %v", manifestPath, phase13Owners, manifest.Owners)
	}
	if !phase11StringSlicesEqual(manifest.Entrypoints, phase13Entrypoints) {
		t.Fatalf("%s entrypoints mismatch\nwant: %v\n got: %v", manifestPath, phase13Entrypoints, manifest.Entrypoints)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase13RequiredChecks) {
		t.Fatalf("%s required_checks mismatch\nwant: %v\n got: %v", manifestPath, phase13RequiredChecks, manifest.RequiredChecks)
	}
	if !phase11StringSlicesEqual(manifest.ClosureRecords, phase13ClosureRecords) {
		t.Fatalf("%s closure_records mismatch\nwant: %v\n got: %v", manifestPath, phase13ClosureRecords, manifest.ClosureRecords)
	}
	if !phase11StringSlicesEqual(manifest.NegativeControls, phase13NegativeControls) {
		t.Fatalf("%s negative_controls mismatch\nwant: %v\n got: %v", manifestPath, phase13NegativeControls, manifest.NegativeControls)
	}
	for _, sourceAnchor := range manifest.SourceAnchors {
		phase11RequireMarkdownReference(t, manifestPath, sourceAnchor)
	}

	phase13RequireOwnerPackages(t)
	phase13RequireDocumentIndexRows(t)
	phase13RequireGovernanceFiles(t)
	requireChangeHistoryArtifactBoundary(t)
	phase13RequireZeroExternalDependencies(t)

	effects := phase13ReadCounterFile(t, filepath.Join(root, "expected", "effects.json"))
	counters := phase13ReadCounterFile(t, filepath.Join(root, "expected", "counters.json"))
	if len(counters) != len(phase13ClosureCounters) {
		t.Fatalf("expected/counters.json must include exactly %d counters, got %d", len(phase13ClosureCounters), len(counters))
	}
	for _, key := range phase13ClosureCounters {
		counter, ok := counters[key]
		if !ok {
			t.Fatalf("expected/counters.json must include %s", key)
		}
		if counter != 0 {
			t.Fatalf("expected/counters.json %s must be 0, got %d", key, counter)
		}
		effect, ok := effects[key]
		if !ok {
			t.Fatalf("expected/effects.json must include %s", key)
		}
		if effect != counter {
			t.Fatalf("expected/effects.json %s must equal counters value %d, got %d", key, counter, effect)
		}
	}

	closureRecords := phase13ReadClosureRecords(t, filepath.Join(root, "records", "closure.jsonl"))
	if len(closureRecords) != 18 {
		t.Fatalf("records/closure.jsonl must include 18 closure records, got %d", len(closureRecords))
	}
	for _, record := range closureRecords {
		if record.Scope != "phase-13-implementation-alignment-quality" {
			t.Fatalf("records/closure.jsonl record %s has unexpected scope %q", record.RecordID, record.Scope)
		}
		if record.Status != "closed" {
			t.Fatalf("records/closure.jsonl record %s must be closed, got %q", record.RecordID, record.Status)
		}
		if len(record.OpenItems) != 0 {
			t.Fatalf("records/closure.jsonl record %s must have no open_items, got %v", record.RecordID, record.OpenItems)
		}
	}
	for _, path := range []string{"records/mutation.jsonl", "records/race.jsonl", "records/fault.jsonl", "records/e2e.jsonl", "records/release.jsonl"} {
		if len(phase13ReadJSONLines(t, filepath.Join(root, filepath.FromSlash(path)))) == 0 {
			t.Fatalf("%s must contain at least one JSONL record", path)
		}
	}

	workflow := phase11MustReadText(t, ".github/workflows/phase13-implementation-alignment-quality.yml")
	if !strings.Contains(workflow, "permissions:\n  contents: read") {
		t.Fatalf("Phase 13 workflow must use minimum contents:read permissions")
	}
	for _, requiredCheck := range phase13RequiredChecks {
		if !strings.Contains(workflow, "\n  "+requiredCheck+":") {
			t.Fatalf("Phase 13 workflow must define required check job %s", requiredCheck)
		}
	}
	if strings.Count(workflow, "timeout-minutes:") < len(phase13RequiredChecks) {
		t.Fatalf("Phase 13 workflow must set timeout-minutes on every required check")
	}
	phase13RequireNodeFreeWorkflow(t, workflow)
}

func TestPhase15ObsidianLocalSyncEvidence(t *testing.T) {
	t.Parallel()

	const root = "testdata/phase15/obsidian-local-sync"
	manifestPath := filepath.Join(root, "manifest.json")
	manifest := phase11ReadManifest(t, manifestPath)
	if manifest.Name != "phase-15-obsidian-local-sync" {
		t.Fatalf("%s name mismatch: %q", manifestPath, manifest.Name)
	}
	if !phase11StringSlicesEqual(manifest.Owners, []string{"obsidian"}) {
		t.Fatalf("%s owners mismatch: %v", manifestPath, manifest.Owners)
	}
	if !phase11StringSlicesEqual(manifest.CollaboratorComponents, []string{"release", "security", "setup", "statefile"}) {
		t.Fatalf("%s collaborators mismatch: %v", manifestPath, manifest.CollaboratorComponents)
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, []string{"phase15-go-format", "phase15-go-test", "phase15-obsidian-sync-fixture", "phase15-sync-atomicity", "phase15-release-setup-distribution", "phase15-document-drift"}) {
		t.Fatalf("%s required checks mismatch: %v", manifestPath, manifest.RequiredChecks)
	}
	for _, rel := range []string{
		"input/options.json",
		"input/sync_state.json",
		"input/project_tree.json",
		"input/vault_tree.json",
		"expected/plan.json",
		"expected/apply_state.json",
		"expected/distribution.json",
		"expected/counters.json",
		"records/closure.jsonl",
	} {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("Phase 15 fixture path missing %s: %v", rel, err)
		}
	}
	counters := phase13ReadCounterFile(t, filepath.Join(root, "expected/counters.json"))
	for _, key := range []string{
		"phase15_sync_plan_mismatch_count",
		"phase15_sync_apply_atomicity_open_count",
		"phase15_sync_conflict_open_count",
		"phase15_sync_tombstone_open_count",
		"phase15_sync_rollback_open_count",
		"phase15_sync_service_dependency_open_count",
		"phase15_distribution_open_count",
		"phase15_fixture_execution_gap_count",
		"final_open_item_count",
	} {
		if counters[key] != 0 {
			t.Fatalf("%s must be 0, got %d", key, counters[key])
		}
	}
	distribution := phase11MustReadText(t, filepath.Join(root, "expected/distribution.json"))
	for _, required := range []string{`"binary_asset":"adlaire-ci-obsidian-linux-amd64"`, `"release_asset_count":9`, `"checksum_line_count":8`, `"setup_mode":"install-obsidian"`} {
		if !strings.Contains(distribution, required) {
			t.Fatalf("expected distribution fixture to contain %s", required)
		}
	}
	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	if !strings.Contains(roadmap, "| Phase 15 | Obsidian local vault 同期。plan / apply / rollback、conflict、tombstone、state schema、atomic write、`adlaire-ci-obsidian` 配布連携の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 15 Obsidian local vault 同期参照](DETAIL_INDEX.md#phase-15-obsidian-local-sync-entry) を参照する。 | 実装済み | Phase 14 |") {
		t.Fatalf("docs/ROADMAP.md must mark Phase 15 implemented")
	}
	if !strings.Contains(roadmap, "| 実装済み | Obsidian 同期 | Phase 15 Obsidian local vault 同期 / plan / apply / rollback / conflict / tombstone / atomic write / `adlaire-ci-obsidian` 配布連携 gate |") {
		t.Fatalf("docs/ROADMAP.md must mark the Phase 15 feature implemented")
	}
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, expected := range []string{
		"| [`testdata/phase15/obsidian-local-sync/`](../testdata/phase15/obsidian-local-sync/) | Phase 15 Obsidian local vault 同期 fixture root | 実在 |",
		"| [`.github/workflows/phase15-obsidian-local-sync.yml`](../.github/workflows/phase15-obsidian-local-sync.yml) | Phase 15 required check workflow | 実在 |",
	} {
		if !strings.Contains(documentIndex, expected) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must contain Phase 15 実在 row: %s", expected)
		}
	}
	workflow := phase11MustReadText(t, ".github/workflows/phase15-obsidian-local-sync.yml")
	for _, requiredCheck := range manifest.RequiredChecks {
		if !strings.Contains(workflow, "\n  "+requiredCheck+":") {
			t.Fatalf("Phase 15 workflow must define required check job %s", requiredCheck)
		}
	}
	phase13RequireNodeFreeWorkflow(t, workflow)
}

type phase11Manifest struct {
	SchemaVersion          int                     `json:"schema_version,omitempty"`
	Name                   string                  `json:"name"`
	Scope                  string                  `json:"scope,omitempty"`
	Section                string                  `json:"section"`
	Feature                string                  `json:"feature"`
	Category               string                  `json:"category"`
	OwnerComponent         string                  `json:"owner_component"`
	CollaboratorComponents []string                `json:"collaborator_components"`
	Components             []string                `json:"components"`
	Owners                 []string                `json:"owners,omitempty"`
	Entrypoints            []string                `json:"entrypoints,omitempty"`
	ProviderTargets        []string                `json:"provider_targets,omitempty"`
	ValidationModes        []string                `json:"validation_modes,omitempty"`
	RequiredChecks         []string                `json:"required_checks,omitempty"`
	ClosureRecords         []string                `json:"closure_records,omitempty"`
	NegativeControls       []string                `json:"negative_controls,omitempty"`
	SourceAnchors          []string                `json:"source_anchors,omitempty"`
	FixtureRoot            string                  `json:"fixture_root,omitempty"`
	References             []string                `json:"references"`
	FakeClock              *string                 `json:"fake_clock"`
	NotApplicable          []phase11ManifestRecord `json:"not_applicable"`
	MissingState           []string                `json:"missing_state"`
	InputFiles             []string                `json:"input_files"`
	Assertions             []string                `json:"assertions"`
}

type phase11ManifestRecord struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

var phase11AllowedCategories = map[string]bool{
	"failure":  true,
	"noop":     true,
	"obsidian": true,
	"partial":  true,
	"security": true,
	"success":  true,
}

var phase11AllowedComponents = map[string]bool{
	"admin":        true,
	"api":          true,
	"archive":      true,
	"builder":      true,
	"commitstatus": true,
	"mcp":          true,
	"obsidian":     true,
	"release":      true,
	"runner":       true,
	"sdk":          true,
	"security":     true,
	"setup":        true,
	"statefile":    true,
	"ui":           true,
}

var phase11AllowedAssertions = map[string]bool{
	"effects":     true,
	"idempotency": true,
	"logs":        true,
	"no-write":    true,
	"order":       true,
	"request":     true,
	"response":    true,
	"secret-mask": true,
	"state":       true,
	"stderr":      true,
	"stdout":      true,
}

var phase11QualityGateFeatures = []string{
	"バグ修正ゼロ化 / 仕様全般完了 gate / cross-owner regression gate",
	"全標準実装 artifact source-code audit / artifact coverage zero gap",
	"意味のあるテスト / test gap inventory / batch closure / race trigger / mutation selection / mutation zero survivor gate",
	"仕様全般不備一括棚卸し / 重複ゼロ / inspection pass / spec-gap closure gate",
	"main dispatch / binary version / output manifest / file tree consistency gate",
	"source-code audit inventory / direct I/O elimination / statefile common persistence gate",
	"builder parser / config / output / generated site regression gate",
	"runner source / pipeline / deploy / notification / corrupt state regression gate",
	"API auth / config / backup / webhook / read model / filesystem regression gate",
	"SDK / UI / admin CLI / MCP bridge / route parity regression gate",
	"setup / release distribution / filesystem / external boundary regression gate",
	"正式 fixture directory harness / fixture identity / fixture manifest / contract drift zero gate",
}

var phase12QualityGateFeatures = []string{
	"Phase 12 契約不整合・状態安全性 gate",
	"Phase 12 実行型 fixture harness / dependency injection / production entrypoint / expected 比較 gate",
	"Phase 12 mutation / race / queue state machine / CI required checks gate",
	"Phase 12 owner package 5 ファイル固定 / 巨大コンポーネント分割 gate",
	"Phase 12 version tag / release notes / release reproducibility gate",
}

var phase13QualityGateFeatures = []string{
	"Phase 13 owner package 5 ファイル責務純度 / 6 ファイル目・サブディレクトリ・空ファイル・ダミー実装ゼロ gate",
	"Phase 13 statefile 直接更新排除 / process lock / atomic write / JSON Lines 破損検出・隔離・復旧 gate",
	"Phase 13 credential / token / Go 標準ライブラリ内製 KDF / secret mask / file safety / required log write gate",
	"Phase 13 API route / Admin / SDK / UI client binding / MCP bridge / setup stdout / credential 初期化契約統一 gate",
	"Phase 13 queue state machine / finalizer / active recovery / at-least-once / backup-restore transaction gate",
	"Phase 13 MCP resendWebhook / subscribe / unsubscribe / sampling 実動作化と未実装 error gate",
	"Phase 13 HTTP timeout / body 上限 / graceful shutdown / Webhook SSRF 防止 / SSH strict / systemd 最小権限 gate",
	"Phase 13 fixture execution / mutation / race / fault injection / integration / E2E / CI required check gate",
	"Phase 13 Git tag / GitHub Release / SHA256SUMS / signature / SBOM / reproducible build evidence / recovery procedure gate",
	"Phase 13 旧 path / 重複仕様 / 実装済み表記 drift ゼロ gate",
}

var phase13Owners = []string{
	"admin",
	"api",
	"archive",
	"builder",
	"commitstatus",
	"mcp",
	"release",
	"runner",
	"sdk",
	"security",
	"setup",
	"statefile",
	"ui",
}

var phase13GoOwners = []string{
	"admin",
	"api",
	"archive",
	"builder",
	"commitstatus",
	"mcp",
	"release",
	"runner",
	"security",
	"setup",
	"statefile",
}

var phase13Entrypoints = []string{
	"adlaire-ci-build",
	"adlaire-ci-runner",
	"adlaire-ci-api",
	"adlaire-ci-admin",
	"adlaire-ci-setup",
	"adlaire-ci-release",
	"adlaire-ci-mcp",
	"admin-ui-browser",
	"phase13-fixture-harness",
}

var phase13RequiredChecks = []string{
	"phase13-go-format",
	"phase13-go-test",
	"phase13-go-race",
	"phase13-go-vet",
	"phase13-stdlib-static-analysis",
	"phase13-dependency-inventory",
	"phase13-secret-boundary-scan",
	"phase13-deno-check-sdk",
	"phase13-owner-shape",
	"phase13-contract-parity",
	"phase13-state-safety",
	"phase13-security-boundary",
	"phase13-archive-commitstatus",
	"phase13-runner-recovery",
	"phase13-mcp-real-behavior",
	"phase13-external-boundary",
	"phase13-executable-fixture",
	"phase13-mutation",
	"phase13-concurrency",
	"phase13-fault-injection",
	"phase13-browser",
	"phase13-integration-e2e",
	"phase13-setup-e2e",
	"phase13-release-e2e",
	"phase13-actions-pinning",
	"phase13-recovery-procedure",
	"phase13-document-drift",
}

var phase13ClosureRecords = []string{
	"records/closure.jsonl",
	"records/mutation.jsonl",
	"records/race.jsonl",
	"records/fault.jsonl",
	"records/e2e.jsonl",
	"records/release.jsonl",
}

var phase13NegativeControls = []string{
	"contract-mismatch",
	"direct-state-mutation",
	"github-actions-unpinned",
	"jsonl-corruption",
	"mutation-survivor",
	"race-trigger",
	"fault-injection-failure",
	"token-argv",
}

var phase13ClosureCounters = []string{
	"phase13_owner_file_violation_count",
	"phase13_dummy_or_empty_file_count",
	"phase13_direct_state_mutation_count",
	"phase13_contract_mismatch_count",
	"phase13_state_safety_open_count",
	"phase13_jsonl_corruption_open_count",
	"phase13_security_kdf_open_count",
	"phase13_token_arg_open_count",
	"phase13_archive_commitstatus_ownership_open_count",
	"phase13_queue_recovery_open_count",
	"phase13_mcp_unimplemented_success_count",
	"phase13_external_boundary_open_count",
	"phase13_required_log_write_ignore_count",
	"phase13_fixture_execution_gap_count",
	"phase13_mutation_survived_count",
	"phase13_race_or_concurrency_open_count",
	"phase13_fault_injection_open_count",
	"phase13_e2e_open_count",
	"phase13_ci_required_check_open_count",
	"phase13_action_pin_open_count",
	"phase13_release_evidence_open_count",
	"phase13_recovery_procedure_open_count",
	"phase13_document_drift_open_count",
	"final_open_item_count",
}

const phase16Scope = "phase-16-quality-evidence-closure"

var phase16Owners = []string{
	"admin",
	"api",
	"archive",
	"builder",
	"commitstatus",
	"mcp",
	"obsidian",
	"release",
	"runner",
	"sdk",
	"security",
	"setup",
	"statefile",
	"ui",
}

var phase16Components = []string{
	"admin",
	"api",
	"archive",
	"builder",
	"commitstatus",
	"mcp",
	"obsidian",
	"release",
	"runner",
	"sdk",
	"security",
	"setup",
	"statefile",
	"ui",
}

var phase16RequiredChecks = []string{
	"phase16-go-format",
	"phase16-go-test",
	"phase16-deno-check",
	"phase16-race",
	"phase16-quality-evidence-fixture",
	"phase16-mutation",
	"phase16-fault-injection",
	"phase16-workflow-hardening",
	"phase16-document-drift",
}

var phase16ClosureRecords = []string{
	"records/closure.jsonl",
	"records/execution.jsonl",
	"records/mutation.jsonl",
	"records/fault.jsonl",
	"records/workflow.jsonl",
}

var phase16ClosureCounters = []string{
	"phase16_declarative_evidence_open_count",
	"phase16_execution_evidence_gap_count",
	"phase16_mutation_survived_count",
	"phase16_fault_injection_open_count",
	"phase16_align_unclassified_count",
	"phase16_align_open_count",
	"phase16_ignored_error_unclassified_count",
	"phase16_required_write_failure_open_count",
	"phase16_skip_open_count",
	"phase16_panic_unclassified_count",
	"phase16_determinism_open_count",
	"phase16_http_boundary_open_count",
	"phase16_filesystem_durability_open_count",
	"phase16_workflow_hardening_open_count",
	"phase16_validation_portability_open_count",
	"phase16_release_rehearsal_open_count",
	"phase16_large_owner_risk_open_count",
	"phase16_document_drift_open_count",
	"final_open_item_count",
}

var phase16RequiredFixtureFiles = []string{
	"manifest.json",
	"input/source_coverage.json",
	"input/evidence_inventory.json",
	"input/align_inventory.json",
	"input/error_inventory.json",
	"input/determinism_inventory.json",
	"input/filesystem_inventory.json",
	"input/workflow_inventory.json",
	"input/large_owner_inventory.json",
	"input/negative_controls.json",
	"expected/counters.json",
	"expected/actions.json",
	"records/closure.jsonl",
	"records/execution.jsonl",
	"records/mutation.jsonl",
	"records/fault.jsonl",
	"records/workflow.jsonl",
}

var phase16InventoryFiles = []string{
	"input/evidence_inventory.json",
	"input/align_inventory.json",
	"input/error_inventory.json",
	"input/determinism_inventory.json",
	"input/filesystem_inventory.json",
	"input/workflow_inventory.json",
	"input/large_owner_inventory.json",
}

type phase16CounterFile struct {
	SchemaVersion int            `json:"schema_version"`
	Scope         string         `json:"scope"`
	Counters      map[string]int `json:"counters"`
}

type phase16InventoryFile struct {
	SchemaVersion int                      `json:"schema_version"`
	Scope         string                   `json:"scope"`
	WorkUnit      string                   `json:"work_unit"`
	Records       []phase16InventoryRecord `json:"records"`
}

type phase16InventoryRecord struct {
	ID               string   `json:"id"`
	WorkUnit         string   `json:"work_unit"`
	OwnerComponent   string   `json:"owner_component"`
	SourceRef        string   `json:"source_ref"`
	Classification   string   `json:"classification"`
	Status           string   `json:"status"`
	CounterKey       string   `json:"counter_key"`
	ClosureRef       string   `json:"closure_ref"`
	NotApplicableRef *string  `json:"not_applicable_ref"`
	FutureRef        *string  `json:"future_ref"`
	EvidenceRecords  []string `json:"evidence_records"`
}

type phase16ActionFile struct {
	SchemaVersion int             `json:"schema_version"`
	Scope         string          `json:"scope"`
	Actions       []phase16Action `json:"actions"`
}

type phase16Action struct {
	ID                 string `json:"id"`
	InventoryID        string `json:"inventory_id"`
	Action             string `json:"action"`
	OwnerComponent     string `json:"owner_component"`
	ExpectedCounterKey string `json:"expected_counter_key"`
	ExpectedStatus     string `json:"expected_status"`
	EvidenceRecord     string `json:"evidence_record"`
}

type phase16NegativeControlFile struct {
	SchemaVersion    int                      `json:"schema_version"`
	Scope            string                   `json:"scope"`
	NegativeControls []phase16NegativeControl `json:"negative_controls"`
}

type phase16NegativeControl struct {
	ID                  string `json:"id"`
	WorkUnit            string `json:"work_unit"`
	OwnerComponent      string `json:"owner_component"`
	TargetContract      string `json:"target_contract"`
	ExpectedFailureCode string `json:"expected_failure_code"`
	IsolationMode       string `json:"isolation_mode"`
	RequiredCheck       string `json:"required_check"`
	CounterKey          string `json:"counter_key"`
	PositiveControlRef  string `json:"positive_control_ref"`
}

type phase16SourceCoverageFile struct {
	SchemaVersion     int                   `json:"schema_version"`
	Scope             string                `json:"scope"`
	GeneratedFrom     string                `json:"generated_from"`
	SourceGroupSHA256 string                `json:"source_group_sha256"`
	GoSum             phase16GoSumRecord    `json:"go_sum"`
	Sources           []phase16SourceRecord `json:"sources"`
}

type phase16GoSumRecord struct {
	Path         string  `json:"path"`
	SourceExists bool    `json:"source_exists"`
	SHA256       *string `json:"sha256"`
	SizeBytes    int64   `json:"size_bytes"`
	Reason       string  `json:"reason"`
	InventoryID  string  `json:"inventory_id"`
}

type phase16SourceRecord struct {
	Path           string                `json:"path"`
	SourceExists   bool                  `json:"source_exists"`
	SHA256         string                `json:"sha256"`
	SizeBytes      int64                 `json:"size_bytes"`
	OwnerComponent string                `json:"owner_component"`
	Area           string                `json:"area"`
	DetectedTerms  []phase16DetectedTerm `json:"detected_terms"`
	InventoryRefs  []string              `json:"inventory_refs"`
}

type phase16DetectedTerm struct {
	TermID             string `json:"term_id"`
	Match              string `json:"match"`
	MatchMode          string `json:"match_mode"`
	Count              int    `json:"count"`
	FirstSourceLocator string `json:"first_source_locator"`
}

type phase16EvidenceRecord struct {
	RecordID               string `json:"record_id"`
	Scope                  string `json:"scope"`
	WorkUnit               string `json:"work_unit"`
	InventoryID            string `json:"inventory_id"`
	OwnerComponent         string `json:"owner_component"`
	SourceRef              string `json:"source_ref"`
	Result                 string `json:"result"`
	CounterKey             string `json:"counter_key"`
	CheckName              string `json:"check_name,omitempty"`
	Command                string `json:"command,omitempty"`
	Runtime                string `json:"runtime,omitempty"`
	InputSHA256            string `json:"input_sha256,omitempty"`
	ExpectedSHA256         string `json:"expected_sha256,omitempty"`
	ActualSHA256           string `json:"actual_sha256,omitempty"`
	DiffResult             string `json:"diff_result,omitempty"`
	ExitCode               *int   `json:"exit_code,omitempty"`
	StdoutSHA256           string `json:"stdout_sha256,omitempty"`
	StderrSHA256           string `json:"stderr_sha256,omitempty"`
	MutationTarget         string `json:"mutation_target,omitempty"`
	MutationClass          string `json:"mutation_class,omitempty"`
	BeforeSHA256           string `json:"before_sha256,omitempty"`
	AfterSHA256            string `json:"after_sha256,omitempty"`
	ExpectedFailure        string `json:"expected_failure,omitempty"`
	ActualFailure          string `json:"actual_failure,omitempty"`
	MutationResult         string `json:"mutation_result,omitempty"`
	FaultClass             string `json:"fault_class,omitempty"`
	InjectionResult        string `json:"injection_result,omitempty"`
	RecoveryResult         string `json:"recovery_result,omitempty"`
	SilentSuccess          *bool  `json:"silent_success,omitempty"`
	WorkflowPath           string `json:"workflow_path,omitempty"`
	JobName                string `json:"job_name,omitempty"`
	ActionPinResult        string `json:"action_pin_result,omitempty"`
	PermissionsResult      string `json:"permissions_result,omitempty"`
	TimeoutResult          string `json:"timeout_result,omitempty"`
	RequiredCheckResult    string `json:"required_check_result,omitempty"`
	DockerFallbackResult   string `json:"docker_fallback_result,omitempty"`
	RuntimeVersionResult   string `json:"runtime_version_result,omitempty"`
	DenoCheckResult        string `json:"deno_check_result,omitempty"`
	ReleaseRehearsalResult string `json:"release_rehearsal_result,omitempty"`
}

var phase11BuilderFormalFixtures = []string{
	"testdata/builder/empty-dir",
	"testdata/builder/safe",
	"testdata/builder/single",
	"testdata/builder/site",
	"testdata/builder/strict",
	"testdata/builder/url-safety",
}

type phase11RootCoverageFixture struct {
	Root       string
	Owner      string
	FixtureDir string
}

var phase11RootCoverageFixtures = []phase11RootCoverageFixture{
	{Root: "testdata/runner/", Owner: "runner", FixtureDir: "testdata/runner/success-runner-phase11-root-coverage"},
	{Root: "testdata/api/", Owner: "api", FixtureDir: "testdata/api/success-api-phase11-root-coverage"},
	{Root: "testdata/sdk/", Owner: "sdk", FixtureDir: "testdata/sdk/success-sdk-phase11-root-coverage"},
	{Root: "testdata/ui/", Owner: "ui", FixtureDir: "testdata/ui/success-ui-phase11-root-coverage"},
	{Root: "testdata/statefile/", Owner: "statefile", FixtureDir: "testdata/statefile/success-statefile-phase11-root-coverage"},
	{Root: "testdata/archive/", Owner: "archive", FixtureDir: "testdata/archive/success-archive-phase11-root-coverage"},
	{Root: "testdata/commitstatus/", Owner: "commitstatus", FixtureDir: "testdata/commitstatus/success-commitstatus-phase11-root-coverage"},
	{Root: "testdata/security/", Owner: "security", FixtureDir: "testdata/security/success-security-phase11-root-coverage"},
}

type phase13ClosureRecord struct {
	RecordID  string   `json:"record_id"`
	Scope     string   `json:"scope"`
	Status    string   `json:"status"`
	OpenItems []string `json:"open_items"`
}

func phase13AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase13/implementation-alignment-quality/manifest.json" &&
		name == "phase-13-implementation-alignment-quality"
}

func phase14AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase14/obsidian-vault-integration/manifest.json" &&
		name == "phase-14-obsidian-vault-integration"
}

func phase15AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase15/obsidian-local-sync/manifest.json" &&
		name == "phase-15-obsidian-local-sync"
}

func phase16AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase16/quality-evidence-closure/manifest.json" &&
		name == phase16Scope
}

func phase16RequireFixtureFiles(t *testing.T, root string) {
	t.Helper()

	info, err := os.Stat(root)
	if err != nil {
		t.Fatalf("Phase 16 fixture root must exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("Phase 16 fixture root must be a directory: %s", root)
	}
	for _, rel := range phase16RequiredFixtureFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Phase 16 fixture path missing %s: %v", rel, err)
		}
		if info.IsDir() {
			t.Fatalf("Phase 16 fixture path must be a file: %s", rel)
		}
		if strings.HasSuffix(rel, ".json") || strings.HasSuffix(rel, ".jsonl") {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			if len(data) == 0 || data[len(data)-1] != '\n' {
				t.Fatalf("%s must end with LF", path)
			}
			if strings.HasSuffix(rel, ".json") {
				phase16RequireNoDuplicateJSONKeys(t, path, data)
			}
		}
	}
}

func phase16ReadCounters(t *testing.T, path string) map[string]int {
	t.Helper()

	var file phase16CounterFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase16Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase16Scope)
	}
	return file.Counters
}

func phase16RequireAllCountersZero(t *testing.T, counters map[string]int) {
	t.Helper()

	if len(counters) != len(phase16ClosureCounters) {
		t.Fatalf("Phase 16 counters must include exactly %d keys, got %d", len(phase16ClosureCounters), len(counters))
	}
	for _, key := range phase16ClosureCounters {
		value, ok := counters[key]
		if !ok {
			t.Fatalf("Phase 16 counters missing %s", key)
		}
		if value != 0 {
			t.Fatalf("Phase 16 counter %s must be 0, got %d", key, value)
		}
	}
}

func phase16ReadInventoryFiles(t *testing.T, root string) map[string]phase16InventoryRecord {
	t.Helper()

	records := map[string]phase16InventoryRecord{}
	for _, rel := range phase16InventoryFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		var file phase16InventoryFile
		phase16ReadStrictJSON(t, path, &file)
		if file.SchemaVersion != 1 || file.Scope != phase16Scope || strings.TrimSpace(file.WorkUnit) == "" {
			t.Fatalf("%s must use schema_version=1, scope=%s, and non-empty work_unit", path, phase16Scope)
		}
		if len(file.Records) == 0 {
			t.Fatalf("%s must include at least one inventory record", path)
		}
		for _, record := range file.Records {
			phase16RequireInventoryRecord(t, path, record, file.WorkUnit)
			if _, exists := records[record.ID]; exists {
				t.Fatalf("duplicate Phase 16 inventory id %s", record.ID)
			}
			records[record.ID] = record
		}
	}
	return records
}

func phase16RequireInventoryRecord(t *testing.T, path string, record phase16InventoryRecord, fileWorkUnit string) {
	t.Helper()

	if record.ID == "" || record.WorkUnit == "" || record.OwnerComponent == "" || record.SourceRef == "" || record.Classification == "" || record.Status == "" || record.CounterKey == "" || record.ClosureRef == "" {
		t.Fatalf("%s inventory record must not contain empty required fields: %+v", path, record)
	}
	if record.WorkUnit != fileWorkUnit {
		t.Fatalf("%s record %s work_unit mismatch: file=%s record=%s", path, record.ID, fileWorkUnit, record.WorkUnit)
	}
	if record.Status != "closed" && record.Status != "not_applicable" && record.Status != "future_plan" {
		t.Fatalf("%s record %s has invalid status %q", path, record.ID, record.Status)
	}
	if !phase16CounterAllowed(record.CounterKey) {
		t.Fatalf("%s record %s uses unknown counter %s", path, record.ID, record.CounterKey)
	}
	if len(record.EvidenceRecords) == 0 {
		t.Fatalf("%s record %s must connect to at least one evidence record file", path, record.ID)
	}
	for _, rel := range record.EvidenceRecords {
		if !phase16ClosureRecordAllowed(rel) {
			t.Fatalf("%s record %s uses unknown evidence record %s", path, record.ID, rel)
		}
	}
}

func phase16ReadActions(t *testing.T, path string) []phase16Action {
	t.Helper()

	var file phase16ActionFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase16Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase16Scope)
	}
	if len(file.Actions) == 0 {
		t.Fatalf("%s must include actions", path)
	}
	return file.Actions
}

func phase16ReadNegativeControls(t *testing.T, path string) []phase16NegativeControl {
	t.Helper()

	var file phase16NegativeControlFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase16Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase16Scope)
	}
	if len(file.NegativeControls) < len(phase16RequiredChecks) {
		t.Fatalf("%s must include at least one negative control per required check", path)
	}
	seenTargets := map[string]bool{}
	for _, control := range file.NegativeControls {
		if control.ID == "" || control.WorkUnit == "" || control.OwnerComponent == "" || control.TargetContract == "" || control.ExpectedFailureCode == "" || control.RequiredCheck == "" || control.CounterKey == "" || control.PositiveControlRef == "" {
			t.Fatalf("%s negative control must not contain empty required fields: %+v", path, control)
		}
		if control.IsolationMode != "temporary_copy" {
			t.Fatalf("%s negative control %s must use isolation_mode=temporary_copy", path, control.ID)
		}
		if !strings.HasPrefix(control.ExpectedFailureCode, "PHASE16_") {
			t.Fatalf("%s negative control %s failure code must start with PHASE16_", path, control.ID)
		}
		if !phase16RequiredCheckAllowed(control.RequiredCheck) {
			t.Fatalf("%s negative control %s uses unknown required check %s", path, control.ID, control.RequiredCheck)
		}
		if !phase16CounterAllowed(control.CounterKey) {
			t.Fatalf("%s negative control %s uses unknown counter %s", path, control.ID, control.CounterKey)
		}
		seenTargets[control.TargetContract] = true
	}
	for _, target := range []string{
		"json-duplicate-key",
		"source-coverage",
		"counter-reaggregation",
		"mutation-survivor",
		"fault-injection",
		"workflow-hardening",
		"future-phase-boundary",
	} {
		if !seenTargets[target] {
			t.Fatalf("%s must include negative control target_contract=%s", path, target)
		}
	}
	return file.NegativeControls
}

func phase16ReadEvidenceRecords(t *testing.T, root string) map[string][]phase16EvidenceRecord {
	t.Helper()

	records := map[string][]phase16EvidenceRecord{}
	for _, rel := range phase16ClosureRecords {
		path := filepath.Join(root, filepath.FromSlash(rel))
		records[rel] = phase16ReadEvidenceRecordFile(t, path, rel)
		if len(records[rel]) == 0 {
			t.Fatalf("%s must include at least one record", path)
		}
	}
	return records
}

func phase16ReadEvidenceRecordFile(t *testing.T, path string, rel string) []phase16EvidenceRecord {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	text := string(data)
	if !strings.HasSuffix(text, "\n") {
		t.Fatalf("%s must end with LF", path)
	}
	var records []phase16EvidenceRecord
	seen := map[string]bool{}
	for index, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			t.Fatalf("%s line %d must not be empty", path, index+1)
		}
		phase16RequireNoDuplicateJSONKeys(t, fmt.Sprintf("%s line %d", path, index+1), []byte(line))
		var record phase16EvidenceRecord
		decoder := json.NewDecoder(strings.NewReader(line))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&record); err != nil {
			t.Fatalf("decode %s line %d: %v", path, index+1, err)
		}
		if record.RecordID == "" || record.Scope != phase16Scope || record.WorkUnit == "" || record.InventoryID == "" || record.OwnerComponent == "" || record.SourceRef == "" || record.Result == "" || record.CounterKey == "" {
			t.Fatalf("%s line %d common fields are invalid: %+v", path, index+1, record)
		}
		if seen[record.RecordID] {
			t.Fatalf("%s repeats record_id %s", path, record.RecordID)
		}
		seen[record.RecordID] = true
		if !strings.HasPrefix(record.RecordID, "phase16.record."+strings.TrimSuffix(filepath.Base(rel), ".jsonl")+".") {
			t.Fatalf("%s line %d record_id does not match record file: %s", path, index+1, record.RecordID)
		}
		if record.Result == "failed" || record.Result == "rejected" {
			t.Fatalf("%s line %d must not leave failed or rejected result: %s", path, index+1, record.Result)
		}
		if !phase16CounterAllowed(record.CounterKey) {
			t.Fatalf("%s line %d uses unknown counter %s", path, index+1, record.CounterKey)
		}
		phase16RequireRecordSpecificFields(t, path, index+1, rel, record)
		records = append(records, record)
	}
	return records
}

func phase16RequireRecordSpecificFields(t *testing.T, path string, line int, rel string, record phase16EvidenceRecord) {
	t.Helper()

	switch rel {
	case "records/execution.jsonl":
		if record.CheckName == "" || record.Command == "" || record.Runtime == "" || record.InputSHA256 == "" || record.ExpectedSHA256 == "" || record.ActualSHA256 == "" || record.DiffResult == "" || record.ExitCode == nil || record.StdoutSHA256 == "" || record.StderrSHA256 == "" {
			t.Fatalf("%s line %d execution record missing required execution fields: %+v", path, line, record)
		}
		if strings.HasPrefix(record.InventoryID, "phase16.negative.") {
			if !strings.HasPrefix(record.DiffResult, "PHASE16_") {
				t.Fatalf("%s line %d negative execution diff_result must be PHASE16_* failure code, got %q", path, line, record.DiffResult)
			}
		} else if record.DiffResult != "empty" {
			t.Fatalf("%s line %d execution diff_result must be empty, got %q", path, line, record.DiffResult)
		}
	case "records/mutation.jsonl":
		if record.MutationTarget == "" || record.MutationClass == "" || record.BeforeSHA256 == "" || record.AfterSHA256 == "" || record.ExpectedFailure == "" || record.ActualFailure == "" || record.MutationResult == "" {
			t.Fatalf("%s line %d mutation record missing required fields: %+v", path, line, record)
		}
		if record.MutationResult == "survived" {
			t.Fatalf("%s line %d mutation must not survive", path, line)
		}
	case "records/fault.jsonl":
		if record.FaultClass == "" || record.InjectionResult == "" || record.RecoveryResult == "" || record.SilentSuccess == nil {
			t.Fatalf("%s line %d fault record missing required fields: %+v", path, line, record)
		}
		if *record.SilentSuccess {
			t.Fatalf("%s line %d fault record must not be silent success", path, line)
		}
	case "records/workflow.jsonl":
		if record.WorkflowPath == "" || record.JobName == "" || record.ActionPinResult == "" || record.PermissionsResult == "" || record.TimeoutResult == "" || record.RequiredCheckResult == "" || record.DockerFallbackResult == "" || record.RuntimeVersionResult == "" || record.DenoCheckResult == "" || record.ReleaseRehearsalResult == "" {
			t.Fatalf("%s line %d workflow record missing required fields: %+v", path, line, record)
		}
		for _, result := range []string{record.ActionPinResult, record.PermissionsResult, record.TimeoutResult, record.RequiredCheckResult, record.DockerFallbackResult, record.RuntimeVersionResult, record.DenoCheckResult, record.ReleaseRehearsalResult} {
			if result != "passed" && result != "not_applicable" {
				t.Fatalf("%s line %d workflow result must be passed or not_applicable, got %q", path, line, result)
			}
		}
	}
}

func phase16RequireSourceCoverage(t *testing.T, path string, inventory map[string]phase16InventoryRecord) {
	t.Helper()

	var file phase16SourceCoverageFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase16Scope || file.GeneratedFrom == "" {
		t.Fatalf("%s must use schema_version=1, scope=%s, and generated_from", path, phase16Scope)
	}
	if file.GoSum.Path != "go.sum" {
		t.Fatalf("%s go_sum path must be go.sum", path)
	}
	_, goSumErr := os.Stat("go.sum")
	if os.IsNotExist(goSumErr) {
		if file.GoSum.SourceExists || file.GoSum.SHA256 != nil || file.GoSum.SizeBytes != 0 || file.GoSum.InventoryID == "" {
			t.Fatalf("%s must record absent go.sum as non-existent source", path)
		}
	} else if goSumErr == nil {
		if !file.GoSum.SourceExists || file.GoSum.SHA256 == nil || *file.GoSum.SHA256 != phase16FileSHA256(t, "go.sum") {
			t.Fatalf("%s go_sum hash must match live go.sum", path)
		}
	} else {
		t.Fatalf("stat go.sum: %v", goSumErr)
	}

	if len(file.Sources) < 10 {
		t.Fatalf("%s must include at least 10 source coverage records", path)
	}
	seen := map[string]bool{}
	var digestMaterial strings.Builder
	for _, source := range file.Sources {
		if source.Path == "" || source.OwnerComponent == "" || source.Area == "" || len(source.InventoryRefs) == 0 || len(source.DetectedTerms) == 0 {
			t.Fatalf("%s source record must not contain empty required fields: %+v", path, source)
		}
		if seen[source.Path] {
			t.Fatalf("%s repeats source path %s", path, source.Path)
		}
		seen[source.Path] = true
		info, err := os.Stat(filepath.FromSlash(source.Path))
		if err != nil {
			t.Fatalf("%s source path %s must exist: %v", path, source.Path, err)
		}
		if info.IsDir() {
			t.Fatalf("%s source path %s must be a file", path, source.Path)
		}
		if !source.SourceExists {
			t.Fatalf("%s source path %s must use source_exists=true", path, source.Path)
		}
		if source.SizeBytes != info.Size() {
			t.Fatalf("%s source %s size mismatch: want %d got %d", path, source.Path, info.Size(), source.SizeBytes)
		}
		if got := phase16FileSHA256(t, source.Path); got != source.SHA256 {
			t.Fatalf("%s source %s sha256 mismatch: want %s got %s", path, source.Path, source.SHA256, got)
		}
		for _, ref := range source.InventoryRefs {
			if _, ok := inventory[ref]; !ok {
				t.Fatalf("%s source %s references unknown inventory id %s", path, source.Path, ref)
			}
		}
		for _, term := range source.DetectedTerms {
			if term.TermID == "" || term.Match == "" || term.MatchMode == "" || term.Count < 1 || term.FirstSourceLocator == "" {
				t.Fatalf("%s source %s contains invalid detected term: %+v", path, source.Path, term)
			}
		}
		digestMaterial.WriteString(source.Path)
		digestMaterial.WriteByte('\n')
		digestMaterial.WriteString(source.SHA256)
		digestMaterial.WriteByte('\n')
	}
	if file.SourceGroupSHA256 != phase16StringSHA256(digestMaterial.String()) {
		t.Fatalf("%s source_group_sha256 mismatch", path)
	}
}

func phase16RequireInventoryActionRecordClosure(t *testing.T, inventory map[string]phase16InventoryRecord, actions []phase16Action, records map[string][]phase16EvidenceRecord, counters map[string]int) {
	t.Helper()

	actionByInventory := map[string]phase16Action{}
	for _, action := range actions {
		if action.ID == "" || action.InventoryID == "" || action.Action == "" || action.OwnerComponent == "" || action.ExpectedCounterKey == "" || action.ExpectedStatus == "" || action.EvidenceRecord == "" {
			t.Fatalf("Phase 16 action must not contain empty required fields: %+v", action)
		}
		if _, ok := inventory[action.InventoryID]; !ok {
			t.Fatalf("Phase 16 action %s references unknown inventory id %s", action.ID, action.InventoryID)
		}
		if _, duplicate := actionByInventory[action.InventoryID]; duplicate {
			t.Fatalf("Phase 16 inventory id %s has multiple actions", action.InventoryID)
		}
		if _, ok := counters[action.ExpectedCounterKey]; !ok {
			t.Fatalf("Phase 16 action %s references unknown counter %s", action.ID, action.ExpectedCounterKey)
		}
		if !phase16ClosureRecordAllowed(action.EvidenceRecord) {
			t.Fatalf("Phase 16 action %s references unknown evidence record %s", action.ID, action.EvidenceRecord)
		}
		actionByInventory[action.InventoryID] = action
	}
	if len(actionByInventory) != len(inventory) {
		t.Fatalf("Phase 16 actions must cover every inventory record: actions=%d inventory=%d", len(actionByInventory), len(inventory))
	}

	recordByInventory := map[string][]phase16EvidenceRecord{}
	for rel, recordSet := range records {
		for _, record := range recordSet {
			if _, ok := inventory[record.InventoryID]; !ok {
				t.Fatalf("%s record %s references unknown inventory id %s", rel, record.RecordID, record.InventoryID)
			}
			if _, ok := counters[record.CounterKey]; !ok {
				t.Fatalf("%s record %s references unknown counter %s", rel, record.RecordID, record.CounterKey)
			}
			recordByInventory[record.InventoryID] = append(recordByInventory[record.InventoryID], record)
		}
	}
	for id, record := range inventory {
		action := actionByInventory[id]
		if action.OwnerComponent != record.OwnerComponent || action.ExpectedCounterKey != record.CounterKey || action.ExpectedStatus != record.Status {
			t.Fatalf("Phase 16 action %s does not match inventory record %+v", action.ID, record)
		}
		var connected bool
		for _, evidence := range recordByInventory[id] {
			if phase16RecordFileForID(evidence.RecordID) == action.EvidenceRecord {
				connected = true
				break
			}
		}
		if !connected {
			t.Fatalf("Phase 16 inventory %s is not connected to action evidence %s", id, action.EvidenceRecord)
		}
	}
}

func phase16RequireClosureRecordSet(t *testing.T, records []phase16EvidenceRecord) {
	t.Helper()

	if len(records) != 18 {
		t.Fatalf("Phase 16 closure record set must include 18 records, got %d", len(records))
	}
	for index, record := range records {
		want := fmt.Sprintf("phase16.record.closure.all.closure-%02d", index+1)
		if record.RecordID != want {
			t.Fatalf("Phase 16 closure record order mismatch: want %s got %s", want, record.RecordID)
		}
		if record.Result != "passed" || record.CounterKey == "" {
			t.Fatalf("Phase 16 closure record %s must be passed and counter-connected", record.RecordID)
		}
	}
}

func phase16RequireRequiredChecks(t *testing.T, manifest phase11Manifest, records map[string][]phase16EvidenceRecord) {
	t.Helper()

	required := map[string]bool{}
	for _, check := range phase16RequiredChecks {
		required[check] = false
	}
	if !phase11StringSlicesEqual(manifest.RequiredChecks, phase16RequiredChecks) {
		t.Fatalf("Phase 16 manifest required checks drift")
	}
	for _, record := range records["records/execution.jsonl"] {
		if _, ok := required[record.CheckName]; ok {
			required[record.CheckName] = true
		}
	}
	for _, record := range records["records/workflow.jsonl"] {
		if _, ok := required[record.JobName]; ok {
			required[record.JobName] = true
		}
	}
	for check, seen := range required {
		if !seen {
			t.Fatalf("Phase 16 required check %s is not connected to execution/workflow records", check)
		}
	}
}

func phase16RequireNegativeControls(t *testing.T, controls []phase16NegativeControl, actions []phase16Action, records map[string][]phase16EvidenceRecord) {
	t.Helper()

	actionByInventory := map[string]phase16Action{}
	for _, action := range actions {
		actionByInventory[action.InventoryID] = action
	}
	executionByInventory := map[string][]phase16EvidenceRecord{}
	for _, record := range records["records/execution.jsonl"] {
		executionByInventory[record.InventoryID] = append(executionByInventory[record.InventoryID], record)
	}
	for _, control := range controls {
		action, ok := actionByInventory[control.ID]
		if !ok {
			t.Fatalf("negative control %s must have action", control.ID)
		}
		if action.Action != "execute_and_close" || action.EvidenceRecord != "records/execution.jsonl" || action.ExpectedStatus != "closed" {
			t.Fatalf("negative control %s action must execute_and_close via records/execution.jsonl", control.ID)
		}
		var detected bool
		for _, record := range executionByInventory[control.ID] {
			if record.CheckName == control.RequiredCheck && record.DiffResult == control.ExpectedFailureCode && record.Result == "passed" {
				detected = true
				break
			}
		}
		if !detected {
			t.Fatalf("negative control %s expected failure %s was not detected", control.ID, control.ExpectedFailureCode)
		}
	}
}

func phase16RequireWorkflow(t *testing.T, path string) {
	t.Helper()

	workflow := phase11MustReadText(t, path)
	if !strings.Contains(workflow, "permissions:\n  contents: read") {
		t.Fatalf("%s must use minimum contents:read permissions", path)
	}
	for _, requiredCheck := range phase16RequiredChecks {
		if !strings.Contains(workflow, "\n  "+requiredCheck+":") {
			t.Fatalf("%s must define required check job %s", path, requiredCheck)
		}
	}
	if strings.Count(workflow, "timeout-minutes:") < len(phase16RequiredChecks) {
		t.Fatalf("%s must set timeout-minutes on every required check", path)
	}
	phase13RequireNodeFreeWorkflow(t, workflow)
}

func phase16RequireDocumentDriftClosed(t *testing.T) {
	t.Helper()

	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	if !strings.Contains(roadmap, "| Phase 16 | 実装済み品質証跡実体化・追加検証候補 closure。") || !strings.Contains(roadmap, "| 実装済み | Phase 15 |") {
		t.Fatalf("docs/ROADMAP.md must mark Phase 16 as 実装済み")
	}
	for _, feature := range []string{
		"Phase 16 source coverage set / detection registry / evidence package manifest / inventory / negative control / expected / record schema / checker implementation artifact / checker 実行入口 / checker 再導出 gate",
		"Phase 16 宣言型証跡の実行型証跡化 / production entrypoint 実行 / actual expected 比較 gate",
		"Phase 16 `ALIGN-*` 追加検証候補分類・完了 / 未分類 0 gate",
		"Phase 16 ignored error 分類 / required write failure / panic / skip closure gate",
		"Phase 16 clock / sleep / timeout / entropy determinism gate",
		"Phase 16 filesystem durability parity / atomic write / fsync / parent directory fsync / symlink 非追従 gate",
		"Phase 16 Phase 12 workflow hardening / Docker 検証手順固定 / Deno stable JavaScript 検証 / release rehearsal gate",
		"Phase 16 巨大 owner file risk ledger / 5 ファイル原則維持 / 内部責務区画検査 gate",
	} {
		if !strings.Contains(roadmap, "| 実装済み |") || !strings.Contains(roadmap, feature) {
			t.Fatalf("docs/ROADMAP.md must mark Phase 16 feature implemented: %s", feature)
		}
		if strings.Contains(roadmap, "| 仕様化済み・未実装 | 検証基盤 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 実装品質 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 状態管理 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | CI / 配布 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 保守性 | "+feature+" |") {
			t.Fatalf("docs/ROADMAP.md must not keep Phase 16 feature unimplemented: %s", feature)
		}
	}
	for _, expected := range []string{
		"| [`testdata/phase16/quality-evidence-closure/`](../testdata/phase16/quality-evidence-closure/) | Phase 16 quality evidence closure fixture root | 実在 |",
		"| [`.github/workflows/phase16-quality-evidence-closure.yml`](../.github/workflows/phase16-quality-evidence-closure.yml) | Phase 16 required check workflow | 実在 |",
		"| `testdata/phase16/quality-evidence-closure/manifest.json` | Phase 16 evidence package manifest | 実在 |",
		"| `testdata/phase16/quality-evidence-closure/records/workflow.jsonl` | Phase 16 workflow hardening 証跡 | 実在 |",
	} {
		if !strings.Contains(documentIndex, expected) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must contain Phase 16実在 row: %s", expected)
		}
	}
	for _, forbidden := range []string{
		"| `testdata/phase16/quality-evidence-closure/` | Phase 16 quality evidence closure fixture root | 未作成 |",
		"| `.github/workflows/phase16-quality-evidence-closure.yml` | Phase 16 required check workflow | 未作成 |",
		"| `testdata/phase16/quality-evidence-closure/manifest.json` | Phase 16 evidence package manifest | 未作成 |",
		"| `testdata/phase16/quality-evidence-closure/records/workflow.jsonl` | Phase 16 workflow hardening 証跡 | 未作成 |",
	} {
		if strings.Contains(documentIndex, forbidden) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must not keep Phase 16 path as 未作成: %s", forbidden)
		}
	}
}

func phase16ReadStrictJSON(t *testing.T, path string, out any) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Fatalf("%s must end with LF", path)
	}
	phase16RequireNoDuplicateJSONKeys(t, path, data)
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		t.Fatalf("%s must contain exactly one JSON value", path)
	}
}

func phase16RequireNoDuplicateJSONKeys(t *testing.T, path string, data []byte) {
	t.Helper()

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := phase16CheckJSONValueForDuplicateKeys(decoder); err != nil {
		t.Fatalf("%s duplicate-key check failed: %v", path, err)
	}
	if _, err := decoder.Token(); err != io.EOF {
		t.Fatalf("%s must contain exactly one JSON value", path)
	}
}

func phase16CheckJSONValueForDuplicateKeys(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("object key is not string")
			}
			if seen[key] {
				return fmt.Errorf("duplicate key %q", key)
			}
			seen[key] = true
			if err := phase16CheckJSONValueForDuplicateKeys(decoder); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim('}') {
			return fmt.Errorf("object not closed")
		}
	case '[':
		for decoder.More() {
			if err := phase16CheckJSONValueForDuplicateKeys(decoder); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim(']') {
			return fmt.Errorf("array not closed")
		}
	default:
		return fmt.Errorf("unexpected delimiter %q", delim)
	}
	return nil
}

func phase16CounterAllowed(key string) bool {
	for _, allowed := range phase16ClosureCounters {
		if key == allowed {
			return true
		}
	}
	return false
}

func phase16ClosureRecordAllowed(rel string) bool {
	for _, allowed := range phase16ClosureRecords {
		if rel == allowed {
			return true
		}
	}
	return false
}

func phase16RequiredCheckAllowed(check string) bool {
	for _, allowed := range phase16RequiredChecks {
		if check == allowed {
			return true
		}
	}
	return false
}

func phase16RecordFileForID(recordID string) string {
	parts := strings.Split(recordID, ".")
	if len(parts) < 3 {
		return ""
	}
	return "records/" + parts[2] + ".jsonl"
}

func phase16FileSHA256(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.FromSlash(path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func phase16StringSHA256(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

const phase17Scope = "phase-17-production-validation"

const phase17OperationWindow = "initial-production-validation"

const phase17MonitorIntervalSeconds = 0

const phase17RequiredMonitorSampleCount = 1

const phase17MaxConsecutiveMonitorFailures = 0

var phase17ProviderTargets = []string{"conoha-vps-primary", "xserver-vps-future"}

var phase17ValidationModes = []string{"trial-production-conoha-vps", "production-equivalent-simulation"}

var phase17RequiredChecks = []string{
	"phase17-conoha-vps-preflight",
	"phase17-trial-production-buildout",
	"phase17-conoha-install-update-rollback",
	"phase17-systemd-lifecycle",
	"phase17-runtime-flow",
	"phase17-trial-operation-loop",
	"phase17-production-simulation",
	"phase17-failure-injection",
	"phase17-security-boundary",
	"phase17-bugfix-spec-first",
	"phase17-document-drift",
}

var phase17FixtureFiles = []string{
	"manifest.json",
	"input/provider_targets.json",
	"input/conoha_vps_environment.json",
	"input/buildout_plan.json",
	"input/trial_operation.json",
	"input/simulation_matrix.json",
	"input/security_boundary.json",
	"expected/counters.json",
	"expected/runtime_flow.json",
	"expected/buildout.json",
	"expected/bugfix_loop.json",
	"records/trial_conoha.jsonl",
	"records/buildout.jsonl",
	"records/operation.jsonl",
	"records/bugfix.jsonl",
	"records/simulation.jsonl",
	"records/failure.jsonl",
	"records/security.jsonl",
	"records/document_drift.jsonl",
}

var phase17RecordFiles = []string{
	"records/trial_conoha.jsonl",
	"records/buildout.jsonl",
	"records/operation.jsonl",
	"records/bugfix.jsonl",
	"records/simulation.jsonl",
	"records/failure.jsonl",
	"records/security.jsonl",
	"records/document_drift.jsonl",
}

var phase17Counters = []string{
	"phase17_conoha_trial_operation_open_count",
	"phase17_buildout_open_count",
	"phase17_production_simulation_open_count",
	"phase17_install_update_rollback_open_count",
	"phase17_systemd_lifecycle_open_count",
	"phase17_runtime_flow_open_count",
	"phase17_failure_injection_open_count",
	"phase17_secret_leak_open_count",
	"phase17_destructive_operation_open_count",
	"phase17_xserver_future_misclassified_count",
	"phase17_document_drift_open_count",
	"phase17_trial_operation_open_count",
	"phase17_bugfix_spec_gap_count",
	"phase17_known_bug_open_count",
	"final_open_item_count",
}

var phase17BuildoutWorkUnits = []string{
	"provider-prerequisite",
	"access-baseline",
	"system-baseline",
	"release-acquisition",
	"install-bootstrap",
	"runtime-baseline",
	"operation-baseline",
	"update-rollback-baseline",
	"simulation-baseline",
	"monitoring-baseline",
}

var phase17ForbiddenProviderOperations = []string{
	"vps-create",
	"vps-delete",
	"plan-change",
	"disk-rebuild",
	"volume-create",
	"volume-delete",
	"volume-attach",
	"volume-detach",
	"firewall-lockout",
	"ssh-lockout",
}

var phase17TrialSteps = []string{
	"preflight",
	"install",
	"runtime-flow",
	"update",
	"rollback",
	"reboot-recovery",
	"cleanup",
}

var phase17OperationSteps = []string{
	"operation-start",
	"monitor",
	"issue-detect",
	"spec-general-update",
	"implementation-fix",
	"redeploy",
	"revalidate",
	"evidence-record",
	"known-bug-zero-check",
}

var phase17FailureClasses = []string{
	"disk full",
	"permission denied",
	"short write",
	"fsync failure",
	"rename failure",
	"process kill",
	"network refusal",
	"DNS failure",
	"SSH failure",
	"log write failure",
	"reboot recovery",
	"interrupted recovery",
}

var phase17AbortConditions = []string{
	"secret-boundary-failure",
	"destructive-operation-boundary-failure",
	"state-corruption",
	"service-unrecoverable",
	"known-critical-bug",
}

var phase17SecretClasses = []string{
	"provider-token",
	"admin-token",
	"ssh-private-key",
	"ip-specific-secret",
	"host-specific-secret",
}

var phase17AllowedSecretInputChannels = []string{
	"operator-local-file",
	"operator-stdin",
	"ssh-agent",
	"systemd-credential",
	"not-required",
}

var phase17ForbiddenSecretInputChannels = []string{
	"repository-file",
	"fixture-file",
	"command-argument",
	"environment-variable",
	"pr-body",
	"issue-comment",
	"chat-message",
	"stdout",
	"stderr",
	"log",
	"screen-shot",
}

var phase17RequiredCheckRecordFiles = map[string][]string{
	"phase17-conoha-vps-preflight":           {"records/trial_conoha.jsonl"},
	"phase17-trial-production-buildout":      {"records/buildout.jsonl"},
	"phase17-conoha-install-update-rollback": {"records/trial_conoha.jsonl"},
	"phase17-systemd-lifecycle":              {"records/trial_conoha.jsonl"},
	"phase17-runtime-flow":                   {"records/trial_conoha.jsonl", "records/operation.jsonl"},
	"phase17-trial-operation-loop":           {"records/operation.jsonl"},
	"phase17-production-simulation":          {"records/simulation.jsonl"},
	"phase17-failure-injection":              {"records/failure.jsonl"},
	"phase17-security-boundary":              {"records/security.jsonl"},
	"phase17-bugfix-spec-first":              {"records/bugfix.jsonl"},
	"phase17-document-drift":                 {"records/document_drift.jsonl"},
}

var phase17RequiredCheckCounterKeys = map[string][]string{
	"phase17-conoha-vps-preflight":           {"phase17_conoha_trial_operation_open_count"},
	"phase17-trial-production-buildout":      {"phase17_buildout_open_count"},
	"phase17-conoha-install-update-rollback": {"phase17_install_update_rollback_open_count"},
	"phase17-systemd-lifecycle":              {"phase17_systemd_lifecycle_open_count"},
	"phase17-runtime-flow":                   {"phase17_runtime_flow_open_count"},
	"phase17-trial-operation-loop":           {"phase17_trial_operation_open_count"},
	"phase17-production-simulation":          {"phase17_production_simulation_open_count"},
	"phase17-failure-injection":              {"phase17_failure_injection_open_count"},
	"phase17-security-boundary":              {"phase17_secret_leak_open_count", "phase17_destructive_operation_open_count"},
	"phase17-bugfix-spec-first":              {"phase17_bugfix_spec_gap_count", "phase17_known_bug_open_count"},
	"phase17-document-drift":                 {"phase17_document_drift_open_count", "phase17_xserver_future_misclassified_count"},
}

type phase17CountersFile struct {
	SchemaVersion int            `json:"schema_version"`
	Scope         string         `json:"scope"`
	Counters      map[string]int `json:"counters"`
}

type phase17ProviderTargetsFile struct {
	SchemaVersion  int                     `json:"schema_version"`
	Scope          string                  `json:"scope"`
	ProviderTarget []phase17ProviderTarget `json:"provider_targets"`
}

type phase17ProviderTarget struct {
	ProviderTarget    string  `json:"provider_target"`
	Classification    string  `json:"classification"`
	CompletionBlocker bool    `json:"completion_blocker"`
	OS                *string `json:"os"`
	MinimumPlanClass  *string `json:"minimum_plan_class"`
	MinimumRAMMB      *int    `json:"minimum_ram_mb"`
	ExcludedPlan      *string `json:"excluded_plan"`
	FutureRef         *string `json:"future_ref"`
}

type phase17EnvironmentFile struct {
	SchemaVersion        int                   `json:"schema_version"`
	Scope                string                `json:"scope"`
	ProviderTarget       string                `json:"provider_target"`
	ValidationMode       string                `json:"validation_mode"`
	OS                   string                `json:"os"`
	MinimumPlanClass     string                `json:"minimum_plan_class"`
	MinimumRAMMB         int                   `json:"minimum_ram_mb"`
	ExcludedPlan         string                `json:"excluded_plan"`
	EnvironmentIdentity  phase17EnvironmentID  `json:"environment_identity"`
	Capabilities         []phase17ResultRecord `json:"capabilities"`
	StateDirectoryParent phase17StateParent    `json:"state_directory_parent"`
	Network              []phase17ResultRecord `json:"network"`
	RequiredCommands     []phase17ResultRecord `json:"required_commands"`
}

type phase17EnvironmentID struct {
	ProviderPlanLabel   string `json:"provider_plan_label"`
	ProviderRegionLabel string `json:"provider_region_label"`
	VPSInstanceLabel    string `json:"vps_instance_label"`
	PublicEndpointLabel string `json:"public_endpoint_label"`
	MetadataPolicy      string `json:"metadata_policy"`
}

type phase17ResultRecord struct {
	Name     string `json:"name"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Result   string `json:"result"`
}

type phase17StateParent struct {
	Path          string `json:"path"`
	OwnerExpected string `json:"owner_expected"`
	ModeExpected  string `json:"mode_expected"`
	Result        string `json:"result"`
}

type phase17BuildoutPlanFile struct {
	SchemaVersion              int      `json:"schema_version"`
	Scope                      string   `json:"scope"`
	ProviderTarget             string   `json:"provider_target"`
	ValidationMode             string   `json:"validation_mode"`
	RequiredCheck              string   `json:"required_check"`
	WorkUnits                  []string `json:"work_units"`
	ForbiddenProviderOperation []string `json:"forbidden_provider_operations"`
}

type phase17TrialOperationFile struct {
	SchemaVersion                 int      `json:"schema_version"`
	Scope                         string   `json:"scope"`
	ProviderTarget                string   `json:"provider_target"`
	ValidationMode                string   `json:"validation_mode"`
	OperationSteps                []string `json:"operation_steps"`
	OperationWindow               string   `json:"operation_window"`
	MonitorIntervalSeconds        int      `json:"monitor_interval_seconds"`
	RequiredMonitorSampleCount    int      `json:"required_monitor_sample_count"`
	MaxConsecutiveMonitorFailures int      `json:"max_consecutive_monitor_failures"`
	AbortConditions               []string `json:"abort_conditions"`
	BugfixSpecFirstRequired       bool     `json:"bugfix_spec_first_required"`
	RevalidationRequired          bool     `json:"revalidation_required"`
	KnownBugZeroRequired          bool     `json:"known_bug_zero_required"`
}

type phase17SimulationMatrixFile struct {
	SchemaVersion                int                   `json:"schema_version"`
	Scope                        string                `json:"scope"`
	ProviderTarget               string                `json:"provider_target"`
	ValidationMode               string                `json:"validation_mode"`
	FailureClasses               []phase17FailureClass `json:"failure_classes"`
	DestructiveOperationBoundary []string              `json:"destructive_operation_boundaries"`
	ProviderExecutionAllowed     bool                  `json:"provider_execution_allowed"`
}

type phase17FailureClass struct {
	FailureClass      string `json:"failure_class"`
	InjectionLocation string `json:"injection_location"`
	ExpectedResult    string `json:"expected_result"`
	RecoveryRequired  bool   `json:"recovery_required"`
}

type phase17SecurityBoundaryFile struct {
	SchemaVersion                  int      `json:"schema_version"`
	Scope                          string   `json:"scope"`
	ProviderTarget                 string   `json:"provider_target"`
	ValidationMode                 string   `json:"validation_mode"`
	RequiredCheck                  string   `json:"required_check"`
	SecretClasses                  []string `json:"secret_classes"`
	AllowedSecretInputChannels     []string `json:"allowed_secret_input_channels"`
	ForbiddenSecretInputChannels   []string `json:"forbidden_secret_input_channels"`
	CredentialStoragePolicy        string   `json:"credential_storage_policy"`
	SecretReferencePolicy          string   `json:"secret_reference_policy"`
	DestructiveOperationBoundaries []string `json:"destructive_operation_boundaries"`
	MaskRequired                   bool     `json:"mask_required"`
}

type phase17RuntimeFlowFile struct {
	SchemaVersion      int      `json:"schema_version"`
	Scope              string   `json:"scope"`
	ProviderTarget     string   `json:"provider_target"`
	ValidationMode     string   `json:"validation_mode"`
	ExpectedSteps      []string `json:"expected_steps"`
	CounterConnections []string `json:"counter_connections"`
}

type phase17BuildoutExpectedFile struct {
	SchemaVersion                   int      `json:"schema_version"`
	Scope                           string   `json:"scope"`
	RequiredCheck                   string   `json:"required_check"`
	ExpectedWorkUnitOrder           []string `json:"expected_work_unit_order"`
	AllowedProviderOperationClasses []string `json:"allowed_provider_operation_classes"`
	ForbiddenProviderOperations     []string `json:"forbidden_provider_operations"`
	CounterKey                      string   `json:"counter_key"`
	CounterValue                    int      `json:"counter_value"`
}

type phase17BugfixLoopFile struct {
	SchemaVersion                   int      `json:"schema_version"`
	Scope                           string   `json:"scope"`
	ExpectedOperationOrder          []string `json:"expected_operation_order"`
	OperationWindow                 string   `json:"operation_window"`
	MonitorIntervalSeconds          int      `json:"monitor_interval_seconds"`
	RequiredMonitorSampleCount      int      `json:"required_monitor_sample_count"`
	MaxConsecutiveMonitorFailures   int      `json:"max_consecutive_monitor_failures"`
	AbortConditions                 []string `json:"abort_conditions"`
	SpecUpdateRequired              bool     `json:"spec_update_required"`
	ResponsibilitySourceRefRequired bool     `json:"responsibility_source_ref_required"`
	RevalidationRequired            bool     `json:"revalidation_required"`
	AllowedKnownBugStatuses         []string `json:"allowed_known_bug_statuses"`
	ForbiddenKnownBugStatuses       []string `json:"forbidden_known_bug_statuses"`
	CounterConnections              []string `json:"counter_connections"`
}

type phase17EvidenceRecord struct {
	RecordID                           string  `json:"record_id"`
	Scope                              string  `json:"scope"`
	ProviderTarget                     string  `json:"provider_target"`
	ValidationMode                     string  `json:"validation_mode"`
	CheckName                          string  `json:"check_name"`
	SourceRef                          string  `json:"source_ref"`
	Result                             string  `json:"result"`
	CounterKey                         string  `json:"counter_key"`
	StartedAt                          *string `json:"started_at"`
	EndedAt                            *string `json:"ended_at"`
	StdoutSHA256                       *string `json:"stdout_sha256"`
	StderrSHA256                       *string `json:"stderr_sha256"`
	SecretScanResult                   string  `json:"secret_scan_result"`
	DestructiveOperationResult         string  `json:"destructive_operation_result"`
	EvidenceRef                        string  `json:"evidence_ref"`
	TrialStep                          string  `json:"trial_step,omitempty"`
	EnvironmentRef                     string  `json:"environment_ref,omitempty"`
	EnvironmentIdentityPolicy          string  `json:"environment_identity_policy,omitempty"`
	RuntimeFlowRef                     string  `json:"runtime_flow_ref,omitempty"`
	ServiceState                       string  `json:"service_state,omitempty"`
	StateDigest                        string  `json:"state_digest,omitempty"`
	LogDigest                          string  `json:"log_digest,omitempty"`
	RollbackState                      string  `json:"rollback_state,omitempty"`
	CleanupResult                      string  `json:"cleanup_result,omitempty"`
	BuildoutStep                       string  `json:"buildout_step,omitempty"`
	Sequence                           *int    `json:"sequence,omitempty"`
	WorkUnitResult                     string  `json:"work_unit_result,omitempty"`
	ProviderOperationClass             string  `json:"provider_operation_class,omitempty"`
	ForbiddenProviderOperationDetected *bool   `json:"forbidden_provider_operation_detected,omitempty"`
	PreconditionRef                    string  `json:"precondition_ref,omitempty"`
	ActualRef                          string  `json:"actual_ref,omitempty"`
	ExpectedRef                        string  `json:"expected_ref,omitempty"`
	OperationStep                      string  `json:"operation_step,omitempty"`
	OperationWindow                    string  `json:"operation_window,omitempty"`
	MonitorIntervalSeconds             *int    `json:"monitor_interval_seconds,omitempty"`
	MonitorSampleIndex                 *int    `json:"monitor_sample_index,omitempty"`
	MonitorObservedAt                  *string `json:"monitor_observed_at,omitempty"`
	AbortCondition                     string  `json:"abort_condition,omitempty"`
	AbortResult                        string  `json:"abort_result,omitempty"`
	HealthResult                       string  `json:"health_result,omitempty"`
	StateWriteResult                   string  `json:"state_write_result,omitempty"`
	LogWriteResult                     string  `json:"log_write_result,omitempty"`
	RevalidationRef                    string  `json:"revalidation_ref,omitempty"`
	BugID                              string  `json:"bug_id,omitempty"`
	SpecUpdateRef                      string  `json:"spec_update_ref,omitempty"`
	ResponsibilitySourceRef            string  `json:"responsibility_source_ref,omitempty"`
	ImplementationRef                  string  `json:"implementation_ref,omitempty"`
	KnownBugStatus                     string  `json:"known_bug_status,omitempty"`
	SimulationCase                     string  `json:"simulation_case,omitempty"`
	FailureClass                       string  `json:"failure_class,omitempty"`
	DestructiveOperationToken          string  `json:"destructive_operation_token,omitempty"`
	ProviderExecutionResult            string  `json:"provider_execution_result,omitempty"`
	RecoveryResult                     string  `json:"recovery_result,omitempty"`
	SimulationEnvironmentRef           string  `json:"simulation_environment_ref,omitempty"`
	InjectionMethod                    string  `json:"injection_method,omitempty"`
	ExpectedFailure                    string  `json:"expected_failure,omitempty"`
	ActualFailure                      string  `json:"actual_failure,omitempty"`
	SilentSuccessDetected              *bool   `json:"silent_success_detected,omitempty"`
	SecretClass                        string  `json:"secret_class,omitempty"`
	InputChannel                       string  `json:"input_channel,omitempty"`
	RedactedReference                  string  `json:"redacted_reference,omitempty"`
	SecretReferencePolicy              string  `json:"secret_reference_policy,omitempty"`
	SecretValuePresent                 *bool   `json:"secret_value_present,omitempty"`
	ScanTarget                         string  `json:"scan_target,omitempty"`
	MaskResult                         string  `json:"mask_result,omitempty"`
	CredentialStorageResult            string  `json:"credential_storage_result,omitempty"`
	BoundaryResult                     string  `json:"boundary_result,omitempty"`
	DocumentRef                        string  `json:"document_ref,omitempty"`
	AnchorRef                          string  `json:"anchor_ref,omitempty"`
	FixturePathRef                     string  `json:"fixture_path_ref,omitempty"`
	WorkflowPathRef                    string  `json:"workflow_path_ref,omitempty"`
	ExpectedState                      string  `json:"expected_state,omitempty"`
	ActualState                        string  `json:"actual_state,omitempty"`
	DriftResult                        string  `json:"drift_result,omitempty"`
}

func phase17AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase17/production-validation/manifest.json" &&
		name == phase17Scope
}

func phase18AllowsScopeName(manifestPath string, name string) bool {
	return filepath.ToSlash(manifestPath) == "testdata/phase18/vps-pull-bootstrap/manifest.json" &&
		name == phase18Scope
}

func phase17RequireFixtureFiles(t *testing.T, root string) {
	t.Helper()

	allowed := map[string]bool{}
	for _, rel := range phase17FixtureFiles {
		allowed[rel] = true
		path := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Lstat(path)
		if err != nil {
			t.Fatalf("Phase 17 fixture path missing %s: %v", rel, err)
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			t.Fatalf("Phase 17 fixture path must be a regular file: %s", rel)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if len(data) == 0 || data[len(data)-1] != '\n' {
			t.Fatalf("%s must end with LF", path)
		}
		if strings.HasSuffix(rel, ".json") {
			phase16RequireNoDuplicateJSONKeys(t, path, data)
		}
		phase17RequireNoSecretLikeFixtureContent(t, path, data)
	}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root || entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !allowed[rel] {
			t.Fatalf("Phase 17 fixture contains unknown file %s", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk Phase 17 fixture: %v", err)
	}
}

func phase17ReadCounters(t *testing.T, path string) map[string]int {
	t.Helper()

	var file phase17CountersFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase17Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase17Scope)
	}
	return file.Counters
}

func phase17RequireCountersClosed(t *testing.T, counters map[string]int) {
	t.Helper()

	if len(counters) != len(phase17Counters) {
		t.Fatalf("Phase 17 counters must include exactly %d keys, got %d", len(phase17Counters), len(counters))
	}
	var sum int
	for _, key := range phase17Counters {
		value, ok := counters[key]
		if !ok {
			t.Fatalf("Phase 17 counters missing %s", key)
		}
		if value < 0 {
			t.Fatalf("Phase 17 counter %s must not be negative", key)
		}
		if key != "final_open_item_count" {
			sum += value
		}
	}
	if counters["final_open_item_count"] != sum {
		t.Fatalf("Phase 17 final_open_item_count must equal sum of open counters: want %d got %d", sum, counters["final_open_item_count"])
	}
	if counters["final_open_item_count"] != 0 {
		t.Fatalf("Phase 17 closure requires final_open_item_count=0, got %d", counters["final_open_item_count"])
	}
}

func phase17RequireInputsAndExpected(t *testing.T, root string) {
	t.Helper()

	var targets phase17ProviderTargetsFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/provider_targets.json"), &targets)
	if targets.SchemaVersion != 1 || targets.Scope != phase17Scope || len(targets.ProviderTarget) != 2 {
		t.Fatalf("Phase 17 provider_targets.json schema mismatch")
	}
	for i, target := range targets.ProviderTarget {
		if target.ProviderTarget != phase17ProviderTargets[i] {
			t.Fatalf("Phase 17 provider target order mismatch")
		}
		if target.ProviderTarget == "conoha-vps-primary" {
			if !target.CompletionBlocker || target.OS == nil || *target.OS != "ubuntu-server-24.04-lts-64bit" || target.MinimumPlanClass == nil || *target.MinimumPlanClass != "conoha-vps-1gb-memory-class" || target.MinimumRAMMB == nil || *target.MinimumRAMMB != 1024 {
				t.Fatalf("ConoHa provider target must fix OS and minimum plan class: %+v", target)
			}
		}
		if target.ProviderTarget == "xserver-vps-future" && (target.CompletionBlocker || target.FutureRef == nil || *target.FutureRef == "") {
			t.Fatalf("xserver-vps-future must be future_plan and non-blocking: %+v", target)
		}
	}

	var env phase17EnvironmentFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/conoha_vps_environment.json"), &env)
	if env.SchemaVersion != 1 || env.Scope != phase17Scope || env.ProviderTarget != "conoha-vps-primary" || env.ValidationMode != "trial-production-conoha-vps" || env.OS != "ubuntu-server-24.04-lts-64bit" || env.MinimumPlanClass != "conoha-vps-1gb-memory-class" || env.MinimumRAMMB != 1024 || env.ExcludedPlan != "512mb" {
		t.Fatalf("Phase 17 ConoHa environment fixed values mismatch: %+v", env)
	}
	if env.EnvironmentIdentity.MetadataPolicy != "opaque-non-secret-labels" || phase17ContainsSecretLikeValue(env.EnvironmentIdentity.ProviderPlanLabel, env.EnvironmentIdentity.ProviderRegionLabel, env.EnvironmentIdentity.VPSInstanceLabel, env.EnvironmentIdentity.PublicEndpointLabel) {
		t.Fatalf("Phase 17 environment identity must be opaque and non-secret: %+v", env.EnvironmentIdentity)
	}
	if env.StateDirectoryParent.ModeExpected != "0700" || len(env.Capabilities) == 0 || len(env.Network) == 0 || len(env.RequiredCommands) == 0 {
		t.Fatalf("Phase 17 environment must include state parent, capabilities, network, and commands")
	}

	var buildout phase17BuildoutPlanFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/buildout_plan.json"), &buildout)
	if buildout.SchemaVersion != 1 || buildout.Scope != phase17Scope || buildout.RequiredCheck != "phase17-trial-production-buildout" || !phase11StringSlicesEqual(buildout.WorkUnits, phase17BuildoutWorkUnits) || !phase11StringSlicesEqual(buildout.ForbiddenProviderOperation, phase17ForbiddenProviderOperations) {
		t.Fatalf("Phase 17 buildout plan mismatch: %+v", buildout)
	}

	var trial phase17TrialOperationFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/trial_operation.json"), &trial)
	if trial.SchemaVersion != 1 || trial.Scope != phase17Scope || !phase11StringSlicesEqual(trial.OperationSteps, phase17OperationSteps) || trial.OperationWindow != phase17OperationWindow || trial.MonitorIntervalSeconds != phase17MonitorIntervalSeconds || trial.RequiredMonitorSampleCount != phase17RequiredMonitorSampleCount || trial.MaxConsecutiveMonitorFailures != phase17MaxConsecutiveMonitorFailures || !phase11StringSlicesEqual(trial.AbortConditions, phase17AbortConditions) || !trial.BugfixSpecFirstRequired || !trial.RevalidationRequired || !trial.KnownBugZeroRequired {
		t.Fatalf("Phase 17 trial operation input mismatch: %+v", trial)
	}

	var simulation phase17SimulationMatrixFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/simulation_matrix.json"), &simulation)
	if simulation.SchemaVersion != 1 || simulation.Scope != phase17Scope || simulation.ProviderExecutionAllowed || !phase11StringSlicesEqual(simulation.DestructiveOperationBoundary, phase17ForbiddenProviderOperations) || len(simulation.FailureClasses) != len(phase17FailureClasses) {
		t.Fatalf("Phase 17 simulation matrix mismatch: %+v", simulation)
	}
	for i, failure := range simulation.FailureClasses {
		if failure.FailureClass != phase17FailureClasses[i] || failure.ExpectedResult == "" || !failure.RecoveryRequired {
			t.Fatalf("Phase 17 failure class mismatch at %d: %+v", i, failure)
		}
	}

	var security phase17SecurityBoundaryFile
	phase16ReadStrictJSON(t, filepath.Join(root, "input/security_boundary.json"), &security)
	if security.SchemaVersion != 1 || security.Scope != phase17Scope || security.RequiredCheck != "phase17-security-boundary" || !phase11StringSlicesEqual(security.SecretClasses, phase17SecretClasses) || !phase11StringSlicesEqual(security.AllowedSecretInputChannels, phase17AllowedSecretInputChannels) || !phase11StringSlicesEqual(security.ForbiddenSecretInputChannels, phase17ForbiddenSecretInputChannels) || security.CredentialStoragePolicy != "repository-forbidden" || security.SecretReferencePolicy != "metadata-only" || !phase11StringSlicesEqual(security.DestructiveOperationBoundaries, phase17ForbiddenProviderOperations) || !security.MaskRequired {
		t.Fatalf("Phase 17 security boundary mismatch: %+v", security)
	}

	var runtime phase17RuntimeFlowFile
	phase16ReadStrictJSON(t, filepath.Join(root, "expected/runtime_flow.json"), &runtime)
	if runtime.SchemaVersion != 1 || runtime.Scope != phase17Scope || !phase11StringSlicesEqual(runtime.ExpectedSteps, phase17TrialSteps) {
		t.Fatalf("Phase 17 runtime flow expected steps mismatch: %+v", runtime)
	}

	var expectedBuildout phase17BuildoutExpectedFile
	phase16ReadStrictJSON(t, filepath.Join(root, "expected/buildout.json"), &expectedBuildout)
	if expectedBuildout.SchemaVersion != 1 || expectedBuildout.Scope != phase17Scope || expectedBuildout.RequiredCheck != "phase17-trial-production-buildout" || !phase11StringSlicesEqual(expectedBuildout.ExpectedWorkUnitOrder, phase17BuildoutWorkUnits) || !phase11StringSlicesEqual(expectedBuildout.ForbiddenProviderOperations, phase17ForbiddenProviderOperations) || expectedBuildout.CounterKey != "phase17_buildout_open_count" {
		t.Fatalf("Phase 17 buildout expected mismatch: %+v", expectedBuildout)
	}

	var bugfix phase17BugfixLoopFile
	phase16ReadStrictJSON(t, filepath.Join(root, "expected/bugfix_loop.json"), &bugfix)
	if bugfix.SchemaVersion != 1 || bugfix.Scope != phase17Scope || !phase11StringSlicesEqual(bugfix.ExpectedOperationOrder, phase17OperationSteps) || bugfix.OperationWindow != trial.OperationWindow || bugfix.MonitorIntervalSeconds != trial.MonitorIntervalSeconds || bugfix.RequiredMonitorSampleCount != trial.RequiredMonitorSampleCount || bugfix.MaxConsecutiveMonitorFailures != trial.MaxConsecutiveMonitorFailures || !phase11StringSlicesEqual(bugfix.AbortConditions, trial.AbortConditions) || !bugfix.SpecUpdateRequired || !bugfix.ResponsibilitySourceRefRequired || !bugfix.RevalidationRequired {
		t.Fatalf("Phase 17 bugfix loop expected mismatch: %+v", bugfix)
	}
}

func phase17ReadRecords(t *testing.T, root string) map[string][]phase17EvidenceRecord {
	t.Helper()

	records := map[string][]phase17EvidenceRecord{}
	for _, rel := range phase17RecordFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
		if len(lines) == 0 {
			t.Fatalf("%s must include at least one JSONL record", path)
		}
		seen := map[string]bool{}
		for index, line := range lines {
			if strings.TrimSpace(line) == "" {
				t.Fatalf("%s line %d must not be empty", path, index+1)
			}
			phase16RequireNoDuplicateJSONKeys(t, fmt.Sprintf("%s line %d", path, index+1), []byte(line))
			var record phase17EvidenceRecord
			decoder := json.NewDecoder(strings.NewReader(line))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&record); err != nil {
				t.Fatalf("decode %s line %d: %v", path, index+1, err)
			}
			phase17RequireRecordCommon(t, rel, index+1, record)
			if seen[record.RecordID] {
				t.Fatalf("%s repeats record_id %s", path, record.RecordID)
			}
			seen[record.RecordID] = true
			phase17RequireRecordSpecific(t, rel, index+1, record)
			records[rel] = append(records[rel], record)
		}
	}
	return records
}

func phase17RequireRecordCommon(t *testing.T, rel string, line int, record phase17EvidenceRecord) {
	t.Helper()

	prefix := "phase17.record." + strings.TrimSuffix(filepath.Base(rel), ".jsonl") + "."
	if record.RecordID == "" || !strings.HasPrefix(record.RecordID, prefix) || record.Scope != phase17Scope || record.ProviderTarget == "" || record.ValidationMode == "" || record.CheckName == "" || record.SourceRef == "" || record.Result == "" || record.CounterKey == "" || record.SecretScanResult == "" || record.DestructiveOperationResult == "" || record.EvidenceRef == "" {
		t.Fatalf("%s line %d common fields invalid: %+v", rel, line, record)
	}
	if !phase17RequiredCheckAllowed(record.CheckName) {
		t.Fatalf("%s line %d unknown check_name %s", rel, line, record.CheckName)
	}
	if !phase17CounterAllowed(record.CounterKey) {
		t.Fatalf("%s line %d unknown counter %s", rel, line, record.CounterKey)
	}
	switch record.Result {
	case "passed", "failed", "not_applicable", "future_plan", "rejected":
	default:
		t.Fatalf("%s line %d invalid result %s", rel, line, record.Result)
	}
	if record.Result == "failed" || record.Result == "rejected" {
		t.Fatalf("%s line %d must not leave failed or rejected result", rel, line)
	}
	if record.Result == "passed" && (record.StartedAt == nil || *record.StartedAt == "" || record.EndedAt == nil || *record.EndedAt == "") {
		t.Fatalf("%s line %d passed record must include started_at and ended_at", rel, line)
	}
	for label, value := range map[string]*string{"stdout_sha256": record.StdoutSHA256, "stderr_sha256": record.StderrSHA256} {
		if value != nil && !phase17IsLowerHexSHA256(*value) {
			t.Fatalf("%s line %d %s must be lowercase SHA-256 hex or null: %v", rel, line, label, value)
		}
	}
	switch record.SecretScanResult {
	case "passed", "not_applicable":
	default:
		t.Fatalf("%s line %d invalid secret_scan_result %s", rel, line, record.SecretScanResult)
	}
	switch record.DestructiveOperationResult {
	case "not_applicable", "not_executed", "blocked":
	default:
		t.Fatalf("%s line %d invalid destructive_operation_result %s", rel, line, record.DestructiveOperationResult)
	}
	if phase17ContainsSecretLikeValue(record.EvidenceRef) {
		t.Fatalf("%s line %d evidence_ref must not contain secret-like or endpoint-like value: %s", rel, line, record.EvidenceRef)
	}
}

func phase17RequireRecordSpecific(t *testing.T, rel string, line int, record phase17EvidenceRecord) {
	t.Helper()

	switch rel {
	case "records/trial_conoha.jsonl":
		if record.TrialStep == "" || record.EnvironmentRef == "" || record.EnvironmentIdentityPolicy != "opaque-non-secret-labels" || record.RuntimeFlowRef == "" || record.ServiceState == "" || record.StateDigest == "" || record.LogDigest == "" || record.RollbackState == "" || record.CleanupResult == "" {
			t.Fatalf("%s line %d missing trial ConoHa fields: %+v", rel, line, record)
		}
	case "records/buildout.jsonl":
		if record.BuildoutStep == "" || record.Sequence == nil || record.WorkUnitResult == "" || record.ProviderOperationClass == "" || record.ForbiddenProviderOperationDetected == nil || *record.ForbiddenProviderOperationDetected || record.PreconditionRef == "" || record.ActualRef == "" || record.ExpectedRef == "" {
			t.Fatalf("%s line %d missing buildout fields or detected forbidden operation: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.BuildoutStep, phase17BuildoutWorkUnits) {
			t.Fatalf("%s line %d uses unknown buildout_step %s", rel, line, record.BuildoutStep)
		}
	case "records/operation.jsonl":
		if record.OperationStep == "" || record.OperationWindow != phase17OperationWindow || record.MonitorIntervalSeconds == nil || *record.MonitorIntervalSeconds != phase17MonitorIntervalSeconds || record.AbortCondition == "" || record.AbortResult == "" || record.ServiceState == "" || record.HealthResult == "" || record.StateWriteResult == "" || record.LogWriteResult == "" || record.RevalidationRef == "" {
			t.Fatalf("%s line %d missing operation fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.OperationStep, phase17OperationSteps) {
			t.Fatalf("%s line %d uses unknown operation_step %s", rel, line, record.OperationStep)
		}
		if record.OperationStep == "monitor" {
			if record.MonitorSampleIndex == nil || *record.MonitorSampleIndex != 1 || record.MonitorObservedAt == nil || *record.MonitorObservedAt == "" {
				t.Fatalf("%s line %d monitor baseline sample must include sample index 1 and observed_at: %+v", rel, line, record)
			}
		}
		if record.OperationStep != "monitor" && (record.MonitorSampleIndex != nil || record.MonitorObservedAt != nil) {
			t.Fatalf("%s line %d non-monitor operation record must not claim monitor sample fields: %+v", rel, line, record)
		}
	case "records/bugfix.jsonl":
		if record.BugID == "" || record.SpecUpdateRef == "" || record.ResponsibilitySourceRef == "" || record.ImplementationRef == "" || record.RevalidationRef == "" || record.KnownBugStatus == "" {
			t.Fatalf("%s line %d missing bugfix fields: %+v", rel, line, record)
		}
		switch record.KnownBugStatus {
		case "closed", "future_plan":
		default:
			t.Fatalf("%s line %d invalid known_bug_status %s", rel, line, record.KnownBugStatus)
		}
	case "records/simulation.jsonl":
		if record.SimulationCase == "" || record.FailureClass == "" || record.DestructiveOperationToken == "" || record.ProviderExecutionResult != "not_executed_on_provider" || record.RecoveryResult == "" || record.SimulationEnvironmentRef == "" {
			t.Fatalf("%s line %d missing simulation fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.FailureClass, phase17FailureClasses) {
			t.Fatalf("%s line %d uses unknown simulation failure_class %s", rel, line, record.FailureClass)
		}
		if record.DestructiveOperationToken != "not_applicable" && !phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) {
			t.Fatalf("%s line %d uses unknown destructive_operation_token %s", rel, line, record.DestructiveOperationToken)
		}
	case "records/failure.jsonl":
		if record.FailureClass == "" || record.InjectionMethod == "" || record.ExpectedFailure == "" || record.ActualFailure == "" || record.RecoveryResult == "" || record.SilentSuccessDetected == nil || *record.SilentSuccessDetected {
			t.Fatalf("%s line %d missing failure fields or silent success: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.FailureClass, phase17FailureClasses) {
			t.Fatalf("%s line %d uses unknown failure_class %s", rel, line, record.FailureClass)
		}
	case "records/security.jsonl":
		if record.SecretClass == "" || record.InputChannel == "" || record.RedactedReference == "" || record.SecretReferencePolicy != "metadata-only" || record.SecretValuePresent == nil || *record.SecretValuePresent || record.ScanTarget == "" || record.MaskResult == "" || record.CredentialStorageResult != "not_stored" || record.DestructiveOperationToken == "" || record.BoundaryResult == "" {
			t.Fatalf("%s line %d missing security fields or leaked secret: %+v", rel, line, record)
		}
		phase17RequireSecurityRecordBoundary(t, rel, line, record)
	case "records/document_drift.jsonl":
		if record.DocumentRef == "" || record.AnchorRef == "" || record.FixturePathRef == "" || record.WorkflowPathRef == "" || record.ExpectedState == "" || record.ActualState == "" || record.DriftResult == "" {
			t.Fatalf("%s line %d missing document drift fields: %+v", rel, line, record)
		}
		switch record.DriftResult {
		case "no_drift", "future_plan_confirmed":
		default:
			t.Fatalf("%s line %d invalid drift_result %s", rel, line, record.DriftResult)
		}
	}
}

func phase17RequireRequiredCheckConnections(t *testing.T, records map[string][]phase17EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for rel, recordSet := range records {
		for _, record := range recordSet {
			if !phase17StringInSlice(rel, phase17RequiredCheckRecordFiles[record.CheckName]) {
				t.Fatalf("Phase 17 record %s places check %s in invalid record file %s", record.RecordID, record.CheckName, rel)
			}
			if !phase17StringInSlice(record.CounterKey, phase17RequiredCheckCounterKeys[record.CheckName]) {
				t.Fatalf("Phase 17 record %s connects check %s to invalid counter %s", record.RecordID, record.CheckName, record.CounterKey)
			}
			seen[record.CheckName] = true
		}
	}
	for _, check := range phase17RequiredChecks {
		if !seen[check] {
			t.Fatalf("Phase 17 required check %s is not connected to records", check)
		}
	}
}

func phase17RequireBoundaryRecords(t *testing.T, counters map[string]int, records map[string][]phase17EvidenceRecord) {
	t.Helper()

	phase17RequireBuildoutOrder(t, records["records/buildout.jsonl"], counters["phase17_buildout_open_count"])
	phase17RequireOperationStepCoverage(t, records["records/operation.jsonl"])

	simulationFailureClasses := map[string]bool{}
	for _, record := range records["records/simulation.jsonl"] {
		simulationFailureClasses[record.FailureClass] = true
		if record.ProviderExecutionResult != "not_executed_on_provider" || record.DestructiveOperationResult != "not_executed" {
			t.Fatalf("Phase 17 simulation must not execute destructive operation on provider: %+v", record)
		}
	}
	for _, failureClass := range phase17FailureClasses {
		if !simulationFailureClasses[failureClass] {
			t.Fatalf("Phase 17 failure class missing from simulation records: %s", failureClass)
		}
	}

	failureClasses := map[string]bool{}
	for _, record := range records["records/failure.jsonl"] {
		failureClasses[record.FailureClass] = true
	}
	for _, failureClass := range phase17FailureClasses {
		if !failureClasses[failureClass] {
			t.Fatalf("Phase 17 failure class missing from failure records: %s", failureClass)
		}
	}

	var blockedForbiddenSecretInput bool
	var blockedDestructiveOperation bool
	for _, record := range records["records/security.jsonl"] {
		if phase17StringInSlice(record.InputChannel, phase17ForbiddenSecretInputChannels) && record.BoundaryResult == "blocked" && !*record.SecretValuePresent {
			blockedForbiddenSecretInput = true
		}
		if phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) && record.DestructiveOperationResult == "blocked" && record.BoundaryResult == "blocked" {
			blockedDestructiveOperation = true
		}
	}
	if !blockedForbiddenSecretInput {
		t.Fatalf("Phase 17 security records must include a blocked forbidden secret input channel")
	}
	if !blockedDestructiveOperation {
		t.Fatalf("Phase 17 security records must include a blocked destructive provider operation")
	}

	var futurePlan bool
	for _, record := range records["records/document_drift.jsonl"] {
		if record.ProviderTarget == "xserver-vps-future" && record.Result == "future_plan" && record.DriftResult == "future_plan_confirmed" {
			futurePlan = true
		}
	}
	if !futurePlan {
		t.Fatalf("Phase 17 document drift records must keep xserver-vps-future as future_plan")
	}

	if counters["phase17_xserver_future_misclassified_count"] != 0 {
		t.Fatalf("Phase 17 xserver-vps-future must remain non-blocking with counter 0")
	}
}

func phase17RequireBuildoutOrder(t *testing.T, records []phase17EvidenceRecord, _ int) {
	t.Helper()

	if len(records) != len(phase17BuildoutWorkUnits) {
		t.Fatalf("Phase 17 buildout records must contain exactly %d work units, got %d", len(phase17BuildoutWorkUnits), len(records))
	}
	seen := map[string]bool{}
	for index, record := range records {
		wantStep := phase17BuildoutWorkUnits[index]
		wantSequence := index + 1
		if record.BuildoutStep != wantStep || record.Sequence == nil || *record.Sequence != wantSequence {
			t.Fatalf("Phase 17 buildout record %d order mismatch: want %s/%d got %+v", index+1, wantStep, wantSequence, record)
		}
		if seen[record.BuildoutStep] {
			t.Fatalf("Phase 17 buildout repeats work unit %s", record.BuildoutStep)
		}
		seen[record.BuildoutStep] = true
		if record.Result != "passed" || record.WorkUnitResult != "passed" || record.CounterKey != "phase17_buildout_open_count" {
			t.Fatalf("Phase 17 buildout work unit must be closed before buildout counter can be 0: %+v", record)
		}
	}
}

func phase17RequireOperationStepCoverage(t *testing.T, records []phase17EvidenceRecord) {
	t.Helper()

	positions := map[string]int{}
	for index, record := range records {
		if record.CheckName != "phase17-trial-operation-loop" {
			continue
		}
		if _, exists := positions[record.OperationStep]; !exists {
			positions[record.OperationStep] = index
		}
	}
	previous := -1
	for _, step := range phase17OperationSteps {
		position, exists := positions[step]
		if !exists {
			t.Fatalf("Phase 17 operation records must include placeholder or evidence for step %s", step)
		}
		if position <= previous {
			t.Fatalf("Phase 17 operation records must keep operation step order; step %s is out of order", step)
		}
		previous = position
	}
}

func phase17RequireSecurityRecordBoundary(t *testing.T, rel string, line int, record phase17EvidenceRecord) {
	t.Helper()

	if !phase17StringInSlice(record.SecretClass, phase17SecretClasses) {
		t.Fatalf("%s line %d unknown secret_class %s", rel, line, record.SecretClass)
	}
	allowedInput := phase17StringInSlice(record.InputChannel, phase17AllowedSecretInputChannels)
	forbiddenInput := phase17StringInSlice(record.InputChannel, phase17ForbiddenSecretInputChannels)
	if allowedInput == forbiddenInput {
		t.Fatalf("%s line %d input_channel must be in exactly one allow/deny set: %s", rel, line, record.InputChannel)
	}
	if phase17ContainsSecretLikeValue(record.RedactedReference, record.ScanTarget) {
		t.Fatalf("%s line %d security metadata must remain opaque and non-secret: %+v", rel, line, record)
	}
	switch record.MaskResult {
	case "masked", "not_present", "not_applicable":
	default:
		t.Fatalf("%s line %d invalid mask_result %s", rel, line, record.MaskResult)
	}
	if record.DestructiveOperationToken != "not_applicable" && !phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) {
		t.Fatalf("%s line %d unknown destructive operation token %s", rel, line, record.DestructiveOperationToken)
	}
	switch record.BoundaryResult {
	case "passed", "blocked", "not_applicable":
	default:
		t.Fatalf("%s line %d invalid boundary_result %s", rel, line, record.BoundaryResult)
	}
	if forbiddenInput && (record.Result != "passed" || record.SecretScanResult != "not_applicable" || record.BoundaryResult != "blocked") {
		t.Fatalf("%s line %d forbidden input channel must be blocked as a passing boundary check: %+v", rel, line, record)
	}
	if phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) && (record.DestructiveOperationResult != "blocked" || record.BoundaryResult != "blocked") {
		t.Fatalf("%s line %d forbidden provider operation must be blocked: %+v", rel, line, record)
	}
}

func phase17RequireClosureEvidenceGates(t *testing.T, counters map[string]int, records map[string][]phase17EvidenceRecord) {
	t.Helper()

	if phase17CountMonitorSamples(records["records/operation.jsonl"]) < phase17RequiredMonitorSampleCount {
		t.Fatalf("Phase 17 trial operation closure requires %d baseline monitor sample", phase17RequiredMonitorSampleCount)
	}
	if phase17AnyPendingExternalEvidence(records) {
		t.Fatalf("Phase 17 closure must not leave pending external evidence")
	}
	if phase17RecordWithResult(records["records/operation.jsonl"], "not_applicable") {
		t.Fatalf("Phase 17 operation records must not leave not_applicable placeholders")
	}
	if phase17RecordWithResult(records["records/trial_conoha.jsonl"], "not_applicable") {
		t.Fatalf("Phase 17 trial ConoHa records must not leave not_applicable placeholders")
	}
}

func phase17RequireDerivedCounters(t *testing.T, counters map[string]int, records map[string][]phase17EvidenceRecord) {
	t.Helper()

	derived := phase17DeriveCounters(records)
	for _, key := range phase17Counters {
		if counters[key] != derived[key] {
			t.Fatalf("Phase 17 counter %s must be derived from records: want %d got %d", key, derived[key], counters[key])
		}
	}
}

func phase17DeriveCounters(records map[string][]phase17EvidenceRecord) map[string]int {
	derived := map[string]int{}
	for _, key := range phase17Counters {
		derived[key] = 0
	}

	for _, record := range records["records/trial_conoha.jsonl"] {
		if phase17TrialRecordOpen(record) {
			switch record.CheckName {
			case "phase17-conoha-vps-preflight":
				derived["phase17_conoha_trial_operation_open_count"] = 1
			case "phase17-conoha-install-update-rollback":
				derived["phase17_install_update_rollback_open_count"] = 1
			case "phase17-systemd-lifecycle":
				derived["phase17_systemd_lifecycle_open_count"] = 1
			case "phase17-runtime-flow":
				derived["phase17_runtime_flow_open_count"] = 1
			}
		}
	}
	if !phase17BuildoutClosed(records["records/buildout.jsonl"]) {
		derived["phase17_buildout_open_count"] = 1
	}
	if !phase17SimulationClosed(records["records/simulation.jsonl"]) {
		derived["phase17_production_simulation_open_count"] = 1
	}
	if !phase17FailureInjectionClosed(records["records/failure.jsonl"]) {
		derived["phase17_failure_injection_open_count"] = 1
	}
	if !phase17DocumentDriftClosed(records["records/document_drift.jsonl"]) {
		derived["phase17_document_drift_open_count"] = 1
	}
	if phase17XserverFutureMisclassified(records["records/document_drift.jsonl"]) {
		derived["phase17_xserver_future_misclassified_count"] = 1
	}
	if !phase17TrialOperationClosed(records["records/operation.jsonl"]) {
		derived["phase17_trial_operation_open_count"] = 1
	}
	for _, record := range records["records/security.jsonl"] {
		if phase17SecurityLeakOpen(record) {
			derived["phase17_secret_leak_open_count"] = 1
		}
		if phase17DestructiveOperationOpen(record) {
			derived["phase17_destructive_operation_open_count"] = 1
		}
	}
	for _, record := range records["records/bugfix.jsonl"] {
		if phase17BugfixSpecGapOpen(record) {
			derived["phase17_bugfix_spec_gap_count"] = 1
		}
		if phase17KnownBugOpen(record) {
			derived["phase17_known_bug_open_count"] = 1
		}
	}
	for _, key := range phase17Counters {
		if key != "final_open_item_count" {
			derived["final_open_item_count"] += derived[key]
		}
	}
	return derived
}

func phase17TrialRecordOpen(record phase17EvidenceRecord) bool {
	if record.Result != "passed" {
		return true
	}
	for _, value := range []string{
		record.ServiceState,
		record.StateDigest,
		record.LogDigest,
		record.RollbackState,
		record.CleanupResult,
	} {
		if value == "pending-external-evidence" || value == "not_started" || value == "" {
			return true
		}
	}
	return false
}

func phase17BuildoutClosed(records []phase17EvidenceRecord) bool {
	if len(records) != len(phase17BuildoutWorkUnits) {
		return false
	}
	for _, record := range records {
		if record.Result != "passed" || record.WorkUnitResult != "passed" || record.ForbiddenProviderOperationDetected == nil || *record.ForbiddenProviderOperationDetected {
			return false
		}
	}
	return true
}

func phase17SimulationClosed(records []phase17EvidenceRecord) bool {
	if len(records) != len(phase17FailureClasses) {
		return false
	}
	for _, record := range records {
		if record.Result != "passed" || record.ProviderExecutionResult != "not_executed_on_provider" || record.RecoveryResult == "" {
			return false
		}
	}
	return true
}

func phase17FailureInjectionClosed(records []phase17EvidenceRecord) bool {
	if len(records) != len(phase17FailureClasses) {
		return false
	}
	for _, record := range records {
		if record.Result != "passed" || record.ActualFailure == "" || record.RecoveryResult == "" || record.SilentSuccessDetected == nil || *record.SilentSuccessDetected {
			return false
		}
	}
	return true
}

func phase17DocumentDriftClosed(records []phase17EvidenceRecord) bool {
	if len(records) == 0 {
		return false
	}
	for _, record := range records {
		switch record.DriftResult {
		case "no_drift", "future_plan_confirmed":
		default:
			return false
		}
	}
	return true
}

func phase17XserverFutureMisclassified(records []phase17EvidenceRecord) bool {
	for _, record := range records {
		if record.ProviderTarget == "xserver-vps-future" {
			return record.Result != "future_plan" || record.ValidationMode != "future_plan" || record.DriftResult != "future_plan_confirmed"
		}
	}
	return true
}

func phase17TrialOperationClosed(records []phase17EvidenceRecord) bool {
	if phase17CountMonitorSamples(records) < phase17RequiredMonitorSampleCount {
		return false
	}
	for _, record := range records {
		if record.CheckName != "phase17-trial-operation-loop" {
			continue
		}
		if record.Result != "passed" || phase17OperationRecordPending(record) || record.AbortResult == "triggered" {
			return false
		}
	}
	return true
}

func phase17SecurityLeakOpen(record phase17EvidenceRecord) bool {
	if record.SecretScanResult == "failed" || record.SecretValuePresent == nil || *record.SecretValuePresent || record.CredentialStorageResult != "not_stored" {
		return true
	}
	return record.BoundaryResult == "failed" || record.MaskResult == ""
}

func phase17DestructiveOperationOpen(record phase17EvidenceRecord) bool {
	if record.DestructiveOperationResult == "boundary_failed" || record.BoundaryResult == "failed" {
		return true
	}
	if phase17StringInSlice(record.DestructiveOperationToken, phase17ForbiddenProviderOperations) {
		return record.DestructiveOperationResult != "blocked" || record.BoundaryResult != "blocked"
	}
	return false
}

func phase17BugfixSpecGapOpen(record phase17EvidenceRecord) bool {
	if record.Result == "failed" || record.Result == "rejected" {
		return true
	}
	if record.KnownBugStatus == "future_plan" {
		return false
	}
	return record.SpecUpdateRef == "" || record.ResponsibilitySourceRef == "" || record.ImplementationRef == "" || record.ImplementationRef == "not_started" || record.RevalidationRef == ""
}

func phase17KnownBugOpen(record phase17EvidenceRecord) bool {
	switch record.KnownBugStatus {
	case "closed", "future_plan":
		return false
	default:
		return true
	}
}

func phase17CountMonitorSamples(records []phase17EvidenceRecord) int {
	count := 0
	seen := map[int]bool{}
	for _, record := range records {
		if record.OperationStep != "monitor" || record.MonitorSampleIndex == nil {
			continue
		}
		if !seen[*record.MonitorSampleIndex] {
			seen[*record.MonitorSampleIndex] = true
			count++
		}
	}
	return count
}

func phase17RecordWithResult(records []phase17EvidenceRecord, result string) bool {
	for _, record := range records {
		if record.Result == result {
			return true
		}
	}
	return false
}

func phase17AnyPendingExternalEvidence(records map[string][]phase17EvidenceRecord) bool {
	for _, recordSet := range records {
		for _, record := range recordSet {
			for _, value := range []string{
				record.ServiceState,
				record.StateDigest,
				record.LogDigest,
				record.RollbackState,
				record.CleanupResult,
				record.HealthResult,
				record.StateWriteResult,
				record.LogWriteResult,
				record.ImplementationRef,
			} {
				if value == "pending-external-evidence" || value == "not_started" {
					return true
				}
			}
		}
	}
	return false
}

func phase17OperationRecordPending(record phase17EvidenceRecord) bool {
	for _, value := range []string{
		record.ServiceState,
		record.HealthResult,
		record.StateWriteResult,
		record.LogWriteResult,
	} {
		if value == "pending-external-evidence" {
			return true
		}
	}
	return false
}

func phase17RequireWorkflow(t *testing.T, path string) {
	t.Helper()

	workflow := phase11MustReadText(t, path)
	if !strings.Contains(workflow, "name: phase17-production-validation") || !strings.Contains(workflow, "permissions:\n  contents: read") {
		t.Fatalf("%s must use the Phase 17 workflow name and contents:read permissions", path)
	}
	for _, requiredCheck := range phase17RequiredChecks {
		if !strings.Contains(workflow, "\n  "+requiredCheck+":") {
			t.Fatalf("%s must define required check job %s", path, requiredCheck)
		}
	}
	if strings.Count(workflow, "timeout-minutes:") < len(phase17RequiredChecks) {
		t.Fatalf("%s must set timeout-minutes on every Phase 17 job", path)
	}
	phase13RequireNodeFreeWorkflow(t, workflow)
}

func phase17RequireDocumentState(t *testing.T) {
	t.Helper()

	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	if !strings.Contains(roadmap, "| Phase 17 | ConoHa VPS 試験本番運用。") || !strings.Contains(roadmap, "| 実装済み | Phase 16 |") {
		t.Fatalf("docs/ROADMAP.md must mark Phase 17 as 実装済み after closure evidence reaches zero")
	}
	if strings.Contains(roadmap, "現在の active Phase は Phase 17") {
		t.Fatalf("docs/ROADMAP.md must not keep Phase 17 as active after closure evidence reaches zero")
	}
	expectedPaths := []string{"testdata/phase17/production-validation/"}
	for _, rel := range phase17FixtureFiles {
		expectedPaths = append(expectedPaths, "testdata/phase17/production-validation/"+rel)
	}
	expectedPaths = append(expectedPaths, ".github/workflows/phase17-production-validation.yml")
	for _, rel := range expectedPaths {
		if !strings.Contains(documentIndex, "| `"+rel+"` |") && !strings.Contains(documentIndex, "| [`"+rel+"`]") {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index Phase 17 path %s", rel)
		}
	}
	for _, forbidden := range []string{
		"| `testdata/phase17/production-validation/` | Phase 17 正式 fixture root | 未作成 |",
		"| `.github/workflows/phase17-production-validation.yml` | Phase 17 required check workflow | 未作成 |",
	} {
		if strings.Contains(documentIndex, forbidden) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must not keep Phase 17 path as 未作成: %s", forbidden)
		}
	}
}

func phase17CounterAllowed(key string) bool {
	return phase17StringInSlice(key, phase17Counters)
}

func phase17RequiredCheckAllowed(check string) bool {
	return phase17StringInSlice(check, phase17RequiredChecks)
}

func phase17RequireNoSecretLikeFixtureContent(t *testing.T, path string, data []byte) {
	t.Helper()

	text := string(data)
	for _, forbidden := range []string{
		"-----BEGIN OPENSSH PRIVATE KEY-----",
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN PRIVATE KEY-----",
		"ghp_",
		"github_pat_",
		"xoxb-",
		"sk-proj-",
		"sk-svcacct-",
		"password=",
		"token=",
		"secret=",
		"http://",
		"https://",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("%s must not contain secret-like literal or endpoint-like value %q", path, forbidden)
		}
	}
	for _, token := range strings.FieldsFunc(text, func(r rune) bool {
		return !(r == '.' || r == '-' || r == '_' || r == ':' || r == '/' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z')
	}) {
		if phase17LooksLikeIPv4(token) {
			t.Fatalf("%s must not contain IPv4-like value %q", path, token)
		}
	}
}

func phase17LooksLikeIPv4(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		if part == "" || len(part) > 3 {
			return false
		}
		n := 0
		for _, r := range part {
			if r < '0' || r > '9' {
				return false
			}
			n = n*10 + int(r-'0')
		}
		if n > 255 {
			return false
		}
	}
	return true
}

func phase17IsLowerHexSHA256(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, r := range value {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}

func phase17StringInSlice(value string, values []string) bool {
	for _, candidate := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func phase17ContainsSecretLikeValue(values ...string) bool {
	for _, value := range values {
		lower := strings.ToLower(value)
		for _, forbidden := range []string{"token", "secret", "key", "password", "credential", "http://", "https://", "@"} {
			if strings.Contains(lower, forbidden) {
				return true
			}
		}
	}
	return false
}

const phase18Scope = "phase-18-vps-pull-bootstrap"

var phase18ProviderTargets = []string{"conoha-vps-primary", "xserver-vps-future"}

var phase18ValidationModes = []string{"vps-pull-bootstrap", "future_plan"}

var phase18RequiredChecks = []string{
	"phase18-bootstrap-prerequisite",
	"phase18-bootstrap-session",
	"phase18-bootstrap-artifact-health",
	"phase18-bootstrap-runtime-flow",
	"phase18-pull-runner-foundation",
	"phase18-bootstrap-update-rollback",
	"phase18-bugfix-spec-first",
	"phase18-security-boundary",
	"phase18-document-drift",
}

var phase18FixtureFiles = []string{
	"manifest.json",
	"input/environment_handover.json",
	"input/session_plan.json",
	"input/operation_scope.json",
	"input/security_boundary.json",
	"expected/counters.json",
	"expected/session_flow.json",
	"expected/bugfix_closure.json",
	"expected/negative_boundaries.json",
	"records/session.jsonl",
	"records/health.jsonl",
	"records/runtime.jsonl",
	"records/update_rollback.jsonl",
	"records/bugfix.jsonl",
	"records/security.jsonl",
	"records/document_drift.jsonl",
}

var phase18RecordFiles = []string{
	"records/session.jsonl",
	"records/health.jsonl",
	"records/runtime.jsonl",
	"records/update_rollback.jsonl",
	"records/bugfix.jsonl",
	"records/security.jsonl",
	"records/document_drift.jsonl",
}

var phase18Counters = []string{
	"phase18_bootstrap_prerequisite_open_count",
	"phase18_bootstrap_session_open_count",
	"phase18_bootstrap_runtime_open_count",
	"phase18_pull_runner_foundation_open_count",
	"phase18_bootstrap_update_rollback_open_count",
	"phase18_bugfix_spec_gap_count",
	"phase18_known_bug_open_count",
	"phase18_secret_leak_open_count",
	"phase18_destructive_operation_open_count",
	"phase18_document_drift_open_count",
	"final_open_item_count",
}

var phase18Phase17RequiredChecks = []string{
	"phase17-conoha-vps-preflight",
	"phase17-trial-production-buildout",
	"phase17-conoha-install-update-rollback",
	"phase17-systemd-lifecycle",
	"phase17-runtime-flow",
	"phase17-trial-operation-loop",
	"phase17-production-simulation",
	"phase17-failure-injection",
	"phase17-security-boundary",
	"phase17-bugfix-spec-first",
	"phase17-document-drift",
}

var phase18WorkUnits = []string{
	"phase17-handover",
	"bootstrap-session-open",
	"bootstrap-artifact-health",
	"bootstrap-runtime-flow",
	"pull-runner-foundation",
	"bootstrap-update-rollback-drill",
	"bootstrap-issue-triage",
	"spec-first-fix",
	"rebootstrap-and-revalidate",
	"bootstrap-session-close",
}

var phase18AbortConditions = []string{
	"secret-boundary-failure",
	"destructive-operation-boundary-failure",
	"state-corruption",
	"service-unrecoverable",
	"known-critical-bug",
	"document-drift-blocker",
}

var phase18OperationTargets = []string{
	"bootstrap-artifact-health",
	"bootstrap-runtime-flow",
	"pull-runner-foundation",
	"bootstrap-update-rollback-drill",
	"bootstrap-issue-triage",
	"spec-first-fix",
	"rebootstrap-and-revalidate",
}

var phase18SourceChannels = []string{"integration-head", "stable-release"}

var phase18TopologyRoles = []string{"ci-cd", "site", "single-node"}

var phase18BootstrapInvocationTokens = []string{
	"sh",
	"adlaire-ci-vps-pull-bootstrap.sh",
	"--source-channel",
	"integration-head|stable-release",
	"--topology-role",
	"ci-cd|site|single-node",
	"--source",
	"absolute-https-url-or-absolute-file-path",
	"--sha256",
	"64-lowercase-hex",
	"--install-dir",
	"/opt/adlaire-builder",
	"--bin-dir",
	"/usr/local/bin",
	"--state-dir",
	"/opt/adlaire-builder",
	"--service-user",
	"root",
}

var phase18BootstrapRequiredCommands = []string{"sh", "mktemp", "mkdir", "chmod", "install", "mv", "rm", "cp", "sha256sum", "systemctl", "sleep", "od"}

var phase18BootstrapHTTPClientOrder = []string{"curl"}

var phase18BootstrapFailureCodes = []string{
	"bootstrap-script-http-client-missing",
	"bootstrap-script-source-unreachable",
	"bootstrap-script-digest-mismatch",
	"bootstrap-script-version-mismatch",
	"bootstrap-script-executable-permission-invalid",
	"bootstrap-script-systemd-unit-invalid",
	"bootstrap-script-state-directory-invalid",
	"bootstrap-script-health-failed",
	"bootstrap-script-rollback-failed",
}

var phase18ForbiddenProviderOperations = []string{
	"vps-create",
	"vps-delete",
	"plan-change",
	"disk-rebuild",
	"volume-create",
	"volume-delete",
	"volume-attach",
	"volume-detach",
	"firewall-lockout",
	"ssh-lockout",
}

var phase18SecretClasses = []string{
	"provider-token",
	"admin-token",
	"ssh-private-key",
	"ip-specific-secret",
	"host-specific-secret",
}

var phase18ForbiddenEvidenceValues = []string{
	"provider-account-id",
	"provider-resource-id",
	"ip-address",
	"hostname",
	"fqdn",
	"credential-file-path",
	"secret-value",
	"secret-hash",
	"token-prefix-suffix",
}

var phase18OperationSamples = []string{"open", "runtime", "close"}

var phase18SampleObservations = []string{
	"service_state",
	"api_health",
	"admin_cli_health",
	"runner_state",
	"statefile_digest",
	"audit_log_write",
	"access_log_write",
	"config_log_write",
	"secret_leakage_boundary",
	"known_bug_status",
	"document_drift_status",
}

var phase18HealthTargets = []string{
	"bootstrap-script",
	"release-asset",
	"binary-version",
	"systemd-service",
	"api-health",
	"admin-cli",
	"sdk-client",
	"ui-runtime",
	"state-directory",
	"log-write",
}

var phase18RuntimeTargets = []string{
	"admin-api-sdk-ui",
	"runner-queue",
	"build-execution",
	"statefile-update",
	"audit-access-config-log",
	"commit-status",
	"deploy-simulation",
}

var phase18UpdateRollbackOperations = []string{"update", "rollback", "cleanup"}

var phase18IssueFlow = []string{
	"issue-detect",
	"spec-general-update",
	"implementation-fix",
	"redeploy",
	"revalidate",
	"evidence-record",
	"known-bug-zero-check",
}

var phase18TriageResults = []string{"fix_required", "environment_only", "not_required", "no_issue_detected"}

var phase18FixModes = []string{"spec-first-fix", "environment-only-revalidate", "not-required-revalidate"}

var phase18RequiredBugfixRefs = []string{
	"issue_ref",
	"spec_general_update_ref",
	"responsibility_source_ref",
	"implementation_fix_ref",
	"revalidation_ref",
	"evidence_ref",
}

var phase18DisallowedFixModes = []string{"hotfix-without-spec", "implementation-only-fix", "counter-only-closure"}

var phase18DocumentDriftTargets = []string{
	"state",
	"path",
	"anchor",
	"fixture-root",
	"workflow",
	"required-check",
	"checker-artifact",
	"closure-counter",
}

const phase18BootstrapScriptExpectedState = "adlaire-ci-vps-pull-bootstrap.sh /bin/sh integration-head single-node absolute-https-url-or-absolute-file-path sha256 /opt/adlaire-builder /usr/local/bin /opt/adlaire-builder root no-command-argument-secret staging-verify-commit-rollback fixed-token-flag-interface sh adlaire-ci-vps-pull-bootstrap.sh --source-channel integration-head|stable-release --topology-role ci-cd|site|single-node --source absolute-https-url-or-absolute-file-path --sha256 64-lowercase-hex --install-dir /opt/adlaire-builder --bin-dir /usr/local/bin --state-dir /opt/adlaire-builder --service-user root sh mktemp mkdir chmod install mv rm cp sha256sum systemctl sleep od curl single-json-object-secret-safe diagnostic-json-lines-on-failure common-cli-exit-code-contract bootstrap-script-http-client-missing bootstrap-script-source-unreachable bootstrap-script-digest-mismatch bootstrap-script-version-mismatch bootstrap-script-executable-permission-invalid bootstrap-script-systemd-unit-invalid bootstrap-script-state-directory-invalid bootstrap-script-health-failed bootstrap-script-rollback-failed"

const phase18BootstrapScriptPath = "adlaire-ci-vps-pull-bootstrap.sh"

var phase18BootstrapInstalledBinaries = []string{
	"adlaire-ci-build",
	"adlaire-ci-runner",
	"adlaire-ci-api",
	"adlaire-ci-setup",
	"adlaire-ci-admin",
	"adlaire-ci-mcp",
	"adlaire-ci-obsidian",
}

var phase18ReferenceLabels = map[string]string{
	"manifest.json":                                                              "Phase 18 manifest",
	"input/environment_handover.json":                                            "Phase 18 environment handover",
	"input/session_plan.json":                                                    "Phase 18 session plan",
	"input/operation_scope.json":                                                 "Phase 18 operation scope",
	"input/security_boundary.json":                                               "Phase 18 security boundary",
	"expected/counters.json":                                                     "Phase 18 closure counters",
	"expected/session_flow.json":                                                 "Phase 18 session flow",
	"expected/bugfix_closure.json":                                               "Phase 18 bugfix closure",
	"records/session.jsonl":                                                      "Phase 18 session records",
	"records/health.jsonl":                                                       "Phase 18 health records",
	"records/runtime.jsonl":                                                      "Phase 18 runtime records",
	"records/update_rollback.jsonl":                                              "Phase 18 update rollback records",
	"records/bugfix.jsonl":                                                       "Phase 18 bugfix records",
	"records/security.jsonl":                                                     "Phase 18 security records",
	"records/document_drift.jsonl":                                               "Phase 18 document drift records",
	"testdata/phase18/vps-pull-bootstrap/":                                       "Phase 18 fixture root",
	"adlaire-ci-vps-pull-bootstrap.sh":                                           "Phase 18 VPS Pull Bootstrap shell artifact",
	".github/workflows/phase18-vps-pull-bootstrap.yml":                           "Phase 18 required check workflow",
	"main_test.go":                                                               "Phase 18 checker implementation artifact",
	"docs/ROADMAP.md#roadmap-initial-phase-plan":                                 "状態・計画責務 §4.1",
	"docs/SPEC.md#policy-spec-maturity":                                          "ポリシー責務 §0a",
	"docs/DETAIL_INDEX.md#phase-18-vps-pull-bootstrap-entry":                     "詳細仕様入口責務 Phase 18 VPS Pull Bootstrap 参照",
	"docs/DETAIL_INDEX.md#phase-18-state-transition-gate":                        "詳細仕様入口責務 Phase 18 状態遷移 gate",
	"docs/DOCUMENT_INDEX.md#phase-18-target-paths":                               "文書・実装ファイル所在の索引責務 Phase 18 target path 所在",
	"docs/details/fixture.md#phase-18-vps-pull-bootstrap-evidence":               "fixture 証跡責務 Phase 18 VPS Pull Bootstrap",
	"docs/details/production-validation.md#phase-18-vps-pull-bootstrap-contract": "本番検証詳細本文責務 Phase 18 VPS Pull Bootstrap 契約",
	"docs/details/fixture.md#phase-17-production-validation-evidence":            "fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡",
	"docs/details/fixture.md#phase-18-negative-boundary-contract":                "fixture 証跡責務 Phase 18 negative boundary 固定契約",
	"docs/details/release.md#release-responsibility-boundary":                    "Release 詳細本文責務 責務境界",
	"docs/details/setup.md#0-責務境界":                                               "Setup 詳細本文責務 責務境界",
	"docs/details/api.md#0-責務境界":                                                 "API 詳細本文責務 責務境界",
	"docs/details/admin.md#0-責務境界":                                               "Admin 詳細本文責務 責務境界",
	"docs/details/sdk.md#0-責務境界":                                                 "SDK 詳細本文責務 責務境界",
	"docs/details/ui.md#0-責務境界":                                                  "UI 詳細本文責務 責務境界",
	"docs/details/statefile.md#0-責務境界":                                           "Statefile 詳細本文責務 責務境界",
	"docs/details/runner.md#0-責務境界":                                              "Runner 詳細本文責務 責務境界",
	"docs/details/builder.md#0-責務境界":                                             "Builder 詳細本文責務 責務境界",
	"docs/details/commitstatus.md#0-責務境界":                                        "Commit Status 詳細本文責務 責務境界",
}

var phase18NegativeBoundarySourceRef = phase18MarkdownRef("docs/details/fixture.md#phase-18-negative-boundary-contract")

type phase18HealthTargetExpectation struct {
	CheckName            string
	SourceRef            string
	OwnerRef             string
	EndpointOrCommandRef string
	State                string
	LogWriteResult       string
	LatencyRequired      bool
}

var phase18HealthTargetExpectations = map[string]phase18HealthTargetExpectation{
	"bootstrap-script": {"phase18-bootstrap-artifact-health", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/production-validation.md#phase-18-vps-pull-bootstrap-contract"), phase18MarkdownRef("input/operation_scope.json"), phase18BootstrapScriptExpectedState, "not_applicable", false},
	"release-asset":    {"phase18-bootstrap-artifact-health", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/release.md#release-responsibility-boundary"), phase18MarkdownRef("docs/details/release.md#release-responsibility-boundary"), "release-artifact-evidence", "not_applicable", false},
	"binary-version":   {"phase18-bootstrap-artifact-health", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/release.md#release-responsibility-boundary"), phase18MarkdownRef("docs/details/release.md#release-responsibility-boundary"), "binary-version-evidence", "not_applicable", false},
	"systemd-service":  {"phase18-bootstrap-artifact-health", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/setup.md#0-責務境界"), phase18MarkdownRef("docs/details/setup.md#0-責務境界"), "service-active", "not_applicable", false},
	"api-health":       {"phase18-bootstrap-artifact-health", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/api.md#0-責務境界"), phase18MarkdownRef("docs/details/api.md#0-責務境界"), "api-health-passed", "not_applicable", true},
	"admin-cli":        {"phase18-bootstrap-runtime-flow", phase18MarkdownRef("expected/session_flow.json"), phase18MarkdownRef("docs/details/admin.md#0-責務境界"), phase18MarkdownRef("docs/details/admin.md#0-責務境界"), "admin-cli-passed", "not_applicable", true},
	"sdk-client":       {"phase18-bootstrap-runtime-flow", phase18MarkdownRef("expected/session_flow.json"), phase18MarkdownRef("docs/details/sdk.md#0-責務境界"), phase18MarkdownRef("docs/details/sdk.md#0-責務境界"), "sdk-client-passed", "not_applicable", true},
	"ui-runtime":       {"phase18-bootstrap-runtime-flow", phase18MarkdownRef("expected/session_flow.json"), phase18MarkdownRef("docs/details/ui.md#0-責務境界"), phase18MarkdownRef("docs/details/ui.md#0-責務境界"), "ui-runtime-passed", "not_applicable", true},
	"state-directory":  {"phase18-bootstrap-artifact-health", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/statefile.md#0-責務境界"), phase18MarkdownRef("docs/details/statefile.md#0-責務境界"), "state-directory-safe", "not_applicable", false},
	"log-write":        {"phase18-bootstrap-artifact-health", phase18MarkdownRef("expected/session_flow.json"), phase18MarkdownRef("docs/details/statefile.md#0-責務境界"), phase18MarkdownRef("docs/details/statefile.md#0-責務境界"), "audit-access-config-log-write", "passed", false},
}

type phase18RuntimeTargetExpectation struct {
	CheckName          string
	CounterKey         string
	SourceRef          string
	OwnerRef           string
	OperationRef       string
	QueueState         string
	BuildResult        string
	StateUpdateResult  string
	AuditResult        string
	CommitStatusResult string
	DeployResult       string
}

var phase18RuntimeTargetExpectations = map[string]phase18RuntimeTargetExpectation{
	"admin-api-sdk-ui":        {"phase18-bootstrap-runtime-flow", "phase18_bootstrap_runtime_open_count", phase18MarkdownRef("expected/session_flow.json"), phase18MarkdownRef("docs/details/api.md#0-責務境界"), phase18MarkdownRef("docs/details/admin.md#0-責務境界"), "not_applicable", "not_applicable", "not_applicable", "not_applicable", "not_applicable", "not_applicable"},
	"runner-queue":            {"phase18-pull-runner-foundation", "phase18_pull_runner_foundation_open_count", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/runner.md#0-責務境界"), phase18MarkdownRef("docs/details/runner.md#0-責務境界"), "succeeded", "not_applicable", "not_applicable", "not_applicable", "not_applicable", "not_applicable"},
	"build-execution":         {"phase18-pull-runner-foundation", "phase18_pull_runner_foundation_open_count", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/builder.md#0-責務境界"), phase18MarkdownRef("docs/details/runner.md#0-責務境界"), "not_applicable", "passed", "not_applicable", "not_applicable", "not_applicable", "not_applicable"},
	"statefile-update":        {"phase18-pull-runner-foundation", "phase18_pull_runner_foundation_open_count", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/statefile.md#0-責務境界"), phase18MarkdownRef("docs/details/statefile.md#0-責務境界"), "not_applicable", "not_applicable", "passed", "not_applicable", "not_applicable", "not_applicable"},
	"audit-access-config-log": {"phase18-pull-runner-foundation", "phase18_pull_runner_foundation_open_count", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/statefile.md#0-責務境界"), phase18MarkdownRef("docs/details/statefile.md#0-責務境界"), "not_applicable", "not_applicable", "not_applicable", "passed", "not_applicable", "not_applicable"},
	"commit-status":           {"phase18-pull-runner-foundation", "phase18_pull_runner_foundation_open_count", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/commitstatus.md#0-責務境界"), phase18MarkdownRef("docs/details/commitstatus.md#0-責務境界"), "not_applicable", "not_applicable", "not_applicable", "not_applicable", "passed", "not_applicable"},
	"deploy-simulation":       {"phase18-pull-runner-foundation", "phase18_pull_runner_foundation_open_count", phase18MarkdownRef("input/operation_scope.json"), phase18MarkdownRef("docs/details/runner.md#0-責務境界"), phase18MarkdownRef("docs/details/runner.md#0-責務境界"), "not_applicable", "not_applicable", "not_applicable", "not_applicable", "not_applicable", "passed"},
}

var phase18RequiredCheckRecordFiles = map[string][]string{
	"phase18-bootstrap-prerequisite":    {"records/session.jsonl"},
	"phase18-bootstrap-session":         {"records/session.jsonl"},
	"phase18-bootstrap-artifact-health": {"records/health.jsonl"},
	"phase18-bootstrap-runtime-flow":    {"records/health.jsonl", "records/runtime.jsonl"},
	"phase18-pull-runner-foundation":    {"records/runtime.jsonl"},
	"phase18-bootstrap-update-rollback": {"records/update_rollback.jsonl"},
	"phase18-bugfix-spec-first":         {"records/bugfix.jsonl"},
	"phase18-security-boundary":         {"records/security.jsonl"},
	"phase18-document-drift":            {"records/document_drift.jsonl"},
}

var phase18RequiredCheckCounterKeys = map[string][]string{
	"phase18-bootstrap-prerequisite":    {"phase18_bootstrap_prerequisite_open_count"},
	"phase18-bootstrap-session":         {"phase18_bootstrap_session_open_count"},
	"phase18-bootstrap-artifact-health": {"phase18_bootstrap_runtime_open_count"},
	"phase18-bootstrap-runtime-flow":    {"phase18_bootstrap_runtime_open_count"},
	"phase18-pull-runner-foundation":    {"phase18_pull_runner_foundation_open_count"},
	"phase18-bootstrap-update-rollback": {"phase18_bootstrap_update_rollback_open_count"},
	"phase18-bugfix-spec-first":         {"phase18_bugfix_spec_gap_count", "phase18_known_bug_open_count"},
	"phase18-security-boundary":         {"phase18_secret_leak_open_count", "phase18_destructive_operation_open_count"},
	"phase18-document-drift":            {"phase18_document_drift_open_count"},
}

type phase18NegativeBoundarySpec struct {
	Condition      string
	Group          string
	DiagnosticCode string
	CounterKey     string
}

var phase18NegativeBoundarySpecs = []phase18NegativeBoundarySpec{
	{"phase17-handover-missing-treated-as-success", "bootstrap-prerequisite", "record-link-missing", "phase18_bootstrap_prerequisite_open_count"},
	{"bootstrap-session-three-samples-missing-treated-as-success", "bootstrap-session", "session-mismatch", "phase18_bootstrap_session_open_count"},
	{"handover-open-sample-counted-as-session-open", "bootstrap-session", "session-mismatch", "phase18_bootstrap_session_open_count"},
	{"twenty-four-hour-monitoring-absence-treated-as-blocker", "bootstrap-session", "session-mismatch", "phase18_bootstrap_session_open_count"},
	{"bootstrap-script-command-contract-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"bootstrap-script-contract-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"bootstrap-script-failure-code-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"bootstrap-script-output-contract-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"release-asset-or-version-evidence-missing-treated-as-success", "bootstrap-runtime", "record-link-missing", "phase18_bootstrap_runtime_open_count"},
	{"runner-evidence-missing-counter-zero", "pull-runner-foundation", "counter-mismatch", "phase18_pull_runner_foundation_open_count"},
	{"update-rollback-missing-treated-as-success", "update-rollback", "record-link-missing", "phase18_bootstrap_update_rollback_open_count"},
	{"hotfix-without-spec-treated-as-success", "bugfix-spec-first", "record-link-missing", "phase18_bugfix_spec_gap_count"},
	{"known-bug-open-treated-as-complete", "known-bug", "schema-mismatch", "phase18_known_bug_open_count"},
	{"provider-identifier-written-to-evidence", "secret-boundary", "secret-boundary-failure", "phase18_secret_leak_open_count"},
	{"secret-written-to-output-or-record", "secret-boundary", "secret-boundary-failure", "phase18_secret_leak_open_count"},
	{"forbidden-provider-operation-treated-as-success", "destructive-operation", "destructive-operation-boundary-failure", "phase18_destructive_operation_open_count"},
	{"checker-artifact-drift-ignored", "document-drift", "document-drift", "phase18_document_drift_open_count"},
	{"closure-counter-drift-ignored", "document-drift", "counter-mismatch", "phase18_document_drift_open_count"},
	{"document-drift-ignored", "document-drift", "document-drift", "phase18_document_drift_open_count"},
	{"screenshot-only-treated-as-success", "document-drift", "record-link-missing", "phase18_document_drift_open_count"},
	{"xserver-future-treated-as-blocker", "document-drift", "document-drift", "phase18_document_drift_open_count"},
	{"xserver-future-used-as-required-check-success", "document-drift", "required-check-mismatch", "phase18_document_drift_open_count"},
}

type phase18CounterFile struct {
	SchemaVersion int            `json:"schema_version"`
	Scope         string         `json:"scope"`
	Counters      map[string]int `json:"counters"`
}

type phase18EnvironmentHandoverFile struct {
	SchemaVersion             int                    `json:"schema_version"`
	Scope                     string                 `json:"scope"`
	Phase17Ref                string                 `json:"phase17_ref"`
	ProviderTarget            string                 `json:"provider_target"`
	OS                        string                 `json:"os"`
	MinimumPlanClass          string                 `json:"minimum_plan_class"`
	MinimumRAMMB              int                    `json:"minimum_ram_mb"`
	EnvironmentIdentityPolicy string                 `json:"environment_identity_policy"`
	HandoverResults           phase18HandoverResults `json:"handover_results"`
}

type phase18HandoverResults struct {
	Phase17Status         string                         `json:"phase17_status"`
	Phase17RequiredChecks phase18Phase17RequiredCheckSet `json:"phase17_required_checks"`
	Phase17OpenItems      int                            `json:"phase17_open_items"`
	BootstrapPrerequisite string                         `json:"bootstrap_prerequisite"`
	Result                string                         `json:"result"`
}

type phase18Phase17RequiredCheckSet struct {
	PassedChecks  []string `json:"passed_checks"`
	FailedChecks  []string `json:"failed_checks"`
	MissingChecks []string `json:"missing_checks"`
}

type phase18SessionPlanFile struct {
	SchemaVersion                int      `json:"schema_version"`
	Scope                        string   `json:"scope"`
	SessionKind                  string   `json:"session_kind"`
	RequiredSessionCount         int      `json:"required_session_count"`
	RequiredOperationSampleCount int      `json:"required_operation_sample_count"`
	WorkUnits                    []string `json:"work_units"`
	AbortConditions              []string `json:"abort_conditions"`
	BugfixSpecFirstRequired      bool     `json:"bugfix_spec_first_required"`
	KnownBugZeroRequired         bool     `json:"known_bug_zero_required"`
}

type phase18OperationScopeFile struct {
	SchemaVersion               int      `json:"schema_version"`
	Scope                       string   `json:"scope"`
	ProviderTarget              string   `json:"provider_target"`
	ValidationMode              string   `json:"validation_mode"`
	BootstrapScriptArtifact     string   `json:"bootstrap_script_artifact"`
	BootstrapScriptRuntime      string   `json:"bootstrap_script_runtime"`
	BootstrapScriptInvocation   string   `json:"bootstrap_script_invocation"`
	SourceChannels              []string `json:"source_channels"`
	ActiveSourceChannel         string   `json:"active_source_channel"`
	TopologyRoles               []string `json:"topology_roles"`
	MinimumTopologyRole         string   `json:"minimum_topology_role"`
	BootstrapSourceIDFormat     string   `json:"bootstrap_source_identifier_format"`
	BootstrapDigestAlgorithm    string   `json:"bootstrap_digest_algorithm"`
	BootstrapInstallDirectory   string   `json:"bootstrap_install_directory"`
	BootstrapBinDirectory       string   `json:"bootstrap_bin_directory"`
	BootstrapStateDirectory     string   `json:"bootstrap_state_directory"`
	BootstrapServiceUser        string   `json:"bootstrap_service_user"`
	BootstrapSecretInputPolicy  string   `json:"bootstrap_secret_input_policy"`
	BootstrapStateModel         string   `json:"bootstrap_state_model"`
	BootstrapCLIContract        string   `json:"bootstrap_cli_contract"`
	BootstrapInvocationTokens   []string `json:"bootstrap_invocation_tokens"`
	BootstrapRequiredCommands   []string `json:"bootstrap_required_commands"`
	BootstrapHTTPClientOrder    []string `json:"bootstrap_http_client_order"`
	BootstrapStdoutContract     string   `json:"bootstrap_stdout_contract"`
	BootstrapStderrContract     string   `json:"bootstrap_stderr_contract"`
	BootstrapExitCodeContract   string   `json:"bootstrap_exit_code_contract"`
	BootstrapFailureCodes       []string `json:"bootstrap_failure_codes"`
	OperationTargets            []string `json:"operation_targets"`
	ForbiddenProviderOperations []string `json:"forbidden_provider_operations"`
	XserverFuturePolicy         string   `json:"xserver_future_policy"`
}

type phase18SecurityBoundaryFile struct {
	SchemaVersion                  int      `json:"schema_version"`
	Scope                          string   `json:"scope"`
	ProviderTarget                 string   `json:"provider_target"`
	ValidationMode                 string   `json:"validation_mode"`
	SecretReferencePolicy          string   `json:"secret_reference_policy"`
	MetadataPolicy                 string   `json:"metadata_policy"`
	SecretClasses                  []string `json:"secret_classes"`
	ForbiddenEvidenceValues        []string `json:"forbidden_evidence_values"`
	DestructiveOperationBoundaries []string `json:"destructive_operation_boundaries"`
	MaskRequired                   bool     `json:"mask_required"`
}

type phase18SessionFlowFile struct {
	SchemaVersion      int      `json:"schema_version"`
	Scope              string   `json:"scope"`
	SessionKind        string   `json:"session_kind"`
	OperationSamples   []string `json:"operation_samples"`
	SampleObservations []string `json:"sample_observations"`
	WorkUnits          []string `json:"work_units"`
	NoOpPolicy         string   `json:"no_op_policy"`
	ClosurePolicy      string   `json:"closure_policy"`
}

type phase18BugfixClosureFile struct {
	SchemaVersion          int      `json:"schema_version"`
	Scope                  string   `json:"scope"`
	IssueFlow              []string `json:"issue_flow"`
	TriageResults          []string `json:"triage_results"`
	FixModes               []string `json:"fix_modes"`
	RequiredRefs           []string `json:"required_refs"`
	TerminalKnownBugStatus string   `json:"terminal_known_bug_status"`
	DisallowedFixModes     []string `json:"disallowed_fix_modes"`
}

type phase18NegativeBoundariesFile struct {
	SchemaVersion int                       `json:"schema_version"`
	Scope         string                    `json:"scope"`
	Cases         []phase18NegativeBoundary `json:"cases"`
}

type phase18NegativeBoundary struct {
	CaseID                    string `json:"case_id"`
	BoundaryGroup             string `json:"boundary_group"`
	ForbiddenSuccessCondition string `json:"forbidden_success_condition"`
	ExpectedDiagnosticCode    string `json:"expected_diagnostic_code"`
	ExpectedCounterKey        string `json:"expected_counter_key"`
	SourceRef                 string `json:"source_ref"`
}

type phase18InputSet struct {
	Handover phase18EnvironmentHandoverFile
	Session  phase18SessionPlanFile
	Scope    phase18OperationScopeFile
	Security phase18SecurityBoundaryFile
	Flow     phase18SessionFlowFile
	Bugfix   phase18BugfixClosureFile
}

type phase18EvidenceRecord struct {
	RecordID                           string   `json:"record_id"`
	Scope                              string   `json:"scope"`
	EvidenceOrigin                     string   `json:"evidence_origin"`
	ExecutionEnvironment               string   `json:"execution_environment"`
	ProviderTarget                     string   `json:"provider_target"`
	ValidationMode                     string   `json:"validation_mode"`
	CheckName                          string   `json:"check_name"`
	SourceRef                          string   `json:"source_ref"`
	Result                             string   `json:"result"`
	CounterKey                         string   `json:"counter_key"`
	StartedAt                          *string  `json:"started_at"`
	EndedAt                            *string  `json:"ended_at"`
	StdoutSHA256                       *string  `json:"stdout_sha256"`
	StderrSHA256                       *string  `json:"stderr_sha256"`
	SecretScanResult                   string   `json:"secret_scan_result"`
	DestructiveOperationResult         string   `json:"destructive_operation_result"`
	EvidenceRef                        string   `json:"evidence_ref"`
	SessionID                          string   `json:"session_id,omitempty"`
	SessionKind                        string   `json:"session_kind,omitempty"`
	OperatorApprovalRef                string   `json:"operator_approval_ref,omitempty"`
	SampleIndex                        *int     `json:"sample_index,omitempty"`
	OperationSample                    string   `json:"operation_sample,omitempty"`
	WorkUnit                           string   `json:"work_unit,omitempty"`
	SampleObservations                 []string `json:"sample_observations,omitempty"`
	ServiceState                       string   `json:"service_state,omitempty"`
	APIHealth                          string   `json:"api_health,omitempty"`
	AdminCLIHealth                     string   `json:"admin_cli_health,omitempty"`
	RunnerState                        string   `json:"runner_state,omitempty"`
	StatefileDigest                    string   `json:"statefile_digest,omitempty"`
	AuditLogWrite                      string   `json:"audit_log_write,omitempty"`
	AccessLogWrite                     string   `json:"access_log_write,omitempty"`
	ConfigLogWrite                     string   `json:"config_log_write,omitempty"`
	SecretLeakageBoundary              string   `json:"secret_leakage_boundary,omitempty"`
	KnownBugStatus                     string   `json:"known_bug_status,omitempty"`
	DocumentDriftStatus                string   `json:"document_drift_status,omitempty"`
	AbortCondition                     string   `json:"abort_condition,omitempty"`
	CloseResult                        string   `json:"close_result,omitempty"`
	HealthTarget                       string   `json:"health_target,omitempty"`
	OwnerRef                           string   `json:"owner_ref,omitempty"`
	EndpointOrCommandRef               string   `json:"endpoint_or_command_ref,omitempty"`
	ExpectedState                      string   `json:"expected_state,omitempty"`
	ActualState                        string   `json:"actual_state,omitempty"`
	HealthResult                       string   `json:"health_result,omitempty"`
	LatencyMS                          *int     `json:"latency_ms,omitempty"`
	LogWriteResult                     string   `json:"log_write_result,omitempty"`
	StateReadResult                    string   `json:"state_read_result,omitempty"`
	SecretMaskResult                   string   `json:"secret_mask_result,omitempty"`
	Diagnostic                         string   `json:"diagnostic,omitempty"`
	ClosureAction                      string   `json:"closure_action,omitempty"`
	RuntimeTarget                      string   `json:"runtime_target,omitempty"`
	OperationRef                       string   `json:"operation_ref,omitempty"`
	QueueState                         string   `json:"queue_state,omitempty"`
	BuildResult                        string   `json:"build_result,omitempty"`
	StateUpdateResult                  string   `json:"state_update_result,omitempty"`
	AuditResult                        string   `json:"audit_result,omitempty"`
	CommitStatusResult                 string   `json:"commit_status_result,omitempty"`
	DeployResult                       string   `json:"deploy_result,omitempty"`
	NoCustomerDataResult               string   `json:"no_customer_data_result,omitempty"`
	OperationType                      string   `json:"operation_type,omitempty"`
	FromVersion                        string   `json:"from_version,omitempty"`
	ToVersion                          string   `json:"to_version,omitempty"`
	ServiceRestartResult               string   `json:"service_restart_result,omitempty"`
	StateMigrationResult               string   `json:"state_migration_result,omitempty"`
	RollbackState                      string   `json:"rollback_state,omitempty"`
	FailureCleanupResult               string   `json:"failure_cleanup_result,omitempty"`
	VersionAfterOperation              string   `json:"version_after_operation,omitempty"`
	IssueID                            string   `json:"issue_id,omitempty"`
	IssueClassification                string   `json:"issue_classification,omitempty"`
	TriageResult                       string   `json:"triage_result,omitempty"`
	FixMode                            string   `json:"fix_mode,omitempty"`
	IssueRef                           string   `json:"issue_ref,omitempty"`
	SpecGeneralUpdateRef               string   `json:"spec_general_update_ref,omitempty"`
	ResponsibilitySourceRef            string   `json:"responsibility_source_ref,omitempty"`
	ImplementationFixRef               string   `json:"implementation_fix_ref,omitempty"`
	RevalidationRef                    string   `json:"revalidation_ref,omitempty"`
	SecretClass                        string   `json:"secret_class,omitempty"`
	ScanTarget                         string   `json:"scan_target,omitempty"`
	MetadataPolicyResult               string   `json:"metadata_policy_result,omitempty"`
	ForbiddenValueDetected             *bool    `json:"forbidden_value_detected,omitempty"`
	ForbiddenProviderOperationToken    string   `json:"forbidden_provider_operation_token,omitempty"`
	ForbiddenProviderOperationDetected *bool    `json:"forbidden_provider_operation_detected,omitempty"`
	BoundaryResult                     string   `json:"boundary_result,omitempty"`
	MaskResult                         string   `json:"mask_result,omitempty"`
	StorageResult                      string   `json:"storage_result,omitempty"`
	DriftTarget                        string   `json:"drift_target,omitempty"`
	DocumentRef                        string   `json:"document_ref,omitempty"`
	AnchorRef                          string   `json:"anchor_ref,omitempty"`
	FixturePathRef                     string   `json:"fixture_path_ref,omitempty"`
	WorkflowPathRef                    string   `json:"workflow_path_ref,omitempty"`
	ExpectedRef                        string   `json:"expected_ref,omitempty"`
	ActualRef                          string   `json:"actual_ref,omitempty"`
	DriftResult                        string   `json:"drift_result,omitempty"`
	CorrectiveActionRef                string   `json:"corrective_action_ref,omitempty"`
	ClosureResult                      string   `json:"closure_result,omitempty"`
}

func phase18RequireFixtureFiles(t *testing.T, root string) {
	t.Helper()

	info, err := os.Lstat(root)
	if err != nil {
		t.Fatalf("Phase 18 fixture root must exist: %v", err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("Phase 18 fixture root must be a directory and not a symlink: %s", root)
	}
	allowed := map[string]bool{}
	for _, rel := range phase18FixtureFiles {
		allowed[rel] = true
		path := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Lstat(path)
		if err != nil {
			t.Fatalf("Phase 18 fixture path missing %s: %v", rel, err)
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			t.Fatalf("Phase 18 fixture path must be a regular file: %s", rel)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if len(data) == 0 || data[len(data)-1] != '\n' {
			t.Fatalf("%s must be non-empty and end with LF", path)
		}
		if strings.HasPrefix(string(data), "\ufeff") {
			t.Fatalf("%s must not include UTF-8 BOM", path)
		}
		if phase18ContainsForbiddenEvidenceLiteral(string(data)) {
			t.Fatalf("%s contains forbidden secret-like, endpoint-like, or provider-identity literal", path)
		}
		if strings.HasSuffix(rel, ".json") {
			phase16RequireNoDuplicateJSONKeys(t, path, data)
		}
		if strings.HasSuffix(rel, ".jsonl") {
			for index, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
				if strings.TrimSpace(line) == "" {
					t.Fatalf("%s line %d must not be empty", path, index+1)
				}
				phase16RequireNoDuplicateJSONKeys(t, fmt.Sprintf("%s line %d", path, index+1), []byte(line))
			}
		}
	}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root || entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !allowed[rel] {
			t.Fatalf("Phase 18 fixture contains unknown file %s", rel)
		}
		return nil
	}); err != nil {
		t.Fatalf("walk Phase 18 fixture: %v", err)
	}
}

func phase18ReadCounters(t *testing.T, path string) map[string]int {
	t.Helper()

	var file phase18CounterFile
	phase16ReadStrictJSON(t, path, &file)
	if file.SchemaVersion != 1 || file.Scope != phase18Scope {
		t.Fatalf("%s must use schema_version=1 and scope=%s", path, phase18Scope)
	}
	return file.Counters
}

func phase18RequireCountersClosed(t *testing.T, counters map[string]int) {
	t.Helper()

	if len(counters) != len(phase18Counters) {
		t.Fatalf("Phase 18 counters must include exactly %d keys, got %d", len(phase18Counters), len(counters))
	}
	var sum int
	for _, key := range phase18Counters {
		value, ok := counters[key]
		if !ok {
			t.Fatalf("Phase 18 counters missing %s", key)
		}
		if value < 0 {
			t.Fatalf("Phase 18 counter %s must not be negative", key)
		}
		if key != "final_open_item_count" {
			sum += value
		}
	}
	if counters["final_open_item_count"] != sum {
		t.Fatalf("Phase 18 final_open_item_count must equal open counter sum: want %d got %d", sum, counters["final_open_item_count"])
	}
	if counters["final_open_item_count"] != 0 {
		t.Fatalf("Phase 18 checker acceptance fixture requires final_open_item_count=0, got %d", counters["final_open_item_count"])
	}
}

func phase18RequireInputsAndExpected(t *testing.T, root string) phase18InputSet {
	t.Helper()

	var inputs phase18InputSet
	phase16ReadStrictJSON(t, filepath.Join(root, "input/environment_handover.json"), &inputs.Handover)
	if inputs.Handover.SchemaVersion != 1 || inputs.Handover.Scope != phase18Scope || inputs.Handover.ProviderTarget != "conoha-vps-primary" || inputs.Handover.OS != "ubuntu-server-24.04-lts-64bit" || inputs.Handover.MinimumPlanClass != "conoha-vps-1gb-memory-class" || inputs.Handover.MinimumRAMMB != 1024 || inputs.Handover.EnvironmentIdentityPolicy != "opaque-non-secret-labels" {
		t.Fatalf("Phase 18 environment handover fixed values mismatch: %+v", inputs.Handover)
	}
	phase18RequireSafeReference(t, "input/environment_handover.json phase17_ref", inputs.Handover.Phase17Ref)
	phase18RequireMarkdownReferenceForm(t, "input/environment_handover.json phase17_ref", inputs.Handover.Phase17Ref)
	phase18RequireResolvableReference(t, "input/environment_handover.json phase17_ref", inputs.Handover.Phase17Ref)
	phase18RequireCanonicalReferenceLabel(t, "input/environment_handover.json phase17_ref", inputs.Handover.Phase17Ref)
	if inputs.Handover.HandoverResults.Phase17Status != "implemented" || inputs.Handover.HandoverResults.Phase17OpenItems != 0 || inputs.Handover.HandoverResults.BootstrapPrerequisite != "passed" || inputs.Handover.HandoverResults.Result != "passed" {
		t.Fatalf("Phase 18 handover must close Phase 17 prerequisites: %+v", inputs.Handover.HandoverResults)
	}
	if !phase11StringSlicesEqual(inputs.Handover.HandoverResults.Phase17RequiredChecks.PassedChecks, phase18Phase17RequiredChecks) || len(inputs.Handover.HandoverResults.Phase17RequiredChecks.FailedChecks) != 0 || len(inputs.Handover.HandoverResults.Phase17RequiredChecks.MissingChecks) != 0 {
		t.Fatalf("Phase 18 handover Phase 17 required checks mismatch: %+v", inputs.Handover.HandoverResults.Phase17RequiredChecks)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "input/session_plan.json"), &inputs.Session)
	if inputs.Session.SchemaVersion != 1 || inputs.Session.Scope != phase18Scope || inputs.Session.SessionKind != "operator-approved-bootstrap-session" || inputs.Session.RequiredSessionCount != 1 || inputs.Session.RequiredOperationSampleCount != 3 || !phase11StringSlicesEqual(inputs.Session.WorkUnits, phase18WorkUnits) || !phase11StringSlicesEqual(inputs.Session.AbortConditions, phase18AbortConditions) || !inputs.Session.BugfixSpecFirstRequired || !inputs.Session.KnownBugZeroRequired {
		t.Fatalf("Phase 18 session plan mismatch: %+v", inputs.Session)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "input/operation_scope.json"), &inputs.Scope)
	if inputs.Scope.SchemaVersion != 1 || inputs.Scope.Scope != phase18Scope || inputs.Scope.ProviderTarget != "conoha-vps-primary" || inputs.Scope.ValidationMode != "vps-pull-bootstrap" || inputs.Scope.BootstrapScriptArtifact != "adlaire-ci-vps-pull-bootstrap.sh" || inputs.Scope.BootstrapScriptRuntime != "/bin/sh" || inputs.Scope.BootstrapScriptInvocation != "operator-runs-sh-script-on-vps" || !phase11StringSlicesEqual(inputs.Scope.SourceChannels, phase18SourceChannels) || inputs.Scope.ActiveSourceChannel != "integration-head" || !phase11StringSlicesEqual(inputs.Scope.TopologyRoles, phase18TopologyRoles) || inputs.Scope.MinimumTopologyRole != "single-node" || inputs.Scope.BootstrapSourceIDFormat != "absolute-https-url-or-absolute-file-path" || inputs.Scope.BootstrapDigestAlgorithm != "sha256" || inputs.Scope.BootstrapInstallDirectory != "/opt/adlaire-builder" || inputs.Scope.BootstrapBinDirectory != "/usr/local/bin" || inputs.Scope.BootstrapStateDirectory != "/opt/adlaire-builder" || inputs.Scope.BootstrapServiceUser != "root" || inputs.Scope.BootstrapSecretInputPolicy != "no-command-argument-secret" || inputs.Scope.BootstrapStateModel != "staging-verify-commit-rollback" || inputs.Scope.BootstrapCLIContract != "fixed-token-flag-interface" || !phase11StringSlicesEqual(inputs.Scope.BootstrapInvocationTokens, phase18BootstrapInvocationTokens) || !phase11StringSlicesEqual(inputs.Scope.BootstrapRequiredCommands, phase18BootstrapRequiredCommands) || !phase11StringSlicesEqual(inputs.Scope.BootstrapHTTPClientOrder, phase18BootstrapHTTPClientOrder) || inputs.Scope.BootstrapStdoutContract != "single-json-object-secret-safe" || inputs.Scope.BootstrapStderrContract != "diagnostic-json-lines-on-failure" || inputs.Scope.BootstrapExitCodeContract != "common-cli-exit-code-contract" || !phase11StringSlicesEqual(inputs.Scope.BootstrapFailureCodes, phase18BootstrapFailureCodes) || !phase11StringSlicesEqual(inputs.Scope.OperationTargets, phase18OperationTargets) || !phase11StringSlicesEqual(inputs.Scope.ForbiddenProviderOperations, phase18ForbiddenProviderOperations) || inputs.Scope.XserverFuturePolicy != "future_plan_not_blocker" {
		t.Fatalf("Phase 18 operation scope mismatch: %+v", inputs.Scope)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "input/security_boundary.json"), &inputs.Security)
	if inputs.Security.SchemaVersion != 1 || inputs.Security.Scope != phase18Scope || inputs.Security.ProviderTarget != "conoha-vps-primary" || inputs.Security.ValidationMode != "vps-pull-bootstrap" || inputs.Security.SecretReferencePolicy != "metadata-only" || inputs.Security.MetadataPolicy != "opaque-non-secret-labels" || !phase11StringSlicesEqual(inputs.Security.SecretClasses, phase18SecretClasses) || !phase11StringSlicesEqual(inputs.Security.ForbiddenEvidenceValues, phase18ForbiddenEvidenceValues) || !phase11StringSlicesEqual(inputs.Security.DestructiveOperationBoundaries, phase18ForbiddenProviderOperations) || !inputs.Security.MaskRequired {
		t.Fatalf("Phase 18 security boundary mismatch: %+v", inputs.Security)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "expected/session_flow.json"), &inputs.Flow)
	if inputs.Flow.SchemaVersion != 1 || inputs.Flow.Scope != phase18Scope || inputs.Flow.SessionKind != inputs.Session.SessionKind || !phase11StringSlicesEqual(inputs.Flow.OperationSamples, phase18OperationSamples) || !phase11StringSlicesEqual(inputs.Flow.SampleObservations, phase18SampleObservations) || !phase11StringSlicesEqual(inputs.Flow.WorkUnits, inputs.Session.WorkUnits) || inputs.Flow.NoOpPolicy != "explicit-not-required-record-required" || inputs.Flow.ClosurePolicy != "single-session-id-required" {
		t.Fatalf("Phase 18 session flow mismatch: %+v", inputs.Flow)
	}

	phase16ReadStrictJSON(t, filepath.Join(root, "expected/bugfix_closure.json"), &inputs.Bugfix)
	if inputs.Bugfix.SchemaVersion != 1 || inputs.Bugfix.Scope != phase18Scope || !phase11StringSlicesEqual(inputs.Bugfix.IssueFlow, phase18IssueFlow) || !phase11StringSlicesEqual(inputs.Bugfix.TriageResults, phase18TriageResults) || !phase11StringSlicesEqual(inputs.Bugfix.FixModes, phase18FixModes) || !phase11StringSlicesEqual(inputs.Bugfix.RequiredRefs, phase18RequiredBugfixRefs) || inputs.Bugfix.TerminalKnownBugStatus != "none" || !phase11StringSlicesEqual(inputs.Bugfix.DisallowedFixModes, phase18DisallowedFixModes) {
		t.Fatalf("Phase 18 bugfix closure mismatch: %+v", inputs.Bugfix)
	}

	var negative phase18NegativeBoundariesFile
	phase16ReadStrictJSON(t, filepath.Join(root, "expected/negative_boundaries.json"), &negative)
	if negative.SchemaVersion != 1 || negative.Scope != phase18Scope || len(negative.Cases) != len(phase18NegativeBoundarySpecs) {
		t.Fatalf("Phase 18 negative boundary file mismatch: %+v", negative)
	}
	seenCondition := map[string]bool{}
	for index, want := range phase18NegativeBoundarySpecs {
		got := negative.Cases[index]
		wantID := "phase18.negative." + want.Group + "." + want.Condition
		if got.CaseID != wantID || got.BoundaryGroup != want.Group || got.ForbiddenSuccessCondition != want.Condition || got.ExpectedDiagnosticCode != want.DiagnosticCode || got.ExpectedCounterKey != want.CounterKey {
			t.Fatalf("Phase 18 negative boundary %d mismatch\nwant: %+v / %s\n got: %+v", index+1, want, wantID, got)
		}
		if seenCondition[got.ForbiddenSuccessCondition] {
			t.Fatalf("Phase 18 negative boundary repeats forbidden_success_condition %s", got.ForbiddenSuccessCondition)
		}
		seenCondition[got.ForbiddenSuccessCondition] = true
		if got.SourceRef != phase18NegativeBoundarySourceRef {
			t.Fatalf("Phase 18 negative boundary %s source_ref mismatch: want %s got %s", got.CaseID, phase18NegativeBoundarySourceRef, got.SourceRef)
		}
		phase18RequireSafeReference(t, "negative boundary source_ref", got.SourceRef)
		phase18RequireMarkdownReferenceForm(t, "negative boundary source_ref", got.SourceRef)
		phase18RequireResolvableReference(t, "negative boundary source_ref", got.SourceRef)
	}
	return inputs
}

func phase18ReadRecords(t *testing.T, root string) map[string][]phase18EvidenceRecord {
	t.Helper()

	records := map[string][]phase18EvidenceRecord{}
	for _, rel := range phase18RecordFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
		if len(lines) == 0 {
			t.Fatalf("%s must include at least one JSONL record", path)
		}
		seen := map[string]bool{}
		for index, line := range lines {
			if strings.TrimSpace(line) == "" {
				t.Fatalf("%s line %d must not be empty", path, index+1)
			}
			phase16RequireNoDuplicateJSONKeys(t, fmt.Sprintf("%s line %d", path, index+1), []byte(line))
			var record phase18EvidenceRecord
			decoder := json.NewDecoder(strings.NewReader(line))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&record); err != nil {
				t.Fatalf("decode %s line %d: %v", path, index+1, err)
			}
			var trailing any
			if err := decoder.Decode(&trailing); err != io.EOF {
				t.Fatalf("%s line %d must contain exactly one JSON object", path, index+1)
			}
			phase18RequireRecordCommon(t, rel, index+1, record)
			if seen[record.RecordID] {
				t.Fatalf("%s repeats record_id %s", path, record.RecordID)
			}
			seen[record.RecordID] = true
			phase18RequireRecordSpecific(t, rel, index+1, record)
			records[rel] = append(records[rel], record)
		}
	}
	return records
}

func phase18RequireRecordCommon(t *testing.T, rel string, line int, record phase18EvidenceRecord) {
	t.Helper()

	prefix := "phase18.record." + strings.TrimSuffix(filepath.Base(rel), ".jsonl") + "."
	if record.RecordID == "" || !strings.HasPrefix(record.RecordID, prefix) || record.Scope != phase18Scope || record.EvidenceOrigin != "checker-acceptance-fixture" || record.ExecutionEnvironment != "repository-fixture" || record.ProviderTarget == "" || record.ValidationMode == "" || record.CheckName == "" || record.SourceRef == "" || record.Result == "" || record.CounterKey == "" || record.SecretScanResult == "" || record.DestructiveOperationResult == "" || record.EvidenceRef == "" {
		t.Fatalf("%s line %d common fields invalid: %+v", rel, line, record)
	}
	if !phase17StringInSlice(record.ProviderTarget, phase18ProviderTargets) {
		t.Fatalf("%s line %d unknown provider_target %s", rel, line, record.ProviderTarget)
	}
	if !phase17StringInSlice(record.ValidationMode, phase18ValidationModes) {
		t.Fatalf("%s line %d unknown validation_mode %s", rel, line, record.ValidationMode)
	}
	if !phase17StringInSlice(record.CheckName, phase18RequiredChecks) {
		t.Fatalf("%s line %d unknown check_name %s", rel, line, record.CheckName)
	}
	if !phase17StringInSlice(record.CounterKey, phase18Counters) {
		t.Fatalf("%s line %d unknown counter_key %s", rel, line, record.CounterKey)
	}
	switch record.Result {
	case "passed", "not_applicable", "future_plan":
	case "failed", "rejected":
		t.Fatalf("%s line %d must not leave failed or rejected result", rel, line)
	default:
		t.Fatalf("%s line %d invalid result %s", rel, line, record.Result)
	}
	if record.Result == "passed" {
		if record.StartedAt == nil || record.EndedAt == nil || *record.StartedAt == "" || *record.EndedAt == "" {
			t.Fatalf("%s line %d passed record must include started_at and ended_at", rel, line)
		}
		if _, err := time.Parse(time.RFC3339, *record.StartedAt); err != nil {
			t.Fatalf("%s line %d invalid started_at: %v", rel, line, err)
		}
		if _, err := time.Parse(time.RFC3339, *record.EndedAt); err != nil {
			t.Fatalf("%s line %d invalid ended_at: %v", rel, line, err)
		}
	}
	for label, value := range map[string]*string{"stdout_sha256": record.StdoutSHA256, "stderr_sha256": record.StderrSHA256} {
		if value != nil && !phase17IsLowerHexSHA256(*value) {
			t.Fatalf("%s line %d %s must be lowercase SHA-256 hex or null", rel, line, label)
		}
	}
	switch record.SecretScanResult {
	case "passed", "not_applicable":
	default:
		t.Fatalf("%s line %d invalid secret_scan_result %s", rel, line, record.SecretScanResult)
	}
	switch record.DestructiveOperationResult {
	case "not_applicable", "not_executed", "blocked":
	default:
		t.Fatalf("%s line %d invalid destructive_operation_result %s", rel, line, record.DestructiveOperationResult)
	}
	phase18RequireSafeReference(t, fmt.Sprintf("%s line %d source_ref", rel, line), record.SourceRef)
	phase18RequireSafeReference(t, fmt.Sprintf("%s line %d evidence_ref", rel, line), record.EvidenceRef)
	phase18RequireRecordReferences(t, rel, line, record)
	if record.ProviderTarget == "xserver-vps-future" && (rel != "records/document_drift.jsonl" || record.ValidationMode != "future_plan" || record.Result != "future_plan" || record.CheckName != "phase18-document-drift") {
		t.Fatalf("%s line %d xserver-vps-future is allowed only as a document_drift future plan classification record: %+v", rel, line, record)
	}
}

func phase18RequireRecordSpecific(t *testing.T, rel string, line int, record phase18EvidenceRecord) {
	t.Helper()

	switch rel {
	case "records/session.jsonl":
		if record.SessionID == "" || record.SessionKind != "operator-approved-bootstrap-session" || record.OperatorApprovalRef == "" || record.SampleIndex == nil || record.OperationSample == "" || record.WorkUnit == "" || len(record.SampleObservations) == 0 || record.ServiceState == "" || record.APIHealth == "" || record.AdminCLIHealth == "" || record.RunnerState == "" || record.StatefileDigest == "" || record.AuditLogWrite == "" || record.AccessLogWrite == "" || record.ConfigLogWrite == "" || record.SecretLeakageBoundary == "" || record.KnownBugStatus == "" || record.DocumentDriftStatus == "" || record.AbortCondition == "" || record.CloseResult == "" {
			t.Fatalf("%s line %d missing session fields: %+v", rel, line, record)
		}
		if !phase17IsLowerHexSHA256(record.StatefileDigest) || !phase11StringSlicesEqual(record.SampleObservations, phase18SampleObservations) || !phase17StringInSlice(record.OperationSample, phase18OperationSamples) || (!phase17StringInSlice(record.AbortCondition, phase18AbortConditions) && record.AbortCondition != "not_triggered") {
			t.Fatalf("%s line %d session fixed fields mismatch: %+v", rel, line, record)
		}
	case "records/health.jsonl":
		if record.SessionID == "" || record.HealthTarget == "" || record.OwnerRef == "" || record.EndpointOrCommandRef == "" || record.ExpectedState == "" || record.ActualState == "" || record.HealthResult == "" || record.LogWriteResult == "" || record.StateReadResult == "" || record.SecretMaskResult == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing health fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.HealthTarget, phase18HealthTargets) || record.CounterKey != "phase18_bootstrap_runtime_open_count" || record.HealthResult != "passed" || record.StateReadResult != "passed" || record.SecretMaskResult != "passed" {
			t.Fatalf("%s line %d health result mismatch: %+v", rel, line, record)
		}
		if record.HealthTarget == "log-write" && record.LogWriteResult != "passed" {
			t.Fatalf("%s line %d log-write target must close log_write_result: %+v", rel, line, record)
		}
		if phase17StringInSlice(record.HealthTarget, []string{"api-health", "admin-cli", "sdk-client", "ui-runtime"}) && record.LatencyMS == nil {
			t.Fatalf("%s line %d interactive health target must include latency_ms", rel, line)
		}
	case "records/runtime.jsonl":
		if record.SessionID == "" || record.RuntimeTarget == "" || record.OwnerRef == "" || record.OperationRef == "" || record.QueueState == "" || record.BuildResult == "" || record.StateUpdateResult == "" || record.AuditResult == "" || record.CommitStatusResult == "" || record.DeployResult == "" || record.NoCustomerDataResult == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing runtime fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.RuntimeTarget, phase18RuntimeTargets) || record.NoCustomerDataResult != "passed" {
			t.Fatalf("%s line %d runtime target mismatch: %+v", rel, line, record)
		}
	case "records/update_rollback.jsonl":
		if record.SessionID == "" || record.OperationType == "" || record.FromVersion == "" || record.ToVersion == "" || record.ServiceRestartResult == "" || record.StateMigrationResult == "" || record.RollbackState == "" || record.FailureCleanupResult == "" || record.VersionAfterOperation == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing update/rollback fields: %+v", rel, line, record)
		}
		if !phase17StringInSlice(record.OperationType, phase18UpdateRollbackOperations) || record.CounterKey != "phase18_bootstrap_update_rollback_open_count" {
			t.Fatalf("%s line %d update/rollback operation mismatch: %+v", rel, line, record)
		}
	case "records/bugfix.jsonl":
		if record.SessionID == "" || record.IssueID == "" || record.IssueClassification == "" || record.TriageResult == "" || record.FixMode == "" || record.IssueRef == "" || record.SpecGeneralUpdateRef == "" || record.ResponsibilitySourceRef == "" || record.ImplementationFixRef == "" || record.RevalidationRef == "" || record.KnownBugStatus == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing bugfix fields: %+v", rel, line, record)
		}
	case "records/security.jsonl":
		if record.SessionID == "" || record.SecretClass == "" || record.ScanTarget == "" || record.MetadataPolicyResult == "" || record.ForbiddenValueDetected == nil || record.ForbiddenProviderOperationToken == "" || record.ForbiddenProviderOperationDetected == nil || record.BoundaryResult == "" || record.MaskResult == "" || record.StorageResult == "" || record.Diagnostic == "" || record.ClosureAction == "" {
			t.Fatalf("%s line %d missing security fields: %+v", rel, line, record)
		}
	case "records/document_drift.jsonl":
		if record.SessionID == "" || record.DriftTarget == "" || record.DocumentRef == "" || record.AnchorRef == "" || record.FixturePathRef == "" || record.WorkflowPathRef == "" || record.ExpectedRef == "" || record.ActualRef == "" || record.DriftResult == "" || record.CorrectiveActionRef == "" || record.ClosureResult == "" {
			t.Fatalf("%s line %d missing document drift fields: %+v", rel, line, record)
		}
	}
}

func phase18RequireRequiredCheckConnections(t *testing.T, records map[string][]phase18EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for rel, recordSet := range records {
		for _, record := range recordSet {
			files := phase18RequiredCheckRecordFiles[record.CheckName]
			if !phase17StringInSlice(rel, files) {
				t.Fatalf("Phase 18 record %s places check %s in invalid record file %s", record.RecordID, record.CheckName, rel)
			}
			if !phase17StringInSlice(record.CounterKey, phase18RequiredCheckCounterKeys[record.CheckName]) {
				t.Fatalf("Phase 18 record %s connects check %s to invalid counter %s", record.RecordID, record.CheckName, record.CounterKey)
			}
			if record.ProviderTarget != "xserver-vps-future" {
				seen[record.CheckName] = true
			}
		}
	}
	for _, check := range phase18RequiredChecks {
		if !seen[check] {
			t.Fatalf("Phase 18 required check %s is not connected to evidence records", check)
		}
	}
}

func phase18RequireSessionRecords(t *testing.T, records []phase18EvidenceRecord, inputs phase18InputSet) {
	t.Helper()

	var sessionID string
	operationSamples := map[string]phase18EvidenceRecord{}
	var handover bool
	for _, record := range records {
		if record.ProviderTarget != "conoha-vps-primary" {
			continue
		}
		if !strings.HasPrefix(record.SessionID, "phase18.session.") {
			t.Fatalf("Phase 18 session_id must use phase18.session.<slug>: %s", record.SessionID)
		}
		if sessionID == "" {
			sessionID = record.SessionID
		}
		if record.SessionID != sessionID {
			t.Fatalf("Phase 18 bootstrap records must reuse a single session_id: %s vs %s", record.SessionID, sessionID)
		}
		if record.CheckName == "phase18-bootstrap-prerequisite" {
			if record.WorkUnit != "phase17-handover" || record.OperationSample != "open" || record.CounterKey != "phase18_bootstrap_prerequisite_open_count" {
				t.Fatalf("Phase 18 handover session record mismatch: %+v", record)
			}
			handover = true
			continue
		}
		if record.CheckName == "phase18-bootstrap-session" {
			operationSamples[record.OperationSample] = record
		}
	}
	if !handover {
		t.Fatalf("Phase 18 session records must include environment handover record")
	}
	if len(operationSamples) != inputs.Session.RequiredOperationSampleCount {
		t.Fatalf("Phase 18 bootstrap session must include exactly %d samples, got %d", inputs.Session.RequiredOperationSampleCount, len(operationSamples))
	}
	for index, sample := range phase18OperationSamples {
		record, ok := operationSamples[sample]
		if !ok {
			t.Fatalf("Phase 18 bootstrap session missing sample %s", sample)
		}
		if record.SampleIndex == nil || *record.SampleIndex != index+1 {
			t.Fatalf("Phase 18 bootstrap session sample %s index mismatch: %+v", sample, record)
		}
	}
	closeRecord := operationSamples["close"]
	if closeRecord.CloseResult != "closed" || closeRecord.SecretLeakageBoundary != "passed" || closeRecord.KnownBugStatus != "none" || closeRecord.DocumentDriftStatus != "no_drift" || closeRecord.AbortCondition != "not_triggered" {
		t.Fatalf("Phase 18 close sample must close session without leaks, known bugs, drift, or abort: %+v", closeRecord)
	}
}

func phase18RequireHealthRecords(t *testing.T, records []phase18EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for _, record := range records {
		expected, ok := phase18HealthTargetExpectations[record.HealthTarget]
		if !ok {
			t.Fatalf("Phase 18 health target has no fixed expectation: %+v", record)
		}
		if record.CheckName != expected.CheckName || record.SourceRef != expected.SourceRef || record.OwnerRef != expected.OwnerRef || record.EndpointOrCommandRef != expected.EndpointOrCommandRef || record.ExpectedState != expected.State || record.ActualState != expected.State || record.LogWriteResult != expected.LogWriteResult {
			t.Fatalf("Phase 18 health target fixed references or state mismatch: %+v", record)
		}
		if record.Result != "passed" || record.HealthResult != "passed" || record.StateReadResult != "passed" || record.SecretMaskResult != "passed" {
			t.Fatalf("Phase 18 health target must pass: %+v", record)
		}
		if record.HealthTarget == "bootstrap-script" {
			if !strings.Contains(record.Diagnostic, "source identifier") || !strings.Contains(record.Diagnostic, "secret policy") || !strings.Contains(record.Diagnostic, "cli token") || !strings.Contains(record.Diagnostic, "required command") || !strings.Contains(record.Diagnostic, "http client") || !strings.Contains(record.Diagnostic, "exit code") || !strings.Contains(record.Diagnostic, "failure code") {
				t.Fatalf("Phase 18 bootstrap script health target must prove exact bootstrap contract: %+v", record)
			}
		}
		if expected.LatencyRequired && record.LatencyMS == nil {
			t.Fatalf("Phase 18 interactive health target must include latency_ms: %+v", record)
		}
		if !expected.LatencyRequired && record.LatencyMS != nil {
			t.Fatalf("Phase 18 non-interactive health target must not include latency_ms: %+v", record)
		}
		seen[record.HealthTarget] = true
	}
	for _, target := range phase18HealthTargets {
		if !seen[target] {
			t.Fatalf("Phase 18 health records missing target %s", target)
		}
	}
}

func phase18RequireRuntimeRecords(t *testing.T, records []phase18EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for _, record := range records {
		seen[record.RuntimeTarget] = true
		expected, ok := phase18RuntimeTargetExpectations[record.RuntimeTarget]
		if !ok {
			t.Fatalf("Phase 18 runtime target has no fixed expectation: %+v", record)
		}
		if record.CheckName != expected.CheckName || record.CounterKey != expected.CounterKey || record.SourceRef != expected.SourceRef || record.OwnerRef != expected.OwnerRef || record.OperationRef != expected.OperationRef || record.QueueState != expected.QueueState || record.BuildResult != expected.BuildResult || record.StateUpdateResult != expected.StateUpdateResult || record.AuditResult != expected.AuditResult || record.CommitStatusResult != expected.CommitStatusResult || record.DeployResult != expected.DeployResult || record.NoCustomerDataResult != "passed" {
			t.Fatalf("Phase 18 runtime target fixed references or result mismatch: %+v", record)
		}
	}
	for _, target := range phase18RuntimeTargets {
		if !seen[target] {
			t.Fatalf("Phase 18 runtime records missing target %s", target)
		}
	}
}

func phase18RequireUpdateRollbackRecords(t *testing.T, records []phase18EvidenceRecord) {
	t.Helper()

	seen := map[string]bool{}
	for _, record := range records {
		seen[record.OperationType] = true
		if record.SourceRef != phase18MarkdownRef("input/operation_scope.json") || record.EvidenceRef != phase18MarkdownRef("records/update_rollback.jsonl") {
			t.Fatalf("Phase 18 update/rollback record must use fixed linked source and evidence refs: %+v", record)
		}
		if record.Result != "passed" || record.FromVersion == "unknown" || record.ToVersion == "unknown" || record.VersionAfterOperation == "unknown" {
			t.Fatalf("Phase 18 update/rollback records must pass with concrete versions: %+v", record)
		}
		switch record.OperationType {
		case "update":
			if record.FromVersion == record.ToVersion || record.VersionAfterOperation != record.ToVersion || record.ServiceRestartResult != "passed" || record.StateMigrationResult != "passed" || record.RollbackState != "committed" {
				t.Fatalf("Phase 18 update record mismatch: %+v", record)
			}
		case "rollback":
			if record.VersionAfterOperation != record.ToVersion || record.ServiceRestartResult != "passed" || record.StateMigrationResult != "passed" || record.RollbackState != "rolled_back" {
				t.Fatalf("Phase 18 rollback record mismatch: %+v", record)
			}
		case "cleanup":
			if record.FailureCleanupResult != "passed" || record.RollbackState != "not_applicable" {
				t.Fatalf("Phase 18 cleanup record mismatch: %+v", record)
			}
		default:
			t.Fatalf("Phase 18 unknown update/rollback operation %s", record.OperationType)
		}
	}
	for _, operation := range phase18UpdateRollbackOperations {
		if !seen[operation] {
			t.Fatalf("Phase 18 update/rollback records missing operation %s", operation)
		}
	}
}

func phase18RequireBugfixRecords(t *testing.T, records []phase18EvidenceRecord, expected phase18BugfixClosureFile) {
	t.Helper()

	var specFirst bool
	var noIssue bool
	for _, record := range records {
		if record.SourceRef != phase18MarkdownRef("expected/bugfix_closure.json") || record.EvidenceRef != phase18MarkdownRef("records/bugfix.jsonl") || record.SpecGeneralUpdateRef != phase18MarkdownRef("docs/ROADMAP.md#roadmap-initial-phase-plan") || record.ResponsibilitySourceRef != phase18MarkdownRef("docs/details/production-validation.md#phase-18-vps-pull-bootstrap-contract") || record.RevalidationRef != phase18MarkdownRef("records/bugfix.jsonl") {
			t.Fatalf("Phase 18 bugfix record must use fixed linked source, evidence, spec, responsibility, and revalidation refs: %+v", record)
		}
		if record.CounterKey != "phase18_bugfix_spec_gap_count" && record.CounterKey != "phase18_known_bug_open_count" {
			t.Fatalf("Phase 18 bugfix record uses invalid counter: %+v", record)
		}
		if !phase17StringInSlice(record.TriageResult, expected.TriageResults) || !phase17StringInSlice(record.FixMode, expected.FixModes) {
			t.Fatalf("Phase 18 bugfix record uses unknown triage or fix mode: %+v", record)
		}
		if record.KnownBugStatus != expected.TerminalKnownBugStatus {
			t.Fatalf("Phase 18 bugfix records must close known bugs to %s: %+v", expected.TerminalKnownBugStatus, record)
		}
		if record.FixMode == "spec-first-fix" {
			specFirst = true
			if record.IssueRef != phase18MarkdownRef("records/bugfix.jsonl") || record.ImplementationFixRef != phase18MarkdownRef("main_test.go") {
				t.Fatalf("Phase 18 spec-first bugfix record must link issue and implementation fix refs: %+v", record)
			}
			for label, ref := range map[string]string{
				"spec_general_update_ref":   record.SpecGeneralUpdateRef,
				"responsibility_source_ref": record.ResponsibilitySourceRef,
				"implementation_fix_ref":    record.ImplementationFixRef,
				"revalidation_ref":          record.RevalidationRef,
				"evidence_ref":              record.EvidenceRef,
			} {
				phase18RequireSafeReference(t, "bugfix "+label, ref)
			}
		}
		if record.IssueClassification == "no_issue" {
			noIssue = true
			if record.IssueID != "not_applicable" || record.TriageResult != "no_issue_detected" || record.FixMode != "not-required-revalidate" || record.IssueRef != "not_applicable" || record.ImplementationFixRef != "not_applicable" {
				t.Fatalf("Phase 18 no_issue bugfix record mismatch: %+v", record)
			}
		}
	}
	if !specFirst || !noIssue {
		t.Fatalf("Phase 18 bugfix records must include spec-first-fix and no_issue closure records")
	}
}

func phase18RequireSecurityRecords(t *testing.T, records []phase18EvidenceRecord, inputs phase18InputSet) {
	t.Helper()

	secretCoverage := map[string]bool{}
	operationCoverage := map[string]bool{}
	for _, record := range records {
		if record.EvidenceRef != phase18MarkdownRef("records/security.jsonl") {
			t.Fatalf("Phase 18 security record must use fixed linked evidence_ref: %+v", record)
		}
		if record.SecretClass != "not_applicable" {
			if record.SourceRef != phase18MarkdownRef("input/security_boundary.json") || record.CounterKey != "phase18_secret_leak_open_count" {
				t.Fatalf("Phase 18 secret coverage record must use security boundary source and secret counter: %+v", record)
			}
			if !phase17StringInSlice(record.SecretClass, inputs.Security.SecretClasses) {
				t.Fatalf("Phase 18 security record unknown secret_class: %+v", record)
			}
			if *record.ForbiddenValueDetected || record.MetadataPolicyResult != "opaque-only" || record.StorageResult != "not_stored" || (record.MaskResult != "masked" && record.MaskResult != "not_present") || record.Result != "passed" {
				t.Fatalf("Phase 18 secret coverage record must prove opaque non-storage and masking: %+v", record)
			}
			secretCoverage[record.SecretClass] = true
		}
		if record.ForbiddenProviderOperationToken != "not_applicable" {
			if record.SourceRef != phase18MarkdownRef("input/operation_scope.json") || record.CounterKey != "phase18_destructive_operation_open_count" {
				t.Fatalf("Phase 18 destructive-operation boundary record must use operation scope source and destructive counter: %+v", record)
			}
			if !phase17StringInSlice(record.ForbiddenProviderOperationToken, inputs.Scope.ForbiddenProviderOperations) {
				t.Fatalf("Phase 18 security record unknown forbidden operation: %+v", record)
			}
			if *record.ForbiddenProviderOperationDetected || record.BoundaryResult != "blocked" || record.DestructiveOperationResult != "blocked" || record.Result != "passed" {
				t.Fatalf("Phase 18 destructive operation must be blocked, not executed: %+v", record)
			}
			operationCoverage[record.ForbiddenProviderOperationToken] = true
		}
	}
	for _, secretClass := range inputs.Security.SecretClasses {
		if !secretCoverage[secretClass] {
			t.Fatalf("Phase 18 security records missing secret class coverage %s", secretClass)
		}
	}
	for _, operation := range inputs.Scope.ForbiddenProviderOperations {
		if !operationCoverage[operation] {
			t.Fatalf("Phase 18 security records missing forbidden provider operation coverage %s", operation)
		}
	}
}

func phase18RequireDocumentDriftRecords(t *testing.T, records []phase18EvidenceRecord) {
	t.Helper()

	targetCoverage := map[string]bool{}
	var futurePlan bool
	var stateTransitionGate bool
	futurePlanSourceRef := phase18MarkdownRef("docs/details/production-validation.md#phase-18-vps-pull-bootstrap-contract")
	futurePlanEvidenceRef := phase18MarkdownRef("docs/ROADMAP.md#roadmap-initial-phase-plan")
	stateTransitionGateRef := phase18MarkdownRef("docs/DETAIL_INDEX.md#phase-18-state-transition-gate")
	for _, record := range records {
		if record.ProviderTarget == "xserver-vps-future" {
			if record.SessionID != "not_applicable" || record.ValidationMode != "future_plan" || record.Result != "future_plan" || record.DriftResult != "future_plan_confirmed" || record.ClosureResult != "not_applicable" || record.CounterKey != "phase18_document_drift_open_count" || record.SourceRef != futurePlanSourceRef || record.EvidenceRef != futurePlanEvidenceRef || record.DocumentRef != futurePlanSourceRef || record.AnchorRef != futurePlanSourceRef || record.FixturePathRef != "not_applicable" || record.WorkflowPathRef != "not_applicable" || record.ExpectedRef != "not_applicable" || record.ActualRef != "not_applicable" || record.CorrectiveActionRef != "not_applicable" {
				t.Fatalf("Phase 18 xserver future plan classification mismatch: %+v", record)
			}
			futurePlan = true
			continue
		}
		if !phase17StringInSlice(record.DriftTarget, phase18DocumentDriftTargets) || record.DriftResult != "no_drift" || record.ClosureResult != "closed" || record.CounterKey != "phase18_document_drift_open_count" {
			t.Fatalf("Phase 18 document drift record mismatch: %+v", record)
		}
		targetCoverage[record.DriftTarget] = true
		if record.RecordID == "phase18.record.document_drift.state-transition-gate" {
			if record.DriftTarget != "anchor" || record.SourceRef != stateTransitionGateRef || record.DocumentRef != stateTransitionGateRef || record.AnchorRef != stateTransitionGateRef || record.ExpectedRef != stateTransitionGateRef || record.ActualRef != stateTransitionGateRef {
				t.Fatalf("Phase 18 state-transition gate document drift record mismatch: %+v", record)
			}
			stateTransitionGate = true
		}
	}
	for _, target := range phase18DocumentDriftTargets {
		if !targetCoverage[target] {
			t.Fatalf("Phase 18 document drift records missing target %s", target)
		}
	}
	if !stateTransitionGate {
		t.Fatalf("Phase 18 document drift records must cover state-transition gate anchor")
	}
	if !futurePlan {
		t.Fatalf("Phase 18 document drift records must keep xserver-vps-future as non-blocking future_plan")
	}
}

func phase18RequireDerivedCounters(t *testing.T, counters map[string]int, records map[string][]phase18EvidenceRecord) {
	t.Helper()

	derived := phase18DeriveCounters(records)
	for _, key := range phase18Counters {
		if counters[key] != derived[key] {
			t.Fatalf("Phase 18 counter %s must be derived from records: want %d got %d", key, derived[key], counters[key])
		}
	}
}

func phase18DeriveCounters(records map[string][]phase18EvidenceRecord) map[string]int {
	derived := map[string]int{}
	for _, key := range phase18Counters {
		derived[key] = 0
	}
	for _, recordSet := range records {
		for _, record := range recordSet {
			if record.Result == "failed" || record.Result == "rejected" {
				derived[record.CounterKey] = 1
			}
		}
	}
	for _, key := range phase18Counters {
		if key != "final_open_item_count" {
			derived["final_open_item_count"] += derived[key]
		}
	}
	return derived
}

func phase18RequireWorkflow(t *testing.T, path string) {
	t.Helper()

	workflow := phase11MustReadText(t, path)
	if !strings.Contains(workflow, "name: phase18-vps-pull-bootstrap") || !strings.Contains(workflow, "permissions:\n  contents: read") {
		t.Fatalf("%s must use the Phase 18 workflow name and contents:read permissions", path)
	}
	if !strings.Contains(workflow, "\n  phase18-vps-pull-bootstrap:") {
		t.Fatalf("%s must define primary Phase 18 job", path)
	}
	if !strings.Contains(workflow, "timeout-minutes: 30") {
		t.Fatalf("%s must set timeout-minutes <= 30 for Phase 18 job", path)
	}
	if !strings.Contains(workflow, "go test ./... -run TestPhase18VPSPullBootstrapEvidence -count=1") {
		t.Fatalf("%s must run the Phase 18 evidence checker directly", path)
	}
	phase13RequireNodeFreeWorkflow(t, workflow)
}

func phase18RequireBootstrapScriptArtifact(t *testing.T, path string, scope phase18OperationScopeFile) {
	t.Helper()

	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("missing-file: %s must exist as Phase 18 VPS Pull Bootstrap shell artifact: %v", path, err)
	}
	if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("schema-mismatch: %s must be a regular non-symlink shell artifact", path)
	}
	if info.Mode().Perm()&0111 == 0 {
		t.Fatalf("schema-mismatch: %s must be executable", path)
	}
	script := phase11MustReadText(t, path)
	if !strings.HasPrefix(script, "#!/bin/sh\n") {
		t.Fatalf("schema-mismatch: %s must start with #!/bin/sh", path)
	}
	if output, err := exec.Command("sh", "-n", path).CombinedOutput(); err != nil {
		t.Fatalf("schema-mismatch: %s must pass /bin/sh syntax validation: %v output=%s", path, err, output)
	}
	for _, forbidden := range []string{
		"#!/usr/bin/env bash",
		"#!/bin/bash",
		"[[",
		"function ",
		"pipefail",
		"git clone",
		"gh ",
		"GITHUB_TOKEN",
		"CONOHA",
		"self-hosted",
		"ssh-key",
		"vps-create",
		"vps-delete",
		"plan-change",
	} {
		if strings.Contains(script, forbidden) {
			t.Fatalf("forbidden-side-effect: %s must not contain forbidden bootstrap token %q", path, forbidden)
		}
	}
	for _, token := range []string{
		scope.BootstrapScriptArtifact,
		scope.BootstrapScriptRuntime,
		scope.ActiveSourceChannel,
		scope.MinimumTopologyRole,
		scope.BootstrapSourceIDFormat,
		scope.BootstrapDigestAlgorithm,
		scope.BootstrapInstallDirectory,
		scope.BootstrapBinDirectory,
		scope.BootstrapStateDirectory,
		scope.BootstrapServiceUser,
		scope.BootstrapSecretInputPolicy,
		scope.BootstrapStateModel,
		scope.BootstrapCLIContract,
		scope.BootstrapStdoutContract,
		scope.BootstrapStderrContract,
		scope.BootstrapExitCodeContract,
		"single-node",
		"stable-release",
		"absolute-https-url",
		"absolute-file-path",
		"/etc/systemd/system/adlaire-ci.service",
		"/etc/systemd/system/adlaire-ci.timer",
		"/etc/systemd/system/adlaire-ci-api.service",
		".admin_credentials",
		"phase18-materialize-admin",
		"--init-credentials",
		"/api/login",
		"/api/change-password",
		"/api/schedule/cooldown",
		"--token-file",
		"phase18-source",
		"phase18-site",
		"run_update_rollback_drill",
		"no_issue_detected",
	} {
		if !strings.Contains(script, token) {
			t.Fatalf("document-drift: %s must contain Phase 18 bootstrap contract token %q", path, token)
		}
	}
	for _, token := range scope.BootstrapInvocationTokens {
		if !strings.Contains(script, token) {
			t.Fatalf("document-drift: %s must contain fixed invocation token %q", path, token)
		}
	}
	for _, command := range append(scope.BootstrapRequiredCommands, scope.BootstrapHTTPClientOrder...) {
		if !strings.Contains(script, command) {
			t.Fatalf("document-drift: %s must contain required command/client token %q", path, command)
		}
	}
	for _, failureCode := range scope.BootstrapFailureCodes {
		if !strings.Contains(script, failureCode) {
			t.Fatalf("document-drift: %s must contain failure code %q", path, failureCode)
		}
	}
	for _, binary := range phase18BootstrapInstalledBinaries {
		if !strings.Contains(script, binary) {
			t.Fatalf("document-drift: %s must install standard binary %q", path, binary)
		}
	}
}

func phase18RequireDocumentState(t *testing.T) {
	t.Helper()

	roadmap := phase11MustReadText(t, "docs/ROADMAP.md")
	spec := phase11MustReadText(t, "docs/SPEC.md")
	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	detailIndex := phase11MustReadText(t, "docs/DETAIL_INDEX.md")
	productionValidation := phase11MustReadText(t, "docs/details/production-validation.md")
	fixture := phase11MustReadText(t, "docs/details/fixture.md")
	phase18Row := "| Phase 18 | VPS Pull Bootstrap。"
	phase18InProgressStatus := "| 実装中・bootstrap 検証未完了 | Phase 17 |"
	phase18ImplementedStatus := "| 実装済み | Phase 17 |"
	if !strings.Contains(roadmap, phase18Row) || !strings.Contains(roadmap, phase18ImplementedStatus) {
		t.Fatalf("docs/ROADMAP.md must mark Phase 18 implemented after repository and live VPS completion evidence close the state-transition gate")
	}
	if !strings.Contains(roadmap, "現在の active Phase は設定しない。") {
		t.Fatalf("docs/ROADMAP.md must leave active Phase unset after Phase 18 completion while Phase 19 awaits implementation approval")
	}
	if strings.Contains(roadmap, phase18InProgressStatus) || strings.Contains(roadmap, "現在の active Phase は Phase 18") || strings.Contains(roadmap, "Phase 18 は `仕様化済み・未実装`") {
		t.Fatalf("docs/ROADMAP.md must not retain a stale Phase 18 active or incomplete state after the state-transition gate closes")
	}
	if err := plannedPhaseMaturityBoundary(roadmap); err != nil {
		t.Fatal(err)
	}
	for _, token := range []string{
		"Phase 11 から Phase 18 の完了証跡条件、closure counter、record schema、required check は ROADMAP 本文で再掲せず",
		"Phase 18 の `実装済み` 遷移可否は本 ROADMAP で本文定義しない",
		"状態遷移判断は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 18 状態遷移 gate](DETAIL_INDEX.md#phase-18-state-transition-gate)",
	} {
		if !strings.Contains(roadmap, token) {
			t.Fatalf("docs/ROADMAP.md must keep Phase 18 state transition as canonical-reference-only text: %q", token)
		}
	}
	for _, forbidden := range []string{
		"Phase 11 は [`docs/details/fixture.md` fixture 証跡責務 Phase 11 完了 closure](details/fixture.md#phase-11-quality-gate-closure) に記録された closure record set の `final_open_item_count=0`",
		"Phase 12 は [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) に記録された closure counter の `final_open_item_count=0`",
		"Phase 13 は [`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence) に記録された全 closure counter の完了値を満たす",
		"Phase 14 は [`docs/details/fixture.md` fixture 証跡責務 Phase 14 Obsidian Vault 連携証跡](details/fixture.md#phase-14-obsidian-vault-integration-evidence) に記録された全 closure counter の完了値を満たす",
		"Phase 15 は [`docs/details/fixture.md` fixture 証跡責務 Phase 15 Obsidian local vault 同期証跡](details/fixture.md#phase-15-obsidian-local-sync-evidence) に記録された全 closure counter の完了値を満たす",
		"Phase 16 は [`docs/details/fixture.md` fixture 証跡責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡](details/fixture.md#phase-16-quality-evidence-closure-evidence) に記録された全 closure counter の完了値を満たす",
		"Phase 17 は [`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](details/fixture.md#phase-17-production-validation-evidence) に記録された全 closure counter の完了値を満たす",
		"Phase 18 は [`docs/details/fixture.md` fixture 証跡責務 Phase 18 VPS Pull Bootstrap](details/fixture.md#phase-18-vps-pull-bootstrap-evidence) に記録された全 closure counter の完了値",
		"Phase 18 を `実装済み` へ遷移する条件は、本 ROADMAP で再定義せず",
		"Phase 18 を `実装済み` へ遷移する場合は、[`docs/details/fixture.md` fixture 証跡責務 Phase 18 VPS Pull Bootstrap](details/fixture.md#phase-18-vps-pull-bootstrap-evidence) の checker acceptance fixture 完了に加え",
		"`live VPS completion summary` を同一 PR 本文へ secret-safe に記録しなければならない",
		"repository 内 fixture、GitHub Actions success、Docker success、または closure counter 完了値だけを、Phase 18 の状態遷移根拠として扱わない",
	} {
		if strings.Contains(roadmap, forbidden) {
			t.Fatalf("docs/ROADMAP.md must not replay closure counter detail owned by fixture/detail specs: %q", forbidden)
		}
	}
	for _, token := range []string{
		"`live VPS completion summary` の接続条件",
		"Phase 18 の完了は、[`docs/details/production-validation.md` 本番検証詳細本文責務 Phase 18 VPS Pull Bootstrap 契約]",
		"repository 内 checker acceptance fixture、GitHub Actions success、Docker success、または closure counter 完了値だけでは認めない",
		"`live VPS completion summary` の summary item、Phase 18 closure counter の key、完了値、record 種別、required check 名を再掲しない",
		"GitHub Actions は repository 側 checker required check としてだけ扱い",
		"VPS bootstrap の実行面、secret delivery、source 取得、deploy 制御、または live VPS completion summary の代替として扱ってはならない",
	} {
		if !strings.Contains(spec, token) {
			t.Fatalf("docs/SPEC.md must keep Phase 18 top-level completion policy connected to live VPS summary token %q", token)
		}
	}
	for _, feature := range []string{
		"Phase 18 VPS Pull Bootstrap / Phase 17 handover / bootstrap script artifact / topology role / source channel / bootstrap session / bootstrap artifact health / bootstrap runtime flow gate",
		"Phase 18 pull runner foundation / bootstrap update / rollback drill / bootstrap 中 issue / 仕様先行バグ修正 / 再検証 / known bug 0 gate",
		"Phase 18 secret / destructive operation / opaque metadata / document drift gate",
	} {
		if strings.Contains(feature, "secret / destructive operation / opaque metadata / document drift") {
			if !strings.Contains(roadmap, "| 実装済み | 本番検証 | "+feature+" |") {
				t.Fatalf("docs/ROADMAP.md must mark the Phase 18 boundary feature implemented in the 本番検証 category: %s", feature)
			}
		} else if !strings.Contains(roadmap, "| 実装済み | VPS Pull Bootstrap | "+feature+" |") {
			t.Fatalf("docs/ROADMAP.md must mark the Phase 18 bootstrap feature implemented in the VPS Pull Bootstrap category: %s", feature)
		}
		if strings.Contains(roadmap, "| 実装済み | CI/CD Bootstrap | "+feature+" |") ||
			strings.Contains(roadmap, "| 実装済み | 試験本番運用 | "+feature+" |") ||
			strings.Contains(roadmap, "| 実装中・bootstrap 検証未完了 | CI/CD Bootstrap | "+feature+" |") ||
			strings.Contains(roadmap, "| 実装中・bootstrap 検証未完了 | VPS Pull Bootstrap | "+feature+" |") ||
			strings.Contains(roadmap, "| 実装中・bootstrap 検証未完了 | 本番検証 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | CI/CD Bootstrap | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 試験本番運用 | "+feature+" |") ||
			strings.Contains(roadmap, "| 仕様化済み・未実装 | 本番検証 | "+feature+" |") {
			t.Fatalf("docs/ROADMAP.md must not classify the completed Phase 18 feature under a stale or incorrect state/category: %s", feature)
		}
	}
	for _, token := range []string{
		"`live VPS completion summary`",
		"`evidence_origin` | `live-vps-bootstrap`",
		"`execution_environment` | `vps`",
		"`provider_target` | `conoha-vps-primary`",
		"`vps_label` | `試験本番VPS`",
		"用途 label `試験本番VPS` を持つ VPS の opaque environment identity",
		"`validation_mode` | `vps-pull-bootstrap`",
		"`bootstrap_session_id` | `phase18.session.<slug>`",
		"`source_channel` | Phase 18 完了判定では `integration-head` 固定",
		"将来 CD 基盤の受け口とは、Phase 19 以降が `stable-release` channel、deploy、rollback、release selection、promotion を実装する際に使用する source channel、topology role、install directory、bin directory、state directory、service unit、log / state boundary を VPS 内に固定すること",
		"Phase 18 では deploy decision、production release selection、CD pipeline execution、customer environment promotion を実装対象、完了条件、または bootstrap 成功条件として扱ってはならない",
		"`topology_role=ci-cd` は VPS 内の配置 role 名",
		"承認済み deploy target、正式本番 deploy、customer environment deploy を成功証跡として扱わない",
		"`.github/workflows/phase18-vps-pull-bootstrap.yml` は repository 側 checker required check の実行経路",
		"VPS bootstrap 実行、secret delivery、source 取得、deploy 制御、または `live VPS completion summary` の代替ではない",
		"`fixture_gate_result`",
		"`checker-acceptance-fixture` を `live-vps-bootstrap` の代替にしない",
		"`closure_counter_result`",
		"`required_check_result`",
		"`final_open_item_result` | `0` 固定",
		"raw stdout",
		"raw stderr",
		"secret-safe summary を作成できない実測結果",
	} {
		if !strings.Contains(productionValidation, token) {
			t.Fatalf("docs/details/production-validation.md must define Phase 18 live VPS completion summary token %q", token)
		}
	}
	if strings.Contains(productionValidation, "deploy simulation または承認済み deploy target") || strings.Contains(fixture, "deploy simulation または承認済み deploy target") {
		t.Fatalf("Phase 18 must not allow approved deploy target as an alternative to deploy simulation evidence")
	}
	if !strings.Contains(roadmap, "Phase 18 は、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 18 状態遷移 gate](DETAIL_INDEX.md#phase-18-state-transition-gate) の完了証跡が同一実装 Pull Request に接続されたため `実装済み` とする。") {
		t.Fatalf("docs/ROADMAP.md must connect the Phase 18 implemented state to the canonical state-transition gate without replaying live evidence")
	}
	ambiguousPhase18VPSLabel := "試験本番運用・" + "開発検証兼用 VPS"
	if strings.Contains(roadmap, ambiguousPhase18VPSLabel) || strings.Contains(productionValidation, ambiguousPhase18VPSLabel) {
		t.Fatalf("Phase 18 must use the 試験本番VPS purpose label instead of an ambiguous trial-production/development VPS label")
	}
	for _, token := range []string{
		"bootstrap 完了判定または `実装済み` 遷移を扱う PR は、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 18 状態遷移 gate](#phase-18-state-transition-gate) に従う",
		"<a id=\"phase-18-state-transition-gate\"></a>",
		"repository 内 checker acceptance fixture の成功だけでは完了しない",
		"`live VPS completion summary` と [`docs/details/fixture.md` fixture 証跡責務 Phase 18 VPS Pull Bootstrap]",
		"summary item、counter key、required check 名、record schema を再掲しない",
	} {
		if !strings.Contains(detailIndex, token) {
			t.Fatalf("docs/DETAIL_INDEX.md must route Phase 18 completion without redefinition token %q", token)
		}
	}
	for _, token := range []string{
		"Phase 18 checker acceptance fixture は repository 内の仕様・schema・record・counter・required check・document drift の受入証跡",
		"bootstrap script artifact の静的契約",
		"`live VPS completion summary` を同一 PR 本文へ接続する",
		"fixture 証跡責務では `live VPS completion summary` の item を再掲しない",
		"Phase 18 Go test checker acceptance 契約",
		"`TestPhase18VPSPullBootstrapEvidence` は standalone CLI ではない",
		"custom stdout / stderr JSON schema",
		"`expected/negative_boundaries.json.expected_diagnostic_code` は、checker が検出すべき failure classification の期待値",
	} {
		if !strings.Contains(fixture, token) {
			t.Fatalf("docs/details/fixture.md must keep Phase 18 fixture/live summary boundary token %q", token)
		}
	}
	for _, forbidden := range []string{
		"Phase 18 checker " + "process 入出力契約",
		"checker " + "process 入出力契約不一致",
	} {
		if strings.Contains(fixture, forbidden) {
			t.Fatalf("docs/details/fixture.md must not describe Phase 18 as a standalone process checker: %q", forbidden)
		}
	}
	expectedPaths := []string{"testdata/phase18/vps-pull-bootstrap/"}
	for _, rel := range phase18FixtureFiles {
		expectedPaths = append(expectedPaths, "testdata/phase18/vps-pull-bootstrap/"+rel)
	}
	expectedPaths = append(expectedPaths, "adlaire-ci-vps-pull-bootstrap.sh")
	expectedPaths = append(expectedPaths, ".github/workflows/phase18-vps-pull-bootstrap.yml")
	for _, rel := range expectedPaths {
		if !strings.Contains(documentIndex, "| `"+rel+"` |") && !strings.Contains(documentIndex, "| [`"+rel+"`]") {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index Phase 18 path %s", rel)
		}
	}
	for _, token := range []string{
		"docs/DETAIL_INDEX.md#phase-18-state-transition-gate",
		"Phase 18 状態遷移 gate",
	} {
		if !strings.Contains(documentIndex, token) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index Phase 18 state transition gate token %q", token)
		}
	}
	phase18IndexStart := strings.Index(documentIndex, "<a id=\"phase-18-target-paths\"></a>")
	if phase18IndexStart < 0 {
		t.Fatalf("docs/DOCUMENT_INDEX.md must define phase-18-target-paths anchor")
	}
	phase18IndexSection := documentIndex[phase18IndexStart:]
	if strings.Contains(phase18IndexSection, "未作成") {
		t.Fatalf("docs/DOCUMENT_INDEX.md must not keep Phase 18 paths as 未作成")
	}
	for _, forbidden := range []string{
		"現在状態と完了可否は " + "[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する",
	} {
		if strings.Contains(documentIndex, forbidden) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must not assign current state and completion possibility to ROADMAP together: %q", forbidden)
		}
	}
	for _, token := range []string{
		"Phase 13 実装整合・品質改善証跡",
		"Phase 14 Obsidian Vault 連携証跡",
		"Phase 15 Obsidian local vault 同期証跡",
		"Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡",
		"Phase 17 ConoHa VPS 試験本番運用証跡",
		"Phase 18 状態遷移 gate",
		"本節は path 所在だけを正本とする",
	} {
		if !strings.Contains(documentIndex, token) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must keep responsibility-split Phase 13-18 index token %q", token)
		}
	}
}

func phase18RequireSafeReference(t *testing.T, label string, value string) {
	t.Helper()

	if strings.TrimSpace(value) == "" {
		t.Fatalf("%s must not be empty", label)
	}
	if value == "not_applicable" || value == "not_required_for_non_open_sample" {
		return
	}
	if filepath.IsAbs(value) || strings.Contains(value, "../") || strings.Contains(value, `..\`) {
		t.Fatalf("%s must not use absolute or parent path reference: %s", label, value)
	}
	if phase18ContainsForbiddenEvidenceLiteral(value) || phase17LooksLikeIPv4(value) {
		t.Fatalf("%s contains forbidden literal: %s", label, value)
	}
}

func phase18RequireRecordReferences(t *testing.T, rel string, line int, record phase18EvidenceRecord) {
	t.Helper()

	references := map[string]string{
		"source_ref":                record.SourceRef,
		"evidence_ref":              record.EvidenceRef,
		"operator_approval_ref":     record.OperatorApprovalRef,
		"owner_ref":                 record.OwnerRef,
		"endpoint_or_command_ref":   record.EndpointOrCommandRef,
		"operation_ref":             record.OperationRef,
		"issue_ref":                 record.IssueRef,
		"spec_general_update_ref":   record.SpecGeneralUpdateRef,
		"responsibility_source_ref": record.ResponsibilitySourceRef,
		"implementation_fix_ref":    record.ImplementationFixRef,
		"revalidation_ref":          record.RevalidationRef,
		"document_ref":              record.DocumentRef,
		"anchor_ref":                record.AnchorRef,
		"fixture_path_ref":          record.FixturePathRef,
		"workflow_path_ref":         record.WorkflowPathRef,
		"expected_ref":              record.ExpectedRef,
		"actual_ref":                record.ActualRef,
		"corrective_action_ref":     record.CorrectiveActionRef,
	}
	for field, value := range references {
		if value == "" {
			continue
		}
		label := fmt.Sprintf("%s line %d %s", rel, line, field)
		phase18RequireSafeReference(t, label, value)
		phase18RequireMarkdownReferenceForm(t, label, value)
		phase18RequireResolvableReference(t, label, value)
		phase18RequireCanonicalReferenceLabel(t, label, value)
	}
}

func phase18RequireMarkdownReferenceForm(t *testing.T, label string, value string) {
	t.Helper()

	if value == "not_applicable" || value == "not_required_for_non_open_sample" {
		return
	}
	if !phase18IsMarkdownReference(value) {
		t.Fatalf("%s must use a responsibility-named Markdown link: %s", label, value)
	}
}

func phase18RequireResolvableReference(t *testing.T, label string, value string) {
	t.Helper()

	if value == "not_applicable" || value == "not_required_for_non_open_sample" {
		return
	}
	rawTarget := phase18ReferenceTarget(value)
	if rawTarget == "" {
		t.Fatalf("%s must include a non-empty reference target: %s", label, value)
	}
	targetParts := strings.SplitN(rawTarget, "#", 2)
	target := targetParts[0]
	if target == "" {
		t.Fatalf("%s must not use an anchor-only reference: %s", label, value)
	}
	candidates := []string{target}
	if phase18LooksFixtureLocalReference(target) {
		candidates = append(candidates, filepath.Join("testdata/phase18/vps-pull-bootstrap", target))
	}
	resolved := ""
	for _, candidate := range candidates {
		if _, err := os.Lstat(candidate); err == nil {
			resolved = candidate
			break
		}
	}
	if resolved == "" {
		t.Fatalf("%s must resolve to an existing repository or Phase 18 fixture path: %s", label, value)
	}
	if len(targetParts) == 2 {
		phase18RequireExplicitAnchor(t, label, value, resolved, targetParts[1])
	}
}

func phase18RequireCanonicalReferenceLabel(t *testing.T, label string, value string) {
	t.Helper()

	if value == "not_applicable" || value == "not_required_for_non_open_sample" {
		return
	}
	target := phase18ReferenceTarget(value)
	if _, ok := phase18ReferenceLabels[target]; !ok {
		return
	}
	expected := phase18MarkdownRef(target)
	if value != expected {
		t.Fatalf("%s must use canonical responsibility reference label: want %s got %s", label, expected, value)
	}
}

func phase18RequireExplicitAnchor(t *testing.T, label string, value string, resolved string, anchor string) {
	t.Helper()

	if strings.TrimSpace(anchor) == "" {
		t.Fatalf("%s must not use an empty anchor reference: %s", label, value)
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		t.Fatalf("%s must read anchor target %s: %v", label, resolved, err)
	}
	body := string(data)
	doubleQuotedID := `id="` + anchor + `"`
	singleQuotedID := `id='` + anchor + `'`
	if strings.Contains(body, doubleQuotedID) || strings.Contains(body, singleQuotedID) {
		return
	}
	t.Fatalf("%s must resolve to an explicit anchor id %q in %s: %s", label, anchor, resolved, value)
}

func phase18ReferenceTarget(value string) string {
	if start := strings.LastIndex(value, "]("); phase18IsMarkdownReference(value) {
		return strings.TrimSpace(value[start+2 : len(value)-1])
	}
	return strings.TrimSpace(value)
}

func phase18IsMarkdownReference(value string) bool {
	start := strings.LastIndex(value, "](")
	return strings.HasPrefix(value, "[") && start >= 0 && strings.HasSuffix(value, ")")
}

func phase18LooksFixtureLocalReference(target string) bool {
	if strings.HasPrefix(target, "docs/") ||
		strings.HasPrefix(target, ".github/") ||
		strings.HasPrefix(target, "testdata/") ||
		strings.HasPrefix(target, "components/") ||
		target == "main.go" ||
		target == "main_test.go" ||
		target == "README.md" ||
		target == "AGENTS.md" {
		return false
	}
	return true
}

func phase18ContainsForbiddenEvidenceLiteral(value string) bool {
	lower := strings.ToLower(value)
	for _, forbidden := range []string{
		"-----begin openssh private key-----",
		"-----begin rsa private key-----",
		"-----begin private key-----",
		"ghp_",
		"github_pat_",
		"xoxb-",
		"sk-proj-",
		"sk-svcacct-",
		"password=",
		"token=",
		"secret=",
		"http://",
		"https://",
	} {
		if strings.Contains(lower, forbidden) {
			return true
		}
	}
	for _, token := range strings.FieldsFunc(value, func(r rune) bool {
		return !(r == '.' || r == '-' || r == '_' || r == ':' || r == '/' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z')
	}) {
		if phase17LooksLikeIPv4(token) {
			return true
		}
	}
	return false
}

func phase18MarkdownRef(target string) string {
	label, ok := phase18ReferenceLabels[target]
	if !ok {
		panic("missing Phase 18 reference label for " + target)
	}
	displayTarget := strings.SplitN(target, "#", 2)[0]
	return "[" + displayTarget + " " + label + "](" + target + ")"
}

func phase13RequireOwnerPackages(t *testing.T) {
	t.Helper()

	for _, owner := range phase13GoOwners {
		dir := filepath.Join("components", owner)
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read owner package %s: %v", dir, err)
		}
		var files []string
		for _, entry := range entries {
			if entry.IsDir() {
				t.Fatalf("%s must not contain subdirectory %s", dir, entry.Name())
			}
			files = append(files, entry.Name())
			info, err := entry.Info()
			if err != nil {
				t.Fatalf("stat %s/%s: %v", dir, entry.Name(), err)
			}
			if info.Size() == 0 {
				t.Fatalf("%s/%s must not be empty", dir, entry.Name())
			}
		}
		sort.Strings(files)
		expected := []string{owner + ".go", owner + "_test.go", "execute.go", "model.go", "validate.go"}
		sort.Strings(expected)
		if !phase11StringSlicesEqual(files, expected) {
			t.Fatalf("%s must contain exactly the five owner files\nwant: %v\n got: %v", dir, expected, files)
		}
	}
}

func phase13RequireDocumentIndexRows(t *testing.T) {
	t.Helper()

	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, expected := range []string{
		"| [`testdata/phase13/implementation-alignment-quality/`](../testdata/phase13/implementation-alignment-quality/) | Phase 13 implementation alignment quality fixture root | 実在 |",
		"| [`.github/workflows/phase13-implementation-alignment-quality.yml`](../.github/workflows/phase13-implementation-alignment-quality.yml) | Phase 13 required check workflow | 実在 |",
		"| [`testdata/phase13/implementation-alignment-quality/manifest.json`](../testdata/phase13/implementation-alignment-quality/manifest.json) | Phase 13 evidence package manifest | 実在 |",
		"| [`testdata/phase13/implementation-alignment-quality/expected/counters.json`](../testdata/phase13/implementation-alignment-quality/expected/counters.json) | Phase 13 closure counter 期待値 | 実在 |",
		"| [`testdata/phase13/implementation-alignment-quality/records/closure.jsonl`](../testdata/phase13/implementation-alignment-quality/records/closure.jsonl) | Phase 13 closure record set | 実在 |",
	} {
		if !strings.Contains(documentIndex, expected) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must contain Phase 13実在 row: %s", expected)
		}
	}
}

func TestChangeHistoryArtifactBoundary(t *testing.T) {
	t.Parallel()
	requireChangeHistoryArtifactBoundary(t)
}

func TestChangeHistoryArtifactBoundaryControls(t *testing.T) {
	t.Parallel()

	for _, path := range []string{
		"CHANGELOG.md", "docs/cHaNgElOg.TXT", "archive/.history.JSON",
		"templates/CHANGE-LOG.template.yaml", "releases/RELEASE_NOTES.v1.txt",
		"metadata/CHANGES", "hidden/.private/HISTORY-2026.md", "docs/変更履歴.rst",
		"nested/CHANGE.LOG.xml", "notes/リリースノート.txt",
		"notes/RELEASENOTES.md", "metadata/CHANGE_HISTORY.json",
		"updates/更新_履歴.2026.txt", "releases/リリース履歴.md",
		"templates/.release notes.template.md",
		"docs/CHANGELOG*.md", "releases/RELEASE_NOTES?.txt",
		"templates/.HISTORY[0-9].md",
	} {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			fullPath := filepath.Join(root, filepath.FromSlash(path))
			if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(fullPath, nil, 0600); err != nil {
				t.Fatal(err)
			}
			got, err := changeHistoryArtifactPaths(root)
			if err != nil || !phase11StringSlicesEqual(got, []string{path}) {
				t.Fatalf("forbidden path must be detected: path=%s got=%v err=%v", path, got, err)
			}
		})
	}
	t.Run("directory-and-symlink", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "exports", "UPDATE_HISTORY"), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "notes.md"), nil, 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("notes.md", filepath.Join(root, "release-history.txt")); err != nil {
			t.Fatal(err)
		}
		got, err := changeHistoryArtifactPaths(root)
		want := []string{"exports/UPDATE_HISTORY", "release-history.txt"}
		if err != nil || !phase11StringSlicesEqual(got, want) {
			t.Fatalf("forbidden directory and symlink must be detected: got=%v err=%v", got, err)
		}
	})
	t.Run("functional-records-and-git", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		for _, path := range []string{".git/CHANGELOG", ".build_history", "state/config_changes.jsonl", "state/events.jsonl", "release/notes.md"} {
			fullPath := filepath.Join(root, filepath.FromSlash(path))
			if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(fullPath, nil, 0600); err != nil {
				t.Fatal(err)
			}
		}
		got, err := changeHistoryArtifactPaths(root)
		if err != nil || len(got) != 0 {
			t.Fatalf("functional records and Git metadata must pass: got=%v err=%v", got, err)
		}
	})
	t.Run("missing-root", func(t *testing.T) {
		t.Parallel()
		if _, err := changeHistoryArtifactPaths(filepath.Join(t.TempDir(), "missing")); err == nil {
			t.Fatal("scan failure must not pass")
		}
	})
	t.Run("non-directory-root", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		file := filepath.Join(root, "notes.md")
		if err := os.WriteFile(file, nil, 0600); err != nil {
			t.Fatal(err)
		}
		link := filepath.Join(root, "linked-root")
		if err := os.Symlink(root, link); err != nil {
			t.Fatal(err)
		}
		for _, path := range []string{file, link, link + string(os.PathSeparator)} {
			if _, err := changeHistoryArtifactPaths(path); err == nil || !strings.Contains(err.Error(), "scan root must be a directory") {
				t.Fatalf("non-directory scan root must fail: path=%s err=%v", path, err)
			}
		}
	})
	t.Run("artifact-declarations", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			name      string
			data      string
			wantError string
		}{
			{"invalid-json", `{`, "parse governance scope"},
			{"include-type", `{"include":true,"release_artifacts":["LICENSE"]}`, "parse governance scope"},
			{"release-artifacts-type", `{"include":["docs/**/*.md"],"release_artifacts":"LICENSE"}`, "parse governance scope"},
			{"missing-include", `{"release_artifacts":["LICENSE"]}`, "include must be a non-empty string array"},
			{"null-include", `{"include":null,"release_artifacts":["LICENSE"]}`, "include must be a non-empty string array"},
			{"empty-include", `{"include":[],"release_artifacts":["LICENSE"]}`, "include must be a non-empty string array"},
			{"missing-release-artifacts", `{"include":["docs/**/*.md"]}`, "release_artifacts must be a non-empty string array"},
			{"null-release-artifacts", `{"include":["docs/**/*.md"],"release_artifacts":null}`, "release_artifacts must be a non-empty string array"},
			{"empty-release-artifacts", `{"include":["docs/**/*.md"],"release_artifacts":[]}`, "release_artifacts must be a non-empty string array"},
			{"include-item-type", `{"include":[1],"release_artifacts":["LICENSE"]}`, "parse governance scope"},
			{"null-include-item", `{"include":[null],"release_artifacts":["LICENSE"]}`, "include must not contain empty entries"},
			{"blank-include-item", `{"include":["  "],"release_artifacts":["LICENSE"]}`, "include must not contain empty entries"},
			{"empty-release-artifacts-item", `{"include":["docs/**/*.md"],"release_artifacts":[""]}`, "release_artifacts must not contain empty entries"},
			{"forbidden-include", `{"include":["docs/**/*.md","CHANGELOG.md"],"release_artifacts":["LICENSE"]}`, "forbidden change history artifact declaration in include: CHANGELOG.md"},
			{"forbidden-release-artifacts", `{"include":["docs/**/*.md"],"release_artifacts":["LICENSE","CHANGELOG.md"]}`, "forbidden change history artifact declaration in release_artifacts: CHANGELOG.md"},
			{"nested-include", `{"include":["templates/HISTORY/**"],"release_artifacts":["LICENSE"]}`, "forbidden change history artifact declaration in include: templates/HISTORY/**"},
			{"hidden-release-artifacts", `{"include":["docs/**/*.md"],"release_artifacts":["hidden/.release_notes.json"]}`, "forbidden change history artifact declaration in release_artifacts: hidden/.release_notes.json"},
			{"duplicate-include", `{"include":["CHANGELOG.md"],"include":["docs/**/*.md"],"release_artifacts":["LICENSE"]}`, `duplicate key "include"`},
			{"duplicate-release-artifacts", `{"include":["docs/**/*.md"],"release_artifacts":["CHANGELOG.md"],"release_artifacts":["LICENSE"]}`, `duplicate key "release_artifacts"`},
			{"escaped-duplicate-include", `{"include":["CHANGELOG.md"],"incl\u0075de":["docs/**/*.md"],"release_artifacts":["LICENSE"]}`, `duplicate key "include"`},
			{"include-case-shadow", `{"include":["CHANGELOG.md"],"Include":["docs/**/*.md"],"release_artifacts":["LICENSE"]}`, "governance scope must use canonical key include"},
			{"release-artifacts-case-shadow", `{"include":["docs/**/*.md"],"release_artifacts":["CHANGELOG.md"],"Release_Artifacts":["LICENSE"]}`, "governance scope must use canonical key release_artifacts"},
			{"noncanonical-include", `{"INCLUDE":["docs/**/*.md"],"release_artifacts":["LICENSE"]}`, "governance scope must use canonical key include"},
			{"noncanonical-release-artifacts", `{"include":["docs/**/*.md"],"RELEASE_ARTIFACTS":["LICENSE"]}`, "governance scope must use canonical key release_artifacts"},
			{"metadata-fields", `{"phase":"Phase 13","owner_packages":["release"],"include":["docs/**/*.md"],"release_artifacts":["LICENSE"]}`, ""},
			{"glob-include", `{"include":["docs/CHANGELOG*.md"],"release_artifacts":["LICENSE"]}`, "forbidden change history artifact declaration in include: docs/CHANGELOG*.md"},
			{"glob-release-artifact", `{"include":["docs/**/*.md"],"release_artifacts":["releases/RELEASE_NOTES?.txt"]}`, "forbidden change history artifact declaration in release_artifacts: releases/RELEASE_NOTES?.txt"},
			{"glob-directory", `{"include":["templates/.HISTORY[0-9]/**"],"release_artifacts":["LICENSE"]}`, "forbidden change history artifact declaration in include: templates/.HISTORY[0-9]/**"},
			{"general-and-functional-patterns", `{"include":["docs/**/*.md","components/*/*.go",".build_history*","state/config_changes?.jsonl"],"release_artifacts":["LICENSE","release/notes[0-9].md"]}`, ""},
			{"functional-records", `{"include":["docs/**/*.md",".build_history","state/config_changes.jsonl"],"release_artifacts":["LICENSE","notes.md"]}`, ""},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				err := validateChangeHistoryArtifactScope([]byte(tc.data))
				if tc.wantError == "" {
					if err != nil {
						t.Fatalf("functional artifact declarations must pass: %v", err)
					}
				} else if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("invalid artifact declarations: want error containing %q, got %v", tc.wantError, err)
				}
			})
		}
	})
}

func requireChangeHistoryArtifactBoundary(t *testing.T) {
	t.Helper()
	paths, err := changeHistoryArtifactPaths(".")
	if err != nil {
		t.Fatalf("scan change history artifacts: %v", err)
	}
	if len(paths) != 0 {
		t.Fatalf("forbidden change history artifacts: %v", paths)
	}
	path := "testdata/phase13/implementation-alignment-quality/input/scope.json"
	if err := validateChangeHistoryArtifactScope([]byte(phase11MustReadText(t, path))); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
}

func validateChangeHistoryArtifactScope(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := phase16CheckJSONValueForDuplicateKeys(decoder); err != nil {
		return fmt.Errorf("parse governance scope keys: %w", err)
	}
	var scope map[string]json.RawMessage
	if err := json.Unmarshal(data, &scope); err != nil {
		return fmt.Errorf("parse governance scope: %w", err)
	}
	for _, name := range []string{"include", "release_artifacts"} {
		for key := range scope {
			if key != name && strings.EqualFold(key, name) {
				return fmt.Errorf("governance scope must use canonical key %s", name)
			}
		}
		var entries []string
		if raw, ok := scope[name]; ok {
			if err := json.Unmarshal(raw, &entries); err != nil {
				return fmt.Errorf("parse governance scope %s: %w", name, err)
			}
		}
		if len(entries) == 0 {
			return fmt.Errorf("%s must be a non-empty string array", name)
		}
		for _, entry := range entries {
			if strings.TrimSpace(entry) == "" {
				return fmt.Errorf("%s must not contain empty entries", name)
			}
			if changeHistoryArtifactPath(entry) {
				return fmt.Errorf("forbidden change history artifact declaration in %s: %s", name, entry)
			}
		}
	}
	return nil
}

func changeHistoryArtifactPaths(root string) ([]string, error) {
	root = filepath.Clean(root)
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			if !entry.IsDir() {
				return fmt.Errorf("scan root must be a directory: %s", root)
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == ".git" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if changeHistoryArtifactPath(entry.Name()) {
			paths = append(paths, filepath.ToSlash(rel))
			if entry.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	sort.Strings(paths)
	return paths, err
}

func changeHistoryArtifactPath(path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if index := strings.IndexAny(part, "*?["); index >= 0 {
			part = part[:index]
		}
		fields := strings.FieldsFunc(strings.ToLower(part), func(r rune) bool {
			return r == '.' || r == '-' || r == '_' || r == ' '
		})
		var prefix string
		for _, field := range fields {
			prefix += field
			switch prefix {
			case "changelog", "changes", "history", "releasenotes", "changehistory", "updatehistory", "releasehistory", "変更履歴", "更新履歴", "リリース履歴", "リリースノート":
				return true
			}
		}
	}
	return false
}

func phase13RequireGovernanceFiles(t *testing.T) {
	t.Helper()

	documentIndex := phase11MustReadText(t, "docs/DOCUMENT_INDEX.md")
	for _, file := range []string{"LICENSE", "SECURITY.md", "CONTRIBUTING.md", "CODEOWNERS"} {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Phase 13 governance file %s must exist: %v", file, err)
		}
		if strings.TrimSpace(string(data)) == "" {
			t.Fatalf("Phase 13 governance file %s must not be empty", file)
		}
		row := "| [`" + file + "`](../" + file + ") | release governance artifact | 実在 |"
		if !strings.Contains(documentIndex, row) {
			t.Fatalf("docs/DOCUMENT_INDEX.md must index governance file as 実在: %s", file)
		}
	}
}

func phase13RequireZeroExternalDependencies(t *testing.T) {
	t.Helper()

	goMod := phase11MustReadText(t, "go.mod")
	for _, forbidden := range []string{"\nrequire ", "\nreplace ", "golang.org/x/"} {
		if strings.Contains(goMod, forbidden) {
			t.Fatalf("go.mod must not contain external dependency marker %q", forbidden)
		}
	}
	for _, path := range []string{"docs/SPEC.md", "docs/details/fixture.md", ".github/workflows/phase13-implementation-alignment-quality.yml"} {
		body := phase11MustReadText(t, path)
		for _, forbidden := range []string{"staticcheck", "govulncheck", "gosec", "Argon2id", "argon2id"} {
			if strings.Contains(body, forbidden) {
				t.Fatalf("%s must not contain forbidden external dependency term %q", path, forbidden)
			}
		}
	}
	implementationFiles, err := filepath.Glob("components/*/*.go")
	if err != nil {
		t.Fatalf("glob implementation files: %v", err)
	}
	implementationFiles = append(implementationFiles, "admin/adlaire-ci-sdk.js", "admin/index.html")
	for _, path := range implementationFiles {
		body := phase11MustReadText(t, path)
		for _, forbidden := range []string{"golang.org/x/", "npm install", "cdn.jsdelivr", "unpkg.com"} {
			if strings.Contains(body, forbidden) {
				t.Fatalf("%s must not contain external dependency marker %q", path, forbidden)
			}
		}
	}
}

func phase13ReadCounterFile(t *testing.T, path string) map[string]int {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var counters map[string]int
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&counters); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return counters
}

func phase13ReadClosureRecords(t *testing.T, path string) []phase13ClosureRecord {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	text := string(data)
	if !strings.HasSuffix(text, "\n") {
		t.Fatalf("%s must end with LF", path)
	}
	var records []phase13ClosureRecord
	for index, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			t.Fatalf("%s line %d must not be empty", path, index+1)
		}
		var record phase13ClosureRecord
		decoder := json.NewDecoder(strings.NewReader(line))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&record); err != nil {
			t.Fatalf("decode %s line %d: %v", path, index+1, err)
		}
		if strings.TrimSpace(record.RecordID) == "" {
			t.Fatalf("%s line %d record_id must not be empty", path, index+1)
		}
		records = append(records, record)
	}
	return records
}

func phase13ReadJSONLines(t *testing.T, path string) []map[string]any {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	text := string(data)
	if !strings.HasSuffix(text, "\n") {
		t.Fatalf("%s must end with LF", path)
	}
	var records []map[string]any
	for index, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			t.Fatalf("%s line %d must not be empty", path, index+1)
		}
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("decode %s line %d: %v", path, index+1, err)
		}
		records = append(records, record)
	}
	return records
}

func phase13RequireNodeFreeWorkflow(t *testing.T, workflow string) {
	t.Helper()
	for _, definition := range developmentActionsDefinitions(t) {
		if strings.HasPrefix(workflow, "name: "+definition.Name+"\n") {
			if err := developmentActionsCheckWorkflow(definition, workflow); err != nil {
				t.Fatal(err)
			}
			return
		}
	}
	t.Fatal("actions-workflow-contract: unknown workflow")
}

func phase11RequireNoUnexpectedFixtureDirectories(t *testing.T) {
	t.Helper()

	var unexpected []string
	err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if phase11HasNumericCopySuffix(d.Name()) {
			unexpected = append(unexpected, filepath.ToSlash(path))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan fixture directories: %v", err)
	}
	formalRoots := []string{"testdata/admin/cli", "testdata/builder", "testdata/setup", "testdata/release"}
	for _, coverage := range phase11RootCoverageFixtures {
		formalRoots = append(formalRoots, strings.TrimSuffix(coverage.Root, "/"))
	}
	for _, formalRoot := range formalRoots {
		entries, err := os.ReadDir(formalRoot)
		if err != nil {
			t.Fatalf("read formal fixture root %s: %v", formalRoot, err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			manifestPath := filepath.Join(formalRoot, entry.Name(), "manifest.json")
			if _, err := os.Stat(manifestPath); err != nil {
				unexpected = append(unexpected, filepath.ToSlash(filepath.Join(formalRoot, entry.Name())))
			}
		}
	}
	sort.Strings(unexpected)
	if len(unexpected) != 0 {
		t.Fatalf("unexpected or unregistered fixture directories remain: %v", unexpected)
	}
}

func phase11HasNumericCopySuffix(name string) bool {
	index := strings.LastIndex(name, " ")
	if index < 0 || index == len(name)-1 {
		return false
	}
	for _, char := range name[index+1:] {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func phase11ReadManifest(t *testing.T, path string) phase11Manifest {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var manifest phase11Manifest
	if err := decoder.Decode(&manifest); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	if decoder.More() {
		t.Fatalf("%s must contain one JSON object", path)
	}
	return manifest
}

func phase11FindManifestFiles(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "manifest.json" {
			paths = append(paths, filepath.ToSlash(path))
		}
		return nil
	})
	sort.Strings(paths)
	return paths, err
}

func phase11FindMarkdownFiles() ([]string, error) {
	roots := []string{"AGENTS.md", "README.md", "docs"}
	var paths []string
	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			if strings.HasSuffix(root, ".md") {
				paths = append(paths, root)
			}
			continue
		}
		err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".md") {
				paths = append(paths, filepath.ToSlash(path))
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func phase11NameHasCategoryPrefix(name string) bool {
	for category := range phase11AllowedCategories {
		if strings.HasPrefix(name, category+"-") {
			return true
		}
	}
	return false
}

func phase11RequireMarkdownReference(t *testing.T, manifestPath string, reference string) {
	t.Helper()

	if strings.Contains(reference, "\\") || strings.Contains(reference, "../") || strings.HasPrefix(reference, "/") {
		t.Fatalf("%s reference %q must be a repository-relative Markdown link target", manifestPath, reference)
	}
	parts := strings.Split(reference, "#")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		t.Fatalf("%s reference %q must include exactly one file path and one fixed anchor", manifestPath, reference)
	}
	if !strings.HasPrefix(parts[0], "docs/") || !strings.HasSuffix(parts[0], ".md") {
		t.Fatalf("%s reference %q must point to a docs/*.md file", manifestPath, reference)
	}
	data := phase11MustReadText(t, parts[0])
	anchor := parts[1]
	if !strings.Contains(data, `id="`+anchor+`"`) && !strings.Contains(data, "id='"+anchor+"'") && !strings.Contains(data, "#"+anchor) {
		t.Fatalf("%s reference %q points to a missing anchor", manifestPath, reference)
	}
}

func phase11RequireManifestRecords(t *testing.T, manifestPath string, field string, records []phase11ManifestRecord) {
	t.Helper()

	seen := make(map[string]bool)
	for _, record := range records {
		if !phase11SafeRelativePath(record.Path) {
			t.Fatalf("%s %s path %q must be a safe relative path", manifestPath, field, record.Path)
		}
		if strings.TrimSpace(record.Reason) == "" {
			t.Fatalf("%s %s path %q must include a reason", manifestPath, field, record.Path)
		}
		if seen[record.Path] {
			t.Fatalf("%s %s repeats path %q", manifestPath, field, record.Path)
		}
		seen[record.Path] = true
	}
}

func phase11RequireSafeRelativePaths(t *testing.T, manifestPath string, field string, paths []string) {
	t.Helper()

	phase11RequireUnique(t, manifestPath, field, paths)
	for _, path := range paths {
		if !phase11SafeRelativePath(path) {
			t.Fatalf("%s %s path %q must be a safe relative path", manifestPath, field, path)
		}
	}
}

func phase11SafeRelativePath(path string) bool {
	if path == "" || strings.HasPrefix(path, "/") || strings.Contains(path, "\\") {
		return false
	}
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		return false
	}
	for _, segment := range strings.Split(path, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
		if strings.ContainsAny(segment, "\x00\r\n") {
			return false
		}
	}
	return true
}

func phase11RequireAssertionEvidence(t *testing.T, fixtureDir string, manifestPath string, assertion string) {
	t.Helper()

	assertionFiles := map[string]string{
		"effects":     "expected/effects.json",
		"request":     "expected/request.json",
		"response":    "expected/response.json",
		"secret-mask": "expected/security.json",
		"stderr":      "expected/stderr.txt",
		"stdout":      "expected/stdout.txt",
	}
	if rel, ok := assertionFiles[assertion]; ok {
		path := filepath.Join(fixtureDir, filepath.FromSlash(rel))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("%s assertion %q requires %s: %v", manifestPath, assertion, rel, err)
		}
		if info.IsDir() {
			t.Fatalf("%s assertion %q requires %s to be a file", manifestPath, assertion, rel)
		}
		if strings.HasSuffix(rel, ".json") {
			phase11RequireValidJSONFile(t, path)
		}
		return
	}

	assertionDirs := map[string]string{
		"logs":  "expected/logs",
		"state": "expected/state",
	}
	if rel, ok := assertionDirs[assertion]; ok {
		path := filepath.Join(fixtureDir, filepath.FromSlash(rel))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("%s assertion %q requires %s: %v", manifestPath, assertion, rel, err)
		}
		if !info.IsDir() {
			t.Fatalf("%s assertion %q requires %s to be a directory", manifestPath, assertion, rel)
		}
		return
	}
}

func phase11RequireValidJSONFile(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read JSON evidence %s: %v", path, err)
	}
	if !json.Valid(data) {
		t.Fatalf("JSON evidence %s is not valid JSON", path)
	}
}

func phase11RequireUnique(t *testing.T, context string, field string, values []string) {
	t.Helper()

	seen := make(map[string]bool)
	for _, value := range values {
		if value == "" {
			t.Fatalf("%s %s must not contain an empty value", context, field)
		}
		if seen[value] {
			t.Fatalf("%s %s repeats %q", context, field, value)
		}
		seen[value] = true
	}
}

func phase11RequireSortedUnique(t *testing.T, context string, field string, values []string) {
	t.Helper()

	phase11RequireUnique(t, context, field, values)
	if !sort.StringsAreSorted(values) {
		t.Fatalf("%s %s must be sorted: %v", context, field, values)
	}
}

func phase11StringSlicesEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type developmentActionsJob struct {
	ID      string `json:"id"`
	Timeout int    `json:"timeout"`
	Runtime string `json:"runtime"`
	Command string `json:"command"`
}

type developmentActionsWorkflow struct {
	Phase int                      `json:"phase"`
	Name  string                   `json:"name"`
	Jobs  []developmentActionsJob  `json:"jobs"`
	Setup *developmentActionsSetup `json:"-"`
}

type developmentActionsSetup struct {
	GoVersion    string `json:"go_version"`
	GoSHA256     string `json:"go_sha256"`
	DenoVersion  string `json:"deno_version"`
	DenoSHA256   string `json:"deno_sha256"`
	HelperSHA256 string `json:"helper_sha256"`
	Checkout     string `json:"checkout"`
}

func developmentActionsSetupDefinition(t *testing.T) *developmentActionsSetup {
	t.Helper()
	text := phase11MustReadText(t, "docs/details/fixture.md")
	_, section, ok := strings.Cut(text, `<a id="development-actions-node-free-contract"></a>`)
	if !ok {
		t.Fatal("actions-setup: missing canonical anchor")
	}
	_, section, ok = strings.Cut(section, "```json\n")
	if !ok {
		t.Fatal("actions-setup: missing definition")
	}
	data, _, ok := strings.Cut(section, "\n```")
	if !ok {
		t.Fatal("actions-setup: missing definition end")
	}
	var setup developmentActionsSetup
	if err := developmentActionsDecodeJSON([]byte(data), &setup, 2<<20, true); err != nil {
		t.Fatal(err)
	}
	for _, digest := range []string{setup.GoSHA256, setup.DenoSHA256, setup.HelperSHA256} {
		decoded, err := hex.DecodeString(digest)
		if err != nil || len(decoded) != 32 || strings.ToLower(digest) != digest {
			t.Fatal("actions-setup: invalid digest")
		}
	}
	if !regexp.MustCompile(`^1\.22\.[0-9]+$`).MatchString(setup.GoVersion) || !regexp.MustCompile(`^2\.[0-9]+\.[0-9]+$`).MatchString(setup.DenoVersion) || setup.Checkout == "" {
		t.Fatal("actions-setup: invalid definition")
	}
	return &setup
}

const developmentActionsGateCommand = `test "$ADLAIRE_ACTIONS_GATE" = required
test -n "$ADLAIRE_ACTIONS_NEEDS"
go test . -run '^TestDevelopmentActions(RequiredGate|WorkflowContract|GateControls|NodeFreeControls|BootstrapControls)$' -count=1 -v | tee "$RUNNER_TEMP/development-actions-gate.log"
for required_test in TestDevelopmentActionsRequiredGate TestDevelopmentActionsWorkflowContract TestDevelopmentActionsGateControls TestDevelopmentActionsNodeFreeControls TestDevelopmentActionsBootstrapControls; do
  grep -Fq -- "--- PASS: $required_test (" "$RUNNER_TEMP/development-actions-gate.log"
done
grep -Fq '"record_type":"development-actions-gate"' "$RUNNER_TEMP/development-actions-gate.log"`

func developmentActionsDefinitions(t *testing.T) []developmentActionsWorkflow {
	t.Helper()
	text := phase11MustReadText(t, "docs/details/fixture.md")
	_, section, ok := strings.Cut(text, `<a id="development-actions-required-gate-contract"></a>`)
	if !ok {
		t.Fatal("actions-definition: missing canonical fixture anchor")
	}
	_, section, ok = strings.Cut(section, "```json\n")
	if !ok {
		t.Fatal("actions-definition: missing JSON fixture")
	}
	data, _, ok := strings.Cut(section, "\n```\n")
	if !ok {
		t.Fatal("actions-definition: unterminated JSON fixture")
	}
	var definitions []developmentActionsWorkflow
	if err := developmentActionsDecodeJSON([]byte(data), &definitions, 2<<20, true); err != nil {
		t.Fatalf("actions-definition: %v", err)
	}
	names := []string{"phase12-quality-gate", "phase13-implementation-alignment-quality", "phase14-obsidian-vault-integration", "phase15-obsidian-local-sync", "phase16-quality-evidence-closure", "phase17-production-validation", "phase18-vps-pull-bootstrap"}
	counts := []int{7, 27, 5, 6, 9, 11, 1}
	if len(definitions) != len(names) {
		t.Fatal("actions-definition: workflow set mismatch")
	}
	for i, definition := range definitions {
		definitions[i].Setup = developmentActionsSetupDefinition(t)
		if definition.Phase != i+12 || definition.Name != names[i] || len(definition.Jobs) != counts[i] {
			t.Fatal("actions-definition: workflow identity or required job count mismatch")
		}
		seen := map[string]bool{}
		for _, job := range definition.Jobs {
			if !strings.HasPrefix(job.ID, fmt.Sprintf("phase%d-", definition.Phase)) || !developmentActionsIdentifier(job.ID) || seen[job.ID] || job.ID == definition.gateID() || job.Timeout < 1 || job.Timeout > 30 || (job.Runtime != "go" && job.Runtime != "deno") || job.Command == "" || strings.ContainsAny(job.Command, "\x00\r") {
				t.Fatal("actions-definition: invalid required job")
			}
			seen[job.ID] = true
		}
	}
	return definitions
}

func (w developmentActionsWorkflow) gateID() string {
	return fmt.Sprintf("phase%d-required-gate", w.Phase)
}

func developmentActionsIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for _, c := range value {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' {
			return false
		}
	}
	return true
}

// This is a fixed workflow serialization oracle, not a general YAML parser.
func developmentActionsWorkflowText(w developmentActionsWorkflow) string {
	if w.Setup == nil {
		return "missing setup definition\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "name: %s\n\non:\n  pull_request:\n    types: [opened, synchronize, reopened, ready_for_review]\n", w.Name)
	if w.Phase == 18 {
		b.WriteString("  workflow_dispatch:\n")
	}
	b.WriteString("  push:\n    branches:\n      - main\n\npermissions:\n  contents: read\n\nconcurrency:\n  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.run_id }}\n  cancel-in-progress: ${{ github.event_name == 'pull_request' }}\n\ndefaults:\n  run:\n    shell: bash\n\njobs:\n")
	checkout := "      - name: Checkout exact source without external actions\n        run: |\n          " + strings.ReplaceAll(w.Setup.Checkout, "\n", "\n          ") + "\n"
	cleanup := "      - name: Clean temporary toolchain\n        if: ${{ always() }}\n        run: |\n          if test -f .github/scripts/development-ci.sh; then\n            bash .github/scripts/development-ci.sh cleanup\n          else\n            test -z \"${ADLAIRE_CI_TOOLS:-}\"\n          fi\n"
	for i, job := range w.Jobs {
		if i != 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "  %s:\n    runs-on: ubuntu-24.04\n    timeout-minutes: %d\n    steps:\n", job.ID, job.Timeout)
		b.WriteString(checkout)
		fmt.Fprintf(&b, "      - name: Prepare verified %s runtime\n        run: bash .github/scripts/development-ci.sh %s\n", job.Runtime, job.Runtime)
		if strings.Contains(job.Command, "\n") {
			b.WriteString("      - run: |\n          ")
			b.WriteString(strings.ReplaceAll(job.Command, "\n", "\n          "))
			b.WriteByte('\n')
		} else {
			fmt.Fprintf(&b, "      - run: %s\n", job.Command)
		}
		b.WriteString(cleanup)
	}
	fmt.Fprintf(&b, "\n  %s:\n    if: ${{ always() }}\n    needs:\n", w.gateID())
	for _, job := range w.Jobs {
		fmt.Fprintf(&b, "      - %s\n", job.ID)
	}
	b.WriteString("    runs-on: ubuntu-24.04\n    timeout-minutes: 10\n    steps:\n")
	b.WriteString(checkout)
	b.WriteString("      - name: Prepare verified go runtime\n        run: bash .github/scripts/development-ci.sh go\n")
	b.WriteString("      - name: Verify required checks and checkout identity\n        env:\n          ADLAIRE_ACTIONS_GATE: required\n          ADLAIRE_ACTIONS_NEEDS: ${{ toJSON(needs) }}\n        run: |\n          ")
	b.WriteString(strings.ReplaceAll(developmentActionsGateCommand, "\n", "\n          "))
	b.WriteByte('\n')
	b.WriteString(cleanup)
	return b.String()
}

func developmentActionsCheckWorkflow(w developmentActionsWorkflow, actual string) error {
	if actual != developmentActionsWorkflowText(w) {
		expected := sha256.Sum256([]byte(developmentActionsWorkflowText(w)))
		got := sha256.Sum256([]byte(actual))
		return fmt.Errorf("actions-workflow-contract: %s expected_sha256=%x actual_sha256=%x", w.Name, expected, got)
	}
	return nil
}

func TestDevelopmentActionsWorkflowContract(t *testing.T) {
	t.Parallel()
	setup := developmentActionsSetupDefinition(t)
	if err := developmentActionsCheckHelper(setup, []byte(phase11MustReadText(t, ".github/scripts/development-ci.sh"))); err != nil {
		t.Fatal(err)
	}
	if err := developmentActionsCheckNodeFreeRepository("."); err != nil {
		t.Fatal(err)
	}
	for _, definition := range developmentActionsDefinitions(t) {
		t.Run("release.development-actions."+definition.Name+".workflow-contract.exact", func(t *testing.T) {
			path := ".github/workflows/" + definition.Name + ".yml"
			if err := developmentActionsCheckWorkflow(definition, phase11MustReadText(t, path)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func developmentActionsDecodeJSON(data []byte, target any, limit int, strict bool) error {
	if len(data) == 0 || len(data) > limit || !utf8.Valid(data) {
		return fmt.Errorf("json-size-or-encoding")
	}
	if !json.Valid(data) {
		return fmt.Errorf("json-syntax-or-depth")
	}
	check := json.NewDecoder(bytes.NewReader(data))
	if err := phase16CheckJSONValueForDuplicateKeys(check); err != nil {
		return fmt.Errorf("json-duplicate-or-syntax")
	}
	if _, err := check.Token(); err != io.EOF {
		return fmt.Errorf("json-trailing-data")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if strict {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("json-shape")
	}
	return nil
}

type developmentActionsGateInput struct {
	Env       map[string]string
	Needs     []byte
	Event     []byte
	TestedSHA string
	Parents   []string
	OSRelease string
}

func developmentActionsSHA(value string) bool {
	if len(value) != 40 || value != strings.ToLower(value) {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func developmentActionsPositiveInteger(value string) bool {
	n, err := strconv.ParseUint(value, 10, 64)
	return err == nil && n > 0 && strconv.FormatUint(n, 10) == value
}

func developmentActionsValidateGate(w developmentActionsWorkflow, input developmentActionsGateInput) (map[string]any, error) {
	env := input.Env
	fail := func(code string) (map[string]any, error) { return nil, fmt.Errorf("actions-gate:%s", code) }
	if env["ADLAIRE_ACTIONS_GATE"] != "required" || env["GITHUB_ACTIONS"] != "true" {
		return fail("gate-mode")
	}
	if env["GITHUB_REPOSITORY"] != "fqwink/Build-Scripts" || env["GITHUB_SERVER_URL"] != "https://github.com" || !developmentActionsPositiveInteger(env["GITHUB_RUN_ID"]) || !developmentActionsPositiveInteger(env["GITHUB_RUN_ATTEMPT"]) || env["GITHUB_JOB"] != w.gateID() || env["GITHUB_WORKFLOW"] != w.Name || env["GITHUB_WORKFLOW_REF"] != env["GITHUB_REPOSITORY"]+"/.github/workflows/"+w.Name+".yml@"+env["GITHUB_REF"] {
		return fail("run-identity")
	}
	if env["RUNNER_ENVIRONMENT"] != "github-hosted" || env["RUNNER_OS"] != "Linux" || env["RUNNER_ARCH"] != "X64" || !strings.Contains("\n"+input.OSRelease, "\nID=ubuntu\n") || !strings.Contains("\n"+input.OSRelease, "\nVERSION_ID=\"24.04\"\n") {
		return fail("runner-identity")
	}
	target := env["GITHUB_SHA"]
	if !developmentActionsSHA(target) || !developmentActionsSHA(input.TestedSHA) || input.TestedSHA != target || env["GITHUB_WORKFLOW_SHA"] != target {
		return fail("tested-sha")
	}
	var event map[string]json.RawMessage
	if err := developmentActionsDecodeJSON(input.Event, &event, 2<<20, false); err != nil || event == nil {
		return fail("event-json")
	}
	source, base := target, ""
	switch env["GITHUB_EVENT_NAME"] {
	case "pull_request":
		var number int
		var pr map[string]json.RawMessage
		var head, baseRef map[string]json.RawMessage
		if json.Unmarshal(event["number"], &number) != nil || number < 1 || json.Unmarshal(event["pull_request"], &pr) != nil || json.Unmarshal(pr["head"], &head) != nil || json.Unmarshal(pr["base"], &baseRef) != nil || json.Unmarshal(head["sha"], &source) != nil || json.Unmarshal(baseRef["sha"], &base) != nil || !developmentActionsSHA(source) || !developmentActionsSHA(base) || env["GITHUB_REF"] != fmt.Sprintf("refs/pull/%d/merge", number) || len(input.Parents) != 2 || input.Parents[0] != base || input.Parents[1] != source {
			return fail("pr-checkout")
		}
	case "push":
		var after, ref string
		if json.Unmarshal(event["after"], &after) != nil || json.Unmarshal(event["ref"], &ref) != nil || after != target || ref != "refs/heads/main" || env["GITHUB_REF"] != ref {
			return fail("push-checkout")
		}
	case "workflow_dispatch":
		ref := env["GITHUB_REF"]
		if w.Phase != 18 || (!strings.HasPrefix(ref, "refs/heads/") && !strings.HasPrefix(ref, "refs/tags/")) || strings.HasSuffix(ref, "/") || strings.ContainsAny(ref, "\x00\r\n ") {
			return fail("dispatch-checkout")
		}
	default:
		return fail("event-type")
	}
	var needs map[string]json.RawMessage
	if err := developmentActionsDecodeJSON(input.Needs, &needs, 256<<10, true); err != nil || needs == nil {
		return fail("needs-json")
	}
	if len(needs) != len(w.Jobs) {
		return fail("needs-count")
	}
	results := make([]map[string]string, 0, len(w.Jobs))
	for _, job := range w.Jobs {
		var fields map[string]json.RawMessage
		if json.Unmarshal(needs[job.ID], &fields) != nil || len(fields) != 2 {
			return fail("needs-job")
		}
		var result string
		var outputs map[string]json.RawMessage
		if json.Unmarshal(fields["result"], &result) != nil || json.Unmarshal(fields["outputs"], &outputs) != nil || outputs == nil || len(outputs) != 0 {
			return fail("needs-job")
		}
		if result != "success" {
			return fail("required-job-not-success")
		}
		results = append(results, map[string]string{"job": job.ID, "result": result})
	}
	definition, err := json.Marshal(w)
	if err != nil {
		return fail("definition-encoding")
	}
	digest := sha256.Sum256(definition)
	attempt, err := strconv.ParseUint(env["GITHUB_RUN_ATTEMPT"], 10, 64)
	if err != nil {
		return fail("run-identity")
	}
	var baseValue any
	if base != "" {
		baseValue = base
	}
	return map[string]any{
		"record_type": "development-actions-gate", "workflow": w.Name, "job": w.gateID(),
		"run_url":     env["GITHUB_SERVER_URL"] + "/" + env["GITHUB_REPOSITORY"] + "/actions/runs/" + env["GITHUB_RUN_ID"],
		"run_attempt": attempt, "event": env["GITHUB_EVENT_NAME"],
		"source_sha": source, "base_sha": baseValue, "tested_sha": input.TestedSHA,
		"definition_sha256": hex.EncodeToString(digest[:]), "job_results": results,
		"runner_environment": env["RUNNER_ENVIRONMENT"], "runner_os": env["RUNNER_OS"], "runner_arch": env["RUNNER_ARCH"],
	}, nil
}

func developmentActionsReadContext(getenv func(string) string, readFile func(string) ([]byte, error), git func(...string) ([]byte, error)) (developmentActionsGateInput, error) {
	input := developmentActionsGateInput{Env: map[string]string{}, Needs: []byte(getenv("ADLAIRE_ACTIONS_NEEDS"))}
	for _, key := range []string{"ADLAIRE_ACTIONS_GATE", "GITHUB_ACTIONS", "GITHUB_REPOSITORY", "GITHUB_SERVER_URL", "GITHUB_RUN_ID", "GITHUB_RUN_ATTEMPT", "GITHUB_JOB", "GITHUB_WORKFLOW", "GITHUB_WORKFLOW_REF", "GITHUB_REF", "GITHUB_SHA", "GITHUB_WORKFLOW_SHA", "GITHUB_EVENT_NAME", "RUNNER_ENVIRONMENT", "RUNNER_OS", "RUNNER_ARCH"} {
		input.Env[key] = getenv(key)
	}
	var err error
	input.Event, err = readFile(getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		return input, fmt.Errorf("actions-gate:event-read")
	}
	osRelease, err := readFile("/etc/os-release")
	if err != nil {
		return input, fmt.Errorf("actions-gate:os-release")
	}
	input.OSRelease = string(osRelease)
	checkout, err := git("rev-parse", "HEAD")
	if err != nil {
		return input, fmt.Errorf("actions-gate:checkout-read")
	}
	input.TestedSHA = strings.TrimSpace(string(checkout))
	commit, err := git("cat-file", "-p", "HEAD")
	if err != nil {
		return input, fmt.Errorf("actions-gate:commit-read")
	}
	for _, line := range strings.Split(string(commit), "\n") {
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "parent ") {
			input.Parents = append(input.Parents, strings.TrimPrefix(line, "parent "))
		}
	}
	return input, nil
}

func TestDevelopmentActionsRequiredGate(t *testing.T) {
	if os.Getenv("ADLAIRE_ACTIONS_GATE") == "" {
		if _, err := developmentActionsValidateGate(developmentActionsWorkflow{}, developmentActionsGateInput{}); err == nil || err.Error() != "actions-gate:gate-mode" {
			t.Fatal("actions-gate: unit invocation must reject absent formal inputs")
		}
		return
	}
	setup := developmentActionsSetupDefinition(t)
	if err := developmentActionsCheckHelper(setup, []byte(phase11MustReadText(t, ".github/scripts/development-ci.sh"))); err != nil {
		t.Fatal(err)
	}
	if err := developmentActionsCheckNodeFreeRepository("."); err != nil {
		t.Fatal(err)
	}
	var definition developmentActionsWorkflow
	for _, candidate := range developmentActionsDefinitions(t) {
		if err := developmentActionsCheckWorkflow(candidate, phase11MustReadText(t, ".github/workflows/"+candidate.Name+".yml")); err != nil {
			t.Fatal(err)
		}
		if candidate.Name == os.Getenv("GITHUB_WORKFLOW") {
			definition = candidate
		}
	}
	if definition.Name == "" {
		t.Fatal("actions-gate: unknown workflow")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	input, err := developmentActionsReadContext(os.Getenv, os.ReadFile, func(args ...string) ([]byte, error) {
		return exec.CommandContext(ctx, "git", args...).Output()
	})
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" || runtime.Version() != "go"+definition.Setup.GoVersion {
		t.Fatal("actions-gate:runtime")
	}
	record, err := developmentActionsValidateGate(definition, input)
	if err != nil {
		t.Fatal(err)
	}
	record["runtime"] = runtime.Version()
	setupJSON, err := json.Marshal(definition.Setup)
	if err != nil {
		t.Fatal("actions-setup: encoding")
	}
	setupSum := sha256.Sum256(setupJSON)
	record["setup_definition_sha256"] = hex.EncodeToString(setupSum[:])
	record["helper_sha256"] = definition.Setup.HelperSHA256
	record["external_action_count"] = 0
	record["log_retention_days"] = os.Getenv("GITHUB_RETENTION_DAYS")
	encoded, err := json.Marshal(record)
	if err != nil {
		t.Fatal("actions-gate:record-encoding")
	}
	fmt.Println(string(encoded))
}

func TestDevelopmentActionsGateControls(t *testing.T) {
	t.Parallel()
	for _, w := range developmentActionsDefinitions(t) {
		w := w
		t.Run(w.Name, func(t *testing.T) {
			t.Parallel()
			for _, job := range w.Jobs {
				if !strings.HasSuffix(job.ID, "-go-format") {
					continue
				}
				for _, fixture := range []struct {
					name, formatter string
					success         bool
				}{
					{"clean", "gofmt() { return 0; }\n", true},
					{"dirty", "gofmt() { printf 'dirty.go\\n'; }\n", false},
					{"failed", "gofmt() { return 7; }\n", false},
				} {
					t.Run("release.development-actions."+job.ID+"."+fixture.name+".exit-code", func(t *testing.T) {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						output, err := exec.CommandContext(ctx, "bash", "--noprofile", "--norc", "-eo", "pipefail", "-c", fixture.formatter+job.Command).CombinedOutput()
						if (err == nil) != fixture.success || len(output) != 0 {
							t.Fatalf("formatter exit propagation expected success=%t, err=%v output=%q", fixture.success, err, output)
						}
					})
				}
			}
			fresh := func() developmentActionsGateInput {
				target, base, source := strings.Repeat("a", 40), strings.Repeat("b", 40), strings.Repeat("c", 40)
				needs := map[string]any{}
				for _, job := range w.Jobs {
					needs[job.ID] = map[string]any{"result": "success", "outputs": map[string]any{}}
				}
				encoded, err := json.Marshal(needs)
				if err != nil {
					t.Fatal(err)
				}
				return developmentActionsGateInput{
					Env:   map[string]string{"ADLAIRE_ACTIONS_GATE": "required", "GITHUB_ACTIONS": "true", "GITHUB_REPOSITORY": "fqwink/Build-Scripts", "GITHUB_SERVER_URL": "https://github.com", "GITHUB_RUN_ID": "123", "GITHUB_RUN_ATTEMPT": "1", "GITHUB_JOB": w.gateID(), "GITHUB_WORKFLOW": w.Name, "GITHUB_WORKFLOW_REF": "fqwink/Build-Scripts/.github/workflows/" + w.Name + ".yml@refs/pull/99/merge", "GITHUB_REF": "refs/pull/99/merge", "GITHUB_SHA": target, "GITHUB_WORKFLOW_SHA": target, "GITHUB_EVENT_NAME": "pull_request", "RUNNER_ENVIRONMENT": "github-hosted", "RUNNER_OS": "Linux", "RUNNER_ARCH": "X64"},
					Needs: encoded, Event: []byte(`{"number":99,"pull_request":{"head":{"sha":"` + source + `"},"base":{"sha":"` + base + `"}}}`), TestedSHA: target, Parents: []string{base, source}, OSRelease: "ID=ubuntu\nVERSION_ID=\"24.04\"\n",
				}
			}
			passLines := "--- PASS: TestDevelopmentActionsRequiredGate (0.01s)\n--- PASS: TestDevelopmentActionsWorkflowContract (0.01s)\n--- PASS: TestDevelopmentActionsGateControls (0.01s)\n--- PASS: TestDevelopmentActionsNodeFreeControls (0.01s)\n--- PASS: TestDevelopmentActionsBootstrapControls (0.01s)\n"
			marker := `{"record_type":"development-actions-gate"}` + "\n"
			for _, fixture := range []struct {
				name, output string
				exit         int
				mode, needs  string
				success      bool
			}{
				{"complete", passLines + marker, 0, "required", "{}", true},
				{"no-tests", "testing: warning: no tests to run\nPASS\n", 0, "required", "{}", false},
				{"missing-root", strings.Replace(passLines, "--- PASS: TestDevelopmentActionsRequiredGate (0.01s)\n", "", 1) + marker, 0, "required", "{}", false},
				{"missing-marker", passLines, 0, "required", "{}", false},
				{"failed-exit", passLines + marker, 7, "required", "{}", false},
				{"missing-mode", passLines + marker, 0, "", "{}", false},
				{"missing-needs", passLines + marker, 0, "required", "", false},
			} {
				t.Run("release.development-actions."+w.Name+".script-"+fixture.name+".exit-code", func(t *testing.T) {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					fake := "go() { printf '%s' '" + fixture.output + "'; return " + strconv.Itoa(fixture.exit) + "; }\n"
					command := exec.CommandContext(ctx, "bash", "--noprofile", "--norc", "-eo", "pipefail", "-c", fake+developmentActionsGateCommand)
					command.Env = append(os.Environ(), "RUNNER_TEMP="+t.TempDir(), "ADLAIRE_ACTIONS_GATE="+fixture.mode, "ADLAIRE_ACTIONS_NEEDS="+fixture.needs)
					output, err := command.CombinedOutput()
					if (err == nil) != fixture.success {
						t.Fatalf("gate script expected success=%t, err=%v output=%q", fixture.success, err, output)
					}
				})
			}
			positive := fresh()
			record, err := developmentActionsValidateGate(w, positive)
			if err != nil || record["tested_sha"] != positive.TestedSHA || record["source_sha"] != positive.Parents[1] || record["base_sha"] != positive.Parents[0] || len(record["job_results"].([]map[string]string)) != len(w.Jobs) {
				t.Fatalf("actions-control: positive PR oracle: %v", err)
			}
			for _, field := range []string{"needs", "event"} {
				t.Run("release.development-actions."+w.Name+"."+field+"-size-limit.accept", func(t *testing.T) {
					input := fresh()
					if field == "needs" {
						input.Needs = append(input.Needs, bytes.Repeat([]byte(" "), (256<<10)-len(input.Needs))...)
					} else {
						input.Event = append(input.Event, bytes.Repeat([]byte(" "), (2<<20)-len(input.Event))...)
					}
					if _, err := developmentActionsValidateGate(w, input); err != nil {
						t.Fatal(err)
					}
				})
			}
			for fault := -1; fault < 4; fault++ {
				fixture := fresh()
				calls := []string{}
				steps := []string{"read:event", "read:os-release", "git:rev-parse HEAD", "git:cat-file -p HEAD"}
				codes := []string{"event-read", "os-release", "checkout-read", "commit-read"}
				call := func(name string, data []byte) ([]byte, error) {
					index := len(calls)
					calls = append(calls, name)
					if index >= len(steps) || steps[index] != name {
						t.Fatalf("assertion=release.development-actions.context-order.exact expected %v got %v", steps, calls)
					}
					if index == fault {
						return nil, fmt.Errorf("fixture-secret-canary")
					}
					return data, nil
				}
				getenv := func(key string) string {
					switch key {
					case "GITHUB_EVENT_PATH":
						return "fixture-event"
					case "ADLAIRE_ACTIONS_NEEDS":
						return string(fixture.Needs)
					default:
						return fixture.Env[key]
					}
				}
				input, err := developmentActionsReadContext(getenv, func(path string) ([]byte, error) {
					if path == "fixture-event" {
						return call("read:event", fixture.Event)
					}
					if path != "/etc/os-release" {
						t.Fatal("actions-control: unexpected file read")
					}
					return call("read:os-release", []byte(fixture.OSRelease))
				}, func(args ...string) ([]byte, error) {
					if strings.Join(args, " ") == "rev-parse HEAD" {
						return call("git:rev-parse HEAD", []byte(fixture.TestedSHA+"\n"))
					}
					return call("git:"+strings.Join(args, " "), []byte("tree "+strings.Repeat("d", 40)+"\nparent "+fixture.Parents[0]+"\nparent "+fixture.Parents[1]+"\n\nparent must-not-parse-body\n"))
				})
				if fault >= 0 {
					if err == nil || err.Error() != "actions-gate:"+codes[fault] || len(calls) != fault+1 {
						t.Fatalf("assertion=release.development-actions.%s.fail-closed err=%v calls=%v", codes[fault], err, calls)
					}
				} else if err != nil {
					t.Fatal(err)
				} else if _, err := developmentActionsValidateGate(w, input); err != nil || len(calls) != 4 {
					t.Fatalf("assertion=release.development-actions.context-positive.accept err=%v calls=%v", err, calls)
				}
			}
			canary := fresh()
			canary.Event = []byte(strings.Replace(string(canary.Event), "{", `{"unknown":"fixture-secret-canary",`, 1))
			canaryRecord, err := developmentActionsValidateGate(w, canary)
			if err != nil {
				t.Fatal(err)
			}
			canaryJSON, err := json.Marshal(canaryRecord)
			if err != nil || bytes.Contains(canaryJSON, []byte("fixture-secret-canary")) {
				t.Fatal("assertion=release.development-actions.secret-safe.no-leak")
			}
			for _, event := range []string{"push", "workflow_dispatch"} {
				input := fresh()
				input.Env["GITHUB_EVENT_NAME"] = event
				input.Env["GITHUB_REF"] = "refs/heads/main"
				input.Env["GITHUB_WORKFLOW_REF"] = "fqwink/Build-Scripts/.github/workflows/" + w.Name + ".yml@refs/heads/main"
				input.Event = []byte(`{"ref":"refs/heads/main","after":"` + input.TestedSHA + `"}`)
				_, err := developmentActionsValidateGate(w, input)
				if (event == "push" || w.Phase == 18) != (err == nil) {
					t.Fatalf("actions-control: event %s applicability: %v", event, err)
				}
			}
			type gateCase struct {
				name string
				edit func(*developmentActionsGateInput)
				code string
			}
			tests := []gateCase{
				{"mode", func(i *developmentActionsGateInput) { delete(i.Env, "ADLAIRE_ACTIONS_GATE") }, "gate-mode"},
				{"not-actions", func(i *developmentActionsGateInput) { i.Env["GITHUB_ACTIONS"] = "false" }, "gate-mode"},
				{"tested-sha", func(i *developmentActionsGateInput) { i.TestedSHA = strings.Repeat("d", 40) }, "tested-sha"},
				{"workflow-sha", func(i *developmentActionsGateInput) { i.Env["GITHUB_WORKFLOW_SHA"] = strings.Repeat("d", 40) }, "tested-sha"},
				{"bad-sha", func(i *developmentActionsGateInput) { i.Env["GITHUB_SHA"] = "invalid" }, "tested-sha"},
				{"uppercase-sha", func(i *developmentActionsGateInput) { i.Env["GITHUB_SHA"] = strings.Repeat("A", 40) }, "tested-sha"},
				{"base", func(i *developmentActionsGateInput) { i.Parents[0] = strings.Repeat("d", 40) }, "pr-checkout"},
				{"head", func(i *developmentActionsGateInput) { i.Parents[1] = strings.Repeat("d", 40) }, "pr-checkout"},
				{"parent-order", func(i *developmentActionsGateInput) { i.Parents[0], i.Parents[1] = i.Parents[1], i.Parents[0] }, "pr-checkout"},
				{"parent-count", func(i *developmentActionsGateInput) { i.Parents = nil }, "pr-checkout"},
				{"event-json", func(i *developmentActionsGateInput) { i.Event = []byte("{") }, "event-json"},
				{"event-empty", func(i *developmentActionsGateInput) { i.Event = []byte("{}") }, "pr-checkout"},
				{"event-null", func(i *developmentActionsGateInput) { i.Event = []byte("null") }, "event-json"},
				{"event-duplicate", func(i *developmentActionsGateInput) {
					i.Event = []byte(strings.Replace(string(i.Event), `"number":99`, `"number":98,"number":99`, 1))
				}, "event-json"},
				{"event-trailing", func(i *developmentActionsGateInput) { i.Event = append(i.Event, []byte(" {}")...) }, "event-json"},
				{"event-oversize", func(i *developmentActionsGateInput) { i.Event = bytes.Repeat([]byte(" "), (2<<20)+1) }, "event-json"},
				{"pr-ref", func(i *developmentActionsGateInput) {
					i.Env["GITHUB_REF"] = "refs/pull/98/merge"
					i.Env["GITHUB_WORKFLOW_REF"] = "fqwink/Build-Scripts/.github/workflows/" + w.Name + ".yml@refs/pull/98/merge"
				}, "pr-checkout"},
				{"event-type", func(i *developmentActionsGateInput) { i.Env["GITHUB_EVENT_NAME"] = "pull_request_target" }, "event-type"},
				{"push-after", func(i *developmentActionsGateInput) { i.Env["GITHUB_EVENT_NAME"] = "push" }, "push-checkout"},
				{"dispatch-ref", func(i *developmentActionsGateInput) { i.Env["GITHUB_EVENT_NAME"] = "workflow_dispatch" }, "dispatch-checkout"},
				{"self-hosted", func(i *developmentActionsGateInput) { i.Env["RUNNER_ENVIRONMENT"] = "self-hosted" }, "runner-identity"},
				{"runner-os", func(i *developmentActionsGateInput) { i.Env["RUNNER_OS"] = "Windows" }, "runner-identity"},
				{"runner-arch", func(i *developmentActionsGateInput) { i.Env["RUNNER_ARCH"] = "ARM64" }, "runner-identity"},
				{"os-version", func(i *developmentActionsGateInput) { i.OSRelease = "ID=ubuntu\nVERSION_ID=\"22.04\"\n" }, "runner-identity"},
				{"needs-missing", func(i *developmentActionsGateInput) { i.Needs = []byte("{}") }, "needs-count"},
				{"needs-extra", func(i *developmentActionsGateInput) {
					i.Needs = []byte(strings.Replace(string(i.Needs), "{", `{"unknown":{"result":"success","outputs":{}},`, 1))
				}, "needs-count"},
				{"needs-renamed", func(i *developmentActionsGateInput) {
					i.Needs = []byte(strings.Replace(string(i.Needs), `"`+w.Jobs[0].ID+`":`, `"unknown":`, 1))
				}, "needs-job"},
				{"needs-duplicate", func(i *developmentActionsGateInput) {
					i.Needs = []byte(strings.Replace(string(i.Needs), "{", `{"`+w.Jobs[0].ID+`":{"result":"success","outputs":{}},`, 1))
				}, "needs-json"},
				{"needs-null", func(i *developmentActionsGateInput) { i.Needs = []byte("null") }, "needs-json"},
				{"needs-malformed", func(i *developmentActionsGateInput) { i.Needs = []byte("{") }, "needs-json"},
				{"needs-depth", func(i *developmentActionsGateInput) {
					i.Needs = []byte(strings.Repeat("[", 10001) + "0" + strings.Repeat("]", 10001))
				}, "needs-json"},
				{"needs-trailing", func(i *developmentActionsGateInput) { i.Needs = append(i.Needs, []byte(" {}")...) }, "needs-json"},
				{"needs-oversize", func(i *developmentActionsGateInput) { i.Needs = bytes.Repeat([]byte(" "), (256<<10)+1) }, "needs-json"},
				{"needs-utf8", func(i *developmentActionsGateInput) { i.Needs = []byte{0xff} }, "needs-json"},
			}
			for _, key := range []string{"GITHUB_REPOSITORY", "GITHUB_SERVER_URL", "GITHUB_RUN_ID", "GITHUB_RUN_ATTEMPT", "GITHUB_JOB", "GITHUB_WORKFLOW", "GITHUB_WORKFLOW_REF"} {
				key := key
				tests = append(tests, gateCase{"identity-" + strings.ReplaceAll(strings.ToLower(key), "_", "-"), func(i *developmentActionsGateInput) { i.Env[key] = "invalid" }, "run-identity"})
			}
			for _, test := range tests {
				t.Run("release.development-actions."+w.Name+"."+test.name+".reject", func(t *testing.T) {
					input := fresh()
					test.edit(&input)
					if _, err := developmentActionsValidateGate(w, input); err == nil || err.Error() != "actions-gate:"+test.code {
						t.Errorf("expected %s, got %v", test.code, err)
					}
				})
			}
			first := w.Jobs[0].ID
			for _, invalid := range []struct{ name, value string }{
				{"failure", `{"result":"failure","outputs":{}}`},
				{"skipped", `{"result":"skipped","outputs":{}}`},
				{"neutral", `{"result":"neutral","outputs":{}}`},
				{"cancelled", `{"result":"cancelled","outputs":{}}`},
				{"timeout", `{"result":"timed_out","outputs":{}}`},
				{"empty-result", `{"result":"","outputs":{}}`},
				{"null-result", `{"result":null,"outputs":{}}`},
				{"case-alias", `{"Result":"success","outputs":{}}`},
				{"null-outputs", `{"result":"success","outputs":null}`},
				{"unexpected-outputs", `{"result":"success","outputs":{"ignored":"value"}}`},
				{"unknown-field", `{"result":"success","outputs":{},"extra":true}`},
				{"duplicate-result", `{"result":"skipped","result":"success","outputs":{}}`},
				{"escaped-duplicate", `{"result":"skipped","resul\u0074":"success","outputs":{}}`},
				{"null-job", `null`}, {"array-job", `[]`}, {"boolean-job", `true`},
			} {
				t.Run("release.development-actions."+w.Name+"."+invalid.name+".reject", func(t *testing.T) {
					input := fresh()
					input.Needs = []byte(strings.Replace(string(input.Needs), `"`+first+`":{"outputs":{},"result":"success"}`, `"`+first+`":`+invalid.value, 1))
					if bytes.Equal(input.Needs, positive.Needs) {
						t.Fatal("actions-control: invalid needs was not injected")
					}
					if _, err := developmentActionsValidateGate(w, input); err == nil {
						t.Error("invalid job result accepted")
					}
				})
			}
			canonical := developmentActionsWorkflowText(w)
			if err := developmentActionsCheckWorkflow(w, canonical); err != nil {
				t.Fatal(err)
			}
			for _, change := range [][2]string{
				{"contents: read", "contents: write"}, {"timeout-minutes: 10", "timeout-minutes: 0"},
				{"${{ always() }}", "${{ success() }}"}, {"test \"$tested_sha\" = \"$GITHUB_SHA\"", "true"},
				{"GIT_TERMINAL_PROMPT=0", "GIT_TERMINAL_PROMPT=1"}, {"--no-tags --depth=1", "--tags --depth=5"},
				{"runs-on: ubuntu-24.04", "runs-on: self-hosted"}, {"cancel-in-progress:", "cancel-disabled:"},
				{"ADLAIRE_ACTIONS_GATE: required", "ADLAIRE_ACTIONS_GATE: optional"},
				{"      - " + first + "\n", ""}, {strings.ReplaceAll(developmentActionsGateCommand, "\n", "\n          "), "true"},
				{"    steps:\n", "    continue-on-error: true\n    steps:\n"},
				{"    steps:\n", "    if: false\n    steps:\n"},
			} {
				mutated := strings.Replace(canonical, change[0], change[1], 1)
				if mutated == canonical || developmentActionsCheckWorkflow(w, mutated) == nil {
					t.Fatalf("actions-control: canonical workflow mutation not rejected: %s", change[0])
				}
			}
		})
	}
}

var developmentActionsNodeCommand = regexp.MustCompile("(?i)(^|[\\s;&|(\"'\\x60])(?:[^\\s;&|(\"'\\x60]*/)?(?:node(?:js)?|npm|npx|pnpm|yarn|corepack)(?:\\.exe)?($|[\\s;&|)\"'\\x60])")
var developmentActionsNodeAPI = regexp.MustCompile("(?m)[\"'\\x60](?:node:|npm:)|\\brequire\\s*\\(|\\bmodule\\s*\\.\\s*exports\\b|\\bexports\\s*[.\\[]|\\bprocess\\s*[.\\[]|\\bBuffer\\s*\\.|\\b__dirname\\b|\\b__filename\\b")

func developmentActionsStaticString(expr ast.Expr) (string, bool) {
	switch expr := expr.(type) {
	case *ast.BasicLit:
		if expr.Kind == token.STRING {
			value, err := strconv.Unquote(expr.Value)
			return value, err == nil
		}
	case *ast.ParenExpr:
		return developmentActionsStaticString(expr.X)
	case *ast.BinaryExpr:
		if expr.Op == token.ADD {
			left, okLeft := developmentActionsStaticString(expr.X)
			right, okRight := developmentActionsStaticString(expr.Y)
			return left + right, okLeft && okRight
		}
	}
	return "", false
}

func developmentActionsCheckNodeContent(path string, data []byte) error {
	ext := filepath.Ext(path)
	if ext == ".go" {
		tree, err := parser.ParseFile(token.NewFileSet(), path, data, parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf("node-ban:go-syntax")
		}
		execNames := map[string]bool{}
		processNames := map[string]bool{}
		for _, imp := range tree.Imports {
			name, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				return fmt.Errorf("node-ban:go-syntax")
			}
			if name == "os/exec" || name == "os" || name == "syscall" {
				alias := "exec"
				if name == "os" || name == "syscall" {
					alias = name
				}
				if imp.Name != nil {
					alias = imp.Name.Name
				}
				if alias == "." {
					return fmt.Errorf("node-ban:command")
				}
				if name == "os/exec" {
					execNames[alias] = true
				} else {
					processNames[alias] = true
				}
			}
		}
		var violation error
		ast.Inspect(tree, func(n ast.Node) bool {
			if violation != nil {
				return false
			}
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}
			spawn := execNames[pkg.Name] && (sel.Sel.Name == "Command" || sel.Sel.Name == "CommandContext" || sel.Sel.Name == "LookPath")
			spawn = spawn || processNames[pkg.Name] && (sel.Sel.Name == "StartProcess" || sel.Sel.Name == "Exec" || sel.Sel.Name == "ForkExec")
			if !spawn {
				return true
			}
			start := 0
			if sel.Sel.Name == "CommandContext" {
				start = 1
			}
			if len(call.Args) <= start {
				violation = fmt.Errorf("node-ban:go-syntax")
				return false
			}
			for _, argument := range call.Args[start:] {
				if value, ok := developmentActionsStaticString(argument); ok && developmentActionsNodeCommand.MatchString(value) {
					violation = fmt.Errorf("node-ban:command")
				}
			}
			return true
		})
		if violation != nil {
			return violation
		}
		if !strings.HasSuffix(path, "_test.go") {
			var found bool
			ast.Inspect(tree, func(n ast.Node) bool {
				if value, ok := n.(*ast.BasicLit); ok && value.Kind == token.STRING {
					text, err := strconv.Unquote(value.Value)
					if err == nil && developmentActionsNodeAPI.MatchString(text) {
						found = true
					}
				}
				return true
			})
			if found {
				return fmt.Errorf("node-ban:api")
			}
		}
		return nil
	}
	if ext == ".js" || ext == ".mjs" || ext == ".html" || ext == ".sh" {
		if developmentActionsNodeAPI.Match(data) {
			return fmt.Errorf("node-ban:api")
		}
		if ext == ".sh" && developmentActionsNodeCommand.Match(data) {
			return fmt.Errorf("node-ban:command")
		}
		if regexp.MustCompile("(?m)^#![^\\n]*\\bnode(?:js)?\\b").Match(data) {
			return fmt.Errorf("node-ban:command")
		}
		if ext != ".sh" {
			imports := regexp.MustCompile("(?:\\bfrom\\s+|\\bimport\\s*(?:\\(\\s*)?)[\"']([^\"']+)[\"']").FindAllSubmatch(data, -1)
			for _, match := range imports {
				name := string(match[1])
				if !strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "/") && !strings.HasPrefix(name, "https://") {
					return fmt.Errorf("node-ban:api")
				}
			}
			spawns := regexp.MustCompile("\\bDeno\\s*\\.\\s*(?:Command|run)\\s*\\(\\s*(?:\\{\\s*cmd\\s*:\\s*\\[\\s*)?[\"']([^\"']+)[\"']").FindAllSubmatch(data, -1)
			for _, match := range spawns {
				if developmentActionsNodeCommand.Match(match[1]) {
					return fmt.Errorf("node-ban:command")
				}
			}
		}
	}
	return nil
}

func developmentActionsCheckNodeFreeRepository(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("node-ban:read")
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("node-ban:read")
		}
		rel = filepath.ToSlash(rel)
		if rel == ".git" {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("node-ban:symlink")
		}
		name := strings.ToLower(d.Name())
		switch name {
		case "node", "node.exe", "nodejs", "nodejs.exe", "node_modules", "package.json", "package-lock.json", "npm-shrinkwrap.json", ".npmrc", ".nvmrc", ".node-version", ".yarnrc", ".yarnrc.yml", "yarn.lock", "pnpm-lock.yaml", "pnpm-workspace.yaml", "bun.lock", "bun.lockb":
			return fmt.Errorf("node-ban:artifact")
		}
		if strings.HasSuffix(name, ".cjs") || strings.HasSuffix(name, ".cts") {
			return fmt.Errorf("node-ban:artifact")
		}
		if rel == ".github/actions" {
			return fmt.Errorf("node-ban:action")
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(rel, ".github/workflows/") {
			valid := false
			for _, name := range []string{"phase12-quality-gate", "phase13-implementation-alignment-quality", "phase14-obsidian-vault-integration", "phase15-obsidian-local-sync", "phase16-quality-evidence-closure", "phase17-production-validation", "phase18-vps-pull-bootstrap"} {
				if rel == ".github/workflows/"+name+".yml" {
					valid = true
				}
			}
			if !valid {
				return fmt.Errorf("node-ban:workflow")
			}
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("node-ban:read")
		}
		if strings.HasPrefix(rel, ".github/") && (filepath.Ext(rel) == ".yml" || filepath.Ext(rel) == ".yaml") {
			if regexp.MustCompile("(?m)(^|[\\s{,])uses\\s*:").Match(data) {
				return fmt.Errorf("node-ban:action")
			}
			if regexp.MustCompile("(?m)^\\s*(?:container|services)\\s*:").Match(data) {
				return fmt.Errorf("node-ban:action")
			}
			if developmentActionsNodeCommand.Match(data) {
				return fmt.Errorf("node-ban:command")
			}
		}
		return developmentActionsCheckNodeContent(rel, data)
	})
}

func TestDevelopmentActionsNodeFreeControls(t *testing.T) {
	t.Parallel()
	spec := phase11MustReadText(t, "docs/SPEC.md")
	_, policy, ok := strings.Cut(spec, `<a id="policy-node-runtime-ban"></a>`)
	if !ok || strings.Count(spec, `<a id="policy-node-runtime-ban"></a>`) != 1 {
		t.Fatal("Node policy canonical anchor missing or duplicated")
	}
	policy, _, _ = strings.Cut(policy, `<a id="policy-runner-secrets"></a>`)
	if !strings.Contains(policy, "いかなる条件下でも禁止する。例外はない。") {
		t.Fatal("unconditional Node prohibition missing")
	}
	for _, path := range []string{"AGENTS.md", "docs/DETAIL_INDEX.md", "docs/DOCUMENT_INDEX.md"} {
		if !strings.Contains(phase11MustReadText(t, path), "#policy-node-runtime-ban") {
			t.Fatal("Node policy route missing: " + path)
		}
	}
	setup := developmentActionsSetupDefinition(t)
	helper := []byte(phase11MustReadText(t, ".github/scripts/development-ci.sh"))
	if err := developmentActionsCheckHelper(setup, helper); err != nil {
		t.Fatal(err)
	}
	if err := developmentActionsCheckHelper(setup, append(append([]byte{}, helper...), '\n')); err == nil {
		t.Fatal("helper drift not rejected")
	}
	wrongDefinition := *setup
	wrongDefinition.DenoVersion = "2.0.0"
	if err := developmentActionsCheckHelper(&wrongDefinition, helper); err == nil {
		t.Fatal("helper runtime definition drift not rejected")
	}
	if err := developmentActionsCheckNodeFreeRepository("."); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"node", "node.exe", "nodejs", "node_modules", "package.json", "package-lock.json", "npm-shrinkwrap.json", ".npmrc", ".nvmrc", ".node-version", ".yarnrc", ".yarnrc.yml", "yarn.lock", "pnpm-lock.yaml", "pnpm-workspace.yaml", "bun.lock", "bun.lockb", "adapter.cjs", "adapter.cts"} {
		label := strings.NewReplacer(".", "-", "_", "-").Replace(strings.ToLower(name))
		t.Run("release.node-ban.artifact-"+label+".reject", func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, ".hidden", "nested")
			if err := os.MkdirAll(dir, 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, name), nil, 0600); err != nil {
				t.Fatal(err)
			}
			if err := developmentActionsCheckNodeFreeRepository(root); err == nil || err.Error() != "node-ban:artifact" {
				t.Fatalf("expected artifact rejection, got %v", err)
			}
		})
	}
	for _, command := range []string{"node --version", "nodejs script.js", "/usr/bin/node script.js", "exec /opt/runtime/node script.js", "npm install", "npx tool", "pnpm install", "yarn install", "corepack enable", "echo ok; node script.js", "NODE.EXE app.js"} {
		label := regexp.MustCompile("[^a-z0-9-]").ReplaceAllString(strings.ToLower(strings.Fields(command)[0]), "-")
		t.Run("release.node-ban.command-"+label+".reject", func(t *testing.T) {
			if err := developmentActionsCheckNodeContent("script.sh", []byte("#!/bin/sh\n"+command+"\n")); err == nil || err.Error() != "node-ban:command" {
				t.Fatalf("expected command rejection, got %v", err)
			}
			source := "package fixture\nimport e \"os/exec\"\nfunc f(){ e.Command(\"sh\", \"-c\", " + strconv.Quote(command) + ") }\n"
			if err := developmentActionsCheckNodeContent("tool.go", []byte(source)); err == nil || err.Error() != "node-ban:command" {
				t.Fatalf("expected Go command rejection, got %v", err)
			}
		})
	}
	for _, source := range []string{"import fs from 'node:fs';", "import('npm:tool');", "require('fs');", "module.exports = {};", "exports.value = 1;", "process.env.TOKEN", "process.exit(0)", "Buffer.from('x')", "Buffer.byteLength('x')", "__dirname", "__filename"} {
		if err := developmentActionsCheckNodeContent("client.js", []byte(source)); err == nil || err.Error() != "node-ban:api" {
			t.Fatalf("expected API rejection, got %v", err)
		}
	}
	for _, source := range []string{"import fs from 'fs';", "import('child_process');", "import 'buffer';"} {
		if err := developmentActionsCheckNodeContent("client.js", []byte(source)); err == nil || err.Error() != "node-ban:api" {
			t.Fatal("bare Node module not rejected")
		}
	}
	for _, source := range []string{"new Deno.Command('node');", "Deno.run({cmd: ['npm', 'install']});", "#!/usr/bin/env node\n"} {
		if err := developmentActionsCheckNodeContent("tool.js", []byte(source)); err == nil || err.Error() != "node-ban:command" {
			t.Fatal("JavaScript process runtime not rejected")
		}
	}
	for _, source := range []string{"package p\nimport \"os/exec\"\nfunc f(){ exec.Command(\"no\"+\"de\", \"--version\") }", "package p\nimport \"os/exec\"\nfunc f(){ exec.CommandContext(nil,\"npm\",\"install\") }"} {
		if err := developmentActionsCheckNodeContent("tool.go", []byte(source)); err == nil {
			t.Fatal("Go command alias/constant not rejected")
		}
	}
	for _, source := range []string{"package p\nimport \"os\"\nfunc f(){ os.StartProcess(\"/usr/bin/node\",nil,nil) }", "package p\nimport \"syscall\"\nfunc f(){ syscall.Exec(\"node\",nil,nil) }", "package p\nimport . \"os/exec\"\nfunc f(){ Command(\"node\") }"} {
		if err := developmentActionsCheckNodeContent("tool.go", []byte(source)); err == nil || err.Error() != "node-ban:command" {
			t.Fatal("alternate process API not rejected")
		}
	}
	if err := developmentActionsCheckNodeContent("tool.go", []byte("package p\nimport \"os/exec\"\nfunc f(){ exec.CommandContext() }")); err == nil || err.Error() != "node-ban:go-syntax" {
		t.Fatal("incomplete process call not rejected")
	}
	for _, fixture := range []struct{ path, text string }{
		{"client.js", "const node = document.createTextNode('ok'); fetch('/api');"},
		{"script.sh", "#!/bin/sh\ngo version\ndeno check client.js\nprintf '%s' single-node\n"},
		{"sample_test.go", "package p\nconst negative = \"import x from 'node:fs';\"\n"},
	} {
		if err := developmentActionsCheckNodeContent(fixture.path, []byte(fixture.text)); err != nil {
			t.Fatalf("positive input rejected: %v", err)
		}
	}
	root := t.TempDir()
	if err := developmentActionsCheckNodeFreeRepository(root); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing", filepath.Join(root, "adapter")); err != nil {
		t.Fatal(err)
	}
	if err := developmentActionsCheckNodeFreeRepository(root); err == nil || err.Error() != "node-ban:symlink" {
		t.Fatalf("expected symlink rejection, got %v", err)
	}
	for _, fixture := range []struct{ path, text, code string }{
		{".github/workflows/extra.yml", "jobs: {}", "workflow"},
		{".github/workflows/phase12-quality-gate.yml", "jobs:\n  bad:\n    uses: org/repo/.github/workflows/ci.yml@main\n", "action"},
		{".github/workflows/phase12-quality-gate.yml", "jobs:\n  bad:\n    steps: [{uses: './local'}]\n", "action"},
		{".github/workflows/phase12-quality-gate.yml", "jobs:\n  bad:\n    steps:\n      - run: node --version\n", "command"},
		{".github/workflows/phase12-quality-gate.yml", "jobs:\n  bad:\n    container: node:24\n", "action"},
		{".github/actions/local/action.yml", "runs:\n  using: node24\n", "action"},
	} {
		root := t.TempDir()
		path := filepath.Join(root, filepath.FromSlash(fixture.path))
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(fixture.text), 0600); err != nil {
			t.Fatal(err)
		}
		if err := developmentActionsCheckNodeFreeRepository(root); err == nil || err.Error() != "node-ban:"+fixture.code {
			t.Fatalf("expected %s, got %v", fixture.code, err)
		}
	}
	if err := developmentActionsCheckNodeFreeRepository(filepath.Join(t.TempDir(), "missing")); err == nil || err.Error() != "node-ban:read" {
		t.Fatalf("expected read rejection, got %v", err)
	}
	worktree := t.TempDir()
	if err := os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: unused"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(worktree, "package.json"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := developmentActionsCheckNodeFreeRepository(worktree); err == nil || err.Error() != "node-ban:artifact" {
		t.Fatal("git metadata file hid later artifact")
	}
}

func developmentActionsCheckHelper(setup *developmentActionsSetup, data []byte) error {
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != setup.HelperSHA256 {
		return fmt.Errorf("actions-setup:helper-digest")
	}
	for _, line := range []string{
		"  version=" + setup.GoVersion, "  expected=" + setup.GoSHA256,
		"  url=https://go.dev/dl/go" + setup.GoVersion + ".linux-amd64.tar.gz",
		"  version=" + setup.DenoVersion, "  expected=" + setup.DenoSHA256,
		"  url=https://github.com/denoland/deno/releases/download/v" + setup.DenoVersion + "/deno-x86_64-unknown-linux-gnu.zip",
	} {
		if !bytes.Contains(data, []byte("\n"+line+"\n")) {
			return fmt.Errorf("actions-setup:helper-definition")
		}
	}
	return nil
}

func TestDevelopmentActionsBootstrapControls(t *testing.T) {
	t.Parallel()
	setup := developmentActionsSetupDefinition(t)
	helper, err := filepath.Abs(".github/scripts/development-ci.sh")
	if err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{"go", "deno"} {
		for _, failure := range []string{"", "download-timeout-recover", "download-http-recover", "download-last-recover", "download-certificate", "download-not-found", "download-invalid-status", "download-unknown-exit", "output-boundary", "staging", "download", "digest-read", "digest", "archive-list", "archive-entry", "archive-type", "extract", "executable", "version", "version-exit", "publish", "platform"} {
			if kind == "deno" && failure == "archive-type" {
				continue
			}
			label := failure
			if label == "" {
				label = "success"
			}
			t.Run("release.bootstrap."+kind+"."+label+".outcome", func(t *testing.T) {
				base := t.TempDir()
				temp, fake := filepath.Join(base, "temp"), filepath.Join(base, "bin")
				for _, path := range []string{temp, fake} {
					if err := os.Mkdir(path, 0700); err != nil {
						t.Fatal(err)
					}
				}
				pathFile, envFile := filepath.Join(base, "path"), filepath.Join(base, "env")
				for _, path := range []string{pathFile, envFile} {
					if err := os.WriteFile(path, nil, 0600); err != nil {
						t.Fatal(err)
					}
				}
				if failure == "output-boundary" {
					if err := os.Remove(envFile); err != nil {
						t.Fatal(err)
					}
					if err := os.Mkdir(envFile, 0700); err != nil {
						t.Fatal(err)
					}
				}
				trace := filepath.Join(base, "trace")
				runtimeFile := filepath.Join(base, "runtime")
				phase18WriteExecutable(t, runtimeFile, "#!/bin/sh\n"+
					"printf 'version\\n' >> \"$FAKE_TRACE\"\n"+
					"if [ \"$FAKE_FAILURE\" = version-exit ]; then if [ \"$FAKE_KIND\" = go ]; then echo 'go version go1.22.12 linux/amd64'; else echo 'deno 2.9.7 (stable, release, x86_64-unknown-linux-gnu)'; fi; echo fixture-secret-canary >&2; exit 7; fi\n"+
					"if [ \"$FAKE_FAILURE\" = publish ]; then /bin/rm \"$GITHUB_PATH\"; mkdir \"$GITHUB_PATH\"; fi\n"+
					"if [ \"$FAKE_FAILURE\" = version ]; then echo invalid; elif [ \"$FAKE_KIND\" = go ]; then echo 'go version go1.22.12 linux/amd64'; else echo 'deno 2.9.7 (stable, release, x86_64-unknown-linux-gnu)'; fi\n")
				phase18WriteExecutable(t, filepath.Join(fake, "uname"), "#!/bin/sh\nif [ \"$FAKE_FAILURE\" = platform ] && [ \"$1\" = -s ]; then echo Linux; exit 7; fi\nif [ \"$1\" = -s ]; then echo Linux; else echo x86_64; fi\n")
				phase18WriteExecutable(t, filepath.Join(fake, "mktemp"), "#!/bin/sh\nif [ \"$FAKE_FAILURE\" = staging ]; then exit 7; fi\nexec /usr/bin/mktemp \"$@\"\n")
				phase18WriteExecutable(t, filepath.Join(fake, "rm"), "#!/bin/sh\nif [ \"$FAKE_FAILURE\" = cleanup-failure ]; then exit 7; fi\nexec /bin/rm \"$@\"\n")
				phase18WriteExecutable(t, filepath.Join(fake, "sleep"), "#!/bin/sh\n[ \"$#\" = 1 ] && [ \"$1\" = 1 ]\n")
				phase18WriteExecutable(t, filepath.Join(fake, "curl"), `#!/bin/sh
printf 'download\n' >> "$FAKE_TRACE"
case "$*" in
  "--disable --fail --silent --location --proto =https --proto-redir =https --tlsv1.2 --connect-timeout 15 --max-time 50 --retry 0 --max-filesize 200000000 --write-out %{http_code} --output "*) ;;
  *) exit 88 ;;
esac
while [ "$#" -gt 0 ]; do if [ "$1" = --output ]; then shift; output=$1; fi; shift; done
[ ! -e "$output" ] || exit 89
count=0
if [ -f "$FAKE_TRACE.count" ]; then read -r count < "$FAKE_TRACE.count"; fi
count=$((count + 1))
printf '%s\n' "$count" > "$FAKE_TRACE.count"
printf partial > "$output"
case "$FAKE_FAILURE" in
 download) echo fixture-secret-canary >&2; printf 000; exit 7 ;;
 download-timeout-recover) if [ "$count" = 1 ]; then printf 000; exit 28; fi ;;
 download-last-recover) if [ "$count" -lt 3 ]; then printf 000; exit 28; fi ;;
 download-http-recover) if [ "$count" = 1 ]; then printf 503; exit 22; fi ;;
 download-certificate) printf 000; exit 60 ;;
 download-not-found) printf 404; exit 22 ;;
 download-invalid-status) printf fixture-secret-canary; exit 0 ;;
 download-unknown-exit) printf 000; exit 88 ;;
esac
printf archive > "$output"
printf 200
`)
				phase18WriteExecutable(t, filepath.Join(fake, "sha256sum"), "#!/bin/sh\nprintf 'digest\\n' >> \"$FAKE_TRACE\"\nif [ \"$FAKE_FAILURE\" = digest-read ]; then echo fixture-secret-canary >&2; exit 7; fi\nif [ \"$FAKE_FAILURE\" = digest ]; then echo invalid; else printf '%s  %s\\n' \"$FAKE_DIGEST\" \"$2\"; fi\n")
				phase18WriteExecutable(t, filepath.Join(fake, "tar"), "#!/bin/sh\n"+
					"case \"$1\" in\n"+
					" -tzf) printf 'list\\n' >> \"$FAKE_TRACE\"; [ \"${3:-}\" = --quoting-style=literal ] || exit 8; if [ \"$FAKE_FAILURE\" = archive-list ]; then exit 7; fi; if [ \"$FAKE_FAILURE\" = archive-entry ]; then printf 'go/bin/go\\ngo/../escape\\n'; else printf 'go/\\ngo/bin/\\ngo/bin/go\\ngo/test/\\303\\236foo.go\\n'; fi;;\n"+
					" -tvzf) printf 'types\\n' >> \"$FAKE_TRACE\"; if [ \"$FAKE_FAILURE\" = archive-type ]; then echo 'lrwx go/bin/go'; else echo '-rwx go/bin/go'; fi;;\n"+
					" --extract) printf 'extract\\n' >> \"$FAKE_TRACE\"; if [ \"$FAKE_FAILURE\" = extract ]; then exit 7; fi; while [ \"$#\" -gt 0 ]; do if [ \"$1\" = --directory ]; then shift; dest=$1; fi; shift; done; mkdir -p \"$dest/go/bin\"; if [ \"$FAKE_FAILURE\" = executable ]; then ln -s \"$FAKE_RUNTIME\" \"$dest/go/bin/go\"; else cp \"$FAKE_RUNTIME\" \"$dest/go/bin/go\"; fi;;\n"+
					" *) exit 8;;\nesac\n")
				phase18WriteExecutable(t, filepath.Join(fake, "unzip"), "#!/bin/sh\n"+
					"if [ \"$1\" = -Z1 ]; then printf 'list\\n' >> \"$FAKE_TRACE\"; if [ \"$FAKE_FAILURE\" = archive-list ]; then exit 7; elif [ \"$FAKE_FAILURE\" = archive-entry ]; then printf 'deno\\n../escape\\n'; else echo deno; fi\n"+
					"elif [ \"$1\" = -q ]; then printf 'extract\\n' >> \"$FAKE_TRACE\"; if [ \"$FAKE_FAILURE\" = extract ]; then exit 7; elif [ \"$FAKE_FAILURE\" = executable ]; then ln -s \"$FAKE_RUNTIME\" \"$4/deno\"; else cp \"$FAKE_RUNTIME\" \"$4/deno\"; fi\nelse exit 8; fi\n")
				digest := setup.GoSHA256
				if kind == "deno" {
					digest = setup.DenoSHA256
				}
				env := append(os.Environ(), "PATH="+fake+":"+os.Getenv("PATH"), "RUNNER_TEMP="+temp, "GITHUB_PATH="+pathFile, "GITHUB_ENV="+envFile, "ADLAIRE_CI_TOOLS=", "FAKE_KIND="+kind, "FAKE_FAILURE="+failure, "FAKE_TRACE="+trace, "FAKE_RUNTIME="+runtimeFile, "FAKE_DIGEST="+digest)
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
				cmd := exec.CommandContext(ctx, "bash", "--noprofile", "--norc", helper, kind)
				cmd.Env = env
				output, err := cmd.CombinedOutput()
				pathData, pathErr := os.ReadFile(pathFile)
				envData, envErr := os.ReadFile(envFile)
				stages, globErr := filepath.Glob(filepath.Join(temp, "adlaire-tools.*"))
				if globErr != nil {
					t.Fatal(globErr)
				}
				if bytes.Contains(output, []byte("fixture-secret-canary")) {
					t.Fatal("bootstrap secret canary leaked")
				}
				traceData, traceErr := os.ReadFile(trace)
				positive := failure == "" || failure == "download-timeout-recover" || failure == "download-http-recover" || failure == "download-last-recover"
				attempts := 1
				if failure == "download" || failure == "download-last-recover" {
					attempts = 3
				} else if failure == "download-timeout-recover" || failure == "download-http-recover" {
					attempts = 2
				}
				if strings.Count(string(traceData), "download\n") != attempts && failure != "output-boundary" && failure != "staging" && failure != "platform" {
					t.Fatalf("download attempt bound: expected=%d trace=%q", attempts, traceData)
				}
				if positive {
					if err != nil || pathErr != nil || envErr != nil || len(stages) != 1 || !bytes.Contains(output, []byte(`"record_type":"development-actions-bootstrap"`)) {
						t.Fatalf("bootstrap positive: err=%v output=%q stages=%v", err, output, stages)
					}
					root := stages[0]
					if strings.Count(string(output), "ci-download:runtime=") != attempts || !bytes.Contains(output, []byte(fmt.Sprintf("attempt=%d curl_exit=0 http_status=200\n", attempts))) {
						t.Fatalf("safe download diagnostic missing: %q", output)
					}
					expectedBin := root + "/go/bin\n"
					expectedTrace := "download\ndigest\nlist\ntypes\nextract\nversion\n"
					if kind == "deno" {
						expectedBin = root + "/deno\n"
						expectedTrace = "download\ndigest\nlist\nextract\nversion\n"
					}
					expectedTrace = strings.Repeat("download\n", attempts-1) + expectedTrace
					if string(pathData) != expectedBin || traceErr != nil || string(traceData) != expectedTrace {
						t.Fatalf("bootstrap order/path expected %q/%q got %q/%q", expectedTrace, expectedBin, traceData, pathData)
					}
					if !strings.Contains(string(envData), "ADLAIRE_CI_TOOLS="+root+"\n") || !strings.Contains(string(envData), "GOTOOLCHAIN=local\n") || !strings.Contains(string(envData), "DENO_NO_PACKAGE_JSON=1\n") {
						t.Fatal("bootstrap environment missing")
					}
					cleanup := exec.CommandContext(ctx, "bash", "--noprofile", "--norc", helper, "cleanup")
					cleanup.Env = append(env, "ADLAIRE_CI_TOOLS="+root, "FAKE_FAILURE=cleanup-failure")
					if out, err := cleanup.CombinedOutput(); err == nil || !bytes.Contains(out, []byte("ci-setup:cleanup\n")) {
						t.Fatal("cleanup failure not propagated")
					}
					if _, err := os.Stat(root); err != nil {
						t.Fatal("failed cleanup must leave identifiable stage")
					}
					cleanup = exec.CommandContext(ctx, "bash", "--noprofile", "--norc", helper, "cleanup")
					cleanup.Env = append(env, "ADLAIRE_CI_TOOLS="+root)
					if out, err := cleanup.CombinedOutput(); err != nil || len(out) != 0 {
						t.Fatalf("cleanup err=%v output=%q", err, out)
					}
					if _, err := os.Lstat(root); !os.IsNotExist(err) {
						t.Fatal("bootstrap stage survived cleanup")
					}
				} else {
					reason := failure
					if strings.HasPrefix(reason, "download-") {
						reason = "download"
					}
					if reason == "version-exit" {
						reason = "version"
					}
					if err == nil || !bytes.Contains(output, []byte("ci-setup:"+reason+"\n")) || bytes.Contains(output, []byte("record_type")) || len(stages) != 0 || len(envData) != 0 || len(pathData) != 0 {
						t.Fatalf("bootstrap expected=%s err=%v output=%q stages=%v env=%q path=%q", reason, err, output, stages, envData, pathData)
					}
					if reason == "download" && (strings.Contains(string(traceData), "digest\n") || strings.Count(string(output), "ci-download:runtime=") != attempts) {
						t.Fatal("download failure continued or lost attempt diagnostic")
					}
					for boundary, forbidden := range map[string]string{"download": "digest", "digest-read": "list", "digest": "list", "archive-list": "extract", "archive-entry": "extract", "archive-type": "extract", "extract": "version", "executable": "version"} {
						if failure == boundary && strings.Contains(string(traceData), forbidden+"\n") {
							t.Fatalf("bootstrap continued after %s: %s", boundary, traceData)
						}
					}
				}
			})
		}
	}
	t.Run("release.bootstrap.cleanup-symlink.reject", func(t *testing.T) {
		temp, target := t.TempDir(), t.TempDir()
		link := filepath.Join(temp, "adlaire-tools.ABCDEFGH")
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command("bash", "--noprofile", "--norc", helper, "cleanup")
		cmd.Env = append(os.Environ(), "RUNNER_TEMP="+temp, "ADLAIRE_CI_TOOLS="+link)
		output, err := cmd.CombinedOutput()
		if err == nil || !bytes.Contains(output, []byte("ci-setup:cleanup-boundary\n")) {
			t.Fatal("cleanup symlink not rejected")
		}
		if info, err := os.Stat(target); err != nil || !info.IsDir() {
			t.Fatal("cleanup touched symlink target")
		}
	})
	t.Run("release.bootstrap.arguments.reject", func(t *testing.T) {
		cmd := exec.Command("bash", "--noprofile", "--norc", helper)
		if output, err := cmd.CombinedOutput(); err == nil || !bytes.Contains(output, []byte("ci-setup:arguments\n")) {
			t.Fatal("missing argument not rejected")
		}
	})
	for _, fixture := range []struct{ name, argument, temp, root, reason string }{
		{"runtime", "other", "", "", "runtime"},
		{"relative-temp", "go", "relative", "", "temp-boundary"},
		{"missing-temp", "go", "/missing-node-free-fixture", "", "temp-boundary"},
		{"cleanup-outside", "cleanup", "", "/outside", "cleanup-boundary"},
	} {
		t.Run("release.bootstrap."+fixture.name+".reject", func(t *testing.T) {
			temp := fixture.temp
			if temp == "" {
				temp = t.TempDir()
			}
			cmd := exec.Command("bash", "--noprofile", "--norc", helper, fixture.argument)
			cmd.Env = append(os.Environ(), "RUNNER_TEMP="+temp, "ADLAIRE_CI_TOOLS="+fixture.root)
			output, err := cmd.CombinedOutput()
			if err == nil || !bytes.Contains(output, []byte("ci-setup:"+fixture.reason+"\n")) {
				t.Fatalf("boundary expected=%s err=%v output=%q", fixture.reason, err, output)
			}
		})
	}
	for _, failure := range []string{"", "init", "remote", "fetch", "checkout", "sha", "dirty", "bad-sha", "repository"} {
		label := failure
		if label == "" {
			label = "success"
		}
		t.Run("release.checkout."+label+".outcome", func(t *testing.T) {
			base := t.TempDir()
			workspace, fake := filepath.Join(base, "workspace"), filepath.Join(base, "bin")
			for _, path := range []string{workspace, fake} {
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
			}
			trace := filepath.Join(base, "trace")
			target := strings.Repeat("a", 40)
			if failure == "dirty" {
				if err := os.WriteFile(filepath.Join(workspace, "existing"), nil, 0600); err != nil {
					t.Fatal(err)
				}
			}
			phase18WriteExecutable(t, filepath.Join(fake, "git"), "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$FAKE_TRACE\"\ncase \"$*\" in\n*' init '*) stage=init;; *' remote '*) stage=remote;; *' fetch '*) stage=fetch;; *' checkout '*) stage=checkout;; 'rev-parse HEAD') if [ \"$FAKE_FAILURE\" = sha ]; then echo invalid; else echo \"$GITHUB_SHA\"; fi; exit 0;; *) exit 8;; esac\n[ \"$FAKE_FAILURE\" != \"$stage\" ] || exit 7\n")
			repository := "fqwink/Build-Scripts"
			if failure == "repository" {
				repository = "unknown/repo"
			}
			if failure == "bad-sha" {
				target = strings.Repeat("A", 40)
			}
			cmd := exec.Command("bash", "--noprofile", "--norc", "-eo", "pipefail", "-c", setup.Checkout)
			cmd.Dir = workspace
			cmd.Env = append(os.Environ(), "PATH="+fake+":"+os.Getenv("PATH"), "GITHUB_REPOSITORY="+repository, "GITHUB_SERVER_URL=https://github.com", "GITHUB_SHA="+target, "GITHUB_WORKSPACE="+workspace, "FAKE_TRACE="+trace, "FAKE_FAILURE="+failure)
			output, err := cmd.CombinedOutput()
			if (err == nil) != (failure == "") {
				t.Fatalf("checkout expected success=%t err=%v output=%q", failure == "", err, output)
			}
			traceData, _ := os.ReadFile(trace)
			count := len(strings.Split(strings.TrimSpace(string(traceData)), "\n"))
			if len(traceData) == 0 {
				count = 0
			}
			expected := map[string]int{"": 5, "init": 1, "remote": 2, "fetch": 3, "checkout": 4, "sha": 5, "dirty": 0, "bad-sha": 0, "repository": 0}[failure]
			if count != expected {
				t.Fatalf("checkout stopped at wrong boundary: want %d got %d trace=%q", expected, count, traceData)
			}
		})
	}
}

func phase11MustReadText(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.HasSuffix(path, ".md") {
		return string(data)
	}
	if strings.Contains(string(data), "\x00") {
		t.Fatalf("%s must not contain NUL bytes", path)
	}
	return string(data)
}
