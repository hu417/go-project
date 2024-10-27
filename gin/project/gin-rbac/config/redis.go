package config

import "fmt"

// RedisConfig redis configuration redis配置
type Redis struct {
	Host     string `yaml:"host" mapstructure:"host"`         // Host: 地址
	Port     int    `yaml:"port" mapstructure:"port"`         // Port: 端口
	Password string `yaml:"password" mapstructure:"password"` // Password: 密码
	DB       int    `yaml:"db" mapstructure:"db"`             // DB: 数据库
}

func (r *Redis) GetRedisAddress() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
