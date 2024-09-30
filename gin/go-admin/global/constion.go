package global

import "time"

var (
	// StaticResource 静态文件目录
	//StaticResource = "E:\\go-admin-static"
	//DbPath = StaticResource + "\\ip2region.xdb"
	StaticResource = "./tmp"
	DbPath         = StaticResource + "/ip2region.xdb"
	// DefaultSize 默认每页查询20条数据
	DefaultSize = 10
	// JwtKey 密钥（建议修改）
	JwtKey = "sys-admin"
	// TokenExpire token 有效期，7天
	TokenExpire = time.Now().Add(time.Second * 3600 * 24 * 7).Unix()
	// RefreshTokenExpire 刷新 token 有效期，14天
	RefreshTokenExpire = time.Now().Add(time.Second * 3600 * 24 * 14).Unix()
	// EmailFrom 邮件发送方
	EmailFrom = "379533177@qq.com"
	// EmailPassWord 邮箱授权码
	EmailPassWord = "ryjxbuztvacacahj"
	// EmailHost 邮箱Host
	EmailHost = "smtp.qq.com"
	// EmailPort 邮箱端口号
	EmailPort = "587"
)
