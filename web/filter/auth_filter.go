package filter

import (
	"net/http"
	"strings"

	"gitee.com/zengtao321/frdocker/utils/logger"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/service/token"

	"github.com/gin-gonic/gin"
)

func checkNoAuth(url string) bool {
	var exceptURLs = []string{
		"/user/login",
		"/user/logout",
		"/swagger",
	}
	for _, _url := range exceptURLs {
		if strings.HasPrefix(url, _url) {
			return true
		}
	}
	return false
}

func UserAuthFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.URL.String()
		if checkNoAuth(uri) {
			c.Next()
			return
		}
		tokenStr, ok := c.Request.Header["Authorization"]
		if !ok {
			c.JSON(http.StatusUnauthorized, R.Error(http.StatusUnauthorized, "Not logged in!", nil))
			c.Abort()
			return
		}
		claims, err := token.ParseToken(tokenStr[0])
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				newToken, _ := token.RefreshToken(claims)
				if newToken != "" {
					logger.Info(nil, "Generate new token: %s\n", newToken)
					c.Request.Header.Set("Authorization", newToken)
					c.Header("refresh-token", newToken)
					c.Next()
					return
				}
			}
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
