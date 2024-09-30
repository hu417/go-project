package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go-admin/api/request"
	"go-admin/dao/db"
	"go-admin/global"
	"go-admin/model"
	"go-admin/pkg/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetUserList 获取管理员列表
func GetUserList(c *gin.Context, in *request.GetUserListRequest) (interface{}, error) {

	var (
		cnt  int64
		list = make([]*request.GetUserListReply, 0)
	)
	err := db.GetUserList(in.Keyword).
		Count(&cnt).
		Offset((in.Page - 1) * in.Size).
		Limit(in.Size).
		Find(&list).Error

	data := struct {
		List  []*request.GetUserListReply `json:"list"`
		Count int64                       `json:"count"`
	}{
		List:  list,
		Count: cnt,
	}

	return data, err
}

// CheckUserByName 判断用户名是否存在
func CheckUserByName(c *gin.Context, name string) (cnt int64, err error) {
	err = global.DB.Model(new(model.SysUser)).Where("username = ?", name).Count(&cnt).Error
	return cnt, err
}

// CheckUserByIdAndName 判断指定用户id和用户名是否存在
func CheckUserByIdAndName(c *gin.Context, id uint, name string) (cnt int64, err error) {
	err = global.DB.Model(new(model.SysUser)).Where("id != ? AND username = ?", id, name).Count(&cnt).Error
	return
}

// AddUser 新增管理员信息
func AddUser(c *gin.Context, in *request.AddUserRequest) error {

	// 2. 保存用户数据
	return global.DB.Create(&model.SysUser{
		Username: in.Username,
		Password: in.Password,
		Phone:    in.Phone,
		Remarks:  in.Remarks,
		RoleId:   in.RoleId,
	}).Error

}

// GetUserDetail 根据ID获取管理员详情信息
func GetUserDetail(c *gin.Context, uId int) (data *request.GetUserDetailReply, err error) {

	// 获取用户基本信息
	sysUser, err := db.GetUserDetail(uint(uId))
	if err != nil {

		return nil, err
	}
	data = &request.GetUserDetailReply{
		ID: sysUser.ID,
		AddUserRequest: request.AddUserRequest{
			Username: sysUser.Username,
			Password: sysUser.Password,
			Phone:    sysUser.Phone,
			Remarks:  sysUser.Remarks,
			RoleId:   sysUser.RoleId,
		},
	}

	return data, nil

}

// UpdateUser 修改管理员信息
func UpdateUser(c *gin.Context, in *request.UpdateUserRequest) error {

	return global.DB.Model(new(model.SysUser)).Where("id = ?", in.ID).Updates(map[string]any{
		"password": in.Password,
		"username": in.Username,
		"phone":    in.Phone,
		"remarks":  in.Remarks,
		"role_id":  in.RoleId,
	}).Error

}

// DeleteUser 删除管理员信息
func DeleteUser(c *gin.Context, id string) error {

	// 删除管理员
	return global.DB.Where("id = ?", id).Delete(new(model.SysUser)).Error

}

// UpdateInfo 更新个人信息
func UpdateInfo(c *gin.Context, id uint, in *request.UpdateUserRequest) error {
	// 更新数据
	return global.DB.Model(new(model.SysUser)).Where("id = ?", id).Updates(map[string]any{
		"sex":      in.Sex,
		"avatar":   in.Avatar,
		"phone":    in.Phone,
		"nickname": in.Nickname,
	}).Error

}

// SendEmail 发送邮件
func SendEmail(c *gin.Context, id uint, toEmail string) error {
	// 用来保存邮件验证码
	session := sessions.Default(c)
	// 登录用户信息
	sysUser, err := db.GetUserDetail(id)
	if err != nil {
		return err
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
	return session.Save()

}

// UpdateEmail 更新个人邮箱
func UpdateEmail(c *gin.Context, id uint, email string) error {

	// 更新数据
	return global.DB.Model(new(model.SysUser)).Where("id = ?", id).Updates(map[string]any{
		"email": email,
	}).Error

}

// UpdatePwd 更新个人密码
func UpdatePwd(c *gin.Context, id uint, in *request.UpdatePwdRequest) error {

	// 根据用户ID获取用户信息
	sysUser, err := db.GetUserDetail(id)
	if err != nil {
		return errors.New("更新失败,该用户不存在！")
	}
	// 判断输入旧密码是否正确
	if sysUser.Password != in.UsedPass {
		return errors.New("更新失败,输入的旧密码不正确！")
	}
	// 更新数据
	return global.DB.Model(new(model.SysUser)).Where("id = ?", id).Updates(map[string]any{
		"password": in.NewPass,
	}).Error

}
