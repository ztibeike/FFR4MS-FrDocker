package router

import (
	"frdocker/settings"
	"frdocker/web/controller/command"
	"frdocker/web/controller/container"
	"frdocker/web/controller/user"
	"frdocker/web/filter"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(settings.RUNNING_MODE)
	r := gin.Default()
	r.Use(filter.UserAuthFilter())
	command.RegisterRouter(r)
	user.RegisterRouter(r)
	container.RegisterRouter(r)
	return r
}
