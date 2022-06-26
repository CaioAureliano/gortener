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
					return &model.Shortener{
						Url: "http://google.com",
					}, nil
				},
			},
		},
		{
			name:    "invalid url",
			gotUrl:  "url",
			wantUrl: "",
			wantErr: ErrInvalidURL,
			repoMock: mockRepository{
				fnCreate: func(shortener *model.Shortener) (*model.Shortener, error) {
					return nil, ErrInvalidURL
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
}
