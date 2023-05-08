package frecovery

import (
	"gitee.com/zengtao321/frdocker/docker"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type FrecoveryApp struct {
	RegistryURL      string
	NetworkInterface string
	DockerCli        *docker.DockerCLI
	Logger           *logrus.Logger
	DbCli            *mongo.Database
	PcapHandle       *pcap.Handle
	Services         map[string]*Service   // key: serviceName, value: service
	Gateways         map[string]*Service   // key: gatewayName, value: gateway
	Containers       map[string]*Container // key: ip:port, value: container
}

func NewFrecoveryApp(registryURL string, networkInterface string, dockerCli *docker.DockerCLI, logger *logrus.Logger, dbCli *mongo.Database) *FrecoveryApp {
	return &FrecoveryApp{
		RegistryURL:      registryURL,
		NetworkInterface: networkInterface,
		DockerCli:        dockerCli,
		Logger:           logger,
		DbCli:            dbCli,
		Services:         make(map[string]*Service),
		Gateways:         make(map[string]*Service),
		Containers:       make(map[string]*Container),
	}
}

// 获取容器
func (app *FrecoveryApp) GetContainer(id string) *Container {
	if _, ok := app.Containers[id]; !ok {
		return nil
	}
	return app.Containers[id]
}

// 获取服务
func (app *FrecoveryApp) GetService(serviceName string) *Service {
	if _, ok := app.Services[serviceName]; !ok {
		return nil
	}
	return app.Services[serviceName]
}

// 获取网关
func (app *FrecoveryApp) GetGateway(gatewayName string) *Service {
	if _, ok := app.Gateways[gatewayName]; !ok {
		return nil
	}
	return app.Gateways[gatewayName]
}

// 获取容器类型
func (app *FrecoveryApp) getContainerType(id string) ContainerType {
	if _, ok := app.Containers[id]; !ok {
		return CTN_INVALID
	}
	serviceName := app.Containers[id].ServiceName
	if _, ok := app.Services[serviceName]; ok {
		return CTN_SERVICE
	}
	if _, ok := app.Gateways[serviceName]; ok {
		return CTN_GATEWAY
	}
	return CTN_INVALID
}
