package setup

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "setup",
		Files: []string{"setup.go", "model.go", "validate.go", "execute.go", "setup_test.go"},
	}
}
