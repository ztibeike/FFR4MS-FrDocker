package command

import (
	"net/http"

	"gitee.com/zengtao321/frdocker/constants"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/entity/dto"

	"github.com/gin-gonic/gin"
)

func UpContainerController(c *gin.Context) {
	var upContainerDTO dto.UpContainerDTO
	if err := c.ShouldBindJSON(&upContainerDTO); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	container := &types.Container{
		IP:     upContainerDTO.ServiceIP,
		Port:   upContainerDTO.ServicePort,
		Group:  upContainerDTO.ServiceGroup,
		Health: true,
	}
	obj, _ := constants.ServiceGroupMap.Get(container.Group)
	serviceGroup := obj.(*types.ServiceGroup)
	serviceGroup.Services = append(serviceGroup.Services, container.IP)
	container.Gateway = serviceGroup.Gateway
	constants.IPAllMSMap.Set(container.IP, "SERVICE:"+container.Group)
	utils.GetServiceContainers([]*types.Container{container})
	constants.IPServiceContainerMap.Set(container.IP, container)
	c.JSON(http.StatusOK, R.OK(nil))
}
