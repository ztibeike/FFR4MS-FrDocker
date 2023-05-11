package docker

import (
	"context"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types"
)

func (cli *DockerCLI) GetContainerStats(containerId string) (*DockerStats, error) {
	stats, err := cli.cli.ContainerStats(context.Background(), containerId, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()
	decoder := json.NewDecoder(stats.Body)
	var data *types.StatsJSON
	err = decoder.Decode(&data)
	if err != nil {
		data = nil
		decoder = json.NewDecoder(io.MultiReader(decoder.Buffered(), stats.Body))
		err = decoder.Decode(&data)
		if err != nil {
			return nil, err
		}
	}
	previousCPU := data.PreCPUStats.CPUUsage.TotalUsage
	previousSystem := data.PreCPUStats.SystemUsage
	cpuPercent := cli.calculateCPUPercentUnix(previousCPU, previousSystem, data)
	mem := cli.calculateMemUsageUnixNoCache(data.MemoryStats)
	memLimit := float64(data.MemoryStats.Limit)
	memPercent := cli.calculateMemPercentUnixNoCache(memLimit, mem)
	netDownload, netUpload := cli.calculateNetwork(data.Networks)
	diskRead, diskWrite := cli.calculateBlockIO(data.BlkioStats)
	return &DockerStats{
		ContainerId: containerId,
		CPUPercent:  cpuPercent,
		MemPercent:  memPercent,
		NetDownload: netDownload,
		NetUpload:   netUpload,
		DiskRead:    float64(diskRead),
		DiskWrite:   float64(diskWrite),
	}, nil
}

func (cli *DockerCLI) calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
		onlineCPUs  = float64(v.CPUStats.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(v.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}

func (cli *DockerCLI) calculateBlockIO(blkio types.BlkioStats) (uint64, uint64) {
	var blkRead, blkWrite uint64
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		if len(bioEntry.Op) == 0 {
			continue
		}
		switch bioEntry.Op[0] {
		case 'r', 'R':
			blkRead = blkRead + bioEntry.Value
		case 'w', 'W':
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	return blkRead, blkWrite
}

func (cli *DockerCLI) calculateMemUsageUnixNoCache(mem types.MemoryStats) float64 {
	// cgroup v1
	if v, isCgroup1 := mem.Stats["total_inactive_file"]; isCgroup1 && v < mem.Usage {
		return float64(mem.Usage - v)
	}
	// cgroup v2
	if v := mem.Stats["inactive_file"]; v < mem.Usage {
		return float64(mem.Usage - v)
	}
	return float64(mem.Usage)
}

func (cli *DockerCLI) calculateMemPercentUnixNoCache(limit float64, usedNoCache float64) float64 {
	// MemoryStats.Limit will never be 0 unless the container is not running and we haven't
	// got any data from cgroup
	if limit != 0 {
		return usedNoCache / limit * 100.0
	}
	return 0
}

func (cli *DockerCLI) calculateNetwork(network map[string]types.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}
