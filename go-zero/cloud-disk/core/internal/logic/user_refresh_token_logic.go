package logic

import (
	"context"

	"strconv"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRefreshTokenLogic {
	return &UserRefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRefreshTokenLogic) UserRefreshToken(req *types.UserRefreshTokenReq, authorization string) (resp *types.UserRefreshTokenResp, err error) {
	// todo: add your logic here and delete this line

	// 根据用户token进行解析，并重新生成
	uc, err := helper.AnalyzeJwtTkoen(authorization)
	if err != nil {
		return
	}

	token, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.TokenExpires)
	if err != nil {
		return
	}
	refreshToken, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.RefreshTokenExpires)

	resp = new(types.UserRefreshTokenResp)
	resp.Token = token
	resp.RefreshToken = refreshToken
	resp.Message = "token 过期时间: " + strconv.Itoa(define.TokenExpires) + "s"

	return
}
