package main

import (
	"go-admin/models"
	"go-admin/router"
)

func main() {
	// 初始化gorm.db
	models.NewGormDB()
	r := router.App()
	r.Run(":8081")
}
