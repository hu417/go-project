package dao

import (
	"context"
	
	"gin-rbac/db/model"

	"gorm.io/gorm"
)

// ImageDao 图片数据访问接口
type ImageDao interface {
	// UploadImage 上传图片
	UploadImage(image *model.ImageModel) (*model.ImageModel, error)
	// GetImage 根据id获取图片
	GetImage(id uint) (*model.ImageModel, error)
	// DeleteImage 根据id删除图片
	DeleteImage(id uint) error
}

// imageDao 图片数据访问接口实现
type imageDao struct {
	db  *gorm.DB
	ctx context.Context
}

// NewImageDAO 创建图片数据访问实例
func NewImageDAO(db *gorm.DB, ctx context.Context) ImageDao {
	return &imageDao{
		db:  db,
		ctx: ctx,
	}
}

// UploadImage 上传图片
func (i *imageDao) UploadImage(image *model.ImageModel) (*model.ImageModel, error) {
	var imageModel *model.ImageModel
	err := Transaction(i.db, func(tx *gorm.DB) error {
		var err error
		imageModel, err = i.uploadImageTx(tx, image)
		return err
	})
	if err != nil {
		return nil, err
	}
	return imageModel, nil
}

// uploadImageTx 上传图片事务内执行
func (i *imageDao) uploadImageTx(tx *gorm.DB, image *model.ImageModel) (*model.ImageModel, error) {
	// 在这里执行所有需要在事务中完成的操作
	if err := tx.WithContext(i.ctx).Create(image).Error; err != nil {
		return nil, err
	}
	// 返回上传成功的图片模型
	return image, nil
}

// GetImage 根据id获取图片
func (i *imageDao) GetImage(id uint) (*model.ImageModel, error) {
	var image model.ImageModel
	err := i.db.WithContext(i.ctx).Where("id = ?", id).First(&image).Error
	return &image, err
}

// DeleteImage 根据id删除图片
func (i *imageDao) DeleteImage(id uint) error {
	return Transaction(i.db, func(tx *gorm.DB) error {
		return tx.WithContext(i.ctx).Delete(&model.ImageModel{}, id).Error
	})
}
