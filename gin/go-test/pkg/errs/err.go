package errs

import "errors"

// Error 自定义错误
type Error struct {
	Code    int
	Message string
}

// errorCodeMap 错误码
var ErrorCodeMap = map[error]int{}

// newError 新建错误
func newError(code int, msg string) error {
	err := errors.New(msg)
	ErrorCodeMap[err] = code
	return err
}

// Error 实现error接口
func (e Error) Error() string {
	return e.Message
}
