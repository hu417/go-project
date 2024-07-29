package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"rbac-v1/dao/base"
	"rbac-v1/model/po"
)

//根据ids获取powerOp列表
func (d *Dao) GetPowerOpByPowerIds(ctx context.Context, ids []uint) (ret []*po.PowerOp, err error) {
	data, err := d.powerOpList(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"power_id in (?)": ids,
		},
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

//删除
func (d *Dao) DelPowerOpByPowerId(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	_, err = d.powerOpDelete(ctx, tx, &base.DBConditions{
		And: map[string]interface{}{
			"power_id = ?": id,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
func (d *Dao) DelPowerOpByOpId(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	_, err = d.powerOpDelete(ctx, tx, &base.DBConditions{
		And: map[string]interface{}{
			"operation_id = ?": id,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

//快捷操作封装
//列表
func (d *Dao) powerOpList(ctx context.Context, conditions *base.DBConditions) (ret []*po.PowerOp, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = []*po.PowerOp{}

	db = db.Table((&po.PowerOp{}).TableName())
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
func (d *Dao) getPowerOp(ctx context.Context, conditions *base.DBConditions) (ret *po.PowerOp, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = &po.PowerOp{}

	db = db.Table((&po.PowerOp{}).TableName())
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
func (d *Dao) PowerOpBatchInsert(ctx context.Context, tx *gorm.DB, row []*po.PowerOp) (err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	return tx.Create(row).Error
}
//修改
func (d *Dao) powerOpUpdate(ctx context.Context, tx *gorm.DB, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Model(&po.PowerOp{}).Where(where).Updates(update)
	if db.Error != nil {
		return 0, err
	}
	affected = db.RowsAffected

	return affected, nil
}
//删除
func (d *Dao) powerOpDelete(ctx context.Context, tx *gorm.DB, conditions *base.DBConditions) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Table((&po.PowerOp{}).TableName())
	db = conditions.Fill(db)

	err = db.Delete(&po.PowerOp{}).Error
	if err != nil {
		return 0, err
	}

	affected = db.RowsAffected

	return affected, nil
}