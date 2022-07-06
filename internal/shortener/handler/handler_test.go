package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CaioAureliano/gortener/internal/shortener/dto"
	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortUrl(t *testing.T) {
	tests := []struct {
		name    string
		request string

		serviceMock service.Shortener

		wantStatusCode   int
		wantSlugLocation string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name:    "should be return created status with valid request body",
			request: `{"url": "example.com"}`,

			serviceMock: mockService{
				fnCreate: func(req *dto.UrlRequest) (*model.Shortener, error) {
					req.AppendProtocolIfNotExists()
					return &model.Shortener{
						Url:  req.Url,
						Slug: "sl8g3",
					}, nil
				},
			},

			wantStatusCode:   http.StatusCreated,
			wantSlugLocation: "sl8g3",
			wantErr:          assert.NoError,
		},
		{
			name:    "should be return bad request status with invalid request body",
			request: `{}`,

			serviceMock: mockService{},

			wantStatusCode:   http.StatusBadRequest,
			wantSlugLocation: "",
			wantErr:          assert.Error,
		},
		{
			name:    "should be return bad request status with invalid request body",
			request: `{"url": ""}`,

			serviceMock: mockService{},

			wantStatusCode:   http.StatusBadRequest,
			wantSlugLocation: "",
			wantErr:          assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidatorMock{validator: validator.New()}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			shortenerService = func() service.Shortener {
				return tt.serviceMock
			}

			err := CreateShortUrl(ctx)

			tt.wantErr(t, err)
			if tt.wantStatusCode == http.StatusCreated {
				assert.Equal(t, tt.wantStatusCode, rec.Code)
			} else {
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.wantStatusCode, httpErr.Code)
			}

			if tt.wantSlugLocation != "" {
				assert.Equal(t, "/"+tt.wantSlugLocation, rec.Header()[echo.HeaderLocation][0])
			}
		})
	}
}
