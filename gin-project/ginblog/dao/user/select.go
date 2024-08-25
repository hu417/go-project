package user

import (
	"ginblog/model"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/errors"
	"gorm.io/gorm"
)

// 查询用户：根据用户名查询用户
func (dao *UserDao) FindUserByName(user *model.User) bool {
	dao.DB.Select("id").Where("username = ?", user.Username).First(&user)

	return user.ID > 0
}

// 查询用户：根据用户id查询用户
func (dao *UserDao) FindUserById(user *model.User) (*model.User, error) {
	err := dao.DB.Where("id = ?", user.ID).Limit(1).First(&user)
	return user, err.Error
}

// 查询所有用户
func (dao *UserDao) FindUserList(user *model.User, pageSize int, pageNum int) (count int64, users []*model.User, err error) {
	// 条件查询
	if user.Username != "" || user.Role != 0 {
		err = dao.DB.Model(&user).
			Where("username LIKE ? AND role LIKE ?", "%"+user.Username+"%", user.Role).
			Count(&count).
			Offset((pageNum - 1) * pageSize).Limit(pageSize).
			Find(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return count, users, errors.New("用户不存在")
			}
		}
		return count, users, err
	}
	// 默认查询所有
	err = dao.DB.Model(&user).Count(&count).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Find(&user).Error
	return count, users, err
}

// 新增用户
func (dao *UserDao) CreateUser(user *model.User) error {
	if user != nil {
		return dao.DB.Create(&user).Error
	}
	return errors.New("db create fail:user struct is nil")
}
