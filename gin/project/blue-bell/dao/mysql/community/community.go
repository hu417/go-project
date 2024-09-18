package community

import (
	"blue-bell/global"
	"blue-bell/model"
	"fmt"
)

// CheckCommunityExist 判断Community是否存在
func CheckCommunityExist(comm *model.Community) (*model.Community, error) {
	res := global.DB.Where("community_name = ?", comm.Name).First(&comm)
	if res.Error != nil {
		return nil, res.Error
	}

	return comm, nil
}

// InsertUser 向数据库中插入一条用户数据
func Insert(comm *model.Community) error {

	return global.DB.Create(comm).Error
}

// GetCommunityList 获取社区列表
func GetCommunityList(name string, page, pagesize int) (count int64, list []*model.Community, err error) {
	//fmt.Printf("name:%s\n page:%v\n pagesize:%v\n", name,page,pagesize)
	if err := global.DB.Model(&model.Community{}).Where("community_name LIKE ?", "%"+name+"%").Count(&count).Limit(pagesize).Offset((page - 1) * pagesize).Find(&list).Error; err != nil {
		fmt.Printf("err:%v\n", err)
		return 0, nil, err
	}
	return count, list, nil
}


//
func GetCommunityDetailByID(id int64) (data []*model.Community,err error){
	err  = global.DB.Where("id = ?", id).Find(&data).Error
	return
}