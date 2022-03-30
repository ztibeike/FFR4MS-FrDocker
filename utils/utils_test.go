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
