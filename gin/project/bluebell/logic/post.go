package logic

import (
	"bluebell/controller/request"
	"bluebell/controller/response"
	"bluebell/dao/mysql/sql"
	"bluebell/dao/redis/keys"
	"bluebell/models"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func CreatePost(ctx context.Context, p *models.Post) (err error) {
	// 1. 生成post id
	p.PostID = uuid.New().String()
	// 2. 保存到数据库
	err = sql.CreatePost(p)
	if err != nil {
		return err
	}
	err = keys.CreatePost(ctx, p.PostID, p.CommunityID)
	// 3. 返回
	return

}

// GetPostById 根据帖子id查询帖子详情数据
func GetPostById(pid int64) (data *response.ApiPostDetail, err error) {
	// 查询并组合我们接口想用的数据
	post, err := sql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := sql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.String("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	community, err := sql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.String("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	// 接口数据拼接
	data = &response.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (total int64, data []*response.ApiPostDetail, err error) {
	total, posts, err := sql.GetPostList(page, size)
	if err != nil {
		return 0, nil, err
	}
	data = make([]*response.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := sql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.String("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := sql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.String("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &response.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(ctx context.Context, p *request.ParamPostList) (data []*response.ApiPostDetail, err error) {
	// 2. 去redis查询id列表
	ids, err := keys.GetPostIDsInOrder(ctx, p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("keys.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 3. 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回
	posts, err := sql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	// 提前查询好每篇帖子的投票数
	voteData, err := keys.GetPostVoteData(ctx, ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := sql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.String("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := sql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.String("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &response.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return

}

func GetCommunityPostList(ctx context.Context, p *request.ParamPostList) (data []*response.ApiPostDetail, err error) {
	// 2. 去redis查询id列表
	ids, err := keys.GetCommunityPostIDsInOrder(ctx, p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("keys.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetCommunityPostIDsInOrder", zap.Any("ids", ids))
	// 3. 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回
	posts, err := sql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	// 提前查询好每篇帖子的投票数
	voteData, err := keys.GetPostVoteData(ctx, ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := sql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.String("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := sql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.String("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &response.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew  将两个查询帖子列表逻辑合二为一的函数
func GetPostListNew(ctx context.Context, p *request.ParamPostList) (data []*response.ApiPostDetail, err error) {
	// 根据请求参数的不同，执行不同的逻辑。
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList2(ctx, p)
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList(ctx, p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
