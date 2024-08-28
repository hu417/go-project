package router

import (
    "github.com/gin-gonic/gin"
    "net/http/pprof"
    "pprof/api/v1"
)

func RegistRoute(engine *gin.Engine) {
    Heap := engine.Group("/heap")
    {
        Heap.GET("/log", v1.Heap.Log)
    }

    systemPprof(engine)
}

// pprof系统性能分析
func systemPprof(engine *gin.Engine) {
    pprofAPI := engine.Group("/pprof")
    {
        pprofAPI.GET("/", gin.WrapF(pprof.Index))
        pprofAPI.GET("/cmdline", gin.WrapF(pprof.Cmdline))
        pprofAPI.GET("/profile", gin.WrapF(pprof.Profile))
        pprofAPI.Any("/symbol", gin.WrapF(pprof.Symbol))
        pprofAPI.GET("/trace", gin.WrapF(pprof.Trace))
        pprofAPI.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
        pprofAPI.GET("/block", gin.WrapH(pprof.Handler("block")))
        pprofAPI.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
        pprofAPI.GET("/heap", gin.WrapH(pprof.Handler("heap")))
        pprofAPI.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
        pprofAPI.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
    }
}
