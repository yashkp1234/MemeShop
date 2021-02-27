package api

import (
	"api/router"
	"config"
	"fmt"
	"log"
	"net/http"
)

//Run runs the server
func Run() {
	config.Load()
	fmt.Printf("Listening on [::]%d ... \n", config.PORT)
	Listen(config.PORT)

}

func Listen(port int) {
	r := router.New()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
