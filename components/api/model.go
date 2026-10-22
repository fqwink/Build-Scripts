package api

import (
	"io"
	"sync"
	"time"

	"github.com/fqwink/build-scripts/components/security"
	"github.com/fqwink/build-scripts/components/statefile"
)

type cicdIdentityBackend interface {
	Login(string, string, string, time.Time, func() (string, bool, error)) (statefile.CICDLoginResult, error)
	ReadService(string) (statefile.CICDIdentity, error)
	Authorize(security.CICDPrincipal, security.CICDAccess, time.Time, func(statefile.CICDIdentity) error) (security.CICDAuthorization, error)
}

type cicdAuthSession struct {
	principal    security.CICDPrincipal
	binding      string
	adminBinding string
}

type cicdAuthWindow struct {
	start time.Time
	count int
}
type cicdAuthFailure struct {
	at, until time.Time
	count     int
}

// This handler is deliberately not registered by the legacy listener.
type cicdAuthHTTP struct {
	legacy    *APIServer
	store     cicdIdentityBackend
	gate      chan struct{}
	entropyMu sync.Mutex
	entropy   io.Reader
	sessions  map[string]cicdAuthSession
	revoked   map[string]time.Time
	windows   map[string]cicdAuthWindow
	failures  map[string]cicdAuthFailure
}

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "api",
		Files: []string{"api.go", "model.go", "validate.go", "execute.go", "api_test.go"},
	}
}
