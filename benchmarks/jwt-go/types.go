package jwt_go

import "github.com/dgrijalva/jwt-go"

type AccessClaims struct {
	jwt.StandardClaims
	V uint `json:"v"`
}

type RefreshClaims struct {
	jwt.StandardClaims
	V string `json:"v"`
}
