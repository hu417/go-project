package logic

import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadReq) (resp *types.FileUploadResp, err error) {
	// todo: add your logic here and delete this line

	// 根据handeler传递过来的req参数，保存db中
	rp := &models.Repository_pool{
		Identity: helper.GetUuid(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}

	_, err = l.svcCtx.Engine.Insert(rp)
	if err != nil {
		return nil, err
	}
	resp = &types.FileUploadResp{
		Identity: rp.Identity,
		Name:     rp.Name,
		Ext:      rp.Ext,
	}

	return
}
