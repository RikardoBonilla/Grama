package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/praxisvr/grama/internal/infrastructure/auth"
	"github.com/praxisvr/grama/internal/infrastructure/http/middleware"
)

// okHandler returns 200 OK to confirm the request passed all middleware.
var okHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
})

// contextWithRole builds a request context that already contains JWT claims,
// simulating a request that has passed through AuthMiddleware.Require.
func contextWithRole(role string) context.Context {
	claims := &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: "user-123"},
		Role:             role,
	}
	return auth.NewContextWithClaims(context.Background(), claims)
}

func TestRequireRole_Allowed(t *testing.T) {
	handler := middleware.RequireRole("owner")(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(contextWithRole("owner"))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	handler := middleware.RequireRole("owner")(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(contextWithRole("client"))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rr.Code)
	}
}

func TestRequireRole_MultipleRoles(t *testing.T) {
	handler := middleware.RequireRole("owner", "operator")(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(contextWithRole("operator"))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRequireRole_NoToken(t *testing.T) {
	handler := middleware.RequireRole("owner")(okHandler)

	// Plain background context — no claims injected.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}
