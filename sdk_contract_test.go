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
	actual := publicSDKMethods(t, sdk)

	if missing := difference(expected, actual); len(missing) > 0 {
		t.Fatalf("SDK is missing API contract methods: %s", strings.Join(missing, ", "))
	}
	if extra := difference(actual, expected); len(extra) > 0 {
		t.Fatalf("SDK has public methods outside API contract: %s", strings.Join(extra, ", "))
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
		"typeof flagged !== \"boolean\"",
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

func publicSDKMethods(t *testing.T, sdk string) []string {
	t.Helper()

	methods := map[string]struct{}{}
	methodRe := regexp.MustCompile(`(?m)^  (?:async\s+)?([A-Za-z][A-Za-z0-9_]*)\(`)
	for _, match := range methodRe.FindAllStringSubmatch(sdk, -1) {
		name := match[1]
		if name == "constructor" || strings.HasPrefix(name, "_") {
			continue
		}
		methods[name] = struct{}{}
	}
	return sortedKeys(methods)
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
