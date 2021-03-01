package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yashkp1234/MemeShop.git/api/middlewares"
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

//SetupRoutesWithMiddlewares sets up the routes for mux router
func SetupRoutesWithMiddlewares(r *mux.Router) *mux.Router {
	for _, route := range Load() {
		r.HandleFunc(route.URI,
			middlewares.SetMiddlewareLogger(
				middlewares.SetMiddlewareJSON(route.Handler)),
		).Methods(route.Method)
	}
	return r
}
