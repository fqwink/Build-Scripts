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
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

func cicdValueBytes(value any) []byte {
	var out bytes.Buffer
	writeCICDCanonical(&out, value)
	return out.Bytes()
}

func cicdValueDigest(value any) string {
	return fmt.Sprintf("%x", sha256.Sum256(cicdValueBytes(value)))
}

func newCICDSnapshotStore(path string) (*cicdSnapshotStore, error) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return nil, ErrCICDCapability
	}
	dir, err := openCICDStateDirectory(path)
	if err != nil {
		return nil, ErrCICDCorrupt
	}
	defer dir.Close()
	return &cicdSnapshotStore{path: path}, nil
}

func openCICDStateDirectory(path string) (*os.File, error) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return nil, ErrCICDCapability
	}
	if !filepath.IsAbs(path) || filepath.Clean(path) != path || len(path) > 4096 || strings.Contains(path, "\\") || cicdHasControl(path) {
		return nil, ErrCICDCorrupt
	}
	dir, err := os.OpenFile("/", os.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, ErrCICDCorrupt
	}
	for _, part := range strings.Split(strings.TrimPrefix(path, "/"), "/") {
		if part == "" {
			dir.Close()
			return nil, ErrCICDCorrupt
		}
		next, err := os.OpenFile(cicdChildPath(dir, part), os.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
		dir.Close()
		if err != nil {
			return nil, ErrCICDCorrupt
		}
		info, err := next.Stat()
		if err != nil {
			next.Close()
			return nil, ErrCICDCorrupt
		}
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok || (stat.Uid != uint32(os.Geteuid()) && stat.Uid != 0) || info.Mode().Perm()&0022 != 0 && !(stat.Uid == 0 && info.Mode()&os.ModeSticky != 0) {
			next.Close()
			return nil, ErrCICDCorrupt
		}
		dir = next
	}
	if !validCICDFile(dir, true) {
		dir.Close()
		return nil, ErrCICDCorrupt
	}
	return dir, nil
}

func cicdChildPath(dir *os.File, name string) string {
	return "/proc/self/fd/" + strconv.FormatUint(uint64(dir.Fd()), 10) + "/" + name
}

