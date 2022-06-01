package handler

import (
	"net/http"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

var authService = service.NewAuthService()

func (h *AuthHandler) Authenticate(c echo.Context) error {
	var req *model.AuthRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	res, err := authService.Login(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
