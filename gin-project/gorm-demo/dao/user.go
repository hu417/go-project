package dao

import (
	"gorm-demo/model"
	"gorm-demo/utils"

	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
    return &UserDao{
		DB: db,
	}
}

// CheckUser 查询用户是否存在
func (u *UserDao)CheckUser(user *model.User) bool {

	u.DB.Select("id").Where("username = ?", user.UserName).First(&user)
	if user.ID > 0 {
		// 用户已存在
		return true
	}
	// 用户不存在
	return false
}
 

// CreateUser 新增用户
func (u *UserDao) CreateUser(user *model.User) error {

	return  u.DB.Create(&user).Error
}
 
// GetUser 查询用户
func (u *UserDao) GetUser(id int) (user *model.User, err error) {
	err = u.DB.Limit(1).Where("ID = ?", id).Find(&user).Error

	if err != nil {
		return nil, err
	}
	return user,nil
}
 
// GetUsers 查询用户列表
func (u *UserDao) GetUsers(username string, pageSize int, pageNum int) (users []*model.User, total int64) {

	if username != "" {
		// 模糊查询
		u.DB.Select("id,username,role,created_at").Where(
			"username LIKE ?", username+"%",
		).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)

		// 统计总数
		u.DB.Model(&users).Where(
			"username LIKE ?", username+"%",
		).Count(&total)
		return 
	}
	// 查询所有
	u.DB.Select("id,username,role,created_at").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
	// 统计总数
	u.DB.Model(&users).Count(&total)
	// 统计页数
	// totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return 
}
 
// EditUser 编辑用户信息
func (u *UserDao) EditUser(id int, user *model.User) error {

	var maps = make(map[string]interface{})
	maps["username"] = user.UserName
	maps["role"] = user.Role

	return u.DB.Model(&user).Where("id = ? ", id).Updates(maps).Error
}
 
// ChangePassword 修改密码
func (u *UserDao) ChangePassword(id int, user *model.User) error {
	// 密码加密
	user = &model.User{
		Password: utils.HashPw(user.Password),
	}
	
	return u.DB.Select("password").Where("id = ?", id).Updates(&user).Error
}
 
// DeleteUser 删除用户
func (u *UserDao) DeleteUser(user *model.User) error {

	return u.DB.Where("id = ? ", user.ID).Delete(&user).Error
}
