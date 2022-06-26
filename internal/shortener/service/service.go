package service

import (
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

const (
	SLUG_LENGTH = 5
)

var (
	repositoryNew = repository.New
)

func (s *shortener) Create(url string) (*model.Shortener, error) {
	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	slug := randutil.RandomString(SLUG_LENGTH)

	shortToCreate := &model.Shortener{
		Url:       url,
		Slug:      slug,
		CreatedAt: time.Now(),
	}

	createdUrl, _ := repositoryNew().Create(shortToCreate)

	return createdUrl, nil
}
