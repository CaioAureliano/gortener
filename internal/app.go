package application

import (
	"os"

	"github.com/CaioAureliano/gortener/internal/shortener/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	port = ":" + os.Getenv("PORT")
)

func Run(e *echo.Echo) {
	e.Use(middleware.Logger())

	routes.Router(e)

	e.Logger.Fatal(e.Start(port))
}
