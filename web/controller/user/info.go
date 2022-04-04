package user

import (
	"net/http"

	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/service/token"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	tokenStr := c.Request.Header["Authorization"]
	claims, _ := token.ParseToken(tokenStr[0])
	roles := []string{claims.UserRole}
	c.JSON(http.StatusOK, R.OK(roles))
}
