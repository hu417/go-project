package logic

import (
	"context"
	"errors"

	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailsLogic {
	return &UserDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailsLogic) UserDetails(req *types.UserDetailsReq) (resp *types.UserDetailsResp, err error) {
	// todo: add your logic here and delete this line

	resp = &types.UserDetailsResp{}
	ub := new(models.User_basic)
	has, err := l.svcCtx.Engine.Where("identity = ?", req.Identity).Get(ub)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("identity 错误")
	}
	resp.Name = ub.Name
	resp.Email = ub.Email
	return
}
