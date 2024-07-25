package service

import "github.com/casbin/casbin/v2"

// RbacDomain 支持多业务(domain)的casbin包
type RbacDomain struct {
	enforcer *casbin.Enforcer
	domain   string
}

func NewRbacDomain(enforcer *casbin.Enforcer, domain string) *RbacDomain {
	return &RbacDomain{
		enforcer: enforcer,
		domain:   domain,
	}
}

// AddRoleForUser 添加用户角色
func (c *RbacDomain) AddRoleForUser(username, role string) (bool, error) {
	return c.enforcer.AddRoleForUser(username, role, c.domain)
}

// DeleteRoleForUser 删除用户角色
func (c *RbacDomain) DeleteRoleForUser(username, role string) (bool, error) {
	return c.enforcer.DeleteRoleForUser(username, role, c.domain)
}

// GetUserRoles 获取用户角色
func (c *RbacDomain) GetUserRoles(username string) ([]string, error) {
	return c.enforcer.GetRolesForUser(username, c.domain)
}

// GetRoleUsers 获取角色下的用户
func (c *RbacDomain) GetRoleUsers(role string) ([]string, error) {
	return c.enforcer.GetUsersForRole(role, c.domain)
}

// AddRoleForUserMulti 添加用户角色-批量
func (c *RbacDomain) AddRoleForUserMulti(username string, roles []string) (bool, error) {
	return c.enforcer.AddGroupingPolicies(c.formatMulti(username, roles))
}

// DeleteRoleForUserMulti 删除用户角色-批量
func (c *RbacDomain) DeleteRoleForUserMulti(username string, roles []string) (bool, error) {
	return c.enforcer.RemoveGroupingPolicies(c.formatMulti(username, roles))
}

// UpdateRoleForUserMulti 批量更新用户角色
func (c *RbacDomain) UpdateRoleForUserMulti(username string, roles []string) (bool, error) {
	//删除旧角色
	_, err := c.enforcer.DeleteRolesForUser(username, c.domain)
	if err != nil {
		return false, err
	}
	//添加新角色
	return c.enforcer.AddGroupingPolicies(c.formatMulti(username, roles))
}

// DeleteUser 删除用户：用户的权限一并删除
func (c *RbacDomain) DeleteUser(user string) (bool, error) {
	var err error
	ret1, err := c.enforcer.RemoveFilteredGroupingPolicy(0, user, "", c.domain)
	if err != nil {
		return ret1, err
	}
	ret2, err := c.enforcer.RemoveFilteredPolicy(0, user, "", c.domain)
	return ret1 || ret2, err
}

// DeleteRole 删除角色：角色下的用户和权限一并删除
func (c *RbacDomain) DeleteRole(role string) (bool, error) {
	var err error
	ret1, err := c.enforcer.RemoveFilteredGroupingPolicy(1, role, c.domain)
	if err != nil {
		return ret1, err
	}
	ret2, err := c.enforcer.RemoveFilteredPolicy(0, role, "", c.domain)
	return ret1 || ret2, err
}

// AddPermission 添加角色or用户权限
func (c *RbacDomain) AddPermission(username, permission string) (bool, error) {
	return c.enforcer.AddPermissionForUser(username, permission, c.domain)
}

// AddPermissionMulti 添加角色or用户权限-批量
func (c *RbacDomain) AddPermissionMulti(username string, permissions []string) (bool, error) {
	return c.enforcer.AddPolicies(c.formatMulti(username, permissions))
}

// DeletePermission 删除角色or用户权限
func (c *RbacDomain) DeletePermission(username, permission string) (bool, error) {
	return c.enforcer.DeletePermissionForUser(username, permission, c.domain)
}

// DeletePermissionMulti 删除角色or用户权限-批量
func (c *RbacDomain) DeletePermissionMulti(username string, permissions []string) (bool, error) {
	return c.enforcer.RemovePolicies(c.formatMulti(username, permissions))
}

// GetPermissionsForRole 获取角色权限
func (c *RbacDomain) GetPermissionsForRole(role string) []string {
	list, err := c.enforcer.GetFilteredNamedPolicy("p", 0, role, "", c.domain)
	if err != nil {
		return nil
	}
	ret := make([]string, 0, len(list))
	for _, v := range list {
		ret = append(ret, v[1]) //TODO 权限目前取第二位
	}
	return ret
}

// UpdatePermissionsForRoleMulti 更新角色权限-批量
func (c *RbacDomain) UpdatePermissionsForRoleMulti(role string, polices []string) (ret bool, err error) {
	//删除角色所有权限
	ret, err = c.enforcer.RemoveFilteredPolicy(0, role, "", c.domain)
	if err != nil {
		return
	}
	//添加角色权限
	ret, err = c.enforcer.AddPolicies(c.formatMulti(role, polices))
	return
}

// RemovePolice 删除指定权限
func (c *RbacDomain) RemovePolice(police string) (bool, error) {
	return c.enforcer.RemoveFilteredPolicy(1, police, c.domain)
}

// HasPermission 用户or角色是否具有某权限
func (c *RbacDomain) HasPermission(user, permission string) (bool, error) {
	return c.enforcer.Enforce(user, permission, c.domain)
}

// GetEnforcer 返回 *casbin.Enforcer，用来执行casbin原生方法
func (c *RbacDomain) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}

func (c *RbacDomain) formatMulti(username string, polices []string) [][]string {
	policeArr := make([][]string, 0, len(polices))
	for _, v := range polices {
		policeArr = append(policeArr, []string{username, v, c.domain})
	}
	return policeArr
}
