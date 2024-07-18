package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	allowExtMap = map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
)

func main() {
	router := gin.Default()
	// 为 multipart forms 类型设置一个较低的内存缓存 (默认是 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 单文件上传
	router.POST("/upload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		//判断后缀
		extName := path.Ext(file.Filename)

		//判断文件是否合法
		_, ok := allowExtMap[extName]
		if !ok {
			c.JSON(200, gin.H{
				"success": false,
				"msg":     "上传文件类型不合法",
			})
			return
		}

		// 保存上传文件到目标目录
		//创建文件保存目录
		day := time.Now().Format("2006/01/02")
		filePath := "./static/img/" + day
		if err := os.MkdirAll(filePath, 0666);err != nil {
			c.JSON(200, gin.H{
				"success": false,
				"msg":     "创建文件失败",
			})
			return
		}
		//按照时间戳保存文件名
		unix := time.Now().UnixNano()
		fileName := strconv.FormatInt(unix, 10) + extName
		filePath = filePath + "/" + fileName

		//上传文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			fmt.Printf("err => %v", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("uploaded fail! err => %v", err))
		}

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// 多文件上传
	router.POST("/uploads", func(c *gin.Context) {
		// 多个 form 类型文件
		form, _ := c.MultipartForm()
		files := form.File["files"]

		for _, file := range files {
			log.Println(file.Filename)

			// 保存到指定文件
			dst := "/Users/lz/go/src/storage/bin/test/dst/" + file.Filename
			if err := c.SaveUploadedFile(file, dst); err != nil {
				fmt.Printf("err => %v", err)
				continue
				// c.String(http.StatusInternalServerError, fmt.Sprintf("uploaded fail! err => %v", err))
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	router.Run(":8080")
}

/*

// 模拟上传行为
curl -X POST http://localhost:8080/upload \
  -F "file=@/Users/lz/go/src/storage/bin/test/src/cover.jpg" \
  -H "Content-Type: multipart/form-data"


curl -X POST http://localhost:8080/upload \
  -F "files=@/Users/lz/go/src/storage/bin/test/src/cover.jpg" \
  -F "files=@/Users/lz/go/src/storage/bin/test/src/cover2.jpg" \
  -H "Content-Type: multipart/form-data"
*/
