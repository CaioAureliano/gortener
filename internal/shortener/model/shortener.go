package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shortener struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Url          string             `json:"url"`
	Slug         string             `json:"slug"`
	CreatorEmail string             `json:"creator_email" bson:"creator_email"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

type ShortenerCreateRequest struct {
	Url string `json:"url"`
}
