package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils"
	"gitee.com/zengtao321/frdocker/utils/logger"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/entity/dto"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

// AddContainerController 添加微服务实例
// @Summary 添加微服务实例
// @Description 添加微服务实例
// @Tags 微服务实例操作
// @Accept application/json
// @Produce application/json
// @Param Authorization		header	string					true	"token"
// @Param addContainerDTO	body	dto.AddContainerDTO		true	"微服务实例信息"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=types.Container} "返回新增的微服务实例信息"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /command/add [post]
func AddContainerController(c *gin.Context) {
	var addContainerDTO dto.AddContainerDTO
	if err := c.ShouldBindJSON(&addContainerDTO); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	container := &types.Container{
		IP:     addContainerDTO.ServiceIP,
		Port:   addContainerDTO.ServicePort,
		Group:  addContainerDTO.ServiceGroup,
		Health: true,
	}
	if commons.IPServiceContainerMap.Has(container.IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Service Already Exists!", nil))
		return
	}
	utils.GetServiceContainers([]*types.Container{container})
	if container.ID == "" {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Incorrect Service Info!", nil))
		return
	}
	obj, _ := commons.ServiceGroupMap.Get(container.Group)
	serviceGroup := obj.(*types.ServiceGroup)
	container.Leaf = serviceGroup.Leaf
	container.Entry = serviceGroup.Entry
	obj, _ = commons.IPServiceContainerMap.Get(serviceGroup.Services[0])
	otherContainer := obj.(*types.Container)
	container.Calls = make([]string, len(otherContainer.Calls))
	copy(container.Calls, otherContainer.Calls)
	serviceGroup.Services = append(serviceGroup.Services, container.IP)
	container.Gateway = serviceGroup.Gateway
	commons.IPAllMSMap.Set(container.IP, "SERVICE:"+container.Group)
	commons.IPServiceContainerMap.Set(container.IP, container)
	commons.AddContainerChan <- container.IP
	logger.Info(container.IP, "[Add New Container] [Group(%s) IP(%s) Port(%s) ID(%s)]\n", container.Group, container.IP, container.Port, container.ID)
	c.JSON(http.StatusOK, R.OK(container))
}

