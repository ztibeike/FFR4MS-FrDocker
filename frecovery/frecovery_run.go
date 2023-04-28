package frecovery

func (app *FrecoveryApp) Run() {
	app.Logger.Info("start frecovery...")
	app.initApp()
	app.Logger.Info("stop frecovery...")
}
