package utils

import (
	"fmt"
	"testing"
)

func TestGetConfigFromEureka(t *testing.T) {
	containers := GetConfigFromEureka("http://localhost:8030/getConf")
	fmt.Println(containers)
}

func TestGetDockerVerison(t *testing.T) {
	version := GetDockerVersion()
	fmt.Println(version)
}

func TestGetContainerLogs(t *testing.T) {
	containerId := "6fccc2a682f5"
	containerLogs, err := GetContainerLogs(containerId, "100")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(containerLogs)
	}
	return
}
