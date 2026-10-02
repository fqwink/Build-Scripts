package runner

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
