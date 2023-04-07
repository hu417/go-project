package logic

import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileUploadPrepareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileUploadPrepareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileUploadPrepareLogic {
	return &UserFileUploadPrepareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileUploadPrepareLogic) UserFileUploadPrepare(req *types.UserFileUploadPrepareReq) (resp *types.UserFileUploadPrepareResp, err error) {
	// todo: add your logic here and delete this line

	// 1、根据md5值进行判断
	rp := new(models.Repository_pool)
	hs, err := l.svcCtx.Engine.Where("hash = ? ", req.Md5).Get(rp)
	if err != nil {
		return
	}

	resp = new(types.UserFileUploadPrepareResp)
	if hs {
		// 秒传成功
		resp.Identity = rp.Identity

	} else {
		// 获取该文件的UploadID、Key,用来进行文件的分片上传
		key, uploadId, err := helper.CosInitPart(req.Ext)
		if err != nil {
			return nil, err
		}
		resp.Key = key
		resp.UploadId = uploadId
	}

	return
}
