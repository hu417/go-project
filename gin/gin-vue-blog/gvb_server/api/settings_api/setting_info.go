package settings_api

import "github.com/gin-gonic/gin"

// 视图函数 - 数据响应
func (SettingsApi) SettingsInfoView(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "xxx",
	})
}
