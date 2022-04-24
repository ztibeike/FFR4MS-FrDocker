package command

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	r.POST("/command/add", AddContainerController)
	r.POST("/command/up", UpContainerController)
	r.POST("/command/delete", DeleteContainerController)
	r.POST("/command/delete/batch", DeleteBatchContainerController)
}
