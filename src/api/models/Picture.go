package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Picture is the model that represents a picture object
type Picture struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	URL       string             `json:"url" bson:"url"`
	Hash      string             `json:"hash" bson:"hash"`
	User      User               `json:"user" bson:"user"`
	ForSale   bool               `json:"for_sale" bson:"for_sale"`
	Title     string             `json:"title" bson:"title"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
