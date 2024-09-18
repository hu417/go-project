package logic

import (
	"errors"

	"blue-bell/controller/e"
	"blue-bell/controller/req"
	"blue-bell/dao/mysql/community"
	"blue-bell/model"

	"gorm.io/gorm"
)

// CreateCommunity 创建社区
func CreateCommunity(p *req.Community) error {
	// 1.校验社区是否存在
	comm := &model.Community{
		Name:         p.Name,
		Introduction: p.Introduction,
	}
	if _, err := community.CheckCommunityExist(comm); err != nil {
		// 判断是否为社区不存在的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// UID := utils.NanoId()
			// 2.保存到数据库
			return community.Insert(comm)
		}
		return err
	}
	// 3.社区已存在
	return e.ErrorCommunityExist
}

// GetCommunityList 获取社区列表
func GetCommunityList(name string, page, pagesize int) (count int64, list []*model.Community, err error) {
	return community.GetCommunityList(name, page, pagesize)
}


//
func GetCommunityDetailByID(id int64) (data interface{},err error){
	data,err =  community.GetCommunityDetailByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil,e.ErrorCommunityNotExist
		}
		return nil,err
	}
	return 
}