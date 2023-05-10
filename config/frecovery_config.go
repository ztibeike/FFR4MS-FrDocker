package config

import "time"

const (
	REGISTRY_METADATA_LEAF_KEY    = "leaf"
	REGISTRY_METADATA_GATEWAY_KEY = "gateway"
	REGISTRY_INFO_URI             = "/frecovery/conf"

	GATEWAY_REPLAY_MESSAGE_URI = "/frecovery/replace"

	CONTAINER_HEALTH_CHECK_URI     = "/actuator/health"
	CONTAINER_HEALTH_CHECK_TIMEOUT = 300 * time.Millisecond
)
