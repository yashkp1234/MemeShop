package contextkey

import "context"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	// ContextKeyUsernameCaller var
	ContextKeyUsernameCaller = contextKey("username")
)

// GetCallerFromContext gets the caller value from the context.
func GetCallerFromContext(ctx context.Context) (string, bool) {
	caller, ok := ctx.Value(ContextKeyUsernameCaller).(string)
	return caller, ok
}
