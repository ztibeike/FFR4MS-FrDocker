package container

import (
	"net/http"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/web/entity/R"

	"github.com/gin-gonic/gin"
)

// GetContainer 根据IP获取微服务容器信息
// @Summary 根据IP获取微服务容器信息
// @Description 根据IP获取微服务容器信息
// @Tags 微服务容器操作
// @Produce application/json
// @Param Authorization	header	string	true	"token"
// @Param ip			query	string	false	"容器的IP地址"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=[]types.Container} "返回对应IP地址的容器/所有容器"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /container [get]
func GetContainer(c *gin.Context) {
	IP := c.Query("ip")
	var containers []*types.Container
	if IP == "" {
		mp := commons.IPServiceContainerMap.Items()
		var containers []*types.Container
		for _, v := range mp {
			_container, _ := v.(*types.Container)
			container := &types.Container{
				IP:      _container.IP,
				Port:    _container.Port,
				Group:   _container.Group,
				Gateway: _container.Gateway,
				Leaf:    _container.Leaf,
				Health:  _container.Health,
				ID:      _container.ID,
				Name:    _container.Name,
				Calls:   _container.Calls,
				Entry:   _container.Entry,
			}
			containers = append(containers, container)
		}
		c.JSON(http.StatusOK, R.OK(containers))
		return
	}
	if !commons.IPServiceContainerMap.Has(IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	obj, _ := commons.IPServiceContainerMap.Get(IP)
	containers = append(containers, obj.(*types.Container))
	c.JSON(http.StatusOK, R.OK(containers))
}

// GetContainerCallChain 获取微服务系统的服务调用链
// @Summary 获取微服务系统的服务调用链
// @Description 获取微服务系统的服务调用链
// @Tags 微服务容器操作
// @Produce application/json
// @Param Authorization	header	string	true	"token"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=map[string][]string} "返回微服务系统的服务调用链"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /container/calls [get]
func GetContainerCallChain(c *gin.Context) {
	serviceGroupMap := commons.ServiceGroupMap.Items()
	var callChainMap = make(map[string][]string)
	for group, obj := range serviceGroupMap {
		ms := obj.(*types.ServiceGroup)
		var calls []string
		if ms.Leaf {
			continue
		}
		for _, IP := range ms.Services {
			obj, _ := commons.IPServiceContainerMap.Get(IP)
			container := obj.(*types.Container)
			if len(container.Calls) != 0 {
				calls = append(calls, container.Calls...)
				break
			}
		}
		callChainMap[group] = calls
	}
	c.JSON(http.StatusOK, R.OK(callChainMap))
}
