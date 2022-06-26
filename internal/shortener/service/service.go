package service

import (
	"errors"
	"strings"
	"time"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/repository"
	"github.com/snapcore/snapd/randutil"
)

type Shortener interface {
	Create(url string) (*model.Shortener, error)
}

type shortener struct {
}

func New() Shortener {
	return &shortener{}
}

var (
	ErrInvalidURL = errors.New("invalid url")

	shortenerRepository = repository.New
)

func (s *shortener) Create(url string) (*model.Shortener, error) {
	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	slug := randutil.RandomString(5)

	shortToCreate := &model.Shortener{
		Url:       url,
		Slug:      slug,
		CreatedAt: time.Now(),
	}

	createdUrl, err := shortenerRepository().Create(shortToCreate)

	if err != nil {
		return nil, err
	}

	return createdUrl, nil
}
