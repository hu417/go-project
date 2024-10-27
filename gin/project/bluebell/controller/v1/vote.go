package v1

import (
	"bluebell/controller/request"
	"bluebell/controller/response"
	"bluebell/global"
	"bluebell/logic"

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

func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(request.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			response.ErrorCode(c, response.CodeInvalidParam)
			return
		}
		errData := request.RemoveTopStruct(errs.Translate(request.Trans)) // 翻译并去除掉错误提示中的结构体标识
		response.ErrorWithCodeMsg(c, response.CodeInvalidParam, errData)
		return
	}
	// 获取当前请求的用户的id
	userID := c.GetString(global.CtxUserIDKey)
	if userID == "" {
		response.ErrorCode(c, response.CodeNeedLogin)
		return
	}
	// 具体投票的业务逻辑
	if err := logic.VoteForPost(c, userID, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		response.ErrorCode(c, response.CodeServerBusy)
		return
	}

	response.Success(c, nil)
}
