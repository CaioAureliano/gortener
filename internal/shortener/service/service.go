package service

import (
	"errors"
	"log"
	"regexp"

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

const LENGTH_RANDOM_SLUG = 6

func (ss *ShortenerService) Create(req *model.ShortenerCreateRequest, user *jwt.Token) (string, error) {
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	randomSlug := randutil.RandomString(LENGTH_RANDOM_SLUG)

	if !isValidUrl(req.Url) {
		return "", errors.New("invalid url")
	}

	s := &model.Shortener{
		Url:          req.Url,
		Slug:         randomSlug,
		CreatorEmail: email,
	}

	if err := shortenerRepository.Create(s); err != nil {
		log.Printf("error to create a shortener: %s", err.Error())
		return "", errors.New("error to create a shortener")
	}

	return randomSlug, nil
}

func isValidUrl(url string) bool {
	valid, err := regexp.MatchString(`[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`, url)
	if err != nil {
		log.Print(err.Error())
	}

	return valid
}
