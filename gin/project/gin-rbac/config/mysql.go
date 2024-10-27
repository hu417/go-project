package config

import "fmt"

// MySQLConfig mysql configuration mysql配置
type MySQL struct {
	Host     string `yaml:"host" mapstructure:"host"`         // Host: 地址
	Port     int    `yaml:"port" mapstructure:"port"`         // Port: 端口
	User     string `yaml:"user" mapstructure:"user"`         // User: 用户名
	Password string `yaml:"password" mapstructure:"password"` // Password: 密码
	DBName   string `yaml:"dbname" mapstructure:"dbname"`     // DBName: 数据库名
}

// Dsn 获取dsn
func (m MySQL) Dsn() string {
	// dsn: username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, m.DBName)
}
