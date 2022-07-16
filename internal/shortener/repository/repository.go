package repository

import (
	"context"
	"log"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
	"github.com/CaioAureliano/gortener/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Shortener interface {
	Create(shortener *model.Shortener) (*model.Shortener, error)
	Get(slug string) (*model.Shortener, error)
	Update(shortener *model.Shortener, id primitive.ObjectID) (*model.Shortener, error)
	AddClick(click model.Click, slug string) (*model.Shortener, error)
}

type shortener struct {
	ctx  context.Context
	coll *mongo.Collection
}

const (
	shortenerCollectionName = "shorteners"
)

func New() Shortener {
	return shortener{
		ctx:  context.Background(),
		coll: database.Connect().Collection(shortenerCollectionName),
	}
}

func (s shortener) Create(shortener *model.Shortener) (*model.Shortener, error) {
	res, err := s.coll.InsertOne(s.ctx, shortener)
	if err != nil {
		log.Printf("error to create a new short url: %s", err.Error())
		return nil, err
	}

	shortener.ID = res.InsertedID.(primitive.ObjectID)

	index := mongo.IndexModel{
		Keys:    bson.M{"slug": 1},
		Options: options.Index().SetUnique(true),
	}

	defer s.coll.Indexes().CreateOne(s.ctx, index)

	return shortener, nil
}

func (s shortener) Get(slug string) (*model.Shortener, error) {
	var res *model.Shortener
	if err := s.coll.FindOne(s.ctx, bson.M{"slug": slug}).Decode(&res); err != nil {
		log.Printf("error to find short url: %s", err.Error())
		return nil, err
	}

	return res, nil
}

func (s shortener) Update(shortener *model.Shortener, id primitive.ObjectID) (*model.Shortener, error) {
	_, err := s.coll.UpdateByID(s.ctx, id, bson.M{"$set": bson.M{
		"clicks": shortener.Click,
	}})

	if err != nil {
		log.Printf("error to update short url: %s", err.Error())
		return nil, err
	}

	return s.Get(shortener.Slug)
}

func (s shortener) AddClick(click model.Click, slug string) (*model.Shortener, error) {
	return nil, nil
}
