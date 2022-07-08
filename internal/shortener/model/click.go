package model

import (
	"time"
)

type Click struct {
	ID        string
	Source    string
	Device    string
	Browser   string
	Language  string
	System    string
	CreatedAt time.Time
}
