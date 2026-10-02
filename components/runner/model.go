package runner

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "runner",
		Files: []string{"runner.go", "model.go", "validate.go", "execute.go", "runner_test.go"},
	}
}
