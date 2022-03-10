package main

import (
	"errors"
	"fmt"
	"frdocker/constants"
	"frdocker/types"
	"frdocker/utils"
	"log"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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
	handler, err = pcap.OpenLive(ifaceName, 1600, true, pcap.BlockForever)
	defer handler.Close()
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
				// fmt.Println(packet)
				if packet == nil || packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
					continue
				}
				tcp := packet.TransportLayer().(*layers.TCP)
				srcIP := packet.NetworkLayer().NetworkFlow().Src().String()
				dstIP := packet.NetworkLayer().NetworkFlow().Dst().String()
				if srcIP == "172.19.0.2" || dstIP == "172.19.0.2" {
					continue
				}
				if len(tcp.Payload) < 16 {
					continue
				}
				if constants.IPServiceContainerMap.Has(srcIP) || constants.IPServiceContainerMap.Has(dstIP) {
					fmt.Printf("%s -> %s\n", srcIP, dstIP)
					fmt.Println(string(tcp.Payload))
					httpType, _ := utils.GetHttpType(tcp.Payload)
					fmt.Printf("Http type: %s\n", httpType)
					if httpType == "REQUEST" {
						traceId := utils.GetTraceId(tcp.Payload)
						fmt.Printf("TraceId: %s\n", traceId)
					}
					fmt.Println("----------------------------------")
				}
			}
		}
	}

}

func main() {
	var ifaceName = "br-eb67d5a21c9c"
	var confPath = "http://localhost:8030/getConf"
	runRecovery(ifaceName, confPath)
}
