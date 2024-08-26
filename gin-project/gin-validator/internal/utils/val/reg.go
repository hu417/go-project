package val

import (
	"regexp"
	"unicode"

	"gin-validator/internal/request"

	"github.com/go-playground/validator/v10"
)

// 自定义用户名验证
func ValidUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	for _, char := range username {
		// 不能包含特殊字符
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

// 自定义密码验证
func ValidatePassword(fl validator.FieldLevel) bool {
	// 内部通过反射获取mobile的值
	password := fl.Field().String()
	// 密码包含至少一位数字，字母和特殊字符,且长度8-16
	// (?![0-9a-zA-Z]+$) 表示匹配后面非数字，非字母的字符
	// (?![a-zA-Z!@#$%^&*]+$) 表示匹配后面非字母，非特殊字符的字符
	// (?![0-9!@#$%^&*]+$) 表示匹配后面非数字，非特殊字符的字符
	// [0-9A-Za-z!@#$%^&*]{8,16} 表示长度8-16，可以是任意的数字，字母和其中包含的特殊字符

	reg := regexp.MustCompile(`^(\?=\.*\d)(\?=.*[a-z])(\?=\.*[A-Z])[a-zA-Z0-9]{8,16}$`)
	return reg.MatchString(password)
}

// 自定义re_password结构体校验
func ValidationRePassword(sl validator.StructLevel) {
	su := sl.Current().Interface().(request.Signup)

	if su.Password != su.RePassword {
		// 输出错误提示信息，最后一个参数就是传递的param
		sl.ReportError(su.RePassword, "re_password", "RePassword", "eqfield", "password")
	}
}

// 自定义手机号验证
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]16[6]|7[1-35-8]|9[189])\d{8}$`, mobile) // 用正则去匹配
	return ok
}

// 自定义邮箱验证
func ValidateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	// fmt.Printf("email => %s", email)
	/*
		^和$分别表示匹配字符串的开始和结束，
		[w.-]+表示匹配一个或多个字母、数字、下划线、点号和短横线，
		[a-zA-Z0-9]+表示匹配一个或多个字母、数字，.表示匹配点号，
		[a-zA-Z]{2,4}表示匹配两到四个字母。
	*/
	reg := regexp.MustCompile(`^[\w\.-]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,4}$`)
	return reg.MatchString(email)
}
