package service

import (
	"casbin-demo/model/req"

	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// 获取所有角色组权限
func (c *CasbinService) GetRolePolicy() (roles []req.RolePolicy, err error) {
	// sql语句
	err = c.adapter.GetDb().Model(&gormadapter.CasbinRule{}).Where("ptype = 'p'").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return
}

// 获取策略中的所有授权规则
func (c *CasbinService) GetPolicy() ([][]string, error) {

	return c.enforcer.GetPolicy()
}

// 获取策略中的所有授权规则，可以指定字段筛选器
func (c *CasbinService) GetFilteredPolicy(fieldIndex int, fieldValues ...string) ([][]string, error) {

	return c.enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
}

// 获取命名策略中的所有授权规则
func (c *CasbinService) GetNamedPolicy(ptype string) ([][]string, error) {

	return c.enforcer.GetNamedPolicy(ptype)
}

// 获取命名策略中的所有授权规则，可以指定字段过滤器。
func (c *CasbinService) GetFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error) {

	return c.enforcer.GetFilteredNamedPolicy(ptype, fieldIndex, fieldValues...)
}

// 创建角色组权限, 已有的会忽略
func (c *CasbinService) CreateRolePolicy(role, url, method string) error {
	// 不直接操作数据库，利用enforcer简化操作
	err := c.enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	// e.AddPolicies(rules) 批量添加，具有原子性
	_, err = c.enforcer.AddPolicy(role, url, method)
	if err != nil {
		return err
	}
	return c.enforcer.SavePolicy()
}

// 修改角色组权限
func (c *CasbinService) UpdateRolePolicy(oldRoleName, oldUrl, oldMethod, newRoleName, newUrl, newMethod string) error {
	_, err := c.enforcer.UpdatePolicy([]string{oldRoleName, oldUrl, oldMethod},
		[]string{newRoleName, newUrl, newMethod})
	if err != nil {
		return err
	}
	return c.enforcer.SavePolicy()
}

// 删除角色组权限
func (c *CasbinService) DeleteRolePolicy(role, url, method string) error {
	_, err := c.enforcer.RemovePolicy(role, url, method)
	if err != nil {
		return err
	}
	return c.enforcer.SavePolicy()
}

// 验证用户/角色权限
func (c *CasbinService) CanAccess(role, url, method string) (ok bool, err error) {
	return c.enforcer.Enforce(role, url, method)
}
