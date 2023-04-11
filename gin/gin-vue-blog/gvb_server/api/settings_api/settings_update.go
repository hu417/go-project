package settings_api

import (
	"gvb_server/config"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models/res"

	"github.com/gin-gonic/gin"
)

// 更新某项配置信息接口
func (SettingsApi) SettingsUpdateView(c *gin.Context) {

	//
	var cf SettingsUri
	err := c.ShouldBindUri(&cf)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	switch cf.Name {
	case "site":
		var site config.SiteInfo
		err := c.ShouldBindJSON(&site)
		if err != nil {
			res.FailWithCode(res.ArgumentError, c)
		}
		global.Config.SiteInfo = site
	case "email":
		var email config.Email
		err := c.ShouldBindJSON(&email)
		if err != nil {
			res.FailWithCode(res.ArgumentError, c)
		}
		global.Config.Email = email
	case "qq":
		var qq config.QQ
		err := c.ShouldBindJSON(&qq)
		if err != nil {
			res.FailWithCode(res.ArgumentError, c)
		}
		global.Config.QQ = qq
	case "qiniu":
		var qiniu config.QiNiu
		err := c.ShouldBindJSON(&qiniu)
		if err != nil {
			res.FailWithCode(res.ArgumentError, c)
		}
		global.Config.QiNiu = qiniu
	case "jwt":
		var jwt config.JWT
		err := c.ShouldBindJSON(&jwt)
		if err != nil {
			res.FailWithCode(res.ArgumentError, c)
		}
		global.Config.JWT = jwt
	default:
		res.FailWithMessage("请求的服务配置不存在", c)
		return
	}

	err = core.WriteConf()
	if err != nil {
		global.Logger.Error(err)
		res.FailWithMessage(err.Error(), c)
		return
	}
	res.OkWith(c)
}
