package logs

import (
	"context"
	"net/http"

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
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=string} "返回异常日志"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /logs/errors [get]
func GetErrorLogs(c *gin.Context) {
	var errLogs []types.ErrorLog
	_ = errorLogMgo.FindAll().All(context.TODO(), &errLogs)
	c.JSON(http.StatusOK, R.OK(errLogs))
}

// SetErrorLogsViewed 设置日志已读
// @Summary 设置日志已读
// @Description 设置日志已读
// @Tags 日志操作
// @Produce application/json
// @Param Authorization	header	string		true	"token"
// @Param errorLogIds	body	[]string	true	"日志Id"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=string} "返回更新数量"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /logs/errors/update [post]
func SetErrorLogsViewed(c *gin.Context) {
	var errorLogIds []string
	if err := c.ShouldBind(&errorLogIds); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	var successCount = 0
	for _, errorLogId := range errorLogIds {
		filter := bson.D{{"id", errorLogId}}
		update := bson.D{{"$set", bson.D{{"viewed", true}}}}
		successCount += int(errorLogMgo.UpdateOne(filter, update).ModifiedCount)
	}
	c.JSON(http.StatusOK, R.OK(successCount))
}

// SetErrorLogsViewed 删除异常日志
// @Summary 删除异常日志
// @Description 删除异常日志
// @Tags 日志操作
// @Produce application/json
// @Param Authorization	header	string		true	"token"
// @Param errorLogIds	body	[]string	true	"日志Id"
// @Security ApiKeyAuth
// @Success 200 {object} R.ResponseEntity{data=string} "返回删除数量"
// @Failure 400 {object} R.ResponseEntity "返回失败信息"
// @Router /logs/errors/delete [post]
func DeleteErrorLogs(c *gin.Context) {
	var errorLogIds []string
	if err := c.ShouldBind(&errorLogIds); err != nil {
		c.JSON(http.StatusBadRequest, R.Error(http.StatusBadRequest, "", nil))
		return
	}
	var successCount = 0
	for _, errorLogId := range errorLogIds {
		filter := bson.D{{"id", errorLogId}}
		successCount += int(errorLogMgo.Delete(filter))
	}
	c.JSON(http.StatusOK, R.OK(successCount))
}
