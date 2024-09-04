package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 上传单个
func Upload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.String(500, "上传文件出错")
	}

	// 上传到指定路径
	c.SaveUploadedFile(file, "./file/"+file.Filename)
	c.String(http.StatusOK, "fileName:", file.Filename)

}

// 上传多个
func Uploads(c *gin.Context) {

	// 获取MultipartForm
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
	}

	// 获取所有文件
	files := form.File["files"]
	for _, file := range files {
		// 逐个存
		fmt.Println(file.Filename)
	}
	c.String(200, fmt.Sprintf("upload ok %d files", len(files)))

}
