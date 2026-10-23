package security

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/url"
	"sort"
	"time"
	"unicode/utf8"
)

func CICDMaskPatterns(secrets [][]byte) ([][]byte, error) {
	if len(secrets) > 16 {
		return nil, ErrCICDSecret
	}
	patterns := make(map[string]bool)
	total := 0
	for _, secret := range secrets {
		if len(secret) < 1 || len(secret) > 4096 || !utf8.Valid(secret) || bytes.ContainsAny(secret, "\x00\r\n") {
			return nil, ErrCICDSecret
		}
		total += len(secret)
		if total > 16384 {
			return nil, ErrCICDSecret
		}
		value := string(secret)
		for _, pattern := range []string{value, base64.StdEncoding.EncodeToString(secret), base64.URLEncoding.EncodeToString(secret), base64.RawURLEncoding.EncodeToString(secret), url.QueryEscape(value), url.PathEscape(value)} {
			patterns[pattern] = true
		}
	}
	result := make([][]byte, 0, len(patterns))
	for pattern := range patterns {
		result = append(result, []byte(pattern))
	}
	sort.Slice(result, func(i, j int) bool {
		if len(result[i]) != len(result[j]) {
			return len(result[i]) > len(result[j])
		}
		return bytes.Compare(result[i], result[j]) < 0
	})
	return result, nil
}

func VerifyCICDPassword(credential CICDCredential, password string) bool {
	ok, err := VerifyPassword(PasswordRecord{PasswordHash: credential.PasswordHash, Salt: credential.Salt, Algorithm: credential.Algorithm, Iterations: credential.Iterations}, password)
	return err == nil && ok
}

func CICDBearerKind(token string) string {
	if isLowerHex(token, 64) {
		return "user"
	}
	if len(token) == 43 {
		decoded, err := base64.RawURLEncoding.Strict().DecodeString(token)
		if err == nil && len(decoded) == 32 && base64.RawURLEncoding.EncodeToString(decoded) == token {
			return "service"
		}
	}
	return ""
}

