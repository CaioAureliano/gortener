package service

import (
	"errors"
	"testing"
	"time"

	"github.com/CaioAureliano/gortener/internal/shortener/dto"
	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/repository"
	"github.com/CaioAureliano/gortener/internal/shortener/repository/cache"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		mockUrl := &dto.UrlRequest{Url: "www.google.com"}
		expectedUrl := "http://" + mockUrl.Url

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
			gotUrl   *dto.UrlRequest
			wantUrl  string
			wantErr  error
			repoMock repository.Shortener
		}{
			{
				name:    "should be return a shortener created with valid URL",
				gotUrl:  &dto.UrlRequest{Url: "google.com"},
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
				gotUrl:  nil,
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

		var stubClicks []model.Click

		shortenerRepository = func() repository.Shortener {
			return mockRepository{
				fnGet: func(slug string) (*model.Shortener, error) {
					return &model.Shortener{
						Slug: slugMock,
						Click: []model.Click{
							{
								ID:      primitive.NewObjectID(),
								Browser: "firefox",
							},
						},
					}, nil
				},
				fnAddClick: func(clicks []model.Click, id primitive.ObjectID) error {
					stubClicks = clicks
					return nil
				},
			}
		}

		shortenerService := New()
		err := shortenerService.AddClick(clickMock, slugMock)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(stubClicks))
		assert.Equal(t, clickMock.Source, stubClicks[1].Source)
		assert.Equal(t, clickMock.Device, stubClicks[1].Device)
		assert.Equal(t, clickMock.Browser, stubClicks[1].Browser)
		assert.Equal(t, clickMock.Language, stubClicks[1].Language)
		assert.Equal(t, clickMock.System, stubClicks[1].System)
	})
}

func TestStats(t *testing.T) {
	t.Run("should be cache a stats", func(t *testing.T) {
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

		clientMock, mock := redismock.NewClientMock()
		cacheRepository = func() cache.Cache {
			return NewMockCache(clientMock)
		}

		keyPrefix := slugMock + slugStatsKeyCachePrefix
		mock.ExpectGet(keyPrefix).RedisNil()
		mock.Regexp().ExpectSet(keyPrefix, `[a-z]+`, statsCacheTime).SetVal("")

		shortenerService := New()
		stats, err := shortenerService.Stats(slugMock)

		if assert.NoError(t, err) && assert.NotNil(t, stats) {
			assert.Equal(t, len(clicksMock), stats.Clicks)
			assert.Equal(t, 2, stats.Browsers["chrome"])
			assert.Equal(t, 1, stats.Browsers["safari"])
		}
	})

	t.Run("should be get stats cached", func(t *testing.T) {
		slug := "xyz01"
		statsEncodedMock := `{
			"browsers": {
				"chrome": 10,
				"safari": 2,
				"firefox": 5
			}
		}`

		shortenerRepository = func() repository.Shortener {
			return mockRepository{}
		}

		clientMock, mock := redismock.NewClientMock()
		cacheRepository = func() cache.Cache {
			return NewMockCache(clientMock)
		}

		mock.ExpectGet(slug + slugStatsKeyCachePrefix).SetVal(statsEncodedMock)

		shortenerService := New()
		res, err := shortenerService.Stats(slug)

		if assert.NoError(t, err) && assert.NotNil(t, res) {
			assert.Equal(t, 10, res.Browsers["chrome"])
			assert.Equal(t, 5, res.Browsers["firefox"])
			assert.Equal(t, 2, res.Browsers["safari"])
		}
	})
}
