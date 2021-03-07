package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/database"
	"github.com/yashkp1234/MemeShop.git/api/gcp"
	"github.com/yashkp1234/MemeShop.git/api/imagecache"
	"github.com/yashkp1234/MemeShop.git/api/router"
	"github.com/yashkp1234/MemeShop.git/config"
)

//Run runs the server
func Run() {
	config.Load()

	database.NewDatabase()
	defer database.CancelDB()

	imagecache.NewImageCacheClient()
	defer imagecache.Cancel()

	gcp.NewStorageClient()
	defer gcp.Disconnect()

	fmt.Printf("Listening on [::]%d ... \n", config.PORT)
	Listen(config.PORT)
}

//Listen creates a server and listens
func Listen(port int) {
	r := router.New()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
