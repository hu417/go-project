package logic

import (
	"bluebell/controller/request"
	"bluebell/dao/mysql/sql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/utils"
	"errors"
)

// 存放业务逻辑的代码

func SignUp(p *request.ParamSignUp) (err error) {
	// 1.判断用户存不存在
	if _, err := sql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2.构造一个User实例
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 3.保存进数据库
	return sql.InsertUser(user)
}

func Login(p *request.ParamLogin) (user *models.User, err error) {
	// 1.判断用户存不存在
	user, _ = sql.CheckUserExist(p.Username)
	if user.Id == 0 {
		return nil, errors.New("用户不存在")
	}

	// 2.判断密码是否正确
	if !utils.BcryptMakeCheck(p.Password, user.Password) {
		return nil, errors.New("密码不正确")
	}

	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
