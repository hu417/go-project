package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-admin/define"
	"go-admin/helper"
	"go-admin/models"
	"gorm.io/gorm"
	"net/http"
)

// LoginPassword 用户登录
func LoginPassword(c *gin.Context) {
	in := new(LoginPasswordRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	// 根据账号、密码查询用户
	sysUser, err := models.GetUserByUsernamePassword(in.Username, in.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或密码错误",
			})
			return
		}
	}
	// 生成 token
	token, err := helper.GenerateToken(sysUser.ID, sysUser.RoleId, sysUser.Username, define.TokenExpire)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	// 刷新token
	refreshToken, err := helper.GenerateToken(sysUser.ID, sysUser.RoleId, sysUser.Username, define.RefreshTokenExpire)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	data := &LoginPasswordReply{
		Token:        token,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "登录成功",
		"result":   data,
		"userInfo": sysUser,
	})
}
