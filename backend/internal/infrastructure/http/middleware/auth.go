package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/praxisvr/grama/internal/infrastructure/auth"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

// AuthMiddleware validates JWT access tokens on protected routes.
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

// Require validates the Bearer token and injects claims into the request context.
// Returns 401 if the token is missing or invalid.
func (m *AuthMiddleware) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if !strings.HasPrefix(bearer, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(bearer, "Bearer ")
		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ClaimsFromContext extracts JWT claims injected by AuthMiddleware.Require.
func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*auth.Claims)
	return claims, ok
}
