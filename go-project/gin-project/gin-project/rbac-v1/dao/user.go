package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"rbac-v1/dao/base"
	"rbac-v1/model/po"
	"rbac-v1/model/vo"
)

func (d *Dao) GetUserByUsername(ctx context.Context, username string) (ret *po.User, err error) {
	data, err := d.getUser(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"is_deleted = ?": 0,
			"username = ?": username,
		},
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

//根据id获取user
func (d *Dao) GetUserById(ctx context.Context, id uint) (ret *po.User, err error) {
	data, err := d.getUser(ctx, &base.DBConditions{
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

func (d *Dao) IsUserExist(ctx context.Context, username string) (ok bool) {
	data, err := d.getUser(ctx, &base.DBConditions{
		And: map[string]interface{}{
			"is_deleted = ?": 0,
			"username = ?": username,
		},
	})
	if err == nil && data != nil {
		return true
	}
	return false
}
//根据id更新
func (d *Dao) UpdateUserById(ctx context.Context, tx *gorm.DB, data *po.User) error {
	where := map[string]interface{}{
		"is_deleted": 0,
		"id": data.Id,
	}
	update := map[string]interface{}{
		"name": data.Name,
		"username": data.Username,
		"mail": data.Mail,
		"phone": data.Phone,
		"password": data.Password,
		"created_at": data.CreatedAt,
		"updated_at": data.UpdatedAt,
	}

	_, err := d.userUpdate(ctx, tx, where, update)
	if err != nil {
		return err
	}
	return nil
}
//删除
func (d *Dao) DelUserById(ctx context.Context, tx *gorm.DB, id uint) (err error) {
	where := map[string]interface{}{
		"id": id,
		"is_deleted": 0,
	}
	update := map[string]interface{}{
		"is_deleted": 1,
	}

	_, err = d.userUpdate(ctx, tx, where, update)
	if err != nil {
		return nil
	}
	return err
}
//列表
func (d *Dao) GetUserList(ctx context.Context, param *vo.UserListRequest) (ret *vo.UserListResponse, err error) {
	where := map[string]interface{}{
		"is_deleted = ?": 0,
	}
	if len(param.Name) != 0 {
		where["name like ?"] = "%" + param.Name + "%"
	}
	if len(param.Username) != 0 {
		where["username like ?"] = "%" + param.Username + "%"
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

	users, err := d.userList(ctx, condition)
	if err != nil {
		return nil, err
	}

	return &vo.UserListResponse{
		Total: condition.Count,
		List:  users,
	}, nil
}

//快捷操作封装
//列表
func (d *Dao) userList(ctx context.Context, conditions *base.DBConditions) (ret []*po.User, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = []*po.User{}

	db = db.Table((&po.User{}).TableName())
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
func (d *Dao) getUser(ctx context.Context, conditions *base.DBConditions) (ret *po.User, err error) {
	db := d.db.WithContext(ctx) //控制连接的超时和取消
	ret = &po.User{}

	db = db.Table((&po.User{}).TableName())
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
func (d *Dao) UserBatchInsert(ctx context.Context, tx *gorm.DB, row []*po.User) (err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	return tx.Create(row).Error
}
//修改
func (d *Dao) userUpdate(ctx context.Context, tx *gorm.DB, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
	if tx == nil {
		tx = d.db.WithContext(ctx)
	}

	db := tx.Model(&po.User{}).Where(where).Updates(update)
	if db.Error != nil {
		return 0, err
	}
	affected = db.RowsAffected

	return affected, nil
}