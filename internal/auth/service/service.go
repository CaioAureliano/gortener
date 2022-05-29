package service

import (
	"log"
	"time"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/golang-jwt/jwt"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

type JwtRespose struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type jwtCustomClaims struct {
	Email          string `json:"email"`
	StandardClaims jwt.StandardClaims
}

var TOKEN_SECRET_KEY = []byte("secret")

func (a *AuthService) Login(userRequest *model.AuthRequest) (*JwtRespose, error) {
	claims := a.createCustomClaimsByEmail(userRequest.Email)

	token, err := a.generateTokenByClaims(claims.StandardClaims)
	if err != nil {
		log.Printf("error to generate token: %s", err.Error())
		return nil, err
	}

	return a.makeJwtResponse(token, claims.StandardClaims.ExpiresAt), nil
}

func (a *AuthService) createCustomClaimsByEmail(email string) *jwtCustomClaims {
	return &jwtCustomClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: a.expiresAt().UnixMilli(),
		},
	}
}

func (a *AuthService) expiresAt() time.Time {
	return time.Now().Add(time.Hour * 12)
}

func (a *AuthService) generateTokenByClaims(claims jwt.StandardClaims) (t string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if t, err = token.SignedString(TOKEN_SECRET_KEY); err != nil {
		return
	}

	return
}

func (a *AuthService) makeJwtResponse(token string, expiresAt int64) *JwtRespose {
	return &JwtRespose{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}
