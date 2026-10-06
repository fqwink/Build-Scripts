package setup

import (
	"encoding/json"
	"fmt"
	"io"
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

type setupAtomicWriteOps struct {
	createTemp func(string, string) (*os.File, error)
	rename     func(string, string) error
	remove     func(string) error
	open       func(string) (*os.File, error)
	sync       func(*os.File) error
}

func defaultSetupAtomicWriteOps() setupAtomicWriteOps {
	return setupAtomicWriteOps{
		createTemp: os.CreateTemp,
		rename:     os.Rename,
		remove:     os.Remove,
		open:       os.Open,
		sync:       func(file *os.File) error { return file.Sync() },
	}
}

func atomicWriteText(path, value string, mode os.FileMode) error {
	return atomicWriteTextWithOps(path, value, mode, defaultSetupAtomicWriteOps())
}

func atomicWriteTextWithOps(path, value string, mode os.FileMode, ops setupAtomicWriteOps) (resultErr error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("target is symlink")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	tmp, err := ops.createTemp(dir, "."+filepath.Base(path)+".tmp-")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	tmpClosed := false
	defer func() {
		if !tmpClosed {
			if err := tmp.Close(); resultErr == nil && err != nil {
				resultErr = err
			}
		}
		if err := ops.remove(tmpPath); resultErr == nil && err != nil && !os.IsNotExist(err) {
			resultErr = err
		}
	}()
	if err := tmp.Chmod(mode); err != nil {
		return err
	}
	if _, err := io.WriteString(tmp, value); err != nil {
		return err
	}
	if err := ops.sync(tmp); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	tmpClosed = true
	if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("target is symlink")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := ops.rename(tmpPath, path); err != nil {
		return err
	}
	directory, err := ops.open(dir)
	if err != nil {
		return err
	}
	if err := ops.sync(directory); err != nil {
		_ = directory.Close()
		return err
	}
	return directory.Close()
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
