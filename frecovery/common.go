package frecovery

type AntsTaskWrapper func()

type MonitorMetricCallback func(metric *ContainerMetric)

type MonitorStateCallBack func(traceId string, state *ContainerState)
