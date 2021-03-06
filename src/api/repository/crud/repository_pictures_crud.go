package crud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" //Image needs this
	_ "image/png"  //Image needs this
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/yashkp1234/MemeShop.git/api/cache"
	"github.com/yashkp1234/MemeShop.git/api/gcp"
	"github.com/yashkp1234/MemeShop.git/api/models"
	"github.com/yashkp1234/MemeShop.git/api/utils/channels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const cacheTime = cache.CacheTime

//RepositoryPicturesCRUD object to store Picture CRUD operations
type RepositoryPicturesCRUD struct {
	db       *mongo.Collection
	cache    *cache.Cache
	imgCloud *gcp.ImageCloudStore
}

//NewRepositoryPicturesCRUD creates a new RepositoryPicturesCRUD object
func NewRepositoryPicturesCRUD(db *mongo.Database, cache *cache.Cache, imgCloud *gcp.ImageCloudStore) *RepositoryPicturesCRUD {
	return &RepositoryPicturesCRUD{db.Collection("pictures"), cache, imgCloud}
}

//Handles processing a file
func handleFile(file *[]byte, handler *multipart.FileHeader) (uint64, string, error) {
	img, _, err := image.Decode(bytes.NewReader(*file))
	if err != nil {
		return 0, "", err
	}

	hash, err := goimagehash.ExtPerceptionHash(img, 8, 8)
	if err != nil {
		return 0, "", err
	}

	return hash.GetHash()[0], strings.ReplaceAll(handler.Filename, " ", ""), nil
}

//Save saves a picture onto database collection
func (r *RepositoryPicturesCRUD) Save(picture models.Picture, file *multipart.File, handler *multipart.FileHeader) (models.Picture, error) {
	var err error

	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)

		var hash uint64
		var id string
		var filename string
		var url string
		var ok bool
		startFile, _ := io.ReadAll(*file)

		hash, filename, err = handleFile(&startFile, handler)
		if err != nil {
			ch <- false
			return
		}

		if ok, err = r.cache.QueryImages(hash); !ok || err != nil {
			err = errors.New("Similar photo already exists")
			ch <- false
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = picture.SetUp()
		if err != nil {
			ch <- false
			return
		}

		id, err = r.cache.AddImage(hash, picture.User+filename)
		if err != nil {
			err = errors.New("Unable to hash photo")
			ch <- false
			return
		}

		picture.HashKey = strings.TrimSpace(fmt.Sprint(id))
		err = r.cache.SyncImages()
		if err != nil {
			ch <- false
			return
		}

		url, err = r.imgCloud.UploadFile(&startFile, picture.ID.Hex())
		if err != nil {
			ch <- false
			return
		}
		picture.URL = url

		_, err = r.db.InsertOne(ctx, picture)
		if err != nil {
			r.cache.DeleteImage(picture.HashKey)
			ch <- false
			return
		}

		err = r.cache.RedisInstance.Set(picture.ID.Hex(), picture, cacheTime).Err()
		if err != nil {
			ch <- false
			return
		}

		err = r.cache.RedisInstance.SAdd(picture.User+"pictures", picture.ID.Hex()).Err()
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
	//Init variables
	var picture models.Picture
	var objID primitive.ObjectID
	var err error
	var str string
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		//Create id from string
		objID, err = primitive.ObjectIDFromHex(pictureID)
		if err != nil {
			ch <- false
			return
		}

		str, err = r.cache.RedisInstance.Get(objID.Hex()).Result()
		if err == nil {
			json.Unmarshal([]byte(str), &picture)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err = r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&picture); err != nil {
				ch <- false
				return
			}
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
			log.Println(err)
			err = errors.New("Error deleting picture, picture not found")
			ch <- false
			return
		}

		if err = r.imgCloud.DeleteFile(picture.ID.Hex()); err != nil {
			ch <- false
			return
		}

		if err = r.cache.RedisInstance.Del(picture.ID.Hex()).Err(); err != nil {
			ch <- false
			return
		}

		r.cache.DeleteImage(picture.HashKey)
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

	var picture models.Picture

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

		if err = r.db.FindOneAndUpdate(ctx, filter, update).Decode(&picture); err != nil {
			err = errors.New("Error deleting picture, picture not found")
			ch <- false
			return
		}

		err = r.cache.RedisInstance.Set(picture.ID.Hex(), picture, cacheTime).Err()
		if err != nil {
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

//FindByUser finds all pictures from a user
func (r *RepositoryPicturesCRUD) FindByUser(username string, setCache bool) ([]models.Picture, error) {
	var pictures []models.Picture
	var err error
	var cursor *mongo.Cursor

	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)

		ids, errz := r.cache.RedisInstance.SMembers(username + "pictures").Result()
		if errz == nil {
			for _, id := range ids {
				pict, errz := r.FindByID(username, id)
				if errz != nil {
					err = errz
					ch <- false
					return
				}
				if setCache || pict.ForSale {
					pictures = append(pictures, pict)
				}
			}
			ch <- true
			return
		}

		filter := bson.M{
			"$or": []interface{}{
				bson.M{"for_sale": true},
				bson.M{"user": username},
			},
		}

		cursor, err = r.db.Find(context.Background(), filter)
		if err != nil {
			ch <- false
			return
		}

		defer cursor.Close(context.Background())
		cursor.All(context.Background(), &pictures)

		if setCache {
			for _, pict := range pictures {
				err := r.cache.RedisInstance.SAdd(username+"pictures", pict.ID.Hex()).Err()
				if err != nil {
					log.Println("Set member error: ", err)
				}
			}
		}

		ch <- true
	}(done)

	if channels.OK(done) {
		return pictures, nil
	}

	return []models.Picture{}, err

}
