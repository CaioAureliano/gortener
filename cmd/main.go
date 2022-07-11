package main

import (
	application "github.com/CaioAureliano/gortener/internal"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	application.Run(e)
}
