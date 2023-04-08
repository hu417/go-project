package routers

import (
	"gvb_server/api"
)

func (r RouterGroup) SettingsGroup() {
	settingsApi := api.ApiGroupApp.SettingsApi
	r.GET("settings", settingsApi.SettingsInfoView)
}
