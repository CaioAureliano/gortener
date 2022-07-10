package model

import "time"

type Shortener struct {
	ID        string    `json:"id,omitempty"`
	Url       string    `json:"url,omitempty"`
	Slug      string    `json:"slug,omitempty"`
	Click     []Click   `json:"clicks,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
