package config

// email 配置
type Email struct {
	Host             string `json:"host" yaml:"host"`
	Port             int    `json:"port" yaml:"port"`
	User             string `json:"user" yaml:"user"` // 发件人邮箱
	Password         string `json:"password" yaml:"password"`
	DefaultFromEmail string `json:"default_from_email" yaml:"default_from_emal"` // 默的发件人名字
	UsessL           bool   `json:"use_ssl" yaml:"use_ssl"`                      // 是否使用ssL
	UserTls          bool   `json:"user_tls" yaml:"user_tis"`                    // 是否加密
}
