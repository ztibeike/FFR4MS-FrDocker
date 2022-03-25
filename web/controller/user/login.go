package user

import (
	"frdocker/db"
	"frdocker/web/entity"
	"frdocker/web/entity/R"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

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
	c.JSON(http.StatusOK, R.OK(nil))
}
