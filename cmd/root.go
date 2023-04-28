package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "frdocker",
	Short: "Frdocker is a docker monitoring and fault localization tool for microservice systems",
	Long: `A docker monitoring tool that monitor communication messages between microservices and
			peformance metrics of docker containers and locate faults for microservice systems`,
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(frecoveryCmd)
}
