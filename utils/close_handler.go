package utils

import (
	"frdocker/constants"
	"frdocker/models"
	"frdocker/types"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func SetupCloseHandler(ifaceName string) {
	sigalChan := make(chan os.Signal, 1)
	signal.Notify(sigalChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigalChan
	for IP, ch := range constants.IPChanMap {
		close(ch)
		delete(constants.IPChanMap, IP)
	}
	var network *models.NetWork
	filter := bson.D{
		{Key: "name", Value: ifaceName},
	}
	networkMgo.FindOne(filter).Decode(&network)
	if network == nil {
		network.Name = ifaceName
		id := uuid.New()
		network.Id = id.String()
		networkMgo.InsertOne(network)
		var dbContainers []interface{}
		for _, obj := range constants.IPServiceContainerMap.Items() {
			container := obj.(*types.Container)
			dbContainer := &models.Container{
				Container: container,
				NetworkId: network.Id,
			}
			dbContainers = append(dbContainers, dbContainer)
		}
		containerMgo.InsertMany(dbContainers)
	} else {
		for _, obj := range constants.IPServiceContainerMap.Items() {
			container := obj.(*types.Container)
			dbContainer := &models.Container{
				Container: container,
				NetworkId: network.Id,
			}
			filter := bson.D{
				{Key: "networkId", Value: dbContainer.NetworkId},
				{Key: "container.ip", Value: dbContainer.Container.IP},
			}
			_ = containerMgo.ReplaceOne(filter, dbContainer)
		}
	}

	os.Exit(1)
}
