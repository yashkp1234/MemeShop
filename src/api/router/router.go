package router

import (
	"github.com/gorilla/mux"
	"github.com/yashkp1234/MemeShop.git/api/router/routes"
)

// New returns a new route handler
func New() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return routes.SetupRoutes(r)
}
