package frecovery

import (
	"gitee.com/zengtao321/frdocker/docker"
	"gitee.com/zengtao321/frdocker/frecovery/entity"
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
	Services         map[string]*entity.Service   // key: serviceName, value: service
	Gateways         map[string]*entity.Service   // key: gatewayName, value: gateway
	Containers       map[string]*entity.Container // key: ip:port, value: container
}

func NewFrecoveryApp(registryURL string, networkInterface string, dockerCli *docker.DockerCLI, logger *logrus.Logger, dbCli *mongo.Database) *FrecoveryApp {
	return &FrecoveryApp{
		RegistryURL:      registryURL,
		NetworkInterface: networkInterface,
		DockerCli:        dockerCli,
		Logger:           logger,
		DbCli:            dbCli,
		Services:         make(map[string]*entity.Service),
		Gateways:         make(map[string]*entity.Service),
		Containers:       make(map[string]*entity.Container),
	}
}
