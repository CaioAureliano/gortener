package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

		wantStatusCode int
		wantErr        assert.ErrorAssertionFunc
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

			wantStatusCode: http.StatusCreated,
			wantErr:        assert.NoError,
		},
		{
			name:    "should be return bad request status with invalid request body",
			request: `{}`,

			serviceMock: mockService{},

			wantStatusCode: http.StatusBadRequest,
			wantErr:        assert.Error,
		},
		{
			name:    "should be return bad request status with invalid request body",
			request: `{"url": ""}`,

			serviceMock: mockService{},

			wantStatusCode: http.StatusBadRequest,
			wantErr:        assert.Error,
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
		})
	}
}

func TestRedirect(t *testing.T) {
	tests := []struct {
		name       string
		gotRequest *http.Request

		wantUrl    string
		wantStatus int
		wantErr    assert.ErrorAssertionFunc

		mockService service.Shortener
	}{
		{
			name: "should be redirect to url with valid slug with moved permanently http status code",

			gotRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
				req.Header.Set("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.53 Safari/537.36`)
				req.Header.Set("Accept-Language", `en-US,en;q=0.5`)
				return req
			}(),

			wantUrl:    "http://example.com",
			wantStatus: http.StatusMovedPermanently,
			wantErr:    assert.NoError,

			mockService: mockService{
				fnGetUrl: func(slug string) (string, error) {
					return "http://example.com", nil
				},
			},
		},
		{
			name: "should be return not found http status code with invalid slug",

			gotRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
				req.Header.Set("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0`)
				req.Header.Set("Accept-Language", `en-US,en;q=0.5`)
				return req
			}(),

			wantUrl:    "",
			wantStatus: http.StatusNotFound,
			wantErr:    assert.Error,

			mockService: mockService{
				fnGetUrl: func(slug string) (string, error) {
					return "", assert.AnError
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			rec := httptest.NewRecorder()
			ctx := e.NewContext(tt.gotRequest, rec)

			shortenerService = func() service.Shortener {
				return tt.mockService
			}

			err := Redirect(ctx)

			tt.wantErr(t, err)

			if tt.wantStatus == http.StatusMovedPermanently {
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.Equal(t, tt.wantUrl, rec.Header()[echo.HeaderLocation][0])
			} else {
				httpErr, ok := err.(*echo.HTTPError)

				assert.True(t, ok)
				assert.Equal(t, tt.wantStatus, httpErr.Code)
			}
		})
	}
}

func TestStats(t *testing.T) {

	tests := []struct {
		name string

		wantStatus int
		wantErr    assert.ErrorAssertionFunc

		gotStatsResponse *model.Stats
		gotStatsErr      error

		gotShortResponse *model.Shortener
		gotShortErr      error
	}{
		{
			name: "should be return OK status with valid slug",

			wantStatus: http.StatusOK,
			wantErr:    assert.NoError,

			gotStatsResponse: &model.Stats{
				Clicks:   10,
				Browsers: map[string]int{"chrome": 5},
			},
			gotShortResponse: &model.Shortener{
				Url:       "http://example.com",
				CreatedAt: time.Now(),
			},
		},
		{
			name: "should be return not found status with invalid slug",

			wantStatus: http.StatusNotFound,
			wantErr:    assert.Error,

			gotShortErr: service.ErrShortenerNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/sl0g3/stats", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			shortenerService = func() service.Shortener {
				return mockService{
					fnStats: func(slug string) (*model.Stats, error) {
						return tt.gotStatsResponse, tt.gotStatsErr
					},
					fnGet: func(slug string) (*model.Shortener, error) {
						return tt.gotShortResponse, tt.gotShortErr
					},
				}
			}

			err := Stats(ctx)

			tt.wantErr(t, err)

			if tt.wantStatus != http.StatusOK {
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.wantStatus, httpErr.Code)
			} else {
				body := map[string]interface{}{
					"stats":      tt.gotStatsResponse,
					"url":        tt.gotShortResponse.Url,
					"created_at": tt.gotShortResponse.CreatedAt,
				}

				expected, err := json.Marshal(body)
				res := strings.TrimSpace(rec.Body.String())

				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.Equal(t, string(expected), res)
			}
		})
	}
}
