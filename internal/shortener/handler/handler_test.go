package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CaioAureliano/gortener/internal/shortener/dto"
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

	slugMock := "sl8g3"
	shortenerService = func() service.Shortener {
		return mockService{
			fnCreate: func(req *dto.UrlRequest) (*model.Shortener, error) {
				req.AppendProtocolIfNotExists()
				return &model.Shortener{
					Url:  req.Url,
					Slug: slugMock,
				}, nil
			},
		}
	}

	err := CreateShortUrl(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header()[echo.HeaderContentType][0])
		assert.Equal(t, "/"+slugMock, rec.Header()[echo.HeaderLocation][0])
	}
}
