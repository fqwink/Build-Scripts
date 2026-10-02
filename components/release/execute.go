package release

import (
	"os"
	"syscall"
	"unicode/utf8"
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

func setupOwnedByCurrentUser(info os.FileInfo) bool {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return false
	}
	return stat.Uid == uint32(os.Geteuid())
}

func isLowerHex(value string, length int) bool {
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
