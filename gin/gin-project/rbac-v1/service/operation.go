package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"rbac-v1/model/po"
	"rbac-v1/model/vo"
	"time"
)

// 列表
func (s *Service) GetOperationList(ctx context.Context, param *vo.OperationListRequest) (ret *vo.OperationListResponse, err error) {
	data, err := s.dao.GetOperationList(ctx, param)
	if err != nil {
		logger.Error("get operation list failed: ", err)
		return nil, err
	}

	return data, nil
}

//创建
func (s *Service) OperationCreate(ctx context.Context, param *vo.OperationCreateRequest) (err error) {
	//检查是否存在
	for _, op := range param.Operations {
		if ok := s.dao.IsOperationExist(ctx, op.Path, op.Method); ok {
			err = errors.New(fmt.Sprintf("create operation failed: duplicate data: %s %s", op.Path, op.Method))
			logger.Error(err.Error())
			return err
		}
	}
	//创建
	//Transaction db的事务方法，一般情况下，我们把变更操作放在里面
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.dao.OperationBatchInsert(ctx, tx, param.Operations); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error("create operation failed: ", err)
		return err
	}

	return nil
}
//更新
func (s *Service) OperationUpdate(ctx context.Context, param *po.Operation) (err error) {
	//获取数据
	op, err := s.dao.GetOperationById(ctx, param.Id)
	if err != nil {
		logger.Error("get operation failed: ", err)
		return err
	}
	//判断是否存在
	if op == nil {
		err = errors.New("no such operation data")
		logger.Error(err.Error())
		return err
	}
	//避免创建时间被人为更改
	param.CreatedAt = op.CreatedAt
	param.UpdatedAt = time.Now().Local()

	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.dao.UpdateOperationById(ctx, tx, param); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error("update operation failed: ", err)
		return err
	}

	return nil
}

//删除
func (s *Service) OperationDelete(ctx context.Context, param *vo.OperationDelRequest) (err error) {
	//检查主表是否存在
	if ok := s.dao.IsOperationExistById(ctx, param.Id); !ok {
		err = errors.New("no such operation data")
		logger.Error(err)
		return err
	}

	//删除
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.dao.DelOperationById(ctx, tx, param.Id); err != nil {
			return err
		}
		if err := s.dao.DelPowerOpByOpId(ctx, tx, param.Id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error("delete operation failed: ", err)
		return err
	}

	return nil
}