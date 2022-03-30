package perf

import (
	"frdocker/constants"
	"frdocker/types"
	"frdocker/utils"
	"frdocker/web/entity/R"
	"frdocker/web/entity/dto"
	"net/http"

	"github.com/docker/go-units"
	"github.com/gin-gonic/gin"
)

func GetContainerPerformance(c *gin.Context) {
	IP := c.Query("ip")
	if IP == "" || !constants.IPServiceContainerMap.Has(IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "No such IP!", nil))
		return
	}
	obj, _ := constants.IPServiceContainerMap.Get(IP)
	container := obj.(*types.Container)
	statsEntry, err := utils.GetContainerStats(container.ID[:10])
	if err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Failed, try again later!", nil))
		return
	}
	containerPerfInfo := &dto.ContainerPerfDTO{
		Memory: &dto.MemoryInfo{
			Total:          units.BytesSize(statsEntry.MemoryLimit),
			Used:           units.BytesSize(statsEntry.Memory),
			Available:      units.BytesSize(statsEntry.MemoryLimit - statsEntry.Memory),
			UsedPercentage: statsEntry.MemoryPercentage,
		},
		Cpu: &dto.CpuInfo{
			Percentage: statsEntry.CPUPercentage,
		},
		NetworkTransfer: &dto.NetworkTransferInfo{
			Upload:   units.BytesSize(statsEntry.NetworkTx),
			Download: units.BytesSize(statsEntry.NetworkRx),
		},
		BlockIO: &dto.BlockIOInfo{
			Read:  units.BytesSize(statsEntry.BlockRead),
			Write: units.BytesSize(statsEntry.BlockWrite),
		},
	}
	c.JSON(http.StatusOK, R.OK(containerPerfInfo))
}
