package handler

import (
	"net/http"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type ShortenerHandler struct {
}

func NewShortenerHandler() *ShortenerHandler {
	return &ShortenerHandler{}
}

var shortenerService = service.NewShortenerService()

func (sh *ShortenerHandler) Create(c echo.Context) error {
	var req *model.ShortenerCreateRequest
	if err := c.Bind(&req); err != nil || req.Url == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	user := c.Get("user").(*jwt.Token)

	slug, err := shortenerService.Create(req, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"slug": slug})
}

func (sh *ShortenerHandler) RedirectBySlug(c echo.Context) error {
	url, err := shortenerService.GetUrlBySlug(c.Param("slug"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "not found"})
	}

	return c.Redirect(http.StatusMovedPermanently, url)
}
