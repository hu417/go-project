package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 自定义验证器注册
		_ = v.RegisterValidation("mobile", validateMobile)

	}
}
