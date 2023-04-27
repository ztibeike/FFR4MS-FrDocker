package docker

import (
	"context"

	"gitee.com/zengtao321/frdocker/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type DockerCLI struct {
	cli        *client.Client
	logger     *logrus.Logger
	containers map[string]types.Container
}

func NewDockerCLI(logger *logrus.Logger) *DockerCLI {
	dockerCli := &DockerCLI{logger: logger, containers: make(map[string]types.Container)}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatalln("Error while creating Docker Client: ", err)
	}
	dockerCli.cli = cli
	dockerCli.GetAllContainers()
	return dockerCli
}

func (dockerCli *DockerCLI) GetAllContainers() {
	containers, err := dockerCli.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		dockerCli.logger.Fatalln("Error while listing containers: ", err)
		return
	}
	for _, container := range containers {
		// 查到的是内网IP
		var ip string
		for _, network := range container.NetworkSettings.Networks {
			ip = network.IPAddress
			break
		}
		port := int(container.Ports[0].PrivatePort)
		key := utils.GenerateIdFromIPAndPort(ip, port)
		dockerCli.containers[key] = container
	}
}

func (dockerCli *DockerCLI) GetContainerInfoByAddr(ip string, port int) *types.Container {
	key := utils.GenerateIdFromIPAndPort(ip, port)
	if container, ok := dockerCli.containers[key]; ok {
		return &container
	}
	// 如果没有找到，重新初始化一次
	dockerCli.GetAllContainers()
	if container, ok := dockerCli.containers[key]; ok {
		return &container
	}
	dockerCli.logger.Fatalln("Can't find container by addr: %s", key)
	return nil
}
