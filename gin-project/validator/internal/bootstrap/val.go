package bootstrap

import (
	"fmt"
	"reflect"
	"strings"

	"gin-validator/internal/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

// 定义一个全局翻译器T
var Trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		//注册一个获取json的自定义方法
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境;后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTrans.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zhTrans.RegisterDefaultTranslations(v, Trans)
		default:
			err = enTrans.RegisterDefaultTranslations(v, Trans)
		}

		// 注册自定义的验证器
		// {
		// 	// 自定义的验证器
		// 	// password
		// 	if err := v.RegisterValidation("password", utils.ValidatePassword); err != nil {
		// 		panic("password验证器注册失败: " + err.Error())
		// 	}

		// 	// 翻译自定义的验证器
		// 	err = v.RegisterTranslation("password", Trans, func(ut ut.Translator) error {
		// 		return ut.Add("password", "{0} 密码格式不正确！", true) // 添加翻译
		// 	}, func(ut ut.Translator, fe validator.FieldError) string {
		// 		t, _ := ut.T("password", fe.Field()) // 使用翻译
		// 		return t
		// 	})
		// }

		// 使用函数来注册
		registerCustomValidator(v, "password", utils.ValidatePassword)
		registerCustomValidator(v, "mobile", utils.ValidateMobile)
		registerCustomValidator(v, "email", utils.ValidateEmail)
		return
	}
	return
}

// 注册自定义的验证器及其翻译
func registerCustomValidator(v *validator.Validate, tag string, validationFunc func(fl validator.FieldLevel) bool) {
	if err := v.RegisterValidation(tag, validationFunc); err != nil {
		panic(fmt.Sprintf("%s验证器注册失败: %s", tag, err.Error()))
	}

	err := v.RegisterTranslation(tag, Trans, func(ut ut.Translator) error {
		return ut.Add(tag, "{0} 格式不正确！", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})

	if err != nil {
		panic(fmt.Sprintf("注册%s翻译失败: %s", tag, err.Error()))
	}
}
