package repository

import "github.com/CaioAureliano/gortener/internal/shortener/model"

type Shortener interface {
	Create(shortener *model.Shortener) (*model.Shortener, error)
}

type shortener struct {
}

func New() Shortener {
	return shortener{}
}

func (s shortener) Create(shortener *model.Shortener) (*model.Shortener, error) {
	return nil, nil
}
