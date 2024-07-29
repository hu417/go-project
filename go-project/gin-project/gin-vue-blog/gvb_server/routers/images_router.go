package routers

import (
	"gvb_server/api"
)

func (r RouterGroup) ImagesGroup() {
	imagesApi := api.ApiGroupApp.ImagesApi

	// 绑定请求url参数
	r.POST("images/upload", imagesApi.ImagesUploadView)
}
