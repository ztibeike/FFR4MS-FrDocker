package frecovery

import (
	"context"
	"strings"
	"sync"

	"gitee.com/zengtao321/frdocker/commons"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils"

	"gitee.com/zengtao321/frdocker/utils/logger"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func RunFrecovery(ifaceName string, confPath string) {
	logger.Info(nil, "Fr-Docker Started!\n")
	defer logger.Info(nil, "Fr-Docker Stopped!\n")
	commons.Network = ifaceName
	commons.RegistryURL = confPath
	var err error
	InitContainers(ifaceName, confPath)
	var wg sync.WaitGroup
	wg.Add(1)
	go SetupCloseHandler(ifaceName, &wg)
	defer wg.Wait()
	pcapHandler, err = pcap.OpenLive(ifaceName, 1600, true, pcap.BlockForever)
	var filter = "tcp"
	if err != nil {
		logger.Fatalln(nil, err)
	}
	if err = pcapHandler.SetBPFFilter(filter); err != nil {
		logger.Fatalln(nil, err)
	}
	logger.Info(nil, "Start Capturing Packets on Interface: %s\n", ifaceName)
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
		// 判断入口微服务组: src非网关和服务 && dst是网关
		if !commons.IPAllMSMap.Has(srcIP) && commons.IPAllMSMap.Has(dstIP) && !commons.IPServiceContainerMap.Has(dstIP) {
			httpInfo, err = utils.GetHttpInfo(packet, tcp)
			if err != nil {
				continue
			}
			// 外界来的请求，类型是REQUEST，且没有TraceId
			if httpInfo.Type == "REQUEST" && httpInfo.TraceId == "" {
				go func() {
					obj, _ := commons.IPAllMSMap.Get(dstIP)
					msType := obj.(string)
					colon := strings.Index(msType, ":")
					obj, _ = commons.ServiceGroupMap.Get(msType[colon+1:])
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
						obj, _ := commons.IPServiceContainerMap.Get(IP)
						container := obj.(*types.Container)
						container.Entry = true
					}
				}()
			}
		}

		if !(commons.IPAllMSMap.Has(srcIP) && commons.IPAllMSMap.Has(dstIP)) {
			continue
		}

		if httpInfo == nil {
			httpInfo, err = utils.GetHttpInfo(packet, tcp)
			if err != nil {
				logger.Errorln(nil, err)
				continue
			}
		}

		var currentIP string // 当前http应该检测的服务IP
		if commons.IPServiceContainerMap.Has(srcIP) {
			currentIP = srcIP
		} else {
			currentIP = dstIP
		}
		if httpInfo.DstIP == currentIP && httpInfo.Type == "REQUEST" {
			trafficChan <- currentIP
		}
		obj, _ := commons.IPServiceContainerMap.Get(currentIP)
		var currentContainer = obj.(*types.Container)
		commons.IPChanMapMutex.Lock()
		if !currentContainer.Health {
			commons.IPChanMapMutex.Unlock()
			continue
		}
		var httpChan chan *types.HttpInfo
		var ok bool
		if httpChan, ok = commons.IPChanMap[currentIP]; ok {
			httpChan <- httpInfo
		} else {
			httpChan = make(chan *types.HttpInfo)
			go StateMonitor(currentIP, httpChan)
			httpChan <- httpInfo
			commons.IPChanMap[currentIP] = httpChan
		}
		commons.IPChanMapMutex.Unlock()
	}
	logger.Info(nil, "Closing All Channels......\n")
	cancel()
	for IP, ch := range commons.IPChanMap {
		close(ch)
		delete(commons.IPChanMap, IP)
	}
	close(trafficChan)
}
