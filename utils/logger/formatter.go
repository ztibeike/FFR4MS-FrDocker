package logger

import (
	"bytes"
	"fmt"
	"strings"

	"gitee.com/zengtao321/frdocker/config"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type LoggerFormatter struct {
	colored bool
}

func NewLoggerFormatter(colored bool) *LoggerFormatter {
	return &LoggerFormatter{colored: colored}
}

func (formatter *LoggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	if entry.Buffer != nil {
		b = entry.Buffer
	}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())[:4]
	var log string
	if entry.HasCaller() {
		log = fmt.Sprintf("[%s] [%s] [%s:%d] %s\n", timestamp, level, entry.Caller.File, entry.Caller.Line, entry.Message)
	} else {
		log = fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, entry.Message)
	}
	if formatter.colored {
		log = formatter.colorMessage(entry, log)
	}

	b.WriteString(log)
	return b.Bytes(), nil
}

func (formatter *LoggerFormatter) colorMessage(entry *logrus.Entry, log string) string {
	var c *color.Color
	switch entry.Level {
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
	c.EnableColor()
	return c.Sprint(log)
}
