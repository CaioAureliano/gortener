package handler

import (
	"net/http"

	"github.com/CaioAureliano/gortener/internal/shortener/dto"
	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type mockService struct {
	fnCreate   func(req *dto.UrlRequest) (*model.Shortener, error)
	fnGet      func(slug string) (*model.Shortener, error)
	fnGetUrl   func(slug string) (string, error)
	fnAddClick func(click model.Click, slug string) (*model.Shortener, error)
	fnStats    func(slug string) (*model.Stats, error)
}

func (m mockService) Create(req *dto.UrlRequest) (*model.Shortener, error) {
	if m.fnCreate != nil {
		return m.fnCreate(req)
	}
	return nil, nil
}

func (m mockService) Get(slug string) (*model.Shortener, error) {
	if m.fnGet != nil {
		return m.fnGet(slug)
	}
	return nil, nil
}

func (m mockService) GetUrl(slug string) (string, error) {
	if m.fnGetUrl != nil {
		return m.fnGetUrl(slug)
	}
	return "", nil
}

func (m mockService) AddClick(click model.Click, slug string) (*model.Shortener, error) {
	if m.fnAddClick != nil {
		return m.fnAddClick(click, slug)
	}
	return nil, nil
}

func (m mockService) Stats(slug string) (*model.Stats, error) {
	if m.fnStats != nil {
		return m.fnStats(slug)
	}
	return nil, nil
}

type CustomValidatorMock struct {
	validator *validator.Validate
}

func (c *CustomValidatorMock) Validate(i interface{}) error {
	if err := c.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
