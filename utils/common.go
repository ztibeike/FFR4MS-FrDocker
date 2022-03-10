package utils

import (
	"encoding/json"
	"frdocker/constants"
	"frdocker/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
		log.Fatalln("bad http url: ", url)
	}
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	err = json.Unmarshal(body, result)
	if err != nil {
		log.Fatalln("bad struct")
	}
}

// 从Eureka注册中心获取配置信息
func GetConfigFromEureka(confPath string) []types.Container {
	resp := types.EurekaConfig{}
	HttpRequest(confPath, "GET", &resp)
	var gatewayMap = make(map[string]string)
	for _, gateway := range resp.ArrayGetWay {
		url := "http://" + gateway + "/actuator/info"
		gatewayInfo := types.GatewayActuatorInfo{}
		HttpRequest(url, "GET", &gatewayInfo)
		gatewayMap[gatewayInfo.Getway] = gateway
		colon := strings.Index(gateway, ":")
		constants.IPAllMSMap.Set(gateway[:colon], "GATEWAY")
	}
	var containers []types.Container
	for idx, service := range resp.ArrayIpPort {
		serviceInfo := types.ServiceActuatorInfo{}
		serviceInfoURL := "http://" + service + "/actuator/info"
		HttpRequest(serviceInfoURL, "GET", &serviceInfo)
		serviceHealth := types.ServiceActuatorHealth{}
		serviceHealthURL := "http://" + service + "/actuator/health"
		HttpRequest(serviceHealthURL, "GET", &serviceHealth)
		colon := strings.Index(service, ":")
		container := types.Container{
			IP:      service[:colon],
			Port:    service[colon+1:],
			Group:   resp.ArrayGroup[idx],
			Gateway: gatewayMap[resp.ArrayGroup[idx]],
			Leaf:    serviceInfo.Leaf == 1,
			Health:  strings.ToUpper(serviceHealth.Status) == "UP",
		}
		containers = append(containers, container)
		constants.IPServiceContainerMap.Set(container.IP, container)
		constants.IPAllMSMap.Set(service[:colon], "SERVICE")
	}
	return containers
}
