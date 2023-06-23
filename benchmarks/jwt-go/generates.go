package jwt_go

import (
	"github.com/dgrijalva/jwt-go"
)

func NewRefreshToken(key []byte, passwordHash string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshClaims{
		jwt.StandardClaims{},
		passwordHash,
	})
	accessToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return accessToken, err
}

func NewAccessToken(id uint, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AccessClaims{
		jwt.StandardClaims{},
		id,
	})
	accessToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return accessToken, err
}
