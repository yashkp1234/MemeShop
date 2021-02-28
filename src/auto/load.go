package auto

import (
	"context"
	"log"

	"github.com/yashkp1234/MemeShop.git/api/database"
	"github.com/yashkp1234/MemeShop.git/api/utils/console"
)

//Load loads data into database
func Load() {
	//Connect to db
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	collection := db.Collection("users")

	for _, user := range users {
		//Hash user password
		err := user.HashPassword()
		if err != nil {
			log.Fatal(err)
		}

		// Insert user into db
		_, err = collection.InsertOne(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}

		//Print out result
		console.Pretty(user)
	}
}
