package req

// 定义请求的参数结构体

// UserSignUp 注册请求参数
type UserSignUp struct {
	Username   string `form:"username" json:"username" binding:"required,min=4,max=20"` // 用户名称,长度为4-20位字符
	Password   string `form:"password" json:"password" binding:"required,password"`     // 必需的，并且应符合一个名为 passwd 的自定义验证规则，密码格式为 6-20 位数字和字母组合
	RePassword string `form:"re_password" json:"re_password" binding:"required,eqfield=Password"`
	// Mobile     string `form:"mobile" json:"mobile,omitempty" binding:"omitempty,required,mobile"` // 必需的，并且是有效的手机号码
}

// GetMessages 自定义错误信息
func (register UserSignUp) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"username.required":   "用户名称不能为空,长度为4-20位字符",
		"password.required":   "用户密码不能为空",
		"password.password":   "密码格式不正确,6-20位数字和大小写字母组合",
		"re_password.eqfield": "两次密码不一致",
		// "email.required":     "邮箱不能为空",
		// "email.email":        "邮箱格式不正确",
		// "mobile.required": "手机号码不能为空",
		// "mobile.mobile":   "手机号码格式不正确",
	}
}

// UserLogin 登录请求参数
type UserLogin struct {
	// Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`       // 必需的，并且是有效的手机号码
	Username string `form:"username" json:"username" binding:"required,min=4,max=20" `
	Password string `form:"password" json:"password" binding:"required,password"` // 必需的，并且应符合一个名为 passwd 的自定义验证规则，密码格式为 6-20 位数字和字母组合
}

// GetMessages 自定义错误信息
func (login UserLogin) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		// "mobile.required":   "用户名称不能为空",
		// "mobile.mobile":     "手机号码格式不正确,11位手机号码",
		"username.required": "用户名称不能为空,长度为4-20位字符",
		"password.required": "密码不能为空",
		"password.password": "密码格式不正确,6-20位数字和大小写字母组合",
	}
}

// 社区
type Community struct {
	Name         string `form:"name" json:"name" binding:"required"`
	Introduction string `form:"introduction" json:"introduction" binding:"required"`
}

// 分页
type Page struct {
	Name  string `json:"name,omitempty"`
	Page  int    `json:"page"`            // 页数
	Size  int    `json:"size"`            // 每页数量
	Order string `json:"order,omitempty"` // 排序
}

// 帖子
type Post struct {
	Title       string `json:"title"`        // 帖子标题
	Content     string `json:"content"`      // 帖子内容
	AuthorID    string  `json:"author_id"`    // 作者ID
	CommunityID string  `json:"community_id"` // 社区ID
}

// ParamPost 获取帖子的参数
type ParamPost struct {
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Order string `json:"order"`
}

// 投票
type Vote struct {
	PostID    int64 `json:"post_id,string" binding:"required"`
	Direction int   `json:"direction"`
}
