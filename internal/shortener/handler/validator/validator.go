package validator

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomRequestValidator struct {
	Validator *validator.Validate
}

func (r *CustomRequestValidator) Validate(i interface{}) error {
	if err := r.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
