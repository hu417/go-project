package v1

import (
	"net/http"

	"go-admin/api/request"
	"go-admin/service"

	"github.com/gin-gonic/gin"
)

// GetLogList 获取日志列表
func GetLogList(c *gin.Context) {
	in := &request.GetLogListRequest{request.NewQueryRequest()}

	// 绑定参数
	if err := c.ShouldBindQuery(in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	// 查询
	data, err := service.GetLogList(c, in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "加载成功",
		"result": data,
	})
}
