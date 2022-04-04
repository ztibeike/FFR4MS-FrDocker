package user

import (
	"frdocker/web/entity/R"
	"frdocker/web/service/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	tokenStr := c.Request.Header["Authorization"]
	claims, _ := token.ParseToken(tokenStr[0])
	roles := []string{claims.UserRole}
	c.JSON(http.StatusOK, R.OK(roles))
}
