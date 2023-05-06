package entity

type ContainerStatus struct {
	Id                    string                     // 容器标识符(IP:Port)
	State                 map[string]*ContainerState // 容器状态
	StateAbnormalHandler  AbnormalHandlerFunc        // 容器状态异常处理函数
	Metric                *ContainerMetric           // 容器指标
	MetricAbnormalHandler AbnormalHandlerFunc        // 容器指标异常处理函数
}

func NewContainerStatus(id string) *ContainerStatus {
	return &ContainerStatus{
		Id:    id,
		State: make(map[string]*ContainerState),
	}
}
