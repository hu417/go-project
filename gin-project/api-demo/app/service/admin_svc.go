package service

import (
	"api-demo/app/model"
	"api-demo/internal/global"
)

type adminService struct {
	baseService
}

func (a *adminService) Profile() (model model.Admin) {

	// a.db.First(&model)
	global.DB.First(&model)
	return model
}
