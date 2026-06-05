package http

import (
	"log"
	stdhttp "net/http"
	"time"

	"github.com/praxisvr/grama/internal/infrastructure/http/handler"
	"github.com/praxisvr/grama/internal/infrastructure/http/middleware"
)

// NewRouter builds the HTTP multiplexer and registers all routes.
// Middleware order: CORS -> rate limit -> auth -> role -> handlers.
func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authMW *middleware.AuthMiddleware,
	allowedOrigins []string,
) stdhttp.Handler {
	mux := stdhttp.NewServeMux()

	// Rate limiters for credential endpoints (5 req/min per IP).
	authLimiter := middleware.NewIPRateLimiter(5, time.Minute)
	loginLimiter := middleware.NewIPRateLimiter(5, time.Minute)

	// Auth endpoints.
	mux.Handle("POST /auth/register",
		authLimiter.Limit(stdhttp.HandlerFunc(authHandler.Register)))
	mux.Handle("POST /auth/login",
		loginLimiter.Limit(stdhttp.HandlerFunc(authHandler.Login)))
	mux.Handle("POST /auth/refresh",
		stdhttp.HandlerFunc(authHandler.Refresh))
	mux.Handle("POST /auth/logout",
		authMW.Require(stdhttp.HandlerFunc(authHandler.Logout)))

	// User profile — all authenticated roles.
	mux.Handle("GET /users/me",
		authMW.Require(stdhttp.HandlerFunc(userHandler.GetMe)))
	mux.Handle("PATCH /users/me",
		authMW.Require(stdhttp.HandlerFunc(userHandler.UpdateMe)))

	// Operator management — owner role only.
	mux.Handle("POST /users/operators",
		authMW.Require(middleware.RequireRole("owner")(stdhttp.HandlerFunc(userHandler.CreateOperator))))
	mux.Handle("GET /users/operators",
		authMW.Require(middleware.RequireRole("owner")(stdhttp.HandlerFunc(userHandler.ListOperators))))

	mux.HandleFunc("GET /health", func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			log.Printf("health write: %v", err)
		}
	})

	// CORS must wrap the entire mux so preflight OPTIONS requests are handled
	// before any auth or rate-limit check rejects them.
	return middleware.CORS(allowedOrigins)(mux)
}
