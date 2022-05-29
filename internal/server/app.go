package server

import (
	"fmt"
	"os"

	"github.com/CaioAureliano/gortener/internal/auth/router"
	"github.com/CaioAureliano/gortener/internal/auth/validator"
	"github.com/labstack/echo/v4"
)

type App struct {
	e *echo.Echo
}

func NewApp(e *echo.Echo) *App {
	return &App{e}
}

var (
	port = fmt.Sprintf(":%s", os.Getenv("PORT"))

	authRouter = router.NewAuthRouter()
)

func (a *App) Run() error {

	a.e.Validator = validator.NewAuthValidator()

	authRouter.Router(a.e)

	return a.e.Start(port)
}
