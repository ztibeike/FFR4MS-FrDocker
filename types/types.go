package types

import (
	"sync"
	"time"
)

// 微服务/容器
type Container struct {
	IP      string              `bson:"ip" json:"ip" binding:"required"`
	Port    string              `bson:"port" json:"port"`
	Group   string              `bson:"group" json:"group"`
	Gateway string              `bson:"gateway" json:"gateway"`
	Leaf    bool                `bson:"leaf" json:"leaf"`
	Health  bool                `bson:"health" json:"health"`
	ID      string              `bson:"id" json:"id"`
	Name    string              `bson:"name" json:"name"`
	States  map[string][]*State `bson:"states" json:"states"`
	Calls   []string            `bson:"calls" json:"calls"`
	Entry   bool                `bson:"entry" json:"entry"`
	// Traffic int64    `bson:"traffic" json:"traffic"`
}

// 微服务状态
type State struct {
	sync.RWMutex
	Id        *StateId       `bson:"id" json:"id"`
	Ecc       float64        `bson:"ecc" json:"ecc"`
	Threshold float64        `bson:"threshold" json:"threshold"`
	Variance  *Vector        `bson:"variance" json:"variance"`
	Sigma     float64        `bson:"sigma" json:"sigma"`
	K         int64          `bson:"k" json:"k"`
	MaxTime   float64        `bson:"maxTime" json:"maxTime"`
	MinTime   float64        `bson:"minTime" json:"minTime"`
	Record    []*StateRecord `bson:"record" json:"record"`
}

type StateRecord struct {
	Ecc       float64 `bson:"ecc" json:"ecc"`
	Threshold float64 `bson:"threshold" json:"threshold"`
}

type StateId struct {
	StartWith *StateEndpointEvent `bson:"startWith" json:"startWith"`
	EndWith   *StateEndpointEvent `bson:"endwith" json:"endwith"`
}

type StateEndpointEvent struct {
	IP       string `bson:"ip" json:"ip"`
	HttpType string `bson:"httpType" json:"httpType"`
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
	URL       string
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
	Leaf     bool
	Entry    bool
}
