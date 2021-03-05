package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/auth"
	"github.com/yashkp1234/MemeShop.git/api/responses"
)

//SetMiddlewareLogger sets up the next middleware logger
func SetMiddlewareLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("")
		log.Printf("\n%s %s%s %s", r.Method, r.Host, r.RequestURI, r.Proto)
		next(w, r)
	}
}

//SetMiddlewareJSON sets up the next JSON middleware
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

//SetMiddlewareAuth authenticates the req
func SetMiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := auth.ValidateToken(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, err)
			return
		}
		next(w, response)
	}
}
