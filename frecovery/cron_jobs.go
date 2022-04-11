package frecovery

import (
	"context"
	"sync/atomic"
	"time"

	"gitee.com/zengtao321/frdocker/constants"
	"gitee.com/zengtao321/frdocker/models"
	"gitee.com/zengtao321/frdocker/settings"
	"gitee.com/zengtao321/frdocker/types"

	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/bson"
)

func CronSaveTraffic(ctx context.Context, trafficChan chan string) {
	var spec string
	if settings.CRON_LEVEL == "HOUR" {
		spec = "0 59 * * * *"
	} else if settings.CRON_LEVEL == "MINUTE" {
		spec = "59 * * * * *"
	} else {
		return
	}
	trafficCountMap := make(map[string]*int64)
	c := cron.New()
	for _, obj := range constants.IPServiceContainerMap.Items() {
		container := obj.(*types.Container)
		var trafficCount int64 = 0
		trafficCountMap[container.IP] = &trafficCount
		c.AddFunc(spec, func() {
			t := time.Now()
			var containerTraffic *models.ContainerTraffic
			var filter = bson.D{
				{Key: "network", Value: constants.Network},
				{Key: "ip", Value: container.IP},
			}
			trafficMgo.FindOne(filter).Decode(&containerTraffic)
			traffic := &models.Traffic{
				Year:   t.Year(),
				Month:  int(t.Month()),
				Day:    t.Day(),
				Hour:   t.Hour(),
				Minute: t.Minute(),
				Level:  settings.CRON_LEVEL,
				Number: atomic.LoadInt64(trafficCountMap[container.IP]),
			}
			if containerTraffic == nil {
				containerTraffic = &models.ContainerTraffic{
					Network: constants.Network,
					IP:      container.IP,
					Port:    container.Port,
					Group:   container.Group,
					Entry:   container.Entry,
					Traffic: []*models.Traffic{traffic},
				}
				trafficMgo.InsertOne(containerTraffic)
			} else {
				_traffics := containerTraffic.Traffic
				start := 0
				if len(_traffics) >= settings.CRON_TRAFFIC_LEN {
					start = len(_traffics) - settings.CRON_TRAFFIC_LEN + 1
				}
				containerTraffic.Traffic = append(containerTraffic.Traffic[start:], traffic)
				containerTraffic.Entry = container.Entry
				trafficMgo.ReplaceOne(filter, containerTraffic)
			}
			atomic.StoreInt64(trafficCountMap[container.IP], 0)
		})
	}
	go c.Start()
	for {
		select {
		case IP := <-trafficChan:
			{
				atomic.AddInt64(trafficCountMap[IP], 1)
			}
		case <-ctx.Done():
			c.Stop()
			return
		}
	}
}

func CronSaveContainerInfo(ctx context.Context, ifaceName string) {
	var spec string
	if settings.CRON_LEVEL == "HOUR" {
		spec = "0 59 * * * *"
	} else if settings.CRON_LEVEL == "MINUTE" {
		spec = "59 * * * * *"
	} else {
		return
	}
	c := cron.New()
	c.AddFunc(spec, func() {
		SaveContainerInfo(ifaceName)
	})
	go c.Start()
	defer c.Stop()
	<-ctx.Done()
}
