package logic

import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (resp *types.UserRegisterResp, err error) {
	// todo: add your logic here and delete this line

	// 1、判断code是否一致
	code, err := l.svcCtx.RDB.Get(l.ctx, req.Email).Result()
	if err != nil {
		return nil, err
	}
	if code != req.Code {
		return nil, errors.New("验证码不一致")
	}

	// 2、判断用户是否存在
	cnt, err := l.svcCtx.Engine.Where("name = ?", req.Name).Count(new(models.User_basic))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("用户已存在")
	}

	// 3、新建用户
	user := &models.User_basic{
		Identity: helper.GetUuid(),
		Name:     req.Name,
		Password: helper.Md5(req.Password),
		Email:    req.Email,
	}
	n, err := l.svcCtx.Engine.Insert(user)
	if err != nil {
		return nil, err
	}
	logx.Infof("sql受影响的行数: %v", n)

	// 4、返回参数
	resp = &types.UserRegisterResp{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		Code:     "******",
	}

	return
}
