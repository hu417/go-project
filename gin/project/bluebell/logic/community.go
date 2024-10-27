package logic

import (
	"bluebell/controller/response"
	"bluebell/dao/mysql/sql"
	"bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库 查找到所有的community 并返回
	return sql.GetCommunityList()
}

func GetCommunityDetail(id string) (*response.CommunityDetail, error) {
	return sql.GetCommunityDetailByID(id)
}
