package curd

import (
	"context"
	"errors"
	"fmt"

	"bluebell/models/table"

	"gorm.io/gorm"
)

type User struct {
	*gorm.DB
	context.Context
}

func NewUserDao(ctx context.Context, db *gorm.DB) *User {
	return &User{
		Context: ctx,
		DB:      db,
	}

}

// CheckUserExist 检查指定用户名的用户是否存在
func (d *User) CheckUserExist(username string) (*table.User, bool, error) {
	var count int64
	user := &table.User{}

	if err := d.DB.WithContext(d.Context).Model(user).Where("username = ?", username).Count(&count).First(&user).Error; err == nil && count > 0 {
		// 用户已存在
		return user, true, nil
	} else if count == 0 {
		// 用户不存在
		return nil, false, nil
	} else {
		// 查询数据库失败
		return nil, false, fmt.Errorf("[dao] username select fail => %w", err)
	}

}

// InsertUser 想数据库中插入一条新的用户记录
func (d *User) InsertUser(user *table.User) (err error) {

	// 执行SQL语句入库
	if err := d.DB.WithContext(d.Context).Model(&table.User{}).Create(&user).Error; err != nil {

		return fmt.Errorf("[dao] username create fail => %w", err)
	}

	return nil
}

// GetUserById 根据id获取用户信息
func (d *User) GetUserById(uid int64) (*table.User, error) {
	user := &table.User{}
	if err := d.DB.WithContext(d.Context).Model(user).Where("user_id = ?", uid).Error; err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, fmt.Errorf("[dao] userid select tail => %w", errors.New("user not exist"))
	}
	return user, nil
}
