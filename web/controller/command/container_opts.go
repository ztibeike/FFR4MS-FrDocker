package command

import (
	"net/http"
	"sync"
	"time"

	"gitee.com/zengtao321/frdocker/constants"
	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils"
	"gitee.com/zengtao321/frdocker/utils/logger"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/entity/dto"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

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
	if constants.IPServiceContainerMap.Has(container.IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Service Already Exists!", nil))
		return
	}
	utils.GetServiceContainers([]*types.Container{container})
	if container.ID == "" {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Incorrect Service Info!", nil))
		return
	}
	obj, _ := constants.ServiceGroupMap.Get(container.Group)
	serviceGroup := obj.(*types.ServiceGroup)
	container.Leaf = serviceGroup.Leaf
	container.Entry = serviceGroup.Entry
	obj, _ = constants.IPServiceContainerMap.Get(serviceGroup.Services[0])
	otherContainer := obj.(*types.Container)
	container.Calls = make([]string, len(otherContainer.Calls))
	copy(container.Calls, otherContainer.Calls)
	serviceGroup.Services = append(serviceGroup.Services, container.IP)
	container.Gateway = serviceGroup.Gateway
	constants.IPAllMSMap.Set(container.IP, "SERVICE:"+container.Group)
	constants.IPServiceContainerMap.Set(container.IP, container)
	constants.AddContainerChan <- container.IP
	logger.Info(container.IP, "[Add New Container] [Group(%s) IP(%s) Port(%s) ID(%s)]\n", container.Group, container.IP, container.Port, container.ID)
	c.JSON(http.StatusOK, R.OK(container))
}

func DeleteContainerController(c *gin.Context) {
	var deleteContainer types.Container
	if err := c.ShouldBindJSON(&deleteContainer); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	if !constants.IPServiceContainerMap.Has(deleteContainer.IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	obj, _ := constants.IPServiceContainerMap.Get(deleteContainer.IP)
	container := obj.(*types.Container)
	ConstantsDelete(container)
	constants.DeleteContainerChan <- deleteContainer.IP
	DataBaseDelete(deleteContainer.IP)
	logger.Warn(container.IP, "[Delete Container] [Group(%s) IP(%s) Port(%s) ID(%s)]\n", container.Group, container.IP, container.Port, container.ID)
	c.JSON(http.StatusOK, R.OK(nil))
}

func DeleteBatchContainerController(c *gin.Context) {
	var deleteContainers []types.Container
	if err := c.ShouldBindJSON(&deleteContainers); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	var matchCount = 0
	var wg sync.WaitGroup
	for _, deleteContainer := range deleteContainers {
		if !constants.IPServiceContainerMap.Has(deleteContainer.IP) {
			continue
		}
		matchCount += 1
		obj, _ := constants.IPServiceContainerMap.Get(deleteContainer.IP)
		container := obj.(*types.Container)
		wg.Add(1)
		go func() {
			ConstantsDelete(container)
			wg.Done()
		}()
	}
	wg.Wait()
	constants.DeleteContainerChan <- "1"
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
	constants.IPServiceContainerMap.Pop(deleteContainer.IP)
	constants.IPChanMapMutex.Lock()
	ch, ok := constants.IPChanMap[deleteContainer.IP]
	if ok {
		close(ch)
	}
	delete(constants.IPChanMap, deleteContainer.IP)
	constants.IPChanMapMutex.Unlock()
	constants.IPAllMSMap.Pop(deleteContainer.IP)
	obj, _ := constants.ServiceGroupMap.Get(deleteContainer.Group)
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
		{Key: "network", Value: constants.Network},
		{Key: "ip", Value: IP},
	}
	trafficMgo.Delete(filter)
	filter = bson.D{
		{Key: "network", Value: constants.Network},
		{Key: "container.ip", Value: IP},
	}
	containerMgo.Delete(filter)
}

func UpContainerController(c *gin.Context) {
	var upContainer types.Container
	if err := c.ShouldBindJSON(&upContainer); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	if !constants.IPServiceContainerMap.Has(upContainer.IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	obj, _ := constants.IPServiceContainerMap.Get(upContainer.IP)
	container := obj.(*types.Container)
	container.Health = true
	container.States = nil
	logger.Info(container.IP, "[Mark Container Health] [Group(%s) IP(%s) ID(%s)]\n", container.Group, container.IP, container.ID)
	c.JSON(http.StatusOK, R.OK(nil))
}