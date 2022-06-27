package service

import (
	"errors"
	"log"
	"regexp"
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

const (
	slugLength = 5

	regexValidURL = `[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
)

func (s *shortener) Create(url string) (*model.Shortener, error) {
	isValid, err := regexp.MatchString(regexValidURL, url)
	if err != nil || !isValid {
		log.Printf("invalid url: %s", url)
		return nil, ErrInvalidURL
	}

	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	slug := randutil.RandomString(slugLength)

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
