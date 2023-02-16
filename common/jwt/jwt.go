package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// 创建jwt token的payload
func createClaims(iat, seconds, userId int64) jwt.MapClaims {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	return claims
}

// 生成jwt token
func GetJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := createClaims(iat, seconds, userId)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func parseJwtToken(secretKey string, tokenString string) (jwt.MapClaims, error) {
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
			return nil, errors.New("parse jwt claims failed")
		}
	}
	return nil, errors.New("invalid jwt token")
}

func ParseUserIdFromJwtToken(secretKey string, tokenString string) (int64, error) {
	claims, err := parseJwtToken(secretKey, tokenString)
	if err != nil {
		return 0, err
	}
	userId := claims["userId"]
	if userId == nil {
		return 0, nil
	}
	return userId.(int64), nil
}
