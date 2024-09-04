package routers

import (
	"gvb_server/api"
)

func (r RouterGroup) SettingsGroup() {
	settingsApi := api.ApiGroupApp.SettingsApi
	// r.GET("settings", settingsApi.SettingsInfoView)
	// r.PUT("settings", settingsApi.SettingsUpdateView)
	// r.GET("settings_email", settingsApi.SettingsEmailView)
	// r.PUT("settings_email", settingsApi.SettingsEmailUpdateView)
	// 绑定请求url参数
	r.GET("settings/:name", settingsApi.SettingsInfoView)
	r.PUT("settings/:name", settingsApi.SettingsUpdateView)
}
