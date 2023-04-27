package frecovery

func (app *FrecoveryApp) Run() {
	app.Logger.Infoln("start frecovery...")
	app.initApp()
	app.Logger.Infoln("stop frecovery...")
}
