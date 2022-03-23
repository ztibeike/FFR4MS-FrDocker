package frecovery

import (
	"frdocker/constants"
	"frdocker/types"
	"frdocker/utils"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func RunFrecovery(ifaceName string, confPath string) {
	log.Println("Fr-Docker Started!")
	var err error
	InitContainers(ifaceName, confPath)
	go SetupCloseHandler(ifaceName)
	handler, err = pcap.OpenLive(ifaceName, 1600, true, pcap.BlockForever)
	var filter = "tcp"
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		log.Println("Fr-Docker Stopped!")
	}()
	if err = handler.SetBPFFilter(filter); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Start Capturing Packets on Interface: %s\n", ifaceName)
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
		constants.IPChanMapMutex.Lock()
		if !currentContainer.Health {
			constants.IPChanMapMutex.Unlock()
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
		constants.IPChanMapMutex.Unlock()
	}
}
