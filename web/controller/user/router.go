package user

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	r.POST("/user/register", RegisterController)
	r.POST("/user/login", LoginController)
}