func CICDCredentialFingerprint(credential CICDCredential) string {
	credential.LoginCount, credential.LastLoginAt = 0, nil
	data, _ := json.Marshal(credential)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// AuthorizeCICD consumes a current snapshot; it does not authenticate credentials.
func AuthorizeCICD(principal CICDPrincipal, snapshot CICDAuthorizationSnapshot, access CICDAccess, now time.Time) (CICDAuthorization, error) {
	deny := CICDAuthorization{}
	subject := snapshot.Subject
	if !validCICDTime(now) || !validCICDTime(principal.ExpiresAt) || !now.Before(principal.ExpiresAt) || !validCICDIdentifier(principal.ID) || principal.ID == "system" || principal.ID == "worker" || principal.ID == "site" || (principal.Kind != "user" && principal.Kind != "service") || subject.Kind != principal.Kind || subject.ID != principal.ID || !subject.Enabled {
		return deny, ErrCICDAuth
	}
	if principal.MustChange != "none" && principal.MustChange != "prompt" && principal.MustChange != "forced" || principal.Kind == "service" && principal.MustChange != "none" {
		return deny, ErrCICDAuth
	}
	expires := principal.ExpiresAt
	if principal.Kind == "service" {
		if subject.ExpiresAt == nil || !validCICDTime(*subject.ExpiresAt) || !now.Before(*subject.ExpiresAt) {
			return deny, ErrCICDAuth
		}
		if subject.ExpiresAt.Before(expires) {
			expires = *subject.ExpiresAt
		}
	}
	if snapshot.Generation < 1 || snapshot.Generation > 9007199254740991 || access.ProjectID != "" && !validCICDIdentifier(access.ProjectID) || access.EnvironmentID != "" && (!validCICDIdentifier(access.EnvironmentID) || access.ProjectID == "") {
		return deny, ErrCICDForbidden
	}
	if principal.MustChange == "forced" && access.Operation != "passwordChange" && access.Operation != "logout" {
		return deny, ErrCICDForbidden
	}
	grants := make([]CICDGrant, 0)
	for _, grant := range snapshot.Grants {
		if !validCICDGrant(grant) {
			return deny, ErrCICDForbidden
		}
		if grant.SubjectKind != principal.Kind || grant.SubjectID != principal.ID || !grant.Enabled || grant.ExpiresAt != nil && !now.Before(*grant.ExpiresAt) {
			continue
		}
		if grant.ExpiresAt != nil {
			copyTime := *grant.ExpiresAt
			grant.ExpiresAt = &copyTime
		}
		grants = append(grants, grant)
	}
	has := func(role string) bool {
		for _, grant := range grants {
			if role == "platform-admin" {
				if grant.Role == role {
					return true
				}
				continue
			}
			if access.ProjectID == "" || grant.ProjectID != access.ProjectID {
				continue
			}
			if role == "project-reader" {
				if grant.Role != "approver" || access.EnvironmentID == "" || grant.EnvironmentID == access.EnvironmentID {
					return true
				}
				continue
			}
			if grant.Role == role && (role != "approver" || access.EnvironmentID != "" && grant.EnvironmentID == access.EnvironmentID) {
				return true
			}
			if role == "viewer" && (grant.Role == "operator" || grant.Role == "project-admin") || role == "operator" && grant.Role == "project-admin" {
				return true
			}
		}
		return false
	}
	allowed := false
	scope := "project"
	switch access.Operation {
	case "logout":
		scope, allowed = "global", true
	case "passwordChange":
		scope, allowed = "global", principal.Kind == "user" && (access.ResourceOwner == "" || access.ResourceOwner == principal.ID)
	case "projectPut", "workerPut", "targetPut", "userPut", "servicePut", "workerDrain", "workerResume", "workersList", "targetsList", "usersList", "servicesList":
		scope, allowed = "global", has("platform-admin")
	case "projectsList":
		scope = "optional"
		allowed = has("platform-admin") || has("project-reader")
		if access.ProjectID == "" {
			allowed = has("platform-admin") || len(grants) > 0
		}
	case "grantPut", "grantsList":
		scope = "optional"
		allowed = has("platform-admin") || has("project-admin")
		if access.Operation == "grantPut" {
			grant := access.ProposedGrant
			allowed = allowed && principal.Kind == "user" && grant != nil && validCICDGrant(*grant)
			if allowed {
				allowed = grant.ProjectID == access.ProjectID && grant.EnvironmentID == access.EnvironmentID && (has("platform-admin") || grant.Role == "viewer" || grant.Role == "operator" || grant.Role == "approver")
			}
		}
	case "operationGet":
		scope = "optional"
		if access.ProjectID == "" {
			allowed = has("platform-admin")
		} else {
			allowed = has("project-admin") || access.ResourceOwner == principal.ID && has("project-reader")
		}
	case "deploymentGet":
		allowed = has("viewer") || has("approver")
	case "deploymentApprove", "deploymentReject":
		allowed = principal.Kind == "user" && has("approver")
	case "plan", "definitionsList", "templatesList", "environmentsList", "runsList", "runGet", "logsGet", "artifactsList", "artifactGet", "artifactDownload", "artifactCompare", "cachesList", "environmentStatus", "deploymentsList", "metricsGet":
		allowed = has("viewer")
	case "runStart", "runRetry", "runResume", "runCancel", "artifactPin", "deploymentPlan", "deploymentStart", "deploymentCancel", "deploymentRollback":
		allowed = has("operator")
	case "definitionPut", "templatePut", "environmentPut", "secretPut", "secretsList", "cacheDelete", "auditsList":
		allowed = has("project-admin")
	case "deploymentReconcile":
		allowed = has("platform-admin")
	default:
		return deny, ErrCICDForbidden
	}
	if principal.Kind == "service" && (access.Operation == "userPut" || access.Operation == "servicePut" || access.Operation == "grantPut") {
		allowed = false
	}
	if scope == "global" && (access.ProjectID != "" || access.EnvironmentID != "") || scope == "project" && (access.ProjectID == "" || !snapshot.ProjectExists) {
		allowed = false
	}
	if access.EnvironmentID != "" && !snapshot.EnvironmentExists {
		allowed = false
	}
	if (access.Operation == "runStart" || access.Operation == "runRetry" || access.Operation == "runResume" || access.Operation == "deploymentStart" || access.Operation == "deploymentRollback") && !snapshot.ProjectEnabled {
		allowed = false
	}
	if !allowed {
		return deny, ErrCICDForbidden
	}
	wireGrants := make([]any, 0, len(grants))
	for _, grant := range grants {
		var project, environment any
		if grant.ProjectID != "" {
			project = grant.ProjectID
		}
		if grant.EnvironmentID != "" {
			environment = grant.EnvironmentID
		}
		wireGrants = append(wireGrants, map[string]any{"subject_kind": grant.SubjectKind, "subject_id": grant.SubjectID, "project_id": project, "environment_id": environment, "role": grant.Role, "enabled": grant.Enabled, "expires_at": grant.ExpiresAt})
	}
	wire, err := json.Marshal(map[string]any{"subject_kind": principal.Kind, "subject_id": principal.ID, "generation": snapshot.Generation, "grants": wireGrants, "expires_at": expires})
	if err != nil || len(wire) > 65536 {
		return deny, ErrCICDQuota
	}
	return CICDAuthorization{Generation: snapshot.Generation, SubjectKind: principal.Kind, SubjectID: principal.ID, ExpiresAt: expires, Grants: grants}, nil
}

func GenerateSalt() (string, error) {
	salt := make([]byte, CredentialSaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func HashPassword(password, saltHex string) (string, error) {
	if err := ValidatePasswordInput(password); err != nil {
		return "", err
	}
	salt, err := hex.DecodeString(saltHex)
	if err != nil || len(salt) != CredentialSaltLength {
		return "", ErrInvalidPasswordRecord
	}
	tag := pbkdf2HMACSHA256([]byte(password), salt, CredentialIterations, CredentialTagLength)
	return hex.EncodeToString(tag), nil
}

func VerifyPassword(record PasswordRecord, password string) (bool, error) {
	if err := ValidatePasswordRecord(record); err != nil {
		return false, err
	}
	hash, err := HashPassword(password, record.Salt)
	if err != nil {
		return false, err
	}
	expected, err := hex.DecodeString(record.PasswordHash)
	if err != nil {
		return false, ErrInvalidPasswordRecord
	}
	actual, err := hex.DecodeString(hash)
	if err != nil {
		return false, ErrInvalidPasswordRecord
	}
	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}

func NewPasswordRecord(password string) (PasswordRecord, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return PasswordRecord{}, err
	}
	hash, err := HashPassword(password, salt)
	if err != nil {
		return PasswordRecord{}, err
	}
	return PasswordRecord{
		PasswordHash: hash,
		Salt:         salt,
		Algorithm:    CredentialAlgorithmPBKDF2HMACSHA256V1,
		Iterations:   CredentialIterations,
	}, nil
}

func pbkdf2HMACSHA256(password, salt []byte, iterations, keyLength int) []byte {
	output := make([]byte, 0, keyLength)
	prf := newHMACSHA256PRF(password)
	var blockIndex uint32 = 1
	for len(output) < keyLength {
		block := pbkdf2Block(prf, salt, iterations, blockIndex)
		output = append(output, block...)
		blockIndex++
	}
	return output[:keyLength]
}

func pbkdf2Block(prf hmacSHA256PRF, salt []byte, iterations int, blockIndex uint32) []byte {
	input := make([]byte, len(salt)+4)
	copy(input, salt)
	binary.BigEndian.PutUint32(input[len(salt):], blockIndex)
	u := prf.sum(input)
	out := append([]byte{}, u[:]...)
	for i := 1; i < iterations; i++ {
		u = prf.sum(u[:])
		for j := range out {
			out[j] ^= u[j]
		}
	}
	return out
}

type hmacSHA256PRF struct {
	ipad [64]byte
	opad [64]byte
}

func newHMACSHA256PRF(key []byte) hmacSHA256PRF {
	if len(key) > 64 {
		sum := sha256.Sum256(key)
		key = sum[:]
	}
	var prf hmacSHA256PRF
	for i := range prf.ipad {
		prf.ipad[i] = 0x36
		prf.opad[i] = 0x5c
	}
	for i, b := range key {
		prf.ipad[i] ^= b
		prf.opad[i] ^= b
	}
	return prf
}

func (p hmacSHA256PRF) sum(data []byte) [sha256.Size]byte {
	inner := sha256.New()
	_, _ = inner.Write(p.ipad[:])
	_, _ = inner.Write(data)
	innerSum := inner.Sum(nil)

	outer := sha256.New()
	_, _ = outer.Write(p.opad[:])
	_, _ = outer.Write(innerSum)
	var out [sha256.Size]byte
	copy(out[:], outer.Sum(nil))
	return out
}

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}
