package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

/*

对于RBAC模型，v0, v1, v2 等变量通常具有以下含义：

v0：通常代表主体（Subject），它可以是用户或角色。
v1：通常代表对象（Object），这是主体想要访问的资源。
v2：通常代表动作（Action），这是主体想要对对象执行的操作，如读取、写入、删除等。
对于更复杂的模型，如ABAC，v3, v4, v5 可能会用来表示额外的属性，如时间、地点或其他与访问决策相关的属性。

在Casbin的模型配置文件中，p 类型的策略（代表permission）和 g 类型的策略（代表grouping）可能会使用这些变量。

*/

func InitCasbin(db *gorm.DB) (*casbin.Enforcer, *gormadapter.Adapter) {

	// a, err := gormadapter.NewAdapterByDBUseTableName(InitDb(), "casbin_rule", "casbin_rule")
	// if err != nil {
	// 	panic(err)
	// }

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("gormadapter.NewAdapterByDB err: %v", err)
	}

	/*
		 Enforcer 构造器
		// 方式一
		m, err := model.NewModelFromString(`
			[request_definition]
			r = sub, dom, obj

			[policy_definition]
			p = sub, dom, obj

			[role_definition]
			g = _, _, _

			[policy_effect]
			e = some(where (p.eft == allow))

			[matchers]
			m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj || checkSuperAdmin(r.sub, "superadmin")`)

		if err != nil {
			panic(err)
		}
		e, err := casbin.NewEnforcer(m, adapter)
		if err != nil {
			fmt.Printf("casbin NewEnforcer error:%s\n", err.Error())
			return
		}

	*/

	// 方式二
	dir, _ := os.Getwd()
	modelPath := dir + "/gin-project/casbin-demo/config/rbac.conf"
	fmt.Println("modelPath:" + modelPath)

	Enforcer, errC := casbin.NewEnforcer(modelPath, adapter)

	if errC != nil {
		//fmt.Println(errC)
		log.Fatalf("SetupCasbinEnforcer err: %v", errC)
		return nil, nil
	}
	// 加载策略
	Enforcer.LoadPolicy()

	//注册超级管理员权限判断
	Enforcer.AddFunction("checkSuperAdmin", func(arguments ...interface{}) (interface{}, error) {
		username := arguments[0].(string)
		role := arguments[1].(string)
		// 检查用户名的角色是否为超级管理员
		return Enforcer.HasRoleForUser(username, role)
	})

	// 开启日志
	Enforcer.EnableLog(true)
	return Enforcer, adapter

}
