package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"go-admin/api/request"
	"go-admin/pkg/jwt"
	"go-admin/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetUserList 获取管理员列表
func GetUserList(c *gin.Context) {
	in := &request.GetUserListRequest{request.NewQueryRequest()}
	err := c.ShouldBindQuery(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	data, err := service.GetUserList(c, in)
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

// AddUser 新增管理员信息
func AddUser(c *gin.Context) {
	in := new(request.AddUserRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	// 1. 判断用户名是否存在

	cnt, err := service.CheckUserByName(c, in.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常，保存失败！",
		})
		return
	}

	// 大于0说明用户已经存在
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户已存在",
		})
		return
	}

	// 2. 保存用户数据
	if err := service.AddUser(c, in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常,保存失败！",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "保存成功",
	})

}

// GetUserDetail 根据ID获取管理员详情信息
func GetUserDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填参数不能为空",
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

	// 获取用户基本信息
	data, err := service.GetUserDetail(c, uId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "获取成功",
		"result": data,
	})

}

// UpdateUser 修改管理员信息
func UpdateUser(c *gin.Context) {
	in := new(request.UpdateUserRequest)
	if err := c.ShouldBindJSON(in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	// 1. 判断用户名是否已存在
	cnt, err := service.CheckUserByIdAndName(c, in.ID, in.Username)
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
			"msg":  "用户已存在",
		})
		return
	}

	// 2. 修改数据
	if err := service.UpdateUser(c, in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

// DeleteUser 删除管理员信息
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填参不能为空",
		})
		return
	}

	// 删除管理员
	if err := service.DeleteUser(c, id); err != nil {
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

// UpdateInfo 更新个人信息
func UpdateInfo(c *gin.Context) {
	// 登录用户信息
	userClaim := c.MustGet("UserClaim").(*jwt.UserClaim)
	in := new(request.UpdateUserRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	if err = service.UpdateInfo(c, userClaim.Id, in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新个人信息成功",
	})

}

// SendEmail 发送邮件
func SendEmail(c *gin.Context) {
	toEmail := c.Query("email")

	// 登录用户信息
	userClaim := c.MustGet("UserClaim").(*jwt.UserClaim)
	if err := service.SendEmail(c, userClaim.Id, toEmail); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"msg":   "发送邮件失败！",
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "邮件已经发送，请注意查收！",
	})
}

// VerifyCode 校验邮箱验证码
func VerifyCode(c *gin.Context) {
	// 初始化 session 对象
	session := sessions.Default(c)
	// 通过 session.Get 读取 session 值
	VCode := session.Get("VCode")
	fmt.Println("VCode", VCode)
	code := c.Query("code")
	if VCode != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "校验失败,输入验证码不正确！",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "校验成功！",
	})
}

// UpdateEmail 更新个人邮箱
func UpdateEmail(c *gin.Context) {

	// 新的邮箱
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败,请输入邮箱！",
		})
		return
	}
	// 新的验证码
	newCode := c.Query("code")
	if newCode == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败,请输入验证码！",
		})
		return
	}
	// 获取保存到session里的验证码
	// 初始化 session 对象
	session := sessions.Default(c)
	// 通过 session.Get 读取 session 值
	VCode := session.Get("VCode")
	if VCode == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败,验证码已过期！",
		})
		return
	}
	if newCode != VCode {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败,输入的验证码不正确！",
		})
		return
	}

	// 登录用户信息
	userClaim := c.MustGet("UserClaim").(*jwt.UserClaim)

	// 更新数据
	if err := service.UpdateEmail(c, userClaim.Id, email); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败！",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新邮箱成功！",
	})

}

// UpdatePwd 更新个人密码
func UpdatePwd(c *gin.Context) {
	// 登录用户信息
	userClaim := c.MustGet("UserClaim").(*jwt.UserClaim)
	in := new(request.UpdatePwdRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败,旧密码和新密码不能为空！",
		})
		return
	}

	// 根据用户ID获取用户信息
	if err := service.UpdatePwd(c, userClaim.Id, in); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"msg":   "更新失败！",
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新密码成功！",
	})
}
