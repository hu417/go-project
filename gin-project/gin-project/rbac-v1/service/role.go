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

//列表
func (s *Service) GetRoleList(ctx context.Context, param *vo.RoleListRequest) (ret *vo.RoleBoListResponse, err error) {
	//查询po的rolelist
	ret = &vo.RoleBoListResponse{}
	roleList, err := s.dao.GetRoleList(ctx, param)
	if err != nil {
		logger.Error("get role list failed: ", err.Error())
		return nil, err
	}
	if roleList.List == nil {
		return ret, nil
	}
	//组装主表ids，查询中间表数据
	roleIds := []uint{}
	for _, role := range roleList.List {
		roleIds = append(roleIds, role.Id)
	}
	rolePos, err := s.dao.GetRolePoByRoleIds(ctx, roleIds)
	if err != nil {
		logger.Error("get rolePo list failed: ", err.Error())
		return nil, err
	}
	//处理无关联数据的情况
	if rolePos == nil {
		//把po的rolelist转为bo的rolelist，里面的关联数据为空值
		roleBos := []*bo.Role{}
		for _, role := range roleList.List {
			roleBo := &bo.Role{
				Role:   role,
			}
			roleBos = append(roleBos, roleBo)
		}

		return &vo.RoleBoListResponse{
			Total: roleList.Total,
			List:  roleBos,
		}, nil
	}
	//组装关联ids
	rolePoMp := map[uint][]uint{}
	poIds := []uint{}
	for _, rolePo := range rolePos {
		if _, ok := rolePoMp[rolePo.RoleId]; !ok {
			rolePoMp[rolePo.RoleId] = []uint{}
		}
		poIds = append(poIds, rolePo.PowerId)
		rolePoMp[rolePo.RoleId] = append(rolePoMp[rolePo.RoleId], rolePo.PowerId)
	}
	//查询关联ids
	powers, err := s.dao.GetPowerByIds(ctx, poIds)
	if err != nil {
		logger.Error("get powers failed: ", err.Error())
		return nil, err
	}
	poMp := map[uint]*po.Power{}
	for _, power := range powers {
		poMp[power.Id] = power
	}
	//组装数据
	roleBos := []*bo.Role{}
	for _, role := range roleList.List {
		roleBo := &bo.Role{
			Role:   role,
			Powers: []*po.Power{},
		}
		if _, ok := rolePoMp[role.Id]; ok {
			for _, powerId := range rolePoMp[role.Id] {
				roleBo.Powers = append(roleBo.Powers, poMp[powerId])
			}
		}
		roleBos = append(roleBos, roleBo)
	}

	return &vo.RoleBoListResponse{
		Total: roleList.Total,
		List:  roleBos,
	}, nil
}

//新增
func (s *Service) RoleCreate(ctx context.Context, param *vo.RoleCreateRequest) (err error) {
	//检查主键是否存在
	poIds := []uint{}
	roles := []*po.Role{}
	for _, role := range param.Roles {
		if ok := s.dao.IsRoleExist(ctx, role.Code); ok {
			err = errors.New(fmt.Sprintf("create role failed, duplicate data: %s", role.Code))
			logger.Error(err.Error())
			return err
		}
		roles = append(roles, role.Role)
		for _, power := range role.Powers {
			poIds = append(poIds, power.Id)
		}
	}
	if len(poIds) != 0 {
		poIds = common.RemoveDuplicates(poIds)
		for _, poId := range poIds {
			if ok := s.dao.IsPowerExistById(ctx, poId); !ok {
				err = errors.New(fmt.Sprintf("create role failed, power not exist: %d", poId))
				logger.Error(err.Error())
				return err
			}
		}
	}
	//事务
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		//新增主数据
		if err := s.dao.RoleBatchInsert(ctx, tx, roles); err != nil {
			return err
		}
		//新增关联数据
		rolePos := []*po.RolePo{}
		for i :=0; i < len(roles); i++ {
			if len(param.Roles[i].Powers) != 0 {
				for _, power := range param.Roles[i].Powers {
					rolePos = append(rolePos, &po.RolePo{
						RoleId:  roles[i].Id,
						PowerId: power.Id,
					})
				}
			}
		}
		if len(rolePos) != 0 {
			if err := s.dao.RolePoBatchInsert(ctx, tx, rolePos); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Error("create role failed: ", err.Error())
		return err
	}
	return nil
}
//更新
func (s *Service) RoleUpdate(ctx context.Context, param *bo.Role) (err error) {
	//查询主表数据
	role, err := s.dao.GetRoleById(ctx, param.Id)
	if err != nil {
		logger.Error("get role failed: ", err)
		return err
	}
	if role == nil {
		err = errors.New("no such role data")
		logger.Error(err.Error())
		return err
	}
	param.CreatedAt = role.CreatedAt
	param.UpdatedAt = time.Now().Local()
	//组装新的关联数据
	rolePos := []*po.RolePo{}
	if len(param.Powers) != 0 {
		for _, power := range param.Powers {
			if ok := s.dao.IsPowerExistById(ctx, power.Id); !ok {
				err = errors.New(fmt.Sprintf("update role failed, power not exist: %d", power.Id))
				logger.Error(err.Error())
				return err
			}
			rolePos = append(rolePos, &po.RolePo{
				RoleId:  param.Id,
				PowerId: power.Id,
			})
		}
	}
	//事务
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		//更新主表数据
		if err := s.dao.UpdateRoleById(ctx, tx, param.Role); err != nil {
			return err
		}
		//删除-新增关联表数据
		if err := s.dao.DelRolePoByRoleId(ctx, tx, param.Id); err != nil {
			return err
		}
		if len(rolePos) != 0 {
			if err := s.dao.RolePoBatchInsert(ctx, tx, rolePos); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Error("update role failed: ", err.Error())
		return err
	}
	return nil
}

//删除
func (s *Service) RoleDelete(ctx context.Context, param *vo.RoleDelRequest) (err error) {
	//查询主表数据
	role, err := s.dao.GetRoleById(ctx, param.Id)
	if err != nil {
		logger.Error("get role failed: ", err.Error())
		return err
	}
	if role == nil {
		err = errors.New(fmt.Sprintf("role is not exist: %d", param.Id))
		logger.Error(err.Error())
		return err
	}
	//事务
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		//删除主表数据
		if err := s.dao.DelRoleById(ctx, tx, param.Id); err != nil {
			return err
		}
		//删除关联表数据
		if err := s.dao.DelRolePoByRoleId(ctx, tx, param.Id); err != nil {
			return err
		}
		if err := s.dao.DelUserRoByRoleId(ctx, tx, param.Id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error("delete role failed: ", err.Error())
		return err
	}

	return nil
}