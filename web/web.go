package web

import (
	"frdocker/settings"
	"frdocker/web/router"
)

func SetupWebHander() {
	r := router.SetupRouter()
	r.Run(":" + settings.WEB_HANDLER_PORT)
}
