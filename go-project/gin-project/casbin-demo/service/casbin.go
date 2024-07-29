package service

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type CasbinService struct {
	enforcer *casbin.Enforcer
	adapter  *gormadapter.Adapter
}

func NewCasbinService(enforcer *casbin.Enforcer, adapter *gormadapter.Adapter) *CasbinService {
	return &CasbinService{
		enforcer: enforcer,
		adapter:  adapter,
	}
}
