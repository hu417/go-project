package dao

import (
	"gorm.io/gorm"
	db2 "rbac-v1/db"
)

type Dao struct {
	db *gorm.DB
}

func NewDao() *Dao {
	db := db2.InitMysql()
	return &Dao{db:db}
}

func (d *Dao) DB() *gorm.DB {
	return d.db
}

func (d *Dao) Close() error {
	return db2.Close(d.db)
}