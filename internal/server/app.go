package server

import (
	"fmt"
	"os"

	"github.com/CaioAureliano/gortener/internal/auth/config"
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

	authConfig = config.NewAuthConfig()
)

func (a *App) Run() error {

	authConfig.Resolve(a.e)

	return a.e.Start(port)
}
