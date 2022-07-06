package service

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/CaioAureliano/gortener/internal/shortener/dto"
	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/repository"
	"github.com/CaioAureliano/gortener/internal/shortener/repository/cache"
	"github.com/go-redis/redis/v8"
	"github.com/snapcore/snapd/randutil"
)

type Shortener interface {
	Create(url *dto.UrlRequest) (*model.Shortener, error)
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
	slugLength     = 5
	urlCacheTime   = time.Hour * 24 * 7
	statsCacheTime = time.Minute * 10

	slugStatsKeyCachePrefix = "_stats"
)

func (s *shortener) Create(req *dto.UrlRequest) (*model.Shortener, error) {
	if req == nil || !req.IsValid() {
		return nil, ErrInvalidURL
	}

	req.AppendProtocolIfNotExists()

	slug := randutil.RandomString(slugLength)

	shortToCreate := &model.Shortener{
		Url:       req.Url,
		Slug:      slug,
		CreatedAt: time.Now(),
	}

	defer s.cacheUrl(shortToCreate.Slug, shortToCreate.Url)
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

		defer s.cacheUrl(shortener.Slug, shortener.Url)
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
	encodedStatsRes, err := cacheRepository().Get(s.getSlugStatsKey(slug))
	if err == redis.Nil {
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

		encodedStats, err := json.Marshal(stats)
		if err != nil {
			return nil, err
		}

		defer cacheRepository().Set(s.getSlugStatsKey(slug), string(encodedStats), statsCacheTime)

		return stats, nil
	}

	var stats *model.Stats
	if err := json.Unmarshal([]byte(encodedStatsRes), &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func (s shortener) cacheUrl(slug, url string) {
	if err := cacheRepository().Set(slug, url, urlCacheTime); err != nil {
		log.Printf("error to cache slug: \"%s\" with url: \"%s\" - error: %s", slug, url, err.Error())
	}
}

func (s shortener) getSlugStatsKey(slug string) string {
	return slug + slugStatsKeyCachePrefix
}
