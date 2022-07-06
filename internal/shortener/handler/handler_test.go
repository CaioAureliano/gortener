package handler

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortUrl(t *testing.T) {
	body := `{"url": "example.com"}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	e := echo.New()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	shortenerService = func() service.Shortener {
		return mockService{
			fnCreate: func(url string) (*model.Shortener, error) {
				return &model.Shortener{Url: url}, nil
			},
		}
	}

	err := CreateShortUrl(c)

	if assert.NoError(t, err) {
		log.Printf("body response: %v", rec.Body.String())
		log.Printf("headers response: %v", rec.Header())

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header()[echo.HeaderContentType][0])
	}
}
