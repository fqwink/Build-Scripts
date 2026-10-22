package runner

func ValidateCICDResources(r CICDResources) error {
	if r.CPUPercent < 10 || r.CPUPercent > 400 ||
		r.MemoryMB < cicdMinimumJobMB || r.MemoryMB > cicdMaximumJobMB ||
		r.DiskMB < 128 || r.DiskMB > 2048 ||
		r.Tasks < 8 || r.Tasks > 256 ||
		r.TimeoutSeconds < 1 || r.TimeoutSeconds > 3600 ||
		r.OutputMB < 1 || r.OutputMB > 16 {
		return ErrCICDSchema
	}
	if r.MemoryMB < r.DiskMB+128 {
		return ErrCICDSchema
	}
	return nil
}

// CheckCICDMemory requires a fresh observation and the caller's claim lock.
// Reservations include unstarted claims and claims awaiting confirmed cleanup.
func CheckCICDMemory(memory CICDHostMemory, reservations []int64, configuredSlots, requestedMB int64) error {
	if configuredSlots < 1 || configuredSlots > 4 || requestedMB < cicdMinimumJobMB || requestedMB > cicdMaximumJobMB {
		return ErrCICDSchema
	}
	if memory.TotalMB < cicdMinimumHostMB || memory.TotalMB > cicdMaximumHostMB ||
		memory.AvailableMB < 0 || memory.AvailableMB > memory.TotalMB {
		return ErrCICDCapability
	}
	if int64(len(reservations)) > configuredSlots {
		return ErrCICDStateCorrupt
	}
	var reservedMB int64
	for _, value := range reservations {
		if value < cicdMinimumJobMB || value > cicdMaximumJobMB {
			return ErrCICDStateCorrupt
		}
		reservedMB += value
	}
	budget := memory.TotalMB - cicdHostReserveMB
	if requestedMB > budget {
		return ErrCICDCapability
	}
	slots := configuredSlots
	if memory.TotalMB < cicdSingleSlotMB {
		slots = 1
	}
	if int64(len(reservations)) >= slots || reservedMB+requestedMB > budget ||
		memory.AvailableMB < reservedMB+requestedMB+cicdMemoryGuardMB {
		return ErrCICDBusy
	}
	return nil
}

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == "runner" && validateExactOwnerFiles(contract.Files, []string{"runner.go", "model.go", "validate.go", "execute.go", "runner_test.go"})
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
