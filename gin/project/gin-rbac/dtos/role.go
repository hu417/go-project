package dtos

import (
	"time"

	"gorm.io/gorm"
)

type GetRoleListReqDTO struct {
	PaginationReqDTO
	Name string `form:"name" binding:"omitempty,max=20"` // 角色名称
}

type CreateRoleReqDTO struct {
	Name        string `json:"name" binding:"required,max=20"`          // 角色名称
	Description string `json:"description" binding:"omitempty,max=255"` // 角色描述
}

type GetRoleReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 角色ID
}

type GetRoleResDTO struct {
	ID          uint           `json:"id"`          // 角色ID
	Name        string         `json:"name"`        // 角色名称
	Description string         `json:"description"` // 角色描述
	CreatedAt   time.Time      `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`  // 更新时间
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`  // 删除时间
}

// GetRoleDTOExample 角色列表返回swagger示例
type GetRoleDTOExample struct {
	ID          uint      `json:"id"`          // 角色ID
	Name        string    `json:"name"`        // 角色名称
	Description string    `json:"description"` // 角色描述
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
	DeletedAt   time.Time `json:"deleted_at"`  // 删除时间
}

type UpdateRoleReqDTO struct {
	Name        string `json:"name" binding:"required,max=20"`         // 角色名称
	Description string `json:"description" binding:"required,max=255"` // 角色描述
}

type DeleteRoleReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 角色ID
}

type RecoverRoleReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 角色ID
}
