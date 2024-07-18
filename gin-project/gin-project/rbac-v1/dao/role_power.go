package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"rbac-v1/dao/base"
	"rbac-v1/model/po"
)

//根据ids获取rolePo列表
func (d *Dao) GetRolePoByRoleIds(ctx context.Context, ids []uint) (ret []*po.RolePo, err error) {
	data, err := d.rolePoList(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"role_id in (?)": ids,
		},
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

//删除
func (d *Dao) DelRolePoByPowerId(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	_, err = d.rolePoDelete(ctx, tx, &base.DBConditions{
		And: map[string]interface{}{
			"power_id = ?": id,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
func (d *Dao) DelRolePoByRoleId(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	_, err = d.rolePoDelete(ctx, tx, &base.DBConditions{
		And: map[string]interface{}{
			"role_id = ?": id,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

//快捷操作封装
//列表
func (d *Dao) rolePoList(ctx context.Context, conditions *base.DBConditions) (ret []*po.RolePo, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = []*po.RolePo{}

	db = db.Table((&po.RolePo{}).TableName())
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
func (d *Dao) getRolePo(ctx context.Context, conditions *base.DBConditions) (ret *po.RolePo, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = &po.RolePo{}

	db = db.Table((&po.RolePo{}).TableName())
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
func (d *Dao) RolePoBatchInsert(ctx context.Context, tx *gorm.DB, row []*po.RolePo) (err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	return tx.Create(row).Error
}
//修改
func (d *Dao) rolePoUpdate(ctx context.Context, tx *gorm.DB, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Model(&po.RolePo{}).Where(where).Updates(update)
	if db.Error != nil {
		return 0, err
	}
	affected = db.RowsAffected

	return affected, nil
}
//删除
func (d *Dao) rolePoDelete(ctx context.Context, tx *gorm.DB, conditions *base.DBConditions) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Table((&po.RolePo{}).TableName())
	db = conditions.Fill(db)

	err = db.Delete(&po.RolePo{}).Error
	if err != nil {
		return 0, err
	}

	affected = db.RowsAffected

	return affected, nil
}