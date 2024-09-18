package e

import "errors"

var (
	ErrorUserNotLogin = errors.New("用户未登录")

	ErrorUserExist                    = errors.New("用户已存在")
	ErrorUserNotExist                 = errors.New("用户不存在")
	ErrorLoginUserNameOrPassWordError = errors.New("用户名或密码错误")
	ErrorCommunityExist                    = errors.New("社区已存在")
	ErrorCommunityNotExist                 = errors.New("社区不存在")
)
