package perf

import (
	"net/http"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/entity/dto"

	"github.com/docker/go-units"
	"github.com/gin-gonic/gin"
)

// GetContainerPerformance 获取微服务容器性能参数
// @Summary 获取微服务容器性能参数
// @Description 获取微服务容器性能参数
// @Tags 系统性能操作
// @Produce application/json
// @Param Authorization	header	string	true	"token"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=dto.SystemPerfDTO} "返回微服务容器性能参数"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /perf/container [get]
func GetContainerPerformance(c *gin.Context) {
	IP := c.Query("ip")
	if IP == "" || !commons.IPServiceContainerMap.Has(IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "No such IP!", nil))
		return
	}
	obj, _ := commons.IPServiceContainerMap.Get(IP)
	container := obj.(*types.Container)
	statsEntry, err := utils.GetContainerStats(container.ID[:10])
	if err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Failed, try again later!", nil))
		return
	}
	containerPerfInfo := &dto.ContainerPerfDTO{
		Memory: &dto.MemoryInfo{
			Total:      units.BytesSize(statsEntry.MemoryLimit),
			Used:       units.BytesSize(statsEntry.Memory),
			Percentage: statsEntry.MemoryPercentage,
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
