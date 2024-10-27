package response

import "bluebell/models"

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"` // 作者
	VoteNum          int64              `json:"vote_num"`    // 投票数
	*models.Post                        // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息
}

// CommunityDetail 社区详情接口的结构体
type CommunityDetail struct {
	CommunityID  string `json:"community_id" db:"community_id"`
	Name         string `json:"community_name" db:"community_name"`
	Introduction string `json:"introduction,omitempty" db:"introduction"`
	CreateTime   int64  `json:"create_at" db:"create_at"`
}
