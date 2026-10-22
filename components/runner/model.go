package runner

import "errors"

var (
	ErrCICDSchema       = errors.New("P19_SCHEMA")
	ErrCICDCapability   = errors.New("P19_CAPABILITY")
	ErrCICDStateCorrupt = errors.New("P19_STATE_CORRUPT")
	ErrCICDBusy         = errors.New("P19_BUSY")
)

type CICDResources struct {
	CPUPercent     int64 `json:"cpu_percent"`
	MemoryMB       int64 `json:"memory_mb"`
	DiskMB         int64 `json:"disk_mb"`
	Tasks          int64 `json:"tasks"`
	TimeoutSeconds int64 `json:"timeout_seconds"`
	OutputMB       int64 `json:"output_mb"`
}

type CICDHostMemory struct {
	TotalMB     int64
	AvailableMB int64
}

const (
	cicdMinimumHostMB  int64 = 896
	cicdMaximumHostMB  int64 = 1048576
	cicdHostReserveMB  int64 = 512
	cicdSingleSlotMB   int64 = 1536
	cicdMemoryGuardMB  int64 = 64
	cicdMinimumJobMB   int64 = 256
	cicdMaximumJobMB   int64 = 8192
	cicdMeminfoMaxSize       = 65536
)

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "runner",
		Files: []string{"runner.go", "model.go", "validate.go", "execute.go", "runner_test.go"},
	}
}
