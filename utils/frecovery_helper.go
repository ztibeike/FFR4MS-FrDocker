package utils

import (
	"fmt"
	"net/http"

	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/types/dto"
	"github.com/go-resty/resty/v2"
)

// 获取注册中心信息
func GetRegistryInfo(addr string) (dto.MSConfig, error) {
	client := resty.New()
	registryConfig := dto.MSConfig{}
	url := fmt.Sprintf("http://%s%s", addr, config.REGISTRY_INFO_URI)
	_, err := client.R().SetHeader("Accept", "application/json").SetResult(&registryConfig).Get(url)
	return registryConfig, err
}

// 容器健康检查
func CheckContainerHealth(ip string, port int) (bool, error) {
	client := resty.New()
	url := fmt.Sprintf("http://%s:%d%s", ip, port, config.CONTAINER_HEALTH_CHECK_URI)
	var health dto.ContainerHealth
	resp, err := client.SetTimeout(config.CONTAINER_HEALTH_CHECK_TIMEOUT).R().SetResult(&health).Get(url)
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != 200 {
		return false, fmt.Errorf("check container health failed, status code: %d", resp.StatusCode())
	}
	return health.Status == "UP", nil
}

// 通知网关进行消息重播
func NotifyGatewayReplayMessage(gatewayAddr, serviceName, containerIP string, containerPort int) error {
	result := dto.CommonResponse[any]{}
	url := fmt.Sprintf("http://%s%s", gatewayAddr, config.GATEWAY_REPLAY_MESSAGE_URI)
	data := dto.ReplayMessageDTO{
		ServiceName:      serviceName,
		DownInstanceHost: containerIP,
		DownInstancePort: containerPort,
	}
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(data).SetResult(&result).Post(url)
	if err != nil {
		return err
	}
	if result.Code != http.StatusOK || resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("error result code %d or status code: %d", result.Code, resp.StatusCode())
	}
	return nil
}
