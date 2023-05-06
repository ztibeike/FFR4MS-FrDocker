package entity

type ContainerMetric struct {
	CPU             float64             // CPU使用率
	Mem             float64             // 内存使用率
	Net             float64             // 网络使用率
	Disk            float64             // 磁盘使用率
	Ecc             float64             // 离心率
	abnormalHandler AbnormalHandlerFunc // 异常处理函数
}

func NewContainerMetric(abnormalHandler AbnormalHandlerFunc) *ContainerMetric {
	return &ContainerMetric{
		CPU:             0.0,
		Mem:             0.0,
		Net:             0.0,
		Disk:            0.0,
		Ecc:             0.0,
		abnormalHandler: abnormalHandler,
	}
}
