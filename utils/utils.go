package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils/logger"
)

func HttpRequest(url string, method string, result interface{}) {
	var response *http.Response
	var err error
	if strings.Compare(strings.ToUpper(method), "GET") == 0 {
		response, err = http.Get(url)
	} else {
		response, err = http.Post(url, "application/json", nil)
	}
	if err != nil {
		logger.Fatalln(nil, "Bad Eureka URL: ", url)
	}
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	err = json.Unmarshal(body, result)
	if err != nil {
		logger.Fatalln(nil, "Bad Eureka Response")
	}
}

// 从Eureka注册中心获取配置信息
func GetConfigFromEureka(confPath string) []*types.Container {
	resp := types.EurekaConfig{}
	HttpRequest(confPath, "GET", &resp)
	for _, gateway := range resp.ArrayGetWay {
		url := "http://" + gateway + "/actuator/info"
		gatewayInfo := types.GatewayActuatorInfo{}
		HttpRequest(url, "GET", &gatewayInfo)
		colon := strings.Index(gateway, ":")
		commons.IPAllMSMap.Set(gateway[:colon], "GATEWAY:"+gatewayInfo.Getway)
		serviceGroup := &types.ServiceGroup{
			Gateway: gateway,
			Entry:   false,
		}
		commons.ServiceGroupMap.Set(gatewayInfo.Getway, serviceGroup)
	}
	var containers []*types.Container
	var obj interface{}
	for idx, service := range resp.ArrayIpPort {
		serviceInfo := types.ServiceActuatorInfo{}
		serviceInfoURL := "http://" + service + "/actuator/info"
		HttpRequest(serviceInfoURL, "GET", &serviceInfo)
		serviceHealth := types.ServiceActuatorHealth{}
		serviceHealthURL := "http://" + service + "/actuator/health"
		HttpRequest(serviceHealthURL, "GET", &serviceHealth)
		colon := strings.Index(service, ":")
		obj, _ = commons.ServiceGroupMap.Get(resp.ArrayGroup[idx])
		serviceGroup := obj.(*types.ServiceGroup)
		serviceGroup.Services = append(serviceGroup.Services, service[:colon])
		serviceGroup.Leaf = (serviceInfo.Leaf == 1)
		container := &types.Container{
			IP:      service[:colon],
			Port:    service[colon+1:],
			Group:   resp.ArrayGroup[idx],
			Gateway: serviceGroup.Gateway,
			Leaf:    serviceInfo.Leaf == 1,
			Health:  strings.ToUpper(serviceHealth.Status) == "UP",
		}
		containers = append(containers, container)
		commons.IPServiceContainerMap.Set(container.IP, container)
		commons.IPAllMSMap.Set(service[:colon], "SERVICE:"+container.Group)
	}
	return containers
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
