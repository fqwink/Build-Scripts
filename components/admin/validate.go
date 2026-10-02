package admin

func validatePhase12Model(model phase12Model) bool {
	return model.owner == "admin"
}
