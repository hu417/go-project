package logic

import (
	"context"
	//"errors"

	"book/service/shorturl/api/internal/svc"
	"book/service/shorturl/api/internal/types"

	//"book/service/shorturl/rpc/types/transform"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExpandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExpandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExpandLogic {
	return &ExpandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExpandLogic) Expand(req *types.ExpandReq) (resp *types.ExpandResp, err error) {
	// todo: add your logic here and delete this line

	// 调用transformer接口的的Expand方法，传递请求req给其进一步处理，并返回结果

	// resp, err = l.svcCtx.Transformer.Expand(l.ctx, &transform.ExpandReq{S})
	if err != nil {
		return &types.ExpandResp{}, err
	}
	return &types.ExpandResp{
		Url: resp.Url,
	}, nil

}
