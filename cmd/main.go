package main

import (
	"github.com/CaioAureliano/gortener/internal/server"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	app := server.NewApp(e)

	e.Logger.Fatal(app.Run())
}
