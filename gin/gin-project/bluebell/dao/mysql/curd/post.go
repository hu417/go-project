package curd

import (
	"context"
	"fmt"

	"bluebell/models/table"

	"gorm.io/gorm"
)

type Post struct {
	*gorm.DB
	context.Context
}

func NewPostDao(ctx context.Context, db *gorm.DB) *Post {
	return &Post{
		Context: ctx,
		DB:      db,
	}

}

// CreatePost 创建帖子
func (d *Post) CreatePost(p *table.Post) (post *table.Post, err error) {
	if err := d.DB.WithContext(d.Context).Model(post).Create(&p).Error; err != nil {
		return nil, fmt.Errorf("[dao] post crate fail => %w", err)
	}
	return p, nil
}

// GetPostById 根据id查询单个贴子数据
func (d *Post) GetPostById(pid int64) (post *table.Post, err error) {

	if err = d.DB.WithContext(d.Context).Model(post).Where("post_id = ?", pid).First(post).Error; err != nil {
		return nil, fmt.Errorf("[dao] post select pid fail => %w", err)
	}
	return post, nil
}

// GetPostList 查询帖子列表函数
func (d *Post) GetPostList(page,size int) (posts []*table.Post, count int64, err error) {
	if err = d.WithContext(d.Context).Model(posts).Find(&posts).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err = d.WithContext(d.Context).Model(posts).Offset(size).Limit(page).Order("id desc").Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

// GetPostListByIDs 根据给定的id列表查询帖子数据
func (d *Post) GetPostListByIDs(ids []string) (posts []*table.Post, count int64, err error) {
	if err = d.WithContext(d.Context).Model(posts).Where("post_id IN ?", ids).Find(&posts).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return posts, count, nil

}

// GetPostListByTitle 根据帖子标题模糊查询帖子数据
func (d *Post) GetPostListByTitles(title string) (posts []*table.Post, count int64, err error) {
	if err = d.WithContext(d.Context).Model(posts).Where("title LIKE ?", "%"+title+"%").Find(&posts).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return posts, count, nil

}
