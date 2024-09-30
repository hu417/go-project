package v1

import (
	"errors"
	"net/http"

	"go-admin/api/request"
	"go-admin/global"
	"go-admin/pkg/jwt"
	"go-admin/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginPassword 用户登录
func LoginPassword(c *gin.Context) {
	in := new(request.LoginPasswordRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	// 根据账号、密码查询用户
	sysUser, err := service.LoginPassword(c, in)
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
	token, err := jwt.GenerateToken(sysUser.ID, sysUser.RoleId, sysUser.Username, global.TokenExpire)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	// 刷新token
	refreshToken, err := jwt.GenerateToken(sysUser.ID, sysUser.RoleId, sysUser.Username, global.RefreshTokenExpire)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	data := &request.LoginPasswordReply{
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
