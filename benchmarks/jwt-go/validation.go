package jwt_go

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

func ParseAccessToken(token string, key []byte) (uint, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid signing method")
		}

		return key, nil
	})
	if err != nil {
		return 0, err
	}

	tokenClaims, ok := parsedToken.Claims.(*AccessClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return tokenClaims.V, nil
}

func ParseRefreshToken(token string, key []byte) (string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid signing method")
		}

		return key, nil
	})
	if err != nil {
		return "", err
	}

	tokenClaims, ok := parsedToken.Claims.(*RefreshClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return tokenClaims.V, nil
}
