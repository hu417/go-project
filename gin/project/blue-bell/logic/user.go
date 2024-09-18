package logic

import (
	"errors"
	"fmt"

	"blue-bell/controller/e"
	"blue-bell/controller/req"
	"blue-bell/dao/mysql/user"
	"blue-bell/global"
	"blue-bell/model"
	"blue-bell/pkg/token"
	"blue-bell/pkg/utils"

	"gorm.io/gorm"
)

// Signup 注册业务处理
func Signup(p *req.UserSignUp) error {
	// 1.校验用户是否存在
	u := &model.User{
		UserName: p.Username,
		PassWord: p.Password,
	}
	if _, err := user.CheckUserExist(u); err != nil {
		// 判断是否为用户不存在的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// UID := utils.NanoId()
			// 2.保存到数据库
			return user.Insert(u)
		}
		return err
	}
	// 3.用户已存在
	return e.ErrorUserExist
}

// Login 登录业务处理
func Login(p *req.UserLogin) (string, error) {
	// 构造一个用户表结构体
	user_req := &model.User{
		UserName: p.Username,
		PassWord: p.Password,
	}
	// 1.判断该用户是否存在，如果不存在则返回错误
	user_dao, err := user.CheckUserExist(user_req)
	fmt.Printf("err: %v\n", err)
	if err == nil {
		if utils.BcryptMakeCheck(p.Password, user_dao.PassWord) {
			token_str, err := token.Newjwt().GenerateJwt(user_dao.UserID, user_dao.UserName, "user", global.Conf.Jwt.Secret, global.Conf.Jwt.JwtTtl)
			if err != nil {
				return "", err
			}
			return token_str, nil

		} else {

			return "", e.ErrorLoginUserNameOrPassWordError
		}
	}
	return "", err
}
