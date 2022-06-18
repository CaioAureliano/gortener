package router

import (
	"github.com/CaioAureliano/gortener/internal/auth/handler"
	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
}

func NewAuthRouter() *AuthRouter {
	return &AuthRouter{}
}

var (
	authHandler = handler.NewAuthHandler()
	userHandler = handler.NewUserHandler()
)

const (
	AUTH_ENDPOINT  = "/login"
	USERS_ENDPOINT = "/users"
)

func (r *AuthRouter) Router(e *echo.Echo) {
	e.POST(AUTH_ENDPOINT, authHandler.Authenticate)
	e.POST(USERS_ENDPOINT, userHandler.Create)
}
