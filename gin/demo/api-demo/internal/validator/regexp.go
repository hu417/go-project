package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 校验手机号
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 手机号为空则不验证，这里只验证有值的
	if mobile == "" {
		return true
	}

	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)

	return ok
}
