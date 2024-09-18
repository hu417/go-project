package e

/*
定义错误码
*/
type ResCode int64

const (
	CodeSuccess             ResCode = 1000 + iota // 默认成功
	CodeInvalidParam                              // 请求参数错误
	CodeUserExist                                 // 用户已存在
	CodeUserNotExist                              // 用户不存在
	CodeUserOrPassWordError                       // 用户名或者密码错误
	CodeNeedLogin                                 // 需要登录才能访问
	CodeInvalidToken                              // Token无效
	CodeServerBusy                                // 服务器繁忙
	CodeVoteSuccess                               // 投票成功
)

var codeMspMap = map[ResCode]string{
	CodeSuccess:             "success",
	CodeInvalidParam:        "请求参数错误",
	CodeUserExist:           "用户已存在",
	CodeUserNotExist:        "用户不存在",
	CodeUserOrPassWordError: "用户名或者密码错误",
	CodeNeedLogin:           "请登录后访问",
	CodeInvalidToken:        "非法Token",
	CodeServerBusy:          "服务器忙",
	CodeVoteSuccess:         "投票成功",
}

func (code ResCode) CodeMsg() string {
	return codeMspMap[code]
}
