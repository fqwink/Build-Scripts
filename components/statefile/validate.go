package statefile

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var cicdIDPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)
var cicdEnvPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,63}$`)

// ValidateCICDValue checks storage shape and primitive grammar, not owner policy.
func ValidateCICDValue(typeName string, data []byte) error {
	value, err := decodeCICDValue(data)
	if err != nil {
		return ErrCICDSchema
	}
	if _, exists := cicdStorageTypes[typeName]; !exists || !validateCICDType(typeName, value, "", 0) {
		return ErrCICDSchema
	}
	return nil
}

func decodeCICDValue(data []byte) (any, error) {
	if len(data) > cicdSnapshotLimit || !validateCICDJSONBytes(data) {
		return nil, ErrCICDSchema
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := readCICDJSONValue(decoder)
	if err != nil {
		return nil, ErrCICDSchema
	}
	if _, err := decoder.Token(); err != io.EOF {
		return nil, ErrCICDSchema
	}
	return value, nil
}

func validateCICDType(expression string, value any, context string, depth int) bool {
	if depth > 64 {
		return false
	}
	if definition, exists := cicdStorageTypes[expression]; exists {
		switch schema := definition.(type) {
		case string:
			return validateCICDType(schema, value, context, depth+1)
		case map[string]string:
			object, ok := value.(map[string]any)
			if !ok || len(object) != len(schema) {
				return false
			}
			for key, childType := range schema {
				child, exists := object[key]
				if !exists || !validateCICDType(childType, child, expression+"."+key, depth+1) {
					return false
				}
			}
			if expression == "Resources" && cicdInt(object["memory_mb"]) < cicdInt(object["disk_mb"])+128 {
				return false
			}
			if expression == "Job" || expression == "PlannedJob" {
				total := 0
				for _, text := range cicdObject(object["env"]) {
					total += len(text.(string))
				}
				if total > 32768 {
					return false
				}
			}
			return true
		}
		return false
	}
	if index := strings.IndexByte(expression, '('); index >= 0 {
		if !strings.HasSuffix(expression, ")") {
			return false
		}
		name, arguments := expression[:index], splitCICDArguments(expression[index+1:len(expression)-1])
		switch name {
		case "nullable":
			return len(arguments) == 1 && (value == nil || validateCICDType(arguments[0], value, context, depth+1))
		case "oneof":
			matches := 0
			for _, option := range arguments {
				if validateCICDType(option, value, context, depth+1) {
					matches++
				}
			}
			return len(arguments) >= 2 && matches == 1
		case "enum":
			text, ok := value.(string)
			for _, option := range arguments {
				if ok && text == option {
					return true
				}
			}
			return false
		case "const":
			if len(arguments) != 1 {
				return false
			}
			return (arguments[0] == "false" && value == false) || (arguments[0] == "1" && value == json.Number("1"))
		case "int", "text", "list", "map", "envmap":
			offset := 0
			if name == "list" || name == "map" || name == "envmap" {
				offset = 1
			}
			if len(arguments) != offset+2 {
				return false
			}
			min, e1 := strconv.ParseInt(arguments[offset], 10, 64)
			max, e2 := strconv.ParseInt(arguments[offset+1], 10, 64)
			if e1 != nil || e2 != nil || min > max {
				return false
			}
			inRange := func(n int64) bool { return n >= min && n <= max }
			switch name {
			case "int":
				number, ok := value.(json.Number)
				n, err := number.Int64()
				return ok && err == nil && validateCICDInteger(number) && inRange(n)
			case "text":
				text, ok := value.(string)
				return ok && inRange(int64(len(text))) && (context == "Log.data" || !strings.ContainsAny(text, "\x00\r")) && (context != "Job.env" && context != "PlannedJob.env" || !strings.ContainsRune(text, '\n'))
			case "list":
				array, ok := value.([]any)
				if !ok || !inRange(int64(len(array))) {
					return false
				}
				for _, child := range array {
					if !validateCICDType(arguments[0], child, context, depth+1) {
						return false
					}
				}
				return true
			default:
				object, ok := value.(map[string]any)
				if !ok || !inRange(int64(len(object))) {
					return false
				}
				for key, child := range object {
					if name == "envmap" && !validCICDEnvKey(key) || name == "map" && !cicdIDPattern.MatchString(key) || !validateCICDType(arguments[0], child, context, depth+1) {
						return false
					}
				}
				return true
			}
		case "legacy":
			return len(arguments) == 1 && validateCICDLegacy(arguments[0], value, depth+1)
		}
		return false
	}
	if expression == "bool" {
		_, ok := value.(bool)
		return ok
	}
	if expression == "scalar" {
		switch scalar := value.(type) {
		case bool:
			return true
		case json.Number:
			return validateCICDType("int(-2147483648,2147483647)", scalar, context, depth+1)
		case string:
			return len(scalar) <= 4096 && !strings.ContainsAny(scalar, "\x00\r\n")
		}
		return false
	}
	text, ok := value.(string)
	if !ok || strings.ContainsAny(text, "\x00\r\n") {
		return false
	}
	switch expression {
	case "id", "jobid":
		return cicdIDPattern.MatchString(text) && (expression != "jobid" || len(text) <= 48)
	case "digest":
		return cicdLowerHex(text, 64)
	case "sha":
		return cicdLowerHex(text, 40)
	case "utc":
		parsed, err := time.Parse("2006-01-02T15:04:05Z", text)
		return err == nil && len(text) == 20 && parsed.Year() >= 1970 && parsed.Year() <= 9999 && parsed.Format("2006-01-02T15:04:05Z") == text
	case "base64":
		decoded, err := base64.StdEncoding.Strict().DecodeString(text)
		return err == nil && len(decoded) <= 262144 && base64.StdEncoding.EncodeToString(decoded) == text
	case "url":
		parsed, err := url.Parse(text)
		return err == nil && len(text) <= 2048 && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Hostname() != "" && parsed.User == nil && parsed.Fragment == "" && !strings.Contains(text, "#")
	case "ref":
		if text == "" || len(text) > 255 || text == "@" || strings.Contains(text, "..") || strings.Contains(text, "@{") || strings.ContainsAny(text, " ~^:?*[\\") || cicdHasControl(text) {
			return false
		}
		for _, part := range strings.Split(text, "/") {
			if part == "" || strings.HasPrefix(part, ".") || strings.HasSuffix(part, ".") || strings.HasSuffix(part, ".lock") {
				return false
			}
		}
		return true
	case "relpath":
		if len(text) > 240 {
			return false
		}
		packageField := context == "GoTestTask.packages" || context == "GoVetTask.packages" || context == "GoBuildTask.package"
		if text == "." && (packageField || context == "BuilderTask.source") || packageField && text == "./..." {
			return true
		}
		if packageField && strings.HasSuffix(text, "/...") {
			text = strings.TrimSuffix(text, "/...")
		}
		return validCICDRelativePath(text)
	case "abspath":
		if !filepath.IsAbs(text) || filepath.Clean(text) != text || len(text) > 4096 || strings.Contains(text, "\\") || cicdHasControl(text) {
			return false
		}
		for path := text; ; path = filepath.Dir(path) {
			info, err := os.Lstat(path)
			if err != nil && !errors.Is(err, os.ErrNotExist) || err == nil && info.Mode()&os.ModeSymlink != 0 {
				return false
			}
			if path == "/" {
				return true
			}
		}
	case "glob":
		if text == "" || len(text) > 240 || cicdHasControl(text) || strings.ContainsAny(text, "\\![]{}:") {
			return false
		}
		parts := strings.Split(text, "/")
		if len(parts) > 32 {
			return false
		}
		for _, part := range parts {
			if part == "" || part == "." || part == ".." || strings.Contains(part, "**") && part != "**" {
				return false
			}
		}
		return true
	}
	return false
}

func splitCICDArguments(text string) []string {
	var result []string
	depth, start := 0, 0
	for index, char := range text {
		switch char {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				result = append(result, text[start:index])
				start = index + 1
			}
		}
	}
	return append(result, text[start:])
}

func cicdLowerHex(text string, length int) bool {
	if len(text) != length {
		return false
	}
	for _, char := range text {
		if !(char >= '0' && char <= '9' || char >= 'a' && char <= 'f') {
			return false
		}
	}
	return true
}

func cicdHasControl(text string) bool {
	for _, char := range text {
		if char < 32 || char == 127 {
			return true
		}
	}
	return false
}

func validCICDRelativePath(text string) bool {
	if text == "" || len(text) > 240 || strings.ContainsAny(text, "\\:") || cicdHasControl(text) {
		return false
	}
	for _, part := range strings.Split(text, "/") {
		if part == "" || len(part) > 64 || part == "." || part == ".." {
			return false
		}
	}
	return true
}

func validCICDEnvKey(key string) bool {
	if !cicdEnvPattern.MatchString(key) {
		return false
	}
	for _, prefix := range []string{"LD_", "DYLD_", "GO", "DENO_", "NODE_", "NPM_"} {
		if strings.HasPrefix(key, prefix) {
			return false
		}
	}
	return key != "PATH" && key != "HOME" && key != "TMPDIR" && key != "LANG" && key != "TZ"
}

func validateCICDLegacy(name string, value any, depth int) bool {
	object, ok := value.(map[string]any)
	if !ok {
		return false
	}
	var fields map[string]string
	if name == "AdminCredentials" {
		fields = map[string]string{"password_hash": "Digest", "salt": "Digest", "algorithm": "enum(pbkdf2_hmac_sha256_v1)", "iterations": "int(600000,600000)", "must_change": "bool", "login_count": "int(0,9007199254740991)", "last_login_at": "nullable(Time)", "updated_at": "Time"}
	} else if name == "TOTPSecret" {
		fields = map[string]string{"enabled": "bool", "secret_base32": "nullable(text(16,64))", "confirmed_at": "nullable(Time)", "last_accepted_step": "nullable(int(0,9007199254740991))"}
	} else {
		return false
	}
	if len(object) != len(fields) {
		return false
	}
	for key, schema := range fields {
		child, exists := object[key]
		if !exists || !validateCICDType(schema, child, "", depth+1) {
			return false
		}
	}
	if name == "TOTPSecret" && object["secret_base32"] != nil {
		secret := object["secret_base32"].(string)
		decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
		if err != nil || base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(decoded) != secret {
			return false
		}
	}
	return true
}

func cicdInt(value any) int64 {
	number, _ := value.(json.Number)
	result, _ := number.Int64()
	return result
}

func cicdObject(value any) map[string]any { return value.(map[string]any) }
func cicdList(value any) []any            { return value.([]any) }

func cicdKey(object map[string]any, fields ...string) string {
	parts := make([]string, len(fields))
	for i, field := range fields {
		parts[i] = fmt.Sprint(object[field])
	}
	return strings.Join(parts, "\x00")
}

func cicdTerminal(state any) bool {
	return state == "succeeded" || state == "failed" || state == "cancelled" || state == "skipped" || state == "rolled_back"
}

func validateCICDRoot(root map[string]any) bool {
	if !validateCICDType("Root", root, "", 0) {
		return false
	}
	keys := map[string][]string{
		"projects": {"id"}, "definitions": {"project_id", "id", "revision"}, "templates": {"project_id", "id", "revision"}, "environments": {"project_id", "id"}, "users": {"id"}, "services": {"id"}, "secrets": {"project_id", "id", "version"}, "grants": {"subject_kind", "subject_id", "project_id", "environment_id", "role"}, "workers": {"id"}, "runs": {"project_id", "id", "attempt"}, "deployments": {"id"}, "artifacts": {"project_id", "digest"}, "caches": {"project_id", "key"}, "operations": {"id"}, "targets": {"id"}, "environment_states": {"project_id", "environment_id"}, "notifications": {"id"}, "trigger_cursors": {"project_id", "definition_id", "definition_revision", "kind", "trigger_id"}, "status_deliveries": {"project_id", "run_id", "attempt", "state"},
	}
	indexes := make(map[string]map[string]map[string]any)
	for collection, fields := range keys {
		index := make(map[string]map[string]any)
		for _, value := range cicdList(root[collection]) {
			object := cicdObject(value)
			key := cicdKey(object, fields...)
			if index[key] != nil {
				return false
			}
			index[key] = object
		}
		indexes[collection] = index
	}
	for collection := range keys {
		for _, object := range indexes[collection] {
			if project := object["project_id"]; project != nil && indexes["projects"][project.(string)] == nil {
				return false
			}
		}
	}
	for id, user := range indexes["users"] {
		if id == "system" || id == "worker" || id == "site" || indexes["services"][id] != nil || id != "admin" && user["credential"] == nil || id != "admin" && user["totp"] != nil {
			return false
		}
	}
	for id := range indexes["services"] {
		if id == "system" || id == "worker" || id == "site" || id == "admin" {
			return false
		}
	}
	admin := false
	for _, grant := range indexes["grants"] {
		kind, id, role := grant["subject_kind"].(string), grant["subject_id"].(string), grant["role"]
		collection := "users"
		if kind == "service" {
			collection = "services"
		}
		subject := indexes[collection][id]
		if subject == nil {
			return false
		}
		if role == "platform-admin" {
			if kind != "user" || grant["project_id"] != nil || grant["environment_id"] != nil {
				return false
			}
			if grant["enabled"] == true && subject["enabled"] == true {
				admin = true
			}
		} else {
			if grant["project_id"] == nil || (role == "approver") != (grant["environment_id"] != nil) {
				return false
			}
			if grant["environment_id"] != nil && indexes["environments"][cicdKey(grant, "project_id", "environment_id")] == nil {
				return false
			}
		}
	}
	if !admin {
		return false
	}
	for _, project := range indexes["projects"] {
		if (project["definition_id"] == nil) != (cicdInt(project["definition_revision"]) == 0) || project["enabled"] == true && project["toolchain_digest"] == nil {
			return false
		}
		if project["definition_id"] != nil {
			key := fmt.Sprint(project["id"]) + "\x00" + fmt.Sprint(project["definition_id"]) + "\x00" + fmt.Sprint(project["definition_revision"])
			if indexes["definitions"][key] == nil {
				return false
			}
		}
	}
	activeRuns, runningJobs := 0, 0
	for _, run := range indexes["runs"] {
		if !validCICDEntityID(run["id"], "r-") || !validateCICDRun(run) {
			return false
		}
		if run["requester"] == "system" {
			if run["trigger"] != "watch" && run["trigger"] != "schedule" {
				return false
			}
		} else if indexes["users"][fmt.Sprint(run["requester"])] == nil && indexes["services"][fmt.Sprint(run["requester"])] == nil {
			return false
		}
		if !cicdTerminal(run["state"]) {
			activeRuns++
		}
		for _, value := range cicdList(run["jobs"]) {
			job := cicdObject(value)
			if job["worker_id"] != nil && indexes["workers"][job["worker_id"].(string)] == nil {
				return false
			}
			if job["execution_id"] != nil && !validCICDEntityID(job["execution_id"], "o-") {
				return false
			}
			if job["state"] == "assigned" || job["state"] == "running" || job["state"] == "cancel_requested" || job["state"] == "recovery_required" {
				runningJobs++
			}
		}
		for attempt := int64(1); attempt < cicdInt(run["attempt"]); attempt++ {
			previous := indexes["runs"][fmt.Sprint(run["project_id"])+"\x00"+fmt.Sprint(run["id"])+"\x00"+strconv.FormatInt(attempt, 10)]
			if previous == nil || !cicdTerminal(previous["state"]) {
				return false
			}
		}
	}
	if activeRuns > 16 || runningJobs > 64 {
		return false
	}
	for _, delivery := range indexes["status_deliveries"] {
		if indexes["runs"][cicdKey(delivery, "project_id", "run_id", "attempt")] == nil || delivery["sent"] == true && delivery["code"] != nil {
			return false
		}
	}
	for _, cursor := range indexes["trigger_cursors"] {
		definition := indexes["definitions"][cicdKey(cursor, "project_id", "definition_id", "definition_revision")]
		if definition == nil || cursor["kind"] == "watch" && cursor["last_due"] != nil || cursor["kind"] == "schedule" && cursor["last_sha"] != nil {
			return false
		}
		field := "watches"
		if cursor["kind"] == "schedule" {
			field = "schedules"
		}
		found := false
		for _, value := range cicdList(definition[field]) {
			if cicdObject(value)["id"] == cursor["trigger_id"] {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	for _, operation := range indexes["operations"] {
		if !validCICDEntityID(operation["id"], "o-") || (operation["state"] == "accepted" || operation["state"] == "succeeded") != (operation["code"] == nil) {
			return false
		}
		created, _ := time.Parse(time.RFC3339, operation["created_at"].(string))
		if created.Add(2592000*time.Second).Format(time.RFC3339) != operation["expires_at"] {
			return false
		}
	}
	for _, value := range cicdList(root["audits"]) {
		audit := cicdObject(value)
		if operation := indexes["operations"][fmt.Sprint(audit["operation_id"])]; operation != nil {
			for _, key := range []string{"subject", "project_id", "target_id"} {
				if audit[key] != operation[key] {
					return false
				}
			}
		}
	}
	return validateCICDOperations(root)
}

func validCICDEntityID(value any, prefix string) bool {
	text, ok := value.(string)
	return ok && strings.HasPrefix(text, prefix) && cicdLowerHex(strings.TrimPrefix(text, prefix), 32)
}

func cicdOperationCode(code string) bool {
	for _, known := range strings.Fields("P19_JSON P19_INPUT P19_AUTH P19_FORBIDDEN P19_NOT_FOUND P19_METHOD P19_REVISION P19_STALE P19_IDEMPOTENCY P19_BUSY P19_APPROVAL P19_FREEZE P19_CURSOR P19_BODY P19_MEDIA P19_SCHEMA P19_DAG P19_MATRIX P19_EMPTY P19_CAPABILITY P19_SOURCE P19_DIGEST P19_ARCHIVE P19_SIGNATURE P19_TRUST P19_QUOTA P19_RATE P19_INTERNAL P19_ENTROPY P19_SECRET P19_EXIT P19_OUTPUT P19_TIMEOUT P19_CLEANUP P19_HEALTH P19_CANCELLED P19_TRANSPORT P19_LEASE P19_STATE_CORRUPT P19_STATE_UNCERTAIN P19_SHUTDOWN P19_STORAGE") {
		if code == known {
			return true
		}
	}
	return false
}

func validCICDMutationReply(kind string, result cicdOperationMutation) bool {
	if kind != "artifactPin" && kind != "cacheDelete" && !cicdIDPattern.MatchString(result.targetID) {
		return false
	}
	if kind == "artifactPin" || kind == "cacheDelete" {
		if !cicdLowerHex(result.targetID, 64) {
			return false
		}
	}
	if strings.HasPrefix(kind, "run") {
		if !validCICDEntityID(result.targetID, "r-") || result.targetAttempt < 1 || result.targetAttempt > 100 {
			return false
		}
	} else if result.targetAttempt != 0 {
		return false
	}
	if strings.HasPrefix(kind, "deployment") && !validCICDEntityID(result.targetID, "d-") {
		return false
	}
	if len(result.result) > 921600 {
		return false
	}
	if len(result.result) > 0 {
		if _, err := decodeCICDValue(result.result); err != nil {
			return false
		}
	}
	return utf8.ValidString(result.oneTime) && len(result.oneTime) <= 4096 && !strings.ContainsAny(result.oneTime, "\x00\r\n") && (result.oneTime == "" || kind == "userPut" || kind == "servicePut")
}

func validCICDOperationTransition(before, after map[string]any) bool {
	data, err := json.Marshal(after)
	if err != nil || ValidateCICDValue("Operation", data) != nil {
		return false
	}
	for _, field := range []string{"id", "project_id", "subject", "kind", "request_digest", "target_id", "target_attempt", "created_at", "expires_at"} {
		if before[field] != after[field] {
			return false
		}
	}
	state, code := after["state"], after["code"]
	if (state == "accepted" || state == "succeeded") != (code == nil) {
		return false
	}
	if code != nil && !cicdOperationCode(code.(string)) {
		return false
	}
	if cicdTerminal(before["state"]) {
		return before["state"] == state && before["code"] == code
	}
	if before["state"] == "accepted" {
		return state == "accepted" || state == "succeeded" || state == "failed" || state == "recovery_required"
	}
	return before["state"] == "recovery_required" && (state == "recovery_required" || state == "succeeded" || state == "failed")
}

func validateCICDOperations(root map[string]any) bool {
	auditsByID := make(map[string][]map[string]any)
	for _, value := range cicdList(root["audits"]) {
		audit := cicdObject(value)
		if audit["code"] != nil && !cicdOperationCode(audit["code"].(string)) {
			return false
		}
		id := audit["operation_id"].(string)
		auditsByID[id] = append(auditsByID[id], audit)
	}
	for _, value := range cicdList(root["operations"]) {
		operation := cicdObject(value)
		kind := operation["kind"].(string)
		if _, known := cicdOperationKinds[kind]; !known {
			return false
		}
		if operation["code"] != nil && !cicdOperationCode(operation["code"].(string)) {
			return false
		}
		target, ok := operation["target_id"].(string)
		if !ok {
			return false
		}
		attempt := 0
		if operation["target_attempt"] != nil {
			attempt = int(cicdInt(operation["target_attempt"]))
		}
		if !validCICDMutationReply(kind, cicdOperationMutation{targetID: target, targetAttempt: attempt}) || !validCICDOperationTarget(root, operation) {
			return false
		}
		accepted, finished := 0, 0
		var environment any
		for _, audit := range auditsByID[operation["id"].(string)] {
			if audit["action"] != kind || audit["subject"] != operation["subject"] || audit["project_id"] != operation["project_id"] || audit["target_id"] != target || audit["at"].(string) < operation["created_at"].(string) {
				return false
			}
			switch audit["phase"] {
			case "accepted":
				accepted++
				environment = audit["environment_id"]
				if audit["at"] != operation["created_at"] || audit["code"] != nil {
					return false
				}
			case "finished":
				finished++
				if audit["code"] != operation["code"] {
					return false
				}
			default:
				return false
			}
		}
		if accepted != 1 || cicdTerminal(operation["state"]) && finished != 1 || !cicdTerminal(operation["state"]) && finished != 0 {
			return false
		}
		for _, audit := range auditsByID[operation["id"].(string)] {
			if audit["environment_id"] != environment {
				return false
			}
		}
	}
	return true
}

func validCICDOperationCompletion(root, operation map[string]any) bool {
	if !cicdTerminal(operation["state"]) {
		return true
	}
	kind := operation["kind"]
	if kind == "cacheDelete" && operation["state"] == "succeeded" {
		for _, value := range cicdList(root["caches"]) {
			cache := cicdObject(value)
			if cache["project_id"] == operation["project_id"] && cache["key"] == operation["target_id"] {
				return false
			}
		}
	}
	if kind == "workerResume" && operation["state"] == "succeeded" {
		for _, value := range cicdList(root["workers"]) {
			worker := cicdObject(value)
			if worker["id"] == operation["target_id"] {
				return worker["enabled"] == true && worker["draining"] == false
			}
		}
		return false
	}
	return true
}

func validCICDOperationTarget(root, operation map[string]any) bool {
	kind, target, project := operation["kind"].(string), operation["target_id"], operation["project_id"]
	collection := map[string]string{"passwordChange": "users", "projectPut": "projects", "workerPut": "workers", "workerDrain": "workers", "workerResume": "workers", "targetPut": "targets", "userPut": "users", "servicePut": "services", "definitionPut": "definitions", "templatePut": "templates", "environmentPut": "environments", "secretPut": "secrets", "artifactPin": "artifacts", "cacheDelete": "caches"}[kind]
	if strings.HasPrefix(kind, "run") {
		collection = "runs"
	}
	if strings.HasPrefix(kind, "deployment") {
		collection = "deployments"
	}
	if kind == "grantPut" {
		for _, name := range []string{"users", "services"} {
			for _, value := range cicdList(root[name]) {
				if cicdObject(value)["id"] == target {
					return true
				}
			}
		}
		return false
	}
	for _, value := range cicdList(root[collection]) {
		object := cicdObject(value)
		field := "id"
		if collection == "artifacts" {
			field = "digest"
		}
		if collection == "caches" {
			field = "key"
		}
		if object[field] != target {
			continue
		}
		if collection == "runs" && (object["project_id"] != project || object["attempt"] != operation["target_attempt"]) {
			continue
		}
		if collection == "deployments" && cicdObject(object["plan"])["project_id"] != project {
			continue
		}
		if collection == "definitions" || collection == "templates" || collection == "environments" || collection == "secrets" || collection == "artifacts" || collection == "caches" {
			if object["project_id"] != project {
				continue
			}
		}
		if kind == "passwordChange" && target != operation["subject"] {
			return false
		}
		return true
	}
	// Content/settings can outlive the retained operation as specified; runs cannot.
	return cicdTerminal(operation["state"]) && (kind == "cacheDelete" || kind == "artifactPin" || kind == "secretPut")
}

func validCICDReservedCapacity(root map[string]any, size int) bool {
	reserved := int64(0)
	for _, value := range cicdList(root["operations"]) {
		state := cicdObject(value)["state"]
		if state == "accepted" || state == "recovery_required" {
			reserved++
		}
	}
	return int64(size)+reserved*262144 <= cicdSnapshotLimit && int64(len(cicdList(root["operations"])))+reserved <= 10000 && int64(len(cicdList(root["audits"])))+reserved <= 20000
}

func validateCICDRun(run map[string]any) bool {
	plan := cicdObject(run["plan"])
	data := cicdObject(plan["data"])
	if run["project_id"] != data["project_id"] || cicdValueDigest(data) != plan["digest"] {
		return false
	}
	if len(cicdValueBytes(data)) > 524288 {
		return false
	}
	source := cicdObject(data["source"])
	files := cicdList(source["files"])
	lastPath, total := "", int64(0)
	for _, value := range files {
		file := cicdObject(value)
		path := file["path"].(string)
		if path <= lastPath {
			return false
		}
		lastPath = path
		total += cicdInt(file["size"])
	}
	if total > 104857600 || cicdValueDigest(files) != source["tree_digest"] {
		return false
	}
	jobs, planned := cicdList(run["jobs"]), cicdList(data["jobs"])
	if len(jobs) != len(planned) {
		return false
	}
	seen := make(map[string]map[string]any)
	for i, value := range planned {
		job, result := cicdObject(value), cicdObject(jobs[i])
		id := job["id"].(string)
		if result["id"] != id || seen[id] != nil {
			return false
		}
		needs := make(map[string]bool)
		for _, dependency := range cicdList(job["needs"]) {
			key := dependency.(string)
			if needs[key] || seen[key] == nil {
				return false
			}
			needs[key] = true
		}
		task := cicdObject(job["task"])
		if task["kind"] == "deploy" {
			key := task["artifact_from"].(string)
			if !needs[key] || cicdObject(seen[key]["task"])["kind"] != "builder" {
				return false
			}
		}
		if job["selected"] == false {
			if result["state"] != "skipped" || result["finished_at"] != run["created_at"] || result["duration_ms"] != json.Number("0") {
				return false
			}
			for _, field := range []string{"execution_id", "worker_id", "started_at", "code", "artifact", "reused_from_attempt"} {
				if result[field] != nil {
					return false
				}
			}
			if cicdInt(result["try_number"]) != 0 || cicdInt(result["lease_generation"]) != 0 || result["log_gap"] != false {
				return false
			}
		}
		if result["state"] == "queued" || result["state"] == "skipped" {
			for _, field := range []string{"execution_id", "worker_id", "started_at", "code", "artifact", "reused_from_attempt"} {
				if result[field] != nil {
					return false
				}
			}
			if cicdInt(result["try_number"]) != 0 || cicdInt(result["lease_generation"]) != 0 || result["log_gap"] != false {
				return false
			}
			if result["state"] == "queued" && result["duration_ms"] != nil || result["state"] == "skipped" && result["duration_ms"] != json.Number("0") {
				return false
			}
		}
		if result["started_at"] != nil && result["finished_at"] != nil && result["finished_at"].(string) < result["started_at"].(string) {
			return false
		}
		if cicdTerminal(result["state"]) != (result["finished_at"] != nil) {
			return false
		}
		if result["started_at"] != nil && result["started_at"].(string) < run["created_at"].(string) || result["finished_at"] != nil && result["finished_at"].(string) < run["created_at"].(string) {
			return false
		}
		seen[id] = job
	}
	if cicdTerminal(run["state"]) != (run["finished_at"] != nil) {
		return false
	}
	if run["started_at"] != nil && run["started_at"].(string) < run["created_at"].(string) || run["finished_at"] != nil && run["finished_at"].(string) < run["created_at"].(string) || run["started_at"] != nil && run["finished_at"] != nil && run["finished_at"].(string) < run["started_at"].(string) {
		return false
	}
	return true
}

func validateCICDJSONBytes(data []byte) bool {
	if !utf8.Valid(data) || !json.Valid(data) {
		return false
	}
	// encoding/json replaces unpaired surrogates; identities must not be repaired.
	for i := 0; i < len(data); {
		if data[i] != '"' {
			i++
			continue
		}
		i++
		for i < len(data) && data[i] != '"' {
			if data[i] != '\\' {
				i++
				continue
			}
			if data[i+1] != 'u' {
				i += 2
				continue
			}
			value, err := strconv.ParseUint(string(data[i+2:i+6]), 16, 16)
			if err != nil {
				return false
			}
			if value >= 0xd800 && value <= 0xdbff {
				if i+12 > len(data) || string(data[i+6:i+8]) != `\u` {
					return false
				}
				low, err := strconv.ParseUint(string(data[i+8:i+12]), 16, 16)
				if err != nil || low < 0xdc00 || low > 0xdfff {
					return false
				}
				i += 12
			} else if value >= 0xdc00 && value <= 0xdfff {
				return false
			} else {
				i += 6
			}
		}
		i++
	}
	return true
}

func validateCICDInteger(number json.Number) bool {
	if number.String() == "-0" {
		return false
	}
	value, err := strconv.ParseInt(number.String(), 10, 64)
	return err == nil && value >= -9007199254740991 && value <= 9007199254740991
}

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == "statefile" && validateExactOwnerFiles(contract.Files, []string{"statefile.go", "model.go", "validate.go", "execute.go", "statefile_test.go"})
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

func validateTargetPath(path string) error {
	if path == "" {
		return errors.New("statefile path must not be empty")
	}
	clean := filepath.Clean(path)
	if clean != path {
		return errors.New("statefile path must be clean")
	}
	parent := filepath.Dir(path)
	if parent == "" || parent == "." {
		return errors.New("statefile parent directory must be explicit")
	}
	if info, err := os.Lstat(parent); err != nil {
		return err
	} else if info.Mode()&os.ModeSymlink != 0 {
		return errors.New("statefile parent directory must not be symlink")
	} else if !info.IsDir() {
		return errors.New("statefile parent path is not directory")
	}
	if info, err := os.Lstat(path); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return errors.New("statefile target must not be symlink")
		}
		if !info.Mode().IsRegular() {
			return errors.New("statefile target must be regular file")
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
