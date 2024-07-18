package logic

import (
	"context"
	"fmt"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileDeleteLogic {
	return &UserFileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileDeleteLogic) UserFileDelete(req *types.UserFileDeleteReq, UserIdentity string) (resp *types.UserFileDeleteResp, err error) {
	// todo: add your logic here and delete this line

	// 软删除，只是删除数据库文件信息，添加deletetime字段值
	_, err = l.svcCtx.Engine.Where("user_identity = ? AND identity = ?", UserIdentity, req.Identity).Delete(new(models.UserRepository))
	if err != nil {
		return
	}
	resp = &types.UserFileDeleteResp{
		Message: fmt.Sprintf("identity = %v 删除成功", req.Identity),
	}
	return
}
