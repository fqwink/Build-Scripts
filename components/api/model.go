package api

type phase12Model struct {
	owner string
}

func newPhase12Model() phase12Model {
	return phase12Model{owner: "api"}
}
