package logic

import (
	"context"
	"errors"

	// "math/rand"
	// "strconv"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MailCodeSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendLogic {
	return &MailCodeSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendLogic) MailCodeSend(req *types.MailCodeReq) (resp *types.MailCodeResp, err error) {
	// todo: add your logic here and delete this line

	// 1、邮箱不存在
	cnt, err := l.svcCtx.Engine.Where("email = ?", req.Email).Count(new(models.User_basic))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("该邮箱已存在")
	}

	// 2、获取随机验证码
	// code := strconv.Itoa(rand.Intn(100000) + 100000)
	code := helper.RandCode()

	// 3、存储验证码
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExoire))
	if err != nil {
		return nil, err
	}

	// 4、发送验证码
	err = helper.SendEmail(req.Email, code)
	if err != nil {
		return nil, err
	}
	resp = &types.MailCodeResp{
		Msg:  "验证码发送成功!",
		Code: code,
	}

	return
}
