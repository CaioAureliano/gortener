package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ShortenerHandler struct {
}

func NewShortenerHandler() *ShortenerHandler {
	return &ShortenerHandler{}
}

func (sh *ShortenerHandler) Create(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}
