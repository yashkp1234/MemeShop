package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/auth"
	"github.com/yashkp1234/MemeShop.git/api/responses"
	"github.com/yashkp1234/MemeShop.git/api/utils/contextkey"
)

//SetMiddlewareLogger sets up the next middleware logger
func SetMiddlewareLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
func SetMiddlewareAuth(next http.HandlerFunc, authReq bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := auth.ValidateToken(r)
		if !authReq && err != nil {
			ctx := context.WithValue(r.Context(), contextkey.ContextKeyUsernameCaller, "")
			ctx2 := context.WithValue(ctx, contextkey.ContextKeyUserIDCaller, "")
			next(w, r.WithContext(ctx2))
			return
		} else if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, err)
			return
		}
		next(w, response)
	}
}
