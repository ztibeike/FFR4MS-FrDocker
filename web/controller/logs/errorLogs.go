package logs

import (
	"context"
	"net/http"
	"strconv"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/web/entity/R"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var errorLogMgo = db.GetErrorLogMgo()

// GetErrorLogs 获取异常日志
// @Summary 获取异常日志
// @Description 获取异常日志
// @Tags 日志操作
// @Produce application/json
// @Param Authorization	header	string	true	"token"
// @Param viewed		query	bool		false	"是否已读"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=string} "返回异常日志"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /logs/errors [get]
func GetErrorLogs(c *gin.Context) {
	viewed := c.Query("viewed")
	if viewed == "" {
		viewed = "false"
	}
	val, err := strconv.ParseBool(viewed)
	if err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	filter := bson.D{{"network", commons.Network}, {"viewed", val}}
	var errLogs []types.ErrorLog
	_ = errorLogMgo.FindMany(filter).All(context.TODO(), &errLogs)
	c.JSON(http.StatusOK, R.OK(errLogs))
}
