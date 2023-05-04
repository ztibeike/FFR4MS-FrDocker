package frecovery

func (app *FrecoveryApp) getContainerType(id string) ContainerType {
	if _, ok := app.Containers[id]; !ok {
		return INVALID
	}
	serviceName := app.Containers[id].ServiceName
	if _, ok := app.Services[serviceName]; ok {
		return SERVICE
	}
	if _, ok := app.Gateways[serviceName]; ok {
		return GATEWAY
	}
	return INVALID
}
