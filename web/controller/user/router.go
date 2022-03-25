package user

import (
	"frdocker/web/filter"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	r.POST("/user/register", filter.AdminAuthFilter(), RegisterController)
	r.POST("/user/login", LoginController)
}
