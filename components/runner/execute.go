package runner

import (
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func DefaultCICDResources() CICDResources {
	return CICDResources{CPUPercent: 100, MemoryMB: 256, DiskMB: 128, Tasks: 64, TimeoutSeconds: 900, OutputMB: 16}
}

// ReadCICDHostMemory observes physical memory, not the complete sandbox capability.
func ReadCICDHostMemory() (CICDHostMemory, error) {
	if runtime.GOOS != "linux" {
		return CICDHostMemory{}, ErrCICDCapability
	}
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return CICDHostMemory{}, ErrCICDCapability
	}
	memory, readErr := ParseCICDHostMemory(file)
	closeErr := file.Close()
	if readErr != nil || closeErr != nil {
		return CICDHostMemory{}, ErrCICDCapability
	}
	return memory, nil
}

func ParseCICDHostMemory(reader io.Reader) (CICDHostMemory, error) {
	if reader == nil {
		return CICDHostMemory{}, ErrCICDCapability
	}
	data, err := io.ReadAll(io.LimitReader(reader, cicdMeminfoMaxSize+1))
	if err != nil || len(data) > cicdMeminfoMaxSize {
		return CICDHostMemory{}, ErrCICDCapability
	}
	values := make(map[string]int64, 2)
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.FieldsFunc(line, func(r rune) bool { return r == ' ' || r == '\t' })
		if len(fields) == 0 || (fields[0] != "MemTotal:" && fields[0] != "MemAvailable:") {
			continue
		}
		key := fields[0]
		if _, duplicate := values[key]; duplicate || len(fields) != 3 || fields[2] != "kB" {
			return CICDHostMemory{}, ErrCICDCapability
		}
		for _, digit := range fields[1] {
			if digit < '0' || digit > '9' {
				return CICDHostMemory{}, ErrCICDCapability
			}
		}
		value, parseErr := strconv.ParseInt(fields[1], 10, 64)
		if parseErr != nil || value < 0 || value > cicdMaximumHostMB*1024 {
			return CICDHostMemory{}, ErrCICDCapability
		}
		values[key] = value
	}
	total, hasTotal := values["MemTotal:"]
	available, hasAvailable := values["MemAvailable:"]
	if !hasTotal || !hasAvailable || total == 0 || available > total {
		return CICDHostMemory{}, ErrCICDCapability
	}
	return CICDHostMemory{TotalMB: total / 1024, AvailableMB: available / 1024}, nil
}

type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string { return e.Msg }

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}
