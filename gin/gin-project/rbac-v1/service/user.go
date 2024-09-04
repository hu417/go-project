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

//获取用户的role和power
func (s *Service) GetUserRoleAndPower(ctx context.Context, param *vo.GetUserRoleAndPowerRequest) (ret *vo.GetUserRoleAndPowerResponse, err error) {
	ret = &vo.GetUserRoleAndPowerResponse{
		Roles:  []string{},
		Powers: []string{},
	}
	//查询关联数据-role
	userRos, err := s.dao.GetUserRoByUserIds(ctx, []uint{param.UserId})
	if err != nil {
		logger.Error("get userRole failed: ", err.Error())
		return
	}
	roIds := []uint{}
	for _, userRo := range userRos {
		roIds = append(roIds, userRo.RoleId)
	}
	if len(roIds) != 0 {
		roles, err := s.dao.GetRoleByIds(ctx, roIds)
		if err != nil {
			logger.Error("get roles failed: ", err.Error())
			return nil, err
		}
		for _, role := range roles {
			ret.Roles = append(ret.Roles, role.Code)
		}
	}
	//查询关联数据-power
	userPos, err := s.dao.GetUserPoByUserIds(ctx, []uint{param.UserId})
	if err != nil {
		logger.Error("get userPower failed: ", err.Error())
		return nil, err
	}
	poIds := []uint{}
	for _, userPo := range userPos {
		poIds = append(poIds, userPo.PowerId)
	}
	if len(poIds) != 0 {
		powers, err := s.dao.GetPowerByIds(ctx, poIds)
		if err != nil {
			logger.Error("get powers failed: ", err.Error())
			return nil, err
		}
		for _, power := range powers {
			ret.Powers = append(ret.Powers, power.Code)
		}
	}

	return ret, nil
}

//列表
func (s *Service) GetUserList(ctx context.Context, param *vo.UserListRequest) (ret *vo.UserBoListResponse, err error) {
	ret = &vo.UserBoListResponse{}
	userList, err := s.dao.GetUserList(ctx, param)
	if err != nil {
		logger.Error("get user list failed: ", err.Error())
		return nil, err
	}
	if userList.List == nil {
		return ret, nil
	}
	userIds := []uint{}
	for _, user := range userList.List {
		userIds = append(userIds, user.Id)
	}
	userRos, err := s.dao.GetUserRoByUserIds(ctx, userIds)
	if err != nil {
		logger.Error("get userRo list failed: ", err.Error())
		return nil, err
	}
	userPos, err := s.dao.GetUserPoByUserIds(ctx, userIds)
	if err != nil {
		logger.Error("get userPo list failed: ", err.Error())
		return nil, err
	}

	if userRos == nil && userPos == nil {
		userBos := []*bo.User{}
		for _, user := range userList.List {
			userBo := &bo.User{
				Id:        user.Id,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				IsDeleted: user.IsDeleted,
				Name:      user.Name,
				Username:  user.Username,
				Mail:      user.Mail,
				Phone:     user.Phone,
			}
			userBos = append(userBos, userBo)
		}

		return &vo.UserBoListResponse{
			Total: userList.Total,
			List:  userBos,
		}, nil
	}

	userRoMp := map[uint][]uint{}
	roIds := []uint{}
	for _, userRo := range userRos {
		if _, ok := userRoMp[userRo.UserId]; !ok {
			userRoMp[userRo.UserId] = []uint{}
		}
		roIds = append(roIds, userRo.RoleId)
		userRoMp[userRo.UserId] = append(userRoMp[userRo.UserId], userRo.RoleId)
	}
	userPoMp := map[uint][]uint{}
	poIds := []uint{}
	for _, userPo := range userPos {
		if _, ok := userPoMp[userPo.UserId]; !ok {
			userPoMp[userPo.UserId] = []uint{}
		}
		poIds = append(poIds, userPo.PowerId)
		userPoMp[userPo.UserId] = append(userPoMp[userPo.UserId], userPo.PowerId)
	}

	roMp := map[uint]*po.Role{}
	if len(roIds) != 0 {
		roles, err := s.dao.GetRoleByIds(ctx, roIds)
		if err != nil {
			logger.Error("get roles failed: ", err.Error())
			return nil, err
		}
		for _, role := range roles {
			roMp[role.Id] = role
		}
	}
	poMp := map[uint]*po.Power{}
	if len(poIds) != 0 {
		powers, err := s.dao.GetPowerByIds(ctx, poIds)
		if err != nil {
			logger.Error("get powers failed: ", err.Error())
			return nil, err
		}
		for _, power := range powers {
			poMp[power.Id] = power
		}
	}

	userBos := []*bo.User{}
	for _, user := range userList.List {
		userBo := &bo.User{
			Id:        user.Id,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			IsDeleted: user.IsDeleted,
			Name:      user.Name,
			Username:  user.Username,
			Mail:      user.Mail,
			Phone:     user.Phone,
			Roles:     []*po.Role{},
			Powers:    []*po.Power{},
		}
		if _, ok := userRoMp[user.Id]; ok {
			for _, roleId := range userRoMp[user.Id] {
				userBo.Roles = append(userBo.Roles, roMp[roleId])
			}
		}
		if _, ok := userPoMp[user.Id]; ok {
			for _, powerId := range userPoMp[user.Id] {
				userBo.Powers = append(userBo.Powers, poMp[powerId])
			}
		}
		userBos = append(userBos, userBo)
	}
	return &vo.UserBoListResponse{
		Total: userList.Total,
		List:  userBos,
	}, nil
}

