package types

import (
	"sync"
	"time"
)

// 微服务/容器
type Container struct {
	IP      string
	Port    string
	Group   string
	Gateway string
	Leaf    bool
	Health  bool
	ID      string
	Name    string
	States  []*State
}

// 微服务状态
type State struct {
	Id       *StateId
	Ecc      float64
	Variance *Vector
	Sigma    float64
	K        int64
	MaxTime  float64
	Mutex    *sync.RWMutex
}

type StateId struct {
	StartWith *StateEndpointEvent
	EndWith   *StateEndpointEvent
}

type StateEndpointEvent struct {
	IP       string
	HttpType string
}

// 服务发现配置: 来自Eureka Server或配置文件
type ServiceDiscoveryConfig struct {
	ServiceDetailList []ServiceDetail
	ServiceGroupList  []string
	GatewayAddrList   []string
}

type ServiceDetail struct {
	Addr    string
	IP      string
	Port    string
	Group   string
	Gateway string
	Health  string
}

type HttpInfo struct {
	Type      string
	SrcIP     string
	SrcPort   string
	DstIP     string
	DstPort   string
	TraceId   string
	Timestamp time.Time
	Internal  bool
}

type ServiceGroup struct {
	Gateway  string
	Services []string
}
