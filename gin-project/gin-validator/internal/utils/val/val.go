package val

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	GetMessages() ValidatorMessages
}

type ValidatorMessages map[string]string

// GetErrorMsg 获取错误信息
func GetErrorMsg(request interface{}, err error) string {
	// 获取错误信息
	if validatorErrs, isValidatorErrors := err.(validator.ValidationErrors); isValidatorErrors {
		var msg string

		// 检查request是否实现了Validator接口
		if validator, isValidator := request.(Validator); isValidator {
			for _, v := range validatorErrs {
				if message, exist := validator.GetMessages()[v.Field()+"."+v.Tag()]; exist {
					return message
				}
			}
		}

		// 如果没有自定义消息，使用默认错误消息
		for _, v := range validatorErrs {
			msg = v.Error()
		}

		return msg

	}

	return "Parameter error"
}

// 去除以"."及其左部分内容
func RemoveTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		// 去掉结构体中的前缀 {User.mobile: "mobile 非法的手机号码!"}
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
