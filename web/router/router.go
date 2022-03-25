package router

import (
	"frdocker/web/controller/command"
	"frdocker/web/controller/user"
	"frdocker/web/filter"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(filter.UserAuthFilter())
	command.RegisterRouter(r)
	user.RegisterRouter(r)
	return r
}
