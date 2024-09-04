package controller

import (
	"bluebell/logic"
	"bluebell/models/req"
	"bluebell/models/common"
	"bluebell/pkg/resp"
	"bluebell/pkg/valida"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票

//type VoteData struct {
//	// UserID 从请求中获取当前的用户
//	PostID    int64 `json:"post_id,string"`   // 贴子id
//	Direction int   `json:"direction,string"` // 赞成票(1)还是反对票(-1)
//}

// PostVoteController 投票
func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(req.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			resp.ResponseError(c, resp.CodeInvalidParam)
			return
		}
		errData := valida.RemoveTopStruct(errs.Translate(valida.Trans)) // 翻译并去除掉错误提示中的结构体标识
		resp.ResponseErrorWithMsg(c, resp.CodeInvalidParam, errData)
		return
	}
	// 获取当前请求的用户的id
	uid, ok := c.Get(common.CtxUserIDKey)
	if !ok {
		resp.ResponseError(c, resp.CodeNeedLogin)
		return
	}
	userID, ok := uid.(uint32)
	if !ok {
		resp.ResponseError(c, resp.CodeNeedLogin)
		return
	}


	// 具体投票的业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		resp.ResponseError(c, resp.CodeServerBusy)
		return
	}

	resp.ResponseSuccess(c, nil)
}
