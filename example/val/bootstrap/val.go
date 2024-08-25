package bootstrap

import (
	"reflect"
	"strings"
	"val/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() {
	// 注册自定义验证器

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		if err := v.RegisterValidation("mobile", utils.ValidateMobile); err != nil {
			panic(err)
		}

		// 注册自定义 json tag 函数
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

}
