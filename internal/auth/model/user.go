package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type UserCreateRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}
