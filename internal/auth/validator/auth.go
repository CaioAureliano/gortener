package validator

import (
	"log"
	"net/http"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/service"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type AuthValidator struct {
	validator   *validator.Validate
	userService *service.UserService
}

func NewAuthValidator(u *service.UserService) *AuthValidator {
	return &AuthValidator{
		validator:   validator.New(),
		userService: u,
	}
}

func (av *AuthValidator) Validate(i interface{}) error {
	exists, err := av.userService.Exists(i.(*model.AuthRequest))
	if !exists {
		log.Printf("user dont exists: %s", err.Error())
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return nil
}
