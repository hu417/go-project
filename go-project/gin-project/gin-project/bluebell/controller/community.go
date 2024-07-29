package controller

import (
	"bluebell/logic"
	"bluebell/pkg/resp"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ---- 跟社区相关的 ----

func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区（community_id, community_name) 以列表的形式返回
	data, err := logic.GetCommunityList(c)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		resp.ResponseError(c, resp.CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	resp.ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区id
	idStr := c.Param("id")                     // 获取URL参数
	id, err := strconv.ParseInt(idStr, 10, 64) // string转int64
	if err != nil {
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}

	// 2. 根据id获取社区详情
	data, err := logic.GetCommunityDetail(c, id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		resp.ResponseError(c, resp.CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	resp.ResponseSuccess(c, data)
}
