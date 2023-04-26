package main

import (
	"gitee.com/zengtao321/frdocker/frecovery"
	"gitee.com/zengtao321/frdocker/utils/docker"
	"gitee.com/zengtao321/frdocker/utils/logger"
)

func main() {
	logger := logger.NewLogger("test.log")
	dockerCli := docker.NewDockerCLI(logger)
	frecoveryApp := frecovery.NewFrecoveryApp("http://localhost:8030/frecovery/conf", "br-7651c77b1278", dockerCli, logger)
	frecoveryApp.Run()
}
