package service

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/repository"
	"github.com/CaioAureliano/gortener/internal/shortener/repository/cache"
	"github.com/go-redis/redis/v8"
	"github.com/snapcore/snapd/randutil"
)

type Shortener interface {
	Create(url string) (*model.Shortener, error)
	Get(slug string) (*model.Shortener, error)
	GetUrl(slug string) (string, error)
	AddClick(click model.Click, slug string) (*model.Shortener, error)
	Stats(slug string) (*model.Stats, error)
}

type shortener struct {
}

func New() Shortener {
	return &shortener{}
}

var (
	ErrInvalidURL        = errors.New("invalid url")
	ErrInvalidSlug       = errors.New("invalid slug")
	ErrShortenerNotFound = errors.New("not found short URL")

	shortenerRepository = repository.New
	cacheRepository     = cache.New
)

const (
	slugLength   = 5
	urlCacheTime = time.Hour * 24 * 8

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

	defer s.CacheUrl(shortToCreate.Slug, shortToCreate.Url)
	createdUrl, err := shortenerRepository().Create(shortToCreate)
	if err != nil {
		return nil, err
	}

	return createdUrl, nil
}

func (s *shortener) Get(slug string) (*model.Shortener, error) {
	reqSlugLenght := len([]byte(slug))
	if reqSlugLenght != slugLength {
		log.Printf("invalid slug length: %s", slug)
		return nil, ErrInvalidSlug
	}

	res, err := shortenerRepository().Get(slug)
	if err != nil {
		log.Printf("shortener not found with slug %s - error: %s", slug, err.Error())
		return nil, ErrShortenerNotFound
	}

	return res, nil
}

func (s *shortener) GetUrl(slug string) (string, error) {
	res, err := cacheRepository().Get(slug)
	if err == redis.Nil {

		log.Printf("not found url from slug: \"%s\" in cache - error: %v", slug, err)

		shortener, err := s.Get(slug)
		if err != nil || shortener == nil {
			return "", ErrShortenerNotFound
		}

		defer s.CacheUrl(shortener.Slug, shortener.Url)
		return shortener.Url, nil
	}

	return res, nil
}

func (s *shortener) AddClick(click model.Click, slug string) (*model.Shortener, error) {
	shortener, err := s.Get(slug)
	if err != nil {
		return nil, err
	}

	clicks := shortener.Click
	clicks = append(clicks, click)
	shortener.Click = clicks

	updated, err := shortenerRepository().Update(shortener, shortener.ID)
	if err != nil {
		log.Printf("error to update shortener: %v with click: %v", shortener, click)
		return nil, errors.New("error to update shortener with click")
	}

	return updated, nil
}

func (s *shortener) Stats(slug string) (*model.Stats, error) {
	shortener, err := s.Get(slug)
	if err != nil {
		return nil, err
	}

	clicks := shortener.Click

	stats := &model.Stats{}
	stats.Initialize()

	stats.Clicks = len(clicks)

	for _, c := range clicks {
		stats.IncrementIfExists(c)
	}

	return stats, nil
}

func (s shortener) CacheUrl(slug, url string) {
	if err := cacheRepository().Set(slug, url, urlCacheTime); err != nil {
		log.Printf("error to cache slug: \"%s\" with url: \"%s\" - error: %s", slug, url, err.Error())
	}
}
