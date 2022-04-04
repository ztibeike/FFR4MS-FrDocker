package frecovery

import (
	"context"
	"strings"

	"gitee.com/zengtao321/frdocker/constants"
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
	logger.Info("Init Containers from DataBase\n")
	filter := bson.D{
		{Key: "network", Value: ifaceName},
	}
	cursor := containerMgo.FindMany(filter)
	defer cursor.Close(context.TODO())
	var dbContainers []*models.Container
	cursor.All(context.TODO(), &dbContainers)
	for _, dbContainer := range dbContainers {
		container := dbContainer.Container
		constants.IPServiceContainerMap.Set(container.IP, container)
		constants.IPAllMSMap.Set(container.IP, "SERVICE:"+container.Group)
		colon := strings.Index(container.Gateway, ":")
		constants.IPAllMSMap.Set(container.Gateway[:colon], "GATEWAY:"+container.Group)
		if !constants.ServiceGroupMap.Has(container.Group) {
			serviceGroup := &types.ServiceGroup{
				Gateway:  container.Gateway,
				Services: []string{container.IP},
				Leaf:     container.Leaf,
				Entry:    container.Entry,
			}
			constants.ServiceGroupMap.Set(container.Group, serviceGroup)
		} else {
			obj, _ := constants.ServiceGroupMap.Get(container.Group)
			serviceGroup := obj.(*types.ServiceGroup)
			serviceGroup.Services = append(serviceGroup.Services, container.IP)
		}
	}
}

func InitFromConfiguration(confPath string) {
	logger.Info("Init Containers from Registry Configuration\n")
	var containers []*types.Container
	if strings.HasPrefix(confPath, "http") {
		containers = utils.GetConfigFromEureka(confPath)
	} else {
		logger.Fatalln("Do not Support File-config Yet!")
	}
	utils.GetServiceContainers(containers)
}
