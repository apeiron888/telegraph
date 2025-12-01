package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"telegraph/internal/acl"
	"telegraph/internal/users"

	"github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	userContextKey contextKey = "user"
)

// LoadUser middleware fetches the full user from database after JWT validation
// Must be used after JWTAuth middleware
func LoadUser(userRepo users.UserRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDStr := users.UserIDFromContext(r.Context())
			if userIDStr == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				http.Error(w, "invalid_user_id", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetByID(r.Context(), userID)
			if err != nil {
				http.Error(w, "user_not_found", http.StatusUnauthorized)
				return
			}

			// Store full user in context
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the full user from context
func GetUserFromContext(ctx context.Context) (*users.User, bool) {
	user, ok := ctx.Value(userContextKey).(*users.User)
	return user, ok
}

// RequireRole middleware enforces RBAC - user must have exact role or higher
func RequireRole(role acl.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if !acl.HasRoleOrHigher(user.Role, role) {
				respondError(w, "insufficient_permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission middleware enforces RBAC permission check
func RequirePermission(permission acl.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if !acl.HasPermission(user.Role, permission) {
				respondError(w, "permission_denied", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// EnforceMAC middleware checks security label clearance
// resourceLabel should be extracted from the resource being accessed
func EnforceMAC(getResourceLabel func(*http.Request) (string, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			resourceLabel, err := getResourceLabel(r)
			if err != nil {
				respondError(w, "internal_error", http.StatusInternalServerError)
				return
			}

			if !acl.CanAccessResource(user.SecurityLabel, resourceLabel) {
				respondError(w, "insufficient_clearance", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// EnforceABAC middleware evaluates attribute-based policies
func EnforceABAC(policies []acl.Policy) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			allowed, err := acl.EvaluatePolicies(user.Attributes, policies)
			if err != nil {
				respondError(w, "policy_evaluation_error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				respondError(w, "policy_violation", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper to send JSON error responses
func respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
