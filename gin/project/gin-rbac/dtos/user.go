package dtos

import (
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type LoginUserReqDTO struct {
	Credential string `json:"credential" binding:"required"`            // 用户名、电话号码或电子邮件
	Password   string `json:"password" binding:"required,min=8,max=20"` // 密码
}

type TokenResDTO struct {
	Token string `json:"token"` // 登录token
}

type CreateUserReqDTO struct {
	Username string `json:"username" binding:"required,min=1,max=15"` // 用户名
	Password string `json:"password" binding:"required,min=8,max=20"` // 密码
	PhoneNum string `json:"phone_num" binding:"omitempty,len=11"`     // 手机号
	Email    string `json:"email" binding:"omitempty,max=256"`        // 邮箱
}

type GetFullUserReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 用户ID
}

type FullUserDTO struct {
	ID         uint           `json:"id"`          // 用户ID
	Username   string         `json:"username"`    // 用户名
	PhoneNum   string         `json:"phone_num"`   // 手机号
	Email      string         `json:"email"`       // 邮箱
	Sex        int8           `json:"sex"`         // 性别
	Intro      string         `json:"intro"`       // 简介
	AvatarPath string         `json:"avatar_path"` // 头像路径
	AvatarName string         `json:"avatar_name"` // 头像文件名
	CreatedAt  time.Time      `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time      `json:"updated_at"`  // 更新时间
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`  // 删除时间
}

type GetFullUserResDTO struct {
	ID              uint   `json:"id"`                // 用户ID
	Username        string `json:"username"`          // 用户名
	PhoneNum        string `json:"phone_num"`         // 手机号
	Email           string `json:"email"`             // 邮箱
	Sex             int8   `json:"sex"`               // 性别
	Intro           string `json:"intro"`             // 简介
	AvatarUrl       string `json:"avatar_url"`        // 头像路径
	CreatedAtFormat string `json:"created_at_format"` // 创建时间式化: "2006-01-02"
	UpdatedAtFormat string `json:"updated_at_format"` // 更新时间式化: "2006-01-02"
	DeletedAt       string `json:"deleted_at"`        // 删除时间
}

// FullUserDTOToGetFullUserResDTO 将FullUserDTO转换为GetFullUserResDTO
func FullUserDTOToGetFullUserResDTO(user FullUserDTO) GetFullUserResDTO {
	return GetFullUserResDTO{
		ID:              user.ID,
		Username:        user.Username,
		PhoneNum:        user.PhoneNum,
		Email:           user.Email,
		Sex:             user.Sex,
		Intro:           user.Intro,
		AvatarUrl:       "http://8.138.88.52/" + filepath.Join(user.AvatarPath, user.AvatarName),
		CreatedAtFormat: user.CreatedAt.Format("2006-01-02"),
		UpdatedAtFormat: user.UpdatedAt.Format("2006-01-02"),
		DeletedAt:       user.DeletedAt.Time.Format("2006-01-02"),
	}
}

type PublicUserListReqDTO struct {
	PaginationReqDTO        // 分页请求
	Username         string `form:"username" binding:"omitempty"` // 关键词
}

type PublicUserReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 用户ID
}

type PublicUserDTO struct {
	ID         uint      `json:"id"`          // 用户ID
	Username   string    `json:"username"`    // 用户名
	Sex        int8      `json:"sex"`         // 性别
	AvatarPath string    `json:"avatar_path"` // 头像路径
	AvatarName string    `json:"avatar_name"` // 头像文件名
	Intro      string    `json:"intro"`       // 简介
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
}

type PublicUserResDTO struct {
	ID              uint   `json:"id"`                // 用户ID
	Username        string `json:"username"`          // 用户名
	Sex             int8   `json:"sex"`               // 性别
	AvatarUrl       string `json:"avatar_url"`        // 头像路径
	Intro           string `json:"intro"`             // 简介
	CreatedAtFormat string `json:"created_at_format"` // 创建时间格式化: "2006-01-02"
}

// PublicUserDTOToPublicUserResDTO 将models.User转换为GetUserResDTO
func PublicUserDTOToPublicUserResDTO(user PublicUserDTO) PublicUserResDTO {
	return PublicUserResDTO{
		ID:              user.ID,
		Username:        user.Username,
		Sex:             user.Sex,
		Intro:           user.Intro,
		AvatarUrl:       "http://8.138.88.52/" + filepath.Join(user.AvatarPath, user.AvatarName),
		CreatedAtFormat: user.CreatedAt.Format("2006-01-02"),
	}
}

type UpdateUserReqDTO struct {
	Username string `json:"username" binding:"omitempty,min=1,max=15"` // 用户名
	PhoneNum string `json:"phone_num" binding:"omitempty,len=11"`      // 手机号
	Email    string `json:"email" binding:"omitempty,max=256"`         // 邮箱
	Sex      int8   `json:"sex" binding:"omitempty,oneof=0 1 2"`       // 性别, 0: 未设置, 1: 男, 2: 女
	Intro    string `json:"intro" binding:"omitempty,max=30"`          // 简介
}

type UpdateUserAvatarReqDTO struct {
	AvatarID uint `json:"avatar_id" binding:"required"` // 头像ID
}

type UpdatePasswordReqDTO struct {
	OldPassword string `json:"old_password" binding:"required"` // 当前密码-旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

type ResetPasswordReqDTO struct {
	Credential  string `json:"credential"`                      // 用户名、电话号码或电子邮件
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

type DeleteUserReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 用户ID
}

type RecoverUserReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 用户ID
}
