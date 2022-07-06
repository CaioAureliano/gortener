package handler

import "github.com/CaioAureliano/gortener/internal/shortener/model"

type mockService struct {
	fnCreate   func(url string) (*model.Shortener, error)
	fnGet      func(slug string) (*model.Shortener, error)
	fnGetUrl   func(slug string) (string, error)
	fnAddClick func(click model.Click, slug string) (*model.Shortener, error)
	fnStats    func(slug string) (*model.Stats, error)
}

func (m mockService) Create(url string) (*model.Shortener, error) {
	if m.fnCreate != nil {
		return m.fnCreate(url)
	}
	return nil, nil
}

func (m mockService) Get(slug string) (*model.Shortener, error) {
	return nil, nil
}

func (m mockService) GetUrl(slug string) (string, error) {
	return "", nil
}

func (m mockService) AddClick(click model.Click, slug string) (*model.Shortener, error) {
	return nil, nil
}

func (m mockService) Stats(slug string) (*model.Stats, error) {
	return nil, nil
}
