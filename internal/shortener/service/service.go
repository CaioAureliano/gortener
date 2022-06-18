package service

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/repository"
	"github.com/golang-jwt/jwt"
	"github.com/snapcore/snapd/randutil"
)

type ShortenerService struct {
}

func NewShortenerService() *ShortenerService {
	return &ShortenerService{}
}

var shortenerRepository = repository.NewShortenerRepository()

const (
	LENGTH_RANDOM_SLUG = 5

	VALID_URL_REGEX = `[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
)

func (ss *ShortenerService) Create(req *model.ShortenerCreateRequest, user *jwt.Token) (string, error) {
	url := req.Url
	randomSlug := randutil.RandomString(LENGTH_RANDOM_SLUG)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	if !isValidUrl(req.Url) {
		return "", errors.New("invalid url: " + url)
	}

	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	if !isValidSite(url) {
		return "", errors.New("invalid site: " + url)
	}

	s := &model.Shortener{
		Url:          url,
		Slug:         randomSlug,
		CreatorEmail: email,
	}

	if err := shortenerRepository.Create(s); err != nil {
		log.Printf("error to create a shortener: %s", err.Error())
		return "", errors.New("error to create a shortener")
	}

	return randomSlug, nil
}

func (ss *ShortenerService) GetUrlBySlug(slug string) (string, error) {
	shortFound, err := shortenerRepository.GetBySlug(slug)
	if err != nil {
		log.Printf("error to find shortener with slug: %s - error: %s", slug, err.Error())
		return "", err
	}
	return shortFound.Url, nil
}

func isValidUrl(url string) bool {
	valid, err := regexp.MatchString(VALID_URL_REGEX, url)
	if err != nil {
		log.Print(err.Error())
	}
	return valid
}

func isValidSite(url string) bool {
	res, err := http.Head(url)
	if err != nil || res.StatusCode != 200 {
		log.Printf("error to head to url: %s - error: %s", url, err.Error())
		return false
	}
	return true
}
