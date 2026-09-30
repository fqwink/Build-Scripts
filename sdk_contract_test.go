package main

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestSDKPublicMethodsMatchAPIContract(t *testing.T) {
	api := readTextFile(t, "docs/details/api.md")
	sdk := readTextFile(t, "admin/adlaire-ci-sdk.js")

	expected := sdkMethodsFromAPIContract(t, api)
	actual := sortedStrings(publicSDKMethodOrder(t, sdk))

	if missing := difference(expected, actual); len(missing) > 0 {
		t.Fatalf("SDK is missing API contract methods: %s", strings.Join(missing, ", "))
	}
	if extra := difference(actual, expected); len(extra) > 0 {
		t.Fatalf("SDK has public methods outside API contract: %s", strings.Join(extra, ", "))
	}
}

func TestSDKPublicMethodOrderMatchesSDKSpec(t *testing.T) {
	spec := readTextFile(t, "docs/details/sdk.md")
	sdk := readTextFile(t, "admin/adlaire-ci-sdk.js")

	expected := sdkMethodOrderFromSDKSpec(t, spec)
	actual := publicSDKMethodOrder(t, sdk)

	if strings.Join(expected, "\n") != strings.Join(actual, "\n") {
		t.Fatalf("SDK public method order differs from docs/details/sdk.md\nexpected:\n%s\nactual:\n%s", strings.Join(expected, "\n"), strings.Join(actual, "\n"))
	}
}

func TestSDKPrivateClassMethodsAreFixed(t *testing.T) {
	sdk := readTextFile(t, "admin/adlaire-ci-sdk.js")
	expected := []string{"_request", "_json", "_query", "_requireToken", "_validateId", "_clearTokenOn401"}
	actual := privateSDKMethodOrder(t, sdk)
	if strings.Join(expected, "\n") != strings.Join(actual, "\n") {
		t.Fatalf("SDK private class methods differ from fixed contract\nexpected:\n%s\nactual:\n%s", strings.Join(expected, "\n"), strings.Join(actual, "\n"))
	}
}

func TestSDKContractGuards(t *testing.T) {
	sdk := readTextFile(t, "admin/adlaire-ci-sdk.js")

	forbidden := []string{
		"URLSearchParams",
		"response.json(",
		"EventSource",
		"localStorage",
		"sessionStorage",
		"window.AdlaireCI",
		"window.AdlaireCIError",
		"globalThis.AdlaireCI",
		"globalThis.AdlaireCIError",
	}
	for _, pattern := range forbidden {
		if strings.Contains(sdk, pattern) {
			t.Fatalf("SDK must not contain %q", pattern)
		}
	}

	required := []string{
		"normalizeBaseUrl(options.baseUrl)",
		"encodeURIComponent(String(value))",
		"failure_category: failureCategory",
		"body: { label, scopes, expires_at: expiresAt }",
		"requireNumber(page, \"page\")",
		"requireNumber(limit, \"limit\")",
		"requireNumber(offset, \"offset\")",
		"requireOptionalString(failureCategory, \"failureCategory\")",
		"typeof flagged !== \"boolean\"",
		"requireStringArray(scopes, \"scopes\", false)",
		"requireStringArray(tags, \"tags\", true)",
		"requireStringArray(allowList, \"allowList\", true)",
		"requireString(owner, \"owner\")",
		"requireString(repo, \"repo\")",
		"/^[A-Za-z0-9_-]{1,64}$/.test(id)",
		"Invalid binary response",
		"Invalid SSE response",
		"Invalid SSE frame",
		"Stream callback failed",
		"new TextDecoder(\"utf-8\", { fatal: true })",
		"done,",
	}
	for _, pattern := range required {
		if !strings.Contains(sdk, pattern) {
			t.Fatalf("SDK contract guard missing %q", pattern)
		}
	}
}

