package req

import (
	"gin-base/utils"
)

type Register struct {
	Name     string `form:"name" json:"name" binding:"required,min=6,max=20"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile,len=10"`
}

// GetMessages 自定义错误信息
func (register Register) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Mame.required":     "用户名称不能为空",
		"password.required": "用户密码不能为空",
		"mobile.required":   "手机号码不能为空",
		"mobile.mobile":     "手机号码格式不正确",
	}
}
