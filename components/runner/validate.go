package runner

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/fqwink/build-scripts/components/statefile"
)

func cicdPlanObject(kind string, data []byte) (map[string]any, error) {
	if statefile.ValidateCICDValue(kind, data) != nil {
		return nil, ErrCICDSchema
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var object map[string]any
	if decoder.Decode(&object) != nil || object == nil {
		return nil, ErrCICDSchema
	}
	return object, nil
}

func cicdPlanBytes(value any) []byte {
	data, _ := json.Marshal(value)
	return data
}

func cicdPlanClone(object map[string]any) map[string]any {
	decoder := json.NewDecoder(bytes.NewReader(cicdPlanBytes(object)))
	decoder.UseNumber()
	var copy map[string]any
	_ = decoder.Decode(&copy)
	return copy
}

func cicdPlanInt(value any) int64 {
	n, _ := value.(json.Number)
	v, _ := n.Int64()
	return v
}

func cicdPlanKeys(object map[string]any) []string {
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func cicdInputValue(declaration map[string]any, value any) bool {
	if statefile.ValidateCICDValue("Scalar", cicdPlanBytes(value)) != nil {
		return false
	}
	min, max := cicdPlanInt(declaration["min"]), cicdPlanInt(declaration["max"])
	switch declaration["kind"] {
	case "string":
		v, ok := value.(string)
		return ok && int64(len(v)) >= min && int64(len(v)) <= max
	case "integer":
		v, ok := value.(json.Number)
		return ok && cicdPlanInt(v) >= min && cicdPlanInt(v) <= max
	case "boolean":
		_, ok := value.(bool)
		return ok
	}
	return false
}

func cicdNormalizeInputs(declarations, supplied map[string]any) (map[string]any, error) {
	result := make(map[string]any)
	for key := range supplied {
		if declarations[key] == nil {
			return nil, ErrCICDSchema
		}
	}
	for _, key := range cicdPlanKeys(declarations) {
		d := declarations[key].(map[string]any)
		min, max := cicdPlanInt(d["min"]), cicdPlanInt(d["max"])
		if min > max || d["required"] == true && d["default"] != nil ||
			d["kind"] == "string" && (min < 0 || max > 4096) ||
			d["kind"] == "boolean" && (min != 0 || max != 1) {
			return nil, ErrCICDSchema
		}
		choices := d["choices"].([]any)
		for i, choice := range choices {
			if !cicdInputValue(d, choice) {
				return nil, ErrCICDSchema
			}
			for _, previous := range choices[:i] {
				if reflect.DeepEqual(previous, choice) {
					return nil, ErrCICDSchema
				}
			}
		}
		allowed := func(value any) bool {
			if !cicdInputValue(d, value) {
				return false
			}
			if len(choices) == 0 {
				return true
			}
			for _, choice := range choices {
				if reflect.DeepEqual(choice, value) {
					return true
				}
			}
			return false
		}
		if d["default"] != nil && !allowed(d["default"]) {
			return nil, ErrCICDSchema
		}
		value, exists := supplied[key]
		if !exists {
			if d["required"] == true {
				return nil, ErrCICDSchema
			}
			value = d["default"]
			if value == nil {
				continue
			}
		}
		if !allowed(value) {
			return nil, ErrCICDSchema
		}
		result[key] = value
	}
	return result, nil
}

// Substitution consumes the original field once and never scans the replacement.
func cicdSubstitute(value any, inputs, matrix map[string]any) (any, error) {
	switch v := value.(type) {
	case string:
		if !strings.Contains(v, "${") {
			return v, nil
		}
		for _, scope := range []struct {
			prefix string
			values map[string]any
		}{{"${input.", inputs}, {"${matrix.", matrix}} {
			if strings.HasPrefix(v, scope.prefix) && strings.HasSuffix(v, "}") {
				key := v[len(scope.prefix) : len(v)-1]
				if statefile.ValidateCICDValue("Id", cicdPlanBytes(key)) != nil {
					return nil, ErrCICDSchema
				}
				if replacement, exists := scope.values[key]; exists {
					return replacement, nil
				}
			}
		}
		return nil, ErrCICDSchema
	case []any:
		copy := make([]any, len(v))
		for i, child := range v {
			replacement, err := cicdSubstitute(child, inputs, matrix)
			if err != nil {
				return nil, err
			}
			copy[i] = replacement
		}
		return copy, nil
	case map[string]any:
		copy := make(map[string]any, len(v))
		for _, key := range cicdPlanKeys(v) {
			replacement, err := cicdSubstitute(v[key], inputs, matrix)
			if err != nil {
				return nil, err
			}
			copy[key] = replacement
		}
		return copy, nil
	default:
		return value, nil
	}
}

func cicdScalarKind(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case bool:
		return "boolean"
	case json.Number:
		return "integer"
	}
	return ""
}

func cicdMatrix(matrix map[string]any) ([]map[string]any, error) {
	axes := matrix["axes"].(map[string]any)
	keys := cicdPlanKeys(axes)
	excluded, included := matrix["exclude"].([]any), matrix["include"].([]any)
	if len(keys) == 0 {
		if len(excluded)+len(included) != 0 {
			return nil, ErrCICDMatrix
		}
		return []map[string]any{{}}, nil
	}
	for _, key := range keys {
		values := axes[key].([]any)
		for i, value := range values {
			if cicdScalarKind(value) != cicdScalarKind(values[0]) {
				return nil, ErrCICDMatrix
			}
			for _, previous := range values[:i] {
				if reflect.DeepEqual(value, previous) {
					return nil, ErrCICDMatrix
				}
			}
		}
	}
	for _, group := range []struct {
		list []any
		full bool
	}{{excluded, false}, {included, true}} {
		for _, value := range group.list {
			tuple := value.(map[string]any)
			if len(tuple) == 0 || group.full && len(tuple) != len(keys) {
				return nil, ErrCICDMatrix
			}
			for key, scalar := range tuple {
				axis, exists := axes[key]
				if !exists || cicdScalarKind(scalar) != cicdScalarKind(axis.([]any)[0]) {
					return nil, ErrCICDMatrix
				}
			}
		}
	}
	result := make([]map[string]any, 0)
	tuple := make(map[string]any)
	var enumerate func(int) error
	enumerate = func(index int) error {
		// A complete partial exclusion prunes its entire remaining subtree.
		for _, value := range excluded {
			match := true
			for key, scalar := range value.(map[string]any) {
				if !reflect.DeepEqual(tuple[key], scalar) {
					match = false
					break
				}
			}
			if match {
				return nil
			}
		}
		if index == len(keys) {
			if len(result) == 256 {
				return ErrCICDMatrix
			}
			result = append(result, cicdPlanClone(tuple))
			return nil
		}
		key := keys[index]
		for _, scalar := range axes[key].([]any) {
			tuple[key] = scalar
			if err := enumerate(index + 1); err != nil {
				return err
			}
		}
		delete(tuple, key)
		return nil
	}
	if err := enumerate(0); err != nil {
		return nil, err
	}
	for _, value := range included {
		add := value.(map[string]any)
		for _, previous := range result {
			if reflect.DeepEqual(previous, add) {
				return nil, ErrCICDMatrix
			}
		}
		result = append(result, cicdPlanClone(add))
		if len(result) > 256 {
			return nil, ErrCICDMatrix
		}
	}
	if len(result) == 0 {
		return nil, ErrCICDMatrix
	}
	return result, nil
}

func cicdMatchGlob(pattern, target string) (bool, error) {
	if statefile.ValidateCICDValue("Glob", cicdPlanBytes(pattern)) != nil || !utf8.ValidString(target) {
		return false, ErrCICDSchema
	}
	p, v := strings.Split(pattern, "/"), strings.Split(target, "/")
	cells := (len(p) + 1) * (len(v) + 1)
	for _, segment := range p {
		if segment == "**" {
			continue
		}
		for _, part := range v {
			cells += (utf8.RuneCountInString(segment) + 1) * (utf8.RuneCountInString(part) + 1)
			if cells > 100000 {
				return false, ErrCICDSchema
			}
		}
	}
	if cells > 100000 {
		return false, ErrCICDSchema
	}
	dp := make([][]bool, len(p)+1)
	for i := range dp {
		dp[i] = make([]bool, len(v)+1)
	}
	dp[0][0] = true
	for i, segment := range p {
		for j := 0; j <= len(v); j++ {
			if segment == "**" {
				dp[i+1][j] = dp[i][j] || j > 0 && dp[i+1][j-1]
			} else if j > 0 && dp[i][j-1] {
				dp[i+1][j] = cicdMatchSegment([]rune(segment), []rune(v[j-1]))
			}
		}
	}
	return dp[len(p)][len(v)], nil
}

func cicdMatchSegment(pattern, target []rune) bool {
	previous := make([]bool, len(target)+1)
	previous[0] = true
	for _, char := range pattern {
		next := make([]bool, len(target)+1)
		if char == '*' {
			next[0] = previous[0]
		}
		for j, value := range target {
			if char == '*' {
				next[j+1] = previous[j+1] || next[j]
			} else {
				next[j+1] = previous[j] && (char == '?' || char == value)
			}
		}
		previous = next
	}
	return previous[len(target)]
}

func ValidateCICDResources(r CICDResources) error {
	if r.CPUPercent < 10 || r.CPUPercent > 400 ||
		r.MemoryMB < cicdMinimumJobMB || r.MemoryMB > cicdMaximumJobMB ||
		r.DiskMB < 128 || r.DiskMB > 2048 ||
		r.Tasks < 8 || r.Tasks > 256 ||
		r.TimeoutSeconds < 1 || r.TimeoutSeconds > 3600 ||
		r.OutputMB < 1 || r.OutputMB > 16 {
		return ErrCICDSchema
	}
	if r.MemoryMB < r.DiskMB+128 {
		return ErrCICDSchema
	}
	return nil
}

// CheckCICDMemory requires a fresh observation and the caller's claim lock.
// Reservations include unstarted claims and claims awaiting confirmed cleanup.
func CheckCICDMemory(memory CICDHostMemory, reservations []int64, configuredSlots, requestedMB int64) error {
	if configuredSlots < 1 || configuredSlots > 4 || requestedMB < cicdMinimumJobMB || requestedMB > cicdMaximumJobMB {
		return ErrCICDSchema
	}
	if memory.TotalMB < cicdMinimumHostMB || memory.TotalMB > cicdMaximumHostMB ||
		memory.AvailableMB < 0 || memory.AvailableMB > memory.TotalMB {
		return ErrCICDCapability
	}
	if int64(len(reservations)) > configuredSlots {
		return ErrCICDStateCorrupt
	}
	var reservedMB int64
	for _, value := range reservations {
		if value < cicdMinimumJobMB || value > cicdMaximumJobMB {
			return ErrCICDStateCorrupt
		}
		reservedMB += value
	}
	budget := memory.TotalMB - cicdHostReserveMB
	if requestedMB > budget {
		return ErrCICDCapability
	}
	slots := configuredSlots
	if memory.TotalMB < cicdSingleSlotMB {
		slots = 1
	}
	if int64(len(reservations)) >= slots || reservedMB+requestedMB > budget ||
		memory.AvailableMB < reservedMB+requestedMB+cicdMemoryGuardMB {
		return ErrCICDBusy
	}
	return nil
}

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == "runner" && validateExactOwnerFiles(contract.Files, []string{"runner.go", "model.go", "validate.go", "execute.go", "runner_test.go"})
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