func validCICDFile(file *os.File, directory bool) bool {
	info, err := file.Stat()
	if err != nil {
		return false
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || stat.Uid != uint32(os.Geteuid()) {
		return false
	}
	if directory {
		return info.IsDir() && info.Mode().Perm() == 0700 && info.Mode()&(os.ModeSetuid|os.ModeSetgid|os.ModeSticky) == 0
	}
	return info.Mode().IsRegular() && info.Mode().Perm() == 0600 && info.Mode()&(os.ModeSetuid|os.ModeSetgid|os.ModeSticky) == 0 && stat.Nlink == 1
}

func openCICDChild(dir *os.File, name string, flags int, directory bool) (*os.File, error) {
	if name == "" || strings.ContainsAny(name, "/\\\x00") || name == "." || name == ".." {
		return nil, ErrCICDCorrupt
	}
	mode := os.FileMode(0600)
	if directory {
		flags |= syscall.O_DIRECTORY
		mode = 0700
	}
	file, err := os.OpenFile(cicdChildPath(dir, name), flags|syscall.O_NOFOLLOW|syscall.O_CLOEXEC|syscall.O_NONBLOCK, mode)
	if err != nil {
		return nil, err
	}
	if !validCICDFile(file, directory) {
		file.Close()
		return nil, ErrCICDCorrupt
	}
	return file, nil
}

func (s *cicdSnapshotStore) directories(create bool) (*os.File, *os.File, error) {
	parent, err := openCICDStateDirectory(s.path)
	if err != nil {
		return nil, nil, err
	}
	defer parent.Close()
	if create {
		if err := os.Mkdir(cicdChildPath(parent, "cicd"), 0700); err != nil && !errors.Is(err, os.ErrExist) {
			return nil, nil, ErrCICDStorage
		}
	}
	dir, err := openCICDChild(parent, "cicd", os.O_RDONLY, true)
	if err != nil {
		return nil, nil, ErrCICDCorrupt
	}
	if create {
		if err := os.Mkdir(cicdChildPath(dir, "snapshots"), 0700); err != nil && !errors.Is(err, os.ErrExist) {
			dir.Close()
			return nil, nil, ErrCICDStorage
		}
		if parent.Sync() != nil || dir.Sync() != nil {
			dir.Close()
			return nil, nil, ErrCICDStorage
		}
	}
	snapshots, err := openCICDChild(dir, "snapshots", os.O_RDONLY, true)
	if err != nil {
		dir.Close()
		return nil, nil, ErrCICDCorrupt
	}
	return dir, snapshots, nil
}

func (s *cicdSnapshotStore) inject(stage string) error {
	if s.fault != nil {
		return s.fault(stage)
	}
	return nil
}

func readCICDBounded(dir *os.File, name string, limit int64) ([]byte, error) {
	file, err := openCICDChild(dir, name, os.O_RDONLY, false)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || info.Size() > limit {
		return nil, ErrCICDCorrupt
	}
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil || int64(len(data)) > limit {
		return nil, ErrCICDCorrupt
	}
	return data, nil
}

func readCICDHead(dir *os.File) (cicdHead, error) {
	data, err := readCICDBounded(dir, "head.json", 4096)
	if err != nil || ValidateCICDValue("Head", data) != nil {
		return cicdHead{}, ErrCICDCorrupt
	}
	var head cicdHead
	if json.Unmarshal(data, &head) != nil {
		return cicdHead{}, ErrCICDCorrupt
	}
	return head, nil
}

func (s *cicdSnapshotStore) read() (map[string]any, cicdHead, error) {
	dir, snapshots, err := s.directories(false)
	if err != nil {
		return nil, cicdHead{}, err
	}
	defer dir.Close()
	defer snapshots.Close()
	return s.readFrom(dir, snapshots)
}

func (s *cicdSnapshotStore) readFrom(dir, snapshots *os.File) (map[string]any, cicdHead, error) {
	s.observed.Lock()
	defer s.observed.Unlock()
	if err := s.inspectFiles(dir, snapshots, 0, false); err != nil {
		return nil, cicdHead{}, ErrCICDCorrupt
	}
	for retry := 0; retry < 3; retry++ {
		head, err := readCICDHead(dir)
		if err != nil {
			return nil, cicdHead{}, ErrCICDCorrupt
		}
		if s.inject("read.snapshot") != nil {
			return nil, cicdHead{}, ErrCICDCorrupt
		}
		data, err := readCICDBounded(snapshots, "g-"+strconv.FormatInt(head.Generation, 10)+".json", cicdSnapshotLimit)
		if errors.Is(err, os.ErrNotExist) {
			latest, headErr := readCICDHead(dir)
			if headErr == nil && latest.Generation > head.Generation {
				continue
			}
			return nil, cicdHead{}, ErrCICDCorrupt
		}
		if err != nil || fmt.Sprintf("%x", sha256.Sum256(data)) != head.SHA256 {
			return nil, cicdHead{}, ErrCICDCorrupt
		}
		value, err := decodeCICDValue(data)
		root, ok := value.(map[string]any)
		if err != nil || !ok || !validateCICDRoot(root) || cicdInt(root["generation"]) != head.Generation {
			return nil, cicdHead{}, ErrCICDCorrupt
		}
		if !validCICDReservedCapacity(root, len(data)) {
			return nil, cicdHead{}, ErrCICDCorrupt
		}
		regressed := head.Generation < s.generation || head.Generation == s.generation && s.digest != "" && head.SHA256 != s.digest
		if !regressed {
			s.generation, s.digest = head.Generation, head.SHA256
		}
		if regressed {
			return nil, cicdHead{}, ErrCICDCorrupt
		}
		return root, head, nil
	}
	return nil, cicdHead{}, ErrCICDBusy
}

func (s *cicdSnapshotStore) lock(create bool) (*os.File, *os.File, func(), error) {
	deadline := time.Now().Add(2 * time.Second)
	for !cicdSnapshotMutex.TryLock() {
		if !time.Now().Before(deadline) {
			return nil, nil, nil, ErrCICDBusy
		}
		time.Sleep(10 * time.Millisecond)
	}
	dir, snapshots, err := s.directories(create)
	if err != nil {
		cicdSnapshotMutex.Unlock()
		return nil, nil, nil, err
	}
	flags := os.O_RDWR
	if create {
		flags |= os.O_CREATE
	}
	lock, err := openCICDChild(dir, "state.lock", flags, false)
	if err != nil {
		dir.Close()
		snapshots.Close()
		cicdSnapshotMutex.Unlock()
		return nil, nil, nil, ErrCICDCorrupt
	}
	unlock := func() {
		_ = syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
		_ = lock.Close()
		_ = snapshots.Close()
		_ = dir.Close()
		cicdSnapshotMutex.Unlock()
	}
	info, err := lock.Stat()
	if err != nil || info.Size() != 0 {
		unlock()
		return nil, nil, nil, ErrCICDCorrupt
	}
	for {
		err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			break
		}
		if err != syscall.EWOULDBLOCK && err != syscall.EAGAIN && err != syscall.EINTR {
			unlock()
			return nil, nil, nil, ErrCICDCorrupt
		}
		if !time.Now().Before(deadline) {
			unlock()
			return nil, nil, nil, ErrCICDBusy
		}
		time.Sleep(10 * time.Millisecond)
	}
	return dir, snapshots, unlock, nil
}

