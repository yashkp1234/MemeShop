package models

import (
	"errors"
	"time"

	"github.com/yashkp1234/MemeShop.git/api/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that represents a user object
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserName  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Balance   int32              `json:"balance" bson:"balance"`
}

const (
	//MinPassLength is Minimum length for a password
	MinPassLength = 6
	//MinUserLength is Minimum length for a username
	MinUserLength = 2
)

//ValidateUpdate validates a user struct on update
func (u *User) ValidateUpdate() error {
	if len(u.Password) < MinPassLength {
		return errors.New("Password is too short")
	}
	return nil
}

//ValidateCreation validates a user struct on creation
func (u *User) ValidateCreation() error {
	if len(u.Password) < MinPassLength {
		return errors.New("Password is too short")
	}
	if len(u.UserName) < MinUserLength {
		return errors.New("Username is too short")
	}
	if u.Balance > 0 {
		return errors.New("Cannot set balance on a user")
	}
	return nil
}

//SetUp sets the users information on creation
func (u *User) SetUp() error {
	err := u.ValidateCreation()
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
