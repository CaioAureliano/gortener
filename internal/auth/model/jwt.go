package model

import "github.com/golang-jwt/jwt"

type JwtRespose struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type JwtCustomClaims struct {
	Email          string `json:"email"`
	StandardClaims jwt.StandardClaims
}
