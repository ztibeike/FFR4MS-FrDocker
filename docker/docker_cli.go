package docker

import (
	"context"
	"fmt"

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

func NewDockerCLI(logger *logrus.Logger) (*DockerCLI, error) {
	dockerCli := &DockerCLI{logger: logger, containers: make(map[string]types.Container)}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	dockerCli.cli = cli
	err = dockerCli.GetAllContainers()
	if err != nil {
		return nil, err
	}
	return dockerCli, nil
}

func (cli *DockerCLI) GetAllContainers() error {
	containers, err := cli.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
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
		cli.containers[key] = container
	}
	return nil
}

func (cli *DockerCLI) GetContainerInfoByAddr(ip string, port int) (*types.Container, error) {
	key := utils.GenerateIdFromIPAndPort(ip, port)
	if container, ok := cli.containers[key]; ok {
		return &container, nil
	}
	// 如果没有找到，重新初始化一次
	err := cli.GetAllContainers()
	if err != nil {
		return nil, err
	}
	if container, ok := cli.containers[key]; ok {
		return &container, nil
	}
	return nil, fmt.Errorf("can not find container by addr: %s", key)
}
