package filter

import (
	"net/http"

	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/service/token"

	"github.com/gin-gonic/gin"
)

func UserAuthFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.URL.String()
		if uri == "/user/login" {
			c.Next()
			return
		}
		tokenStr, ok := c.Request.Header["Authorization"]
		if !ok {
			c.JSON(http.StatusUnauthorized, R.Error(http.StatusUnauthorized, "Not logged in!", nil))
			c.Abort()
			return
		}
		_, err := token.ParseToken(tokenStr[0])
		if err != nil {
			c.JSON(http.StatusUnauthorized, R.Error(http.StatusUnauthorized, err.Error(), nil))
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminAuthFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, ok := c.Request.Header["Authorization"]
		if !ok {
			c.JSON(http.StatusUnauthorized, R.Error(http.StatusUnauthorized, "Not logged in!", nil))
			c.Abort()
			return
		}
		claims, err := token.ParseToken(tokenStr[0])
		if err != nil {
			c.JSON(http.StatusUnauthorized, R.Error(http.StatusUnauthorized, "Not logged in!", nil))
			c.Abort()
			return
		}
		if claims.UserRole != "ADMIN" {
			c.JSON(http.StatusForbidden, R.Error(http.StatusForbidden, http.StatusText(http.StatusForbidden), nil))
			c.Abort()
			return
		}
		c.Next()
	}
}
