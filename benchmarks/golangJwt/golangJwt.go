package golangJwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

func NewAccessToken(id uint, key []byte) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"v": id,
	}).SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

func NewRefreshToken(pass string, key []byte) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"v": pass,
	}).SignedString(key)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return token, nil
}

var errInvalidToken = fmt.Errorf("invalid token")

func ParseAccessToken(tokenString string, key []byte) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidToken
		}

		return key, nil
	})
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return uint(claims["v"].(float64)), nil
	} else {
		return 0, errInvalidToken
	}
}

func ParseRefreshToken(tokenString string, key []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidToken
		}

		return key, nil
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["v"].(string), nil
	} else {
		return "", errInvalidToken
	}
}
