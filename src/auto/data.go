package auto

import (
	"time"

	"github.com/yashkp1234/MemeShop.git/api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GenerateUser generates a user
func GenerateUser(userName string, Password string) models.User {
	return models.User{
		ID:           primitive.NewObjectID(),
		UserName:     userName,
		Password:     Password,
		Token:        "b",
		RefreshToken: "a",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

var users = []models.User{
	{
		ID:           primitive.NewObjectID(),
		UserName:     "asdasd",
		Password:     "testing",
		Token:        "b",
		RefreshToken: "a",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
}
