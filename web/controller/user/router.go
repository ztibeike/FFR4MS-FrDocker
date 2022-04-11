package user

import (
	"gitee.com/zengtao321/frdocker/web/filter"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	r.POST("/user/register", filter.AdminAuthFilter(), RegisterController)
	r.POST("/user/login", LoginController)
	r.POST("/user/logout", LogoutController)
	r.GET("/user/info", GetUserInfo)
}
