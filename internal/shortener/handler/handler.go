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
