package frecovery

import (
	"strconv"
	"strings"

	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/types/dto"
	"gitee.com/zengtao321/frdocker/utils"
	"github.com/google/gopacket/pcap"
)

// 初始化微服务系统中的services和containers
func (app *FrecoveryApp) initMSSystem() {
	app.Logger.Info("init microservice system...")
	registryConfig, err := utils.GetRegistryInfo(app.RegistryAddress)
	if err != nil {
		app.Logger.Fatal("error while getting config from registry: ", err)
		return
	}
	app.initServicesAndGateways(registryConfig.Services, app.Services)
	app.initServicesAndGateways(registryConfig.Gateways, app.Gateways)
	app.setGatewayForServices(registryConfig.Gateways)
	app.Logger.Info("init microservice system success")
}

func (app *FrecoveryApp) initServicesAndGateways(src map[string][]dto.MSInstance, dst map[string]*Service) {
	for key, value := range src {
		key = strings.ToLower(key)
		service := NewService(key)
		for _, msInstance := range value {
			leaf, ok := msInstance.Metadata[config.REGISTRY_METADATA_LEAF_KEY]
			if ok {
				service.IsLeaf, _ = strconv.ParseBool(leaf)
			}
			container, err := NewContainer(app.DockerCli, msInstance.IP, msInstance.Port, service.ServiceName)
			service.Containers = append(service.Containers, container.Id)
			if err != nil {
				app.Logger.Fatalf("error while init container of %s:%s:%d: %s", service.ServiceName, msInstance.IP, msInstance.Port, err)
			}
			app.Containers[container.Id] = container
		}
		dst[service.ServiceName] = service
	}
}

func (app *FrecoveryApp) setGatewayForServices(gateway map[string][]dto.MSInstance) {
	for key, value := range gateway {
		gatewayName := strings.ToLower(key)
		for _, instance := range value {
			serviceName := instance.Metadata[config.REGISTRY_METADATA_GATEWAY_KEY]
			serviceName = strings.ToLower(serviceName)
			service, ok := app.Services[serviceName]
			if !ok {
				continue
			}
			service.Gateway = gatewayName
		}
	}
}

// 初始化网卡监控pcapHandle
func (app *FrecoveryApp) initPcap() {
	app.Logger.Info("init pcap...")
	// 初始化pcapHandle
	handle, err := pcap.OpenLive(app.NetworkInterface, 65535, true, pcap.BlockForever)
	if err != nil {
		app.Logger.Fatal("error while init pcap handle: ", err)
		return
	}
	app.PcapHandle = handle
	// 设置过滤器
	filter := "tcp"
	if err := handle.SetBPFFilter(filter); err != nil {
		app.Logger.Fatal("error while set pcap filter: ", err)
		return
	}
	app.Logger.Info("init pcap success")
}
