package router

import (
	"github.com/CaioAureliano/gortener/internal/auth/handler"
	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	userHandler *handler.UserHandler
	authHandler *handler.AuthHandler
}

func NewAuthRouter(user *handler.UserHandler, auth *handler.AuthHandler) *AuthRouter {
	return &AuthRouter{
		userHandler: user,
		authHandler: auth,
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
