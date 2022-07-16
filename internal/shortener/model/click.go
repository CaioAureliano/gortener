package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Click struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Source    string             `bson:"source,omitempty"`
	Device    string             `bson:"device,omitempty"`
	Browser   string             `bson:"browser,omitempty"`
	Language  string             `bson:"language,omitempty"`
	System    string             `bson:"system,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
}
