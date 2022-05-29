package router

import (
	"github.com/CaioAureliano/gortener/internal/auth/handler"
	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
}

func NewAuthRouter(a *handler.AuthHandler, u *handler.UserHandler) *AuthRouter {
	return &AuthRouter{
		authHandler: a,
		userHandler: u,
	}
}

const (
	AUTH_ENDPOINT  = "/login"
	USERS_ENDPOINT = "/users"
)

func (r *AuthRouter) Router(e *echo.Echo) {
	e.POST(AUTH_ENDPOINT, r.authHandler.Authenticate)

	e.POST(USERS_ENDPOINT, r.userHandler.Create)
}
