package logic

import (
	"context"
	"errors"
	"fmt"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileMoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileMoveLogic {
	return &UserFileMoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileMoveLogic) UserFileMove(req *types.UserFileMoveReq, UserIdentity string) (resp *types.UserFileMoveResp, err error) {
	// todo: add your logic here and delete this line

	// 1、判断文件夹存不存在
	has, err := l.svcCtx.Engine.Where("identity = ? AND user_identity = ? ", req.Identity, UserIdentity).Get(new(models.UserRepository))
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New("文件夹不存在")
	}

	// 2、更新 ParentID,注意，文件信息deletetime未被记录
	_, err = l.svcCtx.Engine.Where("identity = ? ", req.Identity).Update(&models.UserRepository{
		ParentId: req.ParentIdentity,
	})
	if err != nil {
		return
	}

	resp = &types.UserFileMoveResp{
		Message: fmt.Sprintf("文件: %v,文件夹: %v", req.Identity, req.ParentIdentity),
	}
	return
}
