package logic

import (
	"context"
	"errors"
	"fmt"

	"bluebell/dao/mysql"
	"bluebell/dao/mysql/curd"
	"bluebell/models/req"
	"bluebell/models/resp"
	"bluebell/models/table"
	"bluebell/pkg/jwt"
	"bluebell/pkg/utils"
)

// 存放业务逻辑的代码

// 注册
func SignUp(ctx context.Context, p *req.ParamSignUp) error {

	// 1.判断用户存不存在
	_, ok, err := curd.NewUserDao(ctx, mysql.GetDB()).CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("用户已存在")
	}

	// 2.生成UID
	// userID := snowflake.GenID()
	// 构造一个User实例
	user := &table.User{
		Username: p.Username,
	}
	data,err := utils.GeneratePassword(p.Password)
	if err != nil {
		return err
	}
	user.Password = data

	// 3.保存进数据库
	return curd.NewUserDao(ctx, mysql.GetDB()).InsertUser(user)
}

// 登陆
func Login(ctx context.Context, p *req.ParamLogin) (users *resp.ParamToken, err error) {

	// 根据用户名获取密码
	user, ok, err := curd.NewUserDao(ctx, mysql.GetDB()).CheckUserExist(p.Username)
	if err != nil {
		return nil, fmt.Errorf("[svc] login fail => %w", err)
	}
	if ok {
		// 用户已存在
		if !utils.CheckPassword(p.Password, user.Password) {
			return nil, fmt.Errorf("[svc] check_password fail => %w", errors.New("密码错误"))
		}
		// 生成JWT
		token, err := jwt.GenToken(user.UserId, user.Username)
		if err != nil {
			return nil, fmt.Errorf("[svc] gen_token fail => %w", err)
		}
		users = &resp.ParamToken{
			UserId:   user.UserId,
			Username: user.Username,
			Token:    token,
		}
		return users, nil
	}

	return nil, fmt.Errorf("[svc] user fail => %w", errors.New("用户不存在"))
}
