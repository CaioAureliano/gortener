package handler

import (
	"github.com/CaioAureliano/gortener/internal/shortener/service"
	"github.com/labstack/echo/v4"
)

var (
	shortenerService = service.New
)

func CreateShortUrl(c echo.Context) error {
	return nil
}
