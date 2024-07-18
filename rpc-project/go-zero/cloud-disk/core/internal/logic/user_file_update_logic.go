package logic

import (
	"context"
	"errors"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileUpdateLogic {
	return &UserFileUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileUpdateLogic) UserFileUpdate(req *types.UserFileUpdateReq, UserIdentity string) (resp *types.UserFileUpdateResp, err error) {
	// todo: add your logic here and delete this line

	// 判断文件名是否已存在
	cnt, err := l.svcCtx.Engine.Where("name = ? AND parent_id = (SELECT parent_id FROM user_repository ur WHERE ur.identity = ?)", req.Name, req.Identity).Count(new(models.UserRepository))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("文件名已存在")
	}

	// 修改文件名称
	data := &models.UserRepository{
		Name: req.Name,
	}

	_, err = l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.Identity, UserIdentity).Update(data)
	if err != nil {
		return nil, err
	}
	resp = &types.UserFileUpdateResp{
		Identity: req.Identity,
		Name:     req.Name,
	}
	return
}
