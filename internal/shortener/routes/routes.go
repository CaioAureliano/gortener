package routes

import (
	"github.com/CaioAureliano/gortener/internal/shortener/handler"
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo) {
	e.POST("/", handler.CreateShortUrl)
	e.GET("/:slug", handler.Redirect)
	e.GET("/:slug/stats", handler.Stats)
}
