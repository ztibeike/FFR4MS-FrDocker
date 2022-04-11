package logger

import (
	"io"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	isExist, _ := PathExists("/var/log/frdocker")
	if !isExist {
		os.Mkdir("/var/log/frdocker", os.ModePerm)
	}
	file1, _ := os.OpenFile("/var/log/frdocker/test1.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 6)
	mw := io.MultiWriter(os.Stdout, file1)
	log.SetOutput(mw)
	log.Infoln("abc")
	file1.Close()
	file2, _ := os.OpenFile("/var/log/frdocker/test2.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 6)
	mw = io.MultiWriter(os.Stdout, file2)
	log.SetOutput(mw)
	log.Infoln("def")
	file2.Close()
	file1, _ = os.OpenFile("/var/log/frdocker/test1.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 6)
	mw = io.MultiWriter(os.Stdout, file1)
	log.SetOutput(mw)
	log.Infoln("ghi")
	file1.Close()
	file2, _ = os.OpenFile("/var/log/frdocker/test2.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 6)
	mw = io.MultiWriter(os.Stdout, file2)
	log.SetOutput(mw)
	log.Infoln("jkl")
	file2.Close()
}
