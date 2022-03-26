package logger

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
		ForceColors:     true,
	})
}

func Info(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	var output string
	c := color.New(color.FgHiBlue)
	if len(args) == 1 {
		output = c.Sprintf(args[0].(string))
	} else {
		output = c.Sprintf(args[0].(string), args[1:]...)
	}
	log.Infof(output)
}

func Infoln(args ...interface{}) {
	var coloredArgs []interface{}
	c := color.New(color.FgHiBlue)
	for _, arg := range args {
		coloredArgs = append(coloredArgs, c.Sprintf("%v", arg))
	}
	log.Infoln(coloredArgs...)
}

func Trace(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	var output string
	c := color.New(color.FgCyan)
	if len(args) == 1 {
		output = c.Sprintf(args[0].(string))
	} else {
		output = c.Sprintf(args[0].(string), args[1:]...)
	}
	log.Infof(output)
}

func Error(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	var output string
	c := color.New(color.FgRed)
	if len(args) == 1 {
		output = c.Sprintf(args[0].(string))
	} else {
		output = c.Sprintf(args[0].(string), args[1:]...)
	}
	log.Errorf(output)
}

func Errorln(args ...interface{}) {
	var coloredArgs []interface{}
	c := color.New(color.FgRed)
	for _, arg := range args {
		coloredArgs = append(coloredArgs, c.Sprintf("%v", arg))
	}
	log.Errorln(coloredArgs...)
}

func Fatal(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	var output string
	c := color.New(color.FgHiRed)
	if len(args) == 1 {
		output = c.Sprintf(args[0].(string))
	} else {
		output = c.Sprintf(args[0].(string), args[1:]...)
	}
	log.Fatalf(output)
}

func Fatalln(args ...interface{}) {
	var coloredArgs []interface{}
	c := color.New(color.FgHiRed)
	for _, arg := range args {
		coloredArgs = append(coloredArgs, c.Sprintf("%v", arg))
	}
	log.Fatalln(coloredArgs...)
}

func Warn(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	var output string
	c := color.New(color.FgYellow)
	if len(args) == 1 {
		output = c.Sprintf(args[0].(string))
	} else {
		output = c.Sprintf(args[0].(string), args[1:]...)
	}
	log.Warnf(output)
}

func Warnln(args ...interface{}) {
	var coloredArgs []interface{}
	c := color.New(color.FgYellow)
	for _, arg := range args {
		coloredArgs = append(coloredArgs, c.Sprintf("%v", arg))
	}
	log.Warnln(coloredArgs...)
}
