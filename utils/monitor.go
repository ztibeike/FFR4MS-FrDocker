package utils

import (
	"fmt"
	"frdocker/constants"
	"frdocker/types"
	"time"
)

func StateMonitor(IP string, httpChan chan *types.HttpInfo) {
	// var traceIdStateMap = make(map[string]types.State)
	obj, _ := constants.IPServiceContainerMap.Get(IP)
	var container = obj.(*types.Container)
	fmt.Printf("\n[Monitoring Container] Group(%s) IP(%s) ID(%s)\n", container.Group, container.IP, container.ID[:10])
	var TraceMap = make(map[string]chan *types.HttpInfo) // TraceId为key，每个TraceId开启一个go routine
	for httpInfo := range httpChan {
		var channel chan *types.HttpInfo
		var ok bool
		var traceId = httpInfo.TraceId
		if channel, ok = TraceMap[traceId]; ok {
			channel <- httpInfo
			if IP == httpInfo.SrcIP && httpInfo.Type == "RESPONSE" {
				close(channel)
				delete(TraceMap, traceId)
			}
		} else {
			channel = make(chan *types.HttpInfo)
			go CheckingStateByTraceId(traceId, container, channel)
			channel <- httpInfo
			TraceMap[traceId] = channel
		}
		// fmt.Println(*httpInfo)
	}
}

func CheckingStateByTraceId(traceId string, container *types.Container, httpChan chan *types.HttpInfo) {
	for httpInfo := range httpChan {
		var now = time.Now().Format("2006-01-02 15:04:05")
		var otherIP string
		var action = make([]string, 2)
		if container.IP != httpInfo.SrcIP {
			otherIP = httpInfo.SrcIP
			action[0] = "Recieve"
			action[1] = "From"
		} else {
			otherIP = httpInfo.DstIP
			action[0] = "Send"
			action[1] = "to"
		}
		obj, _ := constants.IPAllMSMap.Get(otherIP)
		var ms = obj.(string)
		fmt.Printf("[Checking State] [%s] [TraceId(%s)] Group(%s) IP(%s) ID(%s) %s HTTP %s %s %s(%s)\n",
			now, traceId, container.Group, container.IP, container.ID[:10], action[0], httpInfo.Type, action[1], ms, otherIP)
	}
}
