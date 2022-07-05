package service

import (
	"errors"
	"testing"
	"time"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/repository"
	"github.com/CaioAureliano/gortener/internal/shortener/repository/cache"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	mockRedisCache := make(map[string]interface{})
	cacheRepository = func() cache.Cache {
		return mockCache{
			fnSet: func(key, value string, duration time.Duration) error {
				mockRedisCache[key] = value
				return nil
			},
		}
	}

	t.Run("should be return model response with valid properties values", func(t *testing.T) {
		mockUrl := "www.google.com"
		expectedUrl := "http://" + mockUrl

		shortenerRepository = func() repository.Shortener {
			return mockRepository{
				fnCreate: func(shortener *model.Shortener) (*model.Shortener, error) {
					return shortener, nil
				},
			}
		}

		shortService := New()
		shortCreated, err := shortService.Create(mockUrl)

		if err != nil {
			t.Errorf("error to create a shortener URL: %s", err.Error())
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedUrl, shortCreated.Url)
		assert.Equal(t, slugLength, len([]byte(shortCreated.Slug)))
	})

	t.Run("URL create", func(t *testing.T) {
		tests := []struct {
			name     string
			gotUrl   string
			wantUrl  string
			wantErr  error
			repoMock repository.Shortener
		}{
			{
				name:    "should be return a shortener created with valid URL",
				gotUrl:  "google.com",
				wantUrl: "http://google.com",
				wantErr: nil,
				repoMock: mockRepository{
					fnCreate: func(shortener *model.Shortener) (*model.Shortener, error) {
						return shortener, nil
					},
				},
			},
			{
				name:    "should be return ErrInvalidURL with invalid URL",
				gotUrl:  "url",
				wantUrl: "",
				wantErr: ErrInvalidURL,
				repoMock: mockRepository{
					fnCreate: func(shortener *model.Shortener) (*model.Shortener, error) {
						return shortener, nil
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				shortenerRepository = func() repository.Shortener {
					return tt.repoMock
				}

				shortenerService := New()
				created, err := shortenerService.Create(tt.gotUrl)

				if created != nil {
					assert.Equal(t, tt.wantUrl, created.Url)
				}
				assert.Equal(t, tt.wantErr, err)
			})
		}
	})
}

func TestGet(t *testing.T) {
	tests := []struct {
		name           string
		gotSlug        string
		wantErr        error
		repositoryMock repository.Shortener
	}{
		{
			name:    "should be return shortener response model with valid slug request",
			gotSlug: "SL0G3",
			wantErr: nil,
			repositoryMock: mockRepository{
				fnGet: func(slug string) (*model.Shortener, error) {
					return &model.Shortener{Slug: slug}, nil
				},
			},
		},
		{
			name:           "should be return ErrInvalidSlug with invalid length slug",
			gotSlug:        "5LU6",
			wantErr:        ErrInvalidSlug,
			repositoryMock: mockRepository{},
		},
		{
			name:           "should be return ErrInvalidSlug with invalid length slug",
			gotSlug:        "5LUG30",
			wantErr:        ErrInvalidSlug,
			repositoryMock: mockRepository{},
		},
		{
			name:           "should be return error with empty slug",
			gotSlug:        "",
			wantErr:        ErrInvalidSlug,
			repositoryMock: mockRepository{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortenerRepository = func() repository.Shortener {
				return tt.repositoryMock
			}

			shortenerService := New()
			res, err := shortenerService.Get(tt.gotSlug)

			assert.Equal(t, tt.wantErr, err)
			if res != nil {
				assert.Equal(t, tt.gotSlug, res.Slug)
			}
		})
	}
}

func TestGetUrl(t *testing.T) {
	redisClient, mock := redismock.NewClientMock()
	cacheRepository = func() cache.Cache {
		return NewMockCache(redisClient)
	}

	t.Run("should be return url from cache with valid slug", func(t *testing.T) {
		slugMock := "SL5G0"
		expectedUrl := "http://example.com"

		mock.ExpectGet(slugMock).SetVal(expectedUrl)

		shortenerService := New()
		url, err := shortenerService.GetUrl(slugMock)

		assert.NoError(t, err)
		assert.Equal(t, expectedUrl, url)
	})

	t.Run("should be return url from database with valid slug and cache url", func(t *testing.T) {
		slugMock := "zxy01"
		expectedUrl := "http://example.com"

		mock.ExpectGet(slugMock).RedisNil()
		mock.ExpectSet(slugMock, expectedUrl, urlCacheTime).SetVal(expectedUrl)

		shortenerRepository = func() repository.Shortener {
			return mockRepository{
				fnGet: func(slug string) (*model.Shortener, error) {
					return &model.Shortener{
						Slug: slugMock,
						Url:  expectedUrl,
					}, nil
				},
			}
		}

		shortenerService := New()
		url, err := shortenerService.GetUrl(slugMock)

		assert.NoError(t, err)
		assert.Equal(t, expectedUrl, url)
	})

	t.Run("should be return error from invalid slug and not exists shortener", func(t *testing.T) {
		slugMock := "abcX0"

		mock.ExpectGet(slugMock).RedisNil()

		shortenerRepository = func() repository.Shortener {
			return mockRepository{
				fnGet: func(slug string) (*model.Shortener, error) {
					return nil, errors.New("")
				},
			}
		}

		shortenerService := New()
		url, err := shortenerService.GetUrl(slugMock)

		assert.Empty(t, url)
		assert.EqualError(t, err, ErrShortenerNotFound.Error())
	})
}

func TestAddClick(t *testing.T) {
	t.Run("should be return shortener with new valid click and slug", func(t *testing.T) {

		slugMock := "XZY21"
		clickMock := model.Click{
			Source:   "other",
			Device:   "desktop",
			Browser:  "chrome",
			Language: "en",
			System:   "linux",
		}

		wantShortener := &model.Shortener{
			Slug: slugMock,
			Click: []model.Click{
				clickMock,
			},
		}

		shortenerRepository = func() repository.Shortener {
			return mockRepository{
				fnGet: func(slug string) (*model.Shortener, error) {
					return &model.Shortener{Slug: slugMock}, nil
				},
				fnUpdate: func(shortener *model.Shortener, id string) (*model.Shortener, error) {
					return shortener, nil
				},
			}
		}

		shortenerService := New()
		gotShortener, err := shortenerService.AddClick(clickMock, slugMock)

		assert.NoError(t, err)
		assert.Equal(t, wantShortener, gotShortener)
	})
}

func TestStats(t *testing.T) {
	slugMock := "SL5G3"
	clicksMock := []model.Click{
		{
			Browser: "chrome",
		},
		{
			Browser: "chrome",
		},
		{
			Browser: "safari",
		},
	}

	shortenerRepository = func() repository.Shortener {
		return &mockRepository{
			fnGet: func(slug string) (*model.Shortener, error) {
				return &model.Shortener{
					Click: clicksMock,
				}, nil
			},
		}
	}

	shortenerService := New()
	stats, err := shortenerService.Stats(slugMock)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, len(clicksMock), stats.Clicks)
	assert.Equal(t, 2, stats.Browsers["chrome"])
	assert.Equal(t, 1, stats.Browsers["safari"])
}
