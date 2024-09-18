package logic

import (
	"blue-bell/controller/req"
	"blue-bell/dao/mysql/post"
)

// CreatePost 创建帖子
func CreatePost(p *req.Post) error {
	return post.InseartPost(p)
}

// GetPostList 获取帖子列表
func GetPostList(p *req.Page) (interface{}, error){
	return post.GetPostList(p)

	// return nil,nil
}

// GetPostByID 根据id获取帖子详情
func GetPostByID(id string) (interface{}, error) {
	return post.GetPostById(id)
}