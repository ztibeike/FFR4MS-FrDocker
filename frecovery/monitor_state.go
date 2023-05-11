package frecovery

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"gitee.com/zengtao321/frdocker/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (app *FrecoveryApp) monitorState() {
	app.Logger.Info("start state monitoring...")
	packetSource := gopacket.NewPacketSource(app.PcapHandle, app.PcapHandle.LinkType())
	packets := packetSource.Packets()
	// 监听程序中断信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGABRT, syscall.SIGQUIT)
	for {
		select {
		case packet := <-packets:
			app.handlePacket(packet)
		case <-signalChan:
			return
		}
	}
}

type MonitorStateCallBack func(traceId string, state *ContainerState)

func (app *FrecoveryApp) handlePacket(packet gopacket.Packet) {
	if !app.checkPacketValid(packet) {
		return
	}
	httpInfo, err := NewHttpInfo(packet)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}
	if err = app.setHttpRole(httpInfo); err != nil {
		app.Logger.Error(err.Error())
		return
	}
	container := app.getCheckingContainer(httpInfo)
	if container == nil || !container.IsHealthy {
		return
	}
	container.Monitor.UpdateContainerState(httpInfo, app.monitorStateCallback)
}

func (app *FrecoveryApp) monitorStateCallback(traceId string, state *ContainerState) {
	container := app.GetContainer(state.Id)
	state.mu.RLock()
	defer state.mu.RUnlock()
	ecc, thresh := state.Ecc, state.Thresh
	if ecc > thresh {
		app.Logger.Errorf("[state][%s][%s:%d][%s] ecc: %.4f, thresh: %.4f", container.ServiceName, container.IP, container.Port, traceId, ecc, thresh)
		go app.handleContainerAbnormal(container)
	} else {
		app.Logger.Tracef("[state][%s][%s:%d][%s] ecc: %.4f, thresh: %.4f", container.ServiceName, container.IP, container.Port, traceId, ecc, thresh)
	}
}

func (app *FrecoveryApp) handleContainerAbnormal(container *Container) {
	healthy, _ := utils.CheckContainerHealth(container.IP, container.Port)
	if healthy {
		return
	}
	// 并发控制
	container.mu.Lock()
	defer container.mu.Unlock()
	if !container.IsHealthy {
		return // 避免重复通知
	}
	// 获取容器所属网关
	service := app.GetService(container.ServiceName)
	gateway := app.GetGateway(service.Gateway)
	// TODO 通知网关重播
	for _, ctn := range gateway.Containers {
		err := utils.NotifyGatewayReplayMessage(ctn, container.ServiceName, container.IP, container.Port)
		if err != nil {
			app.Logger.Errorf("[%s][%s:%d] notify gateway %s replay message failed: %s", container.ServiceName, container.IP, container.Port, ctn, err.Error())
		} else {
			container.IsHealthy = false
			app.Logger.Infof("[%s][%s:%d] notify gateway %s replay message success", container.ServiceName, container.IP, container.Port, ctn)
		}
	}
}

// 检查packet是否合理
func (app *FrecoveryApp) checkPacketValid(packet gopacket.Packet) bool {
	// 检查packet有效性
	if packet == nil || packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
		return false
	}
	tcp, ok := packet.TransportLayer().(*layers.TCP)
	if !ok || len(tcp.Payload) < 16 {
		return false
	}
	// 检查packet是否是微服务系统内部通信消息, 微服务间消息通信路径为: service->gateway->service
	srcIP, dstIP, srcPort, dstPort := utils.GetIPAndPort(packet)
	srcType := app.getContainerType(utils.GenerateContainerId(srcIP, srcPort))
	dstType := app.getContainerType(utils.GenerateContainerId(dstIP, dstPort))
	if srcType == CTN_INVALID || dstType == CTN_INVALID {
		return false
	}
	return true
}

func (app *FrecoveryApp) setHttpRole(httpInfo *HttpInfo) error {
	srcId := utils.GenerateContainerId(httpInfo.Src.IP, httpInfo.Src.Port)
	dstId := utils.GenerateContainerId(httpInfo.Dst.IP, httpInfo.Dst.Port)
	srcType := app.getContainerType(srcId)
	dstType := app.getContainerType(dstId)
	if srcType == CTN_INVALID || dstType == CTN_INVALID {
		return errors.New("invalid httpInfo")
	}
	httpInfo.Src.Id = srcId
	httpInfo.Dst.Id = dstId
	httpInfo.Src.Type = srcType
	httpInfo.Dst.Type = dstType
	httpInfo.Src.Name = app.Containers[srcId].ServiceName
	httpInfo.Dst.Name = app.Containers[dstId].ServiceName
	return nil
}

func (app *FrecoveryApp) getCheckingContainer(httpInfo *HttpInfo) *Container {
	var currentId string
	if httpInfo.Src.Type == CTN_SERVICE {
		currentId = httpInfo.Src.Id
	} else {
		currentId = httpInfo.Dst.Id
	}
	return app.GetContainer(currentId)
}
