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
	Services         map[string]*entity.Service
	Gateways         map[string]*entity.Service
	Containers       map[string]*entity.Container
}
