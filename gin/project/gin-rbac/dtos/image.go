package dtos

type UploadImageReqDTO struct {
	ImagePath string `form:"image_path" binding:"required"` // 图片存储路径
	ImageName string `form:"image_name" binding:"required"` // 图片文件名
}

type UploadImageResDTO struct {
	ID        uint   `json:"id"`         // 图片id
	ImageUrl  string `json:"image_url"`  // 图片url
	CreatedAt string `json:"created_at"` // 创建时间
}

type GetImageReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 图片id
}

type GetImageResDTO struct {
	ID        uint   `json:"id"`         // 图片id
	ImageUrl  string `json:"image_url"`  // 图片url
	CreatedAt string `json:"created_at"` // 创建时间
}

type DeleteImageReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 图片id
}
