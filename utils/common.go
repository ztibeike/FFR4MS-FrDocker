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
	if strings.Compare(strings.ToLower(method), "get") == 0 {
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
	// response, err := http.Get(confPath)
	// if err != nil {
	// 	log.Fatalln("bad eureka url")
	// }
	// body, _ := ioutil.ReadAll(response.Body)
	// response.Body.Close()
	// result := struct {
	// 	ArrayIpPort []string
	// 	ArrayGetWay []string
	// 	ArrayGroup  []string
	// }{}
	// err = json.Unmarshal(body, &result)
	// if err != nil {
	// 	log.Fatalln("bad eureka url")
	// }
	resp := struct {
		ArrayIpPort []string
		ArrayGetWay []string
		ArrayGroup  []string
	}{}
	HttpRequest(confPath, "GET", &resp)
	var gatewayMap = make(map[string]string)
	for _, gateway := range resp.ArrayGetWay {
		url := "http://" + gateway + "/actuator/info"
		gatewayInfo := struct {
			Getway string
			Port   string
		}{}
		HttpRequest(url, "GET", &gatewayInfo)
		gatewayMap[gatewayInfo.Getway] = gateway
	}
	var containers []types.Container
	for idx, service := range resp.ArrayIpPort {
		serviceInfo := struct {
			Leaf int
			Port string
		}{}
		serviceInfoURL := "http://" + service + "/actuator/info"
		HttpRequest(serviceInfoURL, "GET", &serviceInfo)
		serviceHealth := struct {
			Status string
		}{}
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
		constants.IPContainerMap.Set(container.IP, container)
	}
	return containers
}
