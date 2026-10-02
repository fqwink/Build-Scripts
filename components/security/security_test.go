package security

import "testing"

func TestOwnerFileContract(t *testing.T) {
	if Owner() != "security" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executeOwnerFileContract(newOwnerFileContract()) {
		t.Fatalf("security owner file contract failed")
	}
}

func TestPBKDF2HMACSHA256Vectors(t *testing.T) {
	cases := []struct {
		name       string
		password   string
		wantHash   string
		iterations int
	}{
		{
			name:       "short-1",
			password:   "password",
			wantHash:   "120fb6cffcf8b32c43e7225256c4f837a86548c92ccc35480805987cb70be17b",
			iterations: 1,
		},
		{
			name:       "short-2",
			password:   "password",
			wantHash:   "ae4d0c95af6b46d32d0adff928f06dd02a303f8ef3c251dfd6e2d85a95474c43",
			iterations: 2,
		},
		{
			name:       "production",
			password:   "password",
			wantHash:   "669cfe52482116fda1aa2cbe409b2f56c8e4563752b7a28f6eaab614ee005178",
			iterations: CredentialIterations,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hash := pbkdf2HMACSHA256([]byte(tc.password), []byte("salt"), tc.iterations, CredentialTagLength)
			if got := toLowerHex(hash); got != tc.wantHash {
				t.Fatalf("hash mismatch: got %s want %s", got, tc.wantHash)
			}
		})
	}
}

func TestPasswordRecordContract(t *testing.T) {
	record, err := NewPasswordRecord("password")
	if err != nil {
		t.Fatalf("NewPasswordRecord failed: %v", err)
	}
	if err := ValidatePasswordRecord(record); err != nil {
		t.Fatalf("record validation failed: %v", err)
	}
	ok, err := VerifyPassword(record, "password")
	if err != nil || !ok {
		t.Fatalf("password verification failed: ok=%v err=%v", ok, err)
	}
	ok, err = VerifyPassword(record, "different-password")
	if err != nil {
		t.Fatalf("unexpected verification error: %v", err)
	}
	if ok {
		t.Fatalf("different password verified")
	}
}

func toLowerHex(data []byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, len(data)*2)
	for i, b := range data {
		out[i*2] = digits[b>>4]
		out[i*2+1] = digits[b&0x0f]
	}
	return string(out)
}
