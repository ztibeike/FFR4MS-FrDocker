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
	var containers []*types.Container
	if strings.HasPrefix(confPath, "http") {
		containers = utils.GetConfigFromEureka(confPath)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(errors.New("do not support file-config yet"))
	}
	utils.GetServiceContainers(containers)
	handler, err = pcap.OpenLive(ifaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalln(err)
	}
	defer handler.Close()
	if err = handler.SetBPFFilter(filter); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Start capturing packets on interface: %s\n", ifaceName)
	packetSource := gopacket.NewPacketSource(handler, handler.LinkType())
	packets := packetSource.Packets()

	// var IPChanMap = make(map[string]chan *types.HttpInfo)

	for packet := range packets {
		// fmt.Println(packet)
		if packet == nil || packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			continue
		}
		tcp := packet.TransportLayer().(*layers.TCP)
		srcIP := packet.NetworkLayer().NetworkFlow().Src().String()
		dstIP := packet.NetworkLayer().NetworkFlow().Dst().String()
		if !(constants.IPAllMSMap.Has(srcIP) && constants.IPAllMSMap.Has(dstIP)) {
			continue
		}
		if len(tcp.Payload) < 16 {
			continue
		}
		var httpInfo *types.HttpInfo
		if httpInfo, err = utils.GetHttpInfo(packet, tcp); err != nil {
			log.Println(err.Error())
			continue
		}
		var currentIP string // 当前http应该检测的服务IP
		if constants.IPServiceContainerMap.Has(srcIP) {
			currentIP = srcIP
		} else {
			currentIP = dstIP
		}
		obj, _ := constants.IPServiceContainerMap.Get(currentIP)
		var currentContainer = obj.(*types.Container)
		if !currentContainer.Health {
			continue
		}
		var httpChan chan *types.HttpInfo
		var ok bool
		if httpChan, ok = constants.IPChanMap[currentIP]; ok {
			httpChan <- httpInfo
		} else {
			httpChan = make(chan *types.HttpInfo)
			go utils.StateMonitor(currentIP, httpChan)
			httpChan <- httpInfo
			constants.IPChanMap[currentIP] = httpChan
		}
		// fmt.Printf("%s -> %s\n", srcIP, dstIP)
		// fmt.Println(string(tcp.Payload))
		// httpType, _ := utils.GetHttpType(tcp.Payload)
		// fmt.Printf("Http type: %s\n", httpType)
		// traceId := utils.GetTraceId(tcp.Payload)
		// fmt.Printf("TraceId: %s\n", traceId)
		// fmt.Println("----------------------------------")
	}
}

func main() {
	var ifaceName = "br-46facbce86c7"
	var confPath = "http://localhost:8030/getConf"
	runRecovery(ifaceName, confPath)
}
