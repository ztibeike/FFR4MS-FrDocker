package frecovery

import (
	"gitee.com/zengtao321/frdocker/db"
	"github.com/robfig/cron/v3"
)

func (app *FrecoveryApp) Run() {
	app.Logger.Info("start frdocker...")
	app.initMSSystem()
	app.initPcap()
	persistecheduler := app.Persist()
	metricScheduler := app.monitorMetric()
	app.monitorState() // 阻塞
	app.Close(persistecheduler, metricScheduler)
	app.Logger.Info("stop frdocker...")
}

func (app *FrecoveryApp) Close(persistScheduler, metricScheduler *cron.Cron) {
	persistScheduler.Stop()
	metricScheduler.Stop()
	db.CloseMongo(app.DbCli)
	app.PcapHandle.Close()
	app.Logger.Writer().Close()
}
