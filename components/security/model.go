package security

import "errors"

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
