package cmd

import (
	"errors"

	"gitee.com/zengtao321/frdocker/frecovery"
	"gitee.com/zengtao321/frdocker/web"
	"github.com/spf13/cobra"
)

var frecoveryCmd = &cobra.Command{
	Use:   "frecovery [ifaceName] [confPath]",
	Short: "Monitoring the containerized microservices-based system.",
	Long:  "By monitoring the containerized microservices-based system, it's able to locate the fault service.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("error arguments")
		}
		return nil
	},
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		executeFrecovery(args)
	},
}

func executeFrecovery(args []string) {
	go web.SetupWebHander()
	frecovery.RunFrecovery(args[0], args[1])
}
