package perf

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	r.GET("/perf/system", GetSystemPerformance)
	r.GET("/perf/container", GetContainerPerformance)
}
