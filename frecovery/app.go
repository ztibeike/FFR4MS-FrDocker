package frecovery

import (
	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/docker"
	"github.com/google/gopacket/pcap"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type FrecoveryApp struct {
	RegistryAddress  string                `json:"registryAddress" bson:"registryAddress"`
	NetworkInterface string                `json:"networkInterface" bson:"networkInterface"`
	DockerCli        *docker.DockerCLI     `json:"-" bson:"-"`
	Logger           *logrus.Logger        `json:"-" bson:"-"`
	DbCli            *mongo.Database       `json:"-" bson:"-"`
	PcapHandle       *pcap.Handle          `json:"-" bson:"-"`
	Services         map[string]*Service   `json:"services" bson:"services"`     // key: serviceName, value: service
	Gateways         map[string]*Service   `json:"gateways" bson:"gateways"`     // key: gatewayName, value: gateway
	Containers       map[string]*Container `json:"containers" bson:"containers"` // key: ip:port, value: container
	Pool             *ants.Pool            `json:"-" bson:"-"`
}

func NewFrecoveryApp(registryAdress string, networkInterface string, dockerCli *docker.DockerCLI, logger *logrus.Logger, dbCli *mongo.Database) *FrecoveryApp {
	app := &FrecoveryApp{
		RegistryAddress:  registryAdress,
		NetworkInterface: networkInterface,
		DockerCli:        dockerCli,
		Logger:           logger,
		DbCli:            dbCli,
		Services:         make(map[string]*Service),
		Gateways:         make(map[string]*Service),
		Containers:       make(map[string]*Container),
	}
	app.Pool, _ = ants.NewPool(config.FRECOVERY_GOROUTINE_POOL_COUNT)
	return app
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
