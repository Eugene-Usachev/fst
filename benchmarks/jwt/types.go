package jwt

type AccessClaims struct {
	V uint `json:"v"`
}

type RefreshClaims struct {
	V string `json:"v"`
}
