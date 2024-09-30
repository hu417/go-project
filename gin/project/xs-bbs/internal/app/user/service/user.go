package service

import (
	"context"

	"xs-bbs/internal/app/user/model"
	"xs-bbs/internal/errs"
	"xs-bbs/pkg/jwt"
	"xs-bbs/pkg/utils"
	"xs-bbs/global"

	"github.com/gogf/gf/util/gconv"
	"github.com/pkg/errors"
)

// Register .
func (u *userService) Register(ctx context.Context, param *model.RegisterReq) (resDto *UserDto, err error) {
	var uModel model.User

	if err = u.repo.CheckUserByUserName(ctx, param.Username); err != nil {
		return
	}

	if err = u.repo.CheckUserByEmail(ctx, param.Email); err != nil {
		return
	}

	if err = gconv.Struct(param, &uModel); err != nil {
		err = errors.Wrap(errs.ErrConvDataErr, err.Error())
		return
	}

	if err = u.repo.Insert(ctx, &uModel); err != nil {
		return
	}

	if err = gconv.Struct(uModel, &resDto); err != nil {
		err = errors.Wrap(errs.ErrConvDataErr, err.Error())
		return
	}

	return
}

// Login 登陆
func (u *userService) Login(ctx context.Context, signIn *model.LoginReq) (token string, err error) {
	var user *model.User
	// 获取用户信息
	if user, err = u.repo.GetUserByName(ctx, signIn.Username); err != nil {
		return
	}

	// 验证密码
	if ! utils.BcryptMakeCheck(signIn.Password, user.Password) {
		return "", errors.New("密码错误")
	}

	// 生成token
	return jwt.GenerateJwt(user.UserID, user.Username,"admin",global.Conf.Jwt.Secret,global.Conf.Jwt.JwtTtl)
}

// Delete 根据用户ID删除用户
func (u *userService) Delete(ctx context.Context, userID string) bool {
	return u.repo.Delete(ctx, userID)
}

// Update 根据用户ID修改用户
func (u *userService) Update(ctx context.Context, user *UserDto) error {
	var uModel model.User
	if err := gconv.Struct(user, &uModel); err != nil {
		return err
	}
	return u.repo.Update(ctx, &uModel)
}

// SelectByName 根据用户名查询用户
func (u *userService) SelectByName(ctx context.Context, userName string) (resDto *UserDto, err error) {
	var uModel *model.User

	if uModel, err = u.repo.GetUserByName(ctx, userName); err != nil {
		return
	}

	if err = gconv.Struct(uModel, &resDto); err != nil {
		err = errors.Wrap(errs.ErrConvDataErr, err.Error())
		return
	}

	return
}

// SelectByID 根据用户ID查询用户
func (u *userService) SelectByID(ctx context.Context, userID string) (resDto *UserDto, err error) {
	var uModel *model.User

	if uModel, err = u.repo.GetUserByID(ctx, userID); err != nil {
		return
	}

	if err = gconv.Struct(uModel, &resDto); err != nil {
		err = errors.Wrap(errs.ErrConvDataErr, err.Error())
		return
	}

	return
}
