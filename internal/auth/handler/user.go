package handler

import (
	"log"
	"net/http"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/internal/auth/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (h *UserHandler) Create(c echo.Context) error {
	var userRequest *model.UserCreateRequest
	if err := c.Bind(&userRequest); err != nil {
		log.Printf("error to bind body: %s", err.Error())
		return err
	}

	if userRequest == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Required request body")
	}

	if err := h.userService.Create(userRequest); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}
