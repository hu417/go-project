package logic

import (
	"context"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserShareBasicInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserShareBasicInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserShareBasicInfoLogic {
	return &UserShareBasicInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserShareBasicInfoLogic) UserShareBasicInfo(req *types.UserShareBasicInfoReq) (resp *types.UserShareBasicInfoResp, err error) {
	// todo: add your logic here and delete this line

	// 1、对分享记录的点击+1

	_, err = l.svcCtx.Engine.Table("share_basic").Exec("Update share_basic Set click_num = click_num + 1 Where identity = ? ", req.Identity)
	if err != nil {
		return nil, err
	}

	// 2、获取分享文件的详情
	resp = new(types.UserShareBasicInfoResp)
	_, err = l.svcCtx.Engine.Table("share_basic").Select("share_basic.repository_identity,user_repository.name,repository_pool.ext,repository_pool.size,repository_pool.path").Join("LEFT", "repository_pool", "share_basic.repository_identity = repository_pool.identity").Join("LEFT", "user_repository", "user_repository.identity = user_repository_identity").Where("share_basic.identity = ? ", req.Identity).Get(resp)
	if err != nil {
		return nil, err
	}

	return
}
