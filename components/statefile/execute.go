package statefile

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func CanonicalCICDJSON(data []byte) ([]byte, error) {
	if !validateCICDJSONBytes(data) {
		return nil, ErrCICDSchema
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := readCICDJSONValue(decoder)
	if err != nil {
		return nil, ErrCICDSchema
	}
	if _, err := decoder.Token(); err != io.EOF {
		return nil, ErrCICDSchema
	}
	var out bytes.Buffer
	writeCICDCanonical(&out, value)
	return out.Bytes(), nil
}

func DigestCICDJSON(data []byte) (string, error) {
	canonical, err := CanonicalCICDJSON(data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(canonical)), nil
}

func readCICDJSONValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, ErrCICDSchema
	}
	switch value := token.(type) {
	case json.Delim:
		switch value {
		case '{':
			object := make(map[string]any)
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return nil, ErrCICDSchema
				}
				key, ok := keyToken.(string)
				if !ok {
					return nil, ErrCICDSchema
				}
				if _, exists := object[key]; exists {
					return nil, ErrCICDSchema
				}
				child, err := readCICDJSONValue(decoder)
				if err != nil {
					return nil, err
				}
				object[key] = child
			}
			if end, err := decoder.Token(); err != nil || end != json.Delim('}') {
				return nil, ErrCICDSchema
			}
			return object, nil
		case '[':
			array := make([]any, 0)
			for decoder.More() {
				child, err := readCICDJSONValue(decoder)
				if err != nil {
					return nil, err
				}
				array = append(array, child)
			}
			if end, err := decoder.Token(); err != nil || end != json.Delim(']') {
				return nil, ErrCICDSchema
			}
			return array, nil
		default:
			return nil, ErrCICDSchema
		}
	case json.Number:
		if !validateCICDInteger(value) {
			return nil, ErrCICDSchema
		}
		return value, nil
	case string, bool, nil:
		return value, nil
	default:
		return nil, ErrCICDSchema
	}
}

func writeCICDCanonical(out *bytes.Buffer, value any) {
	switch value := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		out.WriteByte('{')
		for i, key := range keys {
			if i != 0 {
				out.WriteByte(',')
			}
			writeCICDString(out, key)
			out.WriteByte(':')
			writeCICDCanonical(out, value[key])
		}
		out.WriteByte('}')
	case []any:
		out.WriteByte('[')
		for i, item := range value {
			if i != 0 {
				out.WriteByte(',')
			}
			writeCICDCanonical(out, item)
		}
		out.WriteByte(']')
	case string:
		writeCICDString(out, value)
	case json.Number:
		out.WriteString(value.String())
	case bool:
		if value {
			out.WriteString("true")
		} else {
			out.WriteString("false")
		}
	case nil:
		out.WriteString("null")
	}
}

func writeCICDString(out *bytes.Buffer, value string) {
	const hex = "0123456789abcdef"
	out.WriteByte('"')
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '"', '\\':
			out.WriteByte('\\')
			out.WriteByte(c)
		case '\b':
			out.WriteString(`\b`)
		case '\t':
			out.WriteString(`\t`)
		case '\n':
			out.WriteString(`\n`)
		case '\f':
			out.WriteString(`\f`)
		case '\r':
			out.WriteString(`\r`)
		default:
			if c < 0x20 {
				out.WriteString(`\u00`)
				out.WriteByte(hex[c>>4])
				out.WriteByte(hex[c&15])
			} else {
				out.WriteByte(c)
			}
		}
	}
	out.WriteByte('"')
}

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}

