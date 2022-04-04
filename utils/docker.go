package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"gitee.com/zengtao321/frdocker/types"

	"github.com/ahmetb/dlog"
	dockerTypes "github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
)

var dockerClient, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
var ctx = context.Background()

func GetServiceContainers(containers []*types.Container) {
	originContainers, err := dockerClient.ContainerList(ctx, dockerTypes.ContainerListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	var originContainersMap = make(map[string]dockerTypes.Container)
	for _, originContainer := range originContainers {
		var IP string
		for _, v := range originContainer.NetworkSettings.Networks {
			IP = v.IPAddress
			break
		}
		originContainersMap[IP] = originContainer
	}
	for _, container := range containers {
		if originContainer, ok := originContainersMap[container.IP]; ok {
			container.ID = originContainer.ID
			container.Name = originContainer.Names[0]
			// containers[idx] = container
		}
	}
}

func GetContainerLogs(containerId string, tail string) (string, error) {
	options := dockerTypes.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	}
	resp, err := dockerClient.ContainerLogs(ctx, containerId, options)
	if err != nil {
		return "", err
	}
	defer resp.Close()
	rr := dlog.NewReader(resp)
	s := bufio.NewScanner(rr)
	var result string
	for s.Scan() {
		result = result + s.Text() + "\n"
		if err := s.Err(); err != nil {
			return "", err
		}
	}
	return result, nil
}

func GetDockerVersion() string {
	result, _ := dockerClient.Info(ctx)
	return fmt.Sprintf("Docker Engine %s Community", result.ServerVersion)
}

func GetContainerStats(containerId string) (*types.StatsEntry, error) {
	response, err := dockerClient.ContainerStats(ctx, containerId, false)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	var data *dockerTypes.StatsJSON
	err = dec.Decode(&data)
	if err != nil {
		data = nil
		dec = json.NewDecoder(io.MultiReader(dec.Buffered(), response.Body))
		err = dec.Decode(&data)
		if err != nil {
			return nil, err
		}
	}
	previousCPU := data.PreCPUStats.CPUUsage.TotalUsage
	previousSystem := data.PreCPUStats.SystemUsage
	cpuPercent := calculateCPUPercentUnix(previousCPU, previousSystem, data)
	blkRead, blkWrite := calculateBlockIO(data.BlkioStats)
	mem := calculateMemUsageUnixNoCache(data.MemoryStats)
	memLimit := float64(data.MemoryStats.Limit)
	memPercent := calculateMemPercentUnixNoCache(memLimit, mem)
	pidsStatsCurrent := data.PidsStats.Current
	netRx, netTx := calculateNetwork(data.Networks)
	return &types.StatsEntry{
		Name:             data.Name,
		ID:               data.ID,
		CPUPercentage:    cpuPercent,
		Memory:           mem,
		MemoryPercentage: memPercent,
		MemoryLimit:      memLimit,
		NetworkRx:        netRx,
		NetworkTx:        netTx,
		BlockRead:        float64(blkRead),
		BlockWrite:       float64(blkWrite),
		PidsCurrent:      pidsStatsCurrent,
	}, nil
}

func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *dockerTypes.StatsJSON) float64 {
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

func calculateBlockIO(blkio dockerTypes.BlkioStats) (uint64, uint64) {
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

func calculateMemUsageUnixNoCache(mem dockerTypes.MemoryStats) float64 {
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

func calculateMemPercentUnixNoCache(limit float64, usedNoCache float64) float64 {
	// MemoryStats.Limit will never be 0 unless the container is not running and we haven't
	// got any data from cgroup
	if limit != 0 {
		return usedNoCache / limit * 100.0
	}
	return 0
}

func calculateNetwork(network map[string]dockerTypes.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}
