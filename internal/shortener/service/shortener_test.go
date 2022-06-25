package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	expectedUrl := "http://google.com"
	urlMock := "google.com"

	shortenerService := New()
	shortUrlCreated := shortenerService.Create(urlMock)

	assert.Equal(t, expectedUrl, shortUrlCreated.Url)
}
