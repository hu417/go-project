
## jwt

### 基本信息

`JWT` 由三个部分组成，它们之间用 `.` 分隔，格式如下：`Header.Payload.Signature` → `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJQcm9ncmFtbWVyIiwiaXNzIjoi56iL5bqP5ZGY6ZmI5piO5YuHIiwic3ViIjoiY2hlbm1pbmd5b25nLmNuIn0.uRnH-rUb7lsZtQ11o8wXjIOJnIMBxszkvU1gY6hCGjo`，下面对每个部分的进行简要介绍：

1. `Header（头部）`：`Hedaer` 部分用于描述该 `JWT` 的基本信息，比如其类型（通常是 `JWT`）以及所使用的签名算法（如 `HMAC SHA256` 或 `RSA`）。

2. `Payload（负载）`： `Payload` 部分包含所传递的声明。声明是关于实体（通常是用户）和其他数据的语句。声明可以分为三种类型：**注册声明**、**公共声明** 和 **私有声明**。
	- **注册声明**：这些声明是预定义的，非必须使用的但被推荐使用。官方标准定义的注册声明有 7 个：
		|Claim（声明）|含义|
		|---|---|
		|iss(Issuer)|发行者，标识 JWT 的发行者。|
		|sub(Subject)|主题，标识 JWT 的主题，通常指用户的唯一标识|
		|aud(Audience)|观众，标识 JWT的接收者|
		|exp(Expiration Time)|过期时间。标识 JWT 的过期时间，这个时间必须是将来的|
		|nbf(Not Before)|不可用时间。在此时间之前，JWT 不应被接受处理|
		|iat(Issued At)|发行时间，标识 JWT 的发行时间|
		|jti(JWT ID)|JWT 的唯一标识符，用于防止 JWT 被重放（即重复使用）|

	- **公共声明**：可以由使用 `JWT` 的人自定义，但为了避免冲突，任何新定义的声明都应已在 IANA JSON Web Token Registry 中注册或者是一个 **公共名称**，其中包含了碰撞防抗性名称（`Collision-Resistant Name`）。

	- **私有声明**：发行和使用 `JWT` 的双方共同商定的声明，区别于 **注册声明** 和 **公共声明**。

3. `Signature（签名）`：为了防止数据篡改，将头部和负载的信息进行一定算法处理，加上一个密钥，最后生成签名。如果使用的是 `HMAC SHA256` 算法，那么签名就是将编码后的头部、编码后的负载拼接起来，通过密钥进行`HMAC SHA256` 运算后的结果。

### 使用

#### dgrijalva: jwt-go

安装依赖
```bash
go get -u github.com/dgrijalva/jwt-go

```

##### 使用案例

```go

// middleware.go
package middleware
 
import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)
//Jwtkey 秘钥 可通过配置文件配置
var Jwtkey = []byte("blog_jwt_key")
 
type MyClaims struct {
	UserId int `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}
 
// CreateToken 生成token
func CreateToken(userId int,userName string) (string,error) {
	expireTime := time.Now().Add(2*time.Hour) //过期时间
	nowTime := time.Now() //当前时间
	claims := MyClaims{
		UserId: userId,
		UserName: userName,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt:expireTime.Unix(), //过期时间戳
			IssuedAt: nowTime.Unix(), //当前时间戳
			Issuer: "blogLeo", //颁发者签名
			Subject: "userToken", //签名主题
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return tokenStruct.SignedString(Jwtkey)
}
 
// CheckToken 验证token
func CheckToken(token string) (*MyClaims,bool) {
	tokenObj,_ := jwt.ParseWithClaims(token,&MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Jwtkey,nil
	})
	if key,_ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return key,true
	}else{
		return nil,false
	}
}
 
// JwtMiddleware jwt中间件
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头中获取token
		tokenStr := c.Request.Header.Get("Authorization")
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK,gin.H{"code":0, "msg":"用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//token格式错误
		tokenSlice := strings.SplitN(tokenStr," ",2)
		if len(tokenSlice) != 2 && tokenSlice[0] != "Bearer" {
			c.JSON(http.StatusOK,gin.H{"code":0, "msg":"token格式错误"})
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck,ok := CheckToken(tokenSlice[1])
		if !ok {
			c.JSON(http.StatusOK,gin.H{"code":0, "msg":"token不正确"})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK,gin.H{"code":0, "msg":"token过期"})
			c.Abort() //阻止执行
			return
		}
		c.Set("username",tokenStruck.UserName)
		c.Set("user_id",tokenStruck.UserId)
 
		c.Next()
	}
}

```



#### golang-jwt: v5

安装依赖
```go
go get -u github.com/golang-jwt/jwt/v5

```

##### 生成token

```go
package main

