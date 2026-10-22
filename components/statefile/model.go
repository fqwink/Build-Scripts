package statefile

import (
	"errors"
	"time"
)

var ErrCICDSchema = errors.New("P19_SCHEMA")

type ownerFileContract struct {
	Owner string
	Files []string
}

type LockOptions struct {
	Owner          string
	StaleAfter     time.Duration
	RetryInterval  time.Duration
	AcquireTimeout time.Duration
	Clock          func() time.Time
	Sleep          func(time.Duration)
}

type LockHandle struct {
	path string
	file string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: Owner(),
		Files: []string{"statefile.go", "model.go", "validate.go", "execute.go", "statefile_test.go"},
	}
}
