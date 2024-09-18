package v1

import (
	"blue-bell/controller/e"
	"blue-bell/controller/req"
	"blue-bell/controller/res"
	"blue-bell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(req.Vote)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("c shouldBindJson error ", zap.Error(err))
		res.ResponseError(c, 400, e.CodeInvalidParam, "参数有误")
		return
	}

	// 判断该帖子是否过了投票时间
	// 拿到userID给指定的帖子(postID)投票，zset中记录帖子投票数的有序集合为帖子id（postid）中的元素名称即为userID,分数为投票数
	// 更改该帖子的分数

	// 获取当前用户ID
	userId, err := strconv.ParseInt(c.GetString("user_id"), 10, 64)

	if userId == 0 || err != nil {
		zap.L().Error("getCurrentUserID error ", zap.Error(err))
		res.ResponseError(c, 401, e.CodeNeedLogin, "未登录")
		return
	}

	if err = logic.PostVote(p.PostID, p.Direction, userId); err != nil {
		zap.L().Error("logic PostVote error ", zap.Error(err))
		res.ResponseError(c, 500, e.CodeServerBusy, "服务器繁忙")
		return
	}
	// 返回响应
	res.ResponseSuccess(c, 200, e.CodeVoteSuccess, nil)

}
