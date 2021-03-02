package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Picture is the model that represents a picture object
type Picture struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	URL       string             `json:"url" bson:"url"`
	Hash      string             `json:"hash" bson:"hash"`
	User      string             `json:"user" bson:"user"`
	ForSale   bool               `json:"for_sale" bson:"for_sale"`
	Title     string             `json:"title" bson:"title"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

//Validate validates a picture struct
func (p *Picture) Validate() error {
	if p.Title == "" {
		return errors.New("Title cannot be empty")
	}
	if p.User == "" {
		return errors.New("Picture must be owned by a user")
	}
	if p.URL == "" && p.Hash == "" {
		return errors.New("Picture must have a hash source or url source")
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
	p.UpdatedAt = time.Now()
	return nil
}
