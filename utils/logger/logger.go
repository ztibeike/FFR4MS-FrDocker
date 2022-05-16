package logger

import (
	"fmt"
	"os"
	"strings"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/settings"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type LogFileInfo struct {
	IP    string
	Logs  string
	Level logrus.Level
}

func NewLogFileInfo(ip string, logs string, level logrus.Level) *LogFileInfo {
	return &LogFileInfo{
		IP:    ip,
		Logs:  logs,
		Level: level,
	}
}

var log = logrus.New()
var logChan = make(chan *LogFileInfo, 100)
var IPFileMap = make(map[string]*os.File)

func init() {
	var format = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
		ForceColors:     true,
	}
	log.SetFormatter(format)
	isExist, _ := PathExists(settings.LOG_FILE_DIR)
	if !isExist {
		os.MkdirAll(settings.LOG_FILE_DIR, os.ModePerm)
	}
	go LogToFile()
}

func GenerateOutput(c *color.Color, args ...interface{}) (string, string) {
	var coloredOutput, output string
	if len(args) == 1 {
		coloredOutput = c.Sprintf(args[0].(string))
		output = args[0].(string)
	} else {
		coloredOutput = c.Sprintf(args[0].(string), args[1:]...)
		output = fmt.Sprintf(args[0].(string), args[1:]...)
	}
	return coloredOutput, output
}

func Info(ip interface{}, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	c := color.New(color.FgHiBlue)
	coloredOutput, output := GenerateOutput(c, args...)
	log.Infof(coloredOutput)
	if ip != nil {
		logChan <- NewLogFileInfo(ip.(string), output, logrus.InfoLevel)
	}
}

func Infoln(ip interface{}, args ...interface{}) {
	var coloredArgs []interface{}
	c := color.New(color.FgHiBlue)
	var stringArgs []string
	for _, arg := range args {
		coloredArgs = append(coloredArgs, c.Sprintf("%v", arg))
		stringArgs = append(stringArgs, fmt.Sprintf("%v", arg))
	}
	log.Infoln(coloredArgs...)
	if ip != nil {
		output := strings.Join(stringArgs, " ")
		output += "\n"
		logChan <- NewLogFileInfo(ip.(string), output, logrus.InfoLevel)
	}
}

func Trace(ip interface{}, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	c := color.New(color.FgCyan)
	coloredOutput, output := GenerateOutput(c, args...)
	log.Infof(coloredOutput)
	if ip != nil {
		logChan <- NewLogFileInfo(ip.(string), output, logrus.InfoLevel)
	}
}

func Error(ip interface{}, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	c := color.New(color.FgRed)
	coloredOutput, output := GenerateOutput(c, args...)
	log.Errorf(coloredOutput)
	if ip != nil {
		logChan <- NewLogFileInfo(ip.(string), output, logrus.ErrorLevel)
	}
}

func Errorln(ip interface{}, args ...interface{}) {
	var coloredArgs []interface{}
	c := color.New(color.FgRed)
	var stringArgs []string
	for _, arg := range args {
		coloredArgs = append(coloredArgs, c.Sprintf("%v", arg))
		stringArgs = append(stringArgs, fmt.Sprintf("%v", arg))
	}
	log.Errorln(coloredArgs...)
	if ip != nil {
		output := strings.Join(stringArgs, " ")
		output += "\n"
		logChan <- NewLogFileInfo(ip.(string), output, logrus.ErrorLevel)
	}
}

func Fatal(ip interface{}, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	c := color.New(color.FgHiRed)
	coloredOutput, output := GenerateOutput(c, args...)
	log.Fatalf(coloredOutput)
	if ip != nil {
		logChan <- NewLogFileInfo(ip.(string), output, logrus.FatalLevel)
	}
}

func Fatalln(ip interface{}, args ...interface{}) {
	var coloredArgs []interface{}
	c := color.New(color.FgHiRed)
	var stringArgs []string
	for _, arg := range args {
		coloredArgs = append(coloredArgs, c.Sprintf("%v", arg))
		stringArgs = append(stringArgs, fmt.Sprintf("%v", arg))
	}
	log.Fatalln(coloredArgs...)
	if ip != nil {
		output := strings.Join(stringArgs, " ")
		output += "\n"
		logChan <- NewLogFileInfo(ip.(string), output, logrus.FatalLevel)
	}
}

func Warn(ip interface{}, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	c := color.New(color.FgYellow)
	coloredOutput, output := GenerateOutput(c, args...)
	log.Warnf(coloredOutput)
	if ip != nil {
		logChan <- NewLogFileInfo(ip.(string), output, logrus.WarnLevel)
	}
}

func Warnln(ip interface{}, args ...interface{}) {
	var coloredArgs []interface{}
	c := color.New(color.FgYellow)
	var stringArgs []string
	for _, arg := range args {
		coloredArgs = append(coloredArgs, c.Sprintf("%v", arg))
		stringArgs = append(stringArgs, fmt.Sprintf("%v", arg))
	}
	log.Warnln(coloredArgs...)
	if ip != nil {
		output := strings.Join(stringArgs, " ")
		output += "\n"
		logChan <- NewLogFileInfo(ip.(string), output, logrus.WarnLevel)
	}
}

func LogToFile() {
	logFile := logrus.New()
	var format = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
	}
	logFile.SetFormatter(format)
	for logInfo := range logChan {
		fp, ok := IPFileMap[logInfo.IP]
		if !ok {
			fileName := fmt.Sprintf("%s/%s-%s.log", settings.LOG_FILE_DIR, commons.Network, logInfo.IP)
			fp, _ = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 6)
			IPFileMap[logInfo.IP] = fp
		}
		logFile.SetOutput(fp)
		logFile.Logf(logInfo.Level, "%s", logInfo.Logs)
	}
	for ip, fp := range IPFileMap {
		fp.Close()
		delete(IPFileMap, ip)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Close() {
	close(logChan)
}
