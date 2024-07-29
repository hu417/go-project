package test

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"

	// casbinmodel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// RbacDomain 支持多业务(domain)的casbin包
type RbacDomain struct {
	enforcer *casbin.Enforcer
}

func NewDomainWithGorm(db *gorm.DB) (ret *RbacDomain, err error) {

	// orm, err := gormadapter.NewAdapter("mysql", dbConn, true)
	adapter, err := gormadapter.NewAdapterByDB(db)

	if err != nil {
		fmt.Printf("gormadapter init error:%s\n", err.Error())
		return
	}

	m, err := model.NewModelFromString(`
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act  || checkSuperAdmin(r.sub, "superadmin")`)

	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer(m, adapter, true)
	if err != nil {
		fmt.Printf("casbin NewEnforcer error:%s\n", err.Error())
		return
	}
	// 注册超级管理员权限判断
	e.AddFunction("checkSuperAdmin", func(arguments ...interface{}) (interface{}, error) {
		username := arguments[0].(string)
		role := arguments[1].(string)
		// 检查用户名的角色是否为超级管理员
		return e.HasRoleForUser(username, role)
	})

	// 加载策略
	e.LoadPolicy()

	// e.EnableLog(true)

	ret = &RbacDomain{
		enforcer: e,
	}
	return
}
