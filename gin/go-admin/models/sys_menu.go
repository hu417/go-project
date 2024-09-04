package models

import "gorm.io/gorm"

// SysMenu 菜单结构体
type SysMenu struct {
	gorm.Model
	ParentId      uint   `gorm:"column:parent_id;type:int(11);" json:"parent_id"`
	Name          string `gorm:"column:name;type:varchar(100)" json:"name"`
	WebIcon       string `gorm:"column:web_icon;type:varchar(100);" json:"web_icon"`             // 网页端的图标
	Path          string `gorm:"column:path;type:varchar(255);" json:"path"`                     // 菜单路径
	Sort          int    `gorm:"column:sort;type:int(11);default:0;" json:"sort"`                // 排序规则，默认升序，值越少越靠前
	Level         int    `gorm:"column:level;type:tinyint(1);default:0;" json:"level"`           // 菜单等级，{0：目录，1：菜单，2：按钮}
	ComponentName string `gorm:"column:component_name;type:varchar(100);" json:"component_name"` // 组件名称

}

// TableName 设置数据库表名称
func (table *SysMenu) TableName() string {
	return "sys_menu"
}

// GetMenusList 获取所有菜单列表数据
func GetMenusList() *gorm.DB {
	tx := DB.Model(new(SysMenu)).Select("id,parent_id,name,web_icon,sort,path,level,component_name,created_at,updated_at").Order("sort ASC")
	return tx
}
