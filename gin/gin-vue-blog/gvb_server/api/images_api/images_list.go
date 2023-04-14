package images_api

import (
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"

	"github.com/gin-gonic/gin"
)

// 分页相关数据参数
type Page struct {
	Page  int    `form:"page" description:"第几页"`
	Key   string `form:"key"`
	Limit int    `form:"limit" description:"显示条数"`
	Sort  string `form:"sort" description:"排序"`
}

func (ImagesApi) ImahesListView(c *gin.Context) {

	var cr Page
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	// 返回数据条数
	var imagesList []models.BannerModel
	count := global.DB.Find(&imagesList).RowsAffected
	global.Logger.Infof("数据总数为: %d", count)

	// 分页数据
	offset := (cr.Page - 1) * cr.Limit
	if offset < 0 {
		offset = 0
	}
	global.DB.Limit(cr.Limit).Offset(offset).Find(&imagesList)

	res.OkWithData(gin.H{
		"count": count,
		"list":  imagesList,
	}, c)

}
