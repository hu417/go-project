package bootstrap

import (
	"fmt"
	
	"gin-rbac/config"
)

type System struct {
	system *config.System
}

func NewSystem(sys *config.System) *System {
	return &System{system: sys}
}

// InitSystem 初始化系统配置
func InitSystem(cfg *config.Config) (*System, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	if cfg.System == nil {
		return nil, fmt.Errorf("system configuration cannot be nil")
	}

	system := NewSystem(cfg.System)
	return system, nil
}

// Addr 获取系统监听地址
func (s *System) Addr() string {
	if s.system != nil {
		return s.system.Addr()
	}
	return ""
}
