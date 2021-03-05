package contextkey

import (
	"errors"
	"net/http"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	// ContextKeyUsernameCaller var
	ContextKeyUsernameCaller = contextKey("username")
	// ContextKeyUserIDCaller var
	ContextKeyUserIDCaller = contextKey("userID")
)

// GetUsernameFromContext gets the caller value from the context.
func GetUsernameFromContext(r *http.Request) (string, error) {
	username, err := r.Context().Value(ContextKeyUsernameCaller).(string)
	if !err {
		return "", errors.New("No username found in ctx")
	}
	return username, nil
}

// GetUserIDFromContext gets the caller id value from the context
func GetUserIDFromContext(r *http.Request) (string, error) {
	id, err := r.Context().Value(ContextKeyUserIDCaller).(string)
	if !err {
		return "", errors.New("No username found in ctx")
	}
	return id, nil
}
