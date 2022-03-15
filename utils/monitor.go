package utils

import (
	"fmt"
	"frdocker/constants"
	"frdocker/types"
	"math"
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
	var idx = 0
	var httpInfo_start *types.HttpInfo = nil
	var httpInfo_end *types.HttpInfo = nil
	for httpInfo := range httpChan {
		// var now = time.Now().Format("2006-01-02 15:04:05")
		// var otherIP string
		// var action = make([]string, 2)
		// if container.IP != httpInfo.SrcIP {
		// 	otherIP = httpInfo.SrcIP
		// 	action[0] = "Recieve"
		// 	action[1] = "From"
		// } else {
		// 	otherIP = httpInfo.DstIP
		// 	action[0] = "Send"
		// 	action[1] = "to"
		// }
		// obj, _ := constants.IPAllMSMap.Get(otherIP)
		// var ms = obj.(string)
		// if container.Group == "service-a" {
		// 	fmt.Printf("[Checking State] [%s] [TraceId(%s)] Group(%s) IP(%s) ID(%s) %s HTTP %s %s %s(%s)\n",
		// 		now, traceId, container.Group, container.IP, container.ID[:10], action[0], httpInfo.Type, action[1], ms, otherIP)
		// }
		if httpInfo_start == nil {
			httpInfo_start = httpInfo
		} else {
			httpInfo_end = httpInfo
			if len(container.States) <= idx {
				container.States = append(container.States, &types.State{
					Id: &types.StateId{
						StartWith: &types.StateEndpointEvent{
							IP:       httpInfo_start.SrcIP,
							HttpType: httpInfo_start.Type,
						},
						EndWith: &types.StateEndpointEvent{
							IP:       httpInfo_end.DstIP,
							HttpType: httpInfo_end.Type,
						},
					},
					K:        1,
					Variance: &types.Vector{},
				})
			}
			timeInterval := math.Abs(float64(httpInfo_end.Timestamp.Nanosecond() - httpInfo_start.Timestamp.Nanosecond()))
			data := &types.Vector{
				Data: []float64{timeInterval},
			}
			ecc := TEDA(container.States[idx], data)
			var now = time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("\n[Checking State] [%s] [TraceId(%s)] [Group(%s) IP(%s) ID(%s)] [State(%d) TimeInterval(%dns) Eccentricity(%f)]\n",
				now, traceId, container.Group, container.IP, container.ID[:10], idx, int(timeInterval), ecc)
			idx += 1
			httpInfo_start = nil
			httpInfo_end = nil
		}

	}
}

func TEDA(state *types.State, data *types.Vector) float64 {
	if state.K == 1 {
		state.Lock()
		state.Variance.Data = make([]float64, len(data.Data))
		copy(state.Variance.Data, data.Data)
		state.Sigma = 0.0
		state.Ecc = math.NaN()
		state.K = state.K + 1
		state.Unlock()
		return math.NaN()
	}

	state.RLock()
	variance := state.Variance.Copy()
	sigma := state.Sigma
	k := state.K
	state.RUnlock()

	variance = variance.ScaleVec(float64(k-1) / float64(k)).AddVec(data.ScaleVec(1.0 / float64(k)))
	sigma = sigma*(float64(k-1)/float64(k)) + 1.0/float64(k-1)*math.Pow(data.SubVec(variance).Norm(), 2)
	normalized_ecc := 1.0 / float64(2*k) * (1.0 + data.SubVec(variance).T().MulVec(data.SubVec(variance))/sigma)
	state.Lock()
	defer state.Unlock()
	state.Ecc = normalized_ecc
	state.Sigma = sigma
	state.Ecc = normalized_ecc
	state.K = state.K + 1
	copy(state.Variance.Data, variance.Data)
	return normalized_ecc
}
