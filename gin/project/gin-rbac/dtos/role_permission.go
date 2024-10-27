package dtos

import (
	"time"

	"gorm.io/gorm"
)

type CreateRolePermissionsByIDsReqDTO struct {
	RoleIDList       []uint `json:"role_id_list" binding:"required"`       // 角色ID列表
	PermissionIDList []uint `json:"permission_id_list" binding:"required"` // 权限ID列表
}

type DeleteRolePermissionsByIDsReqDTO struct {
	RoleIDList       []uint `json:"role_id_list" binding:"required"`       // 角色ID列表
	PermissionIDList []uint `json:"permission_id_list" binding:"required"` // 权限ID列表
}

type GetRolePermissionsByPermissionIDReqDTO struct {
	PermissionID uint `uri:"permissionID" binding:"required"` // 权限ID
}

type Role struct {
	RoleID          uint           `json:"role_id"`          // 角色ID
	RoleName        string         `json:"role_name"`        // 角色名称
	RoleDescription string         `json:"role_description"` // 角色描述
	RelativeAt      time.Time      `json:"relative_at"`      // 关联创建时间
	DelRelativeAt   gorm.DeletedAt `json:"del_relative_at"`  // 关联删除时间
}

type GetRolePermissionsByPermissionIDResDTO struct {
	PermissionID          uint   `json:"permission_id"`          // 权限ID
	PermissionName        string `json:"permission_name"`        // 权限名称
	PermissionDescription string `json:"permission_description"` // 权限描述
	PermissionAction      string `json:"permission_action"`      // 权限动作
	Roles                 []Role `json:"roles"`                  // 角色列表
}

type GetRolePermissionsByRoleIDReqDTO struct {
	RoleID uint `uri:"roleID" binding:"required"` // 角色ID
}

type Permission struct {
	PermissionID          uint           `json:"permission_id"`          // 权限ID
	PermissionName        string         `json:"permission_name"`        // 权限名称
	PermissionDescription string         `json:"permission_description"` // 权限描述
	PermissionAction      string         `json:"permission_action"`      // 权限动作
	RelativeAt            time.Time      `json:"relative_at"`            // 关联创建时间
	DelRelativeAt         gorm.DeletedAt `json:"del_relative_at"`        // 关联删除时间
}

type GetRolePermissionsByRoleIDResDTO struct {
	RoleID          uint         `json:"role_id"`          // 角色ID
	RoleName        string       `json:"role_name"`        // 角色名称
	RoleDescription string       `json:"role_description"` // 角色描述
	Permissions     []Permission `json:"permissions"`      // 权限列表
}

type GetRolePermissionByIDReqDTO struct {
	RoleID       uint `uri:"roleID" binding:"required"`       // 角色ID
	PermissionID uint `uri:"permissionID" binding:"required"` // 权限ID
}

type GetRolePermissionResDTO struct {
	RoleID                uint           `json:"role_id"`                // 角色ID
	RoleName              string         `json:"role_name"`              // 角色名称
	RoleDescription       string         `json:"role_description"`       // 角色描述
	PermissionID          uint           `json:"permission_id"`          // 权限ID
	PermissionName        string         `json:"permission_name"`        // 权限名称
	PermissionDescription string         `json:"permission_description"` // 权限描述
	PermissionAction      string         `json:"permission_action"`      // 权限动作
	CreatedAt             time.Time      `json:"created_at"`             // 创建时间
	DeletedAt             gorm.DeletedAt `json:"deleted_at"`             // 删除时间
}

// GetRolePermissionResDTOExample 获取用户角色响应体swagger示例
type GetRolePermissionResDTOExample struct {
	RoleID                uint      `json:"role_id"`                // 角色ID
	RoleName              string    `json:"role_name"`              // 角色名称
	PermissionID          uint      `json:"permission_id"`          // 权限ID
	PermissionName        string    `json:"permission_name"`        // 权限名称
	PermissionDescription string    `json:"permission_description"` // 权限描述
	PermissionAction      string    `json:"permission_action"`      // 权限动作
	CreatedAt             time.Time `json:"created_at"`             // 创建时间
	DeletedAt             time.Time `json:"deleted_at"`             // 删除时间
}
