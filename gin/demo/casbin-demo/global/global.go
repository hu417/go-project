package global

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var (

	// Enforcer casbin 权限控制
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
)
