package service

import (
	"api-demo/app/model"
	"api-demo/internal/global"
)

type adminService struct {
	baseService
}

func (a *adminService) Profile() (model *model.Admin) {
	
	// if err := a.db.First(&model).Error; err != nil {
	// 	global.Logger.Sugar().Errorf("err => %v\n", err.Error())
	// 	return nil
	// }

	global.DB.First(&model)

	return model
}
