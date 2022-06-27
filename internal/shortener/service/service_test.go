package service

import (
	"testing"

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

	t.Run("model response", func(t *testing.T) {
		mockUrl := "www.google.com"

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
		assert.Equal(t, "http://"+mockUrl, shortCreated.Url)
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
