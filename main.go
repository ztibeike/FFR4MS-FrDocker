package main

import (
	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/docker"
	"gitee.com/zengtao321/frdocker/frecovery"
	"gitee.com/zengtao321/frdocker/logger"
)

func main() {
	logger := logger.NewLogger("test.log", config.LOG_COLORED)
	dockerCli, err := docker.NewDockerCLI(logger)
	if err != nil {
		logger.Fatal("docker client init failed: ", err)
	}
	dbCli, err := db.NewMongoDB()
	if err != nil {
		logger.Fatal("database client init failed: ", err)
	}
	frecoveryApp := frecovery.NewFrecoveryApp("http://localhost:8030/frecovery/conf", "br-7651c77b1278", dockerCli, logger, dbCli)
	frecoveryApp.Run()
}
