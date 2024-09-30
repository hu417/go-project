package service

import (
	"go-admin/api/request"
	"go-admin/dao/db"
	"go-admin/model"

	"github.com/gin-gonic/gin"
)

// LoginPassword 用户登录
func LoginPassword(c *gin.Context, in *request.LoginPasswordRequest) (*model.SysUser, error) {

	// 根据账号、密码查询用户
	return db.GetUserByUsernamePassword(in.Username, in.Password)

}
