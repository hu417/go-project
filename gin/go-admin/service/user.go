package service

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-admin/define"
	"go-admin/models"
	"go-admin/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList 获取管理员列表
func GetUserList(c *gin.Context) {
	in := &GetUserListRequest{NewQueryRequest()}
	err := c.ShouldBindQuery(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	var (
		cnt  int64
		list = make([]*GetUserListReply, 0)
	)
	err = models.GetUserList(in.Keyword).Count(&cnt).Offset((in.Page - 1) * in.Size).Limit(in.Size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "加载成功",
		"result": gin.H{
			"list":  list,
			"count": cnt,
		},
	})
}

// AddUser 新增管理员信息
func AddUser(c *gin.Context) {
	in := new(AddUserRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	// 1. 判断用户名是否存在
	var cnt int64
	err = models.DB.Model(new(models.SysUser)).Where("username = ?", in.Username).Count(&cnt).Error
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
	err = models.DB.Create(&models.SysUser{
		Username: in.Username,
		Password: in.Password,
		Phone:    in.Phone,
		Remarks:  in.Remarks,
		RoleId:   in.RoleId,
	}).Error

	if err != nil {
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
	data := new(GetUserDetailReply)
	// 1. 获取用户基本信息
	sysUser, err := models.GetUserDetail(uint(uId))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	data.ID = sysUser.ID
	data.Remarks = sysUser.Remarks
	data.Phone = sysUser.Phone
	data.Username = sysUser.Username
	data.Password = sysUser.Password
	data.RoleId = sysUser.RoleId
	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "获取成功",
		"result": data,
	})

}

// UpdateUser 修改管理员信息
func UpdateUser(c *gin.Context) {
	in := new(UpdateUserRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	// 1. 判断用户名是否已存在
	var cnt int64
	err = models.DB.Model(new(models.SysUser)).Where("id != ? AND username = ?", in.ID, in.Username).Count(&cnt).Error
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
	err = models.DB.Model(new(models.SysUser)).Where("id = ?", in.ID).Updates(map[string]any{
		"password": in.Password,
		"username": in.Username,
		"phone":    in.Phone,
		"remarks":  in.Remarks,
		"role_id":  in.RoleId,
	}).Error
	if err != nil {
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
	err := models.DB.Where("id = ?", id).Delete(new(models.SysUser)).Error
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

// UpdateInfo 更新个人信息
func UpdateInfo(c *gin.Context) {
	// 登录用户信息
	userClaim := c.MustGet("UserClaim").(*define.UserClaim)
	in := new(UpdateUserRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	err = models.DB.Model(new(models.SysUser)).Where("id = ?", userClaim.Id).Updates(map[string]any{
		"sex":      in.Sex,
		"avatar":   in.Avatar,
		"phone":    in.Phone,
		"nickname": in.Nickname,
	}).Error
	if err != nil {
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
	// 用来保存邮件验证码
	session := sessions.Default(c)
	// 登录用户信息
	userClaim := c.MustGet("UserClaim").(*define.UserClaim)
	sysUser, err := models.GetUserDetail(userClaim.Id)
	toEmail := c.Query("email")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发送邮件失败！",
		})
		return
	}
	// 随机生成六位数验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	VCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	content := "验证码：" + VCode + "此验证码用于更换邮箱绑定，请勿将验证码告知他人！"
	if toEmail == "" {
		toEmail = sysUser.Email
	}
	go utils.SendEmail(toEmail, "修改邮箱验证码", content)
	// 设置
	session.Set("VCode", VCode)
	session.Set("hello", "goland")
	// 保存
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "保存验证码失败！",
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

	// 登录用户信息
	userClaim := c.MustGet("UserClaim").(*define.UserClaim)
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

	// 更新数据
	err := models.DB.Model(new(models.SysUser)).Where("id = ?", userClaim.Id).Updates(map[string]any{
		"email": email,
	}).Error
	if err != nil {
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
	userClaim := c.MustGet("UserClaim").(*define.UserClaim)
	in := new(UpdatePwdRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败,旧密码和新密码不能为空！",
		})
		return
	}

	// 根据用户ID获取用户信息
	sysUser, err := models.GetUserDetail(userClaim.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败,该用户不存在！",
		})
		return
	}
	// 判断输入旧密码是否正确
	if sysUser.Password != in.UsedPass {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败,输入的旧密码不正确！",
		})
		return
	}
	// 更新数据
	err = models.DB.Model(new(models.SysUser)).Where("id = ?", userClaim.Id).Updates(map[string]any{
		"password": in.NewPass,
	}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败！",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新密码成功！",
	})
}
