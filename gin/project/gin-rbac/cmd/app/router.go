package app

import (
	"gin-rbac/common/response"
	"gin-rbac/db/dao"
	"gin-rbac/handler"
	"gin-rbac/middleware"
	useRouter "gin-rbac/router"
	"gin-rbac/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// InitRouter 初始化并配置Gin框架的路由。
//
// 描述:
// 函数用于初始化Gin框架的路由，并注册相关的中间件和服务。
// 它还设置了健康检查端点和swagger相关的路由。
func InitRouter(db *gorm.DB, rds *redis.Client) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware())                                    // 使用CORS中间件以支持跨域请求
	router.Use(middleware.LoggerMiddleware(), middleware.RecoveryMiddleware()) // 使用日志记录中间件
	router.Use(middleware.JWTAuthMiddleware(router, db, rds))                  // 使用JWT认证中间件

	// 404
	router.NoRoute(func(ctx *gin.Context) { // 这里只是演示，不要在生产环境中直接返回HTML代码
		ctx.String(http.StatusNotFound, "<h1>404 Page Not Found</h1>")
	})

	// 405
	router.NoMethod(func(ctx *gin.Context) {
		ctx.String(http.StatusMethodNotAllowed, "method not allowed")
	})

	// 注册swagger静态文件路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	/*
			// swagger路由;添加接口认证
		    authorized := r.Group("/swagger", gin.BasicAuth(gin.Accounts{
		        "admin": "666666",
		    }))
		    authorized.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	*/

	// 健康检查端点
	router.GET("/api/v1/health", func(c *gin.Context) {
		response.Ok(c, 0, "The service is healthy", nil)
	})

	// 创建用户服务实例
	userService := service.NewUserService(dao.NewUserDAO(db, &gin.Context{}))
	userHandler := handler.NewUserHandler(userService)

	// 创建角色服务实例
	roleService := service.NewRoleService(dao.NewRoleDAO(db, &gin.Context{}))
	roleHandler := handler.NewRoleHandler(roleService)

	// 创建权限服务实例
	permissionService := service.NewPermissionService(dao.NewPermissionDAO(db, &gin.Context{}))
	permissionHandler := handler.NewPermissionHandler(permissionService)

	// 创建用户角色服务实例
	userRoleService := service.NewUserRoleService(
		dao.NewUserDAO(db, &gin.Context{}),
		dao.NewRoleDAO(db, &gin.Context{}),
		dao.NewUserRoleDAO(db, &gin.Context{}))
	userRoleHandler := handler.NewUserRoleHandler(userRoleService)

	// 创建角色权限服务实例
	rolePermissionService := service.NewRolePermissionService(
		dao.NewRoleDAO(db, &gin.Context{}),
		dao.NewPermissionDAO(db, &gin.Context{}),
		dao.NewRolePermissionDAO(db, &gin.Context{}),
	)
	rolePermissionHandler := handler.NewRolePermissionHandler(rolePermissionService)

	// 创建系统服务实例
	systemService := service.NewSystemService()
	systemHandler := handler.NewSystemHandler(systemService)

	// 创建图片服务实例
	imageService := service.NewImageService(dao.NewImageDAO(db, &gin.Context{}))
	imageHandler := handler.NewImageHandler(imageService)

	// 定义API组
	api := router.Group("/api/v1")
	{
		useRouter.RegisterUserRoutes(api, userHandler)                     // 注册用户相关的路由
		useRouter.RegisterRoleRoutes(api, roleHandler)                     // 注册角色相关的路由
		useRouter.RegisterPermissionRoutes(api, permissionHandler)         // 注册权限相关的路由
		useRouter.RegisterUserRoleRoutes(api, userRoleHandler)             // 注册用户角色相关的路由
		useRouter.RegisterRolePermissionRoutes(api, rolePermissionHandler) // 注册角色权限相关的路由
		useRouter.RegisterSystemRoutes(api, systemHandler, router)         // 注册系统相关的路由
		useRouter.RegisterImageRoutes(api, imageHandler)                   // 注册图片相关的路由
	}

	return router
}
