package service

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

var (
	tokenSecretKey = os.Getenv("SECRET")

	userService = NewUserService()
)

func (a *AuthService) Login(req *model.AuthRequest) (*model.JwtRespose, error) {
	if err := validateUser(req); err != nil {
		return nil, echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	claims := createClaims(req.Email)

	token, err := generateToken(claims.StandardClaims)
	if err != nil {
		log.Printf("error to generate token: %s", err.Error())
		return nil, err
	}

	return buildJwtResponse(token, claims.StandardClaims.ExpiresAt), nil
}

func createClaims(email string) *model.JwtCustomClaims {
	return &model.JwtCustomClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt().UnixMilli(),
		},
	}
}

func expiresAt() time.Time {
	return time.Now().Add(time.Hour * 12)
}

func generateToken(claims jwt.StandardClaims) (t string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(tokenSecretKey))
	return
}

func buildJwtResponse(token string, expiresAt int64) *model.JwtRespose {
	return &model.JwtRespose{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

func validateUser(req *model.AuthRequest) error {
	if exists, err := userService.Exists(req.Email); !exists || err != nil {
		return errors.New("user not exists")
	}

	userFound, err := userService.GetByField(req.Email, "email")
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(req.Password)); err != nil {
		return errors.New("invalid email/password")
	}

	return nil
}
