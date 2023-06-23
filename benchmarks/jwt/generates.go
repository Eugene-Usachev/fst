package jwt

import (
	"github.com/kataras/jwt"
	"time"
)

var (
	refreshKey = []byte(`key1`)
	accessKey  = []byte(`key1`)
)

const (
	AccessTokenLive  = time.Minute * 15
	RefreshTokenLive = time.Hour * 24 * 31
)

func NewRefreshToken(passwordHash string) ([]byte, error) {
	token, err := jwt.Sign(jwt.HS256, refreshKey, RefreshClaims{V: passwordHash})
	return token, err
}

func NewAccessToken(id uint) ([]byte, error) {
	token, err := jwt.Sign(jwt.HS256, accessKey, AccessClaims{V: id})
	return token, err
}
