package utils

import (
	"os"
	"strconv"
)

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

func PathJoin(paths ...string) string {
	var absPath string
	for _, path := range paths {
		absPath += string(os.PathSeparator) + path
	}
	return absPath
}

func GenerateIdFromIPAndPort(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}
