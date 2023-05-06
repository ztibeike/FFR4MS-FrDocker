package entity

import (
	"sync"

	"gitee.com/zengtao321/frdocker/docker"
	"gitee.com/zengtao321/frdocker/utils"
)

type Container struct {
	Id            string       // 容器标识符(IP:Port)
	ContainerID   string       // 容器ID
	ContainerName string       // 容器名称
	IP            string       // 容器IP
	Port          int          // 容器端口
	IsHealthy     bool         // 容器是否健康
	ServiceName   string       // 容器所属服务名称
	mu            sync.RWMutex // 读写锁
}

func NewContainer(dockerCli *docker.DockerCLI, ip string, port int, serviceName string) (*Container, error) {
	container := &Container{
		Id:          utils.GenerateIdFromIPAndPort(ip, port),
		IP:          ip,
		Port:        port,
		ServiceName: serviceName,
		IsHealthy:   true,
	}
	err := setContainerInfoWithDockerCLI(dockerCli, container)
	return container, err
}

func setContainerInfoWithDockerCLI(dockerCli *docker.DockerCLI, container *Container) error {
	dockerContainer, err := dockerCli.GetContainerInfoByAddr(container.IP, container.Port)
	if err != nil {
		return err
	}
	container.ContainerID = dockerContainer.ID
	container.ContainerName = dockerContainer.Names[0]
	return nil
}