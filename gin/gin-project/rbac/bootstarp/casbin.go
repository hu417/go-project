package bootstarp

import (
	"errors"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	PType string `gorm:"size:512;"`
	V0    string `gorm:"size:512;"`
	V1    string `gorm:"size:512;"`
	V2    string `gorm:"size:512;"`
	V3    string `gorm:"size:512;"`
	V4    string `gorm:"size:512;"`
	V5    string `gorm:"size:512;"`
}
//这里是确定该表的表名的部分，具体是什么原理，还未探究
func (CasbinRule) TableName() string {
	return "casbin_rule"
}


type CasbinMethod struct {
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
}

var Casbin *CasbinMethod

// InitCasbinGorm 初始化Casbin Gorm适配器
func InitCasbinGorm(db *gorm.DB) (*CasbinMethod, error) {
	//创建 Gorm适配器
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	// 创建 Casbin Enforcer , 指定自定义的model文件
	e, err := casbin.NewEnforcer("./conf/casbin.conf", a)
	if err != nil {
		return nil, err
	}

	// 开启日志
	e.EnableLog(true)

	// 设置用户root的角色为superadmin
	e.AddRoleForUser("root", "superadmin")
	// 添加自定义函数
	e.AddFunction("checkSuperAdmin", func(arguments ...interface{}) (interface{}, error) {
		// 获取用户名
		username := arguments[0].(string)
		// 检查用户名的角色是否为superadmin
		return e.HasRoleForUser(username, "superadmin")
	})

	if err = e.LoadPolicy(); err != nil {
		return nil, err
	}

	Casbin = &CasbinMethod{
		Enforcer: e,
		Adapter:  a,
	}
	return Casbin, nil
}



/* 授权规则 */
// 增1
func (c *CasbinMethod) AddPolicy(sec string, ptype string, rule []string) error {
	if added, err := c.Enforcer.AddPolicy("linzy", "data1", "read"); err != nil {
		return err
	} else if !added {
		return errors.New("rules are not added")
	}
	return nil
}

// 增2
func (c *CasbinMethod) AddListPolicy(sec string, ptype string, rule []string) error {
	rules := [][]string{
		{"jack", "data4", "read"},
		{"katy", "data4", "write"},
		{"leyo", "data4", "read"},
		{"ham", "data4", "write"},
	}

	if areRulesAdded, err := c.Enforcer.AddPolicies(rules); err != nil {
		return err
	} else if !areRulesAdded {
		return errors.New("rules are not added")
	}

	return nil
}

// 删
func (c *CasbinMethod) DeletePolicy(sec string, ptype string, rule []string) error {
	if isPolicyDeleted, err := c.Enforcer.RemovePolicy(sec, ptype, rule); err != nil {
		return err
	} else if !isPolicyDeleted {
		return errors.New("policy is not deleted")
	}
	return nil
}

// 改
func (c *CasbinMethod) UpdataPolicy(sec string, ptype string, oldRule []string, newRule []string) error {
	if isPolicyUpdated, err := c.Enforcer.UpdatePolicy([]string{"jack", "data4", "read"}, []string{"linzy", "data3", "write"}); err != nil {
		return err
	} else if !isPolicyUpdated {
		return errors.New("policy is not updated")
	}

	return nil
}

// 查
func (c *CasbinMethod) GetPolicy(sec string, ptype string) ([][]string, error) {
	if policies, err := c.Enforcer.GetPolicy(); err != nil {
		return nil, err
	} else {
		return policies, nil
	}
}
