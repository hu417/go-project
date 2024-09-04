package routers

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func InitUserRouter() {
// 	//
// 	RegistRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
// 		//
// 		rgPublic.POST("/login", func(ctx *gin.Context) {
// 			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
// 				"msg": "login successful",
// 			})
// 		})

// 		//
// 		rgAuthUser := rgAuth.Group("user")
// 		rgAuthUser.GET("", func(ctx *gin.Context) {
// 			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
// 				"data": []map[string]any{
// 					{"id": 1, "name": "jd"},
// 					{"id": 2, "name": "tb"},
// 				},
// 			})
// 		})
// 		rgAuthUser.GET("/:id", func(ctx *gin.Context) {
// 			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
// 				"id":   1,
// 				"name": "jd",
// 			})
// 		})

// 	})

// }
// //
