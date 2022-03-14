package utils

import (
	"fmt"
	"frdocker/constants"
	"frdocker/types"
)

func StateMonitor(IP string, httpChan chan *types.HttpInfo) {
	// var traceIdStateMap = make(map[string]types.State)
	obj, _ := constants.IPServiceContainerMap.Get(IP)
	var container = obj.(*types.Container)
	fmt.Printf("[Monitoring] Container: Group(%s) IP(%s) ID(%s)\n", container.Group, container.IP, container.ID[:10])
	for httpInfo := range httpChan {
		fmt.Println(*httpInfo)
	}
}
