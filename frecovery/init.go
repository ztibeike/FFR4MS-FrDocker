package frecovery

import (
	"strconv"
	"strings"

	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/frecovery/entity"
	"gitee.com/zengtao321/frdocker/types/dto"
	"github.com/go-resty/resty/v2"
	"github.com/google/gopacket/pcap"
)

// 初始化微服务系统中的services和containers
func (app *FrecoveryApp) initMSSystem() {
	app.Logger.Info("init microservice system...")
	client := resty.New()
	registryConfig := dto.MSConfig{}
	_, err := client.R().SetHeader("Accept", "application/json").SetResult(&registryConfig).Get(app.RegistryURL)
	if err != nil {
		app.Logger.Fatal("error while getting config from registry: ", err)
		return
	}
	app.initServicesAndGateways(registryConfig.Services, app.Services)
	app.initServicesAndGateways(registryConfig.Gateways, app.Gateways)
	app.Logger.Info("init microservice system success")
}

func (app *FrecoveryApp) initServicesAndGateways(src map[string][]dto.MSInstance, dst map[string]*entity.Service) {
	for key, value := range src {
		key = strings.ToLower(key)
		service := entity.NewService(key)
		for _, msInstance := range value {
			leaf, ok := msInstance.Metadata[config.REGISTRY_METADATA_LEAF_KEY]
			if ok {
				service.IsLeaf, _ = strconv.ParseBool(leaf)
			}
			service.Containers = append(service.Containers, msInstance.Address)
			container, err := entity.NewContainer(app.DockerCli, msInstance.IP, msInstance.Port, service.ServiceName)
			if err != nil {
				app.Logger.Errorf("error while init container of %s:%s:%d: %s", service.ServiceName, msInstance.IP, msInstance.Port, err)
			}
			app.Containers[container.Id] = container
		}
		dst[service.ServiceName] = service
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
