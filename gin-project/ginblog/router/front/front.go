package front

import (
	"ginblog/api/v1/user"

	"github.com/gin-gonic/gin"
)

func FrontRuter(g *gin.RouterGroup) {
	{
		// 用户信息模块
		g.POST("user/add", user.NewUser().Add)
		// g.GET("user/:id", user.GetInfo)
		// g.GET("users", user.GetUsers)

		// // 文章分类信息模块
		// g.GET("category", v1.GetCate)
		// g.GET("category/:id", v1.GetCateInfo)

		// // 文章模块
		// g.GET("article", v1.GetArt)
		// g.GET("article/list/:id", v1.GetCateArt)
		// g.GET("article/info/:id", v1.GetArtInfo)

		// // 登录控制模块
		// g.POST("login", v1.Login)
		// g.POST("loginfront", v1.LoginFront)

		// // 获取个人设置信息
		// g.GET("profile/:id", v1.GetProfile)

		// // 评论模块
		// g.POST("addcomment", v1.AddComment)
		// g.GET("comment/info/:id", v1.GetComment)
		// g.GET("commentfront/:id", v1.GetCommentListFront)
		// g.GET("commentcount/:id", v1.GetCommentCount)
	}
}
