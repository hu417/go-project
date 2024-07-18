package define

// jwt
import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

// 定义用户 jwt token相关参数
type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.StandardClaims
}

// 定义token有效期
var (
	TokenExpires        = 3600
	RefreshTokenExpires = 24 * 60 * 60
)

// 加盐
var Shar = "cloud-disk"

// 邮箱验证码
var Pwd = "MNTITILCKMRCLPSV"

// 定义验证码长度
var CodeLength = 6

// 验证码过期时间
var CodeExoire = 300

// cos存储相关
var COS_SECRETID = os.Getenv("COS_SECRETID")
var COS_SECRETKEY = os.Getenv("COS_SECRETKEY")
var BucketUrl = "https://cloud-1304907914.cos.ap-guangzhou.myqcloud.com/"

// 分页默认参数
var (
	Page     = 1
	PageSize = 20
)

// 定义时间格式
var Datetime = "2006-01-02 15:04:05"
