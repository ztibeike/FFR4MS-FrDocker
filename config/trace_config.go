package config

// 链路追踪配置
const (

	// 链路追踪唯一id
	TRACE_ID_HEADER = "trace-id"

	// 实例id, 格式为: {serviceName}-ip-port
	TRACE_SERVICE_INSTANCE_HEADER = "service-instance-id"

	// 实例所属服务名称
	TRACE_SERVICE_NAME_HEADER = "service-name"

	// 标识该响应被网关拦截和从缓存中返回
	TRACE_CACHED_RESPONSE_HEADER = "cached-response"
)
