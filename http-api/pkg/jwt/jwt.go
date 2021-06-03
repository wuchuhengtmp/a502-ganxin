package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"http-api/pkg/config"
	"http-api/pkg/logger"
	"strconv"
	"time"
)

type MyCustomClaims struct {
	IsDevice bool   `json:"isDevice"`
	Mac      string `json:"mac"`
	Uid      int64  `json:"uid"`
	jwt.StandardClaims
}

// 生成token
func GenerateTokenByUID(uid int64, isDevice bool, mac string) (tokenStr string, err error) {
	privateKey := []byte(config.GetString("jwt.secret"))
	claims := MyCustomClaims{
		isDevice,
		mac,
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + config.GetInt64("jwt.expired"),
			Id:	strconv.FormatInt(uid, 10),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(privateKey)
}

/**
 * jwt token 过期时间
 */
func GetExpiredAt() int64 {
	return time.Now().Unix() + config.GetInt64("jwt.expired")
}

// 解析token
func ParseByTokenStr(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetString("jwt.secret")), nil
	})
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return &MyCustomClaims{}, fmt.Errorf("The token was invalid")
}
