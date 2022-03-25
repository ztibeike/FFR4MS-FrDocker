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
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}
