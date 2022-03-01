package main

import (
	"errors"
	"fmt"
	"frdocker/types"
	"frdocker/utils"
	"log"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func runRecovery(ifaceName string, confPath string) {
	fmt.Println("Fr-Docker Started!")
	var err error
	var handler *pcap.Handle
	var filter = "tcp"
	var containers []types.Container
	if strings.HasPrefix(confPath, "http") {
		containers = utils.GetConfigFromEureka(confPath)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(errors.New("do not support file-config yet"))
	}
	utils.GetServiceContainers(containers)
	fmt.Println(containers)
	handler, err = pcap.OpenLive(ifaceName, 0, true, pcap.BlockForever)
	if err != nil {
		log.Fatalln(err)
	}
	if err = handler.SetBPFFilter(filter); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Start capturing packets on interface: %s\n", ifaceName)
	packetSource := gopacket.NewPacketSource(handler, handler.LinkType())
	packets := packetSource.Packets()

	for {
		select {
		case packet := <-packets:
			{
				fmt.Println(packet)
			}
		}
	}

}

func main() {
	var ifaceName = "br-d3bce3e90ea6"
	var confPath = "http://localhost:8030/getConf"
	runRecovery(ifaceName, confPath)
}
