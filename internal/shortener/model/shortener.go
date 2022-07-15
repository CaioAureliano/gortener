package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shortener struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Url       string             `json:"url,omitempty" bson:"url"`
	Slug      string             `json:"slug,omitempty" bson:"slug"`
	Click     []Click            `json:"clicks,omitempty" bson:"clicks"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
}
