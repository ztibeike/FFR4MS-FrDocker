package frecovery

import (
	"errors"

	"gitee.com/zengtao321/frdocker/docker"
	"github.com/panjf2000/ants/v2"
)

type ContainerMonitor struct {
	Id           string                   `json:"id" bson:"id"`                   // 容器标识符
	ContainerId  string                   `json:"containerId" bson:"containerId"` // 容器ID
	FSMs         map[string]*StateFSM     `json:"fsms" bson:"fsms"`               // 容器状态
	runningState map[string]*StateFSMNode // 正在运行的traceId当前状态
	Metric       *ContainerMetric         `json:"metric" bson:"metric"` // 容器指标
}

func NewContainerMonitor(id, containerId string) *ContainerMonitor {
	return &ContainerMonitor{
		Id:           id,
		ContainerId:  containerId,
		FSMs:         make(map[string]*StateFSM),
		runningState: make(map[string]*StateFSMNode),
		Metric:       NewContainerMetric(id, containerId),
	}
}

func (monitor *ContainerMonitor) UpdateContainerMetric(dockerCli *docker.DockerCLI) error {
	return monitor.Metric.Update(dockerCli)
}

func (monitor *ContainerMonitor) UpdateContainerEcc(ecc, thresh float64) {
	monitor.Metric.UpdateEcc(ecc, thresh)
}

func (monitor *ContainerMonitor) UpdateContainerState(httpInfo *HttpInfo, callback MonitorStateCallBack, pool *ants.Pool) error {
	if monitor.runningState == nil {
		monitor.runningState = make(map[string]*StateFSMNode)
	}
	node := monitor.getStateFSMNode(httpInfo)
	if node != nil {
		monitor.runningState[httpInfo.TraceId] = node
		node.State.EnsureCallback(callback)
		node.State.Update(httpInfo, pool)
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
	if monitor.FSMs == nil {
		monitor.FSMs = make(map[string]*StateFSM)
	}
	if fsm, ok := monitor.FSMs[api]; ok {
		return fsm
	}
	fsm := NewStateFSM(monitor.Id, api)
	monitor.FSMs[api] = fsm
	return fsm
}
