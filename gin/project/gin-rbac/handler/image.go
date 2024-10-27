package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gin-rbac/common/response"
	"gin-rbac/dtos"
	"gin-rbac/global"
	"gin-rbac/service"
	"gin-rbac/utils"

	"github.com/gin-gonic/gin"
)

// ImageHandler 图片处理器
type ImageHandler struct {
	ImageService service.ImageService
}

// NewImageHandler 初始化图片处理器
func NewImageHandler(imageService service.ImageService) *ImageHandler {
	return &ImageHandler{
		ImageService: imageService,
	}
}

// UploadImage 上传图片
//	@Summary		Upload Image
//	@Description	Upload Image
//	@Tags			Image
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			image	formData	file					true	"Upload Image request"
//	@Success		200		{object}	dtos.UploadImageResDTO	"Successfully response with image information"
//	@Router			/images/upload [post]
func (h *ImageHandler) UploadImage(c *gin.Context) {
	// 单个图片的最大大小为8 MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8<<20)
	// 获取上传的图片
	file, err := c.FormFile("image")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "Failed to upload image:"+err.Error())
		return
	}
	// 获取图片扩展名
	ext := filepath.Ext(file.Filename)
	// 定义允许的图片扩展名白名单
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	// 验证文件扩展名是否在白名单中
	isValidExtension := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			isValidExtension = true
			break
		}
	}
	if !isValidExtension {
		response.Fail(c, http.StatusBadRequest, "Invalid image extension")
		return
	}
	// 获取当前时间和年月
	currentTime := time.Now()
	yearMonth := fmt.Sprintf("%d/%02d", currentTime.Year(), currentTime.Month())
	// 生成新的文件名
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	// 获取项目根目录
	projectRoot, err := os.Getwd()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "Failed to get project root directory: "+err.Error())
		return
	}
	parentDir := filepath.Dir(projectRoot)
	// 设置文件保存路径（相对路径）
	relativeBasePath := filepath.Join("uploads", "images", yearMonth)

	// 绝对路径用于文件操作
	dstAbsolute := filepath.Join(parentDir, relativeBasePath, newFileName)
	// 创建 uploads 目录，如果它不存在的话
	err = os.MkdirAll(filepath.Dir(dstAbsolute), os.ModePerm)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "Failed to create upload directory: "+err.Error())
		return
	}
	// 保存上传的文件
	if err := c.SaveUploadedFile(file, dstAbsolute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uploadImageReqDTO := dtos.UploadImageReqDTO{
		ImageName: newFileName,
		ImagePath: relativeBasePath,
	}
	uploadImageResDTO, err := h.ImageService.UploadImage(&uploadImageReqDTO)
	if err != nil {
		response.FailWithData(c, utils.GetStatusCodeFromError(err), "Failed to upload image: "+err.Error(), nil)
		return
	}
	response.OkWithData(c, http.StatusOK, uploadImageResDTO)
}

// GetImage 获取图片
//	@Summary		Get Image
//	@Description	Get Image
//	@Tags			Image
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		uint				true	"Get Image request"
//	@Success		200	{object}	dtos.GetImageResDTO	"Successfully response with image information"
//	@Router			/images/{id} [get]
func (h *ImageHandler) GetImage(c *gin.Context) {
	var getImageReqDTO dtos.GetImageReqDTO
	if err := c.ShouldBindUri(&getImageReqDTO); err != nil {
		global.Log.Warnln("Failed to get image, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	getImageResDTO, err := h.ImageService.GetImage(&getImageReqDTO)
	if err != nil {
		response.FailWithData(c, utils.GetStatusCodeFromError(err), "Failed to get image: "+err.Error(), nil)
		return
	}
	response.OkWithData(c, http.StatusOK, getImageResDTO)
}

// DeleteImage 删除图片
//	@Summary		Delete Image
//	@Description	Delete Image
//	@Tags			Image
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		uint	true	"Delete Image request"
//	@Success		200	{string}	string	"Successfully response with image information"
//	@Router			/images/{id} [delete]
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	var deleteImageReqDTO dtos.DeleteImageReqDTO
	if err := c.ShouldBindUri(&deleteImageReqDTO); err != nil {
		global.Log.Warnln("Failed to delete image, Invalid request data: ", err)
		response.Fail(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	err := h.ImageService.DeleteImage(&deleteImageReqDTO)
	if err != nil {
		response.Fail(c, utils.GetStatusCodeFromError(err), "Failed to delete image: "+err.Error())
		return
	}
	response.OkWithMsg(c, http.StatusOK, "Image deleted successfully")
}
