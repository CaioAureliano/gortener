package router

import (
	authMiddleware "github.com/CaioAureliano/gortener/internal/auth/middleware"
	"github.com/CaioAureliano/gortener/internal/shortener/handler"
	"github.com/labstack/echo/v4"
)

type ShortenerRouter struct {
}

func NewShortenerRouter() *ShortenerRouter {
	return &ShortenerRouter{}
}

const SHORTENER_ENDPOINT = "/shortener"

var shortenerHandler = handler.NewShortenerHandler()

func (sr *ShortenerRouter) Router(e *echo.Echo) {
	s := e.Group(SHORTENER_ENDPOINT)

	s.Use(authMiddleware.ConfigJwt())

	s.POST("", shortenerHandler.Create)
}
