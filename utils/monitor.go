package utils

import (
	"encoding/json"
	"fmt"
	"frdocker/constants"
	"frdocker/types"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
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
	for _, ch := range TraceMap {
		close(ch)
	}
}

func CheckingStateByTraceId(traceId string, container *types.Container, httpChan chan *types.HttpInfo) {
	var idx = 0
	var httpInfo_start *types.HttpInfo = nil
	var httpInfo_end *types.HttpInfo = nil
	// var timeOutIdx = -1
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
			if len(container.States) <= idx {
				container.States = append(container.States, &types.State{
					Id: &types.StateId{
						StartWith: &types.StateEndpointEvent{
							IP:       httpInfo_start.SrcIP,
							HttpType: httpInfo_start.Type,
						},
					},
					K:        1,
					Variance: &types.Vector{},
				})
			}
			go CheckTimeExceedNotEnd(container, traceId, idx, &idx)
		} else {
			httpInfo_end = httpInfo
			currentIdx := idx
			idx += 1
			container.States[currentIdx].Id.EndWith = &types.StateEndpointEvent{
				IP:       httpInfo_end.DstIP,
				HttpType: httpInfo_end.Type,
			}
			timeInterval := math.Abs(float64(httpInfo_end.Timestamp.Nanosecond() - httpInfo_start.Timestamp.Nanosecond()))
			data := &types.Vector{
				Data: []float64{timeInterval},
			}
			ecc := TEDA(container.States[currentIdx], data)
			var now = time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("\n[Checking State] [%s] [TraceId(%s)] [Group(%s) IP(%s) ID(%s)] [State(%d) TimeInterval(%dns) Eccentricity(%f) MinTime(%d) MaxTime(%d) Sigma(%f)]\n",
				now, traceId, container.Group, container.IP, container.ID[:10], currentIdx, int(timeInterval), ecc,
				int(container.States[currentIdx].MinTime), int(container.States[currentIdx].MaxTime), math.Sqrt(container.States[currentIdx].Sigma))

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
		state.MaxTime = state.Variance.Data[0] + constants.NSigma*state.Variance.Data[0]
		state.MinTime = state.Variance.Data[0] - constants.NSigma*state.Variance.Data[0]
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
	state.MaxTime = state.Variance.Data[0] + constants.NSigma*math.Sqrt(state.Sigma)
	state.MinTime = state.Variance.Data[0] - constants.NSigma*math.Sqrt(state.Sigma)
	return normalized_ecc
}

func MarkContainerUnHealthy(container *types.Container) {

}

func CheckTimeExceedNotEnd(container *types.Container, traceId string, currentIdx int, idx *int) {
	state := container.States[currentIdx]
	if state.K == 1 {
		return
	}
	t := time.Duration(state.MaxTime) + 1
	time.Sleep(t)
	if *idx == currentIdx {
		health := CheckHealthByLocalActuator(container, currentIdx)
		log.Printf("\n[Time Exceed] [TraceId(%s)] [Group(%s) IP(%s) ID(%s)] [State(%d) MaxTime(%d)] [Health(%t)]\n",
			traceId, container.Group, container.IP, container.ID[:10], currentIdx, int(state.MaxTime), health)
	}
}

func CheckHealthByLocalActuator(container *types.Container, idx int) bool {
	var IP = container.IP
	var port = container.Port
	var maxTime = container.States[idx].MaxTime
	var client = &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, time.Duration(maxTime)*2)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Duration(maxTime) * 2))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Duration(maxTime) * 2,
		},
	}
	var url = fmt.Sprintf("http://%s:%s/actuator/health", IP, port)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return false
	}
	defer response.Body.Close()
	var resp = &types.ServiceActuatorHealth{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}
	err = json.Unmarshal(body, resp)
	if err != nil || resp.Status != "UP" {
		return false
	}
	return true
}
