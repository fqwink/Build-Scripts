package api

import (
	"time"
	"unicode/utf8"

	"github.com/fqwink/build-scripts/components/statefile"
)

var runnerSleep = time.Sleep

type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string { return e.Msg }

func validCLIArgToken(arg string) bool {
	if !utf8.ValidString(arg) {
		return false
	}
	for _, r := range arg {
		if r == 0 || r == '\n' || r == '\r' || r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

func acquireStateFileLock(path string) (func(), error) {
	now := time.Now()
	sleep := func(duration time.Duration) {
		runnerSleep(duration)
		now = now.Add(duration)
	}
	lock, err := statefile.AcquireLock(path+".lock", statefile.LockOptions{
		Owner:          "api-statefile",
		StaleAfter:     30 * time.Second,
		AcquireTimeout: 10 * time.Second,
		RetryInterval:  100 * time.Millisecond,
		Clock:          func() time.Time { return now },
		Sleep:          sleep,
	})
	if err != nil {
		return func() {}, err
	}
	return func() { _ = lock.Release() }, nil
}

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}
