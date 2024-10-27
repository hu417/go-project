package dtos

import (
	"time"

	"gorm.io/gorm"
)

type GetPermissionListReqDTO struct {
	PaginationReqDTO
	Name        string `form:"name" binding:"omitempty"`        // 权限名称
	Method      string `form:"method" binding:"omitempty"`      // 请求方法
	ApiPath     string `form:"api_path" binding:"omitempty"`    // 请求路径
	Description string `form:"description" binding:"omitempty"` // 权限描述
	ApiGroup    string `form:"api_group" binding:"omitempty"`   // API分组
}

type CreatePermissionReqDTO struct {
	Name        string `json:"name" binding:"required,max=40"`          // 权限名称
	ApiPath     string `json:"api_path" binding:"required,max=255"`     // 请求路径
	Method      string `json:"method" binding:"required,max=10"`        // 请求方法
	Description string `json:"description" binding:"omitempty,max=255"` // 权限描述
	ApiGroup    string `json:"api_group" binding:"required,max=40"`     // API分组
}

type GetPermissionReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 权限ID
}

type GetPermissionResDTO struct {
	ID          uint           `json:"id"`          // 权限ID
	Name        string         `json:"name"`        // 权限名称
	ApiPath     string         `json:"api_path"`    // 请求路径
	Method      string         `json:"method"`      // 请求方法
	Description string         `json:"description"` // 权限描述
	ApiGroup    string         `json:"api_group"`   // API分组
	CreatedAt   time.Time      `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`  // 更新时间
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`  // 删除时间
}

// GetPermissionDTOExample 获取权限响应swagger示例
type GetPermissionDTOExample struct {
	ID          uint      `json:"id"`          // 权限ID
	Name        string    `json:"name"`        // 权限名称
	ApiPath     string    `json:"api_path"`    // 请求路径
	Method      string    `json:"method"`      // 请求方法
	Description string    `json:"description"` // 权限描述
	ApiGroup    string    `json:"api_group"`   // API分组
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
	DeletedAt   time.Time `json:"deleted_at"`  // 删除时间
}

type UpdatePermissionReqDTO struct {
	Name        string `json:"name" binding:"required,max=40"`         // 权限名称
	ApiPath     string `json:"api_path" binding:"required,max=255"`    // 请求路径
	Method      string `json:"method" binding:"required,max=10"`       // 请求方法
	Description string `json:"description" binding:"required,max=255"` // 权限描述
	ApiGroup    string `json:"api_group" binding:"required,max=40"`    // API分组
}

type DeletePermissionReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 权限ID
}

type RecoverPermissionReqDTO struct {
	ID uint `uri:"id" binding:"required"` // 权限ID
}
