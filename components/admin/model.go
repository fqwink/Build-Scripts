package admin

type ownerFileContract struct {
	Owner string
	Files []string
}

func newOwnerFileContract() ownerFileContract {
	return ownerFileContract{
		Owner: "admin",
		Files: []string{"admin.go", "model.go", "validate.go", "execute.go", "admin_test.go"},
	}
}
