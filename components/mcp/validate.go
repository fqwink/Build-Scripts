package mcp

func validateOwnerFileContract(contract ownerFileContract) bool {
	return contract.Owner == "mcp" && validateExactOwnerFiles(contract.Files, []string{"mcp.go", "model.go", "validate.go", "execute.go", "mcp_test.go"})
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
