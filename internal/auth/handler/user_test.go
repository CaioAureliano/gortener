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
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	h := NewUserHandler(userServiceMock{})

	if assert.NoError(t, h.Create(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
