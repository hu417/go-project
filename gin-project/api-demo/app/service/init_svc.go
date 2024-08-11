package service

import "api-demo/internal/global"

var (
	BaseService = &baseService{db: global.DB}

	AdminService = &adminService{baseService: *BaseService}
)
