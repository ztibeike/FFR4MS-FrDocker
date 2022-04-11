package admin

import (
	"gitee.com/zengtao321/frdocker/web/filter"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	r.POST("/admin/user/add", filter.AdminAuthFilter(), AddUser)
	r.POST("/admin/user/list", filter.AdminAuthFilter(), GetUserList)
	r.POST("/admin/user/delete", filter.AdminAuthFilter(), DeleteUser)
	r.POST("/admin/user/update", filter.AdminAuthFilter(), UpdateUser)
}
