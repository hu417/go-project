package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"rbac-v1/dao/base"
	"rbac-v1/model/po"
	"rbac-v1/model/vo"
)
//根据id获取op
func (d *Dao) GetOperationById(ctx context.Context, id uint) (ret *po.Operation, err error) {
	data, err := d.getOperation(ctx, &base.DBConditions{
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
func (d *Dao) GetOperationByIds(ctx context.Context, ids []uint) (ret []*po.Operation, err error) {
	data, err := d.operationList(ctx, &base.DBConditions{
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
func (d *Dao) IsOperationExistById(ctx context.Context, id uint) (ok bool) {
	data, err := d.getOperation(ctx, &base.DBConditions{
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
func (d *Dao) IsOperationExist(ctx context.Context, path, method string) (ok bool) {
	data, err := d.getOperation(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"is_deleted = ?": 0,
			"path = ?": path,
			"method = ?": method,
		},
	})
	if err == nil && data != nil {
		return true
	}
	return false
}
//根据id更新
func (d *Dao) UpdateOperationById(ctx context.Context, tx *gorm.DB, data *po.Operation) error {
	where := map[string]interface{}{
		"is_deleted": 0,
		"id": data.Id,
	}
	update := map[string]interface{}{
		"path": data.Path,
		"type": data.Type,
		"method": data.Method,
		"description": data.Description,
		"created_at": data.CreatedAt,
		"updated_at": data.UpdatedAt,
	}

	_, err := d.operationUpdate(ctx, tx, where, update)
	if err != nil {
		return err
	}
	return nil
}
//删除
func (d *Dao) DelOperationById(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	where := map[string]interface{}{
		"id": id,
		"is_deleted": 0,
	}
	update := map[string]interface{}{
		"is_deleted": 1,
	}

	_, err = d.operationUpdate(ctx, tx, where, update)
	if err != nil {
		return nil
	}
	return err
}
//列表
func (d *Dao) GetOperationList(ctx context.Context, param *vo.OperationListRequest) (ret *vo.OperationListResponse, err error) {
	where := map[string]interface{}{
		"is_deleted = ?": 0,
	}
	if len(param.Ids) != 0 {
		where["id in (?)"] = param.Ids
	}
	if len(param.Path) != 0 {
		where["path like ?"] = "%" + param.Path + "%"
	}
	if len(param.Method) != 0 {
		where["method = ?"] = param.Method
	}
	if param.Type != 0 {
		where["type = ?"] = param.Type
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

	operations, err := d.operationList(ctx, condition)
	if err != nil {
		return nil, err
	}

	return &vo.OperationListResponse{
		Total: condition.Count,
		List:  operations,
	}, nil
}

//快捷操作封装
//列表
func (d *Dao) operationList(ctx context.Context, conditions *base.DBConditions) (ret []*po.Operation, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = []*po.Operation{}

	db = db.Table((&po.Operation{}).TableName())
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
func (d *Dao) getOperation(ctx context.Context, conditions *base.DBConditions) (ret *po.Operation, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = &po.Operation{}

	db = db.Table((&po.Operation{}).TableName())
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
func (d *Dao) OperationBatchInsert(ctx context.Context, tx *gorm.DB, row []*po.Operation) (err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	return tx.Create(row).Error
}
//修改
func (d *Dao) operationUpdate(ctx context.Context, tx *gorm.DB, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Model(&po.Operation{}).Where(where).Updates(update)
	if db.Error != nil {
		return 0, err
	}
	affected = db.RowsAffected

	return affected, nil
}