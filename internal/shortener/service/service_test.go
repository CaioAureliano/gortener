package service

import (
	"testing"
	"time"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/repository"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	fnCreate func(shortener *model.Shortener) (*model.Shortener, error)
}

func (m mockRepository) Create(shortener *model.Shortener) (*model.Shortener, error) {
	if m.fnCreate == nil {
		return nil, nil
	}
	return m.fnCreate(shortener)
}

func TestCreate(t *testing.T) {
	expectedUrl := "http://google.com"
	urlMock := "google.com"

	repositoryNew = func() repository.Shortener {
		return mockRepository{
			fnCreate: func(shortener *model.Shortener) (*model.Shortener, error) {
				return &model.Shortener{
					Url:       expectedUrl,
					Slug:      "ABCD",
					CreatedAt: time.Now(),
				}, nil
			},
		}
	}

	shortenerService := New()
	shortUrlCreated, err := shortenerService.Create(urlMock)

	assert.Equal(t, expectedUrl, shortUrlCreated.Url)
	assert.NoError(t, err)
	assert.NotEmpty(t, shortUrlCreated.Slug)
	assert.NotEmpty(t, shortUrlCreated.CreatedAt)
}
