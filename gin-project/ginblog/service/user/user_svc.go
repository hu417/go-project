package user

import "gorm.io/gorm"

type UserSvc struct {
	*gorm.DB
}

type UserSvcInterface interface {
	CheckUser(username string)
}

func NewUserSvc(db *gorm.DB) *UserSvc {
	return &UserSvc{
		DB: db,
	}
}
