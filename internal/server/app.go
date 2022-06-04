package server

import (
	"fmt"
	"os"

	authRouter "github.com/CaioAureliano/gortener/internal/auth/router"
	shortenerRouter "github.com/CaioAureliano/gortener/internal/shortener/router"
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

	ar = authRouter.NewAuthRouter()
	sr = shortenerRouter.NewShortenerRouter()
)

func (a *App) Run() error {
	ar.Router(a.e)
	sr.Router(a.e)

	return a.e.Start(port)
}
