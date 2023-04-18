package entity

type Service struct {
	ServiceName string   // 服务名称
	Group       string   // 服务组
	Gateway     string   // 网关
	IsLeaf      bool     // 是否叶子节点
	IsRoot      bool     // 是否根节点
	Calls       []string // 调用的服务
	Containers  []string // 服务实例容器
}
