package crud

import (
	"context"
	"errors"
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
		defer close(ch)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = user.SetUp()
		if err != nil {
			ch <- false
			return
		}

		if err = r.db.FindOne(ctx, bson.M{"username": user.UserName}).Err(); err == nil {
			err = errors.New("Username already exists")
			ch <- false
			return
		}

		_, err = r.db.InsertOne(ctx, user)
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return user, nil
	}

	return models.User{}, err
}

//FindByID finds a user by id
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
		defer close(ch)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return user, nil
	}

	return models.User{}, err
}

//Delete deletes a user from db
func (r *RepositoryUsersCRUD) Delete(id string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	var user models.User
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = r.db.FindOneAndDelete(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return id, nil
	}

	return "", err
}

//Update updates a user from db
func (r *RepositoryUsersCRUD) Update(id string, user models.User) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user.HashPassword()

		updateFields := bson.M{
			"password":   user.Password,
			"updated_at": time.Now(),
		}
		update := bson.M{"$set": updateFields}

		if err = r.db.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update).Decode(&user); err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return nil
	}

	return err
}
