package user

import (
	"net/http"

	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/web/entity"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/service/token"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// LoginController 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户操作
// @Accept application/json
// @Produce application/json
// @Param user	body	entity.UserEntity	true	"登录用户信息"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=gin.H} "返回登录Token"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /user/login [post]
func LoginController(c *gin.Context) {
	userMgo := db.GetUserMgo()
	var actualUser, expectUser entity.UserEntity
	if err := c.ShouldBind(&actualUser); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	var filter = bson.D{{Key: "username", Value: actualUser.Username}}
	userMgo.FindOne(filter).Decode(&expectUser)
	if expectUser.Id == "" {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Username or password incorrect!", nil))
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(expectUser.Password), []byte(actualUser.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Username or password incorrect!", nil))
		return
	}
	token, err := token.GenerateToken(expectUser.Id, expectUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, R.Error(http.StatusInternalServerError, "Login failed caused by service error, try again later!", nil))
		return
	}
	c.JSON(http.StatusOK, R.OK(gin.H{"token": token}))
}

// LogoutController 退出登录
// @Summary 退出登录
// @Description 退出登录
// @Tags 用户操作
// @Produce application/json
// @Param Authorization	header	string	true	"token"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=gin.H} "返回登录Token"
// @Router /user/logout [post]
func LogoutController(c *gin.Context) {
	c.JSON(http.StatusOK, R.OK(nil))
}
