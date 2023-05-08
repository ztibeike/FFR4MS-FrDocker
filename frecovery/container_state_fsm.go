package frecovery

import "sync"

// 状态有限机中的工作状态节点;
// 一个状态由http(dst=current)开启, 由http(src=current)关闭
type StateFSMNode struct {
	Id    string          // 容器标识符
	API   string          // 服务API
	from  string          // 进入状态的请求的来源服务/网关(在当前设计下是网关)
	to    string          // 离开状态的请求的目标服务/网关(在当前设计下是网关)
	state *ContainerState // 状态
	next  *StateFSMNode   // 下一个状态
	prev  *StateFSMNode   // 上一个状态
}

func (node *StateFSMNode) IsLeaveState(httpInfo *HttpInfo) bool {
	return httpInfo.IsLeaveContainer(node.Id) && httpInfo.Dst.Name == node.to
}

// 容器状态有限机
type StateFSM struct {
	Id             string // 容器标识符
	API            string // 服务API
	size           int
	head           *StateFSMNode           // 状态链表头
	tail           *StateFSMNode           // 状态链表尾
	mu             sync.RWMutex            // 锁
	monitorHandler StateMonitorHandlerFunc // 状态检测回调函数
}

func NewStateFSM(containerId, api string, monitorHandler StateMonitorHandlerFunc) *StateFSM {
	fsm := &StateFSM{
		Id:             containerId,
		API:            api,
		size:           0,
		head:           &StateFSMNode{},
		tail:           &StateFSMNode{},
		monitorHandler: monitorHandler,
	}
	fsm.head.next = fsm.tail
	fsm.tail.prev = fsm.head
	return fsm
}

// 添加状态节点
func (fsm *StateFSM) AddStateFSMNode(httpInfo *HttpInfo) *StateFSMNode {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	other := httpInfo.GetOtherRole(fsm.Id)
	// 添加状态节点
	node := &StateFSMNode{
		Id:    fsm.Id,
		API:   fsm.API,
		from:  other.Name,
		next:  fsm.tail,
		prev:  fsm.tail.prev,
		state: NewContainerState(fsm.monitorHandler),
	}
	fsm.size += 1
	fsm.tail.prev.next = node
	fsm.tail.prev = node
	return node
}

// 获取第一个节点
func (fsm *StateFSM) GetFirstNode() *StateFSMNode {
	fsm.mu.RLock()
	defer fsm.mu.RUnlock()
	if fsm.size == 0 {
		return nil
	}
	return fsm.head.next
}
