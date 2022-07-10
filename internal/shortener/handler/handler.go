package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/CaioAureliano/gortener/internal/shortener/dto"
	"github.com/CaioAureliano/gortener/internal/shortener/service"
	"github.com/labstack/echo/v4"
)

var (
	shortenerService = service.New

	ErrBadRequestUrl = errors.New("bad request: invalid url")
)

func CreateShortUrl(c echo.Context) error {
	req := new(dto.UrlRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("error to bind body request: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrBadRequestUrl)
	}

	if err := c.Validate(req); err != nil {
		log.Printf("error to validate request URL (%s)", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrBadRequestUrl)
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

func Redirect(c echo.Context) error {
	slug := c.Param("slug")

	url, err := shortenerService().GetUrl(slug)
	if err != nil {
		log.Printf("error to find url by slug: %s [%s]", slug, err.Error())
		return echo.NewHTTPError(http.StatusNotFound, "not found: "+err.Error())
	}

	clicked := new(dto.ClickedRequest)
	defer shortenerService().AddClick(clicked.Set(c.Request()), slug)

	return c.Redirect(http.StatusMovedPermanently, url)
}

func Stats(c echo.Context) error {
	slug := c.Param("slug")

	stats, err := shortenerService().Stats(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	short, err := shortenerService().Get(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"stats":      stats,
		"url":        short.Url,
		"created_at": short.CreatedAt,
	})
}
