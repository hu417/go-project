package router

import (
	"net/http"

	"viper-demo/api/v1"

	"github.com/gin-gonic/gin"
	
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 未知调用方式
	r.NoMethod(NoMethodJson)
	// 未知路由处理
	r.NoRoute(NoRouteJson)

	/*
		后台管理路由接口
	*/
	auth := r.Group("api/v1")
	// auth.Use(middleware.JwtToken())
	{
		// 用户模块的路由接口
		auth.POST("admin/user/add", v1.UserRegister)
		auth.POST("admin/user", v1.UserLogin)
		auth.GET("admin/users", v1.GetUsersList)
		auth.PUT("user/:id", v1.EditUserInfo)
		auth.DELETE("user/:id", v1.DeleteUserById)
		//修改密码
		auth.PUT("admin/changepw/:id", v1.EditUserPassword)
		// 分类模块的路由接口
		// auth.GET("admin/category", v1.GetCate)
		// auth.POST("category/add", v1.AddCategory)
		// auth.PUT("category/:id", v1.EditCate)
		// auth.DELETE("category/:id", v1.DeleteCate)
		// 文章模块的路由接口
		// auth.GET("admin/article/info/:id", v1.GetArtInfo)
		// auth.GET("admin/article", v1.GetArt)
		// auth.POST("article/add", v1.AddArticle)
		// auth.PUT("article/:id", v1.EditArt)
		// auth.DELETE("article/:id", v1.DeleteArt)
		// 上传文件
		// auth.POST("upload", v1.UpLoad)
		// 更新个人设置
		// auth.GET("admin/profile/:id", v1.GetProfile)
		// auth.PUT("profile/:id", v1.UpdateProfile)
		// 评论模块
		// auth.GET("comment/list", v1.GetCommentList)
		// auth.DELETE("delcomment/:id", v1.DeleteComment)
		// auth.PUT("checkcomment/:id", v1.CheckComment)
		// auth.PUT("uncheckcomment/:id", v1.UncheckComment)
	}
 
	/*
		前端展示页面接口
	*/
	// router := r.Group("api/v1")
	{
		// 用户信息模块
		// router.POST("user/add", v1.AddUser)
		// router.GET("user/:id", v1.GetUserInfo)
		// router.GET("users", v1.GetUsers)
 
		// 文章分类信息模块
		// router.GET("category", v1.GetCate)
		// router.GET("category/:id", v1.GetCateInfo)
 
		// 文章模块
		// router.GET("article", v1.GetArt)
		// router.GET("article/list/:id", v1.GetCateArt)
		// router.GET("article/info/:id", v1.GetArtInfo)
 
		// 登录控制模块
		// router.POST("login", v1.Login)
		// router.POST("loginfront", v1.LoginFront)
 
		// 获取个人设置信息
		// router.GET("profile/:id", v1.GetProfile)
 
		// 评论模块
		// router.POST("addcomment", v1.AddComment)
		// router.GET("comment/info/:id", v1.GetComment)
		// router.GET("commentfront/:id", v1.GetCommentListFront)
		// router.GET("commentcount/:id", v1.GetCommentCount)
	}

	return r
}

// 未知路由处理 返回json
func NoRouteJson(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
		"msg":  "path not found",
	})
}

// 未知调用方式 返回json
func NoMethodJson(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"code": http.StatusMethodNotAllowed,
		"msg":  "method not allowed",
	})
}