func TestSDKRequestShapeGuardSnippets(t *testing.T) {
	sdk := readTextFile(t, "admin/adlaire-ci-sdk.js")

	required := []string{
		`this._request("/api/login", { method: "POST", body: { password } })`,
		`this._request("/api/login/totp", { method: "POST", body: { ticket, code } })`,
		`this._request("/api/change-password", {`,
		`body: { current_password: currentPassword, new_password: newPassword }`,
		`this._request("/api/logs", { query: this._query({ n, q }, ["q"]) })`,
		`this._query({ page, per_page: perPage, trigger, tag, flagged, failure_category: failureCategory })`,
		`this._request("/api/access-log", { query: this._query({ limit, offset }) })`,
		`this._request("/api/notify-log", { query: this._query({ limit, offset }) })`,
		`this._request("/api/config-log", { query: this._query({ limit, offset }) })`,
		`this._query({ q, from, to, level }, ["q", "from", "to"])`,
		`this._request("/api/schedule/allowed-hours", { method: "POST", body: { from: null, to: null } })`,
		`this._request("/api/branch-config", { method: "POST", body: { branches } })`,
		`this._request("/api/build-chain-config", { method: "POST", body: { chains } })`,
		"this._request(`/api/approvals/${this._validateId(id)}/approve`, { method: \"POST\" })",
		"this._request(`/api/approvals/${this._validateId(id)}/reject`, { method: \"POST\" })",
		`this._request("/api/repo-config", { method: "POST", body: patch })`,
		"this._request(`/api/snapshots/${this._validateId(id)}/download`, { binary: true })",
		`this._request("/api/smtp-config", { method: "POST", body })`,
	}
	for _, pattern := range required {
		if !strings.Contains(sdk, pattern) {
			t.Fatalf("SDK request shape guard missing %q", pattern)
		}
	}
}

func sdkMethodsFromAPIContract(t *testing.T, api string) []string {
	t.Helper()

	start := strings.Index(api, "<a id=\"sec-22-0e\"></a>")
	if start < 0 {
		t.Fatal("api contract section anchor not found")
	}
	end := strings.Index(api[start:], "**archive / snapshot / rollback")
	if end < 0 {
		t.Fatal("api contract table end not found")
	}
	section := api[start : start+end]

	methods := map[string]struct{}{}
	methodRe := regexp.MustCompile("`([A-Za-z][A-Za-z0-9_]*)\\s*\\(")
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 8 {
			continue
		}
		sdkCell := strings.TrimSpace(cells[6])
		if sdkCell == "none" {
			continue
		}
		for _, match := range methodRe.FindAllStringSubmatch(sdkCell, -1) {
			methods[match[1]] = struct{}{}
		}
	}
	return sortedKeys(methods)
}

func sdkMethodOrderFromSDKSpec(t *testing.T, spec string) []string {
	t.Helper()

	start := strings.Index(spec, "```js\nclass AdlaireCI {")
	if start < 0 {
		t.Fatal("SDK spec class code block not found")
	}
	end := strings.Index(spec[start:], "\nexport { AdlaireCI, AdlaireCIError };")
	if end < 0 {
		t.Fatal("SDK spec class code block end not found")
	}
	block := spec[start : start+end]
	return methodOrderFromBlock(block)
}

func publicSDKMethodOrder(t *testing.T, sdk string) []string {
	t.Helper()

	start := strings.Index(sdk, "class AdlaireCI {")
	if start < 0 {
		t.Fatal("SDK class not found")
	}
	end := strings.Index(sdk[start:], "\n}\n\nfunction runtimeSupported()")
	if end < 0 {
		t.Fatal("SDK class end not found")
	}
	block := sdk[start : start+end]
	return methodOrderFromBlock(block)
}

func privateSDKMethodOrder(t *testing.T, sdk string) []string {
	t.Helper()

	start := strings.Index(sdk, "class AdlaireCI {")
	if start < 0 {
		t.Fatal("SDK class not found")
	}
	end := strings.Index(sdk[start:], "\n}\n\nfunction runtimeSupported()")
	if end < 0 {
		t.Fatal("SDK class end not found")
	}
	block := sdk[start : start+end]
	var methods []string
	methodRe := regexp.MustCompile(`(?m)^  (?:async\s+)?(_[A-Za-z][A-Za-z0-9_]*)\(`)
	for _, match := range methodRe.FindAllStringSubmatch(block, -1) {
		methods = append(methods, match[1])
	}
	return methods
}

func methodOrderFromBlock(block string) []string {
	var methods []string
	methodRe := regexp.MustCompile(`(?m)^  (?:async\s+)?([A-Za-z][A-Za-z0-9_]*)\(`)
	for _, match := range methodRe.FindAllStringSubmatch(block, -1) {
		name := match[1]
		if name == "constructor" {
			continue
		}
		methods = append(methods, name)
	}
	return methods
}

func readTextFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func sortedKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedStrings(values []string) []string {
	copied := append([]string(nil), values...)
	sort.Strings(copied)
	return copied
}

func difference(left, right []string) []string {
	rightSet := make(map[string]struct{}, len(right))
	for _, value := range right {
		rightSet[value] = struct{}{}
	}
	var diff []string
	for _, value := range left {
		if _, ok := rightSet[value]; !ok {
			diff = append(diff, value)
		}
	}
	return diff
}
