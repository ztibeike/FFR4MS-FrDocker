package container

import (
	"context"
	"frdocker/constants"
	"frdocker/db"
	"frdocker/models"
	"frdocker/types"
	"frdocker/web/entity/R"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetContainer(c *gin.Context) {
	IP := c.Query("ip")
	var containers []*types.Container
	if IP == "" {
		mp := constants.IPServiceContainerMap.Items()
		var containers []*types.Container
		for _, v := range mp {
			containers = append(containers, v.(*types.Container))
		}
		c.JSON(http.StatusOK, R.OK(containers))
		return
	}
	if !constants.IPServiceContainerMap.Has(IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	obj, _ := constants.IPServiceContainerMap.Get(IP)
	containers = append(containers, obj.(*types.Container))
	c.JSON(http.StatusOK, R.OK(containers))
}

func GetContainerCallChain(c *gin.Context) {
	serviceGroupMap := constants.ServiceGroupMap.Items()
	var callChainMap = make(map[string][]string)
	for group, obj := range serviceGroupMap {
		ms := obj.(*types.ServiceGroup)
		var calls []string
		if ms.Leaf {
			continue
		}
		for _, IP := range ms.Services {
			obj, _ := constants.IPServiceContainerMap.Get(IP)
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

func GetContainerTraffic(c *gin.Context) {
	IP := c.Query("ip")
	trafficMgo := db.GetTrafficMgo()
	var traffics []*models.Traffic
	if IP == "" {
		var filter = bson.D{
			{Key: "network", Value: constants.Network},
			{Key: "entry", Value: true},
		}
		var containerTraffics []*models.ContainerTraffic
		trafficMgo.FindMany(filter).All(context.TODO(), &containerTraffics)
		if len(containerTraffics) == 0 {
			c.JSON(http.StatusOK, R.OK(traffics))
			return
		}
		sort.Sort(models.ContainerTrafficArray(containerTraffics))
		maxLen := len(containerTraffics[0].Traffic)
		traffics = make([]*models.Traffic, maxLen)
		for k := 1; k <= maxLen; k++ {
			traffic := &models.Traffic{
				Year:   containerTraffics[0].Traffic[maxLen-k].Year,
				Month:  containerTraffics[0].Traffic[maxLen-k].Month,
				Day:    containerTraffics[0].Traffic[maxLen-k].Day,
				Hour:   containerTraffics[0].Traffic[maxLen-k].Hour,
				Minute: containerTraffics[0].Traffic[maxLen-k].Minute,
				Number: 0,
			}
			for _, containerTraffic := range containerTraffics {
				l := len(containerTraffic.Traffic)
				if l < k {
					continue
				}
				traffic.Number += containerTraffic.Traffic[l-k].Number
			}
			traffics[maxLen-k] = traffic
		}
		c.JSON(http.StatusOK, R.OK(traffics))
		return
	}
	if !constants.IPServiceContainerMap.Has(IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	var filter = bson.D{
		{Key: "network", Value: constants.Network},
		{Key: "ip", Value: IP},
	}
	var containerTraffic models.ContainerTraffic
	trafficMgo.FindOne(filter).Decode(&containerTraffic)
	traffics = containerTraffic.Traffic
	c.JSON(http.StatusOK, R.OK(traffics))
}
