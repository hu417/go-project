package db

import (
	"go-admin/global"
	"go-admin/model"

	"gorm.io/gorm"
)

// GetLogList 获取日志列表
func GetLogList(keyword string) *gorm.DB {
	tx := global.DB.Model(&model.SysLog{}).
		Select("id,browser,class_method,http_method,remote_addr,status_code,params,response,use_time,country,region,province,city,isp,created_at,updated_at")
	if keyword != "" {
		tx.Where("title LIKE ?", "%"+keyword+"%")
	}
	return tx
}
