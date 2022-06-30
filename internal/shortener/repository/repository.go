package repository

import "github.com/CaioAureliano/gortener/internal/shortener/model"

type Shortener interface {
	Create(shortener *model.Shortener) (*model.Shortener, error)
	Get(slug string) (*model.Shortener, error)
	Update(shortener *model.Shortener, id string) (*model.Shortener, error)
	AddClick(click model.Click, slug string) (*model.Shortener, error)
}

type shortener struct {
}

func New() Shortener {
	return shortener{}
}

func (s shortener) Create(shortener *model.Shortener) (*model.Shortener, error) {
	return nil, nil
}

func (s shortener) Get(slug string) (*model.Shortener, error) {
	return nil, nil
}

func (s shortener) Update(shortener *model.Shortener, id string) (*model.Shortener, error) {
	return nil, nil
}

func (s shortener) AddClick(click model.Click, slug string) (*model.Shortener, error) {
	return nil, nil
}
