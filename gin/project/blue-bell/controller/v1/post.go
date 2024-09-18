package v1

import (
	"blue-bell/controller/e"
	"blue-bell/controller/req"
	"blue-bell/controller/res"
	"blue-bell/logic"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	var p req.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		res.ResponseError(c, 400, e.CodeInvalidParam, err.Error())
		return
	}

	if err := logic.CreatePost(&p); err != nil {
		res.ResponseError(c, 500, e.CodeServerBusy, err.Error())
		return
	}
	res.ResponseSuccess(c, 200, "帖子创建成功", nil)
}

// PostListHandler 获取帖子列表
func PostListHandler(c *gin.Context) {
	var p req.Page
	if err := c.ShouldBind(&p); err != nil {
		res.ResponseError(c, 400, e.CodeInvalidParam, "参数错误")
		return
	}

	if p.Page < 1 || p.Size < 1 {
		p.Page = 1
		p.Size = 10
	}
	if p.Order == "" {
		p.Order = "id desc"
	}

	data, err := logic.GetPostList(&p)

	if err != nil {
		res.ResponseError(c, 500, e.CodeServerBusy, err.Error())
		return
	}
	res.ResponseSuccess(c, 200, "获取帖子列表成功", data)
}

// PostDetailHandler 获取帖子详情
func PostDetailByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		res.ResponseError(c, 400, e.CodeInvalidParam, "参数错误")
		return
	}
	data, err := logic.GetPostByID(id)
	if err != nil {
		res.ResponseError(c, 500, e.CodeServerBusy, err.Error())
		return
	}
	res.ResponseSuccess(c, 200, "获取帖子详情成功", data)
}
