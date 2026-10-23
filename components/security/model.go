package security

import (
	"errors"
	"time"
)

var ErrCICDAuth = errors.New("P19_AUTH")
var ErrCICDForbidden = errors.New("P19_FORBIDDEN")
var ErrCICDQuota = errors.New("P19_QUOTA")
var ErrCICDSecret = errors.New("P19_SECRET")

// CICDCredential is an internal identity-store value, never an HTTP projection.
type CICDCredential struct {
	PasswordHash string  `json:"password_hash"`
	Salt         string  `json:"salt"`
	Algorithm    string  `json:"algorithm"`
	Iterations   int     `json:"iterations"`
	MustChange   bool    `json:"must_change"`
	LoginCount   int64   `json:"login_count"`
	LastLoginAt  *string `json:"last_login_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type CICDPrincipal struct {
	Kind, ID   string
	ExpiresAt  time.Time
	MustChange string
}

type CICDSubject struct {
	Kind, ID  string
	Enabled   bool
	ExpiresAt *time.Time
}

type CICDGrant struct {
	SubjectKind, SubjectID, ProjectID, EnvironmentID, Role string
	Enabled                                                bool
	ExpiresAt                                              *time.Time
}

type CICDAuthorizationSnapshot struct {
	Generation                                       int64
	Subject                                          CICDSubject
	Grants                                           []CICDGrant
	ProjectExists, ProjectEnabled, EnvironmentExists bool
}

type CICDAccess struct {
	Operation, ProjectID, EnvironmentID, ResourceOwner string
	ProposedGrant                                      *CICDGrant
}

type CICDAuthorization struct {
	Generation             int64
	SubjectKind, SubjectID string
	ExpiresAt              time.Time
	Grants                 []CICDGrant
}

const (
	CredentialAlgorithmPBKDF2HMACSHA256V1 = "pbkdf2_hmac_sha256_v1"
	CredentialIterations                  = 600000
	CredentialSaltLength                  = 32
	CredentialTagLength                   = 32
)

var (
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidPasswordRecord = errors.New("invalid password record")
)

type PasswordRecord struct {
	PasswordHash string
	Salt         string
	Algorithm    string
	Iterations   int
}

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: Owner(),
		Files: []string{"security.go", "model.go", "validate.go", "execute.go", "security_test.go"},
	}
}
