package validator

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type UserValidator struct {
	validator *validator.Validate
}

func (u *UserValidator) Validate(i interface{}) error {
	if err := u.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return nil
}
