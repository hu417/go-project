package resp

import "bluebell/models/table"

// login
type ParamToken struct {
	UserId   int64 `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// PostDetail 帖子详情接口的结构体
type PostDetail struct {
	CommunityName string `json:"community_name"` // 嵌入社区信息
	AuthorName    string `json:"author_name"`    // 作者
	VoteNum       int64  `json:"vote_num"`       // 投票数
	*table.Post          // 嵌入帖子结构体
}

// Community
type Community struct {
	Id            int    `json:"id"`
	CommunityId   int64 `json:"community_id"`
	CommunityName string `json:"community_name"`
}
