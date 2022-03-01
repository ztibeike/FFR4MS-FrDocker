package examples

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ContainerList() {
	var dockerClient *client.Client
	var err error
	// var containers []types.Container
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	originContainers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	for _, container := range originContainers {
		// containers = append(containers, types.Container {
		// 	IP: container.NetworkSettings.Networks[container.HostConfig.NetworkMode].IPAddress,
		// 	ID: container.ID,
		// })
		fmt.Println(container.NetworkSettings.Networks[container.HostConfig.NetworkMode].IPAddress)
	}
}
