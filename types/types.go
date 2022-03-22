package types

import (
	"sync"
	"time"
)

// 微服务/容器
type Container struct {
	IP      string   `bson:"ip"`
	Port    string   `bson:"port"`
	Group   string   `bson:"group"`
	Gateway string   `bson:"gateway"`
	Leaf    bool     `bson:"leaf"`
	Health  bool     `bson:"health"`
	ID      string   `bson:"id"`
	Name    string   `bson:"name"`
	States  []*State `bson:"state"`
}

// 微服务状态
type State struct {
	sync.RWMutex
	Id       *StateId `bson:"id"`
	Ecc      float64  `bson:"ecc"`
	Variance *Vector  `bson:"variance"`
	Sigma    float64  `bson:"sigma"`
	K        int64    `bson:"k"`
	MaxTime  float64  `bson:"maxTime"`
	MinTime  float64  `bson:"minTime"`
}

type StateId struct {
	StartWith *StateEndpointEvent `bson:"startWith"`
	EndWith   *StateEndpointEvent `bson:"endwith"`
}

type StateEndpointEvent struct {
	IP       string `bson:"ip"`
	HttpType string `bson:"httpType"`
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
