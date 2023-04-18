package logger

import (
	"io"
	"os"

	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/utils"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	log        *logrus.Logger
	logFile    *os.File
	logChannel chan LogInfo
	colored    bool
}

func NewLogger(logFile string) *Logger {
	logger := &Logger{
		log:        logrus.New(),
		logChannel: make(chan LogInfo, 100),
		colored:    false,
	}
	logger.initLogger(logFile)
	return logger
}

func (logger *Logger) Info(args ...interface{}) {
	logInfo := NewLogInfo(logrus.InfoLevel, args...)
	logger.logChannel <- logInfo
}

func (logger *Logger) Error(args ...interface{}) {
	logInfo := NewLogInfo(logrus.ErrorLevel, args...)
	logger.logChannel <- logInfo
}

func (logger *Logger) Debug(args ...interface{}) {
	logInfo := NewLogInfo(logrus.DebugLevel, args...)
	logger.logChannel <- logInfo
}

func (logger *Logger) Fatal(args ...interface{}) {
	logInfo := NewLogInfo(logrus.FatalLevel, args...)
	logger.logChannel <- logInfo
}

func (logger *Logger) Trace(args ...interface{}) {
	logInfo := NewLogInfo(logrus.TraceLevel, args...)
	logger.logChannel <- logInfo
}

func (logger *Logger) Warn(args ...interface{}) {
	logInfo := NewLogInfo(logrus.WarnLevel, args...)
	logger.logChannel <- logInfo
}

// 同步记录并发日志
func (logger *Logger) initLogger(logFile string) {
	logger.setDefaultFormatter()
	logger.setDefaultOutput(logFile)
	go func() {
		for logInfo := range logger.logChannel {
			var message string
			if logger.colored {
				message = logInfo.ColoredMessage
			} else {
				message = logInfo.Message
			}
			logger.log.Logln(logInfo.Level, message)
			if logInfo.Level == logrus.FatalLevel {
				os.Exit(1)
			}
		}
	}()
}

func (logger *Logger) setDefaultFormatter() {
	formatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
		ForceColors:     true,
	}
	logger.log.SetFormatter(formatter)
}

func (logger *Logger) SetFormatter(formatter logrus.Formatter) {
	logger.log.SetFormatter(formatter)
}

func (logger *Logger) setDefaultOutput(logFile string) {
	if exist, err := utils.PathExists(config.LOG_FILE_ROOT_PATH); err != nil || !exist {
		os.MkdirAll(config.LOG_FILE_ROOT_PATH, os.ModePerm)
	}
	logFile = utils.PathJoin(config.LOG_FILE_ROOT_PATH, logFile)
	logger.logFile, _ = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logStream := io.MultiWriter(os.Stdout, logger.logFile)
	logger.log.SetOutput(logStream)
}

func (logger *Logger) SetOutput(output io.Writer) {
	logger.log.SetOutput(output)
}

func (logger *Logger) SetColored(colored bool) {
	logger.colored = colored
}

// 关闭日志记录
func (logger *Logger) Close() {
	close(logger.logChannel)
	logger.logFile.Close()
}
