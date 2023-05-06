package frecovery

import (
	"gitee.com/zengtao321/frdocker/db"
)

func (app *FrecoveryApp) Run() {
	app.Logger.Info("start frdocker...")
	app.initMSSystem()
	app.initPcap()
	app.monitorMetric()
	app.monitorMessage() // 阻塞
	app.Close()
	app.Logger.Info("stop frdocker...")
}

func (app *FrecoveryApp) Close() {
	db.CloseMongo(app.DbCli)
	app.PcapHandle.Close()
	app.Logger.Writer().Close()
}
