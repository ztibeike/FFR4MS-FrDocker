package entity

import (
	"gitee.com/zengtao321/frdocker/docker"
	"gitee.com/zengtao321/frdocker/utils"
)

type AbnormalHandlerFunc func()

type Container struct {
	Id            string           // 容器标识符(IP:Port)
	ContainerID   string           // 容器ID
	ContainerName string           // 容器名称
	IP            string           // 容器IP
	Port          int              // 容器端口
	IsHealthy     bool             // 容器是否健康
	ServiceName   string           // 容器所属服务名称
	Status        *ContainerStatus // 容器状态
}

func NewContainer(dockerCli *docker.DockerCLI, ip string, port int, serviceName string) (*Container, error) {
	container := &Container{
		Id:          utils.GenerateIdFromIPAndPort(ip, port),
		IP:          ip,
		Port:        port,
		ServiceName: serviceName,
		IsHealthy:   true,
		Status:      NewContainerStatus(utils.GenerateIdFromIPAndPort(ip, port)),
	}
	container.Status.StateAbnormalHandler = container.setStateAbnormalHandler()
	container.Status.MetricAbnormalHandler = container.setMetricAbnormalHandler()
	err := container.setContainerInfoWithDockerCLI(dockerCli)
	return container, err
}

func (container *Container) setContainerInfoWithDockerCLI(dockerCli *docker.DockerCLI) error {
	dockerContainer, err := dockerCli.GetContainerInfoByAddr(container.IP, container.Port)
	if err != nil {
		return err
	}
	container.ContainerID = dockerContainer.ID
	container.ContainerName = dockerContainer.Names[0]
	return nil
}

func (container *Container) setStateAbnormalHandler() AbnormalHandlerFunc {
	return func() {
		// TODO
	}
}

func (container *Container) setMetricAbnormalHandler() AbnormalHandlerFunc {
	return func() {
		// TODO
	}
}
