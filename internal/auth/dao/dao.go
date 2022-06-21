package dao

import (
	"context"
	"time"

	"github.com/CaioAureliano/gortener/internal/auth/model"
	"github.com/CaioAureliano/gortener/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userDao struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserDao() *userDao {
	return &userDao{
		ctx:        context.Background(),
		collection: database.ConnectDatabase().Collection(USER_COLLECTION_NAME),
	}
}

const USER_COLLECTION_NAME = "users"

func (d *userDao) Create(user *model.User) error {
	user.CreatedAt = time.Now()

	b, err := bson.Marshal(user)
	if err != nil {
		return err
	}

	if _, err := d.collection.InsertOne(d.ctx, b); err != nil {
		return err
	}

	optIndex := options.Index().SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: optIndex}

	if _, err := d.collection.Indexes().CreateOne(d.ctx, index); err != nil {
		return err
	}

	return nil
}

func (d *userDao) GetByField(value, field string) (*model.User, error) {
	var user *model.User
	if err := d.collection.FindOne(d.ctx, bson.M{field: value}, nil).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (d *userDao) ExistsByEmail(email string) (bool, error) {
	if user, err := d.GetByField(email, "email"); user == nil || err != nil {
		return false, err
	}
	return true, nil
}
