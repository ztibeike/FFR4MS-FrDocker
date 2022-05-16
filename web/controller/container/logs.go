package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/settings"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils"
	"gitee.com/zengtao321/frdocker/web/entity/R"

	"github.com/gin-gonic/gin"
)

func GetContainerLogs(c *gin.Context) {
	IP := c.Query("ip")
	if IP == "" || !commons.IPServiceContainerMap.Has(IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "No Such IP!", nil))
		return
	}
	tail := c.Query("tail")
	if tail == "" {
		tail = "100"
	}
	obj, _ := commons.IPServiceContainerMap.Get(IP)
	container := obj.(*types.Container)
	containerLogs, err := utils.GetContainerLogs(container.ID, tail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, R.Error(http.StatusInternalServerError, "Failed to retrieve logs, try again later!", nil))
		return
	}
	c.JSON(http.StatusOK, R.OK(containerLogs))
}

type MonitorLog struct {
	Level string `json:"level"`
	Time  string `json:"time"`
	Msg   string `json:"msg"`
}

func GetMonitorLogs(c *gin.Context) {
	IP := c.Query("ip")
	if IP == "" || !commons.IPServiceContainerMap.Has(IP) {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "No Such IP!", nil))
		return
	}
	tail := c.Query("tail")
	if tail == "" {
		tail = "100"
	}
	fileName := fmt.Sprintf("%s/%s-%s.log", settings.LOG_FILE_DIR, commons.Network, IP)
	var logs string
	if utils.PathExists(fileName) {
		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("tail -n %s %s", tail, fileName))
		out, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err == nil {
			bytes, _ := ioutil.ReadAll(out)
			tempLogs := strings.Split(string(bytes), "\n")
			for _, tempLog := range tempLogs {
				var monitorLog MonitorLog
				_ = json.Unmarshal([]byte(tempLog), &monitorLog)
				logs += fmt.Sprintf("[%s] [%s] %s\n", strings.ToUpper(monitorLog.Level), monitorLog.Time, monitorLog.Msg)
			}
		}
	}
	c.JSON(http.StatusOK, R.OK(logs))
}
