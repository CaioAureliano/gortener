package service

import (
	"context"
	"time"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/internal/shortener/repository/cache"
	"github.com/go-redis/redis/v8"
)

type mockRepository struct {
	fnCreate   func(shortener *model.Shortener) (*model.Shortener, error)
	fnGet      func(slug string) (*model.Shortener, error)
	fnUpdate   func(shortener *model.Shortener, id string) (*model.Shortener, error)
	fnAddClick func(click model.Click, id string) (*model.Shortener, error)
}

func (m mockRepository) Create(shortener *model.Shortener) (*model.Shortener, error) {
	if m.fnCreate == nil {
		return nil, nil
	}
	return m.fnCreate(shortener)
}

func (m mockRepository) Get(slug string) (*model.Shortener, error) {
	if m.fnGet == nil {
		return nil, nil
	}
	return m.fnGet(slug)
}

func (m mockRepository) Update(shortener *model.Shortener, id string) (*model.Shortener, error) {
	if m.fnUpdate == nil {
		return nil, nil
	}
	return m.fnUpdate(shortener, id)
}

func (m mockRepository) AddClick(click model.Click, slug string) (*model.Shortener, error) {
	if m.fnAddClick == nil {
		return nil, nil
	}
	return m.fnAddClick(click, slug)
}

type mockCache struct {
	ctx    context.Context
	client *redis.Client
}

func NewMockCache(client *redis.Client) cache.Cache {
	return mockCache{
		ctx:    context.Background(),
		client: client,
	}
}

func (c mockCache) Set(key, value string, duration time.Duration) error {
	return c.client.Set(c.ctx, key, value, duration).Err()
}

func (c mockCache) Get(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}
