package frecovery

import (
	"frdocker/constants"
	"frdocker/models"
	"frdocker/types"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func SetupCloseHandler(ifaceName string) {
	sigalChan := make(chan os.Signal, 1)
	signal.Notify(sigalChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigalChan
	var wg sync.WaitGroup
	wg.Add(2)
	go CloseChannels(ifaceName, &wg)
	go SaveContainerInfo(ifaceName, &wg)
	wg.Wait()
}

func CloseChannels(ifaceName string, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Printf("\n[%s] Closing All Channels......\n", time.Now().Format("2006-01-02 15:04:05"))
	for IP, ch := range constants.IPChanMap {
		close(ch)
		delete(constants.IPChanMap, IP)
	}
	pcapHandler.Close()
	logger.Printf("[%s] Stop capturing packets on interface: %s\n", time.Now().Format("2006-01-02 15:04:05"), ifaceName)
}

func SaveContainerInfo(ifaceName string, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Printf("\n[%s] Saving All Containers Info & States......\n", time.Now().Format("2006-01-02 15:04:05"))
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
}
