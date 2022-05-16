package commons

import (
	"sync"

	"gitee.com/zengtao321/frdocker/types"

	cmap "github.com/orcaman/concurrent-map"
)

var Network string

var RegistryURL string

// IPServiceContainerMap 功能微服务列表, 服务IP到服务容器详情的映射
var IPServiceContainerMap = cmap.New()

// IPAllMSMap 功能微服务+网关列表, IP到类型的映射，类型有SERVICE:group和GATEWAY:group
var IPAllMSMap = cmap.New()

// ServiceGroupMap 服务实例组，名称到实例组（包括网关和服务）的映射, 其中网关是IP:Port形式，服务是IP形式
var ServiceGroupMap = cmap.New()

// IPChanMap 服务IP与服务监控协程通道的映射
var IPChanMap = make(map[string]chan *types.HttpInfo)
var IPChanMapMutex sync.Mutex

var DeleteContainerChan = make(chan string, 1)
var AddContainerChan = make(chan string, 1)