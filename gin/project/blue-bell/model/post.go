package model

type PostContent struct {
	BaseID
	PostID      int64  `json:"post_id" db:"post_id" gorm:"column:post_id;primaryKey;autoIncrement;"`
	Title       string `json:"title" db:"title" binding:"required" gorm:"column:title;index;size:128;not null;comment:帖子标题"`
	Content     string `json:"content" db:"content" binding:"required" gorm:"column:content;not null;comment:帖子内容"`
	AuthorID    int64  `json:"author_id" db:"author_id" gorm:"column:author_id;index;not null;comment:作者ID"`
	CommunityID int64  `json:"community_id" db:"community_id" binding:"required" gorm:"column:community_id;index;not null;comment:所属社区"`
	Timestamps
}

type ApiPostDetail struct {
	AuthorName               string `json:"author_name" default:"匿名" gorm:"column:author_name;comment:作者名称"`
	VoteNum                  int64  `json:"vote_num" gorm:"column:vote_num;comment:投票数" default:"0" `
	*PostContent                    // 嵌入帖子内容结构体
	*CommunityCategoryDetail        // 嵌入社区分类结构体
}
