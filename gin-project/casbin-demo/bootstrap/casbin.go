package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func InitCasbin(db *gorm.DB,domain string) *casbin.Enforcer {

	// a, err := gormadapter.NewAdapterByDBUseTableName(InitDb(), "casbin_rule", "casbin_rule")
	// if err != nil {
	// 	panic(err)
	// }
	// e, err := casbin.NewEnforcer("./rbac_model.conf", a)
	// if err != nil {
	// 	panic(err)
	// }
	// // 加载策略
	// e.LoadPolicy()
	// return e

	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("gormadapter.NewAdapterByDB err: %v", err)
	}

	dir, _ := os.Getwd()
	modelPath := dir + "/config/rbac_model.conf"
	fmt.Println("modelPath:" + modelPath)

	var errC error
	Enforcer, errC := casbin.NewEnforcer(modelPath, a)
	if errC != nil {
		//fmt.Println(errC)
		log.Fatalf("SetupCasbinEnforcer err: %v", errC)
		return nil
	}
	// 加载策略
	Enforcer.LoadPolicy()

	//注册超级管理员权限判断
	Enforcer.AddFunction("checkSuperAdmin", func(arguments ...interface{}) (interface{}, error) {
		username := arguments[0].(string)
		role := arguments[1].(string)
		// 检查用户名的角色是否为超级管理员
		return Enforcer.HasRoleForUser(username, role, domain)
	})

	// 开启日志
	Enforcer.EnableLog(true)
	return Enforcer

}
