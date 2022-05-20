package container

import (
	"context"
	"net/http"
	"sort"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/models"
	"gitee.com/zengtao321/frdocker/web/entity/R"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetContainerTraffic 获取微服务容器访问流量
// @Summary 获取微服务容器访问流量
// @Description 获取微服务容器访问流量
// @Tags 微服务容器操作
// @Produce application/json
// @Param Authorization	header	string	true	"token"
// @Param ip			query	string	true	"容器的IP地址"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=[]models.Traffic} "返回微服务容器访问流量"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /container/traffic [get]
func GetContainerTraffic(c *gin.Context) {
	IP := c.Query("ip")
	trafficMgo := db.GetTrafficMgo()
	var traffics []*models.Traffic
	if IP == "" {
		var filter = bson.D{
			{Key: "network", Value: commons.Network},
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
				Level:  containerTraffics[0].Traffic[maxLen-k].Level,
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
	if !commons.IPServiceContainerMap.Has(IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	var filter = bson.D{
		{Key: "network", Value: commons.Network},
		{Key: "ip", Value: IP},
	}
	var containerTraffic models.ContainerTraffic
	trafficMgo.FindOne(filter).Decode(&containerTraffic)
	traffics = containerTraffic.Traffic
	c.JSON(http.StatusOK, R.OK(traffics))
}
