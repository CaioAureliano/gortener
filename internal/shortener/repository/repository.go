package repository

import (
	"context"
	"time"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShortenerRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewShortenerRepository() *ShortenerRepository {
	return &ShortenerRepository{
		ctx:        context.Background(),
		collection: database.ConnectDatabase().Collection(SHORTENER_COLLECTION_NAME),
	}
}

var SHORTENER_COLLECTION_NAME = "shortener"

func (sr *ShortenerRepository) Create(s *model.Shortener) error {
	s.CreatedAt = time.Now()

	b, err := bson.Marshal(s)
	if err != nil {
		return err
	}

	if _, err := sr.collection.InsertOne(sr.ctx, b); err != nil {
		return err
	}

	opt := options.Index().SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"slug": 1}, Options: opt}

	if _, err := sr.collection.Indexes().CreateOne(sr.ctx, index); err != nil {
		return err
	}

	return nil
}
