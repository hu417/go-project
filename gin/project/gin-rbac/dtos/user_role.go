package dtos

import (
	"time"

	"gorm.io/gorm"
)

type CreateUserRolesByIDsReqDTO struct {
	UserIDList []uint `json:"user_id_list" binding:"required"` // 用户ID列表
	RoleIDList []uint `json:"role_id_list" binding:"required"` // 角色ID列表
}

type DeleteUserRolesByIDsReqDTO struct {
	UserIDList []uint `json:"user_id_list" binding:"required"` // 用户ID列表
	RoleIDList []uint `json:"role_id_list" binding:"required"` // 角色ID列表
}

type GetUserRolesByUserIDReqDTO struct {
	UserID uint `uri:"userID" binding:"required"` // 用户ID
}

type GetUserRolesByUserIDResDTO struct {
	UserID   uint   `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名
	Roles    []Role `json:"roles"`    // 角色列表
}

type GetUserRolesByRoleIDReqDTO struct {
	RoleID uint `uri:"roleID" binding:"required"` // 角色ID
}

type User struct {
	UserID        uint           `json:"user_id"`         // 用户ID
	Username      string         `json:"username"`        // 用户名
	RelativeAt    time.Time      `json:"relative_at"`     // 关联创建时间
	DelRelativeAt gorm.DeletedAt `json:"del_relative_at"` // 关联删除时间
}

type GetUserRolesByRoleIDResDTO struct {
	RoleID          uint   `json:"role_id"`          // 角色ID
	RoleName        string `json:"role_name"`        // 角色名称
	RoleDescription string `json:"role_description"` // 角色描述
	Users           []User `json:"users"`            // 用户列表
}

type GetUserRoleByIDReqDTO struct {
	UserID uint `uri:"userID" binding:"required" json:"userID"` // 用户ID
	RoleID uint `uri:"roleID" binding:"required" json:"roleID"` // 角色ID
}

type GetUserRoleByIDResDTO struct {
	UserID          uint           `json:"user_id"`          // 用户ID
	Username        string         `json:"username"`         // 用户名
	RoleID          uint           `json:"role_id"`          // 角色ID
	RoleName        string         `json:"role_name"`        // 角色名称
	RoleDescription string         `json:"role_description"` // 角色描述
	CreatedAt       time.Time      `json:"created_at"`       // 创建时间
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`       // 删除时间
}

// GetUserRoleByIDResDTOExample 获取用户角色响应swagger示例
type GetUserRoleByIDResDTOExample struct {
	UserID          uint      `json:"user_id"`          // 用户ID
	Username        string    `json:"username"`         // 用户名
	RoleID          uint      `json:"role_id"`          // 角色ID
	RoleName        string    `json:"role_name"`        // 角色名称
	RoleDescription string    `json:"role_description"` // 角色描述
	CreatedAt       time.Time `json:"created_at"`       // 创建时间
	DeletedAt       time.Time `json:"deleted_at"`       // 删除时间
}
