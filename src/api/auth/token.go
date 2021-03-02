package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/yashkp1234/MemeShop.git/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreateToken creates a JWT token based on user id
func CreateToken(userID primitive.ObjectID) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"used_id":    userID,
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(config.JWTSecret)
}
