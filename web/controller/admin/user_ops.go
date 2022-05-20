package admin

import (
	"context"
	"net/http"

	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/web/entity"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"gitee.com/zengtao321/frdocker/web/service/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var userMgo = db.GetUserMgo()

// AddUser 管理员添加用户接口
// @Summary 管理员添加用户接口
// @Description 管理员添加新用户
// @Tags 管理员操作
// @Accept application/json
// @Produce application/json
// @Param Authorization	header	string					true	"token"
// @Param user			body	entity.UserEntity		true	"用户信息"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=entity.UserEntity} "返回新增用户"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /admin/user/add [post]
func AddUser(c *gin.Context) {
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
	if user.Role == "" {
		user.Role = "USER"
	}
	if user.Role != "USER" && user.Role != "ADMIN" {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "Bad user role!", nil))
		return
	}
	cryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(cryptedPassword)
	user.Id = uuid.New().String()
	userMgo.InsertOne(user)
	c.JSON(http.StatusOK, R.OK(user))
}

// GetUserList 获取用户列表
// @Summary 管理员获取用户列表
// @Description 管理员获取用户列表
// @Tags 管理员操作
// @Produce application/json
// @Param Authorization	header	string	true	"token"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity "返回用户列表"
// @Router /admin/user/list [get]
func GetUserList(c *gin.Context) {
	var users []entity.UserEntity
	userMgo.FindAll().All(context.TODO(), &users)
	c.JSON(http.StatusOK, R.OK(users))
}

// DeleteUser 管理员删除用户
// @Summary 管理员删除用户
// @Description 管理删除用户
// @Tags 管理员操作
// @Accept application/json
// @Produce application/json
// @Param Authorization	header	string					true	"token"
// @Param users			body	[]entity.UserEntity		true	"用户信息"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=int} "返回删除用户数量"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /admin/user/delete [post]
func DeleteUser(c *gin.Context) {
	var users []entity.UserEntity
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	tokenStr := c.Request.Header["Authorization"][0]
	claims, _ := token.ParseToken(tokenStr)
	currentUserId := claims.UserId
	var matchCount = 0
	for _, user := range users {
		if currentUserId == user.Id {
			continue
		}
		var filter = bson.D{{Key: "id", Value: user.Id}}
		matchCount += int(userMgo.Delete(filter))
	}
	c.JSON(http.StatusOK, R.OK(matchCount))
}

// UpdateUser 管理员更新用户信息
// @Summary 管理员更新用户信息
// @Description 管理员更新用户信息
// @Tags 管理员操作
// @Accept application/json
// @Produce application/json
// @Param Authorization	header	string					true	"token"
// @Param user			body	entity.UserEntity		true	"用户信息"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=entity.UserEntity} "返回更新后的用户信息"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /admin/user/update [post]
func UpdateUser(c *gin.Context) {
	var user entity.UserEntity
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	var filter = bson.D{{Key: "id", Value: user.Id}}
	var tempUser *entity.UserEntity
	userMgo.FindOne(filter).Decode(&tempUser)
	if tempUser == nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "No such user!", nil))
		return
	}
	if tempUser.Password != user.Password {
		cryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(cryptedPassword)
	}
	userMgo.ReplaceOne(filter, user)
	c.JSON(http.StatusOK, R.OK(user))
}
