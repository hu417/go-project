package service

// 用户以及组关系的增删改查 //

type User struct {
	UserName  string
	RoleNames []string
}

// 获取所有用户以及关联的角色
func (c *CasbinService) GetUsers() (users []User) {
	p, err := c.enforcer.GetGroupingPolicy()
	if err != nil {
		return nil
	}
	usernameUser := make(map[string]*User, 0)
	for _, _p := range p {
		username, usergroup := _p[0], _p[1]
		if v, ok := usernameUser[username]; ok {
			usernameUser[username].RoleNames = append(v.RoleNames, usergroup)
		} else {
			usernameUser[username] = &User{UserName: username, RoleNames: []string{usergroup}}
		}
	}
	for _, v := range usernameUser {
		users = append(users, *v)
	}
	return
}

// 获取所有角色组
func (c *CasbinService) GetAllRoles() ([]string, error) {
	return c.enforcer.GetAllRoles()
}

// 获取角色下的用户
func (c *CasbinService) GetUsersForRole(rolename string) ([]string, error) {
	return c.enforcer.GetUsersForRole(rolename)
}

// 角色组中添加用户, 没有组默认创建
func (c *CasbinService) AddUserRole(username, rolename string) error {
	_, err := c.enforcer.AddGroupingPolicy(username, rolename)
	if err != nil {
		return err
	}
	return c.enforcer.SavePolicy()
}

// 查询用户角色
func (c *CasbinService) GetUserRoles(username string) ([]string, error) {
	return c.enforcer.GetRolesForUser(username)
}

// 判断用户是否具有某个角色
func (c *CasbinService) HasRoleForUser(username, rolename string) (bool, error) {
	return c.enforcer.HasGroupingPolicy(username, rolename)
}

// 更改用户所属角色
func (c *CasbinService) UpdateUserRole(username, oldRole, newRole string) error {
	_, err := c.enforcer.RemoveGroupingPolicy(username, oldRole)
	if err != nil {
		return err
	}
	_, err = c.enforcer.AddGroupingPolicy(username, newRole)
	if err != nil {
		return err
	}
	return c.enforcer.SavePolicy()
}

// 角色组中删除用户
func (c *CasbinService) DeleteUserRole(username, rolename string) error {
	_, err := c.enforcer.RemoveGroupingPolicy(username, rolename)
	if err != nil {
		return err
	}
	return c.enforcer.SavePolicy()
}
