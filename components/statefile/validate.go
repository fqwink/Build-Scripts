package statefile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"unicode/utf8"
)

func validateCICDJSONBytes(data []byte) bool {
	if !utf8.Valid(data) || !json.Valid(data) {
		return false
	}
	// encoding/json replaces unpaired surrogates; identities must not be repaired.
	for i := 0; i < len(data); {
		if data[i] != '"' {
			i++
			continue
		}
		i++
		for i < len(data) && data[i] != '"' {
			if data[i] != '\\' {
				i++
				continue
			}
			if data[i+1] != 'u' {
				i += 2
				continue
			}
			value, err := strconv.ParseUint(string(data[i+2:i+6]), 16, 16)
			if err != nil {
				return false
			}
			if value >= 0xd800 && value <= 0xdbff {
				if i+12 > len(data) || string(data[i+6:i+8]) != `\u` {
					return false
				}
				low, err := strconv.ParseUint(string(data[i+8:i+12]), 16, 16)
				if err != nil || low < 0xdc00 || low > 0xdfff {
					return false
				}
				i += 12
			} else if value >= 0xdc00 && value <= 0xdfff {
				return false
			} else {
				i += 6
			}
		}
		i++
	}
	return true
}

func validateCICDInteger(number json.Number) bool {
	if number.String() == "-0" {
		return false
	}
	value, err := strconv.ParseInt(number.String(), 10, 64)
	return err == nil && value >= -9007199254740991 && value <= 9007199254740991
}

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
