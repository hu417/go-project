package app

import (
	"gin-rbac/global"

	"github.com/gin-gonic/gin"
)

// App 应用
type App struct {
	Engine *gin.Engine
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{
		Engine: gin.Default(),
	}
}

// Initialize 初始化应用
func (a *App) Initialize() {
	// 初始化配置
	// 在这里注册路由，传入 JWT 和 DB 实例
	a.Engine = InitRouter(global.DB, global.Redis)
	// 注册路由
	//router.SetupRoutes(a.Engine)
}

// Run 启动服务
func (a *App) Run(addr string) error {
	return a.Engine.Run(addr)
}