func (s *cicdSnapshotStore) initialize(data []byte, opID string) (cicdHead, error) {
	value, err := decodeCICDValue(data)
	root, ok := value.(map[string]any)
	if err != nil || !ok || !validCICDEntityID(opID, "o-") || !validateCICDRoot(root) || cicdInt(root["generation"]) != 1 {
		return cicdHead{}, ErrCICDCorrupt
	}
	dir, snapshots, unlock, err := s.lock(true)
	if err != nil {
		return cicdHead{}, err
	}
	defer unlock()
	if _, err := os.Lstat(cicdChildPath(dir, "head.json")); !errors.Is(err, os.ErrNotExist) {
		return cicdHead{}, ErrCICDCorrupt
	}
	if err := s.reclaim(dir, snapshots, 0); err != nil {
		return cicdHead{}, err
	}
	return s.publish(dir, snapshots, root, opID)
}

func (s *cicdSnapshotStore) commit(expectedGeneration int64, opID string, mutate func(map[string]any) error) (cicdHead, error) {
	if !validCICDEntityID(opID, "o-") || mutate == nil {
		return cicdHead{}, ErrCICDCorrupt
	}
	dir, snapshots, unlock, err := s.lock(false)
	if err != nil {
		return cicdHead{}, err
	}
	defer unlock()
	s.observed.Lock()
	uncertain := s.uncertain
	s.observed.Unlock()
	if uncertain {
		return cicdHead{}, ErrCICDUncertain
	}
	root, head, err := s.readFrom(dir, snapshots)
	if err != nil {
		return cicdHead{}, err
	}
	if head.Generation != expectedGeneration || head.Generation >= 9007199254740991 {
		return cicdHead{}, ErrCICDRevision
	}
	if err := s.reclaim(dir, snapshots, head.Generation); err != nil {
		return cicdHead{}, err
	}
	before := cicdValueBytes(root)
	if err := mutate(root); err != nil {
		return cicdHead{}, err
	}
	// Round-trip through the strict decoder before canonical writing caller values.
	data, err := json.Marshal(root)
	if err != nil {
		return cicdHead{}, ErrCICDCorrupt
	}
	value, err := decodeCICDValue(data)
	if err != nil {
		return cicdHead{}, ErrCICDCorrupt
	}
	root = cicdObject(value)
	if !validateCICDRoot(root) || cicdInt(root["generation"]) != head.Generation {
		return cicdHead{}, ErrCICDCorrupt
	}
	if bytes.Equal(before, cicdValueBytes(root)) {
		return head, nil
	}
	root["generation"] = json.Number(strconv.FormatInt(head.Generation+1, 10))
	return s.publish(dir, snapshots, root, opID)
}

