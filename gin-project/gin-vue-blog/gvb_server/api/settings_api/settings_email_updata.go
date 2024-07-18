package settings_api

import (
	"gvb_server/config"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models/res"

	"github.com/gin-gonic/gin"
)

// 更新email配置信息接口
func (SettingsApi) SettingsEmailUpdateView(c *gin.Context) {

	//
	var cf config.Email
	err := c.ShouldBindJSON(&cf)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
	}

	global.Config.Email = cf

	err = core.WriteConf()
	if err != nil {
		global.Logger.Error(err)
		res.FailWithMessage(err.Error(), c)
		return
	}
	res.OkWith(c)
}
