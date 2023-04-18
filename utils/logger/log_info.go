package logger

import (
	"fmt"

	"gitee.com/zengtao321/frdocker/config"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type LogInfo struct {
	Level          logrus.Level
	Args           []interface{}
	Message        string
	ColoredMessage string
}

func NewLogInfo(level logrus.Level, args ...interface{}) LogInfo {
	logInfo := LogInfo{
		Level: level,
		Args:  args,
	}
	generateColoredMessage(&logInfo)
	return logInfo
}

func generateColoredMessage(logInfo *LogInfo) {
	var c *color.Color
	switch logInfo.Level {
	case logrus.InfoLevel:
		c = color.New(config.LOG_INFO_COLOR)
	case logrus.ErrorLevel:
		c = color.New(config.LOG_ERROR_COLOR)
	case logrus.DebugLevel:
		c = color.New(config.LOG_DEBUG_COLOR)
	case logrus.FatalLevel:
		c = color.New(config.LOG_FATAL_COLOR)
	case logrus.TraceLevel:
		c = color.New(config.LOG_TRACE_COLOR)
	case logrus.WarnLevel:
		c = color.New(config.LOG_WARN_COLOR)
	}
	if len(logInfo.Args) == 1 {
		logInfo.Message = fmt.Sprintf(logInfo.Args[0].(string))
		logInfo.ColoredMessage = c.Sprintf(logInfo.Args[0].(string))
	} else {
		logInfo.Message = fmt.Sprintf(logInfo.Args[0].(string), logInfo.Args[1:]...)
		logInfo.ColoredMessage = c.Sprintf(logInfo.Args[0].(string), logInfo.Args[1:]...)
	}
}
