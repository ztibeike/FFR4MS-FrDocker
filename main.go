package main

import (
	"frdocker/frecovery"
	"frdocker/web"
)

func main() {
	var ifaceName = "br-46facbce86c7"
	var confPath = "http://localhost:8030/getConf"

	go web.SetupWebHander()
	frecovery.RunFrecovery(ifaceName, confPath)
}
