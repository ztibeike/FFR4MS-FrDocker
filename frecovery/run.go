package frecovery

import (
	"os"
	"os/signal"
	"syscall"

	"gitee.com/zengtao321/frdocker/db"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (app *FrecoveryApp) Run() {
	app.Logger.Info("start frdocker...")
	app.initMSSystem()
	app.initPcap()
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
			goto close
		}
	}
close:
	app.Close()
	app.Logger.Info("stop frdocker...")
}

func (app *FrecoveryApp) handlePacket(packet gopacket.Packet) {

}

// 检查packet是否合理
func (app *FrecoveryApp) checkValid(packet gopacket.Packet) bool {
	// 检查packet有效性
	if packet == nil || packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
		return false
	}
	tcp, ok := packet.TransportLayer().(*layers.TCP)
	if !ok || len(tcp.Payload) < 16 {
		return false
	}
	// 检查packet是否是微服务系统内部通信消息, 微服务间消息通信路径为: service->gateway->service
	return true
}

func (app *FrecoveryApp) Close() {
	db.CloseMongo(app.DbCli)
	app.PcapHandle.Close()
	app.Logger.Writer().Close()
}
