package perf

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gitee.com/zengtao321/frdocker/constants"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/entity/dto"

	units "github.com/docker/go-units"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func GetSystemPerformance(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(5)
	var memory *mem.VirtualMemoryStat
	go func() {
		memory, _ = mem.VirtualMemory()
		wg.Done()
	}()
	var cpuPhysicalCount, cpuLogicalCount int
	var cpuTotalPercent []float64
	var cpuInfo []cpu.InfoStat
	go func() {
		cpuPhysicalCount, _ = cpu.Counts(false)
		cpuLogicalCount, _ = cpu.Counts(true)
		cpuTotalPercent, _ = cpu.Percent(500*time.Millisecond, false)
		cpuInfo, _ = cpu.Info()
		wg.Done()
	}()
	var diskInfo *disk.UsageStat
	go func() {
		diskInfo, _ = disk.Usage("/var/lib/docker")
		wg.Done()
	}()
	var hostInfo *host.InfoStat
	go func() {
		hostInfo, _ = host.Info()
		wg.Done()
	}()
	var microServiceSytemInfo *dto.MicroServiceSytemInfo
	go func() {
		total := constants.IPServiceContainerMap.Count()
		healthCount := 0
		for _, obj := range constants.IPServiceContainerMap.Items() {
			container := obj.(*types.Container)
			if container.Health {
				healthCount += 1
			}
		}
		microServiceSytemInfo = &dto.MicroServiceSytemInfo{
			Network:       constants.Network,
			Registry:      constants.RegistryURL,
			ServiceGroups: constants.ServiceGroupMap.Count(),
			ServiceInstances: &dto.ServiceInstancesInfo{
				Total:     total,
				Healthy:   healthCount,
				UnHealthy: total - healthCount,
			},
			Gateways: constants.ServiceGroupMap.Count(),
		}
		wg.Done()
	}()
	wg.Wait()
	systemInfo := &dto.SystemPerfDTO{
		Memory: &dto.MemoryInfo{
			Total:          units.BytesSize(float64(memory.Total)),
			Available:      units.BytesSize(float64(memory.Available)),
			Used:           units.BytesSize(float64(memory.Used)),
			UsedPercentage: memory.UsedPercent,
		},
		Cpu: &dto.CpuInfo{
			PhysicalCount: cpuPhysicalCount,
			LogicalCount:  cpuLogicalCount,
			Percentage:    cpuTotalPercent[0],
			ModelName:     cpuInfo[0].ModelName,
		},
		Disk: &dto.DiskInfo{
			Total:          units.BytesSize(float64(diskInfo.Total)),
			Free:           units.BytesSize(float64(diskInfo.Free)),
			Used:           units.BytesSize(float64(diskInfo.Used)),
			UsedPercentage: diskInfo.UsedPercent,
		},
		Host: &dto.HostInfo{
			PlatForm:      fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion),
			Kernel:        fmt.Sprintf("%s %s", hostInfo.OS, hostInfo.KernelVersion),
			DockerVersion: utils.GetDockerVersion(),
		},
		MSSystem: microServiceSytemInfo,
	}
	c.JSON(http.StatusOK, R.OK(systemInfo))
}
