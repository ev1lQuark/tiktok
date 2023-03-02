package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

// 生成jwt
func Create(secretKey string, seconds, userId int64) (string, error) {
	iat := time.Now().Unix()
	claims := jwt.MapClaims{
		"userId": userId,
		"iat":    iat,
		"exp":    iat + seconds,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 解析jwt，取出claims
func parseClaims(secretKey string, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			return claims, nil
		} else {
			return nil, errors.New("invalid jwt token claims")
		}
	} else {
		return nil, errors.New("invalid jwt token")
	}
}

// 检查jwt合法性
func Verify(secretKey string, tokenString string) bool {
	if disable, ok := os.LookupEnv("JWT_DISABLE"); ok && disable == "true" {
		logx.Info("jwt disabled")
		return true
	}
	claims, err := parseClaims(secretKey, tokenString)
	if err != nil {
		return false
	}
	return claims.VerifyExpiresAt(time.Now().Unix(), true)
}

// 检查 jwt 合法性并返回 userId
func GetUserId(secretKey string, tokenString string) (int64, error) {
	if disable, ok := os.LookupEnv("JWT_DISABLE"); ok && disable == "true" {
		logx.Info("jwt disabled")
		return 666666, nil
	}
	claims, err := parseClaims(secretKey, tokenString)
	if err != nil {
		return 0, err
	}
	userId := claims["userId"]
	if userId == nil {
		return 0, nil
	}
	return int64(userId.(float64)), nil
}
