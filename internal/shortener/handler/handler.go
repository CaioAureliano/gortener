package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CaioAureliano/gortener/internal/shortener/dto"
	"github.com/CaioAureliano/gortener/internal/shortener/service"
	"github.com/labstack/echo/v4"
)

var (
	shortenerService = service.New
)

func CreateShortUrl(c echo.Context) error {
	var req *dto.UrlRequest
	if err := c.Bind(&req); err != nil || req == nil {
		log.Printf("error to bind body request: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "bad request: not found url")
	}

	res, err := shortenerService().Create(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().Header().Set(echo.HeaderLocation, "/"+res.Slug)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(res)
}
