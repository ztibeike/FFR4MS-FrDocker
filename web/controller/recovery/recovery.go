package recovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/entity/dto"
	"github.com/gin-gonic/gin"
)

func GetRecoveryList(c *gin.Context) {
	group := c.Query("group")
	if group != "" && !commons.ServiceGroupMap.Has(group) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	serviceGroupMap := commons.ServiceGroupMap.Items()
	var data []dto.RecoveryMessageDTO
	for k, v := range serviceGroupMap {
		if !strings.HasPrefix(k, group) {
			continue
		}
		serviceGroup := v.(*types.ServiceGroup)
		gateway := serviceGroup.Gateway
		var url = fmt.Sprintf("http://%s/zuulApi/getReplayedTraceId", gateway)
		response, err := http.Get(url)
		if err != nil || response == nil || response.StatusCode != http.StatusOK {
			continue
		}
		body, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()
		var result = struct {
			Code    int
			Message string
			Data    []dto.RecoveryMessageDTO
		}{}
		_ = json.Unmarshal(body, &result)
		data = append(data, result.Data...)
	}
	c.JSON(http.StatusOK, R.OK(data))
}
