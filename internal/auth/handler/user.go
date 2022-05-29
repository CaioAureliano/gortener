package handler

import (
	"log"
	"net/http"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(u *service.UserService) *UserHandler {
	return &UserHandler{
		userService: u,
	}
}

func (h *UserHandler) Create(c echo.Context) error {
	var userRequest *model.UserCreateRequest
	if err := c.Bind(&userRequest); err != nil {
		log.Printf("[auth handler] error to bind body: %s", err.Error())
		return err
	}

	if err := h.userService.Create(userRequest); err != nil {
		log.Printf("error to create a new user: %s", err.Error())
		return err
	}

	return c.NoContent(http.StatusCreated)
}
