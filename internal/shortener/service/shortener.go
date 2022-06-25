package service

import "github.com/CaioAureliano/gortener/internal/shortener/model"

type Shortener interface {
	Create(url string) *model.Shortener
}

type shortener struct {
}

func New() Shortener {
	return &shortener{}
}

func (s *shortener) Create(url string) *model.Shortener {

	return &model.Shortener{}
}
