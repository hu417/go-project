package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"rbac-v1/dao/base"
	"rbac-v1/model/po"
	"rbac-v1/model/vo"
)

//根据id获取power
func (d *Dao) GetPowerById(ctx context.Context, id uint) (ret *po.Power, err error) {
	data, err := d.getPower(ctx, &base.DBConditions{
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
func (d *Dao) GetPowerByIds(ctx context.Context, ids []uint) (ret []*po.Power, err error) {
	data, err := d.powerList(ctx, &base.DBConditions{
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
func (d *Dao) IsPowerExistById(ctx context.Context, id uint) (ok bool) {
	data, err := d.getPower(ctx, &base.DBConditions{
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
func (d *Dao) IsPowerExist(ctx context.Context, code string) (ok bool) {
	data, err := d.getPower(ctx, &base.DBConditions{
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
func (d *Dao) UpdatePowerById(ctx context.Context, tx *gorm.DB, data *po.Power) error {
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

	_, err := d.powerUpdate(ctx, tx, where, update)
	if err != nil {
		return err
	}
	return nil
}
//删除
func (d *Dao) DelPowerById(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	where := map[string]interface{}{
		"id": id,
		"is_deleted": 0,
	}
	update := map[string]interface{}{
		"is_deleted": 1,
	}

	_, err = d.powerUpdate(ctx, tx, where, update)
	if err != nil {
		return nil
	}
	return err
}
//列表
func (d *Dao) GetPowerList(ctx context.Context, param *vo.PowerListRequest) (ret *vo.PowerListResponse, err error) {
	where := map[string]interface{}{
		"is_deleted = ?": 0,
	}
	if len(param.Name) != 0 {
		where["name like ?"] = "%" + param.Name + "%"
	}
	if len(param.Code) != 0 {
		where["code like ?"] = "%" + param.Code + "%"
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

	powers, err := d.powerList(ctx, condition)
	if err != nil {
		return nil, err
	}

	return &vo.PowerListResponse{
		Total: condition.Count,
		List:  powers,
	}, nil
}

//快捷操作封装
//列表
func (d *Dao) powerList(ctx context.Context, conditions *base.DBConditions) (ret []*po.Power, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = []*po.Power{}

	db = db.Table((&po.Power{}).TableName())
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
func (d *Dao) getPower(ctx context.Context, conditions *base.DBConditions) (ret *po.Power, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = &po.Power{}

	db = db.Table((&po.Power{}).TableName())
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
func (d *Dao) PowerBatchInsert(ctx context.Context, tx *gorm.DB, row []*po.Power) (err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	return tx.Create(row).Error
}
//修改
func (d *Dao) powerUpdate(ctx context.Context, tx *gorm.DB, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Model(&po.Power{}).Where(where).Updates(update)
	if db.Error != nil {
		return 0, err
	}
	affected = db.RowsAffected

	return affected, nil
}