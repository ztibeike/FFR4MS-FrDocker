package utils

import (
	"fmt"
	"frdocker/constants"
	"os"
	"os/signal"
)

func SetupCloseHandler() {
	sigalChan := make(chan os.Signal, 1)
	signal.Notify(sigalChan, os.Interrupt)
	<-sigalChan
	for IP, ch := range constants.IPChanMap {
		close(ch)
		delete(constants.IPChanMap, IP)
		fmt.Println(IP)
	}
	os.Exit(1)
}
