package handler

import (
	"log"
	"net/http"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

var userService = service.NewUserService()

func (h *UserHandler) Create(c echo.Context) error {
	var userRequest *model.UserCreateRequest
	if err := c.Bind(&userRequest); err != nil {
		log.Printf("error to bind body: %s", err.Error())
		return err
	}

	if err := userService.Create(userRequest); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}
