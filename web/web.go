package web

import (
	"gitee.com/zengtao321/frdocker/settings"
	"gitee.com/zengtao321/frdocker/web/router"
)

func SetupWebHander() {
	r := router.SetupRouter()
	r.Run(":" + settings.WEB_HANDLER_PORT)
}
