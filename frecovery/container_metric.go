package frecovery

import (
	"sync"

	"gitee.com/zengtao321/frdocker/docker"
)

type ContainerMetric struct {
	Id          string  `json:"id" bson:"id"`                   // 容器标识符
	ContainerId string  `json:"containerId" bson:"containerId"` // 容器ID
	CPU         float64 `json:"cpu" bson:"cpu"`                 // CPU使用率
	Mem         float64 `json:"mem" bson:"mem"`                 // 内存使用率
	NetUp       float64 `json:"netUp" bson:"netUp"`             // 网络上传量
	NetDn       float64 `json:"netDn" bson:"netDn"`             // 网络下载量
	DiskR       float64 `json:"diskR" bson:"diskR"`             // 磁盘读取量
	DiskW       float64 `json:"diskW" bson:"diskW"`             // 磁盘写入量
	Ecc         float64 `json:"ecc" bson:"ecc"`                 // 离心率
	Thresh      float64 `json:"thresh" bson:"thresh"`           // 阈值
	mu          sync.RWMutex
}

func NewContainerMetric(id, containerId string) *ContainerMetric {
	return &ContainerMetric{
		Id:          id,
		ContainerId: containerId,
		CPU:         0.0,
		Mem:         0.0,
		NetUp:       0.0,
		NetDn:       0.0,
		DiskR:       0.0,
		DiskW:       0.0,
		Ecc:         0.0,
		Thresh:      0.0,
	}
}

func (metric *ContainerMetric) Update(dockerCli *docker.DockerCLI) error {
	metric.mu.Lock()
	defer metric.mu.Unlock()
	containerStats, err := dockerCli.GetContainerStats(metric.ContainerId)
	if err != nil {
		return err
	}
	metric.CPU = containerStats.CPUPercent
	metric.Mem = containerStats.MemPercent
	metric.NetUp = containerStats.NetUpload
	metric.NetDn = containerStats.NetDownload
	metric.DiskR = containerStats.DiskRead
	metric.DiskW = containerStats.DiskWrite
	return nil
}

func (metric *ContainerMetric) UpdateEcc(ecc, thresh float64) {
	metric.mu.Lock()
	defer metric.mu.Unlock()
	metric.Ecc = ecc
	metric.Thresh = thresh
}
