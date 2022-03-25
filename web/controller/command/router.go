package command

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	r.POST("/command/up", UpContainerController)
}
