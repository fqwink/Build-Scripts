package archive

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: Owner(),
		Files: []string{"archive.go", "model.go", "validate.go", "execute.go", "archive_test.go"},
	}
}
