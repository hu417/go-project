package settings_api

import (
	"gvb_server/global"
	"gvb_server/models/res"

	"github.com/gin-gonic/gin"
)

// 视图函数 - 数据响应
// 获取系统信息接口
func (SettingsApi) SettingsInfoView(c *gin.Context) {
	// c.JSON(200, gin.H{
	// 	"msg": "xxx",
	// })

	// 响应ok
	// res.Ok(map[string]string{}, "xxx", c)
	// 只返回data
	// res.OkWithData(map[string]string{"a": "aa", "b": "bb"}, c)
	// // 只返回mes
	// res.OkWithMessage("xxx", c)
	// 响应error
	//res.Fail(map[string]string{}, "error", c)
	// // 只返回code
	// res.FailWithCode(1000, c)
	// // 只返回mes
	// res.FailWithMessage("系统错误", c)
	res.OkWithData(global.Config.SiteInfo, c)

}
