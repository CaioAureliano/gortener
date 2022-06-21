package handler

import (
	"net/http"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Authenticate(c echo.Context) error {
	var req *model.AuthRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	res, err := h.authService.Login(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
