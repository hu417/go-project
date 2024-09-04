package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 上传单个文件
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file") //  file是表单字段名字
	if err != nil {
		c.String(500, "上传文件出错")
	}

	//  上传到指定路径
	c.SaveUploadedFile(file, "static/"+file.Filename)
	c.String(http.StatusOK, "fileName:", file.Filename)

}

// 上传多个文件
func UploadMultipleFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
	}

	//  获取所有文件; 
	files := form.File["files"]
	for _, file := range files {
		//  逐个存
		fmt.Println(file.Filename)
		//  上传到指定路径
		c.SaveUploadedFile(file, "static/"+file.Filename)
	}
	c.String(200, fmt.Sprintf("upload ok %d files", len(files)))

}
