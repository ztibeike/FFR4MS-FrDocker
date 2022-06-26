package logs

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	r.GET("/logs/container", GetContainerLogs)
	r.GET("/logs/monitor", GetMonitorLogs)
	r.GET("/logs/error", GetErrorLogs)
}
