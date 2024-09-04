package curd

import (
	"context"
	"fmt"

	"bluebell/models/table"

	"gorm.io/gorm"
)

type Community struct {
	*gorm.DB
	context.Context
}

func NewCommunityDao(ctx context.Context, db *gorm.DB) *Community {
	return &Community{
		Context: ctx,
		DB:      db,
	}

}

// GetCommunityList 获取社区列表
func (d *Community)GetCommunityList() (communityList []*table.Community, count int64,err error) {
	if err = d.WithContext(d.Context).Model(communityList).Find(&communityList).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return communityList, count, nil
}

// GetCommunityByID 根据ID查询社区详情
func (d *Community)GetCommunityByID(cid int64) (community *table.Community, err error) {
	if err = d.DB.WithContext(d.Context).Model(community).Where("community_id = ?", cid).First(community).Error; err != nil {
		return nil, fmt.Errorf("[dao] community_id select fail => %w", err)
	}
	return community, nil
}


// GetCommunityByName 根据name查询社区详情
func (d *Community)GetCommunityByName(name string) (community *table.Community,ok bool, err error) {
	
	if  n := d.DB.WithContext(d.Context).Model(community).Where("community_name = ?", name).First(community); n.Error != nil {
		return nil, false, fmt.Errorf("[dao] community_name select fail => %w", err)
	} else if n.RowsAffected == 0 {
		// 用户不存在
		return nil, false, nil
	} else {
		return community, true, nil
	}
	
}

