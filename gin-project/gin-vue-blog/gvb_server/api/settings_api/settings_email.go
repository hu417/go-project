package settings_api

import (
	"gvb_server/global"
	"gvb_server/models/res"

	"github.com/gin-gonic/gin"
)

// 视图函数 - 数据响应
// 获取Email信息接口
func (SettingsApi) SettingsEmailView(c *gin.Context) {
	//
	res.OkWithData(global.Config.Email, c)

}
