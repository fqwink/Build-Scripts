package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"

	"github.com/fqwink/build-scripts/components/security"
)

func hasExactArg(args []string, target string) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}
	return false
}

func safeArgvTokens(args []string) bool {
	for _, arg := range args {
		if !utf8.ValidString(arg) {
			return false
		}
		for _, r := range arg {
			if r < 0x20 || r == 0x7f {
				return false
			}
		}
	}
	return true
}

func validateCredentials(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cred struct {
		PasswordHash string  `json:"password_hash"`
		Salt         string  `json:"salt"`
		Algorithm    string  `json:"algorithm"`
		Iterations   int     `json:"iterations"`
		UpdatedAt    string  `json:"updated_at"`
		LastLoginAt  *string `json:"last_login_at"`
		LoginCount   int64   `json:"login_count"`
	}
	if err := json.Unmarshal(data, &cred); err != nil {
		return err
	}
	if err := security.ValidatePasswordRecord(security.PasswordRecord{
		PasswordHash: cred.PasswordHash,
		Salt:         cred.Salt,
		Algorithm:    cred.Algorithm,
		Iterations:   cred.Iterations,
	}); err != nil {
		return fmt.Errorf("invalid credentials")
	}
	if cred.LoginCount < 0 {
		return fmt.Errorf("invalid credentials")
	}
	if _, err := time.Parse("2006-01-02T15:04:05Z", cred.UpdatedAt); err != nil {
		return fmt.Errorf("invalid credentials")
	}
	if cred.LastLoginAt != nil {
		if _, err := time.Parse("2006-01-02T15:04:05Z", *cred.LastLoginAt); err != nil {
			return fmt.Errorf("invalid credentials")
		}
	}
	return nil
}

func atomicWriteText(path, value string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(value), mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func setupIsLowerHex(value string, length int) bool {
	if len(value) != length {
		return false
	}
	for _, r := range value {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}
