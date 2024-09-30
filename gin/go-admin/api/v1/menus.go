package v1

import (
	"net/http"

	"go-admin/api/request"
	"go-admin/pkg/jwt"
	"go-admin/service"

	"github.com/gin-gonic/gin"
)

// getChildrenMenu 获取子菜单
func getChildrenMenu(parentId int, allMenus []*request.AllMenu) []*request.MenuReply {
	data := make([]*request.MenuReply, 0)
	for _, v := range allMenus {
		if v.ParentId == parentId {
			data = append(data, &request.MenuReply{
				ID:            v.ID,
				Name:          v.Name,
				WebIcon:       v.WebIcon,
				Sort:          v.Sort,
				Path:          v.Path,
				Level:         v.Level,
				ParentId:      v.ParentId,
				ComponentName: v.ComponentName,
				SubMenus:      getChildrenMenu(v.ID, allMenus),
			})
		}
	}
	return data
}

// 生成树形菜单
func allMenuToMenuReply(allMenu []*request.AllMenu) []*request.MenuReply {
	reply := make([]*request.MenuReply, 0)
	// 一层循环，得到顶层菜单
	for _, v := range allMenu {
		if v.ParentId == 0 {
			reply = append(reply, &request.MenuReply{
				ID:            v.ID,
				Name:          v.Name,
				WebIcon:       v.WebIcon,
				Sort:          v.Sort,
				Path:          v.Path,
				Level:         v.Level,
				ComponentName: v.ComponentName,
				SubMenus:      getChildrenMenu(v.ID, allMenu),
			})
		}
	}
	return reply
}

// Menus 获取菜单列表
func Menus(c *gin.Context) {
	// 登录用户信息
	userClaim := c.MustGet("UserClaim").(*jwt.UserClaim)

	// 根据角色获取所有菜单列表数据
	//tx := model.GetMenusList()
	allMenus, err := service.Menus(c, userClaim.RoleId, userClaim.IsAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	data := allMenuToMenuReply(allMenus)
	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "加载成功",
		"result": data,
	})
}

// GetMenuList 设置-获取菜单列表
func GetMenuList(c *gin.Context) {
	Menus(c)
}

// AddMenu 新增菜单
func AddMenu(c *gin.Context) {
	in := new(request.AddMenuRequest)
	if err := c.ShouldBindJSON(in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	// 新增菜单
	if err := service.AddMenu(c, in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "新增成功",
	})

}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	in := new(request.UpdateMenuRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	if in.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填参不能为空",
		})
		return
	}

	// 更新数据
	if err := service.UpdateMenu(c, in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败！",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})

}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除失败！",
		})
		return
	}
	// 删除数据库中的数据
	if err := service.DeleteMenu(c, id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
