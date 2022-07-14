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

type GateWayUpService struct {
	ServiceName    string
	UpInstanceHost string
	UpInstancePort string
}

type RecoveryMessageDTO struct {
	TimeStamp       string `json:"timeStamp"`
	ServiceName     string `json:"serviceName"`
	TraceId         string `json:"traceId"`
	OldInstanceHost string `json:"oldInstanceHost"`
	OldInstancePort int    `json:"oldInstancePort"`
	NewInstanceHost string `json:"newInstanceHost"`
	NewInstancePort int    `json:"newInstancePort"`
}
