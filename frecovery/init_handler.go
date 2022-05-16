package frecovery

import (
	"context"
	"strings"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/models"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils"
	"gitee.com/zengtao321/frdocker/utils/logger"

	"go.mongodb.org/mongo-driver/bson"
)

func InitContainers(ifaceName, confPath string) {
	filter := bson.D{
		{Key: "name", Value: ifaceName},
	}
	var network *models.NetWork
	networkMgo.FindOne(filter).Decode(&network)
	if network == nil {
		InitFromConfiguration(confPath)
	} else {
		InitFromDataBase(ifaceName)
	}
}

func InitFromDataBase(ifaceName string) {
	logger.Info(nil, "Init Containers from DataBase\n")
	filter := bson.D{
		{Key: "network", Value: ifaceName},
	}
	cursor := containerMgo.FindMany(filter)
	defer cursor.Close(context.TODO())
	var dbContainers []*models.Container
	cursor.All(context.TODO(), &dbContainers)
	for _, dbContainer := range dbContainers {
		container := dbContainer.Container
		commons.IPServiceContainerMap.Set(container.IP, container)
		commons.IPAllMSMap.Set(container.IP, "SERVICE:"+container.Group)
		colon := strings.Index(container.Gateway, ":")
		commons.IPAllMSMap.Set(container.Gateway[:colon], "GATEWAY:"+container.Group)
		if !commons.ServiceGroupMap.Has(container.Group) {
			serviceGroup := &types.ServiceGroup{
				Gateway:  container.Gateway,
				Services: []string{container.IP},
				Leaf:     container.Leaf,
				Entry:    container.Entry,
			}
			commons.ServiceGroupMap.Set(container.Group, serviceGroup)
		} else {
			obj, _ := commons.ServiceGroupMap.Get(container.Group)
			serviceGroup := obj.(*types.ServiceGroup)
			serviceGroup.Services = append(serviceGroup.Services, container.IP)
		}
	}
}

func InitFromConfiguration(confPath string) {
	logger.Info(nil, "Init Containers from Registry Configuration\n")
	var containers []*types.Container
	if strings.HasPrefix(confPath, "http") {
		containers = utils.GetConfigFromEureka(confPath)
	} else {
		logger.Fatalln(nil, "Do not Support File-config Yet!")
	}
	utils.GetServiceContainers(containers)
}
