package application

import (
	"os"

	request_validator "github.com/CaioAureliano/gortener/internal/shortener/handler/validator"
	"github.com/CaioAureliano/gortener/internal/shortener/routes"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	port = ":" + os.Getenv("PORT")
)

func Run(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Validator = &request_validator.CustomRequestValidator{Validator: validator.New()}

	routes.Router(e)

	e.Logger.Fatal(e.Start(port))
}
