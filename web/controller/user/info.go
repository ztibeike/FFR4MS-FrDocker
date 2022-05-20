package user

import (
	"net/http"

	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/web/entity"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/entity/dto"
	"gitee.com/zengtao321/frdocker/web/service/token"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags 用户操作
// @Produce application/json
// @Param Authorization	header	string	true	"token"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=dto.UserInfoDTO} "返回用户信息"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /user/info [get]
func GetUserInfo(c *gin.Context) {
	tokenStr := c.Request.Header["Authorization"]
	claims, _ := token.ParseToken(tokenStr[0])
	userId := claims.UserId
	var filter = bson.D{{Key: "id", Value: userId}}
	userMgo := db.GetUserMgo()
	var user entity.UserEntity
	userMgo.FindOne(filter).Decode(&user)
	userInfoDTO := dto.UserInfoDTO{
		Username:    user.Username,
		Permissions: []string{claims.UserRole},
		Avatar:      "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif?imageView2/1/w/80/h/80",
	}
	c.JSON(http.StatusOK, R.OK(userInfoDTO))
}
