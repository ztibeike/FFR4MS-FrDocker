package frecovery

import (
	"os"
	"os/signal"
	"syscall"

	"gitee.com/zengtao321/frdocker/db"
	"gitee.com/zengtao321/frdocker/frecovery/entity"
	"github.com/google/gopacket"
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
	if !app.checkPacketValid(packet) {
		return
	}
	httpInfo, err := entity.NewHttpInfo(packet)
	if err != nil {
		app.Logger.Errorf("%s: %s:%d->%s:%d", err.Error(), httpInfo.Src.IP, httpInfo.Src.Port, httpInfo.Dst.IP, httpInfo.Dst.Port)
		return
	}
	if err = app.setHttpRole(httpInfo); err != nil {
		app.Logger.Errorf("%s: %s:%d->%s:%d", err.Error(), httpInfo.Src.IP, httpInfo.Src.Port, httpInfo.Dst.IP, httpInfo.Dst.Port)
		return
	}
	app.Logger.Infof("[%s][%s]%s:%d -> [%s][%s]%s:%d", httpInfo.Src.Type, httpInfo.Src.Name, httpInfo.Src.IP, httpInfo.Src.Port,
		httpInfo.Dst.Type, httpInfo.Dst.Name, httpInfo.Dst.IP, httpInfo.Dst.Port)
}

func (app *FrecoveryApp) Close() {
	db.CloseMongo(app.DbCli)
	app.PcapHandle.Close()
	app.Logger.Writer().Close()
}
