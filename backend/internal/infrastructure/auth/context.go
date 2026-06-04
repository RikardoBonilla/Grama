package auth

import "context"

type contextKey string

const claimsKey contextKey = "auth_claims"

// NewContextWithClaims stores validated JWT claims in the context.
// Called by AuthMiddleware.Require; used by RequireRole and handlers.
func NewContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// ClaimsFromContext retrieves the claims injected by AuthMiddleware.Require.
// Returns (nil, false) if the request did not pass through the auth middleware.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}
