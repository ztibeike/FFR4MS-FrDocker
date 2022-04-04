package container

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	r.GET("/container", GetContainer)
	r.GET("/container/calls", GetContainerCallChain)
	r.GET("/container/traffic", GetContainerTraffic)
	r.GET("/container/logs", GetContainerLogs)
}
