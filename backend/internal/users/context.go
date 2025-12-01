package users

import (
	"context"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

// ContextWithUserID stores the authenticated user id inside the context.
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext extracts the authenticated user id.
func UserIDFromContext(ctx context.Context) string {
	val := ctx.Value(userIDKey)
	if val == nil {
		return ""
	}
	id, _ := val.(string)
	return id
}
