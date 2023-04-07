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

type UserShareBasicSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserShareBasicSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserShareBasicSaveLogic {
	return &UserShareBasicSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserShareBasicSaveLogic) UserShareBasicSave(req *types.UserShareBasicSaveReq, UserIdentity string) (resp *types.UserShareBasicSaveResp, err error) {
	// todo: add your logic here and delete this line

	// 1、获取资源详情
	rp := new(models.Repository_pool)
	has, err := l.svcCtx.Engine.Where("identity = ?", req.RepositoryIdentity).Get(rp)
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New(req.RepositoryIdentity + " not found")
	}

	// 2、资源保存
	ur := &models.UserRepository{
		Identity:           helper.GetUuid(),
		UserIdentity:       UserIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                rp.Ext,
		Name:               rp.Name,
	}
	_, err = l.svcCtx.Engine.Insert(ur)
	if err != nil {
		return
	}
	resp = &types.UserShareBasicSaveResp{
		Identity: ur.Identity,
	}
	return
}
