package models

import (
	"errors"
	"time"

	"github.com/yashkp1234/MemeShop.git/api/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that represents a user object
type User struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id"`
	UserName     string               `json:"username" bson:"username"`
	Password     string               `json:"password" bson:"password"`
	Token        string               `json:"token" bson:"token"`
	RefreshToken string               `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at" bson:"updated_at"`
	Pictures     []primitive.ObjectID `json:"pictures,omitempty" bson:"pictures,omitempty"`
}

const (
	//MinPassLength is Minimum length for a password
	MinPassLength = 6
	//MinUserLength is Minimum length for a username
	MinUserLength = 2
)

//Validate validates a user struct
func (u *User) Validate(pass bool) error {
	if pass && len(u.Password) < MinPassLength {
		return errors.New("Password is too short")
	}
	if len(u.UserName) < MinUserLength {
		return errors.New("Username is too short")
	}
	return nil
}

//SetUp sets the users information on creation
func (u *User) SetUp() error {
	err := u.Validate(true)
	if err != nil {
		return err
	}

	err = u.HashPassword()
	if err != nil {
		return err
	}

	//Setup user
	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
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
