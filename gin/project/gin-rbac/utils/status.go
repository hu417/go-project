package utils

import (
	"net/http"
	"strings"

	"gin-rbac/common/errs"
)

// GetStatusCodeFromError 根据错误类型获取相应的HTTP状态码
func GetStatusCodeFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	errorName := errs.GetErrorName(err)

	// 根据错误变量名中的关键字来确定状态码
	switch {
	case strings.Contains(errorName, "NotFound"):
		return http.StatusNotFound
	case strings.Contains(errorName, "AlreadyExists"), strings.Contains(errorName, "Conflict"):
		return http.StatusConflict
	case strings.Contains(errorName, "Invalid"), strings.Contains(errorName, "SameAsOld"):
		return http.StatusBadRequest
	case strings.Contains(errorName, "NoFieldsUpdated"):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
