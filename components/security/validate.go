package security

import (
	"strings"
	"unicode/utf8"
)

func ValidatePasswordInput(password string) error {
	if password == "" || !utf8.ValidString(password) {
		return ErrInvalidPassword
	}
	runeCount := utf8.RuneCountInString(password)
	if runeCount < 8 || runeCount > 128 {
		return ErrInvalidPassword
	}
	for _, r := range password {
		if r == 0 || r == '\r' || r == '\n' {
			return ErrInvalidPassword
		}
	}
	return nil
}

func ValidatePasswordRecord(record PasswordRecord) error {
	if record.Algorithm != CredentialAlgorithmPBKDF2HMACSHA256V1 || record.Iterations != CredentialIterations {
		return ErrInvalidPasswordRecord
	}
	if record.Salt != strings.ToLower(record.Salt) || !isLowerHex(record.Salt, CredentialSaltLength*2) {
		return ErrInvalidPasswordRecord
	}
	if record.PasswordHash != strings.ToLower(record.PasswordHash) || !isLowerHex(record.PasswordHash, CredentialTagLength*2) {
		return ErrInvalidPasswordRecord
	}
	return nil
}

func isLowerHex(value string, length int) bool {
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

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == "security" && validateExactOwnerFiles(contract.Files, []string{"security.go", "model.go", "validate.go", "execute.go", "security_test.go"})
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
