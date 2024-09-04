package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 计算时间
func CalcTimeMiddleWare() gin.HandlerFunc {
	fmt.Println(1)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		// 统计时间
		since := time.Since(start)
		fmt.Println("程序用时：", since)
	}
}
