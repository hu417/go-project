package back

import (
	"ginblog/middleware"

	"github.com/gin-gonic/gin"
)

func BackRouter(g *gin.RouterGroup) {
	g.Use(middleware.Jwt())
	{
		// 用户模块的路由接口
		// g.GET("admin/users", v1.GetUsers)
		// g.PUT("user/:id", v1.EditUser)
		// g.DELETE("user/:id", v1.DeleteUser)
		// //修改密码
		// g.PUT("admin/changepw/:id", v1.ChangeUserPassword)
		// // 分类模块的路由接口
		// g.GET("admin/category", v1.GetCate)
		// g.POST("category/add", v1.AddCategory)
		// g.PUT("category/:id", v1.EditCate)
		// g.DELETE("category/:id", v1.DeleteCate)
		// // 文章模块的路由接口
		// g.GET("admin/article/info/:id", v1.GetArtInfo)
		// g.GET("admin/article", v1.GetArt)
		// g.POST("article/add", v1.AddArticle)
		// g.PUT("article/:id", v1.EditArt)
		// g.DELETE("article/:id", v1.DeleteArt)
		// // 上传文件
		// g.POST("upload", v1.UpLoad)
		// // 更新个人设置
		// g.GET("admin/profile/:id", v1.GetProfile)
		// g.PUT("profile/:id", v1.UpdateProfile)
		// // 评论模块
		// g.GET("comment/list", v1.GetCommentList)
		// g.DELETE("delcomment/:id", v1.DeleteComment)
		// g.PUT("checkcomment/:id", v1.CheckComment)
		// g.PUT("uncheckcomment/:id", v1.UncheckComment)
	}

}
