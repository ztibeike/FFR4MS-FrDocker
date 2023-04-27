package frecovery

import (
	"strconv"
	"strings"

	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/frecovery/entity"
	"gitee.com/zengtao321/frdocker/types"
	"github.com/go-resty/resty/v2"
)

// 初始化containers和services
func (app *FrecoveryApp) initApp() {
	app.Logger.Info("init container...")
	client := resty.New()
	registryConfig := types.MSConfig{}
	_, err := client.R().SetHeader("Accept", "application/json").SetResult(&registryConfig).Get(app.RegistryURL)
	if err != nil {
		app.Logger.Fatalln("Error while getting config from registry: ", err)
		return
	}
	for key, value := range registryConfig.Services {
		key = strings.ToLower(key)
		service := entity.NewService(key)
		for _, msInstance := range value {
			leaf, ok := msInstance.Metadata[config.REGISTRY_METADATA_LEAF_KEY]
			if ok {
				service.IsLeaf, _ = strconv.ParseBool(leaf)
			}
			service.Containers = append(service.Containers, msInstance.Address)
			container := entity.NewContainer(app.DockerCli, msInstance.IP, msInstance.Port, service.ServiceName)
			app.Containers[container.Id] = container
		}
		app.Services[service.ServiceName] = service
	}
	for key, value := range registryConfig.Gateways {
		key = strings.ToLower(key)
		gateway := entity.NewService(key)
		for _, msInstance := range value {
			targetService, ok := msInstance.Metadata[config.REGISTRY_METADATA_GATEWAT_KEY]
			if ok {
				app.Services[targetService].Gateway = gateway.ServiceName
			}
			gateway.Containers = append(gateway.Containers, msInstance.Address)
		}
		app.Gateways[gateway.ServiceName] = gateway
	}
	app.Logger.Infoln("init container success")
}
