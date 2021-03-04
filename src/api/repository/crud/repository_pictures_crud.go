package crud

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yashkp1234/MemeShop.git/api/models"
	"github.com/yashkp1234/MemeShop.git/api/utils/channels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//RepositoryPicturesCRUD object to store Picture CRUD operations
type RepositoryPicturesCRUD struct {
	db *mongo.Collection
}

//NewRepositoryPicturesCRUD creates a new RepositoryPicturesCRUD object
func NewRepositoryPicturesCRUD(db *mongo.Database) *RepositoryPicturesCRUD {
	return &RepositoryPicturesCRUD{db.Collection("pictures")}
}

//Save saves a picture onto database collection
func (r *RepositoryPicturesCRUD) Save(picture models.Picture) (models.Picture, error) {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = picture.SetUp()
		if err != nil {
			ch <- false
			return
		}

		//Todo ensure uniqueness

		_, err = r.db.InsertOne(ctx, picture)
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return picture, nil
	}

	return models.Picture{}, err
}

//FindByID saves a picture onto database collection
func (r *RepositoryPicturesCRUD) FindByID(username string, pictureID string) (models.Picture, error) {
	//Create id from string
	objID, err := primitive.ObjectIDFromHex(pictureID)
	if err != nil {
		return models.Picture{}, err
	}

	//Init variables
	var picture models.Picture
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&picture); err != nil {
			ch <- false
			return
		}

		if picture.User != username && !picture.ForSale {
			err = errors.New("No picture found")
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return picture, nil
	}

	return models.Picture{}, err
}

//Delete deletes a picture from db
func (r *RepositoryPicturesCRUD) Delete(username string, idPicture string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(idPicture)
	if err != nil {
		return "", err
	}

	var picture models.Picture
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{
			"_id":  objID,
			"user": username,
		}
		if err = r.db.FindOneAndDelete(ctx, filter).Decode(&picture); err != nil {
			fmt.Println(err)
			err = errors.New("Error deleting picture, picture not found")
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return idPicture, nil
	}

	return "", err
}

//Update updates a picture from db
func (r *RepositoryPicturesCRUD) Update(id string, username string, updates map[string]string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		update := bson.M{"$set": updates}
		filter := bson.M{
			"_id":  objID,
			"user": username,
		}

		if err = r.db.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
			fmt.Println(err)
			err = errors.New("Error deleting picture, picture not found")
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