//新增
func (s *Service) UserCreate(ctx context.Context, param *vo.UserCreateRequest) (err error) {
	roIds := []uint{}
	poIds := []uint{}
	users := []*po.User{}
	for _, user := range param.Users {
		if ok := s.dao.IsUserExist(ctx, user.Username); ok {
			err = errors.New(fmt.Sprintf("create user failed, duplicate data: %s", user.Username))
			logger.Error(err.Error())
			return err
		}
		users = append(users, user.User)
		for _, role := range user.Roles {
			roIds = append(roIds, role.Id)
		}
		for _, power := range user.Powers {
			poIds = append(poIds, power.Id)
		}
	}
	if len(roIds) != 0 {
		roIds = common.RemoveDuplicates(roIds)
		for _, roId := range roIds {
			if ok := s.dao.IsRoleExistById(ctx, roId); !ok {
				err = errors.New(fmt.Sprintf("create user failed, role not exist: %d", roId))
				logger.Error(err.Error())
				return err
			}
		}
	}
	if len(poIds) != 0 {
		poIds = common.RemoveDuplicates(poIds)
		for _, poId := range poIds {
			if ok := s.dao.IsPowerExistById(ctx, poId); !ok {
				err = errors.New(fmt.Sprintf("create user failed, power not exist: %d", poId))
				logger.Error(err.Error())
				return err
			}
		}
	}

	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.dao.UserBatchInsert(ctx, tx, users); err != nil {
			return err
		}
		userRos := []*po.UserRo{}
		userPos := []*po.UserPo{}
		for i := 0; i < len(users); i ++ {
			if len(param.Users[i].Roles) != 0 {
				for _, role := range param.Users[i].Roles {
					userRos = append(userRos, &po.UserRo{
						UserId: users[i].Id, //不可以使用param.Users[i].Id,因为我们使用users进行创建的，创建完成后，这个数据中会有id
						RoleId: role.Id,
					})
				}
			}
			if len(param.Users[i].Powers) != 0 {
				for _, power := range param.Users[i].Powers {
					userPos = append(userPos, &po.UserPo{
						UserId:  users[i].Id, //不可以使用param.Users[i].Id,因为我们使用users进行创建的，创建完成后，这个数据中会有id
						PowerId: power.Id,
					})
				}
			}
		}

		if len(userRos) != 0 {
			if err := s.dao.UserRoBatchInsert(ctx, tx, userRos); err != nil {
				return err
			}
		}
		if len(userPos) != 0 {
			if err := s.dao.UserPoBatchInsert(ctx, tx, userPos); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Error("create user failed: ", err)
		return err
	}
	return nil
}

//更新
func (s *Service) UserUpdate(ctx context.Context, param *bo.UserCreate) (err error) {
	user, err := s.dao.GetUserById(ctx, param.Id)
	if err != nil {
		logger.Error("get user failed: ", err.Error())
		return err
	}
	if user == nil {
		err = errors.New("no such user data")
		logger.Error(err.Error())
		return err
	}
	param.CreatedAt = user.CreatedAt
	param.UpdatedAt = time.Now().Local()

	userRos := []*po.UserRo{}
	if len(param.Roles) != 0 {
		for _, role := range param.Roles {
			if ok := s.dao.IsRoleExistById(ctx, role.Id); !ok {
				err = errors.New(fmt.Sprintf("create user failed, role not exist: %d", role.Id))
				logger.Error(err.Error())
				return err
			}
			userRos = append(userRos, &po.UserRo{
				UserId: param.Id,
				RoleId: role.Id,
			})
		}
	}
	userPos := []*po.UserPo{}
	if len(param.Powers) != 0 {
		for _, power := range param.Powers {
			if ok := s.dao.IsPowerExistById(ctx, power.Id); !ok {
				err = errors.New(fmt.Sprintf("create user failed, power not exist: %d", power.Id))
				logger.Error(err.Error())
				return err
			}
			userPos = append(userPos, &po.UserPo{
				UserId: param.Id,
				PowerId: power.Id,
			})
		}
	}

	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.dao.UpdateUserById(ctx, tx, param.User); err != nil {
			return err
		}
		//role
		if err := s.dao.DelUserRoByUserId(ctx, tx, param.Id); err != nil {
			return err
		}
		if len(userRos) != 0 {
			if err := s.dao.UserRoBatchInsert(ctx, tx, userRos); err != nil {
				return err
			}
		}
		//power
		if err := s.dao.DelUserPoByUserId(ctx, tx, param.Id); err != nil {
			return err
		}
		if len(userPos) != 0 {
			if err := s.dao.UserPoBatchInsert(ctx, tx, userPos); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Error("update user failed: ", err.Error())
		return err
	}

	return nil
}

//删除
func (s *Service) UserDelete(ctx context.Context, param *vo.UserDelRequest) (err error) {
	user, err := s.dao.GetUserById(ctx, param.Id)
	if err != nil {
		logger.Error("get user failed: ", err.Error())
		return err
	}
	if user == nil {
		err = errors.New(fmt.Sprintf("user is not exist: %d", param.Id))
		logger.Error(err.Error())
		return err
	}

	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.dao.DelUserById(ctx, tx, param.Id); err != nil {
			return err
		}
		if err := s.dao.DelUserRoByUserId(ctx, tx, param.Id); err != nil {
			return err
		}
		if err := s.dao.DelUserPoByUserId(ctx, tx, param.Id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error("delete user failed: ", err.Error())
		return err
	}

	return nil
}