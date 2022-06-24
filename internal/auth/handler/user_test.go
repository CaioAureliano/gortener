package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	userBody = `{
		"name": "Test",
		"email": "test@mail.com",
		"password": "pass1234"
	}`
)

func TestCreate(t *testing.T) {

	t.Run("valid request body should be return status code created", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		h := NewUserHandler(userServiceMock{})

		if assert.NoError(t, h.Create(ctx)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("empty request body should return error status code bad request", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		h := NewUserHandler(userServiceMock{})

		assert.Errorf(t, h.Create(ctx), "Required request body")
	})
}
