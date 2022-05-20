package main

import "gitee.com/zengtao321/frdocker/cmd"

// @title FFR4MS Docker监控模块Fr-Docker
// @version 1.0
// @description 监控微服务系统中各个微服务实例的运行时状态以及所在容器的性能参数
// @termsOfService http://swagger.io/terms/
// @contact.name 曾涛
// @contact.email zengtao0618@163.com
func main() {
	// var ifaceName = "br-46facbce86c7"
	// var confPath = "http://localhost:8030/getConf"

	// go web.SetupWebHander()
	// frecovery.RunFrecovery(ifaceName, confPath)
	cmd.Execute()
}
