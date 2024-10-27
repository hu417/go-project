package sql

import (
	"bluebell/controller/response"
	"bluebell/global"
	"bluebell/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	err = global.DB.Select("community_id, community_name").Find(&communityList).Error

	return communityList, err
}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(id string) (communitys *response.CommunityDetail, err error) {
	var community models.Community
	err = global.DB.Select("community_id, community_name, introduction, create_at").Where("community_id = ?", id).Find(&community).Error

	communitys = &response.CommunityDetail{
		CommunityID:  community.CommunityID,
		CreateTime:   community.Timestamps.CreatedAt,
		Introduction: community.Introduction,
		Name:         community.CommunityName,
	}
	return communitys, err
}
