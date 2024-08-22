package fileload

import (
	"gin-api-demo/api/file"

	"github.com/gin-gonic/gin"
)

// 上传单个
func upload(r *gin.Engine) {
	// 给表单限制上传大小 (默认 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", file.Upload)
}

// 上传多个
func uploadMore(r *gin.Engine) {
	// 给表单限制上传大小 (默认 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/uploads", file.Uploads)
}
