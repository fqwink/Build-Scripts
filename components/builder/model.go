package builder

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "builder",
		Files: []string{"builder.go", "model.go", "validate.go", "execute.go", "builder_test.go"},
	}
}
