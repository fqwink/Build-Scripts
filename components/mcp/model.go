package mcp

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "mcp",
		Files: []string{"mcp.go", "model.go", "validate.go", "execute.go", "mcp_test.go"},
	}
}
