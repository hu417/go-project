package logic

import (
	"context"
	"errors"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line

	// 1、从数据库中查询用户
	user := new(models.User_basic)
	has, err := l.svcCtx.Engine.Where("name = ? AND password = ?", req.Name, helper.Md5(req.Password)).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户名/密码错误")
	}
	// 2、生成token
	token, err := helper.GenerateToken(user.Id, user.Identity, user.Name, define.TokenExpires)
	//token, err := helper.GenerateToken(user.Id, user.Identity, user.Name, 20)

	if err != nil {
		return nil, err
	}

	// 3、生成用于刷新token的token
	refreshtoken, err := helper.GenerateToken(user.Id, user.Identity, user.Name, define.RefreshTokenExpires)

	if err != nil {
		return nil, err
	}

	resp = new(types.LoginResponse)
	resp.Token = token
	resp.RefreshToken = refreshtoken

	return
}
