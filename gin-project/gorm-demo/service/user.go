package service

import (
	"gorm-demo/dao"
	"gorm-demo/global"
	"gorm-demo/model"
	"math"
)

type UserSvc struct {
}

func NewUserSvc() *UserSvc {
	return &UserSvc{}
}

// 检测用户是否存在
func (u *UserSvc) CheckUser(user *model.User) bool {
	return dao.NewUserDao(global.DB).CheckUser(user)
}

// 创建用户
func (u *UserSvc) CreateUser(user *model.User) error {
	return dao.NewUserDao(global.DB).CreateUser(user)
}

// 查询用户
func (u *UserSvc) GetUser(id int) (user *model.User, err error) {
	return dao.NewUserDao(global.DB).GetUser(id)
}

// 查询所有用户
func (u *UserSvc) GetUsers(username string, pageSize int, pageNum int) (data interface{}) {

	// 获取所有用户
	users, total := dao.NewUserDao(global.DB).GetUsers(username, pageSize, pageNum)
	// 统计页数
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// 返回数据
	data = struct{
		Users  []*model.User
		Total int64
		TotalPage int
	}{
		Users:    users,
		Total:    total,
		TotalPage: totalPages,
	}
	
	return
}

// 编辑用户
func (u *UserSvc) EditUser(id int, user *model.User) error {
	return dao.NewUserDao(global.DB).EditUser(id, user)
}

// 修改用户密码
func (u *UserSvc) ChangePassword(id int, user *model.User) error {
	return dao.NewUserDao(global.DB).ChangePassword(id, user)
}

// 删除用户
func (u *UserSvc) DeleteUser(user *model.User) error {
	return dao.NewUserDao(global.DB).DeleteUser(user)
}
