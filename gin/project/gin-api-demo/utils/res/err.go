package res

type CustomError struct {
	Code int `json:"code"`
	Msg  interface{} `json:"msg"`
}

type CustomErrors struct {
	BusinessError CustomError
	ValidateError CustomError
}

var Errors = CustomErrors{
	BusinessError: CustomError{40000, "业务错误"},
	ValidateError: CustomError{42200, "请求参数错误"},
}
