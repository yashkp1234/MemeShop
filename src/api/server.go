package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/router"
	"github.com/yashkp1234/MemeShop.git/config"
)

//Run runs the server
func Run() {
	config.Load()
	fmt.Printf("Listening on [::]%d ... \n", config.PORT)
	Listen(config.PORT)

}

//Listen creates a server and listens
func Listen(port int) {
	r := router.New()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
