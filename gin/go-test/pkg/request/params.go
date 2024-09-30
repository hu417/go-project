package request

import "go-test/pkg/responce"

// RegisterRequest 注册参数
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// LoginRequest 登陆参数
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// LoginResponseData 登陆返回参数
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}

// LoginResponse 登陆返回
type LoginResponse struct {
	responce.Response
	Data LoginResponseData
}

// UpdateProfileRequest 更新用户信息
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}

// GetProfileResponseData 获取用户信息
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}

// GetProfileResponse 获取用户信息返回
type GetProfileResponse struct {
	responce.Response
	Data GetProfileResponseData
}
