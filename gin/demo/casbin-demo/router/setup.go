package router

import (
	"casbin-demo/api"
	"casbin-demo/global"
	"casbin-demo/middleware"
	"regexp"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// debug, release, test
	gin.SetMode("debug")
	r := gin.New()
	// 设置信任网络 []string
	// nil 为不计算，避免性能消耗，上线应当设置
	_ = r.SetTrustedProxies(nil)

	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// r.Static("/static", "./web/front/dist/static")
	// r.Static("/admin", "./web/admin/dist")
	// r.StaticFile("/favicon.ico", "/web/front/dist/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "front", nil)
	})

	// 注册路由
	r.POST("/login", api.LoginHandler)
	r.POST("/register", api.RegisterHandler)

	auth := r.Group("/api")
	{
		auth.Use(middleware.CasbinMiddleWare())
		auth.GET("/login", api.LoginHandler)
		auth.GET("/home", api.NewHomeController().Home)
		{
			// 用户以及组关系的增删改查
			// 获取所有用户以及关联的角色
			auth.GET("/casbin/users", api.GetUserList)
			// 获取所有角色组
			auth.GET("/casbin/roles", api.GetRoleList)
			// 获取角色下的用户
			auth.GET("/casbin/role-user", api.GetRoleUser)
			// 角色组中添加用户, 没有组默认创建
			auth.POST("/casbin/user-role", api.AddUserAtRole)
			// 查询用户角色
			auth.GET("/casbin/user-role", api.GetUserRole)
			// 更新用户角色
			auth.PUT("/casbin/user-role", api.UpdateUserRole)
			// 删除角色组中用户
			auth.DELETE("/casbin/user-role", api.DeleteUserAtRole)

		}
		{
			// 获取所有角色组的策略
			auth.GET("/casbin/rolepolicy", api.GetRolePolicyList)
			// 获取策略中的所有授权规则
			auth.GET("/casbin/rolepolicy/filtered", api.GetPolicyList)
			// 获取策略中的所有授权规则，可以指定角色字段筛选器
			auth.GET("/casbin/rolepolicy/filtered/named", api.GetFilteredPolicy)
			// 获取命名策略中的所有授权规则
			auth.GET("/casbin/rolepolicy/named", api.GetNamedPolicy)
			// 获取命名策略中的所有授权规则，可以指定字段过滤器
			auth.GET("/casbin/rolepolicy/named/filtered", api.GetFilteredNamedPolicy)
			// 创建角色组权限
			auth.POST("/casbin/rolepolicy", api.CreateRolePolicy)
			// 修改角色组权限
			auth.PUT("/casbin/rolepolicy", api.UpdateRolePolicy)
			// 删除角色组权限
			auth.DELETE("/casbin/rolepolicy", api.DeleteRolePolicy)
			// 验证用户权限
			auth.GET("/casbin/verify-user-role", api.VerifyUserRole)

		}

	}

	// 在应用启动时，加载所有路由到Casbin
	loadRoutesToCasbin(r)

	return r
}

// 加载所有路由到Casbin
func loadRoutesToCasbin(r *gin.Engine) {
	// 定义正则表达式，用于匹配登录和注册的路由
	loginRegexp := regexp.MustCompile(`/login`)
	registerRegexp := regexp.MustCompile(`/register`)
	homeRegexp := regexp.MustCompile(`/home`)

	// 遍历所有路由
	routes := r.Routes()
	for _, route := range routes {
		// 构建Casbin策略规则
		var policy []string

		// 根据路由路径判断是否为登录或注册接口
		if loginRegexp.MatchString(route.Path) || registerRegexp.MatchString(route.Path) || homeRegexp.MatchString(route.Path) {
			// 对于登录和注册接口，可以设置特定的策略，例如允许任何人访问
			policy = []string{
				"any-user", // 假设任何用户都有访问权限
				route.Method,
				route.Path,
			}
		} else {
			// 对于其他接口，可能需要更严格的权限控制
			policy = []string{
				"admin",
				route.Method,
				route.Path,
			}
		}

		// 将策略添加到Casbin
		if _, err := global.Enforcer.AddPolicy(policy); err != nil {
			panic(err)
		}
		// 保存策略
		global.Enforcer.SavePolicy()
	}
}
