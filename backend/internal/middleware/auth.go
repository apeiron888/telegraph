package middleware

import (
	"net/http"
	"strings"

	"telegraph/internal/users"
)

func JWTAuth(jwtMgr *users.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing_token", 401)
				return
			}

			parts := strings.Split(auth, " ")
			if len(parts) != 2 {
				http.Error(w, "invalid_auth_header", 401)
				return
			}

			userID, err := jwtMgr.Verify(parts[1])
			if err != nil {
				http.Error(w, "invalid_token", 401)
				return
			}

			ctx := users.ContextWithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
