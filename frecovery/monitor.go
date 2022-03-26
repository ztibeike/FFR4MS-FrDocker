package frecovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frdocker/constants"
	"frdocker/settings"
	"frdocker/types"
	"frdocker/utils/logger"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

func StateMonitor(IP string, httpChan chan *types.HttpInfo) {
	// var traceIdStateMap = make(map[string]types.State)
	obj, _ := constants.IPServiceContainerMap.Get(IP)
	var container = obj.(*types.Container)
	logger.Info("\n[Monitoring Container] Group(%s) IP(%s) ID(%s)\n", container.Group, container.IP, container.ID[:10])
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
	for traceId, ch := range TraceMap {
		close(ch)
		delete(TraceMap, traceId)
	}
}

func CheckingStateByTraceId(traceId string, container *types.Container, httpChan chan *types.HttpInfo) {
	var idx = 0
	var httpInfo_start *types.HttpInfo = nil
	var httpInfo_end *types.HttpInfo = nil
	// var timeOutIdx = -1
	var checkTimeExceedMap = make(map[int]chan float64)
	for httpInfo := range httpChan {
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
			ch := make(chan float64, 1)
			checkTimeExceedMap[idx] = ch
			go CheckTimeExceedNotEnd(container, traceId, idx, ch)
		} else {
			httpInfo_end = httpInfo
			currentIdx := idx
			idx += 1
			if container.States[currentIdx].Id.EndWith == nil {
				container.States[currentIdx].Id.EndWith = &types.StateEndpointEvent{
					IP:       httpInfo_end.DstIP,
					HttpType: httpInfo_end.Type,
				}
			}
			timeInterval := math.Abs(float64(httpInfo_end.Timestamp.Nanosecond() - httpInfo_start.Timestamp.Nanosecond()))
			ch := checkTimeExceedMap[currentIdx]
			ch <- timeInterval
			data := &types.Vector{
				Data: []float64{timeInterval},
			}
			ecc := TEDA(container.States[currentIdx], data)
			var threshold = (math.Pow(settings.NSIGMA, 2) + 1) / float64((2 * container.States[currentIdx].K))
			var health = container.Health
			if ecc > threshold && health {
				health = CheckHealthByLocalActuator(container, currentIdx)
				if !health {
					MarkContainerUnHealthy(container)
				}
			}
			logger.Trace("\n[Checking State] [TraceId(%s)] [Group(%s) IP(%s) ID(%s)] [State(%d) TimeInterval(%dns) Eccentricity(%f) MinTime(%d) MaxTime(%d) Health(%t)]\n",
				traceId, container.Group, container.IP, container.ID[:10], currentIdx, int(timeInterval), ecc,
				int(container.States[currentIdx].MinTime), int(container.States[currentIdx].MaxTime), health)
			httpInfo_start = nil
			httpInfo_end = nil
		}
	}
	for idx, ch := range checkTimeExceedMap {
		close(ch)
		delete(checkTimeExceedMap, idx)
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
		state.MaxTime = state.Variance.Data[0] + settings.NSIGMA*state.Variance.Data[0]
		state.MinTime = state.Variance.Data[0] - settings.NSIGMA*state.Variance.Data[0]
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
	state.MaxTime = state.Variance.Data[0] + settings.NSIGMA*math.Sqrt(state.Sigma)
	state.MinTime = state.Variance.Data[0] - settings.NSIGMA*math.Sqrt(state.Sigma)
	return normalized_ecc
}

func CheckTimeExceedNotEnd(container *types.Container, traceId string, currentIdx int, channel chan float64) {
	state := container.States[currentIdx]
	if state.K == 1 {
		return
	}
	maxTime := state.MaxTime
	timeoutChan := make(chan int, 1)
	// 设置超时时间，避免因服务崩溃导致本状态一直未结束，没有计算离心率和状态转变时间，从而导致无法检测出服务故障
	go func() {
		time.Sleep(time.Duration(maxTime))
		timeoutChan <- 1
		close(timeoutChan)
	}()
	/*	两个通道：
			1. channel: 由上层goroutine传递，如果从通道中监听到消息，代表本状态结束了，收到的消息未状态转变时间，计算是否超时，如果超时则进行本地健康检查进行最终判断
			2. timeOutChannel: 超时通道，由下层goroutine传递消息，如果从通道中监听到消息，代表在该状态的最大转变时间内仍未结束该状态，则存在超时异常，进行本地健康检查
		两个通道任选其一执行对应的处理，取决于哪个通道最先收到消息

		为什么要这样写：首先考虑到两种情况：
			1. 状态成功结束但超时
			2. 状态一直不结束
		第一种情况只需在状态结束后计算离心率或者检测超时即可判断异常
		第二种情况由于状态不结束因此无法计算离心率，只能通过在超过最大状态转变时间后进行本地健康检查来判断异常，因此设置了一个超时后就会收到消息的通道timeoutChan
	*/
	for {
		select {
		case t, ok := <-channel:
			{
				if !ok {
					// 通道关闭
					return
				}
				if t > maxTime {
					health := CheckHealthByLocalActuator(container, currentIdx)
					logger.Warn("\n[Time Exceed] [TraceId(%s)] [Group(%s) IP(%s) ID(%s)] [State(%d) MaxTime(%d)] [Health(%t)]\n",
						traceId, container.Group, container.IP, container.ID[:10], currentIdx, int(state.MaxTime), health)
					if !health {
						MarkContainerUnHealthy(container)
					}
				}
				return
			}
		case <-timeoutChan:
			{
				health := CheckHealthByLocalActuator(container, currentIdx)
				logger.Warn("\n[Time Exceed] [TraceId(%s)] [Group(%s) IP(%s) ID(%s)] [State(%d) MaxTime(%d)] [Health(%t)]\n",
					traceId, container.Group, container.IP, container.ID[:10], currentIdx, int(state.MaxTime), health)
				if !health {
					MarkContainerUnHealthy(container)
				}
				return
			}
		}

	}
}

func CheckHealthByLocalActuator(container *types.Container, idx int) bool {
	var IP = container.IP
	var port = container.Port
	var client = &http.Client{
		Timeout: settings.LOCAL_HEALTH_CHECK_TIME_OUT * time.Millisecond,
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
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, resp)
	if err != nil || resp.Status != "UP" {
		return false
	}
	return true
}

func MarkContainerUnHealthy(container *types.Container) {
	var IP = container.IP
	constants.IPChanMapMutex.Lock()
	container.Health = false
	logger.Error("\n[Mark Container Unhealthy] [Group(%s) IP(%s) ID(%s)]\n", container.Group, container.IP, container.ID)
	ch := constants.IPChanMap[IP]
	close(ch)
	delete(constants.IPChanMap, IP)
	constants.IPChanMapMutex.Unlock()
	go GatewayReplaceInstance(container)
}

func GatewayReplaceInstance(container *types.Container) {
	var gateway = container.Gateway
	var url = fmt.Sprintf("http://%s/zuulApi/replaceServiceInstance", gateway)
	var replaceInfo = &types.GateWayReplaceService{
		ServiceName:      container.Group,
		DownInstanceHost: container.IP,
		DownInstancePort: container.Port,
	}
	var requestBody, _ = json.Marshal(replaceInfo)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil || response == nil || response.StatusCode != 200 {
		logger.Info("\n[Gateway Replace Instance] [Gateway(%s) Group(%s) Instance(%s)] Gateway Error!\n",
			gateway, container.Group, container.IP)
	}
	logger.Info("\n[Gateway Replace Instance] [Gateway(%s) Group(%s) Instance(%s)] Service Down!\n",
		gateway, container.Group, container.IP)
}
