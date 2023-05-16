package frecovery

import (
	"sync"

	"gitee.com/zengtao321/frdocker/docker"
	"gitee.com/zengtao321/frdocker/utils"
)

type Container struct {
	Id            string            `json:"id" bson:"id"`                       // 容器标识符
	ContainerID   string            `json:"containerID " bson:"containerID"`    // 容器ID
	ContainerName string            `json:"containerName" bson:"containerName"` // 容器名称
	IP            string            `json:"ip" bson:"ip"`                       // 容器IP
	Port          int               `json:"port" bson:"port"`                   // 容器端口
	IsHealthy     bool              `json:"isHealthy" bson:"isHealthy"`         // 容器是否健康
	ServiceName   string            `json:"serviceName" bson:"serviceName"`     // 容器所属服务名称
	Monitor       *ContainerMonitor `json:"monitor" bson:"monitor"`             // 容器状态
	mu            sync.RWMutex      // 读写锁
}

func NewContainer(dockerCli *docker.DockerCLI, ip string, port int, serviceName string) (*Container, error) {
	id := utils.GenerateContainerId(ip, port)
	container := &Container{
		Id:          id,
		IP:          ip,
		Port:        port,
		ServiceName: serviceName,
		IsHealthy:   true,
	}
	err := container.setContainerInfoWithDockerCLI(dockerCli)
	container.Monitor = NewContainerMonitor(container.Id, container.ContainerID)
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
