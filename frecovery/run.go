package frecovery

import (
	"gitee.com/zengtao321/frdocker/db"
	"github.com/robfig/cron/v3"
)

func (app *FrecoveryApp) Run() {
	app.Logger.Info("start frdocker...")
	app.initMSSystem()
	app.initPcap()
	persistecheduler := app.persist()
	metricScheduler := app.monitorMetric()
	app.monitorState() // 阻塞
	app.Close(persistecheduler, metricScheduler)
	app.Logger.Info("stop frdocker...")
}

func (app *FrecoveryApp) Close(persistScheduler, metricScheduler *cron.Cron) {
	persistScheduler.Stop()
	// 退出时保存状态
	app.persistenceTask()
	metricScheduler.Stop()
	db.CloseMongo(app.DbCli)
	app.PcapHandle.Close()
	app.Logger.Writer().Close()
}
