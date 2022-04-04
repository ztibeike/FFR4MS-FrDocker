package container

import (
	"frdocker/constants"
	"frdocker/types"
	"frdocker/utils"
	"frdocker/web/entity/R"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetContainerLogs(c *gin.Context) {
	IP := c.Query("ip")
	if IP == "" || !constants.IPServiceContainerMap.Has(IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "No Such IP!", nil))
		return
	}
	tail := c.Query("tail")
	if tail == "" {
		tail = "all"
	}
	obj, _ := constants.IPServiceContainerMap.Get(IP)
	container := obj.(*types.Container)
	containerLogs, err := utils.GetContainerLogs(container.ID, tail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, R.Error(http.StatusInternalServerError, "Failed to retrieve logs, try again later!", nil))
		return
	}
	c.JSON(http.StatusOK, R.OK(containerLogs))
}
