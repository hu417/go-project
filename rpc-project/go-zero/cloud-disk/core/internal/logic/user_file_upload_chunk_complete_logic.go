package logic

import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileUploadChunkCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileUploadChunkCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileUploadChunkCompleteLogic {
	return &UserFileUploadChunkCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileUploadChunkCompleteLogic) UserFileUploadChunkComplete(req *types.UserFileUploadChunkCompleteReq) (resp *types.UserFileUploadChunkCompleteResp, err error) {
	// todo: add your logic here and delete this line

	// 分片文件上传完成
	co := make([]cos.Object, 0)
	for _, v := range req.CosObjects {
		co = append(co, cos.Object{
			ETag:       v.Etag,
			PartNumber: v.PartNumber,
		})
	}

	err = helper.CosPartUploadComplete(req.Key, req.UploadId, co)

	return
}
