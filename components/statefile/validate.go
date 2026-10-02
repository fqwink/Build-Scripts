package statefile

import (
	"errors"
	"os"
	"path/filepath"
)

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == "statefile" && validateExactOwnerFiles(contract.Files, []string{"statefile.go", "model.go", "validate.go", "execute.go", "statefile_test.go"})
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

func validateTargetPath(path string) error {
	if path == "" {
		return errors.New("statefile path must not be empty")
	}
	clean := filepath.Clean(path)
	if clean != path {
		return errors.New("statefile path must be clean")
	}
	parent := filepath.Dir(path)
	if parent == "" || parent == "." {
		return errors.New("statefile parent directory must be explicit")
	}
	if info, err := os.Lstat(parent); err != nil {
		return err
	} else if info.Mode()&os.ModeSymlink != 0 {
		return errors.New("statefile parent directory must not be symlink")
	} else if !info.IsDir() {
		return errors.New("statefile parent path is not directory")
	}
	if info, err := os.Lstat(path); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return errors.New("statefile target must not be symlink")
		}
		if !info.Mode().IsRegular() {
			return errors.New("statefile target must be regular file")
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
