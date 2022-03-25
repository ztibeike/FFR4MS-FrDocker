package dto

type UpContainerDTO struct {
	ServiceGroup string `json:"serviceGroup" binding:"required"`
	ServiceIP    string `json:"serviceIP" binding:"required"`
	ServicePort  string `json:"servicePort" binding:"required"`
}
