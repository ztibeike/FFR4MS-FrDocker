package main

import "frdocker/frecovery"

func main() {
	var ifaceName = "br-46facbce86c7"
	var confPath = "http://localhost:8030/getConf"
	frecovery.RunFrecovery(ifaceName, confPath)
}
