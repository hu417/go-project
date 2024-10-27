package sql

import (
	"bluebell/global"
	"bluebell/models"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	// insert into post(post_id, title, content, author_id, community_id) values (?, ?, ?, ?, ?)
	return global.DB.Create(&p).Error
}

// GetPostById 根据id查询单个贴子数据
func GetPostById(pid int64) (post *models.Post, err error) {

	err = global.DB.Select("post_id, title, content, author_id, community_id, create_at").Where("post_id = ?", pid).Find(&post).Error

	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (total int64, posts []*models.Post, err error) {

	err = global.DB.Select("post_id, title, content, author_id, community_id, create_at").
		Count(&total).
		Order("create_at desc").
		Limit(int(size)).
		Offset(int((page - 1) * size)).
		Find(&posts).Error

	return
}

// GetPostListByIDs 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {

	err = global.DB.Select("post_id, title, content, author_id, community_id, create_at").
		Where("post_id IN ?", ids).
		Order("post_id desc").
		Find(&postList).Error

	return
}
