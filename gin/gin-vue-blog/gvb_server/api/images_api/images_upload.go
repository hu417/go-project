package images_api

import (
	"fmt"
	"gvb_server/global"
	"gvb_server/models/res"
	"gvb_server/utils"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// WhiteImageList 图片上传的白名单
	WhiteImageList = []string{
		"jpg",
		"png",
		"jpeg",
		"ico",
		"tiff",
		"gif",
		"svg",
		"webp",
	}
)

type FileUploadResponse struct {
	FileName  string `json:"file_name"`  // 文件名
	InSuccess bool   `json:"in_success"` // 是否上传成功
	Msg       string `json:"msg"`        // 消息
}

// 上传多个图片
func (ImagesApi) ImagesUploadView(c *gin.Context) {
	// 获取请求头参数
	form, err := c.MultipartForm()
	if err != nil {
		res.FailWithMessage(err.Error(), c)
	}

	fileList, ok := form.File["images"]
	if !ok {
		res.FailWithMessage("上传的图片不存在", c)
	}

	// 判断文件路径是否存在
	basePath := global.Config.Upload.Path
	_, err = os.ReadDir(basePath)
	if err != nil {
		err = os.MkdirAll(basePath, fs.ModePerm)
		if err != nil {
			global.Logger.Error(err.Error())
		}
	}

	// 响应消息列表
	var resList []FileUploadResponse

	for _, file := range fileList {
		// global.Logger.Info(file.Header)
		// global.Logger.Info(file.Filename)
		// global.Logger.Info(file.Size)

		// 判断文件格式是否正确
		imageNames := strings.Split(file.Filename, ".")
		suffix := strings.ToLower(imageNames[len(imageNames)-1])
		if !utils.InList(suffix, WhiteImageList) {
			resList = append(resList, FileUploadResponse{
				FileName:  file.Filename,
				InSuccess: false,
				Msg:       fmt.Sprintf("图片格式 .%s 不正确,请重新上传", suffix),
			})
			continue
		}

		// 文件保存路径
		filePath := path.Join(basePath, file.Filename)

		// 判断图片大小
		size := float64(file.Size) / float64(1024*1024)
		if size >= float64(global.Config.Upload.Size) {
			resList = append(resList, FileUploadResponse{
				FileName:  file.Filename,
				InSuccess: false,
				Msg:       fmt.Sprintf("上传的图片大小超出%d MB限制,当前大小为%2.f MB", global.Config.Upload.Size, size),
			})
			continue
		}

		// 上传图片
		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			global.Logger.Error(err.Error())
			// 上传成功
			resList = append(resList, FileUploadResponse{
				FileName:  file.Filename,
				InSuccess: false,
				Msg:       fmt.Sprintf("上传失败,%s", err.Error()),
			})
		}

		// 上传成功
		resList = append(resList, FileUploadResponse{
			FileName:  filePath,
			InSuccess: true,
			Msg:       "上传成功",
		})

	}
	res.OkWithData(resList, c)

}
