package service

// LoginPasswordRequest 登录参数结构体
type LoginPasswordRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// LoginPasswordReply 登录后的token结构体
type LoginPasswordReply struct {
	Token        string `json:"token"`         // token
	RefreshToken string `json:"refresh_token"` // 用于刷新 token 的 token
}

// GetUserListRequest 获取用户列表参数结构体
type GetUserListRequest struct {
	*QueryRequest
}

// QueryRequest 关键字和分页信息结构体
type QueryRequest struct {
	Page    int    `json:"pageIndex" form:"pageIndex"`
	Size    int    `json:"pageSize" form:"pageSize"`
	Keyword string `json:"searchValue" form:"searchValue"`
}

// GetUserListReply 返回管理员信息结构体
type GetUserListReply struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Avatar    string `json:"Avatar"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AddUserRequest 添加管理员结构体
type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Remarks  string `json:"remarks"`
	RoleId   uint   `json:"roleId"`
	Sex      string `json:"sex"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// GetUserDetailReply 获取管理员结构体
type GetUserDetailReply struct {
	ID uint `json:"id"`
	AddUserRequest
}

// UpdateUserRequest 更新管理员信息结构体
type UpdateUserRequest struct {
	ID uint `json:"id"`
	AddUserRequest
}

// GetRoleListRequest 获取角色列表参数结构体
type GetRoleListRequest struct {
	*QueryRequest
}

// GetRoleListReply 返回角色列表结构体
type GetRoleListReply struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Sort      int    `json:"sort"`
	IsAdmin   int    `json:"is_admin"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AddRoleRequest 新增角色参数结构体
type AddRoleRequest struct {
	Name    string `json:"name"`
	Sort    int64  `json:"sort"`
	IsAdmin int8   `json:"isAdmin"`
	Remarks string `json:"remarks"`
	MenuId  []uint `json:"menuId"` // 被授权的菜单列表
}

// GetRoleDetailReply 返回角色详情结构体
type GetRoleDetailReply struct {
	ID uint `json:"id"`
	AddRoleRequest
}

// UpdateRoleRequest 更新角色信息结构体
type UpdateRoleRequest struct {
	ID uint `json:"id"`
	AddRoleRequest
}

// MenuReply 菜单列表返回结构体
type MenuReply struct {
	ID            int          `json:"id"`
	ParentId      int          `json:"parent_id"`
	Name          string       `json:"name"`
	WebIcon       string       `json:"web_icon"`
	Sort          int          `json:"sort"`
	Path          string       `json:"path"`
	Level         int          `json:"level"`          // 菜单等级，{0：目录，1：菜单，2：按钮}
	ComponentName string       `json:"component_name"` // 组件路径
	SubMenus      []*MenuReply `json:"sub_menus"`
}

// AllMenu 所有菜单数据结构体
type AllMenu struct {
	ID            int    `json:"id"`
	ParentId      int    `json:"parent_id"`
	Name          string `json:"name"`
	WebIcon       string `json:"web_icon"`
	Sort          int    `json:"sort"`
	Path          string `json:"path"`
	Level         int    `json:"level"`
	ComponentName string `json:"component_name"` // 组件路径
}

// AddMenuRequest 新增菜单结构体
type AddMenuRequest struct {
	ParentId      uint   `json:"parent_id"`      // 父级唯一标识，不填默认为顶级菜单
	Name          string `json:"name"`           // 菜单名称
	WebIcon       string `json:"web_icon"`       // 网页图标
	Path          string `json:"path"`           // 路径
	Sort          int    `json:"sort"`           // 排序
	Level         int    `json:"level"`          // 菜单等级，{0：目录，1：菜单，2：按钮}
	ComponentName string `json:"component_name"` // 组件路径
}

// UpdateMenuRequest 更新菜单结构体
type UpdateMenuRequest struct {
	ID uint `json:"id"`
	AddMenuRequest
}

// AllListReply 返回角色结构体
type AllListReply struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// UpdatePwdRequest 更信息密码结构体
type UpdatePwdRequest struct {
	UsedPass string `json:"usedPass"`
	NewPass  string `json:"newPass"`
}

// GetLogListRequest 获取日志列表参数结构体
type GetLogListRequest struct {
	*QueryRequest
}

// GetLogListReply 返回日志信息结构体
type GetLogListReply struct {
	ID          uint   `json:"id"`
	Browser     string `json:"browser"`
	ClassMethod string `json:"class_method"`
	UseTime     string `json:"use_time"`
	HttpMethod  string `json:"http_method"`
	StatusCode  string `json:"status_code"`
	Params      string `json:"params"`
	Response    string `json:"response"`
	RemoteAddr  string `json:"remote_addr"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	IpInfo
}

// IpInfo IP归属地详情结构体
type IpInfo struct {

	// 国家
	Country string `json:"country"`
	// 区域
	Region string `json:"region"`
	// 省份
	Province string `json:"province"`
	// 城市
	City string `json:"city"`
	// 运营商
	Isp string `json:"isp"`
}
