package database

import (
	"context"
	"time"

	"github.com/yashkp1234/MemeShop.git/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Connect connects to the mongoDB client
//and returns the database
func Connect() (*mongo.Database, error) {
	config.Load()

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURL))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(config.DBName), nil
}
