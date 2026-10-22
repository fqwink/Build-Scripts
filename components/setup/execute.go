package setup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/fqwink/build-scripts/components/security"
)

func ReadCICDConfig(path string) (CICDConfig, error) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return CICDConfig{}, cicdConfigCapability
	}
	groups, err := os.Getgroups()
	if err != nil {
		return CICDConfig{}, cicdConfigCapability
	}
	var memberships []uint32
	for _, group := range groups {
		memberships = append(memberships, uint32(group))
	}
	return readCICDConfigAt("/", path, 0, uint32(os.Geteuid()), uint32(os.Getegid()), memberships...)
}

func readCICDConfigAt(rootPath, path string, rootUID, uid, gid uint32, groups ...uint32) (CICDConfig, error) {
	empty := CICDConfig{}
	if !cicdConfigPath(path) || filepath.Dir(path) != "/etc/adlaire-ci/cicd" {
		return empty, cicdConfigSchema
	}
	root, err := os.OpenFile(rootPath, os.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err != nil {
		return empty, cicdConfigCapability
	}
	defer root.Close()
	data, before, err := readCICDConfigPath(root, path, rootUID, 0640)
	if err != nil {
		return empty, err
	}
	cfg, err := ParseCICDConfig(data)
	if err != nil {
		return empty, err
	}
	if filepath.Base(path) != cfg.Role+"-"+cfg.ID+".json" {
		return empty, cicdConfigSchema
	}
	configGroup := before[len(before)-1].Sys().(*syscall.Stat_t).Gid
	member := configGroup == gid
	for _, group := range groups {
		member = member || configGroup == group
	}
	if configGroup == 0 || !member {
		return empty, cicdConfigCapability
	}
	for _, info := range before[len(before)-3:] {
		stat := info.Sys().(*syscall.Stat_t)
		if stat.Gid != configGroup || info.IsDir() && info.Mode() != os.ModeDir|0750 {
			return empty, cicdConfigCapability
		}
	}
	passwd, _, err := readCICDConfigPath(root, "/etc/passwd", rootUID, 0)
	if err != nil || !cicdConfigAccount(passwd, cfg, uid, gid) {
		return empty, cicdConfigCapability
	}
	again, after, err := readCICDConfigPath(root, path, rootUID, 0640)
	if err != nil || string(data) != string(again) || len(before) != len(after) {
		return empty, cicdConfigCapability
	}
	for i := range before {
		if !sameCICDConfigFile(before[i], after[i], !before[i].IsDir()) {
			return empty, cicdConfigCapability
		}
	}
	return cfg, nil
}

func cicdConfigAccount(passwd []byte, cfg CICDConfig, uid, gid uint32) bool {
	if uid == 0 || gid == 0 || len(passwd) > 1048576 || !utf8.Valid(passwd) || strings.ContainsRune(string(passwd), '\x00') {
		return false
	}
	names, uids, homes, matches := 0, 0, 0, 0
	for _, line := range strings.Split(string(passwd), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) != 7 {
			return false
		}
		userID, e1 := strconv.ParseUint(fields[2], 10, 32)
		groupID, e2 := strconv.ParseUint(fields[3], 10, 32)
		if e1 != nil || e2 != nil {
			return false
		}
		if fields[0] == cfg.User {
			names++
		}
		if userID == uint64(uid) {
			uids++
		}
		if fields[5] == cfg.StateDir {
			homes++
		}
		shell := fields[6] == "/bin/sh"
		if cfg.Role == "coordinator" {
			shell = fields[6] == "/usr/sbin/nologin" || fields[6] == "/sbin/nologin"
		}
		if fields[0] == cfg.User && userID == uint64(uid) && groupID == uint64(gid) && fields[5] == cfg.StateDir && shell {
			matches++
		}
	}
	return names == 1 && uids == 1 && homes == 1 && matches == 1
}

// Keep all resolution under a held root FD; never follow a path component.
func readCICDConfigPath(root *os.File, path string, rootUID uint32, mode os.FileMode) ([]byte, []os.FileInfo, error) {
	current := root
	var opened []*os.File
	defer func() {
		for _, file := range opened {
			_ = file.Close()
		}
	}()
	info, err := root.Stat()
	if err != nil || !validCICDConfigDirectory(info, rootUID) {
		return nil, nil, cicdConfigCapability
	}
	infos := []os.FileInfo{info}
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	for i, part := range parts {
		flags := os.O_RDONLY | syscall.O_NOFOLLOW | syscall.O_CLOEXEC | syscall.O_NONBLOCK
		last := i == len(parts)-1
		if !last {
			flags |= syscall.O_DIRECTORY
		}
		file, err := os.OpenFile(fmt.Sprintf("/proc/self/fd/%d/%s", current.Fd(), part), flags, 0)
		if err != nil {
			return nil, nil, cicdConfigCapability
		}
		opened = append(opened, file)
		info, err = file.Stat()
		if err != nil {
			return nil, nil, cicdConfigCapability
		}
		if !last {
			if !validCICDConfigDirectory(info, rootUID) {
				return nil, nil, cicdConfigCapability
			}
			infos = append(infos, info)
			current = file
			continue
		}
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok || !info.Mode().IsRegular() || stat.Uid != rootUID || stat.Nlink != 1 || info.Mode()&0022 != 0 || info.Mode()&^os.ModePerm != 0 || mode != 0 && info.Mode() != mode || info.Size() < 0 || info.Size() > 1048576 {
			return nil, nil, cicdConfigCapability
		}
		data, err := io.ReadAll(io.LimitReader(file, 1048577))
		after, statErr := file.Stat()
		if err != nil || statErr != nil || len(data) > 1048576 || int64(len(data)) != info.Size() || !sameCICDConfigFile(info, after, true) {
			return nil, nil, cicdConfigCapability
		}
		return data, append(infos, after), nil
	}
	return nil, nil, cicdConfigCapability
}

func validCICDConfigDirectory(info os.FileInfo, uid uint32) bool {
	stat, ok := info.Sys().(*syscall.Stat_t)
	return ok && info.IsDir() && info.Mode()&0022 == 0 && info.Mode()&^(os.ModeDir|os.ModePerm) == 0 && stat.Uid == uid
}

func sameCICDConfigFile(before, after os.FileInfo, content bool) bool {
	a, okA := before.Sys().(*syscall.Stat_t)
	b, okB := after.Sys().(*syscall.Stat_t)
	if !okA || !okB || !os.SameFile(before, after) || before.Mode() != after.Mode() || a.Uid != b.Uid || a.Gid != b.Gid || a.Nlink != b.Nlink {
		return false
	}
	if !content {
		return true
	}
	ctime := func(stat *syscall.Stat_t) any {
		value := reflect.ValueOf(stat).Elem().FieldByName("Ctim")
		if !value.IsValid() {
			value = reflect.ValueOf(stat).Elem().FieldByName("Ctimespec")
		}
		if !value.IsValid() {
			return nil
		}
		return value.Interface()
	}
	return before.Size() == after.Size() && before.ModTime() == after.ModTime() && ctime(a) != nil && reflect.DeepEqual(ctime(a), ctime(b))
}

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
