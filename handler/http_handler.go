package handler

import (
	"encoding/json"
	"frdocker/constants"
	"frdocker/settings"
	"frdocker/types"
	"frdocker/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

type Resp struct {
	Code    int
	Message string
}

func AddContainerHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) != "POST" {
		resp, _ := json.Marshal(Resp{
			Code:    http.StatusMethodNotAllowed,
			Message: "Request Method Not Allowed!",
		})
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(resp)
		return
	}
	var addContainerDTO types.AddContainerDTO
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &addContainerDTO)
	if err != nil || addContainerDTO.ServiceGroup == "" || addContainerDTO.ServiceIP == "" || addContainerDTO.ServicePort == "" {
		resp, _ := json.Marshal(Resp{
			Code:    http.StatusBadRequest,
			Message: "Wrong Request Params!",
		})
		w.Write(resp)
	}
	container := &types.Container{
		IP:     addContainerDTO.ServiceIP,
		Port:   addContainerDTO.ServicePort,
		Group:  addContainerDTO.ServiceGroup,
		Health: true,
	}
	obj, _ := constants.ServiceGroupMap.Get(container.Group)
	serviceGroup := obj.(*types.ServiceGroup)
	serviceGroup.Services = append(serviceGroup.Services, container.IP)
	container.Gateway = serviceGroup.Gateway
	constants.IPAllMSMap.Set(container.IP, "SERVICE:"+container.Group)
	utils.GetServiceContainers([]*types.Container{container})
	constants.IPServiceContainerMap.Set(container.IP, container)
	resp, _ := json.Marshal(Resp{
		Code:    http.StatusOK,
		Message: "Success!",
	})
	w.Write(resp)
}

func HttpHandler() {
	http.HandleFunc("/add", AddContainerHandler)
	http.ListenAndServe(":"+settings.HTTP_HANDLER_PORT, nil)
}
