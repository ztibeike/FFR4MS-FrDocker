package recovery

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	r.GET("/recovery/list", GetRecoveryList)
}
