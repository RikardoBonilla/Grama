package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/praxisvr/grama/internal/infrastructure/auth"
)

// RequireRole returns 403 if the authenticated user's role is not in the allowed list.
// Must be used after AuthMiddleware.Require (which injects the claims).
// Usage: RequireRole("owner") or RequireRole("owner", "operator")
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := auth.ClaimsFromContext(r.Context())
			if !ok {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
				return
			}
			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		})
	}
}

// writeJSON is a package-level helper shared by all middleware in this package.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
