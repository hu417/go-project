package config

import "fmt"

// System system configuration 系统配置
type System struct {
	Host string `yaml:"host" mapstructure:"host"` // Host: 地址
	Port int    `yaml:"port" mapstructure:"port"` // Port: 端口
	Env  string `yaml:"env" mapstructure:"env"`   // Env: 环境
}

// Addr 获取地址
func (s System) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
