package v1

// Review 评价数据
type Review struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userID"`
	Score       int    `json:"score"`
	Content     string `json:"content"`
	Tags        []Tag  `json:"tags"`
	Status      int    `json:"status"`
	PublishTime int64  `json:"publishDate"`
}

// Tag 评价标签
type Tag struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
}
