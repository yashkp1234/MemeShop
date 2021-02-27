package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	//PORT for the API
	PORT = 0
	//MongoURL for MONGODB
	MongoURL = ""
	//JWTSecret for jwt encryption
	JWTSecret = ""
)

//Load loads the configs from file
func Load() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	PORT, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	MongoURL = os.Getenv("MONGOURL")
	JWTSecret = os.Getenv("JWTSECRET")
}
