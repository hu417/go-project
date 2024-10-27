package bootstrap

import (
	"fmt"
	"reflect"
	"strings"

	"gin-api-demo/pkg/val"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

// InitTrans 初始化翻译器
func InitTrans(locale string) (trans ut.Translator, err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 为SignUpParam注册自定义校验方法
		// v.RegisterStructValidation(val.ValidationRePassword, request.Signup{})

		// 在校验器注册自定义的校验方法
		if err := v.RegisterValidation("mobile", val.ValidateMobile); err != nil {
			return nil, err
		}
		// if err := v.RegisterValidation("email", val.ValidateEmail); err != nil {
		// 	return nil, err
		// }
		// if err := v.RegisterValidation("username", val.ValidUsername); err != nil {
		// 	return nil, err
		// }
		if err := v.RegisterValidation("password", val.ValidatePassword); err != nil {
			return nil, err
		}

		// 注册翻译器
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器

		// 第一个参数是备用（fallback）的语言环境；后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTrans.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTrans.RegisterDefaultTranslations(v, trans)
		default:
			err = enTrans.RegisterDefaultTranslations(v, trans)
		}

		return
	}
	return
}