func (s *cicdSnapshotStore) publish(dir, snapshots *os.File, root map[string]any, opID string) (cicdHead, error) {
	data := cicdValueBytes(root)
	if !validCICDReservedCapacity(root, len(data)) {
		return cicdHead{}, ErrCICDStorage
	}
	var stat syscall.Statfs_t
	if syscall.Fstatfs(int(dir.Fd()), &stat) != nil || stat.Bsize <= 0 || stat.Bavail < uint64((256*1024*1024+int64(stat.Bsize)-1)/int64(stat.Bsize)) {
		return cicdHead{}, ErrCICDStorage
	}
	if s.inject("capacity") != nil {
		return cicdHead{}, ErrCICDStorage
	}
	head := cicdHead{SchemaVersion: 1, Generation: cicdInt(root["generation"]), SHA256: fmt.Sprintf("%x", sha256.Sum256(data))}
	if err := s.writeCICDFile(snapshots, "."+opID+".tmp", "g-"+strconv.FormatInt(head.Generation, 10)+".json", data, "snapshot", false); err != nil {
		return cicdHead{}, err
	}
	headData, err := json.Marshal(head)
	if err != nil {
		return cicdHead{}, ErrCICDStorage
	}
	headData, err = CanonicalCICDJSON(headData)
	if err != nil {
		return cicdHead{}, ErrCICDStorage
	}
	if err := s.writeCICDFile(dir, ".head."+opID+".tmp", "head.json", headData, "head", true); err != nil {
		if err == ErrCICDUncertain {
			s.observed.Lock()
			s.uncertain = true
			s.observed.Unlock()
		}
		return cicdHead{}, err
	}
	s.observed.Lock()
	s.generation, s.digest = head.Generation, head.SHA256
	s.observed.Unlock()
	if err := s.reclaim(dir, snapshots, head.Generation); err != nil && s.diagnostic != nil {
		s.diagnostic(ErrCICDCleanup)
	}
	return head, nil
}

func (s *cicdSnapshotStore) writeCICDFile(dir *os.File, tmp, name string, data []byte, stage string, replace bool) error {
	if s.inject(stage+".create") != nil {
		return ErrCICDStorage
	}
	file, err := openCICDChild(dir, tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, false)
	if err != nil {
		return ErrCICDStorage
	}
	closed := false
	defer func() {
		if !closed {
			_ = file.Close()
		}
		_ = os.Remove(cicdChildPath(dir, tmp))
	}()
	if s.inject(stage+".write") != nil {
		return ErrCICDStorage
	}
	for len(data) > 0 {
		n, err := file.Write(data)
		if err != nil || n == 0 {
			return ErrCICDStorage
		}
		data = data[n:]
	}
	if s.inject(stage+".chmod") != nil || file.Chmod(0600) != nil {
		return ErrCICDStorage
	}
	if s.inject(stage+".sync") != nil || file.Sync() != nil {
		return ErrCICDStorage
	}
	if s.inject(stage+".close") != nil {
		return ErrCICDStorage
	}
	if file.Close() != nil {
		closed = true
		return ErrCICDStorage
	}
	closed = true
	if s.inject(stage+".publish") != nil {
		return ErrCICDStorage
	}
	if replace {
		err = os.Rename(cicdChildPath(dir, tmp), cicdChildPath(dir, name))
	} else {
		err = renameCICDNoReplace(dir, tmp, name)
	}
	if err != nil {
		return ErrCICDStorage
	}
	if s.inject(stage+".directory-sync") != nil || dir.Sync() != nil {
		if replace {
			return ErrCICDUncertain
		}
		return ErrCICDStorage
	}
	return nil
}

