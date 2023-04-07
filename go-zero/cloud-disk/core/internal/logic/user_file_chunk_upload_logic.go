package logic

import (
	"context"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileChunkUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileChunkUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileChunkUploadLogic {
	return &UserFileChunkUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileChunkUploadLogic) UserFileChunkUpload(req *types.UserFileChunkUploadReq, etag string) (resp *types.UserFileChunkUploadResp, err error) {
	// todo: add your logic here and delete this line

	resp = &types.UserFileChunkUploadResp{
		Etag: etag,
	}
	return
}
