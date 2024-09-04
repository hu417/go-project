package service

import "rbac-v1/dao"

var (
	srv *Service
)

type Service struct {
	dao *dao.Dao
}

func New() (s *Service) {
	newDao := dao.NewDao()

	srv = &Service{dao:newDao}

	return srv
}

func (s *Service) Dao() *dao.Dao {
	return s.dao
}

func Srv() *Service {
	return srv
}

func (s *Service) Close() error {
	return s.dao.Close()
}

