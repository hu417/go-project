package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCookie(c *gin.Context) {
	// 获取客户端是否携带cookie
	cookie, err := c.Cookie("key_cookie")
	if err != nil {
		cookie = "cookie"
		c.SetCookie("key_cookie", "value_cookie", //  参数1、2： key & value
			60,          //  参数3： 生存时间(秒);如果是-1,则表示删除
			"/",         //  参数4： 所在目录
			"localhost", //  参数5： 域名
			false,       //  参数6： 安全相关 - 是否智能通过https访问
			true,        //  参数7： 安全相关 - 是否允许别人通过js获取自己的cookie
		)
	}
	fmt.Printf("cookie的值是： %s\n", cookie)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": cookie,
	})

}
