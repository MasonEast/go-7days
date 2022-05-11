package define

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	Id       int
	UserName     string
	jwt.StandardClaims
}

var JwtKey = "cloud-disk-key"
var EmailPassword = os.Getenv("password163")

// CodeLength 验证码长度
var CodeLength = 6

// CodeExpire 验证码过期时间（s）
var CodeExpire = 300

// cosSecretKey 腾讯云对象存储
var CosSecretKey = os.Getenv("cosSecretKey")
var CosSecretID = os.Getenv("cosSecretId")
var CosBucket = "https://getcharzp-1256268070.cos.ap-chengdu.myqcloud.com"

// PageSize 分页的默认参数
var PageSize = 20

var Datetime = "2006-01-02 15:04:05"

var TokenExpire = 3600
var RefreshTokenExpire = 7200
