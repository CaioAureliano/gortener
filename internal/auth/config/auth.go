package config

import (
	"github.com/CaioAureliano/gortener/internal/auth/handler"
	"github.com/CaioAureliano/gortener/internal/auth/repository"
	"github.com/CaioAureliano/gortener/internal/auth/router"
	"github.com/CaioAureliano/gortener/internal/auth/service"
	"github.com/CaioAureliano/gortener/internal/auth/validator"
	"github.com/labstack/echo/v4"
)

type AuthConfig struct {
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{}
}

func (c *AuthConfig) Resolve(e *echo.Echo) {

	userRepository := repository.NewUserRepository()

	authService := service.NewAuthService()
	userService := service.NewUserService(userRepository)

	authHandler := handler.NewAuthHandler(*authService)
	userHandler := handler.NewUserHandler(userService)

	router := router.NewAuthRouter(authHandler, userHandler)

	e.Validator = validator.NewAuthValidator(userService)

	router.Router(e)
}
