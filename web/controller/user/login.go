package user

import (
	"frdocker/db"
	"frdocker/web/entity"
	"frdocker/web/entity/R"
	"frdocker/web/service/token"
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
	token, err := token.GenerateToken(expectUser.Id, expectUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, R.Error(http.StatusInternalServerError, "Login failed caused by service error, try again later!", nil))
		return
	}
	c.JSON(http.StatusOK, R.OK(gin.H{"token": token}))
}
