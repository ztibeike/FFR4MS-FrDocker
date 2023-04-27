package logger

import (
	"io"
	"os"

	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/utils"
	"github.com/sirupsen/logrus"
)

func NewLogger(logFile string, colored bool) *logrus.Logger {
	logger := logrus.New()
	setDefaultFormatter(logger, colored)
	setDefaultOutput(logger, logFile)
	logger.SetLevel(logrus.TraceLevel)
	logger.SetReportCaller(true)
	return logger
}

func setDefaultFormatter(logger *logrus.Logger, colored bool) {
	formatter := NewLoggerFormatter(colored)
	logger.SetFormatter(formatter)
}

func setDefaultOutput(logger *logrus.Logger, logFile string) {
	if exist, err := utils.PathExists(config.LOG_FILE_ROOT_PATH); err != nil || !exist {
		os.MkdirAll(config.LOG_FILE_ROOT_PATH, os.ModePerm)
	}
	logFile = utils.PathJoin(config.LOG_FILE_ROOT_PATH, logFile)
	file, _ := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logStream := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(logStream)
}
