package frecovery

import (
	"errors"
)

type ContainerStatus struct {
	Id                   string                   // 容器标识符
	FSMs                 map[string]*StateFSM     // 容器状态
	stateMonitorHandler  StateMonitorHandlerFunc  // 状态监控回调函数
	runningState         map[string]*StateFSMNode // 正在运行的traceId当前状态
	Metric               *ContainerMetric         // 容器指标
	metricMonitorHandler MetricMonitorHandlerFunc // 指标监控回调函数
}

func NewContainerStatus(id string) *ContainerStatus {
	return &ContainerStatus{
		Id:           id,
		FSMs:         make(map[string]*StateFSM),
		runningState: make(map[string]*StateFSMNode),
	}
}

func (status *ContainerStatus) SetStateMonitorHandler(stateMonitorHandler StateMonitorHandlerFunc) {
	status.stateMonitorHandler = stateMonitorHandler
}

func (status *ContainerStatus) SetMetricMonitorHandler(metricMonitorHandler MetricMonitorHandlerFunc) {
	status.metricMonitorHandler = metricMonitorHandler
}

func (status *ContainerStatus) UpdateContainerState(httpInfo *HttpInfo) error {
	traceId := httpInfo.TraceId
	if httpInfo.IsStartContainerProcess(status.Id) {
		api := httpInfo.URL
		fsm := status.getStateFSM(api)
		node := fsm.GetFirstNode()
		if node == nil {
			node = fsm.AddStateFSMNode(httpInfo)
		}
		status.runningState[traceId] = node
		node.state.Update(httpInfo)
		return nil
	}
	node := status.runningState[traceId]
	fsm := status.getStateFSM(node.API)
	if node.to == "" && httpInfo.IsLeaveContainer(status.Id) {
		node.to = httpInfo.Dst.Name
	}
	if node.IsLeaveState(httpInfo) {
		node.state.Update(httpInfo)
		if httpInfo.IsEndContainerProcess(status.Id) {
			delete(status.runningState, traceId)
		}
		return nil
	}
	if httpInfo.IsEnterContainer(status.Id) {
		if node.next == fsm.tail {
			fsm.AddStateFSMNode(httpInfo)
		}
		status.runningState[traceId] = node.next
		node.next.state.Update(httpInfo)
		return nil
	}
	return errors.New("invalid order of httpInfo")
}

func (status *ContainerStatus) getStateFSM(api string) *StateFSM {
	if fsm, ok := status.FSMs[api]; ok {
		return fsm
	}
	fsm := NewStateFSM(status.Id, api, status.stateMonitorHandler)
	status.FSMs[api] = fsm
	return fsm
}
