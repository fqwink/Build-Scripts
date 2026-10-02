package runner

type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string { return e.Msg }

func executePhase12Model(model phase12Model) bool {
	return validatePhase12Model(model)
}
