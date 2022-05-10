package util

import (
	"cloud-disk/core/define"
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
)

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int, UserName string, second int) (string, error) {
	uc := define.UserClaim{
		Id:       id,
		UserName:     UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(second)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func SendEmail(s , code string) error{
	e := email.NewEmail()
	e.From = "mason <sh941118@163.com>"
	e.To = []string{s}

	e.Subject = "注册网盘邮箱验证"
	e.HTML = []byte("您的验证码为：<h1>"+code+"</h1>")
	err := e.Send("smtp.163.com:25", smtp.PlainAuth("", "sh941118@163.com", "VSCYFUSEUZXKMYMQ", "smtp.163.com"))
	if err != nil {
		return err
	}
	return nil
}

func RandCode() string {
	s := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < define.CodeLength; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}