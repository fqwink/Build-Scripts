package commitstatus

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: Owner(),
		Files: []string{"commitstatus.go", "model.go", "validate.go", "execute.go", "commitstatus_test.go"},
	}
}
