package goJose

import (
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

var (
	refreshKey = []byte(`key1`)
	accessKey  = []byte(`key1`)
)

const (
	AccessTokenLive  = time.Minute * 15
	RefreshTokenLive = time.Hour * 24 * 31
)

type AccessClaims struct {
	V uint `json:"v"`
}

type RefreshClaims struct {
	V string `json:"v"`
}

func NewRefreshToken(passwordHash string) ([]byte, error) {
	claims := RefreshClaims{
		V: passwordHash,
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: refreshKey}, nil)
	if err != nil {
		return nil, err
	}

	builder := jwt.Signed(signer).Claims(claims).Claims(jwt.Claims{})

	token, err := builder.CompactSerialize()
	if err != nil {
		return nil, err
	}

	return []byte(token), nil
}

func NewAccessToken(id uint) ([]byte, error) {
	claims := AccessClaims{
		V: id,
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: accessKey}, nil)
	if err != nil {
		return nil, err
	}

	builder := jwt.Signed(signer).Claims(claims).Claims(jwt.Claims{})

	token, err := builder.CompactSerialize()
	if err != nil {
		return nil, err
	}

	return []byte(token), nil
}

func ParseRefreshToken(token []byte) (string, error) {
	refreshClaims := RefreshClaims{}

	// Parse the token
	parsedToken, err := jwt.ParseSigned(string(token))
	if err != nil {
		return "", err
	}

	// Verify the signature and extract the claims
	err = parsedToken.Claims(refreshKey, &refreshClaims)
	if err != nil {
		return "", err
	}

	return refreshClaims.V, nil
}

func ParseAccessToken(token []byte) (uint, error) {
	accessClaims := AccessClaims{}

	parsedToken, err := jwt.ParseSigned(string(token))
	if err != nil {
		return 0, err
	}

	err = parsedToken.Claims(accessKey, &accessClaims)
	if err != nil {
		return 0, err
	}

	return accessClaims.V, nil
}
