package user

import (
	"net/http"

	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/web/entity"
	"gitee.com/zengtao321/frdocker/web/entity/R"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// RegisterController 用户注册
// @Summary 用户注册
// @Description 用户注册
// @Tags 用户操作
// @Accept application/json
// @Produce application/json
// @Param user			body	entity.UserEntity	true	"注册用户信息"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity "返回成功信息"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /user/register [post]
func RegisterController(c *gin.Context) {
	userMgo := db.GetUserMgo()
	var user entity.UserEntity
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	var filter = bson.D{{Key: "username", Value: user.Username}}
	var tempUser *entity.UserEntity
	userMgo.FindOne(filter).Decode(&tempUser)
	if tempUser != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Username already exists!", nil))
		return
	}
	user.Role = "USER"
	cryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(cryptedPassword)
	user.Id = uuid.New().String()
	userMgo.InsertOne(user)
	c.JSON(http.StatusOK, R.OK(nil))
}
