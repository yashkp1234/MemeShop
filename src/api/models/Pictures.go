package models

import (
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Picture is the model that represents a picture object
type Picture struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	URL       string             `json:"url" bson:"url"`
	HashKey   string             `json:"hash_key" bson:"hash_key"`
	User      string             `json:"user" bson:"user"`
	ForSale   bool               `json:"for_sale" bson:"for_sale"`
	Title     string             `json:"title" bson:"title"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

//ValidatePictureUpdate validates a picture update map
func ValidatePictureUpdate(updates map[string]string) error {
	for k, v := range updates {
		switch k {
		case "title":
			if v == "" {
				return errors.New("Title cannot be empty")
			}
		case "user":
			if v == "" {
				return errors.New("Picture must be owned by a user")
			}
		default:
			return errors.New("Invalid key sent")
		}
	}
	return nil
}

//Validate validates a picture struct
func (p *Picture) Validate() error {
	if p.Title == "" {
		return errors.New("Title cannot be empty")
	}
	if p.User == "" {
		return errors.New("Picture must be owned by a user")
	}
	return nil
}

//SetUp sets the picture information on creation
func (p *Picture) SetUp() error {
	err := p.Validate()
	if err != nil {
		return err
	}

	//Setup user
	p.ID = primitive.NewObjectID()
	p.CreatedAt = time.Now()
	return nil
}

//MarshalBinary marshals a picture into binary
func (p Picture) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}
