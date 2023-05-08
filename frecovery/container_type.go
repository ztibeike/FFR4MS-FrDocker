package frecovery

// 实例类型和枚举常量
type ContainerType int

const (
	// 容器类型: 无效
	CTN_INVALID ContainerType = iota

	// 容器类型: 微服务
	CTN_SERVICE

	// 容器类型: 网关
	CTN_GATEWAY
)

func (containerType ContainerType) String() string {
	switch containerType {
	case CTN_SERVICE:
		return "SERVICE"
	case CTN_GATEWAY:
		return "GATEWAY"
	default:
		return "INVALID"
	}
}
