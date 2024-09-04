package req

import "api-demo/internal/validator"

type SaveAdmin struct {
	ID       uint   `form:"id" json:"id"`
	Name     string `form:"name" json:"name" binding:"required,min=4,max=20"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Age      int    `form:"age" json:"age" binding:"required,min=1,max=100"`
	Gender   int    `form:"gender" json:"gender" binding:"required,oneof=1 2"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile,len=11"`
}

func (SaveAdmin) GetMessages() validator.ValidateErrorMessages {
	return validator.ValidateErrorMessages{
		"Name.required":     "用户名称不能为空",
		"Password.required": "密码不能为空",
		"Age.required":      "年龄不能为空",
		"Age.min":           "年龄最小为1",
		"Age.max":           "年龄最大为100",
		"Gender.required":   "性别不能为空",
		"Gender.oneof":      "性别只能是男或者女",
		"Mobile.required":   "手机号不能为空",
		"Mobile.mobile":     "手机号格式不正确,只能为11位数字",
	}
}
