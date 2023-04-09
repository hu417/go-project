package settings_api

import (
	"gvb_server/config"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models/res"

	"github.com/gin-gonic/gin"
)

// 更新系统配置信息接口
func (SettingsApi) SettingsUpdateView(c *gin.Context) {

	//
	var cf config.SiteInfo
	err := c.ShouldBindJSON(&cf)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
	}

	global.Config.SiteInfo = cf

	err = core.WriteConf()
	if err != nil {
		global.Logger.Error(err)
		res.FailWithMessage(err.Error(), c)
		return
	}
	res.OkWith(c)
}
