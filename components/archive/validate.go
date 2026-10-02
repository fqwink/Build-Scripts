package archive

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == "archive" && validateExactOwnerFiles(contract.Files, []string{"archive.go", "model.go", "validate.go", "execute.go", "archive_test.go"})
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