// DeleteContainerController 下线微服务实例
// @Summary 下线微服务实例
// @Description 下线微服务实例
// @Tags 微服务实例操作
// @Accept application/json
// @Produce application/json
// @Param Authorization		header	string			true	"token"
// @Param deleteContainer	body	types.Container	true	"下线微服务实例信息"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity "返回成功信息"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /command/delete [post]
func DeleteContainerController(c *gin.Context) {
	var deleteContainer types.Container
	if err := c.ShouldBindJSON(&deleteContainer); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	if !commons.IPServiceContainerMap.Has(deleteContainer.IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	obj, _ := commons.IPServiceContainerMap.Get(deleteContainer.IP)
	container := obj.(*types.Container)
	ConstantsDelete(container)
	commons.DeleteContainerChan <- deleteContainer.IP
	DataBaseDelete(deleteContainer.IP)
	logger.Warn(container.IP, "[Delete Container] [Group(%s) IP(%s) Port(%s) ID(%s)]\n", container.Group, container.IP, container.Port, container.ID)
	c.JSON(http.StatusOK, R.OK(nil))
}

// DeleteBatchContainerController 批量下线微服务实例
// @Summary 批量下线微服务实例
// @Description 批量下线微服务实例
// @Tags 微服务实例操作
// @Accept application/json
// @Produce application/json
// @Param Authorization		header	string				true	"token"
// @Param deleteContainers	body	[]types.Container	true	"下线微服务实例信息列表"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=int} "返回成功下线的微服务实例数量"
// @Failure 400 {object} R.ResponseEntity 			"返回失败信息"
// @Router /command/delete/batch [post]
func DeleteBatchContainerController(c *gin.Context) {
	var deleteContainers []types.Container
	if err := c.ShouldBindJSON(&deleteContainers); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	var matchCount = 0
	var wg sync.WaitGroup
	for _, deleteContainer := range deleteContainers {
		if !commons.IPServiceContainerMap.Has(deleteContainer.IP) {
			continue
		}
		matchCount += 1
		obj, _ := commons.IPServiceContainerMap.Get(deleteContainer.IP)
		container := obj.(*types.Container)
		wg.Add(1)
		go func() {
			ConstantsDelete(container)
			wg.Done()
		}()
	}
	wg.Wait()
	commons.DeleteContainerChan <- "1"
	for _, deleteContainer := range deleteContainers {
		wg.Add(1)
		container := deleteContainer
		go func() {
			DataBaseDelete(container.IP)
			wg.Done()
		}()
	}
	wg.Wait()
	c.JSON(http.StatusOK, R.OK(matchCount))
}

func ConstantsDelete(deleteContainer *types.Container) {
	// 等待容器正在处理的请求结束
	deleteContainer.Health = false
	time.Sleep(500 * time.Millisecond)
	commons.IPServiceContainerMap.Pop(deleteContainer.IP)
	commons.IPChanMapMutex.Lock()
	ch, ok := commons.IPChanMap[deleteContainer.IP]
	if ok {
		close(ch)
	}
	delete(commons.IPChanMap, deleteContainer.IP)
	commons.IPChanMapMutex.Unlock()
	commons.IPAllMSMap.Pop(deleteContainer.IP)
	obj, _ := commons.ServiceGroupMap.Get(deleteContainer.Group)
	serviceGroup := obj.(*types.ServiceGroup)
	pos := -1
	for idx, service := range serviceGroup.Services {
		if service == deleteContainer.IP {
			pos = idx
			break
		}
	}
	if pos != -1 {
		serviceGroup.Services = append(serviceGroup.Services[:pos], serviceGroup.Services[pos+1:]...)
	}
}

func DataBaseDelete(IP string) {
	containerMgo := db.GetContainerMgo()
	trafficMgo := db.GetTrafficMgo()
	var filter = bson.D{
		{Key: "network", Value: commons.Network},
		{Key: "ip", Value: IP},
	}
	trafficMgo.Delete(filter)
	filter = bson.D{
		{Key: "network", Value: commons.Network},
		{Key: "container.ip", Value: IP},
	}
	containerMgo.Delete(filter)
}

// UpContainerController 上线微服务实例
// @Summary 上线微服务实例
// @Description 上线微服务实例
// @Tags 微服务实例操作
// @Accept application/json
// @Produce application/json
// @Param Authorization	header	string			true	"token"
// @Param upContainer	body	types.Container	true	"上线的微服务实例信息"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity "返回成功信息"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /command/up [post]
func UpContainerController(c *gin.Context) {
	var upContainer types.Container
	if err := c.ShouldBindJSON(&upContainer); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	if !commons.IPServiceContainerMap.Has(upContainer.IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	obj, _ := commons.IPServiceContainerMap.Get(upContainer.IP)
	container := obj.(*types.Container)
	container.Health = true
	container.States = nil
	logger.Info(container.IP, "[Mark Container Health] [Group(%s) IP(%s) ID(%s)]\n", container.Group, container.IP, container.ID)
	GatewayUpServiceInstance(container)
	c.JSON(http.StatusOK, R.OK(nil))
}

func GatewayUpServiceInstance(container *types.Container) {
	var gateway = container.Gateway
	var url = fmt.Sprintf("http://%s/zuulApi/upServiceInstance", gateway)
	var upServiceInfo = &dto.GateWayUpService{
		ServiceName:    container.Group,
		UpInstanceHost: container.IP,
		UpInstancePort: container.Port,
	}
	var requestBody, _ = json.Marshal(upServiceInfo)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil || response == nil || response.StatusCode != 200 {
		logger.Warn(container.IP, "[Gateway Up Instance] [Gateway(%s) Group(%s) Instance(%s)] Gateway Error!\n",
			gateway, container.Group, container.IP)
	}
	logger.Info(container.IP, "[Gateway Up Instance] [Gateway(%s) Group(%s) Instance(%s)] Service Up!\n",
		gateway, container.Group, container.IP)
}
