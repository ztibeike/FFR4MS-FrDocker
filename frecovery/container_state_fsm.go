package frecovery

import "sync"

// 状态有限机中的工作状态节点;
// 一个状态由http(dst=current)开启, 由http(src=current)关闭
type StateFSMNode struct {
	Id    string          `json:"id" bson:"id"`       // 容器标识符
	API   string          `json:"api" bson:"api"`     // 服务API
	From  string          `json:"from" bson:"from"`   // 进入状态的请求的来源服务/网关(在当前设计下是网关)
	To    string          `json:"to" bson:"to"`       // 离开状态的请求的目标服务/网关(在当前设计下是网关)
	State *ContainerState `json:"state" bson:"state"` // 状态
	Next  *StateFSMNode   `json:"next" bson:"next"`   // 下一个状态
	Prev  *StateFSMNode   `json:"prev" bson:"prev"`   // 上一个状态
}

func (node *StateFSMNode) IsLeaveState(httpInfo *HttpInfo) bool {
	return httpInfo.IsLeaveContainer(node.Id) && httpInfo.Dst.Name == node.To
}

// 容器状态有限机
type StateFSM struct {
	Id       string // 容器标识符
	API      string // 服务API
	Size     int
	Head     *StateFSMNode   // 状态链表头
	Tail     *StateFSMNode   // 状态链表尾
	AllNodes []*StateFSMNode // 所有状态节点
	mu       sync.RWMutex    // 锁
}

func NewStateFSM(id, api string) *StateFSM {
	fsm := &StateFSM{
		Id:   id,
		API:  api,
		Size: 0,
		Head: &StateFSMNode{},
		Tail: &StateFSMNode{},
	}
	fsm.Head.Next = fsm.Tail
	fsm.Tail.Prev = fsm.Head
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
		From:  other.Name,
		Next:  fsm.Tail,
		Prev:  fsm.Tail.Prev,
		State: NewContainerState(fsm.Id),
	}
	fsm.Size += 1
	fsm.AllNodes = append(fsm.AllNodes, node)
	fsm.Tail.Prev.Next = node
	fsm.Tail.Prev = node
	return node
}

// 获取第一个节点
func (fsm *StateFSM) GetFirstNode() *StateFSMNode {
	fsm.mu.RLock()
	defer fsm.mu.RUnlock()
	if fsm.Size == 0 {
		return nil
	}
	return fsm.Head.Next
}
