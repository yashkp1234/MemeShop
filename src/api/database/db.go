package database

import (
	"context"
	"log"

	"github.com/yashkp1234/MemeShop.git/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DBInstance represents a DB instance
var dbInstance *mongo.Database

//NewDatabase setups creating a databse connection
func NewDatabase() {
	var err error
	dbInstance, err = connectToMongo()
	if err != nil {
		log.Fatal(err)
	}
}

//CancelDB disconnects with the database
func CancelDB() {
	dbInstance.Client().Disconnect(context.Background())
}

//Connect returns the connected client
func Connect() *mongo.Database {
	return dbInstance
}

//connectToMongo actually connects to mongoDB
func connectToMongo() (*mongo.Database, error) {
	var err error
	session, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURL))
	log.Println("Connected with db")
	if err != nil {
		return nil, err
	}
	session.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	return session.Database(config.DBName), nil
}
