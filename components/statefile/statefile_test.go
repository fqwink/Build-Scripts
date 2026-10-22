package statefile

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

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
