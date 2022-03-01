package utils

import (
	"context"
	"frdocker/types"
	"log"

	dockerTypes "github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
)

var dockerClient, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
var ctx = context.Background()

func GetServiceContainers(containers []types.Container) {
	originContainers, err := dockerClient.ContainerList(ctx, dockerTypes.ContainerListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	var originContainersMap = make(map[string]dockerTypes.Container)
	for _, originContainer := range originContainers {
		var IP string
		for _, v := range originContainer.NetworkSettings.Networks {
			IP = v.IPAddress
			break
		}
		originContainersMap[IP] = originContainer
	}
	for idx, container := range containers {
		if originContainer, ok := originContainersMap[container.IP]; ok {
			container.ID = originContainer.ID
			container.Name = originContainer.Names[0]
			containers[idx] = container
		}
	}
}
