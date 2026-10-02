package api

func validatePhase12Model(model phase12Model) bool {
	return model.owner == "api"
}
