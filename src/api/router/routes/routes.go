package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

//Route represents a route object
type Route struct {
	URI     string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

//Load loads all routes
func Load() []Route {
	routes := UsersRoutes
	return routes
}

//SetupRoutes sets up the routes for mux router
func SetupRoutes(r *mux.Router) *mux.Router {
	for _, route := range Load() {
		r.HandleFunc(route.URI, route.Handler).Methods(route.Method)
	}
	return r
}
