package req

// 定义请求的参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"` // 代表校验此字段是必须在的
	Password   string `json:"password" binding:"required"` 
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required" `
	Password string `json:"password" binding:"required"`
}

// ParamPost 获取帖子的参数
type ParamPost struct {
	Page  int    `json:"page"`  // 页数
	Size  int    `json:"size"`  // 每页数量
	Order string `json:"order"` // 排序
}
