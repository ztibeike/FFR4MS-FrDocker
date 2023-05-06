package entity

import (
	"context"
	"math"
	"sync"
	"time"

	"gitee.com/zengtao321/frdocker/frecovery/algo"
)

type Pending struct {
	TraceId string
	Ch      chan time.Time
	Start   time.Time
	End     time.Time
}

func NewPending(traceId string, start time.Time) *Pending {
	return &Pending{
		TraceId: traceId,
		Ch:      make(chan time.Time, 1),
		Start:   start,
	}
}

type ContainerState struct {
	API             string              // 服务API
	Var             []float64           // 方差
	Sigma           float64             // 标准差
	Ecc             float64             // 离心率
	Thresh          float64             // 阈值
	MaxTime         int64               // 最大时间
	MinTime         int64               // 最小时间
	Cnt             int64               // 计数
	pending         map[string]*Pending // 挂起等待被监测的状态
	abnormalHandler AbnormalHandlerFunc // 异常处理函数
	mu              sync.RWMutex        // 读写锁
}

func NewContainerState(api string, dataLen int, abnormalHandler AbnormalHandlerFunc) *ContainerState {
	return &ContainerState{
		API:             api,
		Var:             make([]float64, dataLen),
		Sigma:           0.0,
		Ecc:             0.0,
		Thresh:          0.0,
		MaxTime:         0,
		MinTime:         0,
		Cnt:             0,
		pending:         make(map[string]*Pending),
		abnormalHandler: abnormalHandler,
	}
}

// 更新状态，返回更新结果(正常/异常)
func (state *ContainerState) Update(httpInfo *HttpInfo) bool {
	// 如果存在traceId对应的pending
	if pending, ok := state.pending[httpInfo.TraceId]; ok {
		pending.Ch <- httpInfo.Timestamp
	}
	// 如果不存在，新建pending
	state.addPending(httpInfo.TraceId, httpInfo.Timestamp)
	return true
}

func (state *ContainerState) addPending(traceId string, start time.Time) {
	pending := NewPending(traceId, start)
	state.pending[traceId] = pending
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(state.MaxTime))
	go state.watch(traceId, ctx, cancel)
}

// 监测状态
func (state *ContainerState) watch(traceId string, ctx context.Context, cancel context.CancelFunc) {
	pending := state.pending[traceId]
	defer cancel()
	for {
		select {
		case t := <-pending.Ch:
			state.updateState(math.Abs(float64(t.Sub(pending.Start).Nanoseconds())))
			state.check()
			state.removePending(traceId)
			return
		case <-ctx.Done():
			state.updateState(float64(state.MaxTime))
			state.check()
			state.removePending(traceId)
			return
		}
	}
}

func (state *ContainerState) removePending(traceId string) {
	pending := state.pending[traceId]
	close(pending.Ch)
	delete(state.pending, traceId)
}

func (state *ContainerState) updateState(interval float64) {
	// 更新状态加写锁
	state.mu.Lock()
	defer state.mu.Unlock()
	state.Cnt++
	data := []float64{interval}
	if state.Cnt == 1 {
		copy(state.Var, data)
		return
	}
	state.Ecc, state.Thresh, state.Var, state.Sigma = algo.CalculateWithHistory(data, state.Var, state.Sigma, state.Cnt)
}

func (state *ContainerState) check() {
	// 检查状态加读锁
	state.mu.RLock()
	defer state.mu.RUnlock()
	if state.Ecc > state.Thresh {
		state.abnormalHandler()
	}
}
