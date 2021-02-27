package router

import (
	"api/router/routes"

	"github.com/gorilla/mux"
)

// New returns a new route handler
func New() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return routes.SetupRoutes(r)
}
