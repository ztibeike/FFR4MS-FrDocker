package router

import (
	"frdocker/web/controller/command"
	"frdocker/web/controller/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	command.RegisterRouter(r)
	user.RegisterRouter(r)
	return r
}
