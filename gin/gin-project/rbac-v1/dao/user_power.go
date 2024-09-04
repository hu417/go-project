package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"rbac-v1/dao/base"
	"rbac-v1/model/po"
)

//根据ids获取userPo列表
func (d *Dao) GetUserPoByUserIds(ctx context.Context, ids []uint) (ret []*po.UserPo, err error) {
	data, err := d.userPoList(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"user_id in (?)": ids,
		},
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

//删除
func (d *Dao) DelUserPoByPowerId(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	_, err = d.userPoDelete(ctx, tx, &base.DBConditions{
		And: map[string]interface{}{
			"power_id = ?": id,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
func (d *Dao) DelUserPoByUserId(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	_, err = d.userPoDelete(ctx, tx, &base.DBConditions{
		And: map[string]interface{}{
			"user_id = ?": id,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

//快捷操作封装
//列表
func (d *Dao) userPoList(ctx context.Context, conditions *base.DBConditions) (ret []*po.UserPo, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = []*po.UserPo{}

	db = db.Table((&po.UserPo{}).TableName())
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
func (d *Dao) getUserPo(ctx context.Context, conditions *base.DBConditions) (ret *po.UserPo, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = &po.UserPo{}

	db = db.Table((&po.UserPo{}).TableName())
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
func (d *Dao) UserPoBatchInsert(ctx context.Context, tx *gorm.DB, row []*po.UserPo) (err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	return tx.Create(row).Error
}
//修改
func (d *Dao) userPoUpdate(ctx context.Context, tx *gorm.DB, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Model(&po.UserPo{}).Where(where).Updates(update)
	if db.Error != nil {
		return 0, err
	}
	affected = db.RowsAffected

	return affected, nil
}
//删除
func (d *Dao) userPoDelete(ctx context.Context, tx *gorm.DB, conditions *base.DBConditions) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Table((&po.UserPo{}).TableName())
	db = conditions.Fill(db)

	err = db.Delete(&po.UserPo{}).Error
	if err != nil {
		return 0, err
	}

	affected = db.RowsAffected

	return affected, nil
}