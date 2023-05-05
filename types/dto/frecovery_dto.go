package dto

// Fr-Eureka注册中心dto
type MSInstance struct {
	Name     string            `json:"name"`
	IP       string            `json:"ip"`
	Port     int               `json:"port"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

// Fr-Eureka注册中心dto
type MSConfig struct {
	Services map[string][]MSInstance `json:"services"`
	Gateways map[string][]MSInstance `json:"gateways"`
	Groups   []string                `json:"groups"`
}
