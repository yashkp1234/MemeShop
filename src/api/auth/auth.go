package auth

import (
	"context"
	"time"

	"github.com/yashkp1234/MemeShop.git/api/database"
	"github.com/yashkp1234/MemeShop.git/api/models"
	"github.com/yashkp1234/MemeShop.git/api/security"
	"github.com/yashkp1234/MemeShop.git/api/utils/channels"
	"go.mongodb.org/mongo-driver/bson"
)

//SignIn handles the sign in for a user
func SignIn(username, password string) (string, error) {
	var user models.User
	var err error
	db := database.Connect()
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"username": username}
		if err = db.Collection("users").FindOne(ctx, filter).Decode(&user); err != nil {
			ch <- false
			return
		}

		err = security.VerifyPassword(user.Password, password)
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return CreateToken(user.ID.Hex(), username)
	}

	return "", err

}
