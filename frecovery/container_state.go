package frecovery

import (
	"context"
	"math"
	"sync"
	"time"

	"gitee.com/zengtao321/frdocker/config"
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
		Ch:      make(chan time.Time),
		Start:   start,
	}
}

type ContainerState struct {
	Mean           []float64               // 均值
	Sigma          float64                 // 标准差
	Ecc            float64                 // 离心率
	Thresh         float64                 // 阈值
	MaxTime        int64                   // 最大时间
	MinTime        int64                   // 最小时间
	Cnt            int64                   // 计数
	pending        map[string]*Pending     // 挂起等待被监测的状态
	monitorHandler StateMonitorHandlerFunc // 异常处理函数
	mu             sync.RWMutex            // 读写锁
}

func NewContainerState(monitorHandler StateMonitorHandlerFunc) *ContainerState {
	return &ContainerState{
		Mean:           make([]float64, config.TEDA_DATA_LEN),
		Sigma:          0.0,
		Ecc:            0.0,
		Thresh:         0.0,
		MaxTime:        int64(60 * time.Second),
		MinTime:        0,
		Cnt:            0,
		pending:        make(map[string]*Pending),
		monitorHandler: monitorHandler,
	}
}

// 更新状态，返回更新结果(正常/异常)
func (state *ContainerState) Update(httpInfo *HttpInfo) bool {
	// 如果存在traceId对应的pending
	if pending, ok := state.pending[httpInfo.TraceId]; ok {
		pending.Ch <- httpInfo.Timestamp
		state.removePending(httpInfo.TraceId)
	}
	// 如果不存在，新建pending
	state.addPending(httpInfo.TraceId, httpInfo.Timestamp)
	return true
}

func (state *ContainerState) addPending(traceId string, start time.Time) {
	pending := NewPending(traceId, start)
	state.pending[traceId] = pending
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(state.MaxTime*config.TEDA_TIMEOUT_FACTOR))
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
			state.invokeMonitorHandler(traceId)
			return
		case <-ctx.Done():
			state.updateState(float64(state.MaxTime * config.TEDA_TIMEOUT_FACTOR))
			state.invokeMonitorHandler(traceId)
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
		copy(state.Mean, data)
		state.Thresh = (math.Pow(config.TEDA_N_SIGMA, 2) + 1) / float64(2*state.Cnt)
		return
	}
	state.Ecc, state.Thresh, state.Mean, state.Sigma = algo.CalculateWithHistory(data, state.Mean, state.Sigma, state.Cnt)
	state.MaxTime = int64(state.Mean[0] + config.TEDA_N_SIGMA*math.Sqrt(state.Sigma))
	state.MinTime = int64(state.Mean[0] - config.TEDA_N_SIGMA*math.Sqrt(state.Sigma))
}

func (state *ContainerState) invokeMonitorHandler(traceId string) {
	state.monitorHandler(traceId, state)
}
