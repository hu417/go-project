package logic

import (
	"book/service/search/api/internal/svc"
	"book/service/search/api/internal/types"
	"context"
	"encoding/json"
	"fmt"

	"book/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchReply, err error) {
	// todo: add your logic here and delete this line

	// 添加一个log来输出从jwt解析出来的userId
	logx.Infof("userId: %v", l.ctx.Value("userId")) // 这里的key和生成jwt token时传入的key一致

	userIdNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	logx.Infof("userId: %s", userIdNumber)
	userId, err := userIdNumber.Int64()
	if err != nil {
		return nil, err
	}
	// 使用user rpc
	var people *user.UserInfoReply
	people, err = l.svcCtx.UserRpc.GetUser(l.ctx, &user.IdReq{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("user rpc返回的值是: ---> ", people, " <---")
	return &types.SearchReply{
		Name:  people.Name,
		Count: 100,
	}, nil
}
