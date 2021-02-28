package crud

import (
	"context"
	"time"

	"github.com/yashkp1234/MemeShop.git/api/models"
	"github.com/yashkp1234/MemeShop.git/api/utils/channels"
	"go.mongodb.org/mongo-driver/mongo"
)

//RepositoryUsersCRUD object to store user CRUD operations
type RepositoryUsersCRUD struct {
	db *mongo.Collection
}

//NewRepositoryUsersCRUD creates a new RepositoryUsersCRUD object
func NewRepositoryUsersCRUD(db *mongo.Database) *RepositoryUsersCRUD {
	return &RepositoryUsersCRUD{db.Collection("users")}
}

//Save saves a user onto database collection
func (r *RepositoryUsersCRUD) Save(user models.User) (models.User, error) {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		//Create a context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//Hash Password and setup user on creation
		err = user.HashPassword()
		if err != nil {
			ch <- false
			return
		}
		user.SetUp()

		//Insert into db
		_, err = r.db.InsertOne(ctx, user)
		if err != nil {
			ch <- false
			return
		}

		//Return done
		ch <- true
	}(done)

	//Return user if no errors
	if channels.OK(done) {
		return user, nil
	}

	//Return error
	return models.User{}, err
}
