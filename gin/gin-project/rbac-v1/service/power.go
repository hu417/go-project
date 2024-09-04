package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"rbac-v1/common"
	"rbac-v1/model/bo"
	"rbac-v1/model/po"
	"rbac-v1/model/vo"
	"time"
)

func (s *Service) GetPowerList(ctx context.Context, param *vo.PowerListRequest) (ret *vo.PowerBoListResponse, err error) {
	ret = &vo.PowerBoListResponse{}
	//获取powerlist
	powerList, err := s.dao.GetPowerList(ctx, param)
	if err != nil {
		logger.Error("get power list failed: ", err)
		return nil, err
	}
	if powerList.List == nil {
		return ret, nil
	}
	//组装power ids
	powerIds := []uint{}
	for _, power := range powerList.List {
		powerIds = append(powerIds, power.Id)
	}
	//根据power ids 查询powerOp关联
	powerOps, err := s.dao.GetPowerOpByPowerIds(ctx, powerIds)
	if err != nil {
		logger.Error("get powerOp list failed: ", err)
		return nil, err
	}
	//没有关联到operation，就把po power转成bo power返回
	if powerOps == nil {
		powerBos := []*bo.Power{}
		for _, power := range powerList.List {
			powerBo := &bo.Power{
				Power: power,
			}
			powerBos = append(powerBos, powerBo)
		}
		return &vo.PowerBoListResponse{
			Total: powerList.Total,
			List:  powerBos,
		}, nil
	}

	//组装operation ids
	powerOpMp := map[uint][]uint{} //记录power和opration的对应关系
	opIds := []uint{}
	for _, powerOp := range powerOps {
		if _, ok := powerOpMp[powerOp.PowerId]; !ok {
			powerOpMp[powerOp.PowerId] = []uint{}
		}
		opIds = append(opIds, powerOp.OperationId)
		powerOpMp[powerOp.PowerId] = append(powerOpMp[powerOp.PowerId], powerOp.OperationId)
	}
	//查询operation
	operations, err := s.dao.GetOperationByIds(ctx, common.RemoveDuplicates(opIds))
	if err != nil {
		logger.Error("get operations failed: ", err.Error())
		return nil, err
	}
	opMp := map[uint]*po.Operation{}
	for _, operation := range operations {
		opMp[operation.Id] = operation
	}
	//组装到powerbo list中
	powerBos := []*bo.Power{}
	for _, power := range powerList.List {
		powerBo := &bo.Power{
			Power:      power,
			Operations: []*po.Operation{},
		}
		if _, ok := powerOpMp[power.Id];ok {
			for _, operationId := range powerOpMp[power.Id] {
				powerBo.Operations = append(powerBo.Operations, opMp[operationId])
			}
		}
		powerBos = append(powerBos, powerBo)
	}

	return &vo.PowerBoListResponse{
		Total: powerList.Total,
		List:  powerBos,
	}, nil
}

//创建
func (s *Service) PowerCreate(ctx context.Context, param *vo.PowerCreateRequest) (err error) {
	//检查新增业务主键是否存在
	opIds := []uint{}
	powers := []*po.Power{}
	for _, power := range param.Powers {
		if ok := s.dao.IsPowerExist(ctx, power.Code); ok {
			err := errors.New(fmt.Sprintf("create power failed: duplicate data: %s", power.Code))
			logger.Error(err.Error())
			return err
		}
		powers = append(powers, power.Power)
		for _, op := range power.Operations {
			opIds = append(opIds, op.Id)
		}
	}
	//检查关联数据是否存在
	if len(opIds) != 0 {
		opIds = common.RemoveDuplicates(opIds)
		for _, opId := range opIds {
			if ok := s.dao.IsOperationExistById(ctx, opId); !ok {
				err = errors.New(fmt.Sprintf("create power failed: operation not exist: %d", opId))
				logger.Error(err.Error())
				return err
			}
		}
	}
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		//新增数据
		if err := s.dao.PowerBatchInsert(ctx, tx, powers); err != nil {
			return err
		}
		//处理关联
		powerOps := []*po.PowerOp{}
		for i := 0; i < len(powers); i++ {
			if len(param.Powers[i].Operations) != 0 {
				for _, operation := range param.Powers[i].Operations {
					powerOps = append(powerOps, &po.PowerOp{
						PowerId:     powers[i].Id,
						OperationId: operation.Id,
					})
				}
			}
		}
		if len(powerOps) != 0 {
			if err := s.dao.PowerOpBatchInsert(ctx, tx, powerOps); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Error("create power failed: ", err)
	}

	return nil
}

//更新
func (s *Service) PowerUpdate(ctx context.Context, param *bo.Power) (err error) {
	//验证是否存在
	power, err := s.dao.GetPowerById(ctx, param.Id)
	if err != nil {
		logger.Error("get power failed: ", err)
		return err
	}
	if power == nil {
		err = errors.New("no such power data")
		return err
	}
	param.CreatedAt = power.CreatedAt
	param.UpdatedAt = time.Now().Local()
	//组装powerOp
	powerOps := []*po.PowerOp{}
	if len(param.Operations) != 0 {
		for _, op := range param.Operations {
			if ok := s.dao.IsOperationExistById(ctx, op.Id); !ok {
				err = errors.New(fmt.Sprintf("update power failed: operation not exist: %d", op.Id))
				logger.Error(err.Error())
				return err
			}
			powerOps = append(powerOps, &po.PowerOp{
				PowerId:     param.Id,
				OperationId: op.Id,
			})
		}
	}
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		//主表新增
		if err := s.dao.UpdatePowerById(ctx, tx, param.Power); err != nil {
			return err
		}
		//关联表删除
		if err := s.dao.DelPowerOpByPowerId(ctx, tx, param.Id); err != nil {
			return err
		}
		//关联表新增
		if len(powerOps) != 0 {
			if err := s.dao.PowerOpBatchInsert(ctx, tx, powerOps); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Error("update power failed: ", err)
		return err
	}

	return nil
}

//删除
func (s *Service) PowerDelete(ctx context.Context, param *vo.PowerDelRequest) (err error) {
	//验证是否存在
	if ok := s.dao.IsPowerExistById(ctx, param.Id); !ok {
		err = errors.New("power is not exist")
		logger.Error(err)
		return err
	}
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		//删除主表
		if err := s.dao.DelPowerById(ctx, tx, param.Id); err != nil {
			return err
		}
		//删除关联表
		if err := s.dao.DelPowerOpByPowerId(ctx, tx, param.Id); err != nil {
			return err
		}
		if err := s.dao.DelRolePoByPowerId(ctx, tx, param.Id); err != nil {
			return err
		}
		if err := s.dao.DelUserPoByPowerId(ctx, tx, param.Id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error("delete power failed: ", err)
		return err
	}

	return nil
}