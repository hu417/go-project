package logic

import (
	"context"
	"errors"
	"fmt"

	"bluebell/dao/mysql"
	"bluebell/dao/mysql/curd"
	"bluebell/dao/redis"
	"bluebell/models/req"
	"bluebell/models/resp"
	"bluebell/models/table"

	"go.uber.org/zap"
)

// 创建帖子
func CreatePost(ctx context.Context, p *req.ParamPost) (err error) {
	// 1. 根据用户查询id
	// 判断用户存不存在
	u, ok, err := curd.NewUserDao(ctx, mysql.GetDB()).CheckUserExist(p.Username)
	if err != nil {
		return fmt.Errorf("[svc] user select fail => %w", err)
	}
	if !ok {
		return errors.New("用户不存在")
	}
	// 根据社区查询id
	c, ok, err := curd.NewCommunityDao(ctx, mysql.GetDB()).GetCommunityByName(p.Community)
	if err != nil {
		return fmt.Errorf("[svc] Community select fail => %w", err)
	}
	if !ok {
		return errors.New("社区不存在")
	}

	// 2. 保存到数据库
	post := &table.Post{
		Status:      p.Status,
		Title:       p.Title,
		Content:     p.Content,
		CommunityId: c.CommunityId,
		AuthorId:    u.UserId,
	}
	posts, err := curd.NewPostDao(ctx, mysql.GetDB()).CreatePost(post)
	if err != nil {
		return fmt.Errorf("[svc] post insert fail => %w", err)
	}

	// 3. 缓存同步
	err = redis.CreatePost(posts.PostId, c.CommunityId)
	return

}

// GetPostById 根据帖子id查询帖子详情数据
func GetPostById(ctx context.Context, pid int64) (data *resp.PostDetail, err error) {
	// 查询并组合我们接口想用的数据
	post, err := curd.NewPostDao(ctx, mysql.GetDB()).GetPostById(pid)
	if err != nil {
		return nil, fmt.Errorf("[svc] postdetail select fail => %w", err)
	}

	// 根据作者id查询作者信息
	user, err := curd.NewUserDao(ctx, mysql.GetDB()).GetUserById(post.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("[svc] user select fail => %w", err)
	}

	// 根据社区id查询社区详细信息
	community, err := curd.NewCommunityDao(ctx, mysql.GetDB()).GetCommunityByID(post.CommunityId)
	if err != nil {
		return nil, fmt.Errorf("[svc] community select fail => %w", err)
	}

	// 接口数据拼接
	data = &resp.PostDetail{
		CommunityName: community.CommunityName,
		AuthorName:    user.Username,
		Post:          post,
	}
	return data, nil
}

// GetPostList 获取帖子列表
func GetPostList(ctx context.Context, p *req.ParamPage) (data interface{}, err error) {
	page := p.Page
	if page == 0 {
		page = 1
	}
	pageSize := p.PageSize
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	posts, count, err := curd.NewPostDao(ctx, mysql.GetDB()).GetPostList(pageSize, offset)
	if err != nil {
		return nil, err
	}
	postList := make([]*resp.PostDetail, 0, len(posts))

	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := curd.NewUserDao(ctx, mysql.GetDB()).GetUserById(post.AuthorId)
		if err != nil {
			// return nil, fmt.Errorf("[svc] user select fail => %w", err)
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := curd.NewCommunityDao(ctx, mysql.GetDB()).GetCommunityByID(post.CommunityId)
		if err != nil {
			// return nil, fmt.Errorf("[svc] community select fail => %w", err)
			continue
		}
		postDetail := &resp.PostDetail{
			CommunityName: community.CommunityName,
			AuthorName:    user.Username,
			Post:          post,
		}
		postList = append(postList, postDetail)
	}

	// 声明匿名结构体并初始化
	data = struct {
		PostList []*resp.PostDetail
		Count    int64
	}{
		PostList: postList,
		Count:    count,
	}
	return data, nil
}

func GetPostList2(ctx context.Context, p *req.ParamPostList) (data interface{}, err error) {
	// 1. 去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return nil, fmt.Errorf("[svc] redis GetPostID fail => %w", err)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("[svc] redis GetPostID fail => %w", errors.New("return 0 data"))
	}

	// 2. 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回
	posts, count, err := curd.NewPostDao(ctx, mysql.GetDB()).GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	postList := make([]*resp.PostDetail, 0, len(posts))
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := curd.NewUserDao(ctx, mysql.GetDB()).GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Sugar().Errorf("[svc] user select fail => %w", err)
			// return nil, fmt.Errorf("[svc] user select fail => %w", err)
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := curd.NewCommunityDao(ctx, mysql.GetDB()).GetCommunityByID(post.CommunityId)
		if err != nil {
			zap.L().Sugar().Errorf("[svc] community select fail => %w", err)
			// return nil, fmt.Errorf("[svc] community select fail => %w", err)
			continue
		}
		postDetail := &resp.PostDetail{
			CommunityName: community.CommunityName,
			AuthorName:    user.Username,
			VoteNum:       voteData[idx],
			Post:          post,
		}

		postList = append(postList, postDetail)
	}

	// 声明匿名结构体并初始化
	data = struct {
		Count    int64
		PostList []*resp.PostDetail
	}{
		Count:    count,
		PostList: postList,
	}
	return data, nil

}

// GetCommunityPostList 获取社区帖子列表
func GetCommunityPostList(ctx context.Context, p *req.ParamPostList) (data interface{}, err error) {
	// 1. 去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return nil, fmt.Errorf("[svc] redis GetPostID fail => %w", errors.New("return 0 data"))
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 2. 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回
	posts, count, err := curd.NewPostDao(ctx, mysql.GetDB()).GetPostListByIDs(ids)
	if err != nil {
		return nil, fmt.Errorf("[svc] mysql GetPostListByIDs fail => %w", err)
	}

	postList := make([]*resp.PostDetail, 0, len(posts))
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := curd.NewUserDao(ctx, mysql.GetDB()).GetUserById(post.AuthorId)
		if err != nil {
			// return nil, fmt.Errorf("[svc] user select fail => %w", err)
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := curd.NewCommunityDao(ctx, mysql.GetDB()).GetCommunityByID(post.CommunityId)
		if err != nil {
			// return nil, fmt.Errorf("[svc] community select fail => %w", err)
			continue
		}
		postDetail := &resp.PostDetail{
			CommunityName: community.CommunityName,
			AuthorName:    user.Username,
			VoteNum:       voteData[idx],
			Post:          post,
		}

		postList = append(postList, postDetail)
	}

	// 声明匿名结构体并初始化
	data = struct {
		Count    int64
		PostList []*resp.PostDetail
	}{
		Count:    count,
		PostList: postList,
	}
	return data, nil

}

// GetPostListNew  将两个查询帖子列表逻辑合二为一的函数
func GetPostListNew(ctx context.Context, p *req.ParamPostList) (data interface{}, err error) {
	// 根据请求参数的不同，执行不同的逻辑。
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList2(ctx, p)
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList(ctx, p)
	}
	if err != nil {
		return nil, fmt.Errorf("[svc] GetPostListNew fail => %w", err)
	}
	return
}
