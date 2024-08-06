package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
标签			作用
primarykey		主键
unique_index	唯一键
index			键
auto_increment	自增列
column			列名
type			列类型
default			默认值
not null		非空
*/
type User struct {
	//  组合使用  gorm  Model，引用  id、createdAt、updatedAt、deletedAt  等字段
	gorm.Model
	//
	Uid string `gorm:"column:uid;type:varchar(32)"`
	//  列名为  name；列类型字符串；使用该列作为唯一索引
	Name string `gorm:"column:name;type:varchar(15);unique_index"`
	//  该列默认值为  18
	Age int `gorm:"default:18"`
	//  该列值不为空
	Email string `gorm:"not  null"`
	//  该列的数值逐行递增
	Num int `gorm:"auto_increment"`
}

/*
想要明确是显式将字段设置为零值的:
- 使用指针类型
- 使用 sql.Nullxx 类型

*/

func (p *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	//  在创建之前，设置字段的值
	u.Uid = uuid.New().String()
	return
}

// 创建操作前回调
type BeforeCreateInterface interface {
	BeforeCreate(*gorm.DB) error
}

// 创建操作后回调
type AfterCreateInterface interface {
	AfterCreate(*gorm.DB) error
}

// 更新操作前回调
type BeforeUpdateInterface interface {
	BeforeUpdate(*gorm.DB) error
}

// 更新操作后回调
type AfterUpdateInterface interface {
	AfterUpdate(*gorm.DB) error
}

// 保存操作前回调
type BeforeSaveInterface interface {
	BeforeSave(*gorm.DB) error
}

// 保存操作后回调
type AfterSaveInterface interface {
	AfterSave(*gorm.DB) error
}

// 删除操作前回调
type BeforeDeleteInterface interface {
	BeforeDelete(*gorm.DB) error
}

// 删除操作后回调
type AfterDeleteInterface interface {
	AfterDelete(*gorm.DB) error
}

// find    操作后回调
type AfterFindInterface interface {
	AfterFind(*gorm.DB) error
}
