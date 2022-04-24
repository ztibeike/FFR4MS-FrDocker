package dto

type AddContainerDTO struct {
	ServiceGroup string `json:"serviceGroup" binding:"required"`
	ServiceIP    string `json:"serviceIP" binding:"required"`
	ServicePort  string `json:"servicePort" binding:"required"`
}

type UserInfoDTO struct {
	Permissions []string `json:"permissions"`
	Username    string   `json:"username"`
	Avatar      string   `json:"avatar"`
}
