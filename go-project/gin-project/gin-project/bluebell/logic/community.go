package logic

import (
	"context"

	"bluebell/dao/mysql"
	"bluebell/dao/mysql/curd"
	"bluebell/models/resp"
	"bluebell/models/table"
)

// GetCommunityList 获取社区列表
func GetCommunityList(ctx context.Context) (interface{}, error) {
	// 查数据库 查找到所有的community 并返回
	communityList, count, err := curd.NewCommunityDao(ctx, mysql.GetDB()).GetCommunityList()
	if err != nil {
		return nil, err
	}

	data := make([]*resp.Community, 0, len(communityList))
	for _, community := range communityList {
		respCommunity := &resp.Community{
			Id:            community.Id,
			CommunityId:   community.CommunityId,
			CommunityName: community.CommunityName,
		}

		data = append(data, respCommunity)
	}

	// 声明匿名结构体并初始化
	return struct {
		Data  []*resp.Community
		Count int64
	}{
		Data:  data,
		Count: count,
	}, nil

}

// GetCommunityDetail 获取社区详情
func GetCommunityDetail(ctx context.Context, id int64) (*table.Community, error) {
	return curd.NewCommunityDao(ctx, mysql.GetDB()).GetCommunityByID(id)
}
