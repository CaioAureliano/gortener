package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	connection_url = os.Getenv("MONGO_URI")
	db_name        = os.Getenv("DB_NAME")
)

func ConnectDatabase() *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connection_url))
	if err != nil {
		log.Fatalf("error to connect to database: %s", err.Error())
	}

	return client.Database(db_name)
}
