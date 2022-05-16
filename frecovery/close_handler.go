package frecovery

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/models"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils/logger"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func SetupCloseHandler(ifaceName string, wg *sync.WaitGroup) {
	defer wg.Done()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGABRT, syscall.SIGQUIT)
	<-signalChan
	ClosePcapHandler(ifaceName)
	SaveContainerInfo(ifaceName)
	logger.Close()
}

func ClosePcapHandler(ifaceName string) {
	pcapHandler.Close()
	logger.Info(nil, "Stop capturing packets on interface: %s\n", ifaceName)
}

func SaveContainerInfo(ifaceName string) {
	logger.Info(nil, "Saving All Containers Info & States......\n")
	var network *models.NetWork
	filter := bson.D{
		{Key: "name", Value: ifaceName},
	}
	networkMgo.FindOne(filter).Decode(&network)
	if network == nil {
		network = &models.NetWork{}
		network.Name = ifaceName
		id := uuid.New()
		network.Id = id.String()
		networkMgo.InsertOne(network)
		var dbContainers []interface{}
		for _, obj := range commons.IPServiceContainerMap.Items() {
			container := obj.(*types.Container)
			dbContainer := &models.Container{
				Container: container,
				Network:   ifaceName,
			}
			dbContainers = append(dbContainers, dbContainer)
		}
		containerMgo.InsertMany(dbContainers)
	} else {
		for _, obj := range commons.IPServiceContainerMap.Items() {
			container := obj.(*types.Container)
			dbContainer := &models.Container{
				Container: container,
				Network:   ifaceName,
			}
			filter := bson.D{
				{Key: "network", Value: ifaceName},
				{Key: "container.ip", Value: dbContainer.Container.IP},
			}
			_ = containerMgo.ReplaceOne(filter, dbContainer)
		}
	}
}

// func CronSaveContainerInfo
