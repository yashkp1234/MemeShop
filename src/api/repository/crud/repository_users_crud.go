package crud

import (
	"context"
	"time"

	"github.com/yashkp1234/MemeShop.git/api/models"
	"github.com/yashkp1234/MemeShop.git/api/utils/channels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		ch <- true
	}(done)

	//Return user if no errors
	if channels.OK(done) {
		return user, nil
	}

	//Return error
	return models.User{}, err
}

//FindByID saves a user onto database collection
func (r *RepositoryUsersCRUD) FindByID(id string) (models.User, error) {
	//Create id from string
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	//Init variables
	var user models.User
	done := make(chan bool)

	go func(ch chan<- bool) {
		//Create a context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//Find in DB
		if err = r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	//Return user if no errors
	if channels.OK(done) {
		return user, nil
	}

	//Return error
	return models.User{}, err
}

//Delete deletes a user from db
func (r *RepositoryUsersCRUD) Delete(id string) (string, error) {
	//Create id from string
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	//Init variables
	var user models.User
	done := make(chan bool)

	go func(ch chan<- bool) {
		//Create a context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//Find in DB
		if err = r.db.FindOneAndDelete(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	//Return user if no errors
	if channels.OK(done) {
		return id, nil
	}

	//Return error
	return "", err
}

//Update updates a user from db
func (r *RepositoryUsersCRUD) Update(id string, changePassword bool, user models.User) error {
	//Create id from string
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	//Init variables
	done := make(chan bool)

	go func(ch chan<- bool) {
		//Create a context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		updateFields := bson.M{
			"username":  user.UserName,
			"updatedat": time.Now(),
		}
		if changePassword {
			user.HashPassword()
			updateFields["password"] = user.Password
		}
		update := bson.M{"$set": updateFields}

		//Update in DB
		if err = r.db.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update).Decode(&user); err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	//Return user if no errors
	if channels.OK(done) {
		return nil
	}

	//Return error
	return err
}
