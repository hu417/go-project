package request

import (
	"gin-validator/global"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	GetMessages() ValidatorMessages
}

type ValidatorMessages map[string]string

// GetErrorMsg 获取错误信息
func GetErrorMsg(request interface{}, err error) map[string]string {
	// 判断是否是 validator.ValidationErrors 类型
	if ve, isValidatorErrors := err.(validator.ValidationErrors); isValidatorErrors {

		msg := make(ValidatorMessages)
		if _, isValidator := request.(Validator); isValidator {
			for _, v := range ve {
				key := v.Field() + "." + v.Tag()
				if message, exist := request.(Validator).GetMessages()[key]; exist {
					// fmt.Printf("tag: %s, message => %s\n\n", v.Tag(), message)
					msg[v.Field()] = message
				} else {
					msg[v.Field()] = v.Translate(global.Trans)
				}
			}
			return msg
		}
	}

	return map[string]string{
		"err": err.Error(),
	}
}
