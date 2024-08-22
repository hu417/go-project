package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 用户登陆
func UserLogin(c *gin.Context) {
	//  获取客户端是否携带cookie
	_, err := c.Cookie("key_cookie")
	if err != nil {
		c.SetCookie("key_cookie", "value_cookie", //  参数1、2： key & value
			10,          //  参数3： 生存时间（秒）
			"/",         //  参数4： 所在目录
			"localhost", //  参数5： 域名
			false,       //  参数6： 安全相关 - 是否智能通过https访问
			true,        //  参数7： 安全相关 - 是否允许别人通过js获取自己的cookie
		)
		c.String(200, "login success")
		return
	}
	c.String(200, "already login")
}

// @Tags user
// @Summary 用户注册
func CreateUserforForm(c *gin.Context) {
	//  设置默认值
	types := c.DefaultPostForm("type", "post")

	//  使用PostForm获取请求参数
	username := c.PostForm("username")
	password := c.PostForm("password")

	//  返回请求参数
	c.JSON(200, gin.H{
		"types":    types,
		"username": username,
		"password": password,
	})
}

func CreateUserforJson(c *gin.Context) {
	type User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//  将JSON格式请求参数绑定到结构体上
	var user User
	//  c.GetRawData() 表示raw类型的数据,可以理解为特定格式的字符串，如json字符串
	if err := c.BindJSON(&user); err != nil {
		//  返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//  根据req的content type 自动推断如何绑定,form/json/xml等格式
	//  如果发送的不是json格式，那么输出：  "error": "invalid character '-' in numeric literal"
	//  if err := c.ShouldBind(&user); err != nil {
	//  	c.JSON(400, gin.H{"error": err.Error()})
	//  	return
	//  }

	//  返回请求参数
	c.JSON(200, gin.H{
		"username": user.Username,
		"password": user.Password,
	})
}

func GetUserForHeader(c *gin.Context) {
	// 使用Request获取请求头部参数
	username := c.Request.Header.Get("username")
	password := c.Request.Header.Get("password")

	// 返回请求参数
	c.JSON(200, gin.H{
		"username": username,
		"password": password,
	})
}

func GetUserForQueryById(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func GetMathByUrl(c *gin.Context) {
	path := c.Param(".html")
	c.JSON(http.StatusOK, gin.H{
		"path": path,
	})
}

func GetMathByUrlForDefault(c *gin.Context) {
	// 指定默认值
	//  DefaultQuery()若参数不村则，返回默认值，Query()若不存在，返回空串
	//  name := c.DefaultQuery("name", "normal")
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"name": name,
	})
}
