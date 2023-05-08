package frecovery

type ContainerMetric struct {
	Id             string                   // 容器标识符
	CPU            float64                  // CPU使用率
	Mem            float64                  // 内存使用率
	Net            float64                  // 网络使用率
	Disk           float64                  // 磁盘使用率
	Ecc            float64                  // 离心率
	monitorHandler MetricMonitorHandlerFunc // 监控回调函数
}

func NewContainerMetric(containerId string, monitorHandler MetricMonitorHandlerFunc) *ContainerMetric {
	return &ContainerMetric{
		Id:             containerId,
		CPU:            0.0,
		Mem:            0.0,
		Net:            0.0,
		Disk:           0.0,
		Ecc:            0.0,
		monitorHandler: monitorHandler,
	}
}
