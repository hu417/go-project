package user

import (
	"ginblog/dao/user"
	"ginblog/model"
)

func (svc *UserSvc) CheckUser(username string) bool {
	u := model.User{
		Username: username,
	}
	return user.NewUserDao(svc.DB).FindUserByName(&u)
}
