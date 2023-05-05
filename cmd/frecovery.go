package cmd

import (
	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/docker"
	"gitee.com/zengtao321/frdocker/frecovery"
	"gitee.com/zengtao321/frdocker/logger"
	"github.com/spf13/cobra"
)

var (
	registryURL      string
	networkInterface string
	color            bool
)

var frecoveryCmd = &cobra.Command{
	Use:   "frecovery",
	Short: "The entry command of Frdocker",
	Long:  "The entry command of running Frdocker",
	Run: func(cmd *cobra.Command, args []string) {
		runFrecovery()
	},
}

func init() {
	frecoveryCmd.Flags().StringVarP(&registryURL, "registryURL", "r", "http://localhost:8030/frecovery/conf", "The URL of the system registry")
	frecoveryCmd.Flags().StringVarP(&networkInterface, "networkInterface", "n", "br-7651c77b1278", "The network interface of the docker network")
	frecoveryCmd.Flags().BoolVarP(&color, "color", "c", false, "Whether to print colorful logs")
}

func runFrecovery() {
	logger := logger.NewLogger("test.log", color)
	dockerCli, err := docker.NewDockerCLI(logger)
	if err != nil {
		logger.Fatal("docker client init failed: ", err)
	}
	dbCli, err := db.NewMongoDB()
	if err != nil {
		logger.Fatal("database client init failed: ", err)
	}
	frecoveryApp := frecovery.NewFrecoveryApp(registryURL, networkInterface, dockerCli, logger, dbCli)
	frecoveryApp.Run()
}
