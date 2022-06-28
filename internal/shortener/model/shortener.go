package model

import "time"

type Shortener struct {
	ID        string
	Url       string
	Slug      string
	Click     []Click
	CreatedAt time.Time
}
