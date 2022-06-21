package server

import (
	"os"

	"github.com/CaioAureliano/gortener/internal/auth/dao"
	"github.com/CaioAureliano/gortener/internal/auth/handler"
	authRouter "github.com/CaioAureliano/gortener/internal/auth/router"
	"github.com/CaioAureliano/gortener/internal/auth/service"
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
	port = ":" + os.Getenv("PORT")
	sr   = shortenerRouter.NewShortenerRouter()
)

func (a *App) Run() error {

	userRepository := dao.NewUserDao()
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userService)
	authHandler := handler.NewAuthHandler(authService)

	authRoutes := authRouter.NewAuthRouter(userHandler, authHandler)
	authRoutes.Router(a.e)

	sr.Router(a.e)

	return a.e.Start(port)
}
