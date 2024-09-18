package res

import "blue-bell/model"

type ApiPostDetail struct {
	AuthorName       string `json:"author_name" default:"匿名" gorm:"column:author_name;comment:作者名称"`
	VoteNum          int64  `json:"vote_num" gorm:"column:vote_num;comment:投票数" default:"0" `
	*model.Post             // 嵌入帖子内容结构体
	*model.Community        // 嵌入社区分类结构体
}
