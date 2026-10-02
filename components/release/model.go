package release

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "release",
		Files: []string{"release.go", "model.go", "validate.go", "execute.go", "release_test.go"},
	}
}
