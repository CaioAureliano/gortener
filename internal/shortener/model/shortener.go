package model

import "time"

type Shortener struct {
	ID        string
	Url       string
	Slug      string
	CreatedAt time.Time
}
