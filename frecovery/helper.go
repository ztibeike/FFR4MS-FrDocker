package frecovery

import (
	"gitee.com/zengtao321/frdocker/frecovery/entity"
)

// 获取容器类型
func (app *FrecoveryApp) getContainerType(id string) entity.ContainerType {
	if _, ok := app.Containers[id]; !ok {
		return entity.CTN_INVALID
	}
	serviceName := app.Containers[id].ServiceName
	if _, ok := app.Services[serviceName]; ok {
		return entity.CTN_SERVICE
	}
	if _, ok := app.Gateways[serviceName]; ok {
		return entity.CTN_GATEWAY
	}
	return entity.CTN_INVALID
}
