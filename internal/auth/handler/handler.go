package handler

import (
	"log"
	"net/http"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Authenticate(c echo.Context) error {
	var loginRequest *model.AuthRequest
	if err := c.Bind(&loginRequest); err != nil {
		return err
	}

	if err := c.Validate(loginRequest); err != nil {
		log.Printf("error to login: %s", err.Error())
		return err
	}

	res, err := h.authService.Login(loginRequest)
	if err != nil {
		log.Printf("error to authenticate: %s", err.Error())
		return err
	}

	return c.JSON(http.StatusOK, res)
}
