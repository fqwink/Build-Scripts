package api

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == "api" && validateExactOwnerFiles(contract.Files, []string{"api.go", "model.go", "validate.go", "execute.go", "api_test.go"})
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
