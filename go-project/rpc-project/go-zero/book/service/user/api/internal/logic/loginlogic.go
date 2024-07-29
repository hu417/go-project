package logic

import (
	"context"
	//"errors"
	"strings"

	"time"

	"book/common/errorx"
	"book/service/user/api/internal/svc"
	"book/service/user/api/internal/types"
	"book/service/user/model"

	"github.com/golang-jwt/jwt"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginReply, error) {
	// if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
	// 	return nil, errors.New("参数错误")
	// }
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, errorx.NewDefaultError("参数错误")
	}

	userInfo, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	switch err {
	case nil:
	case model.ErrNotFound:
		// return nil, errors.New("用户名不存在")
		return nil, errorx.NewDefaultError("用户名不存在")
	default:
		return nil, err
	}

	if userInfo.Password != req.Password {
		// return nil, errors.New("用户密码不正确")
		return nil, errorx.NewDefaultError("用户密码不正确")
	}

	// ---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, userInfo.Id)
	if err != nil {
		return nil, err
	}
	// ---end---

	return &types.LoginReply{
		Id:           userInfo.Id,
		Name:         userInfo.Name,
		Gender:       userInfo.Gender,
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

// func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginReply, err error) {
// 	// todo: add your logic here and delete this line

// 	return
// }
