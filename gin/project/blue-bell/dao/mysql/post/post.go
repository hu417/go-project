package post

import (
	"blue-bell/controller/e"
	"blue-bell/controller/req"
	"blue-bell/controller/res"
	"blue-bell/global"
	"blue-bell/model"
)

func InseartPost(p *req.Post) error {

	// 验证用户是否存在
	var user model.User
	if err := global.DB.Where("user_id = ?", p.AuthorID).First(&user).Error; err != nil {
		return err
	}
	if user.ID == 0 {
		return e.ErrorUserNotExist
	}

	// 验证社区是否存在
	var community model.Community
	if err := global.DB.Where("community_id = ?", p.CommunityID).First(&community).Error; err != nil {
		return err
	}
	if community.ID == 0 {
		return e.ErrorCommunityNotExist
	}

	// 插入帖子数据
	post := &model.Post{
		AuthorID:    p.AuthorID,
		CommunityID: p.CommunityID,
		Content:     p.Content,
		Title:       p.Title,
	}
	if err := global.DB.Create(post).Error; err != nil {
		return err
	}

	return nil
}

func GetPostList(p *req.Page) (interface{}, error) {
	var data []interface{}
	var count int64
	var posts []model.Post
	// 获取帖子列表
	if err := global.DB.Model(&model.Post{}).Count(&count).Limit(p.Size).Offset((p.Page - 1) * p.Size).Find(&posts).Error; err != nil {
		return nil, err
	}

	// 根据帖子信息查询作者信息, 社区信息
	for _, post := range posts {
		var user model.User
		if err := global.DB.Where("user_id = ?", post.AuthorID).First(&user).Error; err != nil {
			return nil, err
		}
		var community model.Community
		if err := global.DB.Where("community_id = ?", post.CommunityID).First(&community).Error; err != nil {
			return nil, err
		}
		data = append(data, &res.ApiPostDetail{
			AuthorName: user.UserName,
			Community:  &community,
			Post:       &post,
		})
	}

	datas := struct {
		List  []interface{} `json:"list"`
		Page  *req.Page
		Total int64 `json:"total"`
	}{
		List:  data,
		Page:  p,
		Total: count,
	}

	return datas, nil
}

func GetPostById(id string) (interface{}, error) {

	var post model.Post
	if err := global.DB.Where("post_id = ?", id).Find(&post).Error; err != nil {
		return nil, err
	}

	var user model.User
	if err := global.DB.Where("user_id = ?", post.AuthorID).First(&user).Error; err != nil {
		return nil, err
	}
	var community model.Community
	if err := global.DB.Where("community_id = ?", post.CommunityID).First(&community).Error; err != nil {
		return nil, err
	}

	data := struct {
		AuthorName string 
		Community string
		Post model.Post
	}{
		AuthorName: user.UserName,
		Community: community.Name,
		Post: post,
	}

	return data, nil
}
