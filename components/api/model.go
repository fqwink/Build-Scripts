package api

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "api",
		Files: []string{"api.go", "model.go", "validate.go", "execute.go", "api_test.go"},
	}
}
