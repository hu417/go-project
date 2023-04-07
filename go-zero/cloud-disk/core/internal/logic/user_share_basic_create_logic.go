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

type UserShareBasicCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserShareBasicCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserShareBasicCreateLogic {
	return &UserShareBasicCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserShareBasicCreateLogic) UserShareBasicCreate(req *types.UserShareBasicCreateReq, UserIdentity string) (resp *types.UserShareBasicCreateResp, err error) {
	// todo: add your logic here and delete this line

	// 创建文件分享记录
	uuid := helper.GetUuid()
	ur := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Where("identity = ? ", req.UserRepositoryIdentity).Get(ur)
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New(req.UserRepositoryIdentity + " not found or delete")
	}
	data := &models.Share_basic{
		Identity:                 uuid,
		User_identity:            UserIdentity,
		User_Repository_Identity: req.UserRepositoryIdentity,
		Repository_identity:      ur.RepositoryIdentity,
		Expired_time:             req.Expiredtime,
	}

	_, err = l.svcCtx.Engine.Insert(data)
	if err != nil {
		return
	}

	resp = &types.UserShareBasicCreateResp{
		Identity: uuid,
	}
	return
}
