package errs

import "errors"

type CustomError interface {
	error
	GetName() string
}

func NewCustomError(name string, msg string) CustomError {
	return customError{name: name, message: msg}
}

type customError struct {
	name    string
	message string
}

func (e customError) Error() string {
	return e.message
}

func (e customError) GetName() string {
	return e.name
}

// GetErrorName 从错误对象中提取错误变量名
func GetErrorName(err error) string {
	if err == nil {
		return ""
	}

	var ce CustomError
	if errors.As(err, &ce) {
		return ce.GetName()
	}

	return ""
}
