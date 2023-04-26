package frecovery

import (
	"gitee.com/zengtao321/frdocker/frecovery/entity"
	"gitee.com/zengtao321/frdocker/utils/docker"
	"gitee.com/zengtao321/frdocker/utils/logger"
)

type FrecoveryApp struct {
	RegistryURL      string
	NetworkInterface string
	DockerCli        *docker.DockerCLI
	Logger           *logger.Logger
	Services         map[string]*entity.Service   // key: serviceName, value: service
	Gateways         map[string]*entity.Service   // key: gatewayName, value: gateway
	Containers       map[string]*entity.Container // key: ip:port, value: container
}

func NewFrecoveryApp(registryURL string, networkInterface string, dockerCli *docker.DockerCLI, logger *logger.Logger) *FrecoveryApp {
	return &FrecoveryApp{
		RegistryURL:      registryURL,
		NetworkInterface: networkInterface,
		DockerCli:        dockerCli,
		Logger:           logger,
		Services:         make(map[string]*entity.Service),
		Gateways:         make(map[string]*entity.Service),
		Containers:       make(map[string]*entity.Container),
	}
}
