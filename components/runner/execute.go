package runner

type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string { return e.Msg }

func executeOwnerFileContract(contract ownerFileContract) bool {
	return validateOwnerFileContract(contract)
}
