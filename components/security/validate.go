package security

import (
	"strings"
	"time"
	"unicode/utf8"
)

func validCICDIdentifier(id string) bool {
	if len(id) == 0 || len(id) > 64 || id[0] < 'a' || id[0] > 'z' {
		return false
	}
	for _, c := range id {
		if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

func validCICDTime(t time.Time) bool {
	_, offset := t.Zone()
	return t.Year() >= 1970 && t.Year() <= 9999 && t.Nanosecond() == 0 && offset == 0
}

func validCICDGrant(grant CICDGrant) bool {
	if (grant.SubjectKind != "user" && grant.SubjectKind != "service") || !validCICDIdentifier(grant.SubjectID) || grant.ExpiresAt != nil && !validCICDTime(*grant.ExpiresAt) {
		return false
	}
	if grant.Role == "platform-admin" {
		return grant.SubjectKind == "user" && grant.ProjectID == "" && grant.EnvironmentID == ""
	}
	if !validCICDIdentifier(grant.ProjectID) {
		return false
	}
	if grant.Role == "approver" {
		return validCICDIdentifier(grant.EnvironmentID)
	}
	return (grant.Role == "viewer" || grant.Role == "operator" || grant.Role == "project-admin") && grant.EnvironmentID == ""
}

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
