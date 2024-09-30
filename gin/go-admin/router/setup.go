package router

import (
	v1 "go-admin/api/v1"
	"go-admin/global"
	"go-admin/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(
		// 跨域中间件
		middleware.Cors(),
		// 日志中间件
		middleware.LoggerToDb(),
	)

	// 使用session中间件
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Secure:   true,
		SameSite: 4,
		Path:     "/",
		MaxAge:   3600,
	})
	r.Use(sessions.Sessions("admin-session", store))

	// 设置静态资源路由
	r.Static("/uploadFile", global.StaticResource)

	// 根据用户名和密码登录路由
	r.POST("/login/password", v1.LoginPassword)

	// 登录信息认证
	api := r.Group("/api/v1").Use(middleware.LoginAuthCheck())
	{
		// 上传图片文件
		api.POST("/upload/file", v1.UploadFile)

		// 管理员
		{

			// 管理员列表
			api.GET("/user", v1.GetUserList)
			// 新增管理员
			api.POST("/user", v1.AddUser)
			// 管理员详情
			api.GET("/user/detail/:id", v1.GetUserDetail)
			// 修改管理员
			api.PUT("/user", v1.UpdateUser)
			// 删除管理员
			api.DELETE("/user/:id", v1.DeleteUser)
			// 更新个人信息
			api.PUT("/user/updateInfo", v1.UpdateInfo)
			// 发送邮件
			api.GET("/user/sendEmail", v1.SendEmail)
			// 校验邮箱验证码
			api.GET("/user/verifyCode", v1.VerifyCode)
			// 更新邮箱
			api.PUT("/user/updateEmail", v1.UpdateEmail)
			// 更新个人密码
			api.PUT("/user/updatePwd", v1.UpdatePwd)
		}

		// 角色管理
		{
			// 角色列表
			api.GET("/role", v1.GetRoleList)
			// 新增角色
			api.POST("/role", v1.AddRole)
			// 修改角色的管理员身份
			api.PATCH("/role/:id/:isAdmin", v1.PatchRoleAdmin)
			// 角色详情
			api.GET("/role/detail/:id", v1.GetRoleDetail)
			// 修改角色
			api.PUT("/role", v1.UpdateRole)
			// 删除角色
			api.DELETE("/role/:id", v1.DeleteRole)
			// 获取所有角色
			api.GET("/role/all", v1.AllRole)
		}

		// 菜单功能管理
		{
			// 菜单列表
			api.GET("/menu", v1.GetMenuList)
			// 新增菜单
			api.POST("/menu", v1.AddMenu)
			// 修改菜单
			api.PUT("/menu", v1.UpdateMenu)
			// 删除菜单
			api.DELETE("/menu/:id", v1.DeleteMenu)
			// 菜单功能管理 end
		}

		// 日志管理
		{
			api.GET("/log", v1.GetLogList)
		}
	}

	return r
}
