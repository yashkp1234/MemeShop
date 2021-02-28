package models

import (
	"time"

	"github.com/yashkp1234/MemeShop.git/api/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserName     string             `json:"username" bson:"username" validate:"required,min=2,max=100"`
	Password     string             `json:"password" bson:"password" validate:"required,min=6"`
	Token        string             `json:"token" bson:"token"`
	RefreshToken string             `json:"refresh_token" bson:"refreshtoken"`
	CreatedAt    time.Time          `json:"created_at" bson:"createdat"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updatedat"`
}

//SetUp sets the users information on creation
func (u *User) SetUp() {
	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

//HashPassword hashes password of user
func (u *User) HashPassword() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
