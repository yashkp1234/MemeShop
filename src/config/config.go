package config

import (
	"log"
	"os"
	"path/filepath"
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
	//DBName for database
	DBName = ""
)

//Load loads the configs from file
func Load() {
	fp, err := filepath.Abs("../.env")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(fp)
	if err != nil {
		log.Fatal(err)
	}

	PORT, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	MongoURL = os.Getenv("MONGO_URL")
	JWTSecret = os.Getenv("JWTSECRET")
	DBName = os.Getenv("DB_NAME")
}
