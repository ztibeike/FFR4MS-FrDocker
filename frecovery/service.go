package frecovery

type Service struct {
	ServiceName string   `json:"serviceName" bson:"serviceName"` // 服务名称
	Group       string   `json:"group" bson:"group"`             // 服务组
	Gateway     string   `json:"gateway" bson:"gateway"`         // 网关
	IsLeaf      bool     `json:"isLeaf" bson:"isLeaf"`           // 是否叶子节点
	IsRoot      bool     `json:"isRoot" bson:"isRoot"`           // 是否根节点
	Calls       []string `json:"calls" bson:"calls"`             // 调用的服务
	Containers  []string `json:"containers" bson:"containers"`   // 服务实例容器
}

func NewService(serviceName string) *Service {
	return &Service{
		ServiceName: serviceName,
		Group:       serviceName,
		Gateway:     "",
		IsLeaf:      false,
		IsRoot:      false,
	}
}
