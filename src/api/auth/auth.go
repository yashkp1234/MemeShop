package auth

import (
	"context"
	"time"

	"github.com/yashkp1234/MemeShop.git/api/database"
	"github.com/yashkp1234/MemeShop.git/api/models"
	"github.com/yashkp1234/MemeShop.git/api/security"
	"go.mongodb.org/mongo-driver/bson"
)

//SignIn handles the sign in for a user
func SignIn(username, password string) (string, error) {
	db := database.Connect()

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	if err := db.Collection("users").FindOne(ctx, filter).Decode(&user); err != nil {
		return "", err
	}

	err := security.VerifyPassword(user.Password, password)
	if err != nil {
		return "", err
	}

	return CreateToken(user.ID)

}
