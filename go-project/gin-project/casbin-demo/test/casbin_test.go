package test

import (
	"testing"

	"casbin-demo/bootstrap"
)

func Test_casbin(t *testing.T) {

	// 初始化数据库
	db := bootstrap.InitDb()
	if db == nil {
		panic("数据库初始化失败")
	}
	defer func() {
		d, _ := db.DB()
		if err := d.Close(); err != nil {
			panic(err)
		}
	}()

	// 初始化casbin
	rbac, err := NewDomainWithGorm(db)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := rbac.enforcer.SavePolicy(); err != nil {
			panic(err)
		}
	}()

	// 关闭自动加载策略，避免在HasPolicy调用中加载策略

	// 添加角色权限
	// ok, err := rbac.enforcer.AddPolicy("admin", "/user", "create")
	// if err != nil {
	// 	panic(err)
	// }
	// if !ok {
	// 	panic("添加失败")
	// }
	// if ok {
	// 	t.Log("添加成功")
	// }

	// 批量添加角色权限
	// policyRules := [][]string{
	// 	{"admin", "/user", "delete"},
	// 	{"admin", "/user", "updata"},
	// 	{"customer", "/user", "create"},
	// 	{"sale", "/user", "100"},
	// 	// 更多规则...
	// }
	// // 方式一: for循环
	// for _, rule := range policyRules {
	// 	// 添加策略
	// 	ok, err := rbac.enforcer.AddPolicy(rule)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if !ok {
	// 		t.Log("添加失败")
	// 	} else {
	// 		t.Log("添加成功")
	// 	}

	// }
	// 方式二: 使用AddPolicies；如果有已存在，则都不添加
	// ok, err := rbac.enforcer.AddPolicies(policyRules)
	// if err != nil {
	// 	panic(err)
	// }
	// if !ok {
	// 	t.Log("添加成功")
	// } else {
	// 	t.Log("添加失败")
	// }

	// 查询角色权限
	// ret, err := rbac.enforcer.GetPermissionsForUser("admin")
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(ret)

	// 删除角色（相关权限、用户全部删除）
	// ok, err := rbac.enforcer.DeleteRole("admin")
	// if err != nil {
	// 	panic(err)
	// }
	// if !ok {
	// 	t.Log("删除失败")
	// } else {
	// 	t.Log("删除成功")
	// }

	// 添加用户角色
	// ok, err := rbac.enforcer.AddGroupingPolicy("wang", "superadmin")
	// if err != nil {
	// 	panic(err)
	// }
	// if !ok {
	// 	t.Log("添加失败")
	// } else {
	// 	t.Log("添加成功")
	// }

	// 批量添加用户角色
	// polices := make([][]string, 0, 2)
	// polices = append(polices, []string{"wangwu", "customer"})
	// polices = append(polices, []string{"xiaoma", "sale"})
	// ok,err := rbac.enforcer.AddGroupingPolicies(polices)
	// if err != nil {
	// 	panic(err)
	// }
	// if !ok {
	// 	t.Log("添加失败")
	// } else {
	// 	t.Log("添加成功")
	// }

	// 查询用户角色
	// ret, err := rbac.enforcer.GetRolesForUser("zhangsan")
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(ret)

	// 用户是否具有某个角色
	// ok, err := rbac.enforcer.HasRoleForUser("zhangsan", "admin")
	// if err != nil {
	// 	panic(err)
	// }
	// if ok {
	// 	t.Log("zhangsan是admin")
	// } else {
	// 	t.Log("zhangsan不是admin")
	// }

	// 删除用户（同删除角色）
	// ok, err := rbac.enforcer.DeleteUser("zhangsan")
	// if err != nil {
	// 	panic(err)
	// }
	// if !ok {
	// 	t.Log("删除成功")
	// } else {
	// 	t.Log("删除失败")
	// }

	// 验证是否具有权限
	ok, err := rbac.enforcer.Enforce("wang", "/user", "4")
	if err != nil {
		panic(err)
	}
	if ok {
		t.Log("zhangsan具有create权限")
	} else {
		t.Log("zhangsan不具有create权限")
	}

}
