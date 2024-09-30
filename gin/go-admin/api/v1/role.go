package v1

import (
	"net/http"
	"strconv"

	"go-admin/api/request"
	"go-admin/service"

	"github.com/gin-gonic/gin"
)

// GetRoleList 角色列表
func GetRoleList(c *gin.Context) {
	in := &request.GetRoleListRequest{request.NewQueryRequest()}
	if err := c.ShouldBindQuery(in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	// 获取角色
	data, err := service.GetRoleList(c, in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "加载成功",
		"result": data,
	})
}

// AddRole 新增角色
func AddRole(c *gin.Context) {
	in := new(request.AddRoleRequest)
	if err := c.ShouldBindJSON(in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	// 1. 判断角色名称是否存在
	cnt, err := service.CheckRoleByName(c, in.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常，保存失败！",
		})
		return
	}

	// 大于0说明角色名称已经存在
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "角色名称已存在",
		})
		return
	}

	// 2. 新增角色
	if err := service.AddRole(c, in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常,保存失败！",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "新增角色成功",
	})

}

// PatchRoleAdmin 更改管理员身份
func PatchRoleAdmin(c *gin.Context) {
	id := c.Param("id")
	isAdmin := c.Param("isAdmin")
	if id == "" || isAdmin == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填参数不能为空",
		})
		return
	}

	// 更改管理员身份
	err := service.PatchRoleAdmin(c, id, isAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改管理员身份成功",
	})
}

// GetRoleDetail 根据ID获取角色详情
func GetRoleDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填参不能为空",
		})
		return
	}
	uId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	data, err := service.GetRoleDetail(c, uId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取数据失败！",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "获取成功",
		"result": data,
	})
}

// UpdateRole 修改角色信息
func UpdateRole(c *gin.Context) {
	in := new(request.UpdateRoleRequest)
	if err := c.ShouldBindJSON(in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	// 1. 判断角色名称是否已存在
	cnt, err := service.CheckRoleByIdAndName(c, in.ID, in.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "角色名称已存在",
		})
		return
	}

	// 2. 修改数据
	if err := service.UpdateRole(c, in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "修改失败！",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

// DeleteRole 根据ID删除角色
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填参不能为空",
		})
		return
	}
	// 删除角色
	err := service.DeleteRole(c, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})

}

// AllRole 获取所有角色
func AllRole(c *gin.Context) {
	list, err := service.AllRole(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取角色列表失败！",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "加载成功",
		"result": list,
	})
}
