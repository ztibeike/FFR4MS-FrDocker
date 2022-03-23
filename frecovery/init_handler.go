package frecovery

import (
	"context"
	"errors"
	"frdocker/constants"
	"frdocker/models"
	"frdocker/types"
	"frdocker/utils"
	"log"
	"strings"

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
		InitFromDataBase(network.Id)
	}
}

func InitFromDataBase(networkId string) {
	filter := bson.D{
		{Key: "networkId", Value: networkId},
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
	var containers []*types.Container
	if strings.HasPrefix(confPath, "http") {
		containers = utils.GetConfigFromEureka(confPath)
	} else {
		log.Fatalln(errors.New("do not support file-config yet"))
	}
	utils.GetServiceContainers(containers)
}
