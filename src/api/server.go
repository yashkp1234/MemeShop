package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/cache"
	"github.com/yashkp1234/MemeShop.git/api/database"
	"github.com/yashkp1234/MemeShop.git/api/gcp"
	"github.com/yashkp1234/MemeShop.git/api/router"
	"github.com/yashkp1234/MemeShop.git/config"
)

//Run runs the server
func Run() {
	config.Load()

	database.NewDatabase()
	defer database.CancelDB()

	cache.NewCacheClient()
	defer cache.Cancel()

	gcp.NewStorageClient()
	defer gcp.Disconnect()

	log.Printf("Listening on [::]%d ... \n", config.PORT)
	Listen(config.PORT)
}

//Listen creates a server and listens
func Listen(port int) {
	r := router.New()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
