package handler

import (
	"log"
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
	var loginRequest *model.AuthRequest
	if err := c.Bind(&loginRequest); err != nil {
		return err
	}

	if err := c.Validate(loginRequest); err != nil {
		log.Printf("error to login: %s", err.Error())
		return err
	}

	res, err := authService.Login(loginRequest)
	if err != nil {
		log.Printf("error to authenticate: %s", err.Error())
		return err
	}

	return c.JSON(http.StatusOK, res)
}
