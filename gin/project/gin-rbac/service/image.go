package service

import (
	"errors"
	"os"
	"path/filepath"

	"gin-rbac/common/errs"
	"gin-rbac/db/dao"
	"gin-rbac/db/model"
	"gin-rbac/dtos"
	"gin-rbac/global"

	"gorm.io/gorm"
)

// ImageService 图片服务
type ImageService interface {
	// UploadImage 上传图片
	UploadImage(uploadImageReqDTO *dtos.UploadImageReqDTO) (*dtos.UploadImageResDTO, error)
	// GetImage 根据id获取图片
	GetImage(getImageReqDTO *dtos.GetImageReqDTO) (*dtos.GetImageResDTO, error)
	// DeleteImage 删除图片
	DeleteImage(deleteImageReqDTO *dtos.DeleteImageReqDTO) error
}

// imageService 图片服务实现
type imageService struct {
	imageDao dao.ImageDao
}

// NewImageService 创建图片服务
func NewImageService(imageDao dao.ImageDao) ImageService {
	return &imageService{
		imageDao: imageDao,
	}
}

// UploadImage 上传图片
func (s *imageService) UploadImage(uploadImageReqDTO *dtos.UploadImageReqDTO) (*dtos.UploadImageResDTO, error) {
	var image = &model.ImageModel{
		ImageName: uploadImageReqDTO.ImageName,
		ImagePath: uploadImageReqDTO.ImagePath,
	}
	imageModel, err := s.imageDao.UploadImage(image)
	if err != nil {
		global.Log.Errorln("Failed to upload image: ", err)
		return nil, errs.ErrInternalServerError
	}
	return &dtos.UploadImageResDTO{
		CreatedAt: imageModel.CreatedAt.Format("2006-01-02 15:04:05"),
		ID:        imageModel.ID,
		ImageUrl:  filepath.Join(imageModel.ImagePath, imageModel.ImageName),
	}, nil
}

// GetImage 根据id获取图片
func (s *imageService) GetImage(getImageReqDTO *dtos.GetImageReqDTO) (*dtos.GetImageResDTO, error) {
	imageModel, err := s.imageDao.GetImage(getImageReqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warnln("Failed to get image, Image not found: ", err)
			return nil, errs.ErrImageNotFound
		}
		global.Log.Errorln("Failed to get image: ", err)
		return nil, errs.ErrInternalServerError
	}
	return &dtos.GetImageResDTO{
		CreatedAt: imageModel.CreatedAt.Format("2006-01-02 15:04:05"),
		ID:        imageModel.ID,
		ImageUrl:  filepath.Join(imageModel.ImagePath, imageModel.ImageName),
	}, nil
}

// DeleteImage 根据id删除图片
func (s *imageService) DeleteImage(deleteImageReqDTO *dtos.DeleteImageReqDTO) error {
	image, err := s.imageDao.GetImage(deleteImageReqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warnln("Failed to delete image, Image not found: ", err)
			return errs.ErrImageNotFound
		}
		global.Log.Errorln("Failed to delete image, Failed to get image: ", err)
		return errs.ErrInternalServerError
	}
	imageUrl := filepath.Join(image.ImagePath, image.ImageName)
	if err := os.Remove(imageUrl); err != nil {
		global.Log.Errorln("Failed to delete image, Failed to remove image: ", err)
		return errs.ErrInternalServerError
	}
	err = s.imageDao.DeleteImage(deleteImageReqDTO.ID)
	if err != nil {
		global.Log.Errorln("Failed to delete image: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}
