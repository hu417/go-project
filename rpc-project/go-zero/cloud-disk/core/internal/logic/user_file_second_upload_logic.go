package logic

import (
	"context"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileSecondUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileSecondUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileSecondUploadLogic {
	return &UserFileSecondUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileSecondUploadLogic) UserFileSecondUpload(req *types.UserFileSecondUploadReq) (resp *types.UserFileSecondUploadResp, err error) {
	// todo: add your logic here and delete this line

	// 1、根据md5值进行判断
	rp := new(models.Repository_pool)
	hs, err := l.svcCtx.Engine.Where("hash = ? ", req.Md5).Get(rp)
	if err != nil {
		return
	}
	if hs {
		resp = &types.UserFileSecondUploadResp{
			// 秒传成功
			Identity: rp.Identity,
		}

	} else {
		// 进行分片上传
		logx.Info("进行分片上传")

	}

	return
}
