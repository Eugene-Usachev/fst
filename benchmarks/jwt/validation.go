package jwt

import (
	"github.com/kataras/jwt"
)

func ParseAccessToken(token []byte) (uint, error) {
	verifiedToken, err := jwt.Verify(jwt.HS256, accessKey, token)
	if err != nil {
		return 0, err
	}
	var tokenClaims AccessClaims
	err = verifiedToken.Claims(&tokenClaims)
	if err != nil {
		return 0, err
	}
	return tokenClaims.V, nil
}

func ParseRefreshToken(token []byte) (string, error) {
	verifiedToken, err := jwt.Verify(jwt.HS256, refreshKey, token)
	if err != nil {
		return "", err
	}
	var tokenClaims RefreshClaims
	err = verifiedToken.Claims(&tokenClaims)
	if err != nil {
		return "", err
	}
	return tokenClaims.V, nil
}