import (
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// 自定义Claim结构体
type MyCustomClaims struct {
   UserID     int
   Username   string
   GrantScope string
   jwt.RegisteredClaims
}

// 签名密钥
const sign_key = "hello jwt"

// 随机字符串
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(str_len int) string {
   rand_bytes := make([]rune, str_len)
   for i := range rand_bytes {
      rand_bytes[i] = letters[rand.Intn(len(letters))]
   }
   return string(rand_bytes)
}


// 生成token
func GenerateJwt(sign_key any) (string, error) {
   // 初始化claim
   claim := MyCustomClaims{
      UserID:     000001,
      Username:   "Tom",
      GrantScope: "read_user_info",
      RegisteredClaims: jwt.RegisteredClaims{
         Issuer:    "Auth_Server",                                   // 签发者
         Subject:   "Tom",                                           // 签发对象
         Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      //签发受众
         ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   //过期时间
         NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
         IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
         ID:        randStr(10),                                     // wt ID, 类似于盐值
      },
   }

  // 生成token
   return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(sign_key))

}

func main() {
	/* jwtKey := make([]byte, 32) // 生成32字节（256位）的密钥
	if _, err := rand.Read(jwtKey); err != nil {
		panic(err)
	}
	*/
	jwtStr, err := GenerateJwt(sign_key)
	if err != nil {
		panic(err)
	}
	fmt.Println(jwtStr)
}

```

##### 解析token

###### Parse 函数

```go

package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

/*
自定义的 ParseJwt 函数，它负责解析 JWT 字符串，并根据验证结果返回 Claims 数据和一个可能的存在的错误。ParseJwt 函数内部利用 jwt.Parse 解析 JWT 字符串。解析后，函数检查得到的 token 对象的 Valid 属性以确认 Claims 是否有效。有效性检查包括但不限于验证签名、检查 token 是否过期。如果 token 通过所有验证，函数返回 Claims 数据；如果验证失败（如签名不匹配或 token 已过期），则返回错误
*/
func ParseJwt(key any, jwtStr string, options ...jwt.ParserOption) (jwt.Claims, error) {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}, options...)
	if err != nil {
		returnnil, err
	}
	// 校验 Claims 对象是否有效，基于 exp（过期时间），nbf（不早于），iat（签发时间）等进行判断（如果有这些声明的话）。
	if !token.Valid {
		returnnil, errors.New("invalid token")
	}
	return token.Claims, nil
}

func main() {
	/* jwtKey := make([]byte, 32) // 生成32字节（256位）的密钥
	if _, err := rand.Read(jwtKey); err != nil {
		panic(err)
	}
	*/

	jwtStr, err := GenerateJwt(sign_key)
	if err != nil {
		panic(err)
	}

	// 解析 jwt
	claims, err := ParseJwt(sign_key, jwtStr, jwt.WithExpirationRequired())
	if err != nil {
		panic(err)
	}
	fmt.Println(claims)
}
```

###### ParseWithClaims 函数

```go

package main

import (
    "crypto/rand"
    "errors"
    "fmt"
    "github.com/golang-jwt/jwt/v5"
)

/*

ParseJwtWithClaims 函数与 ParseJwt 函数功能类似，都是负责解析 JWT 字符串，并根据验证结果返回 Claims 数据和一个可能的存在的错误。不同之处在于，ParseJwtWithClaims 函数内部使用了 jwt.ParseWithClaims 函数来解析 JWT 字符串，这额外要求我们提供一个 Claims 实例来接收解析后的 claims 数据。在此示例中，通过 jwt.MapClaims 提供了这一实例。
*/
func ParseJwtWithClaims(sign_key any, jwtStr string, options ...jwt.ParserOption) (*MyCustomClaims, error) {

    token, err := jwt.ParseWithClaims(jwtStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
      return []byte(sign_key), nil //返回签名密钥
   }, options...)
   if err != nil {
      return nil, err
   }

    // 校验 Claims 对象是否有效，基于 exp（过期时间），nbf（不早于），iat（签发时间）等进行判断（如果有这些声明的话）。
   if !token.Valid {
      return nil, errors.New("invalid token")
   }

   // 解析claim
   claims, ok := token.Claims.(*MyCustomClaims)
   if !ok {
      return nil, errors.New("invalid claim type")
   }

   return claims, nil

}

func main() {
   /*
   jwtKey := make([]byte, 32) // 生成32字节（256位）的密钥
    if _, err := rand.Read(jwtKey); err != nil {
       panic(err)
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
       "iss": "程序员陈明勇",
       "sub": "chenmingyong.cn",
       "aud": "Programmer",
    })
    jwtStr, err := token.SignedString(jwtKey)
	*/

	jwtStr, err := GenerateJwt(sign_key)
    if err != nil {
       panic(err)
    }

    // 解析 jwt
    claims, err := ParseJwtWithClaims(jwtKey, jwtStr)
    if err != nil {
       panic(err)
    }
    fmt.Println(claims)
}

```

#### golang-jwt: v4

```go

package jwt

//gin框架中使用的jwt工具
import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	// 可根据需要自行添加字段
	Username             string `json:"username"`
	UserID               int64  `json:"userid"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

const TokenExpireDuration = time.Hour * 24

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte("夏天夏天悄悄过去")

// GenToken 生成JWT
func GenToken(userid int64, username string) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		username, // 自定义字段
		userid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "my-project", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}


// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

```


