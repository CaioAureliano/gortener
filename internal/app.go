package internal

import (
	"github.com/CaioAureliano/gortener/internal/shortener/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(e echo.Echo) {
	e.Use(middleware.Logger())

	routes.Router(e)
}