func renameCICDNoReplace(dir *os.File, from, to string) error {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return ErrCICDCapability
	}
	old, err := syscall.BytePtrFromString(from)
	if err != nil {
		return err
	}
	new, err := syscall.BytePtrFromString(to)
	if err != nil {
		return err
	}
	// Linux amd64 renameat2, RENAME_NOREPLACE. Never emulate by check-then-rename.
	_, _, errno := syscall.Syscall6(316, dir.Fd(), uintptr(unsafe.Pointer(old)), dir.Fd(), uintptr(unsafe.Pointer(new)), 1, 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func (s *cicdSnapshotStore) reclaim(dir, snapshots *os.File, generation int64) error {
	if s.inject("cleanup") != nil {
		return ErrCICDCleanup
	}
	return s.inspectFiles(dir, snapshots, generation, true)
}

func (s *cicdSnapshotStore) inspectFiles(dir, snapshots *os.File, generation int64, prune bool) error {
	for _, parent := range []*os.File{dir, snapshots} {
		listing, err := os.Open(cicdChildPath(parent, "."))
		if err != nil {
			return ErrCICDCorrupt
		}
		names, err := listing.Readdirnames(-1)
		listing.Close()
		if err != nil {
			return ErrCICDCorrupt
		}
		for _, name := range names {
			remove := false
			if parent == dir {
				switch name {
				case "snapshots", "secrets", "logs", "artifacts", "cache", "work", "ssh", "site", "locks":
					child, err := openCICDChild(parent, name, os.O_RDONLY, true)
					if err != nil {
						return ErrCICDCorrupt
					}
					child.Close()
					continue
				case "head.json", "state.lock", "master.json", "worker.json", "worker-jobs.json", "setup.json":
					child, err := openCICDChild(parent, name, os.O_RDONLY, false)
					if err != nil {
						return ErrCICDCorrupt
					}
					info, err := child.Stat()
					child.Close()
					if err != nil || name == "state.lock" && info.Size() != 0 {
						return ErrCICDCorrupt
					}
					continue
				}
				if strings.HasPrefix(name, ".head.") && strings.HasSuffix(name, ".tmp") && validCICDEntityID(strings.TrimSuffix(strings.TrimPrefix(name, ".head."), ".tmp"), "o-") {
					remove = true
				} else {
					return ErrCICDCorrupt
				}
			} else if strings.HasPrefix(name, "g-") && strings.HasSuffix(name, ".json") {
				text := strings.TrimSuffix(strings.TrimPrefix(name, "g-"), ".json")
				n, err := strconv.ParseInt(text, 10, 64)
				if err != nil || n < 1 || n > 9007199254740991 || strconv.FormatInt(n, 10) != text {
					return ErrCICDCorrupt
				}
				remove = n > generation || n < generation-2
			} else if strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".tmp") && validCICDEntityID(strings.TrimSuffix(strings.TrimPrefix(name, "."), ".tmp"), "o-") {
				remove = true
			} else {
				return ErrCICDCorrupt
			}
			if !prune {
				continue
			}
			file, err := openCICDChild(parent, name, os.O_RDONLY, false)
			if err != nil {
				return ErrCICDCorrupt
			}
			file.Close()
			if prune && remove && os.Remove(cicdChildPath(parent, name)) != nil {
				return ErrCICDCleanup
			}
		}
		if prune && parent.Sync() != nil {
			return ErrCICDCleanup
		}
	}
	return nil
}

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
