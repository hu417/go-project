package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"rbac-v1/dao/base"
	"rbac-v1/model/po"
	"rbac-v1/model/vo"
)

//根据id获取role
func (d *Dao) GetRoleById(ctx context.Context, id uint) (ret *po.Role, err error) {
	data, err := d.getRole(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"is_deleted = ?": 0,
			"id = ?": id,
		},
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
//根据ids获取op列表
func (d *Dao) GetRoleByIds(ctx context.Context, ids []uint) (ret []*po.Role, err error) {
	data, err := d.roleList(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"is_deleted = ?": 0,
			"id in (?)": ids,
		},
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

//根据id判断op是否存在
func (d *Dao) IsRoleExistById(ctx context.Context, id uint) (ok bool) {
	data, err := d.getRole(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"is_deleted = ?": 0,
			"id = ?": id,
		},
	})
	if err == nil && data != nil {
		return true
	}
	return false
}
//根据path和method判断是否存在
func (d *Dao) IsRoleExist(ctx context.Context, code string) (ok bool) {
	data, err := d.getRole(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"is_deleted = ?": 0,
			"code = ?": code,
		},
	})
	if err == nil && data != nil {
		return true
	}
	return false
}
//根据id更新
func (d *Dao) UpdateRoleById(ctx context.Context, tx *gorm.DB, data *po.Role) error {
	where := map[string]interface{}{
		"is_deleted": 0,
		"id": data.Id,
	}
	update := map[string]interface{}{
		"name": data.Name,
		"code": data.Code,
		"description": data.Description,
		"created_at": data.CreatedAt,
		"updated_at": data.UpdatedAt,
	}

	_, err := d.roleUpdate(ctx, tx, where, update)
	if err != nil {
		return err
	}
	return nil
}
//删除
func (d *Dao) DelRoleById(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	where := map[string]interface{}{
		"id": id,
		"is_deleted": 0,
	}
	update := map[string]interface{}{
		"is_deleted": 1,
	}

	_, err = d.roleUpdate(ctx, tx, where, update)
	if err != nil {
		return nil
	}
	return err
}
//列表
func (d *Dao) GetRoleList(ctx context.Context, param *vo.RoleListRequest) (ret *vo.RoleListResponse, err error) {
	where := map[string]interface{}{
		"is_deleted = ?": 0,
	}
	if len(param.Name) != 0 {
		where["name like ?"] = "%" + param.Name + "%"
	}
	if len(param.Code) != 0 {
		where["code like ?"] = "%" + param.Code + "%"
	}
	if len(param.Ids) != 0 {
		where["id in (?)"] = param.Ids
	}

	condition := &base.DBConditions{
		And: where,
		NeedCount: true,
		Order: "id DESC",
	}

	if param.Page != 0 && param.PageSize != 0 {
		condition.Limit = param.PageSize
		condition.Offset = (param.Page -1) * param.PageSize
	}

	roles, err := d.roleList(ctx, condition)
	if err != nil {
		return nil, err
	}

	return &vo.RoleListResponse{
		Total: condition.Count,
		List:  roles,
	}, nil
}

//快捷操作封装
//列表
func (d *Dao) roleList(ctx context.Context, conditions *base.DBConditions) (ret []*po.Role, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = []*po.Role{}

	db = db.Table((&po.Role{}).TableName())
	db = conditions.Fill(db)

	err = db.Find(&ret).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return ret, nil
}
//详情
func (d *Dao) getRole(ctx context.Context, conditions *base.DBConditions) (ret *po.Role, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = &po.Role{}

	db = db.Table((&po.Role{}).TableName())
	db = conditions.Fill(db)

	err = db.First(&ret).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return ret, nil
}
//新增
//tx *gorm.DB 使用来做事务控制的
func (d *Dao) RoleBatchInsert(ctx context.Context, tx *gorm.DB, row []*po.Role) (err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	return tx.Create(row).Error
}
//修改
func (d *Dao) roleUpdate(ctx context.Context, tx *gorm.DB, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Model(&po.Role{}).Where(where).Updates(update)
	if db.Error != nil {
		return 0, err
	}
	affected = db.RowsAffected

	return affected, nil
}