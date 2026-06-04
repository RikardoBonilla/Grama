package middleware

import "net/http"

// RequireVenueOwner verifies that the venue_id in the request
// belongs to the authenticated owner or is assigned to the operator.
// TODO GRAM-16: implement full IDOR check after venues table
// is available in Sprint 2 (GRAM-19). Currently passes through.
func RequireVenueOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
