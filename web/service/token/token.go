package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("USTBFRDOCKER")

type Claims struct {
	UserId   string `json:"userId"`
	UserRole string `json:"userRole"`
	jwt.StandardClaims
}

func GenerateToken(userId string, userRole string) (string, error) {
	expireAt := time.Now().Add(1 * time.Hour).Unix()
	claims := &Claims{
		userId,
		userRole,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    "frdocker",
			Subject:   "user token",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)
	return tokenStr, err
}

func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	return claims, err
}

func RefreshToken(claims *Claims) (string, error) {
	if withinLimit(claims.ExpiresAt, 600) {
		return GenerateToken(claims.UserId, claims.UserRole)
	}
	return "", errors.New("token is expired")
}

func withinLimit(s int64, l int64) bool {
	e := time.Now().Unix()
	return e-s < l
}