func ReadJSON(path string, out any) error {
	if err := validateTargetPath(path); err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func ReadJSONIfExists(path string, out any) error {
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	return ReadJSON(path, out)
}

func WriteJSONAtomic(path string, value any, mode os.FileMode) error {
	lock, err := AcquireLock(path+".lock", LockOptions{Owner: "json-atomic-write", StaleAfter: 30 * time.Second, AcquireTimeout: 10 * time.Second, RetryInterval: 100 * time.Millisecond})
	if err != nil {
		return err
	}
	defer lock.Release()
	return WriteJSONAtomicLocked(path, value, mode)
}

func WriteJSONAtomicLocked(path string, value any, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	if err := validateTargetPath(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.OpenFile(path+".tmp", os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanupTmp := true
	defer func() {
		if cleanupTmp {
			_ = os.Remove(tmpName)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	cleanupTmp = false
	return fsyncParent(path)
}

func AppendJSONLine(path string, value any, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	lock, err := AcquireLock(path+".lock", LockOptions{Owner: "jsonl-append", StaleAfter: 30 * time.Second, AcquireTimeout: 10 * time.Second, RetryInterval: 100 * time.Millisecond})
	if err != nil {
		return err
	}
	defer lock.Release()
	if err := RecoverJSONLIfCorrupt(path); err != nil {
		return err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, mode)
	if err != nil {
		return err
	}
	if info, err := f.Stat(); err == nil && info.Size() > 0 {
		var last [1]byte
		if _, err := f.ReadAt(last[:], info.Size()-1); err != nil {
			_ = f.Close()
			return err
		}
		if last[0] != '\n' {
			if _, err := f.Write([]byte("\n")); err != nil {
				_ = f.Close()
				return err
			}
		}
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Chmod(mode); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return fsyncParent(path)
}

func WriteJSONLinesAtomic(path string, values []any, mode os.FileMode) error {
	lock, err := AcquireLock(path+".lock", LockOptions{Owner: "jsonl-atomic-write", StaleAfter: 30 * time.Second, AcquireTimeout: 10 * time.Second, RetryInterval: 100 * time.Millisecond})
	if err != nil {
		return err
	}
	defer lock.Release()
	var builder strings.Builder
	for _, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		builder.Write(data)
		builder.WriteByte('\n')
	}
	return writeRawAtomic(path, []byte(builder.String()), mode)
}

func ReadJSONLines(path string) ([]map[string]any, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []map[string]any{}, nil
	}
	if err != nil {
		return nil, err
	}
	records := []map[string]any{}
	for index, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, fmt.Errorf("jsonl corrupt line %d: %w", index+1, err)
		}
		records = append(records, record)
	}
	return records, nil
}

func DetectJSONLCorruption(path string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for index, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		var value any
		if err := json.Unmarshal([]byte(line), &value); err != nil {
			return fmt.Errorf("jsonl corrupt line %d: %w", index+1, err)
		}
	}
	return nil
}

func RecoverJSONLIfCorrupt(path string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	validLines := []string{}
	corrupt := false
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		var value any
		if err := json.Unmarshal([]byte(line), &value); err != nil {
			corrupt = true
			continue
		}
		validLines = append(validLines, line)
	}
	if !corrupt {
		return nil
	}
	quarantinePath := fmt.Sprintf("%s.corrupt.%d", path, time.Now().UTC().UnixNano())
	if err := writeRawAtomic(quarantinePath, data, 0600); err != nil {
		return err
	}
	var builder strings.Builder
	for _, line := range validLines {
		builder.WriteString(line)
		builder.WriteByte('\n')
	}
	if err := writeRawAtomic(path, []byte(builder.String()), 0600); err != nil {
		return err
	}
	return nil
}

func AcquireLock(path string, options LockOptions) (*LockHandle, error) {
	if options.Owner == "" {
		options.Owner = "statefile"
	}
	if options.RetryInterval <= 0 {
		options.RetryInterval = 10 * time.Millisecond
	}
	deadline := lockNow(options).Add(options.AcquireTimeout)
	for {
		handle, err := tryAcquireLock(path, options)
		if err == nil {
			return handle, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		if options.StaleAfter > 0 {
			if removed, staleErr := removeStaleLock(path, options.StaleAfter, lockNow(options)); staleErr != nil {
				return nil, staleErr
			} else if removed {
				continue
			}
		}
		if options.AcquireTimeout <= 0 || !lockNow(options).Add(options.RetryInterval).Before(deadline) {
			return nil, err
		}
		lockSleep(options, options.RetryInterval)
	}
}

func (h *LockHandle) Release() error {
	if h == nil || h.file == "" {
		return nil
	}
	if err := os.Remove(h.file); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return fsyncParent(h.file)
}

func tryAcquireLock(path string, options LockOptions) (*LockHandle, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	data := []byte(fmt.Sprintf("owner=%s\npid=%d\ncreated_at=%s\n", options.Owner, os.Getpid(), lockNow(options).UTC().Format(time.RFC3339Nano)))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return nil, err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return nil, err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(path)
		return nil, err
	}
	if err := fsyncParent(path); err != nil {
		_ = os.Remove(path)
		return nil, err
	}
	return &LockHandle{path: filepath.Dir(path), file: path}, nil
}

func removeStaleLock(path string, staleAfter time.Duration, now time.Time) (bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return false, errors.New("statefile lock must not be symlink")
	}
	if now.Sub(info.ModTime()) <= staleAfter {
		return false, nil
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	return true, fsyncParent(path)
}

func lockNow(options LockOptions) time.Time {
	if options.Clock != nil {
		return options.Clock()
	}
	return time.Now()
}

func lockSleep(options LockOptions, duration time.Duration) {
	if options.Sleep != nil {
		options.Sleep(duration)
		return
	}
	time.Sleep(duration)
}

func fsyncParent(path string) error {
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}

func writeRawAtomic(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	if err := validateTargetPath(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	tmp, err := os.OpenFile(path+".tmp", os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanupTmp := true
	defer func() {
		if cleanupTmp {
			_ = os.Remove(tmpName)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	cleanupTmp = false
	return fsyncParent(path)
}
