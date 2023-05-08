package frecovery

import (
	"gitee.com/zengtao321/frdocker/docker"
	"gitee.com/zengtao321/frdocker/utils"
	"github.com/sirupsen/logrus"
)

type Container struct {
	Id            string           // 容器标识符(IP:Port)
	ContainerID   string           // 容器ID
	ContainerName string           // 容器名称
	IP            string           // 容器IP
	Port          int              // 容器端口
	IsHealthy     bool             // 容器是否健康
	ServiceName   string           // 容器所属服务名称
	Status        *ContainerStatus // 容器状态
	logger        *logrus.Logger   // 日志
}

func NewContainer(dockerCli *docker.DockerCLI, ip string, port int, serviceName string, logger *logrus.Logger) (*Container, error) {
	container := &Container{
		Id:          utils.GenerateContainerId(ip, port),
		IP:          ip,
		Port:        port,
		ServiceName: serviceName,
		IsHealthy:   true,
		Status:      NewContainerStatus(utils.GenerateContainerId(ip, port)),
		logger:      logger,
	}
	container.setStateMonitorHandler()
	container.setMetricMonitorHandler()
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

type StateMonitorHandlerFunc func(traceId string, state *ContainerState)

type MetricMonitorHandlerFunc func(metric *ContainerMetric)

// 设置状态监控回调函数
func (container *Container) setStateMonitorHandler() {
	stateMonitorHandler := func(traceId string, state *ContainerState) {
		// TODO
		state.mu.RLock()
		defer state.mu.RUnlock()
		ecc, thresh := state.Ecc, state.Thresh
		if ecc > thresh {
			container.logger.Errorf("[%s][%s:%d][%s] ecc: %.4f, thresh: %.4f", container.ServiceName, container.IP, container.Port, traceId, ecc, thresh)
			container.IsHealthy = false
		} else {
			container.logger.Tracef("[%s][%s:%d][%s] ecc: %.4f, thresh: %.4f", container.ServiceName, container.IP, container.Port, traceId, ecc, thresh)
		}
	}
	container.Status.SetStateMonitorHandler(stateMonitorHandler)
}

// 设置性能监控回调函数
func (container *Container) setMetricMonitorHandler() {
	metricMonitorHandler := func(metric *ContainerMetric) {
		// TODO
	}
	container.Status.SetMetricMonitorHandler(metricMonitorHandler)
}
