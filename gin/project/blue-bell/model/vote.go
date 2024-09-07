package model

type Vote struct {
	BaseID
	PostID    int64 `json:"post_id,string" binding:"required" gorm:"column:post_id;index;comment:帖子ID"`
	Direction int   `json:"direction" binding:"required,oneof=1 0" gorm:"column:direction;comment:投票类型(0:赞成票(默认),1:反对票)"`
	Timestamps
}
