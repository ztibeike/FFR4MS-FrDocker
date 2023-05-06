package frecovery

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"gitee.com/zengtao321/frdocker/frecovery/entity"
	"gitee.com/zengtao321/frdocker/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (app *FrecoveryApp) monitorMessage() {
	packetSource := gopacket.NewPacketSource(app.PcapHandle, app.PcapHandle.LinkType())
	packets := packetSource.Packets()
	// 监听程序中断信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGABRT, syscall.SIGQUIT)
	app.Logger.Infof("Start Capturing Packets on Interface: %s", app.NetworkInterface)
	for {
		select {
		case packet := <-packets:
			app.handlePacket(packet)
		case <-signalChan:
			return
		}
	}
}

func (app *FrecoveryApp) handlePacket(packet gopacket.Packet) {
	if !app.checkPacketValid(packet) {
		return
	}
	httpInfo, err := entity.NewHttpInfo(packet)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}
	if err = app.setHttpRole(httpInfo); err != nil {
		app.Logger.Error(err.Error())
		return
	}
	app.Logger.Infof("[%s][%s]%s:%d -> [%s][%s]%s:%d", httpInfo.Src.Type, httpInfo.Src.Name, httpInfo.Src.IP, httpInfo.Src.Port,
		httpInfo.Dst.Type, httpInfo.Dst.Name, httpInfo.Dst.IP, httpInfo.Dst.Port)
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
	srcType := app.getContainerType(utils.GenerateIdFromIPAndPort(srcIP, srcPort))
	dstType := app.getContainerType(utils.GenerateIdFromIPAndPort(dstIP, dstPort))
	if srcType == entity.CTN_INVALID || dstType == entity.CTN_INVALID {
		return false
	}
	return true
}

func (app *FrecoveryApp) setHttpRole(httpInfo *entity.HttpInfo) error {
	srcId := utils.GenerateIdFromIPAndPort(httpInfo.Src.IP, httpInfo.Src.Port)
	dstId := utils.GenerateIdFromIPAndPort(httpInfo.Dst.IP, httpInfo.Dst.Port)
	srcType := app.getContainerType(srcId)
	dstType := app.getContainerType(dstId)
	if srcType == entity.CTN_INVALID || dstType == entity.CTN_INVALID {
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
