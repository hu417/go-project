package keys

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 定义一个结构体来接收查询结果
type PermissionResult struct {
	ApiPath string `gorm:"column:api_path"`
	Method  string `gorm:"column:method"`
}

type PermissionList struct {
	db  *gorm.DB
	rds *redis.Client
	ctx context.Context
}

func NewPermissionList(db *gorm.DB, rds *redis.Client, ctx context.Context) *PermissionList {
	return &PermissionList{
		db:  db,
		rds: rds,
		ctx: ctx,
	}
}

// GetPermissionListByUserID 根据用户ID获取权限列表，并将其缓存到Redis中
//
// 描述:
// 函数首先尝试从Redis缓存中获取用户的权限列表。如果存在缓存，则直接解码JSON数据并返回。
// 如果Redis中没有缓存或者缓存失效，函数将从数据库中查询用户的权限信息。
// 查询结果将被转换成一个映射，其中键是API路径，值是另一个映射，其键是HTTP方法，值是一个布尔值，表示该用户是否具有对相应路径和方法的访问权限。
// 最后，查询结果会被序列化为JSON格式并存储回Redis中，以便未来的请求可以直接使用缓存数据。

// 返回:
// - map[string]map[string]bool: 用户的权限列表，其中键是API路径，内部映射的键是HTTP方法，值表示用户是否有此权限
// - error: 如果发生错误，返回错误信息
func (p *PermissionList) GetPermissionListByUserID(userID uint) (map[string]map[string]bool, error) {
	// 生成一个唯一的缓存键，可以根据用户ID和权限列表的特性来生成
	cacheKey := fmt.Sprintf("user:%d:permissions", userID)

	// 尝试从Redis中获取缓存数据
	cachedData, err := p.rds.Get(p.ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 如果Redis中有缓存数据，直接返回
	if err == nil {
		var result map[string]map[string]bool
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return nil, err
		}
		return result, nil
	}

	// 如果Redis中没有缓存数据，从数据库查询
	result := make(map[string]map[string]bool)

	var permissions []PermissionResult

	err = p.db.WithContext(p.ctx).Table("s_user").
		Select("s_permission.api_path, s_permission.method").
		Joins("INNER JOIN r_user_role ON s_user.id = r_user_role.user_id").
		Joins("INNER JOIN r_role_permission ON r_user_role.role_id = r_role_permission.role_id").
		Joins("INNER JOIN s_permission ON r_role_permission.permission_id = s_permission.id").
		Where("s_user.id = ?", userID).
		Scan(&permissions).Error

	if err != nil {
		return nil, err
	}

	for _, perm := range permissions {
		if _, ok := result[perm.ApiPath]; !ok {
			result[perm.ApiPath] = make(map[string]bool)
		}
		result[perm.ApiPath][perm.Method] = true
	}

	// 将查询结果序列化并存储到Redis中
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	if err := p.rds.Set(p.ctx, cacheKey, jsonBytes, 5*time.Minute).Err(); err != nil {
		return nil, err
	}

	return result, nil
}
