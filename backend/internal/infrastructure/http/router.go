package http

import (
	stdhttp "net/http"
	"time"

	"github.com/praxisvr/grama/internal/infrastructure/http/handler"
	"github.com/praxisvr/grama/internal/infrastructure/http/middleware"
)

// NewRouter builds the HTTP multiplexer and registers all routes.
// Middleware order matters: security headers -> rate limit -> auth -> handlers.
func NewRouter(authHandler *handler.AuthHandler, authMW *middleware.AuthMiddleware) stdhttp.Handler {
	mux := stdhttp.NewServeMux()

	// Endpoint-specific rate limiters — tighter limits on credential endpoints.
	authLimiter  := middleware.NewIPRateLimiter(5, time.Minute)
	loginLimiter := middleware.NewIPRateLimiter(5, time.Minute)

	mux.Handle("POST /auth/register",
		authLimiter.Limit(stdhttp.HandlerFunc(authHandler.Register)))

	mux.Handle("POST /auth/login",
		loginLimiter.Limit(stdhttp.HandlerFunc(authHandler.Login)))

	mux.Handle("POST /auth/refresh",
		stdhttp.HandlerFunc(authHandler.Refresh))

	mux.Handle("POST /auth/logout",
		authMW.Require(stdhttp.HandlerFunc(authHandler.Logout)))

	mux.HandleFunc("GET /health", func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}
