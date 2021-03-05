package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/yashkp1234/MemeShop.git/api/utils/contextkey"
	"github.com/yashkp1234/MemeShop.git/config"
)

//CreateToken creates a JWT token based on user id
func CreateToken(userID string, username string) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"username":   username,
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

//ValidateToken validates a token
func ValidateToken(r *http.Request) (*http.Request, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Errors processing JWT claims")
	}

	uid, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("Cannot find user_id in jwt claims")
	}

	ctx := context.WithValue(r.Context(), contextkey.ContextKeyUsernameCaller, claims["username"].(string))
	ctx2 := context.WithValue(ctx, contextkey.ContextKeyUserIDCaller, uid)

	return r.WithContext(ctx2), nil
}

//ExtractToken get the token from the request
func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}

	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}
