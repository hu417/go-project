package flag

import (
	"gvb_server/global"
	"gvb_server/models"
)

func Makmigrations() {
	global.Logger.Info("Migrating xxxxx")

	var err error
	global.DB.SetupJoinTable(&models.UserModel{}, "CollectsModels", &models.UserCollectModel{})
	global.DB.SetupJoinTable(&models.MenuModel{}, "Banners", &models.MenuBannerModel{})
	// 生成四张表的表结构
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.BannerModel{},
		&models.TagModel{},
		&models.MessageModel{},
		&models.AdvertModel{},
		&models.UserModel{},
		&models.CommentModel{},
		&models.ArticleModel{},
		&models.MenuModel{},
		&models.MenuBannerModel{},
		&models.FadeBackModel{},
		&models.LoginDataModel{},
	)
	if err != nil {
		global.Logger.Error("[ error ] 生成数据库表结构失败")
		return
	}
	global.Logger.Info("[ success ] 生成数据库表结构成功！")
}
