package frecovery

import (
	"context"
	"frdocker/constants"
	"frdocker/types"
	"frdocker/utils"
	"strings"
	"sync"

	"frdocker/utils/logger"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func RunFrecovery(ifaceName string, confPath string) {
	logger.Info("Fr-Docker Started!\n")
	defer logger.Info("Fr-Docker Stopped!\n")
	constants.Network = ifaceName
	constants.RegistryURL = confPath
	var err error
	InitContainers(ifaceName, confPath)
	var wg sync.WaitGroup
	wg.Add(1)
	go SetupCloseHandler(ifaceName, &wg)
	defer wg.Wait()
	pcapHandler, err = pcap.OpenLive(ifaceName, 1600, true, pcap.BlockForever)
	var filter = "tcp"
	if err != nil {
		logger.Fatalln(err)
	}
	if err = pcapHandler.SetBPFFilter(filter); err != nil {
		logger.Fatalln(err)
	}
	logger.Info("Start Capturing Packets on Interface: %s\n", ifaceName)
	packetSource := gopacket.NewPacketSource(pcapHandler, pcapHandler.LinkType())
	packets := packetSource.Packets()

	trafficChan := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	go CronSaveTraffic(ctx, trafficChan)
	go CronSaveContainerInfo(ctx, ifaceName)
	// var IPChanMap = make(map[string]chan *types.HttpInfo)

	for packet := range packets {
		// fmt.Println(packet)
		if packet == nil || packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			continue
		}
		tcp := packet.TransportLayer().(*layers.TCP)
		srcIP := packet.NetworkLayer().NetworkFlow().Src().String()
		dstIP := packet.NetworkLayer().NetworkFlow().Dst().String()

		if len(tcp.Payload) < 16 {
			continue
		}

		var httpInfo *types.HttpInfo
		// 判断入口微服务组
		if !constants.IPAllMSMap.Has(srcIP) && constants.IPAllMSMap.Has(dstIP) && !constants.IPServiceContainerMap.Has(dstIP) {
			httpInfo, err = utils.GetHttpInfo(packet, tcp)
			if err != nil {
				continue
			}
			if httpInfo.Type == "REQUEST" {
				go func() {
					obj, _ := constants.IPAllMSMap.Get(dstIP)
					msType := obj.(string)
					colon := strings.Index(msType, ":")
					obj, _ = constants.ServiceGroupMap.Get(msType[colon+1:])
					serviceGroup := obj.(*types.ServiceGroup)
					gateway := serviceGroup.Gateway
					colon = strings.Index(gateway, ":")
					if httpInfo.DstPort != gateway[colon+1:] {
						return
					}
					if serviceGroup.Entry {
						return
					}
					serviceGroup.Entry = true
					for _, IP := range serviceGroup.Services {
						obj, _ := constants.IPServiceContainerMap.Get(IP)
						container := obj.(*types.Container)
						container.Entry = true
					}
				}()
			}
		}

		if !(constants.IPAllMSMap.Has(srcIP) && constants.IPAllMSMap.Has(dstIP)) {
			continue
		}

		if httpInfo == nil {
			httpInfo, err = utils.GetHttpInfo(packet, tcp)
			if err != nil {
				logger.Errorln(err)
				continue
			}
		}

		var currentIP string // 当前http应该检测的服务IP
		if constants.IPServiceContainerMap.Has(srcIP) {
			currentIP = srcIP
		} else {
			currentIP = dstIP
		}
		trafficChan <- currentIP
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
			go StateMonitor(currentIP, httpChan)
			httpChan <- httpInfo
			constants.IPChanMap[currentIP] = httpChan
		}
		constants.IPChanMapMutex.Unlock()
	}
	logger.Info("Closing All Channels......\n")
	cancel()
	for IP, ch := range constants.IPChanMap {
		close(ch)
		delete(constants.IPChanMap, IP)
	}
	close(trafficChan)
}
