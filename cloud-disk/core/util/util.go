package util

import (
	"cloud-disk/core/define"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
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
