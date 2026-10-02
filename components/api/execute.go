package api

import (
	"errors"
	"fmt"
	"os"
	"time"
	"unicode/utf8"
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
	lockPath := path + ".lock"
	var lastErr error
	for attempt := 0; attempt <= 100; attempt++ {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			_, writeErr := fmt.Fprintf(f, "pid=%d\nstarted_at=%s\n", os.Getpid(), time.Now().UTC().Format(time.RFC3339))
			closeErr := f.Close()
			if writeErr != nil || closeErr != nil {
				_ = os.Remove(lockPath)
				if writeErr != nil {
					return func() {}, writeErr
				}
				return func() {}, closeErr
			}
			return func() { _ = os.Remove(lockPath) }, nil
		}
		lastErr = err
		if !errors.Is(err, os.ErrExist) || attempt == 100 {
			break
		}
		runnerSleep(100 * time.Millisecond)
	}
	return func() {}, lastErr
}

func executePhase12Model(model phase12Model) bool {
	return validatePhase12Model(model)
}
