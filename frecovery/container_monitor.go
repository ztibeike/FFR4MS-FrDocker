package frecovery

import (
	"errors"
)

type ContainerMonitor struct {
	Id           string                   // 容器标识符
	FSMs         map[string]*StateFSM     // 容器状态
	runningState map[string]*StateFSMNode // 正在运行的traceId当前状态
	Metric       *ContainerMetric         // 容器指标
}

func NewContainerMonitor(id string) *ContainerMonitor {
	return &ContainerMonitor{
		Id:           id,
		FSMs:         make(map[string]*StateFSM),
		runningState: make(map[string]*StateFSMNode),
	}
}

func (monitor *ContainerMonitor) UpdateContainerState(httpInfo *HttpInfo, callback MonitorStateCallBack) error {
	node := monitor.getStateFSMNode(httpInfo)
	if node != nil {
		monitor.runningState[httpInfo.TraceId] = node
		node.State.EnsureCallback(callback)
		node.State.Update(httpInfo)
		if httpInfo.IsEndContainerProcess(monitor.Id) {
			delete(monitor.runningState, httpInfo.TraceId)
		}
		return nil
	}
	return errors.New("invalid order of httpInfo")
}

func (monitor *ContainerMonitor) getStateFSMNode(httpInfo *HttpInfo) *StateFSMNode {
	traceId := httpInfo.TraceId
	if httpInfo.IsStartContainerProcess(monitor.Id) {
		api := httpInfo.URL
		fsm := monitor.getStateFSM(api)
		node := fsm.GetFirstNode()
		if node == nil {
			node = fsm.AddStateFSMNode(httpInfo)
		}
		return node
	}
	node := monitor.runningState[traceId]
	fsm := monitor.getStateFSM(node.API)
	// FSM未建立完毕, node信息不全
	if node.To == "" && httpInfo.IsLeaveContainer(monitor.Id) {
		node.To = httpInfo.Dst.Name
	}
	if node.IsLeaveState(httpInfo) {
		return node
	}
	if httpInfo.IsEnterContainer(monitor.Id) {
		if node.Next == fsm.Tail {
			fsm.AddStateFSMNode(httpInfo)
		}
		return node.Next
	}
	return nil
}

func (monitor *ContainerMonitor) getStateFSM(api string) *StateFSM {
	if fsm, ok := monitor.FSMs[api]; ok {
		return fsm
	}
	fsm := NewStateFSM(monitor.Id, api)
	monitor.FSMs[api] = fsm
	return fsm
}
