package frecovery

import (
	"context"
	"math"
	"sync"
	"time"

	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/frecovery/algo"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/panjf2000/ants/v2"
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
	Id       string                               // 容器标识符
	Mean     []float64                            // 均值
	Sigma    float64                              // 标准差
	Ecc      float64                              // 离心率
	Thresh   float64                              // 阈值
	MaxTime  int64                                // 最大时间
	MinTime  int64                                // 最小时间
	Cnt      int64                                // 计数
	pending  cmap.ConcurrentMap[string, *Pending] // 挂起等待被监测的状态
	callback MonitorStateCallBack                 // 异常处理函数
	mu       sync.RWMutex                         // 读写锁
}

func NewContainerState(id string) *ContainerState {
	return &ContainerState{
		Id:       id,
		Mean:     make([]float64, config.TEDA_DATA_LEN),
		Sigma:    0.0,
		Ecc:      0.0,
		Thresh:   0.0,
		MaxTime:  int64(60 * time.Second),
		MinTime:  0,
		Cnt:      0,
		pending:  cmap.New[*Pending](),
		callback: nil,
	}
}

func (state *ContainerState) EnsureCallback(callback MonitorStateCallBack) {
	if state.callback == nil {
		state.callback = callback
	}
}

func (state *ContainerState) EnsurePending() {
	state.pending = cmap.New[*Pending]()
}

// 更新状态，返回更新结果(正常/异常)
func (state *ContainerState) Update(httpInfo *HttpInfo, pool *ants.Pool) {
	// 如果存在traceId对应的pending
	if pending, ok := state.pending.Get(httpInfo.TraceId); ok {
		pending.Ch <- httpInfo.Timestamp
		state.removePending(httpInfo.TraceId)
		return
	}
	// 如果不存在，新建pending
	state.addPending(httpInfo.TraceId, httpInfo.Timestamp, pool)
}

func (state *ContainerState) addPending(traceId string, start time.Time, pool *ants.Pool) {
	pending := NewPending(traceId, start)
	state.pending.Set(traceId, pending)
	pool.Submit(func() { state.watch(traceId) })
}

// 监测状态
func (state *ContainerState) watch(traceId string) {
	pending, ok := state.pending.Get(traceId)
	if !ok {
		return
	}
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(state.MaxTime*config.TEDA_TIMEOUT_FACTOR))
	defer cancel()
	for {
		select {
		case t, ok := <-pending.Ch:
			if ok {
				state.updateState(math.Abs(float64(t.Sub(pending.Start).Nanoseconds())))
				state.invokeCallback(traceId)
			}
			// 跳出循环, 实测return不能跳出循环
			goto end
		case <-ctx.Done():
			state.updateState(float64(state.MaxTime * config.TEDA_TIMEOUT_FACTOR))
			state.invokeCallback(traceId)
			// 跳出循环, 实测return不能跳出循环
			goto end
		}
	}
end:
	return
}

func (state *ContainerState) removePending(traceId string) {
	pending, ok := state.pending.Get(traceId)
	if !ok {
		return
	}
	close(pending.Ch)
	state.pending.Remove(traceId)
}

func (state *ContainerState) updateState(interval float64) {
	// 更新状态加写锁
	state.mu.Lock()
	defer state.mu.Unlock()
	state.Cnt++
	data := []float64{interval}
	if state.Cnt == 1 {
		copy(state.Mean, data)
		state.Thresh = (math.Pow(config.TEDA_N_SIGMA, 2) + 1) / float64(2*state.Cnt)
		return
	}
	state.Ecc, state.Thresh, state.Mean, state.Sigma = algo.CalculateWithHistory(data, state.Mean, state.Sigma, state.Cnt)
	state.MaxTime = int64(state.Mean[0] + config.TEDA_N_SIGMA*math.Sqrt(state.Sigma))
	state.MinTime = int64(state.Mean[0] - config.TEDA_N_SIGMA*math.Sqrt(state.Sigma))
}

func (state *ContainerState) invokeCallback(traceId string) {
	if state.callback != nil {
		state.callback(traceId, state)
	}
}
