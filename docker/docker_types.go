package docker

type DockerStats struct {
	ContainerId string
	CPUPercent  float64
	MemPercent  float64
	NetUpload   float64
	NetDownload float64
	DiskRead    float64
	DiskWrite   float64
}
