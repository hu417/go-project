package fileload

import "github.com/gin-gonic/gin"

func Router(e *gin.Engine) {
	upload(e)     // 单个文件
	uploadMore(e) // 多个文件
}