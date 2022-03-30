package dto

type SystemPerfDTO struct {
	Memory   *MemoryInfo            `json:"memory"`
	Cpu      *CpuInfo               `json:"cpu"`
	Disk     *DiskInfo              `json:"disk"`
	Host     *HostInfo              `json:"host"`
	MSSystem *MicroServiceSytemInfo `json:"mssystem"`
}

type ContainerPerfDTO struct {
	Memory          *MemoryInfo          `json:"memory"`
	Cpu             *CpuInfo             `json:"cpu"`
	NetworkTransfer *NetworkTransferInfo `json:"networkTransfer"`
	BlockIO         *BlockIOInfo         `json:"blockIO"`
}

type MemoryInfo struct {
	Total          string  `json:"total"`
	Available      string  `json:"available"`
	Used           string  `json:"used"`
	UsedPercentage float64 `json:"usedPercentage"`
}

type CpuInfo struct {
	PhysicalCount int     `json:"physicalCount"`
	LogicalCount  int     `json:"logicalCount"`
	Percentage    float64 `json:"percentage"`
	ModelName     string  `json:"modelName"`
}

type DiskInfo struct {
	Total          string  `json:"total"`
	Free           string  `json:"free"`
	Used           string  `json:"used"`
	UsedPercentage float64 `json:"usedPercentage"`
}

type HostInfo struct {
	PlatForm      string `json:"platForm"`
	Kernel        string `json:"kernel"`
	DockerVersion string `json:"dockerVersion"`
}

type MicroServiceSytemInfo struct {
	Network          string                `json:"network"`
	Registry         string                `json:"registry"`
	ServiceGroups    int                   `json:"serviceGroups"`
	ServiceInstances *ServiceInstancesInfo `json:"serviceInstances"`
	Gateways         int                   `json:"gateways"`
}

type ServiceInstancesInfo struct {
	Total     int `json:"total"`
	Healthy   int `json:"healthy"`
	UnHealthy int `json:"unHealthy"`
}

type NetworkTransferInfo struct {
	Upload   string `json:"upload"`
	Download string `json:"download"`
}

type BlockIOInfo struct {
	Read  string `json:"read"`
	Write string `json:"write"`
}
