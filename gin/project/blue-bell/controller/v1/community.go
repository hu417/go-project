package v1

import (
	"blue-bell/controller/req"
	"blue-bell/logic"
	"blue-bell/pkg/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

// CreateCommunityHandler 创建社区
func CommunityHandler(c *gin.Context) {
	var comm req.Community
	if err := c.ShouldBindJSON(&comm); err != nil {
		c.JSON(200, gin.H{
			"msg": "参数错误",
		})
		return
	}
	if comm.Name == "" || comm.Introduction == "" {
		c.JSON(200, gin.H{
			"msg": "参数不能为空",
		})
		return
	}

	if err := logic.CreateCommunity(&comm); err != nil {
		c.JSON(200, gin.H{
			"msg": "添加社区失败",
			"err": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "添加社区成功",
	})
}

// CommunityListHandler 社区列表
func CommunityListHandler(c *gin.Context) {
	var p req.Page
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(200, gin.H{
			"msg": "参数错误",
			"err": err.Error(),
		})
		return
	}

	fmt.Printf("p: %+v\n", p)

	if p.Page < 1 || p.Size < 1 {
		p.Page = 1
		p.Size = 10
	} 
	if p.Order == "" {
		p.Order = "id desc"
	}

	total, data, err := logic.GetCommunityList(p.Name, p.Page, p.Size)
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "获取社区列表失败",
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "获取社区列表成功",
		"data": gin.H{
			"total": total,
			"page":  p.Page,
			"size":  p.Size,
			"list":  data,
		},
	})
}

// CommunityDetailByIDHandler 社区详情
func CommunityDetailByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(200, gin.H{
			"msg": "参数错误",
		})
		return
	}
	id_int64 := utils.StrToInt64(id)
	data, err := logic.GetCommunityDetailByID(id_int64)
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "获取社区详情失败",
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "获取社区详情成功",
		"data": data,
	})
}